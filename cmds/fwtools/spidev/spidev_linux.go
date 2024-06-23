// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// spidev communicates with the Linux spidev driver.
//
// Synopsis:
//
//	spidev [OPTIONS] raw < tx.bin > rx.bin
//	spidev [OPTIONS] sfdp
//	spidev [OPTIONS] id
//
// Options:
//
//	-D DEV: spidev device (default /dev/spidev0.0)
//	-s SPEED: max speed in Hz (default whatever the spi package sets)
//
// Description:
//
//	raw: The binary data from stdin is transmitted over the SPI bus.
//	     Received data is printed to stdout.
//	sfdp: Parse and print the parameters in the SFDP.
//	id: print the 3 byte hex id
package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/flash"
	"github.com/u-root/u-root/pkg/flash/chips"
	"github.com/u-root/u-root/pkg/flash/op"
	"github.com/u-root/u-root/pkg/flash/sfdp"
	"github.com/u-root/u-root/pkg/spidev"
)

type spi interface {
	Transfer([]spidev.Transfer) error
	ID() (chips.ID, error)
	Status() (op.Status, error)
	SetSpeedHz(uint32) error
	Close() error
}

var (
	// ErrCommand should be used for any error, including those from flag.Parse()
	ErrCommand = errors.New("usage:spidev [-D device] [-s speed] <raw [bytes]...|sfdp|id>")
	// ErrConvert is for any type of conversion error
	ErrConvert = errors.New("bad syntax")
)

type spiOpenFunc func(dev string) (spi, error)

func openSPIDev(dev string) (spi, error) {
	return spidev.Open(dev)
}

func run(args []string, spiOpen spiOpenFunc, input io.Reader, output io.Writer) error {
	fs := flag.NewFlagSet("spidev <raw [bytes]...|sfdp|id>", flag.ContinueOnError)
	// Usage spews a lot at the wrong time, dirtying up test output.
	// It's also not very controllable. We just print the message in our own way.
	fs.Usage = func() {}
	fs.SetOutput(io.Discard)
	dev := fs.String("D", "/dev/spidev0.0", "spidev device")
	speed := fs.Uint("s", 0, "max speed in Hz")
	if err := fs.Parse(args); err != nil {
		return ErrCommand
	}
	if fs.NArg() == 0 {
		return ErrCommand
	}
	cmd := fs.Arg(0)
	switch cmd {
	case "id", "sfdp":
		if fs.NArg() != 1 {
			return fmt.Errorf("%w: id and sfdp do not require an argument", ErrCommand)
		}
	}
	// Open the spi device.
	s, err := spiOpen(*dev)
	if err != nil {
		return err
	}
	defer s.Close()

	// Note that spidev.Open sets a safe default speed, known to
	// work, that is conservative. In some cases, users might wish
	// to override that speed. Since the speed can be set any number
	// of times, this is a safe operation.
	if *speed != 0 {
		if err := s.SetSpeedHz(uint32(*speed)); err != nil {
			return err
		}
	}

	// Currently, only the raw subcommand is supported.
	switch cmd {
	case "id":
		id, err := s.ID()
		if err != nil {
			return err
		}
		fmt.Fprintf(output, "%02x\n", id)
		return nil

	case "raw":
		for _, a := range fs.Args()[1:] {
			b, err := hex.DecodeString(a)
			if err != nil {
				return fmt.Errorf("%v:%w", err, ErrConvert)
			}
			transfers := []spidev.Transfer{
				{
					Tx: b,
					Rx: make([]byte, len(b)),
				},
			}
			// Perform transfers.
			if err := s.Transfer(transfers); err != nil {
				return err
			}

		}
		// Create transfer from stdin.
		tx, err := io.ReadAll(input)
		if err != nil {
			return err
		}
		if len(tx) == 0 {
			return nil
		}
		transfers := []spidev.Transfer{
			{
				Tx: tx,
				Rx: make([]byte, len(tx)),
			},
		}

		// Perform transfers.
		if err := s.Transfer(transfers); err != nil {
			return err
		}

		_, err = output.Write(transfers[0].Rx)
		return err

	case "sfdp":
		// Create flash device and read SFDP.
		f, err := flash.New(s)
		if err != nil {
			return err
		}

		// Print sfdp.
		return f.SFDP().PrettyPrint(output, sfdp.BasicTableLookup)

	default:
		return fmt.Errorf("%s:%w:%v", cmd, ErrCommand, err)
	}
}

func main() {
	if err := run(os.Args[1:], openSPIDev, os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
