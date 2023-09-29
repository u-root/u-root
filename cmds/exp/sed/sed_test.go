// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// By Ahmed Kamal <email.ahmedkamal@googlemail.com>

package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestTmpWriter(t *testing.T) {
	tmpDir := t.TempDir()
	f1, err := os.CreateTemp(tmpDir, "f1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = f1.WriteString("hix\nnix\n")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	testcases := []struct {
		filename string
		content  string
		err      error
	}{
		{
			filename: "/tmp/tw",
			content:  "foo\nbar",
			err:      nil,
		},
	}

	for idx, tc := range testcases {
		test := tc
		t.Run(fmt.Sprintf("case_%d", idx), func(t *testing.T) {
			tw, err := newTmpWriter(test.filename)
			if err != nil {
				t.Errorf("failed to create tempWriter: %v", err)
			}
			fmt.Fprint(tw, test.content)
			tw.Close()
			fh, _ := os.Open(tc.filename)
			fc, _ := io.ReadAll(fh)
			fcontent := string(fc)
			if fcontent != tc.content {
				t.Errorf("got %#v, want %#v", fcontent, tc.content)
			}
		})
	}
}

func TestTransform(t *testing.T) {

	testcases := []struct {
		re     *regexp.Regexp
		to     string
		input  string
		output string
		global bool
	}{
		{
			re:     regexp.MustCompile(`\d+`),
			to:     "1980",
			input:  "The year \nis 2023\n", // Test fails if line does not end with \n
			output: "The year \nis 1980\n",
		},
		{
			re:     regexp.MustCompile(`\d+`),
			to:     "1980",
			input:  "The year \nis no matches found\n",
			output: "The year \nis no matches found\n",
		},
		{
			re:     regexp.MustCompile(`\d+`),
			to:     "1980",
			input:  "The year \nis 2023 or 2030 not sure\n", // Test fails if line does not end with \n
			output: "The year \nis 1980 or 1980 not sure\n",
			global: true,
		},
		{
			re:     regexp.MustCompile(`\d+`),
			to:     "1980",
			input:  "The year \nis 2023 or 2030 not sure\n", // Test fails if line does not end with \n
			output: "The year \nis 1980 or 2030 not sure\n",
			global: false,
		},
	}

	for idx, tc := range testcases {
		tc := tc
		t.Run(fmt.Sprintf("case_%d", idx), func(t *testing.T) {
			replacer := &transform{tc.re, tc.to, tc.global}
			fhin := strings.NewReader(tc.input)
			fhout := replacer.run(fhin)
			got, _ := io.ReadAll(fhout)
			if string(got) != tc.output {
				t.Errorf("got %#v, want %#v", string(got), tc.output)
			}
		})
	}
}
