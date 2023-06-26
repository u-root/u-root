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

func TestNewPxeBooter(t *testing.T) {
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
			name: "PxeBooter",
			args: args{[]byte(`{"type":"pxeboot","ipv6":"true","ipv4":"false"}`)},
			want: &PxeBooter{Type: "pxeboot", IPV6: "true", IPV4: "false"},
			err:  nil,
		},
		{
			name: "Error PxeBooter",
			args: args{[]byte(`{"type":"boot","ipv6":"true","ipv4":"false"}`)},
			want: nil,
			err:  errWrongType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPxeBooter(tt.args.config, ulog.Log)
			if !errors.Is(err, tt.err) {
				t.Errorf("NewBootBooter() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPxeBooter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPxeBooterTypeName(t *testing.T) {
	type fields struct {
		config []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "PxeBooter TypeName",
			fields: fields{[]byte(`{"type":"pxeboot","ipv6":"true","ipv4":"false"}`)},
			want:   "pxeboot",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nb, err := NewPxeBooter(tt.fields.config, ulog.Log)
			if err != nil {
				t.Errorf("NewPxeBooter() error = %v", err)
				return
			}
			if got := nb.TypeName(); got != tt.want {
				t.Errorf("PxeBooter.TypeName() = %v, want %v", got, tt.want)
			}
		})
	}
}
