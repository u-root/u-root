// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package printf

import (
	"bytes"
	"fmt"
	"io"
)

type printf struct {
	format string
	params []string
	stdout io.Writer
}

func NewPrinterFromArgs(stdout io.Writer, args []string) (*printf, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("%w: %w", ErrPrintf, ErrNotEnoughArguments)
	}
	var params []string
	format := args[0]
	if len(args) > 1 {
		params = args[1:]
	}
	return NewPrinter(stdout, format, params), nil
}

func NewPrinter(stdout io.Writer, format string, params []string) *printf {
	o := &printf{
		stdout: stdout,
		format: format,
		params: params,
	}
	return o
}

func (c *printf) Run() error {
	w := new(bytes.Buffer)
	err := interpret(w, c.format, c.params, false, true)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrPrintf, err)
	}
	// flush on success
	_, err = w.WriteTo(c.stdout)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrPrintf, err)
	}
	return nil
}
