// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/boot/systembooter"
	"github.com/u-root/u-root/pkg/ipmi"
	"github.com/u-root/u-root/pkg/ipmi/ocp"
	"github.com/u-root/u-root/pkg/smbios"
	"github.com/u-root/u-root/pkg/ulog"
	"github.com/u-root/u-root/pkg/vpd"
)

var (
	allowInteractive = flag.Bool("i", true, "Allow user to interrupt boot process and run commands")
	doQuiet          = flag.Bool("q", false, fmt.Sprintf("Disable verbose output. If not specified, read it from VPD var '%s'. Default false", vpdSystembootLogLevel))
	interval         = flag.Int("I", 1, "Interval in seconds before looping to the next boot command")
	noDefaultBoot    = flag.Bool("nodefault", false, "Do not attempt default boot entries if regular ones fail")
)

const (
	// vpdSystembootLogLevel is the name of the VPD variable used to set the log level.
	vpdSystembootLogLevel = "systemboot_log_level"
)

// isFlagPassed checks whether a flag was explicitly passed on the command line
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

var defaultBootsequence = [][]string{
	{"pxeboot", "-ipv6=true", "-ipv4=false"},
	{"boot"},
	{"fbnetboot", "-userclass", "linuxboot"},
}

// VPD variable for enabling IPMI BMC overriding boot order, default is not set
const VpdBmcBootOrderOverride = "bmc_bootorder_override"

var bmcBootOverride bool

// Product list for running IPMI OEM commands
var productList = [5]string{"Tioga Pass", "Mono Lake", "Delta Lake", "Crater Lake", "S9S"}

var selRecorded bool

func isMatched(productName string) bool {
	for _, v := range productList {
		if strings.HasPrefix(productName, v) {
			return true
		}
	}
	return false
}

func getBaseboardProductName(si *smbios.Info) (string, error) {
	t2, err := si.GetBaseboardInfo()
	if err != nil {
		log.Printf("Error getting Baseboard Information: %v", err)
		return "", err
	}
	return t2[0].Product, nil
}

func getSystemFWVersion(si *smbios.Info) (string, error) {
	t0, err := si.GetBIOSInfo()
	if err != nil {
		log.Printf("Error getting BIOS Information: %v", err)
		return "", err
	}
	return t0.Version, nil
}

func checkCMOSClear(ipmi *ipmi.IPMI) error {
	if cmosclear, bootorder, err := ocp.IsCMOSClearSet(ipmi); cmosclear {
		log.Printf("CMOS clear starts")
		if err = cmosClear(); err != nil {
			return err
		}
		if err = vpd.ClearRwVpd(); err != nil {
			return err
		}

		if err = ocp.ClearCMOSClearValidBits(ipmi, bootorder); err != nil {
			return err
		}
		addSEL("cmosclear")
		if err = reboot(); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func runIPMICommands(l ulog.Logger) {
	i, err := ipmi.Open(0)
	if err != nil {
		log.Printf("Failed to open ipmi device %v, watchdog may still be running", err)
		return
	}
	defer i.Close()

	if err = i.ShutoffWatchdog(); err != nil {
		log.Printf("Failed to stop watchdog %v.", err)
	} else {
		log.Printf("Watchdog is stopped.")
	}
	// Try RW_VPD first
	value, err := systembooter.Get(VpdBmcBootOrderOverride, false)
	if err != nil {
		// Try RO_VPD
		value, err = systembooter.Get(VpdBmcBootOrderOverride, true)
	}
	if err == nil && string(value) == "1" {
		bmcBootOverride = true
	}
	log.Printf("VPD %s is %v", VpdBmcBootOrderOverride, string(value))
	// Below IPMI commands would require SMBIOS data
	si, err := smbios.FromSysfs()
	if err != nil {
		log.Printf("Error reading SMBIOS info: %v", err)
		return
	}

	if fwVersion, err := getSystemFWVersion(si); err == nil {
		log.Printf("System firmware version: %s", fwVersion)
		if err = i.SetSystemFWVersion(fwVersion); err != nil {
			log.Printf("Failed to set system firmware version to BMC %v.", err)
		}
	}

	if productName, err := getBaseboardProductName(si); err == nil {
		if isMatched(productName) {
			log.Printf("Running OEM IPMI commands.")
			if err = checkCMOSClear(i); err != nil {
				log.Printf("IPMI CMOS clear err: %v", err)
			}
			if err = ocp.CheckBMCBootOrder(i, bmcBootOverride, l); err != nil {
				log.Printf("Failed to sync BMC Boot Order %v.", err)
			}
			dimmInfo, err := ocp.GetOemIpmiDimmInfo(si)
			if err == nil {
				if err = ocp.SendOemIpmiDimmInfo(i, dimmInfo); err == nil {
					log.Printf("Send the information of DIMMs to BMC.")
				} else {
					log.Printf("Failed to send the information of DIMMs to BMC: %v.", err)
				}
			} else {
				log.Printf("Failed to get the information of DIMMs: %v.", err)
			}

			processorInfo, err := ocp.GetOemIpmiProcessorInfo(si)
			if err == nil {
				if err = ocp.SendOemIpmiProcessorInfo(i, processorInfo); err == nil {
					log.Printf("Send the information of processors to BMC.")
				} else {
					log.Printf("Failed to send the information of processors to BMC: %v.", err)
				}
			} else {
				log.Printf("Failed to get the information of Processors: %v.", err)
			}

			BootDriveInfo, err := ocp.GetOemIpmiBootDriveInfo(si)
			if err == nil {
				if BootDriveInfo != nil {
					if err = ocp.SendOemIpmiBootDriveInfo(i, BootDriveInfo); err == nil {
						log.Printf("Send the information of boot drive to BMC.")
					} else {
						log.Printf("Failed to send the information of boot drive to BMC: %v.", err)
					}
				} else {
					log.Printf("The information of boot drive is not found.")
				}
			} else {
				log.Printf("Failed to get the information of boot drive: %v.", err)
			}

			if err = ocp.SetOemIpmiPostEnd(i); err == nil {
				log.Printf("Send IPMI POST end to BMC")
			} else {
				log.Printf("Failed to send IPMI POST end to BMC: %v.", err)
			}

		} else {
			log.Printf("No product name is matched for OEM commands.")
		}
	}
}

// Add an event to the IPMI System Even Log
func addSEL(sequence string) {
	var bootErr ipmi.Event

	i, err := ipmi.Open(0)
	if err != nil {
		log.Printf("Failed to open ipmi device to send SEL %v", err)
		return
	}
	defer i.Close()

	switch sequence {
	case "netboot":
		fallthrough
	case "fbnetboot", "pxeboot":
		bootErr.RecordID = 0
		bootErr.RecordType = ipmi.OEM_NTS_TYPE
		bootErr.OEMNontsDefinedData[0] = 0x28
		bootErr.OEMNontsDefinedData[5] = 0xf0
		for idx := 6; idx < 13; idx++ {
			bootErr.OEMNontsDefinedData[idx] = 0xff
		}
		if err := i.LogSystemEvent(&bootErr); err != nil {
			log.Printf("SEL recorded: %s fail\n", sequence)
		}
		selRecorded = true
	case "cmosclear":
		bootErr.RecordID = 0
		bootErr.RecordType = ipmi.OEM_NTS_TYPE
		bootErr.OEMNontsDefinedData[0] = 0x28
		bootErr.OEMNontsDefinedData[5] = 0xf1
		for idx := 6; idx < 13; idx++ {
			bootErr.OEMNontsDefinedData[idx] = 0xff
		}
		if err := i.LogSystemEvent(&bootErr); err != nil {
			log.Printf("SEL recorded: %s fail\n", sequence)
		}
	default:
	}
}

// getDebugEnabled checks whether debug output is requested, either via command line or via VPD
// variables.
// If -q was explicitly passed on the command line, will use that value, otherwise will look for
// the VPD variable "systemboot_log_level".
// Valid values are coreboot loglevels https://review.coreboot.org/cgit/coreboot.git/tree/src/commonlib/include/commonlib/loglevel.h,
// either as integer (1, 2) or string ("debug").
// If the VPD variable is missing or it is set to an invalid value, it will use the default.
func getDebugEnabled() bool {
	if isFlagPassed("q") {
		return !*doQuiet
	}

	// -q was not passed, so `doQuiet` contains the default value
	defaultDebugEnabled := !*doQuiet
	// check for the VPD variable "systemboot_log_level". First the read-write, then the read-only
	v, err := vpd.Get(vpdSystembootLogLevel, false)
	if err != nil {
		// TODO do not print warning if file is not found
		log.Printf("Warning: failed to read read-write VPD variable '%s', will try the read-only one. Error was: %v", vpdSystembootLogLevel, err)
		v, err = vpd.Get(vpdSystembootLogLevel, true)
		if err != nil {
			// TODO do not print warning if file is not found
			log.Printf("Warning: failed to read read-only VPD variable '%s', will use the default value. Error was: %v", vpdSystembootLogLevel, err)
			return defaultDebugEnabled
		}
	}
	level := strings.ToLower(strings.TrimSpace(string(v)))
	switch level {
	case "0", "emerg", "1", "alert", "2", "crit", "3", "err", "4", "warning", "9", "never":
		// treat as quiet
		return false
	case "5", "notice", "6", "info", "7", "debug", "8", "spew":
		// treat as debug
		return true
	default:
		log.Printf("Invalid value '%s' for VPD variable '%s', using default", level, vpdSystembootLogLevel)
		return defaultDebugEnabled
	}
}

func main() {
	flag.Parse()

	debugEnabled := getDebugEnabled()

	log.Print(`
                     ____            _                 _                 _
                    / ___| _   _ ___| |_ ___ _ __ ___ | |__   ___   ___ | |_
                    \___ \| | | / __| __/ _ \ '_ ` + "`" + ` _ \| '_ \ / _ \ / _ \| __|
                     ___) | |_| \__ \ ||  __/ | | | | | |_) | (_) | (_) | |_
                    |____/ \__, |___/\__\___|_| |_| |_|_.__/ \___/ \___/ \__|
                           |___/
`)
	l := ulog.Null
	if debugEnabled {
		l = ulog.Log
	}
	runIPMICommands(l)
	sleepInterval := time.Duration(*interval) * time.Second
	if *allowInteractive {
		log.Printf("**************************************************************************")
		log.Print("Starting boot sequence, press CTRL-C within 5 seconds to drop into a shell")
		log.Printf("**************************************************************************")
		time.Sleep(5 * time.Second)
	} else {
		signal.Ignore()
	}

	// Get and show boot entries
	var bootEntries []systembooter.BootEntry
	if bmcBootOverride && ocp.BmcUpdatedBootorder {
		bootEntries = ocp.BootEntries
	} else {
		bootEntries = systembooter.GetBootEntries(l)
	}
	log.Printf("BOOT ENTRIES:")
	for _, entry := range bootEntries {
		log.Printf("    %v : %+v", entry.Name, string(entry.Config))
	}
	for _, entry := range bootEntries {
		log.Printf("Trying boot entry %s: %s", entry.Name, string(entry.Config))
		if err := entry.Booter.Boot(debugEnabled); err != nil {
			log.Printf("Warning: failed to boot with configuration: %s: %s", entry.Name, string(entry.Config))
			addSEL(entry.Booter.TypeName())
		}
		if debugEnabled {
			log.Printf("Sleeping %v before attempting next boot command", sleepInterval)
		}
		time.Sleep(sleepInterval)
	}

	// if boot entries failed, use the default boot sequence
	log.Printf("Boot entries failed")

	if !*noDefaultBoot {
		log.Print("Falling back to the default boot sequence")
		for {
			for _, bootcmd := range defaultBootsequence {
				if _, err := exec.LookPath(bootcmd[0]); err != nil {
					log.Printf("No Path: %v", bootcmd)
					continue
				}
				if debugEnabled {
					bootcmd = append(bootcmd, "-v")
				}
				log.Printf("Running boot command: %v", bootcmd)
				cmd := exec.Command(bootcmd[0], bootcmd[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					// MJ TODO - Need a fix for booters with menues that fail and drop to menu.
					log.Printf("Error executing %v: %v", cmd, err)
					if !selRecorded {
						addSEL(bootcmd[0])
					}
				}
			}
			selRecorded = true

			if debugEnabled {
				log.Printf("Sleeping %v before attempting next boot command", sleepInterval)
			}
			time.Sleep(sleepInterval)
		}
	}
}
