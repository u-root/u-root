// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package recovery

import (
	"log"
	"os/exec"
)

// PermissiveRecoverer properties
// RecoveryCommand: unix command with absolute file path
type PermissiveRecoverer struct {
	RecoveryCommand string
}

// Recover logs error message in panic mode.
// Can jump into a shell for later debugging.
func (pr PermissiveRecoverer) Recover(message string) error {
	log.Print(message)

	if pr.RecoveryCommand != "" {
		cmd := exec.Command(pr.RecoveryCommand)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
