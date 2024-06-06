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

// Devices contains a slice of one or more PCI devices
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
				// Bus: primary=00, secondary=05, subordinate=0c, sec-latency=0
				// I/O behind bridge: 00002000-00002fff [size=4K]
				// Memory behind bridge: f0000000-f1ffffff [size=32M]
				// Prefetchable memory behind bridge: 00000000f2900000-00000000f29fffff [size=1M]
				if _, err := fmt.Fprintf(o, ", Cache Line Size: %d bytes", c[CacheLineSize]); err != nil {
					return err
				}
				if _, err := fmt.Fprintf(o, "\n\tBus: primary=%02x, secondary=%02x, subordinate=%02x, sec-latency=%s",
					pci.Primary, pci.Secondary, pci.Subordinate, pci.SecLatency); err != nil {
					return err
				}
				// I hate this code.
				// I miss Rust tuples at times.
				for _, e := range []struct {
					h, f string
					b, l uint64
				}{
					{h: "\n\tI/O behind bridge: ", f: "%#08x-%#08x [size=%#x]", b: pci.IO.Base, l: pci.IO.Lim},
					{h: "\n\tMemory behind bridge: ", f: "%#08x-%#08x [size=%#x]", b: pci.Mem.Base, l: pci.Mem.Lim},
					{h: "\n\tPrefetchable memory behind bridge: ", f: "%#08x-%#08x [size=%#x]", b: pci.PrefMem.Base, l: pci.PrefMem.Lim},
				} {
					s := e.h + " [disabled]"
					if e.b != 0 {
						sz := e.l - e.b + 1
						s = fmt.Sprintf(e.h+e.f, e.b, e.l, sz)
					}
					if _, err := fmt.Fprint(o, s); err != nil {
						return err
					}
				}
			}
			fmt.Fprintf(o, "\n")
			if pci.IRQPin != 0 {
				if _, err := fmt.Fprintf(o, "\tInterrupt: pin %X routed to IRQ %X\n", 9+pci.IRQPin, pci.IRQLine); err != nil {
					return err
				}
			}
			if !pci.Bridge {
				for _, b := range pci.BARS {
					if _, err := fmt.Fprintf(o, "\t%v\n", b.String()); err != nil {
						return err
					}
				}
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
func (d Devices) SetVendorDeviceName(ids []Vendor) {
	for _, p := range d {
		p.SetVendorDeviceName(ids)
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
