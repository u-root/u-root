// +build !windows

// Binary tpm2-nvread reads data from NVRAM at a specified index. The data is
// printed out hex-encoded.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

var (
	tpmPath = flag.String("tpm-path", "/dev/tpm0", "Path to the TPM device (character device or a Unix socket)")
	index   = flag.Uint("index", 0, "NVRAM index of read")
)

func main() {
	flag.Parse()

	if *index == 0 {
		fmt.Fprintln(os.Stderr, "--index must be set")
		os.Exit(1)
	}

	val, err := nvRead(*tpmPath, uint32(*index))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reading from index 0x%x: %v\n", *index, err)
		os.Exit(1)
	}
	fmt.Printf("NVRAM value at index 0x%x (hex encoded):\n%x\n", *index, val)
}

func nvRead(path string, index uint32) ([]byte, error) {
	rwc, err := tpm2.OpenTPM(path)
	if err != nil {
		return nil, fmt.Errorf("can't open TPM at %q: %v", path, err)
	}
	defer rwc.Close()
	return tpm2.NVRead(rwc, tpmutil.Handle(index))
}
