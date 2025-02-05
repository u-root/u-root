// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tamago

package bzimage

import (
	"bytes"
	"debug/elf"
	"fmt"
	"io"
	"os"

	"github.com/u-root/u-root/pkg/cpio"
)

// AddInitRAMFS adds an initramfs to the BzImage.
func (b *BzImage) AddInitRAMFS(name string) error {
	u, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	// Should we ever want to compress the initramfs this is one
	// way to do it.
	d := u
	if false {
		d, err = compress(u, "--lzma2=,dict=1MiB")
		if err != nil {
			return err
		}
	}
	s, e, err := b.InitRAMFS()
	if err != nil {
		return err
	}
	l := e - s

	if len(d) > l {
		return fmt.Errorf("new initramfs is %d bytes, won't fit in %d byte old one", len(d), l)
	}
	// Do this in a stupid way that is easy to read.
	// What's interesting: the kernel decompressor, if I read it right,
	// finds it easier to skip a bunch of leading nulls. So do that.
	n := make([]byte, l)
	Debug("Offset into n is %d\n", len(n)-len(d))
	copy(n[len(n)-len(d):], d)
	Debug("Install %d byte initramfs in %d bytes of kernel code, @ %d:%d", len(d), len(n), s, e)
	copy(b.KernelCode[s:e], n)
	return nil
}

// InitRAMFS returns a []byte from KernelCode which can be used to save or replace
// an existing InitRAMFS. The fun part is that there are no symbols; what we do instead
// is find the programs what are RW and look for the cpio magic in them. If we find it,
// we see if it can be read as a cpio and, if so, if there is a /dev or /init inside.
// We repeat until we succeed or there's nothing left.
func (b *BzImage) InitRAMFS() (int, int, error) {
	f, err := b.ELF()
	if err != nil {
		return -1, -1, err
	}
	// Find the program header with RWE.
	var dat []byte
	var prog *elf.Prog
	for _, p := range f.Progs {
		if p.Flags&(elf.PF_X|elf.PF_W|elf.PF_R) == elf.PF_X|elf.PF_W|elf.PF_R {
			dat, err = io.ReadAll(p.Open())
			if err != nil {
				return -1, -1, err
			}
			prog = p
			break
		}
	}
	if dat == nil {
		return -1, -1, fmt.Errorf("can't find an RWE prog in kernel")
	}

	archiver, err := cpio.Format("newc")
	if err != nil {
		return -1, -1, fmt.Errorf("format newc not supported: %w", err)
	}
	var cur int
	for cur < len(dat) {
		x := bytes.Index(dat, []byte("070701"))
		if x == -1 {
			return -1, -1, fmt.Errorf("no newc cpio magic found")
		}
		if err != nil {
			return -1, -1, err
		}
		cur = x
		r := bytes.NewReader(dat[cur:])
		rr := archiver.Reader(r)
		Debug("r.Len is %v", r.Len())
		var found bool
		var size int
		for {
			rec, err := rr.ReadRecord()
			Debug("Check %v", rec)
			if err == io.EOF {
				break
			}
			if err != nil {
				Debug("error reading records: %v", err)
				break
			}
			switch rec.Name {
			case "init", "dev", "bin", "usr":
				Debug("Found initramfs at %d, %d bytes", cur, len(dat)-r.Len())
				found = true
			}
			size = int(rec.FilePos) + int(rec.FileSize)
		}
		Debug("Size is %d", size)
		// Add the trailer size.
		y := x + size
		if found {
			// The slice consists of the bytes for cur to the length of initramfs.
			// We can derive the initramfs length by knowing how much is left of the reader.
			Debug("Return %d %#x slice %d:%d from %d byte dat", len(dat[x:y]), len(dat[x:y]), cur, y, len(dat))
			x += int(prog.Off)
			y += int(prog.Off)
			// We need to round y up to the end of the record. We have to do this after we
			// add the prog.Off value to it.
			y = (y + 3) &^ 3
			// and add the size of the trailer record.
			y += 120
			y += 4 // and add at least one word of null
			y = (y + 3) &^ 3
			Debug("InitRAMFS: return %d, %d", x, y)
			return x, y, nil
		}
		cur += 6
	}
	return -1, -1, fmt.Errorf("no cpio found")
}
