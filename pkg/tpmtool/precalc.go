package tpmtool

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"io/ioutil"
	"os"

	"github.com/koding/multiconfig"
	"github.com/systemboot/tpmtool/pkg/tpm"
)

// TPMInterface is a global TPM interface
var TPMInterface tpm.TPM

// CurrentPCRMap is the current used PCR map and a copy of the default map
var CurrentPCRMap map[int][]byte

// TPM1DefaultPCRMap is the TPM 1.2 default PCR map after a power cycle without
// any measurements done
var TPM1DefaultPCRMap = map[int][]byte{
	0:  make([]byte, 20),
	1:  make([]byte, 20),
	2:  make([]byte, 20),
	3:  make([]byte, 20),
	4:  make([]byte, 20),
	5:  make([]byte, 20),
	6:  make([]byte, 20),
	7:  make([]byte, 20),
	8:  make([]byte, 20),
	9:  make([]byte, 20),
	10: make([]byte, 20),
	11: make([]byte, 20),
	12: make([]byte, 20),
	13: make([]byte, 20),
	14: make([]byte, 20),
	15: make([]byte, 20),
	16: make([]byte, 20),
	17: []byte{'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f'},
	18: []byte{'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f'},
	19: []byte{'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f'},
	20: []byte{'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f'},
	21: []byte{'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f'},
	22: []byte{'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f'},
	23: make([]byte, 20),
}

func getMap() map[int][]byte {
	if TPMInterface.Info().Specification == tpm.TPM12 {
		return TPM1DefaultPCRMap
	}

	return nil
}

func hashSum(data []byte, algoID tpm.IAlgHash) ([]byte, error) {
	switch algoID {
	case tpm.TPMAlgSha:
		hash := sha1.Sum(data)
		return hash[:], nil
	case tpm.TPMAlgSha256:
		hash := sha256.Sum256(data)
		return hash[:], nil
	case tpm.TPMAlgSha384:
		hash := sha512.Sum384(data)
		return hash[:], nil
	case tpm.TPMAlgSha512:
		hash := sha512.Sum512(data)
		return hash[:], nil
	case tpm.TPMAlgSm3s256:
		return nil, errors.New("Not implemented yet")
	}

	return nil, errors.New("Hash algorithm not implemented yet")
}

// StaticPCR populates a static PCR into the map
func StaticPCR(pcrIndex int, hash []byte) {
	CurrentPCRMap[pcrIndex] = hash
}

// DynamicPCR gets the current PCR and populates it into the map
func DynamicPCR(pcrIndex int) error {
	hash, err := TPMInterface.ReadPCR(uint32(pcrIndex))
	if err != nil {
		return err
	}

	CurrentPCRMap[pcrIndex] = hash
	return nil
}

// ExtendPCR extends a hash into a current PCR
func ExtendPCR(pcrIndex int, hash []byte, algoID tpm.IAlgHash) error {
	hash, err := hashSum(append(CurrentPCRMap[pcrIndex], hash...), algoID)
	if err != nil {
		return err
	}

	CurrentPCRMap[pcrIndex] = hash
	return nil
}

// MeasurePCR measures a file into a PCR
func MeasurePCR(pcrIndex int, filePath string, algoID tpm.IAlgHash) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var fileHash []byte
	fileHash, err = hashSum(file, algoID)
	if err != nil {
		return err
	}

	hash, err := hashSum(append(CurrentPCRMap[pcrIndex], fileHash...), algoID)
	if err != nil {
		return err
	}

	CurrentPCRMap[pcrIndex] = hash
	return nil
}

// FirmwareLogPCR uses the firmware ACPI log for extending PCRs
func FirmwareLogPCR(pcrIndex int, firmware FirmwareType) error {
	tpmSpec := TPMInterface.Info().Specification
	tcpaLog, err := tpm.ParseLog(string(firmware), tpmSpec)
	if err != nil {
		return err
	}

	for _, event := range tcpaLog.PcrList {
		if event.PcrIndex == pcrIndex {
			for _, digest := range event.Digests {
				if event.PcrIndex == pcrIndex {
					err := ExtendPCR(event.PcrIndex, digest.Digest, digest.DigestAlg)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// LuksPCR extends the hash of a LUKS device into a current PCR
func LuksPCR(pcrIndex int, devicePath string, algoID tpm.IAlgHash) error {
	deviceFD, err := os.Open(devicePath)
	if err != nil {
		return err
	}
	defer deviceFD.Close()

	luksHeader := make([]byte, Luks1HeaderLength)
	_, err = deviceFD.Read(luksHeader)
	if err != nil {
		return err
	}

	hash, err := hashSum(append(CurrentPCRMap[pcrIndex], luksHeader...), algoID)
	if err != nil {
		return err
	}

	CurrentPCRMap[pcrIndex] = hash
	return nil
}

func runCalculations(calculations []PreCalculation, pcrIndex int) error {
	for _, calculation := range calculations {
		if calculation.Method == Static && len(calculations) > 1 {
			return errors.New("Static type: More calculation defined than possible")
		}

		if calculation.Method == Dynamic && len(calculations) > 1 {
			return errors.New("Dynamic type: More calculation defined than possible")
		}
	}

	// For TPM 2.0 we need to get Algo per PCR bank
	var algoID = tpm.TPMAlgError
	if TPMInterface.Info().Specification == tpm.TPM12 {
		algoID = tpm.TPMAlgSha
	}

	CurrentPCRMap[pcrIndex] = getMap()[pcrIndex]
	for _, calculation := range calculations {
		switch calculation.Method {
		case Static:
			hash := calculation.Hash
			if hash == "" {
				return errors.New("Static type: No hash defined")
			}
			StaticPCR(pcrIndex, []byte(hash))
		case Dynamic:
			return DynamicPCR(pcrIndex)
		case Extend:
			if len(calculation.Hashes) <= 0 {
				return errors.New("Extend type: No hashes defined")
			}
			for _, hash := range calculation.Hashes {
				if err := ExtendPCR(pcrIndex, []byte(hash), algoID); err != nil {
					return err
				}
			}
		case Measure:
			if len(calculation.FilePaths) <= 0 {
				return errors.New("Measure type: No paths defined")
			}
			for _, path := range calculation.FilePaths {
				if err := MeasurePCR(pcrIndex, path, algoID); err != nil {
					return err
				}
			}
		case FirmwareLog:
			if calculation.Firmware == "" {
				return errors.New("FirmwareLog type: Firmware not set")
			}
			if err := FirmwareLogPCR(pcrIndex, calculation.Firmware); err != nil {
				return err
			}
		case Luks:
			if calculation.DevicePath == "" {
				return errors.New("Luks type: No path defined")
			}
			return LuksPCR(pcrIndex, calculation.DevicePath, algoID)
		default:
			return errors.New(string(calculation.Method) + " not implemented")
		}
	}

	return nil
}

func executeConfig(sealingConfig *TPM1SealingConfig) error {
	if sealingConfig.Pcr0 != nil {
		if err := runCalculations(sealingConfig.Pcr0, 0); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr1 != nil {
		if err := runCalculations(sealingConfig.Pcr1, 1); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr2 != nil {
		if err := runCalculations(sealingConfig.Pcr2, 2); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr3 != nil {
		if err := runCalculations(sealingConfig.Pcr3, 3); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr4 != nil {
		if err := runCalculations(sealingConfig.Pcr4, 4); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr5 != nil {
		if err := runCalculations(sealingConfig.Pcr5, 5); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr6 != nil {
		if err := runCalculations(sealingConfig.Pcr6, 6); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr7 != nil {
		if err := runCalculations(sealingConfig.Pcr7, 7); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr8 != nil {
		if err := runCalculations(sealingConfig.Pcr8, 8); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr9 != nil {
		if err := runCalculations(sealingConfig.Pcr9, 9); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr10 != nil {
		if err := runCalculations(sealingConfig.Pcr10, 10); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr11 != nil {
		if err := runCalculations(sealingConfig.Pcr11, 11); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr12 != nil {
		if err := runCalculations(sealingConfig.Pcr12, 12); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr13 != nil {
		if err := runCalculations(sealingConfig.Pcr13, 13); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr14 != nil {
		if err := runCalculations(sealingConfig.Pcr14, 14); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr15 != nil {
		if err := runCalculations(sealingConfig.Pcr15, 15); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr16 != nil {
		if err := runCalculations(sealingConfig.Pcr16, 16); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr17 != nil {
		if err := runCalculations(sealingConfig.Pcr17, 17); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr18 != nil {
		if err := runCalculations(sealingConfig.Pcr18, 18); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr19 != nil {
		if err := runCalculations(sealingConfig.Pcr19, 19); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr20 != nil {
		if err := runCalculations(sealingConfig.Pcr20, 20); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr21 != nil {
		if err := runCalculations(sealingConfig.Pcr21, 21); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr22 != nil {
		if err := runCalculations(sealingConfig.Pcr22, 22); err != nil {
			return err
		}
	}

	if sealingConfig.Pcr23 != nil {
		if err := runCalculations(sealingConfig.Pcr23, 23); err != nil {
			return err
		}
	}

	return nil
}

// PreCalculate calculates a PCR map by a given sealing configuration
// doing different types of calculations in the right order
func PreCalculate(tpmInterface tpm.TPM, sealingConfigPath string) (map[int][]byte, error) {
	TPMInterface = tpmInterface

	// Initialize the default values
	var sealingConf *TPM1SealingConfig
	CurrentPCRMap = make(map[int][]byte)
	if TPMInterface.Info().Specification == tpm.TPM12 {
		sealingConf = new(TPM1SealingConfig)
	} else {
		return nil, errors.New("TPM spec not implemented yet")
	}

	config := multiconfig.NewWithPath(sealingConfigPath)
	if config == nil {
		return nil, errors.New("Couldn't load config from disk")
	}

	if err := config.Load(sealingConf); err != nil {
		return nil, err
	}

	if err := executeConfig(sealingConf); err != nil {
		return nil, err
	}

	return CurrentPCRMap, nil
}
