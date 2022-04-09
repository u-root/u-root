// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crypto

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"golang.org/x/crypto/ed25519"
)

var (
	// PubKeyIdentifier is the PEM public key identifier
	PubKeyIdentifier = "PUBLIC KEY"
	// PrivKeyIdentifier is the PEM private key identifier
	PrivKeyIdentifier = "PRIVATE KEY"
	// PEMCipher is the PEM encryption algorithm
	PEMCipher = x509.PEMCipherAES256
	// PubKeyFilePermissions are the public key file perms
	PubKeyFilePermissions os.FileMode = 0o644
	// PrivKeyFilePermissions are the private key file perms
	PrivKeyFilePermissions os.FileMode = 0o600
)

// LoadPublicKeyFromFile loads PEM formatted ED25519 public key from file.
func LoadPublicKeyFromFile(publicKeyPath string) ([]byte, error) {
	x509PEM, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	// Parse x509 PEM file
	var block *pem.Block
	for {
		block, x509PEM = pem.Decode(x509PEM)
		if block == nil {
			return nil, errors.New("can't decode PEM file")
		}
		if block.Type == PubKeyIdentifier {
			break
		}
	}

	return block.Bytes, nil
}

// LoadPrivateKeyFromFile loads PEM formatted ED25519 private key from file.
func LoadPrivateKeyFromFile(privateKeyPath string, password []byte) ([]byte, error) {
	x509PEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	// Parse x509 PEM file
	var block *pem.Block
	for {
		block, x509PEM = pem.Decode(x509PEM)
		if block == nil {
			return nil, errors.New("can't decode PEM file")
		}
		if block.Type == PrivKeyIdentifier {
			break
		}
	}

	// Check for encrypted PEM format
	if x509.IsEncryptedPEMBlock(block) {
		decryptedKey, err := x509.DecryptPEMBlock(block, password)
		if err != nil {
			return nil, err
		}
		return decryptedKey, nil
	}

	return block.Bytes, nil
}

// GeneratED25519Key generates a ED25519 keypair
func GeneratED25519Key(password []byte, privateKeyFilePath string, publicKeyFilePath string) error {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	privBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privKey,
	}

	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKey,
	}

	var privateKey []byte
	if len(password) < 1 {
		encrypted, err := x509.EncryptPEMBlock(rand.Reader, privBlock.Type, privBlock.Bytes, password, PEMCipher)
		if err != nil {
			return err
		}
		privateKey = pem.EncodeToMemory(encrypted)
	} else {
		privateKey = pem.EncodeToMemory(privBlock)
	}

	if err := os.WriteFile(privateKeyFilePath, privateKey, PrivKeyFilePermissions); err != nil {
		return err
	}

	return os.WriteFile(publicKeyFilePath, pem.EncodeToMemory(pubBlock), PubKeyFilePermissions)
}
