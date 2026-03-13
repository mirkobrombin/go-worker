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
	cancel context.CancelFunc
	ctx    context.Context

	mu     sync.RWMutex
	closed bool
	once   sync.Once
}

// NewPool creates a pool with n workers.
func NewPool(n int) *Pool {
	if n <= 0 {
		n = 1
	}

	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{tasks: make(chan Task), cancel: cancel, ctx: ctx}
	for i := 0; i < n; i++ {
		p.wg.Add(1)
		go p.worker()
	}
	return p
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for task := range p.tasks {
		_ = task(p.ctx)
	}
}

// Submit enqueues a task for execution and returns false when the pool is closed.
func (p *Pool) Submit(task Task) bool {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return false
	}
	tasks := p.tasks
	p.mu.RUnlock()

	tasks <- task
	return true
}

// Shutdown stops accepting new tasks, drains queued work, and waits for workers.
func (p *Pool) Shutdown() {
	p.once.Do(func() {
		p.mu.Lock()
		p.closed = true
		close(p.tasks)
		p.mu.Unlock()

		p.wg.Wait()
		p.cancel()
	})
}
