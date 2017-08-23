package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

var (
	remove          = true
	testpath        = "."
	logPrefixLength = len("2009/11/10 23:00:00 ")
)

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

// Test incorrect invocation of id
func TestInvocation(t *testing.T) {
	tempDir, idPath := testutil.CompileInTempDir(t)

	if remove {
		defer os.RemoveAll(tempDir)
	}

	var tests = []test{
		{opt: []string{"-n"}, out: "id: cannot print only names in default format\n"},
		{opt: []string{"-G", "-g"}, out: "id: cannot print \"only\" of more than one choice\n"},
		{opt: []string{"-G", "-u"}, out: "id: cannot print \"only\" of more than one choice\n"},
		{opt: []string{"-g", "-u"}, out: "id: cannot print \"only\" of more than one choice\n"},
		{opt: []string{"-g", "-u", "-G"}, out: "id: cannot print \"only\" of more than one choice\n"},
	}

	for _, test := range tests {
		c := exec.Command(idPath, test.opt...)
		_, e, _ := run(c)

		// Ignore the date and time because we're using Log.Fatalf
		if e[logPrefixLength:] != test.out {
			t.Errorf("id for '%v' failed: got '%s', want '%s'", test.opt, e, test.out)
		}
	}
}
