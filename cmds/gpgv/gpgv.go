// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gpgv validates a signature against a file.
//
// Synopsis:
//     gpgv [-v] KEY SIG CONTENT
//
// Description:
//     It prints "OK\n" to stdout if the check succeeds and exits with 0. It
//     prints an error message and exits with non-0 otherwise.
//
//     The openpgp package ReadKeyRing function does not completely implement
//     RFC4880 in that it can't use a PublicSigningKey with 0 signatures. We
//     use one from Eric Grosse instead.
//
// Options:
//     -v: verbose
//
// Author:
//     grosse@gmail.com.
package main

import (
	"crypto"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/u-root/u-root/pkg/log"

	"golang.org/x/crypto/openpgp/errors"
	"golang.org/x/crypto/openpgp/packet"
)

func main() {
	flag.Parse()
	if flag.NArg() != 3 {
		log.Fatalf("usage: boot-verify [-v] key sig content")
	}

	keyf, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatalf("%v", err)
	}
	sigf, err := os.Open(flag.Args()[1])
	if err != nil {
		log.Fatalf("%v", err)
	}
	contentf, err := os.Open(flag.Args()[2])
	if err != nil {
		log.Fatalf("%v", err)
	}

	key, err := readPublicSigningKey(keyf)
	if err != nil {
		log.Fatalf("key %v", err)
	}

	if err = verifyDetachedSignature(key, contentf, sigf); err != nil {
		log.Fatalf("verify: %v", err)
	}
	fmt.Printf("OK")
}

func readPublicSigningKey(keyf io.Reader) (*packet.PublicKey, error) {
	keypackets := packet.NewReader(keyf)
	p, err := keypackets.Next()
	if err != nil {
		return nil, err
	}
	switch pkt := p.(type) {
	case *packet.PublicKey:
		log.Printf("key: ", pkt)
		return pkt, nil
	default:
		log.Printf("ReadPublicSigningKey: got %T, want *packet.PublicKey", pkt)
	}
	return nil, errors.StructuralError("expected first packet to be PublicKey")
}

func verifyDetachedSignature(key *packet.PublicKey, contentf, sigf io.Reader) error {
	var hashFunc crypto.Hash

	packets := packet.NewReader(sigf)
	p, err := packets.Next()
	if err != nil {
		return err
	}
	switch sig := p.(type) {
	case *packet.Signature:
		hashFunc = sig.Hash
	case *packet.SignatureV3:
		hashFunc = sig.Hash
	default:
		return errors.UnsupportedError("unrecognized signature")
	}

	h := hashFunc.New()
	if _, err := io.Copy(h, contentf); err != nil && err != io.EOF {
		return err
	}
	switch sig := p.(type) {
	case *packet.Signature:
		err = key.VerifySignature(h, sig)
	case *packet.SignatureV3:
		err = key.VerifySignatureV3(h, sig)
	default:
		panic("unreachable")
	}
	return err
}
