package stboot

import (
	"crypto"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	eth                = "eth0"
	BootFilePath       = "root/stboot.zip"
	rootCACertPath     = "/root/LetsEncrypt_Authority_X3.pem"
	entropyAvail       = "/proc/sys/kernel/random/entropy_avail"
	interfaceUpTimeout = 10 * time.Second
)

var debug = func(string, ...interface{}) {}

// VerifySignatureInPath takes path as rootPath and walks
// the directory. Every .cert file it sees, it verifies the .cert
// file with the root certificate, checks if a .signture file
// exists, verify if the signature is correct according to the
// hashValue.
func VerifySignatureInPath(path string, hashValue []byte, rootCert []byte, minAmountValid int) error {
	validSignatures := 0

	// Build up tree
	root := x509.NewCertPool()
	ok := root.AppendCertsFromPEM(rootCert)
	if !ok {
		return errors.New("Failed to parse root certificate")
	}

	opts := x509.VerifyOptions{
		Roots: root,
	}

	// Check certs and signatures
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && (filepath.Ext(info.Name()) == ".cert") {
			// Read cert and verify
			userCert, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("unable to read certificate: %v", err)
			}
			block, _ := pem.Decode(userCert)
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return fmt.Errorf("unable to parse certificate: %v", err)
			}
			// verify certificates with root certificate
			_, err = cert.Verify(opts)
			if err != nil {
				return err
			}
			// Read signature and verify it.
			signatureFilename := strings.TrimSuffix(path, filepath.Ext(path)) + ".signature"
			signatureRaw, err := ioutil.ReadFile(signatureFilename)
			if err != nil {
				return fmt.Errorf("unable to read signature at %s with: %v", signatureFilename, err)
			}
			opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash}
			err = rsa.VerifyPSS(cert.PublicKey.(*rsa.PublicKey), crypto.SHA512, hashValue, signatureRaw, opts)
			if err != nil {
				return fmt.Errorf("signature Verification failed for %s with %v", filepath.Base(signatureFilename), err)
			}
			validSignatures++
			debug(fmt.Sprintf("%s verfied.", signatureFilename))
		}
		return nil
	})
	if validSignatures < minAmountValid {
		return fmt.Errorf("Did not found enough valid signatures. Only %d (%d required) are valid", validSignatures, minAmountValid)
	}

	return nil
}

// DownloardFromHTTPS downloads the stboot.zip file
// to a specific destination via HTTPS.
func DownloadFromHTTPS(url string, destination string) error {

	roots := x509.NewCertPool()
	if err := LoadAndVerifyCertificate(roots); err != nil {
		return fmt.Errorf("Failed to verify root certificate: %v", err)
	}

	// setup https client
	client := http.Client{
		Transport: (&http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: (&tls.Config{
				RootCAs: roots,
			}),
		}),
	}

	// check available kernel entropy
	e, err := ioutil.ReadFile(entropyAvail)
	es := strings.TrimSpace(string(e))
	entr, err := strconv.Atoi(es) // XXX: Insecure?
	if err != nil {
		return fmt.Errorf("Cannot evaluate entropy, %v", err)
	}
	debug("Available kernel entropy: %d", entr)
	if entr < 128 {
		log.Print("WARNING: low entropy!")
		log.Printf("%s : %d", entropyAvail, entr)
	}
	// get remote boot bundle
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("non-200 HTTP status: %d", resp.StatusCode)
	}
	f, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed create boot config file: %v", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write boot config file: %v", err)
	}

	return nil
}

// LoadAndVerifyCertificate loads the certificate needed
// for HTTPS and verifies it.
func LoadAndVerifyCertificate(roots *x509.CertPool) error {
	// load CA certificate
	debug("Load %s as CA certificate", rootCACertPath)
	rootCertBytes, err := ioutil.ReadFile(rootCACertPath)
	if err != nil {
		return fmt.Errorf("Failed to read CA root certificate file: %v", err)
	}
	rootCertPem, _ := pem.Decode(rootCertBytes)
	if rootCertPem.Type != "CERTIFICATE" {
		return fmt.Errorf("Failed decoding certificate: Certificate is of the wrong type. PEM Type is: %s", rootCertPem.Type)
	}
	ok := roots.AppendCertsFromPEM([]byte(rootCertBytes))
	if !ok {

		return fmt.Errorf("Error parsing CA root certificate")
	}

	return nil
}
