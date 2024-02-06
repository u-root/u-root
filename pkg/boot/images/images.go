// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package images

import (
	"io"

	"github.com/u-root/u-root/pkg/boot/linux"
	"github.com/u-root/u-root/pkg/boot/multiboot"
)

type Creator struct {
	LinuxOpts     []linux.Opt
	MultibootOpts []multiboot.Opt
}

type Opt func(*Creator)

func WithLinuxOpt(lo ...linux.Opt) Opt {
	return func(c *Creator) {
		c.LinuxOpts = append(c.LinuxOpts, lo...)
	}
}

func WithMultibootOpt(mo ...multiboot.Opt) Opt {
	return func(c *Creator) {
		c.MultibootOpts = append(c.MultibootOpts, mo...)
	}
}

func NewCreator(opts ...Opt) *Creator {
	var c Creator
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

func (c *Creator) NewLinuxImage(kernel io.ReaderAt, lo ...linux.Opt) *linux.Image {
	return linux.NewImage(kernel, append(c.LinuxOpts, lo...)...)
}

func (c *Creator) NewMultibootImage(kernel io.ReaderAt, lo ...multiboot.Opt) *multiboot.Image {
	return multiboot.NewImage(kernel, append(c.MultibootOpts, lo...)...)
}
