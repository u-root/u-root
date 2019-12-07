package main

import (
	"log"
	"path"

	"github.com/u-root/u-root/pkg/boot/stboot"
)

func packBootBall(config string) error {
	outPath := path.Join(path.Dir(config), stboot.BallName)
	return stboot.ToZip(outPath, config)
}

func addSignatureToBootBall(bootBall, privKey, cert string) error {
	return stboot.AddSignature(bootBall, privKey, cert)
}

func unpackBootBall(bootBall string) error {
	_, outputDir, err := stboot.FromZip(bootBall)
	if err != nil {
		return err
	}
	log.Println("Archive unpacked into: " + outputDir)
	return nil
}
