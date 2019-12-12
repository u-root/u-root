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
	return

}

func addSignatureToBootBall(bootBall, privKey, cert string) (err error) {
	ball, err := stboot.BootBallFromArchie(bootBall)
	if err != nil {
		return
	}
	err = ball.Sign(privKey, cert)
	if err != nil {
		return
	}
	return ball.Pack()
}

func unpackBootBall(bootBall string) (err error) {
	ball, err := stboot.BootBallFromArchie(bootBall)
	if err != nil {
		return err
	}
	log.Println("Archive unpacked into: " + ball.Dir())
	return nil
}
