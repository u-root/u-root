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
// The name onehot may be a misnomer but we're keeping open
// the possibility of adding state such that only one onehot
// file can be open at a time. We will see if we need that.
// Since our first use case is cpio, that seems unlikely.
package onehot

import (
	"io"
	"os"
	"sync"
)

type OneHot struct {
	m sync.Mutex
	name string
	f    *os.File
	dead bool
	err  error
}

// Open returns a new OneHot with the name filled in.
func Open(name string) (io.ReadCloser, error) {
	return &OneHot{name: name}, nil
}

// Read reads from a OneHot. If it is dead, it returns
// -1 and the last error, which is sticky. If it is
// not open, it opens it; if the open fails, it marks
// the OneHot as dead and records the error.
// If the open works, it calls the Read on the
// underlying Reader; if there is an error there
// the OneHot is marked dead and the error is recorded.
func (f *OneHot) Read(b []byte) (int, error) {
	var err error
	f.m.Lock()
	defer f.m.Unlock()
	if f.dead {
		return -1, f.err
	}
	if f.f == nil {
		f.f, err = os.Open(f.name)
		if err != nil {
			f.dead = true
			f.err = err
			return -1, err
		}
	}
	n, err := f.f.Read(b)

	if n < 0 || err != nil {
		f.dead = true
		f.err = err
	}
	return n, err
}

// Close closes the OneNot. It always calls f.f.Close()
// then sets f.f to nil, retains the error if any, and
// makes the OneHot as dead.
func (f *OneHot) Close() error {
	f.m.Lock()
	defer f.m.Unlock()
	// This is a bit tricky. If f.f is nil
	// but it's not dead, it has not been opened.
	// Return no error, but mark it dead.
	if ! f.dead && f.f == nil {
		f.dead = true
		f.err = nil
		return nil
	}
	// calling Close on nil Files is OK.
	// This way we are sure to get the right error
	err := f.f.Close()
	f.dead = true
	f.err = err
	f.f = nil
	return err
}
