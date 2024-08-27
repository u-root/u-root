// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64 || 386

package cmos

import (
	"github.com/u-root/u-root/pkg/memio"
)

const (
	regPort  = 0x70
	dataPort = 0x71
)

type Chip struct {
	memio.PortReadWriter
}

// Read reads a register reg from CMOS into data.
func (c *Chip) Read(reg memio.Uint8, data memio.UintN) error {
	if err := c.PortReadWriter.Out(regPort, &reg); err != nil {
		return err
	}
	return c.PortReadWriter.In(dataPort, data)
}

// Write writes value data into CMOS register reg.
func (c *Chip) Write(reg memio.Uint8, data memio.UintN) error {
	if err := c.PortReadWriter.Out(regPort, &reg); err != nil {
		return err
	}
	return c.PortReadWriter.Out(dataPort, data)
}

// GetCMOS() returns the struct to call Read and Write functions for CMOS
// associated with the correct functions of memio.In and memio.Out
func New() (*Chip, error) {
	pr, err := memio.NewPort()
	if err != nil {
		return nil, err
	}
	return &Chip{
		PortReadWriter: pr,
	}, nil
}
