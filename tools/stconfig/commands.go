package main

import (
	"fmt"
	"io/ioutil"
	"os"
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
	if _, err := os.Stat(*genkeysPrivateKeyFile); os.IsNotExist(err) {
		return fmt.Errorf("private key file does not exist: %v", err)
	}
	if _, err := os.Stat(*genkeysPublicKeyFile); os.IsNotExist(err) {
		return fmt.Errorf("public key file does not exist: %v", err)
	}
	return crypto.GeneratED25519Key([]byte(*genkeysPassphrase), *genkeysPrivateKeyFile, *genkeysPublicKeyFile)
}

// PackBootConfiguration packages a boot configuration containing different
// binaries and a manifest. The files to be included are taken from the
// path specified in the provided manifest.json
func PackBootConfiguration() error {
	if _, err := os.Stat(*createManifest); os.IsNotExist(err) {
		return fmt.Errorf("manifest file does not exist: %v", err)
	}
	return bootconfig.ToZip(*createOutputFilename, *createManifest)
}

// AddSignatureToBootConfiguration TODO:
func AddSignatureToBootConfiguration() error {
	if _, err := os.Stat(*signInputBootfile); os.IsNotExist(err) {
		return fmt.Errorf("boot config file does not exist: %v", err)
	}
	if _, err := os.Stat(*signPrivKeyFile); os.IsNotExist(err) {
		return fmt.Errorf("private key file does not exist: %v", err)
	}
	if _, err := os.Stat(*signCertFile); os.IsNotExist(err) {
		return fmt.Errorf("certifivate file does not exist: %v", err)
	}
	return bootconfig.AddSignature(*signInputBootfile, *signPrivKeyFile, *signCertFile)
}

// UnpackBootConfiguration unpacks a boot configuration file and returns the
// file path of a directory containing the data
func UnpackBootConfiguration() error {
	_, outputDir, err := bootconfig.FromZip(*unpackInputFilename)
	if err != nil {
		return err
	}

	fmt.Println("Boot configuration unpacked into: " + outputDir)

	return nil
}
