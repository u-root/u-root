# Status of implementation

Objects:

## qdisc
- [x] Show/List (library does not provide string functions for structures, hence limited output)
- [x] Add
- [x] Delete
- [x] Replace
- [x] Change
- [ ] Link

### Classless Qdiscs
- [ ] choke
- [x] codel
- [ ] [p|b]FIFO
- [ ] fq
- [ ] fq_codel
- [ ] fq_pie
- [ ] gred
- [ ] hhf
- [x] ingress
- [ ] mqprio
- [ ] multiq
- [ ] netem
- [ ] pfifo_fast
- [ ] pie
- [ ] red
- [ ] sfb
- [ ] sfq
- [ ] tbf
- [x] clsact

## class (untested yet)
- [x] Show/List (library does not provide string functions for structures, hence limited output)
- [x] Add
- [x] Delete
- [ ] Replace
- [ ] Change

### Classfull Qdiscs
- [ ] ATM (not supported for adding byt go-tc library)
- [ ] DRR (not supported for adding byt go-tc library)
- [ ] ETS (not supported for adding byt go-tc library)
- [ ] HFSC
- - [x] qdisc args parsing
- - [x] class args parsing
- [x] HTB
- - [x] qdisc args parsing
- - [x] class args parsing
- [ ] PRIO (not supported for adding by go-tc library)
- [ ] QFQ

## filter (untested yet)
- [x] Show/List (library does not provide string functions for structures, hence limited output)
- [x] Add
- [x] Delete
- [ ] Replace
- [ ] Change

### Filter types
- [x] basic
- - [x] action
- - [ ] match
- - [x] classid/flowid
- [ ] bpf
- - [ ] object-file
- - [ ] section
- - [ ] export
- - [ ] verbose
- - [ ] direct-action (da)
- - [ ] skip-hw | skip-sw
- - [ ] police
- - [ ] action
- - [ ] classid/flowid
- - [ ] bytecode
- - [ ] bytecode-file
- [ ] cgroup
- - [ ] action
- - [ ] match
- [ ] flow,flower
- - [ ] action
- - [ ] baseclass
- - [ ] divisor
- - [ ] hash keys
- - [ ] map key
- - [ ] match
- [ ] fw
- - [ ] classid
- - [ ] action
- [ ] route
- - [ ] action
- - [ ] classid
- - [ ] from
- - [ ] fromif
- - [ ] to
- [ ] u32
- - [ ] handle
- - [ ] offset
- - [ ] hashkey
- - [ ] classid
- - [ ] divisor
- - [ ] order
- - [ ] sample
- - [ ] link
- - [ ] indev
- - [ ] skip-sw
- - [ ] skip-hw
- [ ] matchall
- - [ ] action
- - [ ] classid
- - [ ] skip-sw
- - [ ] skip-hw

## chain
- [ ] add
- [ ] del
- [ ] get
- [ ] show

## monitor
