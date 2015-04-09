package main

import (
	"fmt"
	"log"
	"strings"
	"syscall"
)

func main() {
	var u syscall.Utsname
	if err := syscall.Uname(&u); err != nil {
		log.Fatalf("%v", err)
	}
        fmt.Printf("'%s' '%s' '%s' '%s' '%s' '%s'",
                strings.Trim(string([]uint8(u.Sysname[:])), "\000"),
                strings.Trim(string(u.Nodename[:]), "\000"),
                strings.Trim(string(u.Release[:]), "\000"),
                strings.Trim(string(u.Version[:]), "\000"),
                strings.Trim(string(u.Machine[:]), "\000"),
                strings.Trim(string(u.Domainname[:]), "\000"))
}
