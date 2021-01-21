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
	"io/ioutil"
	"log"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/boot/stboot"
	"github.com/u-root/u-root/pkg/recovery"
)

var debug = func(string, ...interface{}) {}

var (
	dryRun  = flag.Bool("dryrun", false, "Do everything except booting the loaded kernel")
	doDebug = flag.Bool("d", false, "Print debug output")
)

const (
	rootCACertPath          = "/root/LetsEncrypt_Authority_X3.pem"
	rootCertFingerprintPath = "root/signing_rootcert.fingerprint"
	entropyAvail            = "/proc/sys/kernel/random/entropy_avail"
	interfaceUpTimeout      = 6 * time.Second
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

	flag.Parse()
	if *doDebug {
		debug = log.Printf
	}
	log.Print(banner)

	vars, err := stboot.FindHostVarsInInitramfs()
	if err != nil {
		reboot("Cant find Netvars at all: %v", err)
	}

	if *doDebug {
		str, _ := json.MarshalIndent(vars, "", "  ")
		log.Printf("Host variables: %s", str)
	}

	if vars.HostIP != "" {
		err = configureStaticNetwork(vars)
	} else {
		err = configureDHCPNetwork()
	}

	if err != nil {
		reboot("Can not set up IO: %v", err)
	}

	err = validateSystemTime()
	if err != nil {
		reboot("%v", err)
	}

	ballPath := path.Join("root/", stboot.BallName)
	url, err := url.Parse(vars.BootstrapURL)
	if err != nil {
		reboot("Invalid bootstrap URL: %v", err)
	}
	url.Path = path.Join(url.Path, stboot.BallName)
	err = downloadFromHTTPS(url.String(), ballPath)
	if err != nil {
		reboot("Downloading bootball failed: %v", err)
	}

	ball, err := stboot.BootBallFromArchive(ballPath)
	if err != nil {
		reboot("Cannot open bootball: %v", err)
	}

	fp, err := ioutil.ReadFile(rootCertFingerprintPath)
	if err != nil {
		reboot("Cannot read fingerprint: %v", err)
	}

	if *doDebug {
		log.Print("Fingerprint of boot ball's root certificate:")
		log.Print(string(fp))
	}
	if !matchFingerprint(ball.RootCertPEM, string(fp)) {
		reboot("Root certificate of boot ball does not match expacted fingerprint %v", err)
	}

	// Just choose the first Bootconfig for now
	log.Printf("Pick the first boot configuration")
	var index = 0
	bc, err := ball.GetBootConfigByIndex(index)
	if err != nil {
		reboot("Cannot get boot configuration %d: %v", index, err)
	}

	if *doDebug {
		str, _ := json.MarshalIndent(*bc, "", "  ")
		log.Printf("Bootconfig (ID: %s): %s", bc.ID(), str)
	}

	n, valid, err := ball.VerifyBootconfigByID(bc.ID())
	if err != nil {
		reboot("Error verifying bootconfig %d: %v", index, err)
	}
	if valid < vars.MinimalSignaturesMatch {
		reboot("Did not found enough valid signatures: %d found, %d valid, %d required", n, valid, vars.MinimalSignaturesMatch)
	}

	if *doDebug {
		reboot("Signatures: %d found, %d valid, %d required", n, valid, vars.MinimalSignaturesMatch)
	}

	log.Printf("Bootconfig '%s' passed verification", bc.Name)
	log.Print(check)

	if *dryRun {
		debug("Dryrun mode: will not boot")
		return
	}

	log.Println("Starting up new kernel.")

	if err := bc.Boot(); err != nil {
		log.Printf("Failed to boot kernel %s: %v", bc.Kernel, err)
	}
	// if we reach this point, no boot configuration succeeded
	reboot("No boot configuration succeeded")
}

// matchFingerprint returns true if fingerprintHex matches the SHA256
// hash calculated from pem decoded certPEM.
func matchFingerprint(certPEM []byte, fingerprintHex string) bool {
	block, _ := pem.Decode(certPEM)
	fp := sha256.Sum256(block.Bytes)
	str := hex.EncodeToString(fp[:])
	str = strings.TrimSpace(str)

	fingerprintHex = strings.TrimSpace(fingerprintHex)

	return str == fingerprintHex
}

//reboot trys to reboot the system in an infinity loop
func reboot(format string, v ...interface{}) {
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
