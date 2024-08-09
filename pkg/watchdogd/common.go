// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watchdogd

import (
	"errors"
	"fmt"
)

// Common code shared between both tinygo and golang implementations

const (
	opStopPettingTimeoutSeconds = 10

	// oopsBuffSize is the size of the buffer for reading the kernel logs.
	OopsBuffSize = 256 * 1024
)

// Define Custom Error for Watchdog Petting
//
//	OpResultError     = 'E' // Error.
//	OpResultInvalidOp = 'I' // Invalid Op.
type WatchdogError struct {
	Err error
}

func (e WatchdogError) Error() string {
	return fmt.Sprintf("watchdog error: %v", e.Err)
}

func NewWatchdogError(msg string) WatchdogError {
	return WatchdogError{Err: errors.New(msg)}
}

func NewWatchdogErrorf(format string, a ...any) WatchdogError {
	return WatchdogError{Err: fmt.Errorf(format, a...)}
}

var (
	ErrInvalidOp = WatchdogError{Err: errors.New("invalid operation")}
	Error        = WatchdogError{Err: errors.New("error")}
)

type WatchdogOperation uint8

// Define states for watchdog state machine
const (
	OpStop WatchdogOperation = iota
	OpContinue
	OpDisarm
	OpArm
	OpErr
)

// Create a new WatchdogOperation from a byte.
// If the byte does not specify a valid operation, return an error.
func NewWatchdogOperation(b byte) (WatchdogOperation, WatchdogError) {
	if b >= byte(OpStop) && b <= byte(OpArm) {
		return WatchdogOperation(b), WatchdogError{}
	}

	return OpErr, ErrInvalidOp
}
