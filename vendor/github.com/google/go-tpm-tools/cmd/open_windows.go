package cmd

import (
	"io"

	"github.com/google/go-tpm/tpm2"
)

// There is no need for flags on Windows, as there is no concept of a TPM path.
func openImpl() (io.ReadWriteCloser, error) {
	return tpm2.OpenTPM()
}
