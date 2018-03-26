// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sos

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

type RegistryEntryStub struct {
	service string
	port    uint
}

var (
	knownServ1 = RegistryEntryStub{"stub1", 1}
	knownServ2 = RegistryEntryStub{"stub2", 2}
	knownServ3 = RegistryEntryStub{"stub3", 3}
	newServ1   = RegistryEntryStub{"stub4", 4}
	newServ2   = RegistryEntryStub{"stub5", 5}
	newServ3   = RegistryEntryStub{"stub6", 6}
)

func setUpKnownServices() *SosService {
	service := NewSosService()
	service.registry[knownServ1.service] = knownServ1.port
	service.registry[knownServ2.service] = knownServ2.port
	service.registry[knownServ3.service] = knownServ3.port
	return service
}

func TestReadNonExist(t *testing.T) {
	s := NewSosService()
	if _, err := s.Read(knownServ1.service); !reflect.DeepEqual(err, fmt.Errorf("%v is not in the registry", knownServ1.service)) {
		t.Errorf("read(%v)\ngot:(%v)\nwant:(%v)", knownServ1.service, err, fmt.Errorf("%v is not in the registry", knownServ1.service))
	}
}

func TestRead(t *testing.T) {
	s := setUpKnownServices()
	if port, err := s.Read(knownServ1.service); err != nil || port != knownServ1.port {
		t.Errorf("read(%v)\ngot:(%v, %v)\nwant:(%v, %v)", knownServ1.service, port, err, knownServ1.port, nil)
	}
}

func TestRegisterAlreadyExists(t *testing.T) {
	s := setUpKnownServices()
	err := s.Register(knownServ1.service, knownServ1.port)
	if !reflect.DeepEqual(err, fmt.Errorf("%v already exists", knownServ1.service)) {
		t.Errorf("Already Exists Register\ngot:(%v)\nwant:(%v)", err, fmt.Errorf("%v already exists", knownServ1.service))
	}
}

func TestRegisterSuccess(t *testing.T) {
	s := NewSosService()
	s.Register(knownServ1.service, knownServ1.port)
	if port, err := s.Read(knownServ1.service); err != nil || port != knownServ1.port {
		t.Errorf("register(%v)\ngot:(%v, %v)\nwant:(%v, %v)", knownServ1, port, err, knownServ1.port, nil)
	}
}

func TestUnregisterNonExist(t *testing.T) {
	s := NewSosService()
	s.Unregister(knownServ1.service)
	// should not panic
}

func TestUnregister(t *testing.T) {
	s := setUpKnownServices()
	s.Unregister(knownServ1.service)
	if _, err := s.Read(knownServ1.service); !reflect.DeepEqual(err, fmt.Errorf("%v is not in the registry", knownServ1.service)) {
		t.Errorf("unregister(%v)\ngot:(%v)\nwant:(%v)", knownServ1.service, err, fmt.Errorf("%v is not in the registry", knownServ1.service))
	}
}

func TestSnapshot(t *testing.T) {
	s := setUpKnownServices().SnapshotRegistry()
	if port, exists := s[knownServ1.service]; !exists || port != knownServ1.port {
		t.Errorf("%v\ngot:(%v, %v)\nwant:(%v, %v)", knownServ1, port, exists, knownServ1.port, true)
	}
	if port, exists := s[knownServ2.service]; !exists || port != knownServ2.port {
		t.Errorf("%v\ngot:(%v, %v)\nwant:(%v, %v)", knownServ2, port, exists, knownServ2.port, true)
	}
	if port, exists := s[knownServ3.service]; !exists || port != knownServ3.port {
		t.Errorf("%v\ngot:(%v, %v)\nwant:(%v, %v)", knownServ3, port, exists, knownServ1.port, true)
	}
}

func TestRaceCondtion(t *testing.T) {
	//Set Up
	s := NewSosService()
	numReadGoRoutines := 10
	numRegisterGoRoutines := 20
	numUnregisterGoRoutines := 20
	numSnapshotGoRoutines := 10

	var wg sync.WaitGroup
	for i := 0; i < numRegisterGoRoutines; i++ {
		wg.Add(1)
		go func(idx uint) {
			defer wg.Done()
			s.Register(fmt.Sprintf("stub%v", idx), idx)
		}(uint(i))
	}

	for i := 0; i < numReadGoRoutines; i++ {
		wg.Add(1)
		go func(idx uint) {
			defer wg.Done()
			s.Read(fmt.Sprintf("stub%v", idx))
		}(uint(i))
	}

	for i := 0; i < numUnregisterGoRoutines; i++ {
		wg.Add(1)
		go func(idx uint) {
			defer wg.Done()
			s.Unregister(fmt.Sprintf("stub%v", idx))
		}(uint(i))
	}

	for i := 0; i < numSnapshotGoRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.SnapshotRegistry()
		}()
	}

	wg.Wait()
}
