// These tests are unreliable or only pass under certain conditions.
// To run:   go test -v -count=1 -cpu=1,2,4,8,16 -tags=unreliable
// +build unreliable,!race

package pgzip

import (
	"bytes"
	"sync"
	"testing"
	"time"
)

type SlowDiscard time.Duration

func (delay SlowDiscard) Write(p []byte) (int, error) {
	time.Sleep(time.Duration(delay))
	return len(p), nil
}

// Test that the panics catch unsafe concurrent writing (a panic is better than data corruption)
// This test is UNRELIABLE and slow. The more concurrency (GOMAXPROCS), the more likely
// a race condition will be hit. If GOMAXPROCS=1, the condition is never hit.
func TestConcurrentRacePanic(t *testing.T) {
	w := NewWriter(SlowDiscard(2 * time.Millisecond))
	w.SetConcurrency(1000, 1)
	data := bytes.Repeat([]byte("T"), 100000) // varying block splits

	const n = 1000
	recovered := make(chan string, n)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				s, ok := recover().(string)
				if ok {
					recovered <- s
					t.Logf("Recovered from panic: %s", s)
				}
			}()
			// INCORRECT CONCURRENT USAGE!
			<-start
			_, _ = w.Write(data)
		}()
	}
	close(start) // give the start signal

	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()
	hasPanic := false
	select {
	case <-recovered:
		// OK, expected
		hasPanic = true
	case <-timer.C:
		t.Error("Timout")
	}
	wg.Wait()
	if !hasPanic {
		t.Error("Expected a panic, but none happened")
	}
}
