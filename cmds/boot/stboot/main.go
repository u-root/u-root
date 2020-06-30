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
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/stboot"
	"github.com/u-root/u-root/pkg/crypto"
	"github.com/u-root/u-root/pkg/recovery"
	"github.com/u-root/u-root/pkg/ulog"
)

var (
	noMeasuredBoot = flag.Bool("no-measurement", false, "Do not extend PCRs with measurements of the loaded OS")
	doDebug        = flag.Bool("debug", false, "Print additional debug output")
	klog           = flag.Bool("klog", false, "Print output to all attached consoles via the kernel log")
	dryRun         = flag.Bool("dryrun", false, "Do everything except booting the loaded kernel")

	debug = func(string, ...interface{}) {}

	vars         *Hostvars
	txtSupported bool
)

// configuration files form STDATA partition
const (
	provisioningServerFile = "stboot/etc/provisioning-servers.json"
	networkFile            = "stboot/etc/network.json"
	httpsRootsFile         = "stboot/etc/https-root-certificates.pem"
	ntpServerFile          = "stboot/etc/ntp-servers.json"

	newDir              = "stboot/bootballs/new/"
	knownGoodDir        = "stboot/bootballs/known_good/"
	invalidDir          = "stboot/bootballs/invalid/"
	currentBootballFile = "stboot/bootballs/current"
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
	ulog.KernelLog.SetLogLevel(ulog.KLogNotice)
	ulog.KernelLog.SetConsoleLogLevel(ulog.KLogInfo)

	flag.Parse()
	if *doDebug {
		debug = info
	}

	info(banner)

	p := filepath.Join("etc/", hostvarsFile)
	var err error
	vars, err = loadHostvars(p)
	if err != nil {
		reboot("Cannot find hostvars: %v", err)
	}
	if *doDebug {
		str, _ := json.MarshalIndent(vars, "", "  ")
		info("Host variables: %s", str)
	}

	/////////////////
	// Data partition
	/////////////////
	err = findDataPartition()
	if err != nil {
		reboot("%v", err)
	}

	if vars.BootMode == NetworkStatic || vars.BootMode == NetworkDHCP {
		//////////
		// Network
		//////////
		nc, err := getNetConf()
		if err != nil {
			debug("Cannot read network configuration file: %v", err)
			err = configureDHCPNetwork()
			if err != nil {
				reboot("Cannot set up IO: %v", err)
			}
		}

		if nc.HostIP != "" && nc.DefaultGateway != "" {
			if *doDebug {
				str, _ := json.MarshalIndent(nc, "", "  ")
				info("Network configuration: %s", str)
			}
			err = configureStaticNetwork(nc)
		} else {
			debug("no configuration specified in %s", networkFile)
			err = configureDHCPNetwork()
		}
		if err != nil {
			reboot("Cannot set up IO: %v", err)
		}

		////////////////////
		// Time validatition
		////////////////////
		if vars.Timestamp == 0 && *doDebug {
			info("WARNING: No timestamp found in hostvars")
		}
		buildTime := time.Unix(int64(vars.Timestamp), 0)
		err = validateSystemTime(buildTime)
		if err != nil {
			reboot("%v", err)
		}
	}

	////////////////
	// TXT self test
	////////////////
	txtSupported = runTxtTests(*doDebug)
	if !txtSupported {
		info("WARNING: No TXT Support!")
	}
	info("TXT is supported on this platform")

	////////////////
	// Load bootball
	////////////////

	var bootballFiles []string

	switch vars.BootMode {
	case NetworkStatic, NetworkDHCP:
		f, err := loadBallFromNetwork()
		if err != nil {
			reboot("%v", err)
		}
		bootballFiles = append(bootballFiles, f)
	case LocalStorage:
		ff, err := loadBallFromLocalStorage()
		if err != nil {
			reboot("%v", err)
		}
		bootballFiles = append(bootballFiles, ff...)
	default:
		reboot("unknown boot mode: %s", vars.BootMode.string())
	}
	if len(bootballFiles) == 0 {
		reboot("No bootballs found")
	}
	if *doDebug {
		info("Bootballs to be processed:")
		for _, b := range bootballFiles {
			info(b)
		}
	}

	////////////////////
	// Process bootballs
	////////////////////
	var osi boot.OSImage
	for _, b := range bootballFiles {
		osi, err = processBootball(b)
		if err == nil {
			if vars.BootMode == LocalStorage {
				// indicate current bootball
				f := filepath.Join(dataMountPoint, currentBootballFile)
				rel, err := filepath.Rel(filepath.Dir(f), b)
				if err != nil {
					reboot("failed to indicate current bootball: %v", err)
				}
				if err = ioutil.WriteFile(f, []byte(rel), os.ModePerm); err != nil {
					reboot("failed to indicate current bootball: %v", err)
				}
			}
			break
		}
		debug("%v", err)
		if vars.BootMode == LocalStorage {
			// move invalid bootballs
			invalid := filepath.Join(dataMountPoint, invalidDir, filepath.Base(b))
			if err = stboot.CreateAndCopy(b, invalid); err != nil {
				reboot("failed to move invalid bootball: %v", err)
			}
			if err = os.Remove(b); err != nil {
				reboot("failed to move invalid bootball: %v", err)
			}
		}
	}
	info("Operating system: %s", osi.Label())
	debug("%s", osi.String())

	///////////////////////
	// Measure OS into PCRs
	///////////////////////
	if *noMeasuredBoot {
		info("WARNING: measured boot disabled!")
	} else {
		// TODO: measure osi byte stream not its label
		err = crypto.TryMeasureData(crypto.BootConfigPCR, []byte(osi.Label()), osi.Label())
		if err != nil {
			reboot("measured boot failed: %v", err)
		}
		// TODO: measure hostvars.json and files from data partition
	}

	//////////
	// Boot OS
	//////////
	if *dryRun {
		debug("Dryrun mode: will not boot")
		return
	}
	info("Loading operating system \n%s", osi.String())
	err = osi.Load(*doDebug)
	if err != nil {
		reboot("%s", err)
	}
	info("Handing over controll now")
	err = boot.Execute()
	if err != nil {
		reboot("%v", err)
	}

	reboot("unexpected return from kexec")

}

func processBootball(path string) (boot.OSImage, error) {
	info("Try loading bootball %s", path)
	ball, err := stboot.BootballFromArchive(path)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	////////////////////////////////////////
	// Validate bootball's root certificates
	////////////////////////////////////////
	if len(vars.Fingerprints) == 0 {
		return nil, fmt.Errorf("No root certificate fingerprints found in hostvars")
	}
	fp := calculateFingerprint(ball.RootCertPEM)
	info("Fingerprint of boot ball's root certificate:")
	info(fp)
	if !fingerprintIsValid(fp, vars.Fingerprints) {
		return nil, fmt.Errorf("Root certificate of boot ball does not match expacted fingerprint")
	}
	info("OK!")

	//////////////////
	// Verify bootball
	//////////////////
	if *doDebug {
		str, _ := json.MarshalIndent(ball.Config, "", "  ")
		info("Bootball config: %s", str)
	} else {
		info("Label: %s", ball.Config.Label)
	}

	n, valid, err := ball.Verify()
	if err != nil {
		return nil, fmt.Errorf("Error verifying bootball: %v", err)
	}
	if valid < vars.MinimalSignaturesMatch {
		return nil, fmt.Errorf("Not enough valid signatures: %d found, %d valid, %d required", n, valid, vars.MinimalSignaturesMatch)
	}

	debug("Signatures: %d found, %d valid, %d required", n, valid, vars.MinimalSignaturesMatch)
	info("Bootball passed verification")
	info(check)

	/////////////
	// Extract OS
	/////////////
	debug("Try extracting operating system with TXT")
	txt := true
	osiTXT, err := ball.OSImage(txt)
	if err != nil {
		info("%v", err)
	}
	debug("Try extracting non-TXT fallback operating system")
	osiFallback, err := ball.OSImage(!txt)
	if err != nil {
		info("%v", err)
	}

	if osiTXT == nil && osiFallback == nil {
		return nil, fmt.Errorf("Failed to get operating system from bootball")
	}

	var img boot.OSImage
	if txtSupported {
		if osiTXT == nil {
			info("WARNING: TXT will not be used!")
			img = osiFallback
		}
		img = osiTXT
	} else {
		if osiFallback == nil {
			return nil, fmt.Errorf("TXT is not supported by the host and no fallback OS is provided by the bootball")
		}
		info("WARNING: TXT will not be used!")
		img = osiFallback
	}

	return img, nil
}

func loadBallFromNetwork() (string, error) {
	p := filepath.Join(dataMountPoint, provisioningServerFile)
	bytes, err := ioutil.ReadFile(p)
	if err != nil {
		return "", fmt.Errorf("provisioning server URLs: %v", err)
	}
	var urlStrings []string
	if err = json.Unmarshal(bytes, &urlStrings); err != nil {
		return "", fmt.Errorf("provisioning server URLs URLs: %v", err)
	}
	if err = forceHTTPS(urlStrings); err != nil {
		return "", fmt.Errorf("provisioning server URLs URLs: %v", err)
	}

	info("Try downloading individual bootball")
	hwAddr, err := hostHWAddr()
	if err != nil {
		return "", fmt.Errorf("cannot evaluate hardware address: %v", err)
	}
	info("Host's HW address: %s", hwAddr.String())
	prefix := stboot.ComposeIndividualBallPrefix(hwAddr)
	file := prefix + stboot.DefaultBallName
	dest, err := tryDownload(urlStrings, file)
	if err != nil {
		debug("%v", err)
		info("Try downloading general bootball")
		dest, err = tryDownload(urlStrings, stboot.DefaultBallName)
		if err != nil {
			debug("%v", err)
			return "", fmt.Errorf("cannot get appropriate bootball from provisioning servers")
		}
	}

	return dest, nil
}

func loadBallFromLocalStorage() ([]string, error) {
	var bootballs []string
	var newBootballs []string
	var knownGoodBootballs []string

	//new Bootballs
	dir := filepath.Join(dataMountPoint, newDir)
	newBootballs, err := searchBootballFiles(dir)
	if err != nil {
		return nil, err
	}
	bootballs = append(bootballs, newBootballs...)

	// known good bootballs
	dir = filepath.Join(dataMountPoint, knownGoodDir)
	knownGoodBootballs, err = searchBootballFiles(dir)
	if err != nil {
		return nil, err
	}
	bootballs = append(bootballs, knownGoodBootballs...)
	return bootballs, nil
}

func searchBootballFiles(dir string) ([]string, error) {
	var ret []string
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, fi := range fis {
		// take *.stboot files only
		if filepath.Ext(fi.Name()) == ".stboot" {
			b := filepath.Join(dir, fi.Name())
			ret = append(ret, b)
		}
	}
	// reverse order
	for i := 0; i < len(ret)/2; i++ {
		j := len(ret) - i - 1
		ret[i], ret[j] = ret[j], ret[i]
	}
	return ret, nil
}

// fingerprintIsValid returns true if fpHex is equal to on of
// those in expectedHex.
func fingerprintIsValid(fpHex string, expectedHex []string) bool {
	for _, f := range expectedHex {
		f = strings.TrimSpace(f)
		if fpHex == f {
			return true
		}
	}
	return false
}

// calculateFingerprint returns the SHA256 checksum of the
// provided certificate.
func calculateFingerprint(pemBytes []byte) string {
	block, _ := pem.Decode(pemBytes)
	fp := sha256.Sum256(block.Bytes)
	str := hex.EncodeToString(fp[:])
	return strings.TrimSpace(str)
}

//reboot trys to reboot the system in an infinity loop
func reboot(format string, v ...interface{}) {
	if *klog {
		info(format, v...)
		info("REBOOT!")
	}
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

func info(format string, v ...interface{}) {
	if *klog {
		ulog.KernelLog.Printf("stboot: "+format, v...)
	} else {
		log.Printf(format, v...)
	}
}
