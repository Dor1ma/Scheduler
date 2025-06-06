package app

import "github.com/Dor1ma/Scheduler/internal/timer"

type Dispatcher struct {
	tasks chan timer.Task
}

func NewDispatcher(bufferSize int) *Dispatcher {
	return &Dispatcher{tasks: make(chan timer.Task, bufferSize)}
}

func (d *Dispatcher) Dispatch(task timer.Task) {
	d.tasks <- task
}

func (d *Dispatcher) Tasks() <-chan timer.Task {
	return d.tasks
}
