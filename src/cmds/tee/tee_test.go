package main

import (
	"fmt"
	"io/ioutil"
	"os"
	//"os/signal"
	"strings"
	"testing"
)

func Test_copyinput(t *testing.T) {

	var expected, found string

	oflags := os.O_WRONLY | os.O_CREATE

	var buf [8192]byte
	input := "writing with tee"
	reader := strings.NewReader(input)
	n, _ := reader.Read(buf[:])
	fi, _ := ioutil.TempFile("", "testfile.txt")
	files := []*os.File{fi}

	fmt.Printf("Creating a file\n")
	copyinput(files, buf, n)
	expected = input
	bytes, _ := ioutil.ReadFile(fi.Name())
	found = fmt.Sprintf("%s", bytes)
	if expected != found {
		t.Fail()
	}

	fmt.Printf("Appending to file\n")
	oflags |= os.O_APPEND
	copyinput(files, buf, n)
	s := []string{input, input}
	expected = strings.Join(s, "")
	bytes, _ = ioutil.ReadFile(fi.Name())
	found = fmt.Sprintf("%s", bytes)
	if expected != found {
		t.Fail()
	}

	//is this really something that needs testing?
	/*fmt.Printf("Ignoring SIGINT\n")
		signal.Ignore(os.Interrupt)
	        copyinput(files)*/

	os.Remove(fi.Name())
}
