// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package printf

import (
	"bytes"
	"fmt"
	"io"
)

type Printfer interface {
	// Printf processes the printf command with the format and parameters, writing to its writer.
	// it will not have written any bytes to the writer if err is not nil
	// if err is nil, run may or may not write bytes to the writer
	Printf(format string, params ...string) (int64, error)
}

type printf struct {
	writer io.Writer
}

// NewPrinterFromArgs returns a printf using the args provided
// it will error if the length of args is below 1. it will use the first element of args as the format, and the remaining as the params
func Fprintf(writer io.Writer, args ...string) (int64, error) {
	if len(args) < 1 {
		return 0, fmt.Errorf("%w: %w", ErrPrintf, ErrNotEnoughArguments)
	}
	var params []string
	format := args[0]
	if len(args) > 1 {
		params = args[1:]
	}
	return NewPrinter(writer).Printf(format, params...)
}

// NewPrinter returns a printfer
func NewPrinter(writer io.Writer) Printfer {
	o := &printf{
		writer: writer,
	}
	return o
}

// Printf processes the printf command with the format and parameters, writing to its writer.
// it will not have written any bytes to the writer if err is not nil
// if err is nil, run may or may not write bytes to the writer
func (c *printf) Printf(format string, params ...string) (int64, error) {
	w := new(bytes.Buffer)
	err := interpret(w, []byte(format), params, false, true)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrPrintf, err)
	}
	// flush on success
	n, err := w.WriteTo(c.writer)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrPrintf, err)
	}
	return n, nil
}
