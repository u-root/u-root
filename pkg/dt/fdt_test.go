// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dt

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

var tests = []string{
	"testdata/qemu_fdt",
	"testdata/rpi_fdt",
}

func TestRead(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			f, err := os.Open(tt + ".dtb")
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			jsonData, err := ioutil.ReadFile(tt + ".json")
			if err != nil {
				t.Fatal(err)
			}
			testData := &FDT{}
			if err := json.Unmarshal(jsonData, testData); err != nil {
				t.Fatal(err)
			}

			fdt, err := ReadFDT(f)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(fdt, testData) {
				got, err := json.MarshalIndent(fdt, "", "    ")
				if err != nil {
					t.Fatal(err)
				}
				t.Errorf("Read(%q) = %s", tt+".dtb", got)
				t.Errorf("want %s", jsonData)
			}
		})
	}
}

// TestParity tests that the fdt Read+Write operations are compatible with your
// system's fdtdump command.
func TestParity(t *testing.T) {
	// TODO: I'm convinced my system's fdtdump command is broken. CONFIRMED.
	//       fdtdump v1.4.0 is broken
	//       https://bugs.launchpad.net/ubuntu/+source/device-tree-compiler/+bug/1668291
	t.Skip()

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			// Read and write the fdt.
			f, err := os.Open(tt + ".dtb")
			if err != nil {
				t.Fatal(err)
			}
			fdt, err := ReadFDT(f)
			f.Close()
			if err != nil {
				t.Fatal(err)
			}
			f, err = os.Create(tt + ".dtb.output")
			if err != nil {
				t.Fatal(err)
			}
			_, err = fdt.Write(f)
			f.Close()
			if err != nil {
				t.Fatal(err)
			}

			// Run your system's fdtdump command.
			f, err = os.Create(tt + ".dts.output")
			if err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command("fdtdump", tt+".dtb.output")
			cmd.Stdout = f
			err = cmd.Run()
			f.Close()
			if err != nil {
				t.Fatal(err) // TODO: skip if system does not have fdtdump
			}

			cmd = exec.Command("diff", tt+".dts", tt+".dts.output")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		})
	}
}
