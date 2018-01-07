package main

import (
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
	err = printSwap(m, KB)
	if err != nil {
		t.Fatal(err)
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
	err = printSwap(m, KB)
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
SwapTotal:       8265724 kB
SwapFree:        8264956 kB
SReclaimable:     179852 kB`)
	m, err := meminfoFromBytes(input)
	if err != nil {
		t.Fatal(err)
	}
	err = printMem(m, KB)
	if err != nil {
		t.Fatal(err)
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
	err = printMem(m, KB)
	// should error out for the missing field
	if err == nil {
		t.Fatal("printMem: got no error when expecting one")
	}
}
