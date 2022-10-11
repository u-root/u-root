// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// spidev communicates with the Linux spidev driver.
//
// Synopsis:
//
//	spidev [OPTIONS] raw < tx.bin > rx.bin
//	spidev [OPTIONS] sfdp
//
// Options:
//
//	-D DEV: spidev device (default /dev/spidev0.0)
//	-s SPEED: max speed in Hz (default 500000)
//
// Description:
//
//	raw: The binary data from stdin is transmitted over the SPI bus.
//	     Received data is printed to stdout.
//	sfdp: Parse and print the parameters in the SFDP.
package main

import (
	"errors"
	"io"
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/flash"
	"github.com/u-root/u-root/pkg/flash/sfdp"
	"github.com/u-root/u-root/pkg/spidev"
)

type spi interface {
	Transfer([]spidev.Transfer) error
	SetSpeedHz(uint32) error
	Close() error
}

type spiOpenFunc func(dev string) (spi, error)

func openSPIDev(dev string) (spi, error) {
	return spidev.Open(dev)
}

func run(args []string, spiOpen spiOpenFunc, input io.Reader, output io.Writer) error {
	// Parse args.
	fs := flag.NewFlagSet("spidev", flag.ContinueOnError)
	dev := fs.StringP("device", "D", "/dev/spidev0.0", "spidev device")
	speed := fs.Uint32P("speed", "s", 500000, "max speed in Hz")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		flag.Usage()
		return errors.New("expected one subcommand")
	}

	// Open the spi device.
	s, err := spiOpen(*dev)
	if err != nil {
		return err
	}
	defer s.Close()
	if err := s.SetSpeedHz(*speed); err != nil {
		return err
	}

	// Currently, only the raw subcommand is supported.
	switch fs.Args()[0] {
	case "raw":
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
		flag.Usage()
		return errors.New("unknown subcommand")
	}
}

func main() {
	if err := run(os.Args[1:], openSPIDev, os.Stdin, os.Stdout); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
