// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
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

func setupStubServer() (*WifiServer, error) {
	service, err := setupStubService()
	if err != nil {
		return nil, err
	}
	service.Start()
	return NewWifiServer(service), nil
}

func TestConnectHandlerSuccess(t *testing.T) {
	// Set Up
	server, err := setupStubServer()
	if err != nil {
		t.Fatal(err)
	}
	defer server.service.Shutdown()
	router := server.buildRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()
	m := ConnectJsonMsg{EssidStub, PassStub, IdStub}
	b, err := json.Marshal(m)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	req, err := http.NewRequest("POST", ts.URL+"/connect", bytes.NewBuffer(b))
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	// Execute
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	// Assert
	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	var retMsg struct{ Error string }
	if err := decoder.Decode(&retMsg); err != nil {
		t.Errorf("Error Decode JSON Response")
		return
	}
	// nil in response
	if retMsg != struct{ Error string }{} {
		t.Errorf("\ngot:%v\nwant:%v", retMsg, struct{ Error string }{})
		return
	}
	// Check for State change
	state := server.service.GetState()
	if state.CurEssid != EssidStub {
		t.Errorf("\ngot:%v\nwant:%v", state.CurEssid, EssidStub)
	}
}

func TestConnectHandlerFail(t *testing.T) {
	// Set Up
	server, err := setupStubServer()
	if err != nil {
		t.Fatal(err)
	}
	defer server.service.Shutdown()
	router := server.buildRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()
	m := ConnectJsonMsg{EssidStub, "", IdStub}
	b, err := json.Marshal(m)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	req, err := http.NewRequest("POST", ts.URL+"/connect", bytes.NewBuffer(b))
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	// Execute
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	// Assert
	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	var retMsg struct{ Error string }
	if err := decoder.Decode(&retMsg); err != nil {
		t.Errorf("Error Decode JSON Response")
		return
	}
	// Error message in response
	if retMsg != struct{ Error string }{"Invalid user input"} {
		t.Errorf("\ngot:%v\nwant:%v", retMsg, struct{ Error string }{})
		return
	}
}

func TestRefreshHandler(t *testing.T) {
	// Set Up
	server, err := setupStubServer()
	if err != nil {
		t.Fatal(err)
	}
	defer server.service.Shutdown()
	router := server.buildRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()
	req, err := http.NewRequest("POST", ts.URL+"/refresh", nil)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	// Execute
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	// Assert
	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	var retMsg struct{ Error string }
	if err := decoder.Decode(&retMsg); err != nil {
		t.Errorf("Error Decode JSON Response")
		return
	}
	// nil in response
	if retMsg != struct{ Error string }{} {
		t.Errorf("\ngot:%v\nwant:%v", retMsg, struct{ Error string }{})
		return
	}
}

func TestHandlersRace(t *testing.T) {
	// Set Up
	numConnectRoutines, numRefreshGoRoutines, numReadGoRoutines := 10, 10, 100
	server, err := setupStubServer()
	if err != nil {
		t.Fatal(err)
	}
	defer server.service.Shutdown()
	router := server.buildRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()

	essidChoices := []string{"stub1", "stub2", "stub3"}

	// Execute
	var wg sync.WaitGroup

	for i := 0; i < numConnectRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			idx := rand.Intn(len(essidChoices))
			m := ConnectJsonMsg{essidChoices[idx], "", ""}
			b, err := json.Marshal(m)
			if err != nil {
				t.Errorf("error: %v", err)
				return
			}
			req, err := http.NewRequest("POST", ts.URL+"/connect", bytes.NewBuffer(b))
			if err != nil {
				t.Errorf("error: %v", err)
				return
			}
			if _, err = http.DefaultClient.Do(req); err != nil {
				t.Errorf("error: %v", err)
				return
			}
		}()
	}

	for i := 0; i < numRefreshGoRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, err := http.NewRequest("POST", ts.URL+"/refresh", nil)
			if err != nil {
				t.Errorf("error: %v", err)
				return
			}
			if _, err = http.DefaultClient.Do(req); err != nil {
				t.Errorf("error: %v", err)
				return
			}
		}()
	}

	for i := 0; i < numReadGoRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, err := http.NewRequest("GET", ts.URL+"/", nil)
			if err != nil {
				t.Errorf("error: %v", err)
				return
			}
			if _, err = http.DefaultClient.Do(req); err != nil {
				t.Errorf("error: %v", err)
				return
			}
		}()
	}
	wg.Wait()
}
