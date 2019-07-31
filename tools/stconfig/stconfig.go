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
	sign    = kingpin.Command("sign", "Sign the binary inside the provided stboot.zip and add the signatures and certificates")
	unpack  = kingpin.Command("unpack", "Unpack boot configuration file into directory")

	genkeysPrivateKeyFile = genkeys.Arg("privateKey", "File path to write the private key").Required().String()
	genkeysPublicKeyFile  = genkeys.Arg("publicKey", "File path to write the public key").Required().String()
	genkeysPassphrase     = genkeys.Flag("passphrase", "Encrypt keypair in PKCS8 format").String()

	createOutputFilename = create.Flag("output", "Path to output file").PlaceHolder("PATH").Default("stboot.zip").Short('o').String()
	createManifest       = create.Arg("manifest", "Path to the manifest file in JSON format").Required().String()

	signInputBootfile = sign.Arg("input", "stboot.zip file created by 'stconfig create'").Required().String()
	signPrivKeyFile   = sign.Arg("privkey", "private key for signing").Required().String()
	signCertFile      = sign.Arg("certificate", "Certificate to veryfy the signature").Required().String()

	unpackInputFilename = unpack.Arg("bc-file", "Boot configuration file").Required().String()
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
	case sign.FullCommand():
		if err := AddSignatureToBootConfiguration(); err != nil {
			log.Fatalln(err.Error())
		}
	case unpack.FullCommand():
		if err := UnpackBootConfiguration(); err != nil {
			log.Fatalln(err.Error())
		}
	default:
		log.Fatal("Command not found")
	}
}
