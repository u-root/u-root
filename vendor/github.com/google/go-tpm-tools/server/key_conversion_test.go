package server

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/google/go-tpm-tools/internal"
	"github.com/google/go-tpm-tools/tpm2tools"
	"github.com/google/go-tpm/tpm2"
)

func getECCTemplate(curve tpm2.EllipticCurve) tpm2.Public {
	public := tpm2tools.DefaultEKTemplateECC()
	public.ECCParameters.CurveID = curve
	public.ECCParameters.Point.XRaw = nil
	public.ECCParameters.Point.YRaw = nil
	return public
}

func TestCreateEKPublicAreaFromKeyGeneratedKey(t *testing.T) {
	tests := []struct {
		name        string
		template    tpm2.Public
		generateKey func() (crypto.PublicKey, error)
	}{
		{"RSA", tpm2tools.DefaultEKTemplateRSA(), func() (crypto.PublicKey, error) {
			priv, err := rsa.GenerateKey(rand.Reader, 2048)
			return priv.Public(), err
		}},
		{"ECC", tpm2tools.DefaultEKTemplateECC(), func() (crypto.PublicKey, error) {
			priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			return priv.Public(), err
		}},
		{"ECC-P224", getECCTemplate(tpm2.CurveNISTP224), func() (crypto.PublicKey, error) {
			priv, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
			return priv.Public(), err
		}},
		{"ECC-P256", getECCTemplate(tpm2.CurveNISTP256), func() (crypto.PublicKey, error) {
			priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			return priv.Public(), err
		}},
		{"ECC-P384", getECCTemplate(tpm2.CurveNISTP384), func() (crypto.PublicKey, error) {
			priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
			return priv.Public(), err
		}},
		{"ECC-P521", getECCTemplate(tpm2.CurveNISTP521), func() (crypto.PublicKey, error) {
			priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
			return priv.Public(), err
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			key, err := test.generateKey()
			if err != nil {
				t.Fatal(err)
			}
			newArea, err := CreateEKPublicAreaFromKey(key)
			if err != nil {
				t.Fatalf("failed to create public area from public key: %v", err)
			}
			if !newArea.MatchesTemplate(test.template) {
				t.Errorf("public areas did not match. got: %+v want: %+v", newArea, test.template)
			}
		})
	}
}

func TestCreateEKPublicAreaFromKeyTPMKey(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer tpm2tools.CheckedClose(t, rwc)

	tests := []struct {
		name     string
		template tpm2.Public
	}{
		{"RSA", tpm2tools.DefaultEKTemplateRSA()},
		{"ECC", tpm2tools.DefaultEKTemplateECC()},
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
			newArea, err := CreateEKPublicAreaFromKey(ek.PublicKey())
			if err != nil {
				t.Fatalf("failed to create public area from public key: %v", err)
			}
			if matches, err := ek.Name().MatchesPublic(newArea); err != nil || !matches {
				t.Error("public areas did not match or match check failed.")
			}
		})
	}
}
