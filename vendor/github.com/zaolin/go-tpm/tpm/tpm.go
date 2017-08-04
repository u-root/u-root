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

// Package tpm supports direct communication with a tpm device under Linux.
package tpm

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

// OpenTPM opens a channel to the TPM at the given path. If the file is a
// device, then it treats it like a normal TPM device, and if the file is a
// Unix domain socket, then it opens a connection to the socket.
func OpenTPM(path string) (io.ReadWriteCloser, error) {
	// If it's a regular file, then open it
	var rwc io.ReadWriteCloser
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fi.Mode()&os.ModeDevice != 0 {
		var f *os.File
		f, err = os.OpenFile(path, os.O_RDWR, 0600)
		if err != nil {
			return nil, err
		}
		rwc = io.ReadWriteCloser(f)
	} else if fi.Mode()&os.ModeSocket != 0 {
		uc, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: path, Net: "unix"})
		if err != nil {
			return nil, err
		}
		rwc = io.ReadWriteCloser(uc)
	} else {
		return nil, fmt.Errorf("unsupported TPM file mode %s", fi.Mode().String())
	}

	return rwc, nil
}

// GetKeys gets the list of handles for currently-loaded TPM keys.
func GetKeys(rw io.ReadWriter) ([]Handle, error) {
	var b []byte
	subCap, err := pack([]interface{}{rtKey})
	if err != nil {
		return nil, err
	}
	in := []interface{}{capHandle, subCap}
	out := []interface{}{&b}
	if _, err := submitTPMRequest(rw, tagRQUCommand, ordGetCapability, in, out); err != nil {
		return nil, err
	}
	return unpackKeyHandleList(b)
}

// PcrExtend extends a value into the right PCR by index.
func PcrExtend(rw io.ReadWriter, i uint32, p pcrValue) ([]byte, error) {
	in := []interface{}{i, p}
	var d pcrValue
	out := []interface{}{&d}
	_, err := submitTPMRequest(rw, tagRQUCommand, ordExtend, in, out)
	if err != nil {
		return nil, err
	}

	return d[:], nil
}

// ReadPCR reads a PCR value from the TPM.
func ReadPCR(rw io.ReadWriter, i uint32) ([]byte, error) {
	in := []interface{}{i}
	var v pcrValue
	out := []interface{}{&v}
	// There's no need to check the ret value here, since the err value contains
	// all the necessary information.
	if _, err := submitTPMRequest(rw, tagRQUCommand, ordPCRRead, in, out); err != nil {
		return nil, err
	}

	return v[:], nil
}

// FetchPCRValues gets a given sequence of PCR values.
func FetchPCRValues(rw io.ReadWriter, pcrVals []int) ([]byte, error) {
	var pcrs []byte
	for _, v := range pcrVals {
		pcr, err := ReadPCR(rw, uint32(v))
		if err != nil {
			return nil, err
		}

		pcrs = append(pcrs, pcr...)
	}

	return pcrs, nil
}

// GetRandom gets random bytes from the TPM.
func GetRandom(rw io.ReadWriter, size uint32) ([]byte, error) {
	var b []byte
	in := []interface{}{size}
	out := []interface{}{&b}
	// There's no need to check the ret value here, since the err value
	// contains all the necessary information.
	if _, err := submitTPMRequest(rw, tagRQUCommand, ordGetRandom, in, out); err != nil {
		return nil, err
	}

	return b, nil
}

// LoadKey2 loads a key blob (a serialized TPM_KEY or TPM_KEY12) into the TPM
// and returns a handle for this key.
func LoadKey2(rw io.ReadWriter, keyBlob []byte, srkAuth []byte) (Handle, error) {
	// Deserialize the keyBlob as a key
	var k key
	if err := unpack(keyBlob, []interface{}{&k}); err != nil {
		return 0, err
	}

	// Run OSAP for the SRK, reading a random OddOSAP for our initial
	// command and getting back a secret and a handle. LoadKey2 needs an
	// OSAP session for the SRK because the private part of a TPM_KEY or
	// TPM_KEY12 is sealed against the SRK.
	sharedSecret, osapr, err := newOSAPSession(rw, etSRK, khSRK, srkAuth)
	if err != nil {
		return 0, err
	}
	defer osapr.Close(rw)
	defer zeroBytes(sharedSecret[:])

	authIn := []interface{}{ordLoadKey2, k}
	ca, err := newCommandAuth(osapr.AuthHandle, osapr.NonceEven, sharedSecret[:], authIn)
	if err != nil {
		return 0, err
	}

	handle, ra, ret, err := loadKey2(rw, &k, ca)
	if err != nil {
		return 0, err
	}

	// Check the response authentication.
	raIn := []interface{}{ret, ordLoadKey2}
	if err := ra.verify(ca.NonceOdd, sharedSecret[:], raIn); err != nil {
		return 0, err
	}

	return handle, nil
}

// Quote2 performs a quote operation on the TPM for the given data,
// under the key associated with the handle and for the pcr values
// specified in the call.
func Quote2(rw io.ReadWriter, handle Handle, data []byte, pcrVals []int, addVersion byte, aikAuth []byte) ([]byte, error) {
	// Run OSAP for the handle, reading a random OddOSAP for our initial
	// command and getting back a secret and a response.
	sharedSecret, osapr, err := newOSAPSession(rw, etKeyHandle, handle, aikAuth)
	if err != nil {
		return nil, err
	}
	defer osapr.Close(rw)
	defer zeroBytes(sharedSecret[:])

	// Hash the data to get the value to pass to quote2.
	hash := sha1.Sum(data)
	pcrSel, err := newPCRSelection(pcrVals)
	if err != nil {
		return nil, err
	}
	authIn := []interface{}{ordQuote2, hash, pcrSel, addVersion}
	ca, err := newCommandAuth(osapr.AuthHandle, osapr.NonceEven, sharedSecret[:], authIn)
	if err != nil {
		return nil, err
	}

	// TODO(tmroeder): use the returned capVersionInfo.
	pcrShort, _, capBytes, sig, ra, ret, err := quote2(rw, handle, hash, pcrSel, addVersion, ca)
	if err != nil {
		return nil, err
	}

	// Check response authentication.
	raIn := []interface{}{ret, ordQuote2, pcrShort, capBytes, sig}
	if err := ra.verify(ca.NonceOdd, sharedSecret[:], raIn); err != nil {
		return nil, err
	}

	return sig, nil
}

// GetPubKey retrieves an opaque blob containing a public key corresponding to
// a handle from the TPM.
func GetPubKey(rw io.ReadWriter, keyHandle Handle, srkAuth []byte) ([]byte, error) {
	// Run OSAP for the handle, reading a random OddOSAP for our initial
	// command and getting back a secret and a response.
	sharedSecret, osapr, err := newOSAPSession(rw, etKeyHandle, keyHandle, srkAuth)
	if err != nil {
		return nil, err
	}
	defer osapr.Close(rw)
	defer zeroBytes(sharedSecret[:])

	authIn := []interface{}{ordGetPubKey}
	ca, err := newCommandAuth(osapr.AuthHandle, osapr.NonceEven, sharedSecret[:], authIn)
	if err != nil {
		return nil, err
	}

	pk, ra, ret, err := getPubKey(rw, keyHandle, ca)
	if err != nil {
		return nil, err
	}

	// Check response authentication for TPM_GetPubKey.
	raIn := []interface{}{ret, ordGetPubKey, pk}
	if err := ra.verify(ca.NonceOdd, sharedSecret[:], raIn); err != nil {
		return nil, err
	}

	b, err := pack([]interface{}{*pk})
	if err != nil {
		return nil, err
	}
	return b, err
}

// newOSAPSession starts a new OSAP session and derives a shared key from it.
func newOSAPSession(rw io.ReadWriter, entityType uint16, entityValue Handle, srkAuth []byte) ([20]byte, *osapResponse, error) {
	osapc := &osapCommand{
		EntityType:  entityType,
		EntityValue: entityValue,
	}

	var sharedSecret [20]byte
	if _, err := rand.Read(osapc.OddOSAP[:]); err != nil {
		return sharedSecret, nil, err
	}

	osapr, err := osap(rw, osapc)
	if err != nil {
		return sharedSecret, nil, err
	}

	// A shared secret is computed as
	//
	// sharedSecret = HMAC-SHA1(srkAuth, evenosap||oddosap)
	//
	// where srkAuth is the hash of the SRK authentication (which hash is all 0s
	// for the well-known SRK auth value) and even and odd OSAP are the
	// values from the OSAP protocol.
	osapData, err := pack([]interface{}{osapr.EvenOSAP, osapc.OddOSAP})
	if err != nil {
		return sharedSecret, nil, err
	}

	hm := hmac.New(sha1.New, srkAuth)
	hm.Write(osapData)
	// Note that crypto/hash.Sum returns a slice rather than an array, so we
	// have to copy this into an array to make sure that serialization doesn't
	// preprend a length in pack().
	sharedSecretBytes := hm.Sum(nil)
	copy(sharedSecret[:], sharedSecretBytes)
	return sharedSecret, osapr, nil
}

// newCommandAuth creates a new commandAuth structure over the given
// parameters, using the given secret for HMAC computation.
func newCommandAuth(authHandle Handle, nonceEven nonce, key []byte, params []interface{}) (*commandAuth, error) {
	// Auth = HMAC-SHA1(key, SHA1(params) || NonceEven || NonceOdd || ContSession)
	digestBytes, err := pack(params)
	if err != nil {
		return nil, err
	}

	digest := sha1.Sum(digestBytes)
	ca := &commandAuth{AuthHandle: authHandle}
	if _, err := rand.Read(ca.NonceOdd[:]); err != nil {
		return nil, err
	}

	authBytes, err := pack([]interface{}{digest, nonceEven, ca.NonceOdd, ca.ContSession})
	if err != nil {
		return nil, err
	}

	hm2 := hmac.New(sha1.New, key)
	hm2.Write(authBytes)
	auth := hm2.Sum(nil)
	copy(ca.Auth[:], auth[:])
	return ca, nil
}

// verify checks that the response authentication was correct.
// It computes the SHA1 of params, and computes the HMAC-SHA1 of this digest
// with the authentication parameters of ra along with the given odd nonce.
func (ra *responseAuth) verify(nonceOdd nonce, key []byte, params []interface{}) error {
	// Auth = HMAC-SHA1(key, SHA1(params) || ra.NonceEven || NonceOdd || ra.ContSession)
	digestBytes, err := pack(params)
	if err != nil {
		return err
	}

	digest := sha1.Sum(digestBytes)
	authBytes, err := pack([]interface{}{digest, ra.NonceEven, nonceOdd, ra.ContSession})
	if err != nil {
		return err
	}

	hm2 := hmac.New(sha1.New, key)
	hm2.Write(authBytes)
	auth := hm2.Sum(nil)

	if !hmac.Equal(ra.Auth[:], auth) {
		return errors.New("the computed response HMAC didn't match the provided HMAC")
	}

	return nil
}

// zeroBytes zeroes a byte array.
func zeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// Seal encrypts data against a given locality and PCRs and returns the sealed data.
func Seal(rw io.ReadWriter, locality byte, pcrs []int, data []byte, srkAuth []byte) ([]byte, error) {
	pcrInfo, err := newPCRInfoLong(rw, locality, pcrs)
	if err != nil {
		return nil, err
	}

	// Run OSAP for the SRK, reading a random OddOSAP for our initial
	// command and getting back a secret and a handle.
	sharedSecret, osapr, err := newOSAPSession(rw, etSRK, khSRK, srkAuth)
	if err != nil {
		return nil, err
	}
	defer osapr.Close(rw)
	defer zeroBytes(sharedSecret[:])

	// EncAuth for a seal command is computed as
	//
	// encAuth = XOR(srkAuth, SHA1(sharedSecret || <lastEvenNonce>))
	//
	// In this case, the last even nonce is NonceEven from OSAP.
	xorData, err := pack([]interface{}{sharedSecret, osapr.NonceEven})
	if err != nil {
		return nil, err
	}
	defer zeroBytes(xorData)

	encAuthData := sha1.Sum(xorData)
	sc := &sealCommand{KeyHandle: khSRK}
	for i := range sc.EncAuth {
		sc.EncAuth[i] = srkAuth[i] ^ encAuthData[i]
	}

	// The digest input for seal authentication is
	//
	// digest = SHA1(ordSeal || encAuth || binary.Size(pcrInfo) || pcrInfo ||
	//               len(data) || data)
	//
	authIn := []interface{}{ordSeal, sc.EncAuth, uint32(binary.Size(pcrInfo)), pcrInfo, data}
	ca, err := newCommandAuth(osapr.AuthHandle, osapr.NonceEven, sharedSecret[:], authIn)
	if err != nil {
		return nil, err
	}

	sealed, ra, ret, err := seal(rw, sc, pcrInfo, data, ca)
	if err != nil {
		return nil, err
	}

	// Check the response authentication.
	raIn := []interface{}{ret, ordSeal, sealed}
	if err := ra.verify(ca.NonceOdd, sharedSecret[:], raIn); err != nil {
		return nil, err
	}

	sealedBytes, err := pack([]interface{}{*sealed})
	if err != nil {
		return nil, err
	}

	return sealedBytes, nil
}

// Unseal decrypts data encrypted by the TPM.
func Unseal(rw io.ReadWriter, sealed []byte, srkAuth []byte) ([]byte, error) {
	// Run OSAP for the SRK, reading a random OddOSAP for our initial
	// command and getting back a secret and a handle.
	sharedSecret, osapr, err := newOSAPSession(rw, etSRK, khSRK, srkAuth)
	if err != nil {
		return nil, err
	}
	defer osapr.Close(rw)
	defer zeroBytes(sharedSecret[:])

	// The unseal command needs an OIAP session in addition to the OSAP session.
	oiapr, err := oiap(rw)
	if err != nil {
		return nil, err
	}
	defer oiapr.Close(rw)

	// Convert the sealed value into a tpmStoredData.
	var tsd tpmStoredData
	if err := unpack(sealed, []interface{}{&tsd}); err != nil {
		return nil, errors.New("couldn't convert the sealed data into a tpmStoredData struct")
	}

	// The digest for auth1 and auth2 for the unseal command is computed as
	// digest = SHA1(ordUnseal || tsd)
	authIn := []interface{}{ordUnseal, tsd}

	// The first commandAuth uses the shared secret as an HMAC key.
	ca1, err := newCommandAuth(osapr.AuthHandle, osapr.NonceEven, sharedSecret[:], authIn)
	if err != nil {
		return nil, err
	}

	// The second commandAuth is based on OIAP instead of OSAP and uses the
	// SRK auth value as an HMAC key instead of the shared secret.
	ca2, err := newCommandAuth(oiapr.AuthHandle, oiapr.NonceEven, srkAuth, authIn)
	if err != nil {
		return nil, err
	}

	unsealed, ra1, ra2, ret, err := unseal(rw, khSRK, &tsd, ca1, ca2)
	if err != nil {
		return nil, err
	}

	// Check the response authentication.
	raIn := []interface{}{ret, ordUnseal, unsealed}
	if err := ra1.verify(ca1.NonceOdd, sharedSecret[:], raIn); err != nil {
		return nil, err
	}

	if err := ra2.verify(ca2.NonceOdd, srkAuth, raIn); err != nil {
		return nil, err
	}

	return unsealed, nil
}

// Quote produces a TPM quote for the given data under the given PCRs. It uses
// AIK auth and a given AIK handle.
func Quote(rw io.ReadWriter, handle Handle, data []byte, pcrNums []int, aikAuth []byte) ([]byte, []byte, error) {
	// Run OSAP for the handle, reading a random OddOSAP for our initial
	// command and getting back a secret and a response.
	sharedSecret, osapr, err := newOSAPSession(rw, etKeyHandle, handle, aikAuth)
	if err != nil {
		return nil, nil, err
	}
	defer osapr.Close(rw)
	defer zeroBytes(sharedSecret[:])

	// Hash the data to get the value to pass to quote2.
	hash := sha1.Sum(data)
	pcrSel, err := newPCRSelection(pcrNums)
	if err != nil {
		return nil, nil, err
	}
	authIn := []interface{}{ordQuote, hash, pcrSel}
	ca, err := newCommandAuth(osapr.AuthHandle, osapr.NonceEven, sharedSecret[:], authIn)
	if err != nil {
		return nil, nil, err
	}

	pcrc, sig, ra, ret, err := quote(rw, handle, hash, pcrSel, ca)
	if err != nil {
		return nil, nil, err
	}

	// Check response authentication.
	raIn := []interface{}{ret, ordQuote, pcrc, sig}
	if err := ra.verify(ca.NonceOdd, sharedSecret[:], raIn); err != nil {
		return nil, nil, err
	}

	return sig, pcrc.Values, nil
}

// MakeIdentity creates a new AIK with the given new auth value, and the given
// parameters for the privacy CA that will be used to attest to it.
// If both pk and label are nil, then the TPM_CHOSENID_HASH is set to all 0s as
// a special case. MakeIdentity returns a key blob for the newly-created key.
// The caller must be authorized to use the SRK, since the private part of the
// AIK is sealed against the SRK.
// TODO(tmroeder): currently, this code can only create 2048-bit RSA keys.
func MakeIdentity(rw io.ReadWriter, srkAuth []byte, ownerAuth []byte, aikAuth []byte, pk crypto.PublicKey, label []byte) ([]byte, error) {
	// Run OSAP for the SRK, reading a random OddOSAP for our initial command
	// and getting back a secret and a handle.
	sharedSecretSRK, osaprSRK, err := newOSAPSession(rw, etSRK, khSRK, srkAuth)
	if err != nil {
		return nil, err
	}
	defer osaprSRK.Close(rw)
	defer zeroBytes(sharedSecretSRK[:])

	// Run OSAP for the Owner, reading a random OddOSAP for our initial command
	// and getting back a secret and a handle.
	sharedSecretOwn, osaprOwn, err := newOSAPSession(rw, etOwner, khOwner, ownerAuth)
	if err != nil {
		return nil, err
	}
	defer osaprOwn.Close(rw)
	defer zeroBytes(sharedSecretOwn[:])

	// EncAuth for a MakeIdentity command is computed as
	//
	// encAuth = XOR(aikAuth, SHA1(sharedSecretOwn || <lastEvenNonce>))
	//
	// In this case, the last even nonce is NonceEven from OSAP for the Owner.
	xorData, err := pack([]interface{}{sharedSecretOwn, osaprOwn.NonceEven})
	if err != nil {
		return nil, err
	}
	defer zeroBytes(xorData)

	encAuthData := sha1.Sum(xorData)
	var encAuth digest
	for i := range encAuth {
		encAuth[i] = aikAuth[i] ^ encAuthData[i]
	}

	var caDigest digest
	if (pk != nil) != (label != nil) {
		return nil, errors.New("inconsistent null values between the pk and the label")
	}

	if pk != nil {
		pubk, err := convertPubKey(pk)
		if err != nil {
			return nil, err
		}

		// We can't pack the pair of values directly, since the label is
		// included directly as bytes, without any length.
		fullpkb, err := pack([]interface{}{pubk})
		if err != nil {
			return nil, err
		}

		caDigestBytes := append(label, fullpkb...)
		caDigest = sha1.Sum(caDigestBytes)
	}

	rsaAIKParms := rsaKeyParms{
		KeyLength: 2048,
		NumPrimes: 2,
		//Exponent:  big.NewInt(0x10001).Bytes(), // 65537. Implicit?
	}
	packedParms, err := pack([]interface{}{rsaAIKParms})
	if err != nil {
		return nil, err
	}

	aikParms := keyParms{
		AlgID:     algRSA,
		EncScheme: esNone,
		SigScheme: ssRSASaPKCS1v15SHA1,
		Parms:     packedParms,
	}

	aik := &key{
		Version:        0x01010000,
		KeyUsage:       keyIdentity,
		KeyFlags:       0,
		AuthDataUsage:  authAlways,
		AlgorithmParms: aikParms,
	}

	// The digest input for MakeIdentity authentication is
	//
	// digest = SHA1(ordMakeIdentity || encAuth || caDigest || aik)
	//
	authIn := []interface{}{ordMakeIdentity, encAuth, caDigest, aik}
	ca1, err := newCommandAuth(osaprSRK.AuthHandle, osaprSRK.NonceEven, sharedSecretSRK[:], authIn)
	if err != nil {
		return nil, err
	}

	ca2, err := newCommandAuth(osaprOwn.AuthHandle, osaprOwn.NonceEven, sharedSecretOwn[:], authIn)
	if err != nil {
		return nil, err
	}

	k, sig, ra1, ra2, ret, err := makeIdentity(rw, encAuth, caDigest, aik, ca1, ca2)
	if err != nil {
		return nil, err
	}

	// Check response authentication.
	raIn := []interface{}{ret, ordMakeIdentity, k, sig}
	if err := ra1.verify(ca1.NonceOdd, sharedSecretSRK[:], raIn); err != nil {
		return nil, err
	}

	if err := ra2.verify(ca2.NonceOdd, sharedSecretOwn[:], raIn); err != nil {
		return nil, err
	}

	// TODO(tmroeder): check the signature against the pubek.
	blob, err := pack([]interface{}{k})
	if err != nil {
		return nil, err
	}

	return blob, nil
}

// ResetLockValue resets the dictionary-attack value in the TPM; this allows the
// TPM to start working again after authentication errors without waiting for
// the dictionary-attack defenses to time out. This requires owner
// authentication.
func ResetLockValue(rw io.ReadWriter, ownerAuth digest) error {
	// Run OSAP for the Owner, reading a random OddOSAP for our initial command
	// and getting back a secret and a handle.
	sharedSecretOwn, osaprOwn, err := newOSAPSession(rw, etOwner, khOwner, ownerAuth[:])
	if err != nil {
		return err
	}
	defer osaprOwn.Close(rw)
	defer zeroBytes(sharedSecretOwn[:])

	// The digest input for ResetLockValue auth is
	//
	// digest = SHA1(ordResetLockValue)
	//
	authIn := []interface{}{ordResetLockValue}
	ca, err := newCommandAuth(osaprOwn.AuthHandle, osaprOwn.NonceEven, sharedSecretOwn[:], authIn)
	if err != nil {
		return err
	}

	ra, ret, err := resetLockValue(rw, ca)
	if err != nil {
		return err
	}

	// Check response authentication.
	raIn := []interface{}{ret, ordResetLockValue}
	if err := ra.verify(ca.NonceOdd, sharedSecretOwn[:], raIn); err != nil {
		return err
	}

	return nil
}

// ownerReadInternalHelper sets up command auth and checks response auth for
// OwnerReadInternalPub. It's not exported because OwnerReadInternalPub only
// supports two fixed key handles: khEK and khSRK.
func ownerReadInternalHelper(rw io.ReadWriter, kh Handle, ownerAuth digest) (*pubKey, error) {
	// Run OSAP for the Owner, reading a random OddOSAP for our initial command
	// and getting back a secret and a handle.
	sharedSecretOwn, osaprOwn, err := newOSAPSession(rw, etOwner, khOwner, ownerAuth[:])
	if err != nil {
		return nil, err
	}
	defer osaprOwn.Close(rw)
	defer zeroBytes(sharedSecretOwn[:])

	// The digest input for OwnerReadInternalPub is
	//
	// digest = SHA1(ordOwnerReadInternalPub || kh)
	//
	authIn := []interface{}{ordOwnerReadInternalPub, kh}
	ca, err := newCommandAuth(osaprOwn.AuthHandle, osaprOwn.NonceEven, sharedSecretOwn[:], authIn)
	if err != nil {
		return nil, err
	}

	pk, ra, ret, err := ownerReadInternalPub(rw, kh, ca)
	if err != nil {
		return nil, err
	}

	// Check response authentication.
	raIn := []interface{}{ret, ordOwnerReadInternalPub, pk}
	if err := ra.verify(ca.NonceOdd, sharedSecretOwn[:], raIn); err != nil {
		return nil, err
	}

	return pk, nil
}

// OwnerReadSRK uses owner auth to get a blob representing the SRK.
func OwnerReadSRK(rw io.ReadWriter, ownerAuth digest) ([]byte, error) {
	pk, err := ownerReadInternalHelper(rw, khSRK, ownerAuth)
	if err != nil {
		return nil, err
	}

	return pack([]interface{}{pk})
}

// OwnerReadPubEK uses owner auth to get a blob representing the public part of the
// endorsement key.
func OwnerReadPubEK(rw io.ReadWriter, ownerAuth digest) ([]byte, error) {
	pk, err := ownerReadInternalHelper(rw, khEK, ownerAuth)
	if err != nil {
		return nil, err
	}

	return pack([]interface{}{pk})
}

// ReadPubEK reads the public part of the endorsement key when no owner is
// established.
func ReadPubEK(rw io.ReadWriter) ([]byte, error) {
	var n nonce
	if _, err := rand.Read(n[:]); err != nil {
		return nil, err
	}

	pk, d, _, err := readPubEK(rw, n)
	if err != nil {
		return nil, err
	}

	// Recompute the hash of the pk and the nonce to defend against replay
	// attacks.
	b, err := pack([]interface{}{pk, n})
	if err != nil {
		return nil, err
	}

	s := sha1.Sum(b)
	// There's no need for constant-time comparison of these hash values,
	// since no secret is involved.
	if !bytes.Equal(s[:], d[:]) {
		return nil, errors.New("the ReadPubEK operation failed the replay check")
	}

	return pack([]interface{}{pk})
}

// OwnerClear uses owner auth to clear the TPM. After this operation, the TPM
// can change ownership.
func OwnerClear(rw io.ReadWriter, ownerAuth digest) error {
	// Run OSAP for the Owner, reading a random OddOSAP for our initial command
	// and getting back a secret and a handle.
	sharedSecretOwn, osaprOwn, err := newOSAPSession(rw, etOwner, khOwner, ownerAuth[:])
	if err != nil {
		return err
	}
	defer osaprOwn.Close(rw)
	defer zeroBytes(sharedSecretOwn[:])

	// The digest input for OwnerClear is
	//
	// digest = SHA1(ordOwnerClear)
	//
	authIn := []interface{}{ordOwnerClear}
	ca, err := newCommandAuth(osaprOwn.AuthHandle, osaprOwn.NonceEven, sharedSecretOwn[:], authIn)
	if err != nil {
		return err
	}

	ra, ret, err := ownerClear(rw, ca)
	if err != nil {
		return err
	}

	// Check response authentication.
	raIn := []interface{}{ret, ordOwnerClear}
	if err := ra.verify(ca.NonceOdd, sharedSecretOwn[:], raIn); err != nil {
		return err
	}

	return nil
}

// TakeOwnership takes over a TPM and inserts a new owner auth value and
// generates a new SRK, associating it with a new SRK auth value. This
// operation can only be performed if there isn't already an owner for the TPM.
// The pub EK blob can be acquired by calling ReadPubEK if there is no owner, or
// OwnerReadPubEK if there is.
func TakeOwnership(rw io.ReadWriter, newOwnerAuth digest, newSRKAuth digest, pubEK []byte) error {

	// Encrypt the owner and SRK auth with the endorsement key.
	ek, err := UnmarshalPubRSAPublicKey(pubEK)
	if err != nil {
		return err
	}
	encOwnerAuth, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, ek, newOwnerAuth[:], oaepLabel)
	if err != nil {
		return err
	}
	encSRKAuth, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, ek, newSRKAuth[:], oaepLabel)
	if err != nil {
		return err
	}

	// The params for the SRK have very tight requirements:
	// - KeyLength must be 2048
	// - alg must be RSA
	// - Enc must be OAEP SHA1 MGF1
	// - Sig must be None
	// - Key usage must be Storage
	// - Key must not be migratable
	srkRSAParams := rsaKeyParms{
		KeyLength: 2048,
		NumPrimes: 2,
	}
	srkpb, err := pack([]interface{}{srkRSAParams})
	if err != nil {
		return err
	}
	srkParams := keyParms{
		AlgID:     algRSA,
		EncScheme: esRSAEsOAEPSHA1MGF1,
		SigScheme: ssNone,
		Parms:     srkpb,
	}
	srk := &key{
		Version:        0x01010000,
		KeyUsage:       keyStorage,
		KeyFlags:       0,
		AuthDataUsage:  authAlways,
		AlgorithmParms: srkParams,
	}

	// Get command auth using OIAP with the new owner auth.
	oiapr, err := oiap(rw)
	if err != nil {
		return err
	}
	defer oiapr.Close(rw)

	// The digest for TakeOwnership is
	//
	// SHA1(ordTakeOwnership || pidOwner || encOwnerAuth || encSRKAuth || srk)
	authIn := []interface{}{ordTakeOwnership, pidOwner, encOwnerAuth, encSRKAuth, srk}
	ca, err := newCommandAuth(oiapr.AuthHandle, oiapr.NonceEven, newOwnerAuth[:], authIn)
	if err != nil {
		return err
	}

	k, ra, ret, err := takeOwnership(rw, encOwnerAuth, encSRKAuth, srk, ca)
	if err != nil {
		return err
	}

	raIn := []interface{}{ret, ordTakeOwnership, k}
	return ra.verify(ca.NonceOdd, newOwnerAuth[:], raIn)
}
