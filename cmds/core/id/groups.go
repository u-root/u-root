// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Groups manages a group database, typically loaded from /etc/group
type Groups struct {
	gidToGroup map[int]string
	groupToGID map[string]int
	userToGIDs map[string][]int
	gidToUsers map[int][]string
}

// GetGID returns the GID of a group
func (g *Groups) GetGID(name string) (int, error) {
	if v, ok := g.groupToGID[name]; ok {
		return v, nil
	}
	return -1, fmt.Errorf("unknown group name: %s", name)
}

// GetGroup gets the group of a GID
func (g *Groups) GetGroup(gid int) (string, error) {
	if v, ok := g.gidToGroup[gid]; ok {
		return v, nil
	}
	return "", fmt.Errorf("unknown gid: %d", gid)
}

// UserGetGIDs returns a slice of GIDs for a username
func (g *Groups) UserGetGIDs(username string) []int {
	if v, ok := g.userToGIDs[username]; ok {
		return v
	}
	return nil
}

// NewGroups reads the GroupFile for groups.
// It assumes the format "name:passwd:number:groupList".
func NewGroups(file string) (g *Groups, e error) {
	g = &Groups{}

	g.gidToGroup = make(map[int]string)
	g.groupToGID = make(map[string]int)
	g.gidToUsers = make(map[int][]string)
	g.userToGIDs = make(map[string][]int)

	groupFile, err := os.Open(file)
	if err != nil {
		return g, err
	}

	var groupInfo []string

	groupScanner := bufio.NewScanner(groupFile)

	for groupScanner.Scan() {
		txt := groupScanner.Text()
		if len(txt) == 0 || txt[0] == '#' { // skip empty lines and comments
			continue
		}
		groupInfo = strings.Split(txt, ":")
		groupNum, err := strconv.Atoi(groupInfo[2])
		if err != nil {
			return nil, err
		}

		g.gidToGroup[groupNum] = groupInfo[0]
		g.groupToGID[groupInfo[0]] = groupNum

		users := strings.SplitSeq(groupInfo[3], ",")
		for u := range users {
			g.userToGIDs[u] = append(g.userToGIDs[u], groupNum)
			g.gidToUsers[groupNum] = append(g.gidToUsers[groupNum], u)
		}
	}

	return
}
