// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/u-root/iscsinl"
	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/dhclient"
	slaunch "github.com/u-root/u-root/pkg/securelaunch"
	"github.com/u-root/u-root/pkg/securelaunch/policy"
	"github.com/u-root/u-root/pkg/securelaunch/tpm"
)

var (
	slDebug = flag.Bool("d", false, "enable debug logs")
)

func checkDebugFlag() {
	/*
	 * check if uroot.uinitargs=-d is set in kernel cmdline.
	 * if set, slaunch.Debug is set to log.Printf.
	 */
	flag.Parse()

	if flag.NArg() > 1 {
		log.Fatal("Incorrect number of arguments")
	}

	if *slDebug {
		slaunch.Debug = log.Printf
		slaunch.Debug("debug flag is set. Logging Enabled.")
	}
}

/*
 * main parses platform policy file, and based on the inputs,
 * performs measurements and then launches a target kernel.
 *
 * steps followed by sluinit:
 * 1. if debug flag is set, enable logging.
 * 2. gets the TPM handle
 * 3. Gets secure launch policy file entered by user.
 * 4. calls collectors to collect measurements(hashes) a.k.a evidence.
 */
func main() {
	checkDebugFlag()

	err := scanIscsiDrives()
	if err != nil {
		log.Printf("NO ISCSI DRIVES found, err=[%v]", err)
	}

	defer unmountAndExit() // called only on error, on success we kexec
	slaunch.Debug("********Step 1: init completed. starting main ********")
	if err := tpm.New(); err != nil {
		log.Printf("tpm.New() failed. err=%v", err)
		return
	}
	defer tpm.Close()

	slaunch.Debug("********Step 2: locate and parse SL Policy ********")
	p, err := policy.Get()
	if err != nil {
		log.Printf("failed to get policy err=%v", err)
		return
	}
	slaunch.Debug("policy file successfully parsed")

	slaunch.Debug("********Step 3: Collecting Evidence ********")
	for _, c := range p.Collectors {
		slaunch.Debug("Input Collector: %v", c)
		if e := c.Collect(); e != nil {
			log.Printf("Collector %v failed, err = %v", c, e)
		}
	}
	slaunch.Debug("Collectors completed")

	slaunch.Debug("********Step 4: Measuring target kernel, initrd ********")
	if err := p.Launcher.MeasureKernel(); err != nil {
		log.Printf("Launcher.MeasureKernel failed err=%v", err)
		return
	}

	slaunch.Debug("********Step 5: Parse eventlogs *********")
	if err := p.EventLog.Parse(); err != nil {
		log.Printf("EventLog.Parse() failed err=%v", err)
		return
	}

	slaunch.Debug("*****Step 6: Dump logs to disk *******")
	if err := slaunch.ClearPersistQueue(); err != nil {
		log.Printf("ClearPersistQueue failed err=%v", err)
		return
	}

	slaunch.Debug("********Step *: Unmount all ********")
	slaunch.UnmountAll()

	slaunch.Debug("********Step 7: Launcher called to Boot ********")
	if err := p.Launcher.Boot(); err != nil {
		log.Printf("Boot failed. err=%s", err)
		return
	}
}

// unmountAndExit is called on error and unmounts all devices.
// sluinit ends here.
func unmountAndExit() {
	slaunch.UnmountAll()
	time.Sleep(5 * time.Second) // let queued up debug statements get printed
	os.Exit(1)
}

// scanIscsiDrives calls dhcleint to parse cmdline and
// iscsinl to mount iscsi drives.
func scanIscsiDrives() error {
	val, ok := cmdline.Flag("netroot")
	if !ok {
		return errors.New("netroot flag is not set")
	}
	slaunch.Debug("netroot flag is set with val=%s", val)

	target, volume, err := dhclient.ParseISCSIURI(val)
	if err != nil {
		return fmt.Errorf("dhclient ISCSI parser failed err=%v", err)
	}

	slaunch.Debug("resolved ip:port=%s", target)
	slaunch.Debug("resolved vol=%v", volume)

	slaunch.Debug("Scanning kernel cmd line for *rd.iscsi.initiator* flag")
	initiatorName, ok := cmdline.Flag("rd.iscsi.initiator")
	if !ok {
		return errors.New("rd.iscsi.initiator flag is not set")
	}

	devices, err := iscsinl.MountIscsi(
		iscsinl.WithInitiator(initiatorName),
		iscsinl.WithTarget(target.String(), volume),
		iscsinl.WithCmdsMax(128),
		iscsinl.WithQueueDepth(16),
		iscsinl.WithScheduler("noop"),
	)
	if err != nil {
		return err
	}

	for i := range devices {
		slaunch.Debug("Mounted at dev %v", devices[i])
	}
	return nil
}
