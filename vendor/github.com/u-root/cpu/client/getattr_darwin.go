// Copyright 2018 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"syscall"

	"github.com/hugelgupf/p9/p9"
)

// GetAttr implements p9.File.GetAttr.
//
// Not fully implemented.
func (l *CPU9P) GetAttr(req p9.AttrMask) (p9.QID, p9.AttrMask, p9.Attr, error) {
	qid, fi, err := l.info()
	if err != nil {
		return qid, p9.AttrMask{}, p9.Attr{}, err
	}

	stat := fi.Sys().(*syscall.Stat_t)
	attr := p9.Attr{
		Mode:             p9.FileMode(stat.Mode),
		UID:              p9.UID(stat.Uid),
		GID:              p9.GID(stat.Gid),
		NLink:            p9.NLink(stat.Nlink),
		RDev:             p9.Dev(stat.Rdev),
		Size:             uint64(stat.Size),
		BlockSize:        uint64(stat.Blksize),
		Blocks:           uint64(stat.Blocks),
		ATimeSeconds:     uint64(stat.Atimespec.Sec),
		ATimeNanoSeconds: uint64(stat.Atimespec.Nsec),
		MTimeSeconds:     uint64(stat.Mtimespec.Sec),
		MTimeNanoSeconds: uint64(stat.Mtimespec.Nsec),
		CTimeSeconds:     uint64(stat.Ctimespec.Sec),
		CTimeNanoSeconds: uint64(stat.Ctimespec.Nsec),
	}
	valid := p9.AttrMask{
		Mode:   true,
		UID:    true,
		GID:    true,
		NLink:  true,
		RDev:   true,
		Size:   true,
		Blocks: true,
		ATime:  true,
		MTime:  true,
		CTime:  true,
	}

	return qid, valid, attr, nil
}
