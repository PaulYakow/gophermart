package workerpool

import (
	"context"
)

type Worker struct {
	taskCh   <-chan Task
	resultCh chan<- Result
	quit     chan struct{}
}

func NewWorker(ID int, tasks <-chan Task, results chan<- Result) *Worker {
	return &Worker{
		taskCh:   tasks,
		resultCh: results,
		quit:     make(chan struct{}),
	}
}

// StartBackground запускает worker-а в фоне
func (w *Worker) StartBackground(ctx context.Context) {
	for {
		select {
		case task := <-w.taskCh:
			w.resultCh <- task.process(ctx)
		case <-w.quit:
			return
		case <-ctx.Done():
			return
		}
	}
}

// Stop останавливает quit для worker-а
func (w *Worker) Stop() {
	go func() {
		w.quit <- struct{}{}
	}()
}
