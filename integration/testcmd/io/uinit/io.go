// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

// Prints the string "UART TEST\r\n" using IO.
func main() {
	// Writing to the serial port is atomic in QEMU, so no polling or
	// sleeping is needed between characters.
	for _, b := range []byte("UART TEST\r\n") {
		cmd := exec.Command("io", "outb", "0x3f8", fmt.Sprintf("%d", b))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF)
}
