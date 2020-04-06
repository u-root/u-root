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
	create           = kingpin.Command("create", "Create a bootball from the provided stconfig")
	createHWAddr     = create.Flag("mac", "Hardware address of the host if the created stboot.ball needs to be individual for a specific host.").String()
	createConfigFile = create.Arg("stconfig", "Path to the manifest file in JSON format").Required().ExistingFile()

	sign            = kingpin.Command("sign", "Sign the binary inside the provided bootball")
	signPrivKeyFile = sign.Flag("key", "Private key for signing").Required().ExistingFile()
	signCertFile    = sign.Flag("cert", "Certificate corresponding to the private key").Required().ExistingFile()
	signInFile      = sign.Arg("bootball", "Archive created by 'stconfig create'").Required().ExistingFile()

	unpack       = kingpin.Command("unpack", "Unpack boot ball  file into directory")
	unpackInFile = unpack.Arg("bootball", "Archive containing the boot files").Required().ExistingFile()
)

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(goversion).Author(Author)
	kingpin.CommandLine.Help = HelpText

	switch kingpin.Parse() {
	case create.FullCommand():
		if err := packBootBall(*createConfigFile, *createHWAddr); err != nil {
			log.Fatalln(err.Error())
		}
	case sign.FullCommand():
		if err := addSignatureToBootBall(*signInFile, *signPrivKeyFile, *signCertFile); err != nil {
			log.Fatalln(err.Error())
		}
	case unpack.FullCommand():
		if err := unpackBootBall(*unpackInFile); err != nil {
			log.Fatalln(err.Error())
		}
	default:
		log.Fatal("Command not found")
	}
}
