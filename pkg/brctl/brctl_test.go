// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brctl

import (
	"fmt"
	"net"
	"testing"

	"golang.org/x/sys/unix"
)

var test_str_to_tv = []struct {
	name  string
	input string
	tv    unix.Timeval
}{
	{
		name:  "1 second",
		input: "1.000000000",
		tv:    unix.Timeval{Sec: 1, Usec: 0},
	},
	{
		name:  "1.5 seconds",
		input: "1.500000000",
		tv:    unix.Timeval{Sec: 1, Usec: 500000},
	},
	{
		name:  "1.5 seconds",
		input: "1.500000000",
		tv:    unix.Timeval{Sec: 1, Usec: 500000},
	},
}

var test_to_jiffies = []struct {
	name    string
	input   unix.Timeval
	jiffies int
	hz      int
}{
	{
		name:    "1 second, 100Hz",
		input:   unix.Timeval{Sec: 1, Usec: 0},
		jiffies: 100,
		hz:      100,
	},
	{
		name:    "1.5 seconds, 100Hz",
		input:   unix.Timeval{Sec: 1, Usec: 500000},
		jiffies: 150,
		hz:      100,
	},
	{
		name:    "1.5 seconds, 1000Hz",
		input:   unix.Timeval{Sec: 1, Usec: 500000},
		jiffies: 1500,
		hz:      1000,
	},
}

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

func interfacesExist(ifs []string) error {
	for _, iface := range ifs {
		if _, err := net.InterfaceByName(iface); err != nil {
			return err
		}
	}
	return nil
}

// This function should be called at the start of each test to ensure that the environment is clean.
// It removes all bridges that were created during the test.
// It is assumed, that all necessary bridges and interfaces will be added per test case
func clearEnv() error {
	// Check if interfaces exist
	if err := interfacesExist(BRCTL_TEST_IFACES); err != nil {
		return fmt.Errorf("interfacesExist(%v) = %v, want nil", BRCTL_TEST_IFACES, err)
	}

	// Remove all bridges
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Delbr(bridge)
		if err != nil {
			return fmt.Errorf("Delbr(%q) = %v, want nil", bridge, err)
		}
	}

	// Check if bridges were deleted successfully
	for _, iface := range BRCTL_TEST_BRIDGES {
		if _, err := net.InterfaceByName(iface); err == nil {
			return fmt.Errorf("net.InterfaceByName(%q) = nil, want an error", iface)
		}
	}

	return nil
}

func TestStrToTV(t *testing.T) {
	for _, tt := range test_str_to_tv {
		t.Run(tt.name, func(t *testing.T) {
			tv, err := stringToTimeval(tt.input)
			if err != nil {
				t.Errorf("str_to_tv(%q) = %v, want nil", tt.input, err)
			}

			if tv.Sec != tt.tv.Sec || tv.Usec != tt.tv.Usec {
				t.Errorf("str_to_tv(%q) = %v, want %v", tt.input, tv, tt.tv)
			}
		})
	}
}

func TestToJiffies(t *testing.T) {
	for _, tt := range test_to_jiffies {
		t.Run(tt.name, func(t *testing.T) {
			jiffies := timevalToJiffies(tt.input, tt.hz)
			if jiffies != tt.jiffies {
				t.Errorf("to_jiffies(%v, %d) = %d, want %d", tt.input, tt.hz, jiffies, tt.jiffies)
			}
		})
	}
}

func TestFromJiffies(t *testing.T) {
	for _, tt := range test_to_jiffies {
		t.Run(tt.name, func(t *testing.T) {
			tv := jiffiesToTimeval(tt.jiffies, tt.hz)
			if tv.Sec != tt.input.Sec || tv.Usec != tt.input.Usec {
				t.Errorf("from_jiffies(%d, %d) = %v, want %v", tt.jiffies, tt.hz, tv, tt.input)
			}
		})
	}
}

// All the following tests require virtual hardware to work properly.
// Hence they need to be run an VM, re-create the setup accordingly to the defined devices or adjust them to fit your platform.
// The following tests should be executed in a VM with root privileges, which is done in `brctl_integraion_test.go`
func TestAddbrDelbr(t *testing.T) {
	if err := clearEnv(); err != nil {
		t.Fatalf("ClearEnv() = %v, want nil", err)
	}

	// Add bridges
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Addbr(bridge)
		if err != nil {
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
		t.Fatalf("ClearEnv() = %v, want nil", err)
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
		t.Fatalf("ClearEnv() = %v, want nil", err)
	}

	TEST_AGETIME := "1"

	// Add bridges
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Addbr(bridge)
		if err != nil {
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
		t.Fatalf("ClearEnv() = %v, want nil", err)
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
		t.Fatalf("ClearEnv() = %v, want nil", err)
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
		t.Fatalf("ClearEnv() = %v, want nil", err)
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
		t.Fatalf("ClearEnv() = %v, want nil", err)
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
		t.Fatalf("ClearEnv() = %v, want nil", err)
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
		t.Fatalf("ClearEnv() = %v, want nil", err)
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
		t.Fatalf("ClearEnv() = %v, want nil", err)
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
		t.Fatalf("ClearEnv() = %v, want nil", err)
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
