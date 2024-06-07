# Reimplementation of brctl 

## Status

### Implementation Subcommands

- [x] `addbr`: Add bridge
- [x] `delbr`: Delete bridge
- [x] `show`: Show current interfaces
- [x] `addif`: Add interface to bridge
- [x] `delif`: Detach interface from bridge
- [x] `showmacs`: Show list of MAC addresses for a bridge
- [x] `setageingtime`: Set ethernet MAC address ageing time
- [x] `stp`: Control Spanning Tree Protocol (STP)
- [x] `setbridgeprio`: Set bridge priority
- [x] `setfd`: Set 'forward delay' (FD)
- [x] `sethello`: Set 'bridge hello time'
- [x] `setmaxage`: Set 'maximum message age'
- [x] `setpathcost`: Set Port cost
- [x] `setportprio`: Set Port priority
- [x] `hairpin`: Enable/Disable hairpin mode

### Testing

- [x] unit tests for conversions
- [x] integration tests 
- [ ] `vmtest` integration test setup
    * [ ] update kernel image to incorporate necessary configs
    * [ ] make sure the env setup/cleanup works fine

### Tinygo

The long term goal is to build all the commands in u-root with `tinygo`.
This section of the document tracks the current with that:

```
ld.lld: error: undefined symbol: golang.org/x/sys/unix.Syscall
>>> referenced by zsyscall_linux.go:66 ($PROJECT/vendor/golang.org/x/sys/unix/zsyscall_linux.go:66)
>>>               $HOME/.cache/tinygo/thinlto/llvmcache-B03ED774292D0915EAB80DF33F191B7A8402F9E1:(golang.org/x/sys/unix.IoctlIfreq)

ld.lld: error: undefined symbol: golang.org/x/sys/unix.RawSyscall
>>> referenced by zsyscall_linux_amd64.go:480 ($PROJECT/vendor/golang.org/x/sys/unix/zsyscall_linux_amd64.go:480)
>>>               $HOME/.cache/tinygo/thinlto/llvmcache-B03ED774292D0915EAB80DF33F191B7A8402F9E1:(golang.org/x/sys/unix.Socket)
error: failed to link /tmp/tinygo1683407130/main: exit status 1
```

## Additional Information

* To use networking bridges, the kernel needs to have the according featured built into it. This is also necessary for the CI to run the tests.
    * `CONFIG_BRIDGE`
    * `CONFIG_NETLINK`
* The original C implementation offers the ability to issue the bridge configuration in 2 ways: 1. `ioctl` and 2. `sysfs`. Since all modern systems deploy the `sysfs` we use it to configure the bridges whenever possible. The option to manually switch which way to use is disabled.
