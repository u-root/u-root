// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command mkuinitramfs builds a u-root initramfs from the Go module cache,
// without a source checkout.
//
//	mkuinitramfs -o initramfs.cpio github.com/u-root/u-root/cmds/core/...@v0.16.0
//
// It writes a module whose only purpose is to require what the patterns name,
// which is what turns a package pattern into the module that provides it.
// That module stands in for the checkout: it holds no source, only the go.mod
// and go.sum recording what the patterns resolved to, which is enough for the
// go command to list the packages and build them. The source stays in the
// module cache.
//
// The archive is assembled with u-root's own pkg/cpio rather than a copy of
// it, so the fixes that land there reach this tool.
//
// # Installing
//
// This command type-checks your source in-process rather than shelling out
// to the go tool for it, so the Go it is built with has to be at least as
// new as the Go every module in the build declares. Build it with a current
// toolchain:
//
//	GOTOOLCHAIN=go1.27.2 go install github.com/u-root/u-root/tools/mkuinitramfs@latest
//
// GOTOOLCHAIN matters here and nowhere else: it selects the Go that compiles
// the command. Setting it when running the command has no effect on this,
// because by then the type-checker is already compiled in.
//
// The toolchain chosen here is the one the command then uses for everything:
// fetching, listing and compiling all run the go that sits beside the
// type-checker, rather than whatever go is on PATH. So the command needs no
// go on PATH at all, and there is no second toolchain to disagree with the
// one this line picks.
//
// # Arguments
//
// Each argument is a package pattern, optionally followed by @ and a version
// query. The query is anything the go command accepts: a version, a branch, a
// commit, or latest.
//
//	mkuinitramfs -o initramfs.cpio \
//		github.com/u-root/u-root/cmds/core/...@v0.16.0
//
// Commands from other modules come from naming them, and they are built into
// the same binary and linked into /bbin like any other:
//
//	mkuinitramfs -o initramfs.cpio \
//		github.com/u-root/u-root/cmds/core/...@main \
//		golang.org/x/tools/gopls@latest
//
//	mkuinitramfs -o initramfs.cpio \
//		github.com/u-root/u-root/cmds/core/...@main \
//		github.com/microsoft/typescript-go/cmd/tsgo@latest
//
// A trailing /... takes a whole subtree, which can cost more than it looks.
// Against the same cmds/core, naming perkeep's cmd/pk gives a 22 MB archive
// and naming cmd/... gives 64 MB, because five more commands come with it:
//
//	mkuinitramfs -o initramfs.cpio \
//		github.com/u-root/u-root/cmds/core/...@main \
//		perkeep.org/cmd/pk@latest
//
// # When it goes wrong
//
// Four failures account for most of what users hit. The first three are not
// this command refusing to work; they are the Go build rules showing through
// it.
//
// # The command itself is too old
//
//	package requires newer Go version go1.27 (application built with go1.26)
//
// "application" means mkuinitramfs, not the package being built. Check what
// it was built with and what it carries, then rebuild it:
//
//	go version ./mkuinitramfs
//	go version -m ./mkuinitramfs
//
// Running it under a newer GOTOOLCHAIN will not help, for the reason under
// Installing. A binary left lying around from an earlier checkout is the
// usual cause.
//
// # A module in the build declares an old Go version
//
//	implicit function instantiation requires go1.18 or later
//
// Or any other language-version error in a module you did not name. That
// module declares an old go directive in its go.mod and its own source has
// outgrown it. This usually appears when one argument forces a shared
// dependency to a newer version: the go command takes the highest version
// anything asks for, and a module pinned to an old language version cannot
// always compile against it.
//
// The go command rejects the same combination and says so more clearly, so
// ask it directly. Build with -keep-tmp, then in the module it kept:
//
//	go build <the package named in the error>
//
// It names the language version it used, which is the part to act on. The
// fix is to choose versions that agree: a newer version of the module with
// the old directive, or an older version of whatever forced the upgrade.
//
// # A build that worked yesterday now fails
//
// A query of @latest resolves to the newest semver tag, not the tip of the
// default branch. A repository whose newest tag is months behind its main
// branch gives you the tag, including whatever dependencies it had then.
// Use @main when you want the branch.
//
//	go list -m github.com/u-root/u-root@latest
//	go list -m github.com/u-root/u-root@main
//
// These print different versions whenever the tag is not the tip. Only if a
// module has no tags at all does @latest mean the tip.
//
// # The archive has no init
//
//	cannot link /init to "init": no such command was built
//
// Nothing would run. u-root's cmds/core provides init; a set of patterns
// naming only other modules does not. Add a pattern that supplies it, name a
// different command with -initcmd, or pass -initcmd= to leave /init out. The
// same applies to -defaultsh and /bin/sh.
package main
