package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestCopyinput(t *testing.T) {
	oflags := os.O_WRONLY | os.O_CREATE

	var want = []byte("Writing with tee")

	fi, err := ioutil.TempFile("", "testfile.txt")
	if err != nil {
		t.Fatalf("Can't create temp file:%v", err)
	}
	defer os.Remove(fi.Name())
	files := []io.Writer{fi}

	t.Logf("Creating a file\n")
	copyinput(files, want)
	found, err := ioutil.ReadFile(fi.Name())
	if err != nil {
		t.Fatalf("Can't read file: %v", err)
	}
	if !reflect.DeepEqual(want, found) {
		t.Logf("Failed: want %s, got %s", want, found)
	}

	t.Logf("Appending to file\n")
	oflags |= os.O_APPEND
	copyinput(files, want)
	want = append(want, want...)
	found, err = ioutil.ReadFile(fi.Name())
	if err != nil {
		t.Fatalf("Can't read file: %v", err)
	}
	if !reflect.DeepEqual(want, found) {
		t.Logf("Failed: want %s, got %s", want, found)
	}
}

func TestIgnore(t *testing.T) {
	b := make([]byte, 1)
	var stderr bytes.Buffer

	t.Logf("Ignoring SIGINT\n")
	cmd := exec.Command("go", "build", "tee.go")

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build tee: %v", err)
	}
	defer os.Remove("tee")

	cmd = exec.Command("./tee", "-i")
	in, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to create Stdin pipe: %v", err)
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to create Stdout pipe: %v", err)
	}
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to exec tee: %v", err)
	}
	if _, err := in.Write([]byte("!")); err != nil {
		t.Fatalf("Can't write to standard input: %v", err)
	}
	if _, err := out.Read(b); err != nil {
		t.Fatalf("Can't read from stdout: %v", err)
	}
	cmd.Process.Signal(os.Interrupt)
	in.Close()
	if err := cmd.Wait(); err != nil {
		t.Fatalf("%s: %v", stderr.String(), err)
	}
}
