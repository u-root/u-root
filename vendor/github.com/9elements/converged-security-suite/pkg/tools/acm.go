package tools

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

const (
	//ACMChipsetTypeBios as defined in Document 315168-016 Chapter A.1 Table 8. Authenticated Code Module Format
	ACMChipsetTypeBios uint8 = 0x00
	//ACMChipsetTypeSinit as defined in Document 315168-016 Chapter A.1 Table 8. Authenticated Code Module Format
	ACMChipsetTypeSinit uint8 = 0x01
	//ACMChipsetTypeBiosRevoc as defined in Document 315168-016 Chapter A.1 Table 10. Chipset AC Module Information Table
	ACMChipsetTypeBiosRevoc uint8 = 0x08
	//ACMChipsetTypeSinitRevoc as defined in Document 315168-016 Chapter A.1 Table 10. Chipset AC Module Information Table
	ACMChipsetTypeSinitRevoc uint8 = 0x09
	//ACMTypeChipset as defined in Document 315168-016 Chapter A.1 Table 8. Authenticated Code Module Format
	ACMTypeChipset uint16 = 0x02
	//ACMSubTypeReset FIXME
	ACMSubTypeReset uint16 = 0x01
	//ACMVendorIntel as defined in Document 315168-016 Chapter A.1 Table 8. Authenticated Code Module Format
	ACMVendorIntel uint32 = 0x8086

	//TPMExtPolicyIllegal as defined in Document 315168-016 Chapter A.1 Table 16. TPM Capabilities Field
	TPMExtPolicyIllegal uint8 = 0x00
	//TPMExtPolicyAlgAgile as defined in Document 315168-016 Chapter A.1 Table 16. TPM Capabilities Field
	TPMExtPolicyAlgAgile uint8 = 0x01
	//TPMExtPolicyEmbeddedAlgs as defined in Document 315168-016 Chapter A.1 Table 16. TPM Capabilities Field
	TPMExtPolicyEmbeddedAlgs uint8 = 0x10
	//TPMExtPolicyBoth as defined in Document 315168-016 Chapter A.1 Table 16. TPM Capabilities Field
	TPMExtPolicyBoth uint8 = 0x11

	//TPMFamilyIllegal as defined in Document 315168-016 Chapter A.1 Table 16. TPM Capabilities Field
	TPMFamilyIllegal uint16 = 0x0000
	//TPMFamilyDTPM12 as defined in Document 315168-016 Chapter A.1 Table 16. TPM Capabilities Field
	TPMFamilyDTPM12 uint16 = 0x0001
	//TPMFamilyDTPM20 as defined in Document 315168-016 Chapter A.1 Table 16. TPM Capabilities Field
	TPMFamilyDTPM20 uint16 = 0x0010
	//TPMFamilyDTPMBoth combination out of TPMFamilyDTPM12 and TPMFamilyDTPM20
	TPMFamilyDTPMBoth uint16 = 0x0011
	//TPMFamilyPTT20 as defined in Document 315168-016 Chapter A.1 Table 16. TPM Capabilities Field
	TPMFamilyPTT20 uint16 = 0x1000

	//ACMUUIDV3 as defined in Document 315168-016 Chapter A.1 Table 10. Chipset AC Module Information Table
	ACMUUIDV3 string = "7fc03aaa-46a7-18db-ac2e-698f8d417f5a"
	//ACMSizeOffset as defined in Document 315168-016 Chapter A.1 Table 8. Authenticated Code Module Format
	ACMSizeOffset int64 = 24

	//TPMAlgoSHA1 as defined in Document 315168-016 Chapter D.1.3 LCP_POLICY2
	TPMAlgoSHA1 uint16 = 0x0004
	//TPMAlgoSHA256 as defined in Document 315168-016 Chapter D.1.3 LCP_POLICY2
	TPMAlgoSHA256 uint16 = 0x000b
	//TPMAlgoSHA384 FIXME
	TPMAlgoSHA384 uint16 = 0x000c
	//TPMAlgoSHA512 FIXME
	TPMAlgoSHA512 uint16 = 0x000d
	//TPMAlgoNULL as defined in Document 315168-016 Chapter D.1.3 LCP_POLICY2
	TPMAlgoNULL uint16 = 0x0010
	//TPMAlgoSM3_256 as defined in Document 315168-016 Chapter D.1.3 LCP_POLICY2
	TPMAlgoSM3_256 uint16 = 0x0012
	//TPMAlgoRSASSA as defined in Document 315168-016 Chapter D.1.3 LCP_POLICY2
	TPMAlgoRSASSA uint16 = 0x0014
	//TPMAlgoECDSA as defined in Document 315168-016 Chapter D.1.3 LCP_POLICY2
	TPMAlgoECDSA uint16 = 0x0018
	//TPMAlgoSM2 as defined in Document 315168-016 Chapter D.1.3 LCP_POLICY2
	TPMAlgoSM2 uint16 = 0x001B

	//ACMheaderLen as defined in Document 315168-016 Chapter A.1 Table 8. Authenticated Code Module Format (Version 0.0)
	ACMheaderLen uint32 = 161

	//ACMModuleSubtypeSinitACM is an enum
	ACMModuleSubtypeSinitACM uint16 = 0
	//ACMModuleSubtypeCapableOfExecuteAtReset is a flag and enum Based on EDK2 Silicon/Intel/Tools/FitGen/FitGen.c
	ACMModuleSubtypeCapableOfExecuteAtReset uint16 = 1
	//ACMModuleSubtypeAncModule is a flag Based on EDK2 Silicon/Intel/Tools/FitGen/FitGen.c
	ACMModuleSubtypeAncModule uint16 = 2
)

//UUID represents an UUID
type UUID struct {
	Field1 uint32
	Field2 uint16
	Field3 uint16
	Field4 uint16
	Field5 [6]uint8
}

// ACMInfo holds the metadata extracted from the ACM header
type ACMInfo struct {
	UUID                UUID
	ChipsetACMType      uint8
	Version             uint8
	Length              uint16
	ChipsetIDList       uint32
	OSSinitDataVersion  uint32
	MinMleHeaderVersion uint32
	TxtCaps             uint32
	ACMVersion          uint8
	Reserved            [3]uint8
	ProcessorIDList     uint32
	TPMInfoList         uint32
}

//ChipsetID describes the chipset ID found in the ACM header
type ChipsetID struct {
	Flags      uint32
	VendorID   uint16
	DeviceID   uint16
	RevisionID uint16
	Reserved   [3]uint16
}

//Chipsets hold a list of supported chipset IDs as found in the ACM header
type Chipsets struct {
	Count  uint32
	IDList []ChipsetID
}

//ProcessorID describes the processor ID found in the ACM header
type ProcessorID struct {
	FMS          uint32
	FMSMask      uint32
	PlatformID   uint64
	PlatformMask uint64
}

//Processors hold a list of supported processor IDs as found in the ACM header
type Processors struct {
	Count  uint32
	IDList []ProcessorID
}

//TPMs describes the required TPM capabilties and algorithm as found in the ACM header
type TPMs struct {
	Capabilities uint32
	Count        uint16
	AlgID        []uint16
}

// ACMHeader exports the structure of ACM Header found in the firemware interface table
type ACMHeader struct {
	ModuleType      uint16
	ModuleSubType   uint16
	HeaderLen       uint32
	HeaderVersion   uint32
	ChipsetID       uint16
	Flags           uint16
	ModuleVendor    uint32
	Date            uint32
	Size            uint32
	TxtSVN          uint16
	SeSVN           uint16
	CodeControl     uint32
	ErrorEntryPoint uint32
	GDTLimit        uint32
	GDTBase         uint32
	SegSel          uint32
	EntryPoint      uint32
	Reserved2       [64]uint8
	KeySize         uint32
	ScratchSize     uint32
	PubKey          [256]uint8
	PubExp          uint32
	Signatur        [256]uint8
}

// ACM exports the structure of Authenticated Code Modules found in the Firmware Interface Table(FIT)
type ACM struct {
	Header  ACMHeader
	Scratch []byte
	Info    ACMInfo
}

// ACMFlags exports the ACM header flags
type ACMFlags struct {
	Production    bool
	PreProduction bool
	DebugSigned   bool
}

// ParseACMHeader exports the functionality of parsing an ACM Header
func ParseACMHeader(data []byte) (*ACMHeader, error) {
	var acm ACMHeader
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &acm)

	if err != nil {
		return nil, fmt.Errorf("Can't read ACM Header")
	}

	return &acm, nil
}

// ValidateACMHeader validates an ACM Header found in the Firmware Interface Table (FIT)
func ValidateACMHeader(acmheader *ACMHeader) (bool, error) {
	if acmheader.ModuleType != uint16(2) {
		return false, fmt.Errorf("BIOS ACM ModuleType is not 2, this is not specified")
	}
	// Early version of TXT used an enum in ModuleSubType
	// That was changed to flags. Check if unsupported flags are present
	if acmheader.ModuleSubType > (ACMModuleSubtypeAncModule | ACMModuleSubtypeCapableOfExecuteAtReset) {
		return false, fmt.Errorf("BIOS ACM ModuleSubType contains unknown flags")
	}
	if acmheader.HeaderLen < uint32(ACMheaderLen) {
		return false, fmt.Errorf("BIOS ACM HeaderLength is smaller than 4*161 Byte")
	}
	if acmheader.Size == 0 {
		return false, fmt.Errorf("BIOS ACM Size can't be zero")
	}
	if acmheader.ModuleVendor != ACMVendorIntel {
		return false, fmt.Errorf("AC Module Vendor is not Intel. Only Intel as Vendor is allowed")
	}
	if acmheader.KeySize*4 != uint32(len(acmheader.PubKey)) {
		return false, fmt.Errorf("ACM keysize of 0x%x not supported yet", acmheader.KeySize*4)
	}
	if acmheader.ScratchSize > acmheader.Size {
		return false, fmt.Errorf("ACM ScratchSize is bigger than ACM module size")
	}
	return true, nil
}

//ParseACM deconstructs a byte array containing the raw ACM into it's components
func ParseACM(data []byte) (*ACM, *Chipsets, *Processors, *TPMs, error, error) {
	var acmheader ACMHeader
	var acminfo ACMInfo
	var processors Processors
	var chipsets Chipsets
	var tpms TPMs

	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &acmheader)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	scratch := make([]byte, acmheader.ScratchSize*4)

	err = binary.Read(buf, binary.LittleEndian, &scratch)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if (acmheader.ModuleSubType & ACMModuleSubtypeAncModule) > 0 {
		// ANC modules do not have an ACMINFO header
		acm := ACM{acmheader, scratch, acminfo}
		return &acm, &chipsets, &processors, &tpms, nil, nil
	}

	err = binary.Read(buf, binary.LittleEndian, &acminfo)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	acm := ACM{acmheader, scratch, acminfo}

	buf.Seek(int64(acm.Info.ChipsetIDList), io.SeekStart)
	err = binary.Read(buf, binary.LittleEndian, &chipsets.Count)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	chipsets.IDList = make([]ChipsetID, chipsets.Count)
	err = binary.Read(buf, binary.LittleEndian, &chipsets.IDList)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	buf.Seek(int64(acm.Info.ProcessorIDList), io.SeekStart)
	err = binary.Read(buf, binary.LittleEndian, &processors.Count)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	processors.IDList = make([]ProcessorID, processors.Count)
	err = binary.Read(buf, binary.LittleEndian, &processors.IDList)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if acm.Info.ACMVersion >= 5 {
		buf.Seek(int64(acm.Info.TPMInfoList), io.SeekStart)
		err = binary.Read(buf, binary.LittleEndian, &tpms.Capabilities)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}

		err = binary.Read(buf, binary.LittleEndian, &tpms.Count)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}

		tpms.AlgID = make([]uint16, tpms.Count)
		for i := 0; i < int(tpms.Count); i++ {
			err = binary.Read(buf, binary.LittleEndian, &tpms.AlgID[i])
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
		}
	}

	return &acm, &chipsets, &processors, &tpms, nil, nil
}

//LookupSize returns the ACM size
func LookupSize(header []byte) (int64, error) {
	var acmSize uint32

	buf := bytes.NewReader(header[:32])
	buf.Seek(ACMSizeOffset, io.SeekStart)
	err := binary.Read(buf, binary.LittleEndian, &acmSize)
	if err != nil {
		return 0, err
	}

	return int64(acmSize * 4), nil
}

// ParseACMFlags parses the ACM Header flags
func (a *ACMHeader) ParseACMFlags() *ACMFlags {
	var flags ACMFlags
	flags.Production = (a.Flags>>15)&1 == 0 && (a.Flags>>14)&1 == 0
	flags.PreProduction = (a.Flags>>14)&1 != 0
	flags.DebugSigned = (a.Flags>>15)&1 != 0
	return &flags
}

//PrettyPrint prints a human readable representation of the ACMHeader
func (a *ACMHeader) PrettyPrint() {
	log.Println("Authenticated Code Module")

	if a.ModuleVendor == ACMVendorIntel {
		log.Println("Module Vendor: Intel")
	} else {
		log.Println("Module Vendor: Unknown")
	}

	if a.ModuleType == ACMTypeChipset {
		log.Println("Module Type: ACM_TYPE_CHIPSET")
	} else {
		log.Println("Module Type: UNKNOWN")
	}

	if a.ModuleSubType == ACMSubTypeReset {
		log.Println("Module Subtype: Execute at Reset")
	} else if a.ModuleSubType == 0 {
		log.Println("Module Subtype: 0x0")
	} else {
		log.Println("Module Subtype: Unknown")
	}
	log.Printf("Module Date: 0x%02x\n", a.Date)
	log.Printf("Module Size: 0x%x (%d)\n", a.Size*4, a.Size*4)

	log.Printf("Header Length: 0x%x (%d)\n", a.HeaderLen, a.HeaderLen)
	log.Printf("Header Version: %d\n", a.HeaderVersion)
	log.Printf("Chipset ID: 0x%02x\n", a.ChipsetID)
	log.Printf("Flags: 0x%02x\n", a.Flags)
	log.Printf("TXT SVN: 0x%08x\n", a.TxtSVN)
	log.Printf("SE SVN: 0x%08x\n", a.SeSVN)
	log.Printf("Code Control: 0x%02x\n", a.CodeControl)
	log.Printf("Entry Point: 0x%08x:%08x\n", a.SegSel, a.EntryPoint)
	log.Printf("Scratch Size: 0x%x (%d)\n", a.ScratchSize, a.ScratchSize)
}

//PrettyPrint prints a human readable representation of the ACM
func (a *ACM) PrettyPrint() {
	a.Header.PrettyPrint()
	log.Println("Info Table:")

	uuidStr := fmt.Sprintf("%08x-%04x-%04x-%04x-%02x%02x%02x%02x%02x%02x",
		a.Info.UUID.Field1,
		a.Info.UUID.Field2,
		a.Info.UUID.Field3,
		a.Info.UUID.Field4,
		a.Info.UUID.Field5[0],
		a.Info.UUID.Field5[1],
		a.Info.UUID.Field5[2],
		a.Info.UUID.Field5[3],
		a.Info.UUID.Field5[4],
		a.Info.UUID.Field5[5])

	if uuidStr == ACMUUIDV3 {
		log.Println("\tUUID: ACM_UUID_V3")
	}

	switch a.Info.ChipsetACMType {
	case ACMChipsetTypeBios:
		log.Println("\tChipset ACM: BIOS")
		break
	case ACMChipsetTypeBiosRevoc:
		log.Println("\tChipset ACM: BIOS Revocation")
		break
	case ACMChipsetTypeSinit:
		log.Println("\tChipset ACM: SINIT")
		break
	case ACMChipsetTypeSinitRevoc:
		log.Println("\tChipset ACM: SINIT Revocation")
		break
	default:
		log.Println("\tChipset ACM: Unknown")
	}

	log.Printf("\tVersion: %d\n", a.Info.Version)
	log.Printf("\tLength: 0x%x (%d)\n", a.Info.Length, a.Info.Length)
	log.Printf("\tChipset ID List: 0x%02x\n", a.Info.ChipsetIDList)
	log.Printf("\tOS SINIT Data Version: 0x%02x\n", a.Info.OSSinitDataVersion)
	log.Printf("\tMin. MLE Header Version: 0x%08x\n", a.Info.MinMleHeaderVersion)
	log.Printf("\tCapabilities: 0x%08x\n", a.Info.TxtCaps)
	log.Printf("\tACM Version: %d\n", a.Info.ACMVersion)
}

//PrettyPrint prints a human readable representation of the Chipsets
func (c *Chipsets) PrettyPrint() {
	log.Println("Chipset List:")
	log.Printf("\tEntries: %d\n", c.Count)
	for idx, chipset := range c.IDList {
		log.Printf("\tEntry %d:\n", idx)
		log.Printf("\t\tFlags: 0x%02x\n", chipset.Flags)
		log.Printf("\t\tVendor: 0x%02x\n", chipset.VendorID)
		log.Printf("\t\tDevice: 0x%02x\n", chipset.DeviceID)
		log.Printf("\t\tRevision: 0x%02x\n", chipset.RevisionID)
	}
}

//PrettyPrint prints a human readable representation of the Processors
func (p *Processors) PrettyPrint() {
	log.Println("Processor List:")
	log.Printf("\tEntries: %d\n", p.Count)
	for idx, processor := range p.IDList {
		log.Printf("\tEntry %d:\n", idx)
		log.Printf("\t\tFMS: 0x%02x\n", processor.FMS)
		log.Printf("\t\tFMS Maks: 0x%02x\n", processor.FMSMask)
		log.Printf("\t\tPlatform ID: 0x%02x\n", processor.PlatformID)
		log.Printf("\t\tPlatform Mask: 0x%02x\n", processor.PlatformMask)
	}
}

//PrettyPrint prints a human readable representation of the TPMs
func (t *TPMs) PrettyPrint() {
	log.Println("TPM Info List:")
	log.Println("\tCapabilities:")
	log.Printf("\t\tExternal Policy: %02x\n", t.Capabilities)
	log.Printf("\tAlgorithms: %d\n", t.Count)
	for _, algo := range t.AlgID {
		switch algo {
		case TPMAlgoNULL:
			log.Println("\t\tNULL")
			break
		case TPMAlgoSHA1:
			log.Println("\t\tSHA-1")
			break
		case TPMAlgoSHA256:
			log.Println("\t\tSHA-256")
			break
		case TPMAlgoSHA384:
			log.Println("\t\tSHA-384")
			break
		case TPMAlgoSHA512:
			log.Println("\t\tSHA-512")
			break
		case TPMAlgoSM3_256:
			log.Println("\t\tSM3-256")
			break
		case TPMAlgoRSASSA:
			log.Println("\t\tRSA-SSA")
			break
		case TPMAlgoECDSA:
			log.Println("\t\tEC-DSA")
			break
		case TPMAlgoSM2:
			log.Println("\t\tSM2")
			break
		default:
			log.Println("\t\tUnknown")
		}
	}
}
