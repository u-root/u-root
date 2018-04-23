package main

import (
	"flag"
	"log"
	"os/exec"
	"time"
)

// TODO allow to specify boot sequence

var (
	doDebug  = flag.Bool("d", false, "Enable debug mode for subcommands")
	interval = flag.Int("I", 0, "Interval in seconds before looping to the next boot command")
)

var bootsequence = [][]string{
	[]string{"netboot"},
	[]string{"localboot"},
}

func main() {
	flag.Parse()

	log.Print("Starting boot sequence, press CTRL-C within 5 seconds to drop into a shell")
	time.Sleep(5 * time.Second)
	if *doDebug {
		log.Print("Debug mode enabled")
	}

	sleepInterval := time.Duration(*interval) * time.Second
	for {
		for _, bootcmd := range bootsequence {
			if *doDebug {
				bootcmd = append(bootcmd, "-d")
			}
			log.Printf("Running boot command: %v", bootcmd)
			cmd := exec.Command(bootcmd[0], bootcmd[1:]...)
			if err := cmd.Run(); err != nil {
				log.Printf("Error executing %s: %v", cmd, err)
			}
		}
		if *doDebug {
			log.Printf("Sleeping %v before attempting next boot command", sleepInterval)
		}
		time.Sleep(sleepInterval)
	}

}
