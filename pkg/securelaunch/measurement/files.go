// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package measurement

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/u-root/u-root/pkg/mount"
	slaunch "github.com/u-root/u-root/pkg/securelaunch"
	"github.com/u-root/u-root/pkg/securelaunch/tpm"
)

/* describes the "files" portion of policy file */
type FileCollector struct {
	Type  string   `json:"type"`
	Paths []string `json:"paths"`
}

/*
 * NewFileCollector extracts the "files" portion from the policy file.
 * initializes a new FileCollector structure.
 * returns error if unmarshalling of FileCollector fails
 */
func NewFileCollector(config []byte) (Collector, error) {
	slaunch.Debug("New Files Collector initialized\n")
	var fc = new(FileCollector)
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

/*
 * HashFile reads file input by user and calls TPM to measure it and store the hash.
 *
 * inputVal is of format <block device identifier>:<path>
 * E.g sda:/path/to/file _OR UUID:/path/to/file
 * Performs following actions
 * 1. mount device
 * 2. Read file on device into a byte slice.
 * 3. Unmount device
 * 4. Call tpm package which measures byte slice and stores it.
 */
func HashFile(inputVal string) error {
	// inputVal is of type sda:path
	mntFilePath, e := slaunch.GetMountedFilePath(inputVal, mount.MS_RDONLY)
	if e != nil {
		log.Printf("HashFile: GetMountedFilePath err=%v", e)
		return fmt.Errorf("failed to get mount path, err=%v", e)
	}
	slaunch.Debug("File Collector: Reading file=%s", mntFilePath)

	slaunch.Debug("File Collector: fileP=%s\n", mntFilePath)
	d, err := ioutil.ReadFile(mntFilePath)
	if err != nil {
		return fmt.Errorf("failed to read target file: filePath=%s, inputVal=%s, err=%v",
			mntFilePath, inputVal, err)
	}

	eventDesc := fmt.Sprintf("File Collector: measured %s", inputVal)
	return tpm.ExtendPCRDebug(pcr, bytes.NewReader(d), eventDesc)
}

/*
 * Collect satisfies Collector Interface. It loops over all file paths provided by user
 * and for each file path,  calls HashFile. HashFile measures each file on
 * that path and stores the result in TPM.
 */
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
