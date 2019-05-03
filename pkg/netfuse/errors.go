// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netfuse

import (
	"errors"
	"fmt"
	"strconv"
	"syscall"
)

// ErrToString converts an error to a string. If the error is a
// syscall.Errno, it returns a 10-character hex representation.
// Else it returns err.Error().
func ErrToString(err error) (s string) {
	defer SrvDebug("err2string: %v %T -> %v %T", err, err, s, s)
	if err == nil {
		return s
	}
	switch s := err.(type) {
	case syscall.Errno:
		return fmt.Sprintf("%#08x", int(s))
	default:
		return err.Error()
	}
}

// StringToErr converts a string to an error. If the
// string is 10 characters long and starts with 0x,
// and strconv.ParseInt gets no errors, it returns
// syscall.Errno of that number; else it returns
// errors.New of the string.
func StringToErr(s string) (err error) {
	defer ClntDebug("string2err: %v -> %v %T", s, err, err)
	if s == "" {
		return nil
	}
	if len(s) == 10 && s[:2] == "0x" {
		errno, err := strconv.ParseInt(s, 0, 0)
		if err == nil {
			return syscall.Errno(errno)
		}
	}
	return errors.New(s)
}
