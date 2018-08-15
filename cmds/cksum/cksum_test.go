package main

import (
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestCksum(t *testing.T) {
	var testMatrix = []struct {
		data string
		cksum string
	}{
		{"abcdef\n", "2315241002"},
		{"pqra\n", "2999227146"},
	}
	
	for _, testData := testMatrix {
		if testMatrix.cksum != printCksum(testMatrix.data) {
			t.Errorf("Cksum verification failed.")
		}
	}

}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}

