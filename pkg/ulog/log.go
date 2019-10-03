// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ulog exposes logging via a Go interface.
//
// ulog has three implementations of the Logger interface: a Go standard
// library "log" package Logger, a kernel syslog (dmesg) Logger, and a test
// Logger that logs via a test's testing.TB.Logf.
package ulog

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"testing"

	"golang.org/x/sys/unix"
)

// Logger is a log receptacle.
//
// It puts your information somewhere for safekeeping.
type Logger interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})
}

// Log is a Logger that prints to stderr, like the default log package.
var Log = log.New(os.Stderr, "", log.LstdFlags)

// KernelLog is a logger that prints to the kernel syslog buffer.
//
// If the syslog buffer cannot be written to, KernelLog falls back to Log.
var KernelLog = &KLog{
	// Default log level is Info.
	LogLevel: uintptr(KLogInfo),
}

func init() {
	KernelLog.Reinit()
}

// KLog is a logger to the kernel syslog buffer.
type KLog struct {
	// FD for /dev/kmsg if it was openable.
	*os.File

	// LogLevel is the LogLevel to print with.
	//
	// Should only be accessed atomically.
	LogLevel uintptr
}

// Reinit reopens the /dev/kmsg file.
func (k *KLog) Reinit() {
	f, _ := os.OpenFile("/dev/kmsg", os.O_RDWR, 0)
	KernelLog.File = f
}

// writeString returns true iff it was able to write the log to /dev/kmsg.
func (k *KLog) writeString(s string) bool {
	if k.File == nil {
		return false
	}
	if _, err := k.File.WriteString(fmt.Sprintf("<%d>%s", atomic.LoadUintptr(&k.LogLevel), s)); err != nil {
		return false
	}
	return true
}

func (k *KLog) Printf(format string, v ...interface{}) {
	if !k.writeString(fmt.Sprintf(format, v...)) {
		Log.Printf(format, v...)
	}
}

func (k *KLog) Print(v ...interface{}) {
	if !k.writeString(fmt.Sprint(v...)) {
		Log.Print(v...)
	}
}

// KLogLevel are the log levels used by printk.
type KLogLevel uintptr

// These are the log levels used by printk as described in syslog(2).
const (
	KLogEmergency KLogLevel = 0
	KLogAlert     KLogLevel = 1
	KLogCritical  KLogLevel = 2
	KLogError     KLogLevel = 3
	KLogWarning   KLogLevel = 4
	KLogNotice    KLogLevel = 5
	KLogInfo      KLogLevel = 6
	KLogDebug     KLogLevel = 7
)

const (
	_SYSLOG_ACTION_READ_ALL      = 3
	_SYSLOG_ACTION_READ_CLEAR    = 4
	_SYSLOG_ACTION_CLEAR         = 5
	_SYSLOG_ACTION_CONSOLE_LEVEL = 8
)

// SetConsoleLogLevel sets the console level with syslog(2).
//
// After this call, only messages with a level value lower than the one
// specified will be printed to console by the kernel.
func (k *KLog) SetConsoleLogLevel(level KLogLevel) error {
	if _, _, err := unix.Syscall(unix.SYS_SYSLOG, _SYSLOG_ACTION_CONSOLE_LEVEL, 0, uintptr(level)); err != 0 {
		return fmt.Errorf("could not set syslog level to %d: %v", level, err)
	}
	return nil
}

// SetLogLevel sets the level that Printf and Print log to syslog with.
func (k *KLog) SetLogLevel(level KLogLevel) {
	atomic.StoreUintptr(&k.LogLevel, uintptr(level))
}

// ClearLog clears kernel logs back to empty.
func (k *KLog) ClearLog() error {
	_, err := unix.Klogctl(_SYSLOG_ACTION_CLEAR, nil)
	return err
}

// Read reads from the tail of the kernel log.
func (k *KLog) Read(b []byte) (int, error) {
	return unix.Klogctl(_SYSLOG_ACTION_READ_ALL, b)
}

// ReadClear reads from the tail of the kernel log and clears what was read.
func (k *KLog) ReadClear(b []byte) (int, error) {
	return unix.Klogctl(_SYSLOG_ACTION_READ_CLEAR, b)
}

// TestLogger is a Logger implementation that logs to testing.TB.Logf.
type TestLogger struct {
	TB testing.TB
}

func (tl TestLogger) Printf(format string, v ...interface{}) {
	tl.TB.Logf(format, v...)
}

func (tl TestLogger) Print(v ...interface{}) {
	tl.TB.Log(v...)
}
