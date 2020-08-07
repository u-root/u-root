package internal

import (
	"io"
	"sync"
	"testing"

	"github.com/google/go-tpm-tools/simulator"
)

// Only open the TPM device once. Reopening the device causes issues on Linux.
var (
	tpm  io.ReadWriteCloser
	lock sync.Mutex
)

type noClose struct {
	io.ReadWriter
}

func (n noClose) Close() error {
	return nil
}

// GetTPM is a cross-platform testing helper function that retrives the
// appropriate TPM device from the flags passed into "go test".
func GetTPM(tb testing.TB) io.ReadWriteCloser {
	tb.Helper()
	if useRealTPM() {
		lock.Lock()
		defer lock.Unlock()
		if tpm == nil {
			var err error
			if tpm, err = getRealTPM(); err != nil {
				tb.Fatalf("Failed to open TPM: %v", err)
			}
		}
		return noClose{tpm}
	}

	simulator, err := simulator.Get()
	if err != nil {
		tb.Fatalf("Simulator initialization failed: %v", err)
	}
	return simulator
}
