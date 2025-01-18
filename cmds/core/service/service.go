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
	"os/exec"
	"path/filepath"
)

const (
	serviceDir = "/etc/init.d"
)

var (
	doFullRestart bool
	doStatusAll   bool
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: service <SCRIPT> <COMMAND> [OPTIONS]\n")
		flag.PrintDefaults()
	}
	flag.BoolVar(&doFullRestart, "full-restart", false, "Restart all services")
	flag.BoolVar(&doStatusAll, "status-all", false, "Display the status of all services")
	flag.Parse()
	ctx := context.Background()

	if err := run(ctx, flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if doFullRestart {
		return fullRestart(ctx, serviceDir)
	}

	if doStatusAll {
		return statusAll(ctx, serviceDir)
	}

	if len(args) < 2 {
		return errors.New("not enough args: service and command must be specified")
	}

	service := args[0]
	command := args[1]
	extraArgs := args[2:]

	if service == "" || command == "" {
		return errors.New("service and command must be specified")
	}

	return execute(ctx, filepath.Join(serviceDir, service), command, extraArgs...)
}

func statusAll(ctx context.Context, serviceDir string) error {
	entries, err := os.ReadDir(serviceDir)
	if err != nil {
		return fmt.Errorf("failed to read services directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		servicePath := filepath.Join(serviceDir, entry.Name())
		if err := execute(ctx, servicePath, "status"); err != nil {
			fmt.Fprintf(os.Stderr, "failed to check status of service %s: %v\n", entry.Name(), err)
		}
	}
	return nil
}

func fullRestart(ctx context.Context, serviceDir string) error {
	entries, err := os.ReadDir(serviceDir)
	if err != nil {
		return fmt.Errorf("failed to read services directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		servicePath := filepath.Join(serviceDir, entry.Name())
		if err := execute(ctx, servicePath, "restart"); err != nil {
			fmt.Fprintf(os.Stderr, "failed to restart service %s: %v\n", entry.Name(), err)
		}
	}
	return nil
}

func execute(ctx context.Context, service, command string, args ...string) error {
	argv := []string{command}
	argv = append(argv, args...)
	cmd := exec.CommandContext(ctx, service, argv...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute service script: %w", err)
	}
	return nil
}
