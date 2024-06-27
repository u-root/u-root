# IP Command

## Usage

```bash
ip [ OPTIONS ] OBJECT { COMMAND | help }

ip [ -force ] -batch filename

where  OBJECT := { address | addrlabel | amt | fou | help | ila | ioam | l2tp | link | macsec | maddress | monitor | mptcp | mroute | mrule | neighbor | neighbour | netconf | netns | nexthop | ntable | ntbl | route | rule | sr | tap | tcpmetrics | token | tunnel | tuntap | vrf | xfrm }


OPTIONS := { -V[ersion] | -s[tatistics] | -d[etails] | -r[esolve] | -h[uman-readable] | -iec | -j[son] | -p[retty] | -f[amily] { inet | inet6 | mpls | bridge | link } | -4 | -6 | -M | -B | -0 | -l[oops] { maximum-addr-flush-attempts } | -br[ief] | -o[neline] | -t[imestamp] | -ts[hort] | -b[atch] [filename] | -rc[vbuf] [size] | -n[etns] name | -N[umeric] | -a[ll] | -c[olor] | -echo}
```

## Implementation Status

 Commands and options are checked if atleast one subcommand exist.

Options

- [ ] -V[ersion]
- [ ] -s[tatistics] / -stats
- [ ] -d[etails]
- [ ] -r[esolve]
- [ ] -h[uman-readable]
- [ ] -iec
- [ ] -j[son]
- [ ] -p[retty]
- [ ] -f[amily]
- [ ] -4
- [ ] -6
- [ ] -M
- [ ] -B
- [ ] -0
- [ ] -l[oops]
- [ ] -br[ief]
- [ ] -o[neline]
- [ ] -t[imestamp]
- [ ] -ts[hort]
- [ ] -b[atch] [filename]
- [ ] -rc[vbuf] [size]
- [ ] -n[etns] name
- [ ] -N[umeric]
- [ ] -a[ll]
- [ ] -c[olor]
- [ ] -echo

Command

- [x] address
  - [x] add
  - [ ] change
  - [ ] replace
  - [x] del
  - [ ] show
  - [ ] save
  - [ ] flush
  - [ ] showdump
  - [ ] restore

- [ ] addrlabel
  - [ ] add
  - [ ] del
  - [ ] list
  - [ ] flush
  - [ ] help

- [ ] fou
  - [ ] add
  - [ ] del
  - [ ] show

- [ ] ila?
  - [ ] add
  - [ ] del
  - [ ] list
  - [ ] help

- [ ] ioam
  - [ ] namespace
    - [ ] add
    - [ ] delete
    - [ ] show
    - [ ] set
  - [ ] schema
    - [ ] add
    - [ ] delete
    - [ ] show
  - [ ] help

- [ ] l2tp
  - [ ] add
    - [ ] tunnel
    - [ ] session
  - [ ] del
    - [ ] tunnel
    - [ ] session
  - [ ] show
    - [ ] tunnel
    - [ ] session

- [x] link
  - [x] add
  - [ ] delete
  - [x] set
  - [x] show
  - [ ] xstats
  - [ ] afstats
  - [ ] property
    - [ ] add
    - [ ] del
  - [ ] help

- [ ] macsec
  - [ ] add
  - [ ] set
  - [ ] del
  - [ ] show
  - [ ] offload
  - [ ] help

- [ ] maddr(ess)
  - [ ] add
  - [ ] del
  - [ ] show
  - [ ] help

- [ ] monitor
  - [ ] all
  - [ ] address
  - [ ] link
  - [ ] mroute
  - [ ] neigh
  - [ ] netconf
  - [ ] nexthop
  - [ ] nsid
  - [ ] prefix
  - [ ] route
  - [ ] rule
  - [ ] help
  
- [ ] mptcp
  - [ ] endpoint
    - [ ] add
    - [ ] delete
    - [ ] show
    - [ ] flush
  - [ ] limit
    - [ ] set
    - [ ] show
  - [ ] monitor
  - [ ] help

- [ ] mroute
  - [ ] show
  - [ ] help

- [x] neigh(bour/bor)
  - [ ] add
  - [ ] del
  - [ ] change
  - [ ] replace
  - [x] show
  - [ ] flush
  - [ ] get

- [ ] netconf
  - [ ] show
  - [ ] help

- [ ] netns
  - [ ] list
  - [ ] add
  - [ ] attach
  - [ ] set
  - [ ] delete
  - [ ] identify
  - [ ] pids
  - [ ] exec
  - [ ] monitor
  - [ ] list-id
  - [ ] help

- [ ] nexthop
  - [ ] list
  - [ ] flush
  - [ ] add
  - [ ] replace
  - [ ] get
  - [ ] del
  - [ ] bucket
    - [ ] list
    - [ ] get
  - [ ] help

- [ ] ntable
  - [ ] change
  - [ ] show
  - [ ] help

- [x] route
  - [ ] list
  - [ ] flush
  - [ ] save
  - [ ] restore
  - [ ] showdump
  - [ ] get
  - [x] add
  - [x] del
  - [ ] change
  - [ ] append
  - [ ] replace
  - [ ] help

- [ ] rule / mrule
  - [ ] list
  - [ ] add
  - [ ] del
  - [ ] flush
  - [ ] save
  - [ ] restore
  - [ ] help

- [ ] sr  
  - [ ] hmac
    - [ ] set
    - [ ] show
  - [ ] tunsrc
    - [ ] set
    - [ ] show
  - [ ] help

- [ ] tcp_metrics
  - [ ] show
  - [ ] flush
  - [ ] delete
  - [ ] help

- [ ] token
  - [ ] list
  - [ ] set
  - [ ] del
  - [ ] get
  - [ ] help

- [ ] tunnel
  - [ ] add
  - [ ] change
  - [ ] del
  - [ ] show
  - [ ] prl
  - [ ] 6rd
  - [ ] help

- [ ] tuntap / tap
  - [ ] add
  - [ ] del
  - [ ] show
  - [ ] list
  - [ ] lst
  - [ ] help

- [ ] vrf
  - [ ] show
  - [ ] exec
  - [ ] identify
  - [ ] pids
  - [ ] help

- [ ] xfrm
  - [ ] state
    - [ ] add
    - [ ] update
    - [ ] allocspi
    - [ ] delete
    - [ ] get
    - [ ] deleteall
    - [ ] list
    - [ ] flush
    - [ ] count
  - [ ] policy
    - [ ] add
    - [ ] update
    - [ ] delete
    - [ ] get
    - [ ] deleteall
    - [ ] list
    - [ ] flush
    - [ ] count
    - [ ] set
    - [ ] setdefault
    - [ ] getdefault
  - [ ] monitor
    - [ ] all
    - [ ] aquire
    - [ ] expire
    - [ ] SA
    - [ ] aevent
    - [ ] policy
    - [ ] report
  - [ ] help


### Testing

- [ ] unit tests for conversions
- [ ] integration tests 
