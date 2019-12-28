// Copyright 2019 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

package tss

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/google/go-tpm/tpm"
	"github.com/google/go-tpm/tpm2"
)

const (
	tpmRoot = "/sys/class/tpm"

	tpmPtManufacturer = 0x00000100 + 5  // PT_FIXED + offset of 5
	tpmPtVendorString = 0x00000100 + 6  // PT_FIXED + offset of 6
	tpmPtFwVersion1   = 0x00000100 + 11 // PT_FIXED + offset of 11
)

// TCGVendorID TPM manufacturer id
type TCGVendorID uint32

func (id TCGVendorID) String() string {
	return vendors[id]
}

// TPMVersion is used to configure a preference in
// which TPM to use, if multiple are available.
type TPMVersion uint8

// TPMInterface indicates how the client communicates
// with the TPM.
type TPMInterface uint8

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

// TPM versions
const (
	TPMVersionAgnostic TPMVersion = iota
	TPMVersion12
	TPMVersion20
)

// TPM interfaces
const (
	TPMInterfaceDirect TPMInterface = iota
	TPMInterfaceKernelManaged
	TPMInterfaceDaemonManaged
)

// TPM interfaces with a TPM device on the system.
type TPM struct {
	version TPMVersion
	interf  TPMInterface

	sysPath string
	rwc     io.ReadWriteCloser
}

// probedTPM identifies a TPM device on the system, which
// is a candidate for being used.
type probedTPM struct {
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

func probeSystemTPMs() ([]probedTPM, error) {
	var tpms []probedTPM

	tpmDevs, err := ioutil.ReadDir(tpmRoot)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if err == nil {
		for _, tpmDev := range tpmDevs {
			if strings.HasPrefix(tpmDev.Name(), "tpm") {
				tpm := probedTPM{
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

func openTPM(pTPM probedTPM) (*TPM, error) {
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
		version: pTPM.Version,
		interf:  interf,
		sysPath: pTPM.Path,
		rwc:     rwc,
	}, nil
}
