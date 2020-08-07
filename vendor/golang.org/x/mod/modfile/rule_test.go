// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modfile

import (
	"bytes"
	"fmt"
	"testing"

	"golang.org/x/mod/module"
)

var addRequireTests = []struct {
	in   string
	path string
	vers string
	out  string
}{
	{
		`
		module m
		require x.y/z v1.2.3
		`,
		"x.y/z", "v1.5.6",
		`
		module m
		require x.y/z v1.5.6
		`,
	},
	{
		`
		module m
		require x.y/z v1.2.3
		`,
		"x.y/w", "v1.5.6",
		`
		module m
		require (
			x.y/z v1.2.3
			x.y/w v1.5.6
		)
		`,
	},
	{
		`
		module m
		require x.y/z v1.2.3
		require x.y/q/v2 v2.3.4
		`,
		"x.y/w", "v1.5.6",
		`
		module m
		require x.y/z v1.2.3
		require (
			x.y/q/v2 v2.3.4
			x.y/w v1.5.6
		)
		`,
	},
}

var setRequireTests = []struct {
	in   string
	mods []struct {
		path     string
		vers     string
		indirect bool
	}
	out string
}{
	{
		`module m
		require (
			x.y/b v1.2.3

			x.y/a v1.2.3
			x.y/d v1.2.3
		)
		`,
		[]struct {
			path     string
			vers     string
			indirect bool
		}{
			{"x.y/a", "v1.2.3", false},
			{"x.y/b", "v1.2.3", false},
			{"x.y/c", "v1.2.3", false},
		},
		`module m
		require (
			x.y/a v1.2.3
			x.y/b v1.2.3
			x.y/c v1.2.3
		)
		`,
	},
	{
		`module m
		require (
			x.y/a v1.2.3
			x.y/b v1.2.3 //
			x.y/c v1.2.3 //c
			x.y/d v1.2.3 //   c
			x.y/e v1.2.3 // indirect
			x.y/f v1.2.3 //indirect
			x.y/g v1.2.3 //	indirect
		)
		`,
		[]struct {
			path     string
			vers     string
			indirect bool
		}{
			{"x.y/a", "v1.2.3", true},
			{"x.y/b", "v1.2.3", true},
			{"x.y/c", "v1.2.3", true},
			{"x.y/d", "v1.2.3", true},
			{"x.y/e", "v1.2.3", true},
			{"x.y/f", "v1.2.3", true},
			{"x.y/g", "v1.2.3", true},
		},
		`module m
		require (
			x.y/a v1.2.3 // indirect
			x.y/b v1.2.3 // indirect
			x.y/c v1.2.3 // indirect; c
			x.y/d v1.2.3 // indirect; c
			x.y/e v1.2.3 // indirect
			x.y/f v1.2.3 //indirect
			x.y/g v1.2.3 //	indirect
		)
		`,
	},
}

var addGoTests = []struct {
	in      string
	version string
	out     string
}{
	{`module m
		`,
		`1.14`,
		`module m
		go 1.14
		`,
	},
	{`module m
		require x.y/a v1.2.3
		`,
		`1.14`,
		`module m
		go 1.14
		require x.y/a v1.2.3
		`,
	},
	{
		`require x.y/a v1.2.3
		module example.com/inverted
		`,
		`1.14`,
		`require x.y/a v1.2.3
		module example.com/inverted
		go 1.14
		`,
	},
	{
		`require x.y/a v1.2.3
		`,
		`1.14`,
		`require x.y/a v1.2.3
		go 1.14
		`,
	},
}

func TestAddRequire(t *testing.T) {
	for i, tt := range addRequireTests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			f, err := Parse("in", []byte(tt.in), nil)
			if err != nil {
				t.Fatal(err)
			}
			g, err := Parse("out", []byte(tt.out), nil)
			if err != nil {
				t.Fatal(err)
			}
			golden, err := g.Format()
			if err != nil {
				t.Fatal(err)
			}

			if err := f.AddRequire(tt.path, tt.vers); err != nil {
				t.Fatal(err)
			}
			out, err := f.Format()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(out, golden) {
				t.Errorf("have:\n%s\nwant:\n%s", out, golden)
			}
		})
	}
}

func TestSetRequire(t *testing.T) {
	for i, tt := range setRequireTests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			f, err := Parse("in", []byte(tt.in), nil)
			if err != nil {
				t.Fatal(err)
			}
			g, err := Parse("out", []byte(tt.out), nil)
			if err != nil {
				t.Fatal(err)
			}
			golden, err := g.Format()
			if err != nil {
				t.Fatal(err)
			}
			var mods []*Require
			for _, mod := range tt.mods {
				mods = append(mods, &Require{
					Mod: module.Version{
						Path:    mod.path,
						Version: mod.vers,
					},
					Indirect: mod.indirect,
				})
			}

			f.SetRequire(mods)
			out, err := f.Format()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(out, golden) {
				t.Errorf("have:\n%s\nwant:\n%s", out, golden)
			}

			f.Cleanup()
			if len(f.Require) != len(mods) {
				t.Errorf("after Cleanup, len(Require) = %v; want %v", len(f.Require), len(mods))
			}
		})
	}
}

func TestAddGo(t *testing.T) {
	for i, tt := range addGoTests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			f, err := Parse("in", []byte(tt.in), nil)
			if err != nil {
				t.Fatal(err)
			}
			g, err := Parse("out", []byte(tt.out), nil)
			if err != nil {
				t.Fatal(err)
			}
			golden, err := g.Format()
			if err != nil {
				t.Fatal(err)
			}

			if err := f.AddGoStmt(tt.version); err != nil {
				t.Fatal(err)
			}
			out, err := f.Format()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(out, golden) {
				t.Errorf("have:\n%s\nwant:\n%s", out, golden)
			}
		})
	}
}
