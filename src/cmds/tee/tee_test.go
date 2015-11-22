package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestCopyinput(t *testing.T) {
	oflags := os.O_WRONLY | os.O_CREATE

	var buf = []byte("Writing with tee")
	want := "Writing with tee"
	fi, err := ioutil.TempFile("", "testfile.txt")
	if err != nil {
		t.Fatalf("Can't create temp file:%v", err)
	}
	defer os.Remove(fi.Name())
	files := []io.Writer{fi}

	t.Logf("Creating a file\n")
	copyinput(files, buf)
	bytes, err := ioutil.ReadFile(fi.Name())
	if err != nil {
		t.Fatalf("Can't read file: %v", err)
	}
	found := string(bytes)
	if want != found {
		t.Logf("Failed: want %s, got %s", want, found)
	}

	t.Logf("Appending to file\n")
	oflags |= os.O_APPEND
	copyinput(files, buf)
	want += want
	bytes, err = ioutil.ReadFile(fi.Name())
	if err != nil {
		t.Fatalf("Can't read file: %v", err)
	}
	found = string(bytes)
	if want != found {
		t.Logf("Failed: want %s, got %s", want, found)
	}
}

func TestIgnore(t *testing.T) {

	t.Logf("Ignoring SIGINT\n")
	cmd := exec.Command("go", "build", "tee.go")

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build tee: %v", err)
	}
	defer os.Remove("tee")

	cmd = exec.Command("./tee", "-i")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to exec tee: %v", err)
	}
	cmd.Process.Signal(os.Interrupt)
	if err := cmd.Wait(); err != nil {
		t.Fatalf("%v", err)
	}
}
