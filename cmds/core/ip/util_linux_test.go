// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"testing"
)

// TestPrintJSON tests the printJSON function with different scenarios.
func TestPrintJSON(t *testing.T) {
	tests := []struct {
		name    string
		cmd     cmd
		data    VrfJSON
		want    string
		wantErr bool
	}{
		{
			name: "With Prettify",
			cmd: cmd{
				Opts: flags{
					Prettify: true,
				},
				Out: &bytes.Buffer{},
			},
			data:    VrfJSON{Name: "Test", Table: 2},
			want:    "{\n    \"name\": \"Test\",\n    \"table\": 2\n}",
			wantErr: false,
		},
		{
			name: "Without Prettify",
			cmd: cmd{
				Opts: flags{
					Prettify: false,
				},
				Out: &bytes.Buffer{},
			},
			data:    VrfJSON{Name: "Test", Table: 2},
			want:    "{\"name\":\"Test\",\"table\":2}",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := printJSON(tt.cmd, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("printJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got := tt.cmd.Out.(*bytes.Buffer).String(); got != tt.want {
				t.Errorf("printJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}
