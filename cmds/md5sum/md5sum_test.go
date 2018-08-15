package main

import (
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestCksum(t *testing.T) {
	var testMatrix = []struct {
		data []byte
		cksum [16]byte
	}{
		{"abcdef\n", "5ab557c937e38f15291c04b7e99544ad"},
		{"pqra\n", "721d6b135656aa83baca6ebdbd2f6c86"},
	}
	
	for _, testData := testMatrix {
		if testMatrix.cksum != calculateMd5Sum(testMatrix.data) {
			t.Errorf("md5sum verification failed.")
		}
	}

}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}

