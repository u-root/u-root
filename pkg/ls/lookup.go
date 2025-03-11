// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tamago || unix

package ls

import (
	"fmt"
	"os/user"
)

// Without this cache, `ls -l` is orders of magnitude slower.
var (
	uidCache = map[uint32]string{}
	gidCache = map[uint32]string{}
)

// Convert uid to username, or return uid on error.
func lookupUserName(id uint32) string {
	if s, ok := uidCache[id]; ok {
		return s
	}
	s := fmt.Sprint(id)
	if u, err := user.LookupId(s); err == nil {
		s = u.Username
	}
	uidCache[id] = s
	return s
}

// Convert gid to group name, or return gid on error.
func lookupGroupName(id uint32) string {
	if s, ok := gidCache[id]; ok {
		return s
	}
	s := fmt.Sprint(id)
	if g, err := user.LookupGroupId(s); err == nil {
		s = g.Name
	}
	gidCache[id] = s
	return s
}
