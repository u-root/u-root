package internal

import (
	"flag"
	"io"

	"github.com/google/go-tpm/tpm2"
)

var useTBS = flag.Bool("use-tbs", false, "Run the tests against the Windows TBS. Value of false (default) will run tests against the simulator.")

func useRealTPM() bool {
	return *useTBS
}

func getRealTPM() (io.ReadWriteCloser, error) {
	return tpm2.OpenTPM()
}
