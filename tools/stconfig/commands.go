// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/u-root/u-root/pkg/boot/stboot"
)

func packBootBall(config string) (err error) {
	ball, err := stboot.BootBallFromConfig(config)
	if err != nil {
		return
	}

	err = ball.Pack()
	if err != nil {
		return
	}

	log.Printf("Bootball created at: " + ball.Archive)
	return ball.Clean()

}

func addSignatureToBootBall(bootBall, privKey, cert string) (err error) {
	ball, err := stboot.BootBallFromArchive(bootBall)
	if err != nil {
		return
	}

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
