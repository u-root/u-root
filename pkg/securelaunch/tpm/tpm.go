// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tpm reads and extends pcrs with measurements.
package tpm

import (
	"fmt"
	"io"

	"github.com/google/go-tpm/tpm2"
)

/*
 * GetHandle returns a tpm device handle from go-tpm/tpm2
 * that can be used for storing hashes.
 */
func GetHandle() (io.ReadWriteCloser, error) {
	tpm2, err := tpm2.OpenTPM("/dev/tpm0")
	if err != nil {
		return nil, fmt.Errorf("couldn't talk to TPM Device: err=%v", err)
	}

	return tpm2, nil
}
