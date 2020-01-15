// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"path"
	"time"

	"github.com/u-root/u-root/pkg/boot/stboot"
)

var debug = func(string, ...interface{}) {}

var (
	dryRun  = flag.Bool("dryrun", false, "Do everything except booting the loaded kernel")
	doDebug = flag.Bool("d", false, "Print debug output")
)

const (
	rootCACertPath     = "/root/LetsEncrypt_Authority_X3.pem"
	entropyAvail       = "/proc/sys/kernel/random/entropy_avail"
	interfaceUpTimeout = 6 * time.Second
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

	err := validateSystemTime()
	if err != nil {
		log.Fatal(err)
	}

	vars, err := stboot.FindHostVarsInInitramfs()
	if err != nil {
		log.Fatalf("Cant find Netvars at all: %v", err)
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
		log.Fatalf("Can not set up IO: %v", err)
	}

	ballPath := path.Join("root/", stboot.BallName)
	url, err := url.Parse(vars.BootstrapURL)
	if err != nil {
		log.Fatalf("Invalid bootstrap URL: %v", err)
	}
	url.Path = path.Join(url.Path, stboot.BallName)
	err = downloadFromHTTPS(url.String(), ballPath)
	if err != nil {
		log.Fatalf("Downloading bootball failed: %v", err)
	}

	ball, err := stboot.BootBallFromArchive(ballPath)
	if err != nil {
		log.Fatal("Cannot open bootball")
	}

	// Just choose the first Bootconfig for now
	log.Printf("Pick the first boot configuration")
	var index = 0
	bc, err := ball.GetBootConfigByIndex(index)
	if err != nil {
		log.Fatalf("Cannot get boot configuration %d: %v", index, err)
	}

	if *doDebug {
		str, _ := json.MarshalIndent(*bc, "", "  ")
		log.Printf("Bootconfig (ID: %s): %s", bc.ID(), str)
	}

	n, valid, err := ball.VerifyBootconfigByID(bc.ID())
	if err != nil {
		log.Fatalf("Error verifying bootconfig %d: %v", index, err)
	}
	if valid < vars.MinimalSignaturesMatch {
		log.Fatalf("Did not found enough valid signatures: %d found, %d valid, %d required", n, valid, vars.MinimalSignaturesMatch)
	}

	if *doDebug {
		log.Printf("Signatures: %d found, %d valid, %d required", n, valid, vars.MinimalSignaturesMatch)
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
	log.Print("No boot configuration succeeded")
}
