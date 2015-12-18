package main

import (
	"bytes"
	"testing"
)

func TestEcho(t *testing.T) {

	type test struct {
		s         string
		nonewline bool
	}

	tests := []test{{s: "simple test", nonewline: false}, {s: "another one", nonewline: true}}
	bufs := make([]bytes.Buffer, len(tests))

	for i, v := range tests {
		if err := echo(&bufs[i], v.s); err != nil {
			t.Errorf("%s", err)
		}
		if !*nonewline {
			v.s += "\n"
		}
		if string(bufs[i].Bytes()) != v.s {
			t.Fatalf("Want %v, got %v", v.s, string(bufs[i].Bytes()))
		}
	}

}
