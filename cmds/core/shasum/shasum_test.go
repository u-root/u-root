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
		data      []byte
		cksum     string
		algorithm int
	}{
		{[]byte("abcdef\n"), "bdc37c074ec4ee6050d68bc133c6b912f36474df", 1},
		{[]byte("pqra\n"), "e8ed2d487f1dc32152c8590f39c20b7703f9e159", 1},
		{[]byte("abcdef\n"), "ae0666f161fed1a5dde998bbd0e140550d2da0db27db1d0e31e370f2bd366a57", 256},
		{[]byte("pqra\n"), "db296dd0bcb796df9b327f44104029da142c8fff313a25bd1ac7c3b7562caea9", 256},
	}

	for _, testData := range testMatrix {
		if testData.cksum != shaPrinter(testData.algorithm, testData.data) {
			t.Errorf("shasum verification failed.(Expected:%s, Received:%s)", testData.cksum, shaPrinter(testData.algorithm, testData.data))
		}
	}

}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
