// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grub

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWriteTo(t *testing.T) {
	env := &EnvFile{map[string]string{
		"kernel": "bzImage",
		"initrd": "initramfs.cpio",
	}}
	buf := &bytes.Buffer{}
	_, err := env.WriteTo(buf)
	if err != nil {
		t.Errorf("env.WriteTo(%v) error %v", env, err)
	}
	gotFile := buf.String()
	wantFile := `# GRUB Environment Block
initrd=initramfs.cpio
kernel=bzImage
##################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################`
	if diff := cmp.Diff(wantFile, gotFile); diff != "" {
		t.Errorf("env.WriteTo(%v) diff(-want, +got) = \n%s", env, diff)
	}
}

func TestParseEnvFile(t *testing.T) {
	testcases := []struct {
		file    string
		wantEnv *EnvFile
		wantErr error
	}{
		{file: `kernel=bzImage
initrd=initramfs.cpio`, wantEnv: &EnvFile{map[string]string{
			"kernel": "bzImage",
			"initrd": "initramfs.cpio",
		}}},
		{file: `kernel=
initrd=initramfs.cpio`, wantEnv: nil, wantErr: fmt.Errorf(`error parsing "kernel=": either the key or value is empty: "kernel" = ""`)},
		{
			file:    `kernel`,
			wantEnv: nil,
			wantErr: fmt.Errorf(`error parsing "kernel": must find = or # and key + values in each line`),
		},
	}
	for _, tt := range testcases {
		gotEnv, err := ParseEnvFile(bytes.NewBufferString(tt.file))

		if tt.wantErr != nil {
			if err == nil {
				t.Fatalf("expected error: %v but got nil", tt.wantErr)
			}
			if err.Error() != tt.wantErr.Error() {
				t.Errorf("expected error:\n%v but got\n%v", tt.wantErr, err)
			}
		}

		if diff := cmp.Diff(tt.wantEnv, gotEnv); diff != "" {
			t.Errorf("ParseEnvFile(%q) diff(-want, +got) = \n%s", tt.file, diff)
		}
	}
}

func FuzzParseEnvFile(f *testing.F) {
	f.Add([]byte(`kernel=bzImage
		initrd=initramfs.cpio
	`))
	f.Add([]byte("="))
	f.Add([]byte("\r"))
	f.Fuzz(func(t *testing.T, env []byte) {
		readEnv, err := ParseEnvFile(bytes.NewBuffer(env))
		// just return if the given file is not parsable as an env file
		if err != nil {
			return
		}

		writeBuf := &bytes.Buffer{}
		readEnv.WriteTo(writeBuf)

		rereadEnv, err := ParseEnvFile(writeBuf)
		if err != nil {
			t.Fatalf("could not parse previously written env file %#v to %#v: %v", readEnv, writeBuf, err)
		}
		if !reflect.DeepEqual(readEnv, rereadEnv) {
			t.Fatalf("Env files do not match: %#v - %#v", readEnv, rereadEnv)
		}
	})
}
