package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"make-me-reliable/repository"
)

func Test_newJobHandler_invalidJobId(t *testing.T) {
	r := &repository.MockRepository{}

	req, err := http.NewRequest("GET", "/job/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newJobHandler(r))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"message":"JobID is not a valid ID"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_newJobHandler_jobNotFound(t *testing.T) {
	r := &repository.MockRepository{
		GetJobByIdFunc: func() (*repository.Job, error) {
			return nil, &repository.EmptyResponseError{}
		},
	}

	req, err := http.NewRequest("GET", "/job/61b639df6b867164c385ee16", nil)
	if err != nil {
		t.Fatal(err)
	}

	//Hack to fake gorilla/mux vars
	vars := map[string]string{
		"jobId": "61b639df6b867164c385ee16",
	}

	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newJobHandler(r))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	expected := `{"message":"JobID is not found"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_newJobHandler_dbError(t *testing.T) {
	r := &repository.MockRepository{
		GetJobByIdFunc: func() (*repository.Job, error) {
			return nil, fmt.Errorf("something wrong")
		},
	}

	req, err := http.NewRequest("GET", "/job/61b639df6b867164c385ee16", nil)
	if err != nil {
		t.Fatal(err)
	}

	//Hack to fake gorilla/mux vars
	vars := map[string]string{
		"jobId": "61b639df6b867164c385ee16",
	}

	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newJobHandler(r))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := `{"message":"Something is wrong"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_newJobHandler_success(t *testing.T) {
	j := repository.NewJob()
	j.Message = "message"
	j.Finished = false

	r := &repository.MockRepository{
		GetJobByIdFunc: func() (*repository.Job, error) {
			return &j, nil
		},
	}

	req, err := http.NewRequest("GET", "/job", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{
		"jobId": j.ID.Hex(),
	}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newJobHandler(r))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := fmt.Sprintf(`{"id":"%s","message":"message","successful":false,"lastTry":"0001-01-01T00:00:00Z"}`, j.ID.Hex())
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
