// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Synopsis:
//
//	page [file]
//
// Description:
// page prints a page at a time to stdout from either stdin or a named file.
// It stops every x rows, where x is the number of rows determined from gtty.
// Single character commands tell it what to do next. Currently the only ones
// are return and q.
//
// Options:
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/u-root/u-root/pkg/termios"
)

func page(t *termios.TTYIO, r io.Reader, w io.Writer) error {
	rows := int64(24)
	if w, err := t.GetWinSize(); err != nil {
		log.Printf("Could not get win size: %v; continuing assuming %d rows", err, rows)
	} else {
		rows = int64(w.Row)
	}

	l := int64(1)
	scanner := bufio.NewScanner(r)
	for {
		cur := l
		for {
			if !scanner.Scan() {
				return scanner.Err()
			}
			line := scanner.Text()
			if _, err := fmt.Fprintf(w, "%s\r\n", string(line)); err != nil {
				return err
			}
			cur++
			if cur > l+rows {
				break
			}
		}
		if cur == l {
			break
		}
		l = cur
		if _, err := fmt.Fprintf(t, ":"); err != nil {
			return err
		}
		var cmd [1]byte
		if _, err := t.Read(cmd[:]); err != nil {
			return err
		}
		switch cmd[0] {
		default:
			fmt.Printf("%q:unknown\n", cmd[0])
		case '\n', ' ':
		case 'q':
			return nil
		}
		fmt.Fprintf(w, "\r")
	}

	return nil
}

func main() {
	t, err := termios.New()
	if err != nil {
		log.Fatal(err)
	}
	in := os.Stdin

	switch len(os.Args) {
	case 1:
	case 2:
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		in = f
	default:
		log.Fatal("Usage: page [file]")
	}
	c, err := t.Raw()
	if err != nil {
		log.Fatal(err)
	}
	restore := func() {
		if err := t.Set(c); err != nil {
			log.Printf("Restoring modes failed; sorry (%v)", err)
		}
	}
	defer restore()

	cc := make(chan os.Signal, 1)
	signal.Notify(cc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-cc
		restore()
		os.Exit(1)
	}()

	if err := page(t, in, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
