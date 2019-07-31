// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/u-root/u-root/pkg/booter"
)

var (
	allowInteractive = flag.Bool("i", true, "Allow user to interrupt boot process and run commands")
	doQuiet          = flag.Bool("q", false, "Disable verbose output")
	interval         = flag.Int("I", 1, "Interval in seconds before looping to the next boot command")
	noDefaultBoot    = flag.Bool("nodefault", false, "Do not attempt default boot entries if regular ones fail")
)

var defaultBootsequence = [][]string{
	{"fbnetboot", "-userclass", "linuxboot"},
	{"localboot", "-grub"},
}

func main() {
	flag.Parse()

	log.Print(`
                     ____            _                 _                 _   
                    / ___| _   _ ___| |_ ___ _ __ ___ | |__   ___   ___ | |_ 
                    \___ \| | | / __| __/ _ \ '_ ` + "`" + ` _ \| '_ \ / _ \ / _ \| __|
                     ___) | |_| \__ \ ||  __/ | | | | | |_) | (_) | (_) | |_ 
                    |____/ \__, |___/\__\___|_| |_| |_|_.__/ \___/ \___/ \__|
                           |___/
`)

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
	bootEntries := booter.GetBootEntries()
	log.Printf("BOOT ENTRIES:")
	for _, entry := range bootEntries {
		log.Printf("    %v) %+v", entry.Name, string(entry.Config))
	}
	for _, entry := range bootEntries {
		log.Printf("Trying boot entry %s: %s", entry.Name, string(entry.Config))
		if err := entry.Booter.Boot(); err != nil {
			log.Printf("Warning: failed to boot with configuration: %+v", entry)
		}
		if !*doQuiet {
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
				if !*doQuiet {
					bootcmd = append(bootcmd, "-d")
				}
				log.Printf("Running boot command: %v", bootcmd)
				cmd := exec.Command(bootcmd[0], bootcmd[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					log.Printf("Error executing %v: %v", cmd, err)
				}
			}
			if !*doQuiet {
				log.Printf("Sleeping %v before attempting next boot command", sleepInterval)
			}
			time.Sleep(sleepInterval)
		}
	}
}
