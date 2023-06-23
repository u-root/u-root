// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tpm reads and extends pcrs with measurements.
package tpm

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/google/go-tpm/legacy/tpm2"
	slaunch "github.com/u-root/u-root/pkg/securelaunch"
	"github.com/u-root/u-root/pkg/securelaunch/eventlog"
	"github.com/u-root/u-root/pkg/tss"
)

var (
	hashAlgo  = tpm2.AlgSHA256
	tpmHandle *tss.TPM
)

// marshalPcrEvent writes structure fields piecemeal to a buffer.
func marshalPcrEvent(pcr uint32, h []byte, eventDesc []byte) ([]byte, error) {
	const baseTypeTXT = 0x400                       // TXT specification base event value for DRTM values
	const slaunchType = uint32(baseTypeTXT + 0x102) // Secure Launch event log entry type.
	count := uint32(1)
	eventDescLen := uint32(len(eventDesc))
	slaunch.Debug("marshalPcrEvent: pcr=[%v], slaunchType=[%v], count=[%v], hashAlgo=[%v], eventDesc=[%s], eventDescLen=[%v]",
		pcr, slaunchType, count, hashAlgo, eventDesc, eventDescLen)

	endianess := binary.LittleEndian
	var buf bytes.Buffer

	if err := binary.Write(&buf, endianess, pcr); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, endianess, slaunchType); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, endianess, count); err != nil {
		return nil, err
	}

	for i := uint32(0); i < count; i++ {
		if err := binary.Write(&buf, endianess, hashAlgo); err != nil {
			return nil, err
		}

		if err := binary.Write(&buf, endianess, h); err != nil {
			return nil, err
		}
	}

	if err := binary.Write(&buf, endianess, eventDescLen); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, endianess, eventDesc); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// sendEventToSysfs marshals measurement events and writes them to sysfs.
func sendEventToSysfs(pcr uint32, h []byte, eventDesc []byte) {
	b, err := marshalPcrEvent(pcr, h, eventDesc)
	if err != nil {
		log.Println(err)
	}

	if e := eventlog.Add(b); e != nil {
		log.Println(e)
	}
}

// hashReader calculates the sha256 sum of an io reader.
func hashReader(f io.Reader) []byte {
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return h.Sum(nil)
}

// New sets up a TPM device handle that can be used for storing hashes.
func New() error {
	tpm, err := tss.NewTPM()
	if err != nil {
		return fmt.Errorf("couldn't talk to TPM Device: err=%v", err)
	}

	tpmHandle = tpm
	return nil
}

// Close closes the connection of a TPM device handle.
func Close() {
	if tpmHandle != nil {
		tpmHandle.Close()
		tpmHandle = nil
	}
}

// readPCR reads the given PCR and returns the result in a byte slice.
func readPCR(pcr uint32) ([]byte, error) {
	if tpmHandle == nil {
		return nil, errors.New("tpmHandle is nil")
	}

	val, err := tpmHandle.ReadPCR(pcr)
	if err != nil {
		return nil, fmt.Errorf("can't read PCR %d, err= %v", pcr, err)
	}
	return val, nil
}

// extendPCR extends the given PCR with the given hash.
func extendPCR(pcr uint32, hash []byte) error {
	if tpmHandle == nil {
		return errors.New("tpmHandle is nil")
	}

	return tpmHandle.Extend(hash, pcr)
}

// ExtendPCRDebug extends a PCR with the contents of a byte slice and notifies
// the kernel of this measurement by sending an event via sysfs.
//
// In debug mode, it prints:
//  1. The old PCR value before the hash is extended to the PCR
//  2. The new PCR value after the hash is extended to the PCR
func ExtendPCRDebug(pcr uint32, data io.Reader, eventDesc string) error {
	oldPCRValue, err := readPCR(pcr)
	if err != nil {
		return fmt.Errorf("readPCR failed, err=%v", err)
	}
	slaunch.Debug("ExtendPCRDebug: oldPCRValue = [%x]", oldPCRValue)

	hash := hashReader(data)

	slaunch.Debug("Adding hash=[%x] to PCR #%d", hash, pcr)
	if e := extendPCR(pcr, hash); e != nil {
		return fmt.Errorf("can't extend PCR %d, err=%v", pcr, e)
	}
	slaunch.Debug(eventDesc)

	// send event if PCR was successfully extended above.
	sendEventToSysfs(pcr, hash, []byte(eventDesc))

	newPCRValue, err := readPCR(pcr)
	if err != nil {
		return fmt.Errorf("readPCR failed, err=%v", err)
	}
	slaunch.Debug("ExtendPCRDebug: newPCRValue = [%x]", newPCRValue)

	finalPCR := hashReader(bytes.NewReader(append(oldPCRValue, hash...)))
	if !bytes.Equal(finalPCR, newPCRValue) {
		return fmt.Errorf("PCRs not equal, got %x, want %x", finalPCR, newPCRValue)
	}

	return nil
}
