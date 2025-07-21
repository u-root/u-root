// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

// id displays the user id, group id, and groups of the calling process.
//
// Synopsis:
//
//	id [-gGnu]
//
// Description:
//
//	id displays the uid, gid and groups of the calling process
//
// Options:
//
//	-g, --group     print only the effective group ID
//	-G, --groups    print all group IDs
//	-n, --name      print a name instead of a number, for -ugG
//	-u, --user      print only the effective user ID
//	-r, --user      print real ID instead of effective ID
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

const (
	groupFile  = "/etc/group"
	passwdFile = "/etc/passwd"
)

type flags struct {
	group  bool
	groups bool
	name   bool
	user   bool
	real   bool
}

var (
	errOnlyOneChoice     = errors.New("id: cannot print \"only\" of more than one choice")
	errNotOnlyNames      = errors.New("id: cannot print only names in default format")
	errNotOnlyNamesOrIDs = errors.New("id: cannot print only names or real IDs in default format")
)

func correctFlags(flags ...bool) bool {
	n := 0
	for _, v := range flags {
		if v {
			n++
		}
	}
	return n <= 1
}

// User contains user information, as from /etc/passwd
type User struct {
	groups map[int]string
	name   string
	uid    int
	gid    int
}

// UID returns the integer UID for a user
func (u *User) UID() int {
	return u.uid
}

// GID returns the integer GID for a user
func (u *User) GID() int {
	return u.gid
}

// Name returns the name for a user
func (u *User) Name() string {
	return u.name
}

// Groups returns all the groups for a user in a map
func (u *User) Groups() map[int]string {
	return u.groups
}

// GIDName returns the group name for a user's UID
func (u *User) GIDName() string {
	val := u.Groups()[u.UID()]
	return val
}

// NewUser is a factory method for the User type.
func NewUser(flags *flags, username string, users *Users, groups *Groups) (*User, error) {
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
func IDCommand(w io.Writer, flags *flags, u User) {
	if !flags.groups {
		if flags.user {
			if flags.name {
				fmt.Fprintln(w, u.Name())
				return
			}
			fmt.Fprintln(w, u.UID())
			return
		} else if flags.group {
			if flags.name {
				fmt.Fprintln(w, u.GIDName())
				return
			}
			fmt.Fprintln(w, u.GID())
			return
		}

		fmt.Fprintf(w, "uid=%d(%s) ", u.UID(), u.Name())
		fmt.Fprintf(w, "gid=%d(%s) ", u.GID(), u.GIDName())
	}

	if !flags.groups {
		fmt.Fprintf(w, "groups=")
	}

	// handle Go map ordering
	var gids []int
	for gid := range u.Groups() {
		gids = append(gids, gid)
	}
	sort.Ints(gids)

	var groupOutput []string
	for id := range gids {
		gid, name := id, u.Groups()[id]
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

	fmt.Fprintln(w, strings.Join(groupOutput, sep))
}

func run(w io.Writer, name string, f *flags, passwd, group string) error {
	if !correctFlags(f.groups, f.group, f.user) {
		return errOnlyOneChoice
	}
	if f.name && (!f.groups && !f.group && !f.user) {
		return errNotOnlyNames
	}
	if len(name) != 0 && f.real {
		return errNotOnlyNamesOrIDs
	}

	users, err := NewUsers(passwd)
	if err != nil {
		return fmt.Errorf("id: %w", err)
	}
	groups, err := NewGroups(group)
	if err != nil {
		return fmt.Errorf("id: %w", err)
	}

	user, err := NewUser(f, name, users, groups)
	if err != nil {
		return fmt.Errorf("id: %w", err)
	}

	IDCommand(w, f, *user)
	return nil
}

func main() {
	flags := &flags{}
	flag.BoolVar(&flags.group, "g", false, "print only the effective group ID")
	flag.BoolVar(&flags.groups, "G", false, "print all group IDs")
	flag.BoolVar(&flags.name, "n", false, "print a name instead of a number, for -ugG")
	flag.BoolVar(&flags.user, "u", false, "print only the effective user ID")
	flag.BoolVar(&flags.real, "r", false, "print real ID instead of effective ID")

	flag.Parse()
	if err := run(os.Stdout, flag.Arg(0), flags, passwdFile, groupFile); err != nil {
		log.Fatalf("%v", err)
	}
}
