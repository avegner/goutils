package main

import (
	"fmt"
	"os"
	"time"

	"github.com/avegner/utils/stats"
)

func main() {
	foo()
	bar()
	func() {
		defer stats.Took("some operation")()
		time.Sleep(10 * time.Millisecond)
	}()
}

func foo() {
	defer stats.Took("foo")()
	time.Sleep(500 * time.Millisecond)
}

func bar() {
	defer stats.Took("bar", func(msg string) {
		fmt.Fprintf(os.Stderr, "%s\n", msg)
	})()
	time.Sleep(1 * time.Second)
}
