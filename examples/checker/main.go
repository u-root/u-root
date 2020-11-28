// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Checker is a troubleshooting and auto-remediation tool to help you identify
// issues that won't let your system boot.
// It will simply run through a checklist of checks and optional remediations,
// and continue until it either succeeds or fails without recovery.
//
// The checklist is specified via configuration files in JSON format. Run the
// `genconfig` subcommand to get an example configuration file.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/u-root/u-root/pkg/checker"
)

var (
	flagsetRun    = flag.NewFlagSet("run", flag.ExitOnError)
	flagRunConfig = flagsetRun.String("c", "", "Configuration file")
	flagRunParams = flagsetRun.String("p", "", "Template parameters, comma-separated key=value, to use text/template in config files. Example: interface=eth0,linkspeed=40G")

	flagsetGenconfig = flag.NewFlagSet("genconfig", flag.ExitOnError)
)

// CurrentVersion specifies the current version for the configuration format.
const CurrentVersion = "0.1"

var sampleConfig = fmt.Sprintf(`{
    "version": "%s",
	"checks": [
		{
			"check": {
				"name": "printf",
				"description": "example printf",
				"args": ["this is an %%s", "example"]
			},
			"remediation": {
				"name": "cmd",
				"description": "just return true",
				"args": ["/bin/true"]
			}
		},
		{
			"check": {
				"name": "interface_exists",
				"description": "check that interface lo exists",
				"args": ["lo"]
			}
		}
	]
}`, CurrentVersion)

var (
	checkRunners = map[string]interface{}{
		"interface_exists":             checker.InterfaceExists,
		"printf":                       printfWrapper,
		"cmd":                          checker.RunCmd,
		"link_speed":                   checker.LinkSpeed,
		"link_autoneg":                 checker.LinkAutoneg,
		"interface_has_linklocal_addr": checker.InterfaceHasLinkLocalAddress,
		"interface_has_global_addrs":   checker.InterfaceHasGlobalAddresses,
		"interface_can_do_dhcpv6":      checker.InterfaceCanDoDHCPv6,
	}
	remediationRunners = map[string]interface{}{
		"printf":                     printfWrapper,
		"cmd":                        checker.RunCmd,
		"remediate_interface_exists": checker.InterfaceRemediate,
	}
)

func parseTemplateParams(paramsString string) (map[string]string, error) {
	if paramsString == "" {
		return nil, nil
	}
	params := make(map[string]string, 0)
	parts := strings.Split(paramsString, ",")
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		switch len(kv) {
		case 0:
			// empty param
			continue
		case 2:
			params[kv[0]] = kv[1]
		default:
			return nil, fmt.Errorf("invalid parameter '%s', want key=value, got no value", kv)
		}
	}
	return params, nil
}

func runConfig(configFile, paramsString string) error {
	if configFile == "" {
		return errors.New("No configuration file specified")
	}
	params, err := parseTemplateParams(paramsString)
	if err != nil {
		return fmt.Errorf("invalid template parameters: %v", err)
	}
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %v", err)
	}
	tmpl, err := template.New("config").Parse(string(buf))
	if err != nil {
		return fmt.Errorf("failed to parse config template: %v", err)
	}
	var bbuf bytes.Buffer
	if err := tmpl.Execute(&bbuf, params); err != nil {
		return fmt.Errorf("failed to execute config template: %v", err)
	}
	log.Printf("Configuration after parameter substitution: %s", bbuf.Bytes())

	log.Printf("Registered checks:")
	idx := 1
	for cn := range checkRunners {
		log.Printf("% 3d) %s", idx, cn)
		idx++
	}
	log.Printf("Registered remediations:")
	idx = 1
	for cn := range remediationRunners {
		log.Printf("% 3d) %s", idx, cn)
		idx++
	}

	config, err := checker.NewConfig(bbuf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to parse configuration file: %v", err)
	}
	if config.Version != CurrentVersion {
		return fmt.Errorf("cannot handle configuration version '%s', want '%s'", config.Version, CurrentVersion)
	}
	log.Printf("Found %d checks", len(config.Checklist))
	for idx, item := range config.Checklist {
		log.Printf("% 3d) check       : %s (%s) with args %v", idx, item.Check.Name, item.Check.Description, item.Check.Args)
		if item.Remediation != nil {
			log.Printf("     remediation : %s (%s) with args %v", item.Remediation.Name, item.Remediation.Description, item.Remediation.Args)
		} else {
			log.Printf("     remediation : none")
		}
	}

	checklist, err := checker.ResolveChecklist(config, checkRunners, remediationRunners)
	if err != nil {
		return fmt.Errorf("failed to resolve checklist: %v", err)
	}
	return checker.Run(checklist)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  %s <command> [flags]:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Commands:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "run     Run a checklist\n")
		flagsetRun.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "genconfig\n")
		fmt.Fprintf(flag.CommandLine.Output(), "        Print an example configuration\n")
		flagsetGenconfig.PrintDefaults()
	}
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Missing command, see -h")
		os.Exit(1)
	}

	subcmd := flag.Arg(0)
	switch subcmd {
	case "run":
		flagsetRun.Parse(os.Args[2:])
		if err := runConfig(*flagRunConfig, *flagRunParams); err != nil {
			log.Fatal(err)
		}
	case "genconfig":
		flagsetGenconfig.Parse(os.Args[2:])
		fmt.Println(sampleConfig)
		os.Exit(0)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}
