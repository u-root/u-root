// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Performs signature checks on FDT images.
// Currently supports PGP and raw PKCS1v15 RSA signatures.
//
// Expected FDT Format:
//  Node: images
//  Node: image_name
//   P: data
//   Node: signature*
//    P: value
//    P: algo          (ex. 'sha256,rsa4096', 'pgp')
//    P: signer-name   (Optional)
//    P: key-name-hint (Optional)

package fit

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"fmt"
	"strings"
	"unicode"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/u-root/u-root/pkg/dt"
	"github.com/u-root/u-root/pkg/vfile"
)

var algs = map[string]crypto.Hash{
	"MD4":       crypto.MD4,
	"MD5":       crypto.MD5,
	"SHA1":      crypto.SHA1,
	"SHA224":    crypto.SHA224,
	"SHA256":    crypto.SHA256,
	"SHA384":    crypto.SHA384,
	"SHA512":    crypto.SHA512,
	"RIPEMD160": crypto.RIPEMD160,
	"SHA3_224":  crypto.SHA3_224,
	"SHA3_256":  crypto.SHA3_256,
	"SHA3_384":  crypto.SHA3_384,
	"SHA3_512":  crypto.SHA3_512,
}

// Signature defines an extendable interface for verifying images using
// varying signing methods.
type Signature interface {
	fmt.Stringer
	// Warning: If the signature does not exist or does not match the keyring,
	// both the file and a signature error will be returned.
	// Returns a bytes.Reader to the original data array.
	Verify([]byte, openpgp.KeyRing) (*bytes.Reader, error)
}

// PGPSignature implements a OpenPGP signature check.
type PGPSignature struct {
	name  string // Name of signature Node
	value []byte
	// Optional description fields
	signer string
	hint   string
}

// RSASignature implements a PKCS1v15 signature check.
type RSASignature struct {
	name  string // Name of signature Node
	hash  crypto.Hash
	value []byte
	// Optional description fields
	signer string
	hint   string
}

func (s PGPSignature) String() string {
	return fmt.Sprintf("PGP Signature - name: %s, signer: '%s', hint: '%s'", s.name, s.signer, s.hint)
}

// Verify runs a PKCS1v15 check using the RSA keys extracted from the provided
// key ring.
// Warning: If the signature does not exist or does not match the keyring,
// both the file and a signature error will be returned.
func (s PGPSignature) Verify(b []byte, ring openpgp.KeyRing) (*bytes.Reader, error) {
	r := bytes.NewReader(b)
	if signer, err := openpgp.CheckDetachedSignature(ring, bytes.NewReader(b), bytes.NewReader(s.value), nil); err != nil {
		return r, vfile.ErrUnsigned{Path: s.name, Err: err}
	} else if signer == nil {
		return r, vfile.ErrUnsigned{Path: s.name, Err: vfile.ErrWrongSigner{ring}}
	}
	return r, nil
}

// Verify runs a OpenPGP check using the PGP keys extracted from the provided
// key ring.
// Warning: If the signature does not exist or does not match the keyring,
// both the file and a signature error will be returned.
func (s RSASignature) Verify(b []byte, ring openpgp.KeyRing) (*bytes.Reader, error) {
	r := bytes.NewReader(b)
	keys, err := vfile.GetRSAKeysFromRing(ring)
	if err != nil {
		return r, err
	}

	hashed, err := vfile.CalculateHash(bytes.NewReader(b), s.hash.New())
	if err != nil {
		return r, err
	}

	for _, key := range keys {
		if err = rsa.VerifyPKCS1v15(key, s.hash, hashed, s.value); err == nil {
			return r, nil
		}
	}
	return r, vfile.ErrUnsigned{Err: vfile.ErrWrongSigner{ring}}
}

func (s RSASignature) String() string {
	return fmt.Sprintf("RSA Signature - name: %s, signer: '%s', hint: '%s', hash: '%s'", s.name, s.signer, s.hint, s.hash)
}

// parseHash cleans and maps the first detected hash string into a crypto.Hash.
// Expected format: "sha256,rsa4096" or "sha1"
func parseHash(algo string) (crypto.Hash, error) {
	cleaned := strings.TrimFunc(algo, func(r rune) bool {
		return !unicode.IsGraphic(r)
	})
	algoSplit := strings.Split(cleaned, ",")
	if len(algoSplit) == 0 {
		return 0, fmt.Errorf("unrecognized hash algo: '%s'", cleaned)
	}
	for _, alg := range algoSplit {
		if matched, ok := algs[strings.ToUpper(alg)]; ok {
			return matched, nil
		}
	}
	return 0, fmt.Errorf("unrecognized hash algo: '%s'", cleaned)
}

// parseSignatures parses dt.Node to RSASignatures
// Nodes with missing required properties or invalid hash functions are skipped.
// An error is returned if no valid signatures are found
func parseSignatures(n ...*dt.Node) ([]Signature, error) {
	var sigs []Signature
	for _, node := range n {
		v, ok := node.LookProperty("value")
		if !ok {
			fmt.Printf("Skipping signature %s: missing value node", node.Name)
			continue
		}

		a, ok := node.LookProperty("algo")
		if !ok {
			fmt.Printf("Skipping signature %s: missing algo node", node.Name)
			continue
		}

		var signer, hint string
		if signerProp, ok := node.LookProperty("signer-name"); ok {
			signer = string(signerProp.Value)
		}
		if hintProp, ok := node.LookProperty("key-name-hint"); ok {
			hint = string(hintProp.Value)
		}

		// Perform a broad stroke check for algos. RSA is assumed a raw RSA signature
		switch {
		case strings.Contains(string(a.Value), "pgp"):
			sigs = append(sigs, PGPSignature{value: v.Value, signer: signer, hint: hint})
		case strings.Contains(string(a.Value), "rsa"):
			// Parse the hash function used for the signature. ex. 'sha256,rsa4096'
			hf, err := parseHash(string(a.Value))
			if err != nil {
				fmt.Printf("Skipping signature %s: %v", node.Name, err)
				continue
			}
			sigs = append(sigs, RSASignature{value: v.Value, hash: hf, signer: signer, hint: hint})
		}
	}
	if len(sigs) == 0 {
		return nil, fmt.Errorf("failed to parse any valid Signatures")
	}
	return sigs, nil
}

// ReadSignedImage reads an image node from an FDT and verifies the
// content against a key set. Signature information is found in child nodes.
//
// WARNING! Unlike many Go functions, this may return both the file and an
// error.
//
// If the signature does not exist or does not match the keyring, both the file
// and a signature error will be returned.
func (i *Image) ReadSignedImage(image string, ring openpgp.KeyRing) (*bytes.Reader, error) {
	iroot := i.Root.Root().Walk("images").Walk(image)
	b, err := iroot.Property("data").AsBytes()
	if err != nil {
		return nil, err
	}

	br := bytes.NewReader(b)
	sigNodes, err := iroot.FindAll(func(n *dt.Node) bool {
		return strings.HasPrefix(strings.ToLower(n.Name), "signature")
	})
	if err != nil {
		return br, vfile.ErrUnsigned{Path: image, Err: fmt.Errorf("no signature nodes found")}
	}
	sigs, err := parseSignatures(sigNodes...)
	if err != nil {
		return br, vfile.ErrUnsigned{Path: image, Err: err}
	}

	for _, sig := range sigs {
		v, err := sig.Verify(b, ring)
		if err == nil {
			return v, nil
		}
		fmt.Printf("Ignoring failed signature - %s: Failed with %v\n", sig, err)
	}

	return br, vfile.ErrUnsigned{Path: image, Err: vfile.ErrWrongSigner{i.KeyRing}}
}
