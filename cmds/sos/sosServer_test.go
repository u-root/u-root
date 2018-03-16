// Copyright 2018 the u-root Authors. All rights reserved
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

func TestRegisterHandle(t *testing.T) {
	// Set up
	cleanUpForNewTest()
	m := RegisterReqJson{knownServ1.service, knownServ1.port}
	b, err := json.Marshal(m)
	if err != nil {
		t.Error("Setup Fails")
		return
	}
	r := httptest.NewRequest("POST", "localhost:1/register", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	// Execute
	registerHandle(w, r)

	// Assert
	if Registry[knownServ1.service] != knownServ1.port {
		t.Errorf("got:(%v)\nwant:(%v)", Registry[knownServ1.service], knownServ1.port)
	}
}

func TestUnregisterHandle(t *testing.T) {
	// Set up
	cleanUpForNewTest()
	Registry[knownServ1.service] = knownServ1.port
	m := UnRegisterReqJson{knownServ1.service}
	b, err := json.Marshal(m)
	if err != nil {
		t.Error("Setup Fails")
		return
	}
	r := httptest.NewRequest("POST", "localhost:1/unregister", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	// Execute
	unregisterHandle(w, r)

	// Assert
	if _, err := read(knownServ1.service); !reflect.DeepEqual(err, fmt.Errorf("%v is not in the registry", knownServ1.service)) {
		t.Errorf("\ngot:(%v)\nwant:(%v)", err, fmt.Errorf("%v is not in the registry", knownServ1.service))
	}
}

func TestGetService(t *testing.T) {
	// Set up
	cleanUpForNewTest()
	Registry[knownServ1.service] = knownServ1.port
	r := buildRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()
	req, err := http.NewRequest("GET", ts.URL+"/service/"+knownServ1.service, nil)
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
	var retMsg GetServiceResJson
	if err := decoder.Decode(&retMsg); err != nil {
		t.Errorf("Error Decode JSON Response")
		return
	}
	if retMsg.Port != knownServ1.port {
		t.Errorf("\ngot:(%v)\nwant:(%v)", retMsg.Port, knownServ1.port)
	}
}

func TestGetServiceFails(t *testing.T) {
	// Set up
	cleanUpForNewTest()
	r := buildRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()
	req, err := http.NewRequest("GET", ts.URL+"/service/"+knownServ1.service, nil)
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
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("\ngot:(%v)\nwant:(%v)", res.StatusCode, http.StatusNotFound)
	}
}

func TestRace(t *testing.T) {
	// Set Up
	cleanUpForNewTest()
	setUpKnownServices()

	numRegisterGoRoutines, numUnregisterGoRoutines, numReadGoRoutines := 10, 10, 100
	serviceChoices := []RegistryEntryStub{
		knownServ1, knownServ2, knownServ3,
		newServ1, newServ2, newServ3,
	}

	r := buildRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Execute
	var wg sync.WaitGroup

	for i := 0; i < numRegisterGoRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			idx := rand.Intn(len(serviceChoices))
			m := RegisterReqJson{serviceChoices[idx].service, serviceChoices[idx].port}
			b, err := json.Marshal(m)
			if err != nil {
				t.Errorf("Setup Fails")
				return
			}
			req, err := http.NewRequest("POST", ts.URL+"/register", bytes.NewBuffer(b))
			if err != nil {
				t.Errorf("error: %v", err)
				return
			}
			http.DefaultClient.Do(req)
		}()
	}

	for i := 0; i < numUnregisterGoRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			idx := rand.Intn(len(serviceChoices))
			m := UnRegisterReqJson{serviceChoices[idx].service}
			b, err := json.Marshal(m)
			if err != nil {
				t.Errorf("Setup Fails")
				return
			}
			req, err := http.NewRequest("POST", ts.URL+"/unregister", bytes.NewBuffer(b))
			if err != nil {
				t.Errorf("error: %v", err)
				return
			}
			http.DefaultClient.Do(req)
		}()
	}

	for i := 0; i < numReadGoRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			idx := rand.Intn(len(serviceChoices))
			req, err := http.NewRequest("GET", ts.URL+"/service/"+serviceChoices[idx].service, nil)
			if err != nil {
				t.Errorf("error: %v", err)
				return
			}
			http.DefaultClient.Do(req)
		}()
	}

	wg.Wait()
}
