// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mkfifo

import "syscall"

type Mkfifo struct {
	Paths []string
	Mode  uint32
}

func (m *Mkfifo) Exec() error {
	var err error

	for _, path := range m.Paths {
		if err := syscall.Mkfifo(path, m.Mode); err != nil {
			return err
		}
	}
	return err
}
