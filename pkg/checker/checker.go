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
func Run(checklist []Check) {
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
						return
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
	return
}
