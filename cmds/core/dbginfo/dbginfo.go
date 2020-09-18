package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func runCommand(command string, args ...string) {
	c := exec.Command(command, args...)
	argsFmt := fmt.Sprintf("%v", args)
	argsFmt = strings.TrimPrefix(argsFmt, "[")
	argsFmt = strings.TrimSuffix(argsFmt, "]")
	if b, e := c.Output(); e != nil {
		log.Printf("DEBUG: %s %s\n error: %v", command, argsFmt, e)
	} else {
		log.Printf("DEBUG: %s %s\n %s", command, argsFmt, string(b))
	}
}

func main() {
	runCommand("cat", "/proc/filesystems")
	runCommand("ls", "/dev/sd*", "/dev/nvme*")
	runCommand("ls", "/tmp/*/boot/grub/entry-*.cfg")
	runCommand("grep", "/disk", "/proc/mounts")
	runCommand("grep", "image=", "/disk/boot/grub/entry-1.cfg")
	runCommand("ip", "a")
	runCommand("ip", "route")
	runCommand("ip", "-6", "route")
}
