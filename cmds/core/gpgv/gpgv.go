// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gpgv validates a signature against a file.
//
// Synopsis:
//
//	gpgv [-v] KEY SIG CONTENT
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
// Options:
//
//	-v: verbose
//
// Author:
//
//	grosse@gmail.com.
package main

import (
	"crypto"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	gpgerror "github.com/ProtonMail/go-crypto/openpgp/errors"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

var (
	verbose  bool
	debug    = func(string, ...interface{}) {}
	errUsage = errors.New("usage: boot-verify [-v] key sig content")
)

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose")
}

func readPublicSigningKey(keyf io.Reader) (*packet.PublicKey, error) {
	keypackets := packet.NewReader(keyf)
	p, err := keypackets.Next()
	if err != nil {
		return nil, err
	}
	switch pkt := p.(type) {
	case *packet.PublicKey:
		debug("key: ", pkt)
		return pkt, nil
	default:
		log.Printf("ReadPublicSigningKey: got %T, want *packet.PublicKey", pkt)
	}
	return nil, gpgerror.StructuralError("expected first packet to be PublicKey")
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
	default:
		return gpgerror.UnsupportedError("unrecognized signature")
	}

	h := hashFunc.New()
	if _, err := io.Copy(h, contentf); err != nil && err != io.EOF {
		return err
	}
	switch sig := p.(type) {
	case *packet.Signature:
		err = key.VerifySignature(h, sig)
	default:
		return fmt.Errorf("unknown signature")
	}
	return err
}

func runGPGV(w io.Writer, verbose bool, keyfile, sigfile, datafile string) error {
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
	flag.Parse()
	if len(flag.Args()) != 3 {
		log.Fatal(errUsage)
	}
	if err := runGPGV(os.Stdout, verbose, flag.Args()[0], flag.Args()[1], flag.Args()[2]); err != nil {
		log.Fatal(err)
	}
}
