// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tss

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-tpm/legacy/tpm2"
	"github.com/google/go-tpm/tpm"
)

const (
	tpmRoot = "/sys/class/tpm"
)

func probeSystemTPMs() ([]probedTPM, error) {
	var tpms []probedTPM

	tpmDevs, err := os.ReadDir(tpmRoot)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// TPM look up is hardcoded. Taken from googles go-attestation.
	// go-tpm does not support GetCapability with the required subcommand.
	// Implementation will be updated asap this is fixed in Go-tpm
	for _, tpmDev := range tpmDevs {
		if strings.HasPrefix(tpmDev.Name(), "tpm") {
			tpm := probedTPM{
				Path: filepath.Join(tpmRoot, tpmDev.Name()),
			}

			if _, err := os.Stat(filepath.Join(tpm.Path, "caps")); err != nil {
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

	return tpms, nil
}

func newTPM(pTPM probedTPM) (*TPM, error) {
	interf := TPMInterfaceDirect
	var rwc io.ReadWriteCloser
	var err error

	switch pTPM.Version {
	case TPMVersion12:
		devPath := filepath.Join("/dev", filepath.Base(pTPM.Path))
		interf = TPMInterfaceKernelManaged

		rwc, err = tpm.OpenTPM(devPath)
		if err != nil {
			return nil, err
		}
	case TPMVersion20:
		// If the TPM has a kernel-provided resource manager, we should
		// use that instead of communicating directly.
		devPath := filepath.Join("/dev", filepath.Base(pTPM.Path))
		f, err := os.ReadDir(filepath.Join(pTPM.Path, "device", "tpmrm"))
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}
		} else if len(f) > 0 {
			devPath = filepath.Join("/dev", f[0].Name())
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

// MeasurementLog reads the TCPA eventlog in binary format
// from the Linux kernel
func (t *TPM) MeasurementLog() ([]byte, error) {
	return os.ReadFile("/sys/kernel/security/tpm0/binary_bios_measurements")
}
