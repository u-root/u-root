// Copyright (C) 2017 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package tftp // import "pack.ag/tftp"

import (
	"log"
	"os"
)

var (
	debug bool
	trace bool
)

func init() {
	if os.Getenv("TFTP_DEBUG") != "" {
		debug = true
	}
	if os.Getenv("TFTP_TRACE") != "" {
		debug = true
		trace = true
	}
}

type logger struct {
	log *log.Logger
	d   bool
	t   bool
}

func newLogger(name string) *logger {
	prefix := "tftp|"
	if name != "" {
		prefix += name + "|"
	}
	return &logger{log: log.New(os.Stderr, prefix, log.Lshortfile), d: debug, t: trace}
}

func (l *logger) debug(f string, args ...interface{}) {
	if l.d {
		l.log.Printf("[DEBUG] "+f, args...)
	}
}

func (l *logger) trace(f string, args ...interface{}) {
	if l.t {
		l.log.Printf("[TRACE] "+f, args...)
	}
}

func (l *logger) err(f string, args ...interface{}) {
	l.log.Printf("[ERROR] "+f, args...)
}
