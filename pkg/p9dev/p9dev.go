// SPDX-License-Identifier: MIT
// Copyright 2026 Google LLC
//
// Package p9dev wraps a p9 Attacher such that devices can be accessed as
// regular files.
//
// Linux is a pain when it comes to accessing special files, e.g. character or
// block devices, across a mount point.  When you open a special file, you'll
// get the file on the *client*, not the server.  Which can be surprising.
//
// Specifically, Linux's 9p client code (fs/9p/*) calls init_special_inode() on
// special files, and the VFS will override the f_ops with the ops for the
// client's OS's file, based on the major/minor number.
//
// This behavior is fine for something like a locally mounted disk.  Say you
// have an ext4 filesystem on your local HDD and it has some /dev/ files, e.g.
// /dev/console (a character device).  The intent there is that you have a
// /dev/console that points to the console of the machine that is doing the
// mount.  i.e. your root filesystem has these statically made device files
// installed - they gotta come from somewhere!  (at least in Linux).
//
// In this sense, those special files are almost like symlinks.  Open
// /dev/console?  That really means "open your magic local file 5-1".  Some
// systems have /dev/ populated with all sorts of devices too; some that don't
// even exist, but could!  e.g. /dev/sd*.  I've got "sdb15" (a blockdev), but
// certainly don't have that 15 partitions on sdb.  If you try to open it,
// you'll get something like ENXIO.
//
// A similar thing happens across a mount point.  e.g. the 9p server might have
// 100 cpus in /dev/cpu/, but you only have 10.  If you ls /mnt/dev/cpu, you'll
// see 100 cpus!  If you access /mnt/dev/cpu/0/cpuid, you're getting your *own*
// /dev/cpu/0 (since they are both "quasi symlinks" to the magic chardev file).
// If you try to access /mnt/dev/cpu/99/cpuid, you'll get ENXIO, just like with
// /dev/sdb15.  Since the "link" points to a non-existent device.
//
// But when you access a 9P server, the server could intend to offer one of two
// things:
//
//  1. Just like with ext4, it's trying to serve a root filesystem to you.
//     e.g.  you're on a diskless machine and need some sort of /dev/.  Like a
//     distro provided to you over 9p.  In this case, you want the chardevs of
//     the client's OS.
//
//  2. It really wants you to access its character devices, just like in
//     Plan 9.  It's all just reads and writes, and no file is special.  e.g.
//     it wants you to read its cpuid file.  In this case, you don't want the
//     Linux client (kernel module) to do anything special.
//
// To make things more confusing, Linux's mount has the 'nodev' option, which
// prevents accessing special devices over a mount.  This is the client's way of
// saying "i don't want to trust any random nodes on this filesystem".  But
// there isn't really a way to say "i want to access character devices, but i
// want it to be *their* devices".  At least not with the 9p client.  9p2000.u
// has the "nodevmap" option which does that, but we're using 9p2000.L, which
// does not.
//
// We could change the kernel code to add some form of "let me access the
// server's special devices if it exposed them".  And just skip the
// init_special_inode() call in v9fs_init_inode().  Your client would need to
// know that you intend to serve special files, and then mount accordingly.
// (Not a big deal.)  In this case, I think the device would still appear to be
// a chardev, but the f_ops would still be v9fs's.  i.e. ls -l would say
// chardev, I think.  Not sure if you'd need the nodev mount option here to
// access the file.
//
// Alternatively, the server can rewrite the file mode for any special device
// such that the client thinks it's a regular file and just does its reads and
// writes.  Which is what this package does.  In this case, there's no special
// mount option.  The server is explicitly giving access to its special files.
// (By virtue of using this package).  It's in the best position to know whether
// or not it needs this service.  The client can still do nodev (which is a good
// idea in the scenario #2), though not a big deal.  The 'downside' is that the
// client will see all special files as regular, which may be a little confusing
// if you were expecting to ls /dev/console as see a chardev instead of a
// regular file.  (But if you did that, I judge you somewhat).
//
// Finally, note that we explicitly do not allow mknod below.  The server is
// giving access to a set of its device files.  It is not saying "and go ahead
// and make and access *any* device file on my machine!  That's a benefit of the
// server-side approach to this.  Though in general, no 9p server should let you
// mknod on their machine.  It's like being able to make a symlink outside of
// the root of the exported directory.  (And of the p9 servers I looked at, none
// of them let you use mknod).  Also, mknod (with no caching) is the one path in
// the v9fs code that doesn't do a stat before making an inode.  So if you
// really wanted to mknod on the server in this #2 world, you'd need to unmount
// (to destroy the dentry) or wait an unknown amount of time for the dentry to
// die, then walk/open again.
package p9dev

import (
	"syscall"

	"github.com/hugelgupf/p9/p9"
)

type attacher struct {
	p9.Attacher
}

// New wraps a p9.Attacher.
func New(a p9.Attacher) p9.Attacher {
	return &attacher{Attacher: a}
}

// Attach implements p9.Attach.
func (a *attacher) Attach() (p9.File, error) {
	f, err := a.Attacher.Attach()
	if err != nil {
		return nil, err
	}
	return &File{File: f}, nil
}

// File wraps a p9.File such that we can override stat (GetAttr) calls.
// Note we need to override any p9.File interface functions that return a File.
//
// For example, Walk() walks from a directory to a child.  We want to ensure the
// child is a p9Dev.File too, since the child is who we do the GetAttr call on.
// If we just called the underlying p9.File's walk, e.g. func (l *Local) Walk(),
// we'd get a p9.File that was e.g. &Local{path: l.path}.
type File struct {
	p9.File
	// Can't implicitly embed, since p9.File and DefaultWalkGetAttr both
	// implement WalkGetAttr.
	defaultWalker p9.DefaultWalkGetAttr
}

// GetAttr tells the client that special files (S_IFIFO, S_IFBLK, S_IFCHR) are
// regular files, allowing the Linux client to issue regular reads and writes.
func (f *File) GetAttr(req p9.AttrMask) (p9.QID, p9.AttrMask, p9.Attr, error) {
	qid, mask, attr, err := f.File.GetAttr(req)
	if err != nil {
		return qid, mask, attr, err
	}
	switch {
	case attr.Mode.IsCharacterDevice(), attr.Mode.IsBlockDevice(), attr.Mode.IsNamedPipe():
		attr.Mode = (attr.Mode & ^p9.FileModeMask) | p9.ModeRegular
	}
	return qid, mask, attr, err
}

// Mknod always fails, even if the underlying Attacher implements it (none of
// them do btw...).
//
// We explicitly do not want to let the client make arbitrary nodes (chardev,
// blockdev) on our filesystem.  It's one thing for us to export some special
// devices of our machine to the client.  But if we let them mknod, they can
// access any device we can access.
func (f *File) Mknod(name string, mode p9.FileMode, major uint32, minor uint32, _ p9.UID, _ p9.GID) (p9.QID, error) {
	return p9.QID{}, syscall.ENOSYS
}

// Walk wraps the underlying Walk, returning a p9dev.File.
func (f *File) Walk(names []string) ([]p9.QID, p9.File, error) {
	qids, newFile, err := f.File.Walk(names)
	if err != nil {
		return nil, nil, err
	}
	return qids, &File{File: newFile}, nil
}

// Create wraps the underlying Create, returning a p9dev.File.
func (f *File) Create(name string, mode p9.OpenFlags, perms p9.FileMode, uid p9.UID, gid p9.GID) (p9.File, p9.QID, uint32, error) {
	newFile, qid, i, err := f.File.Create(name, mode, perms, uid, gid)
	if err != nil {
		return nil, p9.QID{}, 0, err
	}
	return &File{File: newFile}, qid, i, nil
}

// WalkGetAttr calls the DefaultWalkGetAttr, i.e. return ENOSYS.  The p9 FS's
// we're wrapping don't implement this, so no need for us to either.
func (f *File) WalkGetAttr(names []string) ([]p9.QID, p9.File, p9.AttrMask, p9.Attr, error) {
	return f.defaultWalker.WalkGetAttr(names)
}
