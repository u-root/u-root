package crypto

import (
	"io/ioutil"
	"log"

	"github.com/systemboot/tpmtool/pkg/tpm"
)

const (
	// Blob type in PCR 7
	Blob uint32 = 7
	// BootConfig type in PCR 8
	BootConfig uint32 = 8
	// ConfigData type in PCR 8
	ConfigData uint32 = 8
	// NvramVars type in PCR 9
	NvramVars uint32 = 9
)

// TryMeasureBootConfig measures bootconfig contents
func TryMeasureBootConfig(name string, kernel string, initramfs string, kernelArgs string, deviceTree string) {
	TPMInterface, err := tpm.NewTPM()
	if err != nil {
		log.Printf("Cannot open TPM: %v", err)
		return
	}
	TryMeasureData(BootConfig, []byte(name), &name)
	TryMeasureData(BootConfig, []byte(kernel), &kernel)
	TryMeasureData(BootConfig, []byte(initramfs), &initramfs)
	TryMeasureData(BootConfig, []byte(kernelArgs), &kernelArgs)
	TryMeasureData(BootConfig, []byte(deviceTree), &deviceTree)
	TryMeasureFiles(kernel, initramfs, deviceTree)
	TPMInterface.Close()
}

// TryMeasureData measures a byte array with additional information
func TryMeasureData(pcr uint32, data []byte, info *string) {
	TPMInterface, err := tpm.NewTPM()
	if err != nil {
		log.Printf("Cannot open TPM: %v", err)
		return
	}
	log.Println("Measuring blob: " + *info)
	TPMInterface.Measure(pcr, data)
	TPMInterface.Close()
}

// TryMeasureFiles measures a variable amount of files
func TryMeasureFiles(files ...string) {
	TPMInterface, err := tpm.NewTPM()
	if err != nil {
		log.Printf("Cannot open TPM: %v", err)
		return
	}
	for _, file := range files {
		log.Println("Measuring file: " + file)
		data, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		TPMInterface.Measure(Blob, data)
	}
	TPMInterface.Close()
}
