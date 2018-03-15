// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"reflect"
	"sync"
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
	EssidStub = "stub"
	IdStub    = "stub"
	PassStub  = "123456789"

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

func turnOnTestingMode() {
	t := true
	test = &t
}

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

func TestConnectHandle(t *testing.T) {
	// Set Up
	turnOnTestingMode()
	connectWifiArbitratorSetup("", "", 2)
	defer close(ConnectReqChan)

	m := ConnectJsonMsg{EssidStub, PassStub, IdStub}
	b, err := json.Marshal(m)
	if err != nil {
		t.Errorf("Setup Fails")
		return
	}

	r := httptest.NewRequest("GET", "localhost:"+PortNum+"/connect", bytes.NewBuffer(b))
	w := httptest.NewRecorder()
	connectHandle(w, r)
	if CurEssid != EssidStub {
		t.Errorf("\ngot:%v\nwant:%v", CurEssid, EssidStub)
	}
}

func TestConnectHandleOneAfterAnother(t *testing.T) {
	// Set Up
	turnOnTestingMode()
	connectWifiArbitratorSetup("", "", 2)
	defer close(ConnectReqChan)

	m1 := ConnectJsonMsg{"stub1", PassStub, IdStub}
	b1, err := json.Marshal(m1)
	if err != nil {
		t.Errorf("Setup Fails")
		return
	}

	m2 := ConnectJsonMsg{"stub2", PassStub, IdStub}
	b2, err := json.Marshal(m2)
	if err != nil {
		t.Errorf("Setup Fails")
		return
	}

	r1 := httptest.NewRequest("GET", "localhost:"+PortNum+"/connect", bytes.NewBuffer(b1))
	w1 := httptest.NewRecorder()
	connectHandle(w1, r1)

	r2 := httptest.NewRequest("GET", "localhost:"+PortNum+"/connect", bytes.NewBuffer(b2))
	w2 := httptest.NewRecorder()
	connectHandle(w2, r2)

	if CurEssid != "stub2" {
		t.Errorf("\ngot:%v\nwant:%v", CurEssid, "stub2")
	}
}

func TestConnectHandleRace(t *testing.T) {
	// Set Up
	turnOnTestingMode()
	numGoRoutines := 100
	connectWifiArbitratorSetup("", "", numGoRoutines)
	defer close(ConnectReqChan)

	var wg sync.WaitGroup
	for i := 0; i < numGoRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m := ConnectJsonMsg{EssidStub, PassStub, IdStub}
			b, err := json.Marshal(m)
			if err != nil {
				t.Errorf("Setup Fails")
				return
			}
			r := httptest.NewRequest("GET", "localhost:"+PortNum+"/connect", bytes.NewBuffer(b))
			w := httptest.NewRecorder()
			connectHandle(w, r)
		}()
	}
	wg.Wait()
}
