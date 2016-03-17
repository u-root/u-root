// +build darwin freebsd linux netbsd openbsd

package main

import "syscall"

func sameFile(sys1, sys2 interface{}) bool {
	stat1 := sys1.(*syscall.Stat_t)
	stat2 := sys2.(*syscall.Stat_t)
	return stat1.Dev == stat2.Dev && stat1.Ino == stat2.Ino
}
