// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gpio provides functions for interacting with GPIO pins via the
// GPIO Sysfs Interface for Userspace.
package gpio

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const gpioPath = "/sys/class/gpio"

// Value represents the value of a gpio pin
type Value bool

// Gpio pin values can either be low (0) or high (1)
const (
	Low  Value = false
	High Value = true
)

func (v Value) Dir() string {
	if v == Low {
		return "low"
	}
	return "high"
}

func (v Value) String() string {
	if v == Low {
		return "0"
	}
	return "1"
}

// SetOutputValue configures the gpio as an output pin with the given value.
func SetOutputValue(pin int, val Value) error {
	dir := val.Dir()
	path := filepath.Join(gpioPath, fmt.Sprintf("gpio%d", pin), "direction")
	outFile, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", path, err)
	}
	defer outFile.Close()
	if _, err := outFile.WriteString(dir); err != nil {
		return fmt.Errorf("failed to set gpio %d to %s: %v", pin, dir, err)
	}
	return nil
}

// ReadValue returns the value of the given gpio pin. If the read was
// unsuccessful, it returns a value of Low and the associated error.
func ReadValue(pin int) (Value, error) {
	path := filepath.Join(gpioPath, fmt.Sprintf("gpio%d", pin), "value")
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return Low, fmt.Errorf("failed to read value of gpio %d: %v", pin, err)
	}
	switch string(buf) {
	case "0\n":
		return Low, nil
	case "1\n":
		return High, nil
	}
	return Low, fmt.Errorf("invalid value of gpio %d: %s", pin, string(buf))
}

// Export enables access to the given gpio pin.
func Export(pin int) error {
	path := filepath.Join(gpioPath, "export")
	outFile, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", path, err)
	}
	defer outFile.Close()
	if _, err := outFile.WriteString(strconv.Itoa(pin)); err != nil {
		return fmt.Errorf("failed to export gpio %d: %v", pin, err)
	}
	return nil
}
