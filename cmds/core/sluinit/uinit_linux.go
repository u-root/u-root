// Copyright 2019-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo || tinygo.enable

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/u-root/iscsinl"
	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/dhclient"
	slaunch "github.com/u-root/u-root/pkg/securelaunch"
	"github.com/u-root/u-root/pkg/securelaunch/eventlog"
	"github.com/u-root/u-root/pkg/securelaunch/policy"
	"github.com/u-root/u-root/pkg/securelaunch/tpm"
)

var (
	slDebug = flag.Bool("d", false, "enable debug logs")

	// step keeps track of the current step (e.g., parse policy, measure).
	step = 1
)

// printStep prints a message for the next step.
func printStep(msg string) {
	slaunch.Debug("******** Step %d: %s ********", step, msg)
	step++
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
func initialize() error {
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

	slaunch.Debug("Initialization successfully completed")

	return nil
}

// parsePolicy parses and gets the policy file.
func parsePolicy() (*policy.Policy, error) {
	printStep("Locate and parse SL policy")

	p, err := policy.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to parse policy file: %w", err)
	}

	slaunch.Debug("Policy file successfully parsed")

	return p, nil
}

// collectMeasurements runs any measurements specified in the policy file.
func collectMeasurements(p *policy.Policy) error {
	printStep("Collect evidence")

	for _, collector := range p.Collectors {
		slaunch.Debug("Input Collector: %v", collector)
		if err := collector.Collect(); err != nil {
			log.Printf("Collector %v failed: %v", collector, err)
		}
	}

	slaunch.Debug("Collectors completed")

	return nil
}

// measureFiles measures relevant files (e.g., policy, kernel, initrd).
func measureFiles(p *policy.Policy) error {
	printStep("Measure files")

	if err := policy.Measure(); err != nil {
		return fmt.Errorf("failed to measure policy file: %w", err)
	}

	if p.Launcher.Params["kernel"] != "" {
		if err := p.Launcher.MeasureKernel(); err != nil {
			return fmt.Errorf("failed to measure target kernel: %w", err)
		}
	}

	if p.Launcher.Params["initrd"] != "" {
		if err := p.Launcher.MeasureInitrd(); err != nil {
			return fmt.Errorf("failed to measure target initrd: %w", err)
		}
	}

	slaunch.Debug("Files successfully measured")

	return nil
}

// parseEventLog parses the TPM event log.
func parseEventLog(p *policy.Policy) error {
	printStep("Parse event log")

	if err := p.EventLog.Parse(); err != nil {
		return fmt.Errorf("failed to parse event log: %w", err)
	}

	slaunch.Debug("Event log successfully parsed")

	return nil
}

// dumpLogs writes out any pending logs to a file on disk.
func dumpLogs() error {
	printStep("Dump logs to disk")

	if err := eventlog.ParseEventLog(); err != nil {
		return fmt.Errorf("failed to parse event log: %w", err)
	}

	if err := slaunch.ClearPersistQueue(); err != nil {
		return fmt.Errorf("failed to clear persist queue: %w", err)
	}

	slaunch.Debug("Logs successfully dumped to disk")

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

	if err := initialize(); err != nil {
		exit(err)
	}

	p, err := parsePolicy()
	if err != nil {
		exit(err)
	}

	if err := parseEventLog(p); err != nil {
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
