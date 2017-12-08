// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package find

import (
	"os"
	"path/filepath"
)

type Name struct {
	Name string
	os.FileInfo
	Err error
}

type Finder struct {
	Root     string
	Pattern  string
	Match    func(string, string) (bool, error)
	Mode     os.FileMode
	ModeMask os.FileMode
	Debug    func(string, ...interface{})
	Names    chan *Name
}

type Set func(f *Finder) error

func New(opts ...Set) (*Finder, error) {
	// All of these can be overridden by the opts
	var f = &Finder{Root: "/"}
	f.Debug = func(string, ...interface{}) {}
	f.Names = make(chan *Name, 128)
	f.Match = filepath.Match
	for _, opt := range opts {
		if err := opt(f); err != nil {
			return nil, err
		}
	}
	f.Debug("Create new Finder: %v", f)
	return f, nil
}

func (f *Finder) Find() {
	filepath.Walk(f.Root, func(n string, fi os.FileInfo, err error) error {
		if err != nil {
			f.Names <- &Name{Name: n, Err: err}
			return nil
		}
		// If it matches, then push its name into the result channel,
		// and keep looking.
		f.Debug("Check Pattern '%q' against name '%q'", f.Pattern, fi.Name())
		if f.Pattern != "" {
			m, err := f.Match(f.Pattern, fi.Name())
			if err != nil {
				f.Debug("%s: err on matching: %v", fi.Name(), err)
				return nil
			}
			if !m {
				f.Debug("%s: name does not match", fi.Name())
				return nil
			}
		}
		m := fi.Mode()
		f.Debug("%s fi.Mode %v f.ModeMask %v f.Mode %v", n, m, f.ModeMask, fi.Mode)
		if (m & f.ModeMask) != f.Mode {
			f.Debug("%s: Mode does not match", n)
			return nil
		}
		f.Debug("Found: %v", n)
		f.Names <- &Name{Name: n, FileInfo: fi}
		return nil
	})
	close(f.Names)
}
