// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vpd

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// VpdDir points to the base directory where the VPD sysfs interface is located.
// It is an exported variable to allow for testing
var (
	VpdDir       = "/sys/firmware/vpd"
	MaxBootEntry = 9999
)

func getBaseDir(readOnly bool) string {
	var baseDir string
	if readOnly {
		baseDir = path.Join(VpdDir, "ro")
	} else {
		baseDir = path.Join(VpdDir, "rw")
	}
	return baseDir
}

// Get reads a VPD variable by name and returns its value as a sequence of
// bytes. The `readOnly` flag specifies whether the variable is read-only or
// read-write.
func Get(key string, readOnly bool) ([]byte, error) {
	buf, err := ioutil.ReadFile(path.Join(getBaseDir(readOnly), key))
	if err != nil {
		return []byte{}, err
	}
	return buf, nil
}

// GetAll reads all the VPD variables and returns a map contaiing each
// name:value couple. The `readOnly` flag specifies whether to read RO
// or RW variables.
func GetAll(readOnly bool) (map[string][]byte, error) {
	vpdMap := make(map[string][]byte)
	baseDir := getBaseDir(readOnly)
	err := filepath.Walk(baseDir, func(fpath string, info os.FileInfo, _ error) error {
		key := path.Base(fpath)
		if key == "." || key == "/" || fpath == baseDir {
			// empty or all slashes?
			return nil
		}
		value, err := Get(key, readOnly)
		if err != nil {
			return err
		}
		vpdMap[key] = value
		return nil
	})
	return vpdMap, err
}
