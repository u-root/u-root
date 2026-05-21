// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// base64 - encode and decode base64 from stdin or file to stdout
//
// Synopsis:
//
//	base64 [-d] [FILE]
//
// Description:
//
//	Encode or decode a file to or from base64 encoding.
//	-d   decode data (default is to encode)
//	For stdin, on standard Unix systems, you can use /dev/stdin
package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

type cmd struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	params
	args []string
}

type params struct {
	decode bool
}

var errBadUsage = errors.New("usage: base64 [-d] [file]")

func decodeone(stdin io.Reader, stdout io.Writer) error {
	r := base64.NewDecoder(base64.StdEncoding, stdin)
	if _, err := io.Copy(stdout, r); err != nil {
		return fmt.Errorf("decoding: %w", err)
	}
	return nil
}

func encodeone(stdin io.Reader, stdout io.Writer) error {
	// WriteCloser is important here, from NewEncoder documentation:
	// when finished writing, the caller must Close the returned encoder
	// to flush any partially written blocks.
	wc := base64.NewEncoder(base64.StdEncoding, stdout)
	defer wc.Close()
	if _, err := io.Copy(wc, stdin); err != nil {
		return fmt.Errorf("error encoding: %w", err)
	}
	if err := wc.Close(); err != nil { // flush any remaining data
		return fmt.Errorf("closing encoder: %w", err)
	}
	if _, err := fmt.Fprintln(stdout); err != nil {
		return fmt.Errorf("error encoder writing trailing newline %w", err)
	}
	return nil
}

func run(stdin io.Reader, stdout, stderr io.Writer, args ...string) error {
	c := cmd{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}

	f := flag.NewFlagSet("base64", flag.ExitOnError)
	f.BoolVar(&c.decode, "d", false, "decode or encode the file")

	// ignore error as flag.ExitOnError to not to print it twice
	_ = f.Parse(unixflag.ArgsToGoArgs(args))
	c.args = f.Args()

	switch {
	case len(c.args) > 1:
		return fmt.Errorf("only 0 or 1 arg allowed:%w", errBadUsage)
	case len(c.args) == 0:
	case c.args[0] == "-":
	default:
		stdin, err := os.Open(c.args[0])
		if err != nil {
			return err
		}
		c.stdin = stdin
	}

	if c.decode {
		return decodeone(c.stdin, c.stdout)
	}
	return encodeone(c.stdin, c.stdout)
}

func main() {
	if err := run(os.Stdin, os.Stdout, os.Stderr, os.Args[1:]...); err != nil {
		log.Fatalf("%s: %v", os.Args[0], err)
	}
}
