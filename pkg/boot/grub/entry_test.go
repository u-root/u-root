// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grub

import (
	"bytes"
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
	file := `kernel=bzImage
initrd=initramfs.cpio
`
	gotEnv, err := ParseEnvFile(bytes.NewBufferString(file))
	if err != nil {
		t.Errorf("ParseEnvFile(%q) error %v", file, err)
	}
	wantEnv := &EnvFile{map[string]string{
		"kernel": "bzImage",
		"initrd": "initramfs.cpio",
	}}
	if diff := cmp.Diff(wantEnv, gotEnv); diff != "" {
		t.Errorf("ParseEnvFile(%q) diff(-want, +got) = \n%s", file, diff)
	}
}

func FuzzParseEnvFile(f *testing.F) {
	f.Add([]byte(`kernel=bzImage
		initrd=initramfs.cpio
	`))
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
