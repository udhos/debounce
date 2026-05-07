// Package debounce implements a debouncer.
package debounce

import (
	"sync"
	"time"
)

// Debouncer is a simple implementation of a debounce.
// It coalesces multiple calls to run() within a certain delay window.
// Only the last call is executed after the delay.
type Debouncer struct {
	delay time.Duration
	mutex sync.Mutex
	timer *time.Timer
	f     func()
}

// New creates a new Debouncer.
func New(delay time.Duration) *Debouncer {
	return &Debouncer{delay: delay}
}

// Run schedules a function to be executed after the delay.
// If the function is called multiple times within the delay window,
// only the last call is executed.
func (d *Debouncer) Run(f func()) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.f = f

	if d.timer != nil {
		return
	}

	d.timer = time.AfterFunc(d.delay, func() {
		d.mutex.Lock()
		f := d.f
		d.timer = nil
		d.mutex.Unlock()

		f()
	})
}

// Stop stops the debouncer.
func (d *Debouncer) Stop() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if d.timer != nil {
		d.timer.Stop()
	}
}
