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

import (
	"crypto"
	"crypto/rsa"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/google/go-tpm/tpmutil"
)

// A pcrValue is the fixed-size value of a PCR.
type pcrValue [20]byte

// PCRSize gives the fixed size (20 bytes) of a PCR.
const PCRSize int = 20

// A pcrMask represents a set of PCR choices, one bit per PCR out of the 24
// possible PCR values.
type pcrMask [3]byte

// A pcrSelection is the first element in the input a PCR composition, which is
// A pcrSelection, followed by the combined length of the PCR values,
// followed by the PCR values, all hashed under SHA-1.
type pcrSelection struct {
	Size uint16
	Mask pcrMask
}

// pcrInfoLong stores detailed information about PCRs.
type pcrInfoLong struct {
	Tag              uint16
	LocAtCreation    byte
	LocAtRelease     byte
	PCRsAtCreation   pcrSelection
	PCRsAtRelease    pcrSelection
	DigestAtCreation digest
	DigestAtRelease  digest
}

// pcrInfoShort stores detailed information about PCRs.
type pcrInfoShort struct {
	LocAtRelease    byte
	PCRsAtRelease   pcrSelection
	DigestAtRelease digest
}

type pcrInfo struct {
	PcrSelection     pcrSelection
	DigestAtRelease  digest
	DigestAtCreation digest
}

// A capVersionInfo contains information about the TPM itself. Note that this
// is deserialized specially, since it has a variable-length byte array but no
// length. It is preceeded with a length in the response to the Quote2 command.
type capVersionInfo struct {
	CapVersionFixed capVersionInfoFixed
	VendorSpecific  []byte
}

// A capVersionInfoFixed stores the fixed-length part of capVersionInfo.
type capVersionInfoFixed struct {
	Tag       uint16
	Version   uint32
	SpecLevel uint16
	ErrataRev byte
	VendorID  byte
}

// CloseKey flushes the key associated with the tpmutil.Handle.
func CloseKey(rw io.ReadWriter, h tpmutil.Handle) error {
	return flushSpecific(rw, h, rtKey)
}

// A nonce is a 20-byte value.
type nonce [20]byte

const nonceSize uint32 = 20

// An oiapResponse is a response to an OIAP command.
type oiapResponse struct {
	AuthHandle tpmutil.Handle
	NonceEven  nonce
}

// String returns a string representation of an oiapResponse.
func (opr oiapResponse) String() string {
	return fmt.Sprintf("oiapResponse{AuthHandle: %x, NonceEven: % x}", opr.AuthHandle, opr.NonceEven)
}

// Close flushes the auth handle associated with an OIAP session.
func (opr *oiapResponse) Close(rw io.ReadWriter) error {
	return flushSpecific(rw, opr.AuthHandle, rtAuth)
}

// An osapCommand is a command sent for OSAP authentication.
type osapCommand struct {
	EntityType  uint16
	EntityValue tpmutil.Handle
	OddOSAP     nonce
}

// String returns a string representation of an osapCommand.
func (opc osapCommand) String() string {
	return fmt.Sprintf("osapCommand{EntityType: %x, EntityValue: %x, OddOSAP: % x}", opc.EntityType, opc.EntityValue, opc.OddOSAP)
}

// An osapResponse is a TPM reply to an osapCommand.
type osapResponse struct {
	AuthHandle tpmutil.Handle
	NonceEven  nonce
	EvenOSAP   nonce
}

// String returns a string representation of an osapResponse.
func (opr osapResponse) String() string {
	return fmt.Sprintf("osapResponse{AuthHandle: %x, NonceEven: % x, EvenOSAP: % x}", opr.AuthHandle, opr.NonceEven, opr.EvenOSAP)
}

// Close flushes the AuthHandle associated with an OSAP session.
func (opr *osapResponse) Close(rw io.ReadWriter) error {
	return flushSpecific(rw, opr.AuthHandle, rtAuth)
}

// A Digest is a 20-byte SHA1 value.
type digest [20]byte

const digestSize uint32 = 20

// An AuthValue is a 20-byte value used for authentication.
type authValue [20]byte

const authSize uint32 = 20

// A sealCommand is the command sent to the TPM to seal data.
type sealCommand struct {
	KeyHandle tpmutil.Handle
	EncAuth   authValue
}

// String returns a string representation of a sealCommand.
func (sc sealCommand) String() string {
	return fmt.Sprintf("sealCommand{KeyHandle: %x, EncAuth: % x}", sc.KeyHandle, sc.EncAuth)
}

// commandAuth stores the auth information sent with a command. Commands with
// tagRQUAuth1Command tags use one of these auth structures, and commands with
// tagRQUAuth2Command use two.
type commandAuth struct {
	AuthHandle  tpmutil.Handle
	NonceOdd    nonce
	ContSession byte
	Auth        authValue
}

// String returns a string representation of a sealCommandAuth.
func (ca commandAuth) String() string {
	return fmt.Sprintf("commandAuth{AuthHandle: %x, NonceOdd: % x, ContSession: %x, Auth: % x}", ca.AuthHandle, ca.NonceOdd, ca.ContSession, ca.Auth)
}

// responseAuth contains the auth information returned from a command.
type responseAuth struct {
	NonceEven   nonce
	ContSession byte
	Auth        authValue
}

// String returns a string representation of a responseAuth.
func (ra responseAuth) String() string {
	return fmt.Sprintf("responseAuth{NonceEven: % x, ContSession: %x, Auth: % x}", ra.NonceEven, ra.ContSession, ra.Auth)
}

// These are the parameters of a TPM key.
type keyParms struct {
	AlgID     uint32
	EncScheme uint16
	SigScheme uint16
	Parms     []byte // Serialized rsaKeyParms or symmetricKeyParms.
}

// An rsaKeyParms encodes the length of the RSA prime in bits, the number of
// primes in its factored form, and the exponent used for public-key
// encryption.
type rsaKeyParms struct {
	KeyLength uint32
	NumPrimes uint32
	Exponent  []byte
}

type symmetricKeyParms struct {
	KeyLength uint32
	BlockSize uint32
	IV        []byte
}

// A key is a TPM representation of a key.
type key struct {
	Version        uint32
	KeyUsage       uint16
	KeyFlags       uint32
	AuthDataUsage  byte
	AlgorithmParms keyParms
	PCRInfo        []byte
	PubKey         []byte
	EncData        []byte
}

// A key12 is a newer TPM representation of a key.
type key12 struct {
	Tag            uint16
	Zero           uint16 // Always all 0.
	KeyUsage       uint16
	KeyFlags       uint32
	AuthDataUsage  byte
	AlgorithmParms keyParms
	PCRInfo        []byte // This must be a serialization of a pcrInfoLong.
	PubKey         []byte
	EncData        []byte
}

// A pubKey represents a public key known to the TPM.
type pubKey struct {
	AlgorithmParms keyParms
	Key            []byte
}

// A tpmStoredData holds sealed data from the TPM.
type tpmStoredData struct {
	Version uint32
	Info    []byte
	Enc     []byte
}

// String returns a string representation of a tpmStoredData.
func (tsd tpmStoredData) String() string {
	return fmt.Sprintf("tpmStoreddata{Version: %x, Info: % x, Enc: % x\n", tsd.Version, tsd.Info, tsd.Enc)
}

// A quoteInfo structure is the structure signed by the TPM.
type quoteInfo struct {
	// The Version must be 0x01010000
	Version uint32

	// Fixed is always 'QUOT'.
	Fixed [4]byte

	// The CompositeDigest is computed by ComputePCRComposite.
	CompositeDigest digest

	// The nonce is either a random nonce or the SHA1 hash of data to sign.
	Nonce nonce
}

// A pcrComposite stores a selection of PCRs with the selected PCR values.
type pcrComposite struct {
	Selection pcrSelection
	Values    []byte
}

// convertPubKey converts a public key into TPM form. Currently, this function
// only supports 2048-bit RSA keys.
func convertPubKey(pk crypto.PublicKey) (*pubKey, error) {
	pkRSA, ok := pk.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("the provided Privacy CA public key was not an RSA key")
	}
	if pkRSA.N.BitLen() != 2048 {
		return nil, errors.New("The provided Privacy CA RSA public key was not a 2048-bit key")
	}

	rsakp := rsaKeyParms{
		KeyLength: 2048,
		NumPrimes: 2,
		Exponent:  big.NewInt(int64(pkRSA.E)).Bytes(),
	}
	rsakpb, err := tpmutil.Pack(rsakp)
	if err != nil {
		return nil, err
	}
	kp := keyParms{
		AlgID:     algRSA,
		EncScheme: esNone,
		SigScheme: ssRSASaPKCS1v15SHA1,
		Parms:     rsakpb,
	}
	pubk := pubKey{
		AlgorithmParms: kp,
		Key:            pkRSA.N.Bytes(),
	}

	return &pubk, nil
}
