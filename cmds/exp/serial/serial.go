// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.bug.st/serial"
	"golang.org/x/term"
)

type params struct {
	device   string
	baud     int
	parity   serial.Parity
	databits int
}

var errUsage = errors.New("usage: serial -D=/dev/tty -b=115200 -p=no -d=8")

func parseParams(device, parity string, baud uint, databits int) (params, error) {
	p := params{baud: int(baud)}
	if device == "" {
		return p, errUsage
	}
	p.device = device

	switch parity {
	case "no":
		p.parity = serial.NoParity
	case "odd":
		p.parity = serial.OddParity
	case "even":
		p.parity = serial.EvenParity
	default:
		return p, errUsage
	}

	switch databits {
	case 5, 6, 7, 8:
		p.databits = databits
	default:
		return p, errUsage
	}

	return p, nil
}

func main() {
	device := flag.String("D", "", "device: -D=/dev/tty")
	baud := flag.Uint("b", 115200, "baud: -b=115200")
	parity := flag.String("p", "no", "parity: -p=no|even|odd")
	databits := flag.Int("d", 8, "databits: -d=5|6|7|8")

	flag.Parse()
	p, err := parseParams(*device, *parity, *baud, *databits)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Failed to set terminal to raw mode: %v\n", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), state)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		term.Restore(int(os.Stdin.Fd()), state)
		os.Exit(0)
	}()

	mode := &serial.Mode{
		BaudRate: p.baud,
		DataBits: p.databits,
		Parity:   p.parity,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(p.device, mode)
	if err != nil {
		term.Restore(int(os.Stdin.Fd()), state)
		log.Fatal(err)
	}

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := port.Read(buf)
			if err != nil {
				continue
			}

			fmt.Printf("%v", string(buf[:n]))
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			continue
		}
		if err == io.EOF {
			c <- os.Interrupt
		}

		for i := range n {
			// Control-X to quit
			if buf[i] == 0x18 {
				c <- os.Interrupt
			}
		}

		_, err = port.Write(buf[:n])
		if err != nil {
			log.Printf("write: %v\r\n", err)
			continue
		}
	}
}
