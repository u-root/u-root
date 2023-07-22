// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tss

import (
	"flag"
	"fmt"
	"io"
	"testing"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpm2/transport"
	"github.com/google/go-tpm/tpmutil/mssim"
)

var (
	tpmPort      = flag.Int("tpm_port", -1, "TCP port of running TPM simulator's TPM command server")
	platformPort = flag.Int("platform_port", -1, "TCP port of running TPM simulator's platform command server")
)

// getSimulator connects to the running TCP TPM simulator and restarts it.
func getSimulator(t *testing.T) TPM {
	t.Helper()

	if *tpmPort == -1 {
		t.Skip("No running TPM simulator's TPM command server was provided with the --tpm_port flag. Skipping test.")
	}
	if *platformPort == -1 {
		t.Skip("No running TPM simulator's platform command server was provided with the --platform_port flag. Skipping test.")
	}

	config := mssim.Config{
		CommandAddress:  fmt.Sprintf("127.0.0.1:%v", *tpmPort),
		PlatformAddress: fmt.Sprintf("127.0.0.1:%v", *platformPort),
	}
	// A non-obvious side-effect of mssim.Open() is that it powers the simulator off and on again.
	sim, err := mssim.Open(config)
	if err != nil {
		t.Fatalf("Could not connect to TPM simulator: %v", err)
	}
	t.Cleanup(func() { sim.Close() })

	_, err = tpm2.Startup{
		StartupType: tpm2.TPMSUClear,
	}.Execute(transport.FromReadWriter(sim))
	if err != nil {
		t.Fatalf("Could not start up the TPM: %v", err)
	}

	return TPM{
		Version: TPMVersion20,
		Interf:  TPMInterfaceDirect,
		RWC:     borrowRWC(sim),
	}
}

// borrowedRWC wraps an io.ReadWriter with a no-op Close command (see borrowRWC).
type borrowedRWC struct {
	io.ReadWriter
}

// Close implements the io.Closer interface.
func (borrowedRWC) Close() error { return nil }

// borrowRWC returns an ioReadWriteCloser that wraps the normal ReadWriter
// with a no-op Close (because we don't want to close the simulator connection).
// TODO: Migrate the pkg/tss library to use either transport.TPM (recommended)
// or the legacy API using io.ReadWriter, so that only the code that opens TPM
// connections has to be able to close it (instead of all the code that talks
// to the TPM.
func borrowRWC(rw io.ReadWriter) io.ReadWriteCloser {
	return borrowedRWC{
		ReadWriter: rw,
	}
}

func TestInfo(t *testing.T) {
	tpm := getSimulator(t)

	// Since the TPM simulator is externally provided, we can't assert much
	// about the info. Just make sure it succeeds.
	_, err := tpm.Info()
	if err != nil {
		t.Errorf("tpm.Info() = %v, want nil", err)
	}
}
