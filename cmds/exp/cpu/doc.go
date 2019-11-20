// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cpu - connection to CPU server over SSH protocol
//
// Synopsis:
//     cpu [OPTIONS]
//
// Advisory:
//     cpu connects to a remote machine and serves a namespace to it.
//     That namespace can be a restricted file system (recommended)
//     or everything up to and including your /.
//     We do cover a lot of the simpler adversarial attacks, via
//     private name space mounts and the nonce on the 9p socket, and a
//     quick review of this code by an expert suggest that on systems
//     on which your adversary would only be those running at the same
//     privilege level as you, you may be ok.
//
//     You should not use this command to connect to a host that might
//     have an untrusted, privileged adversary, i.e. someone who might
//     replace the remote version with a corrupted version, or might
//     use a different attack to grab the nonce and mount your file
//     system, changing files before the client cpu exits.
//.
//     cpu is best used to connect to a LinuxBoot machine with only one
//     user, i.e. a non-time-shared machine.
//
//     We are developing a cpu we consider more trustworthy but that is a ways off.
//     You Have Been Warned :-)
//
// Description:
//     On local machine useful flags are all save -remote
//     On remote machines, all save dbg9p, key, hostkey.
//
//     CPU is an ssh client that starts up a shell on a remote machine,
//     as usual; but, further, makes a namespace of the local machine
//     available in a private mount rooted at /tmp/cpu.
//     Wherever you go, there your files are.
//     You can, in the ssh session, do something like this:
//     chroot /tmp/cpu /bin/bash
//     and at that point, you are running a bash, imported from your local
//     machine, on the remote machine; it will use your .profile and
//     all your files are available. You can also do something like
//     cat /tmp/cpu/etc/hosts
//     if your host file is lacking; or
//     cp /etc/hosts /tmp/cpu/tmp
//     to get the /etc/hosts on the remote machine to your local machine.
//
//     The cpu client makes this work by starting a cpu command on the
//     remote machine with a -remote switch and several other arguments.
//     The local cpu starts a 9p server and, using ssh port forwarding,
//     makes that server available to the remote. On the remote side, the
//     cpu command establishes a private, unshared mount of tmpfs on /tmp;
//     creates /tmp/cpu; and mounts the 9p server on that directory.
//
//     CPU has many options, as shown above; most you need not worry about.
//     The most common invocation is
//     cpu -h hostname
//     which will start a shell and mount the 9p server in /tmp/cpu.
//     Note this mount proceeds over the ssh session, and further
//     it mounts in a private /tmp; there is little to see when
//     it is running from outside the ssh session
//
// Options:
//     -bin string
//           path of cpu binary
//     -d    enable debug prints
//     -dbg9p
//           show 9p io
//     -hostkey string
//           host key file
//     -key string
//           key file (default "$HOME/.ssh/cpu_rsa")
//     -network string
//           network to use (default "tcp")
//     -p string
//           port to use (default "22")
//     -port9p string
//           port9p # on remote machine for 9p mount
//     -remote
//           Indicates we are the remote side of the cpu session
//     -srv string
//           what server to run (default none; use internal)
// Examples
// In these examples, cpu runs with warning messages enabled.
// The first message is a warning that cpu could not use overlayfs to build a
// a reasonable union mount. The next are showing you what it is mounting, the
// union mount having failed.
// These mounts are the best we could do for a reasonable compromise of
// wanting local resources visible (e.g. /dev) and using resources from the
// remote machine (e.g. /etc, /lib, /usr and so on).
// u-root doesn't really need /lib and /usr, and u-root's /etc is minimal by design,
// so this works.
// Also note that the user's 9p server running on the local machine is mounted at /tmp/cpu.
// We can turn these off at some point but for now, in early days, we may want them.
// Note that these messages come from the remote side.
// cpu to a machine with bash as your shell and run a command
//   cpu -sp 23 date
//     2019/05/17 16:53:22 Overlayfs mount failed: invalid argument. Proceeding with selective mounts from /tmp/cpu into /
//     2019/05/17 16:53:22 Mounted /tmp/cpu/lib on /lib
//     2019/05/17 16:53:22 Mounted /tmp/cpu/lib64 on /lib64
//     2019/05/17 16:53:22 Warning: mounting /tmp/cpu/lib32 on /lib32 failed: no such file or directory
//     2019/05/17 16:53:22 Mounted /tmp/cpu/usr on /usr
//     2019/05/17 16:53:22 Mounted /tmp/cpu/bin on /bin
//     2019/05/17 16:53:22 Mounted /tmp/cpu/etc on /etc
//     Fri May 17 16:53:23 UTC 2019
// cpu to a machine and run $SHELL (since no arguments were given)
// NOTE: $SHELL is NOT installed on the remote machine! It (and all its .so's and . files)
// come from the local machine.
// cpu sp -23
//    2019/05/17 16:58:04 Overlayfs mount failed: invalid argument. Proceeding with selective mounts from /tmp/cpu into /
//    2019/05/17 16:58:04 Mounted /tmp/cpu/lib on /lib
//    2019/05/17 16:58:04 Mounted /tmp/cpu/lib64 on /lib64
//    2019/05/17 16:58:04 Warning: mounting /tmp/cpu/lib32 on /lib32 failed: no such file or directory
//    2019/05/17 16:58:04 Mounted /tmp/cpu/usr on /usr
//    2019/05/17 16:58:04 Mounted /tmp/cpu/bin on /bin
//    2019/05/17 16:58:05 Mounted /tmp/cpu/etc on /etc
//    root@(none):/# echo ~
//    /tmp/cpu/home/rminnich
//    root@(none):/# ls ~
//    IDAPROPASSWORD  go      ida-7.2  projects            salishan2019random~
//    bin             gopath  papers   salishan2019random  snap
//    root@(none):/#
//    # Now that we are on the node, modprobe something
//    root@(none):/# depmod
//    depmod: ERROR: could not open directory /lib/modules/5.0.0-rc3+: No such file or directory
//    depmod: FATAL: could not search modules: No such file or directory
//    root@(none):/#
//    # Note that, if we had the right modules on our LOCAL machine for this remote machine, we could
//    # insert them. This further means you can build a modular kernel in FLASH and insmod needed modules
//    # later (as long as your core kernel has networking, that is!). Modules could include, e.g., an AHCI
//    # driver.
//    # run the lspci command but redirect output to ~
//    # note it is not installed on the remote machine; it comes from our local machine.
//    root@(none):/# lspci
//    00:00.0 Host bridge: Intel Corporation 82G33/G31/P35/P31 Express DRAM Controller
//    00:01.0 VGA compatible controller: Device 1234:1111 (rev 02)
//    00:02.0 Unclassified device [00ff]: Red Hat, Inc. Virtio RNG
//    00:03.0 Ethernet controller: Intel Corporation 82540EM Gigabit Ethernet Controller (rev 03)
//    00:1f.0 ISA bridge: Intel Corporation 82801IB (ICH9) LPC Interface Controller (rev 02)
//    00:1f.2 SATA controller: Intel Corporation 82801IR/IO/IH (ICH9R/DO/DH) 6 port SATA Controller [AHCI mode] (rev 02)
//    00:1f.3 SMBus: Intel Corporation 82801I (ICH9 Family) SMBus Controller (rev 02)
//    root@(none):/# lspci > ~/xyz
//    root@(none):/# exit
//    # exit and notice that file is on my local machine now:
//    exit
//    rminnich@xcpu:~/gopath/src/github.com/u-root/u-root/xcmds/cpu$ ls -l ~/xyz
//    -rw-r--r-- 1 rminnich rminnich 577 May 17 17:06 /home/rminnich/xyz
//    rminnich@xcpu:~/gopath/src/github.com/u-root/u-root/xcmds/cpu$
//
// A note on sizes of things. We can get an image down to 3 MiB if the only
// binary is cpu. See github.com:linuxboot/mainboards/aeeon/up for an example.
//
// You may see that the remote cpu writes a nonce back to the client, and wonder why
// that is. This gets back to how we make the client name space available to remote processes.
// To export the client name space to the remote cpu proc
// we use a remote ssh forward so the remote kernel can mount
// our built-in server over 9p. But how do we ensure the mount comes
// from the remote cpu proc and not something else?
// We need to ensure that the socket has only
// been opened by the remote cpu process.
// The way we do it here is the result of discussions with Eric Grosse, among others.
// Not that he endorses this solution, but he raised awareness
// of the issues.
//
// The 'remote' invocation of CPU is designed to work under
// cpu acting as init OR an sshd (i.e. ssh host cpu -r ...)
// Anything that makes the cpu init mode special is out; we would also have
// to change sshd to do the same things and we can't change
// every ssh server in the world. For anoher example, using Linux
// private network name spaces for the ssh forward is out,
// as not all ssh daemons support that.
// Trying to use ssh-forwarded Unix domain sockets is out, as we would want them placed in the
// private name space of the remote cpu process, and no ssh daemons support that; further
// the Go ssh packages don't support Unix domain sockets. Even were we to add them (I did)
// the approach fails as we can't specify that they be in a private name space.
// We'd rather not do a full TLS transaction on the forwarded port,
// since we for sure don't want another layer of encryption running
// over our already-encrypted, already-authenticated ssh connection.
//
// In the last year we've experimented with several approaches that all failed
// due to limits in how ssh protocol works, how the ssh packages work, how various ssh daemons work,
// or how Linux manages name spaces. What we describe below is the simplest approach
// that works on everything we've tried.
//
// We only do one accept on the client side.
// For the kernel 9p mount at the remote, we use the fd transport.
// The fd we pass is for a socket that has been verified.
// We verify as follows:
// client generates a nonce and adds it as an environment variable (not argv!).
// We don't put the nonce in argv as an adversary could then see it in ps.
// The remote cpu process reads that variable and removes it from the environment.
// The remote cpu process writes that nonce back on the port forward socket within 10 ms.
// The 10ms requirement is to defend (imperfectly) against a bad person running the remote
// under ptrace control.
// The client only accepts one connect.
//
// We did experiment with writing the nonce to stdin of the remote cpu
// process, but that was not reliable, so we went with the environment.
//
// This is not perfect, but it's a lot better than having the kernel
// mount from a socket.
//
// Note that in the original plan 9 cpu, all the communications over the
// socket were 9p; further, the remote cpu recreates the name space of the
// client by performing mounts as needed, not shuttling 9p requests back to the
// cpu client. In other words, this cpu is a workalike, but is implemented
// very differently. The use of 9p for all operations to the remote cpu sounds
// nice but does not solve all problems; in particular, as delay-bandwidth products
// get larger and larger, 9p's poor behavior becomes more apparent. Further, plan 9
// cpu had weird issues with EOF, manifested in pipelines:
// cpu host ls | cpu otherhost wc
// would hang as often as it worked, which is why plan 9 had a remote execute
// command as well as cpu. This cpu implementation acts a lot more like what ssh
// users are used to.
//
// If you want to learn more read the factotum paper by Grosse et. al. and the original
// Plan 9 papers on cpu, but be aware there are many subtle details visible only
// in code.
package main
