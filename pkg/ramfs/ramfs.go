// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ramfs

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/u-root/u-root/pkg/cpio"
	_ "github.com/u-root/u-root/pkg/cpio/newc"
)

const (
	d = syscall.S_IFDIR
	c = syscall.S_IFCHR
	b = syscall.S_IFBLK
	f = syscall.S_IFREG

	// This is the literal timezone file for GMT-0. Given that we have no idea
	// where we will be running, GMT seems a reasonable guess. If it matters,
	// setup code should download and change this to something else.
	gmt0       = "TZif2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x04\x00\x00\x00\x00\x00\x00GMT\x00\x00\x00TZif2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x04\xf8\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00GMT\x00\x00\x00\nGMT0\n"
	nameserver = "nameserver 8.8.8.8\n"
)

// devCPIOrecords are cpio records as defined in the uroot cpio package.
// Most of the bits can be left unspecified: these all have one link,
// they are mostly root:root, for example.
var DevCPIO = []cpio.Record{
	{Info: cpio.Info{Name: "tcz", Mode: d | 0755}},
	{Info: cpio.Info{Name: "etc", Mode: d | 0755}},
	{Info: cpio.Info{Name: "dev", Mode: d | 0755}},
	{Info: cpio.Info{Name: "ubin", Mode: d | 0755}},
	{Info: cpio.Info{Name: "usr", Mode: d | 0755}},
	{Info: cpio.Info{Name: "usr/lib", Mode: d | 0755}},
	{Info: cpio.Info{Name: "lib64", Mode: d | 0755}},
	{Info: cpio.Info{Name: "bin", Mode: d | 0755}},
	{Info: cpio.Info{Name: "dev/console", Mode: c | 0600, Rmajor: 5, Rminor: 1}},
	{Info: cpio.Info{Name: "dev/tty", Mode: c | 0666, Rmajor: 5, Rminor: 0}},
	{Info: cpio.Info{Name: "dev/null", Mode: c | 0666, Rmajor: 1, Rminor: 3}},
	{Info: cpio.Info{Name: "dev/port", Mode: c | 0640, Rmajor: 1, Rminor: 4}},
	{Info: cpio.Info{Name: "dev/urandom", Mode: c | 0666, Rmajor: 1, Rminor: 9}},
	{Info: cpio.Info{Name: "etc/resolv.conf", Mode: f | 0644, FileSize: uint64(len(nameserver))}, ReadCloser: cpio.NewBytesReadCloser([]byte(nameserver))},
	{Info: cpio.Info{Name: "etc/localtime", Mode: f | 0644, FileSize: uint64(len(gmt0))}, ReadCloser: cpio.NewBytesReadCloser([]byte(gmt0))},
}

type Initramfs struct {
	cpio.Writer
	files map[string]struct{}
}

func NewInitramfs(w cpio.Writer) (*Initramfs, error) {
	// Write devtmpfs records.
	dcpio := DevCPIO[:]
	cpio.MakeAllReproducible(dcpio)
	if err := w.WriteRecords(dcpio); err != nil {
		return nil, err
	}

	return &Initramfs{
		Writer: w,
		files:  make(map[string]struct{}),
	}, nil
}

func (i *Initramfs) WriteRecord(r cpio.Record) error {
	if r.Name == "." || r.Name == "/" {
		return nil
	}

	// Create record for parent directory if needed.
	dir := filepath.Dir(r.Name)
	if _, ok := i.files[dir]; dir != "/" && dir != "." && !ok {
		if err := i.WriteRecord(cpio.Record{
			Info: cpio.Info{
				Name: dir,
				Mode: syscall.S_IFDIR | 0755,
			},
		}); err != nil {
			return err
		}
	}

	i.files[r.Name] = struct{}{}
	return i.Writer.WriteRecord(r)
}

func (i *Initramfs) WriteFile(src string, dest string) error {
	record, err := cpio.GetRecord(src)
	if err != nil {
		return err
	}

	if record.Info.Mode&^0777 == syscall.S_IFDIR {
		return children(src, func(name string) error {
			return i.WriteFile(filepath.Join(src, name), filepath.Join(dest, name))
		})
	} else {
		// Fix the name.
		record.Name = dest
		return i.WriteRecord(cpio.MakeReproducible(record))
	}
}

func children(dir string, fn func(name string) error) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return err
	}

	for _, name := range names {
		if err := fn(name); os.IsNotExist(err) {
			// File was deleted in the meantime.
			continue
		} else if err != nil {
			return err
		}
	}
	return nil
}

// Copy all files relative to `srcDir` to `destDir` in the cpio archive.
func (i *Initramfs) WriteFiles(srcDir string, destDir string, files []string) error {
	for _, file := range files {
		srcPath := filepath.Join(srcDir, file)
		destPath := filepath.Join(destDir, file)
		if err := i.WriteFile(srcPath, destPath); err != nil {
			return err
		}
	}
	return nil
}
