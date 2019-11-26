// Copyright (c) 2014, Google LLC All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tpm

import (
	"fmt"
	"github.com/google/go-tpm/tpmutil"
	"strings"
)

// Supported TPM commands.
const (
	tagPCRInfoLong     uint16 = 0x06
	tagRQUCommand      uint16 = 0x00C1
	tagRQUAuth1Command uint16 = 0x00C2
	tagRQUAuth2Command uint16 = 0x00C3
	tagRSPCommand      uint16 = 0x00C4
	tagRSPAuth1Command uint16 = 0x00C5
	tagRSPAuth2Command uint16 = 0x00C6
)

// Supported TPM operations.
const (
	ordOIAP                 uint32 = 0x0000000A
	ordOSAP                 uint32 = 0x0000000B
	ordTakeOwnership        uint32 = 0x0000000D
	ordExtend               uint32 = 0x00000014
	ordPCRRead              uint32 = 0x00000015
	ordQuote                uint32 = 0x00000016
	ordSeal                 uint32 = 0x00000017
	ordUnseal               uint32 = 0x00000018
	ordCreateWrapKey        uint32 = 0x0000001F
	ordGetPubKey            uint32 = 0x00000021
	ordSign                 uint32 = 0x0000003C
	ordQuote2               uint32 = 0x0000003E
	ordResetLockValue       uint32 = 0x00000040
	ordLoadKey2             uint32 = 0x00000041
	ordGetRandom            uint32 = 0x00000046
	ordOwnerClear           uint32 = 0x0000005B
	ordForceClear           uint32 = 0x0000005D
	ordGetCapability        uint32 = 0x00000065
	ordMakeIdentity         uint32 = 0x00000079
	ordActivateIdentity     uint32 = 0x0000007A
	ordReadPubEK            uint32 = 0x0000007C
	ordOwnerReadInternalPub uint32 = 0x00000081
	ordFlushSpecific        uint32 = 0x000000BA
	ordPcrReset             uint32 = 0x000000C8
	ordNVReadValue          uint32 = 0x000000CF
)

// Capability types.
const (
	capAlg      uint32 = 0x00000002
	capProperty uint32 = 0x00000005
	capFlag     uint32 = 0x00000004
	capNVList   uint32 = 0x0000000D
	capNVIndex  uint32 = 0x00000011
	capHandle   uint32 = 0x00000014
)

// SubCapabilities
const (
	tpmCapPropManufacturer uint32 = 0x00000103
	tpmCapFlagPermanent    uint32 = 0x00000108
)

// permission type
type permission uint32

// NV Permissions and Operations
// Note: Permissions are summable
const (
	nvPerPPWrite    permission = 0x00000001
	nvPerOwnerWrite permission = 0x00000002
	nvPerAuthWrite  permission = 0x00000004
	nvPerWriteAll   permission = 0x00000800
	// Warning: The Value 0x00001000 is
	// defined in the spec as
	// TPM_NV_PER_WRITEDEFINE, but it is
	// not included directly in this
	// code because it locks the given
	// NV Index permanently if used
	// incorrectly. This operation can't
	// be undone in any way. Do not use
	// this value unless you know what
	// you're doing!
	nvPerWriteSTClear permission = 0x00002000
	nvPerGlobalLock   permission = 0x00004000
	nvPerPPRead       permission = 0x00008000
	nvPerOwnerRead    permission = 0x00100000
	nvPerAuthRead     permission = 0x00200000
	nvPerReadSTClear  permission = 0x80000000
)

// permMap : Map of TPM_NV_Permissions to its strings for convenience
var permMap = map[permission]string{
	nvPerPPWrite:      "PPWrite",
	nvPerOwnerWrite:   "OwnerWrite",
	nvPerAuthWrite:    "AuthWrite",
	nvPerWriteAll:     "WriteAll",
	nvPerWriteSTClear: " WriteSTClear",
	nvPerGlobalLock:   "GlobalLock",
	nvPerPPRead:       "PPRead",
	nvPerOwnerRead:    "OwnerRead",
	nvPerAuthRead:     "AuthRead",
	nvPerReadSTClear:  "ReadSTClear",
}

// String returns a textual representation of the set of permissions
func (p permission) String() string {
	var retString strings.Builder
	for iterator, item := range permMap {
		if (p & iterator) != 0 {
			retString.WriteString(item + " + ")
		}
	}
	if retString.String() == "" {
		return "Permission/s not found"
	}
	return strings.TrimSuffix(retString.String(), " + ")

}

// Entity types. The LSB gives the entity type, and the MSB (currently fixed to
// 0x00) gives the ADIP type. ADIP type 0x00 is XOR.
const (
	_ uint16 = iota
	etKeyHandle
	etOwner
	etData
	etSRK
	etKey
	etRevoke
)

// Resource types.
const (
	_ uint32 = iota
	rtKey
	rtAuth
	rtHash
	rtTrans
)

// Locality type
type Locality byte

// Values of locality
// Note: Localities are summable
const (
	locZero Locality = 1 << iota
	locOne
	locTwo
	locThree
	locFour
)

// LocaMap maps Locality values to strings for convenience
var locaMap = map[Locality]string{
	locZero:  "Locality 0",
	locOne:   "Locality 1",
	locTwo:   "Locality 2",
	locThree: "Locality 3",
	locFour:  "Locality 4",
}

// PrettyPrint for Localities
func (l Locality) String() string {
	var retString strings.Builder
	for iterator, item := range locaMap {
		if l&iterator != 0 {
			retString.WriteString(item + " + ")
		}
	}
	if retString.String() == "" {
		return fmt.Sprintf("locality %d", int(l))
	}
	return strings.TrimSuffix(retString.String(), " + ")
}

// Entity values.
const (
	khSRK         tpmutil.Handle = 0x40000000
	khOwner       tpmutil.Handle = 0x40000001
	khRevokeTrust tpmutil.Handle = 0x40000002
	khEK          tpmutil.Handle = 0x40000006
)

// Protocol IDs.
const (
	_ uint16 = iota
	pidOIAP
	pidOSAP
	pidADIP
	pidADCP
	pidOwner
	pidDSAP
	pidTransport
)

// Algorithm type for more convenient handling.
// see Algorithm ID for possible values.
type Algorithm uint32

// Algorithm ID values.
const (
	_ Algorithm = iota
	AlgRSA
	_ // was DES
	_ // was 3DES in EDE mode
	AlgSHA
	AlgHMAC
	AlgAES128
	AlgMGF1
	AlgAES192
	AlgAES256
	AlgXOR
)

// AlgMap Map of Algorithms to Strings for nicer output and comparisons, etc.
var AlgMap = map[Algorithm]string{
	AlgRSA:    "RSA",
	AlgSHA:    "SHA1",
	AlgHMAC:   "HMAC",
	AlgAES128: "AER128",
	AlgMGF1:   "MFG1",
	AlgAES192: "AES192",
	AlgAES256: "AES256",
}

func (a Algorithm) String() string {
	n, ok := AlgMap[a]
	if !ok {
		return "unknown_algorithm"
	}
	return n
}

// Encryption schemes. The values esNone and the two that contain the string
// "RSA" are only valid under AlgRSA. The other two are symmetric encryption
// schemes.
const (
	_ uint16 = iota
	esNone
	esRSAEsPKCSv15
	esRSAEsOAEPSHA1MGF1
	esSymCTR
	esSymOFB
	esSymCBCPKCS5 = 0xff // esSymCBCPKCS5 was taken from go-tspi
)

// Signature schemes. These are only valid under AlgRSA.
const (
	_ uint16 = iota
	ssNone
	ssRSASaPKCS1v15SHA1
	ssRSASaPKCS1v15DER
	ssRSASaPKCS1v15INFO
)

// KeyUsage types for TPM_KEY (the key type).
const (
	keySigning    uint16 = 0x0010
	keyStorage    uint16 = 0x0011
	keyIdentity   uint16 = 0x0012
	keyAuthChange uint16 = 0x0013
	keyBind       uint16 = 0x0014
	keyLegacy     uint16 = 0x0015
	keyMigrate    uint16 = 0x0016
)

const (
	authNever       byte = 0x00
	authAlways      byte = 0x01
	authPrivUseOnly byte = 0x03
)

// fixedQuote is the fixed constant string used in quoteInfo.
var fixedQuote = [4]byte{byte('Q'), byte('U'), byte('O'), byte('T')}

// quoteVersion is the fixed version string for quoteInfo.
const quoteVersion uint32 = 0x01010000

// oaepLabel is the label used for OEAP encryption in esRSAEsOAEPSHA1MGF1
var oaepLabel = []byte{byte('T'), byte('C'), byte('P'), byte('A')}
