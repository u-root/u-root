// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tss

import (
	"crypto/sha1"
	"fmt"
	"io"

	tpm2 "github.com/google/go-tpm/legacy/tpm2"
	tpm1 "github.com/google/go-tpm/tpm"
	tpmutil "github.com/google/go-tpm/tpmutil"
)

func nvRead12(rwc io.ReadWriteCloser, index, offset, length uint32, auth string) ([]byte, error) {
	var ownAuth [20]byte // owner well known
	if auth != "" {
		ownAuth = sha1.Sum([]byte(auth))
	}

	// Get TPMInfo
	indexData, err := tpm1.GetNVIndex(rwc, index)
	if err != nil {
		return nil, err
	}
	if indexData == nil {
		return nil, fmt.Errorf("index not found")
	}

	// Check if authData is needed
	// AuthRead 0x00200000 | OwnerRead 0x00100000
	needAuthData := 1 >> (indexData.Permission.Attributes & (nvPerAuthRead | nvPerOwnerRead))
	authread := 1 >> (indexData.Permission.Attributes & nvPerAuthRead)

	if needAuthData == 0 {
		if authread != 0 {
			return tpm1.NVReadValue(rwc, index, offset, length, ownAuth[:])
		}
		return tpm1.NVReadValueAuth(rwc, index, offset, length, ownAuth[:])
	}
	return tpm1.NVReadValue(rwc, index, offset, length, nil)
}

func nvRead20(rwc io.ReadWriteCloser, index, authHandle tpmutil.Handle, password string, blocksize int) ([]byte, error) {
	return tpm2.NVReadEx(rwc, index, authHandle, password, blocksize)
}
