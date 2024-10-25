// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vfile verifies files against a hash or signature.
//
// vfile is not TOCTTOU-safe against the contents of the file changing.
package vfile

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
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
		return nil, fmt.Errorf("could not open pub key: %w", err)
	}
	defer key.Close()

	ring, err := openpgp.ReadKeyRing(key)
	if err != nil {
		return nil, fmt.Errorf("could not read pub key: %w", err)
	}
	return ring, nil
}

// GetRSAKeysFromRing iterates a PGP Keyring and extracts all rsa.PublicKey.
// An error is returned iff the keyring is not found or no RSA public keys were
// found on it.
func GetRSAKeysFromRing(ring openpgp.KeyRing) ([]*rsa.PublicKey, error) {
	el, ok := ring.(openpgp.EntityList)
	if !ok {
		return nil, fmt.Errorf("failed to assert KeyRing as EntityList to read RSA keys")
	}

	var rsaKeys []*rsa.PublicKey
	for _, entity := range el {
		// Extract Primary Key
		if entity.PrimaryKey != nil {
			pk := (packet.PublicKey)(*entity.PrimaryKey)
			if rsaKey, ok := pk.PublicKey.(*rsa.PublicKey); ok {
				rsaKeys = append(rsaKeys, rsaKey)
			}
		}
		// Extract any subkeys
		for _, subkey := range entity.Subkeys {
			pk := (packet.PublicKey)(*subkey.PublicKey)
			if rsaKey, ok := pk.PublicKey.(*rsa.PublicKey); ok {
				rsaKeys = append(rsaKeys, rsaKey)
			}
		}
	}

	if len(rsaKeys) == 0 {
		return nil, fmt.Errorf("no RSA public keys found on keyring")
	}
	return rsaKeys, nil
}

// OpenSignedSigFile calls OpenSignedFile expecting the signature to be in path.sig.
//
// E.g. if path is /foo/bar, the signature is expected to be in /foo/bar.sig.
func OpenSignedSigFile(keyring openpgp.KeyRing, path string, opts ...OpenSignedFileOption) (*os.File, error) {
	return OpenSignedFile(keyring, path, fmt.Sprintf("%s.sig", path), opts...)
}

// OpenSignedFileOption is an optional argument to OpenSignedFile.
type OpenSignedFileOption func(*openSignedFileOptions)

type openSignedFileOptions struct {
	ignoreTimeConflict bool
}

func WithIgnoreTimeConflict(o *openSignedFileOptions) {
	o.ignoreTimeConflict = true
}

func getEndOfTime() time.Time {
	// number of seconds between Year 1 and 1970 (62135596800 seconds)
	unixToInternal := int64((1969*365 + 1969/4 - 1969/100 + 1969/400) * 24 * 60 * 60)

	// time.Unix adds unixToInternal seconds, subtract them to avoid
	// integer overflow in the internal representation.
	return time.Unix(1<<63-1-unixToInternal, 0)
}

// OpenSignedFile opens a file that is expected to be signed.
//
// WARNING! Unlike many Go functions, this may return both the file and an
// error.
//
// It expects pathSig to be available.
//
// If the signature does not exist or does not match the keyring, both the file
// and a signature error will be returned.
func OpenSignedFile(keyring openpgp.KeyRing, path, pathSig string, opts ...OpenSignedFileOption) (*os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var o openSignedFileOptions
	// Apply options if given.
	for _, opt := range opts {
		opt(&o)
	}

	signaturef, err := os.Open(pathSig)
	if err != nil {
		return f, ErrUnsigned{Path: path, Err: err}
	}
	defer signaturef.Close()

	var config packet.Config
	if o.ignoreTimeConflict {
		config.Time = getEndOfTime
	}

	// After CheckDetachedSignature reads the whole file, seek back to the beginning.
	defer f.Seek(0, io.SeekStart)

	if keyring == nil {
		return f, ErrUnsigned{Path: path, Err: ErrNoKeyRing}
	} else if signer, err := openpgp.CheckDetachedSignature(keyring, f, signaturef, &config); err != nil {
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
// If the contents match, the opened file is returned with no error.
func OpenHashedFile256(path string, wantSHA256Hash []byte) (*os.File, error) {
	return openHashedFile(path, wantSHA256Hash, sha256.New())
}

// OpenHashedFile512 opens path and verifies whether its contents match the
// given sha512 hash.
//
// WARNING! Unlike many Go functions, this may return both the file and an
// error in case the expected hash does not match the contents.
//
// If the contents match, the opened file is returned with no error.
func OpenHashedFile512(path string, wantSHA512Hash []byte) (*os.File, error) {
	return openHashedFile(path, wantSHA512Hash, sha512.New())
}

func openHashedFile(path string, wantHash []byte, h hash.Hash) (*os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if len(wantHash) == 0 {
		return f, ErrInvalidHash{
			Path: path,
			Err:  ErrNoExpectedHash,
		}
	}

	// After io.Copy reads the whole file, Seek back to beginning.
	defer f.Seek(0, io.SeekStart)

	// Hash the file.
	if _, err := io.Copy(h, f); err != nil {
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

// CheckHashedContent verifies a calculated hash against an expected hash array.
//
// WARNING! Unlike many Go functions, this may return both the file and an
// error in case the expected hash does not match the contents.
//
// If the contents match, the contents are returned with no error.
func CheckHashedContent(b *bytes.Reader, wantHash []byte, h hash.Hash) (*bytes.Reader, error) {
	if len(wantHash) == 0 {
		return b, ErrInvalidHash{
			Err: ErrNoExpectedHash,
		}
	}

	got, err := CalculateHash(b, h)
	if err != nil {
		return b, err
	}

	if subtle.ConstantTimeCompare(wantHash, got) == 0 {
		return b, ErrInvalidHash{
			Err: ErrHashMismatch{
				Got:  got,
				Want: wantHash,
			},
		}
	}
	return b, nil
}

// CalculateHash computes the hash of the input data b given a hash function.
func CalculateHash(b *bytes.Reader, h hash.Hash) ([]byte, error) {
	// Hash the file.
	if _, err := io.Copy(h, b); err != nil {
		return nil, ErrInvalidHash{
			Err: err,
		}
	}

	return h.Sum(nil), nil
}
