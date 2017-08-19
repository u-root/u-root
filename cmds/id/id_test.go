package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var testpath = "."

type test struct {
	opt []string
	out string
}

// Run the command, with the optional args, and return a string
// for stdout, stderr, and an error.
func run(c *exec.Cmd) (string, string, error) {
	var o, e bytes.Buffer
	c.Stdout, c.Stderr = &o, &e
	err := c.Run()
	return o.String(), e.String(), err
}

func testInvocation(t *testing.T) {

	var tests = []test{
		{opt: []string{"-n"}, out: "id: cannot print only names in default format\n"},
		{opt: []string{"-G", "-g"}, out: "id: cannot print \"only\" of more than one choice\n"},
		{opt: []string{"-G", "-u"}, out: "id: cannot print \"only\" of more than one choice\n"},
		{opt: []string{"-g", "-u"}, out: "id: cannot print \"only\" of more than one choice\n"},
		{opt: []string{"-g", "-u", "-G"}, out: "id: cannot print \"only\" of more than one choice\n"},
	}

	for _, test := range tests {
		c := exec.Command(testpath, test.opt...)
		_, e, err := run(c)
		// Ignore the date and time because we're using Log.Fatalf
		if e != test.out {
			t.Errorf("id for '%v' failed: got '%s', want '%s'", test.opt, e, test.out)
		} else if err != nil {
			t.Errorf("id for '%v' failed to run", test.opt)
		}
	}
}

func TestMain(m *testing.M) {
	tempDir, err := ioutil.TempDir("", "TestIdSimple")
	if err != nil {
		fmt.Printf("cannot create temporary directory: %v", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	testpath = filepath.Join(tempDir, "testid.exe")
	out, err := exec.Command("go", "build", "-o", testpath, ".").CombinedOutput()
	if err != nil {
		fmt.Printf("go build -o %v cmds/id: %v\n%s", testpath, err, string(out))
		os.Exit(1)
	}
	os.Exit(m.Run())
}
