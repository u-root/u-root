// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tpm reads and extends pcrs with measurements.
package tpm

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
	slaunch "github.com/u-root/u-root/pkg/securelaunch"
	"github.com/u-root/u-root/pkg/securelaunch/eventlog"
)

/*
 * tpm2.ReadPCR and tpm2.ExtendPCR need hashAlgo passed.
 * using sha256 for now
 */
const (
	hashAlgo = tpm2.AlgSHA256
)

// marshalPcrEvent writes structure fields piecemeal to buffer.
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

/*
 * hashReader calculates the sha256 sum of an io reader.
 */
func hashReader(f io.Reader) []byte {

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return h.Sum(nil)
}

/*
 * GetHandle returns a tpm device handle from go-tpm/tpm2
 * returns a tpm handle from go-tpm/tpm2
 * that can be used for storing hashes.
 */
func GetHandle() (io.ReadWriteCloser, error) {
	tpm2, err := tpm2.OpenTPM("/dev/tpm0")
	if err != nil {
		return nil, fmt.Errorf("couldn't talk to TPM Device: err=%v", err)
	}

	return tpm2, nil
}

/*
 * ReadPCR reads pcr#x, where x is provided by 'pcr' arg and returns
 * the result in a byte slice.
 * 'tpmHandle' is the tpm device that owns the 'pcr'.
 * err is returned if read fails.
 */
func ReadPCR(tpmHandle io.ReadWriteCloser, pcr int) ([]byte, error) {
	val, err := tpm2.ReadPCR(tpmHandle, pcr, hashAlgo)
	if err != nil {
		return nil, fmt.Errorf("can't read PCR %d, err= %v", pcr, err)
	}
	return val, nil
}

/*
 * ExtendPCR writes the measurements passed as 'hash' arg to pcr#x,
 * where x is provided by 'pcr' arg.
 *
 * pcr is owned by 'tpm2Handle', a tpm device handle.
 * err is returned if write to pcr fails.
 */
func ExtendPCR(tpmHandle io.ReadWriteCloser, pcr int, hash []byte) error {
	return tpm2.PCRExtend(tpmHandle, tpmutil.Handle(pcr), hashAlgo, hash, "")
}

/*
 * ExtendPCRDebug extends a PCR with the contents of a byte slice
 * and notifies the kernel of this measurement by sending event via sysfs.
 *
 * In debug mode, it prints
 * 1. old pcr value before the hash is written to pcr
 * 2. new pcr values after hash is written to pcr
 * 3. compares old and new pcr values and prints error if they are not
 */
func ExtendPCRDebug(tpmHandle io.ReadWriteCloser, pcr int, data io.Reader, eventDesc string) error {
	oldPCRValue, err := ReadPCR(tpmHandle, pcr)
	if err != nil {
		return fmt.Errorf("ReadPCR failed, err=%v", err)
	}
	slaunch.Debug("ExtendPCRDebug: oldPCRValue = [%x]", oldPCRValue)

	hash := hashReader(data)

	slaunch.Debug("Adding hash=[%x] to PCR #%d", hash, pcr)
	if e := ExtendPCR(tpmHandle, pcr, hash); e != nil {
		return fmt.Errorf("can't extend PCR %d, err=%v", pcr, e)
	}
	slaunch.Debug(eventDesc)

	// send event if PCR was successfully extended above.
	sendEventToSysfs(uint32(pcr), hash, []byte(eventDesc))

	newPCRValue, err := ReadPCR(tpmHandle, pcr)
	if err != nil {
		return fmt.Errorf("ReadPCR failed, err=%v", err)
	}
	slaunch.Debug("ExtendPCRDebug: newPCRValue = [%x]", newPCRValue)

	finalPCR := hashReader(bytes.NewReader(append(oldPCRValue, hash...)))
	if !bytes.Equal(finalPCR, newPCRValue) {
		return fmt.Errorf("PCRs not equal, got %x, want %x", finalPCR, newPCRValue)
	}

	return nil
}
