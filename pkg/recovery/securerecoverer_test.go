// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package recovery

import (
	"errors"
	"testing"
)

type testSyscalls struct{}

func (sc testSyscalls) reboot(cmd int) error {
	return nil
}

var tsc testSyscalls

type testSyscalls2 struct{}

func (sc testSyscalls2) reboot(cmd int) error {
	return errors.New("error")
}

var tsc2 testSyscalls2

func TestSecureRecover(t *testing.T) {
	for _, tt := range []struct {
		name    string
		rec     Recoverer
		wantErr bool
	}{
		{
			name: "recover",
			rec: SecureRecoverer{
				syscalls: tsc,
			},
			wantErr: false,
		},
		{
			name: "recover sync",
			rec: SecureRecoverer{
				Sync:     true,
				syscalls: tsc,
			},
			wantErr: false,
		},
		{
			name: "recover debug",
			rec: SecureRecoverer{
				Debug:    true,
				syscalls: tsc,
			},
			wantErr: false,
		},
		{
			name: "recover rand",
			rec: SecureRecoverer{
				RandWait: true,
				syscalls: tsc,
			},
			wantErr: false,
		},
		{
			name: "recover reboot",
			rec: SecureRecoverer{
				Reboot:   true,
				syscalls: tsc,
			},
			wantErr: false,
		},
		{
			name: "recover power off",
			rec: SecureRecoverer{
				Reboot:   true,
				syscalls: tsc2,
			},
			wantErr: true,
		},
		{
			name: "recover power off failed",
			rec: SecureRecoverer{
				syscalls: tsc2,
			},
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rec.Recover("test"); err != nil && !tt.wantErr {
				t.Errorf("Expected: nil Got: %v", err)
			} else if err == nil && tt.wantErr {
				t.Error("Expected: error Got: nil")
			}
		})
	}
}
