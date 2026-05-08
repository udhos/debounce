// Package main demonstrates the usage of the debouncer.
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/udhos/debounce/debounce"
)

func main() {

	var calls int
	flag.IntVar(&calls, "calls", 100, "number of calls to make to the debouncer")
	flag.Parse()

	if calls < 1 {
		fmt.Println("calls must be greater than 0")
		return
	}

	fmt.Printf("Making %d calls to the debouncer (set with -calls <number>)\n", calls)
	fmt.Println()

	sequential(calls)
	concurrent(calls)
}

func sequential(calls int) {
	debouncer := debounce.New(100 * time.Millisecond)

	var executions int
	var index int

	for i := range calls {
		debouncer.Run(func() {
			executions++
			index = i
		})
	}

	time.Sleep(1 * time.Second)

	expected := 2
	if calls == 1 {
		expected = 1
	}

	fmt.Printf("sequential: executions: %d (expected %d)\n", executions, expected)
	fmt.Printf("sequential: index: %d (expected %d)\n", index, calls-1)
}

func concurrent(calls int) {
	debouncer := debounce.New(100 * time.Millisecond)

	var executions int
	var index int

	for i := range calls {
		go func(i int) {
			debouncer.Run(func() {
				executions++
				index = i
			})
		}(i)
	}

	time.Sleep(1 * time.Second)

	expected := 2
	if calls == 1 {
		expected = 1
	}

	fmt.Printf("concurrent: executions: %d (expected %d)\n", executions, expected)
	fmt.Printf("concurrent: index: %d (expected random value between 0 and %d)\n", index, calls-1)
}
