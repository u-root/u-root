// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestCksum(t *testing.T) {
	var testMatrix = []struct {
		data  []byte
		cksum string
	}{
		{[]byte("abcdef\n"), "5ab557c937e38f15291c04b7e99544ad"},
		{[]byte("pqra\n"), "721d6b135656aa83baca6ebdbd2f6c86"},
	}

	for _, testData := range testMatrix {
		if testData.cksum != calculateMd5Sum("", testData.data) {
			t.Errorf("md5sum verification failed. (Expected: %s, Received: %s)", testData.cksum, calculateMd5Sum("", testData.data))
		}
	}

}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
