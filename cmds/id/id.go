// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Print process information.
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
	GROUP_FILE  = "/etc/group"
	PASSWD_FILE = "/etc/passwd"
	l           = log.New(os.Stderr, "", 0)

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
	Name   string
	Uid    int
	Euid   int
	Gid    int
	Groups map[int]string
}

func (u *User) GetUid() int {
	return u.Uid
}

func (u *User) pullUid() {
	u.Uid = syscall.Getuid()
}

func (u *User) GetEuid() int {
	return u.Euid
}

func (u *User) pullEuid() {
	u.Uid = syscall.Geteuid()
}

func (u *User) GetGid() int {
	return u.Gid
}

func (u *User) pullGid() {
	u.Gid = syscall.Getgid()
}

func (u *User) GetName() string {
	return u.Name
}

// pullName finds the name of the User by reading PASSWD_FILE.
func (u *User) pullName() error {
	passwdFile, err := os.Open(PASSWD_FILE)
	if err != nil {
		return err
	}

	var passwdInfo []string

	passwdScanner := bufio.NewScanner(passwdFile)

	for passwdScanner.Scan() {
		passwdInfo = strings.Split(passwdScanner.Text(), ":")
		if val, err := strconv.Atoi(passwdInfo[2]); err != nil {
			return err
		} else if val == u.GetUid() {
			u.Name = passwdInfo[0]
			return nil
		}
	}

	return fmt.Errorf("User is not in %s", PASSWD_FILE)
}

func (u *User) GetGroups() map[int]string {
	return u.Groups
}

func (u *User) GetGidName() string {
	val := u.GetGroups()[u.GetUid()]
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
			u.Groups[groupNum] = groupName
		} else {
			return fmt.Errorf("Inconsistent %s file", GROUP_FILE)
		}
	}
	return nil
}

// readGroups reads the GROUP_FILE for groups.
// It assumes the format "name:passwd:number:groupList".
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
				fmt.Println(u.GetName())
				return
			}
			fmt.Println(u.GetUid())
			return
		} else if flags.group {
			if flags.name {
				fmt.Println(u.GetGidName())
				return
			}
			fmt.Println(u.GetGid())
			return
		}

		fmt.Printf("uid=%d(%s) ", u.GetUid(), u.GetName())
		fmt.Printf("gid=%d(%s) ", u.GetGid(), u.GetGidName())
	}

	if !flags.groups {
		fmt.Print("groups=")
	}
	n := 0
	length := len(u.Groups)
	for gid, name := range u.Groups {

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
		l.Fatalf("id: %s", err)
	}

	theChosenOne, err := NewUser()
	if err != nil {
		l.Fatalf("id: %s", err)
	}

	IDCommand(*theChosenOne)
}
