// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package launcher boots the target kernel.
package launcher

import (
	"fmt"
	"log"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/mount"
	slaunch "github.com/u-root/u-root/pkg/securelaunch"
	"github.com/u-root/u-root/pkg/securelaunch/measurement"
	"github.com/u-root/uio/uio"
)

// Launcher describes the "launcher" section of policy file.
type Launcher struct {
	Type   string            `json:"type"`
	Params map[string]string `json:"params"`
}

// MeasureKernel hashes the kernel and extends the measurement into a TPM PCR.
func (l *Launcher) MeasureKernel() error {
	kernel := l.Params["kernel"]

	if err := measurement.HashFile(kernel); err != nil {
		return err
	}

	return nil
}

// MeasureInitrd hashes the initrd and extends the measurement into a TPM PCR.
func (l *Launcher) MeasureInitrd() error {
	initrd := l.Params["initrd"]

	if err := measurement.HashFile(initrd); err != nil {
		return err
	}

	return nil
}

// Boot boots the target kernel based on information provided in the "launcher"
// section of the policy file.
//
// Summary of steps:
// - extract the kernel, initrd and cmdline from the "launcher" section of policy file.
// - measure the kernel and initrd file into the tpmDev (tpm device).
// - mount the disks where the kernel and initrd file are located.
// - kexec to boot into the target kernel.
//
// returns error
// - if measurement of kernel and initrd fails
// - if mount fails
// - if kexec fails
func (l *Launcher) Boot() error {
	if l.Type != "kexec" {
		log.Printf("launcher: Unsupported launcher type. Exiting.")
		return fmt.Errorf("launcher: Unsupported launcher type. Exiting")
	}

	slaunch.Debug("Identified Launcher Type = Kexec")

	// TODO: if kernel and initrd are on different devices.
	kernel := l.Params["kernel"]
	initrd := l.Params["initrd"]
	cmdline := l.Params["cmdline"]

	k, e := slaunch.GetMountedFilePath(kernel, mount.MS_RDONLY)
	if e != nil {
		log.Printf("launcher: ERR: kernel input %s couldnt be located, err=%v", kernel, e)
		return e
	}

	i, e := slaunch.GetMountedFilePath(initrd, mount.MS_RDONLY)
	if e != nil {
		log.Printf("launcher: ERR: initrd input %s couldnt be located, err=%v", initrd, e)
		return e
	}

	slaunch.Debug("Calling kexec")
	image := &boot.LinuxImage{
		Kernel:  uio.NewLazyFile(k),
		Initrd:  uio.NewLazyFile(i),
		Cmdline: cmdline,
	}
	err := image.Load()
	if err != nil {
		log.Printf("kexec -l failed. err: %v", err)
		return err
	}

	err = kexec.Reboot()
	if err != nil {
		log.Printf("kexec reboot failed. err=%v", err)
		return err
	}
	return nil
}
