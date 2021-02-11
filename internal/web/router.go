package web

//go:generate swag init -g ../../cmd/2gis-crawler/main.go -o ../../cmd/2gis-crawler/docs

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/gwhite1893/2gis-crawler/cmd/2gis-crawler/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/render"
	"github.com/gwhite1893/2gis-crawler/internal/parser"
	"github.com/pkg/errors"
)

const (
	swaggerVersion = "1.0.0"

	badRequestMessage     = "bad request"
	badRequestDescription = "Не удалось получить параметры запроса"
)

var errBadRequest = errors.New(badRequestMessage)

type sourcePollRequest struct {
	Data []string `json:"data"`
}

func (s *sourcePollRequest) Bind(*http.Request) error {
	if s.Data == nil {
		return errors.WithMessage(errBadRequest, "empty url list")
	}

	return nil
}

type PollResult struct {
	URL   string `json:"url"`
	Body  string `json:"body"`
	Error string `json:"error,omitempty"`
}

type sourcesPollResponse struct {
	Data []*PollResult `json:"data"`
}

func (s *sourcesPollResponse) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

// Poll godoc
// @Summary Sites polling
// @Description Poll by request url
// @Tags sources
// @ID      sources-poll
// @Accept  json
// @Produce json
// @Param	data body sourcePollRequest true "urls array"
// @Success 200 {object}  sourcesPollResponse
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Router /resources/poll [post]
func (s *Server) PollSources(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &sourcePollRequest{}
	if err := render.Bind(r, req); err != nil {
		log.Print(err)

		_ = render.Render(w, r, errInvalidRequest())

		return
	}

	result, err := s.crawler.Crawl(ctx, req.Data)
	if err != nil {
		log.Print(err)

		_ = render.Render(w, r, errRender(
			http.StatusInternalServerError,
			err.Error(),
			"crawl failed",
		))

		return
	}

	const tagName = "title"

	pollResult := make([]*PollResult, len(result))
	for i := range result {
		pollResult[i] = &PollResult{
			URL:   result[i].URL,
			Body:  parser.GetTagValue(bytes.NewReader(result[i].Content), tagName),
			Error: result[i].Err,
		}
	}

	_ = render.Render(w, r, &sourcesPollResponse{
		Data: pollResult,
	})
}

// Swagger godoc
// @Summary swagger
// @Tags swagger
// @Description Описание API
// @ID swagger
// @Produce html
// @Success 200 "swagger html page"
// @Router /swagger [get]
func (s *Server) ServeSwagger() http.HandlerFunc {
	docs.SwaggerInfo.Version = swaggerVersion
	docs.SwaggerInfo.Host = s.httpServer.Addr
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	return httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", s.httpServer.Addr)),
	)
}

func errInvalidRequest() render.Renderer {
	return &errResponse{
		HTTPStatusCode: http.StatusBadRequest,
		ErrorText:      badRequestMessage,
		Description:    badRequestDescription,
	}
}

type errResponse struct {
	HTTPStatusCode int    `json:"status,omitempty"`
	ErrorText      string `json:"error,omitempty"`
	Description    string `json:"description,omitempty"`
}

func (e *errResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)

	return nil
}

func errRender(httpStatusCode int, errorText, description string) render.Renderer {
	return &errResponse{
		HTTPStatusCode: httpStatusCode,
		ErrorText:      errorText,
		Description:    description,
	}
}
