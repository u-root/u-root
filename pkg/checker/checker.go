// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checker

import "fmt"

// CheckRunner is the type of checking functions
type CheckRunner func() error

// RemediationRunner is the type of remediation functions
type RemediationRunner func() error

// Item is a checklist item.
type Item struct {
	Check           Check
	Remediation     *Remediation
	ContinueOnError bool
}

// Check is a check object.
type Check struct {
	Name        string
	Description string
	Run         CheckRunner
}

// Remediation is a remediation object.
type Remediation struct {
	Name        string
	Description string
	Run         RemediationRunner
}

// Run runs the checks and remediations from a check list, in order, and prints the
// check and remediation status.
func Run(checklist []Item) error {
	for idx, item := range checklist {
		fmt.Printf(green("#%d", idx+1)+" Running check '%s'.. ", item.Check.Name)
		if checkErr := item.Check.Run(); checkErr != nil {
			fmt.Println(red("failed: %v", checkErr))
			if item.Remediation != nil {
				fmt.Println(yellow("  -> running remediation"))
				if remErr := item.Remediation.Run(); remErr != nil {
					fmt.Print(red("     Remediation '%s' failed: %v\n", item.Remediation.Name, remErr))
					if !item.ContinueOnError {
						fmt.Println("Exiting because ContinueOnError is false.")
						return remErr
					}
				} else {
					fmt.Printf("     Remediation '%s' succeeded\n", item.Remediation.Name)
				}
			} else {
				fmt.Printf(yellow(" -> no remediation specified for %s\n", item.Check.Name))
				if !item.ContinueOnError {
					fmt.Println("Exiting because ContinueOnError is false.")
					return checkErr
				}
			}
		} else {
			fmt.Println(green("OK"))
		}
	}
	return nil
}

// ResolveChecklist resolves a configuration object into an `[]Item`,
// after looking up the check and remediation functions by name.
// The function mappings are passed via `checkRunners` and `remediationRunners`.
func ResolveChecklist(config *Config, checkRunners map[string]interface{}, remediationRunners map[string]interface{}) ([]Item, error) {
	checklist := make([]Item, 0, len(config.Checklist))
	for _, item := range config.Checklist {
		newItem := Item{
			ContinueOnError: item.ContinueOnError,
		}
		// look up check runner
		cr, ok := checkRunners[item.Check.Name]
		if !ok {
			return nil, fmt.Errorf("check runner '%s' not found", item.Check.Name)
		}
		crf, err := makeCheckRunner(cr, item.Check.Args...)
		if err != nil {
			return nil, fmt.Errorf("failed to get wrapper for '%s': %v", item.Check.Name, err)
		}
		newItem.Check = Check{
			Name:        item.Check.Name,
			Description: item.Check.Description,
			Run:         crf,
		}

		// look up the optional remediation runner
		if item.Remediation != nil {
			rr, ok := remediationRunners[item.Remediation.Name]
			if !ok {
				return nil, fmt.Errorf("remediation runner '%s' not found", item.Remediation.Name)
			}
			rrf, err := makeRemediationRunner(rr, item.Remediation.Args...)
			if err != nil {
				return nil, fmt.Errorf("failed to get wrapper for '%s': %v", item.Remediation.Name, err)
			}
			newItem.Remediation = &Remediation{
				Name:        item.Remediation.Name,
				Description: item.Remediation.Description,
				Run:         rrf,
			}
		}

		checklist = append(checklist, newItem)
	}

	return checklist, nil
}
