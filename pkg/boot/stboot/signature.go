package stboot

import (
	"archive/zip"
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
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// AddSignature signes the content of a ST bootball and inserts the
// signature into the archive along with the respective certificate
func AddSignature(bootball, privKey, certFile string) error {

	cfg, dir, err := FromZip(bootball)
	if err != nil {
		return err
	}

	// collect boot binaries
	// XXX Refactor if we remove bootconfig from manifest
	// Maybe just walk through certs/folders and match do root/bootconfig
	for i := range cfg.BootConfigs {

		bootconfigDir := path.Join(dir, fmt.Sprintf("bootconfig_%d", i))

		bcHash, err := hashBootconfigDir(bootconfigDir)
		if err != nil {
			return fmt.Errorf("failed to hash bootconfig - Err %s", err)
		}

		// Sign hash with Key
		buff, err := ioutil.ReadFile(privKey)
		if err != nil {
			return fmt.Errorf("cannot read key file %s: %v", privKey, err)
		}
		privPem, _ := pem.Decode(buff)
		rsaPrivKey, err := x509.ParsePKCS1PrivateKey(privPem.Bytes)
		if err != nil {
			return fmt.Errorf("cannot parse private key: %v", err)
		}
		if rsaPrivKey == nil {
			return fmt.Errorf("RSA key is nil: %v", err)
		}

		opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash}

		signature, err := rsa.SignPSS(rand.Reader, rsaPrivKey, crypto.SHA512, bcHash, opts)
		if err != nil {
			return fmt.Errorf("signature generation failed: %v", err)
		}
		if signature == nil {
			return fmt.Errorf("signature is nil, %v", err)
		}

		// Create dir for signature
		d := path.Join(dir, fmt.Sprintf("certs/bootconfig_%d/", i))
		err = os.MkdirAll(d, os.ModeDir|os.FileMode(0700))
		if err != nil {
			return fmt.Errorf("creating signatures directories %s failed: %v", dir, err)
		}

		// Extract part of Public Key for identification
		certBytes, err := ioutil.ReadFile(certFile)
		if err != nil {
			return fmt.Errorf("cannot read certificate file %s: %v", certFile, err)
		}

		cert, err := parseCertificate(certBytes)
		if err != nil {
			return fmt.Errorf("failed to parse certificate %s: %v", certFile, err)
		}

		// Write signature to folder
		err = ioutil.WriteFile(path.Join(dir, fmt.Sprintf("certs/bootconfig_%d/%s.signature", i, fmt.Sprintf("%x", cert.PublicKey)[2:18])), signature, 0644)
		if err != nil {
			return fmt.Errorf("writing into %v failed with %v - Check permissions", dir, err)
		}

		// cp cert to folder
		err = ioutil.WriteFile(path.Join(dir, fmt.Sprintf("certs/bootconfig_%d/%s.cert", i, fmt.Sprintf("%x", cert.PublicKey)[2:18])), certBytes, 0644)
		if err != nil {
			return fmt.Errorf("cannot copy certificate %s to archive: %v", certFile, err)
		}
	}

	// Pack it again
	// Create a buffer to write the archive to.
	buf := new(bytes.Buffer)
	// Create a new zip archive.
	z := zip.NewWriter(buf)

	// Walk the directory and pack it.
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			err := tozip(z, strings.Replace(path, dir, "", -1)[1:], path)
			if err != nil {
				return fmt.Errorf(fmt.Sprintf("Error adding file %s to .zip archive again", strings.Replace(path, dir, "", -1)))
			}
		}

		return nil
	})

	z.Close()

	pathToZip := fmt.Sprintf("./.original/%d", time.Now().Unix())
	os.MkdirAll(pathToZip, os.ModePerm)
	os.Rename(bootball, pathToZip+"/stboot.zip")

	err = ioutil.WriteFile(bootball, buf.Bytes(), 0777)
	if err != nil {
		return fmt.Errorf("unable to write new stboot.zip file - recover old from %s", pathToZip)
	}
	os.RemoveAll(pathToZip)

	return nil

}

// HashBootconfigDir hashes every file inside bootconigDir and returns a
// SHA512 hash
func hashBootconfigDir(bootconfigDir string) ([]byte, error) {

	hash := sha512.New()
	hash.Reset()

	files, err := ioutil.ReadDir(bootconfigDir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() {
			p := path.Join(bootconfigDir, file.Name())
			buf, err := ioutil.ReadFile(p)
			if err != nil {
				return nil, err
			}
			hash.Write(buf)

		}
	}
	return hash.Sum(nil), nil
}

// parseCertificate parses certificate from raw certificate
func parseCertificate(rawCertificate []byte) (x509.Certificate, error) {

	block, _ := pem.Decode(rawCertificate)
	if block == nil {
		log.Fatal("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return *pub, fmt.Errorf("failed to parse DER encoded public key: %v", err)
	}

	return *pub, nil
}

func certPool(pem []byte) (*x509.CertPool, error) {
	root := x509.NewCertPool()
	ok := root.AppendCertsFromPEM(pem)
	if !ok {
		return nil, errors.New("Failed to parse root certificate")
	}
	return root, nil
}

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
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
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
				return fmt.Errorf("unable to verify %s with root certificate: %v", path, err)
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
			log.Print(fmt.Sprintf("%s verfied.", signatureFilename))
		}
		return nil
	})
	if err != nil {
		return err
	}
	if validSignatures < minAmountValid {
		return fmt.Errorf("Did not found enough valid signatures. Only %d (%d required) are valid", validSignatures, minAmountValid)
	}

	return nil
}
