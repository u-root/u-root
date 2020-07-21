// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vfile

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/errors"
)

type signedFile struct {
	signers []*openpgp.Entity
	content string
}

func (s signedFile) write(path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write([]byte(s.content)); err != nil {
		return err
	}

	sigf, err := os.OpenFile(fmt.Sprintf("%s.sig", path), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer sigf.Close()
	for _, signer := range s.signers {
		if err := openpgp.DetachSign(sigf, signer, strings.NewReader(s.content), nil); err != nil {
			return err
		}
	}
	return nil
}

type normalFile struct {
	content string
}

func (n normalFile) write(path string) error {
	return ioutil.WriteFile(path, []byte(n.content), 0600)
}

func writeHashedFile(path, content string) ([]byte, error) {
	c := []byte(content)
	if err := ioutil.WriteFile(path, c, 0600); err != nil {
		return nil, err
	}
	hash := sha256.Sum256(c)
	return hash[:], nil
}

func TestOpenSignedFile(t *testing.T) {
	key, err := openpgp.NewEntity("goog", "goog", "goog@goog", nil)
	if err != nil {
		t.Fatal(err)
	}
	ring := openpgp.EntityList{key}

	key2, err := openpgp.NewEntity("goog2", "goog2", "goog@goog", nil)
	if err != nil {
		t.Fatal(err)
	}

	dir, err := ioutil.TempDir("", "opensignedfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	signed := signedFile{
		signers: openpgp.EntityList{key},
		content: "foo",
	}
	signedPath := filepath.Join(dir, "signed_by_key")
	if err := signed.write(signedPath); err != nil {
		t.Fatal(err)
	}

	signed2 := signedFile{
		signers: openpgp.EntityList{key2},
		content: "foo",
	}
	signed2Path := filepath.Join(dir, "signed_by_key2")
	if err := signed2.write(signed2Path); err != nil {
		t.Fatal(err)
	}

	signed12 := signedFile{
		signers: openpgp.EntityList{key, key2},
		content: "foo",
	}
	signed12Path := filepath.Join(dir, "signed_by_both.sig")
	if err := signed12.write(signed12Path); err != nil {
		t.Fatal(err)
	}

	normalPath := filepath.Join(dir, "unsigned")
	if err := ioutil.WriteFile(normalPath, []byte("foo"), 0777); err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		desc             string
		path             string
		keyring          openpgp.KeyRing
		want             error
		isSignatureValid bool
	}{
		{
			desc:             "signed file",
			keyring:          ring,
			path:             signedPath,
			want:             nil,
			isSignatureValid: true,
		},
		{
			desc:             "signed file w/ two signatures (key1 ring)",
			keyring:          ring,
			path:             signed12Path,
			want:             nil,
			isSignatureValid: true,
		},
		{
			desc:             "signed file w/ two signatures (key2 ring)",
			keyring:          openpgp.EntityList{key2},
			path:             signed12Path,
			want:             nil,
			isSignatureValid: true,
		},
		{
			desc:    "nil keyring",
			keyring: nil,
			path:    signed2Path,
			want: ErrUnsigned{
				Path: signed2Path,
				Err:  ErrNoKeyRing,
			},
			isSignatureValid: false,
		},
		{
			desc:    "non-nil empty keyring",
			keyring: openpgp.EntityList{},
			path:    signed2Path,
			want: ErrUnsigned{
				Path: signed2Path,
				Err:  errors.ErrUnknownIssuer,
			},
			isSignatureValid: false,
		},
		{
			desc:    "signed file does not match keyring",
			keyring: openpgp.EntityList{key2},
			path:    signedPath,
			want: ErrUnsigned{
				Path: signedPath,
				Err:  errors.ErrUnknownIssuer,
			},
			isSignatureValid: false,
		},
		{
			desc:    "unsigned file",
			keyring: ring,
			path:    normalPath,
			want: ErrUnsigned{
				Path: normalPath,
				Err: &os.PathError{
					Op:   "open",
					Path: fmt.Sprintf("%s.sig", normalPath),
					Err:  syscall.ENOENT,
				},
			},
			isSignatureValid: false,
		},
		{
			desc:    "file does not exist",
			keyring: ring,
			path:    filepath.Join(dir, "foo"),
			want: &os.PathError{
				Op:   "open",
				Path: filepath.Join(dir, "foo"),
				Err:  syscall.ENOENT,
			},
			isSignatureValid: false,
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			f, gotErr := OpenSignedSigFile(tt.keyring, tt.path)
			if !reflect.DeepEqual(gotErr, tt.want) {
				t.Errorf("openSignedFile(%v, %q) = %v, want %v", tt.keyring, tt.path, gotErr, tt.want)
			}

			if isSignatureValid := (gotErr == nil); isSignatureValid != tt.isSignatureValid {
				t.Errorf("isSignatureValid(%v) = %v, want %v", gotErr, isSignatureValid, tt.isSignatureValid)
			}

			// Make sure that the file is readable from position 0.
			if f != nil {
				content, err := ioutil.ReadAll(f)
				if err != nil {
					t.Errorf("Could not read content: %v", err)
				}
				if got := string(content); got != "foo" {
					t.Errorf("ReadAll = %v, want \"foo\"", got)
				}
			}
		})
	}
}

func TestOpenHashedFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "openhashedfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	hashedPath := filepath.Join(dir, "hashed")
	hash, err := writeHashedFile(hashedPath, "foo")
	if err != nil {
		t.Fatal(err)
	}

	emptyPath := filepath.Join(dir, "empty")
	emptyHash, err := writeHashedFile(emptyPath, "")
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		desc        string
		path        string
		hash        []byte
		want        error
		isHashValid bool
		wantContent string
	}{
		{
			desc:        "correct hash",
			path:        hashedPath,
			hash:        hash,
			want:        nil,
			isHashValid: true,
			wantContent: "foo",
		},
		{
			desc: "wrong hash",
			path: hashedPath,
			hash: []byte{0x99, 0x77},
			want: ErrInvalidHash{
				Path: hashedPath,
				Err: ErrHashMismatch{
					Got:  hash,
					Want: []byte{0x99, 0x77},
				},
			},
			isHashValid: false,
			wantContent: "foo",
		},
		{
			desc: "no hash",
			path: hashedPath,
			hash: []byte{},
			want: ErrInvalidHash{
				Path: hashedPath,
				Err:  ErrNoExpectedHash,
			},
			isHashValid: false,
			wantContent: "foo",
		},
		{
			desc:        "empty file",
			path:        emptyPath,
			hash:        emptyHash,
			want:        nil,
			isHashValid: true,
			wantContent: "",
		},
		{
			desc: "nonexistent file",
			path: filepath.Join(dir, "doesnotexist"),
			hash: nil,
			want: &os.PathError{
				Op:   "open",
				Path: filepath.Join(dir, "doesnotexist"),
				Err:  syscall.ENOENT,
			},
			isHashValid: false,
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			f, err := OpenHashedFile256(tt.path, tt.hash)
			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("OpenHashedFile256(%s, %x) = %v, want %v", tt.path, tt.hash, err, tt.want)
			}

			if isHashValid := (err == nil); isHashValid != tt.isHashValid {
				t.Errorf("isHashValid(%v) = %v, want %v", err, isHashValid, tt.isHashValid)
			}

			// Make sure that the file is readable from position 0.
			if f != nil {
				content, err := ioutil.ReadAll(f)
				if err != nil {
					t.Errorf("Could not read content: %v", err)
				}
				if got := string(content); got != tt.wantContent {
					t.Errorf("ReadAll = %v, want %s", got, tt.wantContent)
				}
			}
		})
	}
}
