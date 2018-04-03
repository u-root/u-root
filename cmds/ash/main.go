// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/u-root/u-root/pkg/complete"
	"github.com/u-root/u-root/pkg/termios"
)

var (
	debug = flag.Bool("d", false, "enable debug prints")
	v     = func(string, ...interface{}) {}
)

func verbose(f string, a ...interface{}) {
	v(f+"\r\n", a...)
}

func output(r io.Reader, w io.Writer) {
	for {
		var b [1]byte
		n, err := r.Read(b[:])
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Printf("output: %v", err)
			return
		}
		if n < len(b) {
			continue
		}
		var s string
		switch b[0] {
		default:
			s = string(b[:])
		case '\b', 127:
			s = "\b \b"
		case '\r', '\n':
			s = "\r\n"
		}
		if _, err := w.Write([]byte(s)); err != nil {
			log.Printf("output write: %v", err)
			return
		}
	}
}
func main() {
	flag.Parse()
	if *debug {
		v = log.Printf
		complete.Debug = verbose
	}
	t, err := termios.New()
	if err != nil {
		log.Fatal(err)
	}
	r, err := t.Raw()
	defer t.Set(r)
	cr, cw, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}
	go output(cr, t)
	for {
		p, err := complete.NewPathCompleter()
		if err != nil {
			log.Fatal(err)
		}
		c := complete.NewMultiCompleter(complete.NewStringCompleter([]string{"exit"}), p)
		l := complete.NewLineReader(c, t, cw)
		s, err := l.ReadOne()
		v("ash: Readone: %v, %v", s, err)
		if err != nil && err != complete.EOL {
			log.Print(err)
			continue
		}
		if len(s) == 0 {
			continue
		}
		if s[0] == "exit" {
			break
		}
		// s[0] is either the match or what they typed so far.
		cw.Write([]byte(" "))
		bin := s[0]
		var args []string
		for err == nil {
			c := complete.NewFileCompleter("")
			l := complete.NewLineReader(c, t, t)
			s, err := l.ReadOne()
			v("ash: l.ReadOne returns %v, %v", s, err)
			args = append(args, s...)
			v("ash: add %v", s)
			if err != nil {
				log.Print(err)
				break
			}
		}
		v("ash: Done reading args")
		cmd := exec.Command(bin, args...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, cw, cw
		if err := cmd.Run(); err != nil {
			log.Print(err)
		}
	}
}
