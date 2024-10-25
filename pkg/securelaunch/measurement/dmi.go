// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package measurement

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	slaunch "github.com/u-root/u-root/pkg/securelaunch"
	"github.com/u-root/u-root/pkg/securelaunch/tpm"
	"github.com/u-root/u-root/pkg/smbios"
)

type fieldCluster struct {
	Label  string   `json:"label"`
	Fields []string `json:"fields"`
}

// DmiCollector describes the "dmi" portion of policy file.
type DmiCollector struct {
	Type     string         `json:"type"`
	Clusters []fieldCluster `json:"events"`
}

// NewDmiCollector extracts the "dmi" portion from the policy file and
// initializes a new DmiCollector structure.
//
// It returns an error if unmarshalling of DmiCollector fails.
func NewDmiCollector(config []byte) (Collector, error) {
	slaunch.Debug("New DMI Collector initialized")
	dc := new(DmiCollector)
	err := json.Unmarshal(config, &dc)
	if err != nil {
		return nil, err
	}
	return dc, nil
}

// Lookup table to convert a string to a DMI type.
// (taken from the dmidecode man page)
var typeTable = map[string]uint8{
	"bios":                             0,
	"system":                           1,
	"base board":                       2,
	"chassis":                          3,
	"processor":                        4,
	"memory controller":                5,
	"memory module":                    6,
	"cache":                            7,
	"port connector":                   8,
	"system slots":                     9,
	"on board devices":                 10,
	"oem strings":                      11,
	"system configuration options":     12,
	"bios language":                    13,
	"group associations":               14,
	"system event log":                 15,
	"physical memory array":            16,
	"memory device":                    17,
	"32-bit memory error":              18,
	"memory array mapped address":      19,
	"memory device mapped address":     20,
	"built-in pointing device":         21,
	"portable battery":                 22,
	"system reset":                     23,
	"hardware security":                24,
	"system power controls":            25,
	"voltage probe":                    26,
	"cooling device":                   27,
	"temperature probe":                28,
	"electrical current probe":         29,
	"out-of-band remote access":        30,
	"boot integrity services":          31,
	"system boot":                      32,
	"64-bit memory error":              33,
	"management device":                34,
	"management device component":      35,
	"management device threshold data": 36,
	"memory channel":                   37,
	"ipmi device":                      38,
	"power supply":                     39,
	"additional information":           40,
	"onboard device":                   41,
}

// parseTypeFilter looks up the DMI types and sets the corresponding entries
// in a map to true.
func parseTypeFilter(typeStrings []string) (map[smbios.TableType]bool, error) {
	types := map[smbios.TableType]bool{}
	for _, ts := range typeStrings {
		if tg, ok := typeTable[strings.ToLower(ts)]; ok {
			types[smbios.TableType(tg)] = true
		}
	}
	return types, nil
}

// Collect gets all smbios data, filters is based on the types provided in the
// policy file, then measures the filtered data into the TPM.
//
// It satisfies the Collector interface.
func (s *DmiCollector) Collect() error {
	slaunch.Debug("DMI Collector: Entering ")
	if s.Type != "dmi" {
		return errors.New("invalid type passed to a DmiCollector method")
	}

	var labels []string // collect all types entered by user in one slice
	for _, fieldCluster := range s.Clusters {
		labels = append(labels, fieldCluster.Label)
	}

	slaunch.Debug("DMI Collector: len(labels)=%d", len(labels))

	// lables would be []{BIOS, Chassis, Processor}
	typeFilter, err := parseTypeFilter(labels)
	if err != nil {
		return fmt.Errorf("invalid --type: %w", err)
	}

	slaunch.Debug("DMI Collector: len(typeFilter)=%d", len(typeFilter))

	si, err := smbios.FromSysfs()
	if err != nil {
		return fmt.Errorf("error parsing data: %w", err)
	}

	slaunch.Debug("DMI Collector: len(si.Tables)=%d", len(si.Tables))

	for _, t := range si.Tables {
		if len(typeFilter) != 0 && !typeFilter[t.Type] {
			continue
		}

		pt, err := smbios.ParseTypedTable(t)
		if err != nil {
			log.Printf("DMI Collector: skipping type %s, err=%v", t.Type, err)
			continue
		}

		slaunch.Debug(pt.String())
		b := []byte(pt.String())
		eventDesc := fmt.Sprintf("DMI Collector: Measured dmi label=[%v]", t.Type)
		if e := tpm.ExtendPCRDebug(pcr, bytes.NewReader(b), eventDesc); e != nil {
			log.Printf("DMI Collector: err =%v", e)
			return e // return error if any single type fails ..
		}
	}

	return nil
}
