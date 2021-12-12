package server

import (
	"net/http"
	"reflect"
	"testing"

	"make-me-reliable/repository"
)

func Test_newJobHandler(t *testing.T) {
	type args struct {
		rdb repository.Repository
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
			if got := newJobHandler(tt.args.rdb); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newJobHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
