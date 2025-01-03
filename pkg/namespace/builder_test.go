// Copyright 2020-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package namespace

import (
	"bytes"
	"io"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
)

type args struct {
	ns Namespace
}

func newTestBuilder(name string) func(t minimock.Tester) *Builder {
	return func(t minimock.Tester) *Builder {
		wd, err := os.Getwd()
		if err != nil {
			t.Errorf(`os.Getwd() = _, %v, want nil`, err)
			return nil
		}
		f, err := os.Open("testdata/" + name)
		if err != nil {
			t.Errorf(`os.Open("testdata/" + name) = _, %v, want nil`, err)
			return nil
		}
		file, err := Parse(f)
		if err != nil {
			t.Errorf(`Parse(f) = _, %v, want nil`, err)
			return nil
		}
		return &Builder{
			dir:  wd,
			file: file,
			open: func(path string) (io.Reader, error) { return bytes.NewBuffer(nil), nil },
		}
	}
}
func mockNSBuilder(t minimock.Tester) args { return args{&noopNS{}} }
func TestBuilder_buildNS(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *Builder
		inspect func(r *Builder, t *testing.T) // inspects *Builder after execution of buildNS

		args func(t minimock.Tester) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) // use for more precise error evaluation
	}{}
	files, err := os.ReadDir("testdata")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		tests = append(tests, struct {
			name    string
			init    func(t minimock.Tester) *Builder
			inspect func(r *Builder, t *testing.T) // inspects *Builder after execution of buildNS

			args func(t minimock.Tester) args

			wantErr    bool
			inspectErr func(err error, t *testing.T) // use for more precise error evaluation
		}{
			name:    file.Name(),
			init:    newTestBuilder(file.Name()),
			args:    mockNSBuilder,
			wantErr: path.Ext(file.Name()) == ".wrong",
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			tArgs := tt.args(mc)
			receiver := tt.init(mc)

			err := receiver.buildNS(tArgs.ns)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if tt.wantErr {
				if err != nil && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else if err != nil {
				t.Errorf(`receiver.buildNS(%v) = %v, want nil`, tArgs.ns, err)
			}
		})
	}
}

type noopNS struct{}

func (m *noopNS) Bind(newname string, oldname string, option mountflag) error { return nil }
func (m *noopNS) Mount(servername, old, spec string, option mountflag) error  { return nil }
func (m *noopNS) Unmount(newname string, oldname string) error                { return nil }
func (m *noopNS) Import(host string, remotepath string, mountpoint string, options mountflag) error {
	return nil
}
func (m *noopNS) Clear() error            { return nil }
func (m *noopNS) Chdir(path string) error { return nil }
