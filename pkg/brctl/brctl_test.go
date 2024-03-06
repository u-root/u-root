package brctl

import (
	"net"
	"testing"

	"golang.org/x/sys/unix"
)

// Check if the list of interfaces exist in the system
func interfacesExist(ifs []string) error {
	for _, iface := range ifs {
		if _, err := net.InterfaceByName(iface); err != nil {
			return err
		}
	}
	return nil
}

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

func TestStrToTV(t *testing.T) {
	for _, tt := range test_str_to_tv {
		t.Run(tt.name, func(t *testing.T) {
			tv, err := str_to_tv(tt.input)
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
			jiffies := to_jiffies(tt.input, tt.hz)
			if jiffies != tt.jiffies {
				t.Errorf("to_jiffies(%v, %d) = %d, want %d", tt.input, tt.hz, jiffies, tt.jiffies)
			}
		})
	}
}

func TestFromJiffies(t *testing.T) {
	for _, tt := range test_to_jiffies {
		t.Run(tt.name, func(t *testing.T) {
			tv := from_jiffies(tt.jiffies, tt.hz)
			if tv.Sec != tt.input.Sec || tv.Usec != tt.input.Usec {
				t.Errorf("from_jiffies(%d, %d) = %v, want %v", tt.jiffies, tt.hz, tv, tt.input)
			}
		})
	}
}

// All the following tests require virtual hardware to work properly, hence they need to be run in a VM.
// This is done by the integration test in brctl_integration_test.go.
func TestAddbrDelbr(t *testing.T) {
	// Check if interfaces exist
	if err := interfacesExist(BRCTL_TEST_IFACES); err != nil {
		t.Fatalf("interfacesExist(%v) = %v, want nil", BRCTL_TEST_IFACES, err)
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

	// Cleanup the VM
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Delbr(bridge)
		if err != nil {
			t.Fatalf("DelBr(%q) = %v, want nil", bridge, err)
		}
	}

	// Check if bridges were deleted successfully
	for _, iface := range BRCTL_TEST_BRIDGES {
		if _, err := net.InterfaceByName(iface); err == nil {
			t.Fatalf("net.InterfaceByName(%q) = nil, want an error", iface)
		}
	}
}

func TestIfDelIf(t *testing.T) {
	// Check if interfaces exist
	if err := interfacesExist(BRCTL_TEST_IFACES); err != nil {
		t.Fatalf("interfacesExist(%v) = %v, want nil", BRCTL_TEST_IFACES, err)
	}

	// Add interface to bridge
	// brrctl addbr br0 eht0
	// brrctl addbr br1 eht1
	for _, bridge := range BRCTL_TEST_BRIDGES {
		err := Addbr(bridge)
		if err != nil {
			t.Fatalf("AddBr(%q) = %v, want nil", bridge, err)
		}
	}
}
