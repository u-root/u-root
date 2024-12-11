// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// mount mounts a filesystem at the specified path.
//
// Synopsis:
//
//	mount [-r] [-o options] [-t FSTYPE] DEV PATH
//
// Options:
//
//	-r: read only
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/loop"
	"golang.org/x/sys/unix"
)

var (
	errUsage     = errors.New("usage")
	errMountPath = errors.New("can not read mount path")
)

type mountOptions []string

func (o *mountOptions) String() string {
	return strings.Join(*o, ",")
}

func (o *mountOptions) Set(value string) error {
	for _, option := range strings.Split(value, ",") {
		*o = append(*o, option)
	}
	return nil
}

type cmd struct {
	stdout          io.Writer
	stderr          io.Writer
	fileSystemsPath string
	fsType          string
	mountsPath      []string
	options         mountOptions
	ro              bool
}

func command(stdout, stderr io.Writer, ro bool, fsType string, opts mountOptions) *cmd {
	return &cmd{
		stdout: stdout,
		stderr: stderr,
		mountsPath: []string{
			"/proc/self/mounts",
			"/proc/mounts",
			"/etc/mtab",
		},
		fileSystemsPath: "/proc/filesystems",
		ro:              ro,
		options:         opts,
		fsType:          fsType,
	}
}

func loopSetup(filename string) (loopDevice string, err error) {
	loopDevice, err = loop.FindDevice()
	if err != nil {
		return "", err
	}
	if err := loop.SetFile(loopDevice, filename); err != nil {
		return "", err
	}
	return loopDevice, nil
}

// extended from boot.go
func (c *cmd) getSupportedFilesystem(originFS string) ([]string, bool, error) {
	var known bool
	var err error
	fs, err := os.ReadFile(c.fileSystemsPath)
	if err != nil {
		return nil, known, err
	}
	var returnValue []string
	for _, f := range strings.Split(string(fs), "\n") {
		n := strings.Fields(f)
		last := len(n) - 1
		if last < 0 {
			continue
		}
		if n[last] == originFS {
			known = true
		}
		returnValue = append(returnValue, n[last])
	}
	return returnValue, known, err
}

func (c *cmd) informIfUnknownFS(originFS string) {
	knownFS, known, err := c.getSupportedFilesystem(originFS)
	if err != nil {
		// just don't make things even worse...
		return
	}
	if !known {
		fmt.Fprintf(c.stderr, "Hint: unknown filesystem %s. Known are: %v", originFS, knownFS)
	}
}

func (c *cmd) run(args ...string) error {
	if len(args) == 0 {
		for _, p := range c.mountsPath {
			if b, err := os.ReadFile(p); err == nil {
				c.stdout.Write(b)
				return nil
			}
		}
		return fmt.Errorf("%w: %v to get namespace", errMountPath, c.mountsPath)
	}

	if len(args) < 2 {
		return errUsage
	}

	dev := args[0]
	path := args[1]
	var flags uintptr
	var data []string
	var err error
	for _, option := range c.options {
		switch option {
		case "loop":
			dev, err = loopSetup(dev)
			if err != nil {
				return fmt.Errorf("error setting loop device: %w", err)
			}
		default:
			if f, ok := opts[option]; ok {
				flags |= f
			} else {
				data = append(data, option)
			}
		}
	}

	if c.ro {
		flags |= unix.MS_RDONLY
	}
	if c.fsType == "" {
		if _, err := mount.TryMount(dev, path, strings.Join(data, ","), flags); err != nil {
			return err
		}
	} else {
		if _, err := mount.Mount(dev, path, c.fsType, strings.Join(data, ","), flags); err != nil {
			c.informIfUnknownFS(c.fsType)
			return err
		}
	}

	return nil
}

func main() {
	ro := flag.Bool("r", false, "Read only mount")
	fsType := flag.String("t", "", "File system type")
	var options mountOptions
	flag.Var(&options, "o", "Comma separated list of mount options")
	flag.Parse()
	cmd := command(os.Stdout, os.Stderr, *ro, *fsType, options)

	err := cmd.run(flag.Args()...)
	if errors.Is(err, errUsage) {
		flag.Usage()
		os.Exit(1)
	} else if err != nil {
		log.Fatal(err)
	}
}
