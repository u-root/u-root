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
	si, err := getSwapInfo(m, &FreeConfig{Unit: KB})
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
	_, err = getSwapInfo(m, &FreeConfig{Unit: KB})
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
	mmi, err := getMainMemInfo(m, &FreeConfig{Unit: KB})
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
	_, err = getMainMemInfo(m, &FreeConfig{Unit: KB})
	// should error out for the missing field
	if err == nil {
		t.Fatal("printMem: got no error when expecting one")
	}
}
