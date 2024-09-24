// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"flag"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/google/go-tpm/tpm"
	"golang.org/x/crypto/ed25519"
)

const (
	tpmDevice  string = "/dev/tpm0"
	mountPath  string = "/mnt/vboot"
	filesystem string = "ext3"
)

var (
	publicKey            = flag.String("pubkey", "/etc/sig.pub", "A public key which should verify the signature.")
	pcr                  = flag.Uint("pcr", 12, "The pcr index used for measuring the kernel before kexec.")
	bootDev              = flag.String("boot-device", "/dev/sda1", "The boot device which is used to kexec into a signed kernel.")
	linuxKernel          = flag.String("kernel", "/mnt/vboot/kernel", "Kernel image file path.")
	linuxKernelSignature = flag.String("kernel-sig", "/mnt/vboot/kernel.sig", "Kernel image signature file path.")
	initrd               = flag.String("initrd", "/mnt/vboot/initrd", "Initrd file path.")
	initrdSignature      = flag.String("initrd-sig", "/mnt/vboot/initrd.sig", "Initrd signature file path.")
	debug                = flag.Bool("debug", false, "Enables debug mode.")
	noTPM                = flag.Bool("no-tpm", false, "Disables tpm measuring process.")
)

func die(err error) {
	if *debug {
		panic(err)
	}
	if err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF); err != nil {
		log.Fatalf("reboot err: %v", err)
	}
}

func main() {
	flag.Parse()

	if err := os.MkdirAll(mountPath, os.ModePerm); err != nil {
		die(err)
	}

	if err := syscall.Mount(*bootDev, mountPath, filesystem, syscall.MS_RDONLY, ""); err != nil {
		die(err)
	}

	paths := []string{*publicKey, *linuxKernel, *linuxKernelSignature, *initrd, *initrdSignature}
	files := make(map[string][]byte)

	for _, element := range paths {
		data, err := os.ReadFile(element)
		if err != nil {
			die(err)
		} else {
			files[element] = data
		}
	}

	kernelDigest := sha256.Sum256(files[*linuxKernel])
	initrdDigest := sha256.Sum256(files[*initrd])

	pcrDigestKernel := sha1.Sum(files[*linuxKernel])
	pcrDigestInitrd := sha1.Sum(files[*initrd])

	kernelSuccess := ed25519.Verify(files[*publicKey], kernelDigest[:], files[*linuxKernelSignature])
	initrdSuccess := ed25519.Verify(files[*publicKey], initrdDigest[:], files[*initrdSignature])

	if !kernelSuccess || !initrdSuccess {
		die(nil)
	}

	if !*noTPM {
		rwc, err := tpm.OpenTPM(tpmDevice)
		if err != nil {
			die(err)
		}

		tpm.PcrExtend(rwc, uint32(*pcr), pcrDigestKernel)
		tpm.PcrExtend(rwc, uint32(*pcr), pcrDigestInitrd)
	}

	binary, lookErr := exec.LookPath("kexec")
	if lookErr != nil {
		die(lookErr)
	}

	args := []string{"kexec", "-initrd", *initrd, *linuxKernel}
	env := os.Environ()

	if execErr := syscall.Exec(binary, args, env); execErr != nil {
		die(execErr)
	}
}
