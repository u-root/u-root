// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/kexec"
	"github.com/u-root/u-root/pkg/uio"
)

// LinuxImage implements OSImage for a Linux kernel + initramfs.
type LinuxImage struct {
	Kernel  io.ReaderAt
	Initrd  io.ReaderAt
	Cmdline string
}

var _ OSImage = &LinuxImage{}

// String prints a human-readable version of this linux image.
func (li *LinuxImage) String() string {
	return fmt.Sprintf("LinuxImage(\n  Kernel: %s\n  Initrd: %s\n  Cmdline: %s\n)\n", li.Kernel, li.Initrd, li.Cmdline)
}

func copyToFile(r io.Reader) (*os.File, error) {
	f, err := ioutil.TempFile("", "nerf-netboot")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return nil, err
	}
	if err := f.Sync(); err != nil {
		return nil, err
	}

	readOnlyF, err := os.Open(f.Name())
	if err != nil {
		return nil, err
	}
	return readOnlyF, nil
}

// Load implements OSImage.Load and kexec_load's the kernel with its initramfs.
func (li *LinuxImage) Load(verbose bool) error {
	if li.Kernel == nil {
		return errors.New("LinuxImage.Kernel must be non-nil")
	}

	k, err := copyToFile(uio.Reader(li.Kernel))
	if err != nil {
		return err
	}
	defer k.Close()

	var i *os.File
	if li.Initrd != nil {
		i, err = copyToFile(uio.Reader(li.Initrd))
		if err != nil {
			return err
		}
		defer i.Close()
	}

	log.Printf("Kernel: %s", k.Name())
	if i != nil {
		log.Printf("Initrd: %s", i.Name())
	}
	log.Printf("Command line: %s", li.Cmdline)
	return kexec.FileLoad(k, i, li.Cmdline)
}
