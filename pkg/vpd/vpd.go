// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vpd

import (
	"os"
	"path"
	"path/filepath"
)

// default variables
var (
	MaxBootEntry  = 9999
	DefaultVpdDir = "/sys/firmware/vpd"
)

var globalReader = NewReader()

// NewReader returns a new VPD Reader.
func NewReader() *Reader {
	return &Reader{
		VpdDir: DefaultVpdDir,
	}
}

// Get wraps globalReader.Get .
func Get(key string, readOnly bool) ([]byte, error) {
	return globalReader.Get(key, readOnly)
}

// Set wraps globalReader.Set .
func Set(key string, value []byte, readOnly bool) error {
	return globalReader.Set(key, value, readOnly)
}

// GetAll wraps globalReader.GetAll .
func GetAll(readOnly bool) (map[string][]byte, error) {
	return globalReader.GetAll(readOnly)
}

func (r *Reader) getBaseDir(readOnly bool) string {
	if readOnly {
		return path.Join(r.VpdDir, "ro")
	}
	return path.Join(r.VpdDir, "rw")
}

// Reader is a VPD reader object.
type Reader struct {
	VpdDir string
}

// Get reads a VPD variable by name and returns its value as a sequence of
// bytes. The `readOnly` flag specifies whether the variable is read-only or
// read-write.
func (r *Reader) Get(key string, readOnly bool) ([]byte, error) {
	buf, err := os.ReadFile(path.Join(r.getBaseDir(readOnly), key))
	if err != nil {
		return []byte{}, err
	}
	return buf, nil
}

// Set sets a VPD variable with `key` as name and `value` as its byte-stream
// value. The `readOnly` flag specifies whether the variable is read-only or
// read-write.
// NOTE Unfortunately Set doesn't currently work, because the sysfs interface
// does not support writing. To write, this library needs a backend able to
// write to flash chips, like the command line tool flashrom or flashtools.
func (r *Reader) Set(key string, value []byte, readOnly bool) error {
	// NOTE this is not implemented yet in the kernel interface, and will always
	// return a permission denied error
	return os.WriteFile(path.Join(r.getBaseDir(readOnly), key), value, 0o644)
}

// GetAll reads all the VPD variables and returns a map contaiing each
// name:value couple. The `readOnly` flag specifies whether the variable is
// read-only or read-write.
func (r *Reader) GetAll(readOnly bool) (map[string][]byte, error) {
	vpdMap := make(map[string][]byte)
	baseDir := r.getBaseDir(readOnly)
	err := filepath.Walk(baseDir, func(fpath string, info os.FileInfo, _ error) error {
		key := path.Base(fpath)
		if key == "." || key == "/" || fpath == baseDir {
			// empty or all slashes?
			return nil
		}
		value, err := r.Get(key, readOnly)
		if err != nil {
			return err
		}
		vpdMap[key] = value
		return nil
	})
	return vpdMap, err
}
