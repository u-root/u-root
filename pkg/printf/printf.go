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
	writer io.Writer
}

// NewPrinterFromArgs returns a printf using the args provided
// it will error if the length of args is below 1. it will use the first element of args as the format, and the remaining as the params
func NewPrinterFromArgs(writer io.Writer, args []string) (*printf, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("%w: %w", ErrPrintf, ErrNotEnoughArguments)
	}
	var params []string
	format := args[0]
	if len(args) > 1 {
		params = args[1:]
	}
	return NewPrinter(writer, format, params), nil
}

// NewPrinter returns a printf
func NewPrinter(writer io.Writer, format string, params []string) *printf {
	o := &printf{
		writer: writer,
		format: format,
		params: params,
	}
	return o
}

// Run processes the printf command with the format and parameters.
// it will not have written any bytes to the writer if err is not nil
// if err is nil, run may or may not write bytes to the writer
func (c *printf) Run() error {
	w := new(bytes.Buffer)
	err := interpret(w, c.format, c.params, false, true)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrPrintf, err)
	}
	// flush on success
	_, err = w.WriteTo(c.writer)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrPrintf, err)
	}
	return nil
}
