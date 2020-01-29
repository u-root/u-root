// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os"

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

	slaunch.Debug("********Step 1: init completed. starting main ********")
	tpmDev, err := tpm.GetHandle()
	if err != nil {
		log.Printf("tpm.getHandle failed. err=%v", err)
		os.Exit(1)
	}
	defer tpmDev.Close()

	slaunch.Debug("********Step 2: locate and parse SL Policy ********")
	p, err := policy.Get()
	if err != nil {
		log.Printf("failed to get policy err=%v", err)
		os.Exit(1)
	}
	slaunch.Debug("policy file successfully parsed")

	slaunch.Debug("********Step 3: Collecting Evidence ********")
	for _, c := range p.Collectors {
		slaunch.Debug("Input Collector: %v", c)
		if e := c.Collect(tpmDev); e != nil {
			log.Printf("Collector %v failed, err = %v", c, e)
		}
	}
	slaunch.Debug("Collectors completed")

	slaunch.Debug("********Step 4: Measuring target kernel, initrd ********")
	if e := p.Launcher.MeasureKernel(tpmDev); e != nil {
		log.Printf("Launcher.MeasureKernel failed err=%v", e)
		os.Exit(1)
	}

	slaunch.Debug("********Step 5: Write eventlog to /boot partition*********")
	if e := p.EventLog.Persist(); e != nil {
		log.Printf("EventLog.Persist() failed err=%v", e)
		os.Exit(1)
	}

	slaunch.Debug("********Step 6: Launcher called to Boot ********")
	err = p.Launcher.Boot(tpmDev)
	log.Printf("Boot failed. err=%s", err)
}
