package rate_limit

import (
	"context"
	"fmt"

	"github.com/juju/ratelimit"

	"make-me-reliable/http"
)

// RateLimit is a rate limit interface definition
type RateLimit interface {
	IsPermitted() bool
	Do(context.Context, map[string]string) (interface{}, error)
}

// Http is a HTTP rate limited client
type Http struct {
	r    *ratelimit.Bucket
	c    http.Caller
	path string
}

// NewHttpRateLimit returns new Http rate limiter
func NewHttpRateLimit(r *ratelimit.Bucket, c http.Caller, path string) *Http {
	return &Http{
		r:    r,
		c:    c,
		path: path,
	}
}

// IsPermitted returns if there is more than 0 tokens available
func (h *Http) IsPermitted() bool {
	return h.r.Available() > 0
}

// Do will do an HTTP call as specified by passed HTTP client
func (h *Http) Do(ctx context.Context, params map[string]string) (interface{}, error) {
	r := h.r.Take(1)
	if r.Seconds() > 0 {
		return nil, fmt.Errorf("not permitted to do call")
	}
	return h.c.CallAPI(ctx, h.path, params)
}
