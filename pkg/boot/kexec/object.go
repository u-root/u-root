// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"debug/elf"
	"debug/plan9obj"
	"fmt"
	"io"

	"github.com/u-root/u-root/pkg/align"
)

// Object is an object file, specific to kexec uses.
// It is used to get Progs and Entry from an object.
// elf.Prog is generic enough that we can use it to define
// loadable segments for non-ELF objects, such as Plan 9 a.out
type Object interface {
	Progs() []*elf.Prog
	Entry() uint64
}

type elfObject struct {
	f *elf.File
}

func (e *elfObject) Progs() []*elf.Prog {
	return e.f.Progs
}

func (e *elfObject) Entry() uint64 {
	return e.f.Entry
}

var _ Object = &elfObject{}

type aout9Object struct {
	f *plan9obj.File
}

func (e *aout9Object) Progs() []*elf.Prog {
	// It's easier in Plan 9; two Progs, that's it.
	// And if there are not two, well, it's not a kernel,
	// so we can make some assumptions.

	/*
		type FileHeader struct {
			Magic       uint32
			Bss         uint32
			Entry       uint64
			PtrSize     int
			LoadAddress uint64
			HdrSize     uint64
		}
	*/
	return []*elf.Prog{
		{
			ProgHeader: elf.ProgHeader{
				Type:   elf.PT_LOAD,
				Flags:  elf.PF_X | elf.PF_R,
				Filesz: uint64(e.f.Sections[0].Size),

				Paddr: e.f.LoadAddress,
			},
			ReaderAt: e.f.Sections[0].ReaderAt,
		},
		{
			ProgHeader: elf.ProgHeader{
				Type:   elf.PT_LOAD,
				Flags:  elf.PF_W | elf.PF_R,
				Filesz: uint64(e.f.Sections[1].Size),
				Paddr:  uint64(align.UpPage(uint(e.f.LoadAddress + uint64(e.f.Sections[0].Size)))),
			},
			ReaderAt: e.f.Sections[1].ReaderAt,
		},
	}
}

func (e *aout9Object) Entry() uint64 {
	return e.f.Entry
}

var _ Object = &aout9Object{}

// ObjectNewFile reads from input stream and returns a specific object. It tries
// elf first, then plan9.
func ObjectNewFile(r io.ReaderAt) (Object, error) {
	f, errELF := elf.NewFile(r)
	if errELF == nil {
		return &elfObject{f: f}, nil
	}
	f9, err9 := plan9obj.NewFile(r)
	if err9 == nil {
		return &aout9Object{f: f9}, nil
	}
	return nil, fmt.Errorf("ELF: %w, plan9obj: %w", errELF, err9)
}
