// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//This program validates a file by verifying a checksum file and a signature file
//Exit status: 0-OK, 1-Any error, 2-Bad signature, 3-Bad checksum
package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	armored = flag.Bool("a", false, "signature is ASCII armored")
	sumfile = flag.String("i", "", "checksum file")
	md      = flag.Bool("md5", false, "use md5sum")
	s1      = flag.Bool("sha1", false, "use sha1sum")
	s256    = flag.Bool("sha256", false, "use sha256")
	s512    = flag.Bool("sha512", false, "use sha512")
)

func checksum(fi *os.File) {
	var csfile *os.File
	var err error
	var hash string

	if *sumfile != "" {
		csfile, err = os.Open(*sumfile)
		if err != nil {
			log.Fatalf("Couldn't open file %s: %v", csfile, err)
		}
		defer csfile.Close()
		parts := strings.Split(csfile.Name(), ".")
		end := parts[len(parts)-1]
		switch {
		case *s256 || end == "sha256sum":
			hash = sha256sum(csfile, fi)
		case *s512 || end == "sha512sum":
			hash = sha512sum(csfile, fi)
		case *md || end == "md5sum":
			hash = md5sum(csfile, fi)
		case *s1 || end == "sha1sum":
			hash = sha1sum(csfile, fi)
		default:
			log.Printf("Couldn't identify checksum type, using default sha256sum")
			hash = sha256sum(csfile, fi)
		}
	} else {
		switch {
		case *s256:
			csfile, err = os.Open(fi.Name() + ".sha256sum")
			if err != nil {
				log.Fatalf("Couldn't open file %s: %v", csfile, err)
			}
			defer csfile.Close()
			hash = sha256sum(csfile, fi)
		case *s512:
			csfile, err = os.Open(fi.Name() + ".sha512sum")
			if err != nil {
				log.Fatalf("Couldn't open file %s: %v", csfile, err)
			}
			defer csfile.Close()
			hash = sha512sum(csfile, fi)
		case *md:
			csfile, err = os.Open(fi.Name() + ".md5sum")
			if err != nil {
				log.Fatalf("Couldn't open file %s: %v", csfile, err)
			}
			defer csfile.Close()
			hash = md5sum(csfile, fi)
		case *s1:
			csfile, err = os.Open(fi.Name() + ".sha1sum")
			if err != nil {
				log.Fatalf("Couldn't open file %s: %v", csfile, err)
			}
			defer csfile.Close()
			hash = sha1sum(csfile, fi)
		default:
			csfile, err = os.Open(fi.Name() + ".sha256sum")
			if err != nil {
				log.Fatalf("Couldn't open file %s: %v", csfile, err)
			}
			defer csfile.Close()
			hash = sha256sum(csfile, fi)
		}
	}

	rd := bufio.NewReader(csfile)

	line, err := rd.ReadString(' ')
	if err != nil {
		log.Fatalf("%v", err)
	}

	if strings.TrimRight(line, " ") != hash {
		log.Printf("Bad checksum\n")
		os.Exit(3)
	}
}

func md5sum(csfile, fi *os.File) string {
	var hash string

	data, err := ioutil.ReadFile(fi.Name())
	if err != nil {
		log.Fatalf("Couldn't read file: %s", fi.Name())
	}
	for _, b := range md5.Sum(data) {
		hex := strconv.FormatInt(int64(b), 16)
		if len(hex) == 1 {
			hex = "0" + hex
		}
		hash += hex
	}

	return hash

}

func sha1sum(csfile, fi *os.File) string {
	var hash string

	data, err := ioutil.ReadFile(fi.Name())
	if err != nil {
		log.Fatalf("Couldn't read file: %s", fi.Name())
	}
	for _, b := range sha1.Sum(data) {
		hex := strconv.FormatInt(int64(b), 16)
		if len(hex) == 1 {
			hex = "0" + hex
		}
		hash += hex
	}

	return hash
}

func sha256sum(csfile, fi *os.File) string {
	var hash string

	data, err := ioutil.ReadFile(fi.Name())
	if err != nil {
		log.Fatalf("Couldn't read file: %s", fi.Name())
	}
	for _, b := range sha256.Sum256(data) {
		hex := strconv.FormatInt(int64(b), 16)
		if len(hex) == 1 {
			hex = "0" + hex
		}
		hash += hex
	}

	return hash
}

func sha512sum(csfile, fi *os.File) string {
	var hash string

	data, err := ioutil.ReadFile(fi.Name())
	if err != nil {
		log.Fatalf("Couldn't read file: %s", fi.Name())
	}
	for _, b := range sha512.Sum512(data) {
		hex := strconv.FormatInt(int64(b), 16)
		if len(hex) == 1 {
			hex = "0" + hex
		}
		hash += hex
	}

	return hash
}

func checksig(fi *os.File, pkfiles []string) {
	var signature *os.File
	var keyring openpgp.EntityList
	var entity *openpgp.Entity
	var err error

	for _, pk := range pkfiles {
		keyRingReader, err := os.Open(pk)
		if err != nil {
			log.Fatalf("%v", err)
		}
		defer keyRingReader.Close()
		signature, err = os.Open(fi.Name() + ".sig")
		if err != nil {
			log.Fatalf("%v", err)
		}
		defer signature.Close()
		keyring, err = openpgp.ReadArmoredKeyRing(keyRingReader)
		if err != nil {
			log.Fatalf("Read Armored Key Ring: %v", err)
		}
	}

	if *armored {
		entity, err = openpgp.CheckArmoredDetachedSignature(keyring, fi, signature)
	} else {
		entity, err = openpgp.CheckDetachedSignature(keyring, fi, signature)
	}

	if err != nil {
		log.Printf("Bad signature: %v", signature.Name(), err)
		os.Exit(2)
	}

	fmt.Printf("Good signature by: %v\n", entity.Identities)
}

func main() {

	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatalf("Need at least a file to be validated and one public key")
	}
	files := flag.Args()
	fi, err := os.Open(files[0])
	if err != nil {
		log.Fatalf("Couldn't open %s: %v", files[0], err)
	}
	defer fi.Close()

	checksig(fi, files[1:])
	checksum(fi)
	os.Exit(0)
}
