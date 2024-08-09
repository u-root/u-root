// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build plan9

package memio

import (
	"fmt"
	"os"
)

const (
	p9pathIOL = "#P/iol"
	p9pathIOW = "#P/iow"
	p9pathIOB = "#P/iob"
)

type Plan9Port struct {
	iol ReadWriteCloser
	iow ReadWriteCloser
	iob ReadWriteCloser
}

var _ PortReadWriter = &Plan9Port{}

// In reads data from the x86 port at address addr. Data must be Uint8, Uint16,
// Uint32, but not Uint64.
func (p *Plan9Port) In(addr uint16, data UintN) error {
	switch data.(type) {
	case *Uint32:
		return p.iol.Read(data, int64(addr))
	case *Uint16:
		return p.iow.Read(data, int64(addr))
	case *Uint8:
		return p.iob.Read(data, int64(addr))
	}
	return fmt.Errorf("port data must be 8, 16 or 32 bits")
}

// Out writes data to the x86 port at address addr. data must be Uint8, Uint16
// uint32, but not Uint64.
func (p *Plan9Port) Out(addr uint16, data UintN) error {
	switch data.(type) {
	case *Uint32:
		return p.iol.Write(data, int64(addr))
	case *Uint16:
		return p.iow.Write(data, int64(addr))
	case *Uint8:
		return p.iob.Write(data, int64(addr))
	}
	return fmt.Errorf("port data must be 8, 16 or 32 bits")
}

func (p *Plan9Port) Close() error {
	if err := p.iol.Close(); err != nil {
		return err
	}
	if err := p.iow.Close(); err != nil {
		return err
	}
	return p.iob.Close()
}

func NewPort() (*Plan9Port, error) {
	f1, err := os.OpenFile(p9pathIOL, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	f2, err := os.OpenFile(p9pathIOW, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	f3, err := os.OpenFile(p9pathIOB, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	return &Plan9Port{
		iol: NewMemIOPort(f1),
		iow: NewMemIOPort(f2),
		iob: NewMemIOPort(f3),
	}, nil
}
