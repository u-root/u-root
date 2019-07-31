// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bootconfig

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/u-root/u-root/pkg/crypto"
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
	crypto.TryMeasureData(crypto.BlobPCR, data, filename)
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
		destination := path.Join(tempDir, f.Name)
		if len(f.Name) == 0 {
			log.Printf("Warning: skipping zero-length file name (flags: %d, mode: %s)", f.Flags, f.Mode())
			continue
		}
		if f.Name[len(f.Name)-1] == '/' {
			// it's a directory, create it
			if err := os.MkdirAll(destination, os.ModeDir|os.FileMode(0700)); err != nil {
				return nil, "", err
			}
			log.Printf("Extracted directory %s (flags: %d)", f.Name, f.Flags)
		} else {
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
	}
	if manifest == nil {
		return nil, "", errors.New("No manifest found")
	}
	return manifest, tempDir, nil
}

// FIXME:
// ToZip tries to pack all files specifoed in the the provided manifest.json
// into a zip archive. An error is returned, if the files (kernel, initrd, etc)
// doesn't exist at the paths written inside the manifest.json relative to its
// location. Optionally , if privkeyfile is not nil, after creating the archive an ed25519 signature is added to the
// archive. A copy of manifest.json is included in the final with adopted path of the
// bootfiles to mach teir location relative to the archive root.
func ToZip(output string, manifest string) error {
	// Get manifest from file. Make sure the file is named accordingliy, since
	// FromZip will search 'manifest.json' while extraction.
	if base := path.Base(manifest); base != "manifest.json" {
		return fmt.Errorf("Invalid manifest name. Want 'manifest.json', got: %s", base)
	}
	manifestBody, err := ioutil.ReadFile(manifest)
	if err != nil {
		return err
	}
	mf, err := ManifestFromBytes(manifestBody)
	if err != nil {
		return err
	} else if !mf.IsValid() {
		return errors.New("Manifest is not valid")
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
			toZip(z, dest, origin)
			cfg.Kernel = dest
		}
		if cfg.Initramfs != "" {
			dest = path.Join(dir, path.Base(cfg.Initramfs))
			origin = path.Join(path.Dir(manifest), cfg.Initramfs)
			toZip(z, dest, origin)
			cfg.Initramfs = dest
		}
		if cfg.DeviceTree != "" {
			dest = path.Join(dir, path.Base(cfg.DeviceTree))
			origin = path.Join(path.Dir(manifest), cfg.DeviceTree)
			toZip(z, dest, origin)
			cfg.DeviceTree = dest
		}
		mf.Configs[i] = cfg
	}

	// Archive root certificate
	z.Create("certs/")
	dest = path.Join("certs/", path.Base(mf.RootCertPath))
	origin = path.Join(path.Dir(manifest), mf.RootCertPath)
	toZip(z, dest, origin)
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
		return fmt.Errorf("Cannot find %s specified in manifest", originPath)
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
	bins := make([]string, 0)
	for i, cfg := range mf.Configs {
		dir := fmt.Sprintf("bootconfig_%d", i)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, file := range files {
			if !file.IsDir() {
				bins = append(bins, file.Name())
			}
		}
	}
	// TODO:
	// read all kernel files in a byte buffer
	// sign it
	// save signature in dir/certs/signatureX
	// save certificate in dir/certs/certificateX
	// update manifest and write it to dir
	// pack everything back tp zip (ToZiP)

	// make use of it in /cmds/boot/stboot/main.go

}
