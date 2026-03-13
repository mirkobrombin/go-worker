# Go Worker

A compact **worker pool** for Go applications that want predictable task execution with graceful shutdown semantics.

## Features

- **Fixed Worker Pool:** Start a pool with a stable number of workers.
- **Graceful Shutdown:** Stop accepting new tasks and wait for queued work to finish.
- **Context-Aware Tasks:** Pass a shared context into each task function.
- **Small API Surface:** Keep submission and shutdown logic easy to reason about.

## Installation

```bash
go get github.com/mirkobrombin/go-worker
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"

    "github.com/mirkobrombin/go-worker/pkg/worker"
)

func main() {
    pool := worker.NewPool(2)
    defer pool.Shutdown()

    pool.Submit(func(ctx context.Context) error {
        fmt.Println("working")
        return nil
    })
}
```

## Documentation

- [Getting Started](docs/getting-started.md)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
