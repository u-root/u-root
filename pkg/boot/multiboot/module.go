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

	// CmdLine is a pointer to a null-terminated ASCII string.
	CmdLine uint32

	// Reserved is always zero.
	Reserved uint32
}

type modules []module

func (m *multiboot) addModules() (uintptr, error) {
	loaded, data, err := loadModules(m.modules)
	if err != nil {
		return 0, err
	}

	cmdlineRange, err := m.mem.AddKexecSegment(data)
	if err != nil {
		return 0, err
	}

	loaded.fix(uint32(cmdlineRange.Start))

	m.loadedModules = loaded

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
	buf := bytes.Buffer{}

	for i, rmod := range rmods {
		if err := loaded[i].setCmdLine(&buf, rmod.CmdLine); err != nil {
			return nil, nil, err
		}
	}

	for i, rmod := range rmods {
		if err := loaded[i].loadModule(&buf, rmod.Module, rmod.Name); err != nil {
			return nil, nil, fmt.Errorf("error adding module %v: %v", rmod.Name, err)
		}
	}

	return loaded, buf.Bytes(), nil
}

// alignUp pads buf to a page boundary.
func alignUp(buf *bytes.Buffer) error {
	mask := (os.Getpagesize() - 1)
	size := (buf.Len() + mask) &^ mask
	_, err := buf.Write(bytes.Repeat([]byte{0}, size-buf.Len()))
	return err
}

func (m *module) loadModule(buf *bytes.Buffer, r io.ReaderAt, name string) error {
	log.Printf("Adding module %v", name)

	// place start of each module to a beginning of a page.
	if err := alignUp(buf); err != nil {
		return err
	}

	m.Start = uint32(buf.Len())

	if _, err := io.Copy(buf, uio.Reader(r)); err != nil {
		return err
	}

	m.End = uint32(buf.Len())

	return nil
}

func (m *module) setCmdLine(buf *bytes.Buffer, cmdLine string) error {
	m.CmdLine = uint32(buf.Len())
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
		m[i].CmdLine += base
	}
}

// marshal writes out the exact bytes of modules to be loaded
// along with the kernel.
func (m modules) marshal() ([]byte, error) {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, ubinary.NativeEndian, m)
	return buf.Bytes(), err
}
