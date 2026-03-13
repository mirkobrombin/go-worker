package main

import (
	"context"
	"fmt"

	"github.com/mirkobrombin/go-worker/pkg/worker"
)

func main() {
	pool := worker.NewPool(2)
	defer pool.Shutdown()

	ok := pool.Submit(func(ctx context.Context) error {
		fmt.Println("task executed")
		return nil
	})
	if !ok {
		panic("task was not accepted")
	}
}
