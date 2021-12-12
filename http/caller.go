package http

import "context"

// Caller is an interface for different HTTP callers
type Caller interface {
	CallAPI(context.Context, string, map[string]string) (interface{}, error)
}
