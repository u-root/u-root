package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/zaolin/go-tpm/tpm"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/sys/unix"
)

const (
	tpmDevice  string = "/dev/tpm0"
	mountPath  string = "/mnt/vboot"
	filesystem string = "ext3"
)

func dieHard() {
	if e1, e2, err := syscall.Syscall6(syscall.SYS_REBOOT, syscall.LINUX_REBOOT_MAGIC1, syscall.LINUX_REBOOT_MAGIC2, syscall.LINUX_REBOOT_CMD_POWER_OFF, 0, 0, 0); err != 0 {
		log.Fatalf("a %v b %v err %v", e1, e2, err)
	}
}

func main() {
	var publicKey = flag.String("pubkey", "/etc/sig.pub", "A public key which should verify the signature.")
	var pcr = flag.Uint("pcr", 12, "The pcr index used for measuring the kernel before kexec.")
	var bootDev = flag.String("boot-device", "/dev/sda1", "The boot device which is used to kexec into a signed kernel.")
	var linuxKernel = flag.String("kernel", "/mnt/vboot/kernel", "Kernel image file path.")
	var linuxKernelSignature = flag.String("kernel-sig", "/mnt/vboot/kernel.sig", "Kernel image signature file path.")
	var initrd = flag.String("initrd", "/mnt/vboot/initrd", "Initrd file path.")
	var initrdSignature = flag.String("initrd-sig", "/mnt/vboot/initrd.sig", "Initrd signature file path.")
	var debug = flag.Bool("debug", false, "Enables debug mode.")
	var noTPM = flag.Bool("no-tpm", false, "Disables tpm measuring process.")

	flag.Parse()

	if err := os.MkdirAll(mountPath, os.ModePerm); err != nil {
		if *debug {
			panic(err)
		} else {
			dieHard()
		}
	}

	if err := unix.Mount(*bootDev, mountPath, filesystem, unix.MS_RDONLY, ""); err != nil {
		if *debug {
			panic(err)
		} else {
			dieHard()
		}
	}

	paths := [...]string{*publicKey, *linuxKernel, *linuxKernelSignature, *initrd, *initrdSignature}
	files := make(map[string][]byte)

	for _, element := range paths {
		data, err := ioutil.ReadFile(element)
		if err != nil {
			if *debug {
				panic(err)
			} else {
				dieHard()
			}
		} else {
			files[element] = data
		}
	}

	kernelDigest := sha256.Sum256(files[*linuxKernel])
	initrdDigest := sha256.Sum256(files[*initrd])

	pcrDigestKernel := sha1.Sum(files[*linuxKernel])
	pcrDigestInitrd := sha1.Sum(files[*initrd])

	kernelSuccess := ed25519.Verify(files[*publicKey], kernelDigest[:], files[*linuxKernelSignature])
	initrdSuccess := ed25519.Verify(files[*publicKey], initrdDigest[:], files[*linuxKernelSignature])

	if kernelSuccess && initrdSuccess {
		if !*noTPM {
			rwc, err := tpm.OpenTPM(tpmDevice)
			if err != nil {
				if *debug {
					panic(err)
				} else {
					dieHard()
				}
			}

			tpm.PcrExtend(rwc, uint32(*pcr), pcrDigestKernel)
			tpm.PcrExtend(rwc, uint32(*pcr), pcrDigestInitrd)
		}

		binary, lookErr := exec.LookPath("kexec")
		if lookErr != nil {
			if *debug {
				panic(lookErr)
			} else {
				dieHard()
			}
		}

		args := []string{"kexec", "-initrd", *initrd, *linuxKernel}
		env := os.Environ()

		execErr := syscall.Exec(binary, args, env)
		if execErr != nil {
			if *debug {
				panic(execErr)
			} else {
				dieHard()
			}
		}
	} else {
		dieHard()
	}
}
