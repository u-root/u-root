// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ramfs

import (
	"fmt"
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

	Path string

	files map[string]struct{}
}

func NewInitramfs(goos string, goarch string) (*Initramfs, error) {
	oname := fmt.Sprintf("/tmp/initramfs.%v_%v.cpio", goos, goarch)
	f, err := os.Create(oname)
	if err != nil {
		return nil, err
	}

	archiver, err := cpio.Format("newc")
	if err != nil {
		return nil, err
	}

	w := archiver.Writer(f)

	// Write devtmpfs records.
	dcpio := DevCPIO[:]
	cpio.MakeAllReproducible(dcpio)
	if err := w.WriteRecords(dcpio); err != nil {
		return nil, err
	}

	return &Initramfs{
		Path:   oname,
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

func (i *Initramfs) WriteFile(src string, dest string, path string) error {
	name, err := filepath.Rel(src, path)
	if err != nil {
		return fmt.Errorf("path %q not relative to src %q: %v", path, src, err)
	}

	record, err := cpio.GetRecord(path)
	if err != nil {
		return err
	}

	// Fix the name.
	record.Name = filepath.Join(dest, name)
	return i.WriteRecord(cpio.MakeReproducible(record))
}

// Walk walks to all files in a tree rooted at `dir`.
//
// filepath.Walk doesn't walk to symlinks, so we can't use it.
func walk(dir string, fn func(path string, fi os.FileInfo) error) error {
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
		path := filepath.Join(dir, name)
		fileInfo, err := os.Lstat(path)
		// File has been deleted, or something.
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return fmt.Errorf("lstat(%q) = %v", path, err)
		}
		if err := fn(path, fileInfo); err != nil {
			return err
		}

		// Visit children of directory.
		if fileInfo.IsDir() {
			if err := walk(path, fn); err != nil {
				return err
			}
		}
	}
	return nil
}

// Copy all files relative to `srcDir` to `destDir` in the cpio archive.
func (i *Initramfs) CopyDir(srcDir string, destDir string) error {
	return walk(srcDir, func(path string, _ os.FileInfo) error {
		return i.WriteFile(srcDir, destDir, path)
	})
}

// Copy all files relative to `srcDir` to `destDir` in the cpio archive.
func (i *Initramfs) WriteFiles(srcDir string, destDir string, files []string) error {
	for _, file := range files {
		path := filepath.Join(srcDir, file)
		fi, err := os.Stat(path)
		if err != nil {
			return err
		}

		switch fi.Mode() &^ 0777 {
		case os.ModeDir:
			dest := filepath.Join(destDir, file)
			// Copy all files in directory.
			if err := i.CopyDir(path, dest); err != nil {
				return err
			}

		default:
			if err := i.WriteFile(srcDir, destDir, path); err != nil {
				return err
			}
		}
	}
	return nil
}
