// Copyright 2019 the u-root Authors. All rights reserved
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
package main
