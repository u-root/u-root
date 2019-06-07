package main

import (
	"log"
	"os"

	"github.com/systemboot/tpmtool/pkg/tpm"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	// Author is the author
	Author = "Philipp Deppenwiese"
	// HelpText is the command line help
	HelpText = "A Linux only tool for TPM interaction"
)

var goversion string

var (
	// TPMDevice which should be used
	tpmDevice = kingpin.Flag("device", "TPM device path").Default("/dev/tpm0").String()

	// TPMInterface is a global TPM interface
	TPMInterface tpm.TPM

	// TPMSpecVersion is the version of the TPM
	TPMSpecVersion string

	// CommandLine Arguments
	status = kingpin.Command("status", "Show the TPM status information")

	ekCommand         = kingpin.Command("ek", "Dump the TPM EK")
	ekCommandPassword = ekCommand.Flag("owner-pass", "TPM owner password").String()
	ekCommandOutfile  = ekCommand.Arg("outfile", "File path to write EK").String()

	ownerCommand         = kingpin.Command("owner", "Management of the TPM")
	ownerCommandPassword = ownerCommand.Flag("owner-pass", "TPM owner password").String()

	ownerCommandTake            = ownerCommand.Command("take", "Take ownership of the TPM")
	ownerCommandTakeSrkPassword = ownerCommandTake.Flag("srk-pass", "TPM SRK password").String()

	ownerCommandClear = ownerCommand.Command("clear", "Clear ownership of the TPM")

	ownerCommandResetLock = ownerCommand.Command("reset-lock", "Reset the TPM lock")

	cryptoCommand               = kingpin.Command("crypto", "Manage TPM data encryption")
	cryptoCommandSrkPassword    = cryptoCommand.Flag("srk-pass", "TPM SRK password").String()
	cryptoCommandSealPlainFile  = cryptoCommandSeal.Arg("plain-file", "Plain text data file path").Required().String()
	cryptoCommandSealCipherFile = cryptoCommandSeal.Arg("sealed-file", "Encrypted data file path").Required().String()

	cryptoCommandSeal         = cryptoCommand.Command("seal", "Seal data against the TPM")
	cryptoCommandSealConfig   = cryptoCommandSeal.Arg("config", "Sealing configuration for PCR pre-calculation").Required().String()
	cryptoCommandSealLocality = cryptoCommandSeal.Flag("locality", "Sets the locality for the sealing operation").Uint8()

	cryptoCommandUnseal           = cryptoCommand.Command("unseal", "Unseal data against the TPM")
	cryptoCommandUnsealCipherFile = cryptoCommandUnseal.Arg("sealed-file", "Encrypted data file path").Required().String()
	cryptoCommandUnsealPlainFile  = cryptoCommandUnseal.Arg("plain-file", "Plain text data file path").Required().String()

	cryptoCommandReseal         = cryptoCommand.Command("reseal", "Reseal already sealed credentials based on a sealing configuration")
	cryptoCommandResealLocality = cryptoCommandReseal.Flag("locality", "Sets the locality for the sealing operation").Uint8()
	cryptoCommandResealConfig   = cryptoCommandReseal.Arg("config", "Sealing configuration for PCR pre-calculation").Required().String()
	cryptoCommandResealKeyfile  = cryptoCommandReseal.Arg("sealed-key-file", "Sealed encrypted key").Required().String()

	pcrCommand = kingpin.Command("pcr", "Manage TPM PCR operations")

	pcrCommandPrint = pcrCommand.Command("list", "Print all PCRs")

	pcrCommandRead      = pcrCommand.Command("read", "Read a specific PCR")
	pcrCommandReadIndex = pcrCommandRead.Flag("pcr", "Set the PCR for the read operation").Required().Uint32()

	pcrCommandMeasure      = pcrCommand.Command("measure", "Measure data into a given PCR")
	pcrCommandMeasureIndex = pcrCommandMeasure.Flag("pcr", "Set the PCR for the measurement operation").Required().Uint32()
	pcrCommandMeasureFile  = pcrCommandMeasure.Arg("measure-file", "File which should be measured into a PCR").Required().String()

	diskCommand            = kingpin.Command("disk", "Manage cryptsetup sealed devices")
	diskCommandSrkPassword = diskCommand.Flag("srk-pass", "TPM SRK password").String()

	diskCommandFormat         = diskCommand.Command("format", "Format cryptsetup partition with sealing")
	diskCommandFormatFile     = diskCommandFormat.Arg("sealed-key-file", "Sealed encrypted key").Required().String()
	diskCommandFormatDevice   = diskCommandFormat.Arg("device", "A device which should be encrypted").Required().String()
	diskCommandFormatConfig   = diskCommandFormat.Arg("config", "Sealing configuration for PCR pre-calculation").Required().String()
	diskCommandFormatLocality = diskCommandFormat.Flag("locality", "Sets the locality for the sealing operation").Uint8()

	diskCommandOpen          = diskCommand.Command("open", "Open cryptsetup partition with sealed key")
	diskCommandOpenSealFile  = diskCommandOpen.Arg("sealed-key-file", "Sealed encrypted key").Required().String()
	diskCommandOpenDevice    = diskCommandOpen.Arg("device", "Device which should be encrypted").Required().String()
	diskCommandOpenMountPath = diskCommandOpen.Arg("mnt-path", "Mount path for mounting unsealed encrypted device").Required().String()

	diskCommandClose     = diskCommand.Command("close", "Close cryptsetup partition")
	diskCommandCloseName = diskCommandClose.Arg("device-name", "cryptsetup device name").Required().String()

	diskCommandExtend       = diskCommand.Command("extend", "Extend luks header into a PCR")
	diskCommandExtendDevice = diskCommandExtend.Arg("device", "Device which should be encrypted").Required().String()
	diskCommandExtendPcr    = diskCommandExtend.Flag("pcr", "Set the PCR for the measurement operation").Required().Uint32()

	eventlog                 = kingpin.Command("eventlog", "TPM eventlog operation")
	eventlogDump             = eventlog.Command("dump", "Dump the eventlog")
	eventlogDumpFirmwareUefi = eventlogDump.Flag("uefi", "Set UEFI firmware").Bool()
	eventlogDumpFirmwareBios = eventlogDump.Flag("bios", "Set BIOS firmware").Bool()
	eventlogDumpTPMSpec1     = eventlogDump.Flag("tpm12", "Set tpm12 specification").Bool()
	eventlogDumpTPMSpec2     = eventlogDump.Flag("tpm20", "Set tpm20 specification").Bool()
	eventlogDumpFile         = eventlogDump.Arg("log", "Custom eventlog file path").String()
)

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(goversion).Author(Author)
	kingpin.CommandLine.Help = HelpText

	if *tpmDevice != "" {
		tpm.TPMDevice = *tpmDevice
	}

	// Check for root user
	file, err := os.Open(tpm.TPMDevice)
	if err != nil {
		log.Fatal("Please run this tool as root user")
	}
	file.Close()

	TPMInterface, err = tpm.NewTPM()
	if err != nil {
		log.Fatal("Can't open TPM interface: " + err.Error())
	}
	defer TPMInterface.Close()

	TPMSpecVersion = TPMInterface.Info().Specification

	switch kingpin.Parse() {
	case "status":
		if err := Status(); err != nil {
			log.Fatalln(err.Error())
		}
	case "ek":
		if err := Ek(); err != nil {
			log.Fatalln(err.Error())
		}
	case "owner take":
		if err := OwnerTake(); err != nil {
			log.Fatalln(err.Error())
		}
	case "owner clear":
		if err := OwnerClear(); err != nil {
			log.Fatalln(err.Error())
		}
	case "owner reset-lock":
		if err := OwnerResetLock(); err != nil {
			log.Fatalln(err.Error())
		}
	case "crypto seal":
		if err := CryptoSeal(); err != nil {
			log.Fatalln(err.Error())
		}
	case "crypto unseal":
		if err := CryptoUnseal(); err != nil {
			log.Fatalln(err.Error())
		}
	case "crypto reseal":
		if err := CryptoReseal(); err != nil {
			log.Fatalln(err.Error())
		}
	case "pcr list":
		if err := PcrList(); err != nil {
			log.Fatalln(err.Error())
		}
	case "pcr read":
		if err := PcrRead(); err != nil {
			log.Fatalln(err.Error())
		}
	case "pcr measure":
		if err := PcrMeasure(); err != nil {
			log.Fatalln(err.Error())
		}
	case "disk format":
		if err := DiskFormat(); err != nil {
			log.Fatalln(err.Error())
		}
	case "disk open":
		if err := DiskOpen(); err != nil {
			log.Fatalln(err.Error())
		}
	case "disk close":
		if err := DiskClose(); err != nil {
			log.Fatalln(err.Error())
		}
	case "disk extend":
		if err := DiskExtend(); err != nil {
			log.Fatalln(err.Error())
		}
	case "eventlog dump":
		if err := EventlogDump(); err != nil {
			log.Fatalln(err.Error())
		}
	default:
		log.Fatal("Command not found")
	}
}
