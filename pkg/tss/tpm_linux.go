// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tss provides TPM 1.2/2.0 core functionality and
// abstraction layer for high-level functions

package tss

import (
	"crypto"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/google/go-tpm/tpm"
	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

const (
	tpmRoot = "/sys/class/tpm"

	tpmPtManufacturer = 0x00000100 + 5  // PT_FIXED + offset of 5
	tpmPtVendorString = 0x00000100 + 6  // PT_FIXED + offset of 6
	tpmPtFwVersion1   = 0x00000100 + 11 // PT_FIXED + offset of 11
)

func readTPM12VendorAttributes(rwc io.ReadWriter) (TPMInfo, error) {
	var vendorInfo string

	_, err := tpm.GetManufacturer(rwc)
	if err != nil {
		return TPMInfo{}, err
	}

	return TPMInfo{
		VendorInfo:           strings.Trim(vendorInfo, "\x00"), // Stubbed
		Manufacturer:         TCGVendorID(uint32(0)),           // Stubbed
		FirmwareVersionMajor: int(0),                           // Stubbed
		FirmwareVersionMinor: int(0),                           // Stubbed
	}, nil
}

func readTPM20VendorAttributes(rwc io.ReadWriter) (TPMInfo, error) {
	var vendorInfo string
	// The Vendor String is split up into 4 sections of 4 bytes,
	// for a maximum length of 16 octets of ASCII text. We iterate
	// through the 4 indexes to get all 16 bytes & construct vendorInfo.
	// See: TPM_PT_VENDOR_STRING_1 in TPM 2.0 Structures reference.
	for i := 0; i < 4; i++ {
		caps, _, err := tpm2.GetCapability(rwc, tpm2.CapabilityTPMProperties, 1, tpmPtVendorString+uint32(i))
		if err != nil {
			return TPMInfo{}, fmt.Errorf("tpm2.GetCapability(PT_VENDOR_STRING_%d) failed: %v", i+1, err)
		}
		subset, ok := caps[0].(tpm2.TaggedProperty)
		if !ok {
			return TPMInfo{}, fmt.Errorf("got capability of type %T, want tpm2.TaggedProperty", caps[0])
		}
		// Reconstruct the 4 ASCII octets from the uint32 value.
		vendorInfo += string(subset.Value&0xFF000000) + string(subset.Value&0xFF0000) + string(subset.Value&0xFF00) + string(subset.Value&0xFF)
	}

	caps, _, err := tpm2.GetCapability(rwc, tpm2.CapabilityTPMProperties, 1, tpmPtManufacturer)
	if err != nil {
		return TPMInfo{}, fmt.Errorf("tpm2.GetCapability(PT_MANUFACTURER) failed: %v", err)
	}
	manu, ok := caps[0].(tpm2.TaggedProperty)
	if !ok {
		return TPMInfo{}, fmt.Errorf("got capability of type %T, want tpm2.TaggedProperty", caps[0])
	}

	caps, _, err = tpm2.GetCapability(rwc, tpm2.CapabilityTPMProperties, 1, tpmPtFwVersion1)
	if err != nil {
		return TPMInfo{}, fmt.Errorf("tpm2.GetCapability(PT_FIRMWARE_VERSION_1) failed: %v", err)
	}
	fw, ok := caps[0].(tpm2.TaggedProperty)
	if !ok {
		return TPMInfo{}, fmt.Errorf("got capability of type %T, want tpm2.TaggedProperty", caps[0])
	}

	return TPMInfo{
		VendorInfo:           strings.Trim(vendorInfo, "\x00"),
		Manufacturer:         TCGVendorID(manu.Value),
		FirmwareVersionMajor: int((fw.Value & 0xffff0000) >> 16),
		FirmwareVersionMinor: int(fw.Value & 0x0000ffff),
	}, nil
}

func probeSystemTPMs() ([]ProbedTPM, error) {
	var tpms []ProbedTPM

	tpmDevs, err := ioutil.ReadDir(tpmRoot)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if err == nil {
		for _, tpmDev := range tpmDevs {
			if strings.HasPrefix(tpmDev.Name(), "tpm") {
				tpm := ProbedTPM{
					Path: path.Join(tpmRoot, tpmDev.Name()),
				}

				if _, err := os.Stat(path.Join(tpm.Path, "caps")); err != nil {
					if !os.IsNotExist(err) {
						return nil, err
					}
					tpm.Version = TPMVersion20
				} else {
					tpm.Version = TPMVersion12
				}
				tpms = append(tpms, tpm)
			}
		}
	}

	return tpms, nil
}

func openTPM(pTPM ProbedTPM) (*TPM, error) {
	interf := TPMInterfaceDirect
	var rwc io.ReadWriteCloser
	var err error

	switch pTPM.Version {
	case TPMVersion12:
		devPath := path.Join("/dev", path.Base(pTPM.Path))
		interf = TPMInterfaceKernelManaged

		rwc, err = tpm.OpenTPM(devPath)
		if err != nil {
			return nil, err
		}
	case TPMVersion20:
		// If the TPM has a kernel-provided resource manager, we should
		// use that instead of communicating directly.
		devPath := path.Join("/dev", path.Base(pTPM.Path))
		f, err := ioutil.ReadDir(path.Join(pTPM.Path, "device", "tpmrm"))
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}
		} else if len(f) > 0 {
			devPath = path.Join("/dev", f[0].Name())
			interf = TPMInterfaceKernelManaged
		}

		rwc, err = tpm2.OpenTPM(devPath)
		if err != nil {
			return nil, err
		}
	}

	return &TPM{
		Version: pTPM.Version,
		Interf:  interf,
		SysPath: pTPM.Path,
		RWC:     rwc,
	}, nil
}

func readAllPCRs20(tpm io.ReadWriter, alg tpm2.Algorithm) (map[uint32][]byte, error) {
	numPCRs := 24
	out := map[uint32][]byte{}

	// The TPM 2.0 spec says that the TPM can partially fulfill the
	// request. As such, we repeat the command up to 8 times to get all
	// 24 PCRs.
	for i := 0; i < numPCRs; i++ {
		// Build a selection structure, specifying all PCRs we do
		// not have the value for.
		sel := tpm2.PCRSelection{Hash: alg}
		for pcr := 0; pcr < numPCRs; pcr++ {
			if _, present := out[uint32(pcr)]; !present {
				sel.PCRs = append(sel.PCRs, pcr)
			}
		}

		// Ask the TPM for those PCR values.
		ret, err := tpm2.ReadPCRs(tpm, sel)
		if err != nil {
			return nil, fmt.Errorf("tpm2.ReadPCRs(%+v) failed with err: %v", sel, err)
		}
		// Keep track of the PCRs we were actually given.
		for pcr, digest := range ret {
			out[uint32(pcr)] = digest
		}
		if len(out) == numPCRs {
			break
		}
	}

	if len(out) != numPCRs {
		return nil, fmt.Errorf("failed to read all PCRs, only read %d", len(out))
	}

	return out, nil
}

func readAllPCRs12(rwc io.ReadWriter) (map[uint32][]byte, error) {
	numPCRs := 24
	out := map[uint32][]byte{}

	for i := 0; i < numPCRs; i++ {
		// Ask the TPM for those PCR values.
		pcr, err := tpm.ReadPCR(rwc, uint32(i))
		if err != nil {
			return nil, fmt.Errorf("tpm.ReadPCR(%d) failed with err: %v", i, err)
		}
		out[uint32(i)] = pcr
		if len(out) == numPCRs {
			break
		}
	}

	if len(out) != numPCRs {
		return nil, fmt.Errorf("failed to read all PCRs, only read %d", len(out))
	}

	return out, nil
}

func (a HashAlg) cryptoHash() crypto.Hash {
	switch a {
	case HashSHA1:
		return crypto.SHA1
	case HashSHA256:
		return crypto.SHA256
	}
	return 0
}

func (a HashAlg) goTPMAlg() tpm2.Algorithm {
	switch a {
	case HashSHA1:
		return tpm2.AlgSHA1
	case HashSHA256:
		return tpm2.AlgSHA256
	}
	return 0
}

func extendPCR12(rwc io.ReadWriter, pcrIndex uint32, hash [20]byte) error {
	_, err := tpm.PcrExtend(rwc, pcrIndex, hash)
	if err != nil {
		return err
	}
	return nil
}

func extendPCR20(rwc io.ReadWriter, pcrIndex uint32, hash []byte, alg HashAlg) error {
	err := tpm2.PCRExtend(rwc, tpmutil.Handle(pcrIndex), alg.goTPMAlg(), hash, "")
	if err != nil {
		return err
	}
	return nil
}

func readPCR12(rwc io.ReadWriter, pcrIndex uint32) ([]byte, error) {
	return tpm.ReadPCR(rwc, pcrIndex)
}

func readPCR20(rwc io.ReadWriter, pcrIndex uint32, alg HashAlg) ([]byte, error) {
	return tpm2.ReadPCR(rwc, int(pcrIndex), alg.goTPMAlg())
}
