package tpmtool

// PreCalculation structure
type PreCalculation struct {
	Method     CalculateType
	Hash       string
	DevicePath string
	Firmware   FirmwareType
	Hashes     []string
	FilePaths  []string
}

// TPM1SealingConfig is a TPM1 sealing configuration
type TPM1SealingConfig struct {
	Pcr0  []PreCalculation
	Pcr1  []PreCalculation
	Pcr2  []PreCalculation
	Pcr3  []PreCalculation
	Pcr4  []PreCalculation
	Pcr5  []PreCalculation
	Pcr6  []PreCalculation
	Pcr7  []PreCalculation
	Pcr8  []PreCalculation
	Pcr9  []PreCalculation
	Pcr10 []PreCalculation
	Pcr11 []PreCalculation
	Pcr12 []PreCalculation
	Pcr13 []PreCalculation
	Pcr14 []PreCalculation
	Pcr15 []PreCalculation
	Pcr16 []PreCalculation
	Pcr17 []PreCalculation
	Pcr18 []PreCalculation
	Pcr19 []PreCalculation
	Pcr20 []PreCalculation
	Pcr21 []PreCalculation
	Pcr22 []PreCalculation
	Pcr23 []PreCalculation
}
