// +build linux darwin

package tpmutil

import (
	"os"
	"testing"
)

func TestPoll(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	if _, err := w.Write([]byte("hi")); err != nil {
		t.Fatalf("error writing to pipe: %v", err)
	}
	if err := poll(r); err != nil {
		t.Errorf("error polling reader side of the pipe: %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("error closing reader side of the pipe: %v", err)
	}
}
