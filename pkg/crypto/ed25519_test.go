// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crypto

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ed25519"
)

const (
	// publicKeyDERFile is a RSA public key in DER format
	publicKeyDERFile string = "tests/public_key.der"
	// publicKeyPEMFile is a RSA public key in PEM format
	publicKeyPEMFile string = "tests/public_key.pem"
	// privateKeyPEMFile is a RSA public key in PEM format
	privateKeyPEMFile string = "tests/private_key.pem"
	// publicKeyPEMFile2 is a RSA public key in PEM format
	publicKeyPEMFile2 string = "tests/public_key2.pem"
	// privateKeyPEMFile2 is a RSA public key in PEM format
	privateKeyPEMFile2 string = "tests/private_key2.pem"
	// testDataFile which should be verified by the good signature
	testDataFile string = "tests/data"
	// signatureGoodFile is a good signature of testDataFile
	signatureGoodFile string = "tests/verify_rsa_pkcs15_sha256.signature"
	// signatureBadFile is a bad signature which does not work with testDataFile
	signatureBadFile string = "tests/verify_rsa_pkcs15_sha256.signature2"
)

var (
	// password is a PEM encrypted passphrase
	password = []byte{'k', 'e', 'i', 'n', 's'}
)

func TestLoadDERPublicKey(t *testing.T) {
	_, err := LoadPublicKeyFromFile(publicKeyDERFile)
	require.Error(t, err)
}

func TestLoadPEMPublicKey(t *testing.T) {
	_, err := LoadPublicKeyFromFile(publicKeyPEMFile)
	require.NoError(t, err)
}

func TestLoadPEMPrivateKey(t *testing.T) {
	_, err := LoadPrivateKeyFromFile(privateKeyPEMFile, password)
	require.NoError(t, err)
}

func TestLoadBadPEMPrivateKey(t *testing.T) {
	_, err := LoadPrivateKeyFromFile(privateKeyPEMFile, []byte{})
	require.Error(t, err)
}

func TestSignVerifyData(t *testing.T) {
	privateKey, err := LoadPrivateKeyFromFile(privateKeyPEMFile, password)
	require.NoError(t, err)

	publicKey, err := LoadPublicKeyFromFile(publicKeyPEMFile)
	require.NoError(t, err)

	testData, err := ioutil.ReadFile(testDataFile)
	require.NoError(t, err)

	signature := ed25519.Sign(privateKey, testData)
	verified := ed25519.Verify(publicKey, testData, signature)
	require.Equal(t, true, verified)
}

func TestGoodSignature(t *testing.T) {
	publicKey, err := LoadPublicKeyFromFile(publicKeyPEMFile)
	require.NoError(t, err)

	testData, err := ioutil.ReadFile(testDataFile)
	require.NoError(t, err)

	signatureGood, err := ioutil.ReadFile(signatureGoodFile)
	require.NoError(t, err)

	verified := ed25519.Verify(publicKey, testData, signatureGood)
	require.Equal(t, true, verified)
}

func TestBadSignature(t *testing.T) {
	publicKey, err := LoadPublicKeyFromFile(publicKeyPEMFile)
	require.NoError(t, err)

	testData, err := ioutil.ReadFile(testDataFile)
	require.NoError(t, err)

	signatureBad, err := ioutil.ReadFile(signatureBadFile)
	require.NoError(t, err)

	verified := ed25519.Verify(publicKey, testData, signatureBad)
	require.Equal(t, false, verified)
}

func TestGenerateKeys(t *testing.T) {
	err := GeneratED25519Key(password, privateKeyPEMFile2, publicKeyPEMFile2)
	require.NoError(t, err)
}

func TestGenerateUnprotectedKeys(t *testing.T) {
	err := GeneratED25519Key(nil, privateKeyPEMFile2, publicKeyPEMFile2)
	require.NoError(t, err)
}
