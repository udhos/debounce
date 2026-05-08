[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/udhos/debounce/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/udhos/debounce)](https://goreportcard.com/report/github.com/udhos/debounce)
[![Go Reference](https://pkg.go.dev/badge/github.com/udhos/debounce.svg)](https://pkg.go.dev/github.com/udhos/debounce)

# debounce

[debounce](https://github.com/udhos/debounce) is a simple implementation of a debouncer.

The debouncer coalesces multiple function calls.

The first function call is executed immediately, and a delay window is started.

Subsequent function calls within the delay window will not reset the timer, but will update the function to be executed.

After the delay has passed, only the most recently provided function will be executed, returning the debouncer to its initial state.

The debouncer is thread-safe, and can be used in concurrent environments. But it is often useful in sequential non-concurrent code as well.

A common use case is to reduce the number of reconciliation operations in response to events. For example, in Kubernetes controllers, events may be sent in bursts, and it is often desirable to only reconcile once after many events have been received.

# Synopsis

In the example below, we create a debouncer with a delay of 100 milliseconds.

We then run the function 100 times in a loop.

Since all calls happen within the delay window, the function will only be executed once.

Thus the output will be:

```bash
executions: 1
```

```go
package main

import (
	"fmt"
	"time"

	"github.com/udhos/debounce/debounce"
)

func main() {
	debouncer := debounce.New(100 * time.Millisecond)

	var executions int

	myFunc := func() {
		executions++
	}

	for range 100 {
		debouncer.Run(myFunc)
	}

	time.Sleep(1 * time.Second)

	fmt.Printf("executions: %d\n", executions)
}
```
