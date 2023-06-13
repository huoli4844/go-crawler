package middlewares

import (
	"bytes"
	"context"
	"errors"
	"github.com/lizongying/go-crawler/pkg"
	"github.com/lizongying/go-crawler/pkg/utils"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type ImageMiddleware struct {
	pkg.UnimplementedMiddleware
	logger pkg.Logger

	spider pkg.Spider
	stats  pkg.StatsWithImage
}

func (m *ImageMiddleware) SpiderStart(_ context.Context, spider pkg.Spider) (err error) {
	m.spider = spider
	m.stats, _ = spider.GetStats().(pkg.StatsWithImage)
	return
}

func (m *ImageMiddleware) ProcessResponse(_ context.Context, response *pkg.Response) (err error) {
	if len(response.BodyBytes) == 0 {
		err = errors.New("BodyBytes empty")
		m.logger.Error(err)
		return
	}

	extra, ok := response.Request.Extra.(pkg.OptionImage)
	if ok {
		img, name, e := image.Decode(bytes.NewReader(response.BodyBytes))
		if e != nil {
			err = e
			m.logger.Error(err)
			return
		}

		rect := img.Bounds()
		extra.SetName(utils.StrMd5(response.Request.URL.String()))
		extra.SetExtension(name)
		extra.SetWidth(rect.Dx())
		extra.SetHeight(rect.Dy())
		if m.stats != nil {
			m.stats.IncImageTotal()
		}
	}

	return
}

func (m *ImageMiddleware) FromCrawler(spider pkg.Spider) pkg.Middleware {
	if m == nil {
		return new(ImageMiddleware).FromCrawler(spider)
	}
	m.logger = spider.GetLogger()
	return m
}
