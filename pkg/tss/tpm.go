package tss

import (
	"errors"
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
