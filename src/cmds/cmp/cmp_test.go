package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func Test_emit(t *testing.T) {

	c1 := make(chan byte, 8192)
	c2 := make(chan byte, 8192)
	c3 := make(chan byte, 8192)
	
	data1 := []byte("hello\nthis is a test\nfile")
	data2 := []byte("hello\nthiz is a text\nFile")
	ioutil.WriteFile("/tmp/dat1", data1, 0644)
	ioutil.WriteFile("/tmp/dat2", data2, 0644)

	f, _ := os.Open("/tmp/dat1")
	err := emit(f, c1, 0)
	if err != io.EOF {
		fmt.Printf("%v\n", err)
		t.Fail()
	}

	f, _ = os.Open("/tmp/dat2")
	err = emit(f, c2, 0)	
	if err != io.EOF {
		fmt.Printf("%v", err)
		t.Fail()
	}

	err = emit(os.Stdin, c3, 0)
	if err != io.EOF {
		fmt.Printf("%v", err)
		t.Fail()
	}

	os.Remove("/tmp/dat1")
	os.Remove("/tmp/dat2")
}

