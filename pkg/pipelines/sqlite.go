package pipelines

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/lizongying/go-crawler/pkg"
	"github.com/lizongying/go-crawler/pkg/items"
	"reflect"
	"strings"
	"time"
)

type SqlitePipeline struct {
	pkg.UnimplementedPipeline
	env     string
	logger  pkg.Logger
	sqlite  *sql.DB
	timeout time.Duration
}

func (m *SqlitePipeline) ProcessItem(itemWithContext pkg.ItemWithContext) (err error) {
	spider := m.GetSpider()
	if itemWithContext == nil {
		err = errors.New("nil item")
		m.logger.Error(err)
		spider.IncItemError()
		return
	}
	if itemWithContext.Name() != pkg.ItemSqlite {
		m.logger.Warn("item not support", pkg.ItemSqlite)
		return
	}
	itemSqlite, ok := itemWithContext.GetItem().(*items.ItemSqlite)
	if !ok {
		m.logger.Warn("item parsing failed with", pkg.ItemSqlite)
		return
	}

	if itemSqlite.GetTable() == "" {
		err = errors.New("table is empty")
		m.logger.Error(err)
		spider.IncItemError()
		return
	}

	data := itemWithContext.Data()
	if data == nil {
		err = errors.New("nil data")
		m.logger.Error(err)
		spider.IncItemError()
		return
	}

	if m.env == "dev" {
		m.logger.Debug("current mode don't need save")
		spider.IncItemIgnore()
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	refType := reflect.TypeOf(itemSqlite.Data()).Elem()
	refValue := reflect.ValueOf(itemSqlite.Data()).Elem()
	var columns []string
	var columns1 []string
	var columns2 []string
	var values []any
	for i := 0; i < refType.NumField(); i++ {
		column := refType.Field(i).Tag.Get("column")
		if column == "" {
			column = refType.Field(i).Name
		}
		columns = append(columns, fmt.Sprintf("%s=?", column))
		columns1 = append(columns1, fmt.Sprintf("`%s`", column))
		columns2 = append(columns2, "?")
		value := refValue.Field(i).Interface()
		values = append(values, value)
	}

	s := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, itemSqlite.GetTable(), strings.Join(columns1, ","), strings.Join(columns2, ","))
	stmt, err := m.sqlite.PrepareContext(ctx, s)
	if err != nil {
		m.logger.Error(err)
		spider.IncItemError()
		return
	}
	res, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		var e *mysql.MySQLError
		o := errors.As(err, &e)
		if !o {
			m.logger.Error(err)
			spider.IncItemError()
			return
		}

		if itemSqlite.GetUpdate() && !reflect.ValueOf(itemSqlite.Id()).IsZero() && e.Number == 1062 {
			s = fmt.Sprintf(`UPDATE %s SET %s WHERE id=?`, itemSqlite.GetTable(), strings.Join(columns, ","))
			values = append(values, itemSqlite.Id())
			stmt, err = m.sqlite.PrepareContext(ctx, s)
			if err != nil {
				m.logger.Error(err)
				spider.IncItemError()
				return
			}

			res, err = stmt.ExecContext(ctx, values...)
			if err != nil {
				m.logger.Error(err)
				spider.IncItemError()
				return
			}

			_, err = res.RowsAffected()
			if err != nil {
				m.logger.Error(err)
				spider.IncItemError()
				return
			}

			m.logger.Info(itemSqlite.GetTable(), "update success", itemSqlite.Id())
		} else {
			m.logger.Error(err)
			spider.IncItemError()
			return
		}
	} else {
		id, e := res.LastInsertId()
		if e != nil {
			m.logger.Error(e)
			spider.IncItemError()
			return
		}

		m.logger.Info(itemSqlite.GetTable(), "insert success", id)
	}

	spider.IncItemSuccess()
	return
}

func (m *SqlitePipeline) FromSpider(spider pkg.Spider) pkg.Pipeline {
	if m == nil {
		return new(SqlitePipeline).FromSpider(spider)
	}

	m.UnimplementedPipeline.FromSpider(spider)
	crawler := spider.GetCrawler()
	m.env = spider.GetConfig().GetEnv()
	m.logger = spider.GetLogger()
	m.sqlite = crawler.GetSqlite().Client()
	m.timeout = time.Minute
	return m
}
