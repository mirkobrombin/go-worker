package worker

import (
	"context"
	"sync"
)

// Task is the unit of work executed by the pool.
type Task func(context.Context) error

// Pool is a fixed-size worker pool with graceful shutdown.
type Pool struct {
	wg     sync.WaitGroup
	tasks  chan Task
	done   chan struct{}
	cancel context.CancelFunc
	ctx    context.Context
	once   sync.Once
}

// NewPool creates a pool with n workers.
func NewPool(n int) *Pool {
	if n <= 0 {
		n = 1
	}

	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{
		tasks:  make(chan Task),
		done:   make(chan struct{}),
		cancel: cancel,
		ctx:    ctx,
	}
	for i := 0; i < n; i++ {
		p.wg.Add(1)
		go p.worker()
	}
	return p
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case task := <-p.tasks:
			_ = task(p.ctx)
		case <-p.done:
			return
		}
	}
}

// Submit enqueues a task for execution. Returns false if the pool is shut down
// or shutting down; the task is not executed in that case.
func (p *Pool) Submit(task Task) bool {
	select {
	case p.tasks <- task:
		return true
	case <-p.done:
		return false
	}
}

// Shutdown stops accepting new tasks and waits for in-progress work to finish.
// The task channel is unbuffered, so no tasks can be pending pickup when done
// is closed — any Submit in progress will observe done and return false.
func (p *Pool) Shutdown() {
	p.once.Do(func() {
		close(p.done)
		p.wg.Wait()
		p.cancel()
	})
}
