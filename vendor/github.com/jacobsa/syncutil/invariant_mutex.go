// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package syncutil

import (
	"sync"
	"sync/atomic"
)

var gEnable uintptr

// Enable checking of invariants when locking and unlocking InvariantMutex.
func EnableInvariantChecking() {
	atomic.StoreUintptr(&gEnable, 1)
}

// Has EnableInvariantChecking previously been called?
func InvariantCheckingEnabled() bool {
	return atomic.LoadUintptr(&gEnable) != 0
}

// A sync.Locker that, when enabled, runs a check for registered invariants at
// times when invariants should hold. This can aid debugging subtle code by
// crashing early as soon as something unexpected happens.
//
// Must be created with NewInvariantMutex. See that function for more details.
//
// A typical use looks like this:
//
//     type myStruct struct {
//       mu syncutil.InvariantMutex
//
//       // INVARIANT: nextGeneration == currentGeneration + 1
//       currentGeneration int // GUARDED_BY(mu)
//       nextGeneration    int // GUARDED_BY(mu)
//     }
//
//     // The constructor function for myStruct sets up the mutex to
//     // call the checkInvariants method.
//     func newMyStruct() *myStruct {
//       s := &myStruct{
//         currentGeneration: 1,
//         nextGeneration:    2,
//       }
//
//       s.mu = syncutil.NewInvariantMutex(s.checkInvariants)
//       return s
//     }
//
//     func (s *myStruct) checkInvariants() {
//       if s.nextGeneration != s.currentGeneration+1 {
//         panic(
//           fmt.Sprintf("%v != %v + 1", s.nextGeneration, s.currentGeneration))
//       }
//     }
//
//     // When EnableInvariantChecking has been called, invariants will be
//     // checked at entry to and exit from this function.
//     func (s *myStruct) setGeneration(n int) {
//       s.mu.Lock()
//       defer s.mu.Unlock()
//
//       s.currentGeneration = n
//       s.nextGeneration = n + 1
//     }
//
type InvariantMutex struct {
	mu    sync.RWMutex
	check func()
}

func (i *InvariantMutex) Lock() {
	i.mu.Lock()
	i.checkIfEnabled()
}

func (i *InvariantMutex) Unlock() {
	i.checkIfEnabled()
	i.mu.Unlock()
}

func (i *InvariantMutex) RLock() {
	i.mu.RLock()
	i.checkIfEnabled()
}

func (i *InvariantMutex) RUnlock() {
	i.checkIfEnabled()
	i.mu.RUnlock()
}

func (i *InvariantMutex) checkIfEnabled() {
	if InvariantCheckingEnabled() {
		i.check()
	}
}

// Create a lock which, when EnableInvariantChecking has been called, will call
// the supplied function at moments when invariants protected by the lock
// should hold (e.g. just after acquiring the lock). The function should crash
// if an invariant is violated. It should not have side effects, as there are
// no guarantees that it will run.
//
// The invariants must hold at the time that NewInvariantMutex is called.
func NewInvariantMutex(check func()) (mu InvariantMutex) {
	if check == nil {
		panic("check must be non-nil.")
	}

	mu.check = check
	mu.checkIfEnabled()

	return
}
