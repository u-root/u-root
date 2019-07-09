// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func ThisCheckAlwaysSucceeds(args CheckArgs) error {
	return nil
}

func ThisCheckAlwaysFails(args CheckArgs) error {
	return fmt.Errorf("TEST FAILURE")
}

func init() {
	registerCheckFun(ThisCheckAlwaysSucceeds)
	registerCheckFun(ThisCheckAlwaysFails)
}

func TestRunSimpleOK(t *testing.T) {
	check := Check{
		Description:     "a_description",
		CheckerFuncName: "ThisCheckAlwaysSucceeds",
		CheckerFuncArgs: nil,
	}
	result := check.Run()
	require.Equal(t, check.Description, result.Description)
	require.Equal(t, check.CheckerFuncName, result.CheckerFuncName)
	require.Equal(t, check.CheckerFuncArgs, result.CheckerFuncArgs)
	require.Equal(t, result.Result, ResultOK)
	require.Equal(t, result.Error, "")
	require.Equal(t, len(result.RemediationResults), 0)
	require.Equal(t, result.StoppedOnFailure, false)
}

func TestRunSimpleError(t *testing.T) {
	check := Check{
		Description:     "a_description",
		CheckerFuncName: "ThisCheckAlwaysFails",
		CheckerFuncArgs: nil,
	}
	result := check.Run()
	require.Equal(t, check.Description, result.Description)
	require.Equal(t, check.CheckerFuncName, result.CheckerFuncName)
	require.Equal(t, check.CheckerFuncArgs, result.CheckerFuncArgs)
	require.Equal(t, ResultError, result.Result)
	require.Equal(t, result.Error, "TEST FAILURE")
	require.Equal(t, len(result.RemediationResults), 0)
	require.Equal(t, result.StoppedOnFailure, false)
}

func TestRunRemediation(t *testing.T) {
	check := Check{
		Description:     "a_description",
		CheckerFuncName: "ThisCheckAlwaysFails",
		CheckerFuncArgs: nil,
		Remediations: []Check{
			Check{
				CheckerFuncName: "ThisCheckAlwaysSucceeds",
			},
		},
	}
	result := check.Run()
	require.Equal(t, check.Description, result.Description)
	require.Equal(t, check.CheckerFuncName, result.CheckerFuncName)
	require.Equal(t, check.CheckerFuncArgs, result.CheckerFuncArgs)
	require.Equal(t, ResultError, result.Result)
	require.Equal(t, result.Error, "TEST FAILURE")
	require.Equal(t, len(result.RemediationResults), 1)
	require.Equal(t, result.RemediationResults[0].CheckerFuncName, "ThisCheckAlwaysSucceeds")
	require.Equal(t, result.RemediationResults[0].Error, "")
	require.Equal(t, result.RemediationResults[0].Result, ResultOK)
	require.Equal(t, result.StoppedOnFailure, false)
}

func TestRunStopOnFailure(t *testing.T) {
	check := Check{
		Description:     "a_description",
		CheckerFuncName: "ThisCheckAlwaysFails",
		CheckerFuncArgs: nil,
		StopOnFailure:   true,
	}
	result := check.Run()
	require.Equal(t, len(result.RemediationResults), 0)
	require.Equal(t, true, result.StoppedOnFailure)
}

func TestRunStopOnFailureWithRemediations(t *testing.T) {
	check := Check{
		Description:     "a_description",
		CheckerFuncName: "ThisCheckAlwaysFails",
		CheckerFuncArgs: nil,
		Remediations: []Check{
			Check{
				CheckerFuncName: "ThisCheckAlwaysSucceeds",
			},
		},
		StopOnFailure: true,
	}
	result := check.Run()
	require.Equal(t, len(result.RemediationResults), 0)
	require.Equal(t, true, result.StoppedOnFailure)
}

func TestRunChecklist(t *testing.T) {
	checklist := []Check{
		Check{
			Description:     "a_description",
			CheckerFuncName: "ThisCheckAlwaysSucceeds",
			CheckerFuncArgs: nil,
		},
	}

	results, numErrors := Run(checklist)
	require.Equal(t, numErrors, 0)
	require.Equal(t, len(results), 1)
	require.Equal(t, checklist[0].Description, results[0].Description)
	require.Equal(t, checklist[0].CheckerFuncName, results[0].CheckerFuncName)
	require.Equal(t, checklist[0].CheckerFuncArgs, results[0].CheckerFuncArgs)
	require.Equal(t, results[0].Result, ResultOK)
}

func TestRunChecklistError(t *testing.T) {
	checklist := []Check{
		Check{
			Description:     "a_description",
			CheckerFuncName: "ThisCheckAlwaysFails",
			CheckerFuncArgs: nil,
		},
	}

	results, numErrors := Run(checklist)
	require.Equal(t, numErrors, 1)
	require.Equal(t, len(results), 1)
	require.Equal(t, checklist[0].Description, results[0].Description)
	require.Equal(t, checklist[0].CheckerFuncName, results[0].CheckerFuncName)
	require.Equal(t, checklist[0].CheckerFuncArgs, results[0].CheckerFuncArgs)
	require.Equal(t, results[0].Result, ResultError)
}
