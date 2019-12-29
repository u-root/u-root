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
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-tpm/tpmutil"
)

var (
	ownerAuthEnvVar = "TPM_OWNER_AUTH"
	srkAuthEnvVar   = "TPM_SRK_AUTH"
	aikAuthEnvVar   = "TPM_AIK_AUTH"
	tpmPathEnvVar   = "TPM_PATH"
)

// getAuth looks in the environment variables to find a given auth input value.
// If the environment variable is not present, then getAuth returns the
// well-known auth value of 20 bytes of zeros.
func getAuth(name string) digest {
	var auth digest
	authInput := os.Getenv(name)
	if authInput != "" {
		aa := sha1.Sum([]byte(authInput))
		copy(auth[:], aa[:])
	}
	return auth
}

func TestGetKeys(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	handles, err := GetKeys(rwc)
	if err != nil {
		t.Fatal("Couldn't enumerate keys in the TPM:", err)
	}

	t.Logf("Got %d keys: % d\n", len(handles), handles)
}

func TestGetManufacturer(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	vendorID, err := GetManufacturer(rwc)
	if err != nil {
		t.Fatal("Couldn't read VendorID from TPM:", err)
	}

	t.Logf("TPM VendorID: %v\n", vendorID)
}

func TestGetAlgs(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	algs, err := GetAlgs(rwc)
	if err != nil {
		t.Fatalf("Couldn't read Algorithms: %v", err)
	}
	want := []Algorithm{AlgRSA, AlgSHA, AlgHMAC, AlgMGF1}
outer:
	for _, alg := range want {
		for _, got := range algs {
			if got == alg {
				continue outer
			}
			t.Errorf("GetAlgs returned %v, which does not contain a required algorithm %v", algs, alg.String())
		}
	}
}

func TestGetNVList(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	nvList, err := GetNVList(rwc)
	if err != nil {
		t.Fatalf("Couldn't read NVList %v", err)
	}
	t.Logf("NVList is: %v", nvList)
}

func TestGetNVIndex(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	nvList, err := GetNVList(rwc)
	if err != nil {
		t.Fatalf("Couldn't read NVList %v", err)
	}
	var nvInfo []NVDataPublic
	for _, nvEntry := range nvList {
		index, err := GetNVIndex(rwc, nvEntry)
		if err != nil {
			t.Fatalf("Can't read NVDataPublic of index: %v with: %v", nvEntry, err)
		}
		nvInfo = append(nvInfo, index)
	}
	t.Logf("NVIndices with Attributes:%v", nvInfo)
}

func TestPcrExtend(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	var pcrValue [20]byte
	var value = "FFFFFFFFFFFFFFFFFFFF"
	copy(pcrValue[:], value)

	oldPcrValue, err := ReadPCR(rwc, 12)
	if err != nil {
		t.Fatal("Couldn't read PCR 12 from the TPM:", err)
	}

	newPcrValue, err := PcrExtend(rwc, 12, pcrValue)
	if err != nil {
		t.Fatal("Couldn't extend PCR 12 from the TPM:", err)
	}

	finalPcr := sha1.Sum(append(oldPcrValue, pcrValue[:]...))

	if bytes.Equal(finalPcr[:], newPcrValue) {
		t.Logf("PCR are equal!\n")
	} else {
		t.Fatal("PCR are not equal! Test failed.\n")
	}
}

func TestReadEKCert(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	ownAuth := getAuth(ownerAuthEnvVar)
	cert, err := ReadEKCert(rwc, ownAuth)
	if err != nil {
		t.Fatal("Unable to read EKCert from NVRAM:", err)
	}

	x509cert, err := x509.ParseCertificate(cert)
	if err != nil {
		t.Logf("Malformed certificate: %v\n", err)
	} else {
		t.Logf("Certificate: %v\n", x509cert)
	}
}

func TestReadPCR(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	res, err := ReadPCR(rwc, 18)
	if err != nil {
		t.Fatal("Couldn't read PCR 18 from the TPM:", err)
	}

	t.Logf("Got PCR 18 value % x\n", res)
}

func TestFetchPCRValues(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	var mask pcrMask
	if err := mask.setPCR(17); err != nil {
		t.Fatal("Couldn't set PCR 17:", err)
	}

	pcrs, err := FetchPCRValues(rwc, []int{17})
	if err != nil {
		t.Fatal("Couldn't get PCRs 17:", err)
	}

	comp, err := createPCRComposite(mask, pcrs)
	if err != nil {
		t.Fatal("Couldn't create PCR composite")
	}

	if len(comp) != int(digestSize) {
		t.Fatal("Invalid PCR composite")
	}

	var locality byte
	_, err = createPCRInfoLong(locality, mask, pcrs)
	if err != nil {
		t.Fatal("Couldn't create a pcrInfoLong structure for these PCRs")
	}
}

func TestGetRandom(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	// Try to get 16 bytes of randomness from the TPM.
	b, err := GetRandom(rwc, 16)
	if err != nil {
		t.Fatal("Couldn't get 16 bytes of randomness from the TPM:", err)
	}

	if len(b) != 16 {
		t.Fatal("Couldn't get 16 bytes of randomness from the TPM")
	}
}

func TestOIAP(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	// Get auth info from OIAP.
	_, err := oiap(rwc)
	if err != nil {
		t.Fatal("Couldn't run OIAP:", err)
	}
}

func TestOSAP(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	// Try to run OSAP for the SRK.
	osapc := &osapCommand{
		EntityType:  etSRK,
		EntityValue: khSRK,
	}

	if _, err := rand.Read(osapc.OddOSAP[:]); err != nil {
		t.Fatal("Couldn't get a random odd OSAP nonce")
	}

	_, err := osap(rwc, osapc)
	if err != nil {
		t.Fatal("Couldn't run OSAP:", err)
	}
}

func TestResizeableSlice(t *testing.T) {
	// Set up an encoded slice with a byte array.
	ra := &responseAuth{
		NonceEven:   [20]byte{},
		ContSession: 1,
		Auth:        [20]byte{},
	}

	b := make([]byte, 322)
	if _, err := rand.Read(b); err != nil {
		t.Fatal("Couldn't read random bytes into the byte array")
	}

	bb, err := tpmutil.Pack(ra, tpmutil.U32Bytes(b))
	if err != nil {
		t.Fatal("Couldn't pack the bytes:", err)
	}

	var ra2 responseAuth
	var b2 tpmutil.U32Bytes
	if _, err := tpmutil.Unpack(bb, &ra2, &b2); err != nil {
		t.Fatal("Couldn't unpack the resizeable values:", err)
	}

	if !bytes.Equal(b2, b) {
		t.Fatal("ResizeableSlice was not resized or copied correctly")
	}
}

func TestSeal(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	data := make([]byte, 64)
	data[0] = 137
	data[1] = 138
	data[2] = 139

	srkAuth := getAuth(srkAuthEnvVar)
	sealed, err := Seal(rwc, 0 /* locality 0 */, []int{17} /* PCR 17 */, data, srkAuth[:])
	if err != nil {
		t.Fatal("Couldn't seal the data:", err)
	}

	data2, err := Unseal(rwc, sealed, srkAuth[:])
	if err != nil {
		t.Fatal("Couldn't unseal the data:", err)
	}

	if !bytes.Equal(data2, data) {
		t.Fatal("Unsealed data doesn't match original data")
	}
}

func TestReseal(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	data := make([]byte, 64)
	data[0] = 137
	data[1] = 138
	data[2] = 139

	pcrMap := make(map[int][]byte)
	pcrMap[23] = make([]byte, 20)
	pcrMap[16] = make([]byte, 20)
	srkAuth := getAuth(srkAuthEnvVar)
	sealed, err := Reseal(rwc, 0 /* locality 0 */, pcrMap, data, srkAuth[:])
	if err != nil {
		t.Fatal("Couldn't seal the data:", err)
	}

	data2, err := Unseal(rwc, sealed, srkAuth[:])
	if err != nil {
		t.Fatal("Couldn't unseal the data:", err)
	}

	if !bytes.Equal(data2, data) {
		t.Fatal("Unsealed data doesn't match original data")
	}
}

func TestLoadKey2(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	// Get the key from aikblob, assuming it exists. Otherwise, skip the test.
	blob, err := ioutil.ReadFile("./aikblob")
	if err != nil {
		t.Skip("No aikblob file; skipping test")
	}

	// We're using the well-known authenticator of 20 bytes of zeros.
	srkAuth := getAuth(srkAuthEnvVar)
	handle, err := LoadKey2(rwc, blob, srkAuth[:])
	if err != nil {
		t.Fatal("Couldn't load the AIK into the TPM and get a handle for it:", err)
	}

	if err := CloseKey(rwc, handle); err != nil {
		t.Fatal("Couldn't flush the AIK from the TPM:", err)
	}
}

func TestQuote2(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	// Get the key from aikblob, assuming it exists. Otherwise, skip the test.
	blob, err := ioutil.ReadFile("./aikblob")
	if err != nil {
		t.Skip("No aikblob file; skipping test")
	}

	// Load the AIK for the quote.
	// We're using the well-known authenticator of 20 bytes of zeros.
	srkAuth := getAuth(srkAuthEnvVar)
	handle, err := LoadKey2(rwc, blob, srkAuth[:])
	if err != nil {
		t.Fatal("Couldn't load the AIK into the TPM and get a handle for it:", err)
	}
	defer CloseKey(rwc, handle)

	// Data to quote.
	data := []byte(`The OS says this test is good`)
	aikAuth := getAuth(aikAuthEnvVar)
	q, err := Quote2(rwc, handle, data, []int{17, 18}, 1 /* addVersion */, aikAuth[:])
	if err != nil {
		t.Fatal("Couldn't quote the data:", err)
	}

	if len(q) == 0 {
		t.Fatal("Couldn't get a quote using an AIK")
	}
}

func TestGetPubKey(t *testing.T) {
	// For testing purposes, use the aikblob if it exists. Otherwise, just skip
	// this test. TODO(tmroeder): implement AIK creation so we can always run
	// this test.
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	// Get the key from aikblob, assuming it exists. Otherwise, skip the test.
	blob, err := ioutil.ReadFile("./aikblob")
	if err != nil {
		t.Skip("No aikblob file; skipping test")
	}

	// Load the AIK for the quote.
	// We're using the well-known authenticator of 20 bytes of zeros.
	srkAuth := getAuth(srkAuthEnvVar)
	handle, err := LoadKey2(rwc, blob, srkAuth[:])
	if err != nil {
		t.Fatal("Couldn't load the AIK into the TPM and get a handle for it:", err)
	}
	defer CloseKey(rwc, handle)

	k, err := GetPubKey(rwc, handle, srkAuth[:])
	if err != nil {
		t.Fatal("Couldn't get the pub key for the AIK")
	}

	if len(k) == 0 {
		t.Fatal("Couldn't get a pubkey blob from an AIK")
	}
}

func TestQuote(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	// Get the key from aikblob, assuming it exists. Otherwise, skip the test.
	blob, err := ioutil.ReadFile("./aikblob")
	if err != nil {
		t.Skip("No aikblob file; skipping test")
	}

	// Load the AIK for the quote.
	srkAuth := getAuth(srkAuthEnvVar)
	handle, err := LoadKey2(rwc, blob, srkAuth[:])
	if err != nil {
		t.Fatal("Couldn't load the AIK into the TPM and get a handle for it:", err)
	}
	defer CloseKey(rwc, handle)

	// Data to quote.
	data := []byte(`The OS says this test is good`)
	pcrNums := []int{17, 18}
	aikAuth := getAuth(aikAuthEnvVar)
	q, values, err := Quote(rwc, handle, data, pcrNums, aikAuth[:])
	if err != nil {
		t.Fatal("Couldn't quote the data:", err)
	}

	// Verify the quote.
	pk, err := UnmarshalRSAPublicKey(blob)
	if err != nil {
		t.Fatal("Couldn't extract an RSA key from the AIK blob:", err)
	}

	if err := VerifyQuote(pk, data, q, pcrNums, values); err != nil {
		t.Fatal("The quote didn't pass verification:", err)
	}
}

func TestUnmarshalRSAPublicKey(t *testing.T) {
	// Get the key from aikblob, assuming it exists. Otherwise, skip the test.
	blob, err := ioutil.ReadFile("./aikblob")
	if err != nil {
		t.Skip("No aikblob file; skipping test")
	}

	if _, err := UnmarshalRSAPublicKey(blob); err != nil {
		t.Fatal("Couldn't extract an RSA key from the AIK blob:", err)
	}
}

func TestMakeIdentity(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	srkAuth := getAuth(srkAuthEnvVar)
	ownerAuth := getAuth(ownerAuthEnvVar)
	aikAuth := getAuth(aikAuthEnvVar)

	// In the simplest case, we pass in nil for the Privacy CA key and the
	// label.
	blob, err := MakeIdentity(rwc, srkAuth[:], ownerAuth[:], aikAuth[:], nil, nil)
	if err != nil {
		t.Fatal("Couldn't make a new AIK in the TPM:", err)
	}

	handle, err := LoadKey2(rwc, blob, srkAuth[:])
	if err != nil {
		t.Fatal("Couldn't load the freshly-generated AIK into the TPM and get a handle for it:", err)
	}
	defer CloseKey(rwc, handle)

	// Data to quote.
	data := []byte(`The OS says this test and new AIK is good`)
	pcrNums := []int{17, 18}
	q, values, err := Quote(rwc, handle, data, pcrNums, aikAuth[:])
	if err != nil {
		t.Fatal("Couldn't quote the data:", err)
	}

	// Verify the quote.
	pk, err := UnmarshalRSAPublicKey(blob)
	if err != nil {
		t.Fatal("Couldn't extract an RSA key from the AIK blob:", err)
	}

	if err := VerifyQuote(pk, data, q, pcrNums, values); err != nil {
		t.Fatal("The quote didn't pass verification:", err)
	}
}

func TestResetLockValue(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	// This test code assumes that the owner auth is the well-known value.
	ownerAuth := getAuth(ownerAuthEnvVar)
	if err := ResetLockValue(rwc, ownerAuth); err != nil {
		t.Fatal("Couldn't reset the lock value:", err)
	}
}

func TestOwnerReadSRK(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	// This test code assumes that the owner auth is the well-known value.
	ownerAuth := getAuth(ownerAuthEnvVar)
	srkb, err := OwnerReadSRK(rwc, ownerAuth)
	if err != nil {
		t.Fatal("Couldn't read the SRK using owner auth:", err)
	}

	if len(srkb) == 0 {
		t.Fatal("Couldn't get an SRK blob from the TPM")
	}
}

func TestOwnerReadPubEK(t *testing.T) {
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	// This test code assumes that the owner auth is the well-known value.
	ownerAuth := getAuth(ownerAuthEnvVar)
	pkb, err := OwnerReadPubEK(rwc, ownerAuth)
	if err != nil {
		t.Fatal("Couldn't read the pub EK using owner auth:", err)
	}

	pk, err := UnmarshalPubRSAPublicKey(pkb)
	if err != nil {
		t.Fatal("Couldn't unmarshal the endorsement key:", err)
	}

	if pk.N.BitLen() != 2048 {
		t.Fatal("Invalid endorsement key: not a 2048-bit RSA key")
	}
}

func TestOwnerClear(t *testing.T) {
	// Only enable this if you know what you're doing.
	t.Skip()
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	// This test code assumes that the owner auth is the well-known value.
	ownerAuth := getAuth(ownerAuthEnvVar)
	if err := OwnerClear(rwc, ownerAuth); err != nil {
		t.Fatal("Couldn't clear the TPM using owner auth:", err)
	}
}

func TestTakeOwnership(t *testing.T) {
	// This only works in limited circumstances, so it's disabled in general.
	t.Skip()
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	ownerAuth := getAuth(ownerAuthEnvVar)
	srkAuth := getAuth(srkAuthEnvVar)

	// This test assumes that the TPM has been cleared using OwnerClear.
	pubEK, err := ReadPubEK(rwc)
	if err != nil {
		t.Fatal("Couldn't read the public endorsement key from the TPM:", err)
	}

	if err := TakeOwnership(rwc, ownerAuth, srkAuth, pubEK); err != nil {
		t.Fatal("Couldn't take ownership of the TPM:", err)
	}
}

func TestForceClear(t *testing.T) {
	// Only enable this if you know what you're doing.
	// TPM force clear clears the ownership of the TPM.
	// Beware of running this test on a production system.
	t.Skip()
	rwc := openTPMOrSkip(t)
	defer rwc.Close()

	if err := ForceClear(rwc); err != nil {
		t.Fatal("Couldn't clear the TPM without owner auth in physical presence mode:", err)
	}
}
