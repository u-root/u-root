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
	"github.com/u-root/u-root/pkg/uio"
)

/* describes the "launcher" section of policy file */
type Launcher struct {
	Type   string            `json:"type"`
	Params map[string]string `json:"params"`
}

// MeasureKernel calls file collector in measurement pkg that
// hashes kernel, initrd files and even store these hashes in tpm pcrs.
func (l *Launcher) MeasureKernel() error {

	kernel := l.Params["kernel"]
	initrd := l.Params["initrd"]

	if e := measurement.HashFile(kernel); e != nil {
		log.Printf("ERR: measure kernel input=%s, err=%v", kernel, e)
		return e
	}

	if e := measurement.HashFile(initrd); e != nil {
		log.Printf("ERR: measure initrd input=%s, err=%v", initrd, e)
		return e
	}
	return nil
}

/*
 * Boot boots the target kernel based on information provided
 * in the "launcher" section of policy file.
 *
 * Summary of steps:
 * - extracts the kernel, initrd and cmdline from the "launcher" section of policy file.
 * - measures the kernel and initrd file into the tpmDev (tpm device).
 * - mounts the disks where the kernel and initrd file are located.
 * - uses kexec to boot into the target kernel.
 * returns error
 * - if measurement of kernel and initrd fails
 * - if mount fails
 * - if kexec fails
 */
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

	slaunch.Debug("********Step 7: kexec called  ********")
	image := &boot.LinuxImage{
		Kernel:  uio.NewLazyFile(k),
		Initrd:  uio.NewLazyFile(i),
		Cmdline: cmdline,
	}
	err := image.Load(false)
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
