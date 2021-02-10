package web

//go:generate swag init -g ../../cmd/2gis-crawler/main.go -o ../../cmd/2gis-crawler/docs

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

const (
	badRequestMessage     = "bad request"
	badRequestDescription = "Не удалось получить параметры запроса"
)

var (
	errBadRequest = errors.New(badRequestMessage)
)

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
	URL  string `json:"url"`
	Body string `json:"body"`
}

type sourcePollResponse struct {
	Data []PollResult `json:"data"`
}

func (s *sourcePollResponse) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

// Poll godoc
// @Summary Sites polling
// @Description Poll by request url
// @Tags sources
// @ID      sources-poll
// @Accept  json
// @Produce json
// @Param	data body sourcePollRequest true "data"
// @Success 200 {object}  sourcePollResponse
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Router /resources/poll [post]
func (s *Server) PollSources(w http.ResponseWriter, r *http.Request) {
	req := &sourcePollRequest{}
	if err := render.Bind(r, req); err != nil {
		log.Print(err)

		_ = render.Render(w, r, errInvalidRequest())

		return
	}

	_ = render.Render(w, r, &sourcePollResponse{
		Data: []PollResult{
			{
				URL:  "url",
				Body: "body",
			},
		},
	})
}

func errInvalidRequest() render.Renderer {
	return &errResponse{
		HTTPStatusCode: http.StatusBadRequest,
		ErrorText:      badRequestMessage,
		Description:    badRequestDescription,
	}
}

type errResponse struct {
	HTTPStatusCode int    `json:"status,omitempty"`      // http response status code
	ErrorText      string `json:"error,omitempty"`       // application-level error message
	Description    string `json:"description,omitempty"` // user-level status message
}

func (e *errResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)

	return nil
}
