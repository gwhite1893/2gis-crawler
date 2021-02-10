package web

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/pkg/errors"

	"github.com/gwhite1893/2gis-crawler/config"
)

type Server struct {
	httpServer http.Server
}

func NewHTTPServer(
	ctx context.Context,
	wg *sync.WaitGroup,
	httpServerCfg *config.HTTPServerCfg,

) (*Server, error) {
	s := &Server{}

	s.setup(httpServerCfg)

	if err := s.serve(ctx); err != nil {
		return nil, errors.WithMessage(err, "unable serve")
	}

	//nolint
	wg.Add(1)

	go func() {
		s.run(ctx)
		wg.Done()
	}()

	log.Printf("http server started with conf: %+v", *httpServerCfg)

	return s, nil
}

func (s *Server) serve(_ context.Context) error {
	listener, err := net.Listen("tcp4", s.httpServer.Addr)
	if err != nil {
		return err
	}

	go func() {
		if err := s.httpServer.Serve(listener); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Println("serve failed", err)
			}
		}
	}()

	return nil
}

func (s *Server) setup(serverCfg *config.HTTPServerCfg) {
	r := chi.NewRouter()
	s.httpServer = http.Server{
		Handler: r,
		Addr:    serverCfg.Addr(),
	}

	r.Use(middleware.URLFormat)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api/crawler/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Route("/sources", func(r chi.Router) {
				r.Post("/poll", s.PollSources)
			})

		})
	})
}

func (s *Server) run(ctx context.Context) {
	<-ctx.Done()

	defaultTimeoutSec := 5

	ctx, cancel := context.WithTimeout(ctx, time.Duration(defaultTimeoutSec)*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		if !errors.Is(err, http.ErrServerClosed) && !errors.Is(err, context.DeadlineExceeded) {
			log.Print("http shutdown")

			return
		}
	}

	log.Println("http shutdown")
}
