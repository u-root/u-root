// Copyright (c) 2014, Google Inc. All rights reserved.
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

import "github.com/google/go-tpm/tpmutil"

func init() {
	// TPM 1.2 spec uses uint32 for length prefix of byte arrays.
	tpmutil.UseTPM12LengthPrefixSize()
}

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
	ordReadPubEK            uint32 = 0x0000007C
	ordOwnerReadInternalPub uint32 = 0x00000081
	ordFlushSpecific        uint32 = 0x000000BA
	ordPcrReset             uint32 = 0x000000C8
)

// Capability types.
const (
	capHandle uint32 = 0x00000014
)

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

// Algorithm ID values.
const (
	_ uint32 = iota
	algRSA
	_ // was DES
	_ // was 3DES in EDE mode
	algSHA
	algHMAC
	algAES128
	algMGF1
	algAES192
	algAES256
	algXOR
)

// Encryption schemes. The values esNone and the two that contain the string
// "RSA" are only valid under algRSA. The other two are symmetric encryption
// schemes.
const (
	_ uint16 = iota
	esNone
	esRSAEsPKCSv15
	esRSAEsOAEPSHA1MGF1
	esSymCTR
	esSymOFB
)

// Signature schemes. These are only valid under algRSA.
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
