// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package measurement

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	slaunch "github.com/u-root/u-root/pkg/securelaunch"
	"github.com/u-root/u-root/pkg/securelaunch/tpm"
)

/* describes the "storage" portion of policy file */
type StorageCollector struct {
	Type  string   `json:"type"`
	Paths []string `json:"paths"`
}

/*
 * NewStorageCollector extracts the "storage" portion from the policy file.
 * initializes a new StorageCollector structure.
 * returns error if unmarshalling of StorageCollector fails
 */
func NewStorageCollector(config []byte) (Collector, error) {
	slaunch.Debug("New Storage Collector initialized\n")
	var sc = new(StorageCollector)
	err := json.Unmarshal(config, &sc)
	if err != nil {
		return nil, err
	}
	return sc, nil
}

/*
 * measureStorageDevice reads the disk path input by user,
 * and then extends the pcr with it.
 *
 * Hashing of buffer is handled by tpm package.
 * - blkDevicePath - string e.g /dev/sda
 * returns
 * - error if Reading the block device fails.
 */
func measureStorageDevice(blkDevicePath string) error {

	log.Printf("Storage Collector: Measuring block device %s\n", blkDevicePath)
	file, err := os.Open(blkDevicePath)
	if err != nil {
		return fmt.Errorf("couldn't open disk=%s err=%v", blkDevicePath, err)
	}

	eventDesc := fmt.Sprintf("Storage Collector: Measured %s", blkDevicePath)
	return tpm.ExtendPCRDebug(pcr, file, eventDesc)
}

/*
 * Collect satisfies Collector Interface. It loops over all storage paths provided
 * by user and calls measureStorageDevice for each storage path. storage path is of
 * form /dev/sda. measureStorageDevice in turn calls tpm
 * package which further hashes this buffer and extends pcr.
 */
func (s *StorageCollector) Collect() error {

	for _, inputVal := range s.Paths {
		device, e := slaunch.GetStorageDevice(inputVal) // inputVal is blkDevicePath e.g UUID or sda
		if e != nil {
			log.Printf("Storage Collector: input = %s, GetStorageDevice: err = %v", inputVal, e)
			return e
		}
		devPath := filepath.Join("/dev", device.Name)
		err := measureStorageDevice(devPath)
		if err != nil {
			log.Printf("Storage Collector: input = %s, err = %v", inputVal, err)
			return err
		}
	}

	return nil
}
