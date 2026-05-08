package debounce

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestDebouncer_Throttling(t *testing.T) {
	delay := 100 * time.Millisecond
	d := New(delay)

	var calls atomic.Int32

	// Fire run() 10 times in a tight loop.
	for range 10 {
		d.Run(func() {
			calls.Add(1)
		})
	}

	// The first call executes immediately.
	// The subsequent 9 calls are coalesced and will execute after the delay.
	if c := calls.Load(); c != 1 {
		t.Errorf("expected 1 call immediately, got %d", c)
	}

	// Wait for the delay to pass
	time.Sleep(delay + 50*time.Millisecond)

	// We should have exactly 2 calls executed (1 immediate, 1 delayed)
	if c := calls.Load(); c != 2 {
		t.Errorf("expected exactly 2 calls after delay, got %d", c)
	}

	// If we run again after the window, it should execute immediately!
	d.Run(func() {
		calls.Add(1)
	})

	if c := calls.Load(); c != 3 {
		t.Errorf("expected exactly 3 calls immediately after new run, got %d", c)
	}

	time.Sleep(delay + 50*time.Millisecond)

	// Still 3, because no subsequent calls were made during the window
	if c := calls.Load(); c != 3 {
		t.Errorf("expected exactly 3 calls total, got %d", c)
	}
}

func TestDebouncer_ContinuousUpdates(t *testing.T) {
	delay := 100 * time.Millisecond
	d := New(delay)

	var calls atomic.Int32

	done := make(chan struct{})
	go func() {
		// send updates every 20ms for 350ms total
		for range 17 {
			d.Run(func() {
				calls.Add(1)
			})
			time.Sleep(20 * time.Millisecond)
		}
		close(done)
	}()

	<-done
	time.Sleep(delay + 50*time.Millisecond) // wait for the last timer to potentially fire

	// A pure debounce would never fire (0 calls) because 20ms < 100ms delay.
	// No debouncing would fire 17 times.
	// With immediate-first throttling (without cooldown), we expect it to fire twice per interval:
	// T=0 (immediate), T=100 (delayed)
	// T=100 (immediate), T=200 (delayed)
	// T=200 (immediate), T=300 (delayed)
	// T=300 (immediate), T=400 (delayed)
	// Total expected: ~7-8 calls.
	c := calls.Load()
	if c < 6 || c > 9 {
		t.Errorf("expected between 6 and 9 calls for continuous updates, got %d", c)
	}
}

func TestDebouncer_LastFunctionCalled(t *testing.T) {
	delay := 100 * time.Millisecond
	d := New(delay)

	var lastCalled atomic.Int32
	var callCount atomic.Int32

	d.Run(func() {
		lastCalled.Store(1)
		callCount.Add(1)
	})

	d.Run(func() {
		lastCalled.Store(2)
		callCount.Add(1)
	})

	d.Run(func() {
		lastCalled.Store(3)
		callCount.Add(1)
	})

	// Immediate execution should have happened for the first call
	if val := lastCalled.Load(); val != 1 {
		t.Errorf("expected exactly the first function to be called immediately (value 1), got %d", val)
	}
	if count := callCount.Load(); count != 1 {
		t.Errorf("expected exactly 1 function to be executed immediately, got %d", count)
	}

	time.Sleep(delay + 50*time.Millisecond)

	// After delay, the last function (3) should have been called
	if val := lastCalled.Load(); val != 3 {
		t.Errorf("expected exactly the last function to be called after delay (value 3), got %d", val)
	}
	if count := callCount.Load(); count != 2 {
		t.Errorf("expected exactly 2 functions to be executed total, got %d", count)
	}
}
