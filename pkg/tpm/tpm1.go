package tpm

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	tspi "github.com/google/go-tpm/tpm"
)

// TPM1 represents a TPM 1.2 device
type TPM1 struct {
	device          io.ReadWriteCloser
	specification   string
	owned           bool
	active          bool
	enabled         bool
	tempDeactivated bool
	ownerPassword   string
	srkPassword     string
}

const (
	// TPMDevice main device path for
	// TSS usage
	TPMDevice = "/dev/tpm0"

	// TpmCapabilities for selecting tpm spec
	TpmCapabilities = "/sys/class/tpm/tpm0/caps"

	// TpmOwnershipState contains owner state
	TpmOwnershipState = "/sys/class/tpm/tpm0/owned"

	// TpmActivatedState contains active state
	TpmActivatedState = "/sys/class/tpm/tpm0/active"

	// TpmEnabledState contains enabled state
	TpmEnabledState = "/sys/class/tpm/tpm0/enabled"

	// TpmTempDeactivatedState contains enabled state
	TpmTempDeactivatedState = "/sys/class/tpm/tpm0/temp_deactivated"

	tpm12      = "1.2"
	tpm20      = "2.0"
	specFilter = "TCG version: "
)

var (
	wellKnownSecret string
)

func getInfo() (string, bool, bool, bool, bool, error) {
	var cap [256]byte
	var owned [1]byte
	var active [1]byte
	var enabled [1]byte
	var tempDeactivated [1]byte

	caps, err := os.Open(TpmCapabilities)
	if err != nil {
		return "", false, false, false, false, err
	}

	ownedState, err := os.Open(TpmOwnershipState)
	if err != nil {
		return "", false, false, false, false, err
	}

	activeState, err := os.Open(TpmActivatedState)
	if err != nil {
		return "", false, false, false, false, err
	}

	enabledState, err := os.Open(TpmEnabledState)
	if err != nil {
		return "", false, false, false, false, err
	}

	tempDeactivatedState, err := os.Open(TpmTempDeactivatedState)
	if err != nil {
		return "", false, false, false, false, err
	}

	caps.Read(cap[:])
	specBytes := bytes.Split(cap[:], []byte(specFilter))
	specBytes = bytes.Split(specBytes[1], []byte("\n"))

	ownedState.Read(owned[:])
	activeState.Read(active[:])
	enabledState.Read(enabled[:])
	tempDeactivatedState.Read(tempDeactivated[:])

	caps.Close()
	ownedState.Close()
	activeState.Close()
	enabledState.Close()
	tempDeactivatedState.Close()

	spec := string(specBytes[0])
	ownedBool, _ := strconv.ParseBool(string(owned[:]))
	activeBool, _ := strconv.ParseBool(string(active[:]))
	enabledBool, _ := strconv.ParseBool(string(enabled[:]))
	tempDeactivatedBool, _ := strconv.ParseBool(string(tempDeactivated[:]))

	return spec, ownedBool, activeBool, enabledBool, tempDeactivatedBool, nil
}

// NewTPM gets a new TPM handle struct with
// io fd and specification string
func NewTPM() (TPM, error) {
	rwc, err := tspi.OpenTPM(TPMDevice)
	if err != nil {
		return nil, err
	}

	// No error checking for spec because of tpm 1.2
	// capability command not being available in deacitvated
	// or disabled state.
	spec, owned, active, enabled, tempDeactivated, err := getInfo()

	if err != nil {
		return nil, err
	}
	if spec == tpm12 {
		return &TPM1{device: rwc, specification: spec, owned: owned, active: active, enabled: enabled, tempDeactivated: tempDeactivated}, nil
	} else if spec == tpm20 {
		return nil, errors.New("TPM 2.0 not supported yet")
	} else {
		return nil, fmt.Errorf("Unknown TPM specification: %s", spec)
	}
}

// OwnerClear clears the TPM and destorys all
// access to existing keys. Afterwards a machine
// power cycle is needed.
func (t *TPM1) OwnerClear(ownerPassword string) error {
	var ownerAuth [20]byte

	if ownerPassword != "" {
		ownerAuth = sha1.Sum([]byte(ownerPassword))
	}

	return tspi.OwnerClear(t.device, ownerAuth)
}

// TakeOwnership takes ownership of the TPM. if no password defined use
// WELL_KNOWN_SECRET aka 20 zero bytes.
func (t *TPM1) TakeOwnership() error {
	var ownerAuth [20]byte
	var srkAuth [20]byte

	if t.ownerPassword != "" {
		ownerAuth = sha1.Sum([]byte(t.ownerPassword))
	}

	if t.srkPassword != "" {
		srkAuth = sha1.Sum([]byte(t.srkPassword))
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
	return tpm12
}

// ClearOwnership clears ownership of the TPM
func (t TPM1) ClearOwnership() error {
	var err error
	if t.specification == tpm12 {
		err = t.OwnerClear(t.ownerPassword)
	}
	return err
}

// SetupTPM enabled, activates and takes
// the ownership of a TPM if it is not in a good
// state
func (t *TPM1) SetupTPM() error {
	if t.owned && t.specification == tpm12 {
		_, err := t.ReadPubEK(wellKnownSecret)
		if err != nil {
			t.ClearOwnership()
			return err
		}
	}

	if !t.owned && t.enabled {
		if err := t.TakeOwnership(); err != nil {
			return err
		}
	}

	if !t.enabled || !t.active || t.tempDeactivated {
		//utils.Die(true, "Please enable the TPM")
	}
	return nil
}

// ReadPCR reads the PCR for the given
// index
func (t *TPM1) ReadPCR(pcr uint32) ([]byte, error) {
	data, err := tspi.ReadPCR(t.device, pcr)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ReadPubEK reads the public Endorsement key part
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

// Close tpm device's file descriptor
func (t *TPM1) Close() {
	if t.device != nil {
		t.device.Close()
		t.device = nil
	}
}

// Info returns TPM information
func (t TPM1) Info() string {
	ret := ""
	ret += fmt.Sprintf("TPM spec:                  %s\n", t.specification)
	ret = fmt.Sprintf("TPM owned:                 %t\n", t.owned)
	ret += fmt.Sprintf("TPM activated:             %t\n", t.active)
	ret += fmt.Sprintf("TPM enabled:               %t\n", t.enabled)
	ret += fmt.Sprintf("TPM temporary deactivated: %t\n", t.tempDeactivated)
	return ret
}
