// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64 && linux

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/u-root/u-root/pkg/pci"
)

const (
	addr  = 0x5c
	dlow  = 0x98
	dhigh = 0x9c
	bad   = 0xffffffff
)

type instanceType uint8

// These constants are target id types
const (
	CCM     instanceType = 0
	GCM     instanceType = 1
	NCM     instanceType = 2
	IOMS    instanceType = 3
	CS      instanceType = 4
	NCS     instanceType = 5
	TCDX    instanceType = 6
	PIE     instanceType = 7
	SPF     instanceType = 8
	LLC     instanceType = 9
	CAKE    instanceType = 0xA
	UNKNOWN instanceType = 0xb
)

type (
	nodeID   uint8
	fabricID uint8
)

const (
	// BROAD is our way of saying "broadcast address"
	BROAD nodeID = 0xff
)

func (i *instanceType) String() string {
	switch *i {
	case CCM:
		return "CCM"
	case GCM:
		return "GCM"
	case NCM:
		return "NCM"
	case IOMS:
		return "IOMS"
	case CS:
		return "CS"
	case NCS:
		return "NCS"
	case TCDX:
		return "TCDX"
	case PIE:
		return "PIE"
	case SPF:
		return "SPF"
	case LLC:
		return "LLC"
	case CAKE:
		return "CAKE"
	case UNKNOWN:
		return "UNKNOWN"
	default:
		return fmt.Sprintf("Invalid type %#x", i)
	}
}

type cfg uint32

var (
	node  = flag.Uint("node", 0, "which node")
	debug = flag.Bool("d", false, "debug prints")
	v     = func(string, ...interface{}) {}
)

type config struct {
	NumCPU uint
}

func (c *config) String() string {
	return fmt.Sprintf("Num CPUS: %d", c.NumCPU)
}

// Unmarshall implements Unmarshall for a cfg.
func (c cfg) Unmarshal() *config {
	i := uint32(c)
	return &config{
		NumCPU: uint((i>>27)&1 + 1),
	}
}

// Fabric is a single component in the fabric.
type Fabric struct {
	InstanceID   nodeID
	InstanceType instanceType
	Enabled      bool
	FabricID     uint8
}

func (f *Fabric) String() string {
	return fmt.Sprintf("InstanceID %#x InstanceType %s Enabled %v FabricID %#x", f.InstanceID, f.InstanceType.String(), f.Enabled, f.FabricID)
}

// DataFabric supports operations like binary.Read and binary.Write.
// N.B. it is NOT safe to implement something like readat/writeat; arbitrary
// byte reads are not acceptable for this interface.
// The most one can do, as for most such interfaces, is 32- or 64-bit reads/writes.
type DataFabric struct {
	Node           uint8
	PCI            *pci.PCI
	Config         *config
	TotalCount     uint
	PIECount       uint
	IOMSCount      uint
	DiesPerSocket  uint
	CCM0InstanceID nodeID
	Components     []*Fabric
}

func (d *DataFabric) String() string {
	s := fmt.Sprintf("Node %d TotalCount %d PIECount %d IOMSCount %d DiesPerSocket %d CCM0InstanceID %d Config %s",
		d.Node, d.TotalCount, d.PIECount, d.IOMSCount, d.DiesPerSocket, d.CCM0InstanceID, d.Config)
	for _, f := range d.Components {
		s += "\n" + f.String()
	}
	return s
}

func (d *DataFabric) address(id nodeID, fun uint8, off uint16) error {
	if fun > 7 {
		return fmt.Errorf("fun is %#x but must be < 8", fun)
	}
	if off&3 != 0 {
		return fmt.Errorf("target id is %#x but must be 8-byte aligned", off)
	}
	if off >= 2048 {
		return fmt.Errorf("off %#x must be < %#x", off, 2048)
	}
	var targ uint32
	switch id {
	case BROAD:
	default:
		targ = 1 | (uint32(id) << 16)
	}
	targ |= (uint32(fun) << 11) | uint32(off)
	v("address: id %#x fun %#x off %#x -> targ %08x", id, fun, off, targ)
	return d.PCI.WriteConfigRegister(addr, 32, uint64(targ))
}

// ReadIndirect uses the indirection registers to read from one thing.
func (d *DataFabric) ReadIndirect(id nodeID, fun uint8, off uint16) (uint64, error) {
	if err := d.address(id, fun, off); err != nil {
		return bad, err
	}
	return d.PCI.ReadConfigRegister(dlow, 32)
}

// ReadBroadcast does the stupid merged read where all bits from all things
// are or'ed. Not sure why anyone ever thought this was useful.
func (d *DataFabric) ReadBroadcast(fun uint8, off uint16) (uint64, error) {
	return d.ReadIndirect(BROAD, fun, off)
}

func New(n uint8) (*DataFabric, error) {
	if n > 1 {
		return nil, fmt.Errorf("node is %d, but can only be 0 or 1", n)
	}
	devName := fmt.Sprintf("0000:00:%02x.4", n+0x18)
	r, err := pci.NewBusReader(devName)
	if err != nil {
		return nil, fmt.Errorf("NewBusReader: %w", err)
	}
	devs, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", devName, err)
	}
	if len(devs) != 1 {
		return nil, fmt.Errorf("%q matches more than one device", devName)
	}

	d := &DataFabric{
		Node:           n,
		PCI:            devs[0],
		TotalCount:     0,
		PIECount:       0,
		IOMSCount:      0,
		DiesPerSocket:  1,
		CCM0InstanceID: 0xff,
	}
	c, err := d.ReadBroadcast(1, 0x200)
	if err != nil {
		return nil, err
	}
	v("config is %#x", c)
	d.Config = cfg(c).Unmarshal()
	if c, err = d.ReadBroadcast(0, 0x40); err != nil {
		return nil, fmt.Errorf("read TotalCount: %w", err)
	}
	d.TotalCount = uint(uint8(c))
	v("Reg 40 is %#x, totalcount %d", c, d.TotalCount)
	for i := nodeID(0); i < 255 && len(d.Components) < int(d.TotalCount); i++ {
		info0, err := d.ReadIndirect(i, 0, 0x44)
		v("%#x: info0 %#x bit 6 %#x err %v", i, info0, info0&(1<<6), err)
		if err != nil {
			continue
		}
		enabled := (info0 & (1 << 6)) != 0

		if info0 == 0 && len(d.Components) > 0 {
			continue
		}
		InstanceType := func() instanceType {
			switch info0 & 0xf {
			case 0:
				if d.CCM0InstanceID == 0xff {
					d.CCM0InstanceID = nodeID(i)
				}
				return CCM
			case 1:
				return GCM
			case 2:
				return NCM
			case 3:
				d.IOMSCount++
				return IOMS
			case 4:
				return CS
			case 5:
				return NCS
			case 6:
				return TCDX
			case 7:
				d.PIECount++
				return PIE
			case 8:
				return SPF
			case 9:
				return LLC
			case 10:
				return CAKE
			default:
				return UNKNOWN
			}
		}()

		ids, err := d.ReadIndirect(i, 0, 0x50)
		if err != nil {
			log.Printf("ReadIndirect(%d, 0, 0x50): %v", i, err)
			continue
		}
		instanceID := uint(ids)
		fabricID := uint8((ids >> 8) & 0x3f)
		if instanceID == 0 && len(d.Components) > 0 {
			log.Printf("WTF")
			continue
		}
		d.Components = append(d.Components, &Fabric{
			InstanceID:   i,
			InstanceType: InstanceType,
			Enabled:      enabled,
			FabricID:     fabricID, // but see below. I am liking Some and None ...
		})
		// result.components.push(FabricComponent { instance_id, instance_type, enabled, fabric_id: if fabric_id != 0 || result.components.len() == 0 { Some(fabric_id) } else { None } }).unwrap();

	}

	return d, nil
}

func main() {
	flag.Parse()
	if *debug {
		v = log.Printf
	}

	df, err := New(uint8(*node))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("df is %s", df)
}
