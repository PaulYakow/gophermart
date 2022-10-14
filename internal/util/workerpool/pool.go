package workerpool

import (
	"context"
)

type Pool struct {
	workers       []*Worker
	workersCount  int
	tasks         chan Task
	results       chan Result
	runBackground chan struct{}
}

func NewPool(workersCount int, tasksBuffer int) *Pool {
	if workersCount < 1 {
		workersCount = 1
	}

	return &Pool{
		workers:      make([]*Worker, workersCount),
		workersCount: workersCount,
		tasks:        make(chan Task, tasksBuffer),
		results:      make(chan Result, workersCount),
	}
}

// RunBackground запускает pool в фоне
func (p *Pool) RunBackground(ctx context.Context) {
	for idx := 1; idx <= p.workersCount; idx++ {
		worker := NewWorker(idx, p.tasks, p.results)
		p.workers[idx-1] = worker
		go worker.StartBackground(ctx)
	}

	p.runBackground = make(chan struct{})
	<-p.runBackground
}

// Stop останавливает запущенных в фоне worker-ов
func (p *Pool) Stop() {
	for idx := range p.workers {
		p.workers[idx].Stop()
	}

	p.runBackground <- struct{}{}

	close(p.results)
}

func (p *Pool) Results() <-chan Result {
	return p.results
}

// AddTask добавляет задачи в pool
func (p *Pool) AddTask(task Task) {
	p.tasks <- task
}

func (p *Pool) GenerateFrom(taskBatch ...Task) {
	for _, task := range taskBatch {
		p.tasks <- task
	}
}
