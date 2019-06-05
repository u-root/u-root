// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type test struct {
	name  string
	flags []string
	in    string
	out   string
}

var stringsTests = []test{
	{
		"empty",
		[]string{}, "", "",
	},
	{
		"sequences are too short",
		[]string{}, "\n\na\nab\n\nabc\nabc\xff\n01\n", "",
	},
	{
		"entire string is too short",
		[]string{}, "abc", "",
	},
	{
		"entire string just fits perfectly",
		[]string{}, "abcd", "abcd\n",
	},
	{
		"entire string is printable",
		[]string{}, "abcdefghijklmnopqrstuvwxyz", "abcdefghijklmnopqrstuvwxyz\n",
	},
	{
		"terminating newline",
		[]string{}, "abcdefghijklmnopqrstuvwxyz\n", "abcdefghijklmnopqrstuvwxyz\n",
	},
	{
		"mix of printable and non-printable sequences",
		[]string{}, "\n\na123456\nab\n\nabc\nabcde\xff\n01\n", "a123456\nabcde\n",
	},
	{
		"spaces are printable",
		[]string{}, " abcdefghijklm nopqrstuvwxyz ", " abcdefghijklm nopqrstuvwxyz \n",
	},
	{
		"shorter value of n",
		[]string{"--n", "1"}, "\n\na\nab\n\nabc\nabc\xff\n01\n", "a\nab\nabc\nabc\n01\n",
	},
	{
		"larger value of n",
		[]string{"--n", "6"}, "\n\na123456\nab\n\nabc\nabcde\xff\n01\n", "a123456\n",
	},
}

// strings < in > out
func TestSortWithPipes(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "strings")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Table-driven testing
	for _, tt := range stringsTests {
		t.Run(fmt.Sprintf("%v", tt.name), func(t *testing.T) {
			cmd := testutil.Command(t, tt.flags...)
			cmd.Stdin = bytes.NewReader([]byte(tt.in))
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("strings(%#v): %v", tt.in, err)
			}
			if string(out) != tt.out {
				t.Errorf("strings(%#v) = %#v; want %#v", tt.in,
					string(out), tt.out)
			}
		})
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
