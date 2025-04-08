// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mkuimage

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/mkuimage/uimage"
	"github.com/u-root/mkuimage/uimage/builder"
	"github.com/u-root/mkuimage/uimage/templates"
	"github.com/u-root/uio/llog"
)

var recommendedVersions = []string{
	"go1.20",
	"go1.21",
	"go1.22",
}

func isRecommendedVersion(v string) bool {
	for _, r := range recommendedVersions {
		if strings.HasPrefix(v, r) {
			return true
		}
	}
	return false
}

func uimageOpts(l *llog.Logger, m []uimage.Modifier, tpl *templates.Templates, f *Flags, conf string, cmds []string) (*uimage.Opts, error) {
	// Evaluate template first -- template settings may always be
	// appended/overridden by further flag-based settings.
	if conf != "" {
		mods, err := tpl.Uimage(conf)
		if err != nil {
			return nil, err
		}
		l.Debugf("Config: %#v", tpl.Configs[conf])
		m = append(m, mods...)
	}

	// Expand command templates.
	if tpl != nil {
		cmds = tpl.CommandsFor(cmds...)
	}

	more, err := f.Modifiers(cmds...)
	if err != nil {
		return nil, err
	}
	return uimage.OptionsFor(append(m, more...)...)
}

func checkAmd64Level(l *llog.Logger, env *golang.Environ) {
	if env.GOARCH != "amd64" {
		return
	}

	// Looking for "amd64.v2" in "env.ToolTags" is unreliable; see
	// <https://github.com/golang/go/issues/72791>. Invoke "go env" instead.
	var bad string
	abiLevel, err := exec.Command("go", "env", "GOAMD64").Output()
	if err == nil {
		if bytes.Equal(abiLevel, []byte("v1\n")) {
			return
		}
		bad = "is not"
	} else {
		if exerr, isExErr := err.(*exec.ExitError); isExErr {
			l.Warnf("\"go env\" failed: %s", exerr.Stderr)
		} else {
			l.Warnf("couldn't execute \"go env\": %s", err)
		}
		bad = "may not be"
	}
	l.Warnf("GOAMD64 %s set to v1; on older CPUs, binaries built into " +
		"the initrd may crash or refuse to run.", bad)
}

// CreateUimage creates a uimage with the given base modifiers and flags, using args as the list of commands.
func CreateUimage(l *llog.Logger, base []uimage.Modifier, tf *TemplateFlags, f *Flags, args []string) error {
	tpl, err := tf.Get()
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	keepTempDir := f.KeepTempDir
	if f.TempDir == nil {
		tempDir, err := os.MkdirTemp("", "u-root")
		if err != nil {
			return err
		}
		f.TempDir = &tempDir
		defer func() {
			if keepTempDir {
				l.Infof("Keeping temp dir %s", tempDir)
			} else {
				os.RemoveAll(tempDir)
			}
		}()
	} else if _, err := os.Stat(*f.TempDir); os.IsNotExist(err) {
		if err := os.MkdirAll(*f.TempDir, 0o755); err != nil {
			return fmt.Errorf("temporary directory %q did not exist; tried to mkdir but failed: %v", *f.TempDir, err)
		}
	}

	opts, err := uimageOpts(l, base, tpl, f, tf.Config, args)
	if err != nil {
		return err
	}

	env := opts.Env

	l.Infof("Build environment: %s", env)
	if env.GOOS != "linux" {
		l.Warnf("GOOS is not linux. Did you mean to set GOOS=linux?")
	}

	checkAmd64Level(l, env);

	v, err := env.Version()
	if err != nil {
		l.Infof("Could not get environment's Go version, using runtime's version: %v", err)
		v = runtime.Version()
	}
	if !isRecommendedVersion(v) {
		l.Warnf(`You are not using one of the recommended Go versions (have = %s, recommended = %v).
			Some packages may not compile.
			Go to https://golang.org/doc/install to find out how to install a newer version of Go,
			or use https://godoc.org/golang.org/dl/%s to install an additional version of Go.`,
			v, recommendedVersions, recommendedVersions[0])
	}

	err = opts.Create(l)
	if errors.Is(err, builder.ErrBusyboxFailed) {
		l.Errorf("Preserving temp dir due to busybox build error")
		keepTempDir = true
	}
	return err
}
