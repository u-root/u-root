// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bootconfig

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
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	urootcrypto "github.com/u-root/u-root/pkg/crypto"
)

// memoryZipReader is used to unpack a zip file from a byte sequence in memory.
type memoryZipReader struct {
	Content []byte
}

func (r *memoryZipReader) ReadAt(p []byte, offset int64) (n int, err error) {
	cLen := int64(len(r.Content))
	if offset > cLen {
		return 0, io.EOF
	}
	if cLen-offset >= int64(len(p)) {
		n = len(p)
		err = nil
	} else {
		err = io.EOF
		n = int(int64(cLen) - offset)
	}
	copy(p, r.Content[offset:int(offset)+n])
	return n, err
}

// FIXME:
// FromZip tries to extract a boot configuration from a ZIP file after verifying
// its signature with the provided public key file. The signature is expected to
// be appended to the ZIP file and have fixed length `ed25519.SignatureSize` .
// The returned string argument is the temporary directory where the files were
// extracted, if successful.
// No decoder (e.g. JSON, ZIP) or other function parsing the input file is called
// before verifying the signature.
func FromZip(filename string) (*Manifest, string, error) {
	// load the whole zip file in memory - we need it anyway for the signature
	// matching.
	// TODO refuse to read if too big?
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, "", err
	}
	urootcrypto.TryMeasureData(urootcrypto.BlobPCR, data, filename)
	zipbytes := data

	r, err := zip.NewReader(&memoryZipReader{Content: zipbytes}, int64(len(zipbytes)))
	if err != nil {
		return nil, "", err
	}
	tempDir, err := ioutil.TempDir(os.TempDir(), "bootconfig")
	if err != nil {
		return nil, "", err
	}
	log.Printf("Created temporary directory %s", tempDir)
	var manifest *Manifest
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			// Dont care - will be handled later
			continue
		}

		destination := path.Join(tempDir, f.Name)
		if len(f.Name) == 0 {
			log.Printf("Warning: skipping zero-length file name (flags: %d, mode: %s)", f.Flags, f.Mode())
			continue
		}
		// Check if folder exists
		if _, err := os.Stat(destination); os.IsNotExist(err) {
			if err := os.MkdirAll(path.Dir(destination), os.ModeDir|os.FileMode(0700)); err != nil {
				return nil, "", err
			}
		}
		fd, err := f.Open()
		if err != nil {
			return nil, "", err
		}
		buf, err := ioutil.ReadAll(fd)
		if err != nil {
			return nil, "", err
		}
		if f.Name == "manifest.json" {
			// make sure it's not a duplicate manifest within the ZIP file
			// and inform the user otherwise
			if manifest != nil {
				log.Printf("Warning: duplicate manifest.json found, the last found wins")
			}
			// parse the Manifest containing the boot configurations
			manifest, err = ManifestFromBytes(buf)
			if err != nil {
				return nil, "", err
			}
		}
		if err := ioutil.WriteFile(destination, buf, f.Mode()); err != nil {
			return nil, "", err
		}
		log.Printf("Extracted file %s (flags: %d, mode: %s)", f.Name, f.Flags, f.Mode())
	}
	if manifest == nil {
		return nil, "", errors.New("No manifest found")
	}
	return manifest, tempDir, nil
}

// FIXME:
// ToZip tries to pack all files specified in the the provided manifest.json
// into a zip archive. An error is returned, if the files (kernel, initrd, etc)
// doesn't exist at the paths written inside the manifest.json relative to its
// location. Optionally , if privkeyfile is not nil, after creating the archive an ed25519 signature is added to the
// archive. A copy of manifest.json is included in the final with adopted path of the
// bootfiles to mach teir location relative to the archive root.
func ToZip(output string, manifest string) error {
	// Get manifest from file. Make sure the file is named accordingliy, since
	// FromZip will search 'manifest.json' while extraction.
	if _, err := os.Stat(manifest); os.IsNotExist(err) {
		return fmt.Errorf("manifest.json not existent: %v", err)
	}
	if base := path.Base(manifest); base != "manifest.json" {
		return fmt.Errorf("invalid manifest name. Want 'manifest.json', got: %s", base)
	}
	manifestBody, err := ioutil.ReadFile(manifest)
	if err != nil {
		return err
	}
	mf, err := ManifestFromBytes(manifestBody)
	if err != nil {
		return fmt.Errorf("error parsing manifest: %v", err)
	} else if !mf.IsValid() {
		return errors.New("manifest is not valid")
	}

	// Create a buffer to write the archive to.
	buf := new(bytes.Buffer)
	// Create a new zip archive.
	z := zip.NewWriter(buf)

	var dest, origin string
	//Archive boot files
	for i, cfg := range mf.Configs {
		dir := fmt.Sprintf("bootconfig_%d/", i)
		z.Create(dir)
		if cfg.Kernel != "" {
			dest = path.Join(dir, path.Base(cfg.Kernel))
			origin = path.Join(path.Dir(manifest), cfg.Kernel)
			if err := toZip(z, dest, origin); err != nil {
				return fmt.Errorf("cant pack kernel: %v", err)
			}
			cfg.Kernel = dest
		}
		if cfg.Initramfs != "" {
			dest = path.Join(dir, path.Base(cfg.Initramfs))
			origin = path.Join(path.Dir(manifest), cfg.Initramfs)
			if err := toZip(z, dest, origin); err != nil {
				return fmt.Errorf("cant pack initramfs: %v", err)
			}
			cfg.Initramfs = dest
		}
		if cfg.DeviceTree != "" {
			dest = path.Join(dir, path.Base(cfg.DeviceTree))
			origin = path.Join(path.Dir(manifest), cfg.DeviceTree)
			if err := toZip(z, dest, origin); err != nil {
				return fmt.Errorf("cant pack device tree: %v", err)
			}
			cfg.DeviceTree = dest
		}
		mf.Configs[i] = cfg
	}

	// Archive root certificate
	z.Create("certs/")
	dest = "certs/root.cert"
	origin = path.Join(path.Dir(manifest), mf.RootCertPath)
	if err = toZip(z, dest, origin); err != nil {
		return fmt.Errorf("cant pack certificate: %v", err)
	}
	mf.RootCertPath = dest

	// Archive manifest
	newManifest, err := ManifestToBytes(mf)
	if err != nil {
		return err
	}
	dst, err := z.Create(path.Base(manifest))
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, bytes.NewReader(newManifest))
	if err != nil {
		return err
	}

	// Write central directory of archive
	err = z.Close()
	if err != nil {
		return err
	}

	// Write buf to disk
	if path.Ext(output) != ".zip" {
		output = output + ".zip"
	}
	err = ioutil.WriteFile(output, buf.Bytes(), 0777)
	if err != nil {
		return err
	}
	return nil
}

func toZip(w *zip.Writer, newPath, originPath string) error {
	dst, err := w.Create(newPath)
	if err != nil {
		return err
	}
	// Copy content from inputpath to new file
	src, err := os.Open(originPath)
	if err != nil {
		return fmt.Errorf("cannot find %s specified in manifest", originPath)
	}
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return src.Close()
}

// AddSignature signes the bootfiles inside an stboot.zip and inserts the
// signatures into the archive along with the respective certificate
func AddSignature(archive, privKey, certificate string) error {

	mf, dir, err := FromZip(archive)
	if err != nil {
		return err
	}

	// collect boot binaries
	// XXX Refactor if we remove bootconfig from manifest
	// Maybe just walk through certs/folders and match do root/bootconfig
	for i := range mf.Configs {

		bootconfigDir := path.Join(dir, fmt.Sprintf("bootconfig_%d", i))

		bcHash, err := HashBootconfigDir(bootconfigDir)
		if err != nil {
			return fmt.Errorf("failed to hash bootconfig - Err %s", err)
		}

		// Sign hash with Key
		buff, err := ioutil.ReadFile(privKey)
		privPem, _ := pem.Decode(buff)
		rsaPrivKey, err := x509.ParsePKCS1PrivateKey(privPem.Bytes)
		if err != nil {
			return fmt.Errorf("Can't parse private key: %v", err)
		}
		if rsaPrivKey == nil {
			return fmt.Errorf("RSA Key is nil: %v", err)
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
		err = os.MkdirAll(path.Join(dir, fmt.Sprintf("certs/bootconfig_%d/", i)), os.ModeDir|os.FileMode(0700))
		if err != nil {
			return fmt.Errorf("Creating signatures directories in %s failed: %v", dir, err)
		}

		// Extract part of Public Key for identification
		certificateString, err := ioutil.ReadFile(certificate)
		if err != nil {
			return fmt.Errorf("failed to read certificate - Err %s", err)
		}

		cert, err := parseCertificate(certificateString)
		if err != nil {
			return fmt.Errorf("failed to parse certificate %s with %v", certificateString, err)
		}

		// Write signature to folder
		err = ioutil.WriteFile(path.Join(dir, fmt.Sprintf("certs/bootconfig_%d/%s.signature", i, fmt.Sprintf("%x", cert.PublicKey)[2:18])), signature, 0644)
		if err != nil {
			return fmt.Errorf("writing into %v failed with %v - Check permissions", dir, err)
		}

		// cp cert to folder
		err = ioutil.WriteFile(path.Join(dir, fmt.Sprintf("certs/bootconfig_%d/%s.cert", i, fmt.Sprintf("%x", cert.PublicKey)[2:18])), certificateString, 0644)
		if err != nil {
			return fmt.Errorf("copying certificate %v to .zip failed with: %v - Check permissions", certificate, err)
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
			err := toZip(z, strings.Replace(path, dir, "", -1)[1:], path)
			if err != nil {
				return fmt.Errorf(fmt.Sprintf("Error adding file %s to .zip archive again", strings.Replace(path, dir, "", -1)))
			}
		}

		return nil
	})

	z.Close()

	pathToZip := fmt.Sprintf("./.original/%d", time.Now().Unix())
	os.MkdirAll(pathToZip, os.ModePerm)
	os.Rename(archive, pathToZip+"/stboot.zip")

	err = ioutil.WriteFile(archive, buf.Bytes(), 0777)
	if err != nil {
		return fmt.Errorf("unable to write new stboot.zip file - recover old from %s", pathToZip)
	}
	os.RemoveAll(pathToZip)

	return nil

}

// parseCertificate parses certificate from raw certificate
func parseCertificate(rawCertificate []byte) (x509.Certificate, error) {

	block, _ := pem.Decode(rawCertificate)
	if block == nil {
		log.Fatal("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return *pub, fmt.Errorf("failed to parse DER encoded public key: ", err)
	}

	return *pub, nil
}

// HashBootconfigDir hashes every file inside bootconigDir and returns a
// SHA512 hash
func HashBootconfigDir(bootconfigDir string) ([]byte, error) {

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
