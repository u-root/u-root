// +build !windows

package internal

import (
	"flag"
	"io"

	"github.com/google/go-tpm/tpm2"
)

// As this package is only included in tests, this flag will not conflict with
// the --tpm-path flag in gotpm/cmd
var tpmPath = flag.String("tpm-path", "", "Path to Linux TPM character device (i.e. /dev/tpm0 or /dev/tpmrm0). Empty value (default) will run tests against the simulator.")

func useRealTPM() bool {
	return *tpmPath != ""
}

func getRealTPM() (io.ReadWriteCloser, error) {
	return tpm2.OpenTPM(*tpmPath)
}
