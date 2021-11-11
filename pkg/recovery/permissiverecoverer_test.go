// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package recovery

import "testing"

func TestPermissiveRecover(t *testing.T) {
	for _, tt := range []struct {
		name    string
		rec     Recoverer
		wantErr bool
	}{
		{
			name:    "recover no cmd",
			rec:     PermissiveRecoverer{},
			wantErr: false,
		},
		{
			name: "recover cmd",
			rec: PermissiveRecoverer{
				RecoveryCommand: "bogus",
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
