package main

import (
	"fmt"
	"os"
)

// fixmynetboot is a troubleshooting tool that can help you identify issues that
// won't let your system boot over the network.

// Checker is the type of checking functions
type Checker func() error

// Remediator is the type of remediation functions
type Remediator func() error

// Check is a type that implements a netboot check
type Check struct {
	Name        string
	Run         Checker
	Remediate   Remediator
	StopOnError bool
}

func main() {
	checklist := []Check{
		Check{"wlp2s0 exists", interfaceExists("wlp2s0"), nil, false},
		Check{"eth0 exists", interfaceExists("eth0"), interfaceRemediate("eth0"), false},
		Check{"eth0 has link-local", interfaceHasLinkLocalAddress("wlp2s0"), nil, false},
		Check{"eth0 has global addresses", interfaceHasGlobalAddresses("wlp2s0"), nil, false},
	}

	for idx, check := range checklist {
		fmt.Printf(green("#%d", idx+1)+" Running check '%s'.. ", check.Name)
		if err := check.Run(); err != nil {
			fmt.Println(red("failed: %v", err))
			if check.Remediate != nil {
				fmt.Println(yellow("  -> running remediation"))
				if err := check.Remediate(); err != nil {
					fmt.Printf(red("     Remediation for '%s' failed: %v\n", check.Name, err))
					if check.StopOnError {
						fmt.Println("Exiting")
						os.Exit(1)
					}
				} else {
					fmt.Printf("     Remediation for '%s' succeeded\n", check.Name)
				}
			} else {
				fmt.Printf(yellow(" -> no remediation found for '%s', skipping\n", check.Name))
			}
		} else {
			fmt.Println(green("OK"))
		}
	}
}
