// Copyright 2014-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mount

import "fmt"

// This code does not use a map. Why?
// The name and magic pairs are stored as simple tuples, rather
// than a map, because these kinds of maps are rarely searched
// more than once, and the size of such a map in the busybox
// is much larger than a simple array of structs.
// These values are from Linux 5.4. FS magics change very
// slowly over time, so there is no generate script.
// This is not tagged for any one kernel as these values
// are (or should be) universal.
var magics = []struct {
	magic uint32
	name  string
}{
	{magic: 0x00c36400, name: "ceph"},
	{magic: 0xa501fcf5, name: "vxfs"},
	{magic: 0x65735543, name: "fuse_ctl"},
	{magic: 0x65735546, name: "fuse"},
	{magic: 0x482b, name: "hfsplus"},
	{magic: 0x20030528, name: "orangefs"},
	{magic: 0x24051905, name: "ubifs"},
	{magic: 0xadf5, name: "adfs"},
	{magic: 0xadff, name: "affs"},
	{magic: 0x5346414f, name: "afs"},
	{magic: 0x0187, name: "autofs"},
	{magic: 0x73757245, name: "coda"},
	{magic: 0xf15f, name: "ecryptfs"},
	{magic: 0x414a53, name: "efs"},
	{magic: 0xe0f5e1e2, name: "erofs"},
	// This is what Linux has: three filesystems with the same magic.
	// HMMM.
	//{magic: 0xef53, name: "ext2"},
	//{magic: 0xef53, name: "ext3"},
	{magic: 0xabba1974, name: "xenfs"},
	{magic: 0xef53, name: "ext4"},
	{magic: 0x9123683e, name: "btrfs"},
	{magic: 0x3434, name: "nilfs"},
	{magic: 0xf2f52010, name: "f2fs"},
	{magic: 0xf995e849, name: "hpfs"},
	{magic: 0x9660, name: "isofs"},
	{magic: 0x72b6, name: "jffs2"},
	{magic: 0x00c0ffee, name: "hostfs"},
	{magic: 0x794c7630, name: "overlayfs"},
	{magic: 0x6969, name: "nfs"},
	{magic: 0x7461636f, name: "ocfs2"},
	{magic: 0x9fa1, name: "openprom"},
	{magic: 0x517b, name: "smb"},
	{magic: 0x27e0eb, name: "cgroup"},
	{magic: 0x63677270, name: "cgroup2"},
	{magic: 0x7655821, name: "rdtgroup"},
	{magic: 0x1cd1, name: "devpts"},
	{magic: 0x6c6f6f70, name: "binderfs"},
	{magic: 0xbad1dea, name: "futexfs"},
	{magic: 0x9fa0, name: "proc"},
	{magic: 0x9fa2, name: "usbdevice"},
	{magic: 0x15013346, name: "udf"},
	{magic: 0x9fa0, name: "proc"},
}

// MagicToName returns a string and error given a
// file system magic number as uint32.
func MagicToName(magic uint32) (string, error) {
	for _, m := range magics {
		if m.magic == magic {
			return m.name, nil
		}
	}
	return "", fmt.Errorf("No file system for %#x", magic)
}

// NameToMagic returns a string and error given a
// file system name.
func NameToMagic(n string) (uint32, error) {
	for _, m := range magics {
		if m.name == n {
			return m.magic, nil
		}
	}
	return 0, fmt.Errorf("No file system for %q", n)
}
