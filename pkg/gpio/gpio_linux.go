// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gpio provides functions for interacting with GPIO pins via the
// GPIO Sysfs Interface for Userspace.
package gpio

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const gpioPath = "/sys/class/gpio"

// Value represents the value of a gpio pin
type Value bool

// Gpio pin values can either be low (0) or high (1)
const (
	Low  Value = false
	High Value = true
)

// Dir returns the representation that sysfs likes to use.
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

func readInt(filename string) (int, error) {
	// Get base offset (the first GPIO managed by this chip)
	buf, err := os.ReadFile(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to read integer out of %s: %w", filename, err)
	}
	baseStr := strings.TrimSpace(string(buf))
	num, err := strconv.Atoi(baseStr)
	if err != nil {
		return 0, fmt.Errorf("could not convert %s contents %s to integer: %w", filename, baseStr, err)
	}
	return num, nil
}

// GetPinID computes the sysfs pin ID for a specific port on a specific GPIO
// controller chip. The controller arg is matched to a gpiochip's label in
// sysfs. GetPinID gets the base offset of that chip, and adds the specific
// pin number.
func GetPinID(controller string, pin uint) (int, error) {
	controllers, err := filepath.Glob(fmt.Sprintf("%s/gpiochip*", gpioPath))
	if err != nil {
		return 0, err
	}

	for _, c := range controllers {
		// Get label (name of the controller)
		buf, err := os.ReadFile(filepath.Join(c, "label"))
		if err != nil {
			return 0, fmt.Errorf("failed to read label of %s: %w", c, err)
		}
		label := strings.TrimSpace(string(buf))

		// Check that this is the controller we want
		if label != controller {
			continue
		}

		// Get base offset (the first GPIO managed by this chip)
		base, err := readInt(filepath.Join(c, "base"))
		if err != nil {
			return 0, fmt.Errorf("failed to read base: %w", err)
		}

		// Get the number of GPIOs managed by this chip.
		ngpio, err := readInt(filepath.Join(c, "ngpio"))
		if err != nil {
			return 0, fmt.Errorf("failed to read number of gpios: %w", err)
		}
		if int(pin) >= ngpio {
			return 0, fmt.Errorf("requested pin %d of controller %s, but controller only has %d pins", pin, controller, ngpio)
		}

		return base + int(pin), nil
	}

	return 0, fmt.Errorf("could not find controller %s", controller)
}

// SetOutputValue configures the gpio as an output pin with the given value.
func SetOutputValue(pin int, val Value) error {
	dir := val.Dir()
	path := filepath.Join(gpioPath, fmt.Sprintf("gpio%d", pin), "direction")
	outFile, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer outFile.Close()
	if _, err := outFile.WriteString(dir); err != nil {
		return fmt.Errorf("failed to set gpio %d to %s: %w", pin, dir, err)
	}
	return nil
}

// ReadValue returns the value of the given gpio pin. If the read was
// unsuccessful, it returns a value of Low and the associated error.
func ReadValue(pin int) (Value, error) {
	path := filepath.Join(gpioPath, fmt.Sprintf("gpio%d", pin), "value")
	buf, err := os.ReadFile(path)
	if err != nil {
		return Low, fmt.Errorf("failed to read value of gpio %d: %w", pin, err)
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
		return fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer outFile.Close()
	if _, err := outFile.WriteString(strconv.Itoa(pin)); err != nil {
		return fmt.Errorf("failed to export gpio %d: %w", pin, err)
	}
	return nil
}
