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

	run(calls, false)

	fmt.Println()

	run(calls, true)
}

func run(calls int, concurrent bool) {
	debouncer := debounce.New(100 * time.Millisecond)

	var executions int
	var index int

	for i := range calls {
		if concurrent {
			go func(i int) {
				debouncer.Run(func() {
					executions++
					index = i
				})
			}(i)
		} else {
			debouncer.Run(func() {
				executions++
				index = i
			})
		}
	}

	time.Sleep(1 * time.Second)

	expected := 2
	if calls == 1 {
		expected = 1
	}

	var label string
	if concurrent {
		label = "concurrent"
	} else {
		label = "sequential"
	}

	fmt.Printf("%s: executions: %d (expected %d)\n", label, executions, expected)
	fmt.Printf("%s: index: %d (expected random value between 0 and %d)\n", label, index, calls-1)
}
