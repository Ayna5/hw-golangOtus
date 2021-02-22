package hw05_parallel_execution //nolint:golint,stylecheck,revive

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Workers struct {
	countGoroutines int
	maxErrCount     int
	tasks           []Task
	mu              *sync.Mutex
	completedTask   int
	errCount        int
}

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, N int, M int) error {
	// Place your code here
	workers := Workers{
		countGoroutines: N,
		maxErrCount:     M,
		tasks:           tasks,
		mu:              &sync.Mutex{},
	}
	err := workers.RunTasks()
	return err
}

func (w *Workers) RunTasks() error {
	var wg sync.WaitGroup
	taskChan := make(chan Task, len(w.tasks))
	errChan := make(chan bool)
	_, cancel := context.WithCancel(context.Background())

	for _, task := range w.tasks {
		taskChan <- task
	}

	for i := 0; i < w.countGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				if w.errLimitCheck() {
					err := task()
					errChan <- err == nil
				}
			}
		}()
	}

	go func() {
		for ok := range errChan {
			w.completedTask++
			if !ok {
				w.mu.Lock()
				w.errCount++
				w.mu.Unlock()
				if !w.errLimitCheck() {
					cancel()
				}
			}
			if w.completedTask == len(w.tasks) {
				cancel()
			}
		}
	}()

	close(taskChan)
	wg.Wait()
	close(errChan)

	if !w.errLimitCheck() {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func (w *Workers) errLimitCheck() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.maxErrCount <= 0 || w.errCount >= w.maxErrCount {
		return false
	}
	return true
}
