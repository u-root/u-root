package tpm

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	tspi "github.com/google/go-tpm/tpm"
)

var (
	// TPMOpener is used to allow unit testing
	TPMOpener = tspi.OpenTPM

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
)

const (
	// TPM12 is the TPM 1.2 identifier
	TPM12 = "1.2"
	// TPM12MaxKeySize is the TPM 1.2 maximum key size
	TPM12MaxKeySize = 256
	// TPM20 is the TPM 2.0 identifier
	TPM20 = "2.0"
)

// Manufactures list of TPM vendors
var Manufactures = map[string]string{
	"0x53544d20": "STMicroelectronics",
}

// TPM is an interface that both TPM1 and TPM2 have to implement. It requires a
// common subset of methods that both TPM versions have to implement.
// Version-specific methods have to be implemented in the relevant object.
type TPM interface {
	Info() Info
	Summary() string
	Version() string
	SetupTPM() error
	TakeOwnership(ownerPassword string, srkPassword string) error
	ClearOwnership(ownerPassword string) error
	Measure(pcr uint32, data []byte) error
	Close()
	ReadPCR(uint32) ([]byte, error)
	ReadPubEK(ownerPassword string) ([]byte, error)
	SealData(locality byte, pcrs []int, data []byte, srkPassword string) ([]byte, error)
	ResealData(locality byte, pcrInfo map[int][]byte, data []byte, srkPassword string) ([]byte, error)
	UnsealData(sealed []byte, srkPassword string) ([]byte, error)
	ResetLock(ownerPassword string) error
}

// Info holds information about a TPM device
type Info struct {
	Manufacturer           string
	Specification          string
	Owned                  bool
	Active                 bool
	Enabled                bool
	TemporarilyDeactivated bool
}

func bytesToBool(data []byte) (bool, error) {
	s := strings.TrimSuffix(string(data), "\n")
	return strconv.ParseBool(s)
}

// NewTPM gets a new TPM handle struct with
// io fd and specification string
func NewTPM() (TPM, error) {
	// It's the caller's responsibility to call TPM.Close()
	rwc, err := TPMOpener(TPMDevice)
	if err != nil {
		return nil, err
	}

	tinfo, err := getInfo()
	if err != nil {
		return nil, err
	}

	if tinfo.Specification == TPM12 {
		return &TPM1{
			device:    rwc,
			tpmInfo:   *tinfo,
			pcrReader: tspi.ReadPCR,
		}, nil
	} else if tinfo.Specification == TPM20 {
		return nil, errors.New("TPM 2.0 not supported yet")
	} else if tinfo.Specification == "" {
		return nil, fmt.Errorf("Invalid empty TPM specification")
	}
	return nil, fmt.Errorf("Unknown TPM specification: %s", tinfo.Specification)
}

// getInfo reads TPM information from various TPM state devices and returns them
// wrapped in an Info structure
func getInfo() (*Info, error) {
	caps, err := ioutil.ReadFile(TpmCapabilities)
	if err != nil {
		return nil, err
	}

	ownedBytes, err := ioutil.ReadFile(TpmOwnershipState)
	if err != nil {
		return nil, err
	}

	activeBytes, err := ioutil.ReadFile(TpmActivatedState)
	if err != nil {
		return nil, err
	}

	enabledBytes, err := ioutil.ReadFile(TpmEnabledState)
	if err != nil {
		return nil, err
	}

	tempDeactivatedBytes, err := ioutil.ReadFile(TpmTempDeactivatedState)
	if err != nil {
		return nil, err
	}

	manufacturerPrefix := "Manufacturer: "
	var manufacturerID string
	for _, lineBytes := range bytes.Split(caps, []byte{'\n'}) {
		line := string(lineBytes)
		if strings.HasPrefix(line, manufacturerPrefix) {
			manufacturerID = line[len(manufacturerPrefix):]
		}
	}

	manufacturer := Manufactures[manufacturerID]
	if manufacturer == "" {
		manufacturer = "< unknown >"
	}

	specPrefix := "TCG version: "
	var tpmVersion string
	for _, lineBytes := range bytes.Split(caps, []byte{'\n'}) {
		line := string(lineBytes)
		if strings.HasPrefix(line, specPrefix) {
			tpmVersion = line[len(specPrefix):]
		}
	}

	owned, err := bytesToBool(ownedBytes)
	if err != nil {
		return nil, err
	}
	active, err := bytesToBool(activeBytes)
	if err != nil {
		return nil, err
	}
	enabled, err := bytesToBool(enabledBytes)
	if err != nil {
		return nil, err
	}
	tempDeactivated, err := bytesToBool(tempDeactivatedBytes)
	if err != nil {
		return nil, err
	}

	tinfo := Info{
		Manufacturer:           manufacturer,
		Specification:          tpmVersion,
		Owned:                  owned,
		Active:                 active,
		Enabled:                enabled,
		TemporarilyDeactivated: tempDeactivated,
	}

	return &tinfo, nil
}
