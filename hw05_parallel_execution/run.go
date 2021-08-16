package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrInvalidWorkersCount = errors.New("errors limit must be positive")
)

type (
	Task func() error
	Pool struct {
		Tasks          []Task
		workersCount   int
		maxErrorsCount int
		tasksChan      chan Task
		wg             sync.WaitGroup
		mu             sync.Mutex
		errCount       int
	}
)

func NewPool(tasks []Task, workersCount int, maxErrorsCount int) *Pool {
	return &Pool{
		Tasks:          tasks,
		workersCount:   workersCount,
		maxErrorsCount: maxErrorsCount,
		tasksChan:      make(chan Task),
	}
}

func (p *Pool) ErrorsLimitExceeded() bool {
	p.mu.Lock()
	errorsLimitExceeded := p.errCount >= p.maxErrorsCount
	p.mu.Unlock()
	return errorsLimitExceeded
}

func (p *Pool) Run() error {
	if p.workersCount <= 0 {
		return ErrInvalidWorkersCount
	}
	for i := 0; i < p.workersCount; i++ {
		go p.work()
	}
	for _, task := range p.Tasks {
		if p.ErrorsLimitExceeded() {
			break
		}
		p.wg.Add(1) // в этот момент другой воркер может вернуть последнюю ошибку
		p.tasksChan <- task

	}
	close(p.tasksChan)
	p.wg.Wait()
	if p.ErrorsLimitExceeded() {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func (p *Pool) work() {
	for task := range p.tasksChan {
		if err := task(); err != nil {
			p.mu.Lock()
			p.errCount++
			p.mu.Unlock()
		}
		p.wg.Done()
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	p := NewPool(tasks, n, m)
	return p.Run()
}
