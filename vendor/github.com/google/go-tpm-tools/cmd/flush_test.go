package cmd

import (
	"testing"

	"github.com/google/go-tpm-tools/internal"
	"github.com/google/go-tpm-tools/tpm2tools"
	"github.com/google/go-tpm/tpm2"
)

func TestFlushNothing(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer tpm2tools.CheckedClose(t, rwc)
	ExternalTPM = rwc

	RootCmd.SetArgs([]string{"flush", "all", "--quiet"})
	if err := RootCmd.Execute(); err != nil {
		t.Error(err)
	}
}

func TestFlush(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer tpm2tools.CheckedClose(t, rwc)
	ExternalTPM = rwc

	RootCmd.SetArgs([]string{"flush", "transient", "--quiet"})

	// Loads then flushes 1, 2, 3 transient handles.
	for numHandles := 1; numHandles <= 3; numHandles++ {
		for i := 0; i < numHandles; i++ {
			internal.LoadRandomExternalKey(t, rwc)
		}

		if err := RootCmd.Execute(); err != nil {
			t.Error(err)
		}

		// Ensure there are no active handles after that.
		h, err := tpm2tools.Handles(rwc, tpm2.HandleTypeTransient)
		if err != nil {
			t.Fatal(err)
		}
		if len(h) != 0 {
			t.Errorf("TPM should be empty of transient handles; got: %d; want: 0", len(h))
		}
	}
}
