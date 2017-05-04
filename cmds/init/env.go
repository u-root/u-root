// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"runtime"
	"strings"
)

type envVar struct {
	name, value string
}

func mkEnv() []envVar {
	var b builder
	b.init()

	env := []envVar{
		{"GOARCH", goarch},
		{"GOBIN", gobin},
		{"GOEXE", exeSuffix},
		{"GOHOSTARCH", runtime.GOARCH},
		{"GOHOSTOS", runtime.GOOS},
		{"GOOS", goos},
		{"GOPATH", buildContext.GOPATH},
		{"GORACE", os.Getenv("GORACE")},
		{"GOROOT", goroot},
		{"GOTOOLDIR", toolDir},

		// disable escape codes in clang errors
		{"TERM", "dumb"},
	}

	if gccgoBin != "" {
		env = append(env, envVar{"GCCGO", gccgoBin})
	} else {
		env = append(env, envVar{"GCCGO", gccgoName})
	}

	switch goarch {
	case "arm":
		env = append(env, envVar{"GOARM", os.Getenv("GOARM")})
	case "386":
		env = append(env, envVar{"GO386", os.Getenv("GO386")})
	}

	cmd := b.gccCmd(".")
	env = append(env, envVar{"CC", cmd[0]})
	env = append(env, envVar{"GOGCCFLAGS", strings.Join(cmd[3:], " ")})
	cmd = b.gxxCmd(".")
	env = append(env, envVar{"CXX", cmd[0]})

	if buildContext.CgoEnabled {
		env = append(env, envVar{"CGO_ENABLED", "1"})
	} else {
		env = append(env, envVar{"CGO_ENABLED", "0"})
	}

	return env
}

func findEnv(env []envVar, name string) string {
	for _, e := range env {
		if e.name == name {
			return e.value
		}
	}
	return ""
}

// extraEnvVars returns environment variables that should not leak into child processes.
func extraEnvVars() []envVar {
	var b builder
	b.init()
	cppflags, cflags, cxxflags, fflags, ldflags := b.cflags(&Package{})
	return []envVar{
		{"PKG_CONFIG", b.pkgconfigCmd()},
		{"CGO_CFLAGS", strings.Join(cflags, " ")},
		{"CGO_CPPFLAGS", strings.Join(cppflags, " ")},
		{"CGO_CXXFLAGS", strings.Join(cxxflags, " ")},
		{"CGO_FFLAGS", strings.Join(fflags, " ")},
		{"CGO_LDFLAGS", strings.Join(ldflags, " ")},
	}
}
