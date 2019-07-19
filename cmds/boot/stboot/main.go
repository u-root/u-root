package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/u-root/u-root/pkg/bootconfig"
	"github.com/u-root/u-root/pkg/storage"
)

// TODO:
// implement booter interface

var (
	dryRun  = flag.Bool("dryrun", false, "Do everything except booting the loaded kernel")
	doDebug = flag.Bool("d", false, "Print debug output")
)

const (
	eth          = "eth0"
	bootFilePath = "root/bz.zip"
	netVarsPath  = "netvars.json"
)

var banner = `
  _____ _______   ____   ____   ____ _______ 
 / ____|__   __| |  _ \ / __ \ / __ \__   __|
| (___    | |    | |_) | |  | | |  | | | |   
 \___ \   | |    |  _ <| |  | | |  | | | |   
 ____) |  | |    | |_) | |__| | |__| | | |   
|_____/   |_|    |____/ \____/ \____/  |_|   
											
`
var debug = func(string, ...interface{}) {}

type netVars struct {
	HostIP         string `json:"host_ip"`
	HostNetmask    string `json:"netmask"`
	DefaultGateway string `json:"gateway"`
	DNSServer      string `json:"dns"`

	HostPrivKey string `json:"host_priv_key"`
	HostPupKey  string `json:"host_pub_key"`

	BootstrapURL    string `json:"bootstrap_url"`
	SignaturePubKey string `json:"signature_pub_key"`
}

func main() {
	flag.Parse()
	if *doDebug {
		debug = log.Printf
	}
	log.Print(banner)

	// get block devices
	devices, err := storage.GetBlockStats()
	if err != nil {
		log.Fatal(err)
	}
	// print partition info
	if *doDebug {
		for _, dev := range devices {
			log.Printf("Device: %+v", dev)
		}
	}

	// get a list of supported file systems for real devices (i.e. skip nodev)
	debug("Getting list of supported filesystems")
	filesystems, err := storage.GetSupportedFilesystems()
	if err != nil {
		log.Fatal(err)
	}
	debug("Supported file systems: %v", filesystems)

	var mounted []storage.Mountpoint
	// try mounting all the available devices, with all the supported file
	// systems
	debug("trying to mount all the available block devices with all the supported file system types")
	mounted = make([]storage.Mountpoint, 0)
	for _, dev := range devices {
		devname := path.Join("/dev", dev.Name)
		mountpath := path.Join("/mnt", dev.Name)
		if mountpoint, err := storage.Mount(devname, mountpath, filesystems); err != nil {
			debug("Failed to mount %s on %s: %v", devname, mountpath, err)
		} else {
			mounted = append(mounted, *mountpoint)
		}
	}
	log.Printf("mounted: %+v", mounted)
	defer func() {
		// clean up
		for _, mountpoint := range mounted {
			syscall.Unmount(mountpoint.Path, syscall.MNT_DETACH)
		}
	}()

	// search for a netvars.json
	var data []byte
	for _, mountpoint := range mounted {
		path := path.Join(mountpoint.Path, netVarsPath)
		log.Printf("Trying to read %s", path)
		data, err = ioutil.ReadFile(path)
		if err == nil {
			break
		}
		log.Printf("cannot open %s: %v", path, err)
	}

	log.Printf("Parse network variables")
	vars := netVars{}
	json.Unmarshal(data, &vars)
	// FIXME: : error handling
	// print network variables
	if *doDebug {
		log.Print("HostIP: " + vars.HostIP)
		log.Print("HostNetmask: " + vars.HostNetmask)
		log.Print("DefaultGateway: " + vars.DefaultGateway)
		log.Print("DNSServer: " + vars.DNSServer)

		log.Print("HostPrivKey: " + vars.HostPrivKey)
		log.Print("HostPubKey: " + vars.HostPupKey)

		log.Print("BootstrapURL: " + vars.BootstrapURL)
		log.Print("SignaturePupKey: " + vars.SignaturePubKey)
	}

	//setup ip
	log.Print("Setup network configuration with IP: " + vars.HostIP)
	cmd := exec.Command("ip", "addr", "add", vars.HostIP, "dev", eth)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing %v: %v", cmd, err)
	}
	cmd = exec.Command("ip", "link", "set", eth, "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing %v: %v", cmd, err)
	}
	cmd = exec.Command("ip", "route", "add", "default", "via", vars.DefaultGateway, "dev", eth)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing %v: %v", cmd, err)
	}

	if *doDebug {
		cmd = exec.Command("ip", "addr")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Error executing %v: %v", cmd, err)
		}
		cmd = exec.Command("ip", "route")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Error executing %v: %v", cmd, err)
		}
	}

	// get remote boot bundle
	log.Print("Get boot files from " + vars.BootstrapURL)
	cmd = exec.Command("wget", "-O", bootFilePath, vars.BootstrapURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing %v: %v", cmd, err)
	}
	// create pup_key.pem
	p := path.Join(os.TempDir(), "pub_key.pem")
	debug("Write public key from netvars.json to %s", p)
	debug("Public key is: %s", vars.SignaturePubKey)
	// If the file doesn't exist, create it, or append to the file
	f, err := os.Create(p)
	if err != nil {
		log.Fatalf("Faild to create public key file: %+v", err)
	}

	_, err = f.WriteString("-----BEGIN PUBLIC KEY-----\n")
	_, err = f.WriteString(vars.SignaturePubKey + "\n")
	_, err = f.WriteString("-----END PUBLIC KEY-----\n")
	if err != nil {
		log.Fatalf("Faild to write public key file: %+v", err)
	}
	f.Close()

	// err = ioutil.WriteFile(p, []byte("-----BEGIN PUBLIC KEY-----"+vars.SignaturePubKey+"-----END PUBLIC KEY-----"), 0555)
	// if err != nil {
	// 	panic(err)
	// }

	// check signature and unpck
	manifest, outputDir, err := bootconfig.FromZip(bootFilePath, &p)
	if err != nil {
		log.Fatal(err)
	}
	debug("Boot files unpacked into: " + outputDir)
	debug("Manifest: %+v", *manifest)
	// get first bootconfig from manifest
	cfg, err := manifest.GetBootConfig(0)
	if err != nil {
		log.Fatal(err)
	}
	debug("Bootconfig: %+v", *cfg)

	// update paths
	cfg.Kernel = path.Join(outputDir, cfg.Kernel)
	if cfg.Initramfs != "" {
		cfg.Initramfs = path.Join(outputDir, cfg.Initramfs)
	}
	if cfg.DeviceTree != "" {
		cfg.Initramfs = path.Join(outputDir, cfg.DeviceTree)
	}
	debug("Adjusted Bootconfig: %+v", *cfg)

	if *dryRun {
		debug("Dryrun mode: will not boot")
		return
	}
	// boot
	if err := cfg.Boot(); err != nil {
		log.Printf("Failed to boot kernel %s: %v", cfg.Kernel, err)
	}
	// if we reach this point, no boot configuration succeeded
	log.Print("No boot configuration succeeded")

	return
}
