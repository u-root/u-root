// Copyright 2017-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crypto

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path"
	"testing"

	"golang.org/x/crypto/ed25519"
)

const (
	// publicKeyDERFile is a RSA public key in DER format
	publicKeyDERFile string = "tests/public_key.der"
	// publicKeyPEMFile is a RSA public key in PEM format
	publicKeyPEMFile string = "tests/public_key.pem"
	// privateKeyPEMFile is a RSA public key in PEM format
	privateKeyPEMFile string = "tests/private_key.pem"
	// testDataFile which should be verified by the good signature
	testDataFile string = "tests/data"
	// signatureGoodFile is a good signature of testDataFile
	signatureGoodFile string = "tests/verify_rsa_pkcs15_sha256.signature"
	// signatureBadFile is a bad signature which does not work with testDataFile
	signatureBadFile string = "tests/verify_rsa_pkcs15_sha256.signature2"
)

// password is a PEM encrypted passphrase
var password = []byte{'k', 'e', 'i', 'n', 's'}

func TestLoadDERPublicKey(t *testing.T) {
	if _, err := LoadPublicKeyFromFile(publicKeyDERFile); err == nil {
		t.Errorf(`LoadPublicKeyFromFile(publicKeyDERFile) = _, %v, want not nil`, err)
	}
}

func TestLoadPEMPublicKey(t *testing.T) {
	if _, err := LoadPublicKeyFromFile(publicKeyPEMFile); err != nil {
		t.Errorf(`LoadPublicKeyFromFile(publicKeyPEMFile) = _, %v, want nil`, err)
	}
}

func TestLoadPEMPrivateKey(t *testing.T) {
	if _, err := LoadPrivateKeyFromFile(privateKeyPEMFile, password); err != nil {
		t.Errorf(`LoadPublicKeyFromFile(privateKeyPEMFile) = _, %v, want nil`, err)
	}
}

func TestLoadBadPEMPrivateKey(t *testing.T) {
	if _, err := LoadPrivateKeyFromFile(privateKeyPEMFile, []byte{}); err == nil {
		t.Errorf(`LoadPrivateKeyFromFile(privateKeyPEMFile, []byte{}) = _, %v, want not nil`, err)
	}
}

func TestSignVerifyData(t *testing.T) {
	privateKey, err := LoadPrivateKeyFromFile(privateKeyPEMFile, password)
	if err != nil {
		t.Errorf(`LoadPrivateKeyFromFile(privateKeyPEMFile, password) = _, %v, want nil`, err)
	}

	publicKey, err := LoadPublicKeyFromFile(publicKeyPEMFile)
	if err != nil {
		t.Errorf(`LoadPublicKeyFromFile(publicKeyPEMFile) = _, %v, want nil`, err)
	}

	testData, err := os.ReadFile(testDataFile)
	if err != nil {
		t.Errorf(`os.ReadFile(testDataFile) = _, %v, want nil`, err)
	}

	signature := ed25519.Sign(privateKey, testData)
	if verified := ed25519.Verify(publicKey, testData, signature); !verified {
		t.Errorf(`ed25519.Verify(publicKey, testData, signature) = %t, want "true"`, verified)
	}
}

func TestGoodSignature(t *testing.T) {
	publicKey, err := LoadPublicKeyFromFile(publicKeyPEMFile)
	if err != nil {
		t.Errorf(`LoadPublicKeyFromFile(publicKeyPEMFile) = _, %v, want nil`, err)
	}

	testData, err := os.ReadFile(testDataFile)
	if err != nil {
		t.Errorf(`os.ReadFile(testDataFile) = _, %v, want nil`, err)
	}

	signatureGood, err := os.ReadFile(signatureGoodFile)
	if err != nil {
		t.Errorf(`os.ReadFile(signatureGoodFile) = _, %v, want nil`, err)
	}

	if verified := ed25519.Verify(publicKey, testData, signatureGood); !verified {
		t.Errorf(`ed25519.Verify(publicKey, testData, signatureGood) = %t, want "true"`, verified)
	}
}

func TestBadSignature(t *testing.T) {
	publicKey, err := LoadPublicKeyFromFile(publicKeyPEMFile)
	if err != nil {
		t.Errorf(`LoadPublicKeyFromFile(publicKeyPEMFile) = _, %v, want nil`, err)
	}

	testData, err := os.ReadFile(testDataFile)
	if err != nil {
		t.Errorf(`os.ReadFile(testDataFile) = _, %v, want nil`, err)
	}

	signatureBad, err := os.ReadFile(signatureBadFile)
	if err != nil {
		t.Errorf(`os.ReadFile(signatureBadFile) = _, %v, want nil`, err)
	}

	if verified := ed25519.Verify(publicKey, testData, signatureBad); verified {
		t.Errorf(`ed25519.Verify(publicKey, testData, signatureBad) = %t, want "false"`, verified)
	}
}

func isEncryptedKeyFile(t *testing.T, privateKeyPath string) bool {
	t.Helper()
	b, err := os.ReadFile(privateKeyPath)
	if err != nil {
		t.Fatalf(`os.ReadFile(%q) = _, %v, want nil`, privateKeyPath, err)
	}
	block, _ := pem.Decode(b)
	if block == nil {
		t.Fatalf(`pem.Decode(%q) = nil, want a PEM block`, privateKeyPath)
	}
	return x509.IsEncryptedPEMBlock(block)
}

func TestGenerateKeys(t *testing.T) {
	tmpdir := t.TempDir()
	privateKeyPath := path.Join(tmpdir, "private_key.pem")
	publicKeyPath := path.Join(tmpdir, "public_key.pem")
	if err := GeneratED25519Key(password, privateKeyPath, publicKeyPath); err != nil {
		t.Fatalf(`GeneratED25519Key(password, privateKeyPath, publicKeyPath) = %v, want nil`, err)
	}

	if !isEncryptedKeyFile(t, privateKeyPath) {
		t.Errorf(`GeneratED25519Key(%q, ...) wrote an unencrypted private key, want an encrypted key`, password)
	}

	privateKey, err := LoadPrivateKeyFromFile(privateKeyPath, password)
	if err != nil {
		t.Fatalf(`LoadPrivateKeyFromFile(privateKeyPath, password) = _, %v, want nil`, err)
	}
	publicKey, err := LoadPublicKeyFromFile(publicKeyPath)
	if err != nil {
		t.Fatalf(`LoadPublicKeyFromFile(publicKeyPath) = _, %v, want nil`, err)
	}
	if got := ed25519.PrivateKey(privateKey).Public(); !bytes.Equal(got.(ed25519.PublicKey), publicKey) {
		t.Errorf(`LoadPrivateKeyFromFile(privateKeyPath, password).Public() = %x, want %x`, got, publicKey)
	}
}

func TestGenerateUnprotectedKeys(t *testing.T) {
	tmpdir := t.TempDir()
	privateKeyPath := path.Join(tmpdir, "private_key.pem")
	if err := GeneratED25519Key(nil, privateKeyPath, path.Join(tmpdir, "public_key.pem")); err != nil {
		t.Fatalf(`GeneratED25519Key(nil, privateKeyPath, path.Join(tmpdir, "public_key.pem")) = %v, want nil`, err)
	}

	if isEncryptedKeyFile(t, privateKeyPath) {
		t.Errorf(`GeneratED25519Key(nil, ...) wrote an encrypted private key, want an unencrypted key`)
	}
}
