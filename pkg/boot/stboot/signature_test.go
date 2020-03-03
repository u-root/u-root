// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stboot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testHash = []byte{184, 125, 17, 19, 148, 144, 119, 239, 215, 18, 102, 79, 239, 214, 17, 133, 71, 108, 137, 131, 240, 22, 15, 165, 145, 215, 69, 91, 64, 30, 69, 244, 75, 50, 66, 44, 239, 120, 166, 25, 48, 230, 19, 190, 162, 4, 72, 244, 43, 62, 207, 124, 163, 49, 79, 215, 70, 55, 210, 29, 167, 131, 37, 43}

// not ready yet
func TestHash(t *testing.T) {
	signer := Sha512PssSigner{}
	hash, err := signer.Hash("testdata/testConfigDir/kernels/kernel1.vmlinuz")
	lenght := len(hash)
	t.Logf("hash: %v", hash)
	require.NoError(t, err)
	require.Equal(t, 64, lenght)
}

//not ready yet
func TestSign(t *testing.T) {
	signer := Sha512PssSigner{}
	sig, err := signer.Sign("testdata/keys/signing-key-1.key", testHash)
	t.Logf("signature: %v", sig)
	require.NoError(t, err)
	require.NotNil(t, sig)
}
