// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tss provides TPM 1.2/2.0 core functionality and
// abstraction layer for high-level functions

package tss

import (
	"crypto"
	"io"

	"github.com/google/go-tpm/tpm2"
)

// TCGVendorID TPM manufacturer id
type TCGVendorID uint32

func (id TCGVendorID) String() string {
	return vendors[id]
}

// HashAlg is the TPM hash algorithm id
type HashAlg uint8

var (
	// HashSHA1 is the TPM 1.2 identifier for SHA1
	HashSHA1 = HashAlg(tpm2.AlgSHA1)
	// HashSHA256 is the TPM 2.0 identifier for SHA256
	HashSHA256 = HashAlg(tpm2.AlgSHA256)
)

var vendors = map[TCGVendorID]string{
	1095582720: "AMD",
	1096043852: "Atmel",
	1112687437: "Broadcom",
	1229081856: "IBM",
	1213220096: "HPE",
	1297303124: "Microsoft",
	1229346816: "Infineon",
	1229870147: "Intel",
	1279610368: "Lenovo",
	1314082080: "National Semiconductor",
	1314150912: "Nationz",
	1314145024: "Nuvoton Technology",
	1363365709: "Qualcomm",
	1397576515: "SMSC",
	1398033696: "ST Microelectronics",
	1397576526: "Samsung",
	1397641984: "Sinosun",
	1415073280: "Texas Instruments",
	1464156928: "Winbond",
	1380926275: "Fuzhou Rockchip",
	1196379975: "Google",
}

// PCR encapsulates the value of a PCR at a point in time.
type PCR struct {
	Index     int
	Digest    []byte
	DigestAlg crypto.Hash
}

// TPM interfaces with a TPM device on the system.
type TPM struct {
	Version TPMVersion
	Interf  TPMInterface

	SysPath string
	RWC     io.ReadWriteCloser
}

// probedTPM identifies a TPM device on the system, which
// is a candidate for being used.
type ProbedTPM struct {
	Version TPMVersion
	Path    string
}

// TPMInfo contains information about the version & interface
// of an open TPM.
type TPMInfo struct {
	Version      TPMVersion
	Interface    TPMInterface
	VendorInfo   string
	Manufacturer TCGVendorID

	// FirmwareVersionMajor and FirmwareVersionMinor describe
	// the firmware version of the TPM, but are only available
	// for TPM 2.0 devices.
	FirmwareVersionMajor int
	FirmwareVersionMinor int
}
