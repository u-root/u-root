// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
)

func test(
	ctx context.Context,
	goVer []string,
	env map[string]string,
	client *dagger.Client,
	opts *TestOpts,
) error {
	client = client.Pipeline("test")
	hostSourceDir := client.Host().Directory(".")

	for _, toolchain := range strings.Split(*opts.Toolchain, ",") {
		for _, version := range goVer {
			var goImage string

			switch toolchain {
			case "std":
				goImage = fmt.Sprintf("golang:%s-alpine", version)
			case "tamago":
				goImage = "debian:11-slim"
			case "tinygo":
				goImage = "tinygo/tinygo:latest" // for now we always choose the latest release
			default:
				continue
			}

			source := client.Container().
				From(goImage).
				WithMountedDirectory("/src", hostSourceDir).
				WithWorkdir("/src")

			artifacts := client.Directory()

			if toolchain == "tamago" {
				source = source.
					WithExec([]string{"tar", "-xvf", "tamago-go1.20.4.linux-amd64.tar.gz", "-C", "/"}).
					WithEnvVariable("TAMAGO", "/usr/local/tamago-go/bin/go")
			}

			path := fmt.Sprintf("test/%s/", version)
			testCmd := []string{"go", "test", "-v", "-failfast", "-timeout=15m"}

			if *opts.Cover {
				source = source.
					WithEnvVariable("UROOT_QEMU_COVERPROFILE", "vmcoverage.txt").
					WithEnvVariable("CGO_ENABLED", "0")
				testCmd = append(testCmd, "-a", "-cover", "-covermode=atomic", "-coverprofile", filepath.Join(path, "coverage.txt"))
			}

			if *opts.Race {
				source = source.
					WithEnvVariable("CGO_ENABLED", "1").
					WithExec([]string{"apk", "add", "build-base"})
				testCmd = append(testCmd, "-race")
			}

			testCmd = append(testCmd, "./cmds/...", "./pkg/...")

			if toolchain == "tinygo" {
				source = source.WithExec([]string{"sudo", "chown", "-R", "tinygo:tinygo", "."})
			}

			test := source.
				WithExec([]string{"mkdir", "-p", path}).
				WithExec(testCmd)

			if *opts.Integration != "" {
				for _, integrationArch := range strings.Split(*opts.Integration, ",") {
					path = fmt.Sprintf("%sintegration/%s/", path, integrationArch)
					integrationTest := source.
						WithEnvVariable("UROOT_QEMU_COVERPROFILE", "coverage.txt").
						WithEnvVariable("UROOT_KERNEL", "/boot/vmlinuz-virt").
						WithEnvVariable("UROOT_QEMU", "").
						WithEnvVariable("UROOT_TESTARCH", integrationArch).
						WithExec([]string{"apk", "add", "linux-virt", "qemu"}).
						WithExec([]string{"go", "test", "-v", "-a", "-timeout=20m", "-failfast", "./integration/..."})
					integrationArtifacts := artifacts.WithDirectory(path, integrationTest.Directory(path))

					_, err := integrationArtifacts.Export(ctx, ".")
					if err != nil {
						return fmt.Errorf("%w: %v", ErrExportFromPipeline, err)
					}
				}
			}

			artifacts = artifacts.WithDirectory(path, test.Directory(path))

			_, err := artifacts.Export(ctx, ".")
			if err != nil {
				return fmt.Errorf("%w: %v", ErrExportFromPipeline, err)
			}
		}
	}

	return nil
}
