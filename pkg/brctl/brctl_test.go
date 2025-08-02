// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brctl

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/hugelgupf/vmtest/guest"
)

var (
	BRCTL_TEST_IFACE_0 = "eth0"
	BRCTL_TEST_IFACE_1 = "eth1"
	BRCTL_TEST_IFACES  = []string{BRCTL_TEST_IFACE_0, BRCTL_TEST_IFACE_1}

	BRCTL_TEST_BR_0    = "br0"
	BRCTL_TEST_BR_1    = "br1"
	BRCTL_TEST_BRIDGES = []string{BRCTL_TEST_BR_0, BRCTL_TEST_BR_1}
)

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
		"1",
		100,
		100,
		false,
		nil,
	},
	{
		"1.5 seconds",
		"1.5",
		100,
		150,
		false,
		nil,
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
	guest.SkipIfNotInVM(t)

	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	for _, bridge := range BRCTL_TEST_BRIDGES {
		if err := Addbr(bridge); err != nil {
			if err.Error() != errno0.Error() {
				t.Fatalf("AddBr(%q) = %v, want nil", bridge, err)
			}
		}
	}

	// Check if bridges were created successfully
	if err := interfacesExist(BRCTL_TEST_BRIDGES); err != nil {
		t.Fatalf("interfacesExist(%v) = %v, want nil", BRCTL_TEST_BRIDGES, err)
	}

	for _, bridge := range BRCTL_TEST_BRIDGES {
		if err := Delbr(bridge); err != nil {
			t.Fatalf("Delbr(%q) = %v, want nil", bridge, err)
		}
	}

	// Check if bridges were deleted successfully
	for _, iface := range BRCTL_TEST_BRIDGES {
		if _, err := net.InterfaceByName(iface); err == nil {
			t.Fatalf("net.InterfaceByName(%q) = nil, want an error", iface)
		}
	}
}

func TestIfDelif(t *testing.T) {
	guest.SkipIfNotInVM(t)
	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	if err := Addbr(BRCTL_TEST_BR_0); err != nil {
		t.Fatalf("AddBr(%q) = %v, want nil", BRCTL_TEST_BR_0, err)
	}

	if err := Addif(BRCTL_TEST_BR_0, BRCTL_TEST_IFACE_0); err != nil {
		t.Fatalf("Addif(%q, %q) = %v, want nil", BRCTL_TEST_BR_0, BRCTL_TEST_IFACE_0, err)
	}

	if err := Delif(BRCTL_TEST_BR_0, BRCTL_TEST_IFACE_0); err != nil {
		t.Fatalf("DelIf(%q, %q) = %v, want nil", BRCTL_TEST_BR_0, BRCTL_TEST_IFACE_0, err)
	}
}

func TestSetageingTime(t *testing.T) {
	guest.SkipIfNotInVM(t)

	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	TEST_AGETIME_STR := "1"
	TEST_AGETIME_INT := "1"

	TEST_AGETIME_JIFFIES_INT, err := stringToJiffies(TEST_AGETIME_STR)
	if err != nil {
		t.Fatalf("stringToJiffies(%q) = %v, want nil", TEST_AGETIME_STR, err)
	}

	TEST_AGETIME_JIFFIES_STR := fmt.Sprintf("%d", TEST_AGETIME_JIFFIES_INT)

	// Add bridges
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Addbr(bridge)
		if err != nil && err != errno0 {
			t.Fatalf("AddBr(%q) = %v, want nil", bridge, err)
		}
	}

	// Set ageing time
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := SetAgeingTime(bridge, TEST_AGETIME_STR)
		if err != nil {
			t.Fatalf("Setageingtime(%q, %q) = '%v', want nil", bridge, TEST_AGETIME_INT, err)
		}
	}

	// Check sysfs for ageing time
	for _, bridge := range BRCTL_TEST_BRIDGES {
		ageingTime, err := getBridgeValue(t, bridge, BRCTL_AGEING_TIME)
		if err != nil {
			t.Fatalf("br_get_val(%q, \"ageing_time\") = %v, want nil", bridge, err)
		}

		if ageingTime != TEST_AGETIME_JIFFIES_STR {
			t.Fatalf("br_get_val(%q, \"ageing_time\") = %q, want %q", bridge, ageingTime, TEST_AGETIME_JIFFIES_STR)
		}
	}
}

func TestShow(t *testing.T) {
	guest.SkipIfNotInVM(t)

	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	if err := Addbr(BRCTL_TEST_BR_0); err != nil {
		t.Fatalf("AddBr(%q) = %v, want nil", BRCTL_TEST_BR_0, err)
	}

	if err := Addif(BRCTL_TEST_BR_0, BRCTL_TEST_IFACE_0); err != nil {
		t.Fatalf("Addif(%q, %q) = %v, want nil", BRCTL_TEST_BR_0, BRCTL_TEST_IFACE_0, err)
	}

	var bufOut bytes.Buffer
	if err := Show(&bufOut, BRCTL_TEST_BR_0); err != nil {
		t.Errorf("Show(%q, %q)= %v", &bufOut, BRCTL_TEST_BR_0, err)
	}

	if err := Show(&bufOut); err != nil {
		t.Errorf("Show(%q, %q)= %v", &bufOut, BRCTL_TEST_BR_0, err)
	}
}

func TestShowSTP(t *testing.T) {
	guest.SkipIfNotInVM(t)

	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	if err := Addbr(BRCTL_TEST_BR_0); err != nil {
		t.Fatalf("AddBr(%q) = %v, want nil", BRCTL_TEST_BR_0, err)
	}

	if err := Addif(BRCTL_TEST_BR_0, BRCTL_TEST_IFACE_0); err != nil {
		t.Fatalf("Addif(%q, %q) = %v, want nil", BRCTL_TEST_BR_0, BRCTL_TEST_IFACE_0, err)
	}

	// Enable STP
	err := SetSTP(BRCTL_TEST_BR_0, "on")
	if err != nil {
		t.Fatalf("Stp(%q, \"on\") = %v, want nil", BRCTL_TEST_BR_0, err)
	}

	var bufOut bytes.Buffer
	if err := ShowStp(&bufOut, BRCTL_TEST_BR_0); err != nil {
		t.Errorf("Show(%q, %q)= %v", &bufOut, BRCTL_TEST_BR_0, err)
	}
}

func TestScpt(t *testing.T) {
	guest.SkipIfNotInVM(t)

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
		stp, err := getBridgeValue(t, bridge, BRCTL_STP_STATE)
		if err != nil {
			t.Fatalf("br_get_val(%q, \"stp_state\") = %v, want nil", bridge, err)
		}

		// By default STP is disabled
		if stp != "0" {
			t.Fatalf("br_get_val(%q, \"stp_state\") = %q, want \"0\"", bridge, stp)
		}

		// Enable STP
		err = SetSTP(bridge, "on")
		if err != nil {
			t.Fatalf("Stp(%q, \"on\") = %v, want nil", bridge, err)
		}

		stp, err = getBridgeValue(t, bridge, BRCTL_STP_STATE)
		if err != nil {
			t.Fatalf("br_get_val(%q, \"stp_state\") = %v, want nil", bridge, err)
		}
		if stp != "1" {
			t.Fatalf("br_get_val(%q, \"stp_state\") = %q, want \"1\"", bridge, stp)
		}
	}
}

func TestSetbridgeprio(t *testing.T) {
	guest.SkipIfNotInVM(t)

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
		err := SetBridgePrio(bridge, "0")
		if err != nil {
			t.Fatalf("Setbridgeprio(%q, \"0\") = %v, want nil", bridge, err)
		}

		prio, err := getBridgeValue(t, bridge, BRCTL_BRIDGE_PRIO)
		if err != nil {
			t.Fatalf("br_get_val(%q, \"bridge_priority\") = %v, want nil", bridge, err)
		}

		if prio != "0" {
			t.Fatalf("br_get_val(%q, \"bridge_priority\") = %q, want \"0\"", bridge, prio)
		}
	}
}

func TestSetfd(t *testing.T) {
	guest.SkipIfNotInVM(t)

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

	TEST_FD := "1"
	TEST_FD_JIFFIES, err := stringToJiffies(TEST_FD)
	if err != nil {
		t.Fatalf("stringToJiffies(%q) = %v, want nil", TEST_FD, err)
	}
	TEST_FD_JIFFIES_STR := strconv.Itoa(TEST_FD_JIFFIES)

	err = SetForwardDelay(BRCTL_TEST_BR_0, TEST_FD)
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}

	// Check sysfs for forward delay
	fd, err := getBridgeValue(t, BRCTL_TEST_BR_0, BRCTL_FORWARD_DELAY)
	if err != nil {
		t.Fatalf("br_get_val(%q, \"forward_delay\") = %v, want nil, err = %v", BRCTL_TEST_BR_0, fd, err)
	}

	if fd != TEST_FD_JIFFIES_STR {
		t.Fatalf("br_get_val(%q, \"forward_delay\") = %q, want %q", BRCTL_TEST_BR_0, fd, 100)
	}
}

func TestSethello(t *testing.T) {
	guest.SkipIfNotInVM(t)

	TEST_SETHELLO_TIME := "1"
	TEST_SETHELLO_JIFFIES, err := stringToJiffies(TEST_SETHELLO_TIME)
	if err != nil {
		t.Fatalf("stringToJiffies(%q) = %v, want nil", TEST_SETHELLO_TIME, err)
	}

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
		err := SetHello(bridge, TEST_SETHELLO_TIME)
		if err != nil {
			t.Fatalf("Sethello(%q, %q) = %v, want nil", bridge, TEST_SETHELLO_TIME, err)
		}

		hello, err := getBridgeValue(t, bridge, BRCTL_HELLO_TIME)
		if err != nil {
			t.Fatalf("br_get_val(%q, \"hello_time\") = %v, want nil", bridge, err)
		}

		jiffies := fmt.Sprintf("%d", TEST_SETHELLO_JIFFIES)
		if hello != jiffies {
			t.Fatalf("br_get_val(%q, \"hello_time\") = %q, want %q", bridge, jiffies, hello)
		}
	}
}

// TODO: also the original brctl returns on a modern linux system -ERANGE which looks as if it is not supported
func TestSetmaxage(t *testing.T) {
	t.Skip()

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
		err := SetMaxAge(bridge, "1s")
		if err != nil {
			t.Fatalf("Setmaxage(%q, \"1\") = %v, want nil", bridge, err)
		}

		maxAge, err := getBridgeValue(t, bridge, BRCTL_MAX_AGE)
		if err != nil {
			t.Fatalf("br_get_val(%q, \"max_age\") = %v, want nil", bridge, err)
		}

		if maxAge != "2000" {
			t.Fatalf("br_get_val(%q, \"max_age\") = %q, want \"1\"", bridge, maxAge)
		}
	}
}

func TestSetpathcost(t *testing.T) {
	guest.SkipIfNotInVM(t)

	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	// Add bridges
	err := Addbr(BRCTL_TEST_BR_0)
	if err != nil {
		t.Fatalf("AddBr(%q) = %v, want nil", BRCTL_TEST_BR_0, err)
	}

	err = Addif(BRCTL_TEST_BR_0, BRCTL_TEST_IFACE_0)
	if err != nil {
		t.Fatalf("Addif(%q, %q) = %v, want nil", BRCTL_TEST_BR_0, BRCTL_TEST_IFACE_0, err)
	}

	// Set Port for test
	TEST_BRIDGE := BRCTL_TEST_BR_0
	TEST_PORT := BRCTL_TEST_IFACE_0
	TEST_COST := "1"

	err = SetPathCost(TEST_BRIDGE, TEST_PORT, TEST_COST)
	if err != nil {
		t.Fatalf("Setpathcost(%q, %q, %v) = %v, want nil", TEST_BRIDGE, TEST_PORT, TEST_COST, err)
	}

	pathCost, err := os.ReadFile(path.Join(BRCTL_SYS_NET, TEST_PORT, "brport", "path_cost"))
	if err != nil {
		t.Fatalf("os.ReadFile(%q) = %v, want nil", path.Join(BRCTL_SYS_NET, TEST_PORT, "brport", "path_cost"), err)
	}

	// trim the '\n' from the output
	if strings.TrimSuffix(string(pathCost), "\n") != TEST_COST {
		t.Fatalf("br_get_val(%q, \"path_cost\") = %q, want %q", TEST_BRIDGE, pathCost, TEST_COST)
	}
}

// /sys/class/net/dummy0/brport
func TestSetportprio(t *testing.T) {
	guest.SkipIfNotInVM(t)

	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	TEST_BRIDGE := BRCTL_TEST_BR_0
	TEST_PORT := BRCTL_TEST_IFACE_0
	TEST_PRIO := "1"

	// Add bridge
	err := Addbr(TEST_BRIDGE)
	if err != nil {
		t.Fatalf("AddBr(%q) = %v, want nil", TEST_BRIDGE, err)
	}
	err = Addif(TEST_BRIDGE, BRCTL_TEST_IFACE_0)
	if err != nil {
		t.Fatalf("Addif(%q, %q) = %v, want nil", TEST_BRIDGE, BRCTL_TEST_IFACE_0, err)
	}

	err = SetPortPrio(TEST_BRIDGE, TEST_PORT, TEST_PRIO)
	if err != nil {
		t.Fatalf("Setportprio(%q, %q, %v) = %v, want nil", TEST_BRIDGE, TEST_PORT, TEST_PRIO, err)
	}

	prio, err := getPortValue(t, TEST_PORT, BRCTL_PRIORITY)
	if err != nil {
		t.Fatalf("br_get_val(%q, \"port_priority\") = %v, want nil", TEST_BRIDGE, err)
	}

	if strings.TrimSuffix(prio, "\n") != TEST_PRIO {
		t.Fatalf("br_get_val(%q, \"port_priority\") = %q, want \"0\"", TEST_BRIDGE, prio)
	}
}

func TestHairpin(t *testing.T) {
	guest.SkipIfNotInVM(t)

	if err := clearEnv(); err != nil {
		t.Skip(err)
	}

	TEST_BRIDGE := BRCTL_TEST_BR_0
	TEST_PORT := BRCTL_TEST_IFACE_0
	TEST_VALUE := "1"

	// Add bridges
	err := Addbr(TEST_BRIDGE)
	if err != nil {
		t.Fatalf("AddBr(%q) = %v, want nil", TEST_BRIDGE, err)
	}

	err = Addif(BRCTL_TEST_BR_0, BRCTL_TEST_IFACE_0)
	if err != nil {
		t.Fatalf("Addif(%q, %q) = %v, want nil", BRCTL_TEST_BR_0, BRCTL_TEST_IFACE_0, err)
	}

	err = Hairpin(TEST_BRIDGE, TEST_PORT, "on")
	if err != nil {
		t.Fatalf("Hairpin(%q, %q, \"on\") = %v, want nil", TEST_BRIDGE, TEST_PORT, err)
	}

	hairpin, err := getPortValue(t, TEST_PORT, BRCTL_HAIRPIN)
	if err != nil {
		t.Fatalf("br_get_val(%q, \"hairpin_mode\") = %v, want nil", TEST_BRIDGE, err)
	}

	if strings.TrimSuffix(hairpin, "\n") != TEST_VALUE {
		t.Fatalf("br_get_val(%q, \"hairpin_mode\") = %q, want %q", TEST_BRIDGE, hairpin, TEST_VALUE)
	}
}

func getBridgeValue(t *testing.T, bridge string, setting string) (string, error) {
	t.Helper()

	out, err := os.ReadFile(filepath.Join(BRCTL_SYS_NET, bridge, "bridge", setting))
	if err != nil {
		return "", err
	}
	// For some reason these values have a '\n' (0x0a) as a suffix, so we need to trim it
	return strings.TrimSuffix(string(out), "\n"), nil
}

func getPortValue(t *testing.T, port string, setting string) (string, error) {
	t.Helper()

	out, err := os.ReadFile(filepath.Join(BRCTL_SYS_NET, port, "brport", setting))
	if err != nil {
		return "", err
	}
	return string(out), nil
}
