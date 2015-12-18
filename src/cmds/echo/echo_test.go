package main

import (
	"bytes"
	"testing"
)

func TestEcho(t *testing.T) {

	var tests = []string{"Simple \ttest\n"}
	bufs := make([]bytes.Buffer, len(tests))

	if err := echo(&bufs[0], "Simple \ttest"); err != nil {
		t.Errorf("%s", err)
	}

	for i, v := range tests {
		if string(bufs[i].Bytes()) != v {
			t.Fatalf("Want %v, got %v", v, string(bufs[i].Bytes()))
		}
	}

}
