// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"
)

func compare(got, want []*Command) error {
	if len(got) != len(want) {
		return fmt.Errorf("got %d commands, want %d commands", len(got), len(want))
	}
	for i := range got {
		g := got[i]
		w := want[i]
		if len(g.Args) != len(w.Args) {
			return fmt.Errorf("%q: Got %d Args, want %d Args", g, len(g.Args), len(w.Args))
		}
		if !reflect.DeepEqual(g.Args, w.Args) {
			return fmt.Errorf("%q: got %q commands, want %q commands", g, g.Args, w.Args)
		}
		if g.Link != w.Link {
			return fmt.Errorf("%q: Link is %q, want %q", g, g.Link, w.Link)
		}
		if g.BG != w.BG {
			return fmt.Errorf("%q: BG is %v, want %v", g, g.BG, w.BG)
		}
		// Just check simple stdin/out redirection.
		if g.fdmap[1] != w.fdmap[1] {
			return fmt.Errorf("%q: fdmap[1] is %v, want %v", g, g.fdmap[1], w.fdmap[1])
		}
		if g.fdmap[0] != w.fdmap[0] {
			return fmt.Errorf("%q: fdmap[0] is %v, want %v", g, g.fdmap[0], w.fdmap[0])
		}
	}
	return nil
}

func TestParsing(t *testing.T) {
	tests := []struct {
		name string
		line string
		c    []*Command
		typ  string
		err  error
	}{
		{name: "EOF", line: "", c: []*Command{}, typ: "EOF", err: io.EOF},
		{name: "singlequote", line: "'gadate a'", c: []*Command{
			{Args: []arg{{"gadate a", "ARG"}}, Link: "", BG: false},
		}, typ: "EOF", err: io.EOF},
		{name: "backslash", line: "\\\\\\gadate", c: []*Command{
			{Args: []arg{{"\\gadate", "ARG"}}, Link: "", BG: false},
		}, typ: "EOF", err: io.EOF},
		{name: "env", line: "adate $a", c: []*Command{
			{Args: []arg{{"adate", "ARG"}, {"a", "ARG"}}, Link: "", BG: false},
		}, typ: "EOF", err: io.EOF},
		{name: "No EOL", line: "adate", c: []*Command{
			{Args: []arg{{"adate", "ARG"}}, Link: "", BG: false},
		}, typ: "EOF", err: io.EOF},
		{name: "redir>/dev/null", line: "redir > /dev/null", c: []*Command{
			{Args: []arg{{"redir", "ARG"}}, Link: "", BG: false, fdmap: map[int]string{1: "/dev/null"}},
		}, typ: "EOF", err: io.EOF},
		{name: "redir</dev/null", line: "redir</dev/null", c: []*Command{
			{Args: []arg{{"redir", "ARG"}}, Link: "", BG: false, fdmap: map[int]string{0: "/dev/null"}},
		}, typ: "EOF", err: io.EOF},
		{name: "Single", line: "adate\n", c: []*Command{
			{Args: []arg{{"adate", "ARG"}}, Link: "", BG: false},
		}, typ: "EOL", err: nil},
		{name: "BG", line: "adate&\n", c: []*Command{
			{Args: []arg{{"adate", "ARG"}}, Link: "", BG: true},
		}, typ: "EOL", err: nil},
		{name: "BG", line: "adate&", c: []*Command{
			{Args: []arg{{"adate", "ARG"}}, Link: "", BG: true},
		}, typ: "EOF", err: nil},
		{name: "one BG one not", line: "adate&bdate\n", c: []*Command{
			{Args: []arg{{"adate", "ARG"}}, Link: "", BG: true},
			{Args: []arg{{"bdate", "ARG"}}, Link: "", BG: false},
		}, typ: "EOL", err: nil},
		{name: "two BG", line: "adate&bdate&\n", c: []*Command{
			{Args: []arg{{"adate", "ARG"}}, Link: "", BG: true},
			{Args: []arg{{"bdate", "ARG"}}, Link: "", BG: true},
		}, typ: "EOL", err: nil},
		{name: "AND", line: "adate&&bdate\n", c: []*Command{
			{Args: []arg{{"adate", "ARG"}}, Link: "&&"},
			{Args: []arg{{"bdate", "ARG"}}, Link: ""},
		}, typ: "EOL", err: nil},
		{name: "OR", line: "adate||bdate\n", c: []*Command{
			{Args: []arg{{"adate", "ARG"}}, Link: "||"},
			{Args: []arg{{"bdate", "ARG"}}, Link: ""},
		}, typ: "EOL", err: nil},
		{name: "BG_OR", line: "zdate & adate||bdate\n", c: []*Command{
			{Args: []arg{{"zdate", "ARG"}}, Link: "", BG: true},
			{Args: []arg{{"adate", "ARG"}}, Link: "||"},
			{Args: []arg{{"bdate", "ARG"}}, Link: ""},
		}, typ: "EOL", err: nil},
		{name: "BG_PIPE", line: "zdate & adate|bdate\n", c: []*Command{
			{Args: []arg{{"zdate", "ARG"}}, Link: "", BG: true},
			{Args: []arg{{"adate", "ARG"}}, Link: "|"},
			{Args: []arg{{"bdate", "ARG"}}, Link: ""},
		}, typ: "EOL", err: nil},
		{name: "BG_PIPE_AND", line: "zdate & adate|bdate&&cdate\n", c: []*Command{
			{Args: []arg{{"zdate", "ARG"}}, Link: "", BG: true},
			{Args: []arg{{"adate", "ARG"}}, Link: "|"},
			{Args: []arg{{"bdate", "ARG"}}, Link: "&&"},
			{Args: []arg{{"cdate", "ARG"}}, Link: ""},
		}, typ: "EOL", err: nil},
	}

	for _, tt := range tests {
		c, typ, err := getCommand(bufio.NewReader(bytes.NewReader([]byte(tt.line))))
		if err != nil && !errors.Is(err, tt.err) {
			t.Errorf("%s: getCommand(%q): %v is not %v", tt.name, tt.line, err, tt.err)
			continue
		}
		if typ != tt.typ {
			t.Errorf("%s: getCommand(%q): got %s want %s", tt.name, tt.line, typ, tt.typ)
			continue
		}
		// We don't test broken parsing here, just that we get some expected
		// arrays
		doArgs(c)
		if err := commands(c); err != nil {
			t.Errorf("commands: %v != nil", err)
			continue
		}
		if err := wire(c); err != nil {
			t.Errorf("wire: %v != nil", err)
			continue
		}
		for _, cmd := range c {
			t.Logf("{Args: %q, Link: %q, BG: %v, fdmap: %v},", cmd.Args, cmd.Link, cmd.BG, cmd.fdmap)
		}
		if err := compare(c, tt.c); err != nil {
			t.Errorf("%s: getCommand(%q): %v", tt.name, tt.line, err)
		}

	}
}
