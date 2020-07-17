package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/avegner/utils/workers"
)

const (
	tasksNumber = 12
	poolSize    = 4
	taskDelay   = 1 * time.Second
	waitDelay   = (tasksNumber/poolSize+1)*taskDelay + 100*time.Millisecond
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed: %v", err)
		os.Exit(1)
	}
}

func run() error {
	p := workers.NewPool(poolSize)
	defer p.Close()

	for i := 0; i < tasksNumber; i++ {
		n := i + 1
		if _, err := p.Dispatch(func(ctx context.Context) error {
			<-time.After(taskDelay)
			fmt.Printf("task %d done\n", n)
			return nil
		}); err != nil {
			return err
		}
	}

	time.Sleep(waitDelay)
	return nil
}
