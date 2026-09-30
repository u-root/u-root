// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kmodule

import (
	"bytes"
	"compress/gzip"
	"errors"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sys/unix"
)

var procModsMock = `hid_generic 16384 0 - Live 0x0000000000000000
usbhid 49152 0 - Live 0x0000000000000000
ccm 20480 6 - Live 0x0000000000000000
`

func TestGenLoadedMods(t *testing.T) {
	m := depMap{
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/hid-generic.ko":   &dependency{},
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/usbhid/usbhid.ko": &dependency{},
		"/lib/modules/6.6.6-generic/kernel/crypto/ccm.ko":                &dependency{},
	}
	br := bytes.NewBufferString(procModsMock)
	err := genLoadedMods(br, m)
	if err != nil {
		t.Fatalf("fail to genLoadedMods: %v\n", err)
	}
	for mod, d := range m {
		if d.state != loaded {
			t.Fatalf("mod %q should have been loaded", path.Base(mod))
		}
	}
}

func TestGenLoadedModsStatesAndOutOfTree(t *testing.T) {
	mock := `out_of_tree_mod 16384 0 - Live 0x0000000000000000
twofish_x86_64 32768 1 - Live 0x0000000000000000
still_loading 16384 0 - Loading 0x0000000000000000
being_unloaded 16384 0 - Unloading 0x0000000000000000
nls_iso8859_1 16384 0 - Live 0x0000000000000000
`
	m := depMap{
		"/lib/modules/6.6/kernel/crypto/twofish-x86_64.ko":  &dependency{},
		"/lib/modules/6.6/kernel/fs/nls/nls_iso8859-1.ko":   &dependency{},
		"/lib/modules/6.6/kernel/drivers/still_loading.ko":  &dependency{},
		"/lib/modules/6.6/kernel/drivers/being-unloaded.ko": &dependency{},
	}
	if err := genLoadedMods(bytes.NewBufferString(mock), m); err != nil {
		t.Fatalf("genLoadedMods failed: %v", err)
	}
	if m["/lib/modules/6.6/kernel/crypto/twofish-x86_64.ko"].state != loaded {
		t.Errorf("expected twofish-x86_64.ko to be marked loaded")
	}
	if m["/lib/modules/6.6/kernel/fs/nls/nls_iso8859-1.ko"].state != loaded {
		t.Errorf("expected nls_iso8859-1.ko to be marked loaded")
	}
	if m["/lib/modules/6.6/kernel/drivers/still_loading.ko"].state != loaded {
		t.Errorf("expected Loading module to be marked loaded so recursive request_module does not deadlock")
	}
	if m["/lib/modules/6.6/kernel/drivers/being-unloaded.ko"].state != unloaded {
		t.Errorf("expected Unloading module to remain unloaded")
	}
}

func TestFindModPathMixedHyphenUnderscore(t *testing.T) {
	m := depMap{
		"/lib/modules/6.6/kernel/crypto/twofish-x86_64.ko.xz": &dependency{},
		"/lib/modules/6.6/kernel/fs/nls/nls_iso8859-1.ko.zst": &dependency{},
		"/lib/modules/6.6/kernel/net/9p/9pnet_fd.ko.gz":       &dependency{},
	}
	for _, tc := range []struct {
		query string
		want  string
	}{
		{"twofish_x86_64", "/lib/modules/6.6/kernel/crypto/twofish-x86_64.ko.xz"},
		{"twofish-x86-64", "/lib/modules/6.6/kernel/crypto/twofish-x86_64.ko.xz"},
		{"twofish-x86_64", "/lib/modules/6.6/kernel/crypto/twofish-x86_64.ko.xz"},
		{"nls_iso8859_1", "/lib/modules/6.6/kernel/fs/nls/nls_iso8859-1.ko.zst"},
		{"nls-iso8859-1", "/lib/modules/6.6/kernel/fs/nls/nls_iso8859-1.ko.zst"},
		{"9pnet-fd", "/lib/modules/6.6/kernel/net/9p/9pnet_fd.ko.gz"},
	} {
		got, err := findModPath(tc.query, m)
		if err != nil {
			t.Errorf("findModPath(%q) unexpected error: %v", tc.query, err)
			continue
		}
		if got != tc.want {
			t.Errorf("findModPath(%q) = %q, want %q", tc.query, got, tc.want)
		}
	}
}

func TestMatchAlias(t *testing.T) {
	for _, tc := range []struct {
		pattern string
		name    string
		want    bool
	}{
		{"fs-9p", "fs-9p", true},
		{"fs-9p", "fs_9p", true},
		{"9p-fd", "9p-fd", true},
		{"block-major-7-*", "block-major-7-0", true},
		{"char-major-10-[0-9]*", "char-major-10-54", true},
		{"char-major-10-[0-9]*", "char-major-10-abc", false},
		{"pci:v00001AF4d00001009sv*sd*bc*sc*i*", "pci:v00001AF4d00001009sv00001AF4sd00000009bc00sc00i00", true},
		{"dmi:*:svnGoogle:pn*", "dmi:bvnSeaBIOS:bvr1.16:bd08/05/2026:svnGoogle:pnComputeEngine:", true},
	} {
		if got := matchAlias(tc.pattern, tc.name); got != tc.want {
			t.Errorf("matchAlias(%q, %q) = %v, want %v", tc.pattern, tc.name, got, tc.want)
		}
	}
}

func TestModprobeConfigsAliasesAndCmdline(t *testing.T) {
	root := t.TempDir()
	kver := "6.6.0-test"
	modDir := filepath.Join(root, "lib/modules", kver)
	etcModprobe := filepath.Join(root, "etc/modprobe.d")
	libModprobe := filepath.Join(root, "lib/modprobe.d")
	procDir := filepath.Join(root, "proc")

	for _, d := range []string{modDir, etcModprobe, libModprobe, procDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	// modules.dep with comments, blank lines, and trailing spaces after
	// ':'.
	modulesDep := `# comment line

kernel/net/9p/9pnet.ko:   
kernel/fs/9p/9p.ko: kernel/net/9p/9pnet.ko
kernel/net/9p/9pnet-fd.ko: kernel/net/9p/9pnet.ko
kernel/drivers/char/bad_mod.ko:
kernel/drivers/misc/override_target.ko:
`
	if err := os.WriteFile(filepath.Join(modDir, "modules.dep"), []byte(modulesDep), 0o644); err != nil {
		t.Fatal(err)
	}

	modulesBuiltin := `kernel/fs/ext4/ext4.ko
`
	if err := os.WriteFile(filepath.Join(modDir, "modules.builtin"), []byte(modulesBuiltin), 0o644); err != nil {
		t.Fatal(err)
	}

	modulesAlias := `# Aliases extracted from modules themselves.
alias fs-9p 9p
alias 9p-fd 9pnet_fd
alias pci:v00001AF4d00001009sv*sd*bc*sc*i* 9p
alias bad-alias bad_mod
alias hw-chained middle_target
`
	if err := os.WriteFile(filepath.Join(modDir, "modules.alias"), []byte(modulesAlias), 0o644); err != nil {
		t.Fatal(err)
	}

	builtinModinfo := "ext4.alias=fs-ext4\x00ext4.license=GPL\x00builtin_only.alias=fs-builtin-only\x00"
	if err := os.WriteFile(filepath.Join(modDir, "modules.builtin.modinfo"), []byte(builtinModinfo), 0o644); err != nil {
		t.Fatal(err)
	}

	// /lib/modprobe.d/override.conf should be shadowed by
	// /etc/modprobe.d/override.conf.
	libOverride := `options 9p from_lib=1
`
	if err := os.WriteFile(filepath.Join(libModprobe, "override.conf"), []byte(libOverride), 0o644); err != nil {
		t.Fatal(err)
	}

	etcOverride := `# Line continuation, quoted #, and options
options 9p \
	cache=loose tag="a#b" # trailing comment
alias custom-alias override_target
alias middle_target override_target
`
	if err := os.WriteFile(filepath.Join(etcModprobe, "override.conf"), []byte(etcOverride), 0o644); err != nil {
		t.Fatal(err)
	}

	cmdline := `/boot/vmlinuz-6.6.0 rd.luks.uuid=123 console=ttyS0 ip=192.168.1.1 9pnet.debug=1 9p.trans="fd" modprobe.blacklist=bad-mod -- ignored.param=1` + "\n"
	if err := os.WriteFile(filepath.Join(procDir, "cmdline"), []byte(cmdline), 0o644); err != nil {
		t.Fatal(err)
	}

	prober, err := NewProber(ProbeOpts{
		RootDir:        root,
		KVer:           kver,
		IgnoreProcMods: true,
	})
	if err != nil {
		t.Fatalf("NewProber failed: %v", err)
	}

	var mu sync.Mutex
	loadedParams := make(map[string]string)
	var loadOrder []string
	prober.loadCB = func(modPath, params string) error {
		mu.Lock()
		defer mu.Unlock()
		base := filepath.Base(modPath)
		loadedParams[base] = params
		loadOrder = append(loadOrder, base)
		return nil
	}

	// 1. Probe fs-9p (kernel autoloading alias for 9p filesystem) with
	// extra user param.
	if err := prober.Probe("fs-9p", "msize=65536"); err != nil {
		t.Fatalf("Probe(fs-9p) failed: %v", err)
	}
	wantOrder := []string{"9pnet.ko", "9p.ko"}
	if !slices.Equal(loadOrder, wantOrder) {
		t.Errorf("loadOrder = %v, want %v", loadOrder, wantOrder)
	}
	if got, want := loadedParams["9pnet.ko"], "debug=1"; got != want {
		t.Errorf("9pnet.ko params = %q, want %q", got, want)
	}
	if got, want := loadedParams["9p.ko"], `cache=loose tag="a#b" trans="fd" msize=65536`; got != want {
		t.Errorf("9p.ko params = %q, want %q", got, want)
	}

	// 2. Probe 9p-fd (alias for 9pnet_fd, file on disk is 9pnet-fd.ko;
	// 9pnet.ko already loaded).
	if err := prober.Probe("9p-fd", ""); err != nil {
		t.Fatalf("Probe(9p-fd) failed: %v", err)
	}
	if len(loadOrder) != 3 || loadOrder[2] != "9pnet-fd.ko" {
		t.Errorf("loadOrder after 9p-fd = %v, want 9pnet-fd.ko appended", loadOrder)
	}

	// 3. Probe builtin aliases (from modules.builtin and
	// modules.builtin.modinfo).
	if err := prober.Probe("fs-ext4", ""); err != nil {
		t.Fatalf("Probe(fs-ext4) failed: %v", err)
	}
	if err := prober.Probe("fs-builtin-only", ""); err != nil {
		t.Fatalf("Probe(fs-builtin-only) failed: %v", err)
	}
	if len(loadOrder) != 3 {
		t.Errorf("builtin probe should not invoke loadCB, got loadOrder = %v", loadOrder)
	}

	// 4. Probe modules.alias entry whose target is overridden by an alias
	// in /etc/modprobe.d (hw-chained -> middle_target -> override_target).
	if err := prober.Probe("hw-chained", ""); err != nil {
		t.Fatalf("Probe(hw-chained) failed: %v", err)
	}
	if len(loadOrder) != 4 || loadOrder[3] != "override_target.ko" {
		t.Errorf("expected override_target.ko to be loaded via chained alias, got %v", loadOrder)
	}

	// 5. Blacklist handling:
	// - Via alias ("bad-alias" -> "bad_mod"): always skipped by blacklist.
	if err := prober.Probe("bad-alias", ""); err != nil {
		t.Fatalf("Probe(bad-alias) failed: %v", err)
	}
	if _, ok := loadedParams["bad_mod.ko"]; ok {
		t.Errorf("blacklisted module bad_mod.ko should not be loaded via alias")
	}
	// - Direct name with UseBlacklist = true: skipped.
	prober.opts.UseBlacklist = true
	if err := prober.Probe("bad_mod", ""); err != nil {
		t.Fatalf("Probe(bad_mod) with UseBlacklist failed: %v", err)
	}
	if _, ok := loadedParams["bad_mod.ko"]; ok {
		t.Errorf("blacklisted module bad_mod.ko should not be loaded when UseBlacklist is true")
	}
	// - Direct name with UseBlacklist = false: loaded.
	prober.opts.UseBlacklist = false
	if err := prober.Probe("bad_mod", ""); err != nil {
		t.Fatalf("Probe(bad_mod) without UseBlacklist failed: %v", err)
	}
	if _, ok := loadedParams["bad_mod.ko"]; !ok {
		t.Errorf("expected bad_mod.ko to be loaded when probed directly with UseBlacklist=false")
	}
}

func TestRecursiveSoftdepsAndAliases(t *testing.T) {
	root := t.TempDir()
	kver := "6.6.0-test"
	modDir := filepath.Join(root, "lib/modules", kver)
	etcModprobe := filepath.Join(root, "etc/modprobe.d")

	for _, d := range []string{modDir, etcModprobe} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	// ipmi_devintf has a hard dependency on ipmi_msghandler, while
	// ipmi_msghandler has a post-softdep on ipmi_devintf and a pre-softdep
	// on pre_helper.
	// Also test a 2-hop post-softdep cycle (top_mod -> mid_mod, where
	// mid_mod has pre: top_mod and post: post_leaf, and post_leaf depends
	// on top_mod) to ensure neither pre- nor post-softdeps deadlock or
	// poison top_mod.
	modulesDep := `kernel/drivers/char/ipmi/pre_helper.ko:
kernel/drivers/char/ipmi/ipmi_msghandler.ko:
kernel/drivers/char/ipmi/ipmi_devintf.ko: kernel/drivers/char/ipmi/ipmi_msghandler.ko
kernel/crypto/sym_provider.ko:
kernel/drivers/misc/mid_mod.ko:
kernel/drivers/misc/top_mod.ko: kernel/drivers/misc/mid_mod.ko
kernel/drivers/misc/post_leaf.ko: kernel/drivers/misc/top_mod.ko
`
	if err := os.WriteFile(filepath.Join(modDir, "modules.dep"), []byte(modulesDep), 0o644); err != nil {
		t.Fatal(err)
	}

	modulesSoftdep := `softdep ipmi_msghandler pre: pre_helper post: ipmi_devintf
softdep mid_mod pre: top_mod post: post_leaf
`
	if err := os.WriteFile(filepath.Join(modDir, "modules.softdep"), []byte(modulesSoftdep), 0o644); err != nil {
		t.Fatal(err)
	}

	modulesAlias := `alias base-ipmi ipmi_msghandler
`
	if err := os.WriteFile(filepath.Join(modDir, "modules.alias"), []byte(modulesAlias), 0o644); err != nil {
		t.Fatal(err)
	}

	modulesSymbols := `alias symbol:my_exported_sym sym_provider
`
	if err := os.WriteFile(filepath.Join(modDir, "modules.symbols"), []byte(modulesSymbols), 0o644); err != nil {
		t.Fatal(err)
	}

	// Chained alias (chained-ipmi -> base-ipmi -> ipmi_msghandler) and a
	// circular alias pair (loop-a <-> loop-b).
	etcConf := `alias chained-ipmi base-ipmi
alias loop-a loop-b
alias loop-b loop-a
`
	if err := os.WriteFile(filepath.Join(etcModprobe, "recursive.conf"), []byte(etcConf), 0o644); err != nil {
		t.Fatal(err)
	}

	prober, err := NewProber(ProbeOpts{
		RootDir:        root,
		KVer:           kver,
		IgnoreProcMods: true,
	})
	if err != nil {
		t.Fatalf("NewProber failed: %v", err)
	}

	var loadOrder []string
	prober.loadCB = func(modPath, params string) error {
		loadOrder = append(loadOrder, filepath.Base(modPath))
		return nil
	}

	// 1. Probing chained-ipmi resolves chained-ipmi -> base-ipmi ->
	// ipmi_msghandler, loads pre-softdep pre_helper.ko first, then
	// ipmi_msghandler.ko, then post-softdep ipmi_devintf.ko (which has a
	// hard dependency on ipmi_msghandler.ko) without deadlocking.
	if err := prober.Probe("chained-ipmi", ""); err != nil {
		t.Fatalf("Probe(chained-ipmi) failed: %v", err)
	}
	wantOrder := []string{"pre_helper.ko", "ipmi_msghandler.ko", "ipmi_devintf.ko"}
	if !slices.Equal(loadOrder, wantOrder) {
		t.Errorf("loadOrder = %v, want %v", loadOrder, wantOrder)
	}

	// 2. Probing symbol:my_exported_sym (used by kernel symbol_request)
	// resolves via modules.symbols.
	if err := prober.Probe("symbol:my_exported_sym", ""); err != nil {
		t.Fatalf("Probe(symbol:my_exported_sym) failed: %v", err)
	}
	if len(loadOrder) != 4 || loadOrder[3] != "sym_provider.ko" {
		t.Errorf("expected sym_provider.ko to be loaded, got %v", loadOrder)
	}

	// 3. Probing top_mod where mid_mod has pre: top_mod and post: post_leaf
	// (which depends on top_mod) skips the cyclic softdeps instead of
	// deadlocking or poisoning top_mod.
	loadOrder = nil
	if err := prober.Probe("top_mod", ""); err != nil {
		t.Fatalf("Probe(top_mod) failed: %v", err)
	}
	if want := []string{"mid_mod.ko", "top_mod.ko"}; !slices.Equal(loadOrder, want) {
		t.Errorf("loadOrder for top_mod = %v, want %v", loadOrder, want)
	}

	// 4. Probing a circular alias returns an error instead of hanging.
	if err := prober.Probe("loop-a", ""); err == nil {
		t.Errorf("expected circular alias loop-a to return an error")
	}
}

func TestRemove(t *testing.T) {
	m := depMap{
		"/lib/modules/6.6/kernel/net/9p/9pnet.ko": &dependency{
			state: loaded,
		},
		"/lib/modules/6.6/kernel/fs/netfs/netfs.ko": &dependency{
			state: loaded,
		},
		"/lib/modules/6.6/kernel/fs/9p/9p.ko": &dependency{
			state: loaded,
			deps: []string{
				"/lib/modules/6.6/kernel/net/9p/9pnet.ko",
				"/lib/modules/6.6/kernel/fs/netfs/netfs.ko",
			},
		},
		"/lib/modules/6.6/kernel/net/9p/9pnet-fd.ko": &dependency{
			state: loaded,
			deps: []string{
				"/lib/modules/6.6/kernel/net/9p/9pnet.ko",
			},
		},
		"/lib/modules/6.6/kernel/fs/ext4/ext4.ko": &dependency{
			state: builtin,
		},
	}

	var deleted []string
	prober := Prober{
		deps: m,
		moduleAliases: []modAlias{
			{pattern: "fs-9p", target: "9p"},
		},
		deleteCB: func(modName string) error {
			deleted = append(deleted, modName)
			return nil
		},
	}

	// Removing 9pnet while 9p and 9pnet-fd are still loaded should fail.
	if err := prober.Remove("9pnet"); err == nil {
		t.Errorf("expected error when removing in-use dependency 9pnet")
	}

	// Removing fs-9p (alias for 9p) should unload 9p and netfs, but keep
	// 9pnet because 9pnet-fd is still loaded and depends on 9pnet.
	if err := prober.Remove("fs-9p"); err != nil {
		t.Fatalf("Remove(fs-9p) failed: %v", err)
	}
	wantDeleted := []string{"9p", "netfs"}
	if !slices.Equal(deleted, wantDeleted) {
		t.Errorf("deleted = %v, want %v", deleted, wantDeleted)
	}

	// Now removing 9pnet-fd should unload 9pnet_fd and then 9pnet.
	deleted = nil
	if err := prober.Remove("9pnet-fd"); err != nil {
		t.Fatalf("Remove(9pnet-fd) failed: %v", err)
	}
	wantDeleted = []string{"9pnet_fd", "9pnet"}
	if !slices.Equal(deleted, wantDeleted) {
		t.Errorf("deleted = %v, want %v", deleted, wantDeleted)
	}

	// Re-probing fs-9p after removal on the same Prober should reload
	// 9pnet, netfs, and 9p.
	var reloadMu sync.Mutex
	var reloaded []string
	prober.loadCB = func(modPath, params string) error {
		reloadMu.Lock()
		defer reloadMu.Unlock()
		reloaded = append(reloaded, filepath.Base(modPath))
		return nil
	}
	if err := prober.Probe("fs-9p", ""); err != nil {
		t.Fatalf("re-Probe(fs-9p) after Remove failed: %v", err)
	}
	if len(reloaded) != 3 {
		t.Errorf("expected 3 modules reloaded after Remove, got %v", reloaded)
	}

	// Removing a builtin module should return an error.
	if err := prober.Remove("ext4"); err == nil {
		t.Errorf("expected error when removing builtin module ext4")
	}

	// If Delete returns EBUSY on the target module, Remove should fail.
	m["/lib/modules/6.6/kernel/net/9p/9pnet.ko"].state = loaded
	m["/lib/modules/6.6/kernel/fs/9p/9p.ko"].state = unloaded
	prober.deleteCB = func(modName string) error {
		return unix.EBUSY
	}
	if err := prober.Remove("9pnet"); !errors.Is(err, unix.EBUSY) {
		t.Errorf("expected EBUSY error when target module is busy, got %v", err)
	}
}

func TestParallelLoad(t *testing.T) {
	loadTime := 100 * time.Millisecond

	m := depMap{
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/hid-generic.ko":   &dependency{},
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/usbhid/usbhid.ko": &dependency{},
		"/lib/modules/6.6.6-generic/kernel/crypto/ccm.ko":                &dependency{},
		"/lib/modules/6.6.6-generic/kernel/tests/depmod.ko": &dependency{
			deps: []string{"/lib/modules/6.6.6-generic/kernel/crypto/ccm.ko",
				"/lib/modules/6.6.6-generic/kernel/drivers/hid/usbhid/usbhid.ko",
				"/lib/modules/6.6.6-generic/kernel/drivers/hid/hid-generic.ko",
			},
		},
		"/lib/modules/6.6.6-generic/kernel/tests/depmod2.ko": &dependency{
			deps: []string{"/lib/modules/6.6.6-generic/kernel/crypto/ccm.ko",
				"/lib/modules/6.6.6-generic/kernel/drivers/hid/usbhid/usbhid.ko",
			},
		},
	}

	var eg errgroup.Group

	prober := Prober{
		deps: m,
		// Wait time encourages racing between dependencies.
		opts: ProbeOpts{DryRunCB: func(path string) { time.Sleep(loadTime) }},
	}

	eg.Go(func() error {
		return prober.Probe("depmod", "")
	})
	eg.Go(func() error {
		return prober.Probe("depmod2", "")
	})

	start := time.Now()
	err := eg.Wait()

	if err != nil {
		t.Fatalf("Probing failed: %v", err)
	}

	// Racy... longest parallel chain is 2, and we load 5 total modules.
	if time.Since(start) > loadTime*3 {
		t.Fatalf("module loading slow")
	}

	for mod, d := range m {
		if d.state != loaded {
			t.Fatalf("mod %q should have been loaded", path.Base(mod))
		}
	}
}

func TestInvalidCircularLoad(t *testing.T) {
	m := depMap{
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/hid-generic.ko":   &dependency{},
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/usbhid/usbhid.ko": &dependency{},
		"/lib/modules/6.6.6-generic/kernel/crypto/ccm.ko":                &dependency{},
		"/lib/modules/6.6.6-generic/kernel/tests/circlemod.ko": &dependency{
			deps: []string{"/lib/modules/6.6.6-generic/kernel/tests/depmod.ko"},
		},
		"/lib/modules/6.6.6-generic/kernel/tests/depmod.ko": &dependency{
			deps: []string{"/lib/modules/6.6.6-generic/kernel/crypto/ccm.ko",
				"/lib/modules/6.6.6-generic/kernel/drivers/hid/usbhid/usbhid.ko",
				"/lib/modules/6.6.6-generic/kernel/tests/circlemod.ko",
			},
		},
	}

	prober := Prober{
		deps: m,
		opts: ProbeOpts{DryRunCB: func(path string) {}},
	}

	var wg sync.WaitGroup
	errs := make([]error, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		errs[0] = prober.Probe("depmod", "")
	}()
	go func() {
		defer wg.Done()
		errs[1] = prober.Probe("circlemod", "")
	}()
	wg.Wait()

	for i, err := range errs {
		if err == nil {
			t.Fatalf("Circular dep goroutine %d should have errored...", i)
		}
	}
}

func TestProbeRemoveOptionsAndEdgeCases(t *testing.T) {
	root := t.TempDir()
	kver := "6.6.0-usr"
	modDir := filepath.Join(root, "usr/lib/modules", kver)
	etcModprobe := filepath.Join(root, "etc/modprobe.d")
	procDir := filepath.Join(root, "proc")

	for _, d := range []string{modDir, etcModprobe, procDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	modulesDep := `kernel/net/9p/9pnet.ko.xz:
kernel/fs/9p/9p.ko.zst: kernel/net/9p/9pnet.ko.xz
kernel/drivers/misc/pre_mod.ko.gz:
kernel/drivers/misc/post_mod.ko:
`
	if err := os.WriteFile(filepath.Join(modDir, "modules.dep"), []byte(modulesDep), 0o644); err != nil {
		t.Fatal(err)
	}

	// Test softdep in modprobe.d with pre:mod and post:mod (no space after
	// colon), self-alias (ignored), short lines, and trailing backslash at
	// EOF.
	conf := `invalid_short_line
alias 9p 9p
softdep 9p pre:pre_mod post:post_mod
options 9p debug=0 \
`
	if err := os.WriteFile(filepath.Join(etcModprobe, "extra.conf"), []byte(conf), 0o644); err != nil {
		t.Fatal(err)
	}

	// Test /proc/cmdline boolean flag (no '=') and parameter override.
	cmdline := `9p.noext 9p.debug=1 modprobe.blacklist=` + "\n"
	if err := os.WriteFile(filepath.Join(procDir, "cmdline"), []byte(cmdline), 0o644); err != nil {
		t.Fatal(err)
	}

	// Test /proc/modules loaded during genDeps when IgnoreProcMods is false.
	procMods := `pre_mod 16384 0 - Live 0x0` + "\n"
	if err := os.WriteFile(filepath.Join(procDir, "modules"), []byte(procMods), 0o644); err != nil {
		t.Fatal(err)
	}

	var dryRunLoaded []string
	opts := ProbeOpts{
		RootDir: root,
		KVer:    kver,
		DryRunCB: func(mp string) {
			dryRunLoaded = append(dryRunLoaded, filepath.Base(mp))
		},
	}

	if err := ProbeOptions("9p", "debug=2", opts); err != nil {
		t.Fatalf("ProbeOptions(9p) failed: %v", err)
	}
	// pre_mod was already marked Live in /proc/modules, so only 9pnet, 9p,
	// and post_mod should be dry-run loaded.
	wantLoaded := []string{"9pnet.ko.xz", "9p.ko.zst", "post_mod.ko"}
	if !slices.Equal(dryRunLoaded, wantLoaded) {
		t.Errorf("ProbeOptions dryRunLoaded = %v, want %v", dryRunLoaded, wantLoaded)
	}

	var dryRunRemoved []string
	opts.DryRunCB = func(mp string) {
		dryRunRemoved = append(dryRunRemoved, filepath.Base(mp))
	}
	if err := RemoveOptions("9p", opts); err != nil {
		t.Fatalf("RemoveOptions(9p) failed: %v", err)
	}
	if !slices.Contains(dryRunRemoved, "9p.ko.zst") {
		t.Errorf("expected 9p.ko.zst in RemoveOptions dryRunRemoved, got %v", dryRunRemoved)
	}

	// Edge cases in normalizeAlias, findModPath, and parseSoftdepFields.
	if !matchAlias(`foo\*bar`, `foo*bar`) {
		t.Errorf("expected escaped wildcard in matchAlias to match")
	}
	if _, err := findModPath("", depMap{}); err == nil {
		t.Errorf("expected error for empty module name in findModPath")
	}
	dupMap := depMap{
		"/b/dup_mod.ko": &dependency{},
		"/a/dup-mod.ko": &dependency{},
	}
	if got, err := findModPath("dup_mod", dupMap); err != nil || got != "/a/dup-mod.ko" {
		t.Errorf("findModPath duplicate deterministic sort = (%q, %v), want (/a/dup-mod.ko, nil)", got, err)
	}
	parseSoftdepFields("", []string{"pre:", "foo"}, map[string]softDep{})
}

// Helper function to generate compression test data for TestCompression.
// Generates a map with the name of the file as key and the compressed data as
// value. The data is compressed using xz, gzip, and zstd. There is also one
// file with a bad extension to test the error handling. Testing compression
// itself is out of scope.
func generateCompressionTestData(data []byte) (map[string][]byte, error) {
	var compressionBuffer bytes.Buffer

	tData := make(map[string][]byte)

	// 0. ko
	tData["test.ko"] = make([]byte, len(data))
	copy(tData["test.xz"], data)

	// 1. xz
	wXZ, err := xz.NewWriter(&compressionBuffer)
	if err != nil {
		return nil, err
	}
	_, err = wXZ.Write(data)
	if err != nil {
		return nil, err
	}
	if err = wXZ.Close(); err != nil {
		return nil, err
	}

	tData["test.xz"] = make([]byte, compressionBuffer.Len())
	copy(tData["test.xz"], compressionBuffer.Bytes())
	compressionBuffer.Reset()

	// 2. gzip
	wGZ := gzip.NewWriter(&compressionBuffer)
	if _, err = wGZ.Write(data); err != nil {
		return nil, err
	}
	if err = wGZ.Close(); err != nil {
		return nil, err
	}

	tData["test.gz"] = make([]byte, compressionBuffer.Len())
	copy(tData["test.gz"], compressionBuffer.Bytes())
	compressionBuffer.Reset()

	// 3. zstd
	wZST, err := zstd.NewWriter(&compressionBuffer)
	if err != nil {
		return nil, err
	}

	if _, err = wZST.Write(data); err != nil {
		return nil, err
	}

	if err = wZST.Close(); err != nil {
		return nil, err
	}

	tData["test.zst"] = make([]byte, compressionBuffer.Len())
	copy(tData["test.zst"], compressionBuffer.Bytes())
	compressionBuffer.Reset()

	// 4. bad
	tData["test.bad"] = []byte{'b', 'a', 'd'}
	return tData, nil
}

// Since we don't need to test the compression function, we just check
// for validity of file extension detection.
func TestCompression(t *testing.T) {
	const compressionTestString = "test\x00"

	tDir := t.TempDir()
	tFd := make(map[string]*os.File, 4)
	tFiles, err := generateCompressionTestData([]byte(compressionTestString))
	if err != nil {
		t.Fatalf("failed to generate test data: '%v'\n", err)
	}

	for name, data := range tFiles {
		tFd[name], err = os.Create(path.Join(tDir, name))
		if err != nil {
			t.Fatalf("failed to create test file %q: '%v'\n", name, err)
		}
		defer tFd[name].Close()

		n, err := tFd[name].Write(data)
		if err != nil {
			t.Fatalf("failed to write to test file %q: '%v'\n", name, err)
		}

		if err = tFd[name].Sync(); err != nil {
			t.Fatalf("failed to sync test file %q: '%v'\n", name, err)
		}

		if _, err := tFd[name].Seek(0, 0); err != nil {
			t.Fatalf("failed to seek to beginning of test file %q: '%v'\n", name, err)
		}

		if n != len(data) {
			t.Fatalf("failed to write all data to test file %q. Expected %d bytes, wrote %d\n", name, len(data), n)
		}
	}

	testCases := map[string]struct {
		file    *os.File
		ext     string
		isError bool
		err     error
	}{
		"test.ko": {
			file:    tFd["test.ko"],
			ext:     ".ko",
			isError: false,
			err:     nil,
		},
		"test.xz": {
			file:    tFd["test.xz"],
			ext:     ".xz",
			isError: false,
			err:     nil,
		},
		"test.gz": {
			file:    tFd["test.gz"],
			ext:     ".gz",
			isError: false,
			err:     nil,
		},
		"test.zst": {
			file:    tFd["test.zst"],
			ext:     ".zst",
			isError: false,
			err:     nil,
		},
		"test.bad": {
			file:    tFd["test.bad"],
			ext:     ".bad",
			isError: true,
			err:     os.ErrNotExist,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := compressionReader(tc.file)
			if tc.isError {
				if errors.Is(err, tc.err) {
					return
				}
				t.Fatalf("expected error %v but got '%v'\n", tc.err, err)
			}
			if !tc.isError && err != nil {
				t.Fatalf("expected no error but got '%v'\n", err)
			}
		})
	}
}
