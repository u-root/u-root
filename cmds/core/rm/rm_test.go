// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setup(t *testing.T) string {
	d := t.TempDir()
	fbody := []byte("Go is cool!")
	for _, f := range []struct {
		name  string
		mode  os.FileMode
		isdir bool
	}{
		{
			name:  "hi",
			mode:  0o755,
			isdir: true,
		},
		{
			name: "hi/one.txt",
			mode: 0o666,
		},
		{
			name: "hi/two.txt",
			mode: 0o777,
		},
		{
			name: "go.txt",
			mode: 0o555,
		},
	} {
		var (
			err      error
			filepath = filepath.Join(d, f.name)
		)
		if f.isdir {
			err = os.Mkdir(filepath, f.mode)
		} else {
			err = os.WriteFile(filepath, fbody, f.mode)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	return d
}

func TestRm(t *testing.T) {
	for _, tt := range []struct {
		name        string
		file        string
		interactive bool
		iString     string
		verbose     bool
		recursive   bool
		force       bool
		want        string
	}{
		{
			name: "no args",
			file: "",
			want: usage,
		},
		{
			name: "rm one file",
			file: "go.txt",
			want: "",
		},
		{
			name:    "rm one file verbose",
			file:    "go.txt",
			verbose: true,
			want:    "",
		},
		{
			name: "fail to rm one file",
			file: "go",
			want: "no such file or directory",
		},
		{
			name:  "fail to rm one file forced to trigger continue",
			file:  "go",
			force: true,
			want:  "",
		},
		{
			name:        "rm one file interactive",
			file:        "go",
			interactive: true,
			iString:     "y\n",
			want:        "",
		},
		{
			name:        "rm one file interactive continue triggered",
			file:        "go",
			interactive: true,
			iString:     "\n",
			want:        "",
		},
		{
			name:      "rm dir recursivly",
			file:      "hi",
			recursive: true,
		},
		{
			name: "rm dir not recursivly",
			file: "hi",
			want: "directory not empty",
		},
	} {
		d := setup(t)

		t.Run(tt.name, func(t *testing.T) {
			var file []string

			*interactive = tt.interactive
			*verbose = tt.verbose
			*recursive = tt.recursive
			*force = tt.force

			buf := &bytes.Buffer{}
			log.SetOutput(buf)
			buf.WriteString(tt.iString)

			if tt.file != "" {
				file = []string{filepath.Join(d, tt.file)}
			}
			if err := rm(buf, file); err != nil {
				if !strings.Contains(err.Error(), tt.want) {
					t.Errorf("rm() = %q, want to contain: %q", err.Error(), tt.want)
				}
			}
		})
	}
}
