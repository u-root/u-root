// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

var ErrUnknownTarget = errors.New("unknown build target")

func build(
	ctx context.Context,
	goVer []string,
	env map[string]string,
	client *dagger.Client,
	opts *BuildOpts,
) error {
	client = client.Pipeline("build")
	hostSourceDir := client.Host().Directory(".")

	for _, toolchain := range strings.Split(*opts.Toolchain, ",") {
		for _, version := range goVer {
			source := client.Container().
				From(fmt.Sprintf("golang:%s-alpine", version)).
				WithMountedDirectory("/src", hostSourceDir).
				WithWorkdir("/src").
				WithExec([]string{"apk", "update"}).
				WithExec([]string{"apk", "add", "bash", "build-base", "gcc-cross-embedded"})

			artifacts := client.Directory()

			for _, goos := range strings.Split(*opts.Platform, ",") {
				for _, goarch := range strings.Split(*opts.Arch, ",") {
					path := fmt.Sprintf("build/%s/%s/%s/", version, goos, goarch)
					build := source.WithEnvVariable("GOOS", goos)
					build = build.WithEnvVariable("GOARCH", goarch)

					switch *opts.Target {
					case "all":
						build = build.WithExec([]string{"go", "build", "-v", "-o", path, "./cmds/..."})

					case "u-root":
						build = build.WithExec([]string{"go", "build", "-v", "-o", path})
					case "templates":
						build = build.WithExec([]string{"go", "build", "-v"})
						runCmd := "./u-root -stats-output-path=stats.json"
						if goos == "plan9" {
							build = build.WithExec([]string{runCmd, "plan9"})
						} else {
							build = build.
								WithExec([]string{runCmd, "minimal"}).
								WithExec([]string{runCmd, "core"}).
								WithExec([]string{runCmd, "coreboot-app"}).
								WithExec([]string{runCmd, "all"}).
								WithExec([]string{runCmd, "world"}).
								WithExec([]string{runCmd, "all"}).
								WithExec([]string{"cat", "stats.json"})
						}

					default:
						return fmt.Errorf("%w: %v", ErrUnknownTarget, *opts.Target)
					}

					artifacts = artifacts.WithDirectory(path, build.Directory(path))
				}
			}

			_, err := artifacts.Export(ctx, ".")
			if err != nil {
				return fmt.Errorf("%w: %v", ErrExportFromPipeline, err)
			}
		}
	}

	return nil
}
