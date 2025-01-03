// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vpd

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func handler(c <-chan os.Signal) {
	for range c {
		log.Printf("ignoring SIGINT during flash write to prevent corruption")
	}
}

// Set RW_VPD key-value via flashrom and vpd executables, remove set to false would set or add the key,
// delete set to true would delete an existing key.
func FlashromRWVpdSet(key string, value []byte, remove bool) error {
	file, err := os.CreateTemp("/tmp", "rwvpd*.bin")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	c := make(chan os.Signal, 1)
	go handler(c)
	defer close(c)

	cmd := exec.Command("flashrom", "-p", "internal:ich_spi_mode=hwseq", "-c", "Opaque flash chip", "--fmap", "-i", "RW_VPD", "-r", file.Name())
	cmd.Stdin, cmd.Stdout = os.Stdin, os.Stdout
	if err = cmd.Run(); err != nil {
		log.Printf("flashrom failed to read RW_VPD: %v", err)
		return err
	}

	if remove {
		cmd = exec.Command("vpd", "-f", file.Name(), "-d", key)
		if err = cmd.Run(); err != nil {
			log.Printf("vpd failed to delete key: %v, err: %v", key, err)
			return err
		}
	} else {
		cmd = exec.Command("vpd", "-f", file.Name(), "-s", key+"="+string(value[:]))
		if err = cmd.Run(); err != nil {
			log.Printf("vpd failed to set key: %v value: %v, err: %v", key, string(value[:]), err)
			return err
		}
	}

	cmd = exec.Command("flashrom", "-p", "internal:ich_spi_mode=hwseq", "-c", "Opaque flash chip", "--fmap", "-i", "RW_VPD", "--noverify-all", "-w", file.Name())
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	signal.Notify(c, syscall.SIGINT)
	defer signal.Reset(syscall.SIGINT)
	if err = cmd.Run(); err != nil {
		log.Printf("flashrom failed to write RW_VPD: %v", err)
		return err
	}
	return nil
}

// ClearRwVpd re-format RW_VPD via flashrom and vpd executables
func ClearRwVpd() error {
	file, err := os.CreateTemp("/tmp", "rwvpd*.bin")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	c := make(chan os.Signal, 1)
	go handler(c)
	defer close(c)

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
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	signal.Notify(c, syscall.SIGINT)
	defer signal.Reset(syscall.SIGINT)
	if err = cmd.Run(); err != nil {
		log.Printf("flashrom failed to write RW_VPD: %v", err)
		return err
	}
	return nil
}

// FlashromVpdDump read and dump all VPD values from RO and RW VPD flash regions directly
func FlashromVpdDump() error {
	file, err := os.CreateTemp("/tmp", "rovpd*.bin")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	log.Printf("RO_VPD values:")
	cmd := exec.Command("flashrom", "-p", "internal:ich_spi_mode=hwseq", "-c", "Opaque flash chip", "--fmap", "-i", "RO_VPD", "-r", file.Name())
	if err = cmd.Run(); err != nil {
		log.Printf("flashrom failed to read RO_VPD: %v", err)
		return err
	}
	cmd = exec.Command("vpd", "-f", file.Name(), "-l")
	cmd.Stdin, cmd.Stdout = os.Stdin, os.Stdout
	if err = cmd.Run(); err != nil {
		log.Printf("vpd failed to print RO_VPD: %v", err)
		return err
	}

	os.Remove(file.Name())
	file, err = os.CreateTemp("/tmp", "rwvpd*.bin")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	log.Printf("RW_VPD values:")
	cmd = exec.Command("flashrom", "-p", "internal:ich_spi_mode=hwseq", "-c", "Opaque flash chip", "--fmap", "-i", "RW_VPD", "-r", file.Name())
	if err = cmd.Run(); err != nil {
		log.Printf("flashrom failed to read RW_VPD: %v", err)
		return err
	}
	cmd = exec.Command("vpd", "-f", file.Name(), "-l")
	cmd.Stdin, cmd.Stdout = os.Stdin, os.Stdout
	if err = cmd.Run(); err != nil {
		log.Printf("vpd failed to print RW_VPD: %v", err)
		return err
	}

	return nil
}
