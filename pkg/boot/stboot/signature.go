// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stboot

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Signature contains the signature bytes and the
// corresponding certificate.
type Signature struct {
	Bytes []byte
	Cert  *x509.Certificate
}

// Write saves the signature and the certificate represented by s to files at
// a path named by dir. The filenames are composed of the first piece of the
// certificate's public key. The file extensions are '.signature' and '.cert'.
func (s *Signature) Write(dir string) error {
	stat, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("not a directory: %s", dir)
	}

	id := fmt.Sprintf("%x", s.Cert.PublicKey)[2:18]
	sigName := fmt.Sprintf("%s.signature", id)
	sigPath := filepath.Join(dir, sigName)
	err = ioutil.WriteFile(sigPath, s.Bytes, os.ModePerm)
	if err != nil {
		return err
	}

	certName := fmt.Sprintf("%s.cert", id)
	certPath := filepath.Join(dir, certName)
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: s.Cert.Raw,
	}
	var certBuf bytes.Buffer
	if err := pem.Encode(&certBuf, block); err != nil {
		return (err)
	}
	return ioutil.WriteFile(certPath, certBuf.Bytes(), os.ModePerm)
}

// Signer is used by BootBall to hash, sign and varify the Bootball.
type Signer interface {
	Hash(files ...string) ([]byte, error)
	Sign(privKey string, data []byte) ([]byte, error)
	Verify(sig Signature, hash []byte) error
}

// DummySigner implements the Signer interface. It creates signatures
// that are always valid.
type DummySigner struct{}

// Hash returns a hash value of just 8 random bytes.
func (DummySigner) Hash(files ...string) ([]byte, error) {
	hash := make([]byte, 8)
	_, err := rand.Read(hash)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// Sign returns a signature containing just 8 random bytes.
func (DummySigner) Sign(privKey string, data []byte) ([]byte, error) {
	sig := make([]byte, 8)
	_, err := rand.Read(sig)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

// Verify will never return an error.
func (DummySigner) Verify(sig Signature, hash []byte) error {
	return nil
}

// Sha512PssSigner implements the Signer interface. It uses SHA512 hashes
// and PSS signatures along with x509 certificates.
type Sha512PssSigner struct{}

// Hash returns the a SHA512 hash value of the provided files.
func (Sha512PssSigner) Hash(files ...string) ([]byte, error) {
	h := sha512.New()
	for _, file := range files {
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		_, err = h.Write(buf)
		if err != nil {
			return nil, err
		}
	}
	return h.Sum(nil), nil
}

// Sign signes the provided data with the key named by privKey. The returned
// byte slice contains a PSS signature value.
func (Sha512PssSigner) Sign(privKey string, data []byte) ([]byte, error) {
	buf, err := ioutil.ReadFile(privKey)
	if err != nil {
		return nil, err
	}

	privPem, _ := pem.Decode(buf)
	key, err := x509.ParsePKCS1PrivateKey(privPem.Bytes)
	if err != nil {
		return nil, err
	}
	if key == nil {
		err = fmt.Errorf("key is empty")
		return nil, err
	}

	opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash}

	sig, err := rsa.SignPSS(rand.Reader, key, crypto.SHA512, data, opts)
	if err != nil {
		return nil, err
	}
	if sig == nil {
		return nil, fmt.Errorf("signature is nil")
	}
	return sig, nil
}

// Verify checks if sig contains a valid signature of hash.
func (Sha512PssSigner) Verify(sig Signature, hash []byte) error {
	opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash}
	err := rsa.VerifyPSS(sig.Cert.PublicKey.(*rsa.PublicKey), crypto.SHA512, hash, sig.Bytes, opts)
	if err != nil {
		return fmt.Errorf("signature verification failed: %v", err)
	}
	return nil
}

// parseCertificate parses a x509 certificate from raw data.
func parseCertificate(raw []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(raw)
	return x509.ParseCertificate(block.Bytes)
}

// certPool returns a x509 certificate pool from PEM encoded data.
func certPool(pem []byte) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(pem)
	if !ok {
		return nil, errors.New("Failed to parse root certificate")
	}
	return certPool, nil
}

// validateCertificate validates cert against certPool. If cert is not signed
// by a certificate of certPool an error is returned.
func validateCertificate(cert *x509.Certificate, rootCertPEM []byte) error {
	certPool, err := certPool(rootCertPEM)
	if err != nil {
		return err
	}
	opts := x509.VerifyOptions{
		Roots: certPool,
	}
	_, err = cert.Verify(opts)
	return err
}
