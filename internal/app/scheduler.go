package app

import (
	"github.com/Dor1ma/Scheduler/internal/storage"
	"github.com/Dor1ma/Scheduler/internal/timer"
	"time"
)

type Scheduler struct {
	tw    *timer.TimeWheel
	store storage.TaskStorage
}

func NewScheduler(tw *timer.TimeWheel, store storage.TaskStorage) *Scheduler {
	return &Scheduler{tw: tw, store: store}
}

func (s *Scheduler) Start() error {
	s.tw.Start()

	tasks, err := s.store.LoadPending(time.Now().Add(5 * time.Minute))
	if err != nil {
		return err
	}

	for _, t := range tasks {
		s.tw.Add(t)
	}

	return nil
}

func (s *Scheduler) Stop() {
	s.tw.Stop()
}

func (s *Scheduler) ScheduleTask(task timer.Task) error {
	if err := s.store.Save(task); err != nil {
		return err
	}
	s.tw.Add(task)
	return nil
}

func (s *Scheduler) CancelTask(id string) error {
	return s.store.Delete(id)
}
