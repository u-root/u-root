package tools

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/google/go-tpm/tpm2"
)

//LCPPol2Hash stores the hashing algorithm used in the LCP policy version 2
type LCPPol2Hash uint16

const (
	//LCPPol2HAlgSHA1 as defined in Document 315168-016 Chapter D.3.2.1 LCP_POLICY_LIST2 Structure. Same as TPMAlgoSHA1
	LCPPol2HAlgSHA1 LCPPol2Hash = 0x04
	//LCPPol2HAlgSHA256 as defined in Document 315168-016 Chapter D.3.2.1 LCP_POLICY_LIST2 Structure. Same as TPMAlgoSHA256
	LCPPol2HAlgSHA256 LCPPol2Hash = 0x0B
	//LCPPol2HAlgSHA384 as defined in Document 315168-016 Chapter D.3.2.1 LCP_POLICY_LIST2 Structure. Same as TPMAlgoSHA384
	LCPPol2HAlgSHA384 LCPPol2Hash = 0x0C
	//LCPPol2HAlgNULL as defined in Document 315168-016 Chapter D.3.2.1 LCP_POLICY_LIST2 Structure. Same as TPMAlgoNULL
	LCPPol2HAlgNULL LCPPol2Hash = 0x10
	//LCPPol2HAlgSM3 as defined in Document 315168-016 Chapter D.3.2.1 LCP_POLICY_LIST2 Structure. Same as TPMAlgoSM3_256
	LCPPol2HAlgSM3 LCPPol2Hash = 0x12
)

const (
	//LCPPolicyVersion2 as defined in Document 315168-016 Chapter 3.2.1 LCP Policy
	LCPPolicyVersion2 uint16 = 0x0204
	//LCPPolicyVersion3 as defined in Document 315168-016 Chapter 3.2.1 LCP Policy
	LCPPolicyVersion3 uint16 = 0x0300
	//LCPPolicyTypeAny as defined in Document 315168-016 Chapter D LCP Data Structures
	LCPPolicyTypeAny uint8 = 1
	//LCPPolicyTypeList as defined in Document 315168-016 Chapter D LCP Data Structures
	LCPPolicyTypeList uint8 = 0
	//LCPMaxLists as defined in Document 315168-016 Chapter D LCP Data Structures
	LCPMaxLists uint = 8
	//SHA1DigestSize as defined in Document 315168-016 Chapter D.1.3 LCP_POLICY2
	SHA1DigestSize uint = 20
	//SHA256DigestSize as defined in Document 315168-016 Chapter D.1.3 LCP_POLICY2
	SHA256DigestSize uint = 32
	//SHA384DigestSize as defined in Document 315168-016 Chapter D.1.3 LCP_POLICY2
	SHA384DigestSize uint = 48
	//SHA512DigestSize FIXME
	SHA512DigestSize uint = 64
	//SM3DigestSize as defined in Document 315168-016 Chapter D.1.3 LCP_POLICY2
	SM3DigestSize uint = 32
	//LCPDataFileSignature as defined in Document 315168-016 Chapter D.2 LCP_POLICY_DATA
	LCPDataFileSignature string = "Intel(R) TXT LCP_POLICY_DATA\x00\x00\x00\x00"

	//LCPSignatureAlgNone FIXME
	LCPSignatureAlgNone uint8 = 0
	//LCPSignatureAlgRSAPKCS15 FIXME
	LCPSignatureAlgRSAPKCS15 uint8 = 1

	//LCPPolicyElementMLE as defined in Document 315168-016 Chapter D.4.4 LCP_MLE_ELEMENT
	LCPPolicyElementMLE uint32 = 0
	//LCPPolicyElementPCONF as defined in Document 315168-016 Chapter D.4.5 LCP_PCONF_ELEMENT
	LCPPolicyElementPCONF uint32 = 1
	//LCPPolicyElementSBIOS FIXME
	LCPPolicyElementSBIOS uint32 = 2
	//LCPPolicyElementCustom as defined in Document 315168-016 Chapter D.4.6 LCP_CUSTOM_ELEMENT
	LCPPolicyElementCustom uint32 = 3
	//LCPPolicyElementMLE2 as defined in Document 315168-016 Chapter D.4.7 LCP_MLE_ELEMENT2
	LCPPolicyElementMLE2 uint32 = 0x10
	//LCPPolicyElementPCONF2 as defined in Document 315168-016 Chapter D.4.8 LCP_PCONF_ELEMENT2
	LCPPolicyElementPCONF2 uint32 = 0x11
	//LCPPolicyElementSBIOS2 FIXME
	LCPPolicyElementSBIOS2 uint32 = 0x12
	//LCPPolicyElementSTM2 as defined in Document 315168-016 Chapter D.4.9 LCP_STM_ELEMENT2
	LCPPolicyElementSTM2 uint32 = 0x14

	// LCPPolHAlgSHA1 Document 315168-016 Chapter D.1 LCP_POLICY
	LCPPolHAlgSHA1 uint8 = 0
)

//LCPHash holds one of the supported hashes
type LCPHash struct {
	sha1   *[SHA1DigestSize]uint8
	sha256 *[SHA256DigestSize]uint8
	sha384 *[SHA384DigestSize]uint8
	sha512 *[SHA512DigestSize]uint8
	sm3    *[SM3DigestSize]uint8
}

//LCPPolicyElement represents a policy element as defined in Document 315168-016 Chapter D.4 LCP_POLICY_ELEMENT
type LCPPolicyElement struct {
	Size             uint32
	Type             uint32
	PolicyEltControl uint32
	MLE              *LCPPolicyMLE
	SBIOS            *LCPPolicySBIOS
	PCONF            *LCPPolicyPCONF
	Custom           *LCPPolicyCustom
}

//LCPPolicyMLE represents a MLE policy element as defined in Document 315168-016 Chapter D.4.4 LCP_MLE_ELEMENT
type LCPPolicyMLE struct {
	SINITMinVersion uint8
	HashAlg         uint8
	NumHashes       uint16
	Hashes          [][20]byte
}

//LCPPolicySBIOS represents a SBIOS policy element
type LCPPolicySBIOS struct {
	HashAlg      uint8
	Reserved1    [3]uint8
	FallbackHash LCPHash
	Reserved2    uint16
	NumHashes    uint16
	Hashes       []LCPHash
}

//LCPPolicyPCONF represents a PCONF policy element
type LCPPolicyPCONF struct {
	NumPCRInfos uint16
	PCRInfos    []TPMPCRInfoShort
}

//TPMPCRInfoShort rFIXME
type TPMPCRInfoShort struct {
	// TPM_PCR_SELECTION
	PCRSelect []int
	// TPM_LOCALITY_SELECTION
	LocalityAtRelease uint8
	// TPM_COMPOSITE_HASH
	DigestAtRelease [20]byte
}

//LCPPolicyCustom represents a custom policy element
type LCPPolicyCustom struct {
	UUID LCPUUID
	Data []byte
}

//LCPUUID represents an UUID
type LCPUUID struct {
	data1 uint32
	data2 uint16
	data3 uint16
	data4 uint16
	data5 [6]uint8
}

//LCPPolicyList2 as defined in Document 315168-016 Chapter D.3.2.1 LCP_POLICY_LIST2 Structure
type LCPPolicyList2 struct {
	Version           uint16
	SignaturAlg       uint16
	PolicyElementSize uint32
	PolicyElements    []LCPPolicyElement
}

//LCPSignature as defined in Document 315168-016 Chapter D.3.2.1 LCP_POLICY_LIST2 Structure
type LCPSignature struct {
	RevocationCounter uint16
	PubkeySize        uint16
	PubkeyValue       []byte
	SigBlock          []byte
}

//LCPPolicyList FIXME not in Document 315168-016
type LCPPolicyList struct {
	Version           uint16
	Reserved          uint8
	SignaturAlg       uint8
	PolicyElementSize uint32
	PolicyElements    []LCPPolicyElement
	Signature         *LCPSignature
}

//LCPList as defined in Document 315168-016 Chapter D.3.2.3 LCP_LIST
type LCPList struct {
	TPM12PolicyList LCPPolicyList
	TPM20PolicyList LCPPolicyList2
}

//PolicyControl as defined in Document 315168-016 Chapter D.1.1 PolicyControl
type PolicyControl struct {
	NPW           bool
	OwnerEnforced bool
	AuxDelete     bool
	SinitCaps     bool
}

//ApprovedHashAlgorithm as defined in Document 315168-016 Chapter D.1.3 LCP_POLICY2
type ApprovedHashAlgorithm struct {
	SHA1   bool
	SHA256 bool
	SHA384 bool
	SM3    bool
}

//ApprovedSignatureAlogrithm as defined in Document 315168-016 Chapter D.1.3 LCP_POLICY2
type ApprovedSignatureAlogrithm struct {
	RSA2048SHA1     bool
	RSA2048SHA256   bool
	RSA3072SHA256   bool
	RSA3072SHA384   bool
	ECDSAP256SHA256 bool
	ECDSAP384SHA384 bool
	SM2SM2CurveSM3  bool
}

//LCPPolicy as defined in Document 315168-016 Chapter D.1.2 LCP_POLICY
type LCPPolicy struct {
	Version                uint16 // < 0x0204
	HashAlg                uint8
	PolicyType             uint8
	SINITMinVersion        uint8
	Reserved               uint8
	DataRevocationCounters [LCPMaxLists]uint16
	PolicyControl          uint32
	MaxSINITMinVersion     uint8
	Reserved1              uint8
	Reserved2              uint16
	Reserved3              uint32
	PolicyHash             [20]byte
}

//LCPPolicy2 as defined in Document 315168-016 Chapter D.1.3 LCP_POLICY2
type LCPPolicy2 struct {
	Version                uint16 // < 0x0302
	HashAlg                LCPPol2Hash
	PolicyType             uint8
	SINITMinVersion        uint8
	DataRevocationCounters [LCPMaxLists]uint16
	PolicyControl          uint32
	MaxSINITMinVersion     uint8 // v2.0 - Only PO index, reserved for PS
	Reserved               uint8 // v2.0 - Only PO index, reserved for PS
	LcpHashAlgMask         uint16
	LcpSignAlgMask         uint32
	Reserved2              uint32
	PolicyHash             LCPHash
}

//LCPPolicyData FIXME
type LCPPolicyData struct {
	FileSignature [32]uint8
	Reserved      [3]uint8
	NumLists      uint8
	PolicyLists   []LCPList
}

// ParsePolicyControl TODO needs to be reverse engineered
func (p *LCPPolicy) ParsePolicyControl() PolicyControl {
	var polCtrl PolicyControl
	polCtrl.NPW = (p.PolicyControl>>1)&1 != 0
	polCtrl.SinitCaps = (p.PolicyControl>>2)&1 != 0
	polCtrl.AuxDelete = (p.PolicyControl>>15)&1 != 0
	polCtrl.OwnerEnforced = (p.PolicyControl>>3)&1 != 0
	return polCtrl
}

// ParsePolicyControl2 TODO needs to be reverse engineered
func (p *LCPPolicy2) ParsePolicyControl2() PolicyControl {
	var polCtrl PolicyControl
	polCtrl.NPW = (p.PolicyControl>>1)&1 != 0
	polCtrl.SinitCaps = (p.PolicyControl>>2)&1 != 0
	polCtrl.AuxDelete = (p.PolicyControl>>15)&1 != 0
	polCtrl.OwnerEnforced = (p.PolicyControl>>3)&1 != 0
	return polCtrl
}

// ParseApprovedHashAlgorithm returns the supported hash algorithms
func (p *LCPPolicy2) ParseApprovedHashAlgorithm() ApprovedHashAlgorithm {
	var hashAlgs ApprovedHashAlgorithm
	hashAlgs.SHA1 = (p.LcpHashAlgMask>>0)&1 != 0
	hashAlgs.SHA256 = (p.LcpHashAlgMask>>3)&1 != 0
	hashAlgs.SHA384 = (p.LcpHashAlgMask>>6)&1 != 0
	hashAlgs.SM3 = (p.LcpHashAlgMask>>5)&1 != 0
	return hashAlgs
}

// ParseApprovedSignatureAlgorithm returns the supported signature algorithms
func (p *LCPPolicy2) ParseApprovedSignatureAlgorithm() ApprovedSignatureAlogrithm {
	var signatureAlgs ApprovedSignatureAlogrithm
	signatureAlgs.RSA2048SHA1 = (p.LcpSignAlgMask>>2)&1 != 0
	signatureAlgs.RSA2048SHA256 = (p.LcpSignAlgMask>>3)&1 != 0
	signatureAlgs.RSA3072SHA256 = (p.LcpSignAlgMask>>6)&1 != 0
	signatureAlgs.RSA3072SHA384 = (p.LcpSignAlgMask>>7)&1 != 0
	signatureAlgs.ECDSAP256SHA256 = (p.LcpSignAlgMask>>12)&1 != 0
	signatureAlgs.ECDSAP384SHA384 = (p.LcpSignAlgMask>>13)&1 != 0
	signatureAlgs.SM2SM2CurveSM3 = (p.LcpSignAlgMask>>16)&1 != 0
	return signatureAlgs
}

func parsePolicy(policy []byte) (*LCPPolicy, error) {
	var pol LCPPolicy
	buf := bytes.NewReader(policy)
	err := binary.Read(buf, binary.LittleEndian, &pol.Version)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol.HashAlg)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol.PolicyType)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol.SINITMinVersion)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol.Reserved)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol.DataRevocationCounters)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol.PolicyControl)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol.MaxSINITMinVersion)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol.Reserved1)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol.Reserved2)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol.Reserved3)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol.PolicyHash)
	if err != nil {
		return nil, err
	}

	return &pol, nil
}

func parsePolicy2(policy []byte) (*LCPPolicy2, error) {
	var pol2 LCPPolicy2
	buf := bytes.NewReader(policy)
	err := binary.Read(buf, binary.LittleEndian, &pol2.Version)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol2.HashAlg)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol2.PolicyType)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol2.SINITMinVersion)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol2.DataRevocationCounters)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol2.PolicyControl)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol2.MaxSINITMinVersion)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol2.Reserved)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol2.LcpHashAlgMask)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol2.LcpSignAlgMask)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &pol2.Reserved2)
	if err != nil {
		return nil, err
	}
	switch pol2.HashAlg {
	case LCPPol2HAlgSHA1:
		var sha1 [SHA1DigestSize]byte
		err = binary.Read(buf, binary.LittleEndian, &sha1)
		if err != nil && err != io.EOF {
			return nil, err
		}
		pol2.PolicyHash.sha1 = &sha1
		break
	case LCPPol2HAlgSHA256:
		var sha256 [SHA256DigestSize]byte
		err = binary.Read(buf, binary.LittleEndian, &sha256)
		if err != nil && err != io.EOF {
			return nil, err
		}
		pol2.PolicyHash.sha256 = &sha256
		break
	case LCPPol2HAlgSHA384:
		var sha384 [SHA384DigestSize]byte
		err = binary.Read(buf, binary.LittleEndian, &sha384)
		if err != nil && err != io.EOF {
			return nil, err
		}
		pol2.PolicyHash.sha384 = &sha384
		break
	case LCPPol2HAlgSM3:
		var sm3 [SM3DigestSize]byte
		err = binary.Read(buf, binary.LittleEndian, &sm3)
		if err != nil && err != io.EOF {
			return nil, err
		}
		pol2.PolicyHash.sm3 = &sm3
		break
	}

	return &pol2, nil
}

//ParsePolicy generates one of LCPPolicy or LCPPolicy2
func ParsePolicy(policy []byte) (*LCPPolicy, *LCPPolicy2, error) {
	var version uint16
	buf := bytes.NewReader(policy)
	err := binary.Read(buf, binary.LittleEndian, &version)
	if err != nil {
		return nil, nil, err
	}
	if version <= LCPPolicyVersion2 {
		pol, err := parsePolicy(policy)
		return pol, nil, err
	} else if version >= LCPPolicyVersion3 {
		pol, err := parsePolicy2(policy)
		return nil, pol, err
	}

	return nil, nil, fmt.Errorf("Can't parse LCP Policy")
}

func parsePolicyElement(buf *bytes.Reader, element *LCPPolicyElement) error {
	err := binary.Read(buf, binary.LittleEndian,
		&element.Size)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.LittleEndian,
		&element.Type)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.LittleEndian,
		&element.PolicyEltControl)
	if err != nil {
		return err
	}

	switch element.Type {
	case LCPPolicyElementMLE:
		var pol LCPPolicyMLE
		err = parsePolicyElementMLE(buf, &pol)
		if err != nil {
			return err
		}
		element.MLE = &pol
	case LCPPolicyElementSBIOS:
		var pol LCPPolicySBIOS
		err = parsePolicyElementSBIOS(buf, &pol)
		if err != nil {
			return err
		}
		element.SBIOS = &pol
	case LCPPolicyElementPCONF:
		var pol LCPPolicyPCONF
		err = parsePolicyElementPCONF(buf, &pol)
		if err != nil {
			return err
		}
		element.PCONF = &pol
	case LCPPolicyElementCustom:
		var pol LCPPolicyCustom
		err = parsePolicyElementCustom(buf, int(element.Size)-16, &pol)
		if err != nil {
			return err
		}
		element.Custom = &pol
	default:
		return fmt.Errorf("Unknown policy element type: %d, See: Intel TXT Software Development Guide, Document: 315168-010, P. 116", element.Type)
	}

	return nil
}

func parsePolicyElementMLE(buf *bytes.Reader, pol *LCPPolicyMLE) error {
	err := binary.Read(buf, binary.LittleEndian,
		&pol.SINITMinVersion)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian,
		&pol.HashAlg)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian,
		&pol.NumHashes)
	if err != nil {
		return err
	}

	pol.Hashes = make([][20]byte, pol.NumHashes)
	for i := 0; i < int(pol.NumHashes); i++ {
		binary.Read(buf, binary.LittleEndian, &pol.Hashes[i])
	}
	return nil
}

func parsePolicyElementSBIOS(buf *bytes.Reader, pol *LCPPolicySBIOS) error {
	err := binary.Read(buf, binary.LittleEndian,
		&pol.HashAlg)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian,
		&pol.Reserved1)
	if err != nil {
		return err
	}

	err = parseLCPHash(buf, &pol.FallbackHash, pol.HashAlg)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian,
		&pol.Reserved2)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian,
		&pol.NumHashes)
	if err != nil {
		return err
	}

	pol.Hashes = make([]LCPHash, pol.NumHashes)
	for i := 0; i < int(pol.NumHashes); i++ {
		err = parseLCPHash(buf, &pol.Hashes[i], pol.HashAlg)
		if err != nil {
			return err
		}
	}

	return nil
}

func parsePolicyElementPCONF(buf *bytes.Reader, pol *LCPPolicyPCONF) error {
	err := binary.Read(buf, binary.LittleEndian,
		&pol.NumPCRInfos)
	if err != nil {
		return err
	}

	pol.PCRInfos = make([]TPMPCRInfoShort, pol.NumPCRInfos)
	for i := 0; i < int(pol.NumPCRInfos); i++ {
		err = parseTPMPCRInfoShort(buf, &pol.PCRInfos[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func parseTPMPCRInfoShort(buf *bytes.Reader, info *TPMPCRInfoShort) error {
	var selSize uint16

	err := binary.Read(buf, binary.BigEndian,
		&selSize)
	if err != nil {
		return err
	}

	for i := 0; i < int(selSize); i++ {
		var b byte

		err = binary.Read(buf, binary.BigEndian,
			&b)
		if err != nil {
			return err
		}

		for j := 0; j < 8; j++ {
			if b&(1<<uint(j)) != 0 {
				info.PCRSelect = append(info.PCRSelect, i*8+j)
			}
		}
	}

	err = binary.Read(buf, binary.BigEndian,
		&info.LocalityAtRelease)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.BigEndian,
		&info.DigestAtRelease)
	if err != nil {
		return err
	}

	return nil
}

func parsePolicyElementCustom(buf *bytes.Reader, size int, pol *LCPPolicyCustom) error {
	err := parseLCPUUID(buf, &pol.UUID)
	if err != nil {
		return err
	}

	pol.Data = make([]byte, size-16)
	err = binary.Read(buf, binary.LittleEndian,
		&pol.Data)
	if err != nil {
		return err
	}

	return nil
}

func parseLCPUUID(buf *bytes.Reader, uuid *LCPUUID) error {
	err := binary.Read(buf, binary.LittleEndian,
		&uuid.data1)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian,
		&uuid.data2)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian,
		&uuid.data3)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian,
		&uuid.data4)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian,
		&uuid.data5)
	if err != nil {
		return err
	}

	return nil
}

func parsePolicyList(buf *bytes.Reader, list *LCPPolicyList) error {
	err := binary.Read(buf, binary.LittleEndian,
		&list.Version)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian,
		&list.Reserved)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian,
		&list.SignaturAlg)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian,
		&list.PolicyElementSize)
	if err != nil {
		return err
	}

	for i := 0; i < int(list.PolicyElementSize); {
		var elt LCPPolicyElement

		err = parsePolicyElement(buf, &elt)
		if err != nil {
			return err
		}

		list.PolicyElements = append(list.PolicyElements, elt)
		i += int(elt.Size)
	}

	switch list.SignaturAlg {
	case LCPSignatureAlgNone:
		// NOP
	case LCPSignatureAlgRSAPKCS15:
		var sig LCPSignature

		err = parseLCPSignature(buf, &sig)
		if err != nil {
			return err
		}
		list.Signature = &sig

	default:
		return fmt.Errorf("unknown signature algorithm: %x", list.SignaturAlg)
	}

	return nil
}

func parsePolicyList2(buf *bytes.Reader, list *LCPPolicyList2) error {
	err := binary.Read(buf, binary.LittleEndian,
		&list.Version)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian,
		&list.SignaturAlg)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian,
		&list.PolicyElementSize)
	if err != nil {
		return err
	}

	list.PolicyElements = make([]LCPPolicyElement, list.PolicyElementSize)
	for i := 0; i < int(list.PolicyElementSize); i++ {
		parsePolicyElement(buf, &list.PolicyElements[i])
	}

	return nil
}

func parseLCPSignature(buf *bytes.Reader, sig *LCPSignature) error {
	err := binary.Read(buf, binary.LittleEndian,
		&sig.RevocationCounter)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian,
		&sig.PubkeySize)
	if err != nil {
		return err
	}

	sig.PubkeyValue = make([]byte, sig.PubkeySize)
	err = binary.Read(buf, binary.LittleEndian,
		&sig.PubkeyValue)
	if err != nil {
		return err
	}

	sig.SigBlock = make([]byte, sig.PubkeySize)
	err = binary.Read(buf, binary.LittleEndian,
		&sig.SigBlock)
	if err != nil {
		return err
	}

	return nil
}

func parseLCPHash(buf *bytes.Reader, hash *LCPHash, alg uint8) error {
	switch alg {
	case LCPPolHAlgSHA1:
		return parseLCPHash2(buf, hash, tpm2.AlgSHA1)
	default:
		return fmt.Errorf("Unsupported hash algorithm: %x", alg)
	}
}

func parseLCPHash2(buf *bytes.Reader, hash *LCPHash, alg tpm2.Algorithm) error {
	switch alg {
	case tpm2.AlgSHA1:
		var sha1 [SHA1DigestSize]byte

		err := binary.Read(buf, binary.LittleEndian, &sha1)
		if err != nil {
			return err
		}
		hash.sha1 = &sha1

	case tpm2.AlgSHA256:
		var sha256 [SHA256DigestSize]byte

		err := binary.Read(buf, binary.LittleEndian, &sha256)
		if err != nil {
			return err
		}
		hash.sha256 = &sha256

	case tpm2.AlgSHA384:
		var sha384 [SHA384DigestSize]byte

		err := binary.Read(buf, binary.LittleEndian, &sha384)
		if err != nil {
			return err
		}
		hash.sha384 = &sha384

	case tpm2.AlgSHA512:
		var sha512 [SHA512DigestSize]byte

		err := binary.Read(buf, binary.LittleEndian, &sha512)
		if err != nil {
			return err
		}
		hash.sha512 = &sha512

	//case tpm2.AlgSM3:
	//	var sm3 [32]byte

	//	err := binary.Read(buf, binary.LittleEndian, &sm3)
	//	if err != nil {
	//		return err
	//	}
	//	hash.sm3 = &sm3

	default:
		return fmt.Errorf("unsupported hash algorithm: %x", alg)
	}

	return nil
}

//ParsePolicyData parses a raw copy of the LCP policy
func ParsePolicyData(policyData []byte) (*LCPPolicyData, error) {
	var polData LCPPolicyData

	buf := bytes.NewReader(policyData)
	err := binary.Read(buf, binary.LittleEndian, &polData.FileSignature)
	if err != nil {
		return nil, err
	}

	err = binary.Read(buf, binary.LittleEndian, &polData.Reserved)
	if err != nil {
		return nil, err
	}

	err = binary.Read(buf, binary.LittleEndian, &polData.NumLists)
	if err != nil {
		return nil, err
	}

	polData.PolicyLists = make([]LCPList, polData.NumLists)
	for i := 0; i < int(polData.NumLists); i++ {
		err = parsePolicyList(buf, &polData.PolicyLists[i].TPM12PolicyList)
		if err != nil {
			err = parsePolicyList2(buf, &polData.PolicyLists[i].TPM20PolicyList)
			if err != nil {
				return nil, err
			}
		}
	}

	return &polData, nil
}

//PrettyPrint prints the LCPHash in a human readable format
func (p *LCPHash) PrettyPrint() string {
	if p.sha1 != nil {
		return fmt.Sprintf("%02x [SHA-1]", *p.sha1)
	} else if p.sha256 != nil {
		return fmt.Sprintf("%02x [SHA-256]", *p.sha256)
	} else if p.sha384 != nil {
		return fmt.Sprintf("%02x [SHA-384]", *p.sha384)
	} else if p.sha512 != nil {
		return fmt.Sprintf("%02x [SHA-512]", *p.sha512)
		//} else if p.sm3 != nil {
		//	return fmt.Sprintf("% 02x [SM3]", *p.sm3)
	} else {
		return fmt.Sprintf("(Invalid)")
	}
}

//PrettyPrint prints the LCPPolicyData in a human readable format
func (pd *LCPPolicyData) PrettyPrint() {
	log.Printf("Launch Control Policy Data\n")

	var fileSig string
	if bytes.Equal(pd.FileSignature[:], []uint8(LCPDataFileSignature)) {
		fileSig = "valid"
	} else {
		fileSig = "invalid"
	}
	log.Printf("File Signature % x (%s)\n", pd.FileSignature, fileSig)

	log.Println("LCP Policy Lists:")
	log.Printf("\tLists: %d\n", pd.NumLists)
	for idx, pol := range pd.PolicyLists {
		log.Printf("\tList %d:\n", idx)
		log.Printf("\t\tVersion: 0x%04x\n", pol.TPM12PolicyList.Version)
		log.Printf("\t\tReserved: % 02x\n", pol.TPM12PolicyList.Reserved)
		log.Printf("\t\tSignature Algorithm: 0x%02x\n", pol.TPM12PolicyList.SignaturAlg)
		log.Printf("\t\tEntries: %d bytes\n", pol.TPM12PolicyList.PolicyElementSize)

		for jdx, ent := range pol.TPM12PolicyList.PolicyElements {
			log.Printf("\t\tPolicy %d:\n", jdx)
			log.Printf("\t\t\tSize: %d bytes\n", ent.Size)
			log.Printf("\t\t\tType: %#v\n", ent.Type)
			log.Printf("\t\t\tPolicyEltControl: %#v\n", ent.PolicyEltControl)

			if ent.MLE != nil {
				log.Printf("\t\t\tSINITMinVersion: %d\n", ent.MLE.SINITMinVersion)
				log.Printf("\t\t\tHashAlg: 0x%04x\n", ent.MLE.HashAlg)
				log.Printf("\t\t\tNumHashes: %d\n", ent.MLE.NumHashes)

				for kdx, h := range ent.MLE.Hashes {
					log.Printf("\t\t\tHash %2d: %02x\n", kdx, h)
				}
			} else if ent.SBIOS != nil {
				log.Printf("\t\t\tHashAlg: 0x%04x\n", ent.SBIOS.HashAlg)
				log.Printf("\t\t\tReserved1: % 02x\n", ent.SBIOS.Reserved1)
				log.Printf("\t\t\tFallbackHash: %s\n", ent.SBIOS.FallbackHash.PrettyPrint())
				log.Printf("\t\t\tReserved2: % 02x\n", ent.SBIOS.Reserved2)
				log.Printf("\t\t\tNumHashes: %d\n", ent.SBIOS.NumHashes)

				for kdx, h := range ent.SBIOS.Hashes {
					log.Printf("\t\t\tHash %2d: %s\n", kdx, h.PrettyPrint())
				}
			} else if ent.PCONF != nil {
				log.Printf("\t\t\tNumPCRInfos: %d\n", ent.PCONF.NumPCRInfos)

				for kdx, info := range ent.PCONF.PCRInfos {
					log.Printf("\t\t\tPCR Info %d:\n", kdx)
					log.Printf("\t\t\t\tPCR Select: %v\n", info.PCRSelect)
					log.Printf("\t\t\t\tLocality: %d\n", info.LocalityAtRelease)
					log.Printf("\t\t\t\tDigest: %02x\n", info.DigestAtRelease)
				}
			} else if ent.Custom != nil {
				log.Printf("\t\t\tUUID: %08x-%04x-%04x-%04x-%02x\n", ent.Custom.UUID.data1, ent.Custom.UUID.data2, ent.Custom.UUID.data3, ent.Custom.UUID.data4, ent.Custom.UUID.data5)
				log.Printf("\t\t\tData: %02x\n", ent.Custom.Data)
			} else {
				log.Printf("\t\t\tError: Unknown Policy Element type\n")
			}
		}

		if pol.TPM12PolicyList.Signature != nil {
			log.Printf("\t\tSignature:\n")
			log.Printf("\t\t\tRevocation Counter: %#v\n", pol.TPM12PolicyList.Signature.RevocationCounter)
			log.Printf("\t\t\tPubkey Size: %d\n", pol.TPM12PolicyList.Signature.PubkeySize)
			log.Printf("\t\t\tPubkey Value: %02x\n", pol.TPM12PolicyList.Signature.PubkeyValue)
			log.Printf("\t\t\tSig Block: %02x\n", pol.TPM12PolicyList.Signature.SigBlock)
		} else {
			log.Printf("\t\tSignature: (None)\n")
		}
	}
}
