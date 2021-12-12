package background

import (
	"context"
	"errors"
	"time"

	rate_limit "make-me-reliable/rate-limit"
	"make-me-reliable/repository"

	"github.com/sirupsen/logrus"
)

// RateLimitedWorker is a worker that will respect rate limitation from RateLimit package
type RateLimitedWorker struct {
	c repository.Repository
	r rate_limit.RateLimit
}

// NewRateLimitedWorker will return new RateLimitedWorker
func NewRateLimitedWorker(c repository.Repository, r rate_limit.RateLimit) *RateLimitedWorker {
	return &RateLimitedWorker{
		c: c,
		r: r,
	}
}

// StartWork will start to process work in infinite loop - this functions is blocking so call it in separate Go routine
func (w *RateLimitedWorker) StartWork(close chan interface{}) {
	logrus.Info("Starting background worker")
	for {
		select {
		case <-time.After(6 * time.Second):
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			logrus.Debug("isPermitted", w.r.IsPermitted())

			if !w.r.IsPermitted() {
				continue
			}

			job, err := w.c.GetUnfinishedJobsCtx(ctx)
			if errors.Is(err, &repository.EmptyResponseError{}) {
				logrus.Debug("No jobs to process")
				continue
			}
			if err != nil {
				logrus.Errorf("Error obtaining latest job from DB: %+v", err)
				continue
			}

			logrus.Debug("Found jobs to work on")

			res, err := w.r.Do(ctx, map[string]string{"message": job.Message})
			if err != nil {
				logrus.Errorf("Request fail: %+v", err)
				continue
			}

			s := res.(string)

			job.Finished = true
			job.Result = s
			err = w.c.UpdateJobCtx(ctx, job)
			if err != nil {
				logrus.Errorf("Error updating job to DB: %+v", err)
			}
		case <-close:
			return
		}
	}
}
