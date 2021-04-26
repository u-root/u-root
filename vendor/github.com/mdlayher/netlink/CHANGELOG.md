# CHANGELOG

## Unreleased

- n/a

## v1.1.1

**This is the last release of package netlink that supports Go 1.11.**

- [Improvement] [#165](https://github.com/mdlayher/netlink/pull/165):
  `netlink.Conn` `SetReadBuffer` and `SetWriteBuffer` methods now attempt the
  `SO_*BUFFORCE` socket options to possibly ignore system limits given elevated
  caller permissions. Thanks @MarkusBauer.
- [Note]
  [commit](https://github.com/mdlayher/netlink/commit/c5f8ab79aa345dcfcf7f14d746659ca1b80a0ecc):
  `netlink.Conn.Close` has had a long-standing bug
  [#162](https://github.com/mdlayher/netlink/pull/162) related to internal
  concurrency handling where a call to `Close` is not sufficient to unblock
  pending reads. To effectively fix this issue, it is necessary to drop support
  for Go 1.11 and below. This will be fixed in a future release, but a
  workaround is noted in the method documentation as of now.

## v1.1.0

- [New API] [#157](https://github.com/mdlayher/netlink/pull/157): the
  `netlink.AttributeDecoder.TypeFlags` method enables retrieval of the type bits
  stored in a netlink attribute's type field, because the existing `Type` method
  masks away these bits. Thanks @ti-mo!
- [Performance] [#157](https://github.com/mdlayher/netlink/pull/157): `netlink.AttributeDecoder`
  now decodes netlink attributes on demand, enabling callers who only need a
  limited number of attributes to exit early from decoding loops. Thanks @ti-mo!
- [Improvement] [#161](https://github.com/mdlayher/netlink/pull/161): `netlink.Conn`
  system calls are now ready for Go 1.14+'s changes to goroutine preemption.
  See the PR for details.

## v1.0.0

- Initial stable commit.
