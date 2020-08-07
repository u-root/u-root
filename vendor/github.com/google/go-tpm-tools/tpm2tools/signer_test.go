package tpm2tools

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/asn1"
	"math/big"
	"testing"

	"github.com/google/go-tpm-tools/internal"
	"github.com/google/go-tpm/tpm2"
)

func templateSSA(hash tpm2.Algorithm) tpm2.Public {
	template := AIKTemplateRSA()
	// Can't sign arbitrary data if restricted.
	template.Attributes &= ^tpm2.FlagRestricted
	template.RSAParameters.Sign.Hash = hash
	return template
}

func templatePSS(hash tpm2.Algorithm) tpm2.Public {
	template := templateSSA(hash)
	template.RSAParameters.Sign.Alg = tpm2.AlgRSAPSS
	return template
}

func templateECC(hash tpm2.Algorithm) tpm2.Public {
	template := AIKTemplateECC()
	template.Attributes &= ^tpm2.FlagRestricted
	template.ECCParameters.Sign.Hash = hash
	return template
}

// Templates that require some sort of (default) authorization
func templateAuthSSA() tpm2.Public {
	template := templateSSA(tpm2.AlgSHA256)
	template.AuthPolicy = defaultEKAuthPolicy()
	template.Attributes |= tpm2.FlagAdminWithPolicy
	template.Attributes &= ^tpm2.FlagUserWithAuth
	return template
}

func templateAuthECC() tpm2.Public {
	template := templateECC(tpm2.AlgSHA256)
	template.AuthPolicy = defaultEKAuthPolicy()
	template.Attributes |= tpm2.FlagAdminWithPolicy
	template.Attributes &= ^tpm2.FlagUserWithAuth
	return template
}

func verifyRSA(pubKey crypto.PublicKey, hash crypto.Hash, digest, sig []byte) bool {
	return rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), hash, digest, sig) == nil
}

func verifyECC(pubKey crypto.PublicKey, _ crypto.Hash, digest, sig []byte) bool {
	var sigStruct struct{ R, S *big.Int }
	asn1.Unmarshal(sig, &sigStruct)
	return ecdsa.Verify(pubKey.(*ecdsa.PublicKey), digest, sigStruct.R, sigStruct.S)
}

func TestSign(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer CheckedClose(t, rwc)

	tests := []struct {
		name     string
		hash     crypto.Hash
		template tpm2.Public
		verify   func(crypto.PublicKey, crypto.Hash, []byte, []byte) bool
	}{
		{"RSA-SHA1", crypto.SHA1, templateSSA(tpm2.AlgSHA1), verifyRSA},
		{"RSA-SHA256", crypto.SHA256, templateSSA(tpm2.AlgSHA256), verifyRSA},
		{"RSA-SHA384", crypto.SHA384, templateSSA(tpm2.AlgSHA384), verifyRSA},
		{"RSA-SHA512", crypto.SHA512, templateSSA(tpm2.AlgSHA512), verifyRSA},
		{"ECC-SHA1", crypto.SHA1, templateECC(tpm2.AlgSHA1), verifyECC},
		{"ECC-SHA256", crypto.SHA256, templateECC(tpm2.AlgSHA256), verifyECC},
		{"ECC-SHA384", crypto.SHA384, templateECC(tpm2.AlgSHA384), verifyECC},
		{"ECC-SHA512", crypto.SHA512, templateECC(tpm2.AlgSHA512), verifyECC},
		{"Auth-RSA", crypto.SHA256, templateAuthSSA(), verifyRSA},
		{"Auth-ECC", crypto.SHA256, templateAuthECC(), verifyECC},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			key, err := NewKey(rwc, tpm2.HandleEndorsement, test.template)
			if err != nil {
				t.Fatal(err)
			}
			defer key.Close()

			hash := test.hash.New()
			hash.Write([]byte("authenticated message"))
			digest := hash.Sum(nil)

			signer, err := key.GetSigner()
			if err != nil {
				t.Fatal(err)
			}

			sig, err := signer.Sign(nil, digest, test.hash)
			if err != nil {
				t.Fatal(err)
			}
			if !test.verify(signer.Public(), test.hash, digest, sig) {
				t.Error(err)
			}
		})
	}
}

func TestSignIncorrectHash(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer CheckedClose(t, rwc)

	key, err := NewKey(rwc, tpm2.HandleEndorsement, templateSSA(tpm2.AlgSHA256))
	if err != nil {
		t.Fatal(err)
	}
	defer key.Close()

	signer, err := key.GetSigner()
	if err != nil {
		t.Fatal(err)
	}

	digestSHA1 := sha1.Sum([]byte("authenticated message"))
	digestSHA256 := sha256.Sum256([]byte("authenticated message"))

	if _, err := signer.Sign(nil, digestSHA1[:], crypto.SHA1); err == nil {
		t.Error("expected failure for digest and hash not matching keys sigScheme.")
	}

	if _, err := signer.Sign(nil, digestSHA1[:], crypto.SHA256); err == nil {
		t.Error("expected failure for correct hash, but incorrect digest.")
	}

	if _, err := signer.Sign(nil, digestSHA256[:], crypto.SHA1); err == nil {
		t.Error("expected failure for correct digest, but incorrect hash.")
	}
}

func TestSignPSS(t *testing.T) {
	rwc := internal.GetTPM(t)
	defer CheckedClose(t, rwc)
	tests := []struct {
		name     string
		opts     crypto.SignerOpts
		template tpm2.Public
		keyBits  uint16
		saltLen  int
	}{
		// saltLen should be (keyBits/8) - digestSize - 2, unless that is less than
		// digestSize in which case, saltLen will be digestSize.
		// The only normal case where saltLen is not digestSize is when using
		// 1024 keyBits with SHA512.
		{"RSA-SHA1", crypto.SHA1, templatePSS(tpm2.AlgSHA1), 1024, 20},
		{"RSA-SHA256", crypto.SHA256, templatePSS(tpm2.AlgSHA256), 1024, 32},
		{"RSA-SHA384", crypto.SHA384, templatePSS(tpm2.AlgSHA384), 1024, 48},
		{"RSA-SHA512", crypto.SHA512, templatePSS(tpm2.AlgSHA512), 1024, 62},
		{"RSA-SHA512", crypto.SHA512, templatePSS(tpm2.AlgSHA512), 2048, 64},
		{"RSA-SHA1", &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA1}, templatePSS(tpm2.AlgSHA1), 1024, 20},
		{"RSA-SHA256", &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA256}, templatePSS(tpm2.AlgSHA256), 1024, 32},
		{"RSA-SHA384", &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA384}, templatePSS(tpm2.AlgSHA384), 1024, 48},
		{"RSA-SHA512", &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA512}, templatePSS(tpm2.AlgSHA512), 1024, 62},
		{"RSA-SHA512", &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA512}, templatePSS(tpm2.AlgSHA512), 2048, 64},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.template.RSAParameters.KeyBits = test.keyBits

			key, err := NewKey(rwc, tpm2.HandleEndorsement, test.template)
			if err != nil {
				t.Fatal(err)
			}
			defer key.Close()

			hash := test.opts.HashFunc().New()
			hash.Write([]byte("authenticated message"))
			digest := hash.Sum(nil)

			signer, err := key.GetSigner()
			if err != nil {
				t.Fatal(err)
			}
			sig, err := signer.Sign(nil, digest[:], test.opts)
			if err != nil {
				t.Fatal(err)
			}
			// Verify with expected salt length.
			err = rsa.VerifyPSS(signer.Public().(*rsa.PublicKey), test.opts.HashFunc(), digest[:], sig, &rsa.PSSOptions{SaltLength: test.saltLen, Hash: test.opts.HashFunc()})
			if err != nil {
				t.Error(err)
			}

			// Verify with default salt length.
			err = rsa.VerifyPSS(signer.Public().(*rsa.PublicKey), test.opts.HashFunc(), digest[:], sig, nil)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
