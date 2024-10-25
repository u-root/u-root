// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package measurement

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/mount"
	slaunch "github.com/u-root/u-root/pkg/securelaunch"
	"github.com/u-root/u-root/pkg/securelaunch/tpm"
)

// Describes the "files" portion of policy file.
type FileCollector struct {
	Type  string   `json:"type"`
	Paths []string `json:"paths"`
}

// NewFileCollector extracts the "files" portion from the policy file and
// initializes a new FileCollector structure.
// It returns an error if unmarshalling of FileCollector fails.
func NewFileCollector(config []byte) (Collector, error) {
	slaunch.Debug("New Files Collector initialized\n")
	fc := new(FileCollector)
	err := json.Unmarshal(config, &fc)
	if err != nil {
		return nil, err
	}
	return fc, nil
}

// HashBytes extends PCR with a byte array and sends an event to sysfs.
// the sent event is described via eventDesc.
func HashBytes(b []byte, eventDesc string) error {
	return tpm.ExtendPCRDebug(pcr, bytes.NewReader(b), eventDesc)
}

// HashFile opens and reads the given file and measures it into the TPM.
//
// inputVal is of format <block device identifier>:<path>
// (e.g., `sda:/path/to/file` or `UUID:/path/to/file`).
func HashFile(inputVal string) error {
	// inputVal is of type sda:path
	mntFilePath, e := slaunch.GetMountedFilePath(inputVal, mount.MS_RDONLY)
	if e != nil {
		log.Printf("HashFile: GetMountedFilePath err=%v", e)
		return fmt.Errorf("failed to get mount path, err=%w", e)
	}
	slaunch.Debug("File Collector: Reading file=%s", mntFilePath)

	slaunch.Debug("File Collector: fileP=%s\n", mntFilePath)
	d, err := os.ReadFile(mntFilePath)
	if err != nil {
		return fmt.Errorf("failed to read target file: filePath=%s, inputVal=%s, err=%w",
			mntFilePath, inputVal, err)
	}

	eventDesc := fmt.Sprintf("File Collector: measured %s", inputVal)
	return tpm.ExtendPCRDebug(pcr, bytes.NewReader(d), eventDesc)
}

// Collect loops over the given file paths and for each file path calls
// HashFile(), which measures a file into the TPM.
//
// It satisfies the Collector interface.
func (s *FileCollector) Collect() error {
	for _, inputVal := range s.Paths {
		// inputVal is of type sda:/path/to/file
		err := HashFile(inputVal)
		if err != nil {
			log.Printf("File Collector: input=%s, err = %v", inputVal, err)
			return err
		}
	}

	return nil
}
