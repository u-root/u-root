package tpm2tools

import (
	"bytes"
	"crypto/sha256"
	"io"
	"testing"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"

	"github.com/google/go-tpm-tools/internal"
)

func TestSeal(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer CheckedClose(t, rwc)

	tests := []struct {
		name   string
		getSRK func(io.ReadWriter) (*Key, error)
	}{
		{"RSA", StorageRootKeyRSA},
		{"ECC", StorageRootKeyECC},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srk, err := test.getSRK(rwc)
			if err != nil {
				t.Fatalf("can't create %s srk from template: %v", test.name, err)
			}
			defer srk.Close()

			secret := []byte("test")
			sel := tpm2.PCRSelection{Hash: tpm2.AlgSHA256, PCRs: []int{7, 23}}
			pcrToExtend := tpmutil.Handle(23)
			sealed, err := srk.Seal(secret, SealCurrent{PCRSelection: sel})
			if err != nil {
				t.Fatalf("failed to seal: %v", err)
			}

			cOpt := CertifyCurrent{
				PCRSelection: tpm2.PCRSelection{
					Hash: tpm2.AlgSHA256,
					PCRs: []int{7},
				},
			}
			unseal, err := srk.Unseal(sealed, cOpt)
			if err != nil {
				t.Fatalf("failed to unseal: %v", err)
			}
			if !bytes.Equal(secret, unseal) {
				t.Fatalf("unsealed (%v) not equal to secret (%v)", unseal, secret)
			}

			extension := bytes.Repeat([]byte{0xAA}, sha256.Size)
			if err = tpm2.PCRExtend(rwc, pcrToExtend, tpm2.AlgSHA256, extension, ""); err != nil {
				t.Fatalf("failed to extend pcr: %v", err)
			}

			// unseal should not succeed.
			if _, err = srk.Unseal(sealed, cOpt); err == nil {
				t.Fatalf("unseal should have caused an error: %v", err)
			}
		})
	}
}

func TestComputeSessionAuth(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer CheckedClose(t, rwc)

	pcrNums := []int{1, 7}

	tests := []struct {
		name    string
		pcrHash tpm2.Algorithm
	}{
		{"sha1", tpm2.AlgSHA1},
		{"sha256", tpm2.AlgSHA256},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sel := tpm2.PCRSelection{Hash: test.pcrHash, PCRs: pcrNums}
			pcrs, err := ReadPCRs(rwc, sel)
			if err != nil {
				t.Fatalf("failed to read PCRs: %v", err)
			}
			computeDigest := ComputePCRSessionAuth(pcrs)

			session, err := newPCRSession(rwc, sel)
			if err != nil {
				t.Fatalf("failed to create PCR session: %v", err)
			}
			defer session.Close()

			auth, err := session.Auth()
			if err != nil {
				t.Fatalf("failed to get PCR auth: %v", err)
			}
			digest, err := tpm2.PolicyGetDigest(rwc, auth.Session)
			if err != nil {
				t.Fatalf("failed to get pcr digest: %v", err)
			}

			if !bytes.Equal(computeDigest, digest) {
				t.Errorf("computed digest (%v) not equal to session digest (%v)", computeDigest, digest)
			}
		})
	}
}

func TestSelfReseal(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer CheckedClose(t, rwc)

	key, err := StorageRootKeyRSA(rwc)
	if err != nil {
		t.Fatalf("can't create srk from template: %v", err)
	}
	defer key.Close()

	secret := []byte("test")
	pcrList := []int{0, 4, 7}
	sOpt := SealCurrent{
		PCRSelection: tpm2.PCRSelection{
			Hash: tpm2.AlgSHA256,
			PCRs: pcrList,
		},
	}

	sealed, err := key.Seal(secret, sOpt)
	if err != nil {
		t.Fatalf("failed to seal: %v", err)
	}

	cOpt := CertifyCurrent{
		PCRSelection: tpm2.PCRSelection{
			Hash: tpm2.AlgSHA256,
			PCRs: []int{7},
		},
	}
	unseal, err := key.Unseal(sealed, cOpt)
	if err != nil {
		t.Fatalf("failed to unseal: %v", err)
	}
	if !bytes.Equal(secret, unseal) {
		t.Errorf("unsealed (%v) not equal to secret (%v)", unseal, secret)
	}

	sealed, err = key.Reseal(sealed, cOpt, sOpt)
	if err != nil {
		t.Fatalf("failed to reseal: %v", err)
	}

	unseal, err = key.Unseal(sealed, cOpt)
	if err != nil {
		t.Fatalf("failed to unseal after resealing: %v", err)
	}
	if !bytes.Equal(secret, unseal) {
		t.Errorf("unsealed (%v) not equal to secret (%v)", unseal, secret)
	}
}

func computePCRValue(base []byte, extensions [][]byte) []byte {
	for _, extension := range extensions {
		sum := sha256.Sum256(append(base, extension...))
		base = sum[:]
	}
	return base
}

func TestComputePCRValue(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer CheckedClose(t, rwc)

	pcrNum := 23
	extensions := [][]byte{
		bytes.Repeat([]byte{0xAA}, sha256.Size),
		bytes.Repeat([]byte{0xAB}, sha256.Size),
		bytes.Repeat([]byte{0xAC}, sha256.Size),
		bytes.Repeat([]byte{0xAD}, sha256.Size),
	}

	pcrBase, err := tpm2.ReadPCR(rwc, pcrNum, tpm2.AlgSHA256)
	if err != nil {
		t.Fatalf("failed to read pcr %v", err)
	}

	for _, extension := range extensions {
		err := tpm2.PCRExtend(rwc, tpmutil.Handle(pcrNum), tpm2.AlgSHA256, extension, "")
		if err != nil {
			t.Fatalf("failed to extend pcr: %v", err)
		}
	}

	pcrVal, err := tpm2.ReadPCR(rwc, pcrNum, tpm2.AlgSHA256)
	if err != nil {
		t.Fatalf("failed to read pcr %v", err)
	}

	computedValue := computePCRValue(pcrBase, extensions)
	if !bytes.Equal(pcrVal, computedValue) {
		t.Fatalf("pcrVal (%v) not equal to computedValue (%v)", pcrVal, computedValue)
	}
}

func TestReseal(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer CheckedClose(t, rwc)

	key, err := StorageRootKeyRSA(rwc)
	if err != nil {
		t.Fatalf("can't create srk from template: %v", err)
	}
	defer key.Close()

	secret := []byte("test")
	pcrToChange := uint32(23)
	sel := tpm2.PCRSelection{Hash: tpm2.AlgSHA256, PCRs: []int{7, 23}}
	sealed, err := key.Seal(secret, SealCurrent{PCRSelection: sel})
	if err != nil {
		t.Fatalf("failed to seal: %v", err)
	}

	cOpt := CertifyCurrent{
		PCRSelection: tpm2.PCRSelection{
			Hash: tpm2.AlgSHA256,
			PCRs: []int{7, 23},
		},
	}
	unseal, err := key.Unseal(sealed, cOpt)
	if err != nil {
		t.Fatalf("failed to unseal: %v", err)
	}
	if !bytes.Equal(secret, unseal) {
		t.Fatalf("unsealed (%v) not equal to secret (%v)", unseal, secret)
	}

	// create a new set of PCRs value for modificiation
	predictedPcrsValue, err := ReadPCRs(rwc, sel)
	if err != nil {
		t.Fatalf("failed to read PCRs value: %v", err)
	}
	// change pcr value to the predicted future value for resealing
	extensions := [][]byte{bytes.Repeat([]byte{0xAA}, sha256.Size)}
	predictedPcrsValue.GetPcrs()[uint32(pcrToChange)] = computePCRValue(predictedPcrsValue.GetPcrs()[uint32(pcrToChange)], extensions)

	resealed, err := key.Reseal(sealed, cOpt, SealTarget{predictedPcrsValue})
	if err != nil {
		t.Fatalf("failed to reseal: %v", err)
	}

	// unseal should not succeed since pcr has not been extended.
	if _, err = key.Unseal(resealed, nil); err == nil {
		t.Fatalf("unseal should have failed: %v", err)
	}

	// save the current PCR value for certification before extend the PCRs
	oldPcrsValue, err := ReadPCRs(rwc, sel)
	if err != nil {
		t.Fatalf("failed to read PCRs value: %v", err)
	}
	for _, extension := range extensions {
		err = tpm2.PCRExtend(rwc, tpmutil.Handle(pcrToChange), tpm2.AlgSHA256, extension, "")
		if err != nil {
			t.Fatalf("failed to extend pcr: %v", err)
		}
	}

	// unseal should fail if certify to current PCRs value, as one PCR has changed
	_, err = key.Unseal(resealed, CertifyCurrent{PCRSelection: sel})
	if err == nil {
		t.Fatalf("unseal should fail since the certify PCRs have changed.")
	}

	// certify to original PCRs value (PCRs value when do the sealing) will work
	unseal, err = key.Unseal(resealed, CertifyExpected{oldPcrsValue})
	if err != nil {
		t.Fatalf("failed to unseal: %v", err)
	}
	if !bytes.Equal(secret, unseal) {
		t.Errorf("unsealed (%v) not equal to secret (%v)", unseal, secret)
	}
}

func TestSealResealWithEmptyPCRs(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer CheckedClose(t, rwc)

	key, err := StorageRootKeyRSA(rwc)
	if err != nil {
		t.Fatalf("can't create srk from template: %v", err)
	}
	defer key.Close()

	secret := []byte("test")
	sealed, err := key.Seal(secret, nil)
	if err != nil {
		t.Fatalf("failed to seal: %v", err)
	}
	cOpt := CertifyCurrent{
		PCRSelection: tpm2.PCRSelection{
			Hash: tpm2.AlgSHA256,
			PCRs: []int{7},
		},
	}
	unseal, err := key.Unseal(sealed, cOpt)
	if err != nil {
		t.Fatalf("failed to unseal: %v", err)
	}
	if !bytes.Equal(secret, unseal) {
		t.Fatalf("unsealed (%v) not equal to secret (%v)", unseal, secret)
	}

	extension := bytes.Repeat([]byte{0xAA}, sha256.Size)
	if err = tpm2.PCRExtend(rwc, 7, tpm2.AlgSHA256, extension, ""); err != nil {
		t.Fatalf("failed to extend pcr: %v", err)
	}

	// unseal should failed as the PCR 7 has changed (not as same as when sealing)
	_, err = key.Unseal(sealed, cOpt)
	if err == nil {
		t.Fatalf("unseal should fail as PCR 7 changed")
	}

	// reseal should success as CertifyOpt is nil
	sealed, err = key.Reseal(sealed, nil, nil)
	if err != nil {
		t.Fatalf("failed to reseal: %v", err)
	}

	// unseal should success as the above Reseal() "refresh" the Ceritfy PCRs.
	unseal, err = key.Unseal(sealed, cOpt)
	if err != nil {
		t.Errorf("failed to unseal: %v", err)
	}
	if !bytes.Equal(secret, unseal) {
		t.Fatalf("unsealed (%v) not equal to secret (%v)", unseal, secret)
	}
}

func BenchmarkSeal(b *testing.B) {
	rwc := internal.GetTPM(b)
	defer CheckedClose(b, rwc)

	pcrSel7 := tpm2.PCRSelection{Hash: tpm2.AlgSHA256, PCRs: []int{7}}
	sOptPCR7 := SealCurrent{PCRSelection: pcrSel7}
	cOptPCR7 := CertifyCurrent{PCRSelection: pcrSel7}
	benchmarks := []struct {
		name   string
		sOpt   SealOpt
		cOpt   CertifyOpt
		getKey func(io.ReadWriter) (*Key, error)
	}{
		{"SRK-ECC-SealPCR7-CertifyPCR7", sOptPCR7, cOptPCR7, StorageRootKeyECC},
		{"SRK-ECC-nil-CertifyPCR7", nil, cOptPCR7, StorageRootKeyECC},
		{"SRK-ECC-SealPCR7-nil", sOptPCR7, nil, StorageRootKeyECC},
		{"SRK-ECC-nil-nil", nil, nil, StorageRootKeyECC},
		{"SRK-RSA-SealPCR7-CertifyPCR7", sOptPCR7, cOptPCR7, StorageRootKeyRSA},
		{"SRK-RSA-nil-CertifyPCR7", nil, cOptPCR7, StorageRootKeyRSA},
		{"SRK-RSA-SealPCR7-nil", sOptPCR7, nil, StorageRootKeyRSA},
		{"SRK-RSA-nil-nil", nil, nil, StorageRootKeyRSA},
	}

	for _, bm := range benchmarks {
		key, err := bm.getKey(rwc)
		if err != nil {
			b.Fatal(err)
		}
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				blob, err := key.Seal([]byte("test123"), bm.sOpt)
				if err != nil {
					b.Fatal(err)
				}
				if _, err = key.Unseal(blob, bm.cOpt); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
