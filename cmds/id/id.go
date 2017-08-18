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

func (u *User) getUid() int {
	return u.Uid
}

func (u *User) pullUid() {
	u.Uid = syscall.Getuid()
}

func (u *User) getEuid() int {
	return u.Euid
}

func (u *User) pullEuid() {
	u.Uid = syscall.Geteuid()
}

func (u *User) getName() string {
	return u.Name
}

func (u *User) pullName() {
	passwdFile, err := os.Open(PASSWD_FILE)
	if err != nil {
		log.Fatal(err)
	}

	var passwdInfo []string

	passwdScanner := bufio.NewScanner(passwdFile)

	for passwdScanner.Scan() {
		passwdInfo = strings.Split(passwdScanner.Text(), ":")
		if val, err := strconv.Atoi(passwdInfo[2]); err != nil {
			log.Fatal(err)
		} else if val == u.getUid() {
			u.Name = passwdInfo[0]
			return
		}
	}

	log.Fatal()
}

func (u *User) getGroups() []*Group {
	return u.Groups
}

func (u *User) pullGroups() {
	groupsNumbers, err := syscall.Getgroups()
	if err != nil {
		log.Fatal(err)
	}

	groupsMap, err := readGroups()
	if err != nil {
		log.Fatal(err)
	}

	for _, groupNum := range groupsNumbers {
		u.Groups = append(u.Groups, &Group{
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

	var groupInfo []string

	groupsMap := make(map[int]string)
	groupScanner := bufio.NewScanner(groupFile)

	for groupScanner.Scan() {
		groupInfo = strings.Split(groupScanner.Text(), ":")
		groupNum, err := strconv.Atoi(groupInfo[2])
		if err != nil {
			log.Fatal(err)
		}

		groupsMap[groupNum] = groupInfo[0]
	}

	return groupsMap, nil
}

func (u *User) Init() {
	u.pullUid()
	u.pullEuid()
	u.pullGroups()
	u.pullName()
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
