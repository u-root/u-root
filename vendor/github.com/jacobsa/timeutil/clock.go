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

package timeutil

import (
	"sync"
	"time"
)

type Clock interface {
	Now() time.Time
}

////////////////////////////////////////////////////////////////////////
// Real clock
////////////////////////////////////////////////////////////////////////

type realClock struct{}

func (c realClock) Now() time.Time {
	return time.Now()
}

// Return a clock that follows the real time, according to the system.
func RealClock() Clock {
	return realClock{}
}

////////////////////////////////////////////////////////////////////////
// Simulated clock
////////////////////////////////////////////////////////////////////////

// A clock that allows for manipulation of the time, which does not change
// unless AdvanceTime is called. The zero value is a clock initialized to the
// zero time.
type SimulatedClock struct {
	Clock

	mu sync.RWMutex
	t  time.Time // GUARDED_BY(mu)
}

func (sc *SimulatedClock) Now() time.Time {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.t
}

// Set the current time according to the clock.
func (sc *SimulatedClock) SetTime(t time.Time) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.t = t
}

// Advance the current time according to the clock by the supplied duration.
func (sc *SimulatedClock) AdvanceTime(d time.Duration) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.t = sc.t.Add(d)
}
