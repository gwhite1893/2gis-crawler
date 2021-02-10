package app

import (
	"context"
	"log"
	"sync"

	"github.com/gwhite1893/2gis-crawler/internal/crawler"

	"github.com/gwhite1893/2gis-crawler/config"
	"github.com/gwhite1893/2gis-crawler/internal/web"

	"github.com/pkg/errors"
)

var errAppFailed = errors.New("app  failed")

type Option func(ctx context.Context, e *app)

type App interface {
	Shutdown()
}

type app struct {
	cancel context.CancelFunc
	wg     *sync.WaitGroup

	crawler    crawler.Crawler
	httpServer *web.Server
}

func NewApp(opts ...Option) (App, error) {
	ctx, cancel := context.WithCancel(context.Background())

	e := &app{
		cancel: cancel,
		wg:     &sync.WaitGroup{},
	}

	for _, opt := range opts {
		opt(ctx, e)
	}

	select {
	case <-ctx.Done():
		e.wg.Wait()

		return nil, errAppFailed
	default:
	}

	log.Println("app started")

	return e, nil
}

func (e *app) Shutdown() {
	e.cancel()
	e.wg.Wait()
	log.Println("app shutdown")
}

func WithHTTPServer(conf *config.HTTPServerCfg) Option {
	return func(ctx context.Context, e *app) {
		if ctx.Err() != nil {
			log.Println("http server failed: ", ctx.Err())

			return
		}

		if conf == nil {
			log.Println("http server configuration failed")
			e.cancel()

			return
		}

		s, err := web.NewHTTPServer(
			ctx,
			e.wg,
			conf,
		)
		if err != nil {
			log.Println("http server start failed", err)
			e.cancel()

			return
		}

		e.httpServer = s
	}
}

func WithCrawler(conf *config.CrawlerCfg) Option {
	return func(ctx context.Context, e *app) {
		if ctx.Err() != nil {
			log.Println("crawler failed: ", ctx.Err())

			return
		}

		if conf == nil {
			log.Println("crawler configuration failed")
			e.cancel()

			return
		}

		c, err := crawler.NewCrawler(
			ctx,
			e.wg,
			conf,
		)
		if err != nil {
			log.Println("crawler start failed", err)
			e.cancel()

			return
		}

		e.crawler = c
	}
}
