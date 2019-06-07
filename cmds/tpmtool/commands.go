package main

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"syscall"
	"time"

	"github.com/systemboot/systemboot/pkg/storage"
	"github.com/systemboot/tpmtool/pkg/tpm"
	"github.com/systemboot/tpmtool/pkg/tpmtool"
)

const (
	// MaxPlatformConfigurationRegister is the maximum number of PCRs
	MaxPlatformConfigurationRegister = 24
	// DefaultFilePermissions is the default write permission
	DefaultFilePermissions = 660
	// LinuxEFIFirmwareDir is the UEFI linux firmware directory
	LinuxEFIFirmwareDir = "/sys/firmware/efi"
	// Delay is used for sealing operations delay
	Delay = 900
)

// Status Dumps the tpm status
func Status() error {
	summary := TPMInterface.Summary()
	fmt.Print(summary)

	if TPMInterface.Info().TemporarilyDeactivated {
		fmt.Println("\nError: Check your BIOS! TPM is temporary deactivated.")
	}

	if (!TPMInterface.Info().Active || !TPMInterface.Info().Enabled) && !TPMInterface.Info().TemporarilyDeactivated {
		fmt.Println("\nError: TPM is inactive or disabled! Check your BIOS physical presence settings.")
	}

	if !TPMInterface.Info().Owned {
		fmt.Println("\nError: TPM is not owned! Please take ownership of the TPM.")
	}

	return nil
}

// Ek dumps the Endorsement Key
func Ek() error {
	var pubEk []byte
	pubEk, err := TPMInterface.ReadPubEK(*ekCommandPassword)
	if err != nil {
		return err
	}

	if *ekCommandOutfile != "" {
		if err := ioutil.WriteFile(*ekCommandOutfile, pubEk, 660); err != nil {
			return err
		}
	}
	fingerprint := sha256.Sum256(pubEk)
	fmt.Printf("EK Pubkey fingerprint: 0x%x\n", fingerprint)

	return nil
}

// OwnerTake takes ownership of the TPM
func OwnerTake() error {
	err := TPMInterface.TakeOwnership(*ownerCommandPassword, *ownerCommandTakeSrkPassword)
	if err != nil {
		return err
	}

	return nil
}

// OwnerClear clears ownership of the TPM
func OwnerClear() error {
	err := TPMInterface.ClearOwnership(*ownerCommandPassword)
	if err != nil {
		return err
	}

	return nil
}

// OwnerResetLock resets the TPM bruteforce lock
func OwnerResetLock() error {
	err := TPMInterface.ResetLock(*ownerCommandPassword)
	if err != nil {
		return err
	}

	return nil
}

// CryptoSeal seals data aganst PCR with TPM
func CryptoSeal() error {
	// No shit sherlock, we are too fast Oo (go-tpm bug ?)
	time.Sleep(Delay * time.Millisecond)

	plainText, err := ioutil.ReadFile(*cryptoCommandSealPlainFile)
	if err != nil {
		return err
	}

	if TPMSpecVersion == tpm.TPM12 && len(plainText) > tpm.TPM12MaxKeySize {
		return errors.New("Plain text file is too big, max 256 bytes")
	}

	pcrInfo, err := tpmtool.PreCalculate(TPMInterface, *cryptoCommandSealConfig)
	if err != nil {
		return err
	}

	sealed, err := TPMInterface.ResealData(*cryptoCommandSealLocality, pcrInfo, plainText, *cryptoCommandSrkPassword)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(*cryptoCommandSealCipherFile, sealed, 660)
}

// CryptoUnseal unseals data by the TPM against PCR
func CryptoUnseal() error {
	// No shit sherlock, we are too fast Oo (go-tpm bug ?)
	time.Sleep(Delay * time.Millisecond)

	cipherText, err := ioutil.ReadFile(*cryptoCommandUnsealCipherFile)
	if err != nil {
		return err
	}

	unsealed, err := TPMInterface.UnsealData(cipherText, *cryptoCommandSrkPassword)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(*cryptoCommandUnsealPlainFile, unsealed, 660)
}

// CryptoReseal reseals a data by given sealing configuration
func CryptoReseal() error {
	// No shit sherlock, we are too fast Oo (go-tpm bug ?)
	time.Sleep(Delay * time.Millisecond)

	sealedFile, err := ioutil.ReadFile(*cryptoCommandResealKeyfile)
	if err != nil {
		return err
	}

	pcrInfo, err := tpmtool.PreCalculate(TPMInterface, *cryptoCommandResealConfig)
	if err != nil {
		return err
	}

	unsealed, err := TPMInterface.UnsealData(sealedFile, *diskCommandSrkPassword)
	if err != nil {
		return err
	}

	sealed, err := TPMInterface.ResealData(tpm.DefaultLocality, pcrInfo, unsealed, *cryptoCommandSrkPassword)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(*cryptoCommandResealKeyfile, sealed, 660)
}

// PcrList dumps all PCRs
func PcrList() error {
	var pcrs string
	for i := uint32(0); i < MaxPlatformConfigurationRegister; i++ {
		hash, err := TPMInterface.ReadPCR(i)
		if err != nil {
			return err
		}
		pcrs += fmt.Sprintf("PCR-%02d: %x\n", i, hash)
	}
	fmt.Print(pcrs)

	return nil
}

// PcrRead reads the value of a PCR
func PcrRead() error {
	if *pcrCommandReadIndex >= MaxPlatformConfigurationRegister || *pcrCommandReadIndex < 0 {
		return errors.New("PCR index is incorrect")
	}

	pcr, err := TPMInterface.ReadPCR(*pcrCommandReadIndex)
	if err != nil {
		return err
	}

	fmt.Printf("PCR-%02d: %x\n", *pcrCommandReadIndex, pcr)

	return nil
}

// PcrMeasure measures a file into a defined PCR
func PcrMeasure() error {
	if *pcrCommandMeasureIndex >= MaxPlatformConfigurationRegister || *pcrCommandMeasureIndex < 0 {
		return errors.New("PCR index is incorrect")
	}

	fileToMeasure, err := ioutil.ReadFile(*pcrCommandMeasureFile)
	if err != nil {
		return err
	}

	err = TPMInterface.Measure(*pcrCommandMeasureIndex, fileToMeasure)
	if err != nil {
		return err
	}

	return nil
}

// DiskFormat formats a device for luks setup.
func DiskFormat() error {
	// No shit sherlock, we are too fast Oo (go-tpm bug ?)
	time.Sleep(Delay * time.Millisecond)

	keystorePath, err := tpmtool.MountKeystore()
	if err != nil {
		return err
	}
	defer tpmtool.UnmountKeystore(keystorePath)

	randBytes := make([]byte, 64)
	if _, err = rand.Read(randBytes); err != nil {
		return err
	}

	if err = ioutil.WriteFile(keystorePath+"/plain", randBytes, 660); err != nil {
		return err
	}

	if err = tpmtool.CryptsetupFormat(keystorePath+"/plain", *diskCommandFormatDevice); err != nil {
		return err
	}

	pcrInfo, err := tpmtool.PreCalculate(TPMInterface, *diskCommandFormatConfig)
	if err != nil {
		return err
	}

	sealed, err := TPMInterface.ResealData(*diskCommandFormatLocality, pcrInfo, randBytes, *diskCommandSrkPassword)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(*diskCommandFormatFile, sealed, 660)
}

// DiskOpen opens a LUKS device
func DiskOpen() error {
	// No shit sherlock, we are too fast Oo (go-tpm bug ?)
	time.Sleep(Delay * time.Millisecond)

	keystorePath, err := tpmtool.MountKeystore()
	if err != nil {
		return err
	}
	defer tpmtool.UnmountKeystore(keystorePath)

	if _, err = os.Stat(*diskCommandOpenMountPath); os.IsNotExist(err) {
		return err
	}

	sealedFile, err := ioutil.ReadFile(*diskCommandOpenSealFile)
	if err != nil {
		return err
	}

	sealed, err := TPMInterface.UnsealData(sealedFile, *diskCommandSrkPassword)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(keystorePath+"/plain", sealed, 660); err != nil {
		return err
	}

	deviceName, err := tpmtool.CryptsetupOpen(keystorePath+"/plain", *diskCommandOpenDevice)
	if err != nil {
		return err
	}

	fmt.Printf("Sealed encrypted device opened with name: %s\n", deviceName)

	return nil
}

// DiskClose closes a LUKS device
func DiskClose() error {
	deviceMapperPath := path.Join(tpmtool.DefaultDevMapperPath, *diskCommandCloseName)
	mountpoint, err := storage.GetMountpointByDevice(deviceMapperPath)
	if err == nil {
		syscall.Unmount(*mountpoint, syscall.MNT_DETACH|syscall.MNT_FORCE)
	}

	return tpmtool.CryptsetupClose(*diskCommandCloseName)
}

// DiskExtend hashes and extends a LUKS header into a PCR
func DiskExtend() error {
	deviceFD, err := os.Open(*diskCommandExtendDevice)
	if err != nil {
		return err
	}
	defer deviceFD.Close()

	luksHeader := make([]byte, tpmtool.Luks1HeaderLength)
	_, err = deviceFD.Read(luksHeader)
	if err != nil {
		return err
	}

	return TPMInterface.Measure(*diskCommandExtendPcr, luksHeader)
}

// EventlogDump dumps the eventlog
func EventlogDump() error {
	if *eventlogDumpFile != "" {
		tpm.DefaultTCPABinaryLog = *eventlogDumpFile
	}

	var firmware tpmtool.FirmwareType
	if *eventlogDumpFirmwareUefi {
		firmware = tpmtool.Uefi
	} else if *eventlogDumpFirmwareBios {
		firmware = tpmtool.Bios
	} else {
		if _, err := os.Stat(LinuxEFIFirmwareDir); os.IsNotExist(err) {
			firmware = tpmtool.Uefi
		} else {
			firmware = tpmtool.Bios
		}
	}

	if *eventlogDumpTPMSpec1 {
		TPMSpecVersion = tpm.TPM12
	} else if *eventlogDumpTPMSpec2 {
		TPMSpecVersion = tpm.TPM20
	}

	tcpaLog, err := tpm.ParseLog(string(firmware), TPMSpecVersion)
	if err != nil {
		return err
	}

	return tpm.DumpLog(tcpaLog)
}
