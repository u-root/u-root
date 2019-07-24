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
	HelpText = "A tool for managing System Transparency boot configurations"
)

var goversion string

var (
	genkeys = kingpin.Command("genkeys", "Generate RSA keypair")
	pack    = kingpin.Command("pack", "Create boot configuration file")
	unpack  = kingpin.Command("unpack", "Unpack boot configuration file into directory")

	genkeysPrivateKeyFile = genkeys.Arg("privateKey", "File path to write the private key").Required().String()
	genkeysPublicKeyFile  = genkeys.Arg("publicKey", "File path to write the public key").Required().String()
	genkeysPassphrase     = genkeys.Flag("passphrase", "Encrypt keypair in PKCS8 format").String()

	packSignPassphrase     = pack.Flag("passphrase", "Passphrase for private key file").String()
	packManifest           = pack.Arg("manifest", "Path to the manifest file in JSON format").Required().String()
	packOutputFilename     = pack.Arg("bc-file", "Path to output file").Required().String()
	packSignPrivateKeyFile = pack.Arg("private-key", "Path to the private key file").Required().String()

	unpackInputFilename       = unpack.Arg("bc-file", "Boot configuration file").Required().String()
	unpackVerifyPublicKeyFile = unpack.Arg("public-key", "Path to the public key file").String()
)

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(goversion).Author(Author)
	kingpin.CommandLine.Help = HelpText

	switch kingpin.Parse() {
	case "genkeys":
		if err := GenKeys(); err != nil {
			log.Fatalln(err.Error())
		}
	case "pack":
		if err := PackBootConfiguration(); err != nil {
			log.Fatalln(err.Error())
		}
	case "unpack":
		if err := UnpackBootConfiguration(); err != nil {
			log.Fatalln(err.Error())
		}
	default:
		log.Fatal("Command not found")
	}
}
