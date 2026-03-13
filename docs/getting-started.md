# Getting Started

`go-worker` provides a small fixed-size pool that runs submitted tasks and shuts down gracefully.

## Lifecycle

- `worker.NewPool(n)` starts `n` workers. Non-positive values are normalized to one worker.
- `Submit` accepts work while the pool is open.
- `Shutdown` closes the queue, waits for in-flight work to finish, and then releases the shared context.

## Shutdown semantics

Shutdown is graceful: already queued tasks are allowed to complete, while future submissions are rejected.
