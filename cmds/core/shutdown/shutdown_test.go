// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"
	"testing"
)

func TestShutdown(t *testing.T) {
	for _, tt := range []struct {
		name    string
		args    []string
		dryrun  bool
		want    uint
		wantErr string
	}{
		{
			name:   "halt",
			args:   []string{},
			dryrun: true,
			want:   1126301404,
		},
		{
			name:   "halt +-0",
			args:   []string{"halt", "+-0"},
			dryrun: true,
			want:   1126301404,
		},
		{
			name:   "halt +0",
			args:   []string{"halt", "+0"},
			dryrun: true,
			want:   1126301404,
		},
		{
			name:    "halt +a",
			args:    []string{"halt", "+a"},
			dryrun:  true,
			wantErr: "invalid duration",
		},
		{
			name:   "halt +1",
			args:   []string{"halt", "+1"},
			dryrun: true,
			want:   1126301404,
		},
		{
			name:   "halt now",
			args:   []string{"halt", "now"},
			dryrun: true,
			want:   1126301404,
		},
		{
			name:   "halt specific date",
			args:   []string{"halt", "2006-01-02T15:04:05Z"},
			dryrun: true,
			want:   1126301404,
		},
		{
			name:    "halt specific date",
			args:    []string{"halt", "2006-01-02T15:04:05Z07:00"},
			dryrun:  true,
			want:    1126301404,
			wantErr: "extra text",
		},
		{
			name:    "halt specific date",
			args:    []string{"halt", "2006-o1-02T15:04:05Z07:00"},
			dryrun:  true,
			want:    1126301404,
			wantErr: "cannot parse",
		},
		{
			name:    "halt police",
			args:    []string{"halt", "police"},
			dryrun:  true,
			wantErr: "cannot parse",
		},
		{
			name:   "-h",
			args:   []string{"-h"},
			dryrun: true,
			want:   1126301404,
		},
		{
			name:   "empty string = halt",
			args:   []string{""},
			dryrun: true,
			want:   0,
		},
		{
			name:   "reboot",
			args:   []string{"reboot"},
			dryrun: true,
			want:   19088743,
		},
		{
			name:   "-r",
			args:   []string{"-r"},
			dryrun: true,
			want:   19088743,
		},
		{
			name:   "suspend",
			args:   []string{"suspend"},
			dryrun: true,
			want:   3489725666,
		},
		{
			name:   "-s",
			args:   []string{"-s"},
			dryrun: true,
			want:   3489725666,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := shutdown(tt.args, tt.dryrun)
			if err != nil {
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("shutdown() = %q, want to contain: %q", err.Error(), tt.wantErr)
				}
			} else {
				if got != tt.want {
					t.Errorf("shutdown() = '%d', want: '%d'", got, tt.want)
				}
			}
		})
	}
}
