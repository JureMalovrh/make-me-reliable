package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"make-me-reliable/repository"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func newJobHandler(rdb repository.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		jobId := vars["jobId"]

		objectId, err := primitive.ObjectIDFromHex(jobId)
		if err != nil {
			responseWithStatus(w, http.StatusBadRequest, errMsg{"JobID is not a valid ID"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		job, err := rdb.GetJobByIdCtx(ctx, objectId)
		if errors.Is(err, &repository.EmptyResponseError{}) {
			responseWithStatus(w, http.StatusNotFound, errMsg{"JobID is not found"})
			return
		}
		if err != nil {
			responseWithStatus(w, http.StatusInternalServerError, errMsg{"Something is wrong"})
			return
		}

		responseWithStatus(w, http.StatusOK, job)
	}
}
