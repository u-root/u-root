// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package onehot implements a "one hot" struct that implements
// Reader but defers the open until the first read. It should be
// used in packages, like cpio, that want to set up lots of
// Readers and read them later. Just reading the files in
// is impractical for many reasons:
// o they might be named pipes on Plan9
// o they might be so large they won't fit in memory
// o there might be so many of them they might not all
//   fit in memory.
// But just opening them and returning an io.Reader is also
// impractical as, seen in practice, we might get EMFILE.
// Hence, we keep track of the file, but don't open it.
// We considered doing a test open but it seems a bit pointless;
// the file might go away by the time we get around to opening
// it for I/O and, on some systems, opening and then closing
// the file can make it go away, even if we don't read it.
package onehot

import (
	"io"
	"os"
	"sync"
)

type OneHot struct {
	Name   string
	Closed bool
	*os.File
}

var (
	m sync.Mutex
	f *OneHot
)

// Open returns a new OneHot with the name filled in
func Open(name string) (io.ReadCloser, error) {
	return &OneHot{Name: name}, nil
}

// Read reads from the file in the OneHot.
// If the File is not the same as our current file,
// we close the current file and open the OneHot.
func (o *OneHot) Read(b []byte) (int, error) {
	if o.Closed {
		return -1, io.EOF
	}
	if f != nil && f != o {
		f.File.Close()
		f.File = nil
		f = nil
	}
	if f == nil {
		var err error
		o.File, err = os.Open(o.Name)
		if err != nil {
			return -1, err
		}
		f = o
	}
	return o.File.Read(b)
}

// Close closes the OneHot. If the File in the o is not the
// currently used file, we forcibly close it and set o.File
// to nil. Closes on nil are allowed but will get errors.
func (o *OneHot) Close() error {
	m.Lock()
	defer m.Unlock()
	if o.Closed {
		return o.File.Close()
	}
	if o.File == nil {
		o.Closed = true
		return nil
	}
	if o != f {
		err := o.File.Close()
		o.File = nil
		o.Closed = true
		return err
	}
	err := o.File.Close()
	f = nil
	o.File = nil
	o.Closed = true
	return err
}
