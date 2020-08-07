// +build !windows

// Binary tpm2-ekcert reads an x509 certificate from a specific NVRAM index.
package main

import (
	"crypto"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

var (
	tpmPath = flag.String("tpm-path", "/dev/tpm0", "Path to the TPM device (character device or a Unix socket)")
	// Default value is defined in section 7.8, "NV Memory" of the latest version pdf on:
	// https://trustedcomputinggroup.org/resource/tcg-tpm-v2-0-provisioning-guidance/
	certIndex = flag.Uint("cert-index", 0x01C00002, "NVRAM index of the certificate file")
	tmplIndex = flag.Uint("template-index", 0, "NVRAM index of the EK template; if zero, default RSA EK template is used")
	outPath   = flag.String("output", "", "File path for output; leave blank to write to stdout")

	// Default EK template defined in:
	// https://trustedcomputinggroup.org/wp-content/uploads/Credential_Profile_EK_V2.0_R14_published.pdf
	defaultEKTemplate = tpm2.Public{
		Type:    tpm2.AlgRSA,
		NameAlg: tpm2.AlgSHA256,
		Attributes: tpm2.FlagFixedTPM | tpm2.FlagFixedParent | tpm2.FlagSensitiveDataOrigin |
			tpm2.FlagAdminWithPolicy | tpm2.FlagRestricted | tpm2.FlagDecrypt,
		AuthPolicy: []byte{
			0x83, 0x71, 0x97, 0x67, 0x44, 0x84,
			0xB3, 0xF8, 0x1A, 0x90, 0xCC, 0x8D,
			0x46, 0xA5, 0xD7, 0x24, 0xFD, 0x52,
			0xD7, 0x6E, 0x06, 0x52, 0x0B, 0x64,
			0xF2, 0xA1, 0xDA, 0x1B, 0x33, 0x14,
			0x69, 0xAA,
		},
		RSAParameters: &tpm2.RSAParams{
			Symmetric: &tpm2.SymScheme{
				Alg:     tpm2.AlgAES,
				KeyBits: 128,
				Mode:    tpm2.AlgCFB,
			},
			KeyBits:    2048,
			ModulusRaw: make([]byte, 256),
		},
	}
)

func main() {
	flag.Parse()

	cert, err := readEKCert(*tpmPath, uint32(*certIndex), uint32(*tmplIndex))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if *outPath == "" {
		fmt.Println(string(cert))
		return
	}
	if err := ioutil.WriteFile(*outPath, cert, os.ModePerm); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func readEKCert(path string, certIdx, tmplIdx uint32) ([]byte, error) {
	rwc, err := tpm2.OpenTPM(path)
	if err != nil {
		return nil, fmt.Errorf("can't open TPM at %q: %v", path, err)
	}
	defer rwc.Close()
	ekCert, err := tpm2.NVRead(rwc, tpmutil.Handle(certIdx))
	if err != nil {
		return nil, fmt.Errorf("reading EK cert: %v", err)
	}
	// Sanity-check that this is a valid certificate.
	cert, err := x509.ParseCertificate(ekCert)
	if err != nil {
		return nil, fmt.Errorf("parsing EK cert: %v", err)
	}

	// Initialize EK and compare public key to ekCert.PublicKey.
	var ekh tpmutil.Handle
	var ekPub crypto.PublicKey
	if tmplIdx != 0 {
		ekTemplate, err := tpm2.NVRead(rwc, tpmutil.Handle(tmplIdx))
		if err != nil {
			return nil, fmt.Errorf("reading EK template: %v", err)
		}
		ekh, ekPub, err = tpm2.CreatePrimaryRawTemplate(rwc, tpm2.HandleEndorsement, tpm2.PCRSelection{}, "", "", ekTemplate)
		if err != nil {
			return nil, fmt.Errorf("creating EK: %v", err)
		}
	} else {
		ekh, ekPub, err = tpm2.CreatePrimary(rwc, tpm2.HandleEndorsement, tpm2.PCRSelection{}, "", "", defaultEKTemplate)
		if err != nil {
			return nil, fmt.Errorf("creating EK: %v", err)
		}
	}
	defer tpm2.FlushContext(rwc, ekh)

	if !reflect.DeepEqual(ekPub, cert.PublicKey) {
		return nil, errors.New("public key in EK certificate differs from public key created via EK template")
	}

	return ekCert, nil
}
