package background

import (
	"fmt"
	"testing"
	"time"

	rate_limit "make-me-reliable/rate-limit"

	"make-me-reliable/repository"
)

func TestRateLimitedWorker_WorkNotPermitted(t *testing.T) {
	closeChan := make(chan interface{})
	r := &repository.MockRepository{}
	rl := &rate_limit.MockRateLimiter{
		IsPermittedFunc: func() bool {
			return false
		},
	}
	w := NewRateLimitedWorker(r, rl)

	go w.StartWork(closeChan)

	<-time.After(10 * time.Second)
	close(closeChan)

	if r.GetCountGet() != 0 {
		t.Errorf("Expected 0, got %d", r.GetCountGet())
	}
	if r.GetCountUnfinished() != 0 {
		t.Errorf("Expected 0, got %d", r.GetCountUnfinished())
	}
}

func TestRateLimitedWorker_WorkPermitted_JobNotFound(t *testing.T) {
	closeChan := make(chan interface{})
	r := &repository.MockRepository{
		GetUnfinishedFunc: func() (*repository.Job, error) {
			return nil, &repository.EmptyResponseError{}
		},
	}
	rl := &rate_limit.MockRateLimiter{
		IsPermittedFunc: func() bool {
			return true
		},
	}
	w := NewRateLimitedWorker(r, rl)

	go w.StartWork(closeChan)

	<-time.After(10 * time.Second)
	close(closeChan)

	if r.GetCountUnfinished() != 1 {
		t.Errorf("Expected 1, got %d", r.GetCountUnfinished())
	}
	if r.GetCountUpdate() != 0 {
		t.Errorf("Expected 0, got %d", r.GetCountUpdate())
	}
}

func TestRateLimitedWorker_WorkPermitted_DBFail(t *testing.T) {
	closeChan := make(chan interface{})
	r := &repository.MockRepository{
		GetUnfinishedFunc: func() (*repository.Job, error) {
			return nil, fmt.Errorf("error")
		},
	}
	rl := &rate_limit.MockRateLimiter{
		IsPermittedFunc: func() bool {
			return true
		},
	}
	w := NewRateLimitedWorker(r, rl)

	go w.StartWork(closeChan)

	<-time.After(10 * time.Second)
	close(closeChan)

	if r.GetCountUnfinished() != 1 {
		t.Errorf("Expected 1, got %d", r.GetCountUnfinished())
	}
	if r.GetCountUpdate() != 0 {
		t.Errorf("Expected 0, got %d", r.GetCountUpdate())
	}
}

func TestRateLimitedWorker_WorkPermitted_CallFail(t *testing.T) {
	closeChan := make(chan interface{})
	r := &repository.MockRepository{
		GetUnfinishedFunc: func() (*repository.Job, error) {
			return nil, fmt.Errorf("error")
		},
	}
	rl := &rate_limit.MockRateLimiter{
		IsPermittedFunc: func() bool {
			return true
		},
		DoFunc: func() (interface{}, error) {
			return nil, fmt.Errorf("Api call failed")
		},
	}
	w := NewRateLimitedWorker(r, rl)

	go w.StartWork(closeChan)

	<-time.After(10 * time.Second)
	close(closeChan)

	if r.GetCountUnfinished() != 1 {
		t.Errorf("Expected 1, got %d", r.GetCountUnfinished())
	}
	if r.GetCountUpdate() != 0 {
		t.Errorf("Expected 0, got %d", r.GetCountUpdate())
	}
}

func TestRateLimitedWorker_WorkPermitted_Success(t *testing.T) {
	closeChan := make(chan interface{})
	r := &repository.MockRepository{
		GetUnfinishedFunc: func() (*repository.Job, error) {
			j := repository.NewJob()
			return &j, nil
		},
		UpdateJobFunc: func() error {
			return nil
		},
	}
	rl := &rate_limit.MockRateLimiter{
		IsPermittedFunc: func() bool {
			return true
		},
		DoFunc: func() (interface{}, error) {
			return "success!", nil
		},
	}
	w := NewRateLimitedWorker(r, rl)

	go w.StartWork(closeChan)

	<-time.After(10 * time.Second)
	close(closeChan)

	if r.GetCountUnfinished() != 1 {
		t.Errorf("Expected 1, got %d", r.GetCountUnfinished())
	}
	if r.GetCountUpdate() != 1 {
		t.Errorf("Expected 1, got %d", r.GetCountUpdate())
	}
}
