package server

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/google/go-tpm-tools/internal"
	"github.com/google/go-tpm-tools/proto"
	"github.com/google/go-tpm-tools/tpm2tools"
	"github.com/google/go-tpm/tpm2"
)

func TestImport(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer tpm2tools.CheckedClose(t, rwc)
	tests := []struct {
		name     string
		template tpm2.Public
	}{
		{"RSA", tpm2tools.DefaultEKTemplateRSA()},
		{"ECC", tpm2tools.DefaultEKTemplateECC()},
		{"SRK-RSA", tpm2tools.SRKTemplateRSA()},
		{"SRK-ECC", tpm2tools.SRKTemplateECC()},
		{"ECC-P224", getECCTemplate(tpm2.CurveNISTP224)},
		{"ECC-P256", getECCTemplate(tpm2.CurveNISTP256)},
		{"ECC-P384", getECCTemplate(tpm2.CurveNISTP384)},
		{"ECC-P521", getECCTemplate(tpm2.CurveNISTP521)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ek, err := tpm2tools.NewKey(rwc, tpm2.HandleEndorsement, test.template)
			if err != nil {
				t.Fatal(err)
			}
			defer ek.Close()
			pub := ek.PublicKey()
			secret := []byte("super secret code")
			blob, err := CreateImportBlob(pub, secret, nil)
			if err != nil {
				t.Fatalf("creating import blob failed: %v", err)
			}

			output, err := ek.Import(blob)
			if err != nil {
				t.Fatalf("import failed: %v", err)
			}
			if !bytes.Equal(output, secret) {
				t.Errorf("got %X, expected %X", output, secret)
			}
		})
	}
}

func TestImportPCRs(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer tpm2tools.CheckedClose(t, rwc)

	ek, err := tpm2tools.EndorsementKeyRSA(rwc)
	if err != nil {
		t.Fatal(err)
	}
	defer ek.Close()
	pcr0, err := tpm2.ReadPCR(rwc, 0, tpm2.AlgSHA256)
	if err != nil {
		t.Fatal(err)
	}
	badPCR := append([]byte(nil), pcr0...)
	// badPCR increments first value so it doesn't match.
	badPCR[0]++
	tests := []struct {
		name          string
		pcrs          *proto.Pcrs
		expectSuccess bool
	}{
		{"No-PCR-nil", nil, true},
		{"No-PCR-empty", &proto.Pcrs{Hash: proto.HashAlgo_SHA256}, true},
		{"Good-PCR", &proto.Pcrs{Hash: proto.HashAlgo_SHA256, Pcrs: map[uint32][]byte{0: pcr0}}, true},
		{"Bad-PCR", &proto.Pcrs{Hash: proto.HashAlgo_SHA256, Pcrs: map[uint32][]byte{0: badPCR}}, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			secret := []byte("super secret code")
			blob, err := CreateImportBlob(ek.PublicKey(), secret, test.pcrs)
			if err != nil {
				t.Fatalf("creating import blob failed: %v", err)
			}
			output, err := ek.Import(blob)
			if test.expectSuccess {
				if err != nil {
					t.Fatalf("import failed: %v", err)
				}
				if !bytes.Equal(output, secret) {
					t.Errorf("got %X, expected %X", output, secret)
				}
			} else if err == nil {
				t.Error("expected Import to fail but it did not")
			}
		})
	}
}

func TestSigningKeyImport(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer tpm2tools.CheckedClose(t, rwc)

	ek, err := tpm2tools.EndorsementKeyRSA(rwc)
	if err != nil {
		t.Fatal(err)
	}
	defer ek.Close()
	signingKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	pcr0, err := tpm2.ReadPCR(rwc, 0, tpm2.AlgSHA256)
	if err != nil {
		t.Fatal(err)
	}
	badPCR := append(make([]byte, 0), pcr0...)
	// badPCR increments first value so it doesn't match.
	badPCR[0]++
	tests := []struct {
		name          string
		pcrs          *proto.Pcrs
		expectSuccess bool
	}{
		{"No-PCR-nil", nil, true},
		{"No-PCR-empty", &proto.Pcrs{Hash: proto.HashAlgo_SHA256}, true},
		{"Good-PCR", &proto.Pcrs{Hash: proto.HashAlgo_SHA256, Pcrs: map[uint32][]byte{0: pcr0}}, true},
		{"Bad-PCR", &proto.Pcrs{Hash: proto.HashAlgo_SHA256, Pcrs: map[uint32][]byte{0: badPCR}}, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			blob, err := CreateSigningKeyImportBlob(ek.PublicKey(), signingKey, test.pcrs)
			if err != nil {
				t.Fatalf("creating import blob failed: %v", err)
			}

			importedKey, err := ek.ImportSigningKey(blob)
			if err != nil {
				t.Fatalf("import failed: %v", err)
			}
			defer importedKey.Close()
			signer, err := importedKey.GetSigner()
			if err != nil {
				t.Fatalf("could not create signer: %v", err)
			}
			var digest [32]byte

			sig, err := signer.Sign(nil, digest[:], crypto.SHA256)
			if test.expectSuccess {
				if err != nil {
					t.Fatalf("import failed: %v", err)
				}
				if err = rsa.VerifyPKCS1v15(&signingKey.PublicKey, crypto.SHA256, digest[:], sig); err != nil {
					t.Error(err)
				}
				return
			} else if err == nil {
				t.Error("expected Import to fail but it did not")
			}
		})
	}
}
