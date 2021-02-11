package crawler

import (
	"context"
	"errors"
)

var errTimeout = errors.New("timeout")

type command struct {
	ctx   context.Context
	value string
	resCh chan<- resourceResponse
}

type resourceResponse struct {
	Content []byte
	Err     error
}
