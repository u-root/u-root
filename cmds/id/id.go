// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Print process information.
//
// Synopsis:
//     id
//
// Description:
//     id displays the uid, guid and groups of the calling process
//
// Options:
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var (
	GROUP_FILE  = "/etc/group"
	PASSWD_FILE = "/etc/passwd"
)

type Group struct {
	Name   string
	Number int
}

type User struct {
	Name   string
	Uid    int
	Euid   int
	Groups []*Group
}

func (u *User) getUid() {
	u.Uid = syscall.Getuid()
}

func (u *User) getEuid() {
	u.Uid = syscall.Getieuid()
}

func (u *User) getGroups() {
	groupsNumbers, err := syscall.Getgroups()
	if err != nil {
		log.Fatal(err)
	}

	groupsMap := readGroups()

	for _, groupNum := range groupsNumbers {
		u.Groups = append(u.Groups, Group{
			Name:   groupsMap[groupNum],
			Number: groupNum,
		})
	}

}

func readGroups() (map[int]string, error) {
	groupFile, err := os.Open(GROUP_FILE)
	if err != nil {
		return nil, err
	}

	var line string
	var groupInfo []string

	groupsMap := make(map[int]string)
	groupScanner := bufio.NewScanner(groupFile)

	for groupScanner.Scan() {
		groupInfo = strings.Split(groupScanner.Text(), ":")
		groupsMap[strconv.Atoi(groupInfo[2])] = groupInfo[0]
	}

	return groupMap, nil
}

func main() {
	uid := syscall.Getuid()
	gid := syscall.Getgid()
	groups, err := syscall.Getgroups()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("uid: %d\n", uid)
	fmt.Printf("gid: %d\n", gid)

	fmt.Print("groups: ")
	for _, group := range groups {
		fmt.Printf("%d ", group)
	}
	fmt.Println()

}
