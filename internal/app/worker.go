package app

import (
	"bytes"
	"fmt"
	"github.com/Dor1ma/Scheduler/internal/timer"
	"github.com/Dor1ma/Scheduler/pkg/logger"
	"net/http"

	"time"
)

type Worker struct {
	id         int
	dispatcher *Dispatcher
}

func NewWorker(id int, dispatcher *Dispatcher) *Worker {
	return &Worker{id: id, dispatcher: dispatcher}
}

func (w *Worker) Start() {
	go func() {
		for task := range w.dispatcher.Tasks() {
			if err := w.execute(task); err != nil {
				logger.Info("[Worker %d] error executing task %s: %v\n", w.id, task.ID, err)
			}
		}
	}()
}

func (w *Worker) execute(task timer.Task) error {
	logger.Info("[Worker %d] waiting for task %s\n", w.id, task.ID)
	req, err := http.NewRequest(task.Method, task.URL, bytes.NewReader(task.Payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("non-2xx response: %d", resp.StatusCode)
	}

	logger.Info("[Worker %d] task %s executed successfully\n", w.id, task.ID)
	return nil
}
