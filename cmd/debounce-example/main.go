// Package main demonstrates the usage of the debouncer.
package main

import (
	"fmt"
	"time"

	"github.com/udhos/debounce/debounce"
)

func main() {
	sequential()
	concurrent()
}

func sequential() {
	debouncer := debounce.New(100 * time.Millisecond)

	var executions int
	var last int

	for i := range 100 {
		debouncer.Run(func() {
			executions++
			last = i
		})
	}

	time.Sleep(1 * time.Second)
	debouncer.Stop()

	fmt.Printf("sequential: executions: %d (expected 1)\n", executions)
	fmt.Printf("sequential: last: %d (expected 0)\n", last)
}

func concurrent() {
	debouncer := debounce.New(100 * time.Millisecond)

	var executions int
	var last int

	for i := range 100 {
		go func(i int) {
			debouncer.Run(func() {
				executions++
				last = i
			})
		}(i)
	}

	time.Sleep(1 * time.Second)
	debouncer.Stop()

	fmt.Printf("concurrent: executions: %d (expected 1)\n", executions)
	fmt.Printf("concurrent: last: %d (expected random value between 0 and 99)\n", last)
}
