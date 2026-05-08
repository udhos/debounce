// Package debounce implements a debouncer.
package debounce

import (
	"sync"
	"time"
)

// Debouncer is a simple implementation of a debounce.
// The first call is executed immediately.
// It coalesces multiple subsequent calls to Run within a certain delay window.
// Only the last call is executed after the delay.
type Debouncer struct {
	delay time.Duration
	mutex sync.Mutex
	timer *time.Timer
	f     func()
}

// New creates a Debouncer.
func New(delay time.Duration) *Debouncer {
	return &Debouncer{delay: delay}
}

// Run schedules a function to be executed.
// The first call is executed immediately.
// If the function is called multiple times within the delay window,
// only the last call is executed after the delay.
func (d *Debouncer) Run(f func()) {
	d.mutex.Lock()

	if d.timer == nil {
		// Schedule the delay timer on the first call.
		d.timer = time.AfterFunc(d.delay, d.timerFired)
		d.mutex.Unlock()

		// First call is executed immediately.
		f()
		return
	}

	// Subsequent calls are coalesced and only the last one
	// is executed when the timer fires.
	d.f = f
	d.mutex.Unlock()
}

// timerFired is called when the timer fires.
// It executes the last function that was called, or none
// if there was only the first immediate call.
// It resets the debounce timer.
func (d *Debouncer) timerFired() {
	d.mutex.Lock()
	f := d.f
	d.f = nil
	d.timer = nil
	d.mutex.Unlock()

	if f != nil {
		f()
	}
}
