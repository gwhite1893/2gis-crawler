package crawler

import (
	"context"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/gwhite1893/2gis-crawler/config"
)

type Result struct {
	Content []byte
	Err     string
	URL     string
}

type Response []*Result

type Crawler interface {
	Crawl(ctx context.Context, links []string) (Response, error)
}

type crawler struct {
	cancel         context.CancelFunc
	httpClient     *http.Client
	timeoutSec     int
	commandChannel chan command
}

func NewCrawler(ctx context.Context, wg *sync.WaitGroup, conf *config.CrawlerCfg) (Crawler, error) {
	ctx, cancel := context.WithCancel(ctx)

	c := &crawler{
		cancel:         cancel,
		timeoutSec:     conf.RequestTimeOutSec,
		httpClient:     http.DefaultClient,
		commandChannel: make(chan command),
	}

	wg.Add(1)

	go func() {
		c.run(ctx)
		wg.Done()
	}()

	return c, nil
}

func (c *crawler) Crawl(ctx context.Context, urls []string) (Response, error) {
	var wg sync.WaitGroup

	resCh := make(chan *Result, len(urls))

	for i := range urls {
		wg.Add(1)

		commandCh := make(chan resourceResponse)
		c.commandChannel <- command{
			ctx:   ctx,
			value: urls[i],
			resCh: commandCh,
		}

		// читаем из канала команды результат запроса
		// и отправляем в результирующий канал
		go func(val string) {
			defer wg.Done()

			res := <-commandCh
			if res.Err != nil {
				resCh <- &Result{
					Err: res.Err.Error(),
					URL: val,
				}

				return
			}

			resCh <- &Result{
				Content: res.Content,
				URL:     val,
			}
		}(urls[i])
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	resp := Response{}

	for res := range resCh {
		resp = append(resp, res)
	}

	return resp, nil
}

func (c *crawler) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.cancel()

			return
		case cmd := <-c.commandChannel:
			{
				res := c.poll(cmd.ctx, cmd.value)

				if res.Err != nil {
					cmd.resCh <- resourceResponse{
						Err: errors.Wrapf(res.Err, "request failed"),
					}

					continue
				}

				cmd.resCh <- resourceResponse{
					Content: res.Content,
				}
			}
		}
	}
}

func (c *crawler) poll(ctx context.Context, url string) resourceResponse {
	resCh := make(chan resourceResponse)

	go func() {
		data, err := makeRequest(ctx, url, c.httpClient)
		if err != nil {
			resCh <- resourceResponse{Err: err}

			return
		}

		resCh <- resourceResponse{Content: data}
	}()

	return waitResult(resCh, c.timeoutSec)
}

func makeRequest(ctx context.Context, url string, client *http.Client) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// wait result обеспечивает механизм таймаута при чтении из канала
func waitResult(resCh chan resourceResponse, maxTimeoutSec int) resourceResponse {
	for {
		select {
		case r := <-resCh:
			return r
		case <-time.After(time.Duration(maxTimeoutSec) * time.Second):
			return resourceResponse{Err: errTimeout}
		}
	}
}
