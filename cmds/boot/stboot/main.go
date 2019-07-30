package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/u-root/u-root/pkg/bootconfig"
	"github.com/u-root/u-root/pkg/storage"
)

// TODO:
// implement booter interface

var (
	verifyZip = flag.Bool("verify-zip", false, "If set, the archive will not be processed without a valid signature appended")
	dryRun    = flag.Bool("dryrun", false, "Do everything except booting the loaded kernel")
	doDebug   = flag.Bool("d", false, "Print debug output")
)

const (
	eth            = "eth0"
	bootFilePath   = "root/bc.zip"
	netVarsPath    = "netvars.json"
	rootCACertPath = "/root/LetsEncrypt_Authority_X3.pem"
	entropyAvail   = "/proc/sys/kernel/random/entropy_avail"
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

	// load CA certificate
	debug("Load %s as CA certificate", rootCACertPath)
	rootCertBytes, err := ioutil.ReadFile(rootCACertPath)
	if err != nil {
		log.Fatalf("Failed to read CA root certificate file: %s\n", err)
	}
	rootCertPem, _ := pem.Decode(rootCertBytes)
	if rootCertPem.Type != "CERTIFICATE" {
		log.Fatalf("Failed decoding certificate: Certificate is of the wrong type. PEM Type is: %s\n", rootCertPem.Type)
	}
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootCertBytes))
	if !ok {
		log.Fatalf("Error parsing CA root certificate")
	}
	debug("CA certificate: \n %s", string(rootCertBytes))

	// setup https client
	client := http.Client{
		Transport: (&http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: (&tls.Config{
				RootCAs: roots,
			}),
		}),
	}
	// check available kernel entropy
	e, err := ioutil.ReadFile(entropyAvail)
	es := strings.TrimSpace(string(e))
	entr, err := strconv.Atoi(es)
	if err != nil {
		log.Fatalf("Cannot evaluate entropy, %v", err)
	}
	debug("Available kernel entropy: %d", entr)
	if entr < 128 {
		log.Print("WARNING: low entropy!")
		log.Printf("%s : %d", entropyAvail, entr)
	}
	// get remote boot bundle
	log.Print("Get boot files from " + vars.BootstrapURL)
	resp, err := client.Get(vars.BootstrapURL)
	if err != nil {
		log.Fatalf("HTTTP GET failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("non-200 HTTP status: %d", resp.StatusCode)
	}
	f, err := os.Create(bootFilePath)
	if err != nil {
		log.Fatalf("Failed create boot config file: %v", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Fatalf("Failed to write boot config file: %v", err)
	}

	var zipPubkeyPath *string
	if *verifyZip {
		// create pup_key.pem
		zipPubkeyPath := path.Join(os.TempDir(), "pub_key.pem")
		debug("Write public key from netvars.json to %s", zipPubkeyPath)
		debug("Public key is: %s", vars.SignaturePubKey)
		t, err := os.Create(zipPubkeyPath)
		if err != nil {
			log.Fatalf("Failed to create public key file: %+v", err)
		}

		_, err = t.WriteString("-----BEGIN PUBLIC KEY-----\n")
		_, err = t.WriteString(vars.SignaturePubKey + "\n")
		_, err = t.WriteString("-----END PUBLIC KEY-----\n")
		if err != nil {
			log.Fatalf("Failed to write public key file: %+v", err)
		}
		t.Close()
	}

	// check signature if necessary and unpck
	manifest, outputDir, err := bootconfig.FromZip(bootFilePath, zipPubkeyPath)
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
