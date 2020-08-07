package tpm2tools

import (
	"reflect"
	"testing"

	"github.com/google/go-tpm-tools/internal"
	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

const (
	// Maximum number of handles to keys tests can create within a simulator.
	maxHandles = 3
)

func TestHandles(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer CheckedClose(t, rwc)

	expected := make([]tpmutil.Handle, 0)
	for i := 0; i < maxHandles; i++ {
		expected = append(expected, internal.LoadRandomExternalKey(t, rwc))

		handles, err := Handles(rwc, tpm2.HandleTypeTransient)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(handles, expected) {
			t.Errorf("Handles mismatch got: %v; want: %v", handles, expected)
		}
	}

	// Don't leak our handles
	for _, handle := range expected {
		if err := tpm2.FlushContext(rwc, handle); err != nil {
			t.Error(err)
		}
	}
}
