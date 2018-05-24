package crypto

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	publicKeyDERFile  string = "tests/public_key.der"
	publicKeyPEMFile  string = "tests/public_key.pem"
	testDataFile      string = "tests/data"
	signatureGoodFile string = "tests/verify_rsa_pkcs15_sha256.signature"
	signatureBadFile  string = "tests/verify_rsa_pkcs15_sha256.signature2"
)

func TestLoadDERPublicKey(t *testing.T) {
	_, err := LoadPublicKeyFromFile(publicKeyDERFile)
	require.Error(t, err)
}

func TestLoadPEMPublicKey(t *testing.T) {
	_, err := LoadPublicKeyFromFile(publicKeyPEMFile)
	require.NoError(t, err)
}

func TestGoodSignature(t *testing.T) {
	publicKey, err := LoadPublicKeyFromFile(publicKeyPEMFile)
	require.NoError(t, err)

	testData, err := ioutil.ReadFile(testDataFile)
	require.NoError(t, err)

	signatureGood, err := ioutil.ReadFile(signatureGoodFile)
	require.NoError(t, err)

	err = VerifyRsaSha256Pkcs1v15Signature(publicKey, testData, signatureGood)
	require.NoError(t, err)
}

func TestBadSignature(t *testing.T) {
	publicKey, err := LoadPublicKeyFromFile(publicKeyPEMFile)
	require.NoError(t, err)

	testData, err := ioutil.ReadFile(testDataFile)
	require.NoError(t, err)

	signatureBad, err := ioutil.ReadFile(signatureBadFile)
	require.NoError(t, err)

	err = VerifyRsaSha256Pkcs1v15Signature(publicKey, testData, signatureBad)
	require.Error(t, err)
}
