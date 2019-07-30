package main

import (
	"io/ioutil"
	"path"

	"github.com/u-root/u-root/pkg/bootconfig"
	"github.com/u-root/u-root/pkg/crypto"
)

func getFilePathsByDir(dirName string) ([]string, error) {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return nil, err
	}

	var listOfFilePaths []string
	for _, file := range files {
		if !file.IsDir() {
			listOfFilePaths = append(listOfFilePaths, path.Join(dirName, file.Name()))
		}
	}

	return listOfFilePaths, nil
}

// GenKeys generates ED25519 keypair and stores it on the harddrive
func GenKeys() error {
	return crypto.GeneratED25519Key([]byte(*genkeysPassphrase), *genkeysPrivateKeyFile, *genkeysPublicKeyFile)
}

// PackBootConfiguration packages a boot configuration containing different
// binaries and a manifest. The files to be included are taken from the
// path specified in the provided manifest.json
func PackBootConfiguration() error {
	return bootconfig.ToZip(*createOutputFilename, *createManifest)
}

// UnpackBootConfiguration unpacks a boot configuration file and returns the
// file path of a directory containing the data
// func UnpackBootConfiguration() error {
// 	if *unpackVerifyPublicKeyFile == "" {
// 		// FromZip expects that no key is provided, only if pointer is nil
// 		unpackVerifyPublicKeyFile = nil
// 	}

// 	_, outputDir, err := bootconfig.FromZip(*unpackInputFilename, unpackVerifyPublicKeyFile)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("Boot configuration unpacked into: " + outputDir)

// 	return nil
// }
