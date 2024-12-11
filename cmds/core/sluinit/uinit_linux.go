// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/u-root/iscsinl"
	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/dhclient"
	slaunch "github.com/u-root/u-root/pkg/securelaunch"
	"github.com/u-root/u-root/pkg/securelaunch/config"
	"github.com/u-root/u-root/pkg/securelaunch/eventlog"
	"github.com/u-root/u-root/pkg/securelaunch/launcher"
	"github.com/u-root/u-root/pkg/securelaunch/policy"
	"github.com/u-root/u-root/pkg/securelaunch/tpm"
)

// booEntryFlag holds the name of the flag pointing to the boot entry to load.
const bootEntryFlag = "securelaunch_entry"

// pcrLogFilename holds the name of the file to dump PCR values to. This is
// used when provisioning the system to seal secrets to the correct values.
const pcrLogFilename = "securelaunch.dat"

// policyFilename is the filename to use for the policy file.
const policyFilename = "securelaunch.policy"

// pubkeyFilename is the filename to use for the policy public key file.
const pubkeyFilename = "securelaunch.pubkey"

// policyFilename is the filename to use for the policy signature file.
const signatureFilename = "securelaunch.sig"

// policyLocationFlag holds the name of the flag to specify where to find
// the securelaunch policy file.
const policyLocationFlag = "securelaunch_policy"

// pubkeyIndex hold the TPM index where the public key hash is stored.
const pubkeyHashIndex = uint32(0x01800180)

// ErrFlagNotSet indicates a required flag is not set.
var ErrFlagNotSet = errors.New("flag is not set")

// ErrInvalidCharacters indicates a name contains an unallowed character.
var ErrInvalidCharacters = errors.New("invalid character in name")

var slDebug = flag.Bool("d", false, "enable debug logs")

// step keeps track of the current step (e.g., parse policy, measure).
var step = 1

// printStep prints a message for the next step.
func printStep(msg string) {
	slaunch.Debug("******** Step %d: %s ********", step, msg)
	step++
}

// printStepDisabled prints a message for a disabled step.
func printStepDisabled(msg string) {
	slaunch.Debug("******** %s disabled in config ********", msg)
}

// checkPolicyLocationFlag checks the kernel cmdline for the policyLocationFlag
// flag. It provides the location of the policy file on disk. If it isn't set,
// an error is returned.
//
// The flag takes an argument formatted as `<block device id>:<path>`
//
//	e.g., sda1:/boot/securelaunch.policy
//	e.g., 4qccd342-12zr-4e99-9ze7-1234cb1234c4:/securelaunch.policy
func checkPolicyLocationFlag() (string, error) {
	location, present := cmdline.Flag(policyLocationFlag)
	if !present {
		return "", fmt.Errorf("%q:%w", policyLocationFlag, ErrFlagNotSet)
	}

	return location, nil
}

// checkDebugFlag checks if `uroot.uinitargs=-d` is set on the kernel cmdline.
// If it is set, slaunch.Debug is set to log.Printf.
func checkDebugFlag() {
	// By default, CommandLine exits on error, but this makes it trivial to get
	// a shell in u-root. Instead, continue on error and let the error handling
	// code here handle it.
	flag.CommandLine.Init(flag.CommandLine.Name(), flag.ContinueOnError)

	flag.Parse()

	if flag.NArg() > 1 {
		log.Fatal("Incorrect number of arguments")
	}

	if *slDebug {
		slaunch.Debug = log.Printf
		slaunch.Debug("debug flag is set. Logging Enabled.")
	}
}

func parseBootEntryFlag() (string, error) {
	bootEntry, present := cmdline.Flag(bootEntryFlag)
	if !present {
		return "", fmt.Errorf("no boot entry found")
	}

	if !launcher.IsValidBootEntry(bootEntry) {
		return "", fmt.Errorf("%q:%w", bootEntry, ErrInvalidCharacters)
	}

	slaunch.Debug("boot entry: %s", bootEntry)

	return bootEntry, nil
}

// iscsiSpecified checks if iscsi has been set on the kernel command line.
func iscsiSpecified() bool {
	return cmdline.ContainsFlag("netroot") && cmdline.ContainsFlag("rd.iscsi.initator")
}

// scanIscsiDrives calls dhcleint to parse cmdline and iscsinl to mount iscsi
// drives.
func scanIscsiDrives() error {
	uri, ok := cmdline.Flag("netroot")
	if !ok {
		return fmt.Errorf("could not get `netroot` argument")
	}
	slaunch.Debug("scanIscsiDrives: netroot flag is set: '%s'", uri)

	initiator, ok := cmdline.Flag("rd.iscsi.initiator")
	if !ok {
		return fmt.Errorf("could not get `rd.iscsi.initiator` argument")
	}
	slaunch.Debug("scanIscsiDrives: rd.iscsi.initiator flag is set: '%s'", initiator)

	target, volume, err := dhclient.ParseISCSIURI(uri)
	if err != nil {
		return fmt.Errorf("dhclient iSCSI parser failed: %w", err)
	}

	slaunch.Debug("scanIscsiDrives: resolved target: '%s'", target)
	slaunch.Debug("scanIscsiDrives: resolved volume: '%s'", volume)

	devices, err := iscsinl.MountIscsi(
		iscsinl.WithInitiator(initiator),
		iscsinl.WithTarget(target.String(), volume),
		iscsinl.WithCmdsMax(128),
		iscsinl.WithQueueDepth(16),
		iscsinl.WithScheduler("noop"),
	)
	if err != nil {
		return fmt.Errorf("could not mount iSCSI drive: %w", err)
	}

	for i := range devices {
		slaunch.Debug("scanIscsiDrives: iSCSI drive mounted at '%s'", devices[i])
	}

	return nil
}

// initialize sets up the environment.
func initialize(policyLocation string) error {
	printStep("Initialization")

	// Check if an iSCSI drive was specified and if so, mount it.
	if iscsiSpecified() {
		if err := scanIscsiDrives(); err != nil {
			return fmt.Errorf("failed to mount iSCSI drive: %w", err)
		}
	}

	if err := tpm.New(); err != nil {
		return fmt.Errorf("failed to get TPM device: %w", err)
	}

	// Write out all the PCRs values.
	pcrSelection := []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
	pcrLogLocation := filepath.Join(policyLocation, pcrLogFilename)

	if err := tpm.LogPCRs(pcrSelection, pcrLogLocation); err != nil {
		return fmt.Errorf("failed to log PCR values: %w", err)
	}

	slaunch.Debug("Initialization successfully completed")

	return nil
}

// verifyAndParsePolicy loads, verifies, and parses and gets the policy file.
func verifyAndParsePolicy(policyLocation string) (*policy.Policy, error) {
	printStep("Verify and parse SL policy")

	policyFileLocation := filepath.Join(policyLocation, policyFilename)
	pubkeyFileLocation := filepath.Join(policyLocation, pubkeyFilename)
	signatureFileLocation := filepath.Join(policyLocation, signatureFilename)

	if err := policy.Load(policyFileLocation, pubkeyFileLocation, signatureFileLocation); err != nil {
		return nil, fmt.Errorf("failed to load policy file: %w", err)
	}

	tpmPubkeyHashBytes, err := tpm.ReadValue32(pubkeyHashIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key hash from TPM: %w", err)
	}

	if err := policy.VerifyPubkey(tpmPubkeyHashBytes); err != nil {
		return nil, fmt.Errorf("failed to verify public key file: %w", err)
	}

	if err := policy.Verify(); err != nil {
		return nil, fmt.Errorf("failed to verify policy file: %w", err)
	}

	policy, err := policy.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse policy file: %w", err)
	}

	slaunch.Debug("Policy file successfully verified and parsed")

	return policy, nil
}

// loadAndVerifyBootEntry loads the specified boot entry and checks the hashes
// of the kernel and initrd files against the expected values.
func loadAndVerifyBootEntry(bootEntry string, bootEntries map[string]launcher.BootEntry) error {
	if err := launcher.MatchBootEntry(bootEntry, bootEntries); err != nil {
		return fmt.Errorf("could not find matching boot entry or hash checks failed: %w", err)
	}

	return nil
}

// collectMeasurements runs any measurements specified in the policy file.
func collectMeasurements(p *policy.Policy) error {
	if config.Conf.Collectors {
		printStep("Collect evidence")

		for _, collector := range p.Collectors {
			slaunch.Debug("Input Collector: %v", collector)
			if err := collector.Collect(); err != nil {
				log.Printf("Collector %v failed: %v", collector, err)
			}
		}

		slaunch.Debug("Collectors completed")
	} else {
		printStepDisabled("Collect evidence")
	}

	return nil
}

// measureFiles measures relevant files (e.g., policy, kernel, initrd).
func measureFiles(p *policy.Policy) error {
	if config.Conf.Measurements {
		printStep("Measure files")

		if err := policy.Measure(); err != nil {
			return fmt.Errorf("failed to measure policy file: %w", err)
		}

		if err := launcher.MeasureKernel(); err != nil {
			return fmt.Errorf("failed to measure target kernel: %w", err)
		}

		// It's possible to not have an initrd (e.g., if it's embedded).
		if launcher.IsInitrdSet() {
			if err := launcher.MeasureInitrd(); err != nil {
				return fmt.Errorf("failed to measure target initrd: %w", err)
			}
		}

		slaunch.Debug("Files successfully measured")
	} else {
		printStepDisabled("Measure files")
	}

	return nil
}

// parseEventLog parses the TPM event log.
func parseEventLog(p *policy.Policy) error {
	if config.Conf.EventLog {
		printStep("Parse event log")

		if err := p.EventLog.Parse(); err != nil {
			return fmt.Errorf("failed to parse event log: %w", err)
		}

		slaunch.Debug("Event log successfully parsed")
	} else {
		printStepDisabled("Parse event log")
	}

	return nil
}

// dumpLogs writes out any pending logs to a file on disk.
func dumpLogs() error {
	if config.Conf.Collectors || config.Conf.EventLog {
		printStep("Dump logs to disk")

		if err := eventlog.ParseEventLog(); err != nil {
			return fmt.Errorf("failed to parse event log: %w", err)
		}

		if err := slaunch.ClearPersistQueue(); err != nil {
			return fmt.Errorf("failed to clear persist queue: %w", err)
		}

		slaunch.Debug("Logs successfully dumped to disk")
	} else {
		printStepDisabled("Dump logs to disk")
	}

	return nil
}

// unmountAll unmounts all mount points.
func unmountAll() error {
	printStep("Unmount all")

	if err := slaunch.UnmountAll(); err != nil {
		return fmt.Errorf("failed to unmount all devices: %w", err)
	}

	slaunch.Debug("Devices successfully unmounted")

	return nil
}

// bootTarget boots the target kernel/initrd.
func bootTarget(p *policy.Policy) error {
	printStep("Boot target")

	if err := p.Launcher.Boot(); err != nil {
		return fmt.Errorf("failed to boot target: %w", err)
	}

	return nil
}

// exit loops forever trying to reboot the system.
func exit(mainErr error) {
	// Print the error.
	fmt.Fprintf(os.Stderr, "ERROR: Failed to boot: %v\n", mainErr)

	// Dump any logs, if possible. This can help figure out what went wrong.
	if err := dumpLogs(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not dump logs: %v\n", err)
	}

	// Umount anything that might be mounted.
	slaunch.UnmountAll()

	// Close the connection to the TPM if it was opened.
	tpm.Close()

	// Loop trying to reboot the system.
	for {
		// Wait 5 seconds.
		time.Sleep(5 * time.Second)

		// Try to reboot the system.
		if err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Failed to reboot: %v\n", err)
		}
	}
}

// main parses platform policy file, and based on the inputs performs
// measurements and then launches a target kernel.
//
// Steps followed by uinit:
// 1. if debug flag is set, enable logging.
// 2. gets the TPM handle
// 3. Gets secure launch policy file entered by user.
// 4. calls collectors to collect measurements(hashes) a.k.a evidence.
func main() {
	// Ignore ctrl+c
	signal.Ignore(syscall.SIGINT)

	checkDebugFlag()

	policyLocation, err := checkPolicyLocationFlag()
	if err != nil {
		exit(err)
	}

	bootEntry, err := parseBootEntryFlag()
	if err != nil {
		exit(err)
	}

	if err := initialize(policyLocation); err != nil {
		exit(err)
	}

	p, err := verifyAndParsePolicy(policyLocation)
	if err != nil {
		exit(err)
	}

	if err := parseEventLog(p); err != nil {
		exit(err)
	}

	if err := loadAndVerifyBootEntry(bootEntry, p.Launcher.BootEntries); err != nil {
		exit(err)
	}

	if err := collectMeasurements(p); err != nil {
		exit(err)
	}

	if err := measureFiles(p); err != nil {
		exit(err)
	}

	if err := dumpLogs(); err != nil {
		exit(err)
	}

	if err := unmountAll(); err != nil {
		exit(err)
	}

	if err := bootTarget(p); err != nil {
		exit(err)
	}
}
