// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sos

import (
	"fmt"
	"sync"
)

type Registry map[string]uint

type SosService struct {
	rWLock   sync.RWMutex
	registry Registry
}

func (s *SosService) Read(serviceName string) (uint, error) {
	s.rWLock.RLock()
	defer s.rWLock.RUnlock()
	port, exists := s.registry[serviceName]
	if !exists {
		return 0, fmt.Errorf("%v is not in the registry", serviceName)
	}
	return port, nil
}

func (s *SosService) Register(serviceName string, portNum uint) error {
	s.rWLock.Lock()
	defer s.rWLock.Unlock()
	_, exists := s.registry[serviceName]
	if exists {
		return fmt.Errorf("%v already exists", serviceName)
	}
	s.registry[serviceName] = portNum
	return nil
}

func (s *SosService) Unregister(serviceName string) {
	s.rWLock.Lock()
	defer s.rWLock.Unlock()
	delete(s.registry, serviceName)
}

func (s *SosService) SnapshotRegistry() Registry {
	s.rWLock.RLock()
	defer s.rWLock.RUnlock()
	snapshot := make(map[string]uint)
	for name, port := range s.registry {
		snapshot[name] = port
	}
	return snapshot
}

func NewSosService() *SosService {
	return &SosService{
		registry: make(map[string]uint),
	}
}
