// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ls

import (
	"fmt"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

const lsregex string = "^([rwxSTstdcb\\-lp?]{10})\\s+(\\d+)?\\s?(\\S+)\\s+(\\S+)\\s+([0-9,]+)?\\s+(\\d+)?(\\D+)?(\\d{1,2}\\D\\d{1,2}\\D\\d{1,2})?(\\D{4})?([\\D|\\d]*)"
const tDelta time.Duration = 1 * time.Second

func TestFileInfo(t *testing.T) {
	u, err := user.Current()
	if err != nil {
		t.Errorf("Failed getting current user: %q", err)
	}
	userid, err := strconv.ParseUint(u.Uid, 0, 32)
	if err != nil {
		t.Errorf("Failed to convert userid to uint64: %q", err)
	}
	groupid, err := strconv.ParseUint(u.Gid, 0, 32)
	if err != nil {
		t.Errorf("Failed to convert groud id to uint64: %q", err)
	}
	gidname, err := user.LookupGroupId(u.Gid)
	if err != nil {
		t.Errorf("Failed look up group id: %q", err)
	}
	fmt.Println(gidname)
	for _, tt := range []struct {
		name          string
		filename      string
		filemode      string
		rdev          uint64
		uid           uint32
		gid           uint32
		size          int64
		mTime         time.Time
		symlinktarget string
		user          string
		group         string
	}{
		{
			name:          "SimpleFile",
			filename:      "testFile-",
			filemode:      "-rw-------",
			rdev:          0,
			uid:           uint32(userid),
			gid:           uint32(groupid),
			size:          0,
			mTime:         time.Now(),
			symlinktarget: "",
			user:          u.Username,
			group:         gidname.Name,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			f, err := os.CreateTemp(tmpDir, tt.filename)
			if err != nil {
				t.Errorf("%q failed with %q", tt.name, err)
			}
			osfi, err := os.Stat(f.Name())
			if err != nil {
				t.Errorf("%q failed with %q", tt.name, err)
			}
			fi := FromOSFileInfo(f.Name(), osfi)

			if !strings.Contains(fi.Name, tt.filename) {
				t.Errorf("%q failed by name. Got: %q, Want: %q", tt.name, tt.filename, strings.TrimPrefix(fi.PrintableName(), tt.name))
			}
			if !strings.Contains(fi.Mode.String(), tt.filemode) {
				t.Errorf("%q failed by filemode. Got: %q, Want: %q", tt.name, tt.filemode, fi.Mode.String())
			}
			if fi.Rdev != tt.rdev {
				t.Errorf("%q failed by Rdev. Got: %q, Want: %q", tt.name, tt.rdev, fi.Rdev)
			}
			if fi.UID != tt.uid {
				t.Errorf("%q failed by UID. Got: %q, Want: %q", tt.name, tt.uid, fi.UID)
			}
			if fi.GID != tt.gid {
				t.Errorf("%q failed by GID. Got: %q, Want: %q", tt.name, tt.gid, fi.GID)
			}
			if fi.Size != tt.size {
				t.Errorf("%q failed by Size. Got: %q, Want: %q", tt.name, tt.size, fi.Size)
			}

			if dt := tt.mTime.Sub(fi.MTime); dt < -tDelta || dt > tDelta {
				t.Errorf("%q failed by MTime. Got: %q, Want: %q", tt.name, tt.mTime, fi.MTime)
			}

			if fi.SymlinkTarget != tt.symlinktarget {
				t.Errorf("%q failed by SymlinkTarget. Got: %q, Want: %q", tt.name, tt.symlinktarget, fi.SymlinkTarget)
			}

			if fi.Name != fi.PrintableName() {
				t.Errorf("%q failed by PrintableName. Got: %q, Want: %q", tt.name, fi.Name, fi.PrintableName())
			}

			if lookupUserName(fi.UID) != tt.user {
				t.Errorf("%q failed by lookupUserName. Got: %q, Want: %q", tt.name, tt.user, lookupUserName(fi.UID))
			}

			if lookupGroupName(fi.GID) != tt.group {
				t.Errorf("%q failed by lookupGroupName. Got: %q, Want: %q", tt.name, tt.group, lookupGroupName(fi.GID))
			}

			ls := LongStringer{
				Human: true,
				Name:  NameStringer{},
			}
			matched, err := regexp.MatchString(lsregex, ls.FileString(fi))
			if err != nil {
				t.Errorf("%q failed at regexp.MatchString: %q", tt.name, err)
			}
			if !matched {
				t.Errorf("%q failed. Output of ls.FileString does not match regular expression", tt.name)
			}
		})
	}
}
