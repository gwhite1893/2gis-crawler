package app

import (
	"context"
	"log"
	"sync"

	"github.com/pkg/errors"

	"github.com/gwhite1893/2gis-crawler/config"
	"github.com/gwhite1893/2gis-crawler/internal/crawler"
	"github.com/gwhite1893/2gis-crawler/internal/web"
)

var errAppFailed = errors.New("app  failed")

type Option func(ctx context.Context, e *app)

type App interface {
	Shutdown()
}

type app struct {
	cancel context.CancelFunc
	wg     *sync.WaitGroup

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

func WithHTTPServer(conf *config.Config) Option {
	return func(ctx context.Context, e *app) {
		if ctx.Err() != nil {
			log.Println("http server failed: ", ctx.Err())

			return
		}

		if conf.HTTPServer == nil {
			log.Println("http server configuration failed")
			e.cancel()

			return
		}

		crwler, err := crawler.NewCrawler(ctx, e.wg, conf.Crawler)
		if err != nil {
			log.Println("crawler failed")
			e.cancel()

			return
		}

		s, err := web.NewHTTPServer(
			ctx,
			e.wg,
			conf.HTTPServer,
			crwler,
		)
		if err != nil {
			log.Println("http server start failed", err)
			e.cancel()

			return
		}

		e.httpServer = s
	}
}
