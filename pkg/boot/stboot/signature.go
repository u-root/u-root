package stboot

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
)

type signature struct {
	Bytes []byte
	Cert  *x509.Certificate
}

// Signer is used by BootBall to hash, sign and varify the BootConfigs
// with appropriate algorithms
type Signer interface {
	hash(files ...string) ([]byte, error)
	sign(privKey string, data []byte) ([]byte, error)
	verify(sig signature, hash []byte) error
}

// DummySigner creates signatures that are always valid.
type DummySigner struct{}

func (DummySigner) hash(files ...string) ([]byte, error) {
	hash := make([]byte, 8)
	_, err := rand.Read(hash)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (DummySigner) sign(privKey string, data []byte) ([]byte, error) {
	sig := make([]byte, 8)
	_, err := rand.Read(sig)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func (DummySigner) verify(sig signature, hash []byte) error {
	return nil
}

// Sha512PssSigner uses SHA512 hashes ans PSS signatures along with
// x509 certificates.
type Sha512PssSigner struct{}

func (Sha512PssSigner) hash(files ...string) ([]byte, error) {
	h := sha512.New()
	h.Reset()

	for _, file := range files {
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		h.Write(buf)
	}
	return h.Sum(nil), nil
}

func (Sha512PssSigner) sign(privKey string, data []byte) ([]byte, error) {
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

func (Sha512PssSigner) verify(sig signature, hash []byte) error {
	opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash}
	err := rsa.VerifyPSS(sig.Cert.PublicKey.(*rsa.PublicKey), crypto.SHA512, hash, sig.Bytes, opts)
	if err != nil {
		return fmt.Errorf("signature verification failed: %v", err)
	}
	return nil
}

// parseCertificate parses certificate from raw certificate
func parseCertificate(raw []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(raw)

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func certPool(pem []byte) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(pem)
	if !ok {
		return nil, errors.New("Failed to parse root certificate")
	}
	return certPool, nil
}

func validateCertificate(cert *x509.Certificate, cerPool *x509.CertPool) error {
	opts := x509.VerifyOptions{
		Roots: cerPool,
	}
	_, err := cert.Verify(opts)
	if err != nil {
		return err
	}
	return nil
}
