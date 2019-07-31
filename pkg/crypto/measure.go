// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crypto

import (
	"io/ioutil"
	"log"

	"github.com/systemboot/tpmtool/pkg/tpm"
)

const (
	// BlobPCR type in PCR 7
	BlobPCR uint32 = 7
	// BootConfigPCR type in PCR 8
	BootConfigPCR uint32 = 8
	// ConfigDataPCR type in PCR 8
	ConfigDataPCR uint32 = 8
	// NvramVarsPCR type in PCR 9
	NvramVarsPCR uint32 = 9
)

// TryMeasureData measures a byte array with additional information
func TryMeasureData(pcr uint32, data []byte, info string) {
	TPMInterface, err := tpm.NewTPM()
	if err != nil {
		log.Printf("Cannot open TPM: %v", err)
		return
	}
	log.Printf("Measuring blob: %v", info)
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
		log.Printf("Measuring file: %v", file)
		data, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		TPMInterface.Measure(BlobPCR, data)
	}
	TPMInterface.Close()
}
