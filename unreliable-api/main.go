package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	log.Fatal(start("0.0.0.0:80"))
}

func newSignHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rn := rand.Intn(100)
		log.Printf("deciding %d\n", rn)
		if rn > 50 {
			w.WriteHeader(200)
			w.Write([]byte("fake successful message"))
			return
		}

		w.WriteHeader(500)
	}
}

func start(addr string) error {
	rand.Seed(time.Now().UnixNano())

	r := http.NewServeMux()
	r.HandleFunc("/crypto/sign", newSignHandler())

	server := http.Server{
		Addr: addr,

		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,

		Handler: r,
	}

	log.Printf("Serving on: %s", server.Addr)
	return server.ListenAndServe()
}
