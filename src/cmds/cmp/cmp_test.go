package main

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func setup() {

	data1 := []byte("hello\nthis is a test\nfile")
	data2 := []byte("hello\nthiz is a text\nFile")
	ioutil.WriteFile("/tmp/dat1", data1, 0644)
	ioutil.WriteFile("/tmp/dat2", data2, 0644)

}

func Test_emit(t *testing.T) {

	c1 := make(chan byte, 8192)
	c2 := make(chan byte, 8192)
	c3 := make(chan byte, 8192)

	setup()

	f, _ := os.Open("/tmp/dat1")
	defer os.Remove("/tmp/dat1")
	err := emit(f, c1, 0)
	if err != io.EOF {
		t.Errorf("%v\n", err)
	}

	f, _ = os.Open("/tmp/dat2")
	defer os.Remove("/tmp/dat2")
	err = emit(f, c2, 0)
	if err != io.EOF {
		t.Errorf("%v", err)
	}

	err = emit(os.Stdin, c3, 0)
	if err != io.EOF {
		t.Errorf("%v", err)
	}

}

func Test_openFile(t *testing.T) {
	setup()

	defer os.Remove("/tmp/dat1")
	defer os.Remove("/tmp/dat2")
	if _, err := openFile("-"); err != nil {
		t.Errorf("Failed to open file %s: %v", err)
	}

	if _, err := openFile("/tmp/dat1"); err != nil {
		t.Errorf("Failed to open file %s: %v", err)
	}
}
