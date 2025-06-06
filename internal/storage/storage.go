package storage

import (
	"github.com/Dor1ma/Scheduler/internal/timer"
	"time"
)

//go:generate mockgen -source=storage.go -destination=storage_mock.go -package=storage
type TaskStorage interface {
	Save(task timer.Task) error
	Delete(id string) error
	LoadPending(until time.Time) ([]timer.Task, error)
}
