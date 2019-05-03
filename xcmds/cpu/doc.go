// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cpu - connection to CPU server over SSH protocol
//
// Synopsis:
//     cpu [OPTIONS]
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
//     -h string
//           host to use (default "localhost")
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
//           what server to run (default "unpfs")
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
package main
