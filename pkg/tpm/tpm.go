package tpm

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
