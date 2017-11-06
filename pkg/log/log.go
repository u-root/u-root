// Copyright 2017 the u-root Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package log implements leveled logging with two levels: nothing or
// everything.
//
// If imported, this package defines the following flags automatically:
//
//  -v               turns on log output
//  -verbose         turns on log output
//  -d               turns on log output
//  -debug           turns on log output
//  -logfile=FILE    uses FILE for log output
package log

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

const timeFormat = "15:04:05.000"

var (
	verbosity bool
	w         fileFlag = fileFlag{os.Stderr}
)

// fileFlag implements flag.Value for a file.
type fileFlag struct {
	*os.File
}

func (ff *fileFlag) Set(s string) error {
	f, err := os.OpenFile(s, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	ff.File = f
	return nil
}

func (ff fileFlag) String() string {
	if ff.File != nil {
		return ff.File.Name()
	}
	return ""
}

// RegisterFlags registers logging flags.
func RegisterFlags(f *flag.FlagSet) {
	f.BoolVar(&verbosity, "v", false, "Turns on log output.")
	f.BoolVar(&verbosity, "verbose", false, "Turns on log output.")
	f.BoolVar(&verbosity, "debug", false, "Turns on log output.")
	f.Var(&w, "logfile", "File to log to if log output is turned on.")
}

// SetVerbosity turns logging on or off.
func SetVerbosity(v bool) {
	verbosity = v
}

// Print emits the given string with time prefixed on the log writer.
func Print(s string) {
	if verbosity {
		fmt.Fprintf(w, "[%s] %s\n", time.Now().Format(timeFormat), s)
	}
}

// Printf emits the given log string with time prefixed on the log writer.
// Arguments are handled as in fmt.Printf.
func Printf(format string, args ...interface{}) {
	if verbosity {
		Print(fmt.Sprintf(format, args...))
	}
}

// PrintObj emits a json-encoded obj on the log writer.
func PrintObj(obj interface{}) {
	if verbosity {
		b, _ := json.Marshal(obj)
		fmt.Fprintf(w, "%s\n", string(b))
	}
}

// Fatalf always emits the given log string with time prefixed on the log
// writer, followed by calling os.Exit(1).
// Arguments are handled as in fmt.Printf.
func Fatalf(format string, args ...interface{}) {
	fmt.Fprintf(w, "[%s] %s\n", time.Now().Format(timeFormat), fmt.Sprintf(format, args...))
	os.Exit(1)
}

func init() {
	RegisterFlags(flag.CommandLine)
}
