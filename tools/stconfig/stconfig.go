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
	genkeys = kingpin.Command("genkeys", "Generate RSA keypair")
	create  = kingpin.Command("create", "Create boot configuration file from manifest.json")
	//(unpack  = kingpin.Command("unpack", "Unpack boot configuration file into directory")

	genkeysPrivateKeyFile = genkeys.Arg("privateKey", "File path to write the private key").Required().String()
	genkeysPublicKeyFile  = genkeys.Arg("publicKey", "File path to write the public key").Required().String()
	genkeysPassphrase     = genkeys.Flag("passphrase", "Encrypt keypair in PKCS8 format").String()

	createSignZipPrivKey    = create.Flag("zip-signing-key", "path tp private key to append additional signature to packed boot file").Default("").String()
	createSignZipPassphrase = create.Flag("passphrase", "Passphrase for private key file").String()
	createManifest          = create.Arg("manifest", "Path to the manifest file in JSON format").Required().String()
	createOutputFilename    = create.Arg("bc-file", "Path to output file").Required().String()

	// unpackInputFilename       = unpack.Arg("bc-file", "Boot configuration file").Required().String()
	// unpackVerifyPublicKeyFile = unpack.Arg("public-key", "Path to the public key file").String()
)

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(goversion).Author(Author)
	kingpin.CommandLine.Help = HelpText

	switch kingpin.Parse() {
	case genkeys.FullCommand():
		if err := GenKeys(); err != nil {
			log.Fatalln(err.Error())
		}
	case create.FullCommand():
		if err := PackBootConfiguration(); err != nil {
			log.Fatalln(err.Error())
		}
	// case "unpack":
	// 	if err := UnpackBootConfiguration(); err != nil {
	// 		log.Fatalln(err.Error())
	// 	}
	default:
		log.Fatal("Command not found")
	}
}
