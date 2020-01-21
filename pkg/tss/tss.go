// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tss provides TPM 1.2/2.0 core functionality and
// abstraction layer for high-level functions
package tss

import (
	"crypto"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/awnumar/memguard"
	"github.com/google/go-tpm/tpm2"
)

// OpenTPM initializes access to the TPM based on the
// config provided.
func OpenTPM() (*TPM, error) {
	candidateTPMs, err := probeSystemTPMs()
	if err != nil {
		return nil, err
	}

	for _, tpm := range candidateTPMs {
		tss, err := openTPM(tpm)
		if err != nil {
			continue
		}
		return tss, nil
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
	memguard.Purge()
	return t.RWC.Close()
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

// String returns a human-friendly representation of the hash algorithm.
func (a HashAlg) String() string {
	switch a {
	case HashSHA1:
		return "SHA1"
	case HashSHA256:
		return "SHA256"
	}
	return fmt.Sprintf("HashAlg<%d>", int(a))
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
func (t *TPM) Extend(hash []byte, pcrIndex uint32, alg HashAlg) error {
	switch t.Version {
	case TPMVersion12:
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
		err := extendPCR20(t.RWC, pcrIndex, hash, alg)
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
		err := extendPCR20(t.RWC, pcrIndex, hash, alg)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported TPM version: %x", t.Version)
	}

	return nil
}

// ReadPCR reads a single PCR value by defining the pcrIndex
func (t *TPM) ReadPCR(pcrIndex uint32, alg HashAlg) ([]byte, error) {
	switch t.Version {
	case TPMVersion12:
		return readPCR12(t.RWC, pcrIndex)
	case TPMVersion20:
		return readPCR20(t.RWC, pcrIndex, alg)
	default:
		return nil, fmt.Errorf("unsupported TPM version: %x", t.Version)
	}
}
