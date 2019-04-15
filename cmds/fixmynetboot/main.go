package main

import (
	"fmt"
	"os"
)

// fixmynetboot is a troubleshooting tool that can help you identify issues that
// won't let your system boot over the network.

// Check is a type that implements a netboot check
type Check struct {
	Name      string
	Run       func() error
	Remediate func() error
}

func main() {
	checklist := []Check{
		Check{"wlp2s0 exists", interfaceExists("wlp2s0"), nil},
		Check{"eth0 exists", interfaceExists("eth0"), interfaceRemediate("eth0")},
	}

	for idx, check := range checklist {
		fmt.Printf(green("#%d", idx+1)+" Running check '%s'.. ", check.Name)
		if err := check.Run(); err != nil {
			fmt.Println(red("failed: %v", err))
			if check.Remediate != nil {
				fmt.Println(yellow("  -> running remediation"))
				if err := check.Remediate(); err != nil {
					fmt.Printf(red("     Remediation for '%s' failed: %v\n", check.Name, err))
					fmt.Println("Exiting")
					os.Exit(1)
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
