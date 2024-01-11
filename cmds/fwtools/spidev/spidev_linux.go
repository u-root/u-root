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
//	-s SPEED: max speed in Hz (default 500000)
//
// Description:
//
//	raw: The binary data from stdin is transmitted over the SPI bus.
//	     Received data is printed to stdout.
//	sfdp: Parse and print the parameters in the SFDP.
//	id: print the 3 byte hex id
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/flash"
	"github.com/u-root/u-root/pkg/flash/op"
	"github.com/u-root/u-root/pkg/flash/sfdp"
	"github.com/u-root/u-root/pkg/spidev"
)

type spi interface {
	Transfer([]spidev.Transfer) error
	SetSpeedHz(uint32) error
	Close() error
}

var (
	errFlag    = errors.New("unknown flag")
	errCommand = errors.New("usage")
)

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
		if err == flag.ErrHelp {
			return fmt.Errorf("%w:<raw|sfdp|id>", errCommand)
		}
		return fmt.Errorf("%w:%v", errFlag, err)
	}

	if fs.NArg() != 1 {
		return fmt.Errorf("%w: %s <raw|sfdp|id>", errCommand, fs.FlagUsages())
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

	cmd := fs.Arg(0)
	// Currently, only the raw subcommand is supported.
	switch cmd {
	case "id":
		// Wake it up, then get the id.
		// PRDRES is not universally handled on all devices, but that's ok.
		// but CE MUST drop, so we structure this as two separate
		// transfers to ensure that happens.
		transfers := []spidev.Transfer{
			{
				Tx:       op.PRDRES.Bytes(),
				Rx:       make([]byte, 1),
				CSChange: true,
			},
			{
				Tx: append(op.ReadJEDECID.Bytes(), 0, 0, 0),
				Rx: make([]byte, 4),
			},
		}

		if err := s.Transfer(transfers); err != nil {
			return err
		}

		fmt.Fprintf(output, "%02x\n", transfers[1].Rx[1:])
		return nil

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
		return fmt.Errorf("%s:%w:%s", cmd, errCommand, fs.FlagUsages())
	}
}

func main() {
	if err := run(os.Args[1:], openSPIDev, os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
