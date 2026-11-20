// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// fakeCgroupMkdir simulates the kernel populating controlFile inside a newly
// created cgroup directory.
func fakeCgroupMkdir(controlFile string) func(string, os.FileMode) error {
	return func(dir string, perm os.FileMode) error {
		if err := os.Mkdir(dir, perm); err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(dir, controlFile), nil, 0o644)
	}
}

func makeFakeV2Root(t *testing.T, controllers, subtreeControl string, scMode os.FileMode) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "cgroup.controllers"), []byte(controllers), 0o444); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "cgroup.subtree_control"), []byte(subtreeControl), scMode); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestSetupCgroupV2(t *testing.T) {
	t.Run("enables missing controllers, deduplicates, and attaches pid", func(t *testing.T) {
		root := makeFakeV2Root(t, "cpuset cpu io memory pids\n", "cpu\n", 0o644)

		if err := setupCgroup(root, "memory,cpu,pids,memory", 4242, fakeCgroupMkdir("cgroup.procs")); err != nil {
			t.Fatalf("setupCgroup() = %v, want nil", err)
		}

		sc, err := os.ReadFile(filepath.Join(root, "cgroup.subtree_control"))
		if err != nil {
			t.Fatalf("read cgroup.subtree_control: %v", err)
		}
		if string(sc) != "+memory +pids" {
			t.Errorf("cgroup.subtree_control = %q, want %q", sc, "+memory +pids")
		}

		procs, err := os.ReadFile(filepath.Join(root, "pflask.4242", "cgroup.procs"))
		if err != nil {
			t.Fatalf("read cgroup.procs: %v", err)
		}
		if string(procs) != "4242" {
			t.Errorf("cgroup.procs = %q, want %q", procs, "4242")
		}
	})

	t.Run("skips subtree_control write when already enabled", func(t *testing.T) {
		root := makeFakeV2Root(t, "cpu memory pids\n", "memory pids\n", 0o444)

		if err := setupCgroup(root, "memory,pids", 101, fakeCgroupMkdir("cgroup.procs")); err != nil {
			t.Fatalf("setupCgroup() = %v, want nil", err)
		}
		sc, err := os.ReadFile(filepath.Join(root, "cgroup.subtree_control"))
		if err != nil {
			t.Fatal(err)
		}
		if string(sc) != "memory pids\n" {
			t.Errorf("cgroup.subtree_control = %q, want unchanged %q", sc, "memory pids\n")
		}
	})

	t.Run("dot attaches without touching subtree_control", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "cgroup.controllers"), nil, 0o444); err != nil {
			t.Fatal(err)
		}

		if err := setupCgroup(root, ".", 99, fakeCgroupMkdir("cgroup.procs")); err != nil {
			t.Fatalf("setupCgroup(.) = %v, want nil", err)
		}
		if _, err := os.Stat(filepath.Join(root, "cgroup.subtree_control")); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("cgroup.subtree_control stat err = %v, want ErrNotExist", err)
		}
		procs, err := os.ReadFile(filepath.Join(root, "pflask.99", "cgroup.procs"))
		if err != nil {
			t.Fatalf("read cgroup.procs: %v", err)
		}
		if string(procs) != "99" {
			t.Errorf("cgroup.procs = %q, want %q", procs, "99")
		}
	})

	t.Run("attaches into preexisting cgroup dir", func(t *testing.T) {
		root := makeFakeV2Root(t, "memory\n", "memory\n", 0o444)
		pre := filepath.Join(root, "pflask.50")
		if err := os.Mkdir(pre, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(pre, "cgroup.procs"), nil, 0o644); err != nil {
			t.Fatal(err)
		}

		if err := setupCgroup(root, "memory", 50, os.Mkdir); err != nil {
			t.Fatalf("setupCgroup() = %v, want nil", err)
		}
		procs, err := os.ReadFile(filepath.Join(pre, "cgroup.procs"))
		if err != nil {
			t.Fatal(err)
		}
		if string(procs) != "50" {
			t.Errorf("cgroup.procs = %q, want %q", procs, "50")
		}
	})

	t.Run("rejects unavailable or malformed controllers without side effects", func(t *testing.T) {
		root := makeFakeV2Root(t, "cpu io memory pids\n", "", 0o644)
		for _, bad := range []string{
			"",
			".,.",
			"memory,.",
			"blkio",
			"memory,",
			"-memory",
			"+memory",
			"memory -cpu",
			"../escape",
		} {
			if err := setupCgroup(root, bad, 10, fakeCgroupMkdir("cgroup.procs")); err == nil {
				t.Errorf("setupCgroup(%q) = nil, want error", bad)
			}
		}
		sc, err := os.ReadFile(filepath.Join(root, "cgroup.subtree_control"))
		if err != nil {
			t.Fatal(err)
		}
		if string(sc) != "" {
			t.Errorf("cgroup.subtree_control = %q, want empty", sc)
		}
		if _, err := os.Stat(filepath.Join(root, "pflask.10")); !errors.Is(err, os.ErrNotExist) {
			t.Error("invalid input created pflask.10 unexpectedly")
		}
	})

	t.Run("rejects empty, whitespace, or nonexistent cgpath", func(t *testing.T) {
		for _, bad := range []string{"", "   "} {
			if err := setupCgroup(bad, "memory", 10, os.Mkdir); err == nil {
				t.Errorf("setupCgroup(cgpath=%q) = nil, want error", bad)
			}
		}
		missing := filepath.Join(t.TempDir(), "missing")
		if err := setupCgroup(missing, ".", 10, os.Mkdir); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("setupCgroup(missing, .) = %v, want ErrNotExist", err)
		}
	})

	t.Run("unwritable subtree_control fails before creating cgroup dir", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("chmod 0o444 does not block root writes")
		}
		root := makeFakeV2Root(t, "memory\n", "", 0o444)
		if err := setupCgroup(root, "memory", 12, fakeCgroupMkdir("cgroup.procs")); !errors.Is(err, os.ErrPermission) {
			t.Fatalf("setupCgroup() = %v, want ErrPermission", err)
		}
		if _, err := os.Stat(filepath.Join(root, "pflask.12")); !errors.Is(err, os.ErrNotExist) {
			t.Error("unwritable cgroup.subtree_control created pflask.12 unexpectedly")
		}
	})

	t.Run("propagates non-ENOENT cgroup.controllers and missing subtree_control", func(t *testing.T) {
		dirCtrl := t.TempDir()
		if err := os.Mkdir(filepath.Join(dirCtrl, "cgroup.controllers"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := setupCgroup(dirCtrl, "memory", 10, os.Mkdir); err == nil {
			t.Error("setupCgroup() = nil when cgroup.controllers is a directory, want error")
		}

		noSC := t.TempDir()
		if err := os.WriteFile(filepath.Join(noSC, "cgroup.controllers"), []byte("memory\n"), 0o444); err != nil {
			t.Fatal(err)
		}
		if err := setupCgroup(noSC, "memory", 10, fakeCgroupMkdir("cgroup.procs")); err == nil {
			t.Error("setupCgroup() = nil when cgroup.subtree_control is missing, want error")
		}
		if _, err := os.Stat(filepath.Join(noSC, "pflask.10")); !errors.Is(err, os.ErrNotExist) {
			t.Error("missing cgroup.subtree_control created pflask.10 unexpectedly")
		}
	})

	t.Run("does not create cgroup.procs on unmounted dir and cleans up", func(t *testing.T) {
		root := makeFakeV2Root(t, "memory\n", "", 0o644)
		// os.Mkdir leaves an empty directory without cgroup.procs.
		if err := setupCgroup(root, "memory", 88, os.Mkdir); err == nil {
			t.Fatal("setupCgroup() = nil on unmounted dir, want error")
		}
		if _, err := os.Stat(filepath.Join(root, "pflask.88")); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("pflask.88 not cleaned up after missing cgroup.procs: %v", err)
		}
	})

	t.Run("preserves preexisting file and empty dir on attach failure", func(t *testing.T) {
		root := makeFakeV2Root(t, "", "", 0o444)

		existingFile := filepath.Join(root, "pflask.77")
		if err := os.WriteFile(existingFile, []byte("keep"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := setupCgroup(root, ".", 77, os.Mkdir); err == nil {
			t.Fatal("setupCgroup() = nil when target is a regular file, want error")
		}
		if _, err := os.Stat(existingFile); err != nil {
			t.Errorf("preexisting file removed on error: %v", err)
		}

		existingDir := filepath.Join(root, "pflask.78")
		if err := os.Mkdir(existingDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := setupCgroup(root, ".", 78, os.Mkdir); err == nil {
			t.Fatal("setupCgroup() = nil when preexisting dir lacks cgroup.procs, want error")
		}
		if _, err := os.Stat(existingDir); err != nil {
			t.Errorf("preexisting dir removed on error: %v", err)
		}
	})
}

func TestSetupCgroupV1(t *testing.T) {
	makeCtrl := func(t *testing.T, root, name string) {
		t.Helper()
		dir := filepath.Join(root, name)
		if err := os.Mkdir(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "tasks"), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("attaches pid to each v1 controller including co-mounted symlinks", func(t *testing.T) {
		root := t.TempDir()
		makeCtrl(t, root, "memory")
		makeCtrl(t, root, "cpu,cpuacct")
		for _, link := range []string{"cpu", "cpuacct"} {
			if err := os.Symlink(filepath.Join(root, "cpu,cpuacct"), filepath.Join(root, link)); err != nil {
				t.Fatal(err)
			}
		}

		if err := setupCgroup(root, "memory,cpu,cpuacct", 314, fakeCgroupMkdir("tasks")); err != nil {
			t.Fatalf("setupCgroup() = %v, want nil", err)
		}
		for _, ctrl := range []string{"memory", "cpu,cpuacct"} {
			got, err := os.ReadFile(filepath.Join(root, ctrl, "pflask.314", "tasks"))
			if err != nil {
				t.Fatalf("read %s tasks: %v", ctrl, err)
			}
			if string(got) != "314" {
				t.Errorf("%s tasks = %q, want %q", ctrl, got, "314")
			}
		}
	})

	t.Run("rejects invalid or unmounted controllers without partial attach", func(t *testing.T) {
		root := t.TempDir()
		makeCtrl(t, root, "memory")
		// Empty dir without tasks simulates an unmounted controller
		// dir on a v1 tmpfs (e.g. CONFIG_CPUSETS_V1=n).
		if err := os.Mkdir(filepath.Join(root, "cpuset"), 0o555); err != nil {
			t.Fatal(err)
		}
		// Directory named tasks is not a mounted cgroup control file.
		if err := os.MkdirAll(filepath.Join(root, "badtasks", "tasks"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "notadir"), nil, 0o644); err != nil {
			t.Fatal(err)
		}

		for _, bad := range []string{
			"",
			"nonexistent",
			"cpuset",
			"badtasks",
			"notadir",
			"memory,nonexistent",
			"memory,cpuset",
			"memory,",
			".",
			"./.",
			"..",
			"../escape",
			"/etc",
		} {
			if err := setupCgroup(root, bad, 55, fakeCgroupMkdir("tasks")); err == nil {
				t.Errorf("setupCgroup(%q) = nil, want error", bad)
			}
		}
		if _, err := os.Stat(filepath.Join(root, "memory", "pflask.55")); !errors.Is(err, os.ErrNotExist) {
			t.Error("failed call partially attached memory/pflask.55")
		}
		if _, err := os.Stat(filepath.Join(root, "nonexistent")); !errors.Is(err, os.ErrNotExist) {
			t.Error("setupCgroup(nonexistent) created fake controller dir")
		}
	})
}
