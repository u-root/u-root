// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"dagger.io/dagger"
)

const DEFAULT_GO_VER = "1.19"

var (
	ErrNoPipelineSpecified = errors.New("no pipeline selected")
	ErrConnectEngine       = errors.New("unable to connect to dagger engine")
	ErrBuildPipeline       = errors.New("error in build pipeline")
	ErrTestPipeline        = errors.New("error in test pipeline")
	ErrExportFromPipeline  = errors.New("failed to export artifacts")
)

type Opts struct {
	Verbose *bool
	GoVer   *string
	Env     *string
	Build   *bool
	BOpts   BuildOpts
	Test    *bool
	TOpts   TestOpts
}

type BuildOpts struct {
	Toolchain *string
	Arch      *string
	Platform  *string
	Target    *string
}

type TestOpts struct {
	Toolchain   *string
	Race        *bool
	Cover       *bool
	Integration *string
}

func main() {
	opts := Opts{
		Verbose: flag.Bool("verbose", false, "Activate logging."),
		GoVer:   flag.String("go", DEFAULT_GO_VER, "Comma separated list of Go versions to use in pipelines."),
		Env:     flag.String("env", "", "Comma separated list of environment variables in the form KEY=VALUE."),
		Build:   flag.Bool("build", false, "Run build pipeline."),
		BOpts: BuildOpts{
			Toolchain: flag.String("buildtoolchain", "std", "Comma separated list of toolchains (std, tamago, tinygo)"),
			Arch:      flag.String("arch", "", "Comma separated list of architectures."),
			Platform:  flag.String("platform", "", "Comma separated list of platforms."),
			Target:    flag.String("target", "all", "Choose from 'all', 'u-root' and 'templates'."),
		},
		Test: flag.Bool("test", false, "Run test pipeline."),
		TOpts: TestOpts{
			Toolchain:   flag.String("testtoolchain", "std", "Comma separated list of toolchains (std, tamago, tinygo)"),
			Race:        flag.Bool("race", false, "Run race condition checks."),
			Cover:       flag.Bool("cover", false, "Generate coverage reports."),
			Integration: flag.String("integration", "", "Comma separated list of architectures to run integration tests for."),
		},
	}

	flag.Parse()

	if err := run(&opts); err != nil {
		log.Fatalln(err)
	}
}

func run(opts *Opts) error {
	if !*opts.Build && !*opts.Test {
		return ErrNoPipelineSpecified
	}

	env := make(map[string]string)
	if *opts.Env != "" {
		for _, e := range strings.Split(*opts.Env, ",") {
			kv := strings.SplitN(e, "=", 2)
			env[kv[0]] = kv[1]
		}
	}

	goVer := strings.Split(*opts.GoVer, ",")

	ctx := context.Background()

	var logWriter io.Writer
	if *opts.Verbose {
		logWriter = os.Stdout
	} else {
		logWriter = io.Discard
	}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(logWriter))
	if err != nil {
		return fmt.Errorf("%w: %v", ErrConnectEngine, err)
	}
	defer client.Close()

	if *opts.Build {
		if err := build(ctx, goVer, env, client, &opts.BOpts); err != nil {
			return fmt.Errorf("%w: %v", ErrBuildPipeline, err)
		}
	}

	if *opts.Test {
		if err := test(ctx, goVer, env, client, &opts.TOpts); err != nil {
			return fmt.Errorf("%w: %v", ErrTestPipeline, err)
		}
	}

	return nil
}
