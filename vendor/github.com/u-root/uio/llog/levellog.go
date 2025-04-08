// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package llog is a dirt-simple leveled text logger.
package llog

import (
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"testing"
)

// Default is the stdlib default log sink.
func Default() *Logger {
	return &Logger{Sink: SinkFor(log.Printf)}
}

// Test is a logger that prints every level to t.Logf.
func Test(tb testing.TB) *Logger {
	tb.Helper()
	return &Logger{Sink: SinkFor(tb.Logf), Level: math.MinInt32}
}

// Debug prints to log.Printf at the debug level.
func Debug() *Logger {
	return &Logger{Sink: SinkFor(log.Printf), Level: slog.LevelDebug}
}

// Printf is a logger printf function.
type Printf func(format string, v ...any)

// Printfer is an interface implementing Printf.
type Printfer interface {
	Printf(format string, v ...any)
}

func lineSprintf(format string, v ...any) string {
	s := fmt.Sprintf(format, v...)
	if strings.HasSuffix(s, "\n") {
		return s
	}
	return s + "\n"
}

// WritePrintf is a Printf that prints lines to w.
func WritePrintf(w io.Writer) Printf {
	return func(format string, v ...any) {
		_, _ = io.WriteString(w, lineSprintf(format, v...))
	}
}

// MultiPrintf is a Printf that prints to all given p.
func MultiPrintf(p ...Printf) Printf {
	return func(format string, v ...any) {
		for _, q := range p {
			if q != nil {
				q(format, v...)
			}
		}
	}
}

// Sink is the output for Logger.
type Sink func(level slog.Level, format string, v ...any)

// SinkFor prepends the log with a log level and outputs to p.
func SinkFor(p Printf) Sink {
	return func(level slog.Level, format string, args ...any) {
		// Prepend log level.
		format = "%s " + format
		args = append([]any{level}, args...)
		p(format, args...)
	}
}

// Logger is a dirt-simple leveled logger.
//
// If the log level is >= Level, it logs to the given Sink.
//
// Logger or Sink may be nil in order to log nothing.
type Logger struct {
	Sink  Sink
	Level slog.Level
}

// New creates a logger from p which prepends the log level to the output and
// uses l as the default log level.
//
// Logs with level >= l will be printed using p.
func New(l slog.Level, p Printf) *Logger {
	return &Logger{
		Sink:  SinkFor(p),
		Level: l,
	}
}

// RegisterLevelFlag registers a flag that sets the given numeric level as the level.
func (l *Logger) RegisterLevelFlag(f *flag.FlagSet, flagName string) {
	f.IntVar((*int)(&l.Level), flagName, int(l.Level), "Level to log at. Lower level emits more logs. -4 = DEBUG, 0 = INFO, 4 = WARN, 8 = ERROR")
}

// RegisterVerboseFlag registers a boolean flag that, if set, assigns verboseLevel as the level.
func (l *Logger) RegisterVerboseFlag(f *flag.FlagSet, flagName string, verboseLevel slog.Level) {
	f.BoolFunc(flagName, fmt.Sprintf("If set, logs at %d level", verboseLevel), func(val string) error {
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		if b {
			l.Level = verboseLevel
		}
		return nil
	})
}

// RegisterDebugFlag registers a boolean flag that, if set, assigns LevelDebug as the level.
func (l *Logger) RegisterDebugFlag(f *flag.FlagSet, flagName string) {
	l.RegisterVerboseFlag(f, flagName, slog.LevelDebug)
}

// AtLevelFunc returns a Printf that can be passed around to log at the given level.
//
// AtLevelFunc never returns nil.
func (l *Logger) AtLevelFunc(level slog.Level) Printf {
	if l == nil || l.Sink == nil {
		return func(fmt string, args ...any) {}
	}
	return func(fmt string, args ...any) {
		l.Logf(level, fmt, args...)
	}
}

type printfer struct {
	printf Printf
}

// Printf implements Printfer.
func (p printfer) Printf(format string, v ...any) {
	p.printf(format, v...)
}

// AtLevel returns a Printfer that can be passed around to log at the given level.
//
// AtLevel never returns nil.
func (l *Logger) AtLevel(level slog.Level) Printfer {
	return printfer{printf: l.AtLevelFunc(level)}
}

// Debugf is a printf function that logs at the Debug level.
func (l *Logger) Debugf(fmt string, args ...any) {
	if l == nil {
		return
	}
	l.Logf(slog.LevelDebug, fmt, args...)
}

// Infof is a printf function that logs at the Info level.
func (l *Logger) Infof(fmt string, args ...any) {
	if l == nil {
		return
	}
	l.Logf(slog.LevelInfo, fmt, args...)
}

// Warnf is a printf function that logs at the Warn level.
func (l *Logger) Warnf(fmt string, args ...any) {
	if l == nil {
		return
	}
	l.Logf(slog.LevelWarn, fmt, args...)
}

// Errorf is a printf function that logs at the Error level.
func (l *Logger) Errorf(fmt string, args ...any) {
	if l == nil {
		return
	}
	l.Logf(slog.LevelError, fmt, args...)
}

// Logf logs at the given level.
func (l *Logger) Logf(level slog.Level, fmt string, args ...any) {
	if l == nil || l.Sink == nil {
		return
	}
	if level >= l.Level {
		l.Sink(level, fmt, args...)
	}
}
