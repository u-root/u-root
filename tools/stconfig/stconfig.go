// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// https://xkcd.com/927/

// stconfig is a configuration tool to create and manage artifacts for
// System Transparency Boot. Artifacts are ment to be uploaded to a
// remote provisioning server.

import (
	"log"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	// Author is the author
	Author = "Philipp Deppenwiese, Jens Drenhaus"
	// HelpText is the command line help
	HelpText = "stconfig can be used for managing System Transparency boot configurations"
)

var goversion string

var (
	create = kingpin.Command("create", "Create a boot ball from stconfig.json")
	sign   = kingpin.Command("sign", "Sign the binary inside the provided stboot.ball and add the signatures and certificates")
	unpack = kingpin.Command("unpack", "Unpack boot ball  file into directory")

	createConfigFile = create.Arg("config", "Path to the manifest file in JSON format").Required().String()

	signInFile      = sign.Arg("bootball", "Archive created by 'stconfig create'").Required().String()
	signPrivKeyFile = sign.Arg("privkey", "Private key for signing").Required().String()
	signCertFile    = sign.Arg("certificate", "Certificate to veryfy the signature").Required().String()

	unpackInFile = unpack.Arg("bootball", "Archive containing the boot files").Required().String()
)

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(goversion).Author(Author)
	kingpin.CommandLine.Help = HelpText

	switch kingpin.Parse() {
	case create.FullCommand():
		if _, err := os.Stat(*createConfigFile); os.IsNotExist(err) {
			log.Fatalf("%s does not exist: %v", *createConfigFile, err)
		}
		if err := packBootBall(*createConfigFile); err != nil {
			log.Fatalln(err.Error())
		}
	case sign.FullCommand():
		if _, err := os.Stat(*signInFile); os.IsNotExist(err) {
			log.Fatalf("%s does not exist: %v", *signInFile, err)
		}
		if _, err := os.Stat(*signPrivKeyFile); os.IsNotExist(err) {
			log.Fatalf("%s does not exist: %v", *signPrivKeyFile, err)
		}
		if _, err := os.Stat(*signCertFile); os.IsNotExist(err) {
			log.Fatalf("%s does not exist: %v", *signCertFile, err)
		}
		if err := addSignatureToBootBall(*signInFile, *signPrivKeyFile, *signCertFile); err != nil {
			log.Fatalln(err.Error())
		}
	case unpack.FullCommand():
		if _, err := os.Stat(*unpackInFile); os.IsNotExist(err) {
			log.Fatalf("%s does not exist: %v", *signInFile, err)
		}
		if err := unpackBootBall(*unpackInFile); err != nil {
			log.Fatalln(err.Error())
		}
	default:
		log.Fatal("Command not found")
	}
}
