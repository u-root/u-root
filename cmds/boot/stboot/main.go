package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/u-root/u-root/pkg/bootconfig"
	"github.com/u-root/u-root/pkg/stboot"
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

	vars, err := stboot.FindNetVars()
	if err != nil {
		log.Fatalf("Cant find Netvars at all: %v", err)
	}

	// search for a netvars.json
	// FIXME if already mounted - cant find netvars.json

	// FIXME: : error handling
	// print network variables
	if *doDebug {
		log.Printf("Parse network variables")
		log.Print("HostIP: " + vars.HostIP)
		log.Print("HostNetmask: " + vars.HostNetmask)
		log.Print("DefaultGateway: " + vars.DefaultGateway)
		log.Print("DNSServer: " + vars.DNSServer)

		log.Print("HostPrivKey: " + vars.HostPrivKey)
		log.Print("HostPubKey: " + vars.HostPupKey)

		log.Print("BootstrapURL: " + vars.BootstrapURL)
		log.Print("SignaturePupKey: " + vars.SignaturePubKey)
		log.Print("MinimalAmountSignatures: ", vars.MinimalAmountSignatures)
	}

	debug("Configuring network interfaces")

	// If we do not have a HostIP we configure it dynamically
	if vars.HostIP != "" {
		// Setup IO from NetVars
		err = stboot.ConfigureStaticNetwork(vars, *doDebug)
	} else {
		err = stboot.ConfigureDHCPNetwork()
	}

	if err != nil {
		log.Println("Can not set up IO.")
		log.Println(err)
		return
	}

	err = stboot.DownloadFromHTTPS(vars.BootstrapURL, stboot.BootFilePath)
	if err != nil {
		log.Printf("Error verifing or download file from %s", vars.BootstrapURL)
		log.Println(err)
		return
	}

	// Unpack
	manifest, outputDir, err := bootconfig.FromZip(stboot.BootFilePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	debug("Boot files unpacked into: " + outputDir)

	// just take the first bootconfig
	// TODO: Should be loop through all bootconfigs?
	// TODO: Make sure 0 exists.

	// hash bootconfig
	dir := path.Join(outputDir, "bootconfig_0")
	hash, err := bootconfig.HashBootconfigDir(dir)
	log.Printf("bootconfig hash is: %x", hash[:])
	if err != nil {
		log.Printf("Error hashing bootconfig files in %s", dir)
		log.Println(err)
		return
	}

	cfg, err := manifest.GetBootConfig(0)
	if err != nil {
		log.Fatal(err)
	}
	if *doDebug {
		str, _ := json.MarshalIndent(*cfg, "", "  ")
		log.Printf("Bootconfig: %s", str)
	}

	// update paths
	cfg.Kernel = path.Join(outputDir, cfg.Kernel)
	kernelInfo, err := os.Stat(cfg.Kernel)
	if err != nil {
		log.Fatalf("cant read kernel file stats: %v", err)
		return
	}
	if kernelInfo.Size() == int64(0) {
		log.Fatalf("kernel file size is zero: %v", kernelInfo.Size())
		return
	}
	if cfg.Initramfs != "" {
		cfg.Initramfs = path.Join(outputDir, cfg.Initramfs)
		initRAMinfo, err := os.Stat(cfg.Initramfs)
		if err != nil {
			log.Fatalf("cant read initramfs file stats: %v", err)
			return
		}
		if initRAMinfo.Size() == int64(0) {
			log.Fatalf("initramfs file size is zero: %v", kernelInfo.Size())
			return
		}
	}
	if cfg.DeviceTree != "" {
		cfg.DeviceTree = path.Join(outputDir, cfg.DeviceTree)
		deviceTreeInfo, err := os.Stat(cfg.DeviceTree)
		if err != nil {
			log.Fatalf("cant read device tree file stats: %v", err)
			return
		}
		if deviceTreeInfo.Size() == int64(0) {
			log.Fatalf("device tree file size is zero: %v", deviceTreeInfo.Size())
			return
		}
	}

	if *doDebug {
		str, _ := json.MarshalIndent(*cfg, "", "  ")
		log.Printf("Adjusted Bootconfig: %s", str)
	}

	certPath := strings.Replace(path.Dir(manifest.Configs[0].Kernel), outputDir, "", -1)
	certPath = path.Join(outputDir, "certs/", certPath)

	if _, err = os.Stat(certPath); os.IsNotExist(err) {
		log.Fatalf("cert path does not exist: %v", err)
		return
	}

	rootCert, err := ioutil.ReadFile(path.Join(outputDir, "certs/root.cert"))
	if err != nil {
		log.Printf("Root Certificate not found: %v", err)
		return
	}
	err = stboot.VerifySignatureInPath(certPath, hash, rootCert, vars.MinimalAmountSignatures)

	if err != nil {
		log.Fatal("The bootconfig seems to be not trustworthy. Err: ", err)
		return
	}

	if *dryRun {
		debug("Dryrun mode: will not boot")
		return
	}

	log.Println("Starting up new kernel.")

	// boot
	if err := cfg.Boot(); err != nil {
		log.Printf("Failed to boot kernel %s: %v", cfg.Kernel, err)
	}
	// if we reach this point, no boot configuration succeeded
	log.Print("No boot configuration succeeded")

	return
}
