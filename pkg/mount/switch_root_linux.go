// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mount

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

// getDev returns the device (as returned by the FSTAT syscall) for the given file descriptor.
func getDev(fd int) (dev uint64, err error) {
	var stat unix.Stat_t
	if err := unix.Fstat(fd, &stat); err != nil {
		return 0, err
	}
	return uint64(stat.Dev), nil
}

// recursiveDelete deletes a directory identified by `fd` and everything in it.
//
// This function allows deleting directories no longer referenceable by
// any file name. This function does not descend into mounts.
//
// It is not an error for this function to fail to delete a file/directory.  That's normal (e.g. mount points).
// The objective is to clear as much memory as possible.
// In general this should be best-effort and should generally warn rather than fail,
// since failure leaves the system in an undefined state.
func recursiveDelete(fd int) error {
	parentDev, err := getDev(fd)
	if err != nil {
		log.Printf("warn: unable to get underlying dev for dir: %v", err)
		return nil
	}

	// The file descriptor is already open, but allocating a os.File
	// here makes reading the files in the dir so much nicer.
	dir := os.NewFile(uintptr(fd), "__ignored__")
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		log.Printf("warn: unable to read dir %s: %v", dir.Name(), err)
		return nil
	}

	for _, name := range names {
		// Loop here, but handle loop in separate function to make defer work as expected.
		if err := recusiveDeleteInner(fd, parentDev, name); err != nil {
			return err
		}
	}
	return nil
}

// recusiveDeleteInner is called from recursiveDelete and either deletes
// or recurses into the given file or directory
//
// There should be no need to call this function directly.
func recusiveDeleteInner(parentFd int, parentDev uint64, childName string) error {
	// O_DIRECTORY and O_NOFOLLOW make this open fail for all files and all symlinks (even when pointing to a dir).
	// We need to filter out symlinks because getDev later follows them.
	childFd, err := unix.Openat(parentFd, childName, unix.O_DIRECTORY|unix.O_NOFOLLOW, unix.O_RDWR)
	if err != nil {
		// childName points to either a file or a symlink, delete in any case.
		if err := unix.Unlinkat(parentFd, childName, 0); err != nil {
			log.Printf("warn: unable to remove file %s: %v", childName, err)
		}
	} else {
		// Open succeeded, which means childName points to a real directory.
		defer unix.Close(childFd)

		// Don't descend into other file systems.
		if childFdDev, err := getDev(childFd); err != nil {
			log.Printf("warn: unable to get underlying dev for dir: %s: %v", childName, err)
			return nil
		} else if childFdDev != parentDev {
			// This means continue in recursiveDelete.
			return nil
		}

		if err := recursiveDelete(childFd); err != nil {
			return err
		}
		// Back from recursion, the directory is now empty, delete.
		if err := unix.Unlinkat(parentFd, childName, unix.AT_REMOVEDIR); err != nil {
			log.Printf("warn: unable to remove dir %s: %v", childName, err)
		}
	}
	return nil
}

// MoveMount moves a mount from oldPath to newPath.
//
// This function is just a wrapper around the MOUNT syscall with the
// MOVE flag supplied.
func MoveMount(oldPath string, newPath string) error {
	return unix.Mount(oldPath, newPath, "", unix.MS_MOVE, "")
}

// addSpecialMounts moves the 'special' mounts to the given target path
//
// 'special' in this context refers to the following non-blockdevice backed
// mounts that are almost always used: /dev, /proc, /sys, and /run.
// This function will create the target directories, if necessary.
// If the target directories already exist, they must be empty.
// This function skips missing mounts.
func addSpecialMounts(newRoot string) error {
	mounts := []string{"/dev", "/proc", "/sys", "/run"}

	for _, mount := range mounts {
		path := filepath.Join(newRoot, mount)
		// Skip all mounting if the directory does not exist.
		if _, err := os.Stat(mount); os.IsNotExist(err) {
			log.Printf("switch_root: Skipping %q as the dir does not exist", mount)
			continue
		} else if err != nil {
			return err
		}
		// Also skip if not currently a mount point
		if same, err := SameFilesystem("/", mount); err != nil {
			return err
		} else if same {
			log.Printf("switch_root: Skipping %q as it is not a mount", mount)
			continue
		}
		// Make sure the target dir exists.
		if err := os.MkdirAll(path, 0o755); err != nil {
			return err
		}
		if err := MoveMount(mount, path); err != nil {
			return err
		}
	}
	return nil
}

// SameFilesystem returns true if both paths reside in the same filesystem.
// This is achieved by comparing Stat_t.Dev, which contains the fs device's
// major/minor numbers.
func SameFilesystem(path1, path2 string) (bool, error) {
	var stat1, stat2 unix.Stat_t
	if err := unix.Stat(path1, &stat1); err != nil {
		return false, err
	}
	if err := unix.Stat(path2, &stat2); err != nil {
		return false, err
	}
	return stat1.Dev == stat2.Dev, nil
}

// SwitchRoot makes newRootDir the new root directory of the system.
//
// To be exact, it makes newRootDir the new root directory of the calling
// process's mount namespace.
//
// It moves special mounts (dev, proc, sys, run) to the new directory, then
// does a chroot, moves the root mount to the new directory and finally
// DELETES EVERYTHING in the old root and execs the given init.
func SwitchRoot(newRootDir string, init string) error {
	err := newRoot(newRootDir)
	if err != nil {
		return err
	}
	return execInit(init)
}

// newRoot is the "first half" of SwitchRoot - that is, it creates special mounts
// in newRoot, chroot's there, and RECURSIVELY DELETES everything in the old root.
func newRoot(newRootDir string) error {
	log.Printf("switch_root: moving mounts")
	if err := addSpecialMounts(newRootDir); err != nil {
		return fmt.Errorf("switch_root: moving mounts failed %w", err)
	}

	log.Printf("switch_root: Changing directory")
	if err := unix.Chdir(newRootDir); err != nil {
		return fmt.Errorf("switch_root: failed change directory to new_root %w", err)
	}

	// Open "/" now, we need the file descriptor later.
	oldRoot, err := os.Open("/")
	if err != nil {
		return err
	}
	defer oldRoot.Close()

	log.Printf("switch_root: Moving /")
	if err := MoveMount(newRootDir, "/"); err != nil {
		return err
	}

	log.Printf("switch_root: Changing root!")
	if err := unix.Chroot("."); err != nil {
		return fmt.Errorf("switch_root: fatal chroot error %w", err)
	}

	log.Printf("switch_root: Deleting old /")
	return recursiveDelete(int(oldRoot.Fd()))
}

// execInit is generally only useful as part of SwitchRoot or similar.
// It exec's the given binary in place of the current binary, necessary so that
// the new binary can be pid 1.
func execInit(init string) error {
	log.Printf("switch_root: executing init")
	if err := unix.Exec(init, []string{init}, []string{}); err != nil {
		return fmt.Errorf("switch_root: exec failed %w", err)
	}
	return nil
}
