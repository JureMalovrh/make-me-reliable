package server

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	rate_limit "make-me-reliable/rate-limit"

	"make-me-reliable/repository"
)

func callExternalApi(ctx context.Context, rl rate_limit.RateLimit, message string) (string, error) {
	res, err := rl.Do(ctx, map[string]string{"message": message})
	if err != nil {
		return "", err
	}

	response := res.(string)
	return response, nil

}

func newSignHandler(rdb repository.Repository, rl rate_limit.RateLimit) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		message := r.URL.Query().Get("message")

		if message == "" {
			responseWithStatus(w, http.StatusBadRequest, errMsg{"message query param empty"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*1)
		defer cancel()

		j := repository.NewJob()
		j.Message = message
		j.Finished = false

		if rl.IsPermitted() {
			res, err := callExternalApi(ctx, rl, message)
			if err != nil {
				logrus.Errorf("Error %v", err)
			}

			if err == nil {
				// only if err == nil means everything was fine and we can proceed with result
				j.Finished = true
				j.Result = res
			}

		}
		if ctx.Err() != nil {
			ctx, cancel = context.WithTimeout(r.Context(), time.Second*1)
			defer cancel()
		}
		_, err := rdb.InsertJobCtx(ctx, &j)
		if err != nil {
			logrus.Errorf("error writing in DB %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		responseWithStatus(w, http.StatusOK, j)
	}
}
