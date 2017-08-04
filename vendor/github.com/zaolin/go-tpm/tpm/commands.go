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
	"encoding/binary"
	"errors"
	"io"
	"strconv"
)

// submitTPMRequest sends a structure to the TPM device file and gets results
// back, interpreting them as a new provided structure.
func submitTPMRequest(rw io.ReadWriter, tag uint16, ord uint32, in []interface{}, out []interface{}) (uint32, error) {
	if rw == nil {
		return 0, errors.New("nil TPM handle")
	}

	ch := commandHeader{tag, 0, ord}
	inb, err := packWithHeader(ch, in)
	if err != nil {
		return 0, err
	}

	if _, err := rw.Write(inb); err != nil {
		return 0, err
	}

	// Try to read the whole thing, but handle the case where it's just a
	// ResponseHeader and not the body, since that's what happens in the error
	// case.
	var rh responseHeader
	rhSize := binary.Size(rh)
	outb := make([]byte, maxTPMResponse)
	outlen, err := rw.Read(outb)
	if err != nil {
		return 0, err
	}

	// Resize the buffer to match the amount read from the TPM.
	outb = outb[:outlen]
	if err := unpack(outb[:rhSize], []interface{}{&rh}); err != nil {
		return 0, err
	}

	// Check success before trying to read the rest of the result.
	// Note that the command tag and its associated response tag differ by 3,
	// e.g., tagRQUCommand == 0x00C1, and tagRSPCommand == 0x00C4.
	if rh.Res != 0 {
		return rh.Res, tpmError(rh.Res)
	}

	if rh.Tag != ch.Tag+3 {
		return 0, errors.New("inconsistent tag returned by TPM. Expected " + strconv.Itoa(int(ch.Tag+3)) + " but got " + strconv.Itoa(int(rh.Tag)))
	}

	if rh.Size > uint32(rhSize) {
		if err := unpack(outb[rhSize:], out); err != nil {
			return 0, err
		}
	}

	return rh.Res, nil
}

// oiap sends an OIAP command to the TPM and gets back an auth value and a
// nonce.
func oiap(rw io.ReadWriter) (*oiapResponse, error) {
	var resp oiapResponse
	out := []interface{}{&resp}
	// In this case, we don't need to check ret, since all the information is
	// contained in err.
	if _, err := submitTPMRequest(rw, tagRQUCommand, ordOIAP, nil, out); err != nil {
		return nil, err
	}

	return &resp, nil
}

// osap sends an OSAPCommand to the TPM and gets back authentication
// information in an OSAPResponse.
func osap(rw io.ReadWriter, osap *osapCommand) (*osapResponse, error) {
	in := []interface{}{osap}
	var resp osapResponse
	out := []interface{}{&resp}
	// In this case, we don't need to check the ret value, since all the
	// information is contained in err.
	if _, err := submitTPMRequest(rw, tagRQUCommand, ordOSAP, in, out); err != nil {
		return nil, err
	}

	return &resp, nil
}

// seal performs a seal operation on the TPM.
func seal(rw io.ReadWriter, sc *sealCommand, pcrs *pcrInfoLong, data []byte, ca *commandAuth) (*tpmStoredData, *responseAuth, uint32, error) {
	pcrsize := binary.Size(pcrs)
	if pcrsize < 0 {
		return nil, nil, 0, errors.New("couldn't compute the size of a pcrInfoLong")
	}

	// TODO(tmroeder): special-case pcrInfoLong in pack/unpack so we don't have
	// to write out the length explicitly here.
	in := []interface{}{sc, uint32(pcrsize), pcrs, data, ca}

	var tsd tpmStoredData
	var ra responseAuth
	out := []interface{}{&tsd, &ra}
	ret, err := submitTPMRequest(rw, tagRQUAuth1Command, ordSeal, in, out)
	if err != nil {
		return nil, nil, 0, err
	}

	return &tsd, &ra, ret, nil
}

// unseal data sealed by the TPM.
func unseal(rw io.ReadWriter, keyHandle Handle, sealed *tpmStoredData, ca1 *commandAuth, ca2 *commandAuth) ([]byte, *responseAuth, *responseAuth, uint32, error) {
	in := []interface{}{keyHandle, sealed, ca1, ca2}
	var outb []byte
	var ra1 responseAuth
	var ra2 responseAuth
	out := []interface{}{&outb, &ra1, &ra2}
	ret, err := submitTPMRequest(rw, tagRQUAuth2Command, ordUnseal, in, out)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	return outb, &ra1, &ra2, ret, nil
}

// flushSpecific removes a handle from the TPM. Note that removing a handle
// doesn't require any authentication.
func flushSpecific(rw io.ReadWriter, handle Handle, resourceType uint32) error {
	// In this case, all the information is in err, so we don't check the
	// specific return-value details.
	_, err := submitTPMRequest(rw, tagRQUCommand, ordFlushSpecific, []interface{}{handle, resourceType}, nil)
	return err
}

// loadKey2 loads a key into the TPM. It's a tagRQUAuth1Command, so it only
// needs one auth parameter.
// TODO(tmroeder): support key12, too.
func loadKey2(rw io.ReadWriter, k *key, ca *commandAuth) (Handle, *responseAuth, uint32, error) {
	// We always load our keys with the SRK as the parent key.
	in := []interface{}{khSRK, k, ca}
	var keyHandle Handle
	var ra responseAuth
	out := []interface{}{&keyHandle, &ra}
	ret, err := submitTPMRequest(rw, tagRQUAuth1Command, ordLoadKey2, in, out)
	if err != nil {
		return 0, nil, 0, err
	}

	return keyHandle, &ra, ret, nil
}

// getPubKey gets a public key from the TPM
func getPubKey(rw io.ReadWriter, keyHandle Handle, ca *commandAuth) (*pubKey, *responseAuth, uint32, error) {
	in := []interface{}{keyHandle, ca}
	var pk pubKey
	var ra responseAuth
	out := []interface{}{&pk, &ra}
	ret, err := submitTPMRequest(rw, tagRQUAuth1Command, ordGetPubKey, in, out)
	if err != nil {
		return nil, nil, 0, err
	}

	return &pk, &ra, ret, nil
}

// quote2 signs arbitrary data under a given set of PCRs and using a key
// specified by keyHandle. It returns information about the PCRs it signed
// under, the signature, auth information, and optionally information about the
// TPM itself. Note that the input to quote2 must be exactly 20 bytes, so it is
// normally the SHA1 hash of the data.
func quote2(rw io.ReadWriter, keyHandle Handle, hash [20]byte, pcrs *pcrSelection, addVersion byte, ca *commandAuth) (*pcrInfoShort, *capVersionInfo, []byte, []byte, *responseAuth, uint32, error) {
	in := []interface{}{keyHandle, hash, pcrs, addVersion, ca}
	var pcrShort pcrInfoShort
	var capInfo capVersionInfo
	var capBytes []byte
	var sig []byte
	var ra responseAuth
	out := []interface{}{&pcrShort, &capBytes, &sig, &ra}
	ret, err := submitTPMRequest(rw, tagRQUAuth1Command, ordQuote2, in, out)
	if err != nil {
		return nil, nil, nil, nil, nil, 0, err
	}

	// Deserialize the capInfo, if any.
	if len(capBytes) == 0 {
		return &pcrShort, nil, capBytes, sig, &ra, ret, nil
	}

	size := binary.Size(capInfo.CapVersionFixed)
	capInfo.VendorSpecific = make([]byte, len(capBytes)-size)
	if err := unpack(capBytes[:size], []interface{}{&capInfo.CapVersionFixed}); err != nil {
		return nil, nil, nil, nil, nil, 0, err
	}

	copy(capInfo.VendorSpecific, capBytes[size:])

	return &pcrShort, &capInfo, capBytes, sig, &ra, ret, nil
}

// quote performs a TPM 1.1 quote operation: it signs data using the
// TPM_QUOTE_INFO structure for the current values of a selectied set of PCRs.
func quote(rw io.ReadWriter, keyHandle Handle, hash [20]byte, pcrs *pcrSelection, ca *commandAuth) (*pcrComposite, []byte, *responseAuth, uint32, error) {
	in := []interface{}{keyHandle, hash, pcrs, ca}
	var pcrc pcrComposite
	var sig []byte
	var ra responseAuth
	out := []interface{}{&pcrc, &sig, &ra}
	ret, err := submitTPMRequest(rw, tagRQUAuth1Command, ordQuote, in, out)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	return &pcrc, sig, &ra, ret, nil
}

// makeIdentity requests that the TPM create a new AIK. It returns the handle to
// this new key.
func makeIdentity(rw io.ReadWriter, encAuth digest, idDigest digest, k *key, ca1 *commandAuth, ca2 *commandAuth) (*key, []byte, *responseAuth, *responseAuth, uint32, error) {
	in := []interface{}{encAuth, idDigest, k, ca1, ca2}
	var aik key
	var sig []byte
	var ra1 responseAuth
	var ra2 responseAuth
	out := []interface{}{&aik, &sig, &ra1, &ra2}
	ret, err := submitTPMRequest(rw, tagRQUAuth2Command, ordMakeIdentity, in, out)
	if err != nil {
		return nil, nil, nil, nil, 0, err
	}

	return &aik, sig, &ra1, &ra2, ret, nil
}

// resetLockValue resets the dictionary-attack lock in the TPM, using owner
// auth.
func resetLockValue(rw io.ReadWriter, ca *commandAuth) (*responseAuth, uint32, error) {
	in := []interface{}{ca}
	var ra responseAuth
	out := []interface{}{&ra}
	ret, err := submitTPMRequest(rw, tagRQUAuth1Command, ordResetLockValue, in, out)
	if err != nil {
		return nil, 0, err
	}

	return &ra, ret, nil
}

// ownerReadInternalPub uses owner auth and OSAP to read either the endorsement
// key (using khEK) or the SRK (using khSRK).
func ownerReadInternalPub(rw io.ReadWriter, kh Handle, ca *commandAuth) (*pubKey, *responseAuth, uint32, error) {
	in := []interface{}{kh, ca}
	var pk pubKey
	var ra responseAuth
	out := []interface{}{&pk, &ra}
	ret, err := submitTPMRequest(rw, tagRQUAuth1Command, ordOwnerReadInternalPub, in, out)
	if err != nil {
		return nil, nil, 0, err
	}

	return &pk, &ra, ret, nil
}

// readPubEK requests the public part of the endorsement key from the TPM. Note
// that this call can only be made when there is no owner in the TPM. Once an
// owner is established, the endorsement key can be retrieved using
// ownerReadInternalPub.
func readPubEK(rw io.ReadWriter, n nonce) (*pubKey, digest, uint32, error) {
	in := []interface{}{n}
	var pk pubKey
	var d digest
	out := []interface{}{&pk, &d}
	ret, err := submitTPMRequest(rw, tagRQUCommand, ordReadPubEK, in, out)
	if err != nil {
		return nil, d, 0, err
	}

	return &pk, d, ret, nil
}

// ownerClear uses owner auth to clear the TPM. After this operation, a caller
// can take ownership of the TPM with TPM_TakeOwnership.
func ownerClear(rw io.ReadWriter, ca *commandAuth) (*responseAuth, uint32, error) {
	in := []interface{}{ca}
	var ra responseAuth
	out := []interface{}{&ra}
	ret, err := submitTPMRequest(rw, tagRQUAuth1Command, ordOwnerClear, in, out)
	if err != nil {
		return nil, 0, err
	}

	return &ra, ret, nil
}

// takeOwnership takes ownership of the TPM and establishes a new SRK and
// owner auth. This operation can only be performed if there is no owner. The
// TPM can be put into this state using TPM_OwnerClear. The encOwnerAuth and
// encSRKAuth values must be encrypted using the endorsement key.
func takeOwnership(rw io.ReadWriter, encOwnerAuth []byte, encSRKAuth []byte, srk *key, ca *commandAuth) (*key, *responseAuth, uint32, error) {
	in := []interface{}{pidOwner, encOwnerAuth, encSRKAuth, srk, ca}
	var k key
	var ra responseAuth
	out := []interface{}{&k, &ra}
	ret, err := submitTPMRequest(rw, tagRQUAuth1Command, ordTakeOwnership, in, out)
	if err != nil {
		return nil, nil, 0, err
	}

	return &k, &ra, ret, nil
}
