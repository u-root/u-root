// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tss provides TPM 1.2/2.0 core functionality and
// abstraction layer for high-level functions
package tss

import (
	"errors"
	"fmt"

	tpmutil "github.com/google/go-tpm/tpmutil"
)

// NewTPM initializes access to the TPM based on the
// config provided.
func NewTPM() (*TPM, error) {
	candidateTPMs, err := probeSystemTPMs()
	if err != nil {
		return nil, err
	}

	for _, tpm := range candidateTPMs {
		tss, err := newTPM(tpm)
		if err != nil {
			continue
		}
		return tss, nil
	}

	return nil, errors.New("TPM device not available")
}

// Info returns information about the TPM.
func (t *TPM) Info() (*TPMInfo, error) {
	var info TPMInfo
	var err error
	switch t.Version {
	case TPMVersion12:
		info, err = readTPM12VendorAttributes(t.RWC)
	case TPMVersion20:
		info, err = readTPM20VendorAttributes(t.RWC)
	default:
		return nil, fmt.Errorf("unsupported TPM version: %x", t.Version)
	}
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// GetVersion returns the TPM version
func (t *TPM) GetVersion() TPMVersion {
	return t.Version
}

// Close closes the TPM socket and wipe locked buffers
func (t *TPM) Close() error {
	return t.RWC.Close()
}

// ReadPCRs reads all PCRs into the PCR structure
func (t *TPM) ReadPCRs(alg HashAlg) ([]PCR, error) {
	var PCRs map[uint32][]byte
	var err error

	switch t.Version {
	case TPMVersion12:
		if alg != HashSHA1 {
			return nil, fmt.Errorf("non-SHA1 algorithm %v is not supported on TPM 1.2", alg)
		}
		PCRs, err = readAllPCRs12(t.RWC)
		if err != nil {
			return nil, fmt.Errorf("failed to read PCRs: %v", err)
		}

	case TPMVersion20:
		PCRs, err = readAllPCRs20(t.RWC, alg.goTPMAlg())
		if err != nil {
			return nil, fmt.Errorf("failed to read PCRs: %v", err)
		}

	default:
		return nil, fmt.Errorf("unsupported TPM version: %x", t.Version)
	}

	out := make([]PCR, len(PCRs))
	for index, digest := range PCRs {
		out[int(index)] = PCR{
			Index:     int(index),
			Digest:    digest,
			DigestAlg: alg.cryptoHash(),
		}
	}

	return out, nil
}

// Extend extends a hash into a pcrIndex with a specific hash algorithm
func (t *TPM) Extend(hash []byte, pcrIndex uint32) error {
	switch t.Version {
	case TPMVersion12:
		var thash [20]byte
		hashlen := len(hash)
		if hashlen != 20 {
			return fmt.Errorf("hash length invalid - need 20, got: %v", hashlen)
		}
		copy(thash[:], hash[:20])
		err := extendPCR12(t.RWC, pcrIndex, thash)
		if err != nil {
			return err
		}
	case TPMVersion20:
		err := extendPCR20(t.RWC, pcrIndex, hash)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported TPM version: %x", t.Version)
	}

	return nil
}

// Measure measures data with a specific hash algorithm and extends it into the pcrIndex
func (t *TPM) Measure(data []byte, pcrIndex uint32, alg HashAlg) error {
	switch t.Version {
	case TPMVersion12:
		hashFunc := HashSHA1.cryptoHash().New()
		hash := hashFunc.Sum(data)
		var thash [20]byte
		hashlen := len(hash)
		if hashlen != 20 {
			return fmt.Errorf("hash length insufficient - need 20, got: %v", hashlen)
		}
		copy(thash[:], hash[:20])
		err := extendPCR12(t.RWC, pcrIndex, thash)
		if err != nil {
			return err
		}
	case TPMVersion20:
		hashFunc := alg.cryptoHash().New()
		hash := hashFunc.Sum(data)
		err := extendPCR20(t.RWC, pcrIndex, hash)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported TPM version: %x", t.Version)
	}

	return nil
}

// ReadPCR reads a single PCR value by defining the pcrIndex
func (t *TPM) ReadPCR(pcrIndex uint32) ([]byte, error) {
	switch t.Version {
	case TPMVersion12:
		return readPCR12(t.RWC, pcrIndex)
	case TPMVersion20:
		return readPCR20(t.RWC, pcrIndex)
	default:
		return nil, fmt.Errorf("unsupported TPM version: %x", t.Version)
	}
}

func (t *TPM) TakeOwnership(newAuth, newSRKAuth string) (bool, error) {
	switch t.Version {
	case TPMVersion12:
		return takeOwnership12(t.RWC, newAuth, newSRKAuth)
	case TPMVersion20:
		return takeOwnership20(t.RWC, newAuth, newSRKAuth)
	}
	return false, fmt.Errorf("unsupported TPM version: %x", t.Version)
}

func (t *TPM) ReadPubEK(ownerPW string) ([]byte, error) {
	switch t.Version {
	case TPMVersion12:
		return readPubEK12(t.RWC, ownerPW)
	case TPMVersion20:
		return readPubEK20(t.RWC, ownerPW)
	}
	return nil, fmt.Errorf("unsupported TPM version: %x", t.Version)
}

func (t *TPM) ResetLockValue(ownerPW string) (bool, error) {
	switch t.Version {
	case TPMVersion12:
		return resetLockValue12(t.RWC, ownerPW)
	case TPMVersion20:
		return resetLockValue20(t.RWC, ownerPW)
	}
	return false, fmt.Errorf("unsupported TPM version: %x", t.Version)
}

// NVReadValue reads a value from a given NVRAM index
// Type and byte order for TPM1.2 interface:
// (offset uint32)
// Type and byte oder for TPM2.0 interface:
// (authhandle uint32)
func (t *TPM) NVReadValue(index uint32, ownerPassword string, size, offhandle uint32) ([]byte, error) {
	switch t.Version {
	case TPMVersion12:
		return nvRead12(t.RWC, index, offhandle, size, ownerPassword)
	case TPMVersion20:
		return nvRead20(t.RWC, tpmutil.Handle(index), tpmutil.Handle(offhandle), ownerPassword, int(size))
	}
	return nil, fmt.Errorf("unsupported TPM version: %x", t.Version)
}
