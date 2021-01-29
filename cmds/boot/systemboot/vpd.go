// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func clearRwVpd() error {
	file, err := ioutil.TempFile("/tmp", "rwvpd*.bin")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	log.Printf("Reading RW_VPD...")
	cmd := exec.Command("flashrom", "-p", "internal:ich_spi_mode=hwseq", "-c", "Opaque flash chip", "--fmap", "-i", "RW_VPD", "-r", file.Name())
	cmd.Stdin, cmd.Stdout = os.Stdin, os.Stdout
	if err = cmd.Run(); err != nil {
		log.Printf("flashrom failed to read RW_VPD: %v", err)
		return err
	}
	cmd = exec.Command("vpd", "-f", file.Name(), "-O")
	if err = cmd.Run(); err != nil {
		log.Printf("vpd failed to re-format RW_VPD: %v", err)
		return err
	}
	cmd = exec.Command("flashrom", "-p", "internal:ich_spi_mode=hwseq", "-c", "Opaque flash chip", "--fmap", "-i", "RW_VPD", "--noverify-all", "-w", file.Name())
	if err = cmd.Run(); err != nil {
		log.Printf("flashrom failed to write RW_VPD: %v", err)
		return err
	}
	return nil
}
