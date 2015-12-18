package main

import (
	"bytes"
	"testing"
)

func Test_echo(t *testing.T) {

	var buf bytes.Buffer
	var buf2 = []byte("Simple \ttest\n")
	
	if err:= echo("Simple \ttest", &buf); err != nil {
		t.Errorf("%s", err)
	}

	if bytes.Compare(buf.Bytes(), buf2) != 0 {
		t.Fatalf("Want %v, got %v", buf2, buf)
	}
}
