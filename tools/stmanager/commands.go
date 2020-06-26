// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net"
	"path/filepath"

	"github.com/u-root/u-root/pkg/boot/stboot"
)

func packBootBall(outDir, label, kernel, initramfs, cmdline, tboot, tbootArgs, rootCert string, acms []string, allowNonTXT bool, mac string) error {
	var individual string
	if mac != "" {
		hwAddr, err := net.ParseMAC(mac)
		if err != nil {
			return err
		}
		individual = stboot.ComposeIndividualBallPrefix(hwAddr)
	}

	ball, err := stboot.InitBootball(outDir, label, kernel, initramfs, cmdline, tboot, tbootArgs, rootCert, acms, allowNonTXT)
	if err != nil {
		return err
	}

	if individual != "" {
		name := filepath.Base(ball.Archive)
		name = individual + name
		ball.Archive = filepath.Join(filepath.Dir(ball.Archive), name)
	}

	err = ball.Pack()
	if err != nil {
		return err
	}

	fmt.Println(filepath.Base(ball.Archive))
	return ball.Clean()

}

func addSignatureToBootBall(bootBall, privKey, cert string) error {
	ball, err := stboot.BootballFromArchive(bootBall)
	if err != nil {
		return err
	}

	log.Print("Signing bootball ...")
	log.Printf("private key: %s", privKey)
	log.Printf("certificate: %s", cert)
	err = ball.Sign(privKey, cert)
	if err != nil {
		return err
	}

	if err = ball.Pack(); err != nil {
		return err
	}

	log.Printf("Signatures included: %d", ball.NumSignatures)
	return ball.Clean()
}

func unpackBootBall(bootBall string) error {
	ball, err := stboot.BootballFromArchive(bootBall)
	if err != nil {
		return err
	}

	log.Println("Archive unpacked into: " + ball.Dir)
	return nil
}
