// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package systembooter

import (
	"errors"
	"reflect"
	"testing"

	"github.com/u-root/u-root/pkg/ulog"
)

func TestNewBootBooter(t *testing.T) {
	type args struct {
		config []byte
	}
	tests := []struct {
		name string
		args args
		want Booter
		err  error
	}{
		{
			name: "Boot Booter",
			args: args{[]byte(`{"type": "boot"}`)},
			want: &BootBooter{Type: "boot"},
			err:  nil,
		},
		{
			name: "Error Boot Booter",
			args: args{[]byte(`{"type": "pxeboot"}`)},
			want: nil,
			err:  errWrongType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBootBooter(tt.args.config, ulog.Log)
			if !errors.Is(err, tt.err) {
				t.Errorf("NewBootBooter() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBootBooter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBootBooterBoot(t *testing.T) {
	type fields struct {
		config []byte
	}
	type args struct {
		debugEnabled bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "Boot Booter Boot",
			fields: fields{[]byte(`{"type": "boot"}`)},
			args:   args{bool(true)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb, err := NewBootBooter(tt.fields.config, ulog.Log)
			if err != nil {
				t.Errorf("NewBootBooter() error = %v", err)
				return
			}
			if err := lb.Boot(tt.args.debugEnabled); err != nil {
				t.Logf("BootBooter.Boot() error = %v", err)
			}
		})
	}
}

func TestBootBooterTypeName(t *testing.T) {
	type fields struct {
		config []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Boot Booter TypeName",
			fields: fields{[]byte(`{"type": "boot"}`)},
			want:   "boot",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb, err := NewBootBooter(tt.fields.config, ulog.Log)
			if err != nil {
				t.Errorf("NewBootBooter() error = %v", err)
				return
			}
			if got := lb.TypeName(); got != tt.want {
				t.Errorf("BootBooter.TypeName() = %v, want %v", got, tt.want)
			}
		})
	}
}
