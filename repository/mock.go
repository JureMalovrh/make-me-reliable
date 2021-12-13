package repository

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//MockRepository is a mock repository implementation. It is safe to be used in multiple go routines
type MockRepository struct {
	CountGet          int
	CountUpdate       int
	CountUnfinished   int
	CountInsert       int
	GetJobByIdFunc    func() (*Job, error)
	InsertJobFunc     func() (*primitive.ObjectID, error)
	UpdateJobFunc     func() error
	GetUnfinishedFunc func() (*Job, error)
	rw                sync.RWMutex
}

// GetJobByIdCtx will call GetJobByIdFunc from MockRepository
func (m *MockRepository) GetJobByIdCtx(context.Context, primitive.ObjectID) (*Job, error) {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.CountGet += 1
	return m.GetJobByIdFunc()
}

// GetUnfinishedJobsCtx will call GetUnfinishedFunc from MockRepository
func (m *MockRepository) GetUnfinishedJobsCtx(context.Context) (*Job, error) {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.CountUnfinished += 1
	return m.GetUnfinishedFunc()
}

// InsertJobCtx will call InsertJobFunc from MockRepository
func (m *MockRepository) InsertJobCtx(context.Context, *Job) (*primitive.ObjectID, error) {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.CountInsert += 1
	return m.InsertJobFunc()
}

// UpdateJobCtx will call UpdateJobFunc from MockRepository
func (m *MockRepository) UpdateJobCtx(context.Context, *Job) error {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.CountUpdate += 1
	return m.UpdateJobFunc()
}

func (m *MockRepository) GetCountGet() int {
	m.rw.Lock()
	defer m.rw.Unlock()
	return m.CountGet
}

func (m *MockRepository) GetCountUpdate() int {
	m.rw.Lock()
	defer m.rw.Unlock()
	return m.CountUpdate
}

func (m *MockRepository) GetCountUnfinished() int {
	m.rw.Lock()
	defer m.rw.Unlock()
	return m.CountUnfinished
}

func (m *MockRepository) GetCountInsert() int {
	m.rw.Lock()
	defer m.rw.Unlock()
	return m.CountInsert
}
