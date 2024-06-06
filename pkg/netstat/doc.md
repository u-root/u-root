# Reimplementation of netstat
## Status quo
Main capabilities (x marks implemented)

- [x] List socket states
- [x] List route table (flags: -r, -r4,-r6,-r46)
- [x] List multicast group membership (flags: -g, -g4, -g6, -g46)
- [x] Print interface information(s) (flags: -i, -I=)
- [x] Print network statistics (flags: -s, -s4, -s6)

Formating flags (For socket IPv4/IPv6)
- [x] `--wide` (don't truncate IP addresses)
- [x] `--extend`(display other/more information)
- [x] `--programs` (display PID/Program name for sockets)
- [x] `--timers` (display timers)
- [ ] `--numeric` (don't resolve names)
- [x] `--numeric-host` (don't resolve host names)
- [x] `--numeric-user` (don't resolve user names)
- [ ] `--numeric-port` (don't resolve port names)
- [x] `--listening` (display listening server sockets)

Formating flags (For socket UNIX)
- [x] `--listening` (display listening server sockets)

Formattig flags (For route printing)
- [x] `--numeric-host` (don't resolve host names)
- [ ] `--numeric-port` (don't resolve port names)
- [ ] `--symbolic` (resolve hardware names)

Route printing source flags
- [ ] `--fib` (display Forwarding Information Base (default). Flag not implemented, because it is default behavior)
- [x] `--cache` (display routing cache instead of FIB (Implemented for IPv6 only))

Address families:
- [x] IPv4
- [x] IPv6

Socket types:
- [x] `IPv4/TCP`
- [x] `IPv4/UDP`
- [x] `IPv4/UDPL`
- [x] `IPv4/RAW`
- [x] `IPv6/TCP`
- [x] `IPv6/UDP`
- [x] `IPv6/UDPL`
- [x] `IPv6/RAW`
- [x] `UNIX`
- [ ] `ax25`
- [ ] `sctp`
- [ ] `ipx`
- [ ] `netrom`


Arbitrary flags:
- [ ] `--verbose`, `-v` (Provides errors for unsupported address families of the kernel)
- [x] `--continuous` (continuous listing)


# Additional information
netstat implements net-tools/netstat functionality. In the process of implementing
it, some legacy functionality was uncovered, hence the implementation deviates in
some places.

## Route cache IPv4/IPv6
### IPv4
net-tools/netstats implements the printing of route cache for IPv4 and IPv6.
The linux kernel does not provide route cache for IPv4 since version 3.6.
>.SH NOTES
>Starting with Linux kernel version 3.6, there is no routing cache for IPv4
>anymore. Hence
>.B "ip route show cached"
>will never print any entries on systems with this or newer kernel versions.

Hence this functionality is not implemented.

- [https:github.com/iproute2/iproute2/blob/39e4b6f5f315680ac914187b1d5526c005b8c0d8/man/man8/ip-route.8.in#L1356C10-L1357C75](https:github.com/iproute2/iproute2/blob/39e4b6f5f315680ac914187b1d5526c005b8c0d8/man/man8/ip-route.8.in#L1356C10-L1357C75)

### IPv6
For IPv6, the provision of route cache has changed over time. Since kernel
verison 4.2 only routes with a PMTU exception are marked as RTF_CACHE.
The printing of these routes is implemented.
>Starting from [Linux 4.2 commit 45e4fd26683c](https:git.kernel.org/pub/scm/linux/kernel/git/davem/net-next.git/commit/?id=45e4fd26683c9a5f88600d91b08a484f7f09226a), only a PMTU exception would create a cache entry.
>A router doesn’t have to handle these exceptions, so only hosts would get cache entries.
>
>And they should be pretty rare. Martin KaFai Lau, from Facebook, explains:
>Out of all IPv6 RTF_CACHE routes that are created, the percentage that has a different MTU is very small.
>In one of our end-user facing proxy server, only 1k out of 80k RTF_CACHE routes have a smaller MTU.
>For our DC traffic, there is no MTU exception.

- [https:vincent.bernat.ch/en/blog/2017-ipv6-route-lookup-linux](https:vincent.bernat.ch/en/blog/2017-ipv6-route-lookup-linux)
- [https:git.kernel.org/pub/scm/linux/kernel/git/davem/net-next.git/commit/?id=45e4fd26683c9a5f88600d91b08a484f7f09226a](https:git.kernel.org/pub/scm/linux/kernel/git/davem/net-next.git/commit/?id=45e4fd26683c9a5f88600d91b08a484f7f09226a)


## Some general thoughts
The net-tools tool collection goes way back until beginning of 1998 and has a lot legacy/compatibility code in it.
Parts of old code is functional, but does not follow guidelines for code quality or readability, hence understanding the code
is quite a challange.

