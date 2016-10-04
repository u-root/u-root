Roadmap (alpha)
======

#### Finish commands and tests

- [x] date
- [x] cat
- [x] rm
- [x] mv
- [x] tee
- [ ] freq
- [ ] ping
- [ ] ps
- [ ] sed
- [ ] which
- [ ] uniq
- [ ] seq
- [ ] sh
- [ ] grep
- [ ] wc
- [ ] cp
- [ ] script
- [ ] ip
- [ ] cmp
- [ ] comm
- [ ] dd
- [ ] dmesg
- [ ] echo
- [ ] hostname
- [ ] init
- [ ] ldd
- [ ] ls
- [ ] mkdir
- [ ] mount
- [ ] printenv
- [x] pwd
- [ ] uname
- [ ] wget
- [ ] dhcp
- [ ] ectool
- [ ] gitclone
- [ ] installcommand
- [ ] kexec
- [ ] losetup
- [ ] pflask
- [ ] srvfiles
- [ ] tcz3
- [ ] unshare



#### New Goal
- [ ] Get enough basic commands working to support a container mechanism.
- [ ] Determine what commands we might need for "New ChromeOS"
- [ ] Bring in Go readline package for the u-root shell
- [ ] Finish implementation of the ip command

#### Figure out a container solution
Options:

* Docker
* Rocket
* wget + unpack (cpio? tar?) + u-root pflask
* implement a gitclone command and use u-root pflask

