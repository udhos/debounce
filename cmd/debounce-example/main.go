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
	var index int

	for i := range 100 {
		debouncer.Run(func() {
			executions++
			index = i
		})
	}

	time.Sleep(1 * time.Second)

	fmt.Printf("sequential: executions: %d (expected 1)\n", executions)
	fmt.Printf("sequential: index: %d (expected 99)\n", index)
}

func concurrent() {
	debouncer := debounce.New(100 * time.Millisecond)

	var executions int
	var index int

	for i := range 100 {
		go func(i int) {
			debouncer.Run(func() {
				executions++
				index = i
			})
		}(i)
	}

	time.Sleep(1 * time.Second)

	fmt.Printf("concurrent: executions: %d (expected 1)\n", executions)
	fmt.Printf("concurrent: index: %d (expected random value between 0 and 99)\n", index)
}
