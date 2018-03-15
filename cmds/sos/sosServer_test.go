// Copyright 2018 the u-root Authors. All rights reserved
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

func TestRegisterHandle(t *testing.T) {
	cleanUpForNewTest()
	m := RegisterReqJson{knownServ1.service, knownServ1.port}
	b, err := json.Marshal(m)
	if err != nil {
		t.Errorf("Setup Fails")
		return
	}
	r := httptest.NewRequest("POST", "localhost:1/register", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	registerHandle(w, r)
	if Registry[knownServ1.service] != knownServ1.port {
		t.Errorf("got:(%v)\nwant:(%v)", Registry[knownServ1.service], knownServ1.port)
	}
}

func TestUnregisterHandle(t *testing.T) {
	cleanUpForNewTest()
	Registry[knownServ1.service] = knownServ1.port
	m := UnRegisterReqJson{knownServ1.service}
	b, err := json.Marshal(m)
	if err != nil {
		t.Errorf("Setup Fails")
		return
	}
	r := httptest.NewRequest("POST", "localhost:1/unregister", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	unregisterHandle(w, r)
	if _, err := read(knownServ1.service); !reflect.DeepEqual(err, fmt.Errorf("%v is not in the registry", knownServ1.service)) {
		t.Errorf("\ngot:(%v)\nwant:(%v)", err, fmt.Errorf("%v is not in the registry", knownServ1.service))
	}
}

func TestRace(t *testing.T) {
	cleanUpForNewTest()
	numRegisterGoRoutines := 10
	numUnregisterGoRoutines := 10

	var wg sync.WaitGroup

	for i := 0; i < numRegisterGoRoutines/2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m := RegisterReqJson{knownServ1.service, knownServ1.port}
			b, err := json.Marshal(m)
			if err != nil {
				t.Errorf("Setup Fails")
				return
			}
			r := httptest.NewRequest("GET", "localhost:1/register", bytes.NewBuffer(b))
			w := httptest.NewRecorder()
			registerHandle(w, r)
		}()
	}

	for i := 0; i < numRegisterGoRoutines/2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m := RegisterReqJson{knownServ2.service, knownServ2.port}
			b, err := json.Marshal(m)
			if err != nil {
				t.Errorf("Setup Fails")
				return
			}
			r := httptest.NewRequest("GET", "localhost:1/register", bytes.NewBuffer(b))
			w := httptest.NewRecorder()
			registerHandle(w, r)
		}()
	}

	for i := 0; i < numUnregisterGoRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m := UnRegisterReqJson{knownServ1.service}
			b, err := json.Marshal(m)
			if err != nil {
				t.Errorf("Setup Fails")
				return
			}
			r := httptest.NewRequest("GET", "localhost:1/unregister", bytes.NewBuffer(b))
			w := httptest.NewRecorder()
			unregisterHandle(w, r)
		}()
	}

	wg.Wait()
}
