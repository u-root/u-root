package tpm

// TPM is an interface that both TPM1 and TPM2 have to implement. It requires a
// common subset of methods that both TPM versions have to implement.
// Version-specific methods have to be implemented in the relevant object.
type TPM interface {
	Version() string
	SetupTPM() error
	TakeOwnership() error
	ClearOwnership() error
	Measure(pcr uint32, data []byte) error
	Info() string
	Close()
	ReadPCR(uint32) ([]byte, error)
}
