package gzip

import (
	"flag"
	"runtime"
	"testing"

	"github.com/klauspost/pgzip"
)

func TestOptions_ParseArgs(t *testing.T) {
	type fields struct {
		Blocksize  int
		Level      int
		Processes  int
		Decompress bool
		Force      bool
		Help       bool
		Keep       bool
		Quiet      bool
		Stdin      bool
		Stdout     bool
		Test       bool
		Verbose    bool
		Suffix     string
	}
	type args struct {
		args    []string
		cmdLine *flag.FlagSet
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				Blocksize:  tt.fields.Blocksize,
				Level:      tt.fields.Level,
				Processes:  tt.fields.Processes,
				Decompress: tt.fields.Decompress,
				Force:      tt.fields.Force,
				Help:       tt.fields.Help,
				Keep:       tt.fields.Keep,
				Quiet:      tt.fields.Quiet,
				Stdin:      tt.fields.Stdin,
				Stdout:     tt.fields.Stdout,
				Test:       tt.fields.Test,
				Verbose:    tt.fields.Verbose,
				Suffix:     tt.fields.Suffix,
			}
			if err := o.ParseArgs(tt.args.args, tt.args.cmdLine); (err != nil) != tt.wantErr {
				t.Errorf("Options.ParseArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestOptions_validate(t *testing.T) {
	type fields struct {
		Blocksize  int
		Level      int
		Processes  int
		Decompress bool
		Force      bool
		Help       bool
		Keep       bool
		Quiet      bool
		Stdin      bool
		Stdout     bool
		Test       bool
		Verbose    bool
		Suffix     string
	}
	type args struct {
		moreArgs bool
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr bool
	}{
		{
			name:    "Default values no args",
			fields:  fields{Blocksize: 128, Level: -1, Processes: runtime.NumCPU(), Decompress: false, Force: false, Help: false},
			args:    args{moreArgs: false},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				Blocksize:  tt.fields.Blocksize,
				Level:      tt.fields.Level,
				Processes:  tt.fields.Processes,
				Decompress: tt.fields.Decompress,
				Force:      tt.fields.Force,
				Help:       tt.fields.Help,
				Keep:       tt.fields.Keep,
				Quiet:      tt.fields.Quiet,
				Stdin:      tt.fields.Stdin,
				Stdout:     tt.fields.Stdout,
				Test:       tt.fields.Test,
				Verbose:    tt.fields.Verbose,
				Suffix:     tt.fields.Suffix,
			}

			if err := o.validate(tt.args.moreArgs); (err != nil) != tt.wantErr {
				t.Errorf("Options.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parseLevels(t *testing.T) {
	type args struct {
		levels [10]bool
	}
	tests := []struct {
		name    string
		args    args
		want    int
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
			name:    "Multuple levels specified",
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
