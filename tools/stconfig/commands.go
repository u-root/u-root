// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net"
	"path/filepath"

	"github.com/u-root/u-root/pkg/boot/stboot"
)

func packBootBall(config string, mac string) (err error) {
	var newName string
	if mac != "" {
		hwAddr, err := net.ParseMAC(mac)
		if err != nil {
			return err
		}
		newName = stboot.ComposeIndividualBallName(hwAddr)
	}

	ball, err := stboot.BootBallFromConfig(config)
	if err != nil {
		return
	}

	if newName != "" {
		ball.Archive = filepath.Join(filepath.Dir(ball.Archive), newName)
	}

	err = ball.Pack()
	if err != nil {
		return
	}

	log.Printf("Bootball created at: %s", ball.Archive)
	return ball.Clean()

}

func addSignatureToBootBall(bootBall, privKey, cert string) (err error) {
	ball, err := stboot.BootBallFromArchive(bootBall)
	if err != nil {
		return
	}

	log.Print("Signing bootball ...")
	log.Printf("private key: %s", privKey)
	log.Printf("certificate: %s", cert)
	err = ball.Sign(privKey, cert)
	if err != nil {
		return
	}

	if err = ball.Pack(); err != nil {
		return
	}

	log.Printf("Signatures included for each bootconfig: %d", ball.NumSignatures)
	return ball.Clean()
}

func unpackBootBall(bootBall string) (err error) {
	ball, err := stboot.BootBallFromArchive(bootBall)
	if err != nil {
		return err
	}

	log.Println("Archive unpacked into: " + ball.Dir())
	return ball.Clean()
}
