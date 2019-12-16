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

type Signer interface {
	hash(files ...string) (hash []byte, err error)
	sign(privKey string, data []byte) (sig []byte, err error)
	verify(sig signature, hash []byte) (err error)
}

type dummySigner struct{}

type sha512PssSigner struct{}

func (dummySigner) hash(files ...string) (hash []byte, err error) {
	hash = make([]byte, 8)
	rand.Read(hash)
	return
}

func (dummySigner) sign(privKey string, data []byte) (sig []byte, err error) {
	sig = make([]byte, 8)
	rand.Read(sig)
	return
}

func (dummySigner) verify(sig signature, hash []byte) (err error) {
	return nil
}

func (sha512PssSigner) hash(files ...string) (hash []byte, err error) {
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

func (sha512PssSigner) sign(privKey string, data []byte) (sig []byte, err error) {
	buf, err := ioutil.ReadFile(privKey)
	if err != nil {
		return
	}

	privPem, _ := pem.Decode(buf)
	key, err := x509.ParsePKCS1PrivateKey(privPem.Bytes)
	if err != nil {
		return
	}
	if key == nil {
		err = fmt.Errorf("key is empty")
		return
	}

	opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash}

	sig, err = rsa.SignPSS(rand.Reader, key, crypto.SHA512, data, opts)
	if err != nil {
		return
	}
	if sig == nil {
		err = fmt.Errorf("signature is nil")
		return
	}
	return
}

func (sha512PssSigner) verify(sig signature, hash []byte) (err error) {
	opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash}
	err = rsa.VerifyPSS(sig.Cert.PublicKey.(*rsa.PublicKey), crypto.SHA512, hash, sig.Bytes, opts)
	if err != nil {
		return fmt.Errorf("signature verification failed: %v", err)
	}
	return
}

// parseCertificate parses certificate from raw certificate
func parseCertificate(raw []byte) (cert *x509.Certificate, err error) {
	block, _ := pem.Decode(raw)

	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return
	}
	return
}

func certPool(pem []byte) (certPool *x509.CertPool, err error) {
	certPool = x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(pem)
	if !ok {
		err = errors.New("Failed to parse root certificate")
		return
	}
	return
}

func validateCertificate(cert *x509.Certificate, cerPool *x509.CertPool) (err error) {
	opts := x509.VerifyOptions{
		Roots: cerPool,
	}
	_, err = cert.Verify(opts)
	if err != nil {
		return
	}
	return
}
