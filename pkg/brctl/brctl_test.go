package brctl

import (
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
