// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"errors"
	"os"
	"testing"
	"time"
)

func TestParseParamsDate(t *testing.T) {
	date := "2021-01-01T00:00:00Z"
	expected, err := time.Parse(time.RFC3339, date)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	p, err := parseParams(date, false, false, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !expected.Equal(p.time) {
		t.Errorf("expected %v, got %v", expected, p.time)
	}

	date = "invalid"
	_, err = parseParams(date, false, false, false)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestParseParams(t *testing.T) {
	var tests = []struct {
		expected     params
		access       bool
		modification bool
		create       bool
	}{
		{
			access:       false,
			modification: false,
			create:       false,
			expected: params{
				access:       true,
				modification: true,
				create:       false,
			},
		},
		{
			access:       true,
			modification: false,
			create:       false,
			expected: params{
				access:       true,
				modification: false,
				create:       false,
			},
		},
		{
			access:       false,
			modification: true,
			create:       true,
			expected: params{
				access:       false,
				modification: true,
				create:       true,
			},
		},
	}

	for _, test := range tests {
		p, err := parseParams("", test.access, test.modification, test.create)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if p.access != test.expected.access {
			t.Errorf("expected %v, got %v", test.expected.access, p.access)
		}
		if p.modification != test.expected.modification {
			t.Errorf("expected %v, got %v", test.expected.modification, p.modification)
		}
		if p.create != test.expected.create {
			t.Errorf("expected %v, got %v", test.expected.create, p.create)
		}
	}
}

var tests = []struct {
	err  error
	p    params
	name string
	args []string
}{
	{
		name: "create is true, no new files created",
		args: []string{"a1", "a2"},
		p: params{
			access:       true,
			modification: true,
			create:       true,
			time:         time.Now(),
		},
	},
	{
		name: "create is false, files should be created",
		args: []string{"a1", "a2"},
		p: params{
			access:       true,
			modification: true,
			create:       false,
			time:         time.Now(),
		},
	},
	{
		name: "no such file or directory",
		args: []string{"no/such/file/or/direcotry"},
		p: params{
			create: false,
			time:   time.Now(),
		},
		err: os.ErrNotExist,
	},
}

func TestTouchEmptyDir(t *testing.T) {
	for _, test := range tests {
		temp := t.TempDir()
		var args []string
		for _, arg := range test.args {
			args = append(args, temp+arg)
		}
		err := command(test.p, args...).run()
		if !errors.Is(err, test.err) {
			t.Fatalf("command() expected %v, got %v", test.err, err)
		}
		if test.err != nil {
			continue
		}

		for _, arg := range args {
			_, err := os.Stat(arg)
			if test.p.create {
				if !os.IsNotExist(err) {
					t.Errorf("expected %s to not exist", arg)
				}
			} else {
				if err != nil {
					t.Errorf("expected %s to exist, got %v", arg, err)
				}

				stat, err := os.Stat(arg)
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if test.p.modification {
					if stat.ModTime().Unix() != test.p.time.Unix() {
						t.Errorf("expected %s to have mod time %v, got %v", arg, test.p.time, stat.ModTime())
					}
				}
			}

		}
	}
}
