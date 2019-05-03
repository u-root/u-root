// Copyright 2015 Google Inc. All Rights Reserved.
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

package netfuse_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"path"
	"reflect"
	"syscall"
	"time"

	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"github.com/u-root/fuse/fusetesting"
	"github.com/u-root/fuse/samples"
	"github.com/u-root/u-root/pkg/netfuse"
)

////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////

type clntFSTest struct {
	samples.SampleTest
	cl netfuse.Clnt
}

func (t *clntFSTest) SetUp(ti *TestInfo) {
	ln, err := net.Listen("tcp", ":")
	AssertEq(nil, err)

	a := ln.Addr().String()
	// If you get lost, enabled this. log.Printf("Now dial %v", a)
	cl, err := rpc.Dial("tcp", a)
	AssertEq(nil, err)
	t.Server = netfuse.NewClntFS(cl)
	// This is a bit complex ...
	// We have to set up a fuse rpc server
	// then a fuse rpc client
	// then run the tests.

	// If you really need debugging this is where you get it.
	if false {
		t.MountConfig.DebugLogger = log.New(os.Stderr, "debug: ", 0)
		t.MountConfig.ErrorLogger = log.New(os.Stderr, "error: ", 0)
	}
	t.MountConfig.FSName = "netfuse"
	t.MountConfig.FSName = "netfuse"

	s, err := ln.Accept()
	AssertEq(nil, err)

	go func() {
		tempDir, err := ioutil.TempDir("", "netfuse")
		if err != nil {
			panic("fix this")
		}
		err = netfuse.Serve(tempDir, s)
		AssertEq(nil, err)
	}()

	t.SampleTest.SetUp(ti)
}

////////////////////////////////////////////////////////////////////////
// Basics
////////////////////////////////////////////////////////////////////////

type ClntFSTest struct {
	clntFSTest
}

// This is purely here to let us skip tests.
type nClntFSTest struct {
	clntFSTest
}

type neverClntFSTest struct {
	clntFSTest
}

func init() {
	// circleci doesn't have fusermount yet.
	if _, err := exec.LookPath("fusermount"); err != nil {
		log.Printf("Skipping FileFSTest since no fusermount")
		return
	}
	RegisterTestSuite(&ClntFSTest{})
}

func (t *ClntFSTest) ContentsOfEmptyFileSystem() {
	entries, err := fusetesting.ReadDirPicky(t.Dir)

	AssertEq(nil, err)
	ExpectThat(entries, ElementsAre())
}

func (t *ClntFSTest) Mkdir_OneLevel() {
	var err error
	var fi os.FileInfo
	var entries []os.FileInfo

	dirName := path.Join(t.Dir, "dir")

	// Create a directory within the root.
	err = os.Mkdir(dirName, 0754)
	AssertEq(nil, err)

	// Stat the directory.
	fi, err = os.Lstat(dirName)

	AssertEq(nil, err)
	ExpectEq("dir", fi.Name())
	ExpectEq(os.ModeDir|applyUmask(0754), fi.Mode())
	ExpectTrue(fi.IsDir())

	// Read the directory.
	entries, err = fusetesting.ReadDirPicky(dirName)

	AssertEq(nil, err)
	ExpectThat(entries, ElementsAre())

	// Read the root.
	entries, err = fusetesting.ReadDirPicky(t.Dir)

	AssertEq(nil, err)
	AssertEq(1, len(entries))

	fi = entries[0]
	ExpectEq("dir", fi.Name())
	ExpectEq(os.ModeDir|applyUmask(0754), fi.Mode())
}

func (t *ClntFSTest) Mkdir_TwoLevels() {
	var err error
	var fi os.FileInfo
	var stat *syscall.Stat_t
	var entries []os.FileInfo

	// Create a directory within the root.
	err = os.Mkdir(path.Join(t.Dir, "parent"), 0700)
	AssertEq(nil, err)

	// Create a child of that directory.
	err = os.Mkdir(path.Join(t.Dir, "parent/dir"), 0754)
	AssertEq(nil, err)

	// Stat the directory.
	fi, err = os.Lstat(path.Join(t.Dir, "parent/dir"))
	stat = fi.Sys().(*syscall.Stat_t)

	AssertEq(nil, err)
	ExpectEq("dir", fi.Name())
	// This is too fs and system-dependent. ExpectEq(0, fi.Size())
	ExpectEq(os.ModeDir|applyUmask(0754), fi.Mode())
	// These time tests make no real sense in a networked world.
	ExpectTrue(fi.IsDir())

	ExpectEq(2, stat.Nlink)

	fi, err = os.Lstat(path.Join(t.Dir, "parent"))
	AssertEq(nil, err)

	// Read the directory.
	entries, err = fusetesting.ReadDirPicky(path.Join(t.Dir, "parent/dir"))

	AssertEq(nil, err)
	ExpectThat(entries, ElementsAre())

	// Read the parent.
	entries, err = fusetesting.ReadDirPicky(path.Join(t.Dir, "parent"))

	AssertEq(nil, err)
	AssertEq(1, len(entries))

	fi = entries[0]
	ExpectEq("dir", fi.Name())
	ExpectEq(os.ModeDir|applyUmask(0754), fi.Mode())
}

func (t *ClntFSTest) Mkdir_AlreadyExists() {
	var err error
	dirName := path.Join(t.Dir, "dir")

	// Create the directory once.
	err = os.Mkdir(dirName, 0754)
	AssertEq(nil, err)

	// Attempt to create it again.
	err = os.Mkdir(dirName, 0754)

	AssertNe(nil, err)
	ExpectThat(err, Error(HasSubstr("exists")))
}

func (t *ClntFSTest) Mkdir_IntermediateIsFile() {
	var err error

	// Create a file.
	fileName := path.Join(t.Dir, "foo")
	err = ioutil.WriteFile(fileName, []byte{}, 0700)
	AssertEq(nil, err)

	// Attempt to create a directory within the file.
	dirName := path.Join(fileName, "dir")
	err = os.Mkdir(dirName, 0754)

	AssertNe(nil, err)
	ExpectThat(err, Error(HasSubstr("not a directory")))
}

func (t *ClntFSTest) Mkdir_IntermediateIsNonExistent() {
	var err error

	// Attempt to create a sub-directory of a non-existent sub-directory.
	dirName := path.Join(t.Dir, "foo/dir")
	err = os.Mkdir(dirName, 0754)

	AssertNe(nil, err)
	ExpectThat(err, Error(HasSubstr("no such file or directory")))
}

func (t *ClntFSTest) Mkdir_PermissionDenied() {
	var err error

	// This test fails for root
	if os.Getuid() == 0 {
		return
	}
	// Create a directory within the root without write permissions.
	err = os.Mkdir(path.Join(t.Dir, "parent"), 0500)
	AssertEq(nil, err)

	// Attempt to create a child of that directory.
	err = os.Mkdir(path.Join(t.Dir, "parent/dir"), 0754)

	AssertNe(nil, err)
	ExpectThat(err, Error(HasSubstr("ermission denied")))
}

func (t *ClntFSTest) CreateNewFile_InRoot() {
	var err error
	var fi os.FileInfo
	var stat *syscall.Stat_t

	// Write a file.
	fileName := path.Join(t.Dir, "foo")
	const contents = "Hello\x00world"

	createTime := time.Now()
	err = ioutil.WriteFile(fileName, []byte(contents), 0400)
	AssertEq(nil, err)

	// Stat it.
	fi, err = os.Lstat(fileName)
	stat = fi.Sys().(*syscall.Stat_t)

	AssertEq(nil, err)
	ExpectEq("foo", fi.Name())
	ExpectEq(len(contents), fi.Size())
	ExpectEq(applyUmask(0400), fi.Mode())
	ExpectThat(fi, fusetesting.MtimeIsWithin(createTime, timeSlop))
	ExpectThat(fi, fusetesting.BirthtimeIsWithin(createTime, timeSlop))
	ExpectFalse(fi.IsDir())

	ExpectNe(0, stat.Ino)
	ExpectEq(1, stat.Nlink)
	ExpectEq(currentUid(), stat.Uid)
	ExpectEq(currentGid(), stat.Gid)
	ExpectEq(len(contents), stat.Size)

	// Read it back.
	slice, err := ioutil.ReadFile(fileName)
	AssertEq(nil, err)
	ExpectEq(contents, string(slice))
}

func (t *ClntFSTest) CreateNewFile_InSubDir() {
	var err error
	var fi os.FileInfo
	var stat *syscall.Stat_t

	// Create a sub-dir.
	dirName := path.Join(t.Dir, "dir")
	err = os.Mkdir(dirName, 0700)
	AssertEq(nil, err)

	// Write a file.
	fileName := path.Join(dirName, "foo")
	const contents = "Hello\x00world"

	//createTime := time.Now()
	err = ioutil.WriteFile(fileName, []byte(contents), 0400)
	AssertEq(nil, err)

	// Stat it.
	fi, err = os.Lstat(fileName)
	stat = fi.Sys().(*syscall.Stat_t)

	AssertEq(nil, err)
	ExpectEq("foo", fi.Name())
	ExpectEq(len(contents), fi.Size())
	ExpectEq(applyUmask(0400), fi.Mode())
	// ExpectThat(fi, fusetesting.MtimeIsWithin(createTime, timeSlop))
	// ExpectThat(fi, fusetesting.BirthtimeIsWithin(createTime, timeSlop))
	ExpectFalse(fi.IsDir())

	ExpectNe(0, stat.Ino)
	ExpectEq(1, stat.Nlink)
	ExpectEq(currentUid(), stat.Uid)
	ExpectEq(currentGid(), stat.Gid)
	ExpectEq(len(contents), stat.Size)

	// Read it back.
	slice, err := ioutil.ReadFile(fileName)
	AssertEq(nil, err)
	ExpectEq(contents, string(slice))
}

func (t *ClntFSTest) ModifyExistingFile_InRoot() {
	var err error
	var n int
	var fi os.FileInfo
	var stat *syscall.Stat_t

	// Write a file.
	fileName := path.Join(t.Dir, "foo")

	createTime := time.Now()
	err = ioutil.WriteFile(fileName, []byte("Hello, world!"), 0600)
	AssertEq(nil, err)

	// Open the file and modify it.
	f, err := os.OpenFile(fileName /*os.O_WRONLY,*/, os.O_RDWR, 0400)
	t.ToClose = append(t.ToClose, f)
	AssertEq(nil, err)

	modifyTime := time.Now()
	n, err = f.WriteAt([]byte("H"), 0)
	AssertEq(nil, err)
	AssertEq(1, n)

	// Stat the file.
	fi, err = os.Lstat(fileName)
	stat = fi.Sys().(*syscall.Stat_t)

	AssertEq(nil, err)
	ExpectEq("foo", fi.Name())
	ExpectEq(len("Hello, world!"), fi.Size())
	ExpectEq(applyUmask(0600), fi.Mode())
	ExpectThat(fi, fusetesting.MtimeIsWithin(modifyTime, timeSlop))
	ExpectThat(fi, fusetesting.BirthtimeIsWithin(createTime, timeSlop))
	ExpectFalse(fi.IsDir())

	ExpectNe(0, stat.Ino)
	ExpectEq(1, stat.Nlink)
	ExpectEq(currentUid(), stat.Uid)
	ExpectEq(currentGid(), stat.Gid)
	ExpectEq(len("Hello, world!"), stat.Size)

	// Read the file back.
	slice, err := ioutil.ReadFile(fileName)
	AssertEq(nil, err)
	ExpectEq("Hello, world!", string(slice))
}

func (t *ClntFSTest) ModifyExistingFile_InSubDir() {
	var err error
	var n int
	var fi os.FileInfo
	var stat *syscall.Stat_t

	// Create a sub-directory.
	dirName := path.Join(t.Dir, "dir")
	err = os.Mkdir(dirName, 0700)
	AssertEq(nil, err)

	// Write a file.
	fileName := path.Join(dirName, "foo")

	createTime := time.Now()
	err = ioutil.WriteFile(fileName, []byte("Hello, world!"), 0600)
	AssertEq(nil, err)

	// Open the file and modify it.
	f, err := os.OpenFile(fileName, os.O_RDWR /*os.O_WRONLY*/, 0400)
	t.ToClose = append(t.ToClose, f)
	AssertEq(nil, err)

	modifyTime := time.Now()
	n, err = f.WriteAt([]byte("H"), 0)
	AssertEq(nil, err)
	AssertEq(1, n)

	// Stat the file.
	fi, err = os.Lstat(fileName)
	stat = fi.Sys().(*syscall.Stat_t)

	AssertEq(nil, err)
	ExpectEq("foo", fi.Name())
	ExpectEq(len("Hello, world!"), fi.Size())
	ExpectEq(applyUmask(0600), fi.Mode())
	ExpectThat(fi, fusetesting.MtimeIsWithin(modifyTime, timeSlop))
	ExpectThat(fi, fusetesting.BirthtimeIsWithin(createTime, timeSlop))
	ExpectFalse(fi.IsDir())

	ExpectNe(0, stat.Ino)
	ExpectEq(1, stat.Nlink)
	ExpectEq(currentUid(), stat.Uid)
	ExpectEq(currentGid(), stat.Gid)
	ExpectEq(len("Hello, world!"), stat.Size)

	// Read the file back.
	slice, err := ioutil.ReadFile(fileName)
	AssertEq(nil, err)
	ExpectEq("Hello, world!", string(slice))
}

func (t *ClntFSTest) UnlinkFile_Exists() {
	var err error

	// Write a file.
	fileName := path.Join(t.Dir, "foo")
	err = ioutil.WriteFile(fileName, []byte("Hello, world!"), 0600)
	AssertEq(nil, err)

	// Unlink it.
	err = os.Remove(fileName)
	AssertEq(nil, err)

	// Statting it should fail.
	_, err = os.Lstat(fileName)

	AssertNe(nil, err)
	ExpectThat(err, Error(HasSubstr("no such file")))

	// Nothing should be in the directory.
	entries, err := fusetesting.ReadDirPicky(t.Dir)
	AssertEq(nil, err)
	ExpectThat(entries, ElementsAre())
}

func (t *ClntFSTest) UnlinkFile_NonExistent() {
	err := os.Remove(path.Join(t.Dir, "foo"))

	AssertNe(nil, err)
	ExpectThat(err, Error(HasSubstr("no such file")))
}

// The idea that this works is actually questionable. If you are mounting from a
// Plan 9 9p server, for example, it won't. We skip this test.
func (t *neverClntFSTest) UnlinkFile_StillOpen() {
	fileName := path.Join(t.Dir, "foo")

	// Create and open a file.
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0600)
	t.ToClose = append(t.ToClose, f)
	AssertEq(nil, err)

	// Write some data into it.
	n, err := f.Write([]byte("taco"))
	AssertEq(nil, err)
	AssertEq(4, n)

	// Unlink it.
	err = os.Remove(fileName)
	AssertEq(nil, err)

	// The directory should no longer contain it.
	entries, err := fusetesting.ReadDirPicky(t.Dir)
	AssertEq(nil, err)
	ExpectThat(entries, ElementsAre())

	// We should be able to stat the file. It should still show as having
	// contents, but with no links.
	fi, err := f.Stat()

	AssertEq(nil, err)
	ExpectEq(4, fi.Size())
	ExpectEq(0, fi.Sys().(*syscall.Stat_t).Nlink)

	// The contents should still be available.
	buf := make([]byte, 1024)
	n, err = f.ReadAt(buf, 0)

	AssertEq(io.EOF, err)
	AssertEq(4, n)
	ExpectEq("taco", string(buf[:4]))

	// Writing should still work, too.
	n, err = f.Write([]byte("burrito"))
	AssertEq(nil, err)
	AssertEq(len("burrito"), n)
}

func (t *ClntFSTest) Rmdir_NonEmpty() {
	var err error

	// Create two levels of directories.
	err = os.MkdirAll(path.Join(t.Dir, "foo/bar"), 0754)
	AssertEq(nil, err)

	// Attempt to remove the parent.
	err = os.Remove(path.Join(t.Dir, "foo"))

	AssertNe(nil, err)
	// Fuse does not always return sensible values? ExpectThat(err, Error(HasSubstr("not empty")))
}

func (t *ClntFSTest) Rmdir_Empty() {
	var err error
	var entries []os.FileInfo

	// Create two levels of directories.
	err = os.MkdirAll(path.Join(t.Dir, "foo/bar"), 0754)
	AssertEq(nil, err)

	// Remove the leaf.
	err = os.Remove(path.Join(t.Dir, "foo/bar"))
	AssertEq(nil, err)

	// There should be nothing left in the parent.
	entries, err = fusetesting.ReadDirPicky(path.Join(t.Dir, "foo"))

	AssertEq(nil, err)
	ExpectThat(entries, ElementsAre())

	// Check the parent's mtime.
	_, err = os.Lstat(path.Join(t.Dir, "foo"))
	AssertEq(nil, err)

	// Remove the parent.
	err = os.Remove(path.Join(t.Dir, "foo"))
	AssertEq(nil, err)

	// Now the root directory should be empty, too.
	entries, err = fusetesting.ReadDirPicky(t.Dir)

	AssertEq(nil, err)
	ExpectThat(entries, ElementsAre())
}

func (t *ClntFSTest) Rmdir_NonExistent() {
	err := os.Remove(path.Join(t.Dir, "blah"))

	AssertNe(nil, err)
	ExpectThat(err, Error(HasSubstr("no such file or directory")))
}

// This is very specific to some Unix file systems. Don't bother.
func (t *neverClntFSTest) Rmdir_OpenedForReading() {
	var err error

	// Create a directory.
	createTime := time.Now()
	err = os.Mkdir(path.Join(t.Dir, "dir"), 0700)
	AssertEq(nil, err)

	// Open the directory for reading.
	f, err := os.Open(path.Join(t.Dir, "dir"))
	defer func() {
		if f != nil {
			ExpectEq(nil, f.Close())
		}
	}()

	AssertEq(nil, err)

	// Remove the directory.
	err = os.Remove(path.Join(t.Dir, "dir"))
	AssertEq(nil, err)

	// Create a new directory, with the same name even, and add some contents
	// within it.
	err = os.MkdirAll(path.Join(t.Dir, "dir/foo"), 0700)
	AssertEq(nil, err)

	err = os.MkdirAll(path.Join(t.Dir, "dir/bar"), 0700)
	AssertEq(nil, err)

	err = os.MkdirAll(path.Join(t.Dir, "dir/baz"), 0700)
	AssertEq(nil, err)

	// We should still be able to stat the open file handle. It should show up as
	// unlinked.
	fi, err := f.Stat()
	AssertEq(nil, err)
	ExpectEq("dir", fi.Name())
	ExpectThat(fi, fusetesting.MtimeIsWithin(createTime, timeSlop))
	ExpectEq(0, fi.Sys().(*syscall.Stat_t).Nlink)

	// Attempt to read from the directory. This shouldn't see any junk from the
	// new directory. It should either succeed with an empty result or should
	// return ENOENT.
	names, err := f.Readdirnames(0)

	if err != nil {
		ExpectThat(err, Error(HasSubstr("no such file")))
	} else {
		ExpectThat(names, ElementsAre())
	}
}

func (t *ClntFSTest) CaseSensitive() {
	var err error

	// Create a file.
	err = ioutil.WriteFile(path.Join(t.Dir, "file"), []byte{}, 0400)
	AssertEq(nil, err)

	// Create a directory.
	err = os.Mkdir(path.Join(t.Dir, "dir"), 0400)
	AssertEq(nil, err)

	// Attempt to stat with the wrong case.
	names := []string{
		"FILE",
		"File",
		"filE",
		"DIR",
		"Dir",
		"dIr",
	}

	for _, name := range names {
		_, err = os.Lstat(path.Join(t.Dir, name))
		AssertNe(nil, err, "Name: %s", name)
		AssertThat(err, Error(HasSubstr("no such file or directory")))
	}
}

func (t *ClntFSTest) WriteOverlapsEndOfFile() {
	var err error
	var n int

	// Create a file.
	f, err := os.Create(path.Join(t.Dir, "foo"))
	t.ToClose = append(t.ToClose, f)
	AssertEq(nil, err)

	// Make it 4 bytes long.
	err = f.Truncate(4)
	AssertEq(nil, err)

	// Write the range [2, 6).
	n, err = f.WriteAt([]byte("taco"), 2)
	AssertEq(nil, err)
	AssertEq(4, n)

	// Read the full contents of the file.
	contents, err := ioutil.ReadAll(f)
	AssertEq(nil, err)
	ExpectEq("\x00\x00taco", string(contents))
}

func (t *ClntFSTest) WriteStartsAtEndOfFile() {
	var err error
	var n int

	// Create a file.
	f, err := os.Create(path.Join(t.Dir, "foo"))
	t.ToClose = append(t.ToClose, f)
	AssertEq(nil, err)

	// Make it 2 bytes long.
	err = f.Truncate(2)
	AssertEq(nil, err)

	// Write the range [2, 6).
	n, err = f.WriteAt([]byte("taco"), 2)
	AssertEq(nil, err)
	AssertEq(4, n)

	// Read the full contents of the file.
	contents, err := ioutil.ReadAll(f)
	AssertEq(nil, err)
	ExpectEq("\x00\x00taco", string(contents))
}

func (t *ClntFSTest) WriteStartsPastEndOfFile() {
	var err error
	var n int

	// Create a file.
	f, err := os.Create(path.Join(t.Dir, "foo"))
	t.ToClose = append(t.ToClose, f)
	AssertEq(nil, err)

	// Write the range [2, 6).
	n, err = f.WriteAt([]byte("taco"), 2)
	AssertEq(nil, err)
	AssertEq(4, n)

	// Read the full contents of the file.
	contents, err := ioutil.ReadAll(f)
	AssertEq(nil, err)
	ExpectEq("\x00\x00taco", string(contents))
}

func (t *ClntFSTest) WriteAtDoesntChangeOffset_NotAppendMode() {
	var err error
	var n int

	// Create a file.
	f, err := os.Create(path.Join(t.Dir, "foo"))
	t.ToClose = append(t.ToClose, f)
	AssertEq(nil, err)

	// Make it 16 bytes long.
	err = f.Truncate(16)
	AssertEq(nil, err)

	// Seek to offset 4.
	_, err = f.Seek(4, 0)
	AssertEq(nil, err)

	// Write the range [10, 14).
	n, err = f.WriteAt([]byte("taco"), 2)
	AssertEq(nil, err)
	AssertEq(4, n)

	// We should still be at offset 4.
	offset, err := getFileOffset(f)
	AssertEq(nil, err)
	ExpectEq(4, offset)
}

func (t *ClntFSTest) WriteAtDoesntChangeOffset_AppendMode() {
	var err error
	var n int

	// Create a file in append mode.
	f, err := os.OpenFile(
		path.Join(t.Dir, "foo"),
		os.O_RDWR|os.O_APPEND|os.O_CREATE,
		0600)

	t.ToClose = append(t.ToClose, f)
	AssertEq(nil, err)

	// Make it 16 bytes long.
	err = f.Truncate(16)
	AssertEq(nil, err)

	// Seek to offset 4.
	_, err = f.Seek(4, 0)
	AssertEq(nil, err)

	// Write the range [10, 14).
	n, err = f.WriteAt([]byte("taco"), 2)
	AssertEq(nil, err)
	AssertEq(4, n)

	// We should still be at offset 4.
	offset, err := getFileOffset(f)
	AssertEq(nil, err)
	ExpectEq(4, offset)
}

func (t *ClntFSTest) LargeFile() {
	var err error

	// Create a file.
	f, err := os.Create(path.Join(t.Dir, "foo"))
	t.ToClose = append(t.ToClose, f)
	AssertEq(nil, err)

	// Copy in large contents.
	const size = 1 << 24
	contents := bytes.Repeat([]byte{0x20}, size)

	_, err = io.Copy(f, bytes.NewReader(contents))
	AssertEq(nil, err)

	// Read the full contents of the file.
	contents, err = ioutil.ReadFile(f.Name())
	AssertEq(nil, err)
	ExpectEq(size, len(contents))
}

func (t *ClntFSTest) AppendMode() {
	var err error
	var n int
	var off int64
	buf := make([]byte, 1024)

	// Create a file with some contents.
	fileName := path.Join(t.Dir, "foo")
	err = ioutil.WriteFile(fileName, []byte("Jello, "), 0600)
	AssertEq(nil, err)

	// Open the file in append mode.
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0600)
	t.ToClose = append(t.ToClose, f)
	AssertEq(nil, err)

	// Seek to somewhere silly and then write.
	off, err = f.Seek(2, 0)
	AssertEq(nil, err)
	AssertEq(2, off)

	n, err = f.Write([]byte("world!"))
	AssertEq(nil, err)
	AssertEq(6, n)

	// The offset should have been updated to point at the end of the file.
	off, err = getFileOffset(f)
	AssertEq(nil, err)
	ExpectEq(13, off)

	// A random write should still work, without updating the offset.
	n, err = f.WriteAt([]byte("H"), 0)
	AssertEq(nil, err)
	AssertEq(1, n)

	off, err = getFileOffset(f)
	AssertEq(nil, err)
	ExpectEq(13, off)

	// Read back the contents of the file, which should be correct even though we
	// seeked to a silly place before writing the world part.
	//
	// Linux's support for pwrite is buggy; the pwrite(2) man page says this:
	//
	//     POSIX requires that opening a file with the O_APPEND flag should have
	//     no affect on the location at which pwrite() writes data.  However, on
	//     Linux,  if  a  file  is opened with O_APPEND, pwrite() appends data to
	//     the end of the file, regardless of the value of offset.
	//
	// So we allow either the POSIX result or the Linux result.
	n, err = f.ReadAt(buf, 0)
	AssertEq(io.EOF, err)
	ExpectThat(string(buf[:n]), AnyOf("Hello, world!", "Jello, world!H"))
}

func (t *ClntFSTest) ReadsPastEndOfFile() {
	var err error
	var n int
	buf := make([]byte, 1024)

	// Create a file.
	f, err := os.Create(path.Join(t.Dir, "foo"))
	t.ToClose = append(t.ToClose, f)
	AssertEq(nil, err)

	// Give it some contents.
	n, err = f.Write([]byte("taco"))
	AssertEq(nil, err)
	AssertEq(4, n)

	// Read a range overlapping EOF.
	n, err = f.ReadAt(buf[:4], 2)
	AssertEq(io.EOF, err)
	ExpectEq(2, n)
	ExpectEq("co", string(buf[:n]))

	// Read a range starting at EOF.
	n, err = f.ReadAt(buf[:4], 4)
	AssertEq(io.EOF, err)
	ExpectEq(0, n)
	ExpectEq("", string(buf[:n]))

	// Read a range starting past EOF.
	n, err = f.ReadAt(buf[:4], 100)
	AssertEq(io.EOF, err)
	ExpectEq(0, n)
	ExpectEq("", string(buf[:n]))
}

func (t *ClntFSTest) Truncate_Smaller() {
	var err error
	fileName := path.Join(t.Dir, "foo")

	// Create a file.
	err = ioutil.WriteFile(fileName, []byte("taco"), 0600)
	AssertEq(nil, err)

	// Open it for modification.
	f, err := os.OpenFile(fileName, os.O_RDWR, 0)
	t.ToClose = append(t.ToClose, f)
	AssertEq(nil, err)

	// Truncate it.
	err = f.Truncate(2)
	AssertEq(nil, err)

	// Stat it.
	fi, err := f.Stat()
	AssertEq(nil, err)
	ExpectEq(2, fi.Size())

	// Read the contents.
	contents, err := ioutil.ReadFile(fileName)
	AssertEq(nil, err)
	ExpectEq("ta", string(contents))
}

func (t *ClntFSTest) Truncate_SameSize() {
	var err error
	fileName := path.Join(t.Dir, "foo")

	// Create a file.
	err = ioutil.WriteFile(fileName, []byte("taco"), 0600)
	AssertEq(nil, err)

	// Open it for modification.
	f, err := os.OpenFile(fileName, os.O_RDWR, 0)
	t.ToClose = append(t.ToClose, f)
	AssertEq(nil, err)

	// Truncate it.
	err = f.Truncate(4)
	AssertEq(nil, err)

	// Stat it.
	fi, err := f.Stat()
	AssertEq(nil, err)
	ExpectEq(4, fi.Size())

	// Read the contents.
	contents, err := ioutil.ReadFile(fileName)
	AssertEq(nil, err)
	ExpectEq("taco", string(contents))
}

func (t *ClntFSTest) Truncate_Larger() {
	var err error
	fileName := path.Join(t.Dir, "foo")

	// Create a file.
	err = ioutil.WriteFile(fileName, []byte("taco"), 0600)
	AssertEq(nil, err)

	// Open it for modification.
	f, err := os.OpenFile(fileName, os.O_RDWR, 0)
	t.ToClose = append(t.ToClose, f)
	AssertEq(nil, err)

	// Truncate it.
	err = f.Truncate(6)
	AssertEq(nil, err)

	// Stat it.
	fi, err := f.Stat()
	AssertEq(nil, err)
	ExpectEq(6, fi.Size())

	// Read the contents.
	contents, err := ioutil.ReadFile(fileName)
	AssertEq(nil, err)
	ExpectEq("taco\x00\x00", string(contents))
}

func (t *ClntFSTest) Chmod() {
	var err error
	fileName := path.Join(t.Dir, "foo")

	// Create a file.
	err = ioutil.WriteFile(fileName, []byte(""), 0600)
	AssertEq(nil, err)

	// Chmod it.
	err = os.Chmod(fileName, 0754)
	AssertEq(nil, err)

	// Stat it.
	fi, err := os.Lstat(fileName)
	AssertEq(nil, err)
	ExpectEq(0754, fi.Mode())
}

func (t *ClntFSTest) Chtimes() {
	var err error
	fileName := path.Join(t.Dir, "foo")

	// Create a file.
	err = ioutil.WriteFile(fileName, []byte(""), 0600)
	AssertEq(nil, err)

	// Chtimes it.
	expectedMtime := time.Now().Add(123 * time.Second).Round(time.Second)
	err = os.Chtimes(fileName, time.Now(), expectedMtime)
	AssertEq(nil, err)

	// Stat it.
	fi, err := os.Lstat(fileName)
	AssertEq(nil, err)
	ExpectThat(fi, fusetesting.MtimeIsWithin(expectedMtime, timeSlop))
}

// This one is also dicey enough for many file systems that we'll let it go
func (t *neverClntFSTest) ReadDirWhileModifying() {
	dirName := path.Join(t.Dir, "dir")
	createFile := func(name string) {
		AssertEq(nil, ioutil.WriteFile(path.Join(dirName, name), []byte{}, 0400))
	}

	// Create a directory.
	err := os.Mkdir(dirName, 0700)
	AssertEq(nil, err)

	// Open the directory.
	d, err := os.Open(dirName)
	t.ToClose = append(t.ToClose, d)
	AssertEq(nil, err)

	// Add four files.
	createFile("foo")
	createFile("bar")
	createFile("baz")
	createFile("qux")

	// Read one entry from the directory.
	names, err := d.Readdirnames(1)
	AssertEq(nil, err)
	AssertThat(names, ElementsAre("foo"))

	// Make two holes in the directory.
	AssertEq(nil, os.Remove(path.Join(dirName, "foo")))
	AssertEq(nil, os.Remove(path.Join(dirName, "baz")))

	// Add a bunch of files to the directory.
	createFile("blah_0")
	createFile("blah_1")
	createFile("blah_2")
	createFile("blah_3")
	createFile("blah_4")

	// Continue reading from the directory, noting the names we see.
	namesSeen := make(map[string]bool)
	for {
		names, err = d.Readdirnames(1)
		for _, n := range names {
			namesSeen[n] = true
		}

		if err == io.EOF {
			break
		}

		AssertEq(nil, err)
	}

	// Posix requires that we should have seen bar and qux, which we didn't
	// delete.
	ExpectTrue(namesSeen["bar"])
	ExpectTrue(namesSeen["qux"])
}

func (t *ClntFSTest) CreateSymlink() {
	var fi os.FileInfo
	var err error

	symlinkName := path.Join(t.Dir, "foo")
	target := "taco/burrito"

	// Create the link.
	err = os.Symlink(target, symlinkName)
	AssertEq(nil, err)

	// Stat the link.
	fi, err = os.Lstat(symlinkName)
	AssertEq(nil, err)

	ExpectEq("foo", fi.Name())
	ExpectEq(os.ModeSymlink, fi.Mode()&os.ModeType)

	// Read the link.
	actual, err := os.Readlink(symlinkName)
	AssertEq(nil, err)
	ExpectEq(target, actual)

	// Read the parent directory.
	entries, err := fusetesting.ReadDirPicky(t.Dir)
	AssertEq(nil, err)
	AssertEq(1, len(entries))

	fi = entries[0]
	ExpectEq("foo", fi.Name())
	ExpectEq(os.ModeSymlink, fi.Mode()&os.ModeType)
}

func (t *ClntFSTest) CreateSymlink_AlreadyExists() {
	var err error

	// Create a file and a directory.
	fileName := path.Join(t.Dir, "foo")
	err = ioutil.WriteFile(fileName, []byte{}, 0400)
	AssertEq(nil, err)

	dirName := path.Join(t.Dir, "bar")
	err = os.Mkdir(dirName, 0700)
	AssertEq(nil, err)

	// Create an existing symlink.
	symlinkName := path.Join(t.Dir, "baz")
	err = os.Symlink("blah", symlinkName)
	AssertEq(nil, err)

	// Symlinking on top of any of them should fail.
	names := []string{
		fileName,
		dirName,
		symlinkName,
	}

	for _, n := range names {
		err = os.Symlink("blah", n)
		AssertNe(nil, err)
		//ExpectThat(err, Error(HasSubstr("exists")))
	}
}

func (t *ClntFSTest) ReadLink_NonExistent() {
	_, err := os.Readlink(path.Join(t.Dir, "foo"))
	ExpectTrue(os.IsNotExist(err), "err: %v", err)
}

func (t *ClntFSTest) ReadLink_NotASymlink() {
	var err error

	// Create a file and a directory.
	fileName := path.Join(t.Dir, "foo")
	err = ioutil.WriteFile(fileName, []byte{}, 0400)
	AssertEq(nil, err)

	dirName := path.Join(t.Dir, "bar")
	err = os.Mkdir(dirName, 0700)
	AssertEq(nil, err)

	// Reading either of them as a symlink should fail.
	names := []string{
		fileName,
		dirName,
	}

	for _, n := range names {
		_, err = os.Readlink(n)
		ExpectThat(err, Error(HasSubstr("invalid argument")))
	}
}

func (t *ClntFSTest) DeleteSymlink() {
	var err error

	symlinkName := path.Join(t.Dir, "foo")
	target := "taco/burrito"

	// Create the link.
	err = os.Symlink(target, symlinkName)
	AssertEq(nil, err)

	// Remove it.
	err = os.Remove(symlinkName)
	AssertEq(nil, err)

	// Statting should now fail.
	_, err = os.Lstat(symlinkName)
	ExpectTrue(os.IsNotExist(err), "err: %v", err)

	// Read the parent directory.
	entries, err := fusetesting.ReadDirPicky(t.Dir)

	AssertEq(nil, err)
	ExpectThat(entries, ElementsAre())
}

func (t *ClntFSTest) CreateHardlink() {
	var fi os.FileInfo
	var err error

	// Create a file.
	fileName := path.Join(t.Dir, "regular_file")
	const contents = "Hello\x00world"

	err = ioutil.WriteFile(fileName, []byte(contents), 0444)
	AssertEq(nil, err)

	// Clean up the file at the end.
	defer func() {
		if false {
			err := os.Remove(fileName)
			AssertEq(nil, err)
		}
	}()

	// Create a link to the file.
	linkName := path.Join(t.Dir, "foo")
	err = os.Link(fileName, linkName)
	AssertEq(nil, err)

	// Clean up the file at the end.
	defer func() {
		if false {
			err := os.Remove(linkName)
			AssertEq(nil, err)
		}
	}()

	// Stat the link.
	fi, err = os.Lstat(linkName)
	AssertEq(nil, err)

	ExpectEq("foo", fi.Name())
	ExpectEq(0444, fi.Mode())

	// Read the parent directory.
	entries, err := fusetesting.ReadDirPicky(t.Dir)
	AssertEq(nil, err)
	AssertEq(2, len(entries))

	fi = entries[0]
	ExpectEq("foo", fi.Name())
	ExpectEq(0444, fi.Mode())

	fi = entries[1]
	ExpectEq("regular_file", fi.Name())
	ExpectEq(0444, fi.Mode())
}

func (t *ClntFSTest) CreateHardlink_AlreadyExists() {
	var err error

	// Create a file and a directory.
	fileName := path.Join(t.Dir, "foo")
	err = ioutil.WriteFile(fileName, []byte{}, 0400)
	AssertEq(nil, err)

	dirName := path.Join(t.Dir, "bar")
	err = os.Mkdir(dirName, 0700)
	AssertEq(nil, err)

	// Create an existing symlink.
	symlinkName := path.Join(t.Dir, "baz")
	err = os.Symlink("blah", symlinkName)
	AssertEq(nil, err)

	// Create another link to the file.
	hardlinkName := path.Join(t.Dir, "qux")
	err = os.Link(fileName, hardlinkName)
	AssertEq(nil, err)

	// Symlinking on top of any of them should fail.
	names := []string{
		fileName,
		dirName,
		symlinkName,
		hardlinkName,
	}

	for _, n := range names {
		err = os.Link(fileName, n)
		ExpectThat(err, Error(HasSubstr("exists")))
	}
}

// This test is fairly bogus. There's not guaranteed reverse map from inumber to file name.
// It's a one way trip.
func (t *neverClntFSTest) DeleteHardlink() {
	var fi os.FileInfo
	var err error

	// Create a file.
	fileName := path.Join(t.Dir, "regular_file")
	const contents = "Hello\x00world"

	err = ioutil.WriteFile(fileName, []byte(contents), 0444)
	AssertEq(nil, err)

	// Step #1: We will create and remove a link and verify that
	// after removal everything is as expected.

	// Create a link to the file.
	linkName := path.Join(t.Dir, "foo")
	err = os.Link(fileName, linkName)
	AssertEq(nil, err)

	// Remove the link.
	err = os.Remove(linkName)
	AssertEq(nil, err)

	// Stat the link.
	fi, err = os.Lstat(linkName)
	AssertEq(nil, fi)
	ExpectThat(err, Error(HasSubstr("no such file")))

	// Read the parent directory.
	entries, err := fusetesting.ReadDirPicky(t.Dir)
	AssertEq(nil, err)
	AssertEq(1, len(entries))

	fi = entries[0]
	ExpectEq("regular_file", fi.Name())
	ExpectEq(0444, fi.Mode())

	// Step #2: We will create a link and remove the original file subsequently
	// and verify that after removal everything is as expected.

	// Create a link to the file.
	linkName = path.Join(t.Dir, "bar")
	err = os.Link(fileName, linkName)
	AssertEq(nil, err)

	// Remove the original file.
	err = os.Remove(fileName)
	AssertEq(nil, err)

	// Stat the link.
	fi, err = os.Lstat(linkName)
	AssertEq(nil, err)
	ExpectEq("bar", fi.Name())
	ExpectEq(0444, fi.Mode())

	// Stat the original file.
	fi, err = os.Lstat(fileName)
	AssertEq(nil, fi)
	ExpectThat(err, Error(HasSubstr("no such file")))

	// Read the parent directory.
	entries, err = fusetesting.ReadDirPicky(t.Dir)
	AssertEq(nil, err)
	AssertEq(1, len(entries))

	fi = entries[0]
	ExpectEq("bar", fi.Name())
	ExpectEq(0444, fi.Mode())

	// Cleanup.
	err = os.Remove(linkName)
	AssertEq(nil, err)
}

// same here. inumber to path just doesn't work this way.
func (t *neverClntFSTest) ReadHardlink() {
	var err error

	// Create a file.
	fileName := path.Join(t.Dir, "regular_file")
	const contents = "Hello\x00world"

	err = ioutil.WriteFile(fileName, []byte(contents), 0444)
	AssertEq(nil, err)

	// Clean up the file at the end.
	defer func() {
		err := os.Remove(fileName)
		AssertEq(nil, err)
	}()

	// Create a link to the file.
	linkName := path.Join(t.Dir, "foo")
	err = os.Link(fileName, linkName)
	AssertEq(nil, err)

	// Clean up the file at the end.
	defer func() {
		err := os.Remove(linkName)
		AssertEq(nil, err)
	}()

	// Read files.
	original, err := ioutil.ReadFile(fileName)
	AssertEq(nil, err)
	linked, err := ioutil.ReadFile(linkName)
	AssertEq(nil, err)

	// Check if the bytes are the same.
	AssertEq(true, reflect.DeepEqual(original, linked))
}

func (t *neverClntFSTest) CreateInParallel_NoTruncate() {
	fusetesting.RunCreateInParallelTest_NoTruncate(t.Ctx, t.Dir)
}

func (t *neverClntFSTest) CreateInParallel_Truncate() {
	fusetesting.RunCreateInParallelTest_Truncate(t.Ctx, t.Dir)
}

func (t *neverClntFSTest) CreateInParallel_Exclusive() {
	fusetesting.RunCreateInParallelTest_Exclusive(t.Ctx, t.Dir)
}

func (t *neverClntFSTest) MkdirInParallel() {
	fusetesting.RunMkdirInParallelTest(t.Ctx, t.Dir)
}

func (t *neverClntFSTest) SymlinkInParallel() {
	fusetesting.RunSymlinkInParallelTest(t.Ctx, t.Dir)
}

func (t *neverClntFSTest) HardlinkInParallel() {
	fusetesting.RunHardlinkInParallelTest(t.Ctx, t.Dir)
}

func (t *neverClntFSTest) RenameWithinDir_File() {
	var err error

	// Create a parent directory.
	parentPath := path.Join(t.Dir, "parent")

	err = os.Mkdir(parentPath, 0700)
	AssertEq(nil, err)

	// And a file within it.
	oldPath := path.Join(parentPath, "foo")

	err = ioutil.WriteFile(oldPath, []byte("taco"), 0400)
	AssertEq(nil, err)

	// Rename it.
	newPath := path.Join(parentPath, "bar")

	err = os.Rename(oldPath, newPath)
	AssertEq(nil, err)

	// The old name shouldn't work.
	_, err = os.Lstat(oldPath)
	ExpectTrue(os.IsNotExist(err), "err: %v", err)

	_, err = ioutil.ReadFile(oldPath)
	ExpectTrue(os.IsNotExist(err), "err: %v", err)

	// The new name should.
	fi, err := os.Lstat(newPath)
	AssertEq(nil, err)
	ExpectEq(len("taco"), fi.Size())
	ExpectEq(os.FileMode(0400), fi.Mode())

	contents, err := ioutil.ReadFile(newPath)
	AssertEq(nil, err)
	ExpectEq("taco", string(contents))

	// There should only be the new entry in the directory.
	entries, err := fusetesting.ReadDirPicky(parentPath)
	AssertEq(nil, err)
	AssertEq(1, len(entries))
	fi = entries[0]

	ExpectEq(path.Base(newPath), fi.Name())
	ExpectEq(os.FileMode(0400), fi.Mode())
}

func (t *neverClntFSTest) RenameWithinDir_Directory() {
	var err error

	// Create a parent directory.
	parentPath := path.Join(t.Dir, "parent")

	err = os.Mkdir(parentPath, 0700)
	AssertEq(nil, err)

	// And a non-empty directory within it.
	oldPath := path.Join(parentPath, "foo")

	err = os.MkdirAll(path.Join(oldPath, "child"), 0700)
	AssertEq(nil, err)

	// Rename it.
	newPath := path.Join(parentPath, "bar")

	err = os.Rename(oldPath, newPath)
	AssertEq(nil, err)

	// The old name shouldn't work.
	_, err = os.Lstat(oldPath)
	ExpectTrue(os.IsNotExist(err), "err: %v", err)

	// The new name should.
	fi, err := os.Lstat(newPath)
	AssertEq(nil, err)
	ExpectEq(os.FileMode(0700)|os.ModeDir, fi.Mode())

	// There should only be the new entry in the parent.
	entries, err := fusetesting.ReadDirPicky(parentPath)
	AssertEq(nil, err)
	AssertEq(1, len(entries))
	fi = entries[0]

	ExpectEq(path.Base(newPath), fi.Name())
	ExpectEq(os.FileMode(0700)|os.ModeDir, fi.Mode())

	// And the child should still be present.
	entries, err = fusetesting.ReadDirPicky(newPath)
	AssertEq(nil, err)
	AssertEq(1, len(entries))
	fi = entries[0]

	ExpectEq("child", fi.Name())
	ExpectEq(os.FileMode(0700)|os.ModeDir, fi.Mode())
}

func (t *neverClntFSTest) RenameWithinDir_SameName() {
	var err error

	// Create a parent directory.
	parentPath := path.Join(t.Dir, "parent")

	err = os.Mkdir(parentPath, 0700)
	AssertEq(nil, err)

	// And a file within it.
	filePath := path.Join(parentPath, "foo")

	err = ioutil.WriteFile(filePath, []byte("taco"), 0400)
	AssertEq(nil, err)

	// Attempt to rename it.
	err = os.Rename(filePath, filePath)
	AssertEq(nil, err)

	// The file should still exist.
	contents, err := ioutil.ReadFile(filePath)
	AssertEq(nil, err)
	ExpectEq("taco", string(contents))

	// There should only be the one entry in the directory.
	entries, err := fusetesting.ReadDirPicky(parentPath)
	AssertEq(nil, err)
	AssertEq(1, len(entries))
	fi := entries[0]

	ExpectEq(path.Base(filePath), fi.Name())
	ExpectEq(os.FileMode(0400), fi.Mode())
}

func (t *neverClntFSTest) RenameAcrossDirs_File() {
	var err error

	// Create two parent directories.
	oldParentPath := path.Join(t.Dir, "old")
	newParentPath := path.Join(t.Dir, "new")

	err = os.Mkdir(oldParentPath, 0700)
	AssertEq(nil, err)

	err = os.Mkdir(newParentPath, 0700)
	AssertEq(nil, err)

	// And a file within the first.
	oldPath := path.Join(oldParentPath, "foo")

	err = ioutil.WriteFile(oldPath, []byte("taco"), 0400)
	AssertEq(nil, err)

	// Rename it.
	newPath := path.Join(newParentPath, "bar")

	err = os.Rename(oldPath, newPath)
	AssertEq(nil, err)

	// The old name shouldn't work.
	_, err = os.Lstat(oldPath)
	ExpectTrue(os.IsNotExist(err), "err: %v", err)

	_, err = ioutil.ReadFile(oldPath)
	ExpectTrue(os.IsNotExist(err), "err: %v", err)

	// The new name should.
	fi, err := os.Lstat(newPath)
	AssertEq(nil, err)
	ExpectEq(len("taco"), fi.Size())
	ExpectEq(os.FileMode(0400), fi.Mode())

	contents, err := ioutil.ReadFile(newPath)
	AssertEq(nil, err)
	ExpectEq("taco", string(contents))

	// Check the old parent.
	entries, err := fusetesting.ReadDirPicky(oldParentPath)
	AssertEq(nil, err)
	AssertEq(0, len(entries))

	// And the new one.
	entries, err = fusetesting.ReadDirPicky(newParentPath)
	AssertEq(nil, err)
	AssertEq(1, len(entries))
	fi = entries[0]

	ExpectEq(path.Base(newPath), fi.Name())
	ExpectEq(os.FileMode(0400), fi.Mode())
}

func (t *neverClntFSTest) RenameAcrossDirs_Directory() {
	var err error

	// Create two parent directories.
	oldParentPath := path.Join(t.Dir, "old")
	newParentPath := path.Join(t.Dir, "new")

	err = os.Mkdir(oldParentPath, 0700)
	AssertEq(nil, err)

	err = os.Mkdir(newParentPath, 0700)
	AssertEq(nil, err)

	// And a non-empty directory within the first.
	oldPath := path.Join(oldParentPath, "foo")

	err = os.MkdirAll(path.Join(oldPath, "child"), 0700)
	AssertEq(nil, err)

	// Rename it.
	newPath := path.Join(newParentPath, "bar")

	err = os.Rename(oldPath, newPath)
	AssertEq(nil, err)

	// The old name shouldn't work.
	_, err = os.Lstat(oldPath)
	ExpectTrue(os.IsNotExist(err), "err: %v", err)

	// The new name should.
	fi, err := os.Lstat(newPath)
	AssertEq(nil, err)
	ExpectEq(os.FileMode(0700)|os.ModeDir, fi.Mode())

	// And the child should still be present.
	entries, err := fusetesting.ReadDirPicky(newPath)
	AssertEq(nil, err)
	AssertEq(1, len(entries))
	fi = entries[0]

	ExpectEq("child", fi.Name())
	ExpectEq(os.FileMode(0700)|os.ModeDir, fi.Mode())

	// Check the old parent.
	entries, err = fusetesting.ReadDirPicky(oldParentPath)
	AssertEq(nil, err)
	AssertEq(0, len(entries))

	// And the new one.
	entries, err = fusetesting.ReadDirPicky(newParentPath)
	AssertEq(nil, err)
	AssertEq(1, len(entries))
	fi = entries[0]

	ExpectEq(path.Base(newPath), fi.Name())
	ExpectEq(os.FileMode(0700)|os.ModeDir, fi.Mode())
}

func (t *neverClntFSTest) RenameOutOfFileSystem() {
	var err error

	// Create a file.
	oldPath := path.Join(t.Dir, "foo")

	err = ioutil.WriteFile(oldPath, []byte("taco"), 0400)
	AssertEq(nil, err)

	// Attempt to move it out of the file system.
	tempDir, err := ioutil.TempDir("", "fs_test")
	AssertEq(nil, err)
	defer os.RemoveAll(tempDir)

	err = os.Rename(oldPath, path.Join(tempDir, "bar"))
	ExpectThat(err, Error(HasSubstr("cross-device")))
}

func (t *ClntFSTest) RenameIntoFileSystem() {
	var err error

	// Create a file outside of our file system.
	f, err := ioutil.TempFile("", "fs_test")
	AssertEq(nil, err)
	defer f.Close()

	oldPath := f.Name()
	defer os.Remove(oldPath)

	// Attempt to move it into the file system.
	err = os.Rename(oldPath, path.Join(t.Dir, "bar"))
	ExpectThat(err, Error(HasSubstr("cross-device")))
}

func (t *ClntFSTest) RenameOverExistingFile() {
	var err error

	// Create two files.
	oldPath := path.Join(t.Dir, "foo")
	err = ioutil.WriteFile(oldPath, []byte("taco"), 0400)
	AssertEq(nil, err)

	newPath := path.Join(t.Dir, "bar")
	err = ioutil.WriteFile(newPath, []byte("burrito"), 0600)
	AssertEq(nil, err)

	// Rename one over the other.
	err = os.Rename(oldPath, newPath)
	AssertEq(nil, err)

	// Check the file contents.
	contents, err := ioutil.ReadFile(newPath)
	AssertEq(nil, err)
	ExpectEq("taco", string(contents))

	// And the parent listing.
	entries, err := fusetesting.ReadDirPicky(t.Dir)
	AssertEq(nil, err)
	AssertEq(1, len(entries))
	fi := entries[0]

	ExpectEq(path.Base(newPath), fi.Name())
	ExpectEq(os.FileMode(0400), fi.Mode())
	ExpectEq(len("taco"), fi.Size())
}

func (t *ClntFSTest) RenameOverExistingDirectory() {
	var err error

	// Create two directories, the first non-empty.
	oldPath := path.Join(t.Dir, "foo")
	err = os.MkdirAll(path.Join(oldPath, "child"), 0700)
	AssertEq(nil, err)

	newPath := path.Join(t.Dir, "bar")
	err = os.Mkdir(newPath, 0600)
	AssertEq(nil, err)

	// Renaming over the non-empty directory shouldn't work.
	err = os.Rename(newPath, oldPath)
	ExpectThat(err, Error(MatchesRegexp("not empty|file exists")))

	// As of Go 1.8 this shouldn't work the other way around either (see
	// https://github.com/golang/go/commit/321c312).
	if atLeastGo18 {
		err = os.Rename(oldPath, newPath)
		ExpectThat(err, Error(HasSubstr("file exists")))

		// Both should still be present in the parent listing.
		entries, err := fusetesting.ReadDirPicky(t.Dir)
		AssertEq(nil, err)
		ExpectEq(2, len(entries))
	}
}

func (t *neverClntFSTest) RenameOverExisting_WrongType() {
	var err error

	// Create a file and a directory.
	filePath := path.Join(t.Dir, "foo")
	err = ioutil.WriteFile(filePath, []byte("taco"), 0400)
	AssertEq(nil, err)

	dirPath := path.Join(t.Dir, "bar")
	err = os.Mkdir(dirPath, 0700)
	AssertEq(nil, err)

	// Renaming either over the other shouldn't work.
	err = os.Rename(filePath, dirPath)
	ExpectThat(err, Error(MatchesRegexp("is a directory|file exists")))

	err = os.Rename(dirPath, filePath)
	ExpectThat(err, Error(HasSubstr("not a directory")))
}

func (t *ClntFSTest) RenameNonExistentFile() {
	var err error

	err = os.Rename(path.Join(t.Dir, "foo"), path.Join(t.Dir, "bar"))
	ExpectThat(err, Error(HasSubstr("no such file")))
}
