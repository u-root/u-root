// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/9elements/converged-security-suite/pkg/hwapi"
	"github.com/9elements/converged-security-suite/pkg/test"
)

func runTxtTests(verbose bool) bool {

	hwAPI := hwapi.GetAPI()

	success, failureMsg, err := test.RunTestsSilent(hwAPI, test.TestsTXTReady)
	if err != nil {
		log.Printf("Error checking for TXT Ready support: %v\n", err)
		return false
	}
	if !success {
		log.Printf("TXT not availabe as: %s\n", failureMsg)
		return false
	}
	return true

}
