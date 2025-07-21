// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ulog exposes logging via a Go interface.
//
// ulog has three implementations of the Logger interface: a Go standard
// library "log" package Logger, a kernel syslog (dmesg) Logger, and a test
// Logger that logs via a test's testing.TB.Logf.
// To use the test logger import "ulog/ulogtest".
package ulog

import (
	"log"
	"os"
)

// Logger is a log receptacle.
//
// It puts your information somewhere for safekeeping.
type Logger interface {
	Printf(format string, v ...any)
}

// Log is a Logger that prints to stderr, like the default log package.
var Log Logger = log.New(os.Stderr, "", log.LstdFlags)

type emptyLogger struct{}

func (emptyLogger) Printf(format string, v ...any) {}

// Null is a logger that prints nothing.
var Null Logger = emptyLogger{}
