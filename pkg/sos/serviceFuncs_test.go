// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sos

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

func TestRegistersNecessaryPatterns(t *testing.T) {
	router := mux.NewRouter()
	RegistersNecessaryPatterns(router)
	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Status Code\ngot:%v\nwant:%v", res.StatusCode, http.StatusOK)
		return
	}

	msg, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	if string(msg) != "pong" {
		t.Errorf("Body\ngot:%v\nwant:%v", string(msg), "pong")
	}
}

func TestRegisterServiceWithSosSuccess(t *testing.T) {
	// Set up
	service := NewSosService()
	server := SosServer{service}
	r := server.buildRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	if err := registerServiceWithSos(knownServ1.service, knownServ1.port, ts.URL); err != nil {
		t.Errorf("error: %v", err)
		return
	}

	if service.registry[knownServ1.service] != knownServ1.port {
		t.Errorf("In Registry\ngot:%v\nwant:%v", service.registry[knownServ1.service], knownServ1.port)
	}
}

func TestRegisterServiceWithSosFail(t *testing.T) {
	// Set up
	service := setUpKnownServices()
	server := SosServer{service}
	r := server.buildRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	err := registerServiceWithSos(knownServ1.service, knownServ1.port, ts.URL)
	if !reflect.DeepEqual(err, fmt.Errorf("%v already exists", knownServ1.service)) {
		t.Errorf("\ngot:%v\nwant:%v", err, fmt.Errorf("%v already exists", knownServ1.service))
	}
}

func TestUnregisterServiceWithSosSuccess(t *testing.T) {
	// Set up
	service := setUpKnownServices()
	server := SosServer{service}
	r := server.buildRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	if err := unregisterServiceWithSos(knownServ1.service, ts.URL); err != nil {
		t.Errorf("error: %v", err)
		return
	}

	if _, err := service.Read(knownServ1.service); !reflect.DeepEqual(err, fmt.Errorf("%v is not in the registry", knownServ1.service)) {
		t.Errorf("unregister(%v)\ngot:(%v)\nwant:(%v)", knownServ1.service, err, fmt.Errorf("%v is not in the registry", knownServ1.service))
	}

}
