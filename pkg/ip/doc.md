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

- [x] -s[tatistics] / -stats
- [x] -d[etails]
- [ ] -r[esolve]
- [x] -h[uman-readable]
- [x] -iec
- [x] -j[son]
- [x] -p[retty]
- [x] -f[amily]
- [x] -4
- [x] -6
- [ ] -M
- [ ] -B
- [x] -0
- [x] -l[oops]
- [x] -br[ief]
- [x] -o[neline]
- [x] -t[imestamp]
- [x] -ts[hort]
- [x] -b[atch] [filename]
- [x] -rc[vbuf] [size]
- [x] -n[etns] name
- [x] -N[umeric]
- [x] -a[ll]
- [ ] -c[olor]

Command

- [x] address
  - [x] add
  - [ ] change
  - [x] replace
  - [x] del
  - [x] show
  - [x] flush

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
  - [x] delete
  - [x] set
  - [x] show
  - [ ] xstats
  - [ ] afstats
  - [ ] property
    - [ ] add
    - [ ] del
  - [x] help

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

- [x] monitor
  - [x] all
  - [x] address
  - [x] link
  - [ ] mroute
  - [x] neigh
  - [ ] netconf
  - [ ] nexthop
  - [ ] nsid
  - [ ] prefix
  - [x] route
  - [ ] rule
  - [x] help
  
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
  - [x] add
  - [x] del
  - [ ] change
  - [x] replace
  - [x] show
  - [x] flush
  - [x] get

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
  - [x] list
  - [x] flush
  - [x] get
  - [x] add
  - [x] del
  - [ ] change
  - [x] append
  - [x] replace
  - [x] help

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

- [x] tcp_metrics
  - [x] show
  - [ ] flush
  - [ ] delete
  - [x] help

- [ ] token
  - [ ] list
  - [ ] set
  - [ ] del
  - [ ] get
  - [ ] help

- [x] tunnel
  - [x] add
  - [ ] change
  - [x] del
  - [ ] show
  - [ ] prl
  - [ ] 6rd
  - [x] help

- [x] tuntap / tap
  - [x] add
  - [x] del
  - [x] show
  - [x] list
  - [x] lst
  - [x] help

- [x] vrf
  - [x] show
  - [ ] exec
  - [ ] identify
  - [ ] pids
  - [x] help

- [x] xfrm
  - [x] state
    - [x] add
    - [x] update
    - [x] allocspi
    - [x] delete
    - [x] get
    - [x] deleteall
    - [x] list
    - [x] flush
    - [x] count
  - [x] policy
    - [x] add
    - [x] update
    - [x] delete
    - [x] get
    - [x] deleteall
    - [x] list
    - [x] flush
    - [x] count
    - [x] set
    - [ ] setdefault
    - [ ] getdefault
  - [x] monitor
    - [x] all
    - [ ] aquire
    - [x] expire
    - [ ] SA
    - [ ] aevent
    - [ ] policy
    - [ ] report
  - [x] help


### Testing

- [ ] unit tests for conversions
- [ ] integration tests 
