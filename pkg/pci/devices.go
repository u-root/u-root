// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

//Devices contains a slice of one or more PCI devices
type Devices []*PCI

// Print prints information to an io.Writer
func (d Devices) Print(o io.Writer, verbose, confSize int) error {
	for _, pci := range d {
		if _, err := fmt.Fprintf(o, "%s\n", pci.String()); err != nil {
			return err
		}
		var extraNL bool
		// Make sure we have read enough config space to satisfy the verbose and confSize requests.
		// If len(pci.Config) is > 64, that's the only test we need.
		if (verbose > 1 || confSize > 64) && len(pci.Config) < 256 {
			return os.ErrPermission
		}
		if verbose >= 1 {
			c := pci.Config
			if _, err := fmt.Fprintf(o, "\tControl: %s\n\tStatus: %s\n\tLatency: %d", pci.Control.String(), pci.Status.String(), pci.Latency); err != nil {
				return err
			}
			if pci.Bridge {
				if _, err := fmt.Fprintf(o, ", Cache Line Size: %d bytes", c[CacheLineSize]); err != nil {
					return err
				}
			}
			fmt.Fprintf(o, "\n")
			if pci.IRQPin != 0 {
				if _, err := fmt.Fprintf(o, "\tInterrupt: pin %X routed to IRQ %s\n", 9+pci.IRQPin, pci.IRQLine); err != nil {
					return err
				}

			}
			if verbose >= 2 {
				if !pci.Bridge {
				} else {
				}

				//Latency: 0, Cache Line Size: 64 bytes
				//Interrupt: pin D routed to IRQ 19
			}
			extraNL = true
		}

		if confSize > 0 {
			r := io.LimitReader(bytes.NewBuffer(pci.Config), int64(confSize))
			e := hex.Dumper(o)
			if _, err := io.Copy(e, r); err != nil {
				return err
			}
			extraNL = true
		}
		// lspci likes that extra line of separation
		if extraNL {
			fmt.Fprintf(o, "\n")
		}
	}
	return nil
}

// SetVendorDeviceName sets all numeric IDs of all the devices
// using the pci device SetVendorDeviceName.
func (d Devices) SetVendorDeviceName() {
	for _, p := range d {
		p.SetVendorDeviceName()
	}
}

// ReadConfig reads the config info for all the devices.
func (d Devices) ReadConfig() error {
	for _, p := range d {
		if err := p.ReadConfig(); err != nil {
			return err
		}
	}
	return nil
}

// ReadConfigRegister reads the config info for all the devices.
func (d Devices) ReadConfigRegister(offset, size int64) ([]uint64, error) {
	var vals []uint64
	for _, p := range d {
		val, err := p.ReadConfigRegister(offset, size)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return vals, nil
}

// WriteConfigRegister writes the config info for all the devices.
func (d Devices) WriteConfigRegister(offset, size int64, val uint64) error {
	for _, p := range d {
		if err := p.WriteConfigRegister(offset, size, val); err != nil {
			return err
		}
	}
	return nil
}
