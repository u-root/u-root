// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// spidev communicates with the Linux spidev driver.
//
// Synopsis:
//     spidev [-D DEV] [-s SPEED] raw < tx.bin > rx.bin
//
// Options:
//     -D DEV: spidev device (default /dev/spidev0.0)
//     -s SPEED: max speed in Hz (default 500000)
//
// Description:
//     With the raw subcommand, the binary data from stdin is transmitted over
//     the SPI bus. Received data is printed to stdout.
package main

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/spi"
)

type spidev interface {
	Transfer([]spi.Transfer) error
	Close() error
}

type spiOpenFunc func(dev string) (spidev, error)

func openRealSpi(dev string) (spidev, error) {
	return spi.Open(dev)
}

func run(args []string, spiOpen spiOpenFunc, input io.Reader, output io.Writer) error {
	// Parse args.
	fs := flag.NewFlagSet("spidev", flag.ContinueOnError)
	dev := fs.StringP("device", "D", "/dev/spidev0.0", "spidev device")
	speed := fs.Uint32P("speed", "s", 500000, "max speed in Hz")
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Currently, only the raw subcommand is supported.
	if fs.NArg() != 1 || fs.Args()[0] != "raw" {
		flag.Usage()
		return errors.New("expected 'raw' subcommand")
	}

	// Open the spi device.
	s, err := spiOpen(*dev)
	if err != nil {
		return err
	}
	defer s.Close()

	// Create transfer from stdin.
	tx, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}
	if len(tx) == 0 {
		return nil
	}
	transfers := []spi.Transfer{
		{
			Tx:       tx,
			Rx:       make([]byte, len(tx)),
			CSChange: true,
			SpeedHz:  *speed,
		},
	}

	// Perform transfers.
	if err := s.Transfer(transfers); err != nil {
		return err
	}

	_, err = output.Write(transfers[0].Rx)
	return err
}

func main() {
	if err := run(os.Args[1:], openRealSpi, os.Stdin, os.Stdout); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
