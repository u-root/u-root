// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package kmodule interfaces with Linux kernel modules.
//
// kmodule allows loading and unloading kernel modules with dependencies, as
// well as locating them through probing.
package kmodule

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sys/unix"
)

// Flags to finit_module(2) / FileInit.
const (
	// Ignore symbol version hashes.
	MODULE_INIT_IGNORE_MODVERSIONS = 0x1

	// Ignore kernel version magic.
	MODULE_INIT_IGNORE_VERMAGIC = 0x2
)

type modState uint8

const (
	unloaded modState = iota
	loading
	loaded
	builtin
)

type dependency struct {
	mu    sync.Mutex
	state modState
	deps  []string
}

type depMap map[string]*dependency

type modAlias struct {
	pattern string
	target  string
}

type softDep struct {
	pre  []string
	post []string
}

// ProbeOpts contains optional parameters to Probe.
//
// An empty ProbeOpts{} should lead to the default behavior.
type ProbeOpts struct {
	DryRunCB       func(string)
	RootDir        string
	KVer           string
	IgnoreProcMods bool
	UseBlacklist   bool
}

// Prober provides a concurrent-safe interface for probing modules.
type Prober struct {
	deps          depMap
	modIndex      map[string]string
	opts          ProbeOpts
	configAliases []modAlias
	moduleAliases []modAlias
	confOpts      map[string][]string
	cmdlineOpts   map[string][]string
	blacklist     map[string]bool
	softDeps      map[string]softDep
	loadCB        func(modPath, params string) error
	deleteCB      func(modName string) error
}

// Init loads the kernel module given by image with the given options.
func Init(image []byte, opts string) error {
	return unix.InitModule(image, opts)
}

// Wrapper for the compression readers.
func CompressionReader(file *os.File) (reader io.Reader, err error) {
	return compressionReader(file)
}

// FileInit loads the kernel module contained by `f` with the given opts and
// flags. Uncompresses modules with a .xz, .gz, and .zst suffix before loading.
//
// FileInit falls back to init_module(2) via Init when the finit_module(2)
// syscall is not available and when loading compressed modules.
func FileInit(f *os.File, opts string, flags uintptr) error {
	var r io.Reader
	var err error

	if r, err = CompressionReader(f); err != nil {
		return err
	}

	if r == nil {
		err := unix.FinitModule(int(f.Fd()), opts, int(flags))
		if err == unix.ENOSYS {
			if flags != 0 {
				return err
			}
			// Fall back to init_module(2).
			r = f
		} else {
			return err
		}
	}

	img, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return Init(img, opts)
}

// Delete removes a kernel module.
func Delete(name string, flags uintptr) error {
	return unix.DeleteModule(normalizeModName(name), int(flags))
}

// hasHardCycle checks whether the hard-dependency subgraph in p.deps rooted at
// startPath contains a directed cycle. Checking this before acquiring
// per-module locks prevents concurrent Probe calls on cyclic graphs from
// deadlocking.
func (p *Prober) hasHardCycle(startPath string) bool {
	const (
		unvisited uint8 = iota
		visiting
		done
	)
	color := make(map[string]uint8)
	var dfs func(curr string) bool
	dfs = func(curr string) bool {
		color[curr] = visiting
		if d, ok := p.deps[curr]; ok {
			for _, next := range d.deps {
				switch color[next] {
				case visiting:
					return true
				case unvisited:
					if dfs(next) {
						return true
					}
				}
			}
		}
		color[curr] = done
		return false
	}
	return dfs(startPath)
}

// reachesAny reports whether startPath can reach any module path in targets by
// following hard dependencies (d.deps) or pre-softdeps. This is used to avoid
// probing softdeps that depend back on an ancestor module currently being
// loaded.
func (p *Prober) reachesAny(startPath string, targets []string) bool {
	if len(targets) == 0 {
		return false
	}
	visited := make(map[string]bool)
	var dfs func(curr string) bool
	dfs = func(curr string) bool {
		if slices.Contains(targets, curr) {
			return true
		}
		if visited[curr] {
			return false
		}
		visited[curr] = true
		d, ok := p.deps[curr]
		if !ok {
			return false
		}
		if slices.ContainsFunc(d.deps, dfs) {
			return true
		}
		if stem, ok := modNameFromPath(curr); ok && p.softDeps != nil {
			sd := p.softDeps[normalizeModName(stem)]
			for _, preName := range sd.pre {
				resolved, _, err := p.resolveTargets(preName, true)
				if err != nil {
					continue
				}
				for _, t := range resolved {
					if dfs(t.modPath) {
						return true
					}
				}
			}
		}
		return false
	}
	return dfs(startPath)
}

// parallelProbeDep loads a module, its hard dependencies, and any pre/post
// soft dependencies. Returns once the module is loaded.
func (p *Prober) parallelProbeDep(modPath, extraConfParams, userParams string, modStack []string) error {
	dep, present := p.deps[modPath]
	if !present {
		return fmt.Errorf("module not in depmap %s", modPath)
	}

	if slices.Contains(modStack, modPath) || p.hasHardCycle(modPath) {
		return fmt.Errorf("circular dependency detected when loading %s", modPath)
	}

	var sd softDep
	if stem, ok := modNameFromPath(modPath); ok && p.softDeps != nil {
		sd = p.softDeps[normalizeModName(stem)]
	}

	nextStack := append(slices.Clip(modStack), modPath)

	dep.mu.Lock()
	if dep.state == builtin || dep.state == loaded {
		dep.mu.Unlock()
		return nil
	}

	// Probe pre-softdeps before loading the module (best-effort, matching
	// libkmod behavior). Skip any pre-softdep that has a cycle or depends
	// back on an ancestor in nextStack.
	for _, preMod := range sd.pre {
		targets, _, err := p.resolveTargets(preMod, true)
		if err != nil {
			continue
		}
		for _, t := range targets {
			if p.hasHardCycle(t.modPath) || p.reachesAny(t.modPath, nextStack) {
				continue
			}
			_ = p.parallelProbeDep(t.modPath, t.extraConfParams, "", nextStack)
		}
	}

	var eg errgroup.Group
	for _, childDep := range dep.deps {
		eg.Go(func() error {
			return p.parallelProbeDep(childDep, "", "", nextStack)
		})
	}

	if err := eg.Wait(); err != nil {
		dep.mu.Unlock()
		return err
	}

	fullParams := p.buildModParams(modPath, extraConfParams, userParams)
	var loadErr error
	if p.loadCB != nil {
		loadErr = p.loadCB(modPath, fullParams)
	} else {
		loadErr = loadModule(modPath, fullParams, p.opts)
	}
	if loadErr != nil {
		dep.mu.Unlock()
		return loadErr
	}
	dep.state = loaded
	dep.mu.Unlock()

	// Probe post-softdeps after unlocking dep.mu so that post-softdeps
	// which depend on modPath see dep.state == loaded without deadlocking.
	// Skip any post-softdep that depends on an still-unfinished ancestor
	// in modStack.
	for _, postMod := range sd.post {
		targets, _, err := p.resolveTargets(postMod, true)
		if err != nil {
			continue
		}
		for _, t := range targets {
			if t.modPath == modPath || p.hasHardCycle(t.modPath) || p.reachesAny(t.modPath, modStack) {
				continue
			}
			_ = p.parallelProbeDep(t.modPath, t.extraConfParams, "", modStack)
		}
	}

	return nil
}

func normalizeParamToken(tok string) string {
	k, v, ok := strings.Cut(tok, "=")
	if !ok {
		return normalizeModName(tok)
	}
	v = strings.Trim(v, `"'`)
	return normalizeModName(k) + "=" + v
}

func (p *Prober) buildModParams(modPath, extraConfParams, userParams string) string {
	var parts []string
	if extraConfParams != "" {
		parts = append(parts, extraConfParams)
	}

	if stem, ok := modNameFromPath(modPath); ok {
		norm := normalizeModName(stem)
		if p.confOpts != nil {
			parts = append(parts, p.confOpts[norm]...)
		}
		if p.cmdlineOpts != nil {
			userTokens := splitCmdlineTokens(userParams)
			normUser := make([]string, len(userTokens))
			for i, u := range userTokens {
				normUser[i] = normalizeParamToken(u)
			}
			for _, opt := range p.cmdlineOpts[norm] {
				if !slices.Contains(normUser, normalizeParamToken(opt)) {
					parts = append(parts, opt)
				}
			}
		}
	}

	if trimmed := strings.TrimSpace(userParams); trimmed != "" {
		parts = append(parts, trimmed)
	}

	return strings.Join(parts, " ")
}

// NewProber looks up current module state and provides Prober.
func NewProber(opts ProbeOpts) (*Prober, error) {
	p := &Prober{
		opts:        opts,
		modIndex:    make(map[string]string),
		confOpts:    make(map[string][]string),
		cmdlineOpts: make(map[string][]string),
		blacklist:   make(map[string]bool),
		softDeps:    make(map[string]softDep),
	}

	deps, modIndex, moduleDir, err := genDeps(opts)
	if err != nil {
		return nil, fmt.Errorf("could not generate dependency map %w", err)
	}
	p.deps = deps
	p.modIndex = modIndex

	if moduleDir != "" {
		parseModuleSoftdeps(moduleDir, p)
		parseModuleAliases(moduleDir, "modules.alias", p)
		parseModuleAliases(moduleDir, "modules.symbols", p)
		parseBuiltinModinfo(moduleDir, p)
	}
	parseModprobeConfigs(opts.RootDir, p)
	parseProcCmdline(opts.RootDir, p)

	return p, nil
}

func (p *Prober) lookupModPath(name string) (string, error) {
	target := normalizeModName(name)
	if target == "" {
		return "", fmt.Errorf("could not find path for module %q", name)
	}
	if p.modIndex != nil {
		if mp, ok := p.modIndex[target]; ok {
			if _, exists := p.deps[mp]; exists {
				return mp, nil
			}
		}
	}
	return findModPath(name, p.deps)
}

type resolvedTarget struct {
	modPath         string
	extraConfParams string
}

func (p *Prober) resolveTargets(name string, honorBlacklistForDirect bool) ([]resolvedTarget, bool, error) {
	return p.resolveTargetsVisited(name, honorBlacklistForDirect, make(map[string]bool))
}

func (p *Prober) resolveTargetsVisited(name string, honorBlacklistForDirect bool, visitedAliases map[string]bool) ([]resolvedTarget, bool, error) {
	var targets []resolvedTarget
	seen := make(map[string]bool)
	matchedAny := false

	normName := normalizeModName(name)
	var aliasExtra string
	if p.confOpts != nil && normName != "" {
		aliasExtra = strings.Join(p.confOpts[normName], " ")
	}

	appendTarget := func(t resolvedTarget) {
		if seen[t.modPath] {
			return
		}
		seen[t.modPath] = true
		if aliasExtra != "" {
			if stem, ok := modNameFromPath(t.modPath); !ok || normalizeModName(stem) != normName {
				if t.extraConfParams != "" {
					t.extraConfParams = aliasExtra + " " + t.extraConfParams
				} else {
					t.extraConfParams = aliasExtra
				}
			}
		}
		targets = append(targets, t)
	}

	resolveAliasTarget := func(targetName string) error {
		normTarget := normalizeModName(targetName)
		if visitedAliases[normTarget] {
			return fmt.Errorf("circular alias detected for %q", targetName)
		}
		subTargets, subMatched, err := p.resolveTargetsVisited(targetName, true, visitedAliases)
		if !subMatched || err != nil {
			if err != nil {
				return err
			}
			return fmt.Errorf("could not find path for module %q", targetName)
		}
		for _, st := range subTargets {
			appendTarget(st)
		}
		return nil
	}

	// 1. Check aliases from modprobe.d configuration files first.
	if !visitedAliases[normName] {
		visitedAliases[normName] = true
		var configErr error
		for _, ca := range p.configAliases {
			if normalizeModName(ca.target) == normName {
				continue
			}
			if matchAlias(ca.pattern, name) {
				matchedAny = true
				if err := resolveAliasTarget(ca.target); err != nil {
					configErr = errors.Join(configErr, err)
				}
			}
		}
		delete(visitedAliases, normName)
		if matchedAny {
			if len(targets) > 0 || configErr != nil {
				return targets, true, configErr
			}
			// If all config alias targets were blacklisted and
			// name is not a direct module, return empty targets.
			if _, err := p.lookupModPath(name); err != nil {
				return nil, true, nil
			}
			matchedAny = false
		}
	}

	// 2. Check direct module name in modules.dep / modules.builtin.
	if modPath, err := p.lookupModPath(name); err == nil {
		if honorBlacklistForDirect && p.blacklist[normName] {
			return nil, true, nil
		}
		return []resolvedTarget{{modPath: modPath}}, true, nil
	}

	// 3. Check aliases from modules.alias, modules.symbols, and
	// modules.builtin.modinfo.
	if !visitedAliases[normName] {
		visitedAliases[normName] = true
		var aliasErr error
		for _, ma := range p.moduleAliases {
			if matchAlias(ma.pattern, name) {
				matchedAny = true
				if err := resolveAliasTarget(ma.target); err != nil {
					aliasErr = errors.Join(aliasErr, err)
				}
			}
		}
		delete(visitedAliases, normName)
		if matchedAny {
			if len(targets) > 0 {
				return targets, true, nil
			}
			return nil, true, aliasErr
		}
	}

	return nil, false, fmt.Errorf("could not find path for module %q", name)
}

// Probe loads the given kernel module (or alias) and its dependencies.
func (p *Prober) Probe(name string, modParams string) error {
	targets, matched, err := p.resolveTargets(name, p.opts.UseBlacklist)
	if !matched || err != nil {
		return err
	}

	var probeErr error
	for _, t := range targets {
		if err := p.parallelProbeDep(t.modPath, t.extraConfParams, modParams, []string{}); err != nil {
			probeErr = errors.Join(probeErr, err)
		}
	}
	return probeErr
}

// Probe loads the given kernel module and its dependencies.
// It calls ProbeOptions with the default ProbeOpts.
func Probe(name string, modParams string) error {
	return ProbeOptions(name, modParams, ProbeOpts{})
}

// ProbeOptions loads the given kernel module and its dependencies.
// This function takes ProbeOpts.
func ProbeOptions(name, modParams string, opts ProbeOpts) error {
	p, err := NewProber(opts)
	if err != nil {
		return err
	}

	return p.Probe(name, modParams)
}

// Remove unloads the given kernel module (or alias) and any unused
// dependencies.
func (p *Prober) Remove(name string) error {
	targets, matched, _ := p.resolveTargets(name, false)
	if !matched || len(targets) == 0 {
		return fmt.Errorf("could not find path for module %q", name)
	}

	var removeErr error
	for _, t := range targets {
		if err := p.removeModPath(t.modPath); err != nil {
			removeErr = errors.Join(removeErr, err)
		}
	}
	return removeErr
}

func (p *Prober) removeModPath(modPath string) error {
	dep, ok := p.deps[modPath]
	if !ok {
		return fmt.Errorf("module not in depmap %s", modPath)
	}
	dep.mu.Lock()
	state := dep.state
	dep.mu.Unlock()
	if state == builtin {
		return fmt.Errorf("module %s is builtin", modPath)
	}

	// Build post-order list of the dependency subgraph rooted at modPath,
	// then reverse it so every dependent is unloaded before its deps.
	var postOrder []string
	visited := make(map[string]bool)
	var visit func(string)
	visit = func(curr string) {
		if visited[curr] {
			return
		}
		visited[curr] = true
		if d, exists := p.deps[curr]; exists {
			for _, child := range d.deps {
				visit(child)
			}
		}
		postOrder = append(postOrder, curr)
	}
	visit(modPath)
	slices.Reverse(postOrder)

	for i, mp := range postOrder {
		d := p.deps[mp]
		if d == nil {
			continue
		}
		isTarget := i == 0

		d.mu.Lock()
		if d.state == builtin {
			d.mu.Unlock()
			continue
		}
		if p.isDepUsed(mp) {
			d.mu.Unlock()
			if isTarget {
				return fmt.Errorf("module %s is in use", mp)
			}
			continue
		}

		if p.opts.DryRunCB != nil {
			if d.state == loaded || isTarget {
				p.opts.DryRunCB(mp)
				d.state = unloaded
			}
			d.mu.Unlock()
			continue
		}

		if !isTarget && d.state != loaded {
			d.mu.Unlock()
			continue
		}

		stem, ok := modNameFromPath(mp)
		if !ok {
			d.mu.Unlock()
			continue
		}
		var err error
		if p.deleteCB != nil {
			err = p.deleteCB(normalizeModName(stem))
		} else {
			err = Delete(stem, unix.O_NONBLOCK)
		}
		if err != nil {
			d.mu.Unlock()
			if isTarget {
				return err
			}
			// Dependency is still in use or already unloaded;
			// ignore.
			continue
		}
		d.state = unloaded
		d.mu.Unlock()
	}

	return nil
}

func (p *Prober) isDepUsed(targetPath string) bool {
	for mp, d := range p.deps {
		if mp == targetPath {
			continue
		}
		d.mu.Lock()
		isLoaded := d.state == loaded
		d.mu.Unlock()
		if isLoaded && slices.Contains(d.deps, targetPath) {
			return true
		}
	}
	return false
}

// Remove unloads the given kernel module and any unused dependencies using the
// default ProbeOpts.
func Remove(name string) error {
	return RemoveOptions(name, ProbeOpts{})
}

// RemoveOptions unloads the given kernel module and any unused dependencies.
func RemoveOptions(name string, opts ProbeOpts) error {
	p, err := NewProber(opts)
	if err != nil {
		return err
	}
	return p.Remove(name)
}

func indexModPath(modPath string, modIndex map[string]string) {
	if modIndex == nil {
		return
	}
	if stem, ok := modNameFromPath(modPath); ok {
		norm := normalizeModName(stem)
		if _, exists := modIndex[norm]; !exists {
			modIndex[norm] = modPath
		}
	}
}

func checkBuiltin(moduleDir string, deps depMap, modIndex map[string]string) error {
	f, err := os.Open(filepath.Join(moduleDir, "modules.builtin"))
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("could not open builtin file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := strings.TrimSpace(scanner.Text())
		if txt == "" || strings.HasPrefix(txt, "#") {
			continue
		}
		modPath := filepath.Join(moduleDir, txt)
		if deps[modPath] == nil {
			deps[modPath] = new(dependency)
		}
		deps[modPath].state = builtin
		indexModPath(modPath, modIndex)
	}

	return scanner.Err()
}

func genDeps(opts ProbeOpts) (depMap, map[string]string, string, error) {
	deps := make(depMap)
	modIndex := make(map[string]string)
	rel := opts.KVer

	if rel == "" {
		var u unix.Utsname
		if err := unix.Uname(&u); err != nil {
			return nil, nil, "", fmt.Errorf("could not get release (uname -r): %w", err)
		}
		rel = unix.ByteSliceToString(u.Release[:])
	}

	var moduleDir string
	for _, n := range []string{"/lib/modules", "/usr/lib/modules"} {
		moduleDir = filepath.Join(opts.RootDir, n, strings.TrimSpace(rel))
		if _, err := os.Stat(moduleDir); err == nil {
			break
		}
	}

	f, err := os.Open(filepath.Join(moduleDir, "modules.dep"))
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not open dependency file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	for scanner.Scan() {
		txt := strings.TrimSpace(scanner.Text())
		if txt == "" || strings.HasPrefix(txt, "#") {
			continue
		}
		modPath, modDeps, ok := strings.Cut(txt, ":")
		if !ok {
			continue
		}
		modPath = filepath.Join(moduleDir, strings.TrimSpace(modPath))

		var dep dependency
		for d := range strings.FieldsSeq(modDeps) {
			dep.deps = append(dep.deps, filepath.Join(moduleDir, d))
		}
		deps[modPath] = &dep
		indexModPath(modPath, modIndex)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, "", err
	}

	if err = checkBuiltin(moduleDir, deps, modIndex); err != nil {
		return nil, nil, "", err
	}

	if !opts.IgnoreProcMods {
		procModsPath := filepath.Join(opts.RootDir, "/proc/modules")
		if fm, err := os.Open(procModsPath); err == nil {
			defer fm.Close()
			_ = genLoadedModsWithIndex(fm, deps, modIndex)
		}
	}

	return deps, modIndex, moduleDir, nil
}

// modNameFromPath strips a known kernel module extension (.ko, .ko.gz, .ko.xz,
// .ko.zst) from the base name of mp and returns the module stem.
func modNameFromPath(mp string) (string, bool) {
	base := path.Base(mp)
	for _, ext := range []string{".ko.gz", ".ko.xz", ".ko.zst", ".ko"} {
		if stem, ok := strings.CutSuffix(base, ext); ok && stem != "" {
			return stem, true
		}
	}
	return "", false
}

func normalizeModName(name string) string {
	return strings.ReplaceAll(strings.TrimSpace(name), "-", "_")
}

// normalizeAlias replaces '-' with '_' outside of '[...]' character ranges so
// that module aliases and patterns match consistently while preserving
// character ranges such as [0-9].
func normalizeAlias(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	inBracket := false
	escaped := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if escaped {
			b.WriteByte(c)
			escaped = false
			continue
		}
		switch c {
		case '\\':
			escaped = true
			b.WriteByte(c)
		case '[':
			inBracket = true
			b.WriteByte(c)
		case ']':
			inBracket = false
			b.WriteByte(c)
		case '-':
			if !inBracket {
				b.WriteByte('_')
			} else {
				b.WriteByte('-')
			}
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

// matchAlias matches a module alias against a shell glob pattern. Unlike file
// paths, module aliases (such as DMI BIOS dates) may contain '/', which '*'
// and '?' must match just like fnmatch(pattern, name, 0).
func matchAlias(pattern, name string) bool {
	p := strings.ReplaceAll(normalizeAlias(pattern), "/", "\x01")
	n := strings.ReplaceAll(normalizeAlias(name), "/", "\x01")
	matched, err := path.Match(p, n)
	return err == nil && matched
}

func findModPath(name string, m depMap) (string, error) {
	// Kernel modules do not have consistency between hyphens and
	// underscores in the module name versus the file name (and a single
	// file name may contain both, e.g. twofish-x86_64.ko). Normalize both
	// sides to '_'.
	target := normalizeModName(name)
	if target == "" {
		return "", fmt.Errorf("could not find path for module %q", name)
	}

	var matches []string
	for mp := range m {
		if stem, ok := modNameFromPath(mp); ok && normalizeModName(stem) == target {
			matches = append(matches, mp)
		}
	}
	if len(matches) == 1 {
		return matches[0], nil
	}
	if len(matches) > 1 {
		slices.Sort(matches)
		return matches[0], nil
	}

	return "", fmt.Errorf("could not find path for module %q", name)
}

func loadModule(path, modParams string, opts ProbeOpts) error {
	if opts.DryRunCB != nil {
		opts.DryRunCB(path)
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := FileInit(f, modParams, 0); err != nil && err != unix.EEXIST {
		return err
	}

	return nil
}

func genLoadedMods(r io.Reader, deps depMap) error {
	return genLoadedModsWithIndex(r, deps, nil)
}

func genLoadedModsWithIndex(r io.Reader, deps depMap, modIndex map[string]string) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		arr := strings.Fields(scanner.Text())
		if len(arr) == 0 {
			continue
		}
		// Field 5 in /proc/modules is the module state ("Live",
		// "Loading", or "Unloading"). Treat both "Live" and "Loading"
		// as loaded: if a module's init routine triggers a recursive
		// request_module() call back into modprobe, re-calling
		// finit_module on the "Loading" parent module would deadlock in
		// the kernel's add_unformed_module(). For sibling processes
		// loading a dependent module, resolve_symbol_wait() in the
		// kernel already waits for "Loading" dependencies to finish.
		if len(arr) >= 5 && arr[4] == "Unloading" {
			continue
		}
		name := arr[0]
		var modPath string
		if modIndex != nil {
			modPath = modIndex[normalizeModName(name)]
		}
		if modPath == "" {
			var err error
			modPath, err = findModPath(name, deps)
			if err != nil {
				// Ignore out-of-tree modules not in modules.dep
				// so the remaining lines in /proc/modules are
				// still processed.
				continue
			}
		}
		if deps[modPath] == nil {
			deps[modPath] = new(dependency)
		}
		deps[modPath].state = loaded
	}
	return scanner.Err()
}

var modprobeConfDirs = []string{
	"/etc/modprobe.d",
	"/run/modprobe.d",
	"/usr/local/lib/modprobe.d",
	"/usr/lib/modprobe.d",
	"/lib/modprobe.d",
}

type confFile struct {
	name string
	path string
}

func parseModprobeConfigs(rootDir string, p *Prober) {
	seen := make(map[string]bool)
	var files []confFile

	for _, dir := range modprobeConfDirs {
		fullDir := filepath.Join(rootDir, dir)
		entries, err := os.ReadDir(fullDir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			name := entry.Name()
			if entry.IsDir() || strings.HasPrefix(name, ".") || !strings.HasSuffix(name, ".conf") {
				continue
			}
			if seen[name] {
				continue
			}
			seen[name] = true
			files = append(files, confFile{
				name: name,
				path: filepath.Join(fullDir, name),
			})
		}
	}

	slices.SortFunc(files, func(a, b confFile) int {
		return strings.Compare(a.name, b.name)
	})

	for _, cf := range files {
		f, err := os.Open(cf.path)
		if err != nil {
			continue
		}
		parseModprobeConfReader(f, p)
		f.Close()
	}
}

func stripComment(line string) string {
	inQuote := false
	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '"':
			inQuote = !inQuote
		case '#':
			if !inQuote {
				return line[:i]
			}
		}
	}
	return line
}

func scanLogicalLines(r io.Reader, fn func(line string)) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	var pending strings.Builder
	for scanner.Scan() {
		line := stripComment(scanner.Text())
		if before, ok := strings.CutSuffix(strings.TrimRight(line, " \t\r"), `\`); ok {
			pending.WriteString(before)
			pending.WriteByte(' ')
			continue
		}
		if pending.Len() > 0 {
			pending.WriteString(line)
			line = pending.String()
			pending.Reset()
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fn(line)
	}
	if pending.Len() > 0 {
		if line := strings.TrimSpace(pending.String()); line != "" {
			fn(line)
		}
	}
}

func parseModprobeConfReader(r io.Reader, p *Prober) {
	scanLogicalLines(r, func(line string) {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return
		}
		switch fields[0] {
		case "alias":
			if len(fields) >= 3 {
				p.configAliases = append(p.configAliases, modAlias{
					pattern: fields[1],
					target:  fields[2],
				})
			}
		case "options":
			if len(fields) >= 3 {
				norm := normalizeModName(fields[1])
				optStr := strings.Join(fields[2:], " ")
				p.confOpts[norm] = append(p.confOpts[norm], optStr)
			}
		case "blacklist":
			norm := normalizeModName(fields[1])
			if norm != "" {
				p.blacklist[norm] = true
			}
		case "softdep":
			if len(fields) >= 3 {
				parseSoftdepFields(fields[1], fields[2:], p.softDeps)
			}
		}
	})
}

func parseModuleSoftdeps(moduleDir string, p *Prober) {
	f, err := os.Open(filepath.Join(moduleDir, "modules.softdep"))
	if err != nil {
		return
	}
	defer f.Close()

	scanLogicalLines(f, func(line string) {
		fields := strings.Fields(line)
		if len(fields) >= 3 && fields[0] == "softdep" {
			parseSoftdepFields(fields[1], fields[2:], p.softDeps)
		}
	})
}

func parseSoftdepFields(modName string, fields []string, softDeps map[string]softDep) {
	norm := normalizeModName(modName)
	if norm == "" {
		return
	}
	var sd softDep
	mode := "pre"
	for _, f := range fields {
		switch f {
		case "pre:":
			mode = "pre"
		case "post:":
			mode = "post"
		default:
			if after, ok := strings.CutPrefix(f, "pre:"); ok {
				mode = "pre"
				if after != "" {
					sd.pre = append(sd.pre, after)
				}
			} else if after, ok := strings.CutPrefix(f, "post:"); ok {
				mode = "post"
				if after != "" {
					sd.post = append(sd.post, after)
				}
			} else if mode == "post" {
				sd.post = append(sd.post, f)
			} else {
				sd.pre = append(sd.pre, f)
			}
		}
	}
	softDeps[norm] = sd
}

func parseModuleAliases(moduleDir, fileName string, p *Prober) {
	f, err := os.Open(filepath.Join(moduleDir, fileName))
	if err != nil {
		return
	}
	defer f.Close()

	scanLogicalLines(f, func(line string) {
		fields := strings.Fields(line)
		if len(fields) >= 3 && fields[0] == "alias" {
			p.moduleAliases = append(p.moduleAliases, modAlias{
				pattern: fields[1],
				target:  fields[2],
			})
		}
	})
}

func parseBuiltinModinfo(moduleDir string, p *Prober) {
	data, err := os.ReadFile(filepath.Join(moduleDir, "modules.builtin.modinfo"))
	if err != nil {
		return
	}

	seenBuiltin := make(map[string]bool)
	for entry := range bytes.SplitSeq(data, []byte{0}) {
		line := strings.TrimSpace(string(entry))
		if line == "" {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		modName, field, ok := strings.Cut(key, ".")
		if !ok || modName == "" {
			continue
		}
		norm := normalizeModName(modName)
		if !seenBuiltin[norm] {
			seenBuiltin[norm] = true
			// Ensure the builtin module is registered in deps even
			// if it was omitted from modules.builtin.
			if _, err := p.lookupModPath(modName); err != nil {
				modPath := filepath.Join(moduleDir, modName+".ko")
				p.deps[modPath] = &dependency{state: builtin}
				indexModPath(modPath, p.modIndex)
			}
		}
		if field == "alias" && val != "" {
			p.moduleAliases = append(p.moduleAliases, modAlias{
				pattern: val,
				target:  modName,
			})
		}
	}
}

func parseProcCmdline(rootDir string, p *Prober) {
	cmdlinePath := filepath.Join(rootDir, "/proc/cmdline")
	data, err := os.ReadFile(cmdlinePath)
	if err != nil {
		return
	}
	parseCmdlineString(string(data), p.cmdlineOpts, p.blacklist)
}

func isValidModIdent(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		isAlpha := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
		isDigit := c >= '0' && c <= '9'
		if !isAlpha && !isDigit && c != '_' && c != '-' {
			return false
		}
	}
	return true
}

func parseCmdlineString(raw string, cmdlineOpts map[string][]string, blacklist map[string]bool) {
	for _, tok := range splitCmdlineTokens(raw) {
		if tok == "--" {
			break
		}
		key, val, hasEq := strings.Cut(tok, "=")
		if normalizeModName(key) == "modprobe.blacklist" {
			if hasEq {
				for m := range strings.SplitSeq(val, ",") {
					if norm := normalizeModName(m); norm != "" {
						blacklist[norm] = true
					}
				}
			}
			continue
		}
		modName, paramName, hasDot := strings.Cut(key, ".")
		if !hasDot || !isValidModIdent(modName) || !isValidModIdent(paramName) {
			continue
		}
		normMod := normalizeModName(modName)
		if hasEq {
			cmdlineOpts[normMod] = append(cmdlineOpts[normMod], paramName+"="+val)
		} else {
			cmdlineOpts[normMod] = append(cmdlineOpts[normMod], paramName)
		}
	}
}

func splitCmdlineTokens(raw string) []string {
	var tokens []string
	var cur strings.Builder
	inQuote := false
	for i := 0; i < len(raw); i++ {
		c := raw[i]
		if c == '"' {
			inQuote = !inQuote
			cur.WriteByte(c)
			continue
		}
		if !inQuote && (c == ' ' || c == '\t' || c == '\n' || c == '\r') {
			if cur.Len() > 0 {
				tokens = append(tokens, cur.String())
				cur.Reset()
			}
			continue
		}
		cur.WriteByte(c)
	}
	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}
	return tokens
}
