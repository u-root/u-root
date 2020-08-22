// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vfile verifies files against a hash or signature.
//
// vfile aims to be TOCTTOU-safe by reading files into memory before verifying.
package vfile

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/openpgp"
)

// ErrUnsigned is returned for a file that failed signature verification.
type ErrUnsigned struct {
	// Path is the file that failed signature verification.
	Path string

	// Err is a nested error, if there was one.
	Err error
}

func (e ErrUnsigned) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("file %q is unsigned: %v", e.Path, e.Err)
	}
	return fmt.Sprintf("file %q is unsigned", e.Path)
}

func (e ErrUnsigned) Unwrap() error {
	return e.Err
}

// ErrNoKeyRing is returned when a nil keyring was given.
var ErrNoKeyRing = errors.New("no keyring given")

// ErrWrongSigner represents a file signed by some key, but not the ones in the given key ring.
type ErrWrongSigner struct {
	// KeyRing is the expected key ring.
	KeyRing openpgp.KeyRing
}

func (e ErrWrongSigner) Error() string {
	return fmt.Sprintf("signed by a key not present in keyring %s", e.KeyRing)
}

// GetKeyRing returns an OpenPGP KeyRing loaded from the specified path.
//
// keyPath must be an already trusted path, e.g. keys are included in the initramfs.
func GetKeyRing(keyPath string) (openpgp.KeyRing, error) {
	key, err := os.Open(keyPath)
	if err != nil {
		return nil, fmt.Errorf("could not open pub key: %v", err)
	}
	defer key.Close()

	ring, err := openpgp.ReadKeyRing(key)
	if err != nil {
		return nil, fmt.Errorf("could not read pub key: %v", err)
	}
	return ring, nil
}

// OpenSignedSigFile calls OpenSignedFile expecting the signature to be in path.sig.
//
// E.g. if path is /foo/bar, the signature is expected to be in /foo/bar.sig.
func OpenSignedSigFile(keyring openpgp.KeyRing, path string) (*File, error) {
	return OpenSignedFile(keyring, path, fmt.Sprintf("%s.sig", path))
}

// File encapsulates a bytes.Reader with the file contents and its name.
type File struct {
	*bytes.Reader

	FileName string
}

// Name returns the file name.
func (f *File) Name() string {
	return f.FileName
}

// OpenSignedFile opens a file that is expected to be signed.
//
// WARNING! Unlike many Go functions, this may return both the file and an
// error.
//
// It expects path.sig to be available.
//
// If the signature does not exist or does not match the keyring, both the file
// and a signature error will be returned.
func OpenSignedFile(keyring openpgp.KeyRing, path, pathSig string) (*File, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f := &File{
		Reader:   bytes.NewReader(content),
		FileName: path,
	}

	signaturef, err := os.Open(pathSig)
	if err != nil {
		return f, ErrUnsigned{Path: path, Err: err}
	}
	defer signaturef.Close()

	if keyring == nil {
		return f, ErrUnsigned{Path: path, Err: ErrNoKeyRing}
	} else if signer, err := openpgp.CheckDetachedSignature(keyring, bytes.NewReader(content), signaturef); err != nil {
		return f, ErrUnsigned{Path: path, Err: err}
	} else if signer == nil {
		return f, ErrUnsigned{Path: path, Err: ErrWrongSigner{keyring}}
	}
	return f, nil
}

// ErrInvalidHash is returned when hash verification failed.
type ErrInvalidHash struct {
	// Path is the path to the file that was supposed to be verified.
	Path string

	// Err is some underlying error.
	Err error
}

func (e ErrInvalidHash) Error() string {
	return fmt.Sprintf("invalid hash for file %q: %v", e.Path, e.Err)
}

func (e ErrInvalidHash) Unwrap() error {
	return e.Err
}

// ErrHashMismatch is returned when the file's hash does not match the expected hash.
type ErrHashMismatch struct {
	Want []byte
	Got  []byte
}

func (e ErrHashMismatch) Error() string {
	return fmt.Sprintf("got hash %#x, expected %#x", e.Got, e.Want)
}

// ErrNoExpectedHash is given when the caller did not specify a hash.
var ErrNoExpectedHash = errors.New("OpenHashedFile: no expected hash given")

// OpenHashedFile256 opens path and verifies whether its contents match the
// given sha256 hash.
//
// WARNING! Unlike many Go functions, this may return both the file and an
// error in case the expected hash does not match the contents.
//
// If the contents match, the contents are returned with no error.
func OpenHashedFile256(path string, wantSHA256Hash []byte) (*File, error) {
	return openHashedFile(path, wantSHA256Hash, sha256.New())
}

// OpenHashedFile512 opens path and verifies whether its contents match the
// given sha512 hash.
//
// WARNING! Unlike many Go functions, this may return both the file and an
// error in case the expected hash does not match the contents.
//
// If the contents match, the contents are returned with no error.
func OpenHashedFile512(path string, wantSHA512Hash []byte) (*File, error) {
	return openHashedFile(path, wantSHA512Hash, sha512.New())
}

func openHashedFile(path string, wantHash []byte, h hash.Hash) (*File, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f := &File{
		Reader:   bytes.NewReader(content),
		FileName: path,
	}

	if len(wantHash) == 0 {
		return f, ErrInvalidHash{
			Path: path,
			Err:  ErrNoExpectedHash,
		}
	}

	// Hash the file.
	if _, err := io.Copy(h, bytes.NewReader(content)); err != nil {
		return f, ErrInvalidHash{
			Path: path,
			Err:  err,
		}
	}

	got := h.Sum(nil)
	if !bytes.Equal(wantHash, got) {
		return f, ErrInvalidHash{
			Path: path,
			Err: ErrHashMismatch{
				Got:  got,
				Want: wantHash,
			},
		}
	}
	return f, nil
}
