package hwapi

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func isReservedType(regionType string) bool {
	switch t := strings.TrimSpace(regionType); t {
	case "reserved":
		return true
	case "Reserved":
		return true
	default:
		return false
	}
}

//IsReservedInE820 reads the e820 table exported via /sys/firmware/memmap and checks whether
// the range [start; end] is marked as reserved. Returns true if it is reserved,
// false if not.
func (t TxtAPI) IsReservedInE820(start uint64, end uint64) (bool, error) {
	if start > end {
		return false, fmt.Errorf("Invalid range")
	}

	dir, err := os.Open("/sys/firmware/memmap")
	if err != nil {
		return false, fmt.Errorf("Cannot access e820 table: %s", err)
	}

	subdirs, err := dir.Readdir(0)
	if err != nil {
		return false, fmt.Errorf("Cannot access e820 table: %s", err)
	}

	for _, subdir := range subdirs {
		if subdir.IsDir() {

			path := fmt.Sprintf("/sys/firmware/memmap/%s/type", subdir.Name())
			buf, err := ioutil.ReadFile(path)
			if err != nil {
				continue
			}

			if isReservedType(string(buf)) {
				path := fmt.Sprintf("/sys/firmware/memmap/%s/start", subdir.Name())
				thisStart, err := readHexInteger(path)
				if err != nil {
					continue
				}

				path = fmt.Sprintf("/sys/firmware/memmap/%s/end", subdir.Name())
				thisEnd, err := readHexInteger(path)
				if err != nil {
					continue
				}

				if thisStart <= start && thisEnd >= end {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func readHexInteger(path string) (uint64, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}

	ret, err := strconv.ParseUint(string(buf[:len(buf)-1]), 0, 64)
	if err != nil {
		return 0, err
	}

	return ret, nil
}
