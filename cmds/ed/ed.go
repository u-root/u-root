// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Ed is a simple line-oriented editor
//
// Synopsis:
//     dd
//
// Description:
//
// Options:
package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"regexp"

	"github.com/u-root/u-root/pkg/log"
)

type editorArg func(Editor) error

var (
	f           Editor = &file{}
	num                = regexp.MustCompile("^[0-9][0-9]*")
	startsearch        = regexp.MustCompile("^/[^/]/")
	endsearch          = regexp.MustCompile("^,/[^/]/")
	editors            = map[string]func(...editorArg) (Editor, error){
		"text": NewTextEditor,
		"bin":  NewBinEditor,
	}
	fileType = flag.String("t", "text", "type of file")
)

func readerio(r io.Reader) editorArg {
	return func(f Editor) error {
		_, err := f.Read(r, 0, 0)
		return err
	}
}

func readFile(n string) editorArg {
	return func(f Editor) error {
		r, err := os.Open(n)
		if err != nil {
			return err
		}

		_, err = f.Read(r, 0, 0)
		return err
	}
}

func main() {
	var (
		args []editorArg
		err  error
	)

	flag.Parse()

	e, ok := editors[*fileType]
	if !ok {
		flag.Usage()
	}

	if len(flag.Args()) == 1 {
		args = append(args, readFile(flag.Args()[0]))
	}

	ed, err := e(args...)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Now just eat the lines, and turn them into commands.
	// The format is a regular language.
	// [start][,end]command[rest of line]
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		if err := DoCommand(ed, s.Text()); err != nil {
			log.Printf(err.Error())
		}
	}
}
