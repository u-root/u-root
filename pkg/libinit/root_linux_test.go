// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libinit

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/u-root/u-root/pkg/cmdline"
)

func TestLoadModule(t *testing.T) {
	var loadedModules []string
	loader := &InitModuleLoader{
		Cmdline: cmdline.NewCmdLine(),
		Prober: func(name, params string) error {
			loadedModules = append(loadedModules, name)
			return nil
		},
	}

	expectedModules := []string{"test", "something-test"}
	InstallModules(loader, expectedModules)
	if diff := cmp.Diff(expectedModules, loadedModules); diff != "" {
		t.Fatalf("unexpected difference of loaded modules (-want, +got): %v", diff)
	}
}

func TestModuleConf(t *testing.T) {
	toBytes := func(s string) []byte {
		return bytes.NewBufferString(s).Bytes()
	}
	files := []struct {
		Name    string
		Content string
		Modules []string
	}{
		{
			Name:    "test.conf",
			Content: `something`,
			Modules: []string{"something"},
		},
		{
			Name: "test2.conf",
			Content: `module1
# not a module
module2`,
			Modules: []string{"module1", "module2"},
		},
	}

	dir := t.TempDir()

	var checkModules []string
	for _, file := range files {
		t.Run(file.Name, func(t *testing.T) {
			p := filepath.Join(dir, file.Name)
			if err := os.WriteFile(p, toBytes(file.Content), 0o644); err != nil {
				t.Fatal(err)
			}
			checkModules = append(checkModules, file.Modules...)
		})
	}

	moduleConfPattern := filepath.Join(dir, "*.conf")
	modules, err := GetModulesFromConf(moduleConfPattern)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(checkModules, modules); diff != "" {
		t.Fatalf("unexpected difference of loaded modules (-want, +got): %v", diff)
	}
}

func TestCmdline(t *testing.T) {
	cline := &cmdline.CmdLine{
		AsMap: map[string]string{
			"modules_load": "test",
			"test.key1":    "value1",
			"test.key2":    "value2",
			"test.key3":    "value3",
		},
	}
	var loadedModules []string
	var moduleParams []string
	loader := &InitModuleLoader{
		Cmdline: cline,
		Prober: func(name, params string) error {
			loadedModules = append(loadedModules, name)
			moduleParams = append(moduleParams, params)
			return nil
		},
	}

	mods, err := GetModulesFromCmdline(loader)
	if err != nil {
		t.Fail()
	}
	InstallModules(loader, mods)
	expectedCmdLine := []string{"key1=value1", "key2=value2", "key3=value3"}
	expectedModules := []string{"test"}

	// Ordering of the parsed cmdline from the package isn't stable
	for _, val := range expectedCmdLine {
		if !strings.Contains(moduleParams[0], val) {
			t.Fatalf("failed cmdline test. Did not find %+v\n", val)
		}
	}

	if diff := cmp.Diff(expectedModules, loadedModules); diff != "" {
		t.Fatalf("unexpected difference of loaded modules (-want, +got): %v", diff)
	}
}

func TestOpenTTYDevices(t *testing.T) {
	tests := []struct {
		name      string
		ttyNames  []string
		numWriter int
		err       error
	}{
		{
			name:      "No TTY devices",
			ttyNames:  []string{},
			err:       nil,
			numWriter: 0,
		},
		{
			name:      "Single TTY device",
			ttyNames:  []string{"tty0"},
			err:       nil,
			numWriter: 1,
		},
		{
			name:      "Multiple TTY devices",
			ttyNames:  []string{"tty0", "ttyS0"},
			err:       nil,
			numWriter: 2,
		},
		{
			name:      "Non-existent TTY device",
			ttyNames:  []string{"nonexistent"},
			err:       os.ErrNotExist,
			numWriter: 0,
		},
		{
			name:      "Existent and non-existent TTY devices",
			ttyNames:  []string{"tty0", "nonexistent"},
			err:       nil,
			numWriter: 1,
		},
		{
			name:      "All TTY devices not existent",
			ttyNames:  []string{"nonexistent1", "nonexistent2"},
			err:       os.ErrNotExist,
			numWriter: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create dummy device files
			for _, tty := range tt.ttyNames {
				if !strings.HasPrefix(tty, "nonexistent") {
					if _, err := os.Create(filepath.Join(tmpDir, tty)); err != nil {
						t.Fatalf("create dummy device file %s: %v", tty, err)
					}
				}
			}

			writers, err := openTTYDevices(tmpDir, tt.ttyNames)
			if !errors.Is(err, tt.err) {
				t.Errorf("got error %v, want %v", err, tt.err)
			}
			if len(writers) != tt.numWriter {
				t.Errorf("got %d writers, want %d", len(writers), tt.numWriter)
			}
		})
	}
}
