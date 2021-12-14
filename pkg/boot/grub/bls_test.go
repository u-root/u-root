// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grub

import (
	"bytes"
	"context"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/u-root/u-root/pkg/curl"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
)

// This test is specifically for the problems reported in issue 2203

func TestBLSParseEnvFile(t *testing.T) {
	file := `# GRUB Environment Block
saved_entry=3370cc3add884b58835edcc89eeb00a0-5.14.1
kernelopts=root=UUID=1b63c044-d73b-4982-8e76-ef3b7b5b65f2 ro crashkernel=auto net.ifnames=0 console=ttyS1,57600n8 nomodeset 
boot_success=0
#######################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################################`

	gotEnv, err := ParseEnvFile(bytes.NewBufferString(file))
	if err != nil {
		t.Errorf("ParseEnvFile(%q) error %v", file, err)
	}
	wantEnv := &EnvFile{map[string]string{
		"saved_entry":  "3370cc3add884b58835edcc89eeb00a0-5.14.1",
		"kernelopts":   "root=UUID=1b63c044-d73b-4982-8e76-ef3b7b5b65f2 ro crashkernel=auto net.ifnames=0 console=ttyS1,57600n8 nomodeset ",
		"boot_success": "0",
	}}
	if diff := cmp.Diff(wantEnv, gotEnv); diff != "" {
		t.Errorf("ParseEnvFile(%q) diff(-want, +got) = \n%s", file, diff)
	}
}

// This is a slightly different set of tests for the issue. It is done this way to avoid bleeding
// constraints into the other tests, or have other tests bleed into this test.
func TestBLSGrubTests(t *testing.T) {
	files, err := filepath.Glob("testdata/*bls*.cfg")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Running BLS test with %q, *update is %v", files, *update)
	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".cfg")
		t.Run(name, func(t *testing.T) {
			golden := strings.TrimSuffix(file, ".cfg") + ".out"
			var out []byte
			// parse with our parser and compare
			var b bytes.Buffer
			wd := &url.URL{
				Scheme: "file",
				Path:   "./testdata",
			}
			mountPool := &mount.Pool{}
			c := newParser(wd, block.BlockDevices{}, mountPool, curl.DefaultSchemes)
			c.W = &b

			script, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("error loading file `%s`, %v", file, err)
			}
			err = c.append(context.Background(), string(script))
			if err != nil {
				t.Fatalf("error parsing file `%s`, %v", file, err)
			}

			if b.String() != string(out) {
				t.Fatalf("wrong script parsing output got `%s` want `%s`", b.String(), string(out))
			}
			// update/create golden file on success
			if *update {
				err := os.WriteFile(golden, out, 0o644)
				if err != nil {
					t.Fatalf("error writing file `%s`, %v", file, err)
				}
			}
		})

	}
}
