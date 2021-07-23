package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type Pool struct {
	Tasks          []Task
	workersCount   int
	maxErrorsCount int
	tasksChan      chan Task
	wg             sync.WaitGroup
	mu             sync.Mutex
	errCount       int
}

func NewPool(tasks []Task, workersCount int, maxErrorsCount int) *Pool {
	return &Pool{
		Tasks:          tasks,
		workersCount:   workersCount,
		maxErrorsCount: maxErrorsCount,
		tasksChan:      make(chan Task),
	}
}

func (p *Pool) Run() error {
	for i := 0; i < p.workersCount; i++ {
		go p.work()
	}

	for _, task := range p.Tasks {
		// put task im queue if not enough errors

		p.wg.Add(1)
		p.tasksChan <- task

	}

	close(p.tasksChan)

	p.wg.Wait()
	if p.errCount >= p.maxErrorsCount {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func (p *Pool) work() {
	for task := range p.tasksChan {
		p.mu.Lock()
		continueAdd := p.errCount < p.maxErrorsCount
		p.mu.Unlock()
		if continueAdd {
			if err := task(); err != nil {
				p.mu.Lock()
				p.errCount++
				p.mu.Unlock()

			}
		}
		p.wg.Done()
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	p := NewPool(tasks, n, m)
	return p.Run()
}
