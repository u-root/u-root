package tpm

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"

	tspi "github.com/google/go-tpm/tpm"
)

// TPM1 represents a TPM 1.2 device
type TPM1 struct {
	device  io.ReadWriteCloser
	tpmInfo Info
	// the following fields are used for unit testing
	// pcrReader emulates go-tpm's ReadPCR
	pcrReader func(io.ReadWriter, uint32) ([]byte, error)
}

const (
	// WellKnownSecret is the 20 bytes zero
	WellKnownSecret = ""
	// DefaultLocality is the TPM locality mostly used
	DefaultLocality byte = 0
)

// Info returns the TPMInfo object associated to this TPM device
func (t TPM1) Info() Info {
	return t.tpmInfo
}

// TakeOwnership takes ownership of the TPM. if no password defined use
// WELL_KNOWN_SECRET aka 20 zero bytes.
func (t *TPM1) TakeOwnership(ownerPassword string, srkPassword string) error {
	var ownerAuth [20]byte
	var srkAuth [20]byte

	if ownerPassword != "" {
		ownerAuth = sha1.Sum([]byte(ownerPassword))
	}

	if srkPassword != "" {
		srkAuth = sha1.Sum([]byte(srkPassword))
	}

	// This test assumes that the TPM has been cleared using OwnerClear.
	pubek, err := tspi.ReadPubEK(t.device)
	if err != nil {
		return err
	}

	return tspi.TakeOwnership(t.device, ownerAuth, srkAuth, pubek)
}

// Version returns the TPM version
func (t TPM1) Version() string {
	return TPM12
}

// ClearOwnership clears ownership of the TPM
func (t TPM1) ClearOwnership(ownerPassword string) error {
	var ownerAuth [20]byte

	if ownerPassword != "" {
		ownerAuth = sha1.Sum([]byte(ownerPassword))
	}

	return tspi.OwnerClear(t.device, ownerAuth)
}

// SetupTPM enabled, activates and takes
// the ownership of a TPM if it is not in a good
// state
func (t *TPM1) SetupTPM() error {
	if t.tpmInfo.Owned && t.tpmInfo.Specification == TPM12 {
		_, err := t.ReadPubEK(WellKnownSecret)
		if err != nil {
			t.ClearOwnership(WellKnownSecret)
			return err
		}
	}

	if !t.tpmInfo.Owned && t.tpmInfo.Enabled {
		if err := t.TakeOwnership(WellKnownSecret, WellKnownSecret); err != nil {
			return err
		}
	}

	if !t.tpmInfo.Enabled || !t.tpmInfo.Active || t.tpmInfo.TemporarilyDeactivated {
		return errors.New("TPM is not enabled")
	}
	return nil
}

// ReadPCR reads the PCR for the given index
func (t *TPM1) ReadPCR(pcr uint32) ([]byte, error) {
	data, err := t.pcrReader(t.device, pcr)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ReadPubEK reads the Public Endorsement Key part
func (t *TPM1) ReadPubEK(ownerPassword string) ([]byte, error) {
	var ownerAuth [20]byte
	if ownerPassword != "" {
		ownerAuth = sha1.Sum([]byte(ownerPassword))
	}

	ek, err := tspi.OwnerReadPubEK(t.device, ownerAuth)
	if err != nil {
		return nil, err
	}

	return ek, nil
}

// Measure hashes data and extends it into
// a TPM 1.2 PCR your choice.
func (t *TPM1) Measure(pcr uint32, data []byte) error {
	hash := sha1.Sum(data)

	if _, err := tspi.PcrExtend(t.device, pcr, hash); err != nil {
		return err
	}

	return nil
}

// SealData seals data at locality with pcrs and srkPassword
func (t *TPM1) SealData(locality byte, pcrs []int, data []byte, srkPassword string) ([]byte, error) {
	var srkAuth [20]byte
	if srkPassword != "" {
		srkAuth = sha1.Sum([]byte(srkPassword))
	}

	sealed, err := tspi.Seal(t.device, locality, pcrs, data, srkAuth[:])
	if err != nil {
		return nil, err
	}

	return sealed, nil
}

// ResealData seals data against a given pcrInfo map and srkPassword
// locality: TPM locality, by default zero.
// pcrInfo: A map of 24 entries. The key is the PCR index and the value is
// a hash.
// data: Data which should be sealed against the PCR of pcrInfo.
// srkPassword: The storage root key password of the TPM.
func (t *TPM1) ResealData(locality byte, pcrInfo map[int][]byte, data []byte, srkPassword string) ([]byte, error) {
	var srkAuth [20]byte
	if srkPassword != "" {
		srkAuth = sha1.Sum([]byte(srkPassword))
	}

	sealed, err := tspi.Reseal(t.device, locality, pcrInfo, data, srkAuth[:])
	if err != nil {
		return nil, err
	}

	return sealed, nil
}

// UnsealData unseals sealed data with srkPassword
func (t *TPM1) UnsealData(sealed []byte, srkPassword string) ([]byte, error) {
	var srkAuth [20]byte
	if srkPassword != "" {
		srkAuth = sha1.Sum([]byte(srkPassword))
	}

	unsealed, err := tspi.Unseal(t.device, sealed, srkAuth[:])
	if err != nil {
		return nil, err
	}
	return unsealed, err
}

// ResetLock resets the TPM brute force protection lock
func (t *TPM1) ResetLock(ownerPassword string) error {
	var ownerAuth [20]byte
	if ownerPassword != "" {
		ownerAuth = sha1.Sum([]byte(ownerPassword))
	}

	return tspi.ResetLockValue(t.device, ownerAuth)
}

// Close tpm device's file descriptor
func (t *TPM1) Close() {
	if t.device != nil {
		t.device.Close()
		t.device = nil
	}
}

// Summary returns a string with formatted TPM information
func (t TPM1) Summary() string {
	ret := ""
	ret += fmt.Sprintf("TPM Manufacturer:          %s\n", t.tpmInfo.Manufacturer)
	ret += fmt.Sprintf("TPM spec:                  %s\n", t.tpmInfo.Specification)
	ret += fmt.Sprintf("TPM owned:                 %t\n", t.tpmInfo.Owned)
	ret += fmt.Sprintf("TPM activated:             %t\n", t.tpmInfo.Active)
	ret += fmt.Sprintf("TPM enabled:               %t\n", t.tpmInfo.Enabled)
	ret += fmt.Sprintf("TPM temporary deactivated: %t\n", t.tpmInfo.TemporarilyDeactivated)
	return ret
}
