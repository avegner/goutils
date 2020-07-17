package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/avegner/utils/workers"
)

const (
	taskDelay = 1 * time.Second
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	p := workers.NewPool()
	defer p.Close()

	tok, err := p.Dispatch(func(ctx context.Context) error {
		<-time.After(taskDelay)
		fmt.Printf("task done\n")
		return nil
	})
	if err != nil {
		return err
	}

	err = tok.WaitWithContext(context.Background())
	fmt.Printf("task returned: %v\n", err)

	return nil
}
