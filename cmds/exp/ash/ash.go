// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/u-root/u-root/pkg/complete"
	"github.com/u-root/u-root/pkg/termios"
)

var (
	debug = flag.Bool("d", false, "enable debug prints")
	test  = flag.Bool("t", false, "test mode -- do completions, don't run commands")
	v     = func(string, ...interface{}) {}
)

func verbose(f string, a ...interface{}) {
	v(f+"\r\n", a...)
}

func main() {
	tty()
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
	if err != nil {
		log.Printf("non-fatal cannot get tty: %v", err)
	}
	defer func() {
		if err := t.Set(r); err != nil {
			log.Print(err)
		}
	}()

	f := complete.NewFileCompleter("")
	p, err := complete.NewPathCompleter()
	if err != nil {
		log.Fatal(err)
	}

	bin := complete.NewMultiCompleter(complete.NewStringCompleter([]string{"exit"}), p, f)
	l := complete.NewNewerLineReader(bin, f)
	l.Prompt = "% "
	for !l.EOF {
		if err := l.ReadLine(t, t); err != nil {
			log.Printf("looperr: %v", err)
			continue
		}
		if _, err := t.Write([]byte("\r\n")); err != nil {
			log.Print(err)
		}
		if l.FullLine == "" {
			continue
		}
		if l.Exact == "" {
			f := strings.Fields(l.FullLine)
			if *test {
				log.Printf("%v", f)
				if _, err := t.Write([]byte("\r\n")); err != nil {
					log.Print(err)
				}
				continue
			}

			if err := t.Set(r); err != nil {
				log.Print(err)
			}
			cmd := exec.Command(f[0], f[1:]...)
			cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

			if err := cmd.Run(); err != nil {
				log.Print(err)
			}

			foreground()
			if _, err := t.Raw(); err != nil {
				log.Print(err)
			}
		}
	}
}
