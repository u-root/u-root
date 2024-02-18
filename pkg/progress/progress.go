// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The progress package starts printing the status of a progress every second
// It takes a pointer to an int64 which holds the amount of
// bytes copied - and prints it.

package progress

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

// ProgressData contains information about a progress printer
type ProgressData struct {
	mode         string // one of: none, xfer, progress
	start        time.Time
	end          time.Time
	endTimeMutex sync.Mutex
	variable     *int64 // must be aligned for atomic operations
	quit         chan struct{}
	w            io.Writer
}

// New creates a new progress struct
func New(w io.Writer, mode string, variable *int64) *ProgressData {
	return &ProgressData{
		mode:         mode,
		start:        time.Now(),
		endTimeMutex: sync.Mutex{},
		variable:     variable,
		w:            w,
	}
}

// Begin begins a progress routine
//
// mode describes in which mode it runs, none, progress or xfer
// variable holds the amount of bytes written
func (p *ProgressData) Begin() {
	if p.mode == "progress" {
		p.print()
		// Print progress in a separate goroutine.
		p.quit = make(chan struct{}, 1)
		go func() {
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					p.print()
				case <-p.quit:
					p.endTimeMutex.Lock()
					p.end = time.Now()
					p.endTimeMutex.Unlock()
					return
				}
			}
		}()
	}
}

// End - Ends the progress and send quit signal to the channel
func (p *ProgressData) End() {
	if p.mode == "progress" {
		// Properly synchronize goroutine.
		p.quit <- struct{}{}
		p.quit <- struct{}{}
	} else {
		p.endTimeMutex.Lock()
		p.end = time.Now()
		p.endTimeMutex.Unlock()
	}
	if p.mode == "progress" || p.mode == "xfer" {
		// Print grand total.
		p.print("\n")
	}
}

// print prints out progress information and any
// extra strings needed at the end.
// With "status=progress", this is called from 3 places:
// - Once at the beginning to appear responsive
// - Every 1s afterwards
// - Once at the end so the final value is accurate
func (p *ProgressData) print(extra ...string) {
	elapse := time.Since(p.start)
	n := atomic.LoadInt64(p.variable)
	d := float64(n)
	const mib = 1024 * 1024
	const mb = 1000 * 1000
	// The ANSI escape may be undesirable to some eyes.
	if p.mode == "progress" {
		p.w.Write([]byte("\033[2K\r"))
	}
	fmt.Fprintf(p.w, "%d bytes (%.3f MB, %.3f MiB) copied, %.3f s, %.3f MB/s",
		n, d/mb, d/mib, elapse.Seconds(), float64(d)/elapse.Seconds()/mb)
	for _, s := range extra {
		fmt.Fprint(p.w, s)
	}
}
