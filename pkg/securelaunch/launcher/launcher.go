// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package launcher boots the target kernel.
package launcher

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"unicode"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/kexec"
	slaunch "github.com/u-root/u-root/pkg/securelaunch"
	"github.com/u-root/u-root/pkg/securelaunch/measurement"
	"github.com/u-root/u-root/pkg/securelaunch/tpm"
)

// BootEntry holds the names and hashes for a kernel and initrd and the command
// line to use.
type BootEntry struct {
	KernelName  string `json:"kernel name"`
	KernelHash  string `json:"kernel hash"`
	KernelBytes []byte
	InitrdName  string `json:"initrd name"`
	InitrdHash  string `json:"initrd hash"`
	InitrdBytes []byte
	Cmdline     string `json:"cmdline"`
}

// Launcher describes the "launcher" section of policy file.
type Launcher struct {
	Type        string               `json:"type"`
	BootEntries map[string]BootEntry `json:"boot entries"`
}

// ErrBootNotSelected means a boot was not selected and it must be
var ErrBootNotSelected = errors.New("boot entry not yet selected")

// bootEntry points to the target BootEntry to use.
var bootEntry *BootEntry

// readAndVerifyFile reads a file and checks its hash against the provided
// expected hash. If the hashes match, it returns a byte slice of the file.
func readAndVerifyFile(fileName string, expectedHash string) ([]byte, error) {
	slaunch.Debug("readAndVerifyFile: reading file '%s'", fileName)

	// Read the file.
	fileBytes, err := slaunch.GetFileBytes(fileName)
	if err != nil {
		slaunch.Debug("readAndVerifyFile: ERR: could not read file '%s': %v", fileName, err)
		return nil, fmt.Errorf("could not read file '%s': %w", fileName, err)
	}

	// Get the file's hash.
	fileHash := hex.EncodeToString(tpm.HashReader(bytes.NewReader(fileBytes)))

	// Compare hashes.
	if expectedHash != fileHash {
		slaunch.Debug("readAndVerifyFile: ERR: file hash (%s) does not match expected (%s)", fileHash, expectedHash)
		return nil, fmt.Errorf("file hash (%s) does not match expected (%s)", fileHash, expectedHash)
	}

	return fileBytes, nil
}

// measureFile extends the measurement of the provided file into the TPM.
func measureFile(fileName string, fileBytes []byte) error {
	if len(fileBytes) == 0 {
		return fmt.Errorf("file not yet loaded or empty")
	}

	eventDesc := fmt.Sprintf("File Collector: measured %s", fileName)
	return measurement.HashBytes(fileBytes, eventDesc)
}

// MeasureKernel hashes the kernel and extends the measurement into a TPM PCR.
func MeasureKernel() error {
	if bootEntry == nil {
		return ErrBootNotSelected
	}

	return measureFile(bootEntry.KernelName, bootEntry.KernelBytes)
}

// MeasureInitrd hashes the initrd and extends the measurement into a TPM PCR.
func MeasureInitrd() error {
	if bootEntry == nil {
		return ErrBootNotSelected
	}

	return measureFile(bootEntry.InitrdName, bootEntry.InitrdBytes)
}

// IsInitrdSet returns whether an initrd has been set or not.
func IsInitrdSet() bool {
	return (bootEntry != nil) && (bootEntry.InitrdName != "")
}

// IsValidBootEntry validates that the provided string compiles to the rules for
// boot entries. Specifically, this means all alphanumeric characters, plus '-',
// '_', and '.'.
func IsValidBootEntry(entry string) bool {
	// Sanitize boot entry; alphanumeric characters and - and _ only.
	for _, rune := range entry {
		if !(unicode.IsLetter(rune) || unicode.IsNumber(rune) ||
			(rune == '-') || (rune == '_') || (rune == '.')) {
			return false
		}
	}

	return true
}

// MatchBootEntry tries to match the given name to a boot entry. If successful,
// the kernel and initrd files are read and the command-line is returned.
func MatchBootEntry(entryName string, bootEntries map[string]BootEntry) error {
	slaunch.Debug("launcher: MatchBootEntry: Looking for entry name = '%s'", entryName)

	for name, entry := range bootEntries {
		slaunch.Debug("launcher: MatchBootEntry: entry name = '%s'", name)
		slaunch.Debug("launcher: matchBootEntry: Found kernel '%s'", entry.KernelName)
		slaunch.Debug("launcher: matchBootEntry: Found initrd '%s'", entry.InitrdName)
		slaunch.Debug("launcher: matchBootEntry: Found cmdline '%s'", entry.Cmdline)

		if name == entryName {
			var err error
			slaunch.Debug("launcher: MatchBootEntry: Found entry '%s'", entryName)

			// Read and verify kernel and initrd.
			entry.KernelBytes, err = readAndVerifyFile(entry.KernelName, entry.KernelHash)
			if err != nil {
				log.Printf("launcher: matchBootEntry: ERR: Could not read kernel '%s': %v", bootEntry.KernelName, err)
				return err
			}

			entry.InitrdBytes, err = readAndVerifyFile(entry.InitrdName, entry.InitrdHash)
			if err != nil {
				log.Printf("launcher: matchBootEntry: ERR: Could not read initrd '%s': %v", bootEntry.InitrdName, err)
				return err
			}

			bootEntry = &entry

			return nil
		}
	}

	log.Printf("launcher: matchBootEntry: ERR: boot entry '%s' not found", entryName)
	return fmt.Errorf("boot entry '%s' not found", entryName)
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
		return fmt.Errorf("unsupported launcher type %q:%w", l.Type, os.ErrInvalid)
	}
	slaunch.Debug("Identified Launcher Type = Kexec")

	if bootEntry == nil {
		return ErrBootNotSelected
	}

	slaunch.Debug("Calling kexec")
	image := &boot.LinuxImage{
		Kernel:  bytes.NewReader(bootEntry.KernelBytes),
		Initrd:  bytes.NewReader(bootEntry.InitrdBytes),
		Cmdline: bootEntry.Cmdline,
	}

	if err := image.Load(); err != nil {
		return fmt.Errorf("kexec -l failed:%w", err)
	}

	if err := kexec.Reboot(); err != nil {
		return fmt.Errorf("kexec reboot failed:%w", err)
	}

	return nil
}
