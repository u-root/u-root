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
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"flag"
	"hash"
	"io"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-tpm-tools/simulator"
	. "github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

var (
	runClear = flag.Bool("run-clear", false, "Set to run tests which will clear hierarchy and lockout authorizations")
)

func init() {
	testing.Init()
	flag.Parse()
}

func openTPM(tb testing.TB) io.ReadWriteCloser {
	tb.Helper()
	if useDeviceTPM() {
		return openDeviceTPM(tb)
	}
	simulator, err := simulator.Get()
	if err != nil {
		tb.Fatalf("Simulator initialization failed: %v", err)
	}
	return simulator
}

var (
	// PCR7 is for SecureBoot.
	pcrSelection0    = PCRSelection{Hash: AlgSHA1, PCRs: []int{0}}
	pcrSelection1    = PCRSelection{Hash: AlgSHA1, PCRs: []int{1}}
	pcrSelection7    = PCRSelection{Hash: AlgSHA1, PCRs: []int{7}}
	pcrSelectionAll  = PCRSelection{Hash: AlgSHA1, PCRs: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}}
	defaultKeyParams = Public{
		Type:       AlgRSA,
		NameAlg:    AlgSHA1,
		Attributes: FlagStorageDefault,
		RSAParameters: &RSAParams{
			Symmetric: &SymScheme{
				Alg:     AlgAES,
				KeyBits: 128,
				Mode:    AlgCFB,
			},
			KeyBits:     2048,
			ExponentRaw: 1<<16 + 1,
		},
	}
	defaultPassword = "\x01\x02\x03\x04"
	emptyPassword   = ""
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestGetRandom(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	if _, err := GetRandom(rw, 16); err != nil {
		t.Fatalf("GetRandom failed: %v", err)
	}
}

func TestReadPCRs(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	pcrs, err := ReadPCRs(rw, pcrSelection7)
	if err != nil {
		t.Errorf("ReadPCRs failed: %s", err)
	}
	if !useDeviceTPM() {
		return // PCR 7 is initialized to 0 in the simulator.
	}
	for pcr, val := range pcrs {
		if empty := make([]byte, len(val)); reflect.DeepEqual(empty, val) {
			t.Errorf("Value of PCR %d is empty", pcr)
		}
	}
}

func TestReadAllPCRs(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	numPCRs := len(pcrSelectionAll.PCRs)
	out := map[uint32][]byte{}

	for i := 0; i < numPCRs; i += 8 {
		// Build a selection structure, specifying 8 PCRs at a time
		end := min(i+8, numPCRs)
		pcrSel := PCRSelection{
			Hash: pcrSelectionAll.Hash,
			PCRs: pcrSelectionAll.PCRs[i:end],
		}

		// Ask the TPM for those PCR values.
		ret, err := ReadPCRs(rw, pcrSel)
		if err != nil {
			t.Errorf("ReadPCRs(%+v) failed: %v", pcrSel, err)
		}

		// Keep track of the PCRs we were actually given.
		for pcr, digest := range ret {
			out[uint32(pcr)] = digest
		}
	}

	if len(out) != numPCRs {
		t.Errorf("Failed to read all PCRs, only read %d", len(out))
	}
}

func TestReadClock(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	if _, _, err := ReadClock(rw); err != nil {
		t.Fatalf("ReadClock failed: %s", err)
	}

}

func TestGetCapability(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	for _, tt := range []struct {
		capa     Capability
		count    uint32
		property uint32
		typ      interface{}
	}{
		{CapabilityHandles, 1, uint32(HandleTypeTransient) << 24, tpmutil.Handle(0)},
		{CapabilityAlgs, 1, 0, AlgorithmDescription{}},
		{CapabilityTPMProperties, 1, uint32(NVMaxBufferSize), TaggedProperty{}},
	} {
		l, _, err := GetCapability(rw, tt.capa, tt.count, tt.property)
		if err != nil {
			t.Fatalf("GetCapability(%v, %d, %d) = _, %v; want _, nil", tt.capa, tt.count, tt.property, err)
		}
		for _, i := range l {
			if reflect.TypeOf(i) != reflect.TypeOf(tt.typ) {
				t.Fatalf("GetCapability(%v, %d, %d) returned an element with the wrong type: %v; want %v", tt.capa, tt.count, tt.property, reflect.TypeOf(i), reflect.TypeOf(tt.typ))
			}
		}
	}
}

func TestCombinedKeyTest(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	parentHandle, _, err := CreatePrimary(rw, HandleOwner, pcrSelection7, emptyPassword, defaultPassword, defaultKeyParams)
	if err != nil {
		t.Fatalf("CreatePrimary failed: %s", err)
	}
	defer FlushContext(rw, parentHandle)

	privateBlob, publicBlob, _, _, _, err := CreateKey(rw, parentHandle, pcrSelection7, defaultPassword, defaultPassword, defaultKeyParams)
	if err != nil {
		t.Fatalf("CreateKey failed: %s", err)
	}

	keyHandle, _, err := Load(rw, parentHandle, defaultPassword, publicBlob, privateBlob)
	if err != nil {
		t.Fatalf("Load failed: %s", err)
	}
	defer FlushContext(rw, keyHandle)

	if _, _, _, err := ReadPublic(rw, keyHandle); err != nil {
		t.Fatalf("ReadPublic failed: %s", err)
	}
}

func TestCombinedEndorsementTest(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	parentHandle, _, err := CreatePrimary(rw, HandleOwner, pcrSelection7, emptyPassword, emptyPassword, defaultKeyParams)
	if err != nil {
		t.Fatalf("CreatePrimary failed: %s", err)
	}
	defer FlushContext(rw, parentHandle)

	privateBlob, publicBlob, _, _, _, err := CreateKey(rw, parentHandle, pcrSelection7, emptyPassword, defaultPassword, defaultKeyParams)
	if err != nil {
		t.Fatalf("CreateKey failed: %s", err)
	}

	keyHandle, _, err := Load(rw, parentHandle, emptyPassword, publicBlob, privateBlob)
	if err != nil {
		t.Fatalf("Load failed: %s", err)
	}
	defer FlushContext(rw, keyHandle)

	_, name, _, err := ReadPublic(rw, keyHandle)
	if err != nil {
		t.Fatalf("ReadPublic failed: %s", err)
	}

	// Generate Credential
	credential := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf, 0x10}
	credBlob, encryptedSecret0, err := MakeCredential(rw, parentHandle, credential, name)
	if err != nil {
		t.Fatalf("MakeCredential failed: %v", err)
	}

	recoveredCredential1, err := ActivateCredential(rw, keyHandle, parentHandle, defaultPassword, emptyPassword, credBlob, encryptedSecret0)
	if err != nil {
		t.Fatalf("ActivateCredential failed: %v", err)
	}
	if !bytes.Equal(credential, recoveredCredential1) {
		t.Fatalf("Credential and recovered credential differ: got %v, want %v", recoveredCredential1, credential)
	}

	recoveredCredential2, err := ActivateCredentialUsingAuth(rw, []AuthCommand{
		{Session: HandlePasswordSession, Attributes: AttrContinueSession, Auth: []byte(defaultPassword)},
		{Session: HandlePasswordSession, Attributes: AttrContinueSession, Auth: []byte(emptyPassword)},
	}, keyHandle, parentHandle, credBlob, encryptedSecret0)
	if err != nil {
		t.Fatalf("ActivateCredentialWithAuth failed: %v", err)
	}
	if !bytes.Equal(credential, recoveredCredential2) {
		t.Errorf("Credential and recovered credential differ: got %v, want %v", recoveredCredential2, credential)
	}

	_, err = ActivateCredentialUsingAuth(rw, []AuthCommand{
		{Session: HandlePasswordSession, Attributes: AttrContinueSession, Auth: []byte("incorrect password")},
		{Session: HandlePasswordSession, Attributes: AttrContinueSession, Auth: []byte(emptyPassword)},
	}, keyHandle, parentHandle, credBlob, encryptedSecret0)
	if err == nil {
		t.Fatal("ActivateCredentialUsingAuth: error == nil, expected authorization failure")
	}
	if !strings.Contains(err.Error(), "the authorization HMAC check failed") {
		t.Errorf("ActivateCredentialUsingAuth: error = %v, expected authorization failure", err)
	}

	_, err = ActivateCredentialUsingAuth(rw, []AuthCommand{
		{Session: HandlePasswordSession, Attributes: AttrContinueSession, Auth: []byte(emptyPassword)},
	}, keyHandle, parentHandle, credBlob, encryptedSecret0)
	if err == nil {
		t.Fatal("ActivateCredentialUsingAuth: error == nil, expected response status 0x98e (authorization failure)")
	}
	if !strings.Contains(err.Error(), "len(auth) = 1, want 2") {
		t.Errorf("ActivateCredentialUsingAuth: error = %v, expected len(auth) = 1, want 2", err)
	}
}

func TestCreatePrimaryEx(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	keyHandle, pub1, creation, _, _, name, err := CreatePrimaryEx(rw, HandleOwner, pcrSelection7, emptyPassword, emptyPassword, defaultKeyParams)
	if err != nil {
		t.Fatalf("CreatePrimary failed: %v", err)
	}
	defer FlushContext(rw, keyHandle)

	pub, _, _, err := ReadPublic(rw, keyHandle)
	if err != nil {
		t.Fatalf("ReadPublic failed: %s", err)
	}
	pub2, err := pub.Encode()
	if err != nil {
		t.Fatalf("Failed to encode public: %v", err)
	}

	if !bytes.Equal(pub1, pub2) {
		t.Error("Mismatch between public returned from CreatePrimaryEx() & ReadPublic()")
		t.Logf("CreatePrimaryEx: %v", pub1)
		t.Logf("ReadPublic:      %v", pub2)
	}

	if _, err := DecodeName(bytes.NewBuffer(name)); err != nil {
		t.Errorf("Failed to decode name: %v", err)
	}
	if _, err := DecodeCreationData(creation); err != nil {
		t.Fatalf("DecodeCreationData() returned err: %v", err)
	}
}

func TestCombinedContextTest(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	rootHandle, _, err := CreatePrimary(rw, HandleOwner, pcrSelection7, emptyPassword, emptyPassword, defaultKeyParams)
	if err != nil {
		t.Fatalf("CreatePrimary failed: %v", err)
	}
	defer FlushContext(rw, rootHandle)

	// CreateKey (Quote Key)
	quotePrivate, quotePublic, _, _, _, err := CreateKey(rw, rootHandle, pcrSelection7, emptyPassword, emptyPassword, defaultKeyParams)
	if err != nil {
		t.Fatalf("CreateKey failed: %v", err)
	}

	quoteHandle, _, err := Load(rw, rootHandle, emptyPassword, quotePublic, quotePrivate)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	defer FlushContext(rw, quoteHandle)

	saveArea, err := ContextSave(rw, quoteHandle)
	if err != nil {
		t.Fatalf("ContextSave failed: %v", err)
	}
	FlushContext(rw, quoteHandle)

	quoteHandle, err = ContextLoad(rw, saveArea)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
}

func TestEvictControl(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	rootHandle, _, err := CreatePrimary(rw, HandleOwner, pcrSelection7, emptyPassword, emptyPassword, defaultKeyParams)
	if err != nil {
		t.Fatalf("CreatePrimary failed: %v", err)
	}
	defer FlushContext(rw, rootHandle)

	// CreateKey (Quote Key)
	quotePrivate, quotePublic, _, _, _, err := CreateKey(rw, rootHandle, pcrSelection7, emptyPassword, emptyPassword, defaultKeyParams)
	if err != nil {
		t.Fatalf("CreateKey failed: %v", err)
	}

	quoteHandle, _, err := Load(rw, rootHandle, emptyPassword, quotePublic, quotePrivate)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	defer FlushContext(rw, quoteHandle)

	persistentHandle := tpmutil.Handle(0x817FFFFF)
	// Evict persistent key, if there is one already (e.g. last test run failed).
	if err := EvictControl(rw, emptyPassword, HandleOwner, persistentHandle, persistentHandle); err != nil {
		t.Logf("(expected) EvictControl failed: %v", err)
	}
	// Make key persistent.
	if err := EvictControl(rw, emptyPassword, HandleOwner, quoteHandle, persistentHandle); err != nil {
		t.Fatalf("EvictControl failed: %v", err)
	}
	// Evict persistent key.
	if err := EvictControl(rw, emptyPassword, HandleOwner, persistentHandle, persistentHandle); err != nil {
		t.Fatalf("EvictControl failed: %v", err)
	}
}

func TestHash(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	run := func(t *testing.T, data []byte, hierarchy tpmutil.Handle, wantValidation bool) {
		gotDigest, gotValidation, err := Hash(rw, AlgSHA256, data, hierarchy)
		if err != nil {
			t.Fatalf("Hash failed: %v", err)
		}
		wantDigest := sha256.Sum256(data)

		if !bytes.Equal(gotDigest, wantDigest[:]) {
			t.Errorf("Hash(%q) returned digest %x, want %x", data, gotDigest, wantDigest)
		}
		if wantValidation && len(gotValidation.Digest) == 0 {
			t.Errorf("Hash(%q) unexpectedly returned empty validation ticket", data)
		}
		if !wantValidation && len(gotValidation.Digest) != 0 {
			t.Errorf("Hash(%q) unexpectedly returned non-empty validation ticket", data)
		}
	}
	t.Run("Null hierarchy", func(t *testing.T) {
		run(t, []byte("foobarbaz"), HandleNull, false)
	})
	t.Run("Owner hierarchy", func(t *testing.T) {
		run(t, []byte("foobarbaz"), HandleOwner, true)
	})
	t.Run("Starts with TPM_GENERATED_VALUE", func(t *testing.T) {
		run(t, []byte("\xffTCGbarbaz"), HandleOwner, false)
	})
}

func skipOnUnsupportedAlg(t testing.TB, rw io.ReadWriter, alg Algorithm) {
	moreData := true
	for i := uint32(0); moreData; i++ {
		var err error
		var descs []interface{}
		descs, moreData, err = GetCapability(rw, CapabilityAlgs, 1, i)
		if err != nil {
			t.Fatalf("Could not get TPM algorithm capability: %v", err)
		}
		for _, desc := range descs {
			if desc.(AlgorithmDescription).ID == alg {
				return
			}
		}
		if !moreData {
			break
		}
	}
	t.Skipf("Algorithm %v is not supported by the TPM", alg)
}

func TestLoadExternalPublicKey(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	run := func(t *testing.T, public Public, private Private) {
		t.Helper()

		h, _, err := LoadExternal(rw, public, private, HandleNull)
		if err != nil {
			t.Fatal(err)
		}
		defer FlushContext(rw, h)
	}

	t.Run("RSA", func(t *testing.T) {
		pk, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}
		rp := Public{
			Type:       AlgRSA,
			NameAlg:    AlgSHA1,
			Attributes: FlagSign | FlagSensitiveDataOrigin | FlagUserWithAuth,
			RSAParameters: &RSAParams{
				Sign: &SigScheme{
					Alg:  AlgRSASSA,
					Hash: AlgSHA1,
				},
				KeyBits:     2048,
				ExponentRaw: uint32(pk.PublicKey.E),
				ModulusRaw:  pk.PublicKey.N.Bytes(),
			},
		}
		private := Private{
			Type:      AlgRSA,
			Sensitive: pk.Primes[0].Bytes(),
		}
		run(t, rp, private)
	})
	t.Run("ECC", func(t *testing.T) {
		skipOnUnsupportedAlg(t, rw, AlgECC)
		pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		public := Public{
			Type:       AlgECC,
			NameAlg:    AlgSHA1,
			Attributes: FlagSign | FlagSensitiveDataOrigin | FlagUserWithAuth,
			ECCParameters: &ECCParams{
				Sign: &SigScheme{
					Alg:  AlgECDSA,
					Hash: AlgSHA1,
				},
				CurveID: CurveNISTP256,
				Point:   ECPoint{XRaw: pk.PublicKey.X.Bytes(), YRaw: pk.PublicKey.Y.Bytes()},
			},
		}
		private := Private{
			Type:      AlgECC,
			Sensitive: pk.D.Bytes(),
		}
		run(t, public, private)
	})
}

func TestCertify(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	params := Public{
		Type:       AlgRSA,
		NameAlg:    AlgSHA256,
		Attributes: FlagSignerDefault,
		RSAParameters: &RSAParams{
			Sign: &SigScheme{
				Alg:  AlgRSASSA,
				Hash: AlgSHA256,
			},
			KeyBits: 2048,
		},
	}
	signerHandle, signerPub, err := CreatePrimary(rw, HandleOwner, pcrSelection7, emptyPassword, defaultPassword, params)
	if err != nil {
		t.Fatalf("CreatePrimary(signer) failed: %s", err)
	}
	defer FlushContext(rw, signerHandle)

	subjectHandle, subjectPub, err := CreatePrimary(rw, HandlePlatform, pcrSelection7, emptyPassword, defaultPassword, params)
	if err != nil {
		t.Fatalf("CreatePrimary(subject) failed: %s", err)
	}
	defer FlushContext(rw, subjectHandle)

	attest, sig, err := Certify(rw, defaultPassword, defaultPassword, subjectHandle, signerHandle, nil)
	if err != nil {
		t.Errorf("Certify failed: %s", err)
		return
	}

	attestHash := sha256.Sum256(attest)
	if err := rsa.VerifyPKCS1v15(signerPub.(*rsa.PublicKey), crypto.SHA256, attestHash[:], sig); err != nil {
		t.Errorf("Signature verification failed: %v", err)
	}

	t.Run("DecodeAttestationData", func(t *testing.T) {
		ad, err := DecodeAttestationData(attest)
		if err != nil {
			t.Fatal("DecodeAttestationData:", err)
		}
		params := Public{
			Type:       AlgRSA,
			NameAlg:    AlgSHA256,
			Attributes: FlagSignerDefault,
			RSAParameters: &RSAParams{
				Sign: &SigScheme{
					Alg:  AlgRSASSA,
					Hash: AlgSHA256,
				},
				KeyBits: 2048,
				// Note: we don't include Exponent because CreatePrimary also
				// returns Public without it.
				ModulusRaw: subjectPub.(*rsa.PublicKey).N.Bytes(),
			},
		}
		matches, err := ad.AttestedCertifyInfo.Name.MatchesPublic(params)
		if err != nil {
			t.Fatalf("AttestedCertifyInfo.Name.MatchesPublic error: %v", err)
		}
		if !matches {
			t.Error("Name in AttestationData doesn't match Public structure of subject")
		}
	})
}

func TestCertifyExternalKey(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	params := Public{
		Type:       AlgRSA,
		NameAlg:    AlgSHA256,
		Attributes: FlagSignerDefault,
		RSAParameters: &RSAParams{
			Sign: &SigScheme{
				Alg:  AlgRSASSA,
				Hash: AlgSHA256,
			},
			KeyBits: 2048,
		},
	}
	signerHandle, signerPub, err := CreatePrimary(rw, HandleOwner, pcrSelection7, emptyPassword, defaultPassword, params)
	if err != nil {
		t.Fatalf("CreatePrimary(signer) failed: %s", err)
	}
	defer FlushContext(rw, signerHandle)

	run := func(t *testing.T, public Public, private Private) {
		t.Helper()
		subjectHandle, _, err := LoadExternal(rw, public, private, HandleNull)
		if err != nil {
			t.Fatalf("LoadExternal: %v", err)
		}
		defer FlushContext(rw, subjectHandle)

		attest, sig, err := Certify(rw, emptyPassword, defaultPassword, subjectHandle, signerHandle, nil)
		if err != nil {
			t.Errorf("Certify failed: %s", err)
			return
		}

		attestHash := sha256.Sum256(attest)
		if err := rsa.VerifyPKCS1v15(signerPub.(*rsa.PublicKey), crypto.SHA256, attestHash[:], sig); err != nil {
			t.Errorf("Signature verification failed: %v", err)
		}
	}
	t.Run("RSA", func(t *testing.T) {
		pk, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}
		public := Public{
			Type:       AlgRSA,
			NameAlg:    AlgSHA1,
			Attributes: FlagSign | FlagSensitiveDataOrigin | FlagUserWithAuth,
			RSAParameters: &RSAParams{
				Sign: &SigScheme{
					Alg:  AlgRSASSA,
					Hash: AlgSHA1,
				},
				KeyBits:     2048,
				ExponentRaw: uint32(pk.PublicKey.E),
				ModulusRaw:  pk.PublicKey.N.Bytes(),
			},
		}
		private := Private{
			Type:      AlgRSA,
			Sensitive: pk.Primes[0].Bytes(),
		}
		run(t, public, private)
	})
	t.Run("ECC", func(t *testing.T) {
		skipOnUnsupportedAlg(t, rw, AlgECC)
		pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		public := Public{
			Type:       AlgECC,
			NameAlg:    AlgSHA1,
			Attributes: FlagSign | FlagSensitiveDataOrigin | FlagUserWithAuth,
			ECCParameters: &ECCParams{
				Sign: &SigScheme{
					Alg:  AlgECDSA,
					Hash: AlgSHA1,
				},
				CurveID: CurveNISTP256,
				Point:   ECPoint{XRaw: pk.PublicKey.X.Bytes(), YRaw: pk.PublicKey.Y.Bytes()},
			},
		}
		private := Private{
			Type:      AlgECC,
			Sensitive: pk.D.Bytes(),
		}
		run(t, public, private)
	})
}

func TestSign(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	run := func(t *testing.T, pub Public) {
		signerHandle, signerPub, err := CreatePrimary(rw, HandleOwner, pcrSelection7, emptyPassword, defaultPassword, pub)
		if err != nil {
			t.Fatalf("CreatePrimary failed: %s", err)
		}
		defer FlushContext(rw, signerHandle)

		digest := sha256.Sum256([]byte("heyo"))

		sig, err := Sign(rw, signerHandle, defaultPassword, digest[:], nil, nil)
		if err != nil {
			t.Fatalf("Sign failed: %s", err)
		}
		switch signerPub := signerPub.(type) {
		case *rsa.PublicKey:
			switch pub.RSAParameters.Sign.Alg {
			case AlgRSASSA:
				if err := rsa.VerifyPKCS1v15(signerPub, crypto.SHA256, digest[:], sig.RSA.Signature); err != nil {
					t.Errorf("Signature verification failed: %v", err)
				}
			case AlgRSAPSS:
				if err := rsa.VerifyPSS(signerPub, crypto.SHA256, digest[:], sig.RSA.Signature, nil); err != nil {
					t.Errorf("Signature verification failed: %v", err)
				}
			default:
				t.Errorf("unsupported signature algorithm 0x%x", pub.RSAParameters.Sign.Alg)
			}
		case *ecdsa.PublicKey:
			if !ecdsa.Verify(signerPub, digest[:], sig.ECC.R, sig.ECC.S) {
				t.Error("Signature verification failed")
			}
		}
	}

	t.Run("RSA SSA", func(t *testing.T) {
		run(t, Public{
			Type:       AlgRSA,
			NameAlg:    AlgSHA256,
			Attributes: FlagSign | FlagSensitiveDataOrigin | FlagUserWithAuth,
			RSAParameters: &RSAParams{
				Sign: &SigScheme{
					Alg:  AlgRSASSA,
					Hash: AlgSHA256,
				},
				KeyBits: 2048,
			},
		})
	})
	t.Run("RSA PSS", func(t *testing.T) {
		run(t, Public{
			Type:       AlgRSA,
			NameAlg:    AlgSHA256,
			Attributes: FlagSign | FlagSensitiveDataOrigin | FlagUserWithAuth,
			RSAParameters: &RSAParams{
				Sign: &SigScheme{
					Alg:  AlgRSAPSS,
					Hash: AlgSHA256,
				},
				KeyBits: 2048,
			},
		})
	})
	t.Run("ECC", func(t *testing.T) {
		skipOnUnsupportedAlg(t, rw, AlgECC)
		run(t, Public{
			Type:       AlgECC,
			NameAlg:    AlgSHA256,
			Attributes: FlagSign | FlagSensitiveDataOrigin | FlagUserWithAuth,
			ECCParameters: &ECCParams{
				Sign: &SigScheme{
					Alg:  AlgECDSA,
					Hash: AlgSHA256,
				},
				CurveID: CurveNISTP256,
			},
		})
	})
}

func TestSignWithAttestationKey(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	run := func(t *testing.T, data []byte, wantErr bool, pub Public) {
		signerHandle, signerPub, err := CreatePrimary(rw, HandleOwner, pcrSelection7, emptyPassword, defaultPassword, pub)
		if err != nil {
			t.Fatalf("CreatePrimary failed: %v", err)
		}
		defer FlushContext(rw, signerHandle)

		digest, validation, err := Hash(rw, AlgSHA256, data, HandleOwner)
		if err != nil {
			t.Fatalf("Hash failed unexpectedly: %v", err)
		}

		sig, err := Sign(rw, signerHandle, defaultPassword, digest, validation, nil)
		if err != nil && !wantErr {
			t.Fatalf("Sign failed unexpectedly: %v", err)
		}
		if err == nil && wantErr {
			t.Fatalf("Sign succeeded unexpectedly: %v", sig)
		}
		if !wantErr {
			switch signerPub := signerPub.(type) {
			case *rsa.PublicKey:
				switch pub.RSAParameters.Sign.Alg {
				case AlgRSASSA:
					if err := rsa.VerifyPKCS1v15(signerPub, crypto.SHA256, digest[:], sig.RSA.Signature); err != nil {
						t.Errorf("Signature verification failed: %v", err)
					}
				case AlgRSAPSS:
					if err := rsa.VerifyPSS(signerPub, crypto.SHA256, digest[:], sig.RSA.Signature, nil); err != nil {
						t.Errorf("Signature verification failed: %v", err)
					}
				default:
					t.Errorf("unsupported signature algorithm 0x%x", pub.RSAParameters.Sign.Alg)
				}
			case *ecdsa.PublicKey:
				if !ecdsa.Verify(signerPub, digest[:], sig.ECC.R, sig.ECC.S) {
					t.Error("Signature verification failed")
				}
			}
		}
	}

	rsassa := Public{
		Type:       AlgRSA,
		NameAlg:    AlgSHA256,
		Attributes: FlagSign | FlagRestricted | FlagSensitiveDataOrigin | FlagUserWithAuth,
		RSAParameters: &RSAParams{
			Sign: &SigScheme{
				Alg:  AlgRSASSA,
				Hash: AlgSHA256,
			},
			KeyBits: 2048,
		},
	}
	rsapss := Public{
		Type:       AlgRSA,
		NameAlg:    AlgSHA256,
		Attributes: FlagSign | FlagRestricted | FlagSensitiveDataOrigin | FlagUserWithAuth,
		RSAParameters: &RSAParams{
			Sign: &SigScheme{
				Alg:  AlgRSAPSS,
				Hash: AlgSHA256,
			},
			KeyBits: 2048,
		},
	}
	ecdsa := Public{
		Type:       AlgECC,
		NameAlg:    AlgSHA256,
		Attributes: FlagSign | FlagRestricted | FlagSensitiveDataOrigin | FlagUserWithAuth,
		ECCParameters: &ECCParams{
			Sign: &SigScheme{
				Alg:  AlgECDSA,
				Hash: AlgSHA256,
			},
			CurveID: CurveNISTP256,
		},
	}

	t.Run("RSA SSA", func(t *testing.T) {
		run(t, []byte("test123"), false, rsassa)
	})
	t.Run("RSA PSS", func(t *testing.T) {
		run(t, []byte("test123"), false, rsapss)
	})
	t.Run("ECC", func(t *testing.T) {
		skipOnUnsupportedAlg(t, rw, AlgECC)
		run(t, []byte("test123"), false, ecdsa)
	})
	t.Run("RSA SSA with TPM_GENERATED_VALUE", func(t *testing.T) {
		run(t, []byte("\xffTCG123"), true, rsassa)
	})
	t.Run("RSA PSS with TPM_GENERATED_VALUE", func(t *testing.T) {
		run(t, []byte("\xffTCG123"), true, rsapss)
	})
	t.Run("ECC with TPM_GENERATED_VALUE", func(t *testing.T) {
		skipOnUnsupportedAlg(t, rw, AlgECC)
		run(t, []byte("\xffTCG123"), true, ecdsa)
	})
}

func TestPCREvent(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()
	debugPCR := uint32(16)
	arbitraryBytes := []byte{1}
	if err := PCREvent(rw, tpmutil.Handle(debugPCR), arbitraryBytes); err != nil {
		t.Fatal(err)
	}
}

func TestPCRExtend(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	tests := []struct {
		desc     string
		hashAlg  Algorithm
		hashSize int
		hashSum  func([]byte) []byte
	}{
		{
			desc:     "SHA1",
			hashAlg:  AlgSHA1,
			hashSize: sha1.Size,
			hashSum: func(in []byte) []byte {
				s := sha1.Sum(in)
				return s[:]
			},
		},
		{
			desc:     "SHA256",
			hashAlg:  AlgSHA256,
			hashSize: sha256.Size,
			hashSum: func(in []byte) []byte {
				s := sha256.Sum256(in)
				return s[:]
			},
		},
	}

	const pcr = int(16)

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			pcrValue := bytes.Repeat([]byte{0xF}, tt.hashSize)
			oldPCRValue, err := ReadPCR(rw, pcr, tt.hashAlg)
			if err != nil {
				t.Fatalf("Can't read PCR %d from the TPM: %s", pcr, err)
			}

			if err = PCRExtend(rw, tpmutil.Handle(pcr), tt.hashAlg, pcrValue, ""); err != nil {
				t.Fatalf("Failed to extend PCR %d: %s", pcr, err)
			}

			newPCRValue, err := ReadPCR(rw, pcr, tt.hashAlg)
			if err != nil {
				t.Fatalf("Can't read PCR %d from the TPM: %s", pcr, err)
			}

			finalPCR := tt.hashSum(append(oldPCRValue, pcrValue...))

			if !bytes.Equal(finalPCR, newPCRValue) {
				t.Fatalf("PCRs not equal, got %x, want %x", finalPCR, newPCRValue)
			}
		})
	}
}

func TestReadPCR(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()
	pcrVal, err := ReadPCR(rw, 16 /*pcr*/, AlgSHA256)
	if err != nil {
		t.Fatal(err)
	}
	if len(pcrVal) != 32 {
		t.Fatalf("Expected a 32 byte PCR value but got: %v", pcrVal)
	}
}

func makeAttestationData() AttestationData {
	signer := tpmutil.Handle(100)
	return AttestationData{
		Magic: 0xff544347,
		QualifiedSigner: Name{
			Handle: &signer,
		},
		ExtraData: []byte("foo"),
		ClockInfo: ClockInfo{
			Clock:        3,
			ResetCount:   4,
			RestartCount: 5,
			Safe:         6,
		},
		FirmwareVersion: 7,
	}
}

func TestEncodeDecodeCertifyAttestationData(t *testing.T) {
	ciQualifiedName := tpmutil.Handle(101)
	ad := makeAttestationData()
	ad.Type = TagAttestCertify
	ad.AttestedCertifyInfo = &CertifyInfo{
		Name: Name{
			Digest: &HashValue{
				Alg:   AlgSHA1,
				Value: make([]byte, crypto.SHA1.Size()),
			},
		},
		QualifiedName: Name{
			Handle: &ciQualifiedName,
		},
	}

	encoded, err := ad.Encode()
	if err != nil {
		t.Fatalf("error encoding AttestationData: %v", err)
	}
	decoded, err := DecodeAttestationData(encoded)
	if err != nil {
		t.Fatalf("error decoding AttestationData: %v", err)
	}

	if !reflect.DeepEqual(*decoded, ad) {
		t.Errorf("got decoded value:\n%v\nwant:\n%v", decoded, ad)
	}
}

func TestEncodeDecodeCreationAttestationData(t *testing.T) {
	ad := makeAttestationData()
	ad.Type = TagAttestCreation
	ad.AttestedCreationInfo = &CreationInfo{
		Name: Name{
			Digest: &HashValue{
				Alg:   AlgSHA1,
				Value: make([]byte, crypto.SHA1.Size()),
			},
		},
		OpaqueDigest: []byte{7, 8, 9},
	}

	encoded, err := ad.Encode()
	if err != nil {
		t.Fatalf("error encoding AttestationData: %v", err)
	}
	decoded, err := DecodeAttestationData(encoded)
	if err != nil {
		t.Fatalf("error decoding AttestationData: %v", err)
	}

	if !reflect.DeepEqual(*decoded, ad) {
		t.Errorf("got decoded value:\n%v\nwant:\n%v", decoded, ad)
	}
}

func TestEncodeDecodePublicDefaultRSAExponent(t *testing.T) {
	p := Public{
		Type:       AlgRSA,
		NameAlg:    AlgSHA1,
		Attributes: FlagSign | FlagSensitiveDataOrigin | FlagUserWithAuth,
		RSAParameters: &RSAParams{
			Sign: &SigScheme{
				Alg:  AlgRSASSA,
				Hash: AlgSHA1,
			},
			KeyBits:     2048,
			ExponentRaw: 1<<16 + 1,
			ModulusRaw:  []byte{1, 2, 3, 4, 7, 8, 9, 9},
		},
	}

	e, err := p.Encode()
	if err != nil {
		t.Fatalf("Public{%+v}.Encode() returned error: %v", p, err)
	}
	d, err := DecodePublic(e)
	if err != nil {
		t.Fatalf("DecodePublic(%v) returned error: %v", e, err)
	}
	if !reflect.DeepEqual(p, d) {
		t.Errorf("RSA TPMT_PUBLIC with default exponent changed after being encoded+decoded")
		t.Logf("\tGot:  %+v", d)
		t.Logf("\tWant: %+v", p)
	}
}

func TestCreateKeyWithSensitive(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()
	parentHandle, _, err := CreatePrimary(rw, HandleOwner, pcrSelection7, emptyPassword, defaultPassword, Public{
		Type:       AlgRSA,
		NameAlg:    AlgSHA256,
		Attributes: FlagRestricted | FlagDecrypt | FlagUserWithAuth | FlagFixedParent | FlagFixedTPM | FlagSensitiveDataOrigin,
		RSAParameters: &RSAParams{
			Symmetric: &SymScheme{
				Alg:     AlgAES,
				KeyBits: 128,
				Mode:    AlgCFB,
			},
			KeyBits: 2048,
		},
	})
	if err != nil {
		t.Fatalf("CreatePrimary failed: %s", err)
	}
	defer FlushContext(rw, parentHandle)

	_, _, _, _, _, err = CreateKeyWithSensitive(rw, parentHandle, pcrSelection7, defaultPassword, defaultPassword, Public{
		Type:       AlgKeyedHash,
		NameAlg:    AlgSHA256,
		Attributes: FlagFixedTPM | FlagFixedParent,
		AuthPolicy: nil,
	}, []byte("test_secret"))
	if err != nil {
		t.Errorf("CreateKeyWithSensitive failed: %s", err)
	}
}

func TestCreateAndCertifyCreation(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()
	params := Public{
		Type:       AlgRSA,
		NameAlg:    AlgSHA256,
		Attributes: FlagSignerDefault | FlagNoDA,
		RSAParameters: &RSAParams{
			Sign: &SigScheme{
				Alg:  AlgRSASSA,
				Hash: AlgSHA256,
			},
			KeyBits: 2048,
		},
	}
	keyHandle, pub, _, creationHash, tix, _, err := CreatePrimaryEx(rw, HandleEndorsement, pcrSelection7, emptyPassword, emptyPassword, params)
	if err != nil {
		t.Fatalf("CreatePrimaryEx failed: %s", err)
	}
	defer FlushContext(rw, keyHandle)

	scheme := SigScheme{Alg: AlgRSASSA, Hash: AlgSHA256, Count: 0}
	attestation, signature, err := CertifyCreation(rw, emptyPassword, keyHandle, keyHandle, nil, creationHash, scheme, tix)
	if err != nil {
		t.Fatalf("CertifyCreation failed: %s", err)
	}
	att, err := DecodeAttestationData(attestation)
	if err != nil {
		t.Fatalf("DecodeAttestationData(%v) failed: %v", attestation, err)
	}
	if att.Type != TagAttestCreation {
		t.Errorf("Got att.Type = %v, want TagAttestCreation", att.Type)
	}
	p, err := DecodePublic(pub)
	if err != nil {
		t.Fatalf("DecodePublic failed: %v", err)
	}
	match, err := att.AttestedCreationInfo.Name.MatchesPublic(p)
	if err != nil {
		t.Fatalf("MatchesPublic failed: %v", err)
	}
	if !match {
		t.Error("Attested name does not match returned public key.")
		t.Logf("Name: %v", att.AttestedCreationInfo.Name)
		t.Logf("Public: %v", p)
	}
	rsaPub := rsa.PublicKey{E: int(p.RSAParameters.Exponent()), N: p.RSAParameters.Modulus()}
	hsh := crypto.SHA256.New()
	hsh.Write(attestation)
	if err := rsa.VerifyPKCS1v15(&rsaPub, crypto.SHA256, hsh.Sum(nil), signature); err != nil {
		t.Errorf("VerifyPKCS1v15 failed: %v", err)
	}
}

func TestCreateAndCertifyCreationECC(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()
	params := Public{
		Type:       AlgECC,
		NameAlg:    AlgSHA256,
		Attributes: FlagSignerDefault,
		ECCParameters: &ECCParams{
			Sign: &SigScheme{
				Alg:  AlgECDSA,
				Hash: AlgSHA256,
			},
			CurveID: CurveNISTP256,
		},
	}
	keyHandle, pub, _, creationHash, tix, _, err := CreatePrimaryEx(rw, HandleEndorsement, pcrSelection7, emptyPassword, emptyPassword, params)
	if err != nil {
		t.Fatalf("CreatePrimaryEx failed: %s", err)
	}
	defer FlushContext(rw, keyHandle)

	scheme := SigScheme{Alg: AlgECDSA, Hash: AlgSHA256, Count: 0}
	attestation, signature, err := CertifyCreation(rw, emptyPassword, keyHandle, keyHandle, nil, creationHash, scheme, tix)
	if err != nil {
		t.Fatalf("CertifyCreation failed: %s", err)
	}

	att, err := DecodeAttestationData(attestation)
	if err != nil {
		t.Fatalf("DecodeAttestationData(%v) failed: %v", attestation, err)
	}
	if att.Type != TagAttestCreation {
		t.Errorf("Got att.Type = %v, want TagAttestCreation", att.Type)
	}
	p, err := DecodePublic(pub)
	if err != nil {
		t.Fatalf("DecodePublic failed: %v", err)
	}
	match, err := att.AttestedCreationInfo.Name.MatchesPublic(p)
	if err != nil {
		t.Fatalf("MatchesPublic failed: %v", err)
	}
	if !match {
		t.Error("Attested name does not match returned public key.")
		t.Logf("Name: %v", att.AttestedCreationInfo.Name)
		t.Logf("Public: %v", p)
	}

	var pkEcdsa ecdsa.PublicKey
	var hsh hash.Hash
	pkEcdsa = ecdsa.PublicKey{Curve: elliptic.P256(), X: p.ECCParameters.Point.X(), Y: p.ECCParameters.Point.Y()}
	signHash, err := p.ECCParameters.Sign.Hash.Hash()
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}
	hsh = signHash.New()
	hsh.Write(attestation)

	r := new(big.Int)
	s := new(big.Int)
	r.SetBytes(signature[:32])
	s.SetBytes(signature[32:])

	if !ecdsa.Verify(&pkEcdsa, hsh.Sum(nil), r, s) {
		t.Fatalf("Verify failed")
	}
}

func TestNVReadWriteAndLocks(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	var (
		idx  tpmutil.Handle = 0x1500000
		data                = []byte("testdata")
		attr                = AttrOwnerWrite | AttrOwnerRead | AttrWriteSTClear | AttrReadSTClear
	)

	// Undefine the space, just in case the previous run of this test failed
	// to clean up.
	if err := NVUndefineSpace(rw, emptyPassword, HandleOwner, idx); err != nil {
		t.Logf("(not a failure) NVUndefineSpace at index 0x%x failed: %v", idx, err)
	}

	// Define space in NV storage and clean up afterwards or subsequent runs will fail.
	if err := NVDefineSpace(rw,
		HandleOwner,
		idx,
		emptyPassword,
		emptyPassword,
		nil,
		attr,
		uint16(len(data)),
	); err != nil {
		t.Fatalf("NVDefineSpace failed: %v", err)
	}
	defer NVUndefineSpace(rw, emptyPassword, HandleOwner, idx)

	// Write the data
	if err := NVWrite(rw, HandleOwner, idx, emptyPassword, data, 0); err != nil {
		t.Fatalf("NVWrite failed: %v", err)
	}

	// Enable write lock
	if err := NVWriteLock(rw, HandleOwner, idx, emptyPassword); err != nil {
		t.Fatalf("NVWriteLock failed: %v", err)
	}

	// Write the data again. Should fail now because it's write-locked.
	err := NVWrite(rw, HandleOwner, idx, emptyPassword, data, 0)
	switch err := err.(type) {
	case nil:
		t.Fatal("NVWrite succeeded after NVWriteLock")
	case Error:
		if err.Code != RCNVLocked {
			t.Fatalf("NVWrite: unexpected error; want RCNVLocked, got %v", err)
		}
	default:
		t.Fatalf("NVWrite: unexpected error; want RCNVLocked, got %v", err)
	}

	// Make sure the public area of the index can be read
	pub, err := NVReadPublic(rw, idx)
	if err != nil {
		t.Fatalf("NVReadPublic failed: %v", err)
	}
	if int(pub.DataSize) != len(data) {
		t.Fatalf("public NV data size mismatch, got %d, want %d, ", pub.DataSize, len(data))
	}

	// Read all of the data with NVReadEx and compare to what was written
	outdata, err := NVReadEx(rw, idx, HandleOwner, emptyPassword, 0)
	if err != nil {
		t.Fatalf("NVReadEx failed: %v", err)
	}
	if !bytes.Equal(data, outdata) {
		t.Fatalf("data read from NV index does not match, got %x, want %x", outdata, data)
	}

	// Enable read lock
	if err := NVReadLock(rw, HandleOwner, idx, emptyPassword); err != nil {
		t.Fatalf("NVReadLock failed: %v (%T)", err, err)
	}

	// Read the data again. Should fail now because it's read-locked.
	if _, err := NVReadEx(rw, idx, HandleOwner, emptyPassword, 0); err == nil {
		t.Fatal("NVRead succeeded after NVReadLock")
	} else if !strings.HasSuffix(err.Error(), ": NV access locked") {
		t.Fatalf("NVRead: unexpected error; want RCNVLocked, got %v", err)
	}
}

func TestPolicySecret(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	sessHandle, _, err := StartAuthSession(rw, HandleNull, HandleNull, make([]byte, 16), nil, SessionPolicy, AlgNull, AlgSHA256)
	if err != nil {
		t.Fatalf("StartAuthSession() failed: %v", err)
	}
	defer FlushContext(rw, sessHandle)

	if _, err := PolicySecret(rw, HandleEndorsement, AuthCommand{Session: HandlePasswordSession, Attributes: AttrContinueSession}, sessHandle, nil, nil, nil, 0); err != nil {
		t.Fatalf("PolicySecret() failed: %v", err)
	}
}

func TestQuote(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()
	params := Public{
		Type:       AlgRSA,
		NameAlg:    AlgSHA256,
		Attributes: FlagSignerDefault | FlagNoDA,
		RSAParameters: &RSAParams{
			Sign: &SigScheme{
				Alg:  AlgRSASSA,
				Hash: AlgSHA256,
			},
			KeyBits: 2048,
		},
	}
	keyHandle, pub, _, _, _, _, err := CreatePrimaryEx(rw, HandleEndorsement, pcrSelection7, emptyPassword, emptyPassword, params)
	if err != nil {
		t.Fatalf("CreatePrimaryEx failed: %s", err)
	}
	defer FlushContext(rw, keyHandle)

	attestation, signature, err := Quote(rw, keyHandle, emptyPassword, emptyPassword, nil, pcrSelection7, AlgNull)
	if err != nil {
		t.Fatalf("Quote failed: %v", err)
	}

	att, err := DecodeAttestationData(attestation)
	if err != nil {
		t.Fatalf("DecodeAttestationData(%v) failed: %v", attestation, err)
	}
	if att.Type != TagAttestQuote {
		t.Errorf("Got att.Type = %v, want TagAttestQuote", att.Type)
	}
	if att.AttestedQuoteInfo == nil {
		t.Error("AttestedQuoteInfo = nil, want non-nil")
	}
	p, err := DecodePublic(pub)
	if err != nil {
		t.Fatalf("DecodePublic failed: %v", err)
	}
	rsaPub := rsa.PublicKey{E: int(p.RSAParameters.Exponent()), N: p.RSAParameters.Modulus()}
	hsh := crypto.SHA256.New()
	hsh.Write(attestation)
	if err := rsa.VerifyPKCS1v15(&rsaPub, crypto.SHA256, hsh.Sum(nil), signature.RSA.Signature); err != nil {
		t.Errorf("VerifyPKCS1v15 failed: %v", err)
	}
}

func TestReadPublicKey(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	run := func(t *testing.T, public Public, private Private, pubKeyIn crypto.PublicKey) {
		t.Helper()

		h, _, err := LoadExternal(rw, public, private, HandleNull)
		if err != nil {
			t.Fatal(err)
		}
		defer FlushContext(rw, h)

		pub, _, _, err := ReadPublic(rw, h)
		if err != nil {
			t.Fatalf("ReadPublic failed: %s", err)
		}

		pubKeyOut, err := pub.Key()
		if err != nil {
			t.Fatalf("Public.Key() failed: %s", err)
		}

		if !reflect.DeepEqual(pubKeyIn, pubKeyOut) {
			t.Fatalf("Public.Key() = %#v; want %#v", pubKeyOut, pubKeyIn)
		}
	}

	t.Run("RSA", func(t *testing.T) {
		pk, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}
		rp := Public{
			Type:       AlgRSA,
			NameAlg:    AlgSHA1,
			Attributes: FlagSign | FlagSensitiveDataOrigin | FlagUserWithAuth,
			RSAParameters: &RSAParams{
				Sign: &SigScheme{
					Alg:  AlgRSASSA,
					Hash: AlgSHA1,
				},
				KeyBits:     2048,
				ExponentRaw: uint32(pk.PublicKey.E),
				ModulusRaw:  pk.PublicKey.N.Bytes(),
			},
		}
		private := Private{
			Type:      AlgRSA,
			Sensitive: pk.Primes[0].Bytes(),
		}
		run(t, rp, private, &pk.PublicKey)
	})
	t.Run("ECC", func(t *testing.T) {
		skipOnUnsupportedAlg(t, rw, AlgECC)
		pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		public := Public{
			Type:       AlgECC,
			NameAlg:    AlgSHA1,
			Attributes: FlagSign | FlagSensitiveDataOrigin | FlagUserWithAuth,
			ECCParameters: &ECCParams{
				Sign: &SigScheme{
					Alg:  AlgECDSA,
					Hash: AlgSHA1,
				},
				CurveID: CurveNISTP256,
				Point:   ECPoint{XRaw: pk.PublicKey.X.Bytes(), YRaw: pk.PublicKey.Y.Bytes()},
			},
		}
		private := Private{
			Type:      AlgECC,
			Sensitive: pk.D.Bytes(),
		}
		run(t, public, private, &pk.PublicKey)
	})
}

func TestEncryptDecrypt(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	parentHandle, _, err := CreatePrimary(rw, HandleOwner, pcrSelection7, emptyPassword, defaultPassword, Public{
		Type:       AlgRSA,
		NameAlg:    AlgSHA256,
		Attributes: FlagRestricted | FlagDecrypt | FlagUserWithAuth | FlagFixedParent | FlagFixedTPM | FlagSensitiveDataOrigin,
		RSAParameters: &RSAParams{
			Symmetric: &SymScheme{
				Alg:     AlgAES,
				KeyBits: 128,
				Mode:    AlgCFB,
			},
			KeyBits: 2048,
		},
	})
	if err != nil {
		t.Fatalf("CreatePrimary failed: %s", err)
	}
	defer FlushContext(rw, parentHandle)
	privateBlob, publicBlob, _, _, _, err := CreateKey(rw, parentHandle, pcrSelection7, defaultPassword, defaultPassword, Public{
		Type:       AlgSymCipher,
		NameAlg:    AlgSHA256,
		Attributes: FlagDecrypt | FlagSign | FlagUserWithAuth | FlagFixedParent | FlagFixedTPM | FlagSensitiveDataOrigin,
		SymCipherParameters: &SymCipherParams{
			Symmetric: &SymScheme{
				Alg:     AlgAES,
				KeyBits: 128,
				Mode:    AlgCFB,
			},
		},
	})
	if err != nil {
		t.Fatalf("CreateKey failed: %s", err)
	}
	key, _, err := Load(rw, parentHandle, defaultPassword, publicBlob, privateBlob)
	if err != nil {
		t.Fatalf("Load failed: %s", err)
	}
	defer FlushContext(rw, key)

	data := bytes.Repeat([]byte("a"), 1e4) // 10KB
	iv := make([]byte, 16)

	encrypted, err := EncryptSymmetric(rw, defaultPassword, key, iv, data)
	if err != nil {
		t.Fatalf("EncryptSymmetric failed: %s", err)
	}
	if bytes.Equal(encrypted, data) {
		t.Error("encrypted blob matches unenecrypted data")
	}
	decrypted, err := DecryptSymmetric(rw, defaultPassword, key, iv, encrypted)
	if err != nil {
		t.Fatalf("DecryptSymmetric failed: %s", err)
	}
	if !bytes.Equal(decrypted, data) {
		t.Errorf("got decrypted data: %q, want: %q", decrypted, data)
	}
}

func TestRSAEncryptDecrypt(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	handle, _, err := CreatePrimary(rw, HandleOwner, pcrSelection7, emptyPassword, defaultPassword, Public{
		Type:       AlgRSA,
		NameAlg:    AlgSHA256,
		Attributes: FlagDecrypt | FlagUserWithAuth | FlagFixedParent | FlagFixedTPM | FlagSensitiveDataOrigin,
		RSAParameters: &RSAParams{
			Sign: &SigScheme{
				Alg:  AlgNull,
				Hash: AlgNull,
			},
			KeyBits: 2048,
		},
	})
	if err != nil {
		t.Fatalf("CreatePrimary failed: %s", err)
	}
	defer FlushContext(rw, handle)

	tests := map[string]struct {
		scheme *AsymScheme
		data   []byte
		label  string
	}{
		"No padding": {
			scheme: &AsymScheme{Alg: AlgNull},
			data:   bytes.Repeat([]byte("a"), 256),
		},
		"RSAES-PKCS1": {
			scheme: &AsymScheme{Alg: AlgRSAES},
			data:   bytes.Repeat([]byte("a"), 245),
		},
		"RSAES-OAEP-SHA1": {
			scheme: &AsymScheme{Alg: AlgOAEP, Hash: AlgSHA1},
			data:   bytes.Repeat([]byte("a"), 214),
		},
		"RSAES-OAEP-SHA256": {
			scheme: &AsymScheme{Alg: AlgOAEP, Hash: AlgSHA256},
			data:   bytes.Repeat([]byte("a"), 190),
		},
		"RSAES-OAEP-SHA256 with label": {
			scheme: &AsymScheme{Alg: AlgOAEP, Hash: AlgSHA256},
			data:   bytes.Repeat([]byte("a"), 190),
			label:  "label",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			encrypted, err := RSAEncrypt(rw, handle, test.data, test.scheme, test.label)
			if err != nil {
				t.Fatal("RSAEncrypt failed:", err)
			}
			decrypted, err := RSADecrypt(rw, handle, defaultPassword, encrypted, test.scheme, test.label)
			if err != nil {
				t.Fatal("RSADecrypt failed:", err)
			}
			if !bytes.Equal(decrypted, test.data) {
				t.Errorf("got decrypted data: %q, want: %q", decrypted, test.data)
			}
		})
	}

}

func TestECDHKeyGen(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	// Generate my key and load the public point into the TPM.
	myPriv, myPubX, myPubY, err := elliptic.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}
	handle, _, err := LoadExternal(rw, Public{
		Type:       AlgECC,
		NameAlg:    AlgSHA256,
		Attributes: FlagDecrypt,
		ECCParameters: &ECCParams{
			CurveID: CurveNISTP256,
			Point:   ECPoint{XRaw: myPubX.Bytes(), YRaw: myPubY.Bytes()},
		},
	}, Private{}, HandleOwner)
	if err != nil {
		t.Fatalf("LoadExternal failed: %v", err)
	}
	defer FlushContext(rw, handle)

	// Ask the TPM to multiply an ephemeral priv by our pub.
	yourZ, yourPub, err := ECDHKeyGen(rw, handle)
	if err != nil {
		t.Fatalf("ECDHKeyGen failed: %v", err)
	}

	// Same calculation on our side: multiply our priv by the TPM's ephemeral pub.
	myZX, myZY := elliptic.P256().ScalarMult(yourPub.X(), yourPub.Y(), myPriv)

	if myZX.Cmp(yourZ.X()) != 0 || myZY.Cmp(yourZ.Y()) != 0 {
		t.Errorf("want (%x, %x) got (%x, %x)", myZX.Bytes(), myZY.Bytes(), yourZ.X().Bytes(), yourZ.Y().Bytes())
	}
}

func TestECDHZGen(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	// Generate our key.
	myPriv, myPubX, myPubY, err := elliptic.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	// Generate a key in the TPM.
	handle, _, err := CreatePrimary(rw, HandleOwner, PCRSelection{}, emptyPassword, defaultPassword, Public{
		Type:       AlgECC,
		NameAlg:    AlgSHA256,
		Attributes: FlagDecrypt | FlagSensitiveDataOrigin | FlagUserWithAuth,
		ECCParameters: &ECCParams{
			CurveID: CurveNISTP256,
		},
	})
	if err != nil {
		t.Fatalf("CreatePrimary failed: %v", err)
	}
	defer FlushContext(rw, handle)

	// Read the public key from the TPM.
	yourPub, _, _, err := ReadPublic(rw, handle)
	if err != nil {
		t.Fatalf("ReadPublic failed: %v", err)
	}

	// Ask the TPM to multiply our public point by its private key.
	yourZ, err := ECDHZGen(rw, handle, defaultPassword, ECPoint{
		XRaw: myPubX.Bytes(),
		YRaw: myPubY.Bytes(),
	})
	if err != nil {
		t.Fatalf("ECDHZGen failed: %v", err)
	}

	// Same calculation on our side: multiply our priv by the TPM's pub.
	myZX, myZY := elliptic.P256().ScalarMult(yourPub.ECCParameters.Point.X(), yourPub.ECCParameters.Point.Y(), myPriv)

	if myZX.Cmp(yourZ.X()) != 0 || myZY.Cmp(yourZ.Y()) != 0 {
		t.Errorf("want (%x, %x) got (%x, %x)", myZX.Bytes(), myZY.Bytes(), yourZ.X().Bytes(), yourZ.Y().Bytes())
	}
}

func TestCreatePrimaryRawTemplate(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	pubRaw, err := defaultKeyParams.Encode()
	if err != nil {
		t.Fatal(err)
	}
	kh, pub, err := CreatePrimaryRawTemplate(rw, HandleEndorsement, PCRSelection{}, "", "", pubRaw)
	if err != nil {
		t.Fatalf("CreatePrimary failed: %v", err)
	}
	defer FlushContext(rw, kh)

	pubRSA, ok := pub.(*rsa.PublicKey)
	if !ok {
		t.Fatalf("got public key type %T, want *rsa.PublicKey", pub)
	}
	gotKeySize, wantKeySize := pubRSA.Size(), int(defaultKeyParams.RSAParameters.KeyBits/8)
	if gotKeySize != wantKeySize {
		t.Errorf("got key size %v, want %v", gotKeySize, wantKeySize)
	}
	if pubRSA.E != int(defaultKeyParams.RSAParameters.Exponent()) {
		t.Errorf("got key exponent %v, want %v", pubRSA.E, defaultKeyParams.RSAParameters.Exponent())
	}
}

func TestMatchesTemplate(t *testing.T) {
	tests := []struct {
		name       string
		makePublic func() Public
		goodChange func(*Public)
		badChange  func(*Public)
	}{
		{
			"RSA",
			func() Public {
				return Public{
					Type:       AlgRSA,
					NameAlg:    AlgSHA256,
					Attributes: FlagSignerDefault,
					RSAParameters: &RSAParams{
						Sign: &SigScheme{
							Alg:  AlgRSASSA,
							Hash: AlgSHA256,
						},
						KeyBits: 2048,
					},
				}
			},
			func(pub *Public) { pub.RSAParameters.ModulusRaw = make([]byte, 256) },
			func(pub *Public) { pub.RSAParameters.KeyBits = 1024 },
		},
		{
			"ECC",
			func() Public {
				return Public{
					Type:       AlgECC,
					NameAlg:    AlgSHA256,
					Attributes: FlagSignerDefault,
					ECCParameters: &ECCParams{
						Sign: &SigScheme{
							Alg:  AlgECDSA,
							Hash: AlgSHA256,
						},
						CurveID: CurveNISTP256,
					},
				}
			},
			func(pub *Public) { pub.ECCParameters.Point.XRaw = make([]byte, 32) },
			func(pub *Public) { pub.ECCParameters.CurveID = CurveNISTP384 },
		},
		{
			"SymCipher",
			func() Public {
				return Public{
					Type:       AlgSymCipher,
					NameAlg:    AlgSHA256,
					Attributes: FlagSignerDefault,
					SymCipherParameters: &SymCipherParams{
						Symmetric: &SymScheme{
							Alg:     AlgAES,
							KeyBits: 128,
							Mode:    AlgCFB,
						},
					},
				}
			},
			func(pub *Public) { pub.SymCipherParameters.Unique = make([]byte, 256) },
			func(pub *Public) { pub.SymCipherParameters.Symmetric.KeyBits = 256 },
		},
		{
			"KeyedHash",
			func() Public {
				return Public{
					Type:       AlgKeyedHash,
					NameAlg:    AlgSHA256,
					Attributes: FlagSignerDefault,
					KeyedHashParameters: &KeyedHashParams{
						Alg:  AlgHMAC,
						Hash: AlgSHA256,
					},
				}
			},
			func(pub *Public) { pub.KeyedHashParameters.Unique = make([]byte, 256) },
			func(pub *Public) { pub.KeyedHashParameters.Hash = AlgSHA1 },
		},
		{
			"TypeMismatch",
			func() Public {
				return Public{
					Type:                AlgKeyedHash,
					NameAlg:             AlgSHA256,
					Attributes:          FlagSignerDefault,
					KeyedHashParameters: &KeyedHashParams{Alg: AlgNull},
				}
			},
			func(pub *Public) { pub.KeyedHashParameters.Unique = make([]byte, 256) },
			func(pub *Public) { pub.Type = AlgRSA },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			template := test.makePublic()
			pub := test.makePublic()

			test.goodChange(&pub)
			if !pub.MatchesTemplate(template) {
				t.Error("Change should not cause template mismatch")
			}

			encTmpl, err := template.Encode()
			if err != nil {
				t.Fatal(err)
			}
			decTmpl, err := DecodePublic(encTmpl)
			if err != nil {
				t.Fatal(err)
			}
			if !pub.MatchesTemplate(decTmpl) {
				t.Error("Encoding/Decoding should not cause template mismatch")
			}

			test.badChange(&pub)
			if pub.MatchesTemplate(template) {
				t.Error("Change should cause template mismatch")
			}
		})
	}
}

func TestPlainImport(t *testing.T) {
	// As this test imports a key without using an inner or outer wrapper, the
	// sensitive data is NOT encrypted. This setup should not actually be used.
	rw := openTPM(t)
	defer rw.Close()

	// Create an EK that doesn't require authorization
	ekHandle, _, err := CreatePrimary(rw, HandleEndorsement, PCRSelection{}, "", "", defaultKeyParams)
	if err != nil {
		t.Fatalf("CreatePrimary failed: %s", err)
	}
	defer FlushContext(rw, ekHandle)
	emptyAuth := AuthCommand{Session: HandlePasswordSession, Attributes: AttrContinueSession}

	inBuff := make([]byte, 18)
	io.ReadFull(rand.Reader, inBuff)
	// Create a private area containing the input
	private := Private{
		Type:      AlgKeyedHash,
		AuthValue: nil,
		SeedValue: make([]byte, 32),
		Sensitive: inBuff,
	}
	io.ReadFull(rand.Reader, private.SeedValue)

	privArea, err := private.Encode()
	if err != nil {
		t.Fatalf("Private encoding error: %s", err)
	}
	// We just use a plain Pack here, as there is no integrity checking.
	duplicate, err := tpmutil.Pack(tpmutil.U16Bytes(privArea))
	if err != nil {
		t.Fatalf("Duplicate encoding error: %s", err)
	}

	// The corresponding Public area contains the hash of the Private area.
	privHash := crypto.SHA256.New()
	privHash.Write(private.SeedValue)
	privHash.Write(private.Sensitive)
	public := Public{
		Type:       AlgKeyedHash,
		NameAlg:    AlgSHA256,
		Attributes: FlagUserWithAuth,
		KeyedHashParameters: &KeyedHashParams{
			Alg:    AlgNull,
			Unique: privHash.Sum(nil),
		},
	}
	pubArea, err := public.Encode()
	if err != nil {
		t.Fatalf("Public encoding error: %s", err)
	}

	privInternal, err := Import(rw, ekHandle, emptyAuth, pubArea, duplicate, nil, nil, nil)
	if err != nil {
		t.Fatalf("Import failed: %s", err)
	}
	newHandle, _, err := Load(rw, ekHandle, "", pubArea, privInternal)
	if err != nil {
		t.Fatalf("Load failed: %s", err)
	}
	defer FlushContext(rw, newHandle)

	outBuff, err := Unseal(rw, newHandle, "")
	if err != nil {
		t.Fatalf("Unseal failed: %s", err)
	}
	if !bytes.Equal(outBuff, inBuff) {
		t.Errorf("Got %X, expected %X", outBuff, inBuff)
	}
}

func TestClear(t *testing.T) {
	if !*runClear {
		t.Skip("Missing flag: run-clear. Test skipped")
	}

	rw := openTPM(t)
	defer rw.Close()

	err := Clear(rw, HandleLockout, AuthCommand{Session: HandlePasswordSession, Attributes: AttrContinueSession})
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}
}

func TestHierarchyChangeAuth(t *testing.T) {
	if !*runClear {
		t.Skip("Missing flag: run-clear. Test skipped")
	}

	rw := openTPM(t)
	defer rw.Close()

	err := Clear(rw, HandleLockout, AuthCommand{Session: HandlePasswordSession, Attributes: AttrContinueSession})
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	err = HierarchyChangeAuth(rw, HandleOwner, AuthCommand{Session: HandlePasswordSession, Attributes: AttrContinueSession}, "pass1")
	if err != nil {
		t.Fatalf("HierarchyChangeAuth failed: %v", err)
	}

	// try to set password again, providing invalid auth
	err = HierarchyChangeAuth(rw, HandleOwner, AuthCommand{Session: HandlePasswordSession, Attributes: AttrContinueSession}, "pass2")
	if err == nil {
		t.Fatal("Expected HierarchyChangeAuth to fail")
	}

	// set password again, providing valid auth
	err = HierarchyChangeAuth(rw, HandleOwner, AuthCommand{Session: HandlePasswordSession, Attributes: AttrContinueSession, Auth: []byte("pass1")}, "pass3")
	if err != nil {
		t.Fatalf("HierarchyChangeAuth failed: %v", err)
	}

	err = Clear(rw, HandleLockout, AuthCommand{Session: HandlePasswordSession, Attributes: AttrContinueSession})
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

}

func TestPolicyPCR(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	sessHandle, _, err := StartAuthSession(rw, HandleNull, HandleNull, make([]byte, 16), nil, SessionPolicy, AlgNull, AlgSHA1)
	if err != nil {
		t.Fatalf("StartAuthSession failed: %v", err)
	}
	defer FlushContext(rw, sessHandle)

	sel := PCRSelection{Hash: AlgSHA1, PCRs: []int{0, 1, 2, 3}}
	pcrs, err := ReadPCRs(rw, sel)
	if err != nil {
		t.Fatalf("ReadPCRs failed: %s", err)
	}

	// From TPM 2.0 Part 1, Selecting Multiple PCR:
	// "The list of selectors is processed in order. The selected PCR are
	// concatenated, with the lowest numbered PCR in the first selector being
	// the first in the list and the highest numbered PCR in the last selector
	// being the last."
	//
	// Concatenate all selected PCR values before hashing. sel.PCRs is already
	// sorted in ascending order.
	var expectedVal []byte
	for _, pcr := range sel.PCRs {
		expectedVal = append(expectedVal, pcrs[pcr]...)
	}
	t.Logf("expectedVal=%x", expectedVal)

	// Hash algorithm must match the one in StartAuthSession.
	expectedDigest := sha1.Sum(expectedVal)
	t.Logf("expectedDigest=%x", expectedDigest)

	tests := []struct {
		desc           string
		expectedDigest []byte
		wantErr        bool
	}{
		// PolicyPCR expects a *digest* of the PCR value, not the value itself.
		// Make sure PolicyPCR with value actually fails.
		{desc: "success: digest passed in expectedDigest", expectedDigest: expectedDigest[:]},
		{desc: "fail: value passed in expectedDigest", expectedDigest: expectedVal, wantErr: true},
		{desc: "success: nil passed in expectedDigest", expectedDigest: nil},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := PolicyPCR(rw, sessHandle, tt.expectedDigest, sel)
			if (err != nil) != tt.wantErr {
				t.Fatalf("got error: %v, want error: %v", err != nil, tt.wantErr)
			}
		})
	}
}

func TestDictionaryAttackParameters(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	testCases := []struct {
		name            string
		maxTries        uint32
		recoveryTime    uint32
		lockoutRecovery uint32
	}{
		{
			name:            "0 0 0",
			maxTries:        0,
			recoveryTime:    0,
			lockoutRecovery: 0,
		},
		{
			name:            "5 10 10",
			maxTries:        5,
			recoveryTime:    10,
			lockoutRecovery: 10,
		},
	}

	auth := AuthCommand{
		Session:    HandlePasswordSession,
		Attributes: AttrContinueSession,
		Auth:       []byte(emptyPassword),
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := DictionaryAttackParameters(rw, auth, tc.maxTries, tc.recoveryTime, tc.lockoutRecovery); err != nil {
				t.Fatalf("DictionaryAttackParameters failed: %v", err)
			}
			caps, _, err := GetCapability(rw, CapabilityTPMProperties, 3, uint32(MaxAuthFail))
			if err != nil {
				t.Fatalf("GetCapability failed: %v", err)
			}
			if caps[0].(TaggedProperty).Value != tc.maxTries {
				t.Fatalf("got %d, expected %d", caps[0].(TaggedProperty).Value, tc.maxTries)
			}
			if caps[1].(TaggedProperty).Value != tc.recoveryTime {
				t.Fatalf("got %d, expected %d", caps[1].(TaggedProperty).Value, tc.recoveryTime)
			}
			if caps[2].(TaggedProperty).Value != tc.lockoutRecovery {
				t.Fatalf("got %d, expected %d", caps[2].(TaggedProperty).Value, tc.lockoutRecovery)
			}
		})
	}
}

func TestDictionaryAttackLockReset(t *testing.T) {
	rw := openTPM(t)
	defer rw.Close()

	auth := AuthCommand{
		Session:    HandlePasswordSession,
		Attributes: AttrContinueSession,
		Auth:       []byte(emptyPassword),
	}
	err := DictionaryAttackParameters(rw, auth, 5, 60, 60)
	if err != nil {
		t.Fatalf("DictionaryAttackParameters failed: %v", err)
	}

	handle, _, err := CreatePrimary(rw, HandleOwner, pcrSelection7, emptyPassword, defaultPassword, Public{
		Type:       AlgRSA,
		NameAlg:    AlgSHA256,
		Attributes: FlagDecrypt | FlagUserWithAuth | FlagFixedParent | FlagFixedTPM | FlagSensitiveDataOrigin,
		RSAParameters: &RSAParams{
			Sign: &SigScheme{
				Alg:  AlgNull,
				Hash: AlgNull,
			},
			KeyBits: 2048,
		},
	})
	if err != nil {
		t.Fatalf("CreatePrimary failed: %s", err)
	}
	defer FlushContext(rw, handle)

	scheme := &AsymScheme{Alg: AlgOAEP, Hash: AlgSHA256}
	label := "label"
	encrypted, err := RSAEncrypt(rw, handle, bytes.Repeat([]byte("a"), 190), scheme, label)
	if err != nil {
		t.Fatalf("RSAEncrypt failed: %v", err)
	}
	// try RSADecrypt with bad password
	if _, err = RSADecrypt(rw, handle, "bad password", encrypted, scheme, label); err != nil {
		if serr, ok := err.(SessionError); !ok || serr.Code != RCAuthFail {
			t.Fatalf("RSADecrypt fails with unexpected error: %v", err)
		}
	}

	caps, _, err := GetCapability(rw, CapabilityTPMProperties, 1, uint32(LockoutCounter))
	if err != nil {
		t.Fatalf("GetCapability failed: %v", err)
	}
	if caps[0].(TaggedProperty).Value != 1 {
		t.Fatalf("got %d, expected 1", caps[0].(TaggedProperty).Value)
	}

	if err = DictionaryAttackLockReset(rw, auth); err != nil {
		t.Fatalf("DictionaryAttackLockReset failed: %v", err)
	}

	caps, _, err = GetCapability(rw, CapabilityTPMProperties, 1, uint32(LockoutCounter))
	if err != nil {
		t.Fatalf("GetCapability failed: %v", err)
	}
	if caps[0].(TaggedProperty).Value != 0 {
		t.Fatalf("got %d, expected 0", caps[0].(TaggedProperty).Value)
	}
}
