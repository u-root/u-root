// Copyright 2014-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package libinit creates the environment and root file system for u-root.
package libinit

import (
	"fmt"
	"os"

	"github.com/u-root/u-root/pkg/ulog"
)

type creator interface {
	create() error
	fmt.Stringer
}

type dir struct {
	Name string
	Mode os.FileMode
}

func (d dir) create() error {
	return os.MkdirAll(d.Name, d.Mode)
}

func (d dir) String() string {
	return fmt.Sprintf("dir %q (mode %#o)", d.Name, d.Mode)
}

type mount struct {
	Source string
	Flag   uint
	Target string
}

func (m mount) create() error {
	return fmt.Errorf("Not yet")
}

func (m mount) String() string {
	return fmt.Sprintf("mount source %q target %q flags %#x", m.Source, m.Target, m.Flag)
}

var (
	// These have to be created / mounted first, so that the logging works correctly.
	preNamespace = []creator{}
	namespace    = []creator{}
)

func goBin() string {
	return "/bin"
}

func create(namespace []creator, optional bool) {
	for _, c := range namespace {
		if err := c.create(); err != nil {
			if optional {
				ulog.KernelLog.Printf("u-root init [optional]: warning creating %s: %v", c, err)
			} else {
				ulog.KernelLog.Printf("u-root init: error creating %s: %v", c, err)
			}
		}
	}
}

// SetEnv sets the default u-root environment.
func SetEnv() {
	env := map[string]string{
		"LD_LIBRARY_PATH": "/usr/local/lib",
		"GOROOT":          "/go",
		"GOPATH":          "/",
		"GOBIN":           "/ubin",
		"CGO_ENABLED":     "0",
	}

	// Not all these paths may be populated or even exist but OTOH they might.
	path := "/ubin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/bin:/usr/local/sbin:/buildbin:/bbin"

	env["PATH"] = fmt.Sprintf("%v:%v", goBin(), path)
	for k, v := range env {
		os.Setenv(k, v)
	}
}

// CreateRootfs creates the default u-root file system.
func CreateRootfs() {
}
