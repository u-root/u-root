package tpm

// Manual testing needed or tpm hardware in VM
/*
import (
	"crypto/sha1"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testReadPcrIndex  uint32 = 23
	testWritePcrIndex uint32 = 16
	testString        string = "teststring"
)

// Beware, TPM testa are a horrible idea because of state transitions
// for PCR. Powercycles are not possible due security architecture.
// Tests must be run as root and follow specific execution flow.
// Also do not run test on a production system if the tpm is used elsewhere.

func TestReadPcrTPM1(t *testing.T) {
	pcrData, err := ReadPcrTPM1(testReadPcrIndex)
	require.NoError(t, err)
	require.Equal(t, pcrData, make([]byte, 20))
}

func TestMeasureTPM1(t *testing.T) {
	oldPcrValue, err := ReadPcrTPM1(testWritePcrIndex)
	require.NoError(t, err)

	pcrValue := sha1.Sum([]byte(testString))
	finalPcr := sha1.Sum(append(oldPcrValue, pcrValue[:]...))

	err = MeasureTPM1(testWritePcrIndex, []byte(testString))
	require.NoError(t, err)

	newPcrValue, err := ReadPcrTPM1(testWritePcrIndex)
	require.NoError(t, err)

	require.Equal(t, finalPcr[:], newPcrValue)
}

func TestOwnerClearTPM1(t *testing.T) {
	err := OwnerClearTPM1("keins")
	require.NoError(t, err)
}

func TestTakeOwnershipTPM1(t *testing.T) {
	err := TakeOwnershipTPM1("", "")
	require.NoError(t, err)
}
*/
