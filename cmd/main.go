package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/juju/ratelimit"
	"github.com/sirupsen/logrus"

	"make-me-reliable/background"
	"make-me-reliable/http"
	"make-me-reliable/internal/config"
	"make-me-reliable/internal/database"
	rate_limit "make-me-reliable/rate-limit"
	"make-me-reliable/repository"
	"make-me-reliable/server"
)

func main() {
	c := config.ParseFromEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.MongoConnectionFromConfigCtx(ctx, c)
	if err != nil {
		log.Println(err)
		panic("can't connect to DB")
	}

	collection := db.Database("reliable-api").Collection("jobs")
	r := repository.NewMongoRepository(collection)

	client := http.NewGetCaller(c.APIAuthHeader, c.APIUrl)
	rl := ratelimit.NewBucketWithQuantum(time.Minute, 10, 10)
	hrl := rate_limit.NewHttpRateLimit(rl, client, "crypto/sign")

	w := background.NewRateLimitedWorker(r, hrl)
	closeChan := make(chan interface{})
	go w.StartWork(closeChan)

	ws := server.NewServer(r, hrl)
	addr := fmt.Sprintf("0.0.0.0:%s", c.ServerPort)
	logrus.Fatal(ws.Start(addr))
}
