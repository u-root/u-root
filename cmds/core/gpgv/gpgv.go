// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gpgv validates a signature against a file.
//
// Synopsis:
//
//	gpgv KEY SIG CONTENT
//
// Description:
//
//	It prints "OK\n" to stdout if the check succeeds and exits with 0. It
//	prints an error message and exits with non-0 otherwise.
//
//	The openpgp package ReadKeyRing function does not completely implement
//	RFC4880 in that it can't use a PublicSigningKey with 0 signatures. We
//	use one from Eric Grosse instead.
//
// Author:
//
//	grosse@gmail.com.
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

var (
	errUsage            = errors.New("usage: gpgv key sig content")
	errUnknownSignature = errors.New("unknown signature")
	errExpectedPacket   = errors.New("expected first packet to be PublicKey")
)

func readPublicSigningKey(keyf io.Reader) (*packet.PublicKey, error) {
	keypackets := packet.NewReader(keyf)
	p, err := keypackets.Next()
	if err != nil {
		return nil, err
	}

	pkt, ok := p.(*packet.PublicKey)
	if !ok {
		return nil, errExpectedPacket
	}
	return pkt, nil
}

func verifyDetachedSignature(key *packet.PublicKey, contentf, sigf io.Reader) error {
	packets := packet.NewReader(sigf)
	p, err := packets.Next()
	if err != nil {
		return err
	}

	sig, ok := p.(*packet.Signature)
	if !ok {
		return errUnknownSignature
	}

	h := sig.Hash.New()
	if _, err := io.Copy(h, contentf); err != nil && err != io.EOF {
		return err
	}

	return key.VerifySignature(h, sig)
}

func runGPGV(w io.Writer, keyfile, sigfile, datafile string) error {
	if keyfile == "" || sigfile == "" || datafile == "" {
		return errUsage
	}

	keyf, err := os.Open(keyfile)
	if err != nil {
		return err
	}
	defer keyf.Close()

	sigf, err := os.Open(sigfile)
	if err != nil {
		return err
	}
	defer sigf.Close()

	contentf, err := os.Open(datafile)
	if err != nil {
		return err
	}
	defer contentf.Close()

	key, err := readPublicSigningKey(keyf)
	if err != nil {
		return fmt.Errorf("key: %w ", err)
	}

	if err = verifyDetachedSignature(key, contentf, sigf); err != nil {
		return fmt.Errorf("verify: %w", err)
	}
	fmt.Fprintf(w, "OK")
	return nil
}

func main() {
	if len(os.Args) != 4 {
		log.Fatal(errUsage)
	}
	if err := runGPGV(os.Stdout, os.Args[1], os.Args[2], os.Args[3]); err != nil {
		log.Fatal(err)
	}
}
