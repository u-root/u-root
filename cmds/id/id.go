// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Synopsis:
//      id [-gGnu]
//
// Description:
//      id displays the uid, guid and groups of the calling process
// Options:
//  		-g, --group     print only the effective group ID
//		  -G, --groups    print all group IDs
//		  -n, --name      print a name instead of a number, for -ugG
//		  -u, --user      print only the effective user ID
//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var (
	GroupFile  = "/etc/group"
	PasswdFile = "/etc/passwd"

	flags struct {
		group  bool
		groups bool
		name   bool
		user   bool
	}
)

func correctFlags(flags ...bool) bool {
	n := 0
	for _, v := range flags {
		if v {
			n += 1
		}
	}
	if n > 1 {
		return false
	} else {
		return true
	}
}

func initFlags() error {
	flag.BoolVar(&flags.group, "g", false, "print only the effective group ID")
	flag.BoolVar(&flags.groups, "G", false, "print all group IDs")
	flag.BoolVar(&flags.name, "n", false, "print a name instead of a number, for -ugG")
	flag.BoolVar(&flags.user, "u", false, "print only the effective user ID")
	flag.Parse()
	if !correctFlags(flags.groups, flags.group, flags.user) {
		return fmt.Errorf("cannot print \"only\" of more than one choice\n")

	}
	if flags.name && !(flags.groups || flags.group || flags.user) {
		return fmt.Errorf("cannot print only names in default format\n")
	}

	return nil
}

type User struct {
	name   string
	uid    int
	euid   int
	gid    int
	groups map[int]string
}

func (u *User) UID() int {
	return u.uid
}

func (u *User) pullUid() {
	u.uid = syscall.Getuid()
}

func (u *User) GetEuid() int {
	return u.euid
}

func (u *User) pullEuid() {
	u.euid = syscall.Geteuid()
}

func (u *User) GID() int {
	return u.gid
}

func (u *User) pullGid() {
	u.gid = syscall.Getgid()
}

func (u *User) Name() string {
	return u.name
}

// pullName finds the name of the User by reading PasswdFile.
func (u *User) pullName() error {
	passwdFile, err := os.Open(PasswdFile)
	if err != nil {
		return err
	}

	var passwdInfo []string

	passwdScanner := bufio.NewScanner(passwdFile)

	for passwdScanner.Scan() {
		passwdInfo = strings.Split(passwdScanner.Text(), ":")
		if val, err := strconv.Atoi(passwdInfo[2]); err != nil {
			return err
		} else if val == u.UID() {
			u.name = passwdInfo[0]
			return nil
		}
	}

	return fmt.Errorf("User is not in %s", PasswdFile)
}

func (u *User) Groups() map[int]string {
	return u.groups
}

func (u *User) GIDName() string {
	val := u.Groups()[u.UID()]
	return val
}

// pullGroups aggregates the groups that the User is in.
func (u *User) pullGroups() error {
	groupsNumbers, err := syscall.Getgroups()
	if err != nil {
		return err
	}

	groupsMap, err := readGroups()
	if err != nil {
		return err
	}

	for _, groupNum := range groupsNumbers {
		if groupName, ok := groupsMap[groupNum]; ok {
			u.groups[groupNum] = groupName
		} else {
			return fmt.Errorf("Inconsistent %s file", GroupFile)
		}
	}
	return nil
}

// readGroups reads the GroupFile for groups.
// It assumes the format "name:passwd:number:groupList".
func readGroups() (map[int]string, error) {
	groupFile, err := os.Open(GroupFile)
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
			return nil, err
		}

		groupsMap[groupNum] = groupInfo[0]
	}

	return groupsMap, nil
}

// NewUser is a factory method for the User type.
func NewUser() (*User, error) {
	emptyMap := make(map[int]string)
	u := &User{"", -1, -1, -1, emptyMap}
	u.pullUid()
	u.pullEuid()
	u.pullGid()
	if err := u.pullGroups(); err != nil {
		return nil, err
	}
	if err := u.pullName(); err != nil {
		return nil, err
	}
	return u, nil
}

// IDCommand runs the "id" with the current user's information.
func IDCommand(u User) {
	if !flags.groups {
		if flags.user {
			if flags.name {
				fmt.Println(u.Name())
				return
			}
			fmt.Println(u.UID())
			return
		} else if flags.group {
			if flags.name {
				fmt.Println(u.GIDName())
				return
			}
			fmt.Println(u.GID())
			return
		}

		fmt.Printf("uid=%d(%s) ", u.UID(), u.Name())
		fmt.Printf("gid=%d(%s) ", u.GID(), u.GIDName())
	}

	if !flags.groups {
		fmt.Print("groups=")
	}
	n := 0
	length := len(u.Groups())
	for gid, name := range u.Groups() {

		if !flags.groups {
			fmt.Printf("%d(%s)", gid, name)

			if n < length-1 {
				fmt.Print(",")
			}
			n += 1
		} else {
			if flags.name {
				fmt.Printf("%s ", name)
			} else {
				fmt.Printf("%d ", gid)
			}
		}
	}
	fmt.Println()
}

func main() {
	if err := initFlags(); err != nil {
		log.Fatalf("id: %s", err)
	}

	theChosenOne, err := NewUser()
	if err != nil {
		log.Fatalf("id: %s", err)
	}

	IDCommand(*theChosenOne)
}
