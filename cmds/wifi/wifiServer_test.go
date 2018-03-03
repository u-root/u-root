// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"reflect"
	"testing"
)

type UserInputValidationTestcase struct {
	name  string
	essid string
	pass  string
	id    string
	exp   []string
	err   error
}

var (
	userInputValidationTestcases = []UserInputValidationTestcase{
		{
			name:  "Essid, passphrase, Id",
			essid: EssidStub,
			pass:  PassStub,
			id:    IdStub,
			exp:   []string{EssidStub, PassStub, IdStub},
			err:   nil,
		},
		{
			name:  "Essid, passphrase",
			essid: EssidStub,
			pass:  PassStub,
			id:    "",
			exp:   []string{EssidStub, PassStub},
			err:   nil,
		},
		{
			name:  "Essid",
			essid: EssidStub,
			pass:  "",
			id:    "",
			exp:   []string{EssidStub},
			err:   nil,
		},
		{
			name:  "No Essid",
			essid: "",
			pass:  PassStub,
			id:    IdStub,
			exp:   nil,
			err:   fmt.Errorf("Invalid user input"),
		},
		{
			name:  "Essid, Id",
			essid: EssidStub,
			pass:  "",
			id:    IdStub,
			exp:   nil,
			err:   fmt.Errorf("Invalid user input"),
		},
	}
)

func TestUserInputValidation(t *testing.T) {
	for _, test := range userInputValidationTestcases {
		out, err := userInputValidation(test.essid, test.pass, test.id)
		if !reflect.DeepEqual(err, test.err) || !reflect.DeepEqual(out, test.exp) {
			t.Logf("TEST %v", test.name)
			fncCall := fmt.Sprintf("userInputValidation(%v, %v, %v)", test.essid, test.pass, test.id)
			t.Errorf("%s\ngot:[%v, %v]\nwant:[%v, %v]", fncCall, out, err, test.exp, test.err)
		}
	}
}
