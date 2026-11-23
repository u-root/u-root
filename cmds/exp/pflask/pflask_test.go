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

// fakeCgroupMkdir simulates the kernel populating cgroup.procs inside a newly
// created cgroup v2 directory.
func fakeCgroupMkdir(dir string, perm os.FileMode) error {
	if err := os.Mkdir(dir, perm); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "cgroup.procs"), nil, 0o644)
}

func makeFakeV2Root(t *testing.T, controllers, subtreeControl string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "cgroup.controllers"), []byte(controllers), 0o444); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "cgroup.subtree_control"), []byte(subtreeControl), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestSetupCgroup(t *testing.T) {
	t.Run("enables missing controllers, deduplicates, and attaches pid", func(t *testing.T) {
		root := makeFakeV2Root(t, "cpuset cpu io memory pids\n", "cpu\n")

		if dir, err := setupCgroup(root, "memory,cpu,pids,memory", 4242, fakeCgroupMkdir); err != nil || dir != "" {
			t.Fatalf("setupCgroup() = (%q, %v), want (\"\", nil)", dir, err)
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
		root := makeFakeV2Root(t, "cpu memory pids\n", "memory pids\n")
		scPath := filepath.Join(root, "cgroup.subtree_control")
		mkdir := func(dir string, perm os.FileMode) error {
			// Remove subtree_control after validation so any write
			// fails even as root (CAP_DAC_OVERRIDE).
			if err := os.Remove(scPath); err != nil {
				return err
			}
			return fakeCgroupMkdir(dir, perm)
		}

		if dir, err := setupCgroup(root, "memory,pids", 101, mkdir); err != nil || dir != "" {
			t.Fatalf("setupCgroup() = (%q, %v), want (\"\", nil)", dir, err)
		}
	})

	t.Run("dot attaches without touching subtree_control", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "cgroup.controllers"), nil, 0o444); err != nil {
			t.Fatal(err)
		}

		if dir, err := setupCgroup(root, ".", 99, fakeCgroupMkdir); err != nil || dir != "" {
			t.Fatalf("setupCgroup(.) = (%q, %v), want (\"\", nil)", dir, err)
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
		root := makeFakeV2Root(t, "memory\n", "memory\n")
		pre := filepath.Join(root, "pflask.50")
		if err := fakeCgroupMkdir(pre, 0o755); err != nil {
			t.Fatal(err)
		}

		if dir, err := setupCgroup(root, "memory", 50, os.Mkdir); err != nil || dir != "" {
			t.Fatalf("setupCgroup() = (%q, %v), want (\"\", nil)", dir, err)
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
		root := makeFakeV2Root(t, "cpu io memory pids\n", "")
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
			if _, err := setupCgroup(root, bad, 10, fakeCgroupMkdir); err == nil {
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

	t.Run("rejects empty, nonexistent, or non-v2 cgpath", func(t *testing.T) {
		// Chdir into a valid v2 root so empty cgpath cannot pass via
		// ENOENT.
		t.Chdir(makeFakeV2Root(t, "memory\n", "memory\n"))
		if _, err := setupCgroup("", "memory", 10, fakeCgroupMkdir); err == nil {
			t.Error(`setupCgroup(cgpath="") = nil, want error`)
		}
		missing := filepath.Join(t.TempDir(), "missing")
		if _, err := setupCgroup(missing, ".", 10, os.Mkdir); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("setupCgroup(missing, .) = %v, want ErrNotExist", err)
		}
		// Non-v2 dir (e.g. legacy v1 without cgroup.controllers).
		v1Root := t.TempDir()
		if err := os.Mkdir(filepath.Join(v1Root, "memory"), 0o755); err != nil {
			t.Fatal(err)
		}
		if _, err := setupCgroup(v1Root, "memory", 10, fakeCgroupMkdir); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("setupCgroup(v1Root, memory) = %v, want ErrNotExist", err)
		}
		if _, err := os.Stat(filepath.Join(v1Root, "memory", "pflask.10")); !errors.Is(err, os.ErrNotExist) {
			t.Error("non-v2 cgpath created memory/pflask.10 unexpectedly")
		}
	})

	t.Run("missing subtree_control fails before creating dir", func(t *testing.T) {
		noSC := t.TempDir()
		if err := os.WriteFile(filepath.Join(noSC, "cgroup.controllers"), []byte("memory\n"), 0o444); err != nil {
			t.Fatal(err)
		}
		if _, err := setupCgroup(noSC, "memory", 10, fakeCgroupMkdir); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("setupCgroup() = %v, want ErrNotExist", err)
		}
		if _, err := os.Stat(filepath.Join(noSC, "pflask.10")); !errors.Is(err, os.ErrNotExist) {
			t.Error("missing cgroup.subtree_control created pflask.10 unexpectedly")
		}
	})

	t.Run("does not create cgroup.procs on unmounted dir and returns cleanupDir", func(t *testing.T) {
		root := makeFakeV2Root(t, "memory\n", "")
		wantDir := filepath.Join(root, "pflask.88")
		// os.Mkdir leaves an empty directory without cgroup.procs.
		dir, err := setupCgroup(root, "memory", 88, os.Mkdir)
		if !errors.Is(err, os.ErrNotExist) || dir != wantDir {
			t.Fatalf("setupCgroup() = (%q, %v), want (%q, ErrNotExist)", dir, err, wantDir)
		}
		// os.Remove succeeds only if writeCgroupFile did not create
		// cgroup.procs inside the unmounted directory.
		if err := os.Remove(dir); err != nil {
			t.Errorf("os.Remove(%q) = %v, want nil (dir should be empty)", dir, err)
		}
	})

	t.Run("preserves preexisting file and empty dir on attach failure", func(t *testing.T) {
		root := makeFakeV2Root(t, "", "")

		existingFile := filepath.Join(root, "pflask.77")
		if err := os.WriteFile(existingFile, []byte("keep"), 0o644); err != nil {
			t.Fatal(err)
		}
		if dir, err := setupCgroup(root, ".", 77, os.Mkdir); err == nil || dir != "" {
			t.Fatalf("setupCgroup() = (%q, %v), want (\"\", non-nil)", dir, err)
		}
		if _, err := os.Stat(existingFile); err != nil {
			t.Errorf("preexisting file removed on error: %v", err)
		}

		existingDir := filepath.Join(root, "pflask.78")
		if err := os.Mkdir(existingDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if dir, err := setupCgroup(root, ".", 78, os.Mkdir); err == nil || dir != "" {
			t.Fatalf("setupCgroup() = (%q, %v), want (\"\", non-nil)", dir, err)
		}
		if _, err := os.Stat(existingDir); err != nil {
			t.Errorf("preexisting dir removed on error: %v", err)
		}
	})
}

func TestSetupCgroupSubtreeControlFailure(t *testing.T) {
	for _, tc := range []struct {
		name        string
		preexisting bool
	}{
		{name: "newly created dir", preexisting: false},
		{name: "preexisting dir", preexisting: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := makeFakeV2Root(t, "cpu memory\n", "cpu\n")
			targetDir := filepath.Join(root, "pflask.4242")
			if tc.preexisting {
				if err := fakeCgroupMkdir(targetDir, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			mkdir := func(dir string, perm os.FileMode) error {
				err := fakeCgroupMkdir(dir, perm)
				// Remove subtree_control after validation so
				// the post-attach write fails cleanly.
				if rmErr := os.Remove(filepath.Join(root, "cgroup.subtree_control")); rmErr != nil {
					t.Fatal(rmErr)
				}
				return err
			}

			cleanupDir, err := setupCgroup(root, "memory", 4242, mkdir)
			if !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("setupCgroup() = (%q, %v), want ErrNotExist", cleanupDir, err)
			}
			wantCleanup := targetDir
			if tc.preexisting {
				wantCleanup = ""
			}
			if cleanupDir != wantCleanup {
				t.Errorf("cleanupDir = %q, want %q", cleanupDir, wantCleanup)
			}
			procs, err := os.ReadFile(filepath.Join(targetDir, "cgroup.procs"))
			if err != nil {
				t.Fatalf("read cgroup.procs: %v", err)
			}
			if string(procs) != "4242" {
				t.Errorf("cgroup.procs = %q, want %q", procs, "4242")
			}
		})
	}
}
