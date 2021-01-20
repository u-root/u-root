// The MIT License (MIT)
//
// Copyright (c) 2018 The Genuinetools Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package container

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/archive"
)

const (
	// DefaultTarballPath holds the default path for the embedded tarball.
	DefaultTarballPath = "image.tar"
)

// UnpackRootfs unpacks the embedded tarball to the rootfs.
func (c *Container) UnpackRootfs(rootfsDir string, asset func(string) ([]byte, error)) error {
	// Make the rootfs directory.
	if err := os.MkdirAll(rootfsDir, 0755); err != nil {
		return err
	}

	// Get the embedded tarball.
	data, err := asset(DefaultTarballPath)
	if err != nil {
		return fmt.Errorf("getting bindata asset image.tar failed: %v", err)
	}

	// Unpack the tarball.
	r := bytes.NewReader(data)
	if err := archive.Untar(r, rootfsDir, &archive.TarOptions{NoLchown: true}); err != nil {
		return err
	}

	// Write a resolv.conf.
	if err := ioutil.WriteFile(filepath.Join(rootfsDir, "etc", "resolv.conf"), []byte("nameserver 8.8.8.8\nnameserver 8.8.4.4"), 0755); err != nil {
		return err
	}

	return nil
}
