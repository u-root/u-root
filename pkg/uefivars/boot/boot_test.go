// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package boot

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/uefivars"
	"github.com/u-root/u-root/pkg/uefivars/vartest"
)

// main is needed to extract the testdata from a zip to temp dir, and to clean
// up the temp dir after
func TestMain(m *testing.M) {
	var setups = []struct {
		name      string
		zippath   string
		testevd   bool
		testevfsd bool
	}{
		{
			name:      "test efivars backend",
			zippath:   "../testdata/sys_fw_efi_vars.zip",
			testevd:   true,
			testevfsd: false,
		},
		{
			name:      "test efivarfs backend",
			zippath:   "../testdata/sys_fw_efivarfs.zip",
			testevd:   false,
			testevfsd: true,
		},
	}
	for _, w := range setups {
		efiVarDir, cleanup, err := vartest.SetupVarZip(w.zippath)
		if err != nil {
			log.Fatalf("couldn't setup boot EFI variable testdata: %v", err)
		}

		// If there is only one directory in the zip, use that as the root.
		if entries, err := os.ReadDir(efiVarDir); err == nil && len(entries) == 1 && entries[0].IsDir() {
			efiVarDir = filepath.Join(efiVarDir, entries[0].Name())
		}

		uefivars.EfiVarDir = "/tmp/invalid"
		uefivars.EfiVarfsDir = "/tmp/invalid"
		if w.testevd {
			uefivars.EfiVarDir = efiVarDir
		}
		if w.testevfsd {
			uefivars.EfiVarfsDir = efiVarDir
		}
		rc := m.Run()

		cleanup()
		if rc != 0 {
			os.Exit(rc)
		}

	}
	os.Exit(0)
}

// func ReadBootVar(num uint16) (b BootVar)
func TestReadBootVar(t *testing.T) {
	var n uint16
	var strs []string
	for n = 0; n < 11; n++ {
		b, err := ReadBootVar(n)
		if err != nil {
			t.Error(err)
		}
		strs = append(strs, b.String())
	}
	if t.Failed() {
		for _, s := range strs {
			t.Log(s)
		}
	}
}

// func AllBootEntryVars() (list []BootEntryVar)
func TestAllBootEntryVars(t *testing.T) {
	bevs := AllBootEntryVars()
	if len(bevs) != 11 {
		for i, e := range bevs {
			t.Logf("#%d: %s", i, e)
		}
		t.Errorf("expected 11 boot vars, got %d", len(bevs))
	}
}

// func AllBootVars() (list EfiVars)
func TestAllBootVars(t *testing.T) {
	n := 14
	bvs := AllBootVars()
	if len(bvs) != n {
		t.Errorf("expected %d boot vars, got %d", n, len(bvs))
	}
	be := bvs.Filter(BootEntryFilter)
	if len(be) != n-3 {
		t.Errorf("expected %d entries, got %d", n-3, len(be))
	}
	// find boot vars that are not boot entries
	nbe := bvs.Filter(uefivars.NotFilter(BootEntryFilter))
	want := []string{"BootCurrent", "BootOptionSupport", "BootOrder"}
	if len(nbe) != len(want) {
		t.Fatalf("want %d got %d", len(want), len(nbe))
	}
	for i, bv := range nbe {
		s := bv.Name
		if i >= len(want) || s != want[i] {
			t.Errorf("%d: %s", i, s)
		}
	}
}

// func ReadCurrentBootVar() (b *BootEntryVar)
func TestReadCurrentBootVar(t *testing.T) {
	v, err := ReadCurrentBootVar()
	if err != nil {
		t.Error(err)
	}

	if v == nil {
		t.Fatal("nil")
	}
	if v.Number != 10 {
		t.Errorf("expected 10, got %d", v.Number)
	}
	if t.Failed() {
		t.Log(v)
	}
}

// func BootCurrent(vars uefivars.EfiVars) *BootCurrentVar
func TestBootCurrent(t *testing.T) {
	bc := BootCurrent(AllBootVars())
	if bc == nil {
		t.Fatal("nil")
	}
	var want uint16 = 10
	if bc.Current != want {
		t.Errorf("want %d got %d", want, bc.Current)
	}
}
