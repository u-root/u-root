// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package namespace

import (
	"fmt"
	"io"
	"os"
)

//go:generate stringer -type mountflag
type mountflag int

//go:generate stringer -type syzcall
type syzcall byte

// File are a collection of namespace modifiers
type File []Modifier

// Namespace is a plan9 namespace. It implmenets the
// http://man.cat-v.org/plan_9/2/bind calls.
//
// Bind and mount modify the file name space of the current
// process and other processes in its name space group
// (see http://man.cat-v.org/plan_9/2/fork).
// For both calls, old is the name of an existing file or directory
// in the current name space where the modification is to be made.
// The name old is evaluated as described in
// http://man.cat-v.org/plan_9/2/intro,
// except that no translation of the final path element is done.
type Namespace interface {
	// Bind binds new on old.
	Bind(name, old string, flag mountflag) error
	// Mount mounts servename on old.
	Mount(servername, old, spec string, flag mountflag) error
	// Unmount unmounts new from old, or everything mounted on old if new is missing.
	Unmount(name, old string) error
	// Clear clears the name space with rfork(RFCNAMEG).
	Clear() error
	// Chdir changes the working directory to dir.
	Chdir(dir string) error
	// Import imports a name space from a remote system
	Import(host, remotepath, mountpoint string, flag mountflag) error
}

// Modifier repesents an individual command that can be applied
// to a plan9 name space which will modify the name space of the process or process group.
type Modifier interface {
	// Modify modifies the namespace
	Modify(ns Namespace, b *Builder) error
	String() string
}

// NewNS builds a name space for user.
// It opens the file nsfile (/lib/namespace is used if nsfile is ""),
// copies the old environment, erases the current name space,
// sets the environment variables user and home, and interprets the commands in nsfile.
// The format of nsfile is described in namespace(6).
func NewNS(nsfile string, user string) error { return buildNS(nil, nsfile, user, true) }

// AddNS also interprets and executes the commands in nsfile.
// Unlike newns it applies the command to the current name
// space rather than starting from scratch.
func AddNS(nsfile string, user string) error { return buildNS(nil, nsfile, user, false) }

func buildNS(ns Namespace, nsfile, user string, clearns bool) error {
	if err := os.Setenv("user", user); err != nil {
		return err
	}
	if ns == nil {
		ns = DefaultNamespace
	}
	if clearns {
		ns.Clear()
	}
	r, err := NewBuilder()
	if err != nil {
		return err
	}
	if err := r.Parse(nsfile); err != nil {
		return err
	}
	return r.buildNS(ns)
}

// OpenFunc opens files for the include or . commands in name space files.
// while the default open function is just os.Open for the file://
// one could extend this to other protocols. Potentially we could
// even import files from the web.eg
//
//	. https://harvey-os.org/lib/namespace@sha256:deadbeef
//
// could be interesting.
type OpenFunc func(path string) (io.Reader, error)

// Builder helps building plan9 name spaces. Builder keeps track of directory changes
// when another name space file is included, another builder will be created for it's
// modifications, and it's final working directory will be set to be the parents working
// directroy after it's modifications are complete.
type Builder struct {
	dir  string
	file File

	open OpenFunc
}

func open1(path string) (io.Reader, error) { return os.Open(path) }

// NewBuilder returns a builder with defaults
func NewBuilder() (*Builder, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &Builder{
		dir:  wd,
		open: open1,
	}, nil
}

func newBuilder(wd string, b OpenFunc) (*Builder, error) {
	return &Builder{
		dir:  wd,
		open: b,
	}, nil
}

// Parse takes a path and parses the namespace file
func (b *Builder) Parse(file string) error {
	f, err := b.open(file)
	if err != nil {
		return err
	}
	b.file, err = Parse(f)
	if err != nil {
		return err
	}
	return nil
}

// Run takes a namespace and runs commands defined in the namespace file
func (b *Builder) buildNS(ns Namespace) error {
	for _, c := range b.file {
		if err := c.Modify(ns, b); err != nil {
			return fmt.Errorf("newns failed to perform %s failed: %w", c, err)
		}
	}
	return nil
}
