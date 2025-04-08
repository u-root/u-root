// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checker

import "fmt"

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

// Run runs the checks and remediations from a check list, in order, and prints the
// check and remediation status.
func Run(checklist []Check) error {
	for idx, check := range checklist {
		fmt.Printf(green("#%d", idx+1)+" Running check '%s'.. ", check.Name)
		if checkErr := check.Run(); checkErr != nil {
			fmt.Println(red("failed: %v", checkErr))
			if check.Remediate != nil {
				fmt.Println(yellow("  -> running remediation"))
				if remErr := check.Remediate(); remErr != nil {
					fmt.Print(red("     Remediation for '%s' failed: %v\n", check.Name, remErr))
					if check.StopOnError {
						fmt.Println("Exiting")
						return remErr
					}
				} else {
					fmt.Printf("     Remediation for '%s' succeeded\n", check.Name)
				}
			} else {
				msg := fmt.Sprintf(" -> no remediation found for %s", check.Name)
				if check.StopOnError {
					fmt.Println(yellow("%s , stop on error requested. Exiting.", msg))
					return checkErr
				}
				fmt.Println(yellow("%s, skipping.", msg))
			}
		} else {
			fmt.Println(green("OK"))
		}
	}
	return nil
}
