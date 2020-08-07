// These tests are skipped when the race detector (-race) is on
// +build !race

package pgzip

import (
	"bytes"
	"io/ioutil"
	"runtime"
	"runtime/debug"
	"testing"
)

// Test that the sync.Pools are working properly and we are not leaking buffers
// Disabled with -race, because the race detector allocates a lot of memory
func TestAllocations(t *testing.T) {

	w := NewWriter(ioutil.Discard)
	w.SetConcurrency(100000, 10)
	data := bytes.Repeat([]byte("TEST"), 41234) // varying block splits

	// Prime the pool to do initial allocs
	for i := 0; i < 10; i++ {
		_, _ = w.Write(data)
	}
	_ = w.Flush()

	allocBytes := allocBytesPerRun(1000, func() {
		_, _ = w.Write(data)
	})
	t.Logf("Allocated %.0f bytes per Write on average", allocBytes)

	// Locally it still allocates 660 bytes, which can probably be further reduced,
	// but it's better than the 175846 bytes before the pool release fix this tests.
	// TODO: Further reduce allocations
	if allocBytes > 10240 {
		t.Errorf("Write allocated too much memory per run (%.0f bytes), Pool used incorrectly?", allocBytes)
	}
}

// allocBytesPerRun returns the average total size of allocations during calls to f.
// The return value is in bytes.
//
// To compute the number of allocations, the function will first be run once as
// a warm-up. The average total size of allocations over the specified number of
// runs will then be measured and returned.
//
// AllocBytesPerRun sets GOMAXPROCS to 1 during its measurement and will restore
// it before returning.
//
// This function is based on testing.AllocsPerRun, which counts the number of
// allocations instead of the total size of them in bytes.
func allocBytesPerRun(runs int, f func()) (avg float64) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))
	// Disable garbage collector, because it could clear our pools during the run
	oldGCPercent := debug.SetGCPercent(-1)
	defer debug.SetGCPercent(oldGCPercent)

	// Warm up the function
	f()

	// Measure the starting statistics
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)
	oldTotal := memstats.TotalAlloc

	// Run the function the specified number of times
	for i := 0; i < runs; i++ {
		f()
	}

	// Read the final statistics
	runtime.ReadMemStats(&memstats)
	allocs := memstats.TotalAlloc - oldTotal

	// Average the mallocs over the runs (not counting the warm-up).
	// We are forced to return a float64 because the API is silly, but do
	// the division as integers so we can ask if AllocsPerRun()==1
	// instead of AllocsPerRun()<2.
	return float64(allocs / uint64(runs))
}
