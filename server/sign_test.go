package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	rate_limit "make-me-reliable/rate-limit"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"make-me-reliable/repository"
)

func Test_newSignHandler_messageNotPresent(t *testing.T) {
	r := &repository.MockRepository{}
	rl := &rate_limit.MockRateLimiter{}
	req, err := http.NewRequest("GET", "/crypto/sign", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newSignHandler(r, rl))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"message":"message query param empty"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_newSignHandler_rateLimitNotPermitted(t *testing.T) {
	r := &repository.MockRepository{
		GetJobByIdFunc: func() (*repository.Job, error) {
			return nil, &repository.EmptyResponseError{}
		},
		InsertJobFunc: func() (*primitive.ObjectID, error) {
			o, _ := primitive.ObjectIDFromHex("00000000000000000000")
			return &o, nil
		},
	}

	rl := &rate_limit.MockRateLimiter{
		IsPermittedFunc: func() bool {
			return false
		},
	}

	req, err := http.NewRequest("GET", "/crypto/sign?message=testtest", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newSignHandler(r, rl))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id":"61b710ecb15de102e23f1284","message":"testtest","successful":false,"lastTry":"0001-01-01T00:00:00Z"}`
	if strings.Contains(rr.Body.String(), `"message":"testtest","successful":false,"lastTry":"0001-01-01T00:00:00Z"`) != true {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_newSignHandler_urlCallFails(t *testing.T) {
	r := &repository.MockRepository{
		GetJobByIdFunc: func() (*repository.Job, error) {
			return nil, &repository.EmptyResponseError{}
		},
		InsertJobFunc: func() (*primitive.ObjectID, error) {
			o, _ := primitive.ObjectIDFromHex("00000000000000000000")
			return &o, nil
		},
	}

	rl := &rate_limit.MockRateLimiter{
		IsPermittedFunc: func() bool {
			return true
		},
		DoFunc: func() (interface{}, error) {
			return "response", nil
		},
	}

	req, err := http.NewRequest("GET", "/crypto/sign?message=testtest", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newSignHandler(r, rl))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id":"61b710ecb15de102e23f1284","message":"testtest","successful":true,"lastTry":"0001-01-01T00:00:00Z"}`
	if strings.Contains(rr.Body.String(), `"message":"testtest","successful":true,"lastTry":"0001-01-01T00:00:00Z","result":"response"`) != true {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_newSignHandler_success(t *testing.T) {
	r := &repository.MockRepository{
		GetJobByIdFunc: func() (*repository.Job, error) {
			return nil, &repository.EmptyResponseError{}
		},
		InsertJobFunc: func() (*primitive.ObjectID, error) {
			o, _ := primitive.ObjectIDFromHex("00000000000000000000")
			return &o, nil
		},
	}

	rl := &rate_limit.MockRateLimiter{
		IsPermittedFunc: func() bool {
			return true
		},
		DoFunc: func() (interface{}, error) {
			return nil, fmt.Errorf("API call failed")
		},
	}

	req, err := http.NewRequest("GET", "/crypto/sign?message=testtest", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newSignHandler(r, rl))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id":"61b710ecb15de102e23f1284","message":"testtest","successful":false,"lastTry":"0001-01-01T00:00:00Z"}`
	if strings.Contains(rr.Body.String(), `"message":"testtest","successful":false,"lastTry":"0001-01-01T00:00:00Z"`) != true {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
