// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package namespace

import (
	"errors"
	"fmt"
	"os"
)

type cmd struct {
	syscall syzcall
	flag    mountflag

	filename string
	line     string

	args []string
}

// valid returns usages for these commands
func (c cmd) valid() error {
	args := c.args
	switch c.syscall {
	case BIND:
		if len(args) < 2 {
			return errors.New("usage: bind [–abcC] new old")
		}
	case MOUNT:
		if len(args) < 2 {
			return errors.New("usage: mount [–abcC] servename old [spec]")
		}
	case UNMOUNT:
		if len(args) < 1 {
			return errors.New("usage: unmount [ new ] old")
		}
	case RFORK:
		// doesn't take args or flags, so always valid even if not.
		return nil
	case CHDIR:
		if len(args) < 1 {
			return errors.New("usage: cd dir")
		}
	case IMPORT:
		if len(args) < 2 {
			return errors.New("usage: import [–abc] host [remotepath] mountpoint")
		}
	case INCLUDE:
		if len(args) < 1 {
			return errors.New("usage: . path")
		}
	default:
		return fmt.Errorf("%d is not implmented", c.syscall)
	}
	return nil
}

func (c cmd) Modify(ns Namespace, b *Builder) error {
	args := []string{}
	for _, arg := range c.args {
		args = append(args, os.ExpandEnv(arg))
	}
	if err := c.valid(); err != nil {
		return err
	}
	switch c.syscall {
	case BIND:
		return ns.Bind(args[0], args[1], c.flag)
	case MOUNT:
		servername := args[0]
		old := args[1]
		spec := ""
		if len(args) == 3 {
			spec = args[2]
		}
		return ns.Mount(servername, old, spec, c.flag)
	case UNMOUNT:
		var newNS, oldNS string
		if len(args) == 2 {
			newNS, oldNS = args[0], args[1]
		} else {
			oldNS = args[0]
		}
		return ns.Unmount(newNS, oldNS)
	case RFORK:
		return ns.Clear()
	case CHDIR:
		if err := ns.Chdir(args[0]); err != nil {
			return err
		}
		b.dir = args[0]
		return nil
	case IMPORT:
		host := ""
		remotepath := ""
		mountpoint := ""
		if len(args) == 2 {
			host = args[0]
			mountpoint = args[1]
		} else if len(args) == 3 {
			host = args[0]
			remotepath = args[1]
			mountpoint = args[2]
		}
		return ns.Import(host, remotepath, mountpoint, c.flag)
	case INCLUDE:
		var nb *Builder
		nb, err := newBuilder(b.dir, b.open)
		if err != nil {
			return err
		}
		if err := nb.Parse(args[0]); err != nil {
			return err
		}
		if err := nb.buildNS(ns); err != nil {
			return err
		}
		b.dir = nb.dir // if the new file has changed the directory we'd like to know
		return nil
	default:
		return fmt.Errorf("%s not implmented", c.syscall)
	}
}

func (c cmd) String() string { return fmt.Sprintf("%s(%v, %d)", c.syscall, c.args, c.flag) }
