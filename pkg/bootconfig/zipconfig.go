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

var (
	signed    byte = 0xff
	notSigned byte = 0x11
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

	// Collect botfiles relative to manifest.json
	files := make([]string, 0)
	files = append(files, path.Base(manifest))
	for _, cfg := range mf.Configs {
		if cfg.Kernel != "" {
			files = append(files, cfg.Kernel)
		}
		if cfg.Initramfs != "" {
			files = append(files, cfg.Initramfs)
		}
		if cfg.DeviceTree != "" {
			files = append(files, cfg.DeviceTree)
		}
	}

	// Create a buffer to write the archive to.
	buf := new(bytes.Buffer)
	// Create a new zip archive.
	w := zip.NewWriter(buf)

	// Archive files
	dir := path.Dir(manifest)
	for _, file := range files {
		// Create directories of each filepath first
		for d := file; ; {
			d, _ = path.Split(d)
			if d != "" {
				w.Create(d)
				d = path.Clean(d)
			} else {
				break
			}
		}
		// Create new file in archive
		dst, err := w.Create(file)
		if err != nil {
			return err
		}
		// Copy content from inputpath to new file
		p := path.Join(dir, file)
		src, err := os.Open(p)
		if err != nil {
			return fmt.Errorf("Cannot find %s specified in %s", p, manifest)
		}
		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}
		err = src.Close()
		if err != nil {
			return err
		}
	}
	// Write central directory of archive
	err = w.Close()
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
