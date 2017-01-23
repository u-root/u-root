// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
)

func create(f *File) error {

	m, err := cpioModetoMode(f.Mode)
	if err != nil {
		return err
	}

	switch m {
	case os.ModeSocket:
		return fmt.Errorf("%v: type %v: not yet", f.Name, m)
	case os.FileMode(0):
		nf, err := os.Create(f.Name)
		if err != nil {
			return err
		}
		_, err = io.Copy(nf, f.Data)
		if err != nil {
			return err
		}
		if err = nf.Chmod(os.FileMode(f.Mode)); err != nil {
			return err
		}
		return nil
	case os.ModeDevice:
		return fmt.Errorf("%v: type %v: not yet", f.Name, m)
	case os.ModeDir:
		err = os.MkdirAll(f.Name, os.FileMode(f.Mode))
		return err
	case os.ModeCharDevice:
		return fmt.Errorf("%v: type %v: not yet", f.Name, m)
	case os.ModeNamedPipe:
		return fmt.Errorf("%v: type %v: not yet", f.Name, m)
	default:
		return fmt.Errorf("WUT?")
	}

}
