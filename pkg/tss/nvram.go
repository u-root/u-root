// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tss

import (
	"crypto/sha1"
	"io"

	tpm1 "github.com/google/go-tpm/tpm"
	tpm2 "github.com/google/go-tpm/tpm2"
	tpmutil "github.com/google/go-tpm/tpmutil"
)

func nvRead12(rwc io.ReadWriteCloser, index, offset, len uint32, ownerPW string) ([]byte, error) {
	var ownAuth [20]byte

	if ownerPW != "" {
		ownAuth = sha1.Sum([]byte(ownerPW))
	}
	return tpm1.NVReadValue(rwc, index, offset, len, []byte(ownAuth[:20]))
}

func nvRead20(rwc io.ReadWriteCloser, index, authHandle tpmutil.Handle, password string, blocksize int) ([]byte, error) {
	return tpm2.NVReadEx(rwc, index, authHandle, password, blocksize)
}
