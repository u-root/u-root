// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// smbios_transfer sends SMBIOS tables to BMC through IPMI blob interfaces.
//
// Synopsis:
//     smbios_tranfer [-num_retries]
//
// Options:
//	-num_retries: number of times to retry transferring SMBIOS tables
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/u-root/u-root/pkg/ipmi"
	"github.com/u-root/u-root/pkg/ipmi/blobs"
)

const (
	maxWriteSize uint32 = 128

	// IPMI blob ID on BMC
	smbiosBlobID = "/smbios\x00"

	sysfsPath = "/sys/firmware/dmi/tables"
)

var (
	retries = flag.Int("num_retries", 2, "Number of times to retry transferring SMBIOS tables")
)

func writeCommitSmbiosBlob(id string, data []uint8, h *blobs.BlobHandler) (rerr error) {
	sessionID, err := h.BlobOpen(id, blobs.BMC_BLOB_OPEN_FLAG_WRITE)
	if err != nil {
		return fmt.Errorf("IPMI BlobOpen for %s failed: %v", id, err)
	}
	defer func() {
		// If the function returned successfully but failed to close the blob,
		// return an error.
		if err := h.BlobClose(sessionID); err != nil {
			err = fmt.Errorf("IPMI BlobClose %s failed: %v", id, err)
			if rerr != nil {
				rerr = fmt.Errorf("%v; %v", rerr, err)
				return
			}
			rerr = err
		}
	}()

	dataLen := uint32(len(data))

	// IPMI max message length defined in ipmi_msgdefs.h as IPMI_MAX_MSG_LENGTH.
	// Read/write longer than the limit will be requested in multiple IPMI
	// commands.
	for offset := uint32(0); offset < dataLen; offset += maxWriteSize {
		end := offset + maxWriteSize
		if end > dataLen {
			end = dataLen
		}
		if err = h.BlobWrite(sessionID, offset, data[offset:end]); err != nil {
			return fmt.Errorf("IPMI BlobWrite %s failed: %v", id, err)
		}
	}

	if err = h.BlobCommit(sessionID, []uint8{}); err != nil {
		return fmt.Errorf("IPMI BlobCommit %s failed: %v", id, err)
	}

	return nil
}

func getSmbiosData() ([]uint8, error) {
	tables, err := ioutil.ReadFile(filepath.Join(sysfsPath, "DMI"))
	if err != nil {
		return nil, fmt.Errorf("error reading DMI data: %v", err)
	}

	return tables, nil
}

func transferSmbiosData() error {
	data, err := getSmbiosData()
	if err != nil {
		return fmt.Errorf("failed to get SMBIOS tables")
	}
	i, err := ipmi.Open(0)
	if err != nil {
		return err
	}
	h := blobs.NewBlobHandler(i)

	blobCount, err := h.BlobGetCount()
	if err != nil {
		return fmt.Errorf("failed to get blob count: %v", err)
	}

	seen := false
	for j := 0; j < blobCount; j++ {
		id, err := h.BlobEnumerate(j)
		if err != nil {
			return fmt.Errorf("failed to enumerate blob %d: %v", j, err)
		}

		if id != smbiosBlobID {
			continue
		}

		seen = true
		if err = writeCommitSmbiosBlob(id, data, h); err != nil {
			return fmt.Errorf("failed to write and commit blob %s: %v", id, err)
		}
		break
	}

	if !seen {
		return fmt.Errorf("no smbios blob found")
	}

	return nil
}

func main() {
	flag.Parse()
	for r := 0; r < *retries; r++ {
		log.Printf("Transferring SMBIOS tables, attempt %d/%d", r+1, *retries)
		if err := transferSmbiosData(); err != nil {
			log.Printf("Error tranferring SMBIOS tables over IPMI: %v", err)
		} else {
			log.Printf("SMBIOS tables are tranferred.")
			break
		}
	}
}
