// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brctl

import (
	"fmt"
	"net"
	"testing"
)

var (
	BRCTL_TEST_IFACE_0 = "eth0"
	BRCTL_TEST_IFACE_1 = "eth1"
	BRCTL_TEST_IFACES  = []string{BRCTL_TEST_IFACE_0, BRCTL_TEST_IFACE_1}

	BRCTL_TEST_BR_0    = "br0"
	BRCTL_TEST_BR_1    = "br1"
	BRCTL_TEST_BRIDGES = []string{BRCTL_TEST_BR_0, BRCTL_TEST_BR_1}
)

var test_fd = []struct {
	name    string
	input   string
	output  string
	wanterr bool
	err     error
}{
	{
		"forward delay 0",
		"0",
		"0",
		false,
		nil,
	},
	{
		"forward delay 1",
		"1",
		"1",
		false,
		nil,
	},
}

var test_str_to_jiffies = []struct {
	name     string
	duration string
	hz       int
	jiffies  int
	wanterr  bool
	err      error
}{
	{
		"1 second",
		"1s",
		100,
		100,
		false,
		nil,
	},
	{
		"1.5 seconds",
		"1.5s",
		100,
		150,
		false,
		nil,
	},
	{
		"1 minute",
		"1m",
		100,
		6000,
		false,
		nil,
	},
	{
		"1.5 minutes err",
		"1.5",
		100,
		0,
		true,
		fmt.Errorf("time: missing unit in duration \"1.5\""),
	},
}

func TestStringToJiffies(t *testing.T) {
	for _, tt := range test_str_to_jiffies {
		t.Run(tt.name, func(t *testing.T) {
			jiffies, err := stringToJiffies(tt.duration)
			if err != nil && !tt.wanterr {
				t.Fatalf("stringToJiffies(%q, %d) = '%v', want nil", tt.duration, tt.hz, err)
			}

			if err != nil && tt.wanterr {
				if err.Error() != tt.err.Error() {
					t.Fatalf("stringToJiffies(%q, %d) = '%v', want '%v'", tt.duration, tt.hz, err, tt.err)
				}
			}

			if jiffies != tt.jiffies {
				t.Fatalf("stringToJiffies(%q, %d) = %d, want %d", tt.duration, tt.hz, jiffies, tt.jiffies)
			}
		})
	}
}

func interfacesExist(ifs []string) error {
	for _, iface := range ifs {
		if _, err := net.InterfaceByName(iface); err != nil {
			return fmt.Errorf("interfacesExist: %w", err)
		}
	}
	return nil
}

// This function should be called at the start of each test to ensure that the environment is clean.
// It removes all bridges that were created during the test.
// It is assumed, that all necessary bridges and interfaces will be added per test case.
func clearEnv() error {
	// Check if interfaces exist
	if err := interfacesExist(BRCTL_TEST_IFACES); err != nil {
		return fmt.Errorf("clearEnv(%v): %w", BRCTL_TEST_IFACES, err)
	}

	// Remove all bridges
	for _, bridge := range BRCTL_TEST_BRIDGES {
		Delbr(bridge)
	}

	return nil
}

// All the following tests require virtual hardware to work properly.
// Hence they need to be run an VM, re-create the setup accordingly to the defined devices or adjust them to fit your platform.
// The following tests should be executed in a VM with root privileges, which is done in `brctl_integraion_test.go`
func TestAddbrDelbr(t *testing.T) {
	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	// Add bridges
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Addbr(bridge)
		if err.Error() != errno0.Error() {
			t.Fatalf("AddBr(%q) = %v, want nil", bridge, err)
		}
	}

	// Check if bridges were created successfully
	if err := interfacesExist(BRCTL_TEST_BRIDGES); err != nil {
		t.Fatalf("interfacesExist(%v) = %v, want nil", BRCTL_TEST_BRIDGES, err)
	}

	// Check if bridges were deleted successfully
	for _, iface := range BRCTL_TEST_BRIDGES {
		if _, err := net.InterfaceByName(iface); err == nil {
			t.Fatalf("net.InterfaceByName(%q) = nil, want an error", iface)
		}
	}
}

func TestIfDelif(t *testing.T) {
	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Addbr(bridge)
		if err != nil {
			t.Fatalf("AddBr(%q) = %v, want nil", bridge, err)
		}
	}
}

func TestSetageingTime(t *testing.T) {
	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	TEST_AGETIME := "1"

	// Add bridges
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Addbr(bridge)
		if err != errno0 {
			t.Fatalf("AddBr(%q) = %v, want nil", bridge, err)
		}
	}

	// Set ageing time
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Setageingtime(bridge, "1")
		if err != nil {
			t.Fatalf("Setageingtime(%q, \"1\") = %v, want nil", bridge, err)
		}
	}

	// Check sysfs for ageing time
	for _, bridge := range BRCTL_TEST_BRIDGES {
		ageingTime, err := getBridgeValue(bridge, BRCTL_AGEING_TIME)
		if err != nil {
			t.Fatalf("br_get_val(%q, \"ageing_time\") = %v, want nil", bridge, err)
		}

		if ageingTime != TEST_AGETIME {
			t.Fatalf("br_get_val(%q, \"ageing_time\") = %q, want \"1\"", bridge, ageingTime)
		}
	}
}

func TestShow(t *testing.T) {

}

func TestScpt(t *testing.T) {
	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	// Add bridges
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Addbr(bridge)
		if err != nil {
			t.Fatalf("AddBr(%q) = %v, want nil", bridge, err)
		}
	}

	for _, bridge := range BRCTL_TEST_BRIDGES {
		stp, err := getBridgeValue(bridge, BRCTL_STP_STATE)
		if err != nil {
			t.Fatalf("br_get_val(%q, \"stp_state\") = %v, want nil", bridge, err)
		}

		// By default STP is disabled
		if stp != "0" {
			t.Fatalf("br_get_val(%q, \"stp_state\") = %q, want \"0\"", bridge, stp)
		}

		// Enable STP
		err = Stp(bridge, "on")
		if err != nil {
			t.Fatalf("Stp(%q, \"on\") = %v, want nil", bridge, err)
		}

		stp, err = getBridgeValue(bridge, BRCTL_STP_STATE)
		if err != nil {
			t.Fatalf("br_get_val(%q, \"stp_state\") = %v, want nil", bridge, err)
		}
		if stp != "1" {
			t.Fatalf("br_get_val(%q, \"stp_state\") = %q, want \"1\"", bridge, stp)
		}
	}
}

func TestSetbridgeprio(t *testing.T) {
	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	// Add bridges
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Addbr(bridge)
		if err != nil {
			t.Fatalf("AddBr(%q) = %v, want nil", bridge, err)
		}
	}

	for _, bridge := range BRCTL_TEST_BRIDGES {
		// Set it to 0
		err := Setbridgeprio(bridge, "0")
		if err != nil {
			t.Fatalf("Setbridgeprio(%q, \"0\") = %v, want nil", bridge, err)
		}

		prio, err := getBridgeValue(bridge, BRCTL_BRIDGE_PRIO)
		if err != nil {
			t.Fatalf("br_get_val(%q, \"bridge_priority\") = %v, want nil", bridge, err)
		}

		if prio != "0" {
			t.Fatalf("br_get_val(%q, \"bridge_priority\") = %q, want \"0\"", bridge, prio)
		}
	}
}

func TestSetfd(t *testing.T) {
	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	// Add bridges
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Addbr(bridge)
		if err != nil {
			t.Fatalf("AddBr(%q) = %v, want nil", bridge, err)
		}
	}

	for _, bridge := range BRCTL_TEST_BRIDGES {
		for _, tt := range test_fd {
			t.Run(tt.name, func(t *testing.T) {
				err := Setfd(bridge, tt.input)
				if err != nil {
					t.Fatalf("Setfd(%q, %q) = %v, want nil", tt.input, tt.output, err)
				}

				// Check sysfs for forward delay
				for _, bridge := range BRCTL_TEST_BRIDGES {
					fd, err := getBridgeValue(bridge, BRCTL_FORWARD_DELAY)
					if err != nil {
						t.Fatalf("br_get_val(%q, \"forward_delay\") = %v, want nil", bridge, err)
					}

					if fd != tt.output {
						t.Fatalf("br_get_val(%q, \"forward_delay\") = %q, want %q", bridge, fd, tt.output)
					}
				}
			})
		}
	}
}

func TestSethello(t *testing.T) {
	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	// Add bridges
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Addbr(bridge)
		if err != nil {
			t.Fatalf("AddBr(%q) = %v, want nil", bridge, err)
		}
	}

	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Sethello(bridge, "1")
		if err != nil {
			t.Fatalf("Sethello(%q, \"1\") = %v, want nil", bridge, err)
		}

		hello, err := getBridgeValue(bridge, BRCTL_HELLO_TIME)
		if err != nil {
			t.Fatalf("br_get_val(%q, \"hello_time\") = %v, want nil", bridge, err)
		}

		if hello != "1" {
			t.Fatalf("br_get_val(%q, \"hello_time\") = %q, want \"1\"", bridge, hello)
		}
	}
}

func TestSetmaxage(t *testing.T) {
	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	// Add bridges
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Addbr(bridge)
		if err != nil {
			t.Fatalf("AddBr(%q) = %v, want nil", bridge, err)
		}
	}

	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Setmaxage(bridge, "0")
		if err != nil {
			t.Fatalf("Setmaxage(%q, \"1\") = %v, want nil", bridge, err)
		}

		maxAge, err := getBridgeValue(bridge, BRCTL_MAX_AGE)
		if err != nil {
			t.Fatalf("br_get_val(%q, \"max_age\") = %v, want nil", bridge, err)
		}

		if maxAge != "0" {
			t.Fatalf("br_get_val(%q, \"max_age\") = %q, want \"1\"", bridge, maxAge)
		}
	}
}

func TestSetpathcost(t *testing.T) {
	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	// Add bridges
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Addbr(bridge)
		if err != nil {
			t.Fatalf("AddBr(%q) = %v, want nil", bridge, err)
		}
	}

	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Setpathcost(bridge, "0", "1")
		if err != nil {
			t.Fatalf("Setpathcost(%q, \"0\", \"1\", \"1\") = %v, want nil", bridge, err)
		}

		pathCost, err := getBridgePort(bridge, "1", BRCTL_PATH_COST)
		if err != nil {
			t.Fatalf("br_get_val(%q, \"path_cost\") = %v, want nil", bridge, err)
		}

		if pathCost != "0" {
			t.Fatalf("br_get_val(%q, \"path_cost\") = %q, want \"0\"", bridge, pathCost)
		}
	}
}

func TestSetportprio(t *testing.T) {
	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	// Add bridges
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Addbr(bridge)
		if err != nil {
			t.Fatalf("AddBr(%q) = %v, want nil", bridge, err)
		}
	}

	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Setportprio(bridge, "1", "0")
		if err != nil {
			t.Fatalf("Setportprio(%q, \"1\", \"0\") = %v, want nil", bridge, err)
		}

		prio, err := getBridgePort(bridge, "1", BRCTL_PRIORITY)
		if err != nil {
			t.Fatalf("br_get_val(%q, \"port_priority\") = %v, want nil", bridge, err)
		}

		if prio != "0" {
			t.Fatalf("br_get_val(%q, \"port_priority\") = %q, want \"0\"", bridge, prio)
		}
	}
}

func TestHairpin(t *testing.T) {
	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	// Add bridges
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Addbr(bridge)
		if err != nil {
			t.Fatalf("AddBr(%q) = %v, want nil", bridge, err)
		}
	}

	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Hairpin(bridge, "1", "on")
		if err != nil {
			t.Fatalf("Hairpin(%q, \"1\", \"on\") = %v, want nil", bridge, err)
		}

		hairpin, err := getBridgePort(bridge, "1", BRCTL_HAIRPIN)
		if err != nil {
			t.Fatalf("br_get_val(%q, \"hairpin_mode\") = %v, want nil", bridge, err)
		}

		if hairpin != "1" {
			t.Fatalf("br_get_val(%q, \"hairpin_mode\") = %q, want \"1\"", bridge, hairpin)
		}
	}
}
