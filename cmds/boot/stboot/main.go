// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/boot/stboot"
	"github.com/u-root/u-root/pkg/recovery"
	"github.com/u-root/u-root/pkg/ulog"
)

var (
	dryRun  = flag.Bool("dryrun", false, "Do everything except booting the loaded kernel")
	doDebug = flag.Bool("debug", false, "Print additional debug output")
	klog    = flag.Bool("klog", false, "Print output to all attached consoles via the kernel log")

	debug = func(string, ...interface{}) {}

	data dataPartition
)

const (
	provisioningServerFile = "provisioning-servers.json"
	networkFile            = "network.json"
	httpsRootsFile         = "https-root-certificates.pem"
	ntpServerFile          = "ntp-servers.json"
)

var banner = `
  _____ _______   _____   ____   ____________
 / ____|__   __|  |  _ \ / __ \ / __ \__   __|
| (___    | |     | |_) | |  | | |  | | | |   
 \___ \   | |     |  _ <| |  | | |  | | | |   
 ____) |  | |     | |_) | |__| | |__| | | |   
|_____/   |_|     |____/ \____/ \____/  |_|   

`

var check = `           
           //\\
OS is     //  \\
valid    //   //
        //   //
 //\\  //   //
//  \\//   //
\\        //
 \\      //
  \\    //
   \\__//
`

func main() {
	log.SetPrefix("stboot: ")
	ulog.KernelLog.SetLogLevel(ulog.KLogNotice)
	ulog.KernelLog.SetConsoleLogLevel(ulog.KLogInfo)

	flag.Parse()
	if *doDebug {
		debug = info
	}

	info(banner)

	vars, err := loadHostvars()
	if err != nil {
		reboot("Cannot find hostvars: %v", err)
	}
	if *doDebug {
		str, _ := json.MarshalIndent(vars, "", "  ")
		info("Host variables: %s", str)
	}

	/////////////////
	// Data partition
	/////////////////

	data, err = findDataPartition()
	if err != nil {
		reboot("%v", err)
	}

	//////////
	// Network
	//////////
	nc, err := getNetConf()
	if err != nil {
		debug("Cannot read network configuration file: %v", err)
		err = configureDHCPNetwork()
		if err != nil {
			reboot("Cannot set up IO: %v", err)
		}
	}

	if nc.HostIP != "" && nc.DefaultGateway != "" {
		if *doDebug {
			str, _ := json.MarshalIndent(nc, "", "  ")
			info("Network configuration: %s", str)
		}
		err = configureStaticNetwork(nc)
	} else {
		debug("no configuration specified in %s", networkFile)
		err = configureDHCPNetwork()
	}
	if err != nil {
		reboot("Cannot set up IO: %v", err)
	}

	hwAddr, err := hostHWAddr()
	if err != nil {
		reboot("%v", err)
	}
	info("Host's HW address: %s", hwAddr.String())

	////////////////////
	// Time validatition
	////////////////////
	if vars.Timestamp == 0 && *doDebug {
		info("WARNING: No timestamp found in hostvars")
	}
	buildTime := time.Unix(int64(vars.Timestamp), 0)
	err = validateSystemTime(buildTime)
	if err != nil {
		reboot("%v", err)
	}

	////////////////////
	// Download bootball
	////////////////////

	bytes, err := data.get(provisioningServerFile)
	if err != nil {
		reboot("Bootstrap URLs: %v", err)
	}
	var urlStrings []string
	if err = json.Unmarshal(bytes, &urlStrings); err != nil {
		reboot("Bootstrap URLs: %v", err)
	}
	if err = forceHTTPS(urlStrings); err != nil {
		reboot("Bootstrap URLs: %v", err)
	}

	info("Try downloading individual bootball")
	file := stboot.ComposeIndividualBallName(hwAddr)
	dest, err := tryDownload(urlStrings, file)
	if err != nil {
		debug("%v", err)
		info("Try downloading general bootball")
		dest, err = tryDownload(urlStrings, stboot.BallName)
		if err != nil {
			debug("%v", err)
			reboot("Cannot get appropriate bootball from provisioning servers")
		}
	}

	ball, err := stboot.BootBallFromArchive(dest)
	if err != nil {
		reboot("%v", err)
	}

	////////////////////////////////////////////////
	// Validate bootball's signing root certificates
	////////////////////////////////////////////////
	if len(vars.Fingerprints) == 0 {
		reboot("No root certificate fingerprints found in hostvars")
	}
	fp := calculateFingerprint(ball.RootCertPEM)
	info("Fingerprint of boot ball's root certificate:")
	info(fp)
	if !fingerprintIsValid(fp, vars.Fingerprints) {
		reboot("Root certificate of boot ball does not match expacted fingerprint")
	}
	info("OK!")

	////////////////////////////
	// Verify boot configuration
	////////////////////////////
	info("Pick the first boot configuration")
	var index = 0 // Just choose the first Bootconfig for now
	bc, err := ball.GetBootConfigByIndex(index)
	if err != nil {
		reboot("Cannot get boot configuration %d: %v", index, err)
	}

	if *doDebug {
		str, _ := json.MarshalIndent(*bc, "", "  ")
		info("Bootconfig (ID: %s): %s", bc.ID(), str)
	}

	n, valid, err := ball.VerifyBootconfigByID(bc.ID())
	if err != nil {
		reboot("Error verifying bootconfig %d: %v", index, err)
	}
	if valid < vars.MinimalSignaturesMatch {
		reboot("Did not found enough valid signatures: %d found, %d valid, %d required", n, valid, vars.MinimalSignaturesMatch)
	}

	debug("Signatures: %d found, %d valid, %d required", n, valid, vars.MinimalSignaturesMatch)
	info("Bootconfig '%s' passed verification", bc.Name)
	info(check)

	if *dryRun {
		debug("Dryrun mode: will not boot")
		return
	}
	//////////
	// Boot OS
	//////////
	info("Starting up new kernel.")

	if err := bc.Boot(); err != nil {
		reboot("Failed to boot kernel %s: %v", bc.Kernel, err)
	}
	// if we reach this point, no boot configuration succeeded
	reboot("No boot configuration succeeded")
}

// fingerprintIsValid returns true if fpHex is equal to on of
// those in expectedHex.
func fingerprintIsValid(fpHex string, expectedHex []string) bool {
	for _, f := range expectedHex {
		f = strings.TrimSpace(f)
		if fpHex == f {
			return true
		}
	}
	return false
}

// calculateFingerprint returns the SHA256 checksum of the
// provided certificate.
func calculateFingerprint(pemBytes []byte) string {
	block, _ := pem.Decode(pemBytes)
	fp := sha256.Sum256(block.Bytes)
	str := hex.EncodeToString(fp[:])
	return strings.TrimSpace(str)
}

//reboot trys to reboot the system in an infinity loop
func reboot(format string, v ...interface{}) {
	if *klog {
		info(format, v...)
		info("REBOOT!")
	}
	for {
		recover := recovery.SecureRecoverer{
			Reboot:   true,
			Debug:    true,
			RandWait: true,
		}
		err := recover.Recover(fmt.Sprintf(format, v...))
		if err != nil {
			continue
		}
	}
}

func info(format string, v ...interface{}) {
	if *klog {
		ulog.KernelLog.Printf("stboot: "+format, v...)
	} else {
		log.Printf(format, v...)
	}
}
