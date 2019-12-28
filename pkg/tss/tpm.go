// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tss

import (
	"errors"
	"fmt"
	"io/ioutil"
)

// MatchesConfig returns true if the TPM satisfies the constraints
// specified by the given config.
func (t *probedTPM) matchesConfig(config OpenConfig) bool {
	return config.TPMVersion == TPMVersionAgnostic || t.Version == config.TPMVersion
}

// OpenTPM initializes access to the TPM based on the
// config provided.
func OpenTPM(config *OpenConfig) (*TPM, error) {
	if config == nil {
		config = &OpenConfig{}
	}
	candidateTPMs, err := probeSystemTPMs()
	if err != nil {
		return nil, err
	}

	for _, tpm := range candidateTPMs {
		if tpm.matchesConfig(*config) {
			return openTPM(tpm)
		}
	}

	return nil, errors.New("TPM device not available")
}

// MeasurementLog reads the TCPA eventlog in binary format
// from the Linux kernel
func (t *TPM) MeasurementLog() ([]byte, error) {
	return ioutil.ReadFile("/sys/kernel/security/tpm0/binary_bios_measurements")
}

// Info returns information about the TPM.
func (t *TPM) Info() (*TPMInfo, error) {
	var info TPMInfo
	var err error
	switch t.version {
	case TPMVersion12:
		info, err = readTPM12VendorAttributes(t.rwc)
	case TPMVersion20:
		info, err = readTPM20VendorAttributes(t.rwc)
	default:
		return nil, fmt.Errorf("unsupported TPM version: %x", t.version)
	}
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// Version returns the TPM version
func (t *TPM) Version() TPMVersion {
	return t.version
}

// Close closes the TPM socket
func (t *TPM) Close() error {
	return t.rwc.Close()
}
