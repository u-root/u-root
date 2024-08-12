// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package health_test

// This is from a very basic Google Cloud HPC cluster.
var data = `
[
	{
		"Hostname": "hpcslurm-debugnodeset-0",
		"Info": {
			"memory": {
				"total_physical_bytes": 8589934592,
				"total_usable_bytes": 8056090624,
				"supported_page_sizes": [
					1073741824,
					2097152
				],
				"modules": null
			},
			"block": {
				"total_size_bytes": 53687091200,
				"disks": [
					{
						"name": "sda",
						"size_bytes": 53687091200,
						"physical_block_size_bytes": 4096,
						"drive_type": "hdd",
						"removable": false,
						"storage_controller": "scsi",
						"bus_path": "pci-0000:00:03.0-scsi-0:0:1:0",
						"vendor": "Google",
						"model": "PersistentDisk",
						"serial_number": "persistent-disk-0",
						"wwn": "unknown",
						"partitions": [
							{
								"name": "sda1",
								"label": "EFI\\x20System\\x20Partition",
								"mount_point": "/boot/efi",
								"size_bytes": 209715200,
								"type": "vfat",
								"read_only": false,
								"uuid": "a407d4b7-cfe4-4f7e-b9fc-ee7799ba3b84",
								"filesystem_label": "unknown"
							},
							{
								"name": "sda2",
								"label": "unknown",
								"mount_point": "/",
								"size_bytes": 53475328000,
								"type": "xfs",
								"read_only": false,
								"uuid": "144c8c6f-9c84-47c9-b637-8b7723fdb3ef",
								"filesystem_label": "root"
							}
						]
					}
				]
			},
			"cpu": {
				"total_cores": 1,
				"total_threads": 1,
				"processors": [
					{
						"id": 0,
						"total_cores": 1,
						"total_threads": 1,
						"vendor": "GenuineIntel",
						"model": "Intel(R) Xeon(R) CPU @ 2.80GHz",
						"capabilities": [
							"fpu",
							"vme",
							"de",
							"pse",
							"tsc",
							"msr",
							"pae",
							"mce",
							"cx8",
							"apic",
							"sep",
							"mtrr",
							"pge",
							"mca",
							"cmov",
							"pat",
							"pse36",
							"clflush",
							"mmx",
							"fxsr",
							"sse",
							"sse2",
							"ss",
							"ht",
							"syscall",
							"nx",
							"pdpe1gb",
							"rdtscp",
							"lm",
							"constant_tsc",
							"rep_good",
							"nopl",
							"xtopology",
							"nonstop_tsc",
							"cpuid",
							"tsc_known_freq",
							"pni",
							"pclmulqdq",
							"ssse3",
							"fma",
							"cx16",
							"pcid",
							"sse4_1",
							"sse4_2",
							"x2apic",
							"movbe",
							"popcnt",
							"aes",
							"xsave",
							"avx",
							"f16c",
							"rdrand",
							"hypervisor",
							"lahf_lm",
							"abm",
							"3dnowprefetch",
							"invpcid_single",
							"ssbd",
							"ibrs",
							"ibpb",
							"stibp",
							"ibrs_enhanced",
							"fsgsbase",
							"tsc_adjust",
							"bmi1",
							"hle",
							"avx2",
							"smep",
							"bmi2",
							"erms",
							"invpcid",
							"rtm",
							"avx512f",
							"avx512dq",
							"rdseed",
							"adx",
							"smap",
							"clflushopt",
							"clwb",
							"avx512cd",
							"avx512bw",
							"avx512vl",
							"xsaveopt",
							"xsavec",
							"xgetbv1",
							"xsaves",
							"arat",
							"avx512_vnni",
							"md_clear",
							"arch_capabilities"
						],
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						]
					}
				]
			},
			"topology": {
				"architecture": "smp",
				"nodes": [
					{
						"id": 0,
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						],
						"caches": [
							{
								"level": 1,
								"type": "instruction",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 1,
								"type": "data",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 2,
								"type": "unified",
								"size_bytes": 1048576,
								"logical_processors": [
									0
								]
							},
							{
								"level": 3,
								"type": "unified",
								"size_bytes": 34603008,
								"logical_processors": [
									0
								]
							}
						],
						"distances": [
							10
						],
						"memory": {
							"total_physical_bytes": 8589934592,
							"total_usable_bytes": 8056090624,
							"supported_page_sizes": [
								1073741824,
								2097152
							],
							"modules": null
						}
					}
				]
			},
			"network": {
				"nics": [
					{
						"name": "eth0",
						"mac_address": "42:01:0a:00:00:d0",
						"is_virtual": false,
						"capabilities": [
							{
								"name": "auto-negotiation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "pause-frame-use",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-checksumming",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-checksumming",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv4",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-ip-generic",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv6",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-fcoe-crc",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-sctp",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather-fraglist",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tcp-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-ecn-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tcp-mangleid-segmentation",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tx-tcp6-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-receive-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "large-receive-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "ntuple-filters",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "receive-hashing",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "highdma",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "rx-vlan-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "vlan-challenged",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-lockless",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "netns-local",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-robust",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-fcoe-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip4-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip6-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-partial",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tunnel-remcsum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-sctp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-esp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-list",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp-gro-forwarding",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "rx-gro-list",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tls-hw-rx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "fcoe-mtu",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-nocache-copy",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "loopback",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-fcs",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-all",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-stag-hw-insert",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-hw-parse",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "l2-fwd-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "hw-tc-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-tx-csum-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp_tunnel-port-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-tx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-gro-hw",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-record",
								"is_enabled": false,
								"can_enable": false
							}
						],
						"speed": "Unknown!",
						"duplex": "Unknown!(255)"
					}
				]
			},
			"gpu": {
				"cards": null
			},
			"chassis": {
				"asset_tag": "",
				"serial_number": "unknown",
				"type": "1",
				"type_description": "Other",
				"vendor": "Google",
				"version": ""
			},
			"bios": {
				"vendor": "Google",
				"version": "Google",
				"date": "06/07/2024"
			},
			"baseboard": {
				"asset_tag": "79271F5B-5EDA-F01E-EA13-FEE7DF20AF71",
				"serial_number": "unknown",
				"vendor": "Google",
				"version": "",
				"product": "Google Compute Engine"
			},
			"product": {
				"family": "",
				"name": "Google Compute Engine",
				"vendor": "Google",
				"serial_number": "unknown",
				"uuid": "unknown",
				"sku": "",
				"version": ""
			},
			"pci": {
				"Devices": [
					{
						"driver": "",
						"address": "0000:00:00.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "1237",
							"name": "440FX - 82441FX PMC [Natoma]"
						},
						"revision": "0x02",
						"subsystem": {
							"id": "1100",
							"name": "Qemu virtual machine"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "00",
							"name": "Host bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7110",
							"name": "82371AB/EB/MB PIIX4 ISA"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "01",
							"name": "ISA bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.3",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7113",
							"name": "82371AB/EB/MB PIIX4 ACPI"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "80",
							"name": "Bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:03.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1004",
							"name": "Virtio SCSI"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0008",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "00",
							"name": "Non-VGA unclassified device"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:04.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1000",
							"name": "Virtio network device"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0001",
							"name": "unknown"
						},
						"class": {
							"id": "02",
							"name": "Network controller"
						},
						"subclass": {
							"id": "00",
							"name": "Ethernet controller"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:05.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1005",
							"name": "Virtio RNG"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0004",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "ff",
							"name": "unknown"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					}
				]
			}
		},
		"Kernel": {
			"Version": "Linux version 4.18.0-513.24.1.el8_9.x86_64 (mockbuild@iad1-prod-build001.bld.equ.rockylinux.org) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-20) (GCC)) #1 SMP Thu Apr 4 18:13:02 UTC 2024\n",
			"Modules": "tcp_diag 16384 0 - Live 0x0000000000000000\ninet_diag 24576 1 tcp_diag, Live 0x0000000000000000\nbinfmt_misc 24576 1 - Live 0x0000000000000000\nnfsv3 57344 1 - Live 0x0000000000000000\nrpcsec_gss_krb5 45056 0 - Live 0x0000000000000000\nnfsv4 917504 2 - Live 0x0000000000000000\ndns_resolver 16384 1 nfsv4, Live 0x0000000000000000\nnfs 425984 4 nfsv3,nfsv4, Live 0x0000000000000000\nfscache 389120 1 nfs, Live 0x0000000000000000\nintel_rapl_msr 16384 0 - Live 0x0000000000000000\nintel_rapl_common 24576 1 intel_rapl_msr, Live 0x0000000000000000\nintel_uncore_frequency_common 16384 0 - Live 0x0000000000000000\nisst_if_common 16384 0 - Live 0x0000000000000000\nnfit 65536 0 - Live 0x0000000000000000\nlibnvdimm 200704 1 nfit, Live 0x0000000000000000\nrapl 20480 0 - Live 0x0000000000000000\ni2c_piix4 24576 0 - Live 0x0000000000000000\nvfat 20480 1 - Live 0x0000000000000000\nfat 86016 1 vfat, Live 0x0000000000000000\npcspkr 16384 0 - Live 0x0000000000000000\nnfsd 548864 13 - Live 0x0000000000000000\nauth_rpcgss 139264 2 rpcsec_gss_krb5,nfsd, Live 0x0000000000000000\nnfs_acl 16384 2 nfsv3,nfsd, Live 0x0000000000000000\nlockd 126976 3 nfsv3,nfs,nfsd, Live 0x0000000000000000\ngrace 16384 2 nfsd,lockd, Live 0x0000000000000000\nxfs 1593344 1 - Live 0x0000000000000000\nlibcrc32c 16384 1 xfs, Live 0x0000000000000000\nsd_mod 57344 2 - Live 0x0000000000000000\nsg 40960 0 - Live 0x0000000000000000\nnvme_tcp 36864 0 - Live 0x0000000000000000 (X)\nnvme_fabrics 24576 1 nvme_tcp, Live 0x0000000000000000\ncrct10dif_pclmul 16384 1 - Live 0x0000000000000000\ncrc32_pclmul 16384 0 - Live 0x0000000000000000\ncrc32c_intel 24576 1 - Live 0x0000000000000000\nvirtio_net 61440 0 - Live 0x0000000000000000\nghash_clmulni_intel 16384 0 - Live 0x0000000000000000\nserio_raw 16384 0 - Live 0x0000000000000000\nnet_failover 24576 1 virtio_net, Live 0x0000000000000000\nfailover 16384 1 net_failover, Live 0x0000000000000000\nvirtio_scsi 20480 2 - Live 0x0000000000000000\nnvme 45056 0 - Live 0x0000000000000000\nnvme_core 139264 3 nvme_tcp,nvme_fabrics,nvme, Live 0x0000000000000000\nt10_pi 16384 2 sd_mod,nvme_core, Live 0x0000000000000000\nsunrpc 585728 32 nfsv3,rpcsec_gss_krb5,nfsv4,nfs,nfsd,auth_rpcgss,nfs_acl,lockd, Live 0x0000000000000000\ndm_mirror 28672 0 - Live 0x0000000000000000\ndm_region_hash 20480 1 dm_mirror, Live 0x0000000000000000\ndm_log 20480 2 dm_mirror,dm_region_hash, Live 0x0000000000000000\ndm_mod 155648 2 dm_mirror,dm_log, Live 0x0000000000000000\n",
			"Drivers": "Character devices:\n  1 mem\n  4 /dev/vc/0\n  4 tty\n  4 ttyS\n  5 /dev/tty\n  5 /dev/console\n  5 /dev/ptmx\n  7 vcs\n 10 misc\n 13 input\n 21 sg\n 29 fb\n128 ptm\n136 pts\n162 raw\n180 usb\n188 ttyUSB\n189 usb_device\n202 cpu/msr\n203 cpu/cpuid\n240 dimmctl\n241 ndctl\n242 nvme-generic\n243 nvme\n244 hidraw\n245 ttyDBC\n246 usbmon\n247 bsg\n248 watchdog\n249 ptp\n250 pps\n251 rtc\n252 dax\n253 tpm\n254 gpiochip\n\nBlock devices:\n  8 sd\n  9 md\n 65 sd\n 66 sd\n 67 sd\n 68 sd\n 69 sd\n 70 sd\n 71 sd\n128 sd\n129 sd\n130 sd\n131 sd\n132 sd\n133 sd\n134 sd\n135 sd\n253 device-mapper\n254 mdp\n259 blkext\n"
		},
		"Stderr": "WARNING: \n/sys/class/drm does not exist on this system (likely the host system is a\nvirtual machine or container with no graphics). Therefore,\nGPUInfo.GraphicsCards will be an empty array.\nWARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied\nWARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied\nWARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied\nWARNING: Unable to read product_uuid: open /sys/class/dmi/id/product_uuid: permission denied\n",
		"Err": ""
	},
	{
		"Hostname": "hpcslurm-debugnodeset-16",
		"Info": {
			"memory": {
				"total_physical_bytes": 8589934592,
				"total_usable_bytes": 8056098816,
				"supported_page_sizes": [
					1073741824,
					2097152
				],
				"modules": null
			},
			"block": {
				"total_size_bytes": 53687091200,
				"disks": [
					{
						"name": "sda",
						"size_bytes": 53687091200,
						"physical_block_size_bytes": 4096,
						"drive_type": "hdd",
						"removable": false,
						"storage_controller": "scsi",
						"bus_path": "pci-0000:00:03.0-scsi-0:0:1:0",
						"vendor": "Google",
						"model": "PersistentDisk",
						"serial_number": "persistent-disk-0",
						"wwn": "unknown",
						"partitions": [
							{
								"name": "sda1",
								"label": "EFI\\x20System\\x20Partition",
								"mount_point": "/boot/efi",
								"size_bytes": 209715200,
								"type": "vfat",
								"read_only": false,
								"uuid": "a407d4b7-cfe4-4f7e-b9fc-ee7799ba3b84",
								"filesystem_label": "unknown"
							},
							{
								"name": "sda2",
								"label": "unknown",
								"mount_point": "/",
								"size_bytes": 53475328000,
								"type": "xfs",
								"read_only": false,
								"uuid": "144c8c6f-9c84-47c9-b637-8b7723fdb3ef",
								"filesystem_label": "root"
							}
						]
					}
				]
			},
			"cpu": {
				"total_cores": 1,
				"total_threads": 1,
				"processors": [
					{
						"id": 0,
						"total_cores": 1,
						"total_threads": 1,
						"vendor": "GenuineIntel",
						"model": "Intel(R) Xeon(R) CPU @ 2.80GHz",
						"capabilities": [
							"fpu",
							"vme",
							"de",
							"pse",
							"tsc",
							"msr",
							"pae",
							"mce",
							"cx8",
							"apic",
							"sep",
							"mtrr",
							"pge",
							"mca",
							"cmov",
							"pat",
							"pse36",
							"clflush",
							"mmx",
							"fxsr",
							"sse",
							"sse2",
							"ss",
							"ht",
							"syscall",
							"nx",
							"pdpe1gb",
							"rdtscp",
							"lm",
							"constant_tsc",
							"rep_good",
							"nopl",
							"xtopology",
							"nonstop_tsc",
							"cpuid",
							"tsc_known_freq",
							"pni",
							"pclmulqdq",
							"ssse3",
							"fma",
							"cx16",
							"pcid",
							"sse4_1",
							"sse4_2",
							"x2apic",
							"movbe",
							"popcnt",
							"aes",
							"xsave",
							"avx",
							"f16c",
							"rdrand",
							"hypervisor",
							"lahf_lm",
							"abm",
							"3dnowprefetch",
							"invpcid_single",
							"ssbd",
							"ibrs",
							"ibpb",
							"stibp",
							"ibrs_enhanced",
							"fsgsbase",
							"tsc_adjust",
							"bmi1",
							"hle",
							"avx2",
							"smep",
							"bmi2",
							"erms",
							"invpcid",
							"rtm",
							"avx512f",
							"avx512dq",
							"rdseed",
							"adx",
							"smap",
							"clflushopt",
							"clwb",
							"avx512cd",
							"avx512bw",
							"avx512vl",
							"xsaveopt",
							"xsavec",
							"xgetbv1",
							"xsaves",
							"arat",
							"avx512_vnni",
							"md_clear",
							"arch_capabilities"
						],
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						]
					}
				]
			},
			"topology": {
				"architecture": "smp",
				"nodes": [
					{
						"id": 0,
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						],
						"caches": [
							{
								"level": 1,
								"type": "instruction",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 1,
								"type": "data",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 2,
								"type": "unified",
								"size_bytes": 1048576,
								"logical_processors": [
									0
								]
							},
							{
								"level": 3,
								"type": "unified",
								"size_bytes": 34603008,
								"logical_processors": [
									0
								]
							}
						],
						"distances": [
							10
						],
						"memory": {
							"total_physical_bytes": 8589934592,
							"total_usable_bytes": 8056098816,
							"supported_page_sizes": [
								1073741824,
								2097152
							],
							"modules": null
						}
					}
				]
			},
			"network": {
				"nics": [
					{
						"name": "eth0",
						"mac_address": "42:01:0a:00:00:d7",
						"is_virtual": false,
						"capabilities": [
							{
								"name": "auto-negotiation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "pause-frame-use",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-checksumming",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-checksumming",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv4",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-ip-generic",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv6",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-fcoe-crc",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-sctp",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather-fraglist",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tcp-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-ecn-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tcp-mangleid-segmentation",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tx-tcp6-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-receive-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "large-receive-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "ntuple-filters",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "receive-hashing",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "highdma",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "rx-vlan-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "vlan-challenged",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-lockless",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "netns-local",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-robust",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-fcoe-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip4-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip6-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-partial",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tunnel-remcsum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-sctp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-esp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-list",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp-gro-forwarding",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "rx-gro-list",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tls-hw-rx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "fcoe-mtu",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-nocache-copy",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "loopback",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-fcs",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-all",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-stag-hw-insert",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-hw-parse",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "l2-fwd-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "hw-tc-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-tx-csum-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp_tunnel-port-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-tx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-gro-hw",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-record",
								"is_enabled": false,
								"can_enable": false
							}
						],
						"speed": "Unknown!",
						"duplex": "Unknown!(255)"
					}
				]
			},
			"gpu": {
				"cards": null
			},
			"chassis": {
				"asset_tag": "",
				"serial_number": "unknown",
				"type": "1",
				"type_description": "Other",
				"vendor": "Google",
				"version": ""
			},
			"bios": {
				"vendor": "Google",
				"version": "Google",
				"date": "06/07/2024"
			},
			"baseboard": {
				"asset_tag": "5F90E76A-DC6D-14B4-B245-54526F84DFE5",
				"serial_number": "unknown",
				"vendor": "Google",
				"version": "",
				"product": "Google Compute Engine"
			},
			"product": {
				"family": "",
				"name": "Google Compute Engine",
				"vendor": "Google",
				"serial_number": "unknown",
				"uuid": "unknown",
				"sku": "",
				"version": ""
			},
			"pci": {
				"Devices": [
					{
						"driver": "",
						"address": "0000:00:00.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "1237",
							"name": "440FX - 82441FX PMC [Natoma]"
						},
						"revision": "0x02",
						"subsystem": {
							"id": "1100",
							"name": "Qemu virtual machine"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "00",
							"name": "Host bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7110",
							"name": "82371AB/EB/MB PIIX4 ISA"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "01",
							"name": "ISA bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.3",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7113",
							"name": "82371AB/EB/MB PIIX4 ACPI"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "80",
							"name": "Bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:03.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1004",
							"name": "Virtio SCSI"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0008",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "00",
							"name": "Non-VGA unclassified device"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:04.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1000",
							"name": "Virtio network device"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0001",
							"name": "unknown"
						},
						"class": {
							"id": "02",
							"name": "Network controller"
						},
						"subclass": {
							"id": "00",
							"name": "Ethernet controller"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:05.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1005",
							"name": "Virtio RNG"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0004",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "ff",
							"name": "unknown"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					}
				]
			}
		},
		"Kernel": {
			"Version": "Linux version 4.18.0-513.24.1.el8_9.x86_64 (mockbuild@iad1-prod-build001.bld.equ.rockylinux.org) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-20) (GCC)) #1 SMP Thu Apr 4 18:13:02 UTC 2024\n",
			"Modules": "tcp_diag 16384 0 - Live 0x0000000000000000\ninet_diag 24576 1 tcp_diag, Live 0x0000000000000000\nbinfmt_misc 24576 1 - Live 0x0000000000000000\nnfsv3 57344 1 - Live 0x0000000000000000\nrpcsec_gss_krb5 45056 0 - Live 0x0000000000000000\nnfsv4 917504 2 - Live 0x0000000000000000\ndns_resolver 16384 1 nfsv4, Live 0x0000000000000000\nnfs 425984 4 nfsv3,nfsv4, Live 0x0000000000000000\nfscache 389120 1 nfs, Live 0x0000000000000000\nintel_rapl_msr 16384 0 - Live 0x0000000000000000\nintel_rapl_common 24576 1 intel_rapl_msr, Live 0x0000000000000000\nintel_uncore_frequency_common 16384 0 - Live 0x0000000000000000\nisst_if_common 16384 0 - Live 0x0000000000000000\nnfit 65536 0 - Live 0x0000000000000000\nlibnvdimm 200704 1 nfit, Live 0x0000000000000000\nvfat 20480 1 - Live 0x0000000000000000\nfat 86016 1 vfat, Live 0x0000000000000000\ncrct10dif_pclmul 16384 1 - Live 0x0000000000000000\ncrc32_pclmul 16384 0 - Live 0x0000000000000000\nghash_clmulni_intel 16384 0 - Live 0x0000000000000000\nrapl 20480 0 - Live 0x0000000000000000\ni2c_piix4 24576 0 - Live 0x0000000000000000\npcspkr 16384 0 - Live 0x0000000000000000\nnfsd 548864 13 - Live 0x0000000000000000\nauth_rpcgss 139264 2 rpcsec_gss_krb5,nfsd, Live 0x0000000000000000\nnfs_acl 16384 2 nfsv3,nfsd, Live 0x0000000000000000\nlockd 126976 3 nfsv3,nfs,nfsd, Live 0x0000000000000000\ngrace 16384 2 nfsd,lockd, Live 0x0000000000000000\nsunrpc 585728 32 nfsv3,rpcsec_gss_krb5,nfsv4,nfs,nfsd,auth_rpcgss,nfs_acl,lockd, Live 0x0000000000000000\nxfs 1593344 1 - Live 0x0000000000000000\nlibcrc32c 16384 1 xfs, Live 0x0000000000000000\nsd_mod 57344 2 - Live 0x0000000000000000\nsg 40960 0 - Live 0x0000000000000000\nvirtio_net 61440 0 - Live 0x0000000000000000\ncrc32c_intel 24576 1 - Live 0x0000000000000000\nserio_raw 16384 0 - Live 0x0000000000000000\nnet_failover 24576 1 virtio_net, Live 0x0000000000000000\nfailover 16384 1 net_failover, Live 0x0000000000000000\nvirtio_scsi 20480 2 - Live 0x0000000000000000\nnvme 45056 0 - Live 0x0000000000000000\nnvme_core 139264 1 nvme, Live 0x0000000000000000\nt10_pi 16384 2 sd_mod,nvme_core, Live 0x0000000000000000\n",
			"Drivers": "Character devices:\n  1 mem\n  4 /dev/vc/0\n  4 tty\n  4 ttyS\n  5 /dev/tty\n  5 /dev/console\n  5 /dev/ptmx\n  7 vcs\n 10 misc\n 13 input\n 21 sg\n 29 fb\n128 ptm\n136 pts\n162 raw\n180 usb\n188 ttyUSB\n189 usb_device\n202 cpu/msr\n203 cpu/cpuid\n240 dimmctl\n241 ndctl\n242 nvme-generic\n243 nvme\n244 hidraw\n245 ttyDBC\n246 usbmon\n247 bsg\n248 watchdog\n249 ptp\n250 pps\n251 rtc\n252 dax\n253 tpm\n254 gpiochip\n\nBlock devices:\n  8 sd\n  9 md\n 65 sd\n 66 sd\n 67 sd\n 68 sd\n 69 sd\n 70 sd\n 71 sd\n128 sd\n129 sd\n130 sd\n131 sd\n132 sd\n133 sd\n134 sd\n135 sd\n254 mdp\n259 blkext\n"
		},
		"Stderr": "WARNING: \n/sys/class/drm does not exist on this system (likely the host system is a\nvirtual machine or container with no graphics). Therefore,\nGPUInfo.GraphicsCards will be an empty array.\nWARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied\nWARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied\nWARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied\nWARNING: Unable to read product_uuid: open /sys/class/dmi/id/product_uuid: permission denied\n",
		"Err": ""
	},
	{
		"Hostname": "hpcslurm-debugnodeset-1",
		"Info": {
			"memory": {
				"total_physical_bytes": 8589934592,
				"total_usable_bytes": 8056090624,
				"supported_page_sizes": [
					1073741824,
					2097152
				],
				"modules": null
			},
			"block": {
				"total_size_bytes": 53687091200,
				"disks": [
					{
						"name": "sda",
						"size_bytes": 53687091200,
						"physical_block_size_bytes": 4096,
						"drive_type": "hdd",
						"removable": false,
						"storage_controller": "scsi",
						"bus_path": "pci-0000:00:03.0-scsi-0:0:1:0",
						"vendor": "Google",
						"model": "PersistentDisk",
						"serial_number": "persistent-disk-0",
						"wwn": "unknown",
						"partitions": [
							{
								"name": "sda1",
								"label": "EFI\\x20System\\x20Partition",
								"mount_point": "/boot/efi",
								"size_bytes": 209715200,
								"type": "vfat",
								"read_only": false,
								"uuid": "a407d4b7-cfe4-4f7e-b9fc-ee7799ba3b84",
								"filesystem_label": "unknown"
							},
							{
								"name": "sda2",
								"label": "unknown",
								"mount_point": "/",
								"size_bytes": 53475328000,
								"type": "xfs",
								"read_only": false,
								"uuid": "144c8c6f-9c84-47c9-b637-8b7723fdb3ef",
								"filesystem_label": "root"
							}
						]
					}
				]
			},
			"cpu": {
				"total_cores": 1,
				"total_threads": 1,
				"processors": [
					{
						"id": 0,
						"total_cores": 1,
						"total_threads": 1,
						"vendor": "GenuineIntel",
						"model": "Intel(R) Xeon(R) CPU @ 2.80GHz",
						"capabilities": [
							"fpu",
							"vme",
							"de",
							"pse",
							"tsc",
							"msr",
							"pae",
							"mce",
							"cx8",
							"apic",
							"sep",
							"mtrr",
							"pge",
							"mca",
							"cmov",
							"pat",
							"pse36",
							"clflush",
							"mmx",
							"fxsr",
							"sse",
							"sse2",
							"ss",
							"ht",
							"syscall",
							"nx",
							"pdpe1gb",
							"rdtscp",
							"lm",
							"constant_tsc",
							"rep_good",
							"nopl",
							"xtopology",
							"nonstop_tsc",
							"cpuid",
							"tsc_known_freq",
							"pni",
							"pclmulqdq",
							"ssse3",
							"fma",
							"cx16",
							"pcid",
							"sse4_1",
							"sse4_2",
							"x2apic",
							"movbe",
							"popcnt",
							"aes",
							"xsave",
							"avx",
							"f16c",
							"rdrand",
							"hypervisor",
							"lahf_lm",
							"abm",
							"3dnowprefetch",
							"invpcid_single",
							"ssbd",
							"ibrs",
							"ibpb",
							"stibp",
							"ibrs_enhanced",
							"fsgsbase",
							"tsc_adjust",
							"bmi1",
							"hle",
							"avx2",
							"smep",
							"bmi2",
							"erms",
							"invpcid",
							"rtm",
							"avx512f",
							"avx512dq",
							"rdseed",
							"adx",
							"smap",
							"clflushopt",
							"clwb",
							"avx512cd",
							"avx512bw",
							"avx512vl",
							"xsaveopt",
							"xsavec",
							"xgetbv1",
							"xsaves",
							"arat",
							"avx512_vnni",
							"md_clear",
							"arch_capabilities"
						],
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						]
					}
				]
			},
			"topology": {
				"architecture": "smp",
				"nodes": [
					{
						"id": 0,
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						],
						"caches": [
							{
								"level": 1,
								"type": "instruction",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 1,
								"type": "data",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 2,
								"type": "unified",
								"size_bytes": 1048576,
								"logical_processors": [
									0
								]
							},
							{
								"level": 3,
								"type": "unified",
								"size_bytes": 34603008,
								"logical_processors": [
									0
								]
							}
						],
						"distances": [
							10
						],
						"memory": {
							"total_physical_bytes": 8589934592,
							"total_usable_bytes": 8056090624,
							"supported_page_sizes": [
								1073741824,
								2097152
							],
							"modules": null
						}
					}
				]
			},
			"network": {
				"nics": [
					{
						"name": "eth0",
						"mac_address": "42:01:0a:00:00:ce",
						"is_virtual": false,
						"capabilities": [
							{
								"name": "auto-negotiation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "pause-frame-use",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-checksumming",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-checksumming",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv4",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-ip-generic",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv6",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-fcoe-crc",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-sctp",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather-fraglist",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tcp-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-ecn-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tcp-mangleid-segmentation",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tx-tcp6-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-receive-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "large-receive-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "ntuple-filters",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "receive-hashing",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "highdma",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "rx-vlan-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "vlan-challenged",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-lockless",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "netns-local",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-robust",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-fcoe-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip4-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip6-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-partial",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tunnel-remcsum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-sctp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-esp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-list",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp-gro-forwarding",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "rx-gro-list",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tls-hw-rx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "fcoe-mtu",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-nocache-copy",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "loopback",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-fcs",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-all",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-stag-hw-insert",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-hw-parse",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "l2-fwd-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "hw-tc-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-tx-csum-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp_tunnel-port-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-tx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-gro-hw",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-record",
								"is_enabled": false,
								"can_enable": false
							}
						],
						"speed": "Unknown!",
						"duplex": "Unknown!(255)"
					}
				]
			},
			"gpu": {
				"cards": null
			},
			"chassis": {
				"asset_tag": "",
				"serial_number": "unknown",
				"type": "1",
				"type_description": "Other",
				"vendor": "Google",
				"version": ""
			},
			"bios": {
				"vendor": "Google",
				"version": "Google",
				"date": "06/07/2024"
			},
			"baseboard": {
				"asset_tag": "F89BD669-7C29-7733-7C14-BCDA0501EF3D",
				"serial_number": "unknown",
				"vendor": "Google",
				"version": "",
				"product": "Google Compute Engine"
			},
			"product": {
				"family": "",
				"name": "Google Compute Engine",
				"vendor": "Google",
				"serial_number": "unknown",
				"uuid": "unknown",
				"sku": "",
				"version": ""
			},
			"pci": {
				"Devices": [
					{
						"driver": "",
						"address": "0000:00:00.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "1237",
							"name": "440FX - 82441FX PMC [Natoma]"
						},
						"revision": "0x02",
						"subsystem": {
							"id": "1100",
							"name": "Qemu virtual machine"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "00",
							"name": "Host bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7110",
							"name": "82371AB/EB/MB PIIX4 ISA"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "01",
							"name": "ISA bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.3",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7113",
							"name": "82371AB/EB/MB PIIX4 ACPI"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "80",
							"name": "Bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:03.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1004",
							"name": "Virtio SCSI"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0008",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "00",
							"name": "Non-VGA unclassified device"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:04.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1000",
							"name": "Virtio network device"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0001",
							"name": "unknown"
						},
						"class": {
							"id": "02",
							"name": "Network controller"
						},
						"subclass": {
							"id": "00",
							"name": "Ethernet controller"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:05.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1005",
							"name": "Virtio RNG"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0004",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "ff",
							"name": "unknown"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					}
				]
			}
		},
		"Kernel": {
			"Version": "Linux version 4.18.0-513.11.1.el8_9.cloud.0.1.x86_64 (mockbuild@iad1-prod-build001.bld.equ.rockylinux.org) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-20) (GCC)) #1 SMP Thu Jan 11 13:51:41 UTC 2024\n",
			"Modules": "tcp_diag 16384 0 - Live 0x0000000000000000\nbinfmt_misc 24576 1 - Live 0x0000000000000000\ninet_diag 24576 1 tcp_diag, Live 0x0000000000000000\nnfsv3 57344 1 - Live 0x0000000000000000\nrpcsec_gss_krb5 45056 0 - Live 0x0000000000000000\nnfsv4 917504 2 - Live 0x0000000000000000\ndns_resolver 16384 1 nfsv4, Live 0x0000000000000000\nnfs 425984 4 nfsv3,nfsv4, Live 0x0000000000000000\nfscache 389120 1 nfs, Live 0x0000000000000000\nintel_rapl_msr 16384 0 - Live 0x0000000000000000\nintel_rapl_common 24576 1 intel_rapl_msr, Live 0x0000000000000000\nintel_uncore_frequency_common 16384 0 - Live 0x0000000000000000\nisst_if_common 16384 0 - Live 0x0000000000000000\nnfit 65536 0 - Live 0x0000000000000000\nlibnvdimm 200704 1 nfit, Live 0x0000000000000000\ncrct10dif_pclmul 16384 1 - Live 0x0000000000000000\ncrc32_pclmul 16384 0 - Live 0x0000000000000000\nghash_clmulni_intel 16384 0 - Live 0x0000000000000000\nrapl 20480 0 - Live 0x0000000000000000\ni2c_piix4 24576 0 - Live 0x0000000000000000\nvfat 20480 1 - Live 0x0000000000000000\nfat 86016 1 vfat, Live 0x0000000000000000\npcspkr 16384 0 - Live 0x0000000000000000\nnfsd 548864 13 - Live 0x0000000000000000\nauth_rpcgss 139264 2 rpcsec_gss_krb5,nfsd, Live 0x0000000000000000\nnfs_acl 16384 2 nfsv3,nfsd, Live 0x0000000000000000\nlockd 126976 3 nfsv3,nfs,nfsd, Live 0x0000000000000000\ngrace 16384 2 nfsd,lockd, Live 0x0000000000000000\nsunrpc 585728 32 nfsv3,rpcsec_gss_krb5,nfsv4,nfs,nfsd,auth_rpcgss,nfs_acl,lockd, Live 0x0000000000000000\nxfs 1593344 1 - Live 0x0000000000000000\nlibcrc32c 16384 1 xfs, Live 0x0000000000000000\nsd_mod 57344 2 - Live 0x0000000000000000\nsg 40960 0 - Live 0x0000000000000000\nvirtio_net 61440 0 - Live 0x0000000000000000\ncrc32c_intel 24576 1 - Live 0x0000000000000000\nserio_raw 16384 0 - Live 0x0000000000000000\nnet_failover 24576 1 virtio_net, Live 0x0000000000000000\nvirtio_scsi 20480 2 - Live 0x0000000000000000\nfailover 16384 1 net_failover, Live 0x0000000000000000\nnvme 45056 0 - Live 0x0000000000000000\nnvme_core 139264 1 nvme, Live 0x0000000000000000\nt10_pi 16384 2 sd_mod,nvme_core, Live 0x0000000000000000\n",
			"Drivers": "Character devices:\n  1 mem\n  4 /dev/vc/0\n  4 tty\n  4 ttyS\n  5 /dev/tty\n  5 /dev/console\n  5 /dev/ptmx\n  7 vcs\n 10 misc\n 13 input\n 21 sg\n 29 fb\n128 ptm\n136 pts\n162 raw\n180 usb\n188 ttyUSB\n189 usb_device\n202 cpu/msr\n203 cpu/cpuid\n240 dimmctl\n241 ndctl\n242 nvme-generic\n243 nvme\n244 hidraw\n245 ttyDBC\n246 usbmon\n247 bsg\n248 watchdog\n249 ptp\n250 pps\n251 rtc\n252 dax\n253 tpm\n254 gpiochip\n\nBlock devices:\n  8 sd\n  9 md\n 65 sd\n 66 sd\n 67 sd\n 68 sd\n 69 sd\n 70 sd\n 71 sd\n128 sd\n129 sd\n130 sd\n131 sd\n132 sd\n133 sd\n134 sd\n135 sd\n254 mdp\n259 blkext\n"
		},
		"Stderr": "WARNING: \n/sys/class/drm does not exist on this system (likely the host system is a\nvirtual machine or container with no graphics). Therefore,\nGPUInfo.GraphicsCards will be an empty array.\nWARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied\nWARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied\nWARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied\nWARNING: Unable to read product_uuid: open /sys/class/dmi/id/product_uuid: permission denied\n",
		"Err": ""
	},
	{
		"Hostname": "hpcslurm-debugnodeset-13",
		"Info": {
			"memory": {
				"total_physical_bytes": 8589934592,
				"total_usable_bytes": 8056098816,
				"supported_page_sizes": [
					1073741824,
					2097152
				],
				"modules": null
			},
			"block": {
				"total_size_bytes": 53687091200,
				"disks": [
					{
						"name": "sda",
						"size_bytes": 53687091200,
						"physical_block_size_bytes": 4096,
						"drive_type": "hdd",
						"removable": false,
						"storage_controller": "scsi",
						"bus_path": "pci-0000:00:03.0-scsi-0:0:1:0",
						"vendor": "Google",
						"model": "PersistentDisk",
						"serial_number": "persistent-disk-0",
						"wwn": "unknown",
						"partitions": [
							{
								"name": "sda1",
								"label": "EFI\\x20System\\x20Partition",
								"mount_point": "/boot/efi",
								"size_bytes": 209715200,
								"type": "vfat",
								"read_only": false,
								"uuid": "a407d4b7-cfe4-4f7e-b9fc-ee7799ba3b84",
								"filesystem_label": "unknown"
							},
							{
								"name": "sda2",
								"label": "unknown",
								"mount_point": "/",
								"size_bytes": 53475328000,
								"type": "xfs",
								"read_only": false,
								"uuid": "144c8c6f-9c84-47c9-b637-8b7723fdb3ef",
								"filesystem_label": "root"
							}
						]
					}
				]
			},
			"cpu": {
				"total_cores": 1,
				"total_threads": 1,
				"processors": [
					{
						"id": 0,
						"total_cores": 1,
						"total_threads": 1,
						"vendor": "GenuineIntel",
						"model": "Intel(R) Xeon(R) CPU @ 2.80GHz",
						"capabilities": [
							"fpu",
							"vme",
							"de",
							"pse",
							"tsc",
							"msr",
							"pae",
							"mce",
							"cx8",
							"apic",
							"sep",
							"mtrr",
							"pge",
							"mca",
							"cmov",
							"pat",
							"pse36",
							"clflush",
							"mmx",
							"fxsr",
							"sse",
							"sse2",
							"ss",
							"ht",
							"syscall",
							"nx",
							"pdpe1gb",
							"rdtscp",
							"lm",
							"constant_tsc",
							"rep_good",
							"nopl",
							"xtopology",
							"nonstop_tsc",
							"cpuid",
							"tsc_known_freq",
							"pni",
							"pclmulqdq",
							"ssse3",
							"fma",
							"cx16",
							"pcid",
							"sse4_1",
							"sse4_2",
							"x2apic",
							"movbe",
							"popcnt",
							"aes",
							"xsave",
							"avx",
							"f16c",
							"rdrand",
							"hypervisor",
							"lahf_lm",
							"abm",
							"3dnowprefetch",
							"invpcid_single",
							"ssbd",
							"ibrs",
							"ibpb",
							"stibp",
							"ibrs_enhanced",
							"fsgsbase",
							"tsc_adjust",
							"bmi1",
							"hle",
							"avx2",
							"smep",
							"bmi2",
							"erms",
							"invpcid",
							"rtm",
							"avx512f",
							"avx512dq",
							"rdseed",
							"adx",
							"smap",
							"clflushopt",
							"clwb",
							"avx512cd",
							"avx512bw",
							"avx512vl",
							"xsaveopt",
							"xsavec",
							"xgetbv1",
							"xsaves",
							"arat",
							"avx512_vnni",
							"md_clear",
							"arch_capabilities"
						],
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						]
					}
				]
			},
			"topology": {
				"architecture": "smp",
				"nodes": [
					{
						"id": 0,
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						],
						"caches": [
							{
								"level": 1,
								"type": "instruction",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 1,
								"type": "data",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 2,
								"type": "unified",
								"size_bytes": 1048576,
								"logical_processors": [
									0
								]
							},
							{
								"level": 3,
								"type": "unified",
								"size_bytes": 34603008,
								"logical_processors": [
									0
								]
							}
						],
						"distances": [
							10
						],
						"memory": {
							"total_physical_bytes": 8589934592,
							"total_usable_bytes": 8056098816,
							"supported_page_sizes": [
								1073741824,
								2097152
							],
							"modules": null
						}
					}
				]
			},
			"network": {
				"nics": [
					{
						"name": "eth0",
						"mac_address": "42:01:0a:00:00:ec",
						"is_virtual": false,
						"capabilities": [
							{
								"name": "auto-negotiation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "pause-frame-use",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-checksumming",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-checksumming",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv4",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-ip-generic",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv6",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-fcoe-crc",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-sctp",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather-fraglist",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tcp-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-ecn-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tcp-mangleid-segmentation",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tx-tcp6-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-receive-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "large-receive-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "ntuple-filters",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "receive-hashing",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "highdma",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "rx-vlan-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "vlan-challenged",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-lockless",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "netns-local",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-robust",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-fcoe-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip4-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip6-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-partial",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tunnel-remcsum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-sctp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-esp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-list",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp-gro-forwarding",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "rx-gro-list",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tls-hw-rx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "fcoe-mtu",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-nocache-copy",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "loopback",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-fcs",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-all",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-stag-hw-insert",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-hw-parse",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "l2-fwd-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "hw-tc-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-tx-csum-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp_tunnel-port-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-tx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-gro-hw",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-record",
								"is_enabled": false,
								"can_enable": false
							}
						],
						"speed": "Unknown!",
						"duplex": "Unknown!(255)"
					}
				]
			},
			"gpu": {
				"cards": null
			},
			"chassis": {
				"asset_tag": "",
				"serial_number": "unknown",
				"type": "1",
				"type_description": "Other",
				"vendor": "Google",
				"version": ""
			},
			"bios": {
				"vendor": "Google",
				"version": "Google",
				"date": "06/07/2024"
			},
			"baseboard": {
				"asset_tag": "A99E7F52-07DD-2D81-7685-70ABA074CB3D",
				"serial_number": "unknown",
				"vendor": "Google",
				"version": "",
				"product": "Google Compute Engine"
			},
			"product": {
				"family": "",
				"name": "Google Compute Engine",
				"vendor": "Google",
				"serial_number": "unknown",
				"uuid": "unknown",
				"sku": "",
				"version": ""
			},
			"pci": {
				"Devices": [
					{
						"driver": "",
						"address": "0000:00:00.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "1237",
							"name": "440FX - 82441FX PMC [Natoma]"
						},
						"revision": "0x02",
						"subsystem": {
							"id": "1100",
							"name": "Qemu virtual machine"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "00",
							"name": "Host bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7110",
							"name": "82371AB/EB/MB PIIX4 ISA"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "01",
							"name": "ISA bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.3",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7113",
							"name": "82371AB/EB/MB PIIX4 ACPI"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "80",
							"name": "Bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:03.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1004",
							"name": "Virtio SCSI"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0008",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "00",
							"name": "Non-VGA unclassified device"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:04.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1000",
							"name": "Virtio network device"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0001",
							"name": "unknown"
						},
						"class": {
							"id": "02",
							"name": "Network controller"
						},
						"subclass": {
							"id": "00",
							"name": "Ethernet controller"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:05.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1005",
							"name": "Virtio RNG"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0004",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "ff",
							"name": "unknown"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					}
				]
			}
		},
		"Kernel": {
			"Version": "Linux version 4.18.0-513.24.1.el8_9.x86_64 (mockbuild@iad1-prod-build001.bld.equ.rockylinux.org) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-20) (GCC)) #1 SMP Thu Apr 4 18:13:02 UTC 2024\n",
			"Modules": "tcp_diag 16384 0 - Live 0x0000000000000000\ninet_diag 24576 1 tcp_diag, Live 0x0000000000000000\nbinfmt_misc 24576 1 - Live 0x0000000000000000\nnfsv3 57344 1 - Live 0x0000000000000000\nrpcsec_gss_krb5 45056 0 - Live 0x0000000000000000\nnfsv4 917504 2 - Live 0x0000000000000000\ndns_resolver 16384 1 nfsv4, Live 0x0000000000000000\nnfs 425984 4 nfsv3,nfsv4, Live 0x0000000000000000\nfscache 389120 1 nfs, Live 0x0000000000000000\nintel_rapl_msr 16384 0 - Live 0x0000000000000000\nintel_rapl_common 24576 1 intel_rapl_msr, Live 0x0000000000000000\nintel_uncore_frequency_common 16384 0 - Live 0x0000000000000000\nisst_if_common 16384 0 - Live 0x0000000000000000\nnfit 65536 0 - Live 0x0000000000000000\nlibnvdimm 200704 1 nfit, Live 0x0000000000000000\ncrct10dif_pclmul 16384 1 - Live 0x0000000000000000\ncrc32_pclmul 16384 0 - Live 0x0000000000000000\nghash_clmulni_intel 16384 0 - Live 0x0000000000000000\nrapl 20480 0 - Live 0x0000000000000000\ni2c_piix4 24576 0 - Live 0x0000000000000000\nvfat 20480 1 - Live 0x0000000000000000\nfat 86016 1 vfat, Live 0x0000000000000000\npcspkr 16384 0 - Live 0x0000000000000000\nnfsd 548864 13 - Live 0x0000000000000000\nauth_rpcgss 139264 2 rpcsec_gss_krb5,nfsd, Live 0x0000000000000000\nnfs_acl 16384 2 nfsv3,nfsd, Live 0x0000000000000000\nlockd 126976 3 nfsv3,nfs,nfsd, Live 0x0000000000000000\ngrace 16384 2 nfsd,lockd, Live 0x0000000000000000\nsunrpc 585728 32 nfsv3,rpcsec_gss_krb5,nfsv4,nfs,nfsd,auth_rpcgss,nfs_acl,lockd, Live 0x0000000000000000\nxfs 1593344 1 - Live 0x0000000000000000\nlibcrc32c 16384 1 xfs, Live 0x0000000000000000\nsd_mod 57344 2 - Live 0x0000000000000000\nsg 40960 0 - Live 0x0000000000000000\nvirtio_net 61440 0 - Live 0x0000000000000000\ncrc32c_intel 24576 1 - Live 0x0000000000000000\nvirtio_scsi 20480 2 - Live 0x0000000000000000\nserio_raw 16384 0 - Live 0x0000000000000000\nnet_failover 24576 1 virtio_net, Live 0x0000000000000000\nfailover 16384 1 net_failover, Live 0x0000000000000000\nnvme 45056 0 - Live 0x0000000000000000\nnvme_core 139264 1 nvme, Live 0x0000000000000000\nt10_pi 16384 2 sd_mod,nvme_core, Live 0x0000000000000000\n",
			"Drivers": "Character devices:\n  1 mem\n  4 /dev/vc/0\n  4 tty\n  4 ttyS\n  5 /dev/tty\n  5 /dev/console\n  5 /dev/ptmx\n  7 vcs\n 10 misc\n 13 input\n 21 sg\n 29 fb\n128 ptm\n136 pts\n162 raw\n180 usb\n188 ttyUSB\n189 usb_device\n202 cpu/msr\n203 cpu/cpuid\n240 dimmctl\n241 ndctl\n242 nvme-generic\n243 nvme\n244 hidraw\n245 ttyDBC\n246 usbmon\n247 bsg\n248 watchdog\n249 ptp\n250 pps\n251 rtc\n252 dax\n253 tpm\n254 gpiochip\n\nBlock devices:\n  8 sd\n  9 md\n 65 sd\n 66 sd\n 67 sd\n 68 sd\n 69 sd\n 70 sd\n 71 sd\n128 sd\n129 sd\n130 sd\n131 sd\n132 sd\n133 sd\n134 sd\n135 sd\n254 mdp\n259 blkext\n"
		},
		"Stderr": "WARNING: \n/sys/class/drm does not exist on this system (likely the host system is a\nvirtual machine or container with no graphics). Therefore,\nGPUInfo.GraphicsCards will be an empty array.\nWARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied\nWARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied\nWARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied\nWARNING: Unable to read product_uuid: open /sys/class/dmi/id/product_uuid: permission denied\n",
		"Err": ""
	},
	{
		"Hostname": "hpcslurm-debugnodeset-5",
		"Info": {
			"memory": {
				"total_physical_bytes": 8589934592,
				"total_usable_bytes": 8056098816,
				"supported_page_sizes": [
					1073741824,
					2097152
				],
				"modules": null
			},
			"block": {
				"total_size_bytes": 53687091200,
				"disks": [
					{
						"name": "sda",
						"size_bytes": 53687091200,
						"physical_block_size_bytes": 4096,
						"drive_type": "hdd",
						"removable": false,
						"storage_controller": "scsi",
						"bus_path": "pci-0000:00:03.0-scsi-0:0:1:0",
						"vendor": "Google",
						"model": "PersistentDisk",
						"serial_number": "persistent-disk-0",
						"wwn": "unknown",
						"partitions": [
							{
								"name": "sda1",
								"label": "EFI\\x20System\\x20Partition",
								"mount_point": "/boot/efi",
								"size_bytes": 209715200,
								"type": "vfat",
								"read_only": false,
								"uuid": "a407d4b7-cfe4-4f7e-b9fc-ee7799ba3b84",
								"filesystem_label": "unknown"
							},
							{
								"name": "sda2",
								"label": "unknown",
								"mount_point": "/",
								"size_bytes": 53475328000,
								"type": "xfs",
								"read_only": false,
								"uuid": "144c8c6f-9c84-47c9-b637-8b7723fdb3ef",
								"filesystem_label": "root"
							}
						]
					}
				]
			},
			"cpu": {
				"total_cores": 1,
				"total_threads": 1,
				"processors": [
					{
						"id": 0,
						"total_cores": 1,
						"total_threads": 1,
						"vendor": "GenuineIntel",
						"model": "Intel(R) Xeon(R) CPU @ 2.80GHz",
						"capabilities": [
							"fpu",
							"vme",
							"de",
							"pse",
							"tsc",
							"msr",
							"pae",
							"mce",
							"cx8",
							"apic",
							"sep",
							"mtrr",
							"pge",
							"mca",
							"cmov",
							"pat",
							"pse36",
							"clflush",
							"mmx",
							"fxsr",
							"sse",
							"sse2",
							"ss",
							"ht",
							"syscall",
							"nx",
							"pdpe1gb",
							"rdtscp",
							"lm",
							"constant_tsc",
							"rep_good",
							"nopl",
							"xtopology",
							"nonstop_tsc",
							"cpuid",
							"tsc_known_freq",
							"pni",
							"pclmulqdq",
							"ssse3",
							"fma",
							"cx16",
							"pcid",
							"sse4_1",
							"sse4_2",
							"x2apic",
							"movbe",
							"popcnt",
							"aes",
							"xsave",
							"avx",
							"f16c",
							"rdrand",
							"hypervisor",
							"lahf_lm",
							"abm",
							"3dnowprefetch",
							"invpcid_single",
							"ssbd",
							"ibrs",
							"ibpb",
							"stibp",
							"ibrs_enhanced",
							"fsgsbase",
							"tsc_adjust",
							"bmi1",
							"hle",
							"avx2",
							"smep",
							"bmi2",
							"erms",
							"invpcid",
							"rtm",
							"avx512f",
							"avx512dq",
							"rdseed",
							"adx",
							"smap",
							"clflushopt",
							"clwb",
							"avx512cd",
							"avx512bw",
							"avx512vl",
							"xsaveopt",
							"xsavec",
							"xgetbv1",
							"xsaves",
							"arat",
							"avx512_vnni",
							"md_clear",
							"arch_capabilities"
						],
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						]
					}
				]
			},
			"topology": {
				"architecture": "smp",
				"nodes": [
					{
						"id": 0,
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						],
						"caches": [
							{
								"level": 1,
								"type": "instruction",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 1,
								"type": "data",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 2,
								"type": "unified",
								"size_bytes": 1048576,
								"logical_processors": [
									0
								]
							},
							{
								"level": 3,
								"type": "unified",
								"size_bytes": 34603008,
								"logical_processors": [
									0
								]
							}
						],
						"distances": [
							10
						],
						"memory": {
							"total_physical_bytes": 8589934592,
							"total_usable_bytes": 8056098816,
							"supported_page_sizes": [
								1073741824,
								2097152
							],
							"modules": null
						}
					}
				]
			},
			"network": {
				"nics": [
					{
						"name": "eth0",
						"mac_address": "42:01:0a:00:00:d1",
						"is_virtual": false,
						"capabilities": [
							{
								"name": "auto-negotiation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "pause-frame-use",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-checksumming",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-checksumming",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv4",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-ip-generic",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv6",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-fcoe-crc",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-sctp",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather-fraglist",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tcp-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-ecn-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tcp-mangleid-segmentation",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tx-tcp6-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-receive-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "large-receive-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "ntuple-filters",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "receive-hashing",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "highdma",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "rx-vlan-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "vlan-challenged",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-lockless",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "netns-local",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-robust",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-fcoe-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip4-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip6-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-partial",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tunnel-remcsum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-sctp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-esp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-list",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp-gro-forwarding",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "rx-gro-list",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tls-hw-rx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "fcoe-mtu",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-nocache-copy",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "loopback",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-fcs",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-all",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-stag-hw-insert",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-hw-parse",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "l2-fwd-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "hw-tc-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-tx-csum-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp_tunnel-port-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-tx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-gro-hw",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-record",
								"is_enabled": false,
								"can_enable": false
							}
						],
						"speed": "Unknown!",
						"duplex": "Unknown!(255)"
					}
				]
			},
			"gpu": {
				"cards": null
			},
			"chassis": {
				"asset_tag": "",
				"serial_number": "unknown",
				"type": "1",
				"type_description": "Other",
				"vendor": "Google",
				"version": ""
			},
			"bios": {
				"vendor": "Google",
				"version": "Google",
				"date": "06/07/2024"
			},
			"baseboard": {
				"asset_tag": "E6692EF0-F9AC-FEB4-5D65-26182058BAC0",
				"serial_number": "unknown",
				"vendor": "Google",
				"version": "",
				"product": "Google Compute Engine"
			},
			"product": {
				"family": "",
				"name": "Google Compute Engine",
				"vendor": "Google",
				"serial_number": "unknown",
				"uuid": "unknown",
				"sku": "",
				"version": ""
			},
			"pci": {
				"Devices": [
					{
						"driver": "",
						"address": "0000:00:00.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "1237",
							"name": "440FX - 82441FX PMC [Natoma]"
						},
						"revision": "0x02",
						"subsystem": {
							"id": "1100",
							"name": "Qemu virtual machine"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "00",
							"name": "Host bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7110",
							"name": "82371AB/EB/MB PIIX4 ISA"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "01",
							"name": "ISA bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.3",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7113",
							"name": "82371AB/EB/MB PIIX4 ACPI"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "80",
							"name": "Bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:03.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1004",
							"name": "Virtio SCSI"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0008",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "00",
							"name": "Non-VGA unclassified device"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:04.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1000",
							"name": "Virtio network device"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0001",
							"name": "unknown"
						},
						"class": {
							"id": "02",
							"name": "Network controller"
						},
						"subclass": {
							"id": "00",
							"name": "Ethernet controller"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:05.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1005",
							"name": "Virtio RNG"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0004",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "ff",
							"name": "unknown"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					}
				]
			}
		},
		"Kernel": {
			"Version": "Linux version 4.18.0-513.24.1.el8_9.x86_64 (mockbuild@iad1-prod-build001.bld.equ.rockylinux.org) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-20) (GCC)) #1 SMP Thu Apr 4 18:13:02 UTC 2024\n",
			"Modules": "tcp_diag 16384 0 - Live 0x0000000000000000\ninet_diag 24576 1 tcp_diag, Live 0x0000000000000000\nbinfmt_misc 24576 1 - Live 0x0000000000000000\nnfsv3 57344 1 - Live 0x0000000000000000\nrpcsec_gss_krb5 45056 0 - Live 0x0000000000000000\nnfsv4 917504 2 - Live 0x0000000000000000\ndns_resolver 16384 1 nfsv4, Live 0x0000000000000000\nnfs 425984 4 nfsv3,nfsv4, Live 0x0000000000000000\nfscache 389120 1 nfs, Live 0x0000000000000000\nintel_rapl_msr 16384 0 - Live 0x0000000000000000\nintel_rapl_common 24576 1 intel_rapl_msr, Live 0x0000000000000000\nintel_uncore_frequency_common 16384 0 - Live 0x0000000000000000\nisst_if_common 16384 0 - Live 0x0000000000000000\nnfit 65536 0 - Live 0x0000000000000000\nlibnvdimm 200704 1 nfit, Live 0x0000000000000000\ncrct10dif_pclmul 16384 1 - Live 0x0000000000000000\ncrc32_pclmul 16384 0 - Live 0x0000000000000000\nghash_clmulni_intel 16384 0 - Live 0x0000000000000000\nrapl 20480 0 - Live 0x0000000000000000\nvfat 20480 1 - Live 0x0000000000000000\nfat 86016 1 vfat, Live 0x0000000000000000\ni2c_piix4 24576 0 - Live 0x0000000000000000\npcspkr 16384 0 - Live 0x0000000000000000\nnfsd 548864 13 - Live 0x0000000000000000\nauth_rpcgss 139264 2 rpcsec_gss_krb5,nfsd, Live 0x0000000000000000\nnfs_acl 16384 2 nfsv3,nfsd, Live 0x0000000000000000\nlockd 126976 3 nfsv3,nfs,nfsd, Live 0x0000000000000000\ngrace 16384 2 nfsd,lockd, Live 0x0000000000000000\nsunrpc 585728 32 nfsv3,rpcsec_gss_krb5,nfsv4,nfs,nfsd,auth_rpcgss,nfs_acl,lockd, Live 0x0000000000000000\nxfs 1593344 1 - Live 0x0000000000000000\nlibcrc32c 16384 1 xfs, Live 0x0000000000000000\nsd_mod 57344 2 - Live 0x0000000000000000\nsg 40960 0 - Live 0x0000000000000000\nvirtio_net 61440 0 - Live 0x0000000000000000\ncrc32c_intel 24576 1 - Live 0x0000000000000000\nserio_raw 16384 0 - Live 0x0000000000000000\nnet_failover 24576 1 virtio_net, Live 0x0000000000000000\nvirtio_scsi 20480 2 - Live 0x0000000000000000\nfailover 16384 1 net_failover, Live 0x0000000000000000\nnvme 45056 0 - Live 0x0000000000000000\nnvme_core 139264 1 nvme, Live 0x0000000000000000\nt10_pi 16384 2 sd_mod,nvme_core, Live 0x0000000000000000\n",
			"Drivers": "Character devices:\n  1 mem\n  4 /dev/vc/0\n  4 tty\n  4 ttyS\n  5 /dev/tty\n  5 /dev/console\n  5 /dev/ptmx\n  7 vcs\n 10 misc\n 13 input\n 21 sg\n 29 fb\n128 ptm\n136 pts\n162 raw\n180 usb\n188 ttyUSB\n189 usb_device\n202 cpu/msr\n203 cpu/cpuid\n240 dimmctl\n241 ndctl\n242 nvme-generic\n243 nvme\n244 hidraw\n245 ttyDBC\n246 usbmon\n247 bsg\n248 watchdog\n249 ptp\n250 pps\n251 rtc\n252 dax\n253 tpm\n254 gpiochip\n\nBlock devices:\n  8 sd\n  9 md\n 65 sd\n 66 sd\n 67 sd\n 68 sd\n 69 sd\n 70 sd\n 71 sd\n128 sd\n129 sd\n130 sd\n131 sd\n132 sd\n133 sd\n134 sd\n135 sd\n254 mdp\n259 blkext\n"
		},
		"Stderr": "WARNING: \n/sys/class/drm does not exist on this system (likely the host system is a\nvirtual machine or container with no graphics). Therefore,\nGPUInfo.GraphicsCards will be an empty array.\nWARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied\nWARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied\nWARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied\nWARNING: Unable to read product_uuid: open /sys/class/dmi/id/product_uuid: permission denied\n",
		"Err": ""
	},
	{
		"Hostname": "hpcslurm-debugnodeset-9",
		"Info": {
			"memory": {
				"total_physical_bytes": 8589934592,
				"total_usable_bytes": 8056090624,
				"supported_page_sizes": [
					1073741824,
					2097152
				],
				"modules": null
			},
			"block": {
				"total_size_bytes": 53687091200,
				"disks": [
					{
						"name": "sda",
						"size_bytes": 53687091200,
						"physical_block_size_bytes": 4096,
						"drive_type": "hdd",
						"removable": false,
						"storage_controller": "scsi",
						"bus_path": "pci-0000:00:03.0-scsi-0:0:1:0",
						"vendor": "Google",
						"model": "PersistentDisk",
						"serial_number": "persistent-disk-0",
						"wwn": "unknown",
						"partitions": [
							{
								"name": "sda1",
								"label": "EFI\\x20System\\x20Partition",
								"mount_point": "/boot/efi",
								"size_bytes": 209715200,
								"type": "vfat",
								"read_only": false,
								"uuid": "a407d4b7-cfe4-4f7e-b9fc-ee7799ba3b84",
								"filesystem_label": "unknown"
							},
							{
								"name": "sda2",
								"label": "unknown",
								"mount_point": "/",
								"size_bytes": 53475328000,
								"type": "xfs",
								"read_only": false,
								"uuid": "144c8c6f-9c84-47c9-b637-8b7723fdb3ef",
								"filesystem_label": "root"
							}
						]
					}
				]
			},
			"cpu": {
				"total_cores": 1,
				"total_threads": 1,
				"processors": [
					{
						"id": 0,
						"total_cores": 1,
						"total_threads": 1,
						"vendor": "GenuineIntel",
						"model": "Intel(R) Xeon(R) CPU @ 2.80GHz",
						"capabilities": [
							"fpu",
							"vme",
							"de",
							"pse",
							"tsc",
							"msr",
							"pae",
							"mce",
							"cx8",
							"apic",
							"sep",
							"mtrr",
							"pge",
							"mca",
							"cmov",
							"pat",
							"pse36",
							"clflush",
							"mmx",
							"fxsr",
							"sse",
							"sse2",
							"ss",
							"ht",
							"syscall",
							"nx",
							"pdpe1gb",
							"rdtscp",
							"lm",
							"constant_tsc",
							"rep_good",
							"nopl",
							"xtopology",
							"nonstop_tsc",
							"cpuid",
							"tsc_known_freq",
							"pni",
							"pclmulqdq",
							"ssse3",
							"fma",
							"cx16",
							"pcid",
							"sse4_1",
							"sse4_2",
							"x2apic",
							"movbe",
							"popcnt",
							"aes",
							"xsave",
							"avx",
							"f16c",
							"rdrand",
							"hypervisor",
							"lahf_lm",
							"abm",
							"3dnowprefetch",
							"invpcid_single",
							"ssbd",
							"ibrs",
							"ibpb",
							"stibp",
							"ibrs_enhanced",
							"fsgsbase",
							"tsc_adjust",
							"bmi1",
							"hle",
							"avx2",
							"smep",
							"bmi2",
							"erms",
							"invpcid",
							"rtm",
							"avx512f",
							"avx512dq",
							"rdseed",
							"adx",
							"smap",
							"clflushopt",
							"clwb",
							"avx512cd",
							"avx512bw",
							"avx512vl",
							"xsaveopt",
							"xsavec",
							"xgetbv1",
							"xsaves",
							"arat",
							"avx512_vnni",
							"md_clear",
							"arch_capabilities"
						],
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						]
					}
				]
			},
			"topology": {
				"architecture": "smp",
				"nodes": [
					{
						"id": 0,
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						],
						"caches": [
							{
								"level": 1,
								"type": "instruction",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 1,
								"type": "data",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 2,
								"type": "unified",
								"size_bytes": 1048576,
								"logical_processors": [
									0
								]
							},
							{
								"level": 3,
								"type": "unified",
								"size_bytes": 34603008,
								"logical_processors": [
									0
								]
							}
						],
						"distances": [
							10
						],
						"memory": {
							"total_physical_bytes": 8589934592,
							"total_usable_bytes": 8056090624,
							"supported_page_sizes": [
								1073741824,
								2097152
							],
							"modules": null
						}
					}
				]
			},
			"network": {
				"nics": [
					{
						"name": "eth0",
						"mac_address": "42:01:0a:00:00:d4",
						"is_virtual": false,
						"capabilities": [
							{
								"name": "auto-negotiation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "pause-frame-use",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-checksumming",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-checksumming",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv4",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-ip-generic",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv6",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-fcoe-crc",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-sctp",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather-fraglist",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tcp-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-ecn-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tcp-mangleid-segmentation",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tx-tcp6-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-receive-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "large-receive-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "ntuple-filters",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "receive-hashing",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "highdma",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "rx-vlan-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "vlan-challenged",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-lockless",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "netns-local",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-robust",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-fcoe-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip4-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip6-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-partial",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tunnel-remcsum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-sctp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-esp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-list",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp-gro-forwarding",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "rx-gro-list",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tls-hw-rx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "fcoe-mtu",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-nocache-copy",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "loopback",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-fcs",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-all",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-stag-hw-insert",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-hw-parse",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "l2-fwd-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "hw-tc-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-tx-csum-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp_tunnel-port-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-tx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-gro-hw",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-record",
								"is_enabled": false,
								"can_enable": false
							}
						],
						"speed": "Unknown!",
						"duplex": "Unknown!(255)"
					}
				]
			},
			"gpu": {
				"cards": null
			},
			"chassis": {
				"asset_tag": "",
				"serial_number": "unknown",
				"type": "1",
				"type_description": "Other",
				"vendor": "Google",
				"version": ""
			},
			"bios": {
				"vendor": "Google",
				"version": "Google",
				"date": "06/27/2024"
			},
			"baseboard": {
				"asset_tag": "6CC588E4-F210-7C45-43EB-51DC3CF8E1D9",
				"serial_number": "unknown",
				"vendor": "Google",
				"version": "",
				"product": "Google Compute Engine"
			},
			"product": {
				"family": "",
				"name": "Google Compute Engine",
				"vendor": "Google",
				"serial_number": "unknown",
				"uuid": "unknown",
				"sku": "",
				"version": ""
			},
			"pci": {
				"Devices": [
					{
						"driver": "",
						"address": "0000:00:00.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "1237",
							"name": "440FX - 82441FX PMC [Natoma]"
						},
						"revision": "0x02",
						"subsystem": {
							"id": "1100",
							"name": "Qemu virtual machine"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "00",
							"name": "Host bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7110",
							"name": "82371AB/EB/MB PIIX4 ISA"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "01",
							"name": "ISA bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.3",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7113",
							"name": "82371AB/EB/MB PIIX4 ACPI"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "80",
							"name": "Bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:03.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1004",
							"name": "Virtio SCSI"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0008",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "00",
							"name": "Non-VGA unclassified device"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:04.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1000",
							"name": "Virtio network device"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0001",
							"name": "unknown"
						},
						"class": {
							"id": "02",
							"name": "Network controller"
						},
						"subclass": {
							"id": "00",
							"name": "Ethernet controller"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:05.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1005",
							"name": "Virtio RNG"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0004",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "ff",
							"name": "unknown"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					}
				]
			}
		},
		"Kernel": {
			"Version": "Linux version 4.18.0-513.24.1.el8_9.x86_64 (mockbuild@iad1-prod-build001.bld.equ.rockylinux.org) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-20) (GCC)) #1 SMP Thu Apr 4 18:13:02 UTC 2024\n",
			"Modules": "tcp_diag 16384 0 - Live 0x0000000000000000\nbinfmt_misc 24576 1 - Live 0x0000000000000000\ninet_diag 24576 1 tcp_diag, Live 0x0000000000000000\nnfsv3 57344 1 - Live 0x0000000000000000\nrpcsec_gss_krb5 45056 0 - Live 0x0000000000000000\nnfsv4 917504 2 - Live 0x0000000000000000\ndns_resolver 16384 1 nfsv4, Live 0x0000000000000000\nnfs 425984 4 nfsv3,nfsv4, Live 0x0000000000000000\nfscache 389120 1 nfs, Live 0x0000000000000000\nintel_rapl_msr 16384 0 - Live 0x0000000000000000\nintel_rapl_common 24576 1 intel_rapl_msr, Live 0x0000000000000000\nintel_uncore_frequency_common 16384 0 - Live 0x0000000000000000\nisst_if_common 16384 0 - Live 0x0000000000000000\nnfit 65536 0 - Live 0x0000000000000000\nlibnvdimm 200704 1 nfit, Live 0x0000000000000000\ncrct10dif_pclmul 16384 1 - Live 0x0000000000000000\ncrc32_pclmul 16384 0 - Live 0x0000000000000000\nghash_clmulni_intel 16384 0 - Live 0x0000000000000000\ni2c_piix4 24576 0 - Live 0x0000000000000000\nrapl 20480 0 - Live 0x0000000000000000\npcspkr 16384 0 - Live 0x0000000000000000\nvfat 20480 1 - Live 0x0000000000000000\nfat 86016 1 vfat, Live 0x0000000000000000\nnfsd 548864 13 - Live 0x0000000000000000\nauth_rpcgss 139264 2 rpcsec_gss_krb5,nfsd, Live 0x0000000000000000\nnfs_acl 16384 2 nfsv3,nfsd, Live 0x0000000000000000\nlockd 126976 3 nfsv3,nfs,nfsd, Live 0x0000000000000000\ngrace 16384 2 nfsd,lockd, Live 0x0000000000000000\nsunrpc 585728 32 nfsv3,rpcsec_gss_krb5,nfsv4,nfs,nfsd,auth_rpcgss,nfs_acl,lockd, Live 0x0000000000000000\nxfs 1593344 1 - Live 0x0000000000000000\nlibcrc32c 16384 1 xfs, Live 0x0000000000000000\nsd_mod 57344 2 - Live 0x0000000000000000\nsg 40960 0 - Live 0x0000000000000000\nvirtio_net 61440 0 - Live 0x0000000000000000\ncrc32c_intel 24576 1 - Live 0x0000000000000000\nserio_raw 16384 0 - Live 0x0000000000000000\nnet_failover 24576 1 virtio_net, Live 0x0000000000000000\nvirtio_scsi 20480 2 - Live 0x0000000000000000\nfailover 16384 1 net_failover, Live 0x0000000000000000\nnvme 45056 0 - Live 0x0000000000000000\nnvme_core 139264 1 nvme, Live 0x0000000000000000\nt10_pi 16384 2 sd_mod,nvme_core, Live 0x0000000000000000\n",
			"Drivers": "Character devices:\n  1 mem\n  4 /dev/vc/0\n  4 tty\n  4 ttyS\n  5 /dev/tty\n  5 /dev/console\n  5 /dev/ptmx\n  7 vcs\n 10 misc\n 13 input\n 21 sg\n 29 fb\n128 ptm\n136 pts\n162 raw\n180 usb\n188 ttyUSB\n189 usb_device\n202 cpu/msr\n203 cpu/cpuid\n240 dimmctl\n241 ndctl\n242 nvme-generic\n243 nvme\n244 hidraw\n245 ttyDBC\n246 usbmon\n247 bsg\n248 watchdog\n249 ptp\n250 pps\n251 rtc\n252 dax\n253 tpm\n254 gpiochip\n\nBlock devices:\n  8 sd\n  9 md\n 65 sd\n 66 sd\n 67 sd\n 68 sd\n 69 sd\n 70 sd\n 71 sd\n128 sd\n129 sd\n130 sd\n131 sd\n132 sd\n133 sd\n134 sd\n135 sd\n254 mdp\n259 blkext\n"
		},
		"Stderr": "WARNING: \n/sys/class/drm does not exist on this system (likely the host system is a\nvirtual machine or container with no graphics). Therefore,\nGPUInfo.GraphicsCards will be an empty array.\nWARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied\nWARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied\nWARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied\nWARNING: Unable to read product_uuid: open /sys/class/dmi/id/product_uuid: permission denied\n",
		"Err": ""
	},
	{
		"Hostname": "hpcslurm-debugnodeset-3",
		"Info": {
			"memory": {
				"total_physical_bytes": 8589934592,
				"total_usable_bytes": 8056090624,
				"supported_page_sizes": [
					1073741824,
					2097152
				],
				"modules": null
			},
			"block": {
				"total_size_bytes": 53687091200,
				"disks": [
					{
						"name": "sda",
						"size_bytes": 53687091200,
						"physical_block_size_bytes": 4096,
						"drive_type": "hdd",
						"removable": false,
						"storage_controller": "scsi",
						"bus_path": "pci-0000:00:03.0-scsi-0:0:1:0",
						"vendor": "Google",
						"model": "PersistentDisk",
						"serial_number": "persistent-disk-0",
						"wwn": "unknown",
						"partitions": [
							{
								"name": "sda1",
								"label": "EFI\\x20System\\x20Partition",
								"mount_point": "/boot/efi",
								"size_bytes": 209715200,
								"type": "vfat",
								"read_only": false,
								"uuid": "a407d4b7-cfe4-4f7e-b9fc-ee7799ba3b84",
								"filesystem_label": "unknown"
							},
							{
								"name": "sda2",
								"label": "unknown",
								"mount_point": "/",
								"size_bytes": 53475328000,
								"type": "xfs",
								"read_only": false,
								"uuid": "144c8c6f-9c84-47c9-b637-8b7723fdb3ef",
								"filesystem_label": "root"
							}
						]
					}
				]
			},
			"cpu": {
				"total_cores": 1,
				"total_threads": 1,
				"processors": [
					{
						"id": 0,
						"total_cores": 1,
						"total_threads": 1,
						"vendor": "GenuineIntel",
						"model": "Intel(R) Xeon(R) CPU @ 2.80GHz",
						"capabilities": [
							"fpu",
							"vme",
							"de",
							"pse",
							"tsc",
							"msr",
							"pae",
							"mce",
							"cx8",
							"apic",
							"sep",
							"mtrr",
							"pge",
							"mca",
							"cmov",
							"pat",
							"pse36",
							"clflush",
							"mmx",
							"fxsr",
							"sse",
							"sse2",
							"ss",
							"ht",
							"syscall",
							"nx",
							"pdpe1gb",
							"rdtscp",
							"lm",
							"constant_tsc",
							"rep_good",
							"nopl",
							"xtopology",
							"nonstop_tsc",
							"cpuid",
							"tsc_known_freq",
							"pni",
							"pclmulqdq",
							"ssse3",
							"fma",
							"cx16",
							"pcid",
							"sse4_1",
							"sse4_2",
							"x2apic",
							"movbe",
							"popcnt",
							"aes",
							"xsave",
							"avx",
							"f16c",
							"rdrand",
							"hypervisor",
							"lahf_lm",
							"abm",
							"3dnowprefetch",
							"invpcid_single",
							"ssbd",
							"ibrs",
							"ibpb",
							"stibp",
							"ibrs_enhanced",
							"fsgsbase",
							"tsc_adjust",
							"bmi1",
							"hle",
							"avx2",
							"smep",
							"bmi2",
							"erms",
							"invpcid",
							"rtm",
							"avx512f",
							"avx512dq",
							"rdseed",
							"adx",
							"smap",
							"clflushopt",
							"clwb",
							"avx512cd",
							"avx512bw",
							"avx512vl",
							"xsaveopt",
							"xsavec",
							"xgetbv1",
							"xsaves",
							"arat",
							"avx512_vnni",
							"md_clear",
							"arch_capabilities"
						],
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						]
					}
				]
			},
			"topology": {
				"architecture": "smp",
				"nodes": [
					{
						"id": 0,
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						],
						"caches": [
							{
								"level": 1,
								"type": "instruction",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 1,
								"type": "data",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 2,
								"type": "unified",
								"size_bytes": 1048576,
								"logical_processors": [
									0
								]
							},
							{
								"level": 3,
								"type": "unified",
								"size_bytes": 34603008,
								"logical_processors": [
									0
								]
							}
						],
						"distances": [
							10
						],
						"memory": {
							"total_physical_bytes": 8589934592,
							"total_usable_bytes": 8056090624,
							"supported_page_sizes": [
								1073741824,
								2097152
							],
							"modules": null
						}
					}
				]
			},
			"network": {
				"nics": [
					{
						"name": "eth0",
						"mac_address": "42:01:0a:00:00:f7",
						"is_virtual": false,
						"capabilities": [
							{
								"name": "auto-negotiation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "pause-frame-use",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-checksumming",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-checksumming",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv4",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-ip-generic",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv6",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-fcoe-crc",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-sctp",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather-fraglist",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tcp-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-ecn-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tcp-mangleid-segmentation",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tx-tcp6-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-receive-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "large-receive-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "ntuple-filters",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "receive-hashing",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "highdma",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "rx-vlan-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "vlan-challenged",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-lockless",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "netns-local",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-robust",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-fcoe-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip4-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip6-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-partial",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tunnel-remcsum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-sctp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-esp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-list",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp-gro-forwarding",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "rx-gro-list",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tls-hw-rx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "fcoe-mtu",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-nocache-copy",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "loopback",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-fcs",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-all",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-stag-hw-insert",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-hw-parse",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "l2-fwd-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "hw-tc-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-tx-csum-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp_tunnel-port-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-tx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-gro-hw",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-record",
								"is_enabled": false,
								"can_enable": false
							}
						],
						"speed": "Unknown!",
						"duplex": "Unknown!(255)"
					}
				]
			},
			"gpu": {
				"cards": null
			},
			"chassis": {
				"asset_tag": "",
				"serial_number": "unknown",
				"type": "1",
				"type_description": "Other",
				"vendor": "Google",
				"version": ""
			},
			"bios": {
				"vendor": "Google",
				"version": "Google",
				"date": "06/07/2024"
			},
			"baseboard": {
				"asset_tag": "366C3950-A143-FA80-FD92-642B9326B452",
				"serial_number": "unknown",
				"vendor": "Google",
				"version": "",
				"product": "Google Compute Engine"
			},
			"product": {
				"family": "",
				"name": "Google Compute Engine",
				"vendor": "Google",
				"serial_number": "unknown",
				"uuid": "unknown",
				"sku": "",
				"version": ""
			},
			"pci": {
				"Devices": [
					{
						"driver": "",
						"address": "0000:00:00.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "1237",
							"name": "440FX - 82441FX PMC [Natoma]"
						},
						"revision": "0x02",
						"subsystem": {
							"id": "1100",
							"name": "Qemu virtual machine"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "00",
							"name": "Host bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7110",
							"name": "82371AB/EB/MB PIIX4 ISA"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "01",
							"name": "ISA bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.3",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7113",
							"name": "82371AB/EB/MB PIIX4 ACPI"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "80",
							"name": "Bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:03.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1004",
							"name": "Virtio SCSI"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0008",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "00",
							"name": "Non-VGA unclassified device"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:04.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1000",
							"name": "Virtio network device"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0001",
							"name": "unknown"
						},
						"class": {
							"id": "02",
							"name": "Network controller"
						},
						"subclass": {
							"id": "00",
							"name": "Ethernet controller"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:05.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1005",
							"name": "Virtio RNG"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0004",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "ff",
							"name": "unknown"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					}
				]
			}
		},
		"Kernel": {
			"Version": "Linux version 4.18.0-513.24.1.el8_9.x86_64 (mockbuild@iad1-prod-build001.bld.equ.rockylinux.org) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-20) (GCC)) #1 SMP Thu Apr 4 18:13:02 UTC 2024\n",
			"Modules": "tcp_diag 16384 0 - Live 0x0000000000000000\ninet_diag 24576 1 tcp_diag, Live 0x0000000000000000\nbinfmt_misc 24576 1 - Live 0x0000000000000000\nnfsv3 57344 1 - Live 0x0000000000000000\nrpcsec_gss_krb5 45056 0 - Live 0x0000000000000000\nnfsv4 917504 2 - Live 0x0000000000000000\ndns_resolver 16384 1 nfsv4, Live 0x0000000000000000\nnfs 425984 4 nfsv3,nfsv4, Live 0x0000000000000000\nfscache 389120 1 nfs, Live 0x0000000000000000\nintel_rapl_msr 16384 0 - Live 0x0000000000000000\nintel_rapl_common 24576 1 intel_rapl_msr, Live 0x0000000000000000\nintel_uncore_frequency_common 16384 0 - Live 0x0000000000000000\nisst_if_common 16384 0 - Live 0x0000000000000000\nnfit 65536 0 - Live 0x0000000000000000\nlibnvdimm 200704 1 nfit, Live 0x0000000000000000\ncrct10dif_pclmul 16384 1 - Live 0x0000000000000000\ncrc32_pclmul 16384 0 - Live 0x0000000000000000\nghash_clmulni_intel 16384 0 - Live 0x0000000000000000\ni2c_piix4 24576 0 - Live 0x0000000000000000\nvfat 20480 1 - Live 0x0000000000000000\nrapl 20480 0 - Live 0x0000000000000000\nfat 86016 1 vfat, Live 0x0000000000000000\npcspkr 16384 0 - Live 0x0000000000000000\nnfsd 548864 13 - Live 0x0000000000000000\nauth_rpcgss 139264 2 rpcsec_gss_krb5,nfsd, Live 0x0000000000000000\nnfs_acl 16384 2 nfsv3,nfsd, Live 0x0000000000000000\nlockd 126976 3 nfsv3,nfs,nfsd, Live 0x0000000000000000\ngrace 16384 2 nfsd,lockd, Live 0x0000000000000000\nsunrpc 585728 32 nfsv3,rpcsec_gss_krb5,nfsv4,nfs,nfsd,auth_rpcgss,nfs_acl,lockd, Live 0x0000000000000000\nxfs 1593344 1 - Live 0x0000000000000000\nlibcrc32c 16384 1 xfs, Live 0x0000000000000000\nsd_mod 57344 2 - Live 0x0000000000000000\nsg 40960 0 - Live 0x0000000000000000\nvirtio_net 61440 0 - Live 0x0000000000000000\ncrc32c_intel 24576 1 - Live 0x0000000000000000\nserio_raw 16384 0 - Live 0x0000000000000000\nnet_failover 24576 1 virtio_net, Live 0x0000000000000000\nfailover 16384 1 net_failover, Live 0x0000000000000000\nvirtio_scsi 20480 2 - Live 0x0000000000000000\nnvme 45056 0 - Live 0x0000000000000000\nnvme_core 139264 1 nvme, Live 0x0000000000000000\nt10_pi 16384 2 sd_mod,nvme_core, Live 0x0000000000000000\n",
			"Drivers": "Character devices:\n  1 mem\n  4 /dev/vc/0\n  4 tty\n  4 ttyS\n  5 /dev/tty\n  5 /dev/console\n  5 /dev/ptmx\n  7 vcs\n 10 misc\n 13 input\n 21 sg\n 29 fb\n128 ptm\n136 pts\n162 raw\n180 usb\n188 ttyUSB\n189 usb_device\n202 cpu/msr\n203 cpu/cpuid\n240 dimmctl\n241 ndctl\n242 nvme-generic\n243 nvme\n244 hidraw\n245 ttyDBC\n246 usbmon\n247 bsg\n248 watchdog\n249 ptp\n250 pps\n251 rtc\n252 dax\n253 tpm\n254 gpiochip\n\nBlock devices:\n  8 sd\n  9 md\n 65 sd\n 66 sd\n 67 sd\n 68 sd\n 69 sd\n 70 sd\n 71 sd\n128 sd\n129 sd\n130 sd\n131 sd\n132 sd\n133 sd\n134 sd\n135 sd\n254 mdp\n259 blkext\n"
		},
		"Stderr": "WARNING: \n/sys/class/drm does not exist on this system (likely the host system is a\nvirtual machine or container with no graphics). Therefore,\nGPUInfo.GraphicsCards will be an empty array.\nWARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied\nWARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied\nWARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied\nWARNING: Unable to read product_uuid: open /sys/class/dmi/id/product_uuid: permission denied\n",
		"Err": ""
	},
	{
		"Hostname": "hpcslurm-debugnodeset-15",
		"Info": {
			"memory": {
				"total_physical_bytes": 8589934592,
				"total_usable_bytes": 8056098816,
				"supported_page_sizes": [
					1073741824,
					2097152
				],
				"modules": null
			},
			"block": {
				"total_size_bytes": 53687091200,
				"disks": [
					{
						"name": "sda",
						"size_bytes": 53687091200,
						"physical_block_size_bytes": 4096,
						"drive_type": "hdd",
						"removable": false,
						"storage_controller": "scsi",
						"bus_path": "pci-0000:00:03.0-scsi-0:0:1:0",
						"vendor": "Google",
						"model": "PersistentDisk",
						"serial_number": "persistent-disk-0",
						"wwn": "unknown",
						"partitions": [
							{
								"name": "sda1",
								"label": "EFI\\x20System\\x20Partition",
								"mount_point": "/boot/efi",
								"size_bytes": 209715200,
								"type": "vfat",
								"read_only": false,
								"uuid": "a407d4b7-cfe4-4f7e-b9fc-ee7799ba3b84",
								"filesystem_label": "unknown"
							},
							{
								"name": "sda2",
								"label": "unknown",
								"mount_point": "/",
								"size_bytes": 53475328000,
								"type": "xfs",
								"read_only": false,
								"uuid": "144c8c6f-9c84-47c9-b637-8b7723fdb3ef",
								"filesystem_label": "root"
							}
						]
					}
				]
			},
			"cpu": {
				"total_cores": 1,
				"total_threads": 1,
				"processors": [
					{
						"id": 0,
						"total_cores": 1,
						"total_threads": 1,
						"vendor": "GenuineIntel",
						"model": "Intel(R) Xeon(R) CPU @ 2.80GHz",
						"capabilities": [
							"fpu",
							"vme",
							"de",
							"pse",
							"tsc",
							"msr",
							"pae",
							"mce",
							"cx8",
							"apic",
							"sep",
							"mtrr",
							"pge",
							"mca",
							"cmov",
							"pat",
							"pse36",
							"clflush",
							"mmx",
							"fxsr",
							"sse",
							"sse2",
							"ss",
							"ht",
							"syscall",
							"nx",
							"pdpe1gb",
							"rdtscp",
							"lm",
							"constant_tsc",
							"rep_good",
							"nopl",
							"xtopology",
							"nonstop_tsc",
							"cpuid",
							"tsc_known_freq",
							"pni",
							"pclmulqdq",
							"ssse3",
							"fma",
							"cx16",
							"pcid",
							"sse4_1",
							"sse4_2",
							"x2apic",
							"movbe",
							"popcnt",
							"aes",
							"xsave",
							"avx",
							"f16c",
							"rdrand",
							"hypervisor",
							"lahf_lm",
							"abm",
							"3dnowprefetch",
							"invpcid_single",
							"ssbd",
							"ibrs",
							"ibpb",
							"stibp",
							"ibrs_enhanced",
							"fsgsbase",
							"tsc_adjust",
							"bmi1",
							"hle",
							"avx2",
							"smep",
							"bmi2",
							"erms",
							"invpcid",
							"rtm",
							"avx512f",
							"avx512dq",
							"rdseed",
							"adx",
							"smap",
							"clflushopt",
							"clwb",
							"avx512cd",
							"avx512bw",
							"avx512vl",
							"xsaveopt",
							"xsavec",
							"xgetbv1",
							"xsaves",
							"arat",
							"avx512_vnni",
							"md_clear",
							"arch_capabilities"
						],
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						]
					}
				]
			},
			"topology": {
				"architecture": "smp",
				"nodes": [
					{
						"id": 0,
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						],
						"caches": [
							{
								"level": 1,
								"type": "instruction",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 1,
								"type": "data",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 2,
								"type": "unified",
								"size_bytes": 1048576,
								"logical_processors": [
									0
								]
							},
							{
								"level": 3,
								"type": "unified",
								"size_bytes": 34603008,
								"logical_processors": [
									0
								]
							}
						],
						"distances": [
							10
						],
						"memory": {
							"total_physical_bytes": 8589934592,
							"total_usable_bytes": 8056098816,
							"supported_page_sizes": [
								1073741824,
								2097152
							],
							"modules": null
						}
					}
				]
			},
			"network": {
				"nics": [
					{
						"name": "eth0",
						"mac_address": "42:01:0a:00:00:d6",
						"is_virtual": false,
						"capabilities": [
							{
								"name": "auto-negotiation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "pause-frame-use",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-checksumming",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-checksumming",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv4",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-ip-generic",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv6",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-fcoe-crc",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-sctp",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather-fraglist",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tcp-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-ecn-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tcp-mangleid-segmentation",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tx-tcp6-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-receive-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "large-receive-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "ntuple-filters",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "receive-hashing",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "highdma",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "rx-vlan-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "vlan-challenged",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-lockless",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "netns-local",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-robust",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-fcoe-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip4-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip6-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-partial",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tunnel-remcsum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-sctp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-esp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-list",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp-gro-forwarding",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "rx-gro-list",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tls-hw-rx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "fcoe-mtu",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-nocache-copy",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "loopback",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-fcs",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-all",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-stag-hw-insert",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-hw-parse",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "l2-fwd-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "hw-tc-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-tx-csum-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp_tunnel-port-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-tx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-gro-hw",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-record",
								"is_enabled": false,
								"can_enable": false
							}
						],
						"speed": "Unknown!",
						"duplex": "Unknown!(255)"
					}
				]
			},
			"gpu": {
				"cards": null
			},
			"chassis": {
				"asset_tag": "",
				"serial_number": "unknown",
				"type": "1",
				"type_description": "Other",
				"vendor": "Google",
				"version": ""
			},
			"bios": {
				"vendor": "Google",
				"version": "Google",
				"date": "04/02/2024"
			},
			"baseboard": {
				"asset_tag": "9A54A8E2-A740-F525-E58D-FF110F58A9AC",
				"serial_number": "unknown",
				"vendor": "Google",
				"version": "",
				"product": "Google Compute Engine"
			},
			"product": {
				"family": "",
				"name": "Google Compute Engine",
				"vendor": "Google",
				"serial_number": "unknown",
				"uuid": "unknown",
				"sku": "",
				"version": ""
			},
			"pci": {
				"Devices": [
					{
						"driver": "",
						"address": "0000:00:00.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "1237",
							"name": "440FX - 82441FX PMC [Natoma]"
						},
						"revision": "0x02",
						"subsystem": {
							"id": "1100",
							"name": "Qemu virtual machine"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "00",
							"name": "Host bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7110",
							"name": "82371AB/EB/MB PIIX4 ISA"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "01",
							"name": "ISA bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.3",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7113",
							"name": "82371AB/EB/MB PIIX4 ACPI"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "80",
							"name": "Bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:03.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1004",
							"name": "Virtio SCSI"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0008",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "00",
							"name": "Non-VGA unclassified device"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:04.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1000",
							"name": "Virtio network device"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0001",
							"name": "unknown"
						},
						"class": {
							"id": "02",
							"name": "Network controller"
						},
						"subclass": {
							"id": "00",
							"name": "Ethernet controller"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:05.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1005",
							"name": "Virtio RNG"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0004",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "ff",
							"name": "unknown"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					}
				]
			}
		},
		"Kernel": {
			"Version": "Linux version 4.18.0-513.24.1.el8_9.x86_64 (mockbuild@iad1-prod-build001.bld.equ.rockylinux.org) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-20) (GCC)) #1 SMP Thu Apr 4 18:13:02 UTC 2024\n",
			"Modules": "nfsv3 57344 1 - Live 0x0000000000000000\nrpcsec_gss_krb5 45056 0 - Live 0x0000000000000000\nnfsv4 917504 2 - Live 0x0000000000000000\ndns_resolver 16384 1 nfsv4, Live 0x0000000000000000\nnfs 425984 4 nfsv3,nfsv4, Live 0x0000000000000000\nfscache 389120 1 nfs, Live 0x0000000000000000\ntcp_diag 16384 0 - Live 0x0000000000000000\ninet_diag 24576 1 tcp_diag, Live 0x0000000000000000\nbinfmt_misc 24576 1 - Live 0x0000000000000000\nintel_rapl_msr 16384 0 - Live 0x0000000000000000\nintel_rapl_common 24576 1 intel_rapl_msr, Live 0x0000000000000000\nintel_uncore_frequency_common 16384 0 - Live 0x0000000000000000\nisst_if_common 16384 0 - Live 0x0000000000000000\nnfit 65536 0 - Live 0x0000000000000000\nlibnvdimm 200704 1 nfit, Live 0x0000000000000000\ncrct10dif_pclmul 16384 1 - Live 0x0000000000000000\ncrc32_pclmul 16384 0 - Live 0x0000000000000000\nghash_clmulni_intel 16384 0 - Live 0x0000000000000000\nrapl 20480 0 - Live 0x0000000000000000\nvfat 20480 1 - Live 0x0000000000000000\nfat 86016 1 vfat, Live 0x0000000000000000\npcspkr 16384 0 - Live 0x0000000000000000\ni2c_piix4 24576 0 - Live 0x0000000000000000\nnfsd 548864 13 - Live 0x0000000000000000\nauth_rpcgss 139264 2 rpcsec_gss_krb5,nfsd, Live 0x0000000000000000\nnfs_acl 16384 2 nfsv3,nfsd, Live 0x0000000000000000\nlockd 126976 3 nfsv3,nfs,nfsd, Live 0x0000000000000000\ngrace 16384 2 nfsd,lockd, Live 0x0000000000000000\nsunrpc 585728 32 nfsv3,rpcsec_gss_krb5,nfsv4,nfs,nfsd,auth_rpcgss,nfs_acl,lockd, Live 0x0000000000000000\nxfs 1593344 1 - Live 0x0000000000000000\nlibcrc32c 16384 1 xfs, Live 0x0000000000000000\nsd_mod 57344 2 - Live 0x0000000000000000\nsg 40960 0 - Live 0x0000000000000000\nvirtio_net 61440 0 - Live 0x0000000000000000\ncrc32c_intel 24576 1 - Live 0x0000000000000000\nnet_failover 24576 1 virtio_net, Live 0x0000000000000000\nserio_raw 16384 0 - Live 0x0000000000000000\nfailover 16384 1 net_failover, Live 0x0000000000000000\nvirtio_scsi 20480 2 - Live 0x0000000000000000\nnvme 45056 0 - Live 0x0000000000000000\nnvme_core 139264 1 nvme, Live 0x0000000000000000\nt10_pi 16384 2 sd_mod,nvme_core, Live 0x0000000000000000\n",
			"Drivers": "Character devices:\n  1 mem\n  4 /dev/vc/0\n  4 tty\n  4 ttyS\n  5 /dev/tty\n  5 /dev/console\n  5 /dev/ptmx\n  7 vcs\n 10 misc\n 13 input\n 21 sg\n 29 fb\n128 ptm\n136 pts\n162 raw\n180 usb\n188 ttyUSB\n189 usb_device\n202 cpu/msr\n203 cpu/cpuid\n240 dimmctl\n241 ndctl\n242 nvme-generic\n243 nvme\n244 hidraw\n245 ttyDBC\n246 usbmon\n247 bsg\n248 watchdog\n249 ptp\n250 pps\n251 rtc\n252 dax\n253 tpm\n254 gpiochip\n\nBlock devices:\n  8 sd\n  9 md\n 65 sd\n 66 sd\n 67 sd\n 68 sd\n 69 sd\n 70 sd\n 71 sd\n128 sd\n129 sd\n130 sd\n131 sd\n132 sd\n133 sd\n134 sd\n135 sd\n254 mdp\n259 blkext\n"
		},
		"Stderr": "WARNING: \n/sys/class/drm does not exist on this system (likely the host system is a\nvirtual machine or container with no graphics). Therefore,\nGPUInfo.GraphicsCards will be an empty array.\nWARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied\nWARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied\nWARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied\nWARNING: Unable to read product_uuid: open /sys/class/dmi/id/product_uuid: permission denied\n",
		"Err": ""
	},
	{
		"Hostname": "hpcslurm-debugnodeset-6",
		"Info": {
			"memory": {
				"total_physical_bytes": 8589934592,
				"total_usable_bytes": 8056090624,
				"supported_page_sizes": [
					1073741824,
					2097152
				],
				"modules": null
			},
			"block": {
				"total_size_bytes": 53687091200,
				"disks": [
					{
						"name": "sda",
						"size_bytes": 53687091200,
						"physical_block_size_bytes": 4096,
						"drive_type": "hdd",
						"removable": false,
						"storage_controller": "scsi",
						"bus_path": "pci-0000:00:03.0-scsi-0:0:1:0",
						"vendor": "Google",
						"model": "PersistentDisk",
						"serial_number": "persistent-disk-0",
						"wwn": "unknown",
						"partitions": [
							{
								"name": "sda1",
								"label": "EFI\\x20System\\x20Partition",
								"mount_point": "/boot/efi",
								"size_bytes": 209715200,
								"type": "vfat",
								"read_only": false,
								"uuid": "a407d4b7-cfe4-4f7e-b9fc-ee7799ba3b84",
								"filesystem_label": "unknown"
							},
							{
								"name": "sda2",
								"label": "unknown",
								"mount_point": "/",
								"size_bytes": 53475328000,
								"type": "xfs",
								"read_only": false,
								"uuid": "144c8c6f-9c84-47c9-b637-8b7723fdb3ef",
								"filesystem_label": "root"
							}
						]
					}
				]
			},
			"cpu": {
				"total_cores": 1,
				"total_threads": 1,
				"processors": [
					{
						"id": 0,
						"total_cores": 1,
						"total_threads": 1,
						"vendor": "GenuineIntel",
						"model": "Intel(R) Xeon(R) CPU @ 2.80GHz",
						"capabilities": [
							"fpu",
							"vme",
							"de",
							"pse",
							"tsc",
							"msr",
							"pae",
							"mce",
							"cx8",
							"apic",
							"sep",
							"mtrr",
							"pge",
							"mca",
							"cmov",
							"pat",
							"pse36",
							"clflush",
							"mmx",
							"fxsr",
							"sse",
							"sse2",
							"ss",
							"ht",
							"syscall",
							"nx",
							"pdpe1gb",
							"rdtscp",
							"lm",
							"constant_tsc",
							"rep_good",
							"nopl",
							"xtopology",
							"nonstop_tsc",
							"cpuid",
							"tsc_known_freq",
							"pni",
							"pclmulqdq",
							"ssse3",
							"fma",
							"cx16",
							"pcid",
							"sse4_1",
							"sse4_2",
							"x2apic",
							"movbe",
							"popcnt",
							"aes",
							"xsave",
							"avx",
							"f16c",
							"rdrand",
							"hypervisor",
							"lahf_lm",
							"abm",
							"3dnowprefetch",
							"invpcid_single",
							"ssbd",
							"ibrs",
							"ibpb",
							"stibp",
							"ibrs_enhanced",
							"fsgsbase",
							"tsc_adjust",
							"bmi1",
							"hle",
							"avx2",
							"smep",
							"bmi2",
							"erms",
							"invpcid",
							"rtm",
							"avx512f",
							"avx512dq",
							"rdseed",
							"adx",
							"smap",
							"clflushopt",
							"clwb",
							"avx512cd",
							"avx512bw",
							"avx512vl",
							"xsaveopt",
							"xsavec",
							"xgetbv1",
							"xsaves",
							"arat",
							"avx512_vnni",
							"md_clear",
							"arch_capabilities"
						],
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						]
					}
				]
			},
			"topology": {
				"architecture": "smp",
				"nodes": [
					{
						"id": 0,
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						],
						"caches": [
							{
								"level": 1,
								"type": "instruction",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 1,
								"type": "data",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 2,
								"type": "unified",
								"size_bytes": 1048576,
								"logical_processors": [
									0
								]
							},
							{
								"level": 3,
								"type": "unified",
								"size_bytes": 34603008,
								"logical_processors": [
									0
								]
							}
						],
						"distances": [
							10
						],
						"memory": {
							"total_physical_bytes": 8589934592,
							"total_usable_bytes": 8056090624,
							"supported_page_sizes": [
								1073741824,
								2097152
							],
							"modules": null
						}
					}
				]
			},
			"network": {
				"nics": [
					{
						"name": "eth0",
						"mac_address": "42:01:0a:00:00:e7",
						"is_virtual": false,
						"capabilities": [
							{
								"name": "auto-negotiation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "pause-frame-use",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-checksumming",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-checksumming",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv4",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-ip-generic",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv6",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-fcoe-crc",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-sctp",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather-fraglist",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tcp-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-ecn-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tcp-mangleid-segmentation",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tx-tcp6-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-receive-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "large-receive-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "ntuple-filters",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "receive-hashing",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "highdma",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "rx-vlan-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "vlan-challenged",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-lockless",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "netns-local",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-robust",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-fcoe-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip4-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip6-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-partial",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tunnel-remcsum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-sctp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-esp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-list",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp-gro-forwarding",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "rx-gro-list",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tls-hw-rx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "fcoe-mtu",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-nocache-copy",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "loopback",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-fcs",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-all",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-stag-hw-insert",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-hw-parse",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "l2-fwd-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "hw-tc-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-tx-csum-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp_tunnel-port-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-tx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-gro-hw",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-record",
								"is_enabled": false,
								"can_enable": false
							}
						],
						"speed": "Unknown!",
						"duplex": "Unknown!(255)"
					}
				]
			},
			"gpu": {
				"cards": null
			},
			"chassis": {
				"asset_tag": "",
				"serial_number": "unknown",
				"type": "1",
				"type_description": "Other",
				"vendor": "Google",
				"version": ""
			},
			"bios": {
				"vendor": "Google",
				"version": "Google",
				"date": "06/07/2024"
			},
			"baseboard": {
				"asset_tag": "0E9FB30B-90E8-5D79-4415-3623E007BFB8",
				"serial_number": "unknown",
				"vendor": "Google",
				"version": "",
				"product": "Google Compute Engine"
			},
			"product": {
				"family": "",
				"name": "Google Compute Engine",
				"vendor": "Google",
				"serial_number": "unknown",
				"uuid": "unknown",
				"sku": "",
				"version": ""
			},
			"pci": {
				"Devices": [
					{
						"driver": "",
						"address": "0000:00:00.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "1237",
							"name": "440FX - 82441FX PMC [Natoma]"
						},
						"revision": "0x02",
						"subsystem": {
							"id": "1100",
							"name": "Qemu virtual machine"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "00",
							"name": "Host bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7110",
							"name": "82371AB/EB/MB PIIX4 ISA"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "01",
							"name": "ISA bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.3",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7113",
							"name": "82371AB/EB/MB PIIX4 ACPI"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "80",
							"name": "Bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:03.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1004",
							"name": "Virtio SCSI"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0008",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "00",
							"name": "Non-VGA unclassified device"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:04.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1000",
							"name": "Virtio network device"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0001",
							"name": "unknown"
						},
						"class": {
							"id": "02",
							"name": "Network controller"
						},
						"subclass": {
							"id": "00",
							"name": "Ethernet controller"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:05.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1005",
							"name": "Virtio RNG"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0004",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "ff",
							"name": "unknown"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					}
				]
			}
		},
		"Kernel": {
			"Version": "Linux version 4.18.0-513.24.1.el8_9.x86_64 (mockbuild@iad1-prod-build001.bld.equ.rockylinux.org) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-20) (GCC)) #1 SMP Thu Apr 4 18:13:02 UTC 2024\n",
			"Modules": "tcp_diag 16384 0 - Live 0x0000000000000000\ninet_diag 24576 1 tcp_diag, Live 0x0000000000000000\nbinfmt_misc 24576 1 - Live 0x0000000000000000\nnfsv3 57344 1 - Live 0x0000000000000000\nrpcsec_gss_krb5 45056 0 - Live 0x0000000000000000\nnfsv4 917504 2 - Live 0x0000000000000000\ndns_resolver 16384 1 nfsv4, Live 0x0000000000000000\nnfs 425984 4 nfsv3,nfsv4, Live 0x0000000000000000\nfscache 389120 1 nfs, Live 0x0000000000000000\nintel_rapl_msr 16384 0 - Live 0x0000000000000000\nintel_rapl_common 24576 1 intel_rapl_msr, Live 0x0000000000000000\nintel_uncore_frequency_common 16384 0 - Live 0x0000000000000000\nisst_if_common 16384 0 - Live 0x0000000000000000\nnfit 65536 0 - Live 0x0000000000000000\nlibnvdimm 200704 1 nfit, Live 0x0000000000000000\ncrct10dif_pclmul 16384 1 - Live 0x0000000000000000\ncrc32_pclmul 16384 0 - Live 0x0000000000000000\nghash_clmulni_intel 16384 0 - Live 0x0000000000000000\nrapl 20480 0 - Live 0x0000000000000000\npcspkr 16384 0 - Live 0x0000000000000000\ni2c_piix4 24576 0 - Live 0x0000000000000000\nvfat 20480 1 - Live 0x0000000000000000\nfat 86016 1 vfat, Live 0x0000000000000000\nnfsd 548864 13 - Live 0x0000000000000000\nauth_rpcgss 139264 2 rpcsec_gss_krb5,nfsd, Live 0x0000000000000000\nnfs_acl 16384 2 nfsv3,nfsd, Live 0x0000000000000000\nlockd 126976 3 nfsv3,nfs,nfsd, Live 0x0000000000000000\ngrace 16384 2 nfsd,lockd, Live 0x0000000000000000\nsunrpc 585728 32 nfsv3,rpcsec_gss_krb5,nfsv4,nfs,nfsd,auth_rpcgss,nfs_acl,lockd, Live 0x0000000000000000\nxfs 1593344 1 - Live 0x0000000000000000\nlibcrc32c 16384 1 xfs, Live 0x0000000000000000\nsd_mod 57344 2 - Live 0x0000000000000000\nsg 40960 0 - Live 0x0000000000000000\nvirtio_net 61440 0 - Live 0x0000000000000000\ncrc32c_intel 24576 1 - Live 0x0000000000000000\nserio_raw 16384 0 - Live 0x0000000000000000\nnet_failover 24576 1 virtio_net, Live 0x0000000000000000\nfailover 16384 1 net_failover, Live 0x0000000000000000\nvirtio_scsi 20480 2 - Live 0x0000000000000000\nnvme 45056 0 - Live 0x0000000000000000\nnvme_core 139264 1 nvme, Live 0x0000000000000000\nt10_pi 16384 2 sd_mod,nvme_core, Live 0x0000000000000000\n",
			"Drivers": "Character devices:\n  1 mem\n  4 /dev/vc/0\n  4 tty\n  4 ttyS\n  5 /dev/tty\n  5 /dev/console\n  5 /dev/ptmx\n  7 vcs\n 10 misc\n 13 input\n 21 sg\n 29 fb\n128 ptm\n136 pts\n162 raw\n180 usb\n188 ttyUSB\n189 usb_device\n202 cpu/msr\n203 cpu/cpuid\n240 dimmctl\n241 ndctl\n242 nvme-generic\n243 nvme\n244 hidraw\n245 ttyDBC\n246 usbmon\n247 bsg\n248 watchdog\n249 ptp\n250 pps\n251 rtc\n252 dax\n253 tpm\n254 gpiochip\n\nBlock devices:\n  8 sd\n  9 md\n 65 sd\n 66 sd\n 67 sd\n 68 sd\n 69 sd\n 70 sd\n 71 sd\n128 sd\n129 sd\n130 sd\n131 sd\n132 sd\n133 sd\n134 sd\n135 sd\n254 mdp\n259 blkext\n"
		},
		"Stderr": "WARNING: \n/sys/class/drm does not exist on this system (likely the host system is a\nvirtual machine or container with no graphics). Therefore,\nGPUInfo.GraphicsCards will be an empty array.\nWARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied\nWARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied\nWARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied\nWARNING: Unable to read product_uuid: open /sys/class/dmi/id/product_uuid: permission denied\n",
		"Err": ""
	},
	{
		"Hostname": "hpcslurm-debugnodeset-12",
		"Info": {
			"memory": {
				"total_physical_bytes": 8589934592,
				"total_usable_bytes": 8056098816,
				"supported_page_sizes": [
					1073741824,
					2097152
				],
				"modules": null
			},
			"block": {
				"total_size_bytes": 53687091200,
				"disks": [
					{
						"name": "sda",
						"size_bytes": 53687091200,
						"physical_block_size_bytes": 4096,
						"drive_type": "hdd",
						"removable": false,
						"storage_controller": "scsi",
						"bus_path": "pci-0000:00:03.0-scsi-0:0:1:0",
						"vendor": "Google",
						"model": "PersistentDisk",
						"serial_number": "persistent-disk-0",
						"wwn": "unknown",
						"partitions": [
							{
								"name": "sda1",
								"label": "EFI\\x20System\\x20Partition",
								"mount_point": "/boot/efi",
								"size_bytes": 209715200,
								"type": "vfat",
								"read_only": false,
								"uuid": "a407d4b7-cfe4-4f7e-b9fc-ee7799ba3b84",
								"filesystem_label": "unknown"
							},
							{
								"name": "sda2",
								"label": "unknown",
								"mount_point": "/",
								"size_bytes": 53475328000,
								"type": "xfs",
								"read_only": false,
								"uuid": "144c8c6f-9c84-47c9-b637-8b7723fdb3ef",
								"filesystem_label": "root"
							}
						]
					}
				]
			},
			"cpu": {
				"total_cores": 1,
				"total_threads": 1,
				"processors": [
					{
						"id": 0,
						"total_cores": 1,
						"total_threads": 1,
						"vendor": "GenuineIntel",
						"model": "Intel(R) Xeon(R) CPU @ 2.80GHz",
						"capabilities": [
							"fpu",
							"vme",
							"de",
							"pse",
							"tsc",
							"msr",
							"pae",
							"mce",
							"cx8",
							"apic",
							"sep",
							"mtrr",
							"pge",
							"mca",
							"cmov",
							"pat",
							"pse36",
							"clflush",
							"mmx",
							"fxsr",
							"sse",
							"sse2",
							"ss",
							"ht",
							"syscall",
							"nx",
							"pdpe1gb",
							"rdtscp",
							"lm",
							"constant_tsc",
							"rep_good",
							"nopl",
							"xtopology",
							"nonstop_tsc",
							"cpuid",
							"tsc_known_freq",
							"pni",
							"pclmulqdq",
							"ssse3",
							"fma",
							"cx16",
							"pcid",
							"sse4_1",
							"sse4_2",
							"x2apic",
							"movbe",
							"popcnt",
							"aes",
							"xsave",
							"avx",
							"f16c",
							"rdrand",
							"hypervisor",
							"lahf_lm",
							"abm",
							"3dnowprefetch",
							"invpcid_single",
							"ssbd",
							"ibrs",
							"ibpb",
							"stibp",
							"ibrs_enhanced",
							"fsgsbase",
							"tsc_adjust",
							"bmi1",
							"hle",
							"avx2",
							"smep",
							"bmi2",
							"erms",
							"invpcid",
							"rtm",
							"avx512f",
							"avx512dq",
							"rdseed",
							"adx",
							"smap",
							"clflushopt",
							"clwb",
							"avx512cd",
							"avx512bw",
							"avx512vl",
							"xsaveopt",
							"xsavec",
							"xgetbv1",
							"xsaves",
							"arat",
							"avx512_vnni",
							"md_clear",
							"arch_capabilities"
						],
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						]
					}
				]
			},
			"topology": {
				"architecture": "smp",
				"nodes": [
					{
						"id": 0,
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						],
						"caches": [
							{
								"level": 1,
								"type": "instruction",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 1,
								"type": "data",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 2,
								"type": "unified",
								"size_bytes": 1048576,
								"logical_processors": [
									0
								]
							},
							{
								"level": 3,
								"type": "unified",
								"size_bytes": 34603008,
								"logical_processors": [
									0
								]
							}
						],
						"distances": [
							10
						],
						"memory": {
							"total_physical_bytes": 8589934592,
							"total_usable_bytes": 8056098816,
							"supported_page_sizes": [
								1073741824,
								2097152
							],
							"modules": null
						}
					}
				]
			},
			"network": {
				"nics": [
					{
						"name": "eth0",
						"mac_address": "42:01:0a:00:00:ee",
						"is_virtual": false,
						"capabilities": [
							{
								"name": "auto-negotiation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "pause-frame-use",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-checksumming",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-checksumming",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv4",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-ip-generic",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv6",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-fcoe-crc",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-sctp",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather-fraglist",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tcp-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-ecn-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tcp-mangleid-segmentation",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tx-tcp6-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-receive-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "large-receive-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "ntuple-filters",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "receive-hashing",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "highdma",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "rx-vlan-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "vlan-challenged",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-lockless",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "netns-local",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-robust",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-fcoe-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip4-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip6-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-partial",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tunnel-remcsum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-sctp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-esp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-list",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp-gro-forwarding",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "rx-gro-list",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tls-hw-rx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "fcoe-mtu",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-nocache-copy",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "loopback",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-fcs",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-all",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-stag-hw-insert",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-hw-parse",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "l2-fwd-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "hw-tc-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-tx-csum-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp_tunnel-port-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-tx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-gro-hw",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-record",
								"is_enabled": false,
								"can_enable": false
							}
						],
						"speed": "Unknown!",
						"duplex": "Unknown!(255)"
					}
				]
			},
			"gpu": {
				"cards": null
			},
			"chassis": {
				"asset_tag": "",
				"serial_number": "unknown",
				"type": "1",
				"type_description": "Other",
				"vendor": "Google",
				"version": ""
			},
			"bios": {
				"vendor": "Google",
				"version": "Google",
				"date": "06/27/2024"
			},
			"baseboard": {
				"asset_tag": "363B6AF1-333E-9974-2FA4-3ED2CAF80D42",
				"serial_number": "unknown",
				"vendor": "Google",
				"version": "",
				"product": "Google Compute Engine"
			},
			"product": {
				"family": "",
				"name": "Google Compute Engine",
				"vendor": "Google",
				"serial_number": "unknown",
				"uuid": "unknown",
				"sku": "",
				"version": ""
			},
			"pci": {
				"Devices": [
					{
						"driver": "",
						"address": "0000:00:00.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "1237",
							"name": "440FX - 82441FX PMC [Natoma]"
						},
						"revision": "0x02",
						"subsystem": {
							"id": "1100",
							"name": "Qemu virtual machine"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "00",
							"name": "Host bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7110",
							"name": "82371AB/EB/MB PIIX4 ISA"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "01",
							"name": "ISA bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.3",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7113",
							"name": "82371AB/EB/MB PIIX4 ACPI"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "80",
							"name": "Bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:03.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1004",
							"name": "Virtio SCSI"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0008",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "00",
							"name": "Non-VGA unclassified device"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:04.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1000",
							"name": "Virtio network device"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0001",
							"name": "unknown"
						},
						"class": {
							"id": "02",
							"name": "Network controller"
						},
						"subclass": {
							"id": "00",
							"name": "Ethernet controller"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:05.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1005",
							"name": "Virtio RNG"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0004",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "ff",
							"name": "unknown"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					}
				]
			}
		},
		"Kernel": {
			"Version": "Linux version 4.18.0-513.24.1.el8_9.x86_64 (mockbuild@iad1-prod-build001.bld.equ.rockylinux.org) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-20) (GCC)) #1 SMP Thu Apr 4 18:13:02 UTC 2024\n",
			"Modules": "tcp_diag 16384 0 - Live 0x0000000000000000\ninet_diag 24576 1 tcp_diag, Live 0x0000000000000000\nbinfmt_misc 24576 1 - Live 0x0000000000000000\nnfsv3 57344 1 - Live 0x0000000000000000\nrpcsec_gss_krb5 45056 0 - Live 0x0000000000000000\nnfsv4 917504 2 - Live 0x0000000000000000\ndns_resolver 16384 1 nfsv4, Live 0x0000000000000000\nnfs 425984 4 nfsv3,nfsv4, Live 0x0000000000000000\nfscache 389120 1 nfs, Live 0x0000000000000000\nintel_rapl_msr 16384 0 - Live 0x0000000000000000\nintel_rapl_common 24576 1 intel_rapl_msr, Live 0x0000000000000000\nintel_uncore_frequency_common 16384 0 - Live 0x0000000000000000\nisst_if_common 16384 0 - Live 0x0000000000000000\nnfit 65536 0 - Live 0x0000000000000000\nlibnvdimm 200704 1 nfit, Live 0x0000000000000000\ncrct10dif_pclmul 16384 1 - Live 0x0000000000000000\nvfat 20480 1 - Live 0x0000000000000000\nfat 86016 1 vfat, Live 0x0000000000000000\ncrc32_pclmul 16384 0 - Live 0x0000000000000000\nghash_clmulni_intel 16384 0 - Live 0x0000000000000000\nrapl 20480 0 - Live 0x0000000000000000\ni2c_piix4 24576 0 - Live 0x0000000000000000\npcspkr 16384 0 - Live 0x0000000000000000\nnfsd 548864 13 - Live 0x0000000000000000\nauth_rpcgss 139264 2 rpcsec_gss_krb5,nfsd, Live 0x0000000000000000\nnfs_acl 16384 2 nfsv3,nfsd, Live 0x0000000000000000\nlockd 126976 3 nfsv3,nfs,nfsd, Live 0x0000000000000000\ngrace 16384 2 nfsd,lockd, Live 0x0000000000000000\nsunrpc 585728 32 nfsv3,rpcsec_gss_krb5,nfsv4,nfs,nfsd,auth_rpcgss,nfs_acl,lockd, Live 0x0000000000000000\nxfs 1593344 1 - Live 0x0000000000000000\nlibcrc32c 16384 1 xfs, Live 0x0000000000000000\nsd_mod 57344 2 - Live 0x0000000000000000\nsg 40960 0 - Live 0x0000000000000000\nvirtio_net 61440 0 - Live 0x0000000000000000\ncrc32c_intel 24576 1 - Live 0x0000000000000000\nserio_raw 16384 0 - Live 0x0000000000000000\nnet_failover 24576 1 virtio_net, Live 0x0000000000000000\nvirtio_scsi 20480 2 - Live 0x0000000000000000\nfailover 16384 1 net_failover, Live 0x0000000000000000\nnvme 45056 0 - Live 0x0000000000000000\nnvme_core 139264 1 nvme, Live 0x0000000000000000\nt10_pi 16384 2 sd_mod,nvme_core, Live 0x0000000000000000\n",
			"Drivers": "Character devices:\n  1 mem\n  4 /dev/vc/0\n  4 tty\n  4 ttyS\n  5 /dev/tty\n  5 /dev/console\n  5 /dev/ptmx\n  7 vcs\n 10 misc\n 13 input\n 21 sg\n 29 fb\n128 ptm\n136 pts\n162 raw\n180 usb\n188 ttyUSB\n189 usb_device\n202 cpu/msr\n203 cpu/cpuid\n240 dimmctl\n241 ndctl\n242 nvme-generic\n243 nvme\n244 hidraw\n245 ttyDBC\n246 usbmon\n247 bsg\n248 watchdog\n249 ptp\n250 pps\n251 rtc\n252 dax\n253 tpm\n254 gpiochip\n\nBlock devices:\n  8 sd\n  9 md\n 65 sd\n 66 sd\n 67 sd\n 68 sd\n 69 sd\n 70 sd\n 71 sd\n128 sd\n129 sd\n130 sd\n131 sd\n132 sd\n133 sd\n134 sd\n135 sd\n254 mdp\n259 blkext\n"
		},
		"Stderr": "WARNING: \n/sys/class/drm does not exist on this system (likely the host system is a\nvirtual machine or container with no graphics). Therefore,\nGPUInfo.GraphicsCards will be an empty array.\nWARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied\nWARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied\nWARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied\nWARNING: Unable to read product_uuid: open /sys/class/dmi/id/product_uuid: permission denied\n",
		"Err": ""
	},
	{
		"Hostname": "hpcslurm-debugnodeset-14",
		"Info": {
			"memory": {
				"total_physical_bytes": 8589934592,
				"total_usable_bytes": 8056090624,
				"supported_page_sizes": [
					1073741824,
					2097152
				],
				"modules": null
			},
			"block": {
				"total_size_bytes": 53687091200,
				"disks": [
					{
						"name": "sda",
						"size_bytes": 53687091200,
						"physical_block_size_bytes": 4096,
						"drive_type": "hdd",
						"removable": false,
						"storage_controller": "scsi",
						"bus_path": "pci-0000:00:03.0-scsi-0:0:1:0",
						"vendor": "Google",
						"model": "PersistentDisk",
						"serial_number": "persistent-disk-0",
						"wwn": "unknown",
						"partitions": [
							{
								"name": "sda1",
								"label": "EFI\\x20System\\x20Partition",
								"mount_point": "/boot/efi",
								"size_bytes": 209715200,
								"type": "vfat",
								"read_only": false,
								"uuid": "a407d4b7-cfe4-4f7e-b9fc-ee7799ba3b84",
								"filesystem_label": "unknown"
							},
							{
								"name": "sda2",
								"label": "unknown",
								"mount_point": "/",
								"size_bytes": 53475328000,
								"type": "xfs",
								"read_only": false,
								"uuid": "144c8c6f-9c84-47c9-b637-8b7723fdb3ef",
								"filesystem_label": "root"
							}
						]
					}
				]
			},
			"cpu": {
				"total_cores": 1,
				"total_threads": 1,
				"processors": [
					{
						"id": 0,
						"total_cores": 1,
						"total_threads": 1,
						"vendor": "GenuineIntel",
						"model": "Intel(R) Xeon(R) CPU @ 2.80GHz",
						"capabilities": [
							"fpu",
							"vme",
							"de",
							"pse",
							"tsc",
							"msr",
							"pae",
							"mce",
							"cx8",
							"apic",
							"sep",
							"mtrr",
							"pge",
							"mca",
							"cmov",
							"pat",
							"pse36",
							"clflush",
							"mmx",
							"fxsr",
							"sse",
							"sse2",
							"ss",
							"ht",
							"syscall",
							"nx",
							"pdpe1gb",
							"rdtscp",
							"lm",
							"constant_tsc",
							"rep_good",
							"nopl",
							"xtopology",
							"nonstop_tsc",
							"cpuid",
							"tsc_known_freq",
							"pni",
							"pclmulqdq",
							"ssse3",
							"fma",
							"cx16",
							"pcid",
							"sse4_1",
							"sse4_2",
							"x2apic",
							"movbe",
							"popcnt",
							"aes",
							"xsave",
							"avx",
							"f16c",
							"rdrand",
							"hypervisor",
							"lahf_lm",
							"abm",
							"3dnowprefetch",
							"invpcid_single",
							"ssbd",
							"ibrs",
							"ibpb",
							"stibp",
							"ibrs_enhanced",
							"fsgsbase",
							"tsc_adjust",
							"bmi1",
							"hle",
							"avx2",
							"smep",
							"bmi2",
							"erms",
							"invpcid",
							"rtm",
							"avx512f",
							"avx512dq",
							"rdseed",
							"adx",
							"smap",
							"clflushopt",
							"clwb",
							"avx512cd",
							"avx512bw",
							"avx512vl",
							"xsaveopt",
							"xsavec",
							"xgetbv1",
							"xsaves",
							"arat",
							"avx512_vnni",
							"md_clear",
							"arch_capabilities"
						],
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						]
					}
				]
			},
			"topology": {
				"architecture": "smp",
				"nodes": [
					{
						"id": 0,
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						],
						"caches": [
							{
								"level": 1,
								"type": "instruction",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 1,
								"type": "data",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 2,
								"type": "unified",
								"size_bytes": 1048576,
								"logical_processors": [
									0
								]
							},
							{
								"level": 3,
								"type": "unified",
								"size_bytes": 34603008,
								"logical_processors": [
									0
								]
							}
						],
						"distances": [
							10
						],
						"memory": {
							"total_physical_bytes": 8589934592,
							"total_usable_bytes": 8056090624,
							"supported_page_sizes": [
								1073741824,
								2097152
							],
							"modules": null
						}
					}
				]
			},
			"network": {
				"nics": [
					{
						"name": "eth0",
						"mac_address": "42:01:0a:00:00:d3",
						"is_virtual": false,
						"capabilities": [
							{
								"name": "auto-negotiation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "pause-frame-use",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-checksumming",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-checksumming",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv4",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-ip-generic",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv6",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-fcoe-crc",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-sctp",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather-fraglist",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tcp-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-ecn-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tcp-mangleid-segmentation",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tx-tcp6-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-receive-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "large-receive-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "ntuple-filters",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "receive-hashing",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "highdma",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "rx-vlan-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "vlan-challenged",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-lockless",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "netns-local",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-robust",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-fcoe-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip4-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip6-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-partial",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tunnel-remcsum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-sctp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-esp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-list",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp-gro-forwarding",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "rx-gro-list",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tls-hw-rx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "fcoe-mtu",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-nocache-copy",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "loopback",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-fcs",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-all",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-stag-hw-insert",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-hw-parse",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "l2-fwd-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "hw-tc-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-tx-csum-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp_tunnel-port-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-tx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-gro-hw",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-record",
								"is_enabled": false,
								"can_enable": false
							}
						],
						"speed": "Unknown!",
						"duplex": "Unknown!(255)"
					}
				]
			},
			"gpu": {
				"cards": null
			},
			"chassis": {
				"asset_tag": "",
				"serial_number": "unknown",
				"type": "1",
				"type_description": "Other",
				"vendor": "Google",
				"version": ""
			},
			"bios": {
				"vendor": "Google",
				"version": "Google",
				"date": "04/02/2024"
			},
			"baseboard": {
				"asset_tag": "E00C2954-1973-DAEF-E954-BAB0DF14097C",
				"serial_number": "unknown",
				"vendor": "Google",
				"version": "",
				"product": "Google Compute Engine"
			},
			"product": {
				"family": "",
				"name": "Google Compute Engine",
				"vendor": "Google",
				"serial_number": "unknown",
				"uuid": "unknown",
				"sku": "",
				"version": ""
			},
			"pci": {
				"Devices": [
					{
						"driver": "",
						"address": "0000:00:00.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "1237",
							"name": "440FX - 82441FX PMC [Natoma]"
						},
						"revision": "0x02",
						"subsystem": {
							"id": "1100",
							"name": "Qemu virtual machine"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "00",
							"name": "Host bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7110",
							"name": "82371AB/EB/MB PIIX4 ISA"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "01",
							"name": "ISA bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.3",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7113",
							"name": "82371AB/EB/MB PIIX4 ACPI"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "80",
							"name": "Bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:03.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1004",
							"name": "Virtio SCSI"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0008",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "00",
							"name": "Non-VGA unclassified device"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:04.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1000",
							"name": "Virtio network device"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0001",
							"name": "unknown"
						},
						"class": {
							"id": "02",
							"name": "Network controller"
						},
						"subclass": {
							"id": "00",
							"name": "Ethernet controller"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:05.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1005",
							"name": "Virtio RNG"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0004",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "ff",
							"name": "unknown"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					}
				]
			}
		},
		"Kernel": {
			"Version": "Linux version 4.18.0-513.24.1.el8_9.x86_64 (mockbuild@iad1-prod-build001.bld.equ.rockylinux.org) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-20) (GCC)) #1 SMP Thu Apr 4 18:13:02 UTC 2024\n",
			"Modules": "nfsv3 57344 1 - Live 0x0000000000000000\nrpcsec_gss_krb5 45056 0 - Live 0x0000000000000000\nnfsv4 917504 2 - Live 0x0000000000000000\ndns_resolver 16384 1 nfsv4, Live 0x0000000000000000\nnfs 425984 4 nfsv3,nfsv4, Live 0x0000000000000000\nfscache 389120 1 nfs, Live 0x0000000000000000\ntcp_diag 16384 0 - Live 0x0000000000000000\ninet_diag 24576 1 tcp_diag, Live 0x0000000000000000\nbinfmt_misc 24576 1 - Live 0x0000000000000000\nintel_rapl_msr 16384 0 - Live 0x0000000000000000\nintel_rapl_common 24576 1 intel_rapl_msr, Live 0x0000000000000000\nintel_uncore_frequency_common 16384 0 - Live 0x0000000000000000\nisst_if_common 16384 0 - Live 0x0000000000000000\nnfit 65536 0 - Live 0x0000000000000000\nlibnvdimm 200704 1 nfit, Live 0x0000000000000000\ncrct10dif_pclmul 16384 1 - Live 0x0000000000000000\ncrc32_pclmul 16384 0 - Live 0x0000000000000000\nghash_clmulni_intel 16384 0 - Live 0x0000000000000000\nrapl 20480 0 - Live 0x0000000000000000\nvfat 20480 1 - Live 0x0000000000000000\nfat 86016 1 vfat, Live 0x0000000000000000\npcspkr 16384 0 - Live 0x0000000000000000\ni2c_piix4 24576 0 - Live 0x0000000000000000\nnfsd 548864 13 - Live 0x0000000000000000\nauth_rpcgss 139264 2 rpcsec_gss_krb5,nfsd, Live 0x0000000000000000\nnfs_acl 16384 2 nfsv3,nfsd, Live 0x0000000000000000\nlockd 126976 3 nfsv3,nfs,nfsd, Live 0x0000000000000000\ngrace 16384 2 nfsd,lockd, Live 0x0000000000000000\nsunrpc 585728 32 nfsv3,rpcsec_gss_krb5,nfsv4,nfs,nfsd,auth_rpcgss,nfs_acl,lockd, Live 0x0000000000000000\nxfs 1593344 1 - Live 0x0000000000000000\nlibcrc32c 16384 1 xfs, Live 0x0000000000000000\nsd_mod 57344 2 - Live 0x0000000000000000\nsg 40960 0 - Live 0x0000000000000000\nvirtio_net 61440 0 - Live 0x0000000000000000\ncrc32c_intel 24576 1 - Live 0x0000000000000000\nserio_raw 16384 0 - Live 0x0000000000000000\nnet_failover 24576 1 virtio_net, Live 0x0000000000000000\nfailover 16384 1 net_failover, Live 0x0000000000000000\nvirtio_scsi 20480 2 - Live 0x0000000000000000\nnvme 45056 0 - Live 0x0000000000000000\nnvme_core 139264 1 nvme, Live 0x0000000000000000\nt10_pi 16384 2 sd_mod,nvme_core, Live 0x0000000000000000\n",
			"Drivers": "Character devices:\n  1 mem\n  4 /dev/vc/0\n  4 tty\n  4 ttyS\n  5 /dev/tty\n  5 /dev/console\n  5 /dev/ptmx\n  7 vcs\n 10 misc\n 13 input\n 21 sg\n 29 fb\n128 ptm\n136 pts\n162 raw\n180 usb\n188 ttyUSB\n189 usb_device\n202 cpu/msr\n203 cpu/cpuid\n240 dimmctl\n241 ndctl\n242 nvme-generic\n243 nvme\n244 hidraw\n245 ttyDBC\n246 usbmon\n247 bsg\n248 watchdog\n249 ptp\n250 pps\n251 rtc\n252 dax\n253 tpm\n254 gpiochip\n\nBlock devices:\n  8 sd\n  9 md\n 65 sd\n 66 sd\n 67 sd\n 68 sd\n 69 sd\n 70 sd\n 71 sd\n128 sd\n129 sd\n130 sd\n131 sd\n132 sd\n133 sd\n134 sd\n135 sd\n254 mdp\n259 blkext\n"
		},
		"Stderr": "WARNING: \n/sys/class/drm does not exist on this system (likely the host system is a\nvirtual machine or container with no graphics). Therefore,\nGPUInfo.GraphicsCards will be an empty array.\nWARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied\nWARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied\nWARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied\nWARNING: Unable to read product_uuid: open /sys/class/dmi/id/product_uuid: permission denied\n",
		"Err": ""
	},
	{
		"Hostname": "hpcslurm-debugnodeset-7",
		"Info": {
			"memory": {
				"total_physical_bytes": 8589934592,
				"total_usable_bytes": 8056098816,
				"supported_page_sizes": [
					1073741824,
					2097152
				],
				"modules": null
			},
			"block": {
				"total_size_bytes": 53687091200,
				"disks": [
					{
						"name": "sda",
						"size_bytes": 53687091200,
						"physical_block_size_bytes": 4096,
						"drive_type": "hdd",
						"removable": false,
						"storage_controller": "scsi",
						"bus_path": "pci-0000:00:03.0-scsi-0:0:1:0",
						"vendor": "Google",
						"model": "PersistentDisk",
						"serial_number": "persistent-disk-0",
						"wwn": "unknown",
						"partitions": [
							{
								"name": "sda1",
								"label": "EFI\\x20System\\x20Partition",
								"mount_point": "/boot/efi",
								"size_bytes": 209715200,
								"type": "vfat",
								"read_only": false,
								"uuid": "a407d4b7-cfe4-4f7e-b9fc-ee7799ba3b84",
								"filesystem_label": "unknown"
							},
							{
								"name": "sda2",
								"label": "unknown",
								"mount_point": "/",
								"size_bytes": 53475328000,
								"type": "xfs",
								"read_only": false,
								"uuid": "144c8c6f-9c84-47c9-b637-8b7723fdb3ef",
								"filesystem_label": "root"
							}
						]
					}
				]
			},
			"cpu": {
				"total_cores": 1,
				"total_threads": 1,
				"processors": [
					{
						"id": 0,
						"total_cores": 1,
						"total_threads": 1,
						"vendor": "GenuineIntel",
						"model": "Intel(R) Xeon(R) CPU @ 2.80GHz",
						"capabilities": [
							"fpu",
							"vme",
							"de",
							"pse",
							"tsc",
							"msr",
							"pae",
							"mce",
							"cx8",
							"apic",
							"sep",
							"mtrr",
							"pge",
							"mca",
							"cmov",
							"pat",
							"pse36",
							"clflush",
							"mmx",
							"fxsr",
							"sse",
							"sse2",
							"ss",
							"ht",
							"syscall",
							"nx",
							"pdpe1gb",
							"rdtscp",
							"lm",
							"constant_tsc",
							"rep_good",
							"nopl",
							"xtopology",
							"nonstop_tsc",
							"cpuid",
							"tsc_known_freq",
							"pni",
							"pclmulqdq",
							"ssse3",
							"fma",
							"cx16",
							"pcid",
							"sse4_1",
							"sse4_2",
							"x2apic",
							"movbe",
							"popcnt",
							"aes",
							"xsave",
							"avx",
							"f16c",
							"rdrand",
							"hypervisor",
							"lahf_lm",
							"abm",
							"3dnowprefetch",
							"invpcid_single",
							"ssbd",
							"ibrs",
							"ibpb",
							"stibp",
							"ibrs_enhanced",
							"fsgsbase",
							"tsc_adjust",
							"bmi1",
							"hle",
							"avx2",
							"smep",
							"bmi2",
							"erms",
							"invpcid",
							"rtm",
							"avx512f",
							"avx512dq",
							"rdseed",
							"adx",
							"smap",
							"clflushopt",
							"clwb",
							"avx512cd",
							"avx512bw",
							"avx512vl",
							"xsaveopt",
							"xsavec",
							"xgetbv1",
							"xsaves",
							"arat",
							"avx512_vnni",
							"md_clear",
							"arch_capabilities"
						],
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						]
					}
				]
			},
			"topology": {
				"architecture": "smp",
				"nodes": [
					{
						"id": 0,
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						],
						"caches": [
							{
								"level": 1,
								"type": "instruction",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 1,
								"type": "data",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 2,
								"type": "unified",
								"size_bytes": 1048576,
								"logical_processors": [
									0
								]
							},
							{
								"level": 3,
								"type": "unified",
								"size_bytes": 34603008,
								"logical_processors": [
									0
								]
							}
						],
						"distances": [
							10
						],
						"memory": {
							"total_physical_bytes": 8589934592,
							"total_usable_bytes": 8056098816,
							"supported_page_sizes": [
								1073741824,
								2097152
							],
							"modules": null
						}
					}
				]
			},
			"network": {
				"nics": [
					{
						"name": "eth0",
						"mac_address": "42:01:0a:00:00:cc",
						"is_virtual": false,
						"capabilities": [
							{
								"name": "auto-negotiation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "pause-frame-use",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-checksumming",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-checksumming",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv4",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-ip-generic",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv6",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-fcoe-crc",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-sctp",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather-fraglist",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tcp-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-ecn-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tcp-mangleid-segmentation",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tx-tcp6-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-receive-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "large-receive-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "ntuple-filters",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "receive-hashing",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "highdma",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "rx-vlan-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "vlan-challenged",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-lockless",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "netns-local",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-robust",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-fcoe-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip4-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip6-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-partial",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tunnel-remcsum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-sctp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-esp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-list",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp-gro-forwarding",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "rx-gro-list",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tls-hw-rx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "fcoe-mtu",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-nocache-copy",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "loopback",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-fcs",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-all",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-stag-hw-insert",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-hw-parse",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "l2-fwd-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "hw-tc-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-tx-csum-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp_tunnel-port-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-tx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-gro-hw",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-record",
								"is_enabled": false,
								"can_enable": false
							}
						],
						"speed": "Unknown!",
						"duplex": "Unknown!(255)"
					}
				]
			},
			"gpu": {
				"cards": null
			},
			"chassis": {
				"asset_tag": "",
				"serial_number": "unknown",
				"type": "1",
				"type_description": "Other",
				"vendor": "Google",
				"version": ""
			},
			"bios": {
				"vendor": "Google",
				"version": "Google",
				"date": "06/27/2024"
			},
			"baseboard": {
				"asset_tag": "8D2AB171-2610-5008-05F0-313250D123BF",
				"serial_number": "unknown",
				"vendor": "Google",
				"version": "",
				"product": "Google Compute Engine"
			},
			"product": {
				"family": "",
				"name": "Google Compute Engine",
				"vendor": "Google",
				"serial_number": "unknown",
				"uuid": "unknown",
				"sku": "",
				"version": ""
			},
			"pci": {
				"Devices": [
					{
						"driver": "",
						"address": "0000:00:00.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "1237",
							"name": "440FX - 82441FX PMC [Natoma]"
						},
						"revision": "0x02",
						"subsystem": {
							"id": "1100",
							"name": "Qemu virtual machine"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "00",
							"name": "Host bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7110",
							"name": "82371AB/EB/MB PIIX4 ISA"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "01",
							"name": "ISA bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.3",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7113",
							"name": "82371AB/EB/MB PIIX4 ACPI"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "80",
							"name": "Bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:03.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1004",
							"name": "Virtio SCSI"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0008",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "00",
							"name": "Non-VGA unclassified device"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:04.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1000",
							"name": "Virtio network device"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0001",
							"name": "unknown"
						},
						"class": {
							"id": "02",
							"name": "Network controller"
						},
						"subclass": {
							"id": "00",
							"name": "Ethernet controller"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:05.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1005",
							"name": "Virtio RNG"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0004",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "ff",
							"name": "unknown"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					}
				]
			}
		},
		"Kernel": {
			"Version": "Linux version 4.18.0-513.24.1.el8_9.x86_64 (mockbuild@iad1-prod-build001.bld.equ.rockylinux.org) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-20) (GCC)) #1 SMP Thu Apr 4 18:13:02 UTC 2024\n",
			"Modules": "tcp_diag 16384 0 - Live 0x0000000000000000\ninet_diag 24576 1 tcp_diag, Live 0x0000000000000000\nbinfmt_misc 24576 1 - Live 0x0000000000000000\nnfsv3 57344 1 - Live 0x0000000000000000\nrpcsec_gss_krb5 45056 0 - Live 0x0000000000000000\nnfsv4 917504 2 - Live 0x0000000000000000\ndns_resolver 16384 1 nfsv4, Live 0x0000000000000000\nnfs 425984 4 nfsv3,nfsv4, Live 0x0000000000000000\nfscache 389120 1 nfs, Live 0x0000000000000000\nintel_rapl_msr 16384 0 - Live 0x0000000000000000\nintel_rapl_common 24576 1 intel_rapl_msr, Live 0x0000000000000000\nintel_uncore_frequency_common 16384 0 - Live 0x0000000000000000\nisst_if_common 16384 0 - Live 0x0000000000000000\nnfit 65536 0 - Live 0x0000000000000000\nlibnvdimm 200704 1 nfit, Live 0x0000000000000000\ncrct10dif_pclmul 16384 1 - Live 0x0000000000000000\ncrc32_pclmul 16384 0 - Live 0x0000000000000000\nghash_clmulni_intel 16384 0 - Live 0x0000000000000000\nrapl 20480 0 - Live 0x0000000000000000\ni2c_piix4 24576 0 - Live 0x0000000000000000\nvfat 20480 1 - Live 0x0000000000000000\nfat 86016 1 vfat, Live 0x0000000000000000\npcspkr 16384 0 - Live 0x0000000000000000\nnfsd 548864 13 - Live 0x0000000000000000\nauth_rpcgss 139264 2 rpcsec_gss_krb5,nfsd, Live 0x0000000000000000\nnfs_acl 16384 2 nfsv3,nfsd, Live 0x0000000000000000\nlockd 126976 3 nfsv3,nfs,nfsd, Live 0x0000000000000000\ngrace 16384 2 nfsd,lockd, Live 0x0000000000000000\nsunrpc 585728 32 nfsv3,rpcsec_gss_krb5,nfsv4,nfs,nfsd,auth_rpcgss,nfs_acl,lockd, Live 0x0000000000000000\nxfs 1593344 1 - Live 0x0000000000000000\nlibcrc32c 16384 1 xfs, Live 0x0000000000000000\nsd_mod 57344 2 - Live 0x0000000000000000\nsg 40960 0 - Live 0x0000000000000000\nvirtio_net 61440 0 - Live 0x0000000000000000\ncrc32c_intel 24576 1 - Live 0x0000000000000000\nserio_raw 16384 0 - Live 0x0000000000000000\nnet_failover 24576 1 virtio_net, Live 0x0000000000000000\nfailover 16384 1 net_failover, Live 0x0000000000000000\nvirtio_scsi 20480 2 - Live 0x0000000000000000\nnvme 45056 0 - Live 0x0000000000000000\nnvme_core 139264 1 nvme, Live 0x0000000000000000\nt10_pi 16384 2 sd_mod,nvme_core, Live 0x0000000000000000\n",
			"Drivers": "Character devices:\n  1 mem\n  4 /dev/vc/0\n  4 tty\n  4 ttyS\n  5 /dev/tty\n  5 /dev/console\n  5 /dev/ptmx\n  7 vcs\n 10 misc\n 13 input\n 21 sg\n 29 fb\n128 ptm\n136 pts\n162 raw\n180 usb\n188 ttyUSB\n189 usb_device\n202 cpu/msr\n203 cpu/cpuid\n240 dimmctl\n241 ndctl\n242 nvme-generic\n243 nvme\n244 hidraw\n245 ttyDBC\n246 usbmon\n247 bsg\n248 watchdog\n249 ptp\n250 pps\n251 rtc\n252 dax\n253 tpm\n254 gpiochip\n\nBlock devices:\n  8 sd\n  9 md\n 65 sd\n 66 sd\n 67 sd\n 68 sd\n 69 sd\n 70 sd\n 71 sd\n128 sd\n129 sd\n130 sd\n131 sd\n132 sd\n133 sd\n134 sd\n135 sd\n254 mdp\n259 blkext\n"
		},
		"Stderr": "WARNING: \n/sys/class/drm does not exist on this system (likely the host system is a\nvirtual machine or container with no graphics). Therefore,\nGPUInfo.GraphicsCards will be an empty array.\nWARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied\nWARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied\nWARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied\nWARNING: Unable to read product_uuid: open /sys/class/dmi/id/product_uuid: permission denied\n",
		"Err": ""
	},
	{
		"Hostname": "hpcslurm-debugnodeset-10",
		"Info": {
			"memory": {
				"total_physical_bytes": 8589934592,
				"total_usable_bytes": 8056098816,
				"supported_page_sizes": [
					1073741824,
					2097152
				],
				"modules": null
			},
			"block": {
				"total_size_bytes": 53687091200,
				"disks": [
					{
						"name": "sda",
						"size_bytes": 53687091200,
						"physical_block_size_bytes": 4096,
						"drive_type": "hdd",
						"removable": false,
						"storage_controller": "scsi",
						"bus_path": "pci-0000:00:03.0-scsi-0:0:1:0",
						"vendor": "Google",
						"model": "PersistentDisk",
						"serial_number": "persistent-disk-0",
						"wwn": "unknown",
						"partitions": [
							{
								"name": "sda1",
								"label": "EFI\\x20System\\x20Partition",
								"mount_point": "/boot/efi",
								"size_bytes": 209715200,
								"type": "vfat",
								"read_only": false,
								"uuid": "a407d4b7-cfe4-4f7e-b9fc-ee7799ba3b84",
								"filesystem_label": "unknown"
							},
							{
								"name": "sda2",
								"label": "unknown",
								"mount_point": "/",
								"size_bytes": 53475328000,
								"type": "xfs",
								"read_only": false,
								"uuid": "144c8c6f-9c84-47c9-b637-8b7723fdb3ef",
								"filesystem_label": "root"
							}
						]
					}
				]
			},
			"cpu": {
				"total_cores": 1,
				"total_threads": 1,
				"processors": [
					{
						"id": 0,
						"total_cores": 1,
						"total_threads": 1,
						"vendor": "GenuineIntel",
						"model": "Intel(R) Xeon(R) CPU @ 2.80GHz",
						"capabilities": [
							"fpu",
							"vme",
							"de",
							"pse",
							"tsc",
							"msr",
							"pae",
							"mce",
							"cx8",
							"apic",
							"sep",
							"mtrr",
							"pge",
							"mca",
							"cmov",
							"pat",
							"pse36",
							"clflush",
							"mmx",
							"fxsr",
							"sse",
							"sse2",
							"ss",
							"ht",
							"syscall",
							"nx",
							"pdpe1gb",
							"rdtscp",
							"lm",
							"constant_tsc",
							"rep_good",
							"nopl",
							"xtopology",
							"nonstop_tsc",
							"cpuid",
							"tsc_known_freq",
							"pni",
							"pclmulqdq",
							"ssse3",
							"fma",
							"cx16",
							"pcid",
							"sse4_1",
							"sse4_2",
							"x2apic",
							"movbe",
							"popcnt",
							"aes",
							"xsave",
							"avx",
							"f16c",
							"rdrand",
							"hypervisor",
							"lahf_lm",
							"abm",
							"3dnowprefetch",
							"invpcid_single",
							"ssbd",
							"ibrs",
							"ibpb",
							"stibp",
							"ibrs_enhanced",
							"fsgsbase",
							"tsc_adjust",
							"bmi1",
							"hle",
							"avx2",
							"smep",
							"bmi2",
							"erms",
							"invpcid",
							"rtm",
							"avx512f",
							"avx512dq",
							"rdseed",
							"adx",
							"smap",
							"clflushopt",
							"clwb",
							"avx512cd",
							"avx512bw",
							"avx512vl",
							"xsaveopt",
							"xsavec",
							"xgetbv1",
							"xsaves",
							"arat",
							"avx512_vnni",
							"md_clear",
							"arch_capabilities"
						],
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						]
					}
				]
			},
			"topology": {
				"architecture": "smp",
				"nodes": [
					{
						"id": 0,
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						],
						"caches": [
							{
								"level": 1,
								"type": "instruction",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 1,
								"type": "data",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 2,
								"type": "unified",
								"size_bytes": 1048576,
								"logical_processors": [
									0
								]
							},
							{
								"level": 3,
								"type": "unified",
								"size_bytes": 34603008,
								"logical_processors": [
									0
								]
							}
						],
						"distances": [
							10
						],
						"memory": {
							"total_physical_bytes": 8589934592,
							"total_usable_bytes": 8056098816,
							"supported_page_sizes": [
								1073741824,
								2097152
							],
							"modules": null
						}
					}
				]
			},
			"network": {
				"nics": [
					{
						"name": "eth0",
						"mac_address": "42:01:0a:00:00:f5",
						"is_virtual": false,
						"capabilities": [
							{
								"name": "auto-negotiation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "pause-frame-use",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-checksumming",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-checksumming",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv4",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-ip-generic",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv6",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-fcoe-crc",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-sctp",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather-fraglist",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tcp-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-ecn-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tcp-mangleid-segmentation",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tx-tcp6-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-receive-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "large-receive-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "ntuple-filters",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "receive-hashing",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "highdma",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "rx-vlan-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "vlan-challenged",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-lockless",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "netns-local",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-robust",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-fcoe-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip4-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip6-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-partial",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tunnel-remcsum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-sctp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-esp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-list",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp-gro-forwarding",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "rx-gro-list",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tls-hw-rx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "fcoe-mtu",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-nocache-copy",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "loopback",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-fcs",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-all",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-stag-hw-insert",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-hw-parse",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "l2-fwd-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "hw-tc-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-tx-csum-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp_tunnel-port-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-tx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-gro-hw",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-record",
								"is_enabled": false,
								"can_enable": false
							}
						],
						"speed": "Unknown!",
						"duplex": "Unknown!(255)"
					}
				]
			},
			"gpu": {
				"cards": null
			},
			"chassis": {
				"asset_tag": "",
				"serial_number": "unknown",
				"type": "1",
				"type_description": "Other",
				"vendor": "Google",
				"version": ""
			},
			"bios": {
				"vendor": "Google",
				"version": "Google",
				"date": "06/07/2024"
			},
			"baseboard": {
				"asset_tag": "88231395-864A-4399-08F7-726557DA4AB0",
				"serial_number": "unknown",
				"vendor": "Google",
				"version": "",
				"product": "Google Compute Engine"
			},
			"product": {
				"family": "",
				"name": "Google Compute Engine",
				"vendor": "Google",
				"serial_number": "unknown",
				"uuid": "unknown",
				"sku": "",
				"version": ""
			},
			"pci": {
				"Devices": [
					{
						"driver": "",
						"address": "0000:00:00.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "1237",
							"name": "440FX - 82441FX PMC [Natoma]"
						},
						"revision": "0x02",
						"subsystem": {
							"id": "1100",
							"name": "Qemu virtual machine"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "00",
							"name": "Host bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7110",
							"name": "82371AB/EB/MB PIIX4 ISA"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "01",
							"name": "ISA bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.3",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7113",
							"name": "82371AB/EB/MB PIIX4 ACPI"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "80",
							"name": "Bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:03.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1004",
							"name": "Virtio SCSI"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0008",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "00",
							"name": "Non-VGA unclassified device"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:04.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1000",
							"name": "Virtio network device"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0001",
							"name": "unknown"
						},
						"class": {
							"id": "02",
							"name": "Network controller"
						},
						"subclass": {
							"id": "00",
							"name": "Ethernet controller"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:05.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1005",
							"name": "Virtio RNG"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0004",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "ff",
							"name": "unknown"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					}
				]
			}
		},
		"Kernel": {
			"Version": "Linux version 4.18.0-513.24.1.el8_9.x86_64 (mockbuild@iad1-prod-build001.bld.equ.rockylinux.org) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-20) (GCC)) #1 SMP Thu Apr 4 18:13:02 UTC 2024\n",
			"Modules": "tcp_diag 16384 0 - Live 0x0000000000000000\nbinfmt_misc 24576 1 - Live 0x0000000000000000\ninet_diag 24576 1 tcp_diag, Live 0x0000000000000000\nnfsv3 57344 1 - Live 0x0000000000000000\nrpcsec_gss_krb5 45056 0 - Live 0x0000000000000000\nnfsv4 917504 2 - Live 0x0000000000000000\ndns_resolver 16384 1 nfsv4, Live 0x0000000000000000\nnfs 425984 4 nfsv3,nfsv4, Live 0x0000000000000000\nfscache 389120 1 nfs, Live 0x0000000000000000\nintel_rapl_msr 16384 0 - Live 0x0000000000000000\nintel_rapl_common 24576 1 intel_rapl_msr, Live 0x0000000000000000\nintel_uncore_frequency_common 16384 0 - Live 0x0000000000000000\nisst_if_common 16384 0 - Live 0x0000000000000000\nnfit 65536 0 - Live 0x0000000000000000\nlibnvdimm 200704 1 nfit, Live 0x0000000000000000\ncrct10dif_pclmul 16384 1 - Live 0x0000000000000000\ncrc32_pclmul 16384 0 - Live 0x0000000000000000\nghash_clmulni_intel 16384 0 - Live 0x0000000000000000\nrapl 20480 0 - Live 0x0000000000000000\ni2c_piix4 24576 0 - Live 0x0000000000000000\npcspkr 16384 0 - Live 0x0000000000000000\nvfat 20480 1 - Live 0x0000000000000000\nfat 86016 1 vfat, Live 0x0000000000000000\nnfsd 548864 13 - Live 0x0000000000000000\nauth_rpcgss 139264 2 rpcsec_gss_krb5,nfsd, Live 0x0000000000000000\nnfs_acl 16384 2 nfsv3,nfsd, Live 0x0000000000000000\nlockd 126976 3 nfsv3,nfs,nfsd, Live 0x0000000000000000\ngrace 16384 2 nfsd,lockd, Live 0x0000000000000000\nsunrpc 585728 32 nfsv3,rpcsec_gss_krb5,nfsv4,nfs,nfsd,auth_rpcgss,nfs_acl,lockd, Live 0x0000000000000000\nxfs 1593344 1 - Live 0x0000000000000000\nlibcrc32c 16384 1 xfs, Live 0x0000000000000000\nsd_mod 57344 2 - Live 0x0000000000000000\nsg 40960 0 - Live 0x0000000000000000\nvirtio_net 61440 0 - Live 0x0000000000000000\ncrc32c_intel 24576 1 - Live 0x0000000000000000\nserio_raw 16384 0 - Live 0x0000000000000000\nnet_failover 24576 1 virtio_net, Live 0x0000000000000000\nvirtio_scsi 20480 2 - Live 0x0000000000000000\nfailover 16384 1 net_failover, Live 0x0000000000000000\nnvme 45056 0 - Live 0x0000000000000000\nnvme_core 139264 1 nvme, Live 0x0000000000000000\nt10_pi 16384 2 sd_mod,nvme_core, Live 0x0000000000000000\n",
			"Drivers": "Character devices:\n  1 mem\n  4 /dev/vc/0\n  4 tty\n  4 ttyS\n  5 /dev/tty\n  5 /dev/console\n  5 /dev/ptmx\n  7 vcs\n 10 misc\n 13 input\n 21 sg\n 29 fb\n128 ptm\n136 pts\n162 raw\n180 usb\n188 ttyUSB\n189 usb_device\n202 cpu/msr\n203 cpu/cpuid\n240 dimmctl\n241 ndctl\n242 nvme-generic\n243 nvme\n244 hidraw\n245 ttyDBC\n246 usbmon\n247 bsg\n248 watchdog\n249 ptp\n250 pps\n251 rtc\n252 dax\n253 tpm\n254 gpiochip\n\nBlock devices:\n  8 sd\n  9 md\n 65 sd\n 66 sd\n 67 sd\n 68 sd\n 69 sd\n 70 sd\n 71 sd\n128 sd\n129 sd\n130 sd\n131 sd\n132 sd\n133 sd\n134 sd\n135 sd\n254 mdp\n259 blkext\n"
		},
		"Stderr": "WARNING: \n/sys/class/drm does not exist on this system (likely the host system is a\nvirtual machine or container with no graphics). Therefore,\nGPUInfo.GraphicsCards will be an empty array.\nWARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied\nWARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied\nWARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied\nWARNING: Unable to read product_uuid: open /sys/class/dmi/id/product_uuid: permission denied\n",
		"Err": ""
	},
	{
		"Hostname": "hpcslurm-debugnodeset-2",
		"Info": {
			"memory": {
				"total_physical_bytes": 8589934592,
				"total_usable_bytes": 8056098816,
				"supported_page_sizes": [
					1073741824,
					2097152
				],
				"modules": null
			},
			"block": {
				"total_size_bytes": 53687091200,
				"disks": [
					{
						"name": "sda",
						"size_bytes": 53687091200,
						"physical_block_size_bytes": 4096,
						"drive_type": "hdd",
						"removable": false,
						"storage_controller": "scsi",
						"bus_path": "pci-0000:00:03.0-scsi-0:0:1:0",
						"vendor": "Google",
						"model": "PersistentDisk",
						"serial_number": "persistent-disk-0",
						"wwn": "unknown",
						"partitions": [
							{
								"name": "sda1",
								"label": "EFI\\x20System\\x20Partition",
								"mount_point": "/boot/efi",
								"size_bytes": 209715200,
								"type": "vfat",
								"read_only": false,
								"uuid": "a407d4b7-cfe4-4f7e-b9fc-ee7799ba3b84",
								"filesystem_label": "unknown"
							},
							{
								"name": "sda2",
								"label": "unknown",
								"mount_point": "/",
								"size_bytes": 53475328000,
								"type": "xfs",
								"read_only": false,
								"uuid": "144c8c6f-9c84-47c9-b637-8b7723fdb3ef",
								"filesystem_label": "root"
							}
						]
					}
				]
			},
			"cpu": {
				"total_cores": 1,
				"total_threads": 1,
				"processors": [
					{
						"id": 0,
						"total_cores": 1,
						"total_threads": 1,
						"vendor": "GenuineIntel",
						"model": "Intel(R) Xeon(R) CPU @ 2.80GHz",
						"capabilities": [
							"fpu",
							"vme",
							"de",
							"pse",
							"tsc",
							"msr",
							"pae",
							"mce",
							"cx8",
							"apic",
							"sep",
							"mtrr",
							"pge",
							"mca",
							"cmov",
							"pat",
							"pse36",
							"clflush",
							"mmx",
							"fxsr",
							"sse",
							"sse2",
							"ss",
							"ht",
							"syscall",
							"nx",
							"pdpe1gb",
							"rdtscp",
							"lm",
							"constant_tsc",
							"rep_good",
							"nopl",
							"xtopology",
							"nonstop_tsc",
							"cpuid",
							"tsc_known_freq",
							"pni",
							"pclmulqdq",
							"ssse3",
							"fma",
							"cx16",
							"pcid",
							"sse4_1",
							"sse4_2",
							"x2apic",
							"movbe",
							"popcnt",
							"aes",
							"xsave",
							"avx",
							"f16c",
							"rdrand",
							"hypervisor",
							"lahf_lm",
							"abm",
							"3dnowprefetch",
							"invpcid_single",
							"ssbd",
							"ibrs",
							"ibpb",
							"stibp",
							"ibrs_enhanced",
							"fsgsbase",
							"tsc_adjust",
							"bmi1",
							"hle",
							"avx2",
							"smep",
							"bmi2",
							"erms",
							"invpcid",
							"rtm",
							"avx512f",
							"avx512dq",
							"rdseed",
							"adx",
							"smap",
							"clflushopt",
							"clwb",
							"avx512cd",
							"avx512bw",
							"avx512vl",
							"xsaveopt",
							"xsavec",
							"xgetbv1",
							"xsaves",
							"arat",
							"avx512_vnni",
							"md_clear",
							"arch_capabilities"
						],
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						]
					}
				]
			},
			"topology": {
				"architecture": "smp",
				"nodes": [
					{
						"id": 0,
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						],
						"caches": [
							{
								"level": 1,
								"type": "instruction",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 1,
								"type": "data",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 2,
								"type": "unified",
								"size_bytes": 1048576,
								"logical_processors": [
									0
								]
							},
							{
								"level": 3,
								"type": "unified",
								"size_bytes": 34603008,
								"logical_processors": [
									0
								]
							}
						],
						"distances": [
							10
						],
						"memory": {
							"total_physical_bytes": 8589934592,
							"total_usable_bytes": 8056098816,
							"supported_page_sizes": [
								1073741824,
								2097152
							],
							"modules": null
						}
					}
				]
			},
			"network": {
				"nics": [
					{
						"name": "eth0",
						"mac_address": "42:01:0a:00:00:f2",
						"is_virtual": false,
						"capabilities": [
							{
								"name": "auto-negotiation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "pause-frame-use",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-checksumming",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-checksumming",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv4",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-ip-generic",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv6",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-fcoe-crc",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-sctp",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather-fraglist",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tcp-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-ecn-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tcp-mangleid-segmentation",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tx-tcp6-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-receive-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "large-receive-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "ntuple-filters",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "receive-hashing",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "highdma",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "rx-vlan-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "vlan-challenged",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-lockless",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "netns-local",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-robust",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-fcoe-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip4-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip6-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-partial",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tunnel-remcsum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-sctp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-esp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-list",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp-gro-forwarding",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "rx-gro-list",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tls-hw-rx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "fcoe-mtu",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-nocache-copy",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "loopback",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-fcs",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-all",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-stag-hw-insert",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-hw-parse",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "l2-fwd-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "hw-tc-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-tx-csum-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp_tunnel-port-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-tx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-gro-hw",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-record",
								"is_enabled": false,
								"can_enable": false
							}
						],
						"speed": "Unknown!",
						"duplex": "Unknown!(255)"
					}
				]
			},
			"gpu": {
				"cards": null
			},
			"chassis": {
				"asset_tag": "",
				"serial_number": "unknown",
				"type": "1",
				"type_description": "Other",
				"vendor": "Google",
				"version": ""
			},
			"bios": {
				"vendor": "Google",
				"version": "Google",
				"date": "06/27/2024"
			},
			"baseboard": {
				"asset_tag": "8D129F9A-B9EE-300F-86CB-58B9CCE59FB1",
				"serial_number": "unknown",
				"vendor": "Google",
				"version": "",
				"product": "Google Compute Engine"
			},
			"product": {
				"family": "",
				"name": "Google Compute Engine",
				"vendor": "Google",
				"serial_number": "unknown",
				"uuid": "unknown",
				"sku": "",
				"version": ""
			},
			"pci": {
				"Devices": [
					{
						"driver": "",
						"address": "0000:00:00.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "1237",
							"name": "440FX - 82441FX PMC [Natoma]"
						},
						"revision": "0x02",
						"subsystem": {
							"id": "1100",
							"name": "Qemu virtual machine"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "00",
							"name": "Host bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7110",
							"name": "82371AB/EB/MB PIIX4 ISA"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "01",
							"name": "ISA bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.3",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7113",
							"name": "82371AB/EB/MB PIIX4 ACPI"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "80",
							"name": "Bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:03.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1004",
							"name": "Virtio SCSI"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0008",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "00",
							"name": "Non-VGA unclassified device"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:04.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1000",
							"name": "Virtio network device"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0001",
							"name": "unknown"
						},
						"class": {
							"id": "02",
							"name": "Network controller"
						},
						"subclass": {
							"id": "00",
							"name": "Ethernet controller"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:05.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1005",
							"name": "Virtio RNG"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0004",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "ff",
							"name": "unknown"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					}
				]
			}
		},
		"Kernel": {
			"Version": "Linux version 4.18.0-513.24.1.el8_9.x86_64 (mockbuild@iad1-prod-build001.bld.equ.rockylinux.org) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-20) (GCC)) #1 SMP Thu Apr 4 18:13:02 UTC 2024\n",
			"Modules": "tcp_diag 16384 0 - Live 0x0000000000000000\ninet_diag 24576 1 tcp_diag, Live 0x0000000000000000\nbinfmt_misc 24576 1 - Live 0x0000000000000000\nnfsv3 57344 1 - Live 0x0000000000000000\nrpcsec_gss_krb5 45056 0 - Live 0x0000000000000000\nnfsv4 917504 2 - Live 0x0000000000000000\ndns_resolver 16384 1 nfsv4, Live 0x0000000000000000\nnfs 425984 4 nfsv3,nfsv4, Live 0x0000000000000000\nfscache 389120 1 nfs, Live 0x0000000000000000\nintel_rapl_msr 16384 0 - Live 0x0000000000000000\nintel_rapl_common 24576 1 intel_rapl_msr, Live 0x0000000000000000\nintel_uncore_frequency_common 16384 0 - Live 0x0000000000000000\nisst_if_common 16384 0 - Live 0x0000000000000000\nnfit 65536 0 - Live 0x0000000000000000\nlibnvdimm 200704 1 nfit, Live 0x0000000000000000\ncrct10dif_pclmul 16384 1 - Live 0x0000000000000000\ncrc32_pclmul 16384 0 - Live 0x0000000000000000\nghash_clmulni_intel 16384 0 - Live 0x0000000000000000\ni2c_piix4 24576 0 - Live 0x0000000000000000\nvfat 20480 1 - Live 0x0000000000000000\nfat 86016 1 vfat, Live 0x0000000000000000\nrapl 20480 0 - Live 0x0000000000000000\npcspkr 16384 0 - Live 0x0000000000000000\nnfsd 548864 13 - Live 0x0000000000000000\nauth_rpcgss 139264 2 rpcsec_gss_krb5,nfsd, Live 0x0000000000000000\nnfs_acl 16384 2 nfsv3,nfsd, Live 0x0000000000000000\nlockd 126976 3 nfsv3,nfs,nfsd, Live 0x0000000000000000\ngrace 16384 2 nfsd,lockd, Live 0x0000000000000000\nsunrpc 585728 32 nfsv3,rpcsec_gss_krb5,nfsv4,nfs,nfsd,auth_rpcgss,nfs_acl,lockd, Live 0x0000000000000000\nxfs 1593344 1 - Live 0x0000000000000000\nlibcrc32c 16384 1 xfs, Live 0x0000000000000000\nsd_mod 57344 2 - Live 0x0000000000000000\nsg 40960 0 - Live 0x0000000000000000\nvirtio_net 61440 0 - Live 0x0000000000000000\ncrc32c_intel 24576 1 - Live 0x0000000000000000\nserio_raw 16384 0 - Live 0x0000000000000000\nnet_failover 24576 1 virtio_net, Live 0x0000000000000000\nvirtio_scsi 20480 2 - Live 0x0000000000000000\nfailover 16384 1 net_failover, Live 0x0000000000000000\nnvme 45056 0 - Live 0x0000000000000000\nnvme_core 139264 1 nvme, Live 0x0000000000000000\nt10_pi 16384 2 sd_mod,nvme_core, Live 0x0000000000000000\n",
			"Drivers": "Character devices:\n  1 mem\n  4 /dev/vc/0\n  4 tty\n  4 ttyS\n  5 /dev/tty\n  5 /dev/console\n  5 /dev/ptmx\n  7 vcs\n 10 misc\n 13 input\n 21 sg\n 29 fb\n128 ptm\n136 pts\n162 raw\n180 usb\n188 ttyUSB\n189 usb_device\n202 cpu/msr\n203 cpu/cpuid\n240 dimmctl\n241 ndctl\n242 nvme-generic\n243 nvme\n244 hidraw\n245 ttyDBC\n246 usbmon\n247 bsg\n248 watchdog\n249 ptp\n250 pps\n251 rtc\n252 dax\n253 tpm\n254 gpiochip\n\nBlock devices:\n  8 sd\n  9 md\n 65 sd\n 66 sd\n 67 sd\n 68 sd\n 69 sd\n 70 sd\n 71 sd\n128 sd\n129 sd\n130 sd\n131 sd\n132 sd\n133 sd\n134 sd\n135 sd\n254 mdp\n259 blkext\n"
		},
		"Stderr": "WARNING: \n/sys/class/drm does not exist on this system (likely the host system is a\nvirtual machine or container with no graphics). Therefore,\nGPUInfo.GraphicsCards will be an empty array.\nWARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied\nWARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied\nWARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied\nWARNING: Unable to read product_uuid: open /sys/class/dmi/id/product_uuid: permission denied\n",
		"Err": ""
	},
	{
		"Hostname": "hpcslurm-debugnodeset-4",
		"Info": {
			"memory": {
				"total_physical_bytes": 8589934592,
				"total_usable_bytes": 8056090624,
				"supported_page_sizes": [
					1073741824,
					2097152
				],
				"modules": null
			},
			"block": {
				"total_size_bytes": 53687091200,
				"disks": [
					{
						"name": "sda",
						"size_bytes": 53687091200,
						"physical_block_size_bytes": 4096,
						"drive_type": "hdd",
						"removable": false,
						"storage_controller": "scsi",
						"bus_path": "pci-0000:00:03.0-scsi-0:0:1:0",
						"vendor": "Google",
						"model": "PersistentDisk",
						"serial_number": "persistent-disk-0",
						"wwn": "unknown",
						"partitions": [
							{
								"name": "sda1",
								"label": "EFI\\x20System\\x20Partition",
								"mount_point": "/boot/efi",
								"size_bytes": 209715200,
								"type": "vfat",
								"read_only": false,
								"uuid": "a407d4b7-cfe4-4f7e-b9fc-ee7799ba3b84",
								"filesystem_label": "unknown"
							},
							{
								"name": "sda2",
								"label": "unknown",
								"mount_point": "/",
								"size_bytes": 53475328000,
								"type": "xfs",
								"read_only": false,
								"uuid": "144c8c6f-9c84-47c9-b637-8b7723fdb3ef",
								"filesystem_label": "root"
							}
						]
					}
				]
			},
			"cpu": {
				"total_cores": 1,
				"total_threads": 1,
				"processors": [
					{
						"id": 0,
						"total_cores": 1,
						"total_threads": 1,
						"vendor": "GenuineIntel",
						"model": "Intel(R) Xeon(R) CPU @ 2.80GHz",
						"capabilities": [
							"fpu",
							"vme",
							"de",
							"pse",
							"tsc",
							"msr",
							"pae",
							"mce",
							"cx8",
							"apic",
							"sep",
							"mtrr",
							"pge",
							"mca",
							"cmov",
							"pat",
							"pse36",
							"clflush",
							"mmx",
							"fxsr",
							"sse",
							"sse2",
							"ss",
							"ht",
							"syscall",
							"nx",
							"pdpe1gb",
							"rdtscp",
							"lm",
							"constant_tsc",
							"rep_good",
							"nopl",
							"xtopology",
							"nonstop_tsc",
							"cpuid",
							"tsc_known_freq",
							"pni",
							"pclmulqdq",
							"ssse3",
							"fma",
							"cx16",
							"pcid",
							"sse4_1",
							"sse4_2",
							"x2apic",
							"movbe",
							"popcnt",
							"aes",
							"xsave",
							"avx",
							"f16c",
							"rdrand",
							"hypervisor",
							"lahf_lm",
							"abm",
							"3dnowprefetch",
							"invpcid_single",
							"ssbd",
							"ibrs",
							"ibpb",
							"stibp",
							"ibrs_enhanced",
							"fsgsbase",
							"tsc_adjust",
							"bmi1",
							"hle",
							"avx2",
							"smep",
							"bmi2",
							"erms",
							"invpcid",
							"rtm",
							"avx512f",
							"avx512dq",
							"rdseed",
							"adx",
							"smap",
							"clflushopt",
							"clwb",
							"avx512cd",
							"avx512bw",
							"avx512vl",
							"xsaveopt",
							"xsavec",
							"xgetbv1",
							"xsaves",
							"arat",
							"avx512_vnni",
							"md_clear",
							"arch_capabilities"
						],
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						]
					}
				]
			},
			"topology": {
				"architecture": "smp",
				"nodes": [
					{
						"id": 0,
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						],
						"caches": [
							{
								"level": 1,
								"type": "instruction",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 1,
								"type": "data",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 2,
								"type": "unified",
								"size_bytes": 1048576,
								"logical_processors": [
									0
								]
							},
							{
								"level": 3,
								"type": "unified",
								"size_bytes": 34603008,
								"logical_processors": [
									0
								]
							}
						],
						"distances": [
							10
						],
						"memory": {
							"total_physical_bytes": 8589934592,
							"total_usable_bytes": 8056090624,
							"supported_page_sizes": [
								1073741824,
								2097152
							],
							"modules": null
						}
					}
				]
			},
			"network": {
				"nics": [
					{
						"name": "eth0",
						"mac_address": "42:01:0a:00:00:cb",
						"is_virtual": false,
						"capabilities": [
							{
								"name": "auto-negotiation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "pause-frame-use",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-checksumming",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-checksumming",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv4",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-ip-generic",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv6",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-fcoe-crc",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-sctp",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather-fraglist",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tcp-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-ecn-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tcp-mangleid-segmentation",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tx-tcp6-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-receive-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "large-receive-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "ntuple-filters",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "receive-hashing",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "highdma",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "rx-vlan-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "vlan-challenged",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-lockless",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "netns-local",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-robust",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-fcoe-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip4-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip6-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-partial",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tunnel-remcsum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-sctp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-esp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-list",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp-gro-forwarding",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "rx-gro-list",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tls-hw-rx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "fcoe-mtu",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-nocache-copy",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "loopback",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-fcs",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-all",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-stag-hw-insert",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-hw-parse",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "l2-fwd-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "hw-tc-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-tx-csum-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp_tunnel-port-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-tx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-gro-hw",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-record",
								"is_enabled": false,
								"can_enable": false
							}
						],
						"speed": "Unknown!",
						"duplex": "Unknown!(255)"
					}
				]
			},
			"gpu": {
				"cards": null
			},
			"chassis": {
				"asset_tag": "",
				"serial_number": "unknown",
				"type": "1",
				"type_description": "Other",
				"vendor": "Google",
				"version": ""
			},
			"bios": {
				"vendor": "Google",
				"version": "Google",
				"date": "06/07/2024"
			},
			"baseboard": {
				"asset_tag": "60D0B3D6-4AB7-A513-16B0-4AE13317D4AE",
				"serial_number": "unknown",
				"vendor": "Google",
				"version": "",
				"product": "Google Compute Engine"
			},
			"product": {
				"family": "",
				"name": "Google Compute Engine",
				"vendor": "Google",
				"serial_number": "unknown",
				"uuid": "unknown",
				"sku": "",
				"version": ""
			},
			"pci": {
				"Devices": [
					{
						"driver": "",
						"address": "0000:00:00.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "1237",
							"name": "440FX - 82441FX PMC [Natoma]"
						},
						"revision": "0x02",
						"subsystem": {
							"id": "1100",
							"name": "Qemu virtual machine"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "00",
							"name": "Host bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7110",
							"name": "82371AB/EB/MB PIIX4 ISA"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "01",
							"name": "ISA bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.3",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7113",
							"name": "82371AB/EB/MB PIIX4 ACPI"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "80",
							"name": "Bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:03.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1004",
							"name": "Virtio SCSI"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0008",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "00",
							"name": "Non-VGA unclassified device"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:04.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1000",
							"name": "Virtio network device"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0001",
							"name": "unknown"
						},
						"class": {
							"id": "02",
							"name": "Network controller"
						},
						"subclass": {
							"id": "00",
							"name": "Ethernet controller"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:05.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1005",
							"name": "Virtio RNG"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0004",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "ff",
							"name": "unknown"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					}
				]
			}
		},
		"Kernel": {
			"Version": "Linux version 4.18.0-513.24.1.el8_9.x86_64 (mockbuild@iad1-prod-build001.bld.equ.rockylinux.org) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-20) (GCC)) #1 SMP Thu Apr 4 18:13:02 UTC 2024\n",
			"Modules": "tcp_diag 16384 0 - Live 0x0000000000000000\ninet_diag 24576 1 tcp_diag, Live 0x0000000000000000\nbinfmt_misc 24576 1 - Live 0x0000000000000000\nnfsv3 57344 1 - Live 0x0000000000000000\nrpcsec_gss_krb5 45056 0 - Live 0x0000000000000000\nnfsv4 917504 2 - Live 0x0000000000000000\ndns_resolver 16384 1 nfsv4, Live 0x0000000000000000\nnfs 425984 4 nfsv3,nfsv4, Live 0x0000000000000000\nfscache 389120 1 nfs, Live 0x0000000000000000\nintel_rapl_msr 16384 0 - Live 0x0000000000000000\nintel_rapl_common 24576 1 intel_rapl_msr, Live 0x0000000000000000\nintel_uncore_frequency_common 16384 0 - Live 0x0000000000000000\nisst_if_common 16384 0 - Live 0x0000000000000000\nnfit 65536 0 - Live 0x0000000000000000\nlibnvdimm 200704 1 nfit, Live 0x0000000000000000\ncrct10dif_pclmul 16384 1 - Live 0x0000000000000000\ncrc32_pclmul 16384 0 - Live 0x0000000000000000\nghash_clmulni_intel 16384 0 - Live 0x0000000000000000\nrapl 20480 0 - Live 0x0000000000000000\nvfat 20480 1 - Live 0x0000000000000000\nfat 86016 1 vfat, Live 0x0000000000000000\ni2c_piix4 24576 0 - Live 0x0000000000000000\npcspkr 16384 0 - Live 0x0000000000000000\nnfsd 548864 13 - Live 0x0000000000000000\nauth_rpcgss 139264 2 rpcsec_gss_krb5,nfsd, Live 0x0000000000000000\nnfs_acl 16384 2 nfsv3,nfsd, Live 0x0000000000000000\nlockd 126976 3 nfsv3,nfs,nfsd, Live 0x0000000000000000\ngrace 16384 2 nfsd,lockd, Live 0x0000000000000000\nsunrpc 585728 32 nfsv3,rpcsec_gss_krb5,nfsv4,nfs,nfsd,auth_rpcgss,nfs_acl,lockd, Live 0x0000000000000000\nxfs 1593344 1 - Live 0x0000000000000000\nlibcrc32c 16384 1 xfs, Live 0x0000000000000000\nsd_mod 57344 2 - Live 0x0000000000000000\nsg 40960 0 - Live 0x0000000000000000\nvirtio_net 61440 0 - Live 0x0000000000000000\nnet_failover 24576 1 virtio_net, Live 0x0000000000000000\ncrc32c_intel 24576 1 - Live 0x0000000000000000\nserio_raw 16384 0 - Live 0x0000000000000000\nfailover 16384 1 net_failover, Live 0x0000000000000000\nvirtio_scsi 20480 2 - Live 0x0000000000000000\nnvme 45056 0 - Live 0x0000000000000000\nnvme_core 139264 1 nvme, Live 0x0000000000000000\nt10_pi 16384 2 sd_mod,nvme_core, Live 0x0000000000000000\n",
			"Drivers": "Character devices:\n  1 mem\n  4 /dev/vc/0\n  4 tty\n  4 ttyS\n  5 /dev/tty\n  5 /dev/console\n  5 /dev/ptmx\n  7 vcs\n 10 misc\n 13 input\n 21 sg\n 29 fb\n128 ptm\n136 pts\n162 raw\n180 usb\n188 ttyUSB\n189 usb_device\n202 cpu/msr\n203 cpu/cpuid\n240 dimmctl\n241 ndctl\n242 nvme-generic\n243 nvme\n244 hidraw\n245 ttyDBC\n246 usbmon\n247 bsg\n248 watchdog\n249 ptp\n250 pps\n251 rtc\n252 dax\n253 tpm\n254 gpiochip\n\nBlock devices:\n  8 sd\n  9 md\n 65 sd\n 66 sd\n 67 sd\n 68 sd\n 69 sd\n 70 sd\n 71 sd\n128 sd\n129 sd\n130 sd\n131 sd\n132 sd\n133 sd\n134 sd\n135 sd\n254 mdp\n259 blkext\n"
		},
		"Stderr": "WARNING: \n/sys/class/drm does not exist on this system (likely the host system is a\nvirtual machine or container with no graphics). Therefore,\nGPUInfo.GraphicsCards will be an empty array.\nWARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied\nWARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied\nWARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied\nWARNING: Unable to read product_uuid: open /sys/class/dmi/id/product_uuid: permission denied\n",
		"Err": ""
	},
	{
		"Hostname": "hpcslurm-debugnodeset-8",
		"Info": {
			"memory": {
				"total_physical_bytes": 8589934592,
				"total_usable_bytes": 8056086528,
				"supported_page_sizes": [
					1073741824,
					2097152
				],
				"modules": null
			},
			"block": {
				"total_size_bytes": 53687091200,
				"disks": [
					{
						"name": "sda",
						"size_bytes": 53687091200,
						"physical_block_size_bytes": 4096,
						"drive_type": "hdd",
						"removable": false,
						"storage_controller": "scsi",
						"bus_path": "pci-0000:00:03.0-scsi-0:0:1:0",
						"vendor": "Google",
						"model": "PersistentDisk",
						"serial_number": "persistent-disk-0",
						"wwn": "unknown",
						"partitions": [
							{
								"name": "sda1",
								"label": "EFI\\x20System\\x20Partition",
								"mount_point": "/boot/efi",
								"size_bytes": 209715200,
								"type": "vfat",
								"read_only": false,
								"uuid": "a407d4b7-cfe4-4f7e-b9fc-ee7799ba3b84",
								"filesystem_label": "unknown"
							},
							{
								"name": "sda2",
								"label": "unknown",
								"mount_point": "/",
								"size_bytes": 53475328000,
								"type": "xfs",
								"read_only": false,
								"uuid": "144c8c6f-9c84-47c9-b637-8b7723fdb3ef",
								"filesystem_label": "root"
							}
						]
					}
				]
			},
			"cpu": {
				"total_cores": 1,
				"total_threads": 1,
				"processors": [
					{
						"id": 0,
						"total_cores": 1,
						"total_threads": 1,
						"vendor": "GenuineIntel",
						"model": "Intel(R) Xeon(R) CPU @ 2.80GHz",
						"capabilities": [
							"fpu",
							"vme",
							"de",
							"pse",
							"tsc",
							"msr",
							"pae",
							"mce",
							"cx8",
							"apic",
							"sep",
							"mtrr",
							"pge",
							"mca",
							"cmov",
							"pat",
							"pse36",
							"clflush",
							"mmx",
							"fxsr",
							"sse",
							"sse2",
							"ss",
							"ht",
							"syscall",
							"nx",
							"pdpe1gb",
							"rdtscp",
							"lm",
							"constant_tsc",
							"rep_good",
							"nopl",
							"xtopology",
							"nonstop_tsc",
							"cpuid",
							"tsc_known_freq",
							"pni",
							"pclmulqdq",
							"ssse3",
							"fma",
							"cx16",
							"pcid",
							"sse4_1",
							"sse4_2",
							"x2apic",
							"movbe",
							"popcnt",
							"aes",
							"xsave",
							"avx",
							"f16c",
							"rdrand",
							"hypervisor",
							"lahf_lm",
							"abm",
							"3dnowprefetch",
							"invpcid_single",
							"ssbd",
							"ibrs",
							"ibpb",
							"stibp",
							"ibrs_enhanced",
							"fsgsbase",
							"tsc_adjust",
							"bmi1",
							"hle",
							"avx2",
							"smep",
							"bmi2",
							"erms",
							"invpcid",
							"rtm",
							"avx512f",
							"avx512dq",
							"rdseed",
							"adx",
							"smap",
							"clflushopt",
							"clwb",
							"avx512cd",
							"avx512bw",
							"avx512vl",
							"xsaveopt",
							"xsavec",
							"xgetbv1",
							"xsaves",
							"arat",
							"avx512_vnni",
							"md_clear",
							"arch_capabilities"
						],
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						]
					}
				]
			},
			"topology": {
				"architecture": "smp",
				"nodes": [
					{
						"id": 0,
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						],
						"caches": [
							{
								"level": 1,
								"type": "instruction",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 1,
								"type": "data",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 2,
								"type": "unified",
								"size_bytes": 1048576,
								"logical_processors": [
									0
								]
							},
							{
								"level": 3,
								"type": "unified",
								"size_bytes": 34603008,
								"logical_processors": [
									0
								]
							}
						],
						"distances": [
							10
						],
						"memory": {
							"total_physical_bytes": 8589934592,
							"total_usable_bytes": 8056086528,
							"supported_page_sizes": [
								1073741824,
								2097152
							],
							"modules": null
						}
					}
				]
			},
			"network": {
				"nics": [
					{
						"name": "eth0",
						"mac_address": "42:01:0a:00:00:cd",
						"is_virtual": false,
						"capabilities": [
							{
								"name": "auto-negotiation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "pause-frame-use",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-checksumming",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-checksumming",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv4",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-ip-generic",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv6",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-fcoe-crc",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-sctp",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather-fraglist",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tcp-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-ecn-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tcp-mangleid-segmentation",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tx-tcp6-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-receive-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "large-receive-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "ntuple-filters",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "receive-hashing",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "highdma",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "rx-vlan-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "vlan-challenged",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-lockless",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "netns-local",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-robust",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-fcoe-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip4-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip6-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-partial",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tunnel-remcsum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-sctp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-esp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-list",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp-gro-forwarding",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "rx-gro-list",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tls-hw-rx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "fcoe-mtu",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-nocache-copy",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "loopback",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-fcs",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-all",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-stag-hw-insert",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-hw-parse",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "l2-fwd-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "hw-tc-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-tx-csum-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp_tunnel-port-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-tx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-gro-hw",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-record",
								"is_enabled": false,
								"can_enable": false
							}
						],
						"speed": "Unknown!",
						"duplex": "Unknown!(255)"
					}
				]
			},
			"gpu": {
				"cards": null
			},
			"chassis": {
				"asset_tag": "",
				"serial_number": "unknown",
				"type": "1",
				"type_description": "Other",
				"vendor": "Google",
				"version": ""
			},
			"bios": {
				"vendor": "Google",
				"version": "Google",
				"date": "06/07/2024"
			},
			"baseboard": {
				"asset_tag": "D2907AF3-24F5-8D76-0A8D-A7FB057B6EC4",
				"serial_number": "unknown",
				"vendor": "Google",
				"version": "",
				"product": "Google Compute Engine"
			},
			"product": {
				"family": "",
				"name": "Google Compute Engine",
				"vendor": "Google",
				"serial_number": "unknown",
				"uuid": "unknown",
				"sku": "",
				"version": ""
			},
			"pci": {
				"Devices": [
					{
						"driver": "",
						"address": "0000:00:00.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "1237",
							"name": "440FX - 82441FX PMC [Natoma]"
						},
						"revision": "0x02",
						"subsystem": {
							"id": "1100",
							"name": "Qemu virtual machine"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "00",
							"name": "Host bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7110",
							"name": "82371AB/EB/MB PIIX4 ISA"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "01",
							"name": "ISA bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.3",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7113",
							"name": "82371AB/EB/MB PIIX4 ACPI"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "80",
							"name": "Bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:03.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1004",
							"name": "Virtio SCSI"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0008",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "00",
							"name": "Non-VGA unclassified device"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:04.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1000",
							"name": "Virtio network device"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0001",
							"name": "unknown"
						},
						"class": {
							"id": "02",
							"name": "Network controller"
						},
						"subclass": {
							"id": "00",
							"name": "Ethernet controller"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:05.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1005",
							"name": "Virtio RNG"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0004",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "ff",
							"name": "unknown"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					}
				]
			}
		},
		"Kernel": {
			"Version": "Linux version 4.18.0-513.24.1.el8_9.x86_64 (mockbuild@iad1-prod-build001.bld.equ.rockylinux.org) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-20) (GCC)) #1 SMP Thu Apr 4 18:13:02 UTC 2024\n",
			"Modules": "tcp_diag 16384 0 - Live 0x0000000000000000\ninet_diag 24576 1 tcp_diag, Live 0x0000000000000000\nbinfmt_misc 24576 1 - Live 0x0000000000000000\nnfsv3 57344 1 - Live 0x0000000000000000\nrpcsec_gss_krb5 45056 0 - Live 0x0000000000000000\nnfsv4 917504 2 - Live 0x0000000000000000\ndns_resolver 16384 1 nfsv4, Live 0x0000000000000000\nnfs 425984 4 nfsv3,nfsv4, Live 0x0000000000000000\nfscache 389120 1 nfs, Live 0x0000000000000000\nintel_rapl_msr 16384 0 - Live 0x0000000000000000\nintel_rapl_common 24576 1 intel_rapl_msr, Live 0x0000000000000000\nintel_uncore_frequency_common 16384 0 - Live 0x0000000000000000\nisst_if_common 16384 0 - Live 0x0000000000000000\nnfit 65536 0 - Live 0x0000000000000000\nlibnvdimm 200704 1 nfit, Live 0x0000000000000000\ncrct10dif_pclmul 16384 1 - Live 0x0000000000000000\ncrc32_pclmul 16384 0 - Live 0x0000000000000000\nghash_clmulni_intel 16384 0 - Live 0x0000000000000000\nrapl 20480 0 - Live 0x0000000000000000\nvfat 20480 1 - Live 0x0000000000000000\nfat 86016 1 vfat, Live 0x0000000000000000\ni2c_piix4 24576 0 - Live 0x0000000000000000\npcspkr 16384 0 - Live 0x0000000000000000\nnfsd 548864 13 - Live 0x0000000000000000\nauth_rpcgss 139264 2 rpcsec_gss_krb5,nfsd, Live 0x0000000000000000\nnfs_acl 16384 2 nfsv3,nfsd, Live 0x0000000000000000\nlockd 126976 3 nfsv3,nfs,nfsd, Live 0x0000000000000000\ngrace 16384 2 nfsd,lockd, Live 0x0000000000000000\nsunrpc 585728 32 nfsv3,rpcsec_gss_krb5,nfsv4,nfs,nfsd,auth_rpcgss,nfs_acl,lockd, Live 0x0000000000000000\nxfs 1593344 1 - Live 0x0000000000000000\nlibcrc32c 16384 1 xfs, Live 0x0000000000000000\nsd_mod 57344 2 - Live 0x0000000000000000\nsg 40960 0 - Live 0x0000000000000000\nvirtio_net 61440 0 - Live 0x0000000000000000\ncrc32c_intel 24576 1 - Live 0x0000000000000000\nserio_raw 16384 0 - Live 0x0000000000000000\nnet_failover 24576 1 virtio_net, Live 0x0000000000000000\nfailover 16384 1 net_failover, Live 0x0000000000000000\nvirtio_scsi 20480 2 - Live 0x0000000000000000\nnvme 45056 0 - Live 0x0000000000000000\nnvme_core 139264 1 nvme, Live 0x0000000000000000\nt10_pi 16384 2 sd_mod,nvme_core, Live 0x0000000000000000\n",
			"Drivers": "Character devices:\n  1 mem\n  4 /dev/vc/0\n  4 tty\n  4 ttyS\n  5 /dev/tty\n  5 /dev/console\n  5 /dev/ptmx\n  7 vcs\n 10 misc\n 13 input\n 21 sg\n 29 fb\n128 ptm\n136 pts\n162 raw\n180 usb\n188 ttyUSB\n189 usb_device\n202 cpu/msr\n203 cpu/cpuid\n240 dimmctl\n241 ndctl\n242 nvme-generic\n243 nvme\n244 hidraw\n245 ttyDBC\n246 usbmon\n247 bsg\n248 watchdog\n249 ptp\n250 pps\n251 rtc\n252 dax\n253 tpm\n254 gpiochip\n\nBlock devices:\n  8 sd\n  9 md\n 65 sd\n 66 sd\n 67 sd\n 68 sd\n 69 sd\n 70 sd\n 71 sd\n128 sd\n129 sd\n130 sd\n131 sd\n132 sd\n133 sd\n134 sd\n135 sd\n254 mdp\n259 blkext\n"
		},
		"Stderr": "WARNING: \n/sys/class/drm does not exist on this system (likely the host system is a\nvirtual machine or container with no graphics). Therefore,\nGPUInfo.GraphicsCards will be an empty array.\nWARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied\nWARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied\nWARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied\nWARNING: Unable to read product_uuid: open /sys/class/dmi/id/product_uuid: permission denied\n",
		"Err": ""
	},
	{
		"Hostname": "hpcslurm-debugnodeset-11",
		"Info": {
			"memory": {
				"total_physical_bytes": 8589934592,
				"total_usable_bytes": 8056090624,
				"supported_page_sizes": [
					1073741824,
					2097152
				],
				"modules": null
			},
			"block": {
				"total_size_bytes": 53687091200,
				"disks": [
					{
						"name": "sda",
						"size_bytes": 53687091200,
						"physical_block_size_bytes": 4096,
						"drive_type": "hdd",
						"removable": false,
						"storage_controller": "scsi",
						"bus_path": "pci-0000:00:03.0-scsi-0:0:1:0",
						"vendor": "Google",
						"model": "PersistentDisk",
						"serial_number": "persistent-disk-0",
						"wwn": "unknown",
						"partitions": [
							{
								"name": "sda1",
								"label": "EFI\\x20System\\x20Partition",
								"mount_point": "/boot/efi",
								"size_bytes": 209715200,
								"type": "vfat",
								"read_only": false,
								"uuid": "a407d4b7-cfe4-4f7e-b9fc-ee7799ba3b84",
								"filesystem_label": "unknown"
							},
							{
								"name": "sda2",
								"label": "unknown",
								"mount_point": "/",
								"size_bytes": 53475328000,
								"type": "xfs",
								"read_only": false,
								"uuid": "144c8c6f-9c84-47c9-b637-8b7723fdb3ef",
								"filesystem_label": "root"
							}
						]
					}
				]
			},
			"cpu": {
				"total_cores": 1,
				"total_threads": 1,
				"processors": [
					{
						"id": 0,
						"total_cores": 1,
						"total_threads": 1,
						"vendor": "GenuineIntel",
						"model": "Intel(R) Xeon(R) CPU @ 2.80GHz",
						"capabilities": [
							"fpu",
							"vme",
							"de",
							"pse",
							"tsc",
							"msr",
							"pae",
							"mce",
							"cx8",
							"apic",
							"sep",
							"mtrr",
							"pge",
							"mca",
							"cmov",
							"pat",
							"pse36",
							"clflush",
							"mmx",
							"fxsr",
							"sse",
							"sse2",
							"ss",
							"ht",
							"syscall",
							"nx",
							"pdpe1gb",
							"rdtscp",
							"lm",
							"constant_tsc",
							"rep_good",
							"nopl",
							"xtopology",
							"nonstop_tsc",
							"cpuid",
							"tsc_known_freq",
							"pni",
							"pclmulqdq",
							"ssse3",
							"fma",
							"cx16",
							"pcid",
							"sse4_1",
							"sse4_2",
							"x2apic",
							"movbe",
							"popcnt",
							"aes",
							"xsave",
							"avx",
							"f16c",
							"rdrand",
							"hypervisor",
							"lahf_lm",
							"abm",
							"3dnowprefetch",
							"invpcid_single",
							"ssbd",
							"ibrs",
							"ibpb",
							"stibp",
							"ibrs_enhanced",
							"fsgsbase",
							"tsc_adjust",
							"bmi1",
							"hle",
							"avx2",
							"smep",
							"bmi2",
							"erms",
							"invpcid",
							"rtm",
							"avx512f",
							"avx512dq",
							"rdseed",
							"adx",
							"smap",
							"clflushopt",
							"clwb",
							"avx512cd",
							"avx512bw",
							"avx512vl",
							"xsaveopt",
							"xsavec",
							"xgetbv1",
							"xsaves",
							"arat",
							"avx512_vnni",
							"md_clear",
							"arch_capabilities"
						],
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						]
					}
				]
			},
			"topology": {
				"architecture": "smp",
				"nodes": [
					{
						"id": 0,
						"cores": [
							{
								"id": 0,
								"total_threads": 1,
								"logical_processors": [
									0
								]
							}
						],
						"caches": [
							{
								"level": 1,
								"type": "instruction",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 1,
								"type": "data",
								"size_bytes": 32768,
								"logical_processors": [
									0
								]
							},
							{
								"level": 2,
								"type": "unified",
								"size_bytes": 1048576,
								"logical_processors": [
									0
								]
							},
							{
								"level": 3,
								"type": "unified",
								"size_bytes": 34603008,
								"logical_processors": [
									0
								]
							}
						],
						"distances": [
							10
						],
						"memory": {
							"total_physical_bytes": 8589934592,
							"total_usable_bytes": 8056090624,
							"supported_page_sizes": [
								1073741824,
								2097152
							],
							"modules": null
						}
					}
				]
			},
			"network": {
				"nics": [
					{
						"name": "eth0",
						"mac_address": "42:01:0a:00:00:e6",
						"is_virtual": false,
						"capabilities": [
							{
								"name": "auto-negotiation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "pause-frame-use",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-checksumming",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-checksumming",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv4",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-ip-generic",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-checksum-ipv6",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-fcoe-crc",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-checksum-sctp",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-scatter-gather-fraglist",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tcp-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "tx-tcp-ecn-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tcp-mangleid-segmentation",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tx-tcp6-segmentation",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-segmentation-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "generic-receive-offload",
								"is_enabled": true,
								"can_enable": true
							},
							{
								"name": "large-receive-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "ntuple-filters",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "receive-hashing",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "highdma",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "rx-vlan-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "vlan-challenged",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-lockless",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "netns-local",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-robust",
								"is_enabled": true,
								"can_enable": false
							},
							{
								"name": "tx-fcoe-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gre-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip4-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-ipxip6-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp_tnl-csum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-partial",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-tunnel-remcsum-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-sctp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-esp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-udp-segmentation",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-gso-list",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp-gro-forwarding",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "rx-gro-list",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "tls-hw-rx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "fcoe-mtu",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-nocache-copy",
								"is_enabled": false,
								"can_enable": true
							},
							{
								"name": "loopback",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-fcs",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-all",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tx-vlan-stag-hw-insert",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-hw-parse",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-vlan-stag-filter",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "l2-fwd-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "hw-tc-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "esp-tx-csum-hw-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-udp_tunnel-port-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-tx-offload",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "rx-gro-hw",
								"is_enabled": false,
								"can_enable": false
							},
							{
								"name": "tls-hw-record",
								"is_enabled": false,
								"can_enable": false
							}
						],
						"speed": "Unknown!",
						"duplex": "Unknown!(255)"
					}
				]
			},
			"gpu": {
				"cards": null
			},
			"chassis": {
				"asset_tag": "",
				"serial_number": "unknown",
				"type": "1",
				"type_description": "Other",
				"vendor": "Google",
				"version": ""
			},
			"bios": {
				"vendor": "Google",
				"version": "Google",
				"date": "06/27/2024"
			},
			"baseboard": {
				"asset_tag": "92341A9B-0DA8-75C9-C055-C2EDFD2EAE64",
				"serial_number": "unknown",
				"vendor": "Google",
				"version": "",
				"product": "Google Compute Engine"
			},
			"product": {
				"family": "",
				"name": "Google Compute Engine",
				"vendor": "Google",
				"serial_number": "unknown",
				"uuid": "unknown",
				"sku": "",
				"version": ""
			},
			"pci": {
				"Devices": [
					{
						"driver": "",
						"address": "0000:00:00.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "1237",
							"name": "440FX - 82441FX PMC [Natoma]"
						},
						"revision": "0x02",
						"subsystem": {
							"id": "1100",
							"name": "Qemu virtual machine"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "00",
							"name": "Host bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.0",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7110",
							"name": "82371AB/EB/MB PIIX4 ISA"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "01",
							"name": "ISA bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "",
						"address": "0000:00:01.3",
						"vendor": {
							"id": "8086",
							"name": "Intel Corporation"
						},
						"product": {
							"id": "7113",
							"name": "82371AB/EB/MB PIIX4 ACPI"
						},
						"revision": "0x03",
						"subsystem": {
							"id": "0000",
							"name": "unknown"
						},
						"class": {
							"id": "06",
							"name": "Bridge"
						},
						"subclass": {
							"id": "80",
							"name": "Bridge"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:03.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1004",
							"name": "Virtio SCSI"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0008",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "00",
							"name": "Non-VGA unclassified device"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:04.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1000",
							"name": "Virtio network device"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0001",
							"name": "unknown"
						},
						"class": {
							"id": "02",
							"name": "Network controller"
						},
						"subclass": {
							"id": "00",
							"name": "Ethernet controller"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					},
					{
						"driver": "virtio-pci",
						"address": "0000:00:05.0",
						"vendor": {
							"id": "1af4",
							"name": "Red Hat, Inc."
						},
						"product": {
							"id": "1005",
							"name": "Virtio RNG"
						},
						"revision": "0x00",
						"subsystem": {
							"id": "0004",
							"name": "unknown"
						},
						"class": {
							"id": "00",
							"name": "Unclassified device"
						},
						"subclass": {
							"id": "ff",
							"name": "unknown"
						},
						"programming_interface": {
							"id": "00",
							"name": "unknown"
						}
					}
				]
			}
		},
		"Kernel": {
			"Version": "Linux version 4.18.0-513.24.1.el8_9.x86_64 (mockbuild@iad1-prod-build001.bld.equ.rockylinux.org) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-20) (GCC)) #1 SMP Thu Apr 4 18:13:02 UTC 2024\n",
			"Modules": "tcp_diag 16384 0 - Live 0x0000000000000000\ninet_diag 24576 1 tcp_diag, Live 0x0000000000000000\nbinfmt_misc 24576 1 - Live 0x0000000000000000\nnfsv3 57344 1 - Live 0x0000000000000000\nrpcsec_gss_krb5 45056 0 - Live 0x0000000000000000\nnfsv4 917504 2 - Live 0x0000000000000000\ndns_resolver 16384 1 nfsv4, Live 0x0000000000000000\nnfs 425984 4 nfsv3,nfsv4, Live 0x0000000000000000\nfscache 389120 1 nfs, Live 0x0000000000000000\nintel_rapl_msr 16384 0 - Live 0x0000000000000000\nintel_rapl_common 24576 1 intel_rapl_msr, Live 0x0000000000000000\nintel_uncore_frequency_common 16384 0 - Live 0x0000000000000000\nisst_if_common 16384 0 - Live 0x0000000000000000\nnfit 65536 0 - Live 0x0000000000000000\nlibnvdimm 200704 1 nfit, Live 0x0000000000000000\ncrct10dif_pclmul 16384 1 - Live 0x0000000000000000\ncrc32_pclmul 16384 0 - Live 0x0000000000000000\nghash_clmulni_intel 16384 0 - Live 0x0000000000000000\nrapl 20480 0 - Live 0x0000000000000000\ni2c_piix4 24576 0 - Live 0x0000000000000000\nvfat 20480 1 - Live 0x0000000000000000\nfat 86016 1 vfat, Live 0x0000000000000000\npcspkr 16384 0 - Live 0x0000000000000000\nnfsd 548864 13 - Live 0x0000000000000000\nauth_rpcgss 139264 2 rpcsec_gss_krb5,nfsd, Live 0x0000000000000000\nnfs_acl 16384 2 nfsv3,nfsd, Live 0x0000000000000000\nlockd 126976 3 nfsv3,nfs,nfsd, Live 0x0000000000000000\ngrace 16384 2 nfsd,lockd, Live 0x0000000000000000\nsunrpc 585728 32 nfsv3,rpcsec_gss_krb5,nfsv4,nfs,nfsd,auth_rpcgss,nfs_acl,lockd, Live 0x0000000000000000\nxfs 1593344 1 - Live 0x0000000000000000\nlibcrc32c 16384 1 xfs, Live 0x0000000000000000\nsd_mod 57344 2 - Live 0x0000000000000000\nsg 40960 0 - Live 0x0000000000000000\nvirtio_net 61440 0 - Live 0x0000000000000000\ncrc32c_intel 24576 1 - Live 0x0000000000000000\nserio_raw 16384 0 - Live 0x0000000000000000\nnet_failover 24576 1 virtio_net, Live 0x0000000000000000\nfailover 16384 1 net_failover, Live 0x0000000000000000\nvirtio_scsi 20480 2 - Live 0x0000000000000000\nnvme 45056 0 - Live 0x0000000000000000\nnvme_core 139264 1 nvme, Live 0x0000000000000000\nt10_pi 16384 2 sd_mod,nvme_core, Live 0x0000000000000000\n",
			"Drivers": "Character devices:\n  1 mem\n  4 /dev/vc/0\n  4 tty\n  4 ttyS\n  5 /dev/tty\n  5 /dev/console\n  5 /dev/ptmx\n  7 vcs\n 10 misc\n 13 input\n 21 sg\n 29 fb\n128 ptm\n136 pts\n162 raw\n180 usb\n188 ttyUSB\n189 usb_device\n202 cpu/msr\n203 cpu/cpuid\n240 dimmctl\n241 ndctl\n242 nvme-generic\n243 nvme\n244 hidraw\n245 ttyDBC\n246 usbmon\n247 bsg\n248 watchdog\n249 ptp\n250 pps\n251 rtc\n252 dax\n253 tpm\n254 gpiochip\n\nBlock devices:\n  8 sd\n  9 md\n 65 sd\n 66 sd\n 67 sd\n 68 sd\n 69 sd\n 70 sd\n 71 sd\n128 sd\n129 sd\n130 sd\n131 sd\n132 sd\n133 sd\n134 sd\n135 sd\n254 mdp\n259 blkext\n"
		},
		"Stderr": "WARNING: \n/sys/class/drm does not exist on this system (likely the host system is a\nvirtual machine or container with no graphics). Therefore,\nGPUInfo.GraphicsCards will be an empty array.\nWARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied\nWARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied\nWARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied\nWARNING: Unable to read product_uuid: open /sys/class/dmi/id/product_uuid: permission denied\n",
		"Err": ""
	}
]
`
