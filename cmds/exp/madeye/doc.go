// Copyright 2013-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// u-root was intended to be capable of function as a universal root, i.e. a root
// file system that you could boot from different architectures.
// We call this Multiple Architecture Device Image, or MADI, pronounced Mad-Eye.
// Apologies to Harry Potter.
// This command implements MADI and works as follows:
// Given a set of images, e.g.
// initramfs.linux_<arch>.cpio it derives the architecture from the name.
// It then reads the cpio in.
// For all symlinks which resolve to absolute paths, it records the file they point to.
// For a distinguished set of directories, it relocates them from / to /<arch>/, a la Plan 9.
// If there is a /init, it moves to /<arch>/init.
// It adjusts absolute path symlinks.
// future: it looks for conflicting dev entries.
// It then writes it out.
// To boot a kernel with a MadEye, one must adjust the init= arg to prepend the architecture.
// For example, /init on arm would become /arm/init
// For now, this only works for bb mode.
package main
