// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// This JSON is from an APU2 running coreboot.
var apu2JSON = `{
	"Memory": {
		"Tag": 1,
		"Size": 108,
		"Maps": [
			{
				"Start": 0,
				"Size": 4096,
				"Mtype": 16
			},
			{
				"Start": 4096,
				"Size": 651264,
				"Mtype": 1
			},
			{
				"Start": 786432,
				"Size": 2012143616,
				"Mtype": 1
			},
			{
				"Start": 2012930048,
				"Size": 335872,
				"Mtype": 16
			}
		]
	},
	"MemConsole": {
		"Tag": 23,
		"Size": 131064,
		"Address": 2013130752,
		"CSize": 0,
		"Cursor": 240,
		"Data": "PCEngines apu2\r\ncoreboot build 20170228\r\n2032 MB DRAM\r\n\r\n\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000"
	},
	"Consoles": [
		""
	],
	"TimeStampsTable": {
		"Tag": 0,
		"Size": 0,
		"Addr": 0
	},
	"TimeStamps": null,
	"UART": [
		{
			"Tag": 15,
			"Size": 20,
			"Type": 1,
			"BaseAddr": 1016,
			"Baud": 115200,
			"RegWidth": 16
		}
	],
	"MainBoard": {
		"Tag": 3,
		"Size": 40,
		"Vendor": "PC Engines",
		"PartNumber": "PCEngines apu2"
	},
	"Hwrpb": {
		"Tag": 0,
		"Size": 0,
		"HwrPB": 0
	},
	"CBMemory": null,
	"BoardID": {
		"Tag": 37,
		"Size": 16,
		"BoardID": 2012962816
	},
	"StringVars": {
		"LB_TAG_BUILD": "Tue Feb 28 22:34:13 UTC 2017",
		"LB_TAG_COMPILE_BY": "root",
		"LB_TAG_COMPILE_DOMAIN": "",
		"LB_TAG_COMPILE_HOST": "3aa919ff57dc",
		"LB_TAG_COMPILE_TIME": "22:34:13",
		"LB_TAG_EXTRA_VERSION": "-4.0.7",
		"LB_TAG_VERSION": "8b10004"
	},
	"BootMediaParams": {
		"Tag": 0,
		"Size": 0,
		"FMAPOffset": 0,
		"CBFSOffset": 0,
		"CBFSSize": 0,
		"BootMediaSize": 0
	},
	"VersionTimeStamp": 38,
	"Unknown": null,
	"Ignored": null
}`
