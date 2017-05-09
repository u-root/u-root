// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// newc implements the interface for new type cpio files.
package newc

import (
	"fmt"
	"io"
	"reflect"

	"github.com/u-root/u-root/cmds/cpio/pkg"
)

const (
	newcMagic = "070701"
	magicLen  = len(newcMagic)
	headerLen = 13*8 + magicLen
)

type Header struct {
	Ino      uint64
	Mode     uint64
	UID      uint64
	GID      uint64
	Nlink    uint64
	Mtime    uint64
	FileSize uint64
	Major    uint64
	Minor    uint64
	Rmajor   uint64
	Rminor   uint64
	NameSize uint64
	_        uint64
}

func (f *Header) String() string {

	return fmt.Sprintf("Ino %d Mode %#o UID %d GID %d Nlink %d Mtime %#x FileSize %d Major %d Minor %d RMajor %d Rminor %d NameSize %d",
		f.Ino,
		f.Mode,
		f.UID,
		f.GID,
		f.Nlink,
		// what a mess. This fails on travis.
		//time.Unix(int64(f.Mtime), 0).String(),
		f.Mtime,
		f.FileSize,
		f.Major,
		f.Minor,
		f.Rmajor,
		f.Rminor,
		f.NameSize)
}

type newcReader cpio.CpioReader
type newcWriter cpio.CpioWriter

func Header2Info(h *Header, n string, i *cpio.Info) {
	i.Ino = h.Ino
	i.Mode = h.Mode
	i.UID = h.UID
	i.GID = h.GID
	i.Nlink = h.Nlink
	i.Mtime = h.Mtime
	i.FileSize = h.FileSize
	i.Major = h.Major
	i.Minor = h.Minor
	i.Rmajor = h.Rmajor
	i.Rminor = h.Rminor
	i.Name = n
}

// round4 is intended to be use with offsets to WriteAt and ReadAt
func round4(n ...int64) (ret int64) {
	for _, v := range n {
		ret += v
	}

	ret = ((ret + 3) / 4) * 4
	return
}

// Write implements the write interface for newcWriter.
// It allows us to track the position and round it up
// if needed. This allows us to use function such as
// io.copy and Fprintf
func (t *newcWriter) Write(b []byte) (int, error) {
	amt, err := t.Writer.Write(b)
	if err != nil {
		return -1, err
	}
	t.Pos += int64(amt)
	return amt, err
}

func (t *newcWriter) advance(amt int64) error {
	pad := make([]byte, 5)
	o := round4(t.Pos, amt)
	if o == t.Pos {
		return nil
	}
	_, err := t.Write(pad[:o-t.Pos])
	return err
}

func Writer(w io.Writer) (cpio.RecWriter, error) {
	return &newcWriter{Writer: w}, nil
}

func info2Header(i *cpio.Info, h *Header) {

	h.Ino = i.Ino
	h.Mode = i.Mode
	h.UID = i.UID
	h.GID = i.GID
	h.Nlink = i.Nlink
	h.Mtime = i.Mtime
	h.FileSize = i.FileSize
	h.Major = i.Major
	h.Minor = i.Minor
	h.Rmajor = i.Rmajor
	h.Rminor = i.Rminor
	h.NameSize = uint64(len(i.Name) + 1)
}

// Write writes newc cpio records. It pads the header+name write to
// 4 byte alignment and pads the data write as well.
func (t *newcWriter) RecWrite(f *cpio.File) (int, error) {
	var h = &Header{}
	pos := t.Pos

	_, err := t.Write([]byte(newcMagic))
	if err != nil {
		return -1, err
	}

	info2Header(&f.Info, h)
	v := reflect.ValueOf(h)
	for i := 0; i < 13; i++ {
		n := v.Elem().Field(i)
		if _, err := fmt.Fprintf(t, "%08x", n.Uint()); err != nil {
			return -1, err
		}
	}

	if _, err := t.Write([]byte(f.Name)); err != nil {
		return -1, err
	}
	// round to at least one byte past the name.
	if err := t.advance(1); err != nil {
		return -1, err
	}

	if f.Data != nil {
		_, err := io.Copy(t, f.Data)
		if err != nil {
			return -1, err
		}
		if err := t.advance(0); err != nil {
			return -1, err
		}
	}

	return int(t.Pos - pos), nil
}

// Finish does all remaining work to finish an archive.
// It is optional, since sometimes it is desirable to write an archive
// that can be appended to.
func (w *newcWriter) Finish() error {
	if _, err := w.RecWrite(TrailerRecord); err != nil {
		return fmt.Errorf("TestReadWrite: error writing TRAILER: %v", err)
	}
	return nil
}

func Reader(r io.ReaderAt) (cpio.RecReader, error) {
	m := io.NewSectionReader(r, 0, 6)
	var magic [6]byte
	if _, err := m.Read(magic[:]); err != nil {
		return nil, fmt.Errorf("newcReader: unable to read magic: %v", err)
	}
	if string(magic[:]) != newcMagic {
		return nil, fmt.Errorf("newcReader: magic is '%s' and must be '%s'", magic, newcMagic)
	}
	return &newcReader{ReaderAt: r}, nil
}

func (t *newcReader) RecRead() (*cpio.File, error) {
	// There's almost certainly a better way to do this but this
	// will do for now.
	var f = &cpio.File{}
	var h = make([]byte, headerLen)

	cpio.Debug("Next record: pos is %d\n", t.Pos)

	if count, err := t.ReadAt(h[:], t.Pos); count != len(h) || err != nil {
		return nil, fmt.Errorf("Header: at %v got %d of %d bytes, error %v", t.Pos, count, len(h), err)
	}
	t.Pos += int64(len(h))
	// Make sure it's right.
	magic := string(h[:6])
	if magic != newcMagic {
		return nil, fmt.Errorf("Reader: magic '%s' not a newc file", magic)
	}

	cpio.Debug("Header is %v\n", h)
	var hdr Header
	v := reflect.ValueOf(&hdr)
	for i := 0; i < 12; i++ {
		var n uint64
		f := v.Elem().Field(i)
		_, err := fmt.Sscanf(string(h[i*8+6:(i+1)*8+6]), "%x", &n)
		if err != nil {
			return nil, err
		}
		f.SetUint(n)
	}
	cpio.Debug("f is %s\n", (&hdr).String())
	var n = make([]byte, hdr.NameSize)
	if l, err := t.ReadAt(n, t.Pos); l != int(hdr.NameSize) || err != nil {
		return nil, fmt.Errorf("Reading name: got %d of %d bytes, err was %v", l, hdr.NameSize, err)
	}

	// we have to seek to hdr.NameSize + len(h) rounded up to a multiple of 4.
	t.Pos = int64(round4(t.Pos, int64(hdr.NameSize)))

	name := string(n[:hdr.NameSize-1])
	if name == "TRAILER!!!" {
		cpio.Debug("AT THE TRAILER!!!\n")
		return nil, io.EOF
	}
	Header2Info(&hdr, name, &f.Info)
	f.Data = io.NewSectionReader(t, t.Pos, int64(hdr.FileSize))
	t.Pos = int64(round4(t.Pos + int64(hdr.FileSize)))
	return f, nil
}

func init() {
	cpio.AddMap("newc", Reader, Writer)
}
