package rate_limit

import "context"

// MockRateLimiter is a mock used for testing
type MockRateLimiter struct {
	IsPermittedFunc func() bool
	DoFunc          func() (interface{}, error)
}

// IsPermitted will call IsPermittedFunc from MockRateLimiter
func (rl *MockRateLimiter) IsPermitted() bool {
	return rl.IsPermittedFunc()
}

// Do will call DoFunc from MockRateLimiter
func (rl *MockRateLimiter) Do(context.Context, map[string]string) (interface{}, error) {
	return rl.DoFunc()
}
