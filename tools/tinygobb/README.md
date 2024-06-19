# Status of tinygo in u-root

Currently, many of the busybox commands fail to build using tinygo.
This document aims to cover the current process of enabling all subcommands to be built using tinygo.
While enabling more and more commands the list of working commands in this document will be updated.

## building
Since currently the tinygo support is lacking, many commands got the `!tinygo` build tag.
This prohibits building, so anyone recreating the build steps should be aware that this has to be removed for each command.
The list below is the result of building each subcommand individually for x86_64 linux.

The neccessary additions to tinygo will be tracked in the according issue that this document is based on [#2979](https://github.com/u-root/u-root/issues/2979).

## command list
### core
41 tinygo build errors found.
 - [ ] `bind`
 - [ ] `df`
 - [ ] `dhclient`
 - [ ] `dmesg`
 - [ ] `fusermount`
 - [ ] `gosh`
 - [ ] `gpgv`
 - [ ] `gzip`
 - [ ] `hostname`
 - [ ] `hwclock`
 - [ ] `init`
 - [ ] `insmod`
 - [ ] `ip`
 - [ ] `kexec`
 - [ ] `lockmsrs`
 - [ ] `losetup`
 - [ ] `mkfifo`
 - [ ] `mknod`
 - [ ] `mount`
 - [ ] `msr`
 - [ ] `netcat`
 - [ ] `ntpdate`
 - [ ] `ping`
 - [ ] `poweroff`
 - [ ] `rmmod`
 - [ ] `shutdown`
 - [ ] `sluinit`
 - [ ] `sshd`
 - [ ] `strace`
 - [ ] `stty`
 - [ ] `switch_root`
 - [ ] `tee`
 - [ ] `truncate`
 - [ ] `umount`
 - [ ] `uname`
 - [ ] `watchdog`
 - [ ] `watchdogd`
 - [ ] `wget`
 - [ ] `which`
 - [ ] `nohup`

64 cmds build successful
- [x] `backoff`
- [x] `base64`
- [x] `basename`
- [x] `blkid`
- [x] `cat`
- [x] `chmod`
- [x] `chroot`
- [x] `cmp`
- [x] `comm`
- [x] `cp`
- [x] `cpio`
- [x] `date`
- [x] `dd`
- [x] `dirname`
- [x] `echo`
- [x] `false`
- [x] `find`
- [x] `free`
- [x] `gpt`
- [x] `grep`
- [x] `hexdump`
- [x] `id`
- [x] `io`
- [x] `kill`
- [x] `lddfiles`
- [x] `ln`
- [x] `ls`
- [x] `lsdrivers`
- [x] `lsmod`
- [x] `man`
- [x] `md5sum`
- [x] `mkdir`
- [x] `mktemp`
- [x] `more`
- [x] `mv`
- [x] `pci`
- [x] `printenv`
- [x] `ps`
- [x] `pwd`
- [x] `readlink`
- [x] `rm`
- [x] `rsdp`
- [x] `scp`
- [x] `seq`
- [x] `shasum`
- [x] `sleep`
- [x] `sort`
- [x] `strings`
- [x] `sync`
- [x] `tail`
- [x] `tar`
- [x] `time`
- [x] `timeout`
- [x] `touch`
- [x] `tr`
- [x] `true`
- [x] `ts`
- [x] `uniq`
- [x] `unmount`
- [x] `unshare`
- [x] `uptime`
- [x] `wc`
- [x] `xargs`
- [x] `yes`

### exp
30 tinygo build errors found.
 - [ ] `bootvars`
 - [ ] `bzimage`
 - [ ] `console`
 - [ ] `disk_unlock`
 - [ ] `efivarfs`
 - [ ] `esxiboot`
 - [ ] `getty`
 - [ ] `hdparm`
 - [ ] `ipmidump`
 - [ ] `kconf`
 - [ ] `modprobe`
 - [ ] `netbootxyz`
 - [ ] `newsshd`
 - [ ] `nvme_unlock`
 - [ ] `page`
 - [ ] `partprobe`
 - [ ] `pflask`
 - [ ] `pox`
 - [ ] `pxeserver`
 - [ ] `run`
 - [ ] `smbios_transfer`
 - [ ] `ssh`
 - [ ] `syscallfilter`
 - [ ] `uefiboot`
 - [ ] `vboot`
 - [ ] `vmboot`
 - [ ] `dumpmemmap`
 - [ ] `fbnetboot`
 - [ ] `localboot`
 - [ ] `systemboot`

27 cmds build successful
- [x] `acpicat`
- [x] `acpigrep`
- [x] `ansi`
- [x] `cbmem`
- [x] `crc`
- [x] `dmidecode`
- [x] `dumpebda`
- [x] `ectool`
- [x] `ed`
- [x] `fbsplash`
- [x] `fdtdump`
- [x] `field`
- [x] `fixrsdp`
- [x] `forth`
- [x] `freq`
- [x] `lsfabric`
- [x] `madeye`
- [x] `readelf`
- [x] `readpe`
- [x] `rush`
- [x] `smn`
- [x] `srvfiles`
- [x] `tac`
- [x] `tcz`
- [x] `watch`
- [x] `zbi`
- [x] `zimage`

