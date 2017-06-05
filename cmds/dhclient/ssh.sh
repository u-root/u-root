#!/bin/bash
set -e
go build .
scp ./dhclient xchenan@nuc:~/dhclient/
