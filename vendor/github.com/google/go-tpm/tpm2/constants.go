// Copyright (c) 2018, Google LLC All rights reserved.
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

package tpm2

import (
	"crypto"
	"crypto/elliptic"
	"fmt"

	// Register the relevant hash implementations to prevent a runtime failure.
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"

	"github.com/google/go-tpm/tpmutil"
)

var hashMapping = map[Algorithm]crypto.Hash{
	AlgSHA1:   crypto.SHA1,
	AlgSHA256: crypto.SHA256,
	AlgSHA384: crypto.SHA384,
	AlgSHA512: crypto.SHA512,
}

// MAX_DIGEST_BUFFER is the maximum size of []byte request or response fields.
// Typically used for chunking of big blobs of data (such as for hashing or
// encryption).
const maxDigestBuffer = 1024

// Algorithm represents a TPM_ALG_ID value.
type Algorithm uint16

// IsNull returns true if a is AlgNull or zero (unset).
func (a Algorithm) IsNull() bool {
	return a == AlgNull || a == AlgUnknown
}

// UsesCount returns true if a signature algorithm uses count value.
func (a Algorithm) UsesCount() bool {
	return a == AlgECDAA
}

// UsesHash returns true if the algorithm requires the use of a hash.
func (a Algorithm) UsesHash() bool {
	return a == AlgOAEP
}

// Hash returns a crypto.Hash based on the given TPM_ALG_ID.
// An error is returned if the given algorithm is not a hash algorithm or is not available.
func (a Algorithm) Hash() (crypto.Hash, error) {
	hash, ok := hashMapping[a]
	if !ok {
		return crypto.Hash(0), fmt.Errorf("hash algorithm not supported: 0x%x", a)
	}
	if !hash.Available() {
		return crypto.Hash(0), fmt.Errorf("go hash algorithm #%d not available", hash)
	}
	return hash, nil
}

// Supported Algorithms.
const (
	AlgUnknown   Algorithm = 0x0000
	AlgRSA       Algorithm = 0x0001
	AlgSHA1      Algorithm = 0x0004
	AlgHMAC      Algorithm = 0x0005
	AlgAES       Algorithm = 0x0006
	AlgKeyedHash Algorithm = 0x0008
	AlgXOR       Algorithm = 0x000A
	AlgSHA256    Algorithm = 0x000B
	AlgSHA384    Algorithm = 0x000C
	AlgSHA512    Algorithm = 0x000D
	AlgNull      Algorithm = 0x0010
	AlgRSASSA    Algorithm = 0x0014
	AlgRSAES     Algorithm = 0x0015
	AlgRSAPSS    Algorithm = 0x0016
	AlgOAEP      Algorithm = 0x0017
	AlgECDSA     Algorithm = 0x0018
	AlgECDH      Algorithm = 0x0019
	AlgECDAA     Algorithm = 0x001A
	AlgKDF2      Algorithm = 0x0021
	AlgECC       Algorithm = 0x0023
	AlgSymCipher Algorithm = 0x0025
	AlgCTR       Algorithm = 0x0040
	AlgOFB       Algorithm = 0x0041
	AlgCBC       Algorithm = 0x0042
	AlgCFB       Algorithm = 0x0043
	AlgECB       Algorithm = 0x0044
)

// HandleType defines a type of handle.
type HandleType uint8

// Supported handle types
const (
	HandleTypePCR           HandleType = 0x00
	HandleTypeNVIndex       HandleType = 0x01
	HandleTypeHMACSession   HandleType = 0x02
	HandleTypeLoadedSession HandleType = 0x02
	HandleTypePolicySession HandleType = 0x03
	HandleTypeSavedSession  HandleType = 0x03
	HandleTypePermanent     HandleType = 0x40
	HandleTypeTransient     HandleType = 0x80
	HandleTypePersistent    HandleType = 0x81
)

// SessionType defines the type of session created in StartAuthSession.
type SessionType uint8

// Supported session types.
const (
	SessionHMAC   SessionType = 0x00
	SessionPolicy SessionType = 0x01
	SessionTrial  SessionType = 0x03
)

// SessionAttributes represents an attribute of a session.
type SessionAttributes byte

// Session Attributes (Structures 8.4 TPMA_SESSION)
const (
	AttrContinueSession SessionAttributes = 1 << iota
	AttrAuditExclusive
	AttrAuditReset
	_ // bit 3 reserved
	_ // bit 4 reserved
	AttrDecrypt
	AttrEcrypt
	AttrAudit
)

// EmptyAuth represents the empty authorization value.
var EmptyAuth []byte

// KeyProp is a bitmask used in Attributes field of key templates. Individual
// flags should be OR-ed to form a full mask.
type KeyProp uint32

// Key properties.
const (
	FlagFixedTPM            KeyProp = 0x00000002
	FlagFixedParent         KeyProp = 0x00000010
	FlagSensitiveDataOrigin KeyProp = 0x00000020
	FlagUserWithAuth        KeyProp = 0x00000040
	FlagAdminWithPolicy     KeyProp = 0x00000080
	FlagNoDA                KeyProp = 0x00000400
	FlagRestricted          KeyProp = 0x00010000
	FlagDecrypt             KeyProp = 0x00020000
	FlagSign                KeyProp = 0x00040000

	FlagSealDefault   = FlagFixedTPM | FlagFixedParent
	FlagSignerDefault = FlagSign | FlagRestricted | FlagFixedTPM |
		FlagFixedParent | FlagSensitiveDataOrigin | FlagUserWithAuth
	FlagStorageDefault = FlagDecrypt | FlagRestricted | FlagFixedTPM |
		FlagFixedParent | FlagSensitiveDataOrigin | FlagUserWithAuth
)

// TPMProp represents a Property Tag (TPM_PT) used with calls to GetCapability(CapabilityTPMProperties).
type TPMProp uint32

// TPM Capability Properties, see TPM 2.0 Spec, Rev 1.38, Table 23.
// Fixed TPM Properties (PT_FIXED)
const (
	FamilyIndicator TPMProp = 0x100 + iota
	SpecLevel
	SpecRevision
	SpecDayOfYear
	SpecYear
	Manufacturer
	VendorString1
	VendorString2
	VendorString3
	VendorString4
	VendorTPMType
	FirmwareVersion1
	FirmwareVersion2
	InputMaxBufferSize
	TransientObjectsMin
	PersistentObjectsMin
	LoadedObjectsMin
	ActiveSessionsMax
	PCRCount
	PCRSelectMin
	ContextGapMax
	_ // (PT_FIXED + 21) is skipped
	NVCountersMax
	NVIndexMax
	MemoryMethod
	ClockUpdate
	ContextHash
	ContextSym
	ContextSymSize
	OrderlyCount
	CommandMaxSize
	ResponseMaxSize
	DigestMaxSize
	ObjectContextMaxSize
	SessionContextMaxSize
	PSFamilyIndicator
	PSSpecLevel
	PSSpecRevision
	PSSpecDayOfYear
	PSSpecYear
	SplitSigningMax
	TotalCommands
	LibraryCommands
	VendorCommands
	NVMaxBufferSize
	TPMModes
	CapabilityMaxBufferSize
)

// Variable TPM Properties (PT_VAR)
const (
	TPMAPermanent TPMProp = 0x200 + iota
	TPMAStartupClear
	HRNVIndex
	HRLoaded
	HRLoadedAvail
	HRActive
	HRActiveAvail
	HRTransientAvail
	CurrentPersistent
	AvailPersistent
	NVCounters
	NVCountersAvail
	AlgorithmSet
	LoadedCurves
	LockoutCounter
	MaxAuthFail
	LockoutInterval
	LockoutRecovery
	NVWriteRecovery
	AuditCounter0
	AuditCounter1
)

// Allowed ranges of different kinds of Handles (TPM_HANDLE)
// These constants have type TPMProp for backwards compatibility.
const (
	PCRFirst           TPMProp = 0x00000000
	HMACSessionFirst   TPMProp = 0x02000000
	LoadedSessionFirst TPMProp = 0x02000000
	PolicySessionFirst TPMProp = 0x03000000
	ActiveSessionFirst TPMProp = 0x03000000
	TransientFirst     TPMProp = 0x80000000
	PersistentFirst    TPMProp = 0x81000000
	PersistentLast     TPMProp = 0x81FFFFFF
	PlatformPersistent TPMProp = 0x81800000
	NVIndexFirst       TPMProp = 0x01000000
	NVIndexLast        TPMProp = 0x01FFFFFF
	PermanentFirst     TPMProp = 0x40000000
	PermanentLast      TPMProp = 0x4000010F
)

// Reserved Handles.
const (
	HandleOwner tpmutil.Handle = 0x40000001 + iota
	HandleRevoke
	HandleTransport
	HandleOperator
	HandleAdmin
	HandleEK
	HandleNull
	HandleUnassigned
	HandlePasswordSession
	HandleLockout
	HandleEndorsement
	HandlePlatform
)

// Capability identifies some TPM property or state type.
type Capability uint32

// TPM Capabilities.
const (
	CapabilityAlgs Capability = iota
	CapabilityHandles
	CapabilityCommands
	CapabilityPPCommands
	CapabilityAuditCommands
	CapabilityPCRs
	CapabilityTPMProperties
	CapabilityPCRProperties
	CapabilityECCCurves
	CapabilityAuthPolicies
)

// TPM Structure Tags. Tags are used to disambiguate structures, similar to Alg
// values: tag value defines what kind of data lives in a nested field.
const (
	TagNull           tpmutil.Tag = 0x8000
	TagNoSessions     tpmutil.Tag = 0x8001
	TagSessions       tpmutil.Tag = 0x8002
	TagAttestCertify  tpmutil.Tag = 0x8017
	TagAttestQuote    tpmutil.Tag = 0x8018
	TagAttestCreation tpmutil.Tag = 0x801a
	TagHashCheck      tpmutil.Tag = 0x8024
)

// StartupType instructs the TPM on how to handle its state during Shutdown or
// Startup.
type StartupType uint16

// Startup types
const (
	StartupClear StartupType = iota
	StartupState
)

// EllipticCurve identifies specific EC curves.
type EllipticCurve uint16

// ECC curves supported by TPM 2.0 spec.
const (
	CurveNISTP192 = EllipticCurve(iota + 1)
	CurveNISTP224
	CurveNISTP256
	CurveNISTP384
	CurveNISTP521

	CurveBNP256 = EllipticCurve(iota + 10)
	CurveBNP638

	CurveSM2P256 = EllipticCurve(0x0020)
)

var toGoCurve = map[EllipticCurve]elliptic.Curve{
	CurveNISTP224: elliptic.P224(),
	CurveNISTP256: elliptic.P256(),
	CurveNISTP384: elliptic.P384(),
	CurveNISTP521: elliptic.P521(),
}

// Supported TPM operations.
const (
	cmdEvictControl               tpmutil.Command = 0x00000120
	cmdUndefineSpace              tpmutil.Command = 0x00000122
	cmdClear                      tpmutil.Command = 0x00000126
	cmdHierarchyChangeAuth        tpmutil.Command = 0x00000129
	cmdDefineSpace                tpmutil.Command = 0x0000012A
	cmdCreatePrimary              tpmutil.Command = 0x00000131
	cmdIncrementNVCounter         tpmutil.Command = 0x00000134
	cmdWriteNV                    tpmutil.Command = 0x00000137
	cmdWriteLockNV                tpmutil.Command = 0x00000138
	cmdDictionaryAttackLockReset  tpmutil.Command = 0x00000139
	cmdDictionaryAttackParameters tpmutil.Command = 0x0000013A
	cmdPCREvent                   tpmutil.Command = 0x0000013C
	cmdStartup                    tpmutil.Command = 0x00000144
	cmdShutdown                   tpmutil.Command = 0x00000145
	cmdActivateCredential         tpmutil.Command = 0x00000147
	cmdCertify                    tpmutil.Command = 0x00000148
	cmdCertifyCreation            tpmutil.Command = 0x0000014A
	cmdReadNV                     tpmutil.Command = 0x0000014E
	cmdReadLockNV                 tpmutil.Command = 0x0000014F
	// CmdPolicySecret is a command code for TPM2_PolicySecret.
	// It's exported for computing of default AuthPolicy value.
	CmdPolicySecret     tpmutil.Command = 0x00000151
	cmdCreate           tpmutil.Command = 0x00000153
	cmdImport           tpmutil.Command = 0x00000156
	cmdLoad             tpmutil.Command = 0x00000157
	cmdQuote            tpmutil.Command = 0x00000158
	cmdRSADecrypt       tpmutil.Command = 0x00000159
	cmdSign             tpmutil.Command = 0x0000015D
	cmdUnseal           tpmutil.Command = 0x0000015E
	cmdContextLoad      tpmutil.Command = 0x00000161
	cmdContextSave      tpmutil.Command = 0x00000162
	cmdEncryptDecrypt   tpmutil.Command = 0x00000164
	cmdFlushContext     tpmutil.Command = 0x00000165
	cmdLoadExternal     tpmutil.Command = 0x00000167
	cmdMakeCredential   tpmutil.Command = 0x00000168
	cmdReadPublicNV     tpmutil.Command = 0x00000169
	cmdReadPublic       tpmutil.Command = 0x00000173
	cmdRSAEncrypt       tpmutil.Command = 0x00000174
	cmdStartAuthSession tpmutil.Command = 0x00000176
	cmdGetCapability    tpmutil.Command = 0x0000017A
	cmdGetRandom        tpmutil.Command = 0x0000017B
	cmdHash             tpmutil.Command = 0x0000017D
	cmdPCRRead          tpmutil.Command = 0x0000017E
	// CmdPolicyPCR is the command code for TPM2_PolicyPCR.
	// It's exported for computing AuthPolicy values for PCR-based sessions.
	CmdPolicyPCR       tpmutil.Command = 0x0000017F
	cmdReadClock       tpmutil.Command = 0x00000181
	cmdPCRExtend       tpmutil.Command = 0x00000182
	cmdPolicyGetDigest tpmutil.Command = 0x00000189
	cmdPolicyPassword  tpmutil.Command = 0x0000018C
	cmdEncryptDecrypt2 tpmutil.Command = 0x00000193
)

// Regular TPM 2.0 devices use 24-bit mask (3 bytes) for PCR selection.
const sizeOfPCRSelect = 3

const defaultRSAExponent = 1<<16 + 1

// NVAttr is a bitmask used in Attributes field of NV indexes. Individual
// flags should be OR-ed to form a full mask.
type NVAttr uint32

// NV Attributes
const (
	AttrPPWrite        NVAttr = 0x00000001
	AttrOwnerWrite     NVAttr = 0x00000002
	AttrAuthWrite      NVAttr = 0x00000004
	AttrPolicyWrite    NVAttr = 0x00000008
	AttrPolicyDelete   NVAttr = 0x00000400
	AttrWriteLocked    NVAttr = 0x00000800
	AttrWriteAll       NVAttr = 0x00001000
	AttrWriteDefine    NVAttr = 0x00002000
	AttrWriteSTClear   NVAttr = 0x00004000
	AttrGlobalLock     NVAttr = 0x00008000
	AttrPPRead         NVAttr = 0x00010000
	AttrOwnerRead      NVAttr = 0x00020000
	AttrAuthRead       NVAttr = 0x00040000
	AttrPolicyRead     NVAttr = 0x00080000
	AttrNoDA           NVAttr = 0x02000000
	AttrOrderly        NVAttr = 0x04000000
	AttrClearSTClear   NVAttr = 0x08000000
	AttrReadLocked     NVAttr = 0x10000000
	AttrWritten        NVAttr = 0x20000000
	AttrPlatformCreate NVAttr = 0x40000000
	AttrReadSTClear    NVAttr = 0x80000000
)
