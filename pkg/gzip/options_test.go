// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gzip

import (
	"errors"
	"flag"
	"reflect"
	"runtime"
	"testing"

	"github.com/klauspost/pgzip"
)

func TestOptionsParseArgs(t *testing.T) {
	type args struct {
		cmdLine *flag.FlagSet
		args    []string
	}
	tests := []struct {
		name       string
		args       args
		wantOption Options
		wantErr    error
	}{
		{
			name: "default values no flags",
			args: args{
				cmdLine: flag.NewFlagSet("test", flag.ContinueOnError),
				args:    []string{"gzip", "file.txt"},
			},
			wantOption: Options{
				Blocksize: 128,
				Level:     -1,
				Processes: runtime.NumCPU(),
				Suffix:    ".gz",
			},
		},
		{
			name: "set level 7",
			args: args{
				cmdLine: flag.NewFlagSet("test", flag.ContinueOnError),
				args:    []string{"gzip", "-7", "file.txt"},
			},
			wantOption: Options{
				Blocksize: 128,
				Level:     7,
				Processes: runtime.NumCPU(),
				Suffix:    ".gz",
			},
		},
		{
			name: "stdin/stdout with force",
			args: args{
				cmdLine: flag.NewFlagSet("test", flag.ContinueOnError),
				args:    []string{"gzip", "-f"},
			},
			wantOption: Options{
				Blocksize: 128,
				Level:     -1,
				Processes: runtime.NumCPU(),
				Suffix:    ".gz",
				Force:     true,
				Stdin:     true,
				Stdout:    true,
			},
		},
		{
			name: "with test decompress should be true",
			args: args{
				cmdLine: flag.NewFlagSet("test", flag.ContinueOnError),
				args:    []string{"gzip", "-t", "file.txt"},
			},
			wantOption: Options{
				Blocksize:  128,
				Level:      -1,
				Processes:  runtime.NumCPU(),
				Suffix:     ".gz",
				Test:       true,
				Decompress: true,
			},
		},
		{
			name: "symlink to gunzip",
			args: args{
				cmdLine: flag.NewFlagSet("test", flag.ContinueOnError),
				args:    []string{"gunzip", "file.gz"},
			},
			wantOption: Options{
				Blocksize:  128,
				Level:      -1,
				Processes:  runtime.NumCPU(),
				Suffix:     ".gz",
				Decompress: true,
			},
		},
		{
			name: "symlink to gunzip",
			args: args{
				cmdLine: flag.NewFlagSet("test", flag.ContinueOnError),
				args:    []string{"gunzip", "file.gz"},
			},
			wantOption: Options{
				Blocksize:  128,
				Level:      -1,
				Processes:  runtime.NumCPU(),
				Suffix:     ".gz",
				Decompress: true,
			},
		},
		{
			name: "symlink to gzcat",
			args: args{
				cmdLine: flag.NewFlagSet("test", flag.ContinueOnError),
				args:    []string{"gzcat", "file.gz"},
			},
			wantOption: Options{
				Blocksize:  128,
				Level:      -1,
				Processes:  runtime.NumCPU(),
				Suffix:     ".gz",
				Decompress: true,
				Stdout:     true,
			},
		},
		{
			name: "no args and no force",
			args: args{
				cmdLine: flag.NewFlagSet("test", flag.ContinueOnError),
				args:    []string{"gzip"},
			},
			wantErr: ErrStdoutNoForce,
		},
		{
			name: "request for help",
			args: args{
				cmdLine: flag.NewFlagSet("test", flag.ContinueOnError),
				args:    []string{"gzip", "-h"},
			},
			wantErr: ErrHelp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{}
			err := o.ParseArgs(tt.args.args, tt.args.cmdLine)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Options.ParseArgs() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr != nil {
				return
			}

			if !reflect.DeepEqual(*o, tt.wantOption) {
				t.Errorf("Options.ParseArgs() = \n%+v, want \n%+v", *o, tt.wantOption)
			}
		})
	}
}

func TestParseLevels(t *testing.T) {
	type args struct {
		levels [10]bool
	}
	tests := []struct {
		name    string
		want    int
		args    args
		wantErr bool
	}{
		{
			name:    "No level specified",
			args:    args{levels: [10]bool{}},
			want:    pgzip.DefaultCompression,
			wantErr: false,
		},
		{
			name:    "Level 1",
			args:    args{levels: [10]bool{false, true, false, false, false, false, false, false, false, false}},
			want:    1,
			wantErr: false,
		},
		{
			name:    "Level 9",
			args:    args{levels: [10]bool{false, false, false, false, false, false, false, false, false, true}},
			want:    9,
			wantErr: false,
		},
		{
			name:    "Multiple levels specified",
			args:    args{levels: [10]bool{false, true, false, false, false, false, false, false, false, true}},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLevels(tt.args.levels)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLevels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseLevels() = %v, want %v", got, tt.want)
			}
		})
	}
}
