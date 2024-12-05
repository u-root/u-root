// Copyright 2018-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestMeminfoFromBytes(t *testing.T) {
	input := []byte(`MemTotal:        8052976 kB
MemFree:          721716 kB
MemAvailable:    2774100 kB
Buffers:          244880 kB
Cached:          3462124 kB
SwapTotal:       8265724 kB
SwapFree:        8264956 kB
SReclaimable:     179852 kB`)
	m, err := meminfoFromBytes(input)
	if err != nil {
		t.Fatal(err)
	}
	if m["MemFree"] != 721716 {
		t.Fatalf("MemFree: got %v, want 721716", m["MemFree"])
	}
	if m["MemAvailable"] != 2774100 {
		t.Fatalf("MemAvailable: got %v, want 2774100", m["MemAvailable"])
	}
	if m["Buffers"] != 244880 {
		t.Fatalf("Buffers: got %v, want 244880", m["Buffers"])
	}
	if m["Cached"] != 3462124 {
		t.Fatalf("Cached: got %v, want 3462124", m["Cached"])
	}
	if m["SwapTotal"] != 8265724 {
		t.Fatalf("SwapTotal: got %v, want 8265724", m["SwapTotal"])
	}
	if m["SwapFree"] != 8264956 {
		t.Fatalf("SwapFree: got %v, want 8264956", m["SwapFree"])
	}
	if m["SReclaimable"] != 179852 {
		t.Fatalf("SReclaimable: got %v, want 179852", m["SReclaimable"])
	}
}

func TestPrintSwap(t *testing.T) {
	input := []byte(`MemTotal:        8052976 kB
MemFree:          721716 kB
MemAvailable:    2774100 kB
Buffers:          244880 kB
Cached:          3462124 kB
SwapTotal:       8265724 kB
SwapFree:        8264956 kB
SReclaimable:     179852 kB`)
	m, err := meminfoFromBytes(input)
	if err != nil {
		t.Fatal(err)
	}
	si, err := getSwapInfo(m)
	if err != nil {
		t.Fatal(err)
	}
	if si.Total != 8464101376 {
		t.Fatalf("Swap.Total: got %v, want 8464101376", si.Total)
	}
	if si.Used != 786432 {
		t.Fatalf("Swap.Used: got %v, want 786432", si.Used)
	}
	if si.Free != 8463314944 {
		t.Fatalf("Swap.Free: got %v, want 8463314944", si.Free)
	}
}

func TestPrintSwapMissingFields(t *testing.T) {
	input := []byte(`MemTotal:        8052976 kB
MemFree:          721716 kB
MemAvailable:    2774100 kB
Buffers:          244880 kB
Cached:          3462124 kB
SwapTotal:       8265724 kB
SReclaimable:     179852 kB`)
	m, err := meminfoFromBytes(input)
	if err != nil {
		t.Fatal(err)
	}
	_, err = getSwapInfo(m)
	// should error out for the missing field
	if err == nil {
		t.Fatal("printSwap: got no error when expecting one")
	}
}

func TestPrintMem(t *testing.T) {
	input := []byte(`MemTotal:        8052976 kB
MemFree:          721716 kB
MemAvailable:    2774100 kB
Buffers:          244880 kB
Cached:          3462124 kB
Shmem:           1617788 kB
SwapTotal:       8265724 kB
SwapFree:        8264956 kB
SReclaimable:     179852 kB`)
	m, err := meminfoFromBytes(input)
	if err != nil {
		t.Fatal(err)
	}
	mmi, err := getMainMemInfo(m)
	if err != nil {
		t.Fatal(err)
	}
	if mmi.Total != 8246247424 {
		t.Fatalf("MainMem.Total: got %v, want 8246247424", mmi.Total)
	}
	if mmi.Free != 739037184 {
		t.Fatalf("MainMem.Free: got %v, want 739037184", mmi.Free)
	}
	if mmi.Used != 3527069696 {
		t.Fatalf("MainMem.Used: got %v, want 3527069696", mmi.Used)
	}
	if mmi.Shared != 1656614912 {
		t.Fatalf("MainMem.Shared: got %v, want 1656614912", mmi.Shared)
	}
	if mmi.Cached != 3729383424 {
		t.Fatalf("MainMem.Cached: got %v, want 3729383424", mmi.Cached)
	}
	if mmi.Buffers != 250757120 {
		t.Fatalf("MainMem.Buffers: got %v, want 250757120", mmi.Buffers)
	}
	if mmi.Available != 2840678400 {
		t.Fatalf("MainMem.Available: got %v, want 2840678400", mmi.Available)
	}
}

func TestPrintMemMissingFields(t *testing.T) {
	input := []byte(`MemTotal:        8052976 kB
MemFree:          721716 kB
MemAvailable:    2774100 kB
Buffers:          244880 kB
SwapTotal:       8265724 kB
SwapFree:        8264956 kB
SReclaimable:     179852 kB`)
	m, err := meminfoFromBytes(input)
	if err != nil {
		t.Fatal(err)
	}
	_, err = getMainMemInfo(m)
	// should error out for the missing field
	if err == nil {
		t.Fatal("printMem: got no error when expecting one")
	}
}

func TestParse(t *testing.T) {
	input := []byte(`MemTotal:        8052976 kB
MemFree:          721716 kB
MemAvailable:    2774100 kB
Buffers:          244880 kB
Cached:          3462124 kB
Shmem:           1617788 kB
SwapTotal:       8265724 kB
SwapFree:        8264956 kB
SReclaimable:     179852 kB`)
	m, err := meminfoFromBytes(input)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		expectedTotalMem  string
		expectedTotalSwap string
		o                 options
	}{
		{
			o:                 options{bytes: true},
			expectedTotalMem:  "8246247424",
			expectedTotalSwap: "8464101376",
		},
		{
			expectedTotalMem:  "8052976",
			expectedTotalSwap: "8265724",
		},
		{
			o:                 options{mbytes: true},
			expectedTotalMem:  "7864",
			expectedTotalSwap: "8071",
		},
		{
			o:                 options{gbytes: true},
			expectedTotalMem:  "7",
			expectedTotalSwap: "7",
		},
		{
			o:                 options{tbytes: true},
			expectedTotalMem:  "0",
			expectedTotalSwap: "0",
		},
		{
			o:                 options{human: true},
			expectedTotalMem:  "7.7GiB",
			expectedTotalSwap: "7.9GiB",
		},
	}

	for _, test := range tests {
		var stdout bytes.Buffer
		cmd, err := command(&stdout, test.o)
		if err != nil {
			t.Fatal(err)
		}

		err = cmd.parse(m)
		if err != nil {
			t.Fatal(err)
		}

		lines := strings.Split(stdout.String(), "\n")
		memFields := strings.Fields(lines[1])
		if memFields[1] != test.expectedTotalMem {
			t.Errorf("expected total %s, got %s", test.expectedTotalMem, memFields[1])
		}
		swapFields := strings.Fields(lines[2])
		if swapFields[1] != test.expectedTotalSwap {
			t.Errorf("expected total %s, got %s", test.expectedTotalMem, swapFields[1])
		}
	}
}

func TestOptionError(t *testing.T) {
	_, err := command(nil, options{kbytes: true, mbytes: true})
	if err != errMultipleUnits {
		t.Errorf("expected error: %v, got %v", errMultipleUnits, err)
	}
}
