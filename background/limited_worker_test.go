package background

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"make-me-reliable/repository"
)

type mockRepository struct {
	CountGet    int
	CountUpdate int
}

func (m *mockRepository) GetJobByIdCtx(context.Context, primitive.ObjectID) (*repository.Job, error) {
	m.CountGet += 1
	return nil, nil
}

func (m *mockRepository) GetUnfinishedJobsCtx(context.Context) (*repository.Job, error) {
	panic("unimplemented")
}

func (m *mockRepository) InsertJobCtx(context.Context, *repository.Job) (*primitive.ObjectID, error) {
	panic("unimplemented")
}

func (m *mockRepository) UpdateJobCtx(context.Context, *repository.Job) error {
	panic("unimplemented")
}

type mockRateLimiter struct {
}

func (rl *mockRateLimiter) IsPermitted() bool {
	return true
}

func (rl *mockRateLimiter) Do(context.Context, map[string]string) (interface{}, error) {
	return nil, nil
}

func TestRateLimitedWorker_StartWork(t *testing.T) {
	closeChan := make(chan interface{})
	q := &mockRepository{}
	rl := &mockRateLimiter{}
	w := NewRateLimitedWorker(q, rl)

	go w.StartWork(closeChan)

	select {
	case <-time.After(10 * time.Second):
		close(closeChan)
	}

	if q.CountGet != 1 {
		t.Errorf("Expected 1, got %d", q.CountGet)
	}
}

/*
// StartWork will start to process work in infinite loop - this functions is blocking so call it in separate Go routine
func (w *RateLimitedWorker) StartWork() {
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
		}
	}
}




*/
