// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package measurement

import (
	"log"
	"testing"

	slaunch "github.com/u-root/u-root/pkg/securelaunch"
)

var storage_config = `
{
	"type": "storage",
	"paths": [ "sda1" ]
}`

var dmi_config = `
{
	"type": "dmi",
	"events": [
	  {
		"label": "BIOS",
		"fields": []
	  },
	  {
		"label": "System",
		"fields": []
	  },
	  {
		"label": "Processor",
		"fields": []
	  }
	]
}`

var files_config = `
{
	"type": "files",
	"paths": [ "sda1:/opc/foo" ]
}`

var cpuid_config = `
{
	"type": "cpuid",
	"location": "sda2:/cpuid.txt"
}`

func TestGetCollector(t *testing.T) {
	slaunch.Debug = log.Printf

	collector, err := GetCollector([]byte(storage_config))
	if err != nil {
		t.Fatalf(`GetCollector([]byte(storage_config)) = %v, not nil; collector = %v`, err, collector)
	}

	collector, err = GetCollector([]byte(dmi_config))
	if err != nil {
		t.Fatalf(`GetCollector([]byte(dmi_config)) = %v, not nil; collector = %v`, err, collector)
	}

	collector, err = GetCollector([]byte(files_config))
	if err != nil {
		t.Fatalf(`GetCollector([]byte(files_config)) = %v, not nil; collector = %v`, err, collector)
	}

	collector, err = GetCollector([]byte(cpuid_config))
	if err != nil {
		t.Fatalf(`GetCollector([]byte(cpuid_config)) = %v, not nil; collector = %v`, err, collector)
	}
}
