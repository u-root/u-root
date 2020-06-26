// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stboot

import (
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/multiboot"
	"github.com/u-root/u-root/pkg/uio"
	"github.com/u-root/u-root/pkg/uzip"
)

const (
	bootfilesDir  string = "boot"
	acmDir        string = "boot/acms"
	signaturesDir string = "signatures"
	rootCertPath  string = "signatures/root.cert"
)

// Bootball contains data to operate on the system transparency
// bootball archive. There is an underlying temporary directory
// representing the extracted archive.
type Bootball struct {
	Archive           string
	Dir               string
	Config            *Stconfig
	FilesToBeMeasured []string
	RootCertPEM       []byte
	Signatures        []Signature
	NumSignatures     int
	HashValue         []byte
	Signer            Signer
}

// BootballFromArchive constructs a Bootball from a zip file at archive.
func BootballFromArchive(archive string) (*Bootball, error) {
	var ball = &Bootball{}

	if _, err := os.Stat(archive); err != nil {
		return nil, fmt.Errorf("Bootball: %v", err)
	}

	dir, err := ioutil.TempDir("", "bootball")
	if err != nil {
		return nil, fmt.Errorf("Bootball: cannot create tmp dir: %v", err)
	}

	err = uzip.FromZip(archive, dir)
	if err != nil {
		return nil, fmt.Errorf("Bootball: cannot unzip %s: %v", archive, err)
	}

	cfg, err := StconfigFromFile(filepath.Join(dir, ConfigName))
	if err != nil {
		return nil, fmt.Errorf("Bootball: getting configuration faild: %v", err)
	}
	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("Bootball: invalid config: %v", err)
	}

	ball.Archive = archive
	ball.Dir = dir
	ball.Config = cfg

	err = ball.init()
	if err != nil {
		return ball, err
	}

	return ball, nil
}

// InitBootball constructs a Bootball from the parsed files. The underlying
// tmporary directory is created with standardized paths and names.
func InitBootball(outDir, label, kernel, initramfs, cmdline, tboot, tbootArgs, rootCert string, acms []string, allowNonTXT bool) (*Bootball, error) {
	var ball = &Bootball{}

	t := time.Now()
	tstr := fmt.Sprintf("%04d-%02d-%02d-%02d-%02d-%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	name := "ball-" + tstr + BootballExt
	ball.Archive = filepath.Join(outDir, name)

	dir, cfg, err := createFileTree(kernel, initramfs, tboot, rootCert, acms)
	if err != nil {
		return nil, fmt.Errorf("Bootball: creating standard file tree faild: %v", err)
	}

	cfg.Label = label
	cfg.Cmdline = cmdline
	cfg.AllowNonTXT = allowNonTXT
	cfg.Write(dir)

	ball.Dir = dir
	ball.Config = &cfg

	err = ball.init()
	if err != nil {
		return nil, err
	}

	return ball, nil
}

func (ball *Bootball) init() error {
	certPEM, err := ioutil.ReadFile(filepath.Join(ball.Dir, rootCertPath))
	if err != nil {
		return fmt.Errorf("Bootball: reading root certificate faild: %v", err)
	}
	ball.RootCertPEM = certPEM

	err = ball.getFilesToBeHashed()
	if err != nil {
		return fmt.Errorf("Bootball: collecting files for measurement failed: %v", err)
	}

	ball.Signer = Sha512PssSigner{}

	err = ball.getSignatures()
	if err != nil {
		return fmt.Errorf("Bootball: getting signatures: %v", err)
	}

	ball.NumSignatures = len(ball.Signatures)
	return nil
}

// Clean removes the underlying temporary directory.
func (ball *Bootball) Clean() error {
	err := os.RemoveAll(ball.Dir)
	if err != nil {
		return err
	}
	ball.Dir = ""
	return nil
}

// Pack archives the all contents of the underlying temporary
// directory using zip.
func (ball *Bootball) Pack() error {
	if ball.Archive == "" {
		return errors.New("Booball.Archive is not set")
	}
	if ball.Dir == "" {
		return errors.New("Cannot locate underlying directory")
	}
	return uzip.ToZip(ball.Dir, ball.Archive)
}

// Hash calculates hashes of all boot configurations in Bootball using the
// Bootball.Signer's hash function.
func (ball *Bootball) Hash() error {
	hash, err := ball.Signer.Hash(ball.FilesToBeMeasured...)
	if err != nil {
		return err
	}
	ball.HashValue = hash
	return nil
}

// Sign signes ball.HashValue using ball.Signer with the private key named by
// privKeyFile. The certificate named by certFile is supposed to correspond
// to the private key. Both, the signature and the certificate are stored into
// the Bootball.
func (ball *Bootball) Sign(privKeyFile, certFile string) error {
	if _, err := os.Stat(privKeyFile); err != nil {
		return err
	}

	buf, err := ioutil.ReadFile(certFile)
	if err != nil {
		return err
	}

	cert, err := parseCertificate(buf)
	if err != nil {
		return err
	}

	err = validateCertificate(cert, ball.RootCertPEM)
	if err != nil {
		return err
	}

	if ball.HashValue == nil {
		err = ball.Hash()
		if err != nil {
			return err
		}
	}

	// check for dublicate certificates
	for _, sig := range ball.Signatures {
		if sig.Cert.Equal(cert) {
			return fmt.Errorf("certificate has already been used: %v", certFile)
		}
	}
	// sign with private key
	s, err := ball.Signer.Sign(privKeyFile, ball.HashValue)
	if err != nil {
		return err
	}
	sig := Signature{
		Bytes: s,
		Cert:  cert}
	// check certificate's public key
	err = ball.Signer.Verify(sig, ball.HashValue)
	if err != nil {
		return fmt.Errorf("public key in %s does not match the private key %s", filepath.Base(certFile), filepath.Base(privKeyFile))
	}
	// save
	ball.Signatures = append(ball.Signatures, sig)
	dir := filepath.Join(ball.Dir, signaturesDir)
	if err = sig.Write(dir); err != nil {
		return err
	}

	ball.NumSignatures++
	return nil
}

// Verify first validates the certificates stored together with the signatures
// and the verifies the signatures. The number of found signatures and the
// number of valid signatures are returned. A signature is valid if:
// * Its certificate was signed by balls's root certificate
// * Verification is passed
// * No previous signature has the same certificate
func (ball *Bootball) Verify() (found, verified int, err error) {
	if ball.HashValue == nil {
		err := ball.Hash()
		if err != nil {
			return 0, 0, err
		}
	}

	found = 0
	verified = 0
	var certsUsed []*x509.Certificate
	for i, sig := range ball.Signatures {
		found++
		err := validateCertificate(sig.Cert, ball.RootCertPEM)
		if err != nil {
			log.Printf("skip signature %d: invalid certificate: %v", i+1, err)
			continue
		}
		var dublicate bool
		for _, c := range certsUsed {
			if c.Equal(sig.Cert) {
				dublicate = true
				break
			}
		}
		if dublicate {
			log.Printf("skip signature %d: dublicate", i+1)
			continue
		}
		certsUsed = append(certsUsed, sig.Cert)
		err = ball.Signer.Verify(sig, ball.HashValue)
		if err != nil {
			log.Printf("skip signature %d: verification failed: %v", i+1, err)
			continue
		}
		verified++
	}
	return found, verified, nil
}

// OSImage retunrns a boot.OSImage generated from ball's configuration
func (ball *Bootball) OSImage(txt bool) (boot.OSImage, error) {
	err := ball.Config.Validate()
	if err != nil {
		return nil, err
	}

	if txt && ball.Config.Tboot == "" {
		return nil, errors.New("Bootball does not contain a TXT-ready configuration")
	}

	if !txt && !ball.Config.AllowNonTXT {
		return nil, errors.New("Bootball requires the use of TXT")
	}

	var osi boot.OSImage
	if !txt {
		osi = &boot.LinuxImage{
			Name:    ball.Config.Label,
			Kernel:  uio.NewLazyFile(filepath.Join(ball.Dir, ball.Config.Kernel)),
			Initrd:  uio.NewLazyFile(filepath.Join(ball.Dir, ball.Config.Initramfs)),
			Cmdline: ball.Config.Cmdline,
		}
		return osi, nil
	}

	var modules []multiboot.Module
	kernel := multiboot.Module{
		Module:  uio.NewLazyFile(filepath.Join(ball.Dir, ball.Config.Kernel)),
		Name:    "OS-Kernel",
		CmdLine: ball.Config.Cmdline,
	}
	modules = append(modules, kernel)

	initramfs := multiboot.Module{
		Module: uio.NewLazyFile(filepath.Join(ball.Dir, ball.Config.Initramfs)),
		Name:   "OS-Initramfs",
	}
	modules = append(modules, initramfs)

	for n, a := range ball.Config.ACMs {
		acm := multiboot.Module{
			Module: uio.NewLazyFile(filepath.Join(ball.Dir, a)),
			Name:   fmt.Sprintf("ACM%d", n+1),
		}
		modules = append(modules, acm)
	}

	osi = &boot.MultibootImage{
		Name:    ball.Config.Label,
		Kernel:  uio.NewLazyFile(filepath.Join(ball.Dir, ball.Config.Tboot)),
		Cmdline: ball.Config.TbootArgs,
		Modules: modules,
	}
	return osi, nil
}

// getFilesToBeHashed the paths of the bootball' files that are supposed
// to be hashed for signing and varifiaction. These are:
// * stconfig.json
// * root.cert
// * files defined in stconfig.json if they are present
func (ball *Bootball) getFilesToBeHashed() error {
	var f []string

	// these files must be present
	config := filepath.Join(ball.Dir, ConfigName)
	kernel := filepath.Join(ball.Dir, ball.Config.Kernel)
	rootCert := filepath.Join(ball.Dir, rootCertPath)
	_, err := os.Stat(config)
	if err != nil {
		return errors.New("files to be measured: missing stconfig.json")
	}
	_, err = os.Stat(kernel)
	if err != nil {
		return errors.New("files to be measured: missing kernel")
	}
	_, err = os.Stat(rootCert)
	if err != nil {
		return errors.New("files to be measured: missing root certificate")
	}
	f = append(f, config, kernel, rootCert)

	// following files are measured if present
	if ball.Config.Initramfs != "" {
		initramfs := filepath.Join(ball.Dir, ball.Config.Initramfs)
		_, err = os.Stat(initramfs)
		if err == nil {
			f = append(f, initramfs)
		}
	}
	if ball.Config.Tboot != "" {
		tboot := filepath.Join(ball.Dir, ball.Config.Tboot)
		_, err = os.Stat(tboot)
		if err == nil {
			f = append(f, tboot)
		}
	}
	for _, acm := range ball.Config.ACMs {
		a := filepath.Join(ball.Dir, acm)
		_, err = os.Stat(a)
		if err == nil {
			f = append(f, a)
		}
	}

	ball.FilesToBeMeasured = f
	return nil
}

// getSignatures initializes ball.signatures with the corresponding signatures
// and certificates found in the signatures folder of ball's underlying tmpDir.
// An error is returned if one of the files cannot be read or parsed.
func (ball *Bootball) getSignatures() error {
	root := filepath.Join(ball.Dir, signaturesDir)

	sigs := make([]Signature, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(info.Name())

		if !info.IsDir() && (ext == ".signature") {
			sigBytes, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			certFile := strings.TrimSuffix(path, filepath.Ext(path)) + ".cert"
			certBytes, err := ioutil.ReadFile(certFile)
			if err != nil {
				return err
			}

			cert, err := parseCertificate(certBytes)
			if err != nil {
				return err
			}

			sig := Signature{
				Bytes: sigBytes,
				Cert:  cert,
			}
			sigs = append(sigs, sig)
			ball.Signatures = sigs
		}
		return nil
	})
	return err
}

// createFileTree copies the provided files to a well known tree inside
// the bootball's underlying tmpDir. The created tmpDir and a Stconfig
// initialized with corresponding paths is retruned.
func createFileTree(kernel, initramfs, tboot, rootCert string, acms []string) (dir string, cfg Stconfig, err error) {
	dir, err = ioutil.TempDir(os.TempDir(), "bootball")
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		}
	}()

	var dst, rel string

	// Kernel
	if kernel == "" {
		err = errors.New("kernel missing")
	}
	rel = filepath.Join(bootfilesDir, filepath.Base(kernel))
	dst = filepath.Join(dir, rel)
	if err = createAndCopy(kernel, dst); err != nil {
		return
	}
	cfg.Kernel = rel

	// Initramfs
	if initramfs != "" {
		rel = filepath.Join(bootfilesDir, filepath.Base(initramfs))
		dst = filepath.Join(dir, rel)
		if err = createAndCopy(initramfs, dst); err != nil {
			return
		}
		cfg.Initramfs = rel
	}

	// tboot
	if tboot != "" {
		rel = filepath.Join(bootfilesDir, filepath.Base(tboot))
		dst = filepath.Join(dir, rel)
		if err = createAndCopy(tboot, dst); err != nil {
			return
		}
		cfg.Tboot = rel
	}

	// Root Certificate
	if rootCert == "" {
		err = errors.New("root certificate missing")
	}
	dst = filepath.Join(dir, rootCertPath)
	if err = createAndCopy(rootCert, dst); err != nil {
		return
	}

	// ACMs
	if len(acms) > 0 {
		for _, acm := range acms {
			rel = filepath.Join(acmDir, filepath.Base(acm))
			dst = filepath.Join(dir, rel)
			if err = createAndCopy(acm, dst); err != nil {
				return
			}
			cfg.ACMs = append(cfg.ACMs, rel)
		}
	}

	return
}

// createAndCopy copies the content of the file at src to dst. If dst does not
// exist it is created. In case case src does not exist, creation of dst
// or copying fails and error is returned.
func createAndCopy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err = os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}
