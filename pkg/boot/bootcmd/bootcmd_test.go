// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bootcmd

import (
	"errors"
	"slices"
	"testing"

	"github.com/u-root/u-root/pkg/boot/menu"
)

type entry struct {
	label string
	exec  error
	execs *[]string
}

func (e entry) Label() string                    { return e.label }
func (e entry) Edit(func(cmdline string) string) {}
func (e entry) Load() error                      { return nil }
func (e entry) IsDefault() bool                  { return true }

func (e entry) Exec() error {
	*e.execs = append(*e.execs, e.label)
	return e.exec
}

// loader returns the given entries one per call, then nil.
func loader(chosen ...menu.Entry) func(bool, ...menu.Entry) menu.Entry {
	return func(bool, ...menu.Entry) menu.Entry {
		if len(chosen) == 0 {
			return nil
		}
		e := chosen[0]
		chosen = chosen[1:]
		return e
	}
}

// shell is an entry that returns to the menu.
type shell struct {
	entry
}

func (shell) ReturnsToMenu() {}

func TestShowMenuAndBootReturnsToMenu(t *testing.T) {
	var execs []string
	sh := shell{entry{label: "shell", execs: &execs}}

	err := showMenuAndBoot(loader(sh, sh), nil, nil, false)
	if err == nil {
		t.Errorf("showMenuAndBoot() = nil, want an error once nothing is chosen")
	}
	if want := []string{"shell", "shell"}; !slices.Equal(execs, want) {
		t.Errorf("executed %v, want %v", execs, want)
	}
}

func TestShowMenuAndBootStops(t *testing.T) {
	for _, tt := range []struct {
		name  string
		chose func(execs *[]string) menu.Entry
	}{
		{
			name: "kexec failed",
			chose: func(execs *[]string) menu.Entry {
				return entry{label: "kexec", exec: errors.New("kexec failed"), execs: execs}
			},
		},
		{
			name: "kexec returned",
			chose: func(execs *[]string) menu.Entry {
				return entry{label: "kexec", execs: execs}
			},
		},
		{
			name: "shell failed",
			chose: func(execs *[]string) menu.Entry {
				return shell{entry{label: "shell", exec: errors.New("no shell"), execs: execs}}
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var execs []string
			e := tt.chose(&execs)
			if err := showMenuAndBoot(loader(e, e), nil, nil, false); err == nil {
				t.Errorf("showMenuAndBoot() = nil, want an error")
			}
			if want := []string{e.Label()}; !slices.Equal(execs, want) {
				t.Errorf("executed %v, want %v", execs, want)
			}
		})
	}
}

func TestShowMenuAndBootNoExec(t *testing.T) {
	var execs []string
	e := entry{label: "kexec", execs: &execs}

	if err := showMenuAndBoot(loader(e), nil, nil, true); err != nil {
		t.Errorf("showMenuAndBoot() = %v, want nil", err)
	}
	if len(execs) != 0 {
		t.Errorf("executed %v with noExec", execs)
	}
}
