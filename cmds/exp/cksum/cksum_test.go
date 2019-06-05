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
		cksum uint32
	}{
		{[]byte("abcdef\n"), 3512391007},
		{[]byte("pqra\n"), 1063566492},
		{[]byte("abcdef\nafdsfsfgdglfdgkdlvcmdfposvpomfvmlcvdfk;lgkd'f;k;lvcvcv\n" +
			"afdsfsfgdglfdgkdlvcmdfposvpomfvmlcvdfk;lgkd'f;k;lvcvcv\nafdsfsfg" +
			"dglfdgkdlvcmdfposvpomfvmlcvdfk;lgkd'f;k;lvcvcv\nafdsfsfgdglfdgkd" +
			"lvcmdfposvpomfvmlcvdfk;lgkd'f;k;lvcvcv\nsdddsfsfsdfsdfsdasaarwre" +
			"mazadsfssfsfsfsafsadfsfdsadfsafsafsfsafdsfsdfsfdsdf"), 689622513},
	}

	for _, testData := range testMatrix {
		if testData.cksum != calculateCksum(testData.data) {
			t.Errorf("Cksum verification failed. (Expected: %d, Received: %d)", testData.cksum, calculateCksum(testData.data))
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
