package server

import (
	"net/http"
	"reflect"
	"testing"

	rate_limit "make-me-reliable/rate-limit"
	"make-me-reliable/repository"
)

func Test_newSignHandler(t *testing.T) {
	type args struct {
		rdb repository.Repository
		rl  rate_limit.RateLimit
	}
	tests := []struct {
		name string
		args args
		want func(w http.ResponseWriter, r *http.Request)
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSignHandler(tt.args.rdb, tt.args.rl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSignHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
