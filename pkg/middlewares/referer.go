package middlewares

import (
	"github.com/lizongying/go-crawler/pkg"
)

type RefererMiddleware struct {
	pkg.UnimplementedMiddleware
	logger pkg.Logger

	refererPolicy pkg.ReferrerPolicy
}

func (m *RefererMiddleware) ProcessRequest(c *pkg.Context) (err error) {
	m.logger.Debug("enter ProcessRequest")
	defer func() {
		m.logger.Debug("exit ProcessRequest")
	}()

	request := c.Request

	if m.refererPolicy == pkg.NoReferrerPolicy && request.Header != nil && request.Header.Get("Referer") != "" {
		//request.Header.Del("Referer")
	}

	if m.refererPolicy == pkg.DefaultReferrerPolicy && request.Header != nil && request.Referer != "" {
		request.SetHeader("Referer", request.Referer)
	}

	err = c.NextRequest()
	if err != nil {
		m.logger.Debug(err)
		return
	}

	return
}

func (m *RefererMiddleware) FromCrawler(spider pkg.Spider) pkg.Middleware {
	m.logger = spider.GetLogger()
	m.refererPolicy = spider.GetConfig().GetReferrerPolicy()
	return m
}

func NewRefererMiddleware() pkg.Middleware {
	return &RefererMiddleware{}
}
