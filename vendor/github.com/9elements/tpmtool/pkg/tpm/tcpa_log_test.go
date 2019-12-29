package tpm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTPM12BiosEventLog(t *testing.T) {
	DefaultTCPABinaryLog = "tests/binary_bios_measurements_tpm12_bios"
	_, err := ParseLog(Bios, TPM12)
	require.NoError(t, err)
}

func TestParseTPM12UefiEventLog(t *testing.T) {
	DefaultTCPABinaryLog = "tests/binary_bios_measurements_tpm12_efi"
	_, err := ParseLog(Uefi, TPM12)
	require.NoError(t, err)
}

func TestParseTPM20UefiNonAgileEventLog(t *testing.T) {
	DefaultTCPABinaryLog = "tests/binary_bios_measurements_tpm20_efi_non_agile"
	_, err := ParseLog(Uefi, TPM20)
	require.NoError(t, err)
}

func TestParseTPM20UefiAgileEventLog(t *testing.T) {
	DefaultTCPABinaryLog = "tests/binary_bios_measurements_tpm20_efi_agile"
	_, err := ParseLog(Uefi, TPM20)
	require.NoError(t, err)
}
