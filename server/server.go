package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	rate_limit "make-me-reliable/rate-limit"
	"make-me-reliable/repository"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server is a web server struct definition
type Server struct {
	r  repository.Repository
	rl rate_limit.RateLimit
}

type errMsg struct {
	Message string `json:"message"`
}

// NewServer returns new instance of server
func NewServer(r repository.Repository, rl rate_limit.RateLimit) *Server {
	return &Server{
		r:  r,
		rl: rl,
	}
}

// Start will register routes and start a web server on provided addr
// Note addr must be in form: IP:PORT e.g. 0.0.0.0:8000
func (s *Server) Start(addr string) error {

	r := mux.NewRouter()
	r.HandleFunc("/crypto/sign", newSignHandler(s.r, s.rl)).Methods("GET")
	r.HandleFunc("/jobs/{jobId}", newJobHandler(s.r)).Methods("GET")

	server := http.Server{
		Addr: addr,

		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,

		Handler: r,
	}

	logrus.Info("Serving on: ", server.Addr)
	return server.ListenAndServe()
}

func responseWithStatus(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if response != nil {
		r, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.Warnf("Error formatting response %+v", err)
			if _, err := fmt.Fprintf(w, "Unexpected error"); err != nil {
				logrus.Warnf("Error writing response %s", err)
			}
			return
		}
		if _, err := w.Write(r); err != nil {
			logrus.Warnf("Error writing response %s", err)
		}
	}
}
