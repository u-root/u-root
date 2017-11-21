// Copyright 2016 The Netstack Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tmutex provides the implementation of a mutex that implements an
// efficient TryLock function in addition to Lock and Unlock.
package tmutex

import (
	"sync/atomic"
)

// Mutex is a mutual exclusion primitive that implements TryLock in addition
// to Lock and Unlock.
type Mutex struct {
	v  int32
	ch chan struct{}
}

// Init initializes the mutex.
func (m *Mutex) Init() {
	m.v = 1
	m.ch = make(chan struct{}, 1)
}

// Lock acquires the mutex. If it is currently held by another goroutine, Lock
// will wait until it has a chance to acquire it.
func (m *Mutex) Lock() {
	// Uncontended case.
	if atomic.AddInt32(&m.v, -1) == 0 {
		return
	}

	for {
		// Try to acquire the mutex again, at the same time making sure
		// that m.v is negative, which indicates to the owner of the
		// lock that it is contended, which will force it to try to wake
		// someone up when it releases the mutex.
		if v := atomic.LoadInt32(&m.v); v >= 0 && atomic.SwapInt32(&m.v, -1) == 1 {
			return
		}

		// Wait for the mutex to be released before trying again.
		<-m.ch
	}
}

// TryLock attempts to acquire the mutex without blocking. If the mutex is
// currently held by another goroutine, it fails to acquire it and returns
// false.
func (m *Mutex) TryLock() bool {
	v := atomic.LoadInt32(&m.v)
	if v <= 0 {
		return false
	}
	return atomic.CompareAndSwapInt32(&m.v, 1, 0)
}

// Unlock releases the mutex.
func (m *Mutex) Unlock() {
	if atomic.SwapInt32(&m.v, 1) == 0 {
		// There were no pending waiters.
		return
	}

	// Wake some waiter up.
	select {
	case m.ch <- struct{}{}:
	default:
	}
}
