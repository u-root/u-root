// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

// id displays the user id, group id, and groups of the calling process.
//
// Synopsis:
//      id [-gGnu]
//
// Description:
//      id displays the uid, gid and groups of the calling process
//
// Options:
//	-g, --group     print only the effective group ID
//	-G, --groups    print all group IDs
//	-n, --name      print a name instead of a number, for -ugG
//	-u, --user      print only the effective user ID
//	-r, --user      print real ID instead of effective ID
package main

import (
	"flag"
	"fmt"
	"log"
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
		real   bool
	}
)

func correctFlags(flags ...bool) bool {
	n := 0
	for _, v := range flags {
		if v {
			n++
		}
	}
	return !(n > 1)
}

func init() {
	flag.BoolVar(&flags.group, "g", false, "print only the effective group ID")
	flag.BoolVar(&flags.groups, "G", false, "print all group IDs")
	flag.BoolVar(&flags.name, "n", false, "print a name instead of a number, for -ugG")
	flag.BoolVar(&flags.user, "u", false, "print only the effective user ID")
	flag.BoolVar(&flags.real, "r", false, "print real ID instead of effective ID")
}

type User struct {
	name   string
	uid    int
	gid    int
	groups map[int]string
}

func (u *User) UID() int {
	return u.uid
}

func (u *User) GID() int {
	return u.gid
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Groups() map[int]string {
	return u.groups
}

func (u *User) GIDName() string {
	val := u.Groups()[u.UID()]
	return val
}

// NewUser is a factory method for the User type.
func NewUser(username string, users *Users, groups *Groups) (*User, error) {
	var groupsNumbers []int

	u := &User{groups: make(map[int]string)}
	if len(username) == 0 { // no username provided, get current
		if flags.real {
			u.uid = syscall.Geteuid()
			u.gid = syscall.Getegid()
		} else {
			u.uid = syscall.Getuid()
			u.gid = syscall.Getgid()
		}
		groupsNumbers, _ = syscall.Getgroups()
		if v, err := users.GetUser(u.uid); err == nil {
			u.name = v
		} else {
			u.name = strconv.Itoa(u.uid)
		}
	} else {
		if v, err := users.GetUID(username); err == nil { // user is username
			u.name = username
			u.uid = v
		} else {
			if uid, err := strconv.Atoi(username); err == nil { // user is valid int
				if v, err := users.GetUser(uid); err == nil { // user is valid uid
					u.name = v
					u.uid = uid
				}
			} else {
				return nil, fmt.Errorf("id: no such user or uid: %s", username)
			}
		}
		u.gid, _ = users.GetGID(u.uid)
		groupsNumbers = append([]int{u.gid}, groups.UserGetGIDs(u.name)...)
		// FIXME: not yet implemented group listing lookups
	}

	for _, groupNum := range groupsNumbers {
		if groupName, err := groups.GetGroup(groupNum); err == nil {
			u.groups[groupNum] = groupName
		} else {
			u.groups[groupNum] = strconv.Itoa(groupNum)
		}
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

	var groupOutput []string

	for gid, name := range u.Groups() {

		if !flags.groups {
			groupOutput = append(groupOutput, fmt.Sprintf("%d(%s)", gid, name))

		} else {
			if flags.name {
				groupOutput = append(groupOutput, fmt.Sprintf("%s ", name))
			} else {
				groupOutput = append(groupOutput, fmt.Sprintf("%d ", gid))
			}
		}
	}

	sep := ","
	if flags.groups {
		sep = ""
	}

	fmt.Println(strings.Join(groupOutput, sep))
}

func main() {
	flag.Parse()
	if !correctFlags(flags.groups, flags.group, flags.user) {
		log.Fatalf("id: cannot print \"only\" of more than one choice")

	}
	if flags.name && !(flags.groups || flags.group || flags.user) {
		log.Fatalf("id: cannot print only names in default format")
	}
	if len(flag.Arg(0)) != 0 && flags.real {
		log.Fatalf("id: cannot print only names or real IDs in default format")
	}

	users, err := NewUsers(PasswdFile)
	if err != nil {
		log.Printf("id: unable to read %s: %v", PasswdFile, err)
	}
	groups, err := NewGroups(GroupFile)
	if err != nil {
		log.Printf("id: unable to read %s: %v", PasswdFile, err)
	}

	user, err := NewUser(flag.Arg(0), users, groups)
	if err != nil {
		log.Fatalf("id: %s", err)
	}

	IDCommand(*user)
}
