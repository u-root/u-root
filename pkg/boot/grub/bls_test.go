// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grub

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
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
