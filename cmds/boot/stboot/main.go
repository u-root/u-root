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

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/stboot"
	"github.com/u-root/u-root/pkg/recovery"
	"github.com/u-root/u-root/pkg/ulog"
)

var (
	noMeasuredBoot = flag.Bool("insecure", false, "Do not extend PCRs with measurements of the loaded OS")
	doDebug        = flag.Bool("debug", false, "Print additional debug output")
	klog           = flag.Bool("klog", false, "Print output to all attached consoles via the kernel log")
	dryRun         = flag.Bool("dryrun", false, "Do everything except booting the loaded kernel")

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

	////////////////
	// TXT self test
	////////////////
	txtSupported := runTxtTests(*doDebug)
	if !txtSupported {
		info("WARNING: No TXT Support!")
	}
	info("TXT is supported on this platform")

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

	ball, err := stboot.BootballFromArchive(dest)
	if err != nil {
		reboot("%v", err)
	}

	////////////////////////////////////////
	// Validate bootball's root certificates
	////////////////////////////////////////
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
	if *doDebug {
		str, _ := json.MarshalIndent(ball.Config, "", "  ")
		info("Bootball config: %s", str)
	} else {
		info("Label: %s", ball.Config.Label)
	}

	n, valid, err := ball.Verify()
	if err != nil {
		reboot("Error verifying bootball: %v", err)
	}
	if valid < vars.MinimalSignaturesMatch {
		reboot("Not enough valid signatures: %d found, %d valid, %d required", n, valid, vars.MinimalSignaturesMatch)
	}

	debug("Signatures: %d found, %d valid, %d required", n, valid, vars.MinimalSignaturesMatch)
	info("Bootball passed verification")
	info(check)

	/////////////////////////////
	// Measure bootball into PCRs
	/////////////////////////////
	// if !*noMeasuredBoot {
	// 	err = crypto.TryMeasureData(crypto.BootConfigPCR, ball.HashValue, ball.Config.Label)
	// 	if err != nil {
	// 		reboot("measured boot failed: %v", err)
	// 	}
	// 	// TODO: measure hostvars.json and files from data partition
	// }

	//////////
	// Boot OS
	//////////
	debug("Try extracting operating system with TXT")
	txt := true
	osiTXT, err := ball.OSImage(txt)
	if err != nil {
		debug("%v", err)
	}
	debug("Try extracting non-TXT fallback operating system")
	osiFallback, err := ball.OSImage(!txt)
	if err != nil {
		debug("%s", err)
	}

	if osiTXT == nil && osiFallback == nil {
		reboot("Failed to get operating system from bootball")
	}

	var osi boot.OSImage
	if txtSupported {
		if osiTXT == nil {
			info("WARNING: TXT will not be used!")
			osi = osiFallback
		}
		osi = osiTXT
	} else {
		if osiFallback == nil {
			reboot("TXT is not supported by the host and no fallback OS is provided by the bootball")
		}
		info("WARNING: TXT will not be used!")
		osi = osiFallback
	}

	info("Loading operating system \n%s", osi.String())
	err = osi.Load(*doDebug)
	if err != nil {
		reboot("%s", err)
	}

	if *dryRun {
		debug("Dryrun mode: will not boot")
		return
	}
	info("Handing over controll now")
	err = boot.Execute()
	if err != nil {
		reboot("%v", err)
	}

	reboot("unexpected return from kexec")
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
