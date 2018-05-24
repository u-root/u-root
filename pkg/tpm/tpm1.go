package tpm

import (
	"crypto/sha1"

	tspi "github.com/google/go-tpm/tpm"
)

// OwnerClearTPM1 clears the TPM and destorys all
// access to existing keys. Afterwards a machine
// power cycle is needed.
func (handle *TPM) OwnerClearTPM1(ownerPassword string) error {
	var ownerAuth [20]byte

	if ownerPassword != "" {
		ownerAuth = sha1.Sum([]byte(ownerPassword))
	}

	return tspi.OwnerClear(handle.device, ownerAuth)
}

// TakeOwnershipTPM1 takes ownership of the TPM. if no password defined use
// WELL_KNOWN_SECRET aka 20 zero bytes.
func (handle *TPM) TakeOwnershipTPM1(ownerPassword string, srkPassword string) error {
	var ownerAuth [20]byte
	var srkAuth [20]byte

	if ownerPassword != "" {
		ownerAuth = sha1.Sum([]byte(ownerPassword))
	}

	if srkPassword != "" {
		srkAuth = sha1.Sum([]byte(srkPassword))
	}

	// This test assumes that the TPM has been cleared using OwnerClear.
	pubek, err := tspi.ReadPubEK(handle.device)
	if err != nil {
		return err
	}

	return tspi.TakeOwnership(handle.device, ownerAuth, srkAuth, pubek)
}

// ReadPcrTPM1 reads the PCR for the given
// index
func (handle *TPM) ReadPcrTPM1(pcr uint32) ([]byte, error) {
	data, err := tspi.ReadPCR(handle.device, pcr)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ReadPubEKTPM1 reads the public Endorsement key part
func (handle *TPM) ReadPubEKTPM1(ownerPassword string) ([]byte, error) {
	var ownerAuth [20]byte

	if ownerPassword != "" {
		ownerAuth = sha1.Sum([]byte(ownerPassword))
	}

	ek, err := tspi.OwnerReadPubEK(handle.device, ownerAuth)
	if err != nil {
		return nil, err
	}

	return ek, nil
}

// MeasureTPM1 hashes data and extends it into
// a TPM 1.2 PCR your choice.
func (handle *TPM) MeasureTPM1(pcr uint32, data []byte) error {
	hash := sha1.Sum(data)

	if _, err := tspi.PcrExtend(handle.device, pcr, hash); err != nil {
		return err
	}

	return nil
}
