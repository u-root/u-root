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
	"golang.org/x/sys/unix"
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

func TestSelectCgroupMode(t *testing.T) {
	const (
		bothFS = "nodev\tsysfs\nnodev\tcgroup\nnodev\tcgroup2\n\text4\n"
		v2Only = "nodev\tsysfs\nnodev\tcgroup2\n\text4\n"
		v1Only = "nodev\tsysfs\nnodev\tcgroup\n\text4\n"
		noCgFS = "nodev\tsysfs\n\text4\n\tcgroup2\nnodev\tcgroup2\textra\n"
	)

	tests := []struct {
		name      string
		initFlags string
		procFS    string
		want      cgroupMode
	}{
		{
			name:   "auto prefers v2 when both present",
			procFS: bothFS,
			want:   cgroupV2,
		},
		{
			name:   "auto selects v2 when only cgroup2 present",
			procFS: v2Only,
			want:   cgroupV2,
		},
		{
			name:   "auto falls back to v1 on pre-v2 kernel",
			procFS: v1Only,
			want:   cgroupV1,
		},
		{
			name:   "auto skips when cgroups disabled in kernel",
			procFS: noCgFS,
			want:   cgroupNone,
		},
		{
			name: "auto tries v2 when /proc/filesystems unreadable",
			want: cgroupV2,
		},
		{
			name:      "explicit cgroups=v2",
			initFlags: "cgroups=v2",
			procFS:    v1Only,
			want:      cgroupV2,
		},
		{
			name:      "explicit cgroups=v1",
			initFlags: "cgroups=v1",
			procFS:    bothFS,
			want:      cgroupV1,
		},
		{
			name:      "explicit cgroups=none",
			initFlags: "cgroups=none",
			procFS:    bothFS,
			want:      cgroupNone,
		},
		{
			name:      "systemd skips cgroups",
			initFlags: "systemd",
			procFS:    bothFS,
			want:      cgroupNone,
		},
		{
			name:      "systemd takes precedence over cgroups=v2",
			initFlags: "systemd cgroups=v2",
			procFS:    bothFS,
			want:      cgroupNone,
		},
		{
			name:      "systemd=false honors cgroups=v1",
			initFlags: "systemd=false cgroups=v1",
			procFS:    bothFS,
			want:      cgroupV1,
		},
		{
			name:      "unknown cgroups value falls back to auto",
			initFlags: "cgroups=bogus",
			procFS:    v1Only,
			want:      cgroupV1,
		},
		{
			name:      "bare cgroups flag falls back to auto",
			initFlags: "cgroups",
			procFS:    noCgFS,
			want:      cgroupNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			procFSPath := filepath.Join(t.TempDir(), "filesystems")
			if tt.procFS != "" {
				if err := os.WriteFile(procFSPath, []byte(tt.procFS), 0o444); err != nil {
					t.Fatal(err)
				}
			}
			cl := &cmdline.CmdLine{
				AsMap: map[string]string{"uroot.initflags": tt.initFlags},
			}
			if got := selectCgroupMode(cl.GetInitFlagMap(), procFSPath); got != tt.want {
				t.Errorf("selectCgroupMode(%q) = %v, want %v", tt.initFlags, got, tt.want)
			}
		})
	}
}

func TestMountCgroup2(t *testing.T) {
	tests := []struct {
		name     string
		failWith map[string]error
		wantOpts []string
		wantErr  error
	}{
		{
			name:     "modern kernel accepts full options",
			wantOpts: []string{"nsdelegate,memory_recursiveprot"},
		},
		{
			name: "falls back to nsdelegate on EINVAL",
			failWith: map[string]error{
				"nsdelegate,memory_recursiveprot": unix.EINVAL,
			},
			wantOpts: []string{"nsdelegate,memory_recursiveprot", "nsdelegate"},
		},
		{
			name: "falls back to empty opts on EINVAL",
			failWith: map[string]error{
				"nsdelegate,memory_recursiveprot": unix.EINVAL,
				"nsdelegate":                      unix.EINVAL,
			},
			wantOpts: []string{"nsdelegate,memory_recursiveprot", "nsdelegate", ""},
		},
		{
			name: "returns EINVAL when all attempts return EINVAL",
			failWith: map[string]error{
				"nsdelegate,memory_recursiveprot": unix.EINVAL,
				"nsdelegate":                      unix.EINVAL,
				"":                                unix.EINVAL,
			},
			wantOpts: []string{"nsdelegate,memory_recursiveprot", "nsdelegate", ""},
			wantErr:  unix.EINVAL,
		},
		{
			name: "does not retry non-EINVAL error",
			failWith: map[string]error{
				"nsdelegate,memory_recursiveprot": unix.EPERM,
			},
			wantOpts: []string{"nsdelegate,memory_recursiveprot"},
			wantErr:  unix.EPERM,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotOpts []string
			err := mountCgroup2("/sys/fs/cgroup", func(source, target, fstype string, flags uintptr, data string) error {
				if source != "cgroup2" || target != "/sys/fs/cgroup" || fstype != "cgroup2" {
					t.Errorf("mount(%q, %q, %q), want (cgroup2, /sys/fs/cgroup, cgroup2)", source, target, fstype)
				}
				if flags != cgroup2MountFlags {
					t.Errorf("flags = %#x, want %#x", flags, cgroup2MountFlags)
				}
				gotOpts = append(gotOpts, data)
				return tt.failWith[data]
			})
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.wantOpts, gotOpts); diff != "" {
				t.Errorf("mount opts mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIsMountPoint(t *testing.T) {
	if isMountPoint(filepath.Join(t.TempDir(), "nonexistent")) {
		t.Error("isMountPoint(nonexistent) = true, want false")
	}
	if isMountPoint(t.TempDir()) {
		t.Error("isMountPoint(t.TempDir()) = true, want false")
	}

	var st unix.Statfs_t
	if err := unix.Statfs("/proc", &st); err == nil && int64(st.Type) == unix.PROC_SUPER_MAGIC {
		for _, p := range []string{"/proc", "/proc/"} {
			if !isMountPoint(p) {
				t.Errorf("isMountPoint(%q) = false, want true", p)
			}
		}
	}
}

func TestSetupCgroups(t *testing.T) {
	unmounted := t.TempDir()

	tests := []struct {
		name      string
		mode      cgroupMode
		target    string
		mountErr  error
		skipMount bool
		wantV2    int
		wantV1    int
	}{
		{
			name:   "v2 mounts on existing unmounted dir",
			mode:   cgroupV2,
			target: unmounted,
			wantV2: 1,
		},
		{
			name:     "v2 mount failure is optional",
			mode:     cgroupV2,
			target:   unmounted,
			mountErr: unix.EPERM,
			wantV2:   1,
		},
		{
			name:   "v1 calls mountV1",
			mode:   cgroupV1,
			target: unmounted,
			wantV1: 1,
		},
		{
			name:   "none skips all",
			mode:   cgroupNone,
			target: unmounted,
		},
		{
			name:      "already mounted skips v2",
			mode:      cgroupV2,
			target:    cgroupMountPoint,
			skipMount: true,
		},
		{
			name:      "already mounted skips v1",
			mode:      cgroupV1,
			target:    cgroupMountPoint,
			skipMount: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipMount && !isMountPoint(tt.target) {
				t.Skipf("%s is not a mount point", tt.target)
			}
			var v2Calls, v1Calls int
			setupCgroups(tt.mode, tt.target, func(_, target, fstype string, _ uintptr, _ string) error {
				v2Calls++
				if target != tt.target || fstype != "cgroup2" {
					t.Errorf("unexpected mount(%q, %q)", target, fstype)
				}
				return tt.mountErr
			}, func() { v1Calls++ })
			if v2Calls != tt.wantV2 || v1Calls != tt.wantV1 {
				t.Errorf("v2Calls=%d v1Calls=%d, want %d, %d", v2Calls, v1Calls, tt.wantV2, tt.wantV1)
			}
		})
	}
}
