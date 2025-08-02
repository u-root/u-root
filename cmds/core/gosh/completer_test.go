// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/knz/bubbline/editline"
	"mvdan.cc/sh/v3/syntax"
)

func TestAutocompleteBubb(t *testing.T) {
	for _, tt := range []struct {
		name        string
		input       string
		completions []string
	}{
		{
			name:        "echo",
			input:       "ech",
			completions: []string{"echo"},
		},
		{
			name:        "cwd",
			input:       "./",
			completions: []string{"completer_test.go"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			val := make([][]rune, 1)
			val[0] = append(val[0], []rune(tt.input)...)

			_, completions := autocompleteBubb(val, 0, len(tt.input))

			if !completionsEqual(0, len(tt.completions), tt.completions, completions) {
				t.Errorf("want: %v, got: %v", tt.completions, completions)
			}
		})
	}
}

func completionsEqual(numCat, numEnt int, want []string, got editline.Completions) bool {
	if got.NumCategories() < numCat {
		return false
	}

	for i := range numCat {
		if got.NumEntries(i) < numEnt {
			return false
		}
	}

	for j := range numCat {
		for i := range numEnt {
			found := false
			for _, entry := range want {
				if entry == got.Entry(j, i).Title() {
					found = true
				}
			}
			if !found {
				return false
			}
		}
	}

	return true
}

func TestAutocomplete(t *testing.T) {
	parser := syntax.NewParser()

	p := t.TempDir()
	os.WriteFile(filepath.Join(p, "echo"), nil, 0o777)
	os.WriteFile(filepath.Join(p, "cat"), nil, 0o777)
	t.Setenv("PATH", p+":"+t.TempDir())

	dir2Files := t.TempDir()
	os.WriteFile(filepath.Join(dir2Files, "foo.txt"), nil, 0o777)
	os.WriteFile(filepath.Join(dir2Files, "bar.txt"), nil, 0o777)

	t.Setenv("GOSH_TEST", "1")

	for _, tt := range []struct {
		name  string
		input string
		cwd   string
		want  []string
	}{
		{
			name:  "echo",
			input: "ech",
			want:  []string{"echo"},
		},
		{
			name:  "cwd",
			input: "./c",
			want: []string{
				"./completer.go",
				"./completer_common.go",
				"./completer_liner.go",
				"./completer_nobuild.go",
				"./completer_test.go",
			},
		},
		{
			name:  "path",
			input: "../gosh/c",
			want: []string{
				"../gosh/completer.go",
				"../gosh/completer_common.go",
				"../gosh/completer_liner.go",
				"../gosh/completer_nobuild.go",
				"../gosh/completer_test.go",
			},
		},
		{
			name:  "directory",
			input: "./testdata/",
			want: []string{
				"./testdata/fuzz",
				"./testdata/setenv.sh",
			},
		},
		{
			name:  "directory",
			input: "./testdata",
			want: []string{
				"./testdata",
			},
		},
		{
			name:  "empty",
			input: "",
			want: []string{
				"cat",
				"echo",
			},
		},
		{
			name:  "spaces",
			input: "  ",
			want: []string{
				"  cat",
				"  echo",
			},
		},
		{
			name:  "args",
			input: "cat ",
			cwd:   dir2Files,
			want: []string{
				"cat ./bar.txt",
				"cat ./foo.txt",
			},
		},
		{
			name:  "envs",
			input: "FOO=bar",
		},
		{
			name:  "envs with space",
			input: "FOO=bar ",
			want: []string{
				"FOO=bar cat",
				"FOO=bar echo",
			},
		},
		{
			name:  "multistatement empty with spaces",
			input: `echo "foo";  `,
			want: []string{
				`echo "foo";  cat`,
				`echo "foo";  echo`,
			},
		},
		{
			name:  "multistatement empty",
			input: `echo "foo";`,
			want: []string{
				`echo "foo";cat`,
				`echo "foo";echo`,
			},
		},
		{
			name:  "expansion",
			cwd:   dir2Files,
			input: `cat "./bar"`,
			want: []string{
				`cat ./bar.txt`,
			},
		},
		{
			name:  "expansion env var",
			cwd:   dir2Files,
			input: `cat $GOSH_TEST`,
		},
		// We could allow this in the future by expanding to `echo
		// foobar foofoo`, which is what zsh does.
		{
			name:  "expansion multi curly",
			input: `echo foo{bar,foo}`,
		},
		{
			name:  "expansion curly",
			input: `echo foo{bar}`,
		},

		// TODO: improve redirect support.
		{
			name:  "redirect",
			input: `echo "foo" >`,
			// want: []string{`echo "foo" > {files}`}
		},
		{
			name:  "redirect",
			input: `echo "foo" > ./testdata/`,
			// want: ignore directories like fuzz
			// want:  []string{`echo "foo" > ./testdata/fuzz`},
		},
		{
			name:  "multistatement with &",
			input: `echo "foo" &`,
		},
		// TODO: parser.Stmts returns an error in this case
		/*{
			name:  "multistatement with && empty",
			input: `echo "foo" && `,
			want: []string{
				`echo "foo" && echo`,
			},
		},*/
		{
			name:  "multistatement",
			input: `echo "foo"; ech`,
			want:  []string{`echo "foo"; echo`},
		},
		{
			name:  "nocomplete",
			input: `echo "./co`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cwd != "" {
				pwd, _ := os.Getwd()
				os.Chdir(tt.cwd)
				defer os.Chdir(pwd)
			}
			got := autocompleteLiner(parser)(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("autocomplete %q = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}
