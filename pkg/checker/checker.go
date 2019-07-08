package checker

import "fmt"

// ResultOK is the string to appear in CheckResult.Result if the check was successful
const ResultOK = "OK"

//ResultError is the string to appear in CheckResult.Result if the check failed
const ResultError = "ERROR"

// CheckFun is the type of checking functions
type CheckFun func(CheckArgs) error

// CheckArgs is the argument map to be passed to CheckFuns
type CheckArgs map[string]string

// A Check describes a call to a CheckFun with possible remediations
type Check struct {
	Description string `json:"description"`

	CheckerFuncName string    `json:"checkFun"`
	CheckerFuncArgs CheckArgs `json:"checkArgs"`

	Remediations  []Check `json:"remediations"`
	StopOnFailure bool    `json:"stopOnFailure"`
}

// A CheckResult describes the result of Running a Check
type CheckResult struct {
	Description string `json:"description"`

	CheckerFuncName string    `json:"checkFun"`
	CheckerFuncArgs CheckArgs `json:"checkArgs"`

	Result             string        `json:"result"`
	Error              string        `json:"error"`
	RemediationResults []CheckResult `json:"remediationResults,omitempty"`
	StoppedOnFailure   bool          `json:"-"`
}

// Run a Check
func (check *Check) Run() CheckResult {
	return check.run(0)
}

func (check *Check) run(lvl int) CheckResult {
	result := CheckResult{
		Description:     check.Description,
		CheckerFuncName: check.CheckerFuncName,
		CheckerFuncArgs: check.CheckerFuncArgs,
		Result:          ResultOK,
	}

	fmt.Printf(indent(lvl)+"Running check '%s' (%s(%#v)).. ", check.Description, check.CheckerFuncName, check.CheckerFuncArgs)

	// Call check function and get (possible) error
	checkErr := Call(check.CheckerFuncName, check.CheckerFuncArgs)

	if checkErr == nil {
		fmt.Print(green("OK\n"))
	} else {
		result.Result = ResultError
		result.Error = checkErr.Error()
		result.StoppedOnFailure = check.StopOnFailure
		fmt.Print(red("failed: %s\n", result.Error))
	}

	// If the check failed (and StopOnFailure is false), run OnFailure callbacks
	if result.StoppedOnFailure {
		fmt.Printf(indent(lvl)+"Check '%s' failed with stopOnFailure=true, bailing out...\n", check.Description)
		return result
	}
	if len(check.Remediations) > 0 {
		fmt.Printf(indent(lvl) + yellow("Running remediations for '%s':\n", check.Description))
		for _, c := range check.Remediations {
			res := c.run(lvl + 1)
			result.RemediationResults = append(result.RemediationResults, res)

			if res.StoppedOnFailure {
				result.StoppedOnFailure = true
				break
			}
		}
	}

	return result
}

// Run a list of Checks
func Run(checklist []Check) ([]CheckResult, int) {
	results := make([]CheckResult, 0)
	numErrors := 0
	for _, check := range checklist {
		res := check.Run()
		results = append(results, res)

		if res.Error != "" {
			numErrors++
		}

		if res.StoppedOnFailure {
			break
		}
	}
	return results, numErrors
}

func indent(lvl int) string {
	s := ""
	for i := 0; i < lvl; i++ {
		s += "  "
	}
	return s
}
