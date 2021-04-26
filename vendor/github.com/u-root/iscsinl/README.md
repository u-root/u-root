# iscsinl

[![CircleCI](https://circleci.com/gh/u-root/iscsinl.svg?style=svg)](https://circleci.com/gh/u-root/iscsinl)
[![Go Report Card](https://goreportcard.com/badge/github.com/u-root/iscsinl)](https://goreportcard.com/report/github.com/u-root/iscsinl)
[![GoDoc](https://godoc.org/github.com/u-root/iscsinl?status.svg)](https://godoc.org/github.com/u-root/iscsinl)

Go iSCSI netlink library

## TODO

Currently, after establishing a successful iscsi session with target, iscsinl scans all LUNS i.e
uses wild card `- - -` while writing to `/sys/class/scsi_host/host%d/scan.`
The three `- - -` stand for channel, SCSI target ID, and LUN, where - means all.

In future we would like the iscsnl initiator code to accept LUN as an input argument
just like initiatorName, so that user can customize which LUN s(he) wants to be scanned.
