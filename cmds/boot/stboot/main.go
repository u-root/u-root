package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/u-root/u-root/pkg/bootconfig"
	"github.com/u-root/u-root/pkg/loop"
	"github.com/u-root/u-root/pkg/storage"
	"golang.org/x/sys/unix"
)

// TODO:
// implement booter interface

var (
	dryRun  = flag.Bool("dryrun", false, "Do everything except booting the loaded kernel")
	doDebug = flag.Bool("d", false, "Print debug output")
)

const (
	eth            = "eth0"
	bootFilePath   = "root/stboot.zip"
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

	MinimalAmountSignatures int `json:"minimal-amount-signatures"`
}

func stbootSetupIOFromNetVars(vars netVars) error {
	//setup ip
	log.Print("Setup network configuration with IP: " + vars.HostIP)
	cmd := exec.Command("ip", "addr", "add", vars.HostIP, "dev", eth)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing %v: %v", cmd, err)
		return err
	}
	cmd = exec.Command("ip", "link", "set", eth, "up")
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing %v: %v", cmd, err)
		return err
	}
	cmd = exec.Command("ip", "route", "add", "default", "via", vars.DefaultGateway, "dev", eth)
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing %v: %v", cmd, err)
		return err
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

	return nil
}

func stbootVerifySignatureInPath(path string, hashValue []byte, opts x509.VerifyOptions, minAmountValid int) error {
	validSignatures := 0

	// Check certs and signatures
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && (filepath.Ext(info.Name()) == ".cert") {
			// Read cert and verify
			userCert, err := ioutil.ReadFile(path)
			if err == nil {
				block, _ := pem.Decode(userCert)
				cert, err := x509.ParseCertificate(block.Bytes)
				if err == nil {
					// verify certificates with root certificate
					_, err = cert.Verify(opts)
					if err == nil {
						// Read signature and verify it.
						signatureFilename := strings.TrimSuffix(path, filepath.Ext(path)) + ".signature"
						signatureRaw, err := ioutil.ReadFile(signatureFilename)

						if err != nil {
							log.Println(fmt.Sprintf("Unable to read signature at %s. Erroring.", signatureFilename))
							return err
						}

						opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash}
						err = rsa.VerifyPSS(cert.PublicKey.(*rsa.PublicKey), crypto.SHA512, hashValue, signatureRaw, opts)
						if err != nil {
							log.Println(fmt.Sprintf("Signature Verification failed for %s.", filepath.Base(signatureFilename)))
						} else {
							validSignatures++
							debug(fmt.Sprintf("%s verfied.", signatureFilename))
						}
					} else {
						log.Fatal(err)
					}
				} else {
					log.Fatal(err)
				}
			} else {
				log.Fatal(fmt.Sprintf("Unable to read user certificate %s", path))
			}
		}

		return nil
	})

	if validSignatures < minAmountValid {
		log.Fatalf("Did not found enough valid signatures. Only %d (%d required) are valid.", validSignatures, minAmountValid)
		return errors.New(("Not enough valid signatures found."))
	}

	return nil
}

func stbootMountIso(pathToIso string, mountPoint string) (string, string, error) {

	// Mount the iso
	log.Println(fmt.Sprintf("Trying to mount %s in /tmp/iso", pathToIso))
	os.MkdirAll(mountPoint, os.ModeDir|os.FileMode(0700))
	var flags = uintptr(unix.UMOUNT_NOFOLLOW)
	flags |= unix.MNT_FORCE

	device, err := loop.New(pathToIso, mountPoint, "iso9660", flags, "")
	if err != nil {
		log.Println(fmt.Sprintf("%v", err))
		return "", "", err
	}
	device.Mount()
	log.Println("Mounted.")

	kernelPath := path.Join(mountPoint, "vmlinuz")
	initramfsPath := path.Join(mountPoint, "initramf")

	return kernelPath, initramfsPath, nil
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
	// FIXME if already mounted - cant find netvars.json
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
		log.Print("MinimalAmountSignatures: ", vars.MinimalAmountSignatures)
	}

	// Setup IO from NetVars
	err = stbootSetupIOFromNetVars(vars)
	if err != nil {
		log.Fatal("Can not set up IO from NetVars")
		log.Fatal(err)
		return
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

	// check signature if necessary and unpck
	manifest, outputDir, err := bootconfig.FromZip(bootFilePath)
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

	// Hash payload -  right now it is the kernel
	hash := sha512.New()
	hash.Reset()
	kernelRaw, _ := ioutil.ReadFile(cfg.Kernel)
	hash.Write(kernelRaw)
	kernelHash := hash.Sum(nil)

	// Verify Boot Config
	// Search for signatures to this bootconfig
	certPath := strings.Replace(path.Dir(manifest.Configs[0].Kernel), outputDir, "", -1)
	certPath = path.Join(outputDir, "certs/", certPath)

	// TODO: Check if path really exists

	rootCert, err := ioutil.ReadFile(path.Join(outputDir, "certs/root.cert"))
	if err != nil {
		log.Println("Root Certificate not found.")
		return
	}

	root := x509.NewCertPool()
	ok = root.AppendCertsFromPEM(rootCert)
	if !ok {
		log.Fatalln("Failed to parse root certificate")
		return
	}

	opts := x509.VerifyOptions{
		Roots: root,
	}

	err = stbootVerifySignatureInPath(certPath, kernelHash, opts, vars.MinimalAmountSignatures)

	if err != nil {
		log.Fatal("The bootconfig seems to be not trustworthy.")
		return
	}

	if *dryRun {
		debug("Dryrun mode: will not boot")
		return
	}

	tmpPath, err := ioutil.TempDir(os.TempDir(), "iso")
	if err != nil {
		log.Fatalf("Unable to create temporary dir in %v", err)
		return
	}
	kernelPath, initramfsPath, err := stbootMountIso(cfg.Kernel, tmpPath)

	if err != nil || kernelPath == "" || initramfsPath == "" {
		log.Fatalln("Error Mounting Iso.")
		return
	}

	// Extend arguments.
	cfg.KernelArgs = cfg.KernelArgs + " root=/var/squashfs/filesystem.squashfs"
	cfg.Kernel = kernelPath
	cfg.Initramfs = initramfsPath

	log.Printf("%v", cfg)

	log.Println("Starting up new kernel.")

	// boot
	if err := cfg.Boot(); err != nil {
		log.Printf("Failed to boot kernel %s: %v", cfg.Kernel, err)
	}
	// if we reach this point, no boot configuration succeeded
	log.Print("No boot configuration succeeded")

	return
}
