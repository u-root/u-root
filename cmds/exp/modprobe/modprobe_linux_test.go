// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestModules(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	kver := "6.6.0-test"
	modDir := filepath.Join(root, "lib/modules", kver)
	etcModprobe := filepath.Join(root, "etc/modprobe.d")
	procDir := filepath.Join(root, "proc")

	for _, d := range []string{modDir, etcModprobe, procDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	modulesDep := `kernel/net/9p/9pnet.ko:
kernel/fs/9p/9p.ko: kernel/net/9p/9pnet.ko
kernel/net/9p/9pnet_fd.ko: kernel/net/9p/9pnet.ko
kernel/drivers/char/blacklisted_mod.ko:
`
	if err := os.WriteFile(filepath.Join(modDir, "modules.dep"), []byte(modulesDep), 0o644); err != nil {
		t.Fatal(err)
	}

	modulesAlias := `alias fs-9p 9p
alias 9p-fd 9pnet_fd
`
	if err := os.WriteFile(filepath.Join(modDir, "modules.alias"), []byte(modulesAlias), 0o644); err != nil {
		t.Fatal(err)
	}

	conf := `blacklist blacklisted_mod
`
	if err := os.WriteFile(filepath.Join(etcModprobe, "blacklist.conf"), []byte(conf), 0o644); err != nil {
		t.Fatal(err)
	}

	// Empty /proc/modules so host modules do not interfere.
	if err := os.WriteFile(filepath.Join(procDir, "modules"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	return root, kver
}

func TestModprobeCLI(t *testing.T) {
	root, kver := setupTestModules(t)

	t.Run("kernel call_modprobe invocation (-q -- fs-9p)", func(t *testing.T) {
		var stderr bytes.Buffer
		err := run([]string{"-n", "-d", root, "-S", kver, "-q", "--", "fs-9p"}, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if stderr.Len() != 0 {
			t.Errorf("expected empty stderr with -q, got %q", stderr.String())
		}

		err = run([]string{"-n", "-d", root, "-S", kver, "--", "fs-9p"}, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		out := stderr.String()
		if !strings.Contains(out, "9pnet.ko") || !strings.Contains(out, "9p.ko") {
			t.Errorf("expected 9pnet.ko and 9p.ko in dry-run output, got:\n%s", out)
		}
	})

	t.Run("combined flags -va and deduplicated prober", func(t *testing.T) {
		var stderr bytes.Buffer
		err := run([]string{"-n", "-va", "-d", root, "-S", kver, "fs-9p", "9p-fd"}, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		out := stderr.String()
		if got := strings.Count(out, "9pnet.ko"); got != 1 {
			t.Errorf("expected 9pnet.ko to be loaded once across -a, got %d times in:\n%s", got, out)
		}
		if !strings.Contains(out, "9p.ko") || !strings.Contains(out, "9pnet_fd.ko") {
			t.Errorf("expected both 9p.ko and 9pnet_fd.ko in output, got:\n%s", out)
		}
	})

	t.Run("combined flags after -d and -S", func(t *testing.T) {
		var stderr bytes.Buffer
		err := run([]string{"-d", root, "-S", kver, "-nva", "fs-9p", "9p-fd"}, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		out := stderr.String()
		if got := strings.Count(out, "9pnet.ko"); got != 1 {
			t.Errorf("expected 9pnet.ko to be loaded once across -a, got %d times in:\n%s", got, out)
		}
	})

	t.Run("-a returns error when module missing", func(t *testing.T) {
		var stderr bytes.Buffer
		err := run([]string{"-n", "-a", "-d", root, "-S", kver, "fs-9p", "nonexistent_mod"}, &stderr)
		if err == nil {
			t.Fatalf("expected error for nonexistent_mod with -a")
		}
	})

	t.Run("-q suppresses error output", func(t *testing.T) {
		var stderr bytes.Buffer
		err := run([]string{"-q", "-d", root, "-S", kver, "--", "nonexistent_mod"}, &stderr)
		if err == nil {
			t.Fatalf("expected error for nonexistent_mod")
		}
		if stderr.Len() != 0 {
			t.Errorf("expected empty stderr with -q, got %q", stderr.String())
		}
	})

	t.Run("-b honors blacklist on direct module name", func(t *testing.T) {
		var stderr bytes.Buffer
		err := run([]string{"-n", "-b", "-d", root, "-S", kver, "blacklisted_mod"}, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if strings.Contains(stderr.String(), "blacklisted_mod.ko") {
			t.Errorf("expected blacklisted_mod.ko to be skipped with -b, got:\n%s", stderr.String())
		}
	})

	t.Run("short flag with equals (-nd=root)", func(t *testing.T) {
		var stderr bytes.Buffer
		err := run([]string{"-vd=" + root, "-S=" + kver, "-n", "fs-9p"}, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(stderr.String(), "9p.ko") {
			t.Errorf("expected 9p.ko in output, got:\n%s", stderr.String())
		}
	})

	t.Run("-r removes module in dry-run and fails on missing", func(t *testing.T) {
		var stderr bytes.Buffer
		err := run([]string{"-n", "-r", "-d", root, "-S", kver, "fs-9p"}, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(stderr.String(), "9p.ko") {
			t.Errorf("expected 9p.ko in remove dry-run output, got:\n%s", stderr.String())
		}

		stderr.Reset()
		if err := run([]string{"-n", "-r", "-d", root, "-S", kver, "nonexistent_mod"}, &stderr); err == nil {
			t.Fatalf("expected error removing nonexistent_mod")
		}
	})

	t.Run("missing modules.dep directory returns error", func(t *testing.T) {
		var stderr bytes.Buffer
		if err := run([]string{"-d", root, "-S", "missing-kver", "fs-9p"}, &stderr); err == nil {
			t.Fatalf("expected error for missing kernel version dir")
		}
	})

	t.Run("invalid flag returns error", func(t *testing.T) {
		var stderr bytes.Buffer
		if err := run([]string{"--no-such-flag", "fs-9p"}, &stderr); err == nil {
			t.Fatalf("expected error for unknown flag")
		}
	})

	t.Run("no args returns error", func(t *testing.T) {
		var stderr bytes.Buffer
		if err := run([]string{}, &stderr); err == nil {
			t.Fatalf("expected error when no args given")
		}
	})
}
