// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiboot

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/ubinary"
	"github.com/u-root/u-root/pkg/uio"
)

// A module represents a module to be loaded along with the kernel.
type module struct {
	// Start is the inclusive start of the Module memory location
	Start uint32

	// End is the exclusive end of the Module memory location.
	End uint32

	// Cmdline is a pointer to a null-terminated ASCII string.
	Cmdline uint32

	// Reserved is always zero.
	Reserved uint32
}

type modules []module

func (m *multiboot) loadModules() (modules, error) {
	loaded, data, err := loadModules(m.modules)
	if err != nil {
		return nil, err
	}

	cmdlineRange, err := m.mem.AddKexecSegment(data)
	if err != nil {
		return nil, err
	}

	loaded.fix(uint32(cmdlineRange.Start))
	m.loadedModules = loaded

	for i, mod := range loaded {
		log.Printf("Added module %s at [%#x, %#x)", m.modules[i].Name(), mod.Start, mod.End)
	}

	return loaded, nil
}

func (m *multiboot) addMultibootModules() (uintptr, error) {
	loaded, err := m.loadModules()
	if err != nil {
		return 0, err
	}
	b, err := loaded.marshal()
	if err != nil {
		return 0, err
	}
	modRange, err := m.mem.AddKexecSegment(b)
	if err != nil {
		return 0, err
	}
	return modRange.Start, nil
}

// loadModules loads module files.
// Returns loaded modules description and buffer storing loaded modules.
// Memory layout of the loaded modules is following:
//			cmdLine_1
//			cmdLine_2
//			...
//			cmdLine_n
//			<padding>
//			modules_1
//			<padding>
//			modules_2
//			...
//			<padding>
//			modules_n
//
// <padding> aligns the start of each module to a page beginning.
func loadModules(rmods []Module) (loaded modules, data []byte, err error) {
	loaded = make(modules, len(rmods))
	var buf bytes.Buffer

	for i, rmod := range rmods {
		if err := loaded[i].setCmdline(&buf, rmod.Cmdline); err != nil {
			return nil, nil, err
		}
	}

	for i, rmod := range rmods {
		if err := loaded[i].loadModule(&buf, rmod.Module); err != nil {
			return nil, nil, fmt.Errorf("error adding module %v: %v", rmod.Name(), err)
		}
	}
	return loaded, buf.Bytes(), nil
}

func pageAlign(val uint32) uint32 {
	mask := uint32(os.Getpagesize() - 1)
	return (val + mask) &^ mask
}

// pageAlignBuf pads buf to a page boundary.
func pageAlignBuf(buf *bytes.Buffer) error {
	mask := (os.Getpagesize() - 1)
	size := (buf.Len() + mask) &^ mask
	_, err := buf.Write(bytes.Repeat([]byte{0}, size-buf.Len()))
	return err
}

func (m *module) loadModule(buf *bytes.Buffer, r io.ReaderAt) error {
	// place start of each module to a beginning of a page.
	if err := pageAlignBuf(buf); err != nil {
		return err
	}

	m.Start = uint32(buf.Len())

	if _, err := io.Copy(buf, uio.Reader(r)); err != nil {
		return err
	}

	m.End = uint32(buf.Len())
	return nil
}

func (m *module) setCmdline(buf *bytes.Buffer, cmdLine string) error {
	m.Cmdline = uint32(buf.Len())
	if _, err := buf.WriteString(cmdLine); err != nil {
		return err
	}
	return buf.WriteByte(0)
}

// fix fixes pointers converting relative values to absolute values.
func (m modules) fix(base uint32) {
	for i := range m {
		m[i].Start += base
		m[i].End += base
		m[i].Cmdline += base
	}
}

// marshal writes out the module list in multiboot info format, as described in
// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html#Boot-information-format
func (m modules) marshal() ([]byte, error) {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, ubinary.NativeEndian, m)
	return buf.Bytes(), err
}

// elems adds mutiboot info elements describing where to find each module and
// its cmdline.
func (m modules) elems() []elem {
	var e []elem
	for _, mm := range m {
		e = append(e, &mutibootModule{
			cmdline:    uint64(mm.Cmdline),
			moduleSize: uint64(mm.End - mm.Start),
			ranges: []mutibootModuleRange{
				{
					startPageNum: uint64(mm.Start / uint32(os.Getpagesize())),
					numPages:     pageAlign(mm.End-mm.Start) / uint32(os.Getpagesize()),
				},
			},
		})
	}
	return e
}
