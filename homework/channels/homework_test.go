package main

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

// buffer size: workersNumber * bufferSizeMultiply
const bufferSizeMultiply = 2

// Error returned when the pool cannot accept more tasks
var ErrPoolFull = errors.New("pool is full")

// Error returned when the pool is closed
var ErrPoolClosed = errors.New("pool is closed")

type WorkerPool struct {
	workersNumber int
	buff          chan func()
	// Use sync.WaitGroup to wait for all tasks to complete
	workerGroup *sync.WaitGroup
	// Use a channel to signal when the pool is closed
	close chan struct{}
}

func NewWorkerPool(workersNumber int) *WorkerPool {
	wp := &WorkerPool{
		workersNumber: workersNumber,
		buff:          make(chan func(), workersNumber*bufferSizeMultiply),
		workerGroup:   &sync.WaitGroup{},
		close:         make(chan struct{}),
	}
	wp.workerGroup.Add(workersNumber)
	for i := 0; i < wp.workersNumber; i++ {
		go func() {
			defer wp.workerGroup.Done()
			for task := range wp.buff {
				// We can catch the potential panic here, but I don't implement it for simplicity
				task()
			}
		}()
	}
	return wp
}

// Return an error if the pool is full
func (wp *WorkerPool) AddTask(task func()) error {
	select {
	case <-wp.close:
		return ErrPoolClosed
	default:
	}

	select {
	case wp.buff <- task:
		return nil
	default:
		return ErrPoolFull
	}
}

// Shutdown all workers and workerGroup for all
// tasks in the pool to complete
func (wp *WorkerPool) Shutdown() {
	// Use select to check if close channel is already closed
	select {
	case <-wp.close:
		// Pool already shut down, do nothing
		return
	default:
		// Signal that the pool is closed to prevent new tasks
		close(wp.close)
	}
	// Now close the buffer and wait for workers to complete
	close(wp.buff)
	wp.workerGroup.Wait()
}

func TestWorkerPool(t *testing.T) {
	var counter atomic.Int32
	task := func() {
		time.Sleep(time.Millisecond * 500)
		counter.Add(1)
	}

	pool := NewWorkerPool(2)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(2), counter.Load())

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(3), counter.Load())

	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	pool.Shutdown() // wait tasks

	assert.Equal(t, int32(6), counter.Load())
}

func TestWorkerPoolFullCapacity(t *testing.T) {
	pool := NewWorkerPool(1)
	_ = pool.AddTask(func() {
		time.Sleep(time.Second * 1)
	})
	_ = pool.AddTask(func() {
		time.Sleep(time.Second * 1)
	})
	err := pool.AddTask(func() {
		time.Sleep(time.Second * 1)
	})
	assert.Equal(t, ErrPoolFull, err)
	pool.Shutdown()
	assert.True(t, true, "All tasks completed")
}

func TestWorkerPoolMultiplyShutdown(t *testing.T) {
	pool := NewWorkerPool(1)
	_ = pool.AddTask(func() {
		time.Sleep(time.Second * 1)
	})
	pool.Shutdown()
	pool.Shutdown()
	pool.Shutdown()
	pool.Shutdown()

	assert.True(t, true, "All tasks completed")
}

func TestWorkerPoolAddTaskAfterShutdown(t *testing.T) {
	pool := NewWorkerPool(1)
	_ = pool.AddTask(func() {
	})
	pool.Shutdown()
	err := pool.AddTask(func() {})

	assert.Equal(t, ErrPoolClosed, err)
}
