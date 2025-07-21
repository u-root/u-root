// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/lsb"
)

const (
	serviceDir = "/etc/init.d"
	etc        = "/etc"
)

var commands = map[string]subcommand{
	"defaults": {
		usage:   "update-rc.d <SCRIPT> defaults",
		handler: defaults,
	},
	"defaults-disable": {
		usage:   "update-rc.d <SCRIPT> defaults-disabled",
		handler: defaultsDisable,
	},
	"disable": {
		usage:   "update-rc.d <SCRIPT> disable [ S|2|3|4|5 ]",
		handler: disable,
	},
	"enable": {
		usage:   "update-rc.d <SCRIPT> enable [ S|2|3|4|5 ]",
		handler: enable,
	},
	"remove": {
		usage:   "update-rc.d [-f] <SCRIPT> remove",
		handler: remove,
	},
}

var force bool

type options struct {
	etc        string
	serviceDir string
	force      bool
	extraArgs  []string
}

type subcommand struct {
	usage   string
	handler func(ctx context.Context, script string, opts options) error
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n")
		for _, cmd := range commands {
			fmt.Fprintf(flag.CommandLine.Output(), "  %s\n", cmd.usage)
		}
		flag.PrintDefaults()
	}
	flag.BoolVar(&force, "f", false, "Force removal of symlinks even if /etc/init.d/<SCRIPT> still exists.")
	flag.Parse()

	if err := run(context.Background(), flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return errors.New("not enough args: script and command must be specified")
	}

	script := args[0]
	command := args[1]
	extraArgs := args[2:]

	subcmd, ok := commands[command]
	if !ok {
		return fmt.Errorf("unknown command: %q", command)
	}

	return subcmd.handler(ctx, script, options{
		force:      force,
		extraArgs:  extraArgs,
		serviceDir: serviceDir,
		etc:        etc,
	})
}

func readLSBScriptMeta(script string) (*lsb.InitScript, error) {
	f, err := os.Open(script)
	if err != nil {
		return nil, fmt.Errorf("failed to open script file %q: %w", script, err)
	}
	meta := &lsb.InitScript{}
	if err := meta.Unmarshal(f); err != nil {
		return nil, fmt.Errorf("failed to parse LSB metadata: %w", err)
	}
	return meta, nil
}

func parseRunlevels(args []string) []int {
	var levels []int
	for _, arg := range args {
		switch arg {
		case "S":
			levels = append(levels, 0)
		case "2", "3", "4", "5":
			if lvl, err := strconv.Atoi(arg); err == nil {
				levels = append(levels, lvl)
			}
		}
	}
	return levels
}

// defaults makes links named /etc/rc<RUN_LEVEL>.d/[SK]<NN><SCRIPT> that point to the
// script /etc/init.d/<SCRIPT>, using runlevel and dependency information from the
// init.d script LSB comment header.
func defaults(ctx context.Context, script string, opts options) error {
	scriptPath := filepath.Join(opts.serviceDir, script)

	meta, err := readLSBScriptMeta(scriptPath)
	if err != nil {
		return err
	}

	for _, runlevel := range meta.DefaultStart {
		dir := fmt.Sprintf("%s/rc%d.d", opts.etc, runlevel)
		link := filepath.Join(dir, fmt.Sprintf("S%02d%s", meta.SequenceNumber(), script))
		if err := os.Symlink(scriptPath, link); err != nil {
			return fmt.Errorf("failed to create symlink for %q: %w", script, err)
		}
	}
	return nil
}

// defaultsDisable makes links named /etc/rc<RUN_LEVEL>.d/K<NN><SCRIPT> that point to
// the script /etc/init.d/name, using dependency information from the init.d
// script LSB comment header. This means that the init.d script will be disabled.
func defaultsDisable(ctx context.Context, script string, opts options) error {
	scriptPath := filepath.Join(opts.serviceDir, script)

	meta, err := readLSBScriptMeta(scriptPath)
	if err != nil {
		return err
	}

	for _, runlevel := range meta.DefaultStart {
		dir := fmt.Sprintf("%s/rc%d.d", opts.etc, runlevel)
		link := filepath.Join(dir, fmt.Sprintf("K%02d%s", meta.SequenceNumber(), script))
		if err := os.Symlink(scriptPath, link); err != nil {
			return fmt.Errorf("failed to create symlink for %q: %w", script, err)
		}
	}
	return nil
}

// disable [ S|2|3|4|5 ] modifies existing runlevel links for the script
// /etc/init.d/<SCRIPT> by renaming start links to stop links with a sequence
// number equal to the difference of 100 minus the original sequence number.
//
// Only operate on start runlevel links of S, 2, 3, 4 or 5.
// If no start runlevel is specified the script will attempt to modify links
// in all start runlevels.
func disable(ctx context.Context, script string, opts options) error {
	scriptPath := filepath.Join(opts.serviceDir, script)

	meta, err := readLSBScriptMeta(scriptPath)
	if err != nil {
		return err
	}

	levels := parseRunlevels(opts.extraArgs)
	if len(levels) == 0 {
		levels = []int{0, 2, 3, 4, 5}
	}

	for _, runlevel := range levels {
		dir := fmt.Sprintf("%s/rc%d.d", opts.etc, runlevel)
		startLink := filepath.Join(dir, fmt.Sprintf("S%02d%s", meta.SequenceNumber(), script))
		stopLink := filepath.Join(dir, fmt.Sprintf("K%02d%s", 100-meta.SequenceNumber(), script))
		if err := os.Rename(startLink, stopLink); err != nil {
			return fmt.Errorf("failed to disable script for runlevel %d: %w", runlevel, err)
		}
	}
	return nil
}

// enable [ S|2|3|4|5 ] modifies existing runlevel links for the script
// /etc/init.d/<SCRIPT> by renaming stop links to start links with a sequence
// number equal to the positive difference of current sequence number minus 100,
// thus returning to the original sequence number that the script had been installed with
// before disabling it.
//
// Only operate on start runlevel links of S, 2, 3, 4 or 5.
// If no start runlevel is specified the script will attempt to modify links
// in all start runlevels.
func enable(ctx context.Context, script string, opts options) error {
	scriptPath := filepath.Join(opts.serviceDir, script)

	meta, err := readLSBScriptMeta(scriptPath)
	if err != nil {
		return err
	}

	levels := parseRunlevels(opts.extraArgs)
	if len(levels) == 0 {
		levels = []int{0, 2, 3, 4, 5}
	}

	for _, runlevel := range levels {
		dir := fmt.Sprintf("%s/rc%d.d", opts.etc, runlevel)
		stopLink := filepath.Join(dir, fmt.Sprintf("K%02d%s", 100-meta.SequenceNumber(), script))
		startLink := filepath.Join(dir, fmt.Sprintf("S%02d%s", meta.SequenceNumber(), script))
		if err := os.Rename(stopLink, startLink); err != nil {
			return fmt.Errorf("failed to enable script for runlevel %d: %w", runlevel, err)
		}
	}
	return nil
}

// remove removes any links in the /etc/rc<RUN_LEVEL>.d directories to the script
// /etc/init.d/<SCRIPT>. The script must have been deleted already. If the script
// is still present then it aborts with an error message.
//
// This is the only subcommand that respects -f flag.
func remove(ctx context.Context, script string, opts options) error {
	scriptPath := filepath.Join(opts.serviceDir, script)

	if _, err := os.Stat(scriptPath); err == nil && !opts.force {
		return fmt.Errorf("script %q still exists, use -f to force removal", script)
	}

	for runlevel := 0; runlevel <= 6; runlevel++ {
		dir := fmt.Sprintf("%s/rc%d.d", opts.etc, runlevel)
		pattern := fmt.Sprintf("%s%s", strings.Repeat("?", 3), script)

		entries, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			return fmt.Errorf("failed to list links for runlevel %d: %w", runlevel, err)
		}

		for _, entry := range entries {
			if err := os.Remove(entry); err != nil {
				return fmt.Errorf("failed to remove link %q: %w", entry, err)
			}
		}
	}
	return nil
}
