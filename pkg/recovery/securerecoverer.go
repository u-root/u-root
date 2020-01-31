// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package recovery

import (
	"log"
	"math/rand"
	"syscall"
	"time"
)

// DebugTimeout sets the timeout for how long
// the debug message is shown before power cycle.
const DebugTimeout time.Duration = 10

// SecureRecoverer properties
// Reboot: does a reboot if true
// Sync: sync file descriptors and devices
// Debug: enables debug messages
type SecureRecoverer struct {
	Reboot   bool
	Sync     bool
	Debug    bool
	RandWait bool
}

// Recover by reboot or poweroff without or with sync
func (sr SecureRecoverer) Recover(message string) error {
	if sr.Sync {
		syscall.Sync()
	}

	if sr.Debug {
		if message != "" {
			log.SetPrefix("RECOVERY: ")
			log.Print(message)
		}
		time.Sleep(DebugTimeout * time.Second)
	}

	if sr.RandWait {
		rd := time.Duration(rand.Intn(15))
		time.Sleep(rd * time.Second)
		log.SetPrefix("RECOVERY: ")
		log.Printf("Reboot in %s seconds", rd)
	}

	if sr.Reboot {
		if err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART); err != nil {
			return err
		}
	} else {
		if err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF); err != nil {
			return err
		}
	}

	return nil
}
