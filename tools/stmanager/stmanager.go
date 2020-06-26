// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// https://xkcd.com/927/

// stconfig is a configuration tool to create and manage artifacts for
// System Transparency Boot. Artifacts are ment to be uploaded to a
// remote provisioning server.

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	// Author is the author
	Author = "Jens Drenhaus"
	// HelpText is the command line help
	HelpText = "stmanager can be used for managing System Transparency bootballs"
)

var goversion string

var (
	create            = kingpin.Command("create", "Create a bootball from the provided config")
	createOut         = create.Flag("out", "Output directory of the bootball. If not set, current directory is used").ExistingDir()
	createLabel       = create.Flag("label", "Name of the boot configuration. Defaults to 'System Tarnsparency Bootball <kernel>'").String()
	createKernel      = create.Flag("kernel", "Operation system kernel").Required().ExistingFile()
	createInitramfs   = create.Flag("initramfs", "Operation system initramfs").ExistingFile()
	createCmdline     = create.Flag("cmd", "Kernel command line").String()
	createTboot       = create.Flag("tboot", "Pre-execution module that sets up TXT").ExistingFile()
	createTbootArgs   = create.Flag("tcmd", "tboot command line").String()
	createRootCert    = create.Flag("cert", "Root certificate of certificates used for signing").Required().ExistingFile()
	createACM         = create.Flag("acm", "Authenticated Code Module for TXT. This can be a path to single ACM or directory containig multiple ACMs.").ExistingFileOrDir()
	createAllowNonTXT = create.Flag("unsave", "Allow booting without TXT").Bool()
	createHWAddr      = create.Flag("mac", "Hardware address of the host if the created bootball needs to be individual for a specific host.").String()

	sign            = kingpin.Command("sign", "Sign the binary inside the provided bootball")
	signPrivKeyFile = sign.Flag("key", "Private key for signing").Required().ExistingFile()
	signCertFile    = sign.Flag("cert", "Certificate corresponding to the private key").Required().ExistingFile()
	signBootball    = sign.Arg("bootball", "Archive created by 'stconfig create'").Required().ExistingFile()

	unpack         = kingpin.Command("unpack", "Unpack boot ball  file into directory")
	unpackBootball = unpack.Arg("bootball", "Archive containing the boot files").Required().ExistingFile()
)

func main() {
	log.SetPrefix("stmanager: ")
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(goversion).Author(Author)
	kingpin.CommandLine.Help = HelpText

	switch kingpin.Parse() {
	case create.FullCommand():
		var label string
		if *createLabel != "" {
			label = *createLabel
		} else {
			k := filepath.Base(*createKernel)
			label = fmt.Sprintf("System Tarnsparency Bootball %s", k)
		}
		var acms []string
		if *createACM != "" {
			stat, err := os.Stat(*createACM)
			if err != nil {
				log.Fatal(err)
			}
			if stat.IsDir() {
				err := filepath.Walk(*createACM, func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						log.Fatalf("%s must contain acm files only", *createACM)
					}
					acms = append(acms, path)
					return nil
				})
				if err != nil {
					panic(err)
				}
			} else {
				acms = append(acms, *createACM)
			}
		}
		if err := packBootBall(*createOut, label, *createKernel, *createInitramfs, *createCmdline, *createTboot, *createTbootArgs, *createRootCert, acms, *createAllowNonTXT, *createHWAddr); err != nil {
			log.Fatalln(err.Error())
		}
	case sign.FullCommand():
		if err := addSignatureToBootBall(*signBootball, *signPrivKeyFile, *signCertFile); err != nil {
			log.Fatalln(err.Error())
		}
	case unpack.FullCommand():
		if err := unpackBootBall(*unpackBootball); err != nil {
			log.Fatalln(err.Error())
		}
	default:
		log.Fatal("Command not found")
	}
}
