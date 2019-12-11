package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"path"

	"github.com/u-root/u-root/pkg/boot/stboot"
)

var debug = func(string, ...interface{}) {}

var (
	dryRun  = flag.Bool("dryrun", false, "Do everything except booting the loaded kernel")
	doDebug = flag.Bool("d", false, "Print debug output")
)

var banner = `
  _____ _______   _____   ____   ____________
 / ____|__   __|  |  _ \ / __ \ / __ \__   __|
| (___    | |     | |_) | |  | | |  | | | |   
 \___ \   | |     |  _ <| |  | | |  | | | |   
 ____) |  | |     | |_) | |__| | |__| | | |   
|_____/   |_|     |____/ \____/ \____/  |_|   

`

func main() {
	flag.Parse()
	if *doDebug {
		debug = log.Printf
	}
	log.Print(banner)

	vars, err := stboot.FindHostVarsInInitramfs()
	if err != nil {
		log.Fatalf("Cant find Netvars at all: %v", err)
	}

	if *doDebug {
		str, _ := json.MarshalIndent(vars, "", "  ")
		log.Printf("Host variables: %s", str)
	}

	debug("Configuring network interfaces")

	if vars.HostIP != "" {
		err = stboot.ConfigureStaticNetwork(vars, *doDebug)
	} else {
		err = stboot.ConfigureDHCPNetwork()
	}

	if err != nil {
		log.Println("Can not set up IO.")
		log.Println(err)
		return
	}

	ballPath := path.Join("root/", stboot.BallName)
	url, err := url.Parse(vars.BootstrapURL)
	if err != nil {
		log.Printf("Invalid bootstrap URL: %v", err)
		return
	}
	url.Path = path.Join(url.Path, stboot.BallName)
	err = stboot.DownloadFromHTTPS(url.String(), ballPath)
	if err != nil {
		log.Printf("Error downloading bootball from %s", url)
		log.Println(err)
		return
	}

	ball, err := stboot.BootBallFromArchie(ballPath)
	if err != nil {
		log.Fatal("Cannot open bootball")
	}

	if err = ball.Verify(); err != nil {
		log.Fatal("The bootconfig seems to be not trustworthy. Err: ", err)
	}

	var index = 0
	bc, err := ball.GetBootConfigByIndex(index)
	if err != nil {
		log.Fatalf("Cannot get boot configuration %d: %v", index, err)
	}

	if *doDebug {
		str, _ := json.MarshalIndent(*bc, "", "  ")
		log.Printf("Bootconfig: %s", str)
	}

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

	return
}
