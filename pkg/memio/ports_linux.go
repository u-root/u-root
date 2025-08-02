// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && amd64) || (linux && 386)

package memio

import (
	"fmt"
	"os"
)

var linuxPath = "/dev/port"

// LinuxPort implements ReadWriteCloser for Linux devices.
type LinuxPort struct {
	ReadWriteCloser
}

var _ PortReadWriter = &LinuxPort{}

// In reads data from the x86 port at address addr. Data must be Uint8, Uint16,
// Uint32, but not Uint64.
func (p *LinuxPort) In(addr uint16, data UintN) error {
	if _, ok := data.(*Uint8); !ok {
		return fmt.Errorf("/dev/port data must be 8 bits on Linux")
	}
	return p.ReadWriteCloser.Read(data, int64(addr))
}

// Out writes data to the x86 port at address addr. data must be Uint8, Uint16
// uint32, but not Uint64.
func (p *LinuxPort) Out(addr uint16, data UintN) error {
	if _, ok := data.(*Uint8); !ok {
		return fmt.Errorf("/dev/port data must be 8 bits on Linux")
	}
	return p.ReadWriteCloser.Write(data, int64(addr))
}

// Close implements Close.
func (p *LinuxPort) Close() error {
	return p.ReadWriteCloser.Close()
}

// NewPort returns a new instance of LinuxPort for read/write operations on /dev/port
func NewPort() (*LinuxPort, error) {
	var _ PortReadWriter = &LinuxPort{}
	f, err := os.OpenFile(linuxPath, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	memPort := NewMemIOPort(f)
	return &LinuxPort{
		ReadWriteCloser: memPort,
	}, nil
}

// In is deprecated. Only here for compatibility. Use NewPort() and the interface functions instead.
func In(addr uint16, data UintN) error {
	port, err := NewPort()
	if err != nil {
		return err
	}
	defer port.Close()
	return port.In(addr, data)
}

// Out is deprecated. Only here for compatibility. Use NewPort() and the interface functions instead.
func Out(addr uint16, data UintN) error {
	port, err := NewPort()
	if err != nil {
		return err
	}
	defer port.Close()
	return port.Out(addr, data)
}
