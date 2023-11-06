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

// Fprintf will output to writer, using first of args as format
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

func Sprintf(format string, args ...string) (string, error) {
	o := &bytes.Buffer{}
	_, err := Fprintf(o, append([]string{format}, args...)...)
	if err != nil {
		return o.String(), err
	}
	return o.String(), nil
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
