package tpm

import (
	"bytes"
	"os"
	"strconv"
)

// TPM is an interface that both TPM1 and TPM2 have to implement. It requires a
// common subset of methods that both TPM versions have to implement.
// Version-specific methods have to be implemented in the relevant object.
type TPM interface {
	Info() Info
	Summary() string
	Version() string
	SetupTPM() error
	TakeOwnership() error
	ClearOwnership() error
	Measure(pcr uint32, data []byte) error
	Close()
	ReadPCR(uint32) ([]byte, error)
}

// Info holds information about a TPM device
type Info struct {
	Specification          string
	Owned                  bool
	Active                 bool
	Enabled                bool
	TemporarilyDeactivated bool
}

// getInfo reads TPM information from various TPM state devices and returns them
// wrapped in an Info structure
func getInfo() (*Info, error) {
	var cap [256]byte
	var owned [1]byte
	var active [1]byte
	var enabled [1]byte
	var tempDeactivated [1]byte

	caps, err := os.Open(TpmCapabilities)
	if err != nil {
		return nil, err
	}
	defer caps.Close()

	ownedState, err := os.Open(TpmOwnershipState)
	if err != nil {
		return nil, err
	}
	defer ownedState.Close()

	activeState, err := os.Open(TpmActivatedState)
	if err != nil {
		return nil, err
	}
	defer activeState.Close()

	enabledState, err := os.Open(TpmEnabledState)
	if err != nil {
		return nil, err
	}
	defer enabledState.Close()

	tempDeactivatedState, err := os.Open(TpmTempDeactivatedState)
	if err != nil {
		return nil, err
	}
	defer tempDeactivatedState.Close()

	caps.Read(cap[:])
	specBytes := bytes.Split(cap[:], []byte(specFilter))
	specBytes = bytes.Split(specBytes[1], []byte("\n"))

	ownedState.Read(owned[:])
	activeState.Read(active[:])
	enabledState.Read(enabled[:])
	tempDeactivatedState.Read(tempDeactivated[:])

	spec := string(specBytes[0])
	ownedBool, _ := strconv.ParseBool(string(owned[:]))
	activeBool, _ := strconv.ParseBool(string(active[:]))
	enabledBool, _ := strconv.ParseBool(string(enabled[:]))
	tempDeactivatedBool, _ := strconv.ParseBool(string(tempDeactivated[:]))

	tinfo := Info{
		Specification:          spec,
		Owned:                  ownedBool,
		Active:                 activeBool,
		Enabled:                enabledBool,
		TemporarilyDeactivated: tempDeactivatedBool,
	}

	return &tinfo, nil
}
