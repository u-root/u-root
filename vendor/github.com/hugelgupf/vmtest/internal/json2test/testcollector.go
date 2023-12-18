// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json2test

import (
	"fmt"
	"log"
	"sync"
)

// TestState are the possible Go test states.
type TestState string

// These states are taken from Go.
const (
	StateSkip    TestState = "skip"
	StateFail    TestState = "fail"
	StatePass    TestState = "pass"
	StatePaused  TestState = "paused"
	StateRunning TestState = "running"
)

var actionToState = map[Action]TestState{
	Skip:     StateSkip,
	Fail:     StateFail,
	Pass:     StatePass,
	Pause:    StatePaused,
	Run:      StateRunning,
	Continue: StateRunning,
}

// TestKind are the Go test types.
type TestKind int

// The two Go test types, test and benchmark.
const (
	KindTest TestKind = iota
	KindBenchmark
)

// TestResult is an individual tests' outcome.
type TestResult struct {
	Kind       TestKind
	State      TestState
	FullOutput string
}

// TestCollector holds Go test result information.
type TestCollector struct {
	mu sync.Mutex

	// Package collects all output for a particular package.
	Packages map[string]string

	// Tests are indexed by fully-qualified packageName.TestName strings.
	Tests map[string]*TestResult
}

// NewTestCollector returns a Handler that collects test results.
func NewTestCollector() *TestCollector {
	return &TestCollector{
		Packages: make(map[string]string),
		Tests:    make(map[string]*TestResult),
	}
}

// Handle implements Handler.
func (tc *TestCollector) Handle(e TestEvent) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if _, ok := tc.Packages[e.Package]; !ok {
		tc.Packages[e.Package] = ""
	}
	tc.Packages[e.Package] += e.Output

	if len(e.Test) == 0 {
		return
	}

	testName := fmt.Sprintf("%s.%s", e.Package, e.Test)
	t, ok := tc.Tests[testName]
	if !ok {
		t = &TestResult{
			Kind: KindTest,
		}
		tc.Tests[testName] = t
	}

	switch e.Action {
	case Benchmark:
		t.Kind = KindBenchmark
	case Output:
	default:
		s, ok := actionToState[e.Action]
		if !ok {
			log.Printf("Unknown action %q in event %v", e.Action, e)
		}
		t.State = s
	}
	t.FullOutput += e.Output
}
