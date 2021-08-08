package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("every n-th task has error", func(t *testing.T) {
		tasksCount := 50 // number of task
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32
		n := rand.Intn(10) + 1 // every n-th task has error
		for i := 0; i < tasksCount; i++ {
			var err error
			if i%n == n-1 {
				err = fmt.Errorf("error from task %d", i)
			}
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(10))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := rand.Intn(4) + 1
		maxErrorsCount := rand.Intn(3) + 1
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		// если бы у нас все таски выполнялись последовательно, то выполнилось бы maxErrorsCount*n - потому что каждый n-ый таск генерит ошибку
		// у нас же исполнение конкурентно - это значит что число запущенных тасков может быть скорректировано на workersCount (если уж совсем упарываться то на workersCount-1)
		// важно еще и то что время выполнения у нас фиксировано - иначе может быть одна длинная задача с ошибкой, а остальные воркеры все делают
		// на самом деле это обощение первого теста где по сути n=1 и число тасков не более 1*maxErrorsCount + workersCount
		// простите что на русском, но на английском пояснить этот не самый тривиальный момент не смог
		require.LessOrEqual(t, runTasksCount, int32(n*maxErrorsCount+workersCount-1), "extra tasks were started")
		require.GreaterOrEqual(t, runTasksCount, int32(n*maxErrorsCount-workersCount+1), "some tasks were not started")
	})

	t.Run("Zero errors (no one task started)", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32
		n := rand.Intn(10) + 1
		for i := 0; i < tasksCount; i++ {
			err := error(nil)
			if i%n == n-1 {
				err = fmt.Errorf("error from task %d", i)
			}
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(10))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := rand.Intn(4) + 1
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.Equal(t, runTasksCount, int32(0), "tasks were started")
	})
}
