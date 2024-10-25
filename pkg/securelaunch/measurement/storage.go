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

// StorageCollector describes the "storage" portion of the policy file.
type StorageCollector struct {
	Type  string   `json:"type"`
	Paths []string `json:"paths"`
}

// NewStorageCollector extracts the "storage" portion from the policy file and
// initializes a new StorageCollector structure.
//
// It returns an error if unmarshalling of StorageCollector fails.
func NewStorageCollector(config []byte) (Collector, error) {
	slaunch.Debug("New Storage Collector initialized\n")
	sc := new(StorageCollector)
	err := json.Unmarshal(config, &sc)
	if err != nil {
		return nil, err
	}
	return sc, nil
}

// measureStorageDevice reads the given disk path and measures it into the TPM.
//
// blkDevicePath is a string to a block device (e.g., /dev/sda).
// It returns and error if Reading the block device fails.
func measureStorageDevice(blkDevicePath string) error {
	log.Printf("Storage Collector: Measuring block device %s\n", blkDevicePath)
	file, err := os.Open(blkDevicePath)
	if err != nil {
		return fmt.Errorf("couldn't open disk=%s err=%w", blkDevicePath, err)
	}

	eventDesc := fmt.Sprintf("Storage Collector: Measured %s", blkDevicePath)
	return tpm.ExtendPCRDebug(pcr, file, eventDesc)
}

// Collect loops over the given storage paths and for each storage path calls
// measureStorageDevice(), which measures a storage device into the TPM.
//
// It satisfies the Collector interface.
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
