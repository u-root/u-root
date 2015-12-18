package main

import (
	"bytes"
	"os"
	"testing"
)

func Test_hostname(t *testing.T) {
	var buf bytes.Buffer
	var host string
	var err error

	if err = hostname(&buf); err != nil {
		t.Errorf("%v", err)
	}

	if host, err = os.Hostname(); err != nil {
		t.Errorf("%v", err)
	}

	buf2 := []byte(host)

	if bytes.Compare(buf.Bytes(), buf2) != 0 {
		t.Fatalf("want %v, got %v", buf2, buf)
	}
}
