/*
 * Copyright 2018 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License. You may obtain a copy of
 * the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

// Package simulator provides a go interface to the Microsoft TPM2 simulator.
package simulator

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/google/go-tpm-tools/simulator/internal"
	"github.com/google/go-tpm/tpm2"
)

// Simulator represents a go-tpm compatible interface to the IBM TPM2 simulator.
// Similar to the file-based (for linux) or syscall-based (for Windows) TPM
// handles, no synchronization is provided; the same simulator handle should not
// be used from multiple threads.
type Simulator struct {
	buf bytes.Buffer
}

// ErrSimulatorInUse indicates another open Simulator already exists.
var ErrSimulatorInUse = errors.New("simulator is being used by another caller")

// The simulator is a global resource, so we use the variables below to make
// sure we only ever have one open reference to the Simulator at a time.
var (
	lock  sync.Mutex
	inUse bool
)

// Get the pointer to an initialized, powered on, and started simulator. As only
// one simulator may be running at a time, a second call to Get() will return
// ErrSimulatorInUse until the first Simulator is Closed.
func Get() (*Simulator, error) {
	lock.Lock()
	defer lock.Unlock()
	if inUse {
		return nil, ErrSimulatorInUse
	}

	simulator := &Simulator{}
	internal.Reset(true)
	if err := simulator.on(true); err != nil {
		return nil, err
	}
	inUse = true
	return simulator, nil
}

// GetWithFixedSeedInsecure behaves like Get() expect that all of the internal
// hierarchy seeds are derived from the input seed. Note that this function
// compromises the security of the keys/seeds and should only be used for tests.
func GetWithFixedSeedInsecure(seed int64) (*Simulator, error) {
	s, err := Get()
	if err != nil {
		return nil, err
	}

	internal.SetSeeds(rand.New(rand.NewSource(seed)))
	return s, nil
}

// Reset the TPM as if the host computer had rebooted.
func (s *Simulator) Reset() error {
	if err := s.off(); err != nil {
		return err
	}
	internal.Reset(false)
	return s.on(false)
}

// ManufactureReset behaves like Reset() except that the TPM is complete wiped.
// All data (NVData, Hierarchy seeds, etc...) is cleared or reset.
func (s *Simulator) ManufactureReset() error {
	if err := s.off(); err != nil {
		return err
	}
	internal.Reset(true)
	return s.on(true)
}

// Write executes the command specified by commandBuffer. The command response
// can be retrieved with a subsequent call to Read().
func (s *Simulator) Write(commandBuffer []byte) (int, error) {
	resp, err := internal.RunCommand(commandBuffer)
	if err != nil {
		return 0, err
	}
	return s.buf.Write(resp)
}

// Read gets the response of a command previously issued by calling Write().
func (s *Simulator) Read(responseBuffer []byte) (int, error) {
	return s.buf.Read(responseBuffer)
}

// Close cleans up and stops the simulator, Close() should always be called when
// the Simulator is no longer needed, freeing up other callers to use Get().
func (s *Simulator) Close() error {
	lock.Lock()
	defer lock.Unlock()
	inUse = false
	return s.off()
}

func (s *Simulator) on(manufactureReset bool) error {
	// TPM2_Startup must be the first command the TPM receives.
	if err := tpm2.Startup(s, tpm2.StartupClear); err != nil {
		return fmt.Errorf("startup: %v", err)
	}
	return nil
}

func (s *Simulator) off() error {
	// TPM2_Shutdown must be the last command the TPM receives. We call
	// Shutdown with StartupClear to simulate a full reboot.
	if err := tpm2.Shutdown(s, tpm2.StartupClear); err != nil {
		return fmt.Errorf("shutdown: %v", err)
	}
	return nil
}
