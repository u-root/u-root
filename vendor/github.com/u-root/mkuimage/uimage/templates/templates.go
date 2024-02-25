// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package templates defines a uimage template configuration file parser.
package templates

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/mkuimage/uimage"
	"gopkg.in/yaml.v3"
)

// ErrTemplateNotExist is returned when the given config name did not exist.
var ErrTemplateNotExist = errors.New("config template does not exist")

func findConfigFile(name string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for dir != "/" {
		p := filepath.Join(dir, name)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
		dir = filepath.Dir(dir)
	}
	return "", fmt.Errorf("%w: could not find %s in current directory or any parent", os.ErrNotExist, name)
}

// Command represents commands to build.
type Command struct {
	// Builder is bb, gbb, or binary.
	//
	// Defaults to bb if not given.
	Builder string

	// Commands are commands or template names.
	Commands []string
}

// Config is a mkuimage build configuration.
type Config struct {
	GOOS      string
	GOARCH    string
	BuildTags []string `yaml:"build_tags"`
	Commands  []Command
	Files     []string
	Init      *string
	Uinit     *string
	Shell     *string
}

// Templates are a set of mkuimage build configs and command templates.
type Templates struct {
	Configs map[string]Config

	// Commands defines a set of command template name -> commands to expand.
	Commands map[string][]string
}

// Uimage returns the uimage modifiers for the given templated config name.
func (t *Templates) Uimage(config string) ([]uimage.Modifier, error) {
	if config == "" {
		return nil, nil
	}
	if t == nil {
		return nil, fmt.Errorf("%w: no templates parsed", ErrTemplateNotExist)
	}
	c, ok := t.Configs[config]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrTemplateNotExist, config)
	}
	m := []uimage.Modifier{
		uimage.WithFiles(c.Files...),
		uimage.WithEnv(
			golang.WithGOOS(c.GOOS),
			golang.WithGOARCH(c.GOARCH),
			golang.WithBuildTag(c.BuildTags...),
		),
	}
	if c.Init != nil {
		m = append(m, uimage.WithInit(*c.Init))
	}
	if c.Uinit != nil {
		m = append(m, uimage.WithUinitCommand(*c.Uinit))
	}
	if c.Shell != nil {
		m = append(m, uimage.WithShell(*c.Shell))
	}
	for _, cmds := range c.Commands {
		switch cmds.Builder {
		case "binary":
			m = append(m, uimage.WithBinaryCommands(t.CommandsFor(cmds.Commands...)...))
		case "bb", "gbb":
			fallthrough
		default:
			m = append(m, uimage.WithBusyboxCommands(t.CommandsFor(cmds.Commands...)...))
		}
	}
	return m, nil
}

// CommandsFor expands commands according to command templates.
func (t *Templates) CommandsFor(names ...string) []string {
	if t == nil {
		return names
	}
	var c []string
	for _, name := range names {
		cmds, ok := t.Commands[name]
		if ok {
			c = append(c, cmds...)
		} else {
			c = append(c, name)
		}
	}
	return c
}

// TemplateFrom parses a template from bytes.
func TemplateFrom(b []byte) (*Templates, error) {
	var tpl Templates
	if err := yaml.Unmarshal(b, &tpl); err != nil {
		return nil, err
	}
	return &tpl, nil
}

// Template parses the first file named .mkuimage.yaml in the current directory or any of its parents.
func Template() (*Templates, error) {
	p, err := findConfigFile(".mkuimage.yaml")
	if err != nil {
		return nil, fmt.Errorf("%w: no templates found", os.ErrNotExist)
	}
	return TemplateFromFile(p)
}

// TemplateFromFile parses a template from the given file.
func TemplateFromFile(p string) (*Templates, error) {
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return TemplateFrom(b)
}
