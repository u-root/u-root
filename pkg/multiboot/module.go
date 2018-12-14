// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Module defines modules to be loaded along with the kernel.
// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html#Boot-information-format.

package multiboot

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/ubinary"
)

// A Module represents a module to be loaded along with the kernel.
type Module struct {
	// Start is the inclusive start of the Module memory location
	Start uint32
	// End is the exclusive end of the Module memory location.
	End uint32

	// CmdLine is a pointer to a null-terminated ASCII string.
	CmdLine uint32

	// Reserved is always zero.
	Reserved uint32
}

type modules []Module

func (m *Multiboot) addModules() (uintptr, error) {
	loaded, data, err := loadModules(m.modules)
	if err != nil {
		return 0, err
	}

	addr, err := m.mem.AddKexecSegment(data)
	if err != nil {
		return 0, err
	}

	loaded.fix(uint32(addr))

	b, err := loaded.marshal()
	if err != nil {
		return 0, err
	}
	return m.mem.AddKexecSegment(b)
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
func loadModules(cmds []string) (loaded modules, data []byte, err error) {
	loaded = make(modules, len(cmds))
	buf := bytes.Buffer{}

	for i, cmd := range cmds {
		if err := loaded[i].setCmdLine(&buf, cmd); err != nil {
			return nil, nil, err
		}
	}

	for i, cmd := range cmds {
		args := strings.Fields(cmd)
		if err := loaded[i].loadModule(&buf, args[0], cmd); err != nil {
			return nil, nil, fmt.Errorf("error adding module %v: %v", args[0], err)
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

func (m *Module) loadModule(buf *bytes.Buffer, name, cmdLine string) error {
	log.Printf("Adding module %v", name)

	b, err := readModule(name)
	if err != nil {
		return err
	}

	// place start of each module to a beginning of a page.
	if err := alignUp(buf); err != nil {
		return err
	}

	m.Start = uint32(buf.Len())
	if _, err := buf.Write(b); err != nil {
		return err
	}
	m.End = uint32(buf.Len())

	return nil
}

func (m *Module) setCmdLine(buf *bytes.Buffer, cmdLine string) error {
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

func readGzip(r io.Reader) ([]byte, error) {
	z, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer z.Close()
	return ioutil.ReadAll(z)
}

func readModule(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := readGzip(f)
	if err == nil {
		return b, err
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("cannot rewind file: %v", err)
	}

	return ioutil.ReadAll(f)
}
