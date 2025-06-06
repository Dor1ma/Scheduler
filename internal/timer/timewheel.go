package timer

import (
	"sync"
	"time"
)

type slot struct {
	tasks []Task
	lock  sync.Mutex
}

type TimeWheel struct {
	interval time.Duration
	slots    []slot
	current  int
	ticker   *time.Ticker
	stopCh   chan struct{}
	handler  func(Task)
	size     int
}

func NewTimeWheel(interval time.Duration, slotCount int, handler func(Task)) *TimeWheel {
	slots := make([]slot, slotCount)
	return &TimeWheel{
		interval: interval,
		slots:    slots,
		ticker:   time.NewTicker(interval),
		stopCh:   make(chan struct{}),
		handler:  handler,
		size:     slotCount,
	}
}

func (tw *TimeWheel) Start() {
	go func() {
		for {
			select {
			case <-tw.ticker.C:
				tw.tick()
			case <-tw.stopCh:
				return
			}
		}
	}()
}

func (tw *TimeWheel) Stop() {
	close(tw.stopCh)
	tw.ticker.Stop()
}

func (tw *TimeWheel) Add(task Task) {
	delay := task.ExecuteAt.Sub(time.Now())
	if delay < 0 {
		go tw.handler(task)
		return
	}
	ticks := int(delay / tw.interval)
	slotIndex := (tw.current + ticks) % tw.size

	slot := &tw.slots[slotIndex]
	slot.lock.Lock()
	slot.tasks = append(slot.tasks, task)
	slot.lock.Unlock()
}

func (tw *TimeWheel) tick() {
	slot := &tw.slots[tw.current]

	slot.lock.Lock()
	tasks := slot.tasks
	slot.tasks = nil
	slot.lock.Unlock()

	now := time.Now()
	for _, task := range tasks {
		if task.ExecuteAt.Before(now) || task.ExecuteAt.Equal(now) {
			go tw.handler(task)
		} else {
			tw.Add(task)
		}
	}

	tw.current = (tw.current + 1) % tw.size
}
