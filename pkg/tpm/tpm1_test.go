package tpm

// Manual testing needed or tpm hardware in VM
import (
	"io"
	"os"
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

func TestTPM1NewTPM(t *testing.T) {
	// TODO use a fake TPM. Unfortunately go-tpm checks if it's a device or a
	// socket, and fails otherwise. Need to use github.com/stefanberger/swtpm or
	// mock go-tpm entirely
	TPMOpener = func(string) (io.ReadWriteCloser, error) {
		fd, err := os.Open("tests/fake_tpm")
		if err != nil {
			return nil, err
		}
		return io.ReadWriteCloser(fd), nil
	}
	TpmCapabilities = "tests/fake_caps_tpm12"
	TpmOwnershipState = "tests/fake_owned_1"
	TpmActivatedState = "tests/fake_active_1"
	TpmEnabledState = "tests/fake_enabled_1"
	TpmTempDeactivatedState = "tests/fake_temp_deactivated_0"

	tpm, err := NewTPM()
	require.NoError(t, err)
	require.Equal(t, tpm12, tpm.Version())
}

func TestTPM1ReadPcr(t *testing.T) {
	tpm, err := NewTPM()
	tpm.(*TPM1).pcrReader = func(io.ReadWriter, uint32) ([]byte, error) {
		return make([]byte, 20), nil
	}
	require.NoError(t, err)
	pcrData, err := tpm.ReadPCR(testReadPcrIndex)
	require.NoError(t, err)
	require.Equal(t, pcrData, make([]byte, 20))
}

/*
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
