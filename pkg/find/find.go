// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package find

import (
	"os"
	"path/filepath"
	"regexp"
)

type Finder struct {
	Root     string
	Match    *regexp.Regexp
	Mode     os.FileMode
	ModeMask os.FileMode
	Name     chan string
	Err      chan error
}

type Set func(f *Finder) error

func New(rules ...Set) (*Finder, error) {
	var s = &Finder{Root: "/", Match: regexp.MustCompile(".*")}
	for _, r := range rules {
		if err := r(s); err != nil {
			return nil, err
		}
	}
	s.Name = make(chan string, 128)
	s.Err = make(chan error, 128)
	return s, nil
}

func (f *Finder) Find() {
	filepath.Walk(f.Root, func(n string, fi os.FileInfo, err error) error {
		if err != nil {
			f.Err <- err
			return err
		}
		// If it matches, then push its name into the result channel,
		// and keep looking.
		if !f.Match.Match([]byte(n)) {
			return nil
		}
		if (fi.Mode() & f.ModeMask) != f.Mode {
			return nil
		}
		f.Name <- n
		return nil
	})
}
