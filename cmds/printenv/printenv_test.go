package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestPrintenv(t *testing.T) {
	var buf bytes.Buffer

	want := os.Environ()

	printenv(&buf)

	found := strings.Split(buf.String(), "\n")

	for i, v := range want {
		if v != found[i] {
			t.Fatalf("want %s, got %s", v, found[i])
		}
	}
}
