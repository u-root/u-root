// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//This program validates a file by verifying a checksum file and a signature file
//Exit status: 0-OK, 1-Any error, 2-Bad signature, 3-Bad checksum
package main

import (
	"bufio"
	"crypto/sha1"
	"flag"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func checksum(fi *os.File) {

	sha1file, err := os.Open(fi.Name() + ".sha1sum")
	if err != nil {
		log.Fatalf("Coudn't open sha1 file %s: %v", sha1file.Name(), err)
	}
	rd := bufio.NewReader(sha1file)
	defer sha1file.Close()
	var data []byte

	data, err = ioutil.ReadFile(fi.Name())
	if err != nil {
		log.Fatalf("Couldn't read file: %s", fi.Name())
	}

	line, err := rd.ReadString(' ')
	if err != nil {
		log.Fatalf("%v", err)
	}

	sum := ""
	for _, b := range sha1.Sum(data) {
		hex := strconv.FormatInt(int64(b), 16)
		//We need to preserve the leading zero
		if len(hex) == 1 {
			hex = "0" + hex
		}
		sum += hex
	}

	if strings.TrimRight(line, " ") != sum {
		fmt.Printf("%s and %s\n", strings.TrimRight(line, " "), sum)
		log.Printf("Bad checksum\n")
		os.Exit(3)
	}
}

func checksig(fi *os.File, pkfile string) {
	keyRingReader, err := os.Open(pkfile)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	signature, err := os.Open(fi.Name() + ".sig")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	keyring, err := openpgp.ReadArmoredKeyRing(keyRingReader)
	if err != nil {
		log.Fatalf("Read Armored Key Ring: %v", err)
		return
	}
	entity, err := openpgp.CheckDetachedSignature(keyring, fi, signature)
	if err != nil {
		log.Printf("Bad signature: %v", err)
		os.Exit(2)
	}

	fmt.Printf("Good signature by: %v\n", entity)
}

func main() {

	flag.Parse()
	files := flag.Args()
	fi, err := os.Open(files[0])
	if err != nil {
		log.Fatalf("Couldn't open %s: %v", files[0], err)
	}
	defer fi.Close()

	checksig(fi, files[1])
	checksum(fi)
	os.Exit(0)
}
