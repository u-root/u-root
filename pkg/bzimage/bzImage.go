// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package bzImage implements encoding.UnmarshalBinary for bzImage files.
// The bzImage struct contains all the information about the file and can
// be used to create a new bzImage.
package bzimage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"reflect"
)

var Debug = func(string, ...interface{}) {}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// For now, it hardwires the KernelBase to 0x100000.
func (b *BzImage) UnmarshalBinary(d []byte) error {
	Debug("Processing %d byte image", len(d))
	r := bytes.NewReader(d)
	if err := binary.Read(r, binary.LittleEndian, &b.Header); err != nil {
		return err
	}
	Debug("Header was %d bytes", len(d)-r.Len())
	Debug("bzImage header %v", b.Header)
	Debug("magic %x switch %v", b.Header.HeaderMagic, b.Header.RealModeSwitch)
	if b.Header.HeaderMagic != HeaderMagic {
		return fmt.Errorf("Not a bzImage: magic should be %02x, and is %02x", HeaderMagic, b.Header.HeaderMagic)
	}
	Debug("RamDisk image %x size %x", b.Header.RamDiskImage, b.Header.RamDiskSize)
	Debug("StartSys %x", b.Header.StartSys)
	Debug("Boot type: %s(%x)", LoaderType[boottype(b.Header.TypeOfLoader)], b.Header.TypeOfLoader)
	Debug("SetupSects %d", b.Header.SetupSects)

	b.KernelOffset = (uintptr(b.Header.SetupSects) + 1) * 512
	Debug("Kernel offset is %d bytes", b.KernelOffset)
	if b.KernelOffset > uintptr(len(d)) {
		return fmt.Errorf("len(b) is %d, b.Header.SetupSects+1 * 512 is %d: too small?", len(d), b.KernelOffset)
	}
	b.BootCode = d[len(d)-r.Len() : b.KernelOffset]
	Debug("%d bytes of BootCode", len(b.BootCode))
	b.Kernel = d[b.KernelOffset:]
	Debug("Kernel at %d, %d bytes", b.KernelOffset, len(b.Kernel))
	b.KernelBase = uintptr(0x100000)
	b.InitRAMFS = d[b.Header.RamDiskImage : b.Header.RamDiskImage+b.Header.RamDiskSize]
	Debug("Ramdisk at %d, %d bytes", b.Header.RamDiskImage, b.Header.RamDiskSize)
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (b *BzImage) MarshalBinary() ([]byte, error) {
	var w bytes.Buffer
	w.Grow(int(b.KernelOffset) + len(b.Kernel))
	Debug("Grew output buffer to %d bytes", w.Len())
	if err := binary.Write(&w, binary.LittleEndian, &b.Header); err != nil {
		return nil, err
	}
	Debug("Wrote %d bytes of header", w.Len())
	if _, err := w.Write(b.BootCode); err != nil {
		return nil, err
	}
	if _, err := w.Write(b.Kernel); err != nil {
		return nil, err
	}
	Debug("Finished writing, len is now %d bytes", w.Len())
	return w.Bytes(), nil
}

func Equal(a, b []byte) error {
	if len(a) != len(b) {
		return fmt.Errorf("Images differ in len: %d bytes and %d bytes", len(a), len(b))
	}
	var ba BzImage
	if err := ba.UnmarshalBinary(a); err != nil {
		return err
	}
	var bb BzImage
	if err := bb.UnmarshalBinary(b); err != nil {
		return err
	}
	if !reflect.DeepEqual(ba.Header, bb.Header) {
		return fmt.Errorf("Headers do not match")
	}
	// this is overkill, I can't see any way it can happen.
	if len(ba.Kernel) != len(bb.Kernel) {
		return fmt.Errorf("Kernel lengths differ: %d vs %d bytes", len(ba.Kernel), len(bb.Kernel))
	}
	if len(ba.BootCode) != len(bb.BootCode) {
		return fmt.Errorf("BootCode lengths differ: %d vs %d bytes", len(ba.Kernel), len(bb.Kernel))
	}

	if !reflect.DeepEqual(ba.BootCode, bb.BootCode) {
		return fmt.Errorf("BootCode does not match")
	}
	if !reflect.DeepEqual(ba.Kernel, bb.Kernel) {
		return fmt.Errorf("Kernels do not match")
	}

	return nil
}

func (b *BzImage) AddInitRAMFS(s string) error {
	d, err := ioutil.ReadFile(s)
	if err != nil {
		return err
	}
	b.InitRAMFS = d
	return nil
}

// MakeLinuxHeader marshals a LinuxHeader into a []byte.
func MakeLinuxHeader(h *LinuxHeader) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, h)
	return buf.Bytes(), err
}
