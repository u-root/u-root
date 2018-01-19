package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
)

const meminfoFile = "/proc/meminfo"

// meminfo returns a mapping that represents the fields contained in
// /proc/meminfo
func meminfo() (meminfomap, error) {
	buf, err := ioutil.ReadFile(meminfoFile)
	if err != nil {
		return nil, err
	}
	return meminfoFromBytes(buf)
}

// meminfoFromBytes returns a mapping that represents the fields contained in a
// byte stream with a content compatible with /proc/meminfo
func meminfoFromBytes(buf []byte) (meminfomap, error) {
	ret := make(meminfomap, 0)
	for _, line := range bytes.Split(buf, []byte{'\n'}) {
		kv := bytes.SplitN(line, []byte{':'}, 2)
		if len(kv) != 2 {
			// invalid line?
			continue
		}
		key := string(kv[0])
		tokens := bytes.SplitN(bytes.TrimSpace(kv[1]), []byte{' '}, 2)
		if len(tokens) > 0 {
			value, err := strconv.ParseUint(string(tokens[0]), 10, 64)
			if err != nil {
				return nil, err
			}
			ret[key] = value
		}
	}
	return ret, nil
}

// getMainMemInfo prints the physical memory information in the specified units. Only
// the relevant fields will be used from the input map.
func getMainMemInfo(m meminfomap, config *FreeConfig) (*mainMemInfo, error) {
	fields := []string{
		"MemTotal",
		"MemFree",
		"Buffers",
		"Cached",
		"Shmem",
		"SReclaimable",
		"MemAvailable",
	}
	if missingRequiredFields(m, fields) {
		return nil, fmt.Errorf("Missing required fields from meminfo")
	}

	// These values are expressed in kibibytes, convert to the desired unit
	memTotal := m["MemTotal"] << KB
	memFree := m["MemFree"] << KB
	memShared := m["Shmem"] << KB
	memCached := (m["Cached"] + m["SReclaimable"]) << KB
	memBuffers := (m["Buffers"]) << KB
	memUsed := memTotal - memFree - memCached - memBuffers
	if memUsed < 0 {
		memUsed = memTotal - memFree
	}
	memAvailable := m["MemAvailable"] << KB

	mmi := mainMemInfo{
		Total:     memTotal,
		Used:      memUsed,
		Free:      memFree,
		Shared:    memShared,
		Cached:    memCached,
		Buffers:   memBuffers,
		Available: memAvailable,
	}
	return &mmi, nil
}

// getSwapInfo prints the swap space information in the specified units. Only the
// relevant fields will be used from the input map.
func getSwapInfo(m meminfomap, config *FreeConfig) (*swapInfo, error) {
	fields := []string{
		"SwapTotal",
		"SwapFree",
	}
	if missingRequiredFields(m, fields) {
		return nil, fmt.Errorf("Missing required fields from meminfo")
	}
	// These values are expressed in kibibytes, convert to the desired unit
	swapTotal := m["SwapTotal"] << KB
	swapUsed := (m["SwapTotal"] - m["SwapFree"]) << KB
	swapFree := m["SwapFree"] << KB

	si := swapInfo{
		Total: swapTotal,
		Used:  swapUsed,
		Free:  swapFree,
	}
	return &si, nil
}
