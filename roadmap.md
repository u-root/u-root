Roadmap (alpha)
======

#### Finish commands and tests

![Progress](http://progressed.io/bar/48)

- [x] cat
- [x] cmp
- [ ] comm
- [x] cp
- [x] date
- [ ] dd
- [ ] dhcp
- [ ] dmesg
- [x] echo
- [ ] ectool
- [x] freq
- [ ] gitclone
- [x] grep
- [x] hostname
- [ ] init
- [ ] installcommand
- [ ] ip
- [ ] kexec
- [ ] ldd
- [x] ln
- [ ] losetup
- [x] ls
- [ ] mkdir
- [ ] mount
- [x] mv
- [ ] pflask
- [ ] ping
- [x] printenv
- [x] ps
- [x] pwd
- [x] rm
- [ ] rush
- [ ] script
- [x] seq
- [ ] srvfiles
- [ ] tcz
- [ ] tee
- [x] uname
- [x] uniq
- [ ] unshare
- [x] wc
- [x] wget
- [x] which



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

