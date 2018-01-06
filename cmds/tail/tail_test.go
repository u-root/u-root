package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestTailReadBackwards(t *testing.T) {
	input, err := os.Open("./test_samples/read_backwards.txt")
	if err != nil {
		t.Error(err)
	}
	output := &bytes.Buffer{}
	err = readLastLinesBackwards(input, output, 2)
	if err != nil {
		t.Error(err)
	}
	expected := []byte("second\nthird\n")
	got := output.Bytes()
	if !bytes.Equal(got, expected) {
		t.Fatalf("Invalid result reading backwards. Got %v; want %v", got, expected)
	}
	// try reading more, which should return EOF
	buf := make([]byte, 16)
	n, err := input.Read(buf)
	if err == nil {
		t.Fatalf("Expected EOF, got more bytes instead: %v", string(buf[:n]))
	}
	if err != io.EOF {
		t.Fatalf("Expected EOF, got another error instead: %v", err)
	}
}

func TestTailReadFromBeginning(t *testing.T) {
	input, err := os.Open("./test_samples/read_from_beginning.txt")
	if err != nil {
		t.Error(err)
	}
	output := &bytes.Buffer{}
	err = readLastLinesFromBeginning(input, output, 3)
	if err != nil {
		t.Error(err)
	}
	expected := []byte("eight\nnine\nten\n")
	got := make([]byte, 4096) // anything larger than the expected result
	n, err := output.Read(got)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got[:n], expected) {
		t.Fatalf("Invalid data while reading from the beginning. Got %v; want %v", string(got[:n]), string(expected))
	}
	// try reading more, which should return EOF
	buf := make([]byte, 16)
	n, err = input.Read(buf)
	if err == nil {
		t.Fatalf("Expected EOF, got more bytes instead: %v", string(buf[:n]))
	}
	if err != io.EOF {
		t.Fatalf("Expected EOF, got another error instead: %v", err)
	}
}
