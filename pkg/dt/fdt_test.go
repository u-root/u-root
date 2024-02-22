// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dt

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadFDT(t *testing.T) {
	jsonData, err := os.ReadFile("testdata/fdt.json")
	if err != nil {
		t.Fatal(err)
	}
	testData := &FDT{}
	if err := json.Unmarshal(jsonData, testData); err != nil {
		t.Fatal(err)
	}

	// 1. Load by path given and succeed.
	dtb, err := os.Open("testdata/fdt.dtb")
	if err != nil {
		t.Fatal(err)
	}
	fdt, err := LoadFDT(dtb)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(fdt, testData) {
		got, err := json.MarshalIndent(fdt, "", "    ")
		if err != nil {
			t.Fatal(err)
		}
		t.Errorf(`Read("fdt.dtb") = %s \n, want %s`, got, jsonData)
	}

	dir := t.TempDir()
	nonexistDTB := filepath.Join(dir, "xxx")
	// 2. Fallback to read from sys fs, and sys fs reading also failed.
	if _, err = LoadFDT(nil, nonexistDTB); !errors.Is(err, ErrNoValidReaders) {
		t.Errorf("LoadFDT(%s) got %v, want %v", nonexistDTB, err, ErrNoValidReaders)
	}

	// 3. Fallback to read from sys fs, and succeed.
	fdt, err = LoadFDT(nil, "testdata/fdt.dtb")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(fdt, testData) {
		got, err := json.MarshalIndent(fdt, "", "    ")
		if err != nil {
			t.Fatal(err)
		}
		t.Errorf(`Read("fdt.dtb") = %s \n, want %s`, got, jsonData)
	}
}

func TestRead(t *testing.T) {
	f, err := os.Open("testdata/fdt.dtb")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	jsonData, err := os.ReadFile("testdata/fdt.json")
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
		t.Errorf(`Read("fdt.dtb") = %s`, got)
		t.Errorf(`want %s`, jsonData)
	}
}

// TestParity tests that the fdt Read+Write operations are compatible with your
// system's fdtdump command.
func TestParity(t *testing.T) {
	// TODO: I'm convinced my system's fdtdump command is broken.
	t.Skip()

	// Read and write the fdt.
	fdt, err := New(WithFileName("testdata/fdt.dtb"))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	dtb := filepath.Join(dir, "fdt2.dtb")
	f, err := os.Create(dtb)
	if err != nil {
		t.Fatal(err)
	}
	_, err = fdt.Write(f)
	f.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Run your system's fdtdump command.
	dts := filepath.Join(dir, "fdt2.dts")
	f, err = os.Create(dts)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("fdtdump", dtb)
	cmd.Stdout = f
	err = cmd.Run()
	f.Close()
	if err != nil {
		t.Fatal(err) // TODO: skip if system does not have fdtdump
	}

	// This used to run diff, has to be done better. It's not even working now so.
	if false {
		cmd = exec.Command("diff", "testdata/fdt.dts", "testdata/fdt2.dts")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

func TestFindNode(t *testing.T) {
	f, err := os.Open("testdata/fdt.dtb")
	if err != nil {
		t.Fatal(err)
	}

	fdt, err := ReadFDT(f)
	f.Close()
	if err != nil {
		t.Fatal(err)
	}

	n, ok := fdt.NodeByName("psci")
	if !ok {
		t.Fatalf("Finding psci in %s: got false, want true", fdt)
	}
	t.Logf("Got the node: %s", n)
}

func TestFindAllNode(t *testing.T) {
	f, err := os.Open("testdata/fdt.dtb")
	if err != nil {
		t.Fatal(err)
	}

	fdt, err := ReadFDT(f)
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	const expectedVirtNodes = 32
	nodes, err := fdt.Root().FindAll(func(n *Node) bool {
		return strings.HasPrefix(n.Name, "virtio_mmio")
	})
	if err != nil {
		t.Fatalf("Finding all virtio_mmio in %s: got err %v, want nil", fdt, err)
	}

	if len(nodes) != expectedVirtNodes {
		t.Fatalf("Finding all virtio_mmio in %s: got returned %d nodes, want %d", fdt, len(nodes), expectedVirtNodes)
	}
}

func TestFindProperty(t *testing.T) {
	f, err := os.Open("testdata/fdt.dtb")
	if err != nil {
		t.Fatal(err)
	}
	fdt, err := ReadFDT(f)
	f.Close()
	if err != nil {
		t.Fatal(err)
	}

	n, ok := fdt.NodeByName("psci")
	if !ok {
		t.Fatalf("Finding psci in %s: got false, want true", fdt)
	}
	t.Logf("Got the node: %s", n)
	l := "migrate"
	p, ok := n.LookProperty(l)
	if !ok {
		t.Fatalf("Find property %q in %s: got false, want true", l, n)
	}
	v := []byte{0x84, 0, 0, 0x5}
	if !bytes.Equal(p.Value, v) {
		t.Fatalf("Checking value of %s: got %q, want %q", p.Name, p.Value, v)
	}
	l = "bogosity"
	if _, ok = n.LookProperty(l); ok {
		t.Fatalf("Find property %q in %s: got true, want false", l, n)
	}
}

func TestWalk(t *testing.T) {
	f, err := os.Open("testdata/fdt.dtb")
	if err != nil {
		t.Fatal(err)
	}
	fdt, err := ReadFDT(f)
	f.Close()
	if err != nil {
		t.Fatal(err)
	}

	b, err := fdt.Root().Walk("psci").Property("migrate").AsBytes()
	if err != nil {
		t.Fatalf("Walk to psci/migrate: got %v, want nil", err)
	}
	v := []byte{0x84, 0, 0, 0x5}
	if !bytes.Equal(b, v) {
		t.Fatalf("Checking value of psci/migrate: got %q, want %q", b, v)
	}
}
