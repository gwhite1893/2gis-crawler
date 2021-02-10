package crawler

import (
	"context"
	"net/http"
	"net/url"
	"sync"

	"github.com/pkg/errors"

	"github.com/gwhite1893/2gis-crawler/config"
)

const maxRequests = 3

var (
	errCrawlerFailed     = errors.New("crawler  failed")
	errBadResponseStatus = errors.New("response failed")
)

type Response struct {
	Data string
}

type Crawler interface {
	Crawl(u *url.URL) (Response, error)
	Shutdown()
}

type crawler struct {
	cancel context.CancelFunc
	wg     *sync.WaitGroup

	httpClient http.Client
	maxRequest int
	timeoutSec int

	commandChannel chan command
}

func (c *crawler) Crawl(u *url.URL) (Response, error) {
	panic("implement me")
}

func (c *crawler) Shutdown() {
	panic("implement me")
}

type command struct {
	ctx        context.Context
	value      string
	resChannel chan result
}

type result struct {
	string string
	err    error
}

func NewCrawler(ctx context.Context, wg *sync.WaitGroup, conf *config.CrawlerCfg) (Crawler, error) {
	ctx, cancel := context.WithCancel(ctx)

	c := &crawler{
		cancel:     cancel,
		maxRequest: conf.MaxRequests,
		timeoutSec: conf.RequestTimeOutSec,

		commandChannel: make(chan command),
	}

	// nolint
	wg.Add(1)

	go func() {
		c.run(ctx)
		wg.Done()
	}()

	return c, nil
}

type Option func(ctx context.Context, c *crawler)

func (c *crawler) run(ctx context.Context) {

}
