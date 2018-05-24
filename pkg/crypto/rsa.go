package crypto

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

const (
	pubKeyIdentifier string = "PUBLIC KEY"
)

// LoadPublicKeyFromFile loads DER formatted RSA public key from file.
func LoadPublicKeyFromFile(publicKeyPath string) (*rsa.PublicKey, error) {
	x509PEM, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	// Parse x509 PEM file
	block, _ := pem.Decode(x509PEM)
	if block == nil || block.Type != pubKeyIdentifier {
		return nil, errors.New("Can't decode PEM file")
	}

	// parse Public Key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, err
	}
}

// VerifyRsaSha256Pkcs1v15Signature verifies a PKCSv1.5 signature made by
// a SHA-256 checksum. Public key must be a RSA key in PEM format.
func VerifyRsaSha256Pkcs1v15Signature(publicKey *rsa.PublicKey, data []byte, signature []byte) error {
	if publicKey == nil {
		return errors.New("Public key is nil")
	}

	hash := sha256.Sum256(data)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
}
