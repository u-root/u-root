// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// By Ahmed Kamal <email.ahmedkamal@googlemail.com>

package main

import (
	"bytes"
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

// SedTest is a table-driven which spawns sed with a variety of options and inputs.
func TestStdinSed(t *testing.T) {
	tests := []struct {
		input  string
		output string
		err    error
		p      params
		args   []string
	}{
		// BEWARE: the IO package seems to want this to be newline terminated.
		// If you just use hix with no newline the test will fail. Yuck.
		{
			input:  "hix\n",
			output: "hix\n",
			err:    nil,
		},
		{
			input:  "hix\n",
			output: "nix\n",
			err:    nil,
			p:      params{expr: []string{"s/hix/nix/"}},
		},
		{
			input:  "hix\n",
			output: "nix\n",
			err:    fmt.Errorf("error parsing expressions"),
			p:      params{expr: []string{"s/hix"}},
		},
		{
			input:  "hix\n",
			output: "nix\n",
			err:    nil,
			p:      params{expr: []string{"s/hix/nix/"}, inplace: true},
		},
		{
			input:  "foo and foo\nbar and bar\n",
			output: "FOO and FOO\nBAR and bar\n",
			err:    nil,
			p:      params{expr: []string{"s/foo/FOO/g", "s@bar@BAR@"}},
		},
	}

	for idx, te := range tests {
		test := te
		t.Run(fmt.Sprintf("case_%d", idx), func(t *testing.T) {
			var stdout bytes.Buffer
			rc := io.NopCloser(strings.NewReader(test.input))
			cmd := command(rc, &stdout, nil, test.p, test.args)
			err := cmd.run()
			if (err != nil) != (test.err != nil) {
				t.Errorf("unexpected err %v", err)
			}
			if err == nil {
				res := stdout.String()
				if res != test.output {
					t.Errorf("got out %q, want %q", res, test.output)
				}
			}
		})
	}
}
