package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/u-root/u-root/pkg/boot/stboot"
	"github.com/u-root/u-root/pkg/bootconfig"
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

	dest := path.Join("root/", stboot.BallName)
	url, err := url.Parse(vars.BootstrapURL)
	if err != nil {
		log.Printf("Invalid bootstrap URL: %v", err)
		return
	}
	url.Path = path.Join(url.Path, stboot.BallName)
	err = stboot.DownloadFromHTTPS(url.String(), dest)
	if err != nil {
		log.Printf("Error downloading bootball from %s", url)
		log.Println(err)
		return
	}

	// Unpack
	cfg, outputDir, err := stboot.FromZip(dest)
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

	bc, err := cfg.GetBootConfig(0)
	if err != nil {
		log.Fatal(err)
	}
	if *doDebug {
		str, _ := json.MarshalIndent(*bc, "", "  ")
		log.Printf("Bootconfig: %s", str)
	}

	// update paths
	bc.Kernel = path.Join(outputDir, bc.Kernel)
	kernelInfo, err := os.Stat(bc.Kernel)
	if err != nil {
		log.Fatalf("cannot read kernel file stats: %v", err)
		return
	}
	if kernelInfo.Size() == int64(0) {
		log.Fatalf("kernel file size is zero: %v", kernelInfo.Size())
		return
	}
	if bc.Initramfs != "" {
		bc.Initramfs = path.Join(outputDir, bc.Initramfs)
		initRAMinfo, err := os.Stat(bc.Initramfs)
		if err != nil {
			log.Fatalf("cannot read initramfs file stats: %v", err)
			return
		}
		if initRAMinfo.Size() == int64(0) {
			log.Fatalf("initramfs file size is zero: %v", kernelInfo.Size())
		}
	}
	if bc.DeviceTree != "" {
		bc.DeviceTree = path.Join(outputDir, bc.DeviceTree)
		deviceTreeInfo, err := os.Stat(bc.DeviceTree)
		if err != nil {
			log.Fatalf("cannot read device tree file stats: %v", err)
			return
		}
		if deviceTreeInfo.Size() == int64(0) {
			log.Fatalf("device tree file size is zero: %v", deviceTreeInfo.Size())
			return
		}
	}

	if *doDebug {
		str, _ := json.MarshalIndent(*bc, "", "  ")
		log.Printf("Adjusted Bootconfig: %s", str)
	}

	certPath := strings.Replace(path.Dir(cfg.BootConfigs[0].Kernel), outputDir, "", -1)
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
	err = stboot.VerifySignatureInPath(certPath, hash, rootCert, vars.MinimalSignaturesMatch)

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
	if err := bc.Boot(); err != nil {
		log.Printf("Failed to boot kernel %s: %v", bc.Kernel, err)
	}
	// if we reach this point, no boot configuration succeeded
	log.Print("No boot configuration succeeded")

	return
}
