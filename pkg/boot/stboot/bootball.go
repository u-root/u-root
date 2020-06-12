// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stboot

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/boot/jsonboot"
	"github.com/u-root/u-root/pkg/uzip"
)

const (
	signaturesDirName string = "signatures"
	rootCertName      string = "root.cert"
)

// BootBall contains data to operate on the system transparency
// bootball archive. There is an underlying temporary directory
// representing the extracted archive.
type BootBall struct {
	Archive        string
	dir            string
	config         *Stconfig
	numBootConfigs int
	bootFiles      map[string][]string
	RootCertPEM    []byte
	signatures     map[string][]Signature
	NumSignatures  int
	hashes         map[string][]byte
	Signer         Signer
}

// BootBallFromArchive constructs a BootBall zip file at archive
func BootBallFromArchive(archive string) (*BootBall, error) {
	var ball = new(BootBall)

	dir, err := ioutil.TempDir("", "bootball")
	if err != nil {
		return ball, fmt.Errorf("BootBall: cannot create tmp dir: %v", err)
	}

	err = uzip.FromZip(archive, dir)
	if err != nil {
		return ball, fmt.Errorf("BootBall: cannot unzip %s: %v", archive, err)
	}

	cfg, err := getConfig(filepath.Join(dir, ConfigName))
	if err != nil {
		return ball, fmt.Errorf("BootBall: getting configuration faild: %v", err)
	}

	ball.Archive = archive
	ball.dir = dir
	ball.config = cfg

	err = ball.init()
	if err != nil {
		return ball, err
	}

	return ball, nil
}

// BootBallFromConfig constructs a BootBall from a stconfig.json at configFile.
// the underlying tmporary directory is created with standardized paths and an
// updated copy of stconfig.json
func BootBallFromConfig(configFile string) (*BootBall, error) {
	var ball = new(BootBall)

	archive := filepath.Join(filepath.Dir(configFile), BallName)

	cfg, err := getConfig(configFile)
	if err != nil {
		return ball, fmt.Errorf("BootBall: getting configuration faild: %v", err)
	}

	dir, err := makeConfigDir(cfg, filepath.Dir(configFile))
	if err != nil {
		return ball, fmt.Errorf("BootBall: creating standard configuration directory faild: %v", err)
	}

	ball.Archive = archive
	ball.dir = dir
	ball.config = cfg

	err = ball.init()
	if err != nil {
		return ball, err
	}

	return ball, nil
}

func (ball *BootBall) init() error {
	certPEM, err := ioutil.ReadFile(filepath.Join(ball.dir, signaturesDirName, rootCertName))
	if err != nil {
		return fmt.Errorf("BootBall: reading root certificate faild: %v", err)
	}

	bootFiles, err := getBootFiles(ball.config, ball.dir)
	if err != nil {
		return fmt.Errorf("BootBall: getting boot files faild: %v", err)
	}

	ball.RootCertPEM = certPEM
	ball.numBootConfigs = len(ball.config.BootConfigs)
	ball.bootFiles = bootFiles
	ball.Signer = Sha512PssSigner{}

	err = ball.getSignatures()
	if err != nil {
		return fmt.Errorf("BootBall: getting signatures: %v", err)
	}

	var x int = 0
	for _, sigPool := range ball.signatures {
		if x == 0 {
			x = len(sigPool)
			continue
		}
		if len(sigPool) != x {
			return errors.New("BootBall: invalid map of signatures")
		}
	}
	ball.NumSignatures = x
	return nil
}

// Clean removes the underlying temporary directory.
func (ball *BootBall) Clean() error {
	err := os.RemoveAll(ball.dir)
	if err != nil {
		return err
	}
	ball.dir = ""
	return nil
}

// Pack archives the all contents of the underlying temporary
// directory using zip.
func (ball *BootBall) Pack() error {
	if ball.Archive == "" || ball.dir == "" {
		return errors.New("BootBall.Pacstandak: booball.archive and bootball.dir must be set")
	}
	return uzip.ToZip(ball.dir, ball.Archive)
}

// Dir returns the temporary directory associated with BootBall.
func (ball *BootBall) Dir() string {
	return ball.dir
}

// GetBootConfigByIndex returns the Bootconfig at index from the BootBall's configs arrey.
func (ball *BootBall) GetBootConfigByIndex(index int) (*jsonboot.BootConfig, error) {
	bc, err := ball.config.getBootConfig(index)
	if err != nil {
		return nil, err
	}
	bc.SetFilePathsPrefix(ball.dir)
	return bc, nil
}

// Hash calculates hashes of all boot configurations in BootBall using the
// BootBall.Signer's hash function
func (ball *BootBall) Hash() error {
	ball.hashes = make(map[string][]byte)
	for key, files := range ball.bootFiles {
		hash, herr := ball.Signer.Hash(files...)
		if herr != nil {
			return herr
		}
		ball.hashes[key] = hash
	}
	return nil
}

// Sign signes the hashes of all boot configurations in BootBall using the
// BootBall.Signer's hash function with the provided privKeyFile. The signature
// is stored along with the provided certFile inside the BootBall.
func (ball *BootBall) Sign(privKeyFile, certFile string) error {
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

	log.Printf("Signing with: %s", privKeyFile)

	if ball.hashes == nil {
		err = ball.Hash()
		if err != nil {
			return err
		}
	}

	for key, hash := range ball.hashes {
		s, err := ball.Signer.Sign(privKeyFile, hash)
		if err != nil {
			return err
		}
		sig := Signature{
			Bytes: s,
			Cert:  cert}
		ball.signatures[key] = append(ball.signatures[key], sig)
		d := filepath.Join(ball.dir, signaturesDirName, key)
		if err = writeSignature(d, certFile, sig); err != nil {
			return err
		}
	}

	ball.NumSignatures++
	return nil
}

// VerifyBootconfigByID validates the certificates stored together with the
// signatures of BootConfig id and verifies the signatures. The number of
// valid signatures is returned.
func (ball *BootBall) VerifyBootconfigByID(id string) (found, verified int, err error) {
	if ball.hashes == nil {
		err := ball.Hash()
		if err != nil {
			return 0, 0, err
		}
	}

	found = 0
	verified = 0
	for _, sig := range ball.signatures[id] {
		err := validateCertificate(sig.Cert, ball.RootCertPEM)
		if err != nil {
			return found, verified, err
		}
		found++
		err = ball.Signer.Verify(sig, ball.hashes[id])
		if err != nil {
			log.Print(err)
		}
		verified++
	}
	return found, verified, nil
}

// getConfig returns a Stconfig struct from a JSON file at src
func getConfig(src string) (*Stconfig, error) {
	cfgBytes, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}
	cfg, err := stconfigFromBytes(cfgBytes)
	if err != nil {
		return nil, err
	}
	if !(cfg.IsValid()) {
		return nil, errors.New("invalid configuration")
	}
	return cfg, nil
}

// getBootFiles returns the file paths of all files of a u-root bootconfig
// for all bootconfigs in cfg.BootConfigs. Prefix is added in front of each
// file path. The map's keys are set to the respective bootconfig's name.
// An error is returned in case one of the files does not exist.
func getBootFiles(cfg *Stconfig, prefix string) (map[string][]string, error) {
	bootFiles := make(map[string][]string)
	for _, bc := range cfg.BootConfigs {
		files := make([]string, 0)
		for _, file := range bc.FileNames() {
			file = filepath.Join(prefix, file)
			if _, err := os.Stat(file); err != nil {
				return nil, err
			}
			files = append(files, file)
		}
		bootFiles[bc.ID()] = files
	}
	return bootFiles, nil
}

// getSignatures initializes ball.signatures with the corresponding signatures
// and certificates found in the signatures folder (stboot.signaturesDirName)
// of ball's underlying tmpDir (ball.dir). An error is returned if one of the
// files cannot be read or parsed.
func (ball *BootBall) getSignatures() error {
	ball.signatures = make(map[string][]Signature)
	path := filepath.Join(ball.dir, signaturesDirName)

	sigPool := make([]Signature, 0)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
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
			sigPool = append(sigPool, sig)
			key := filepath.Base(filepath.Dir(path))
			ball.signatures[key] = sigPool
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// writeSignature writes the signature represented by sig to a file in
// dir along with a copy of certFile. The filenames are composed of the
// first piece of the public key of the certificate.
func writeSignature(dir, certFile string, sig Signature) error {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("%x", sig.Cert.PublicKey)[2:18]
	sigName := fmt.Sprintf("%s.signature", id)
	sigPath := filepath.Join(dir, sigName)
	err = ioutil.WriteFile(sigPath, sig.Bytes, 0644)
	if err != nil {
		return err
	}

	certName := fmt.Sprintf("%s.cert", id)
	certPath := filepath.Join(dir, certName)
	return copyFile(certFile, certPath)
}

// makeConfigDir copies the files named in cfg to well known directory tree
// inside the bootball's underlying tmpDir since the files in user's cfg can
// reside anywhere in the file system. The created tmpDir is returned.
// If one of the files in cfg does not exist or copying fails an error is
// returned.
func makeConfigDir(cfg *Stconfig, origDir string) (string, error) {
	dir, err := ioutil.TempDir(os.TempDir(), "bootball")
	if err != nil {
		return "", err
	}

	dstPath := filepath.Join(dir, signaturesDirName, rootCertName)
	srcPath := filepath.Join(origDir, cfg.RootCertPath)
	if err := copyFile(srcPath, dstPath); err != nil {
		return "", err
	}

	for i, bc := range cfg.BootConfigs {
		dirName := bc.ID()
		for _, file := range bc.FileNames() {
			fileName := filepath.Base(file)
			dstPath := filepath.Join(dir, dirName, fileName)
			srcPath := filepath.Join(origDir, file)
			if err := copyFile(srcPath, dstPath); err != nil {
				return "", err
			}
		}

		bc.ChangeFilePaths(dirName)
		cfg.BootConfigs[i] = bc
	}

	dstPath = filepath.Join(dir, ConfigName)
	bytes, err := cfg.bytes()
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(dstPath, bytes, os.ModePerm)
	if err != nil {
		return "", err
	}

	return dir, nil
}

// copyFiles copies the content of the file at src to dst. If dst does not
// exist it is created. In case case src does not exist, creation of dst
// or copying fails an error is returned
func copyFile(src, dst string) error {
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
