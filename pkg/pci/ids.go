// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A slice of bytes of the literal plain text of pci.ids has been found
// to produce the smallest binary compared to a native Go map, marshalled
// JSON, and go-bindata (gzip'ed bytes). Further runtime of parsing the plain
// text pci.ids is lower than all options compared. The pciids in this package
// is stripped of all comments, empty lines, sub-devices, and classes to save
// on binary size.

package pci

type idMap map[uint16]Vendor

var ids idMap

// newIDs contains the plain text contents of pci.ids. It returns
// a map to be used as lookup from hex ID to human readable lable.
// We do not admit of the possibility of error, any failure
// should be caught by the test. We might just want to just always
// create ids since the most common use of pci will be with names,
// not numbers.
func newIDs() idMap {
	if ids != nil {
		return ids
	}

	pciids := []byte(`0001  SafeNet (wrong ID)
0010  Allied Telesis, Inc (Wrong ID)
	8139  AT-2500TX V3 Ethernet
001c  PEAK-System Technik GmbH
	0001  PCAN-PCI CAN-Bus controller
003d  Lockheed Martin-Marietta Corp
0059  Tiger Jet Network Inc. (Wrong ID)
0070  Hauppauge computer works Inc.
	7801  WinTV HVR-1800 MCE
0071  Nebula Electronics Ltd.
0095  Silicon Image, Inc. (Wrong ID)
	0680  Ultra ATA/133 IDE RAID CONTROLLER CARD
00a7  Teles AG (Wrong ID)
0100  Thales e-Security
0123  General Dynamics
018a  LevelOne
	0106  FPC-0106TX misprogrammed [RTL81xx]
021b  Compaq Computer Corporation
	8139  HNE-300 (RealTek RTL8139c) [iPaq Networking]
0270  Hauppauge computer works Inc. (Wrong ID)
0291  Davicom Semiconductor, Inc. (Wrong ID)
02ac  SpeedStream
	1012  1012 PCMCIA 10/100 Ethernet Card [RTL81xx]
02e0  XFX Pine Group Inc. (Wrong ID)
0303  Hewlett-Packard Company (Wrong ID)
0308  ZyXEL Communications Corporation (Wrong ID)
0315  SK-Electronics Co., Ltd.
0357  TTTech Computertechnik AG (Wrong ID)
	000a  TTP-Monitoring Card V2.0
0432  SCM Microsystems, Inc.
	0001  Pluto2 DVB-T Receiver for PCMCIA [EasyWatch MobilSet]
0675  Dynalink
	1700  IS64PH ISDN Adapter
	1702  IS64PH ISDN Adapter
	1703  ISDN Adapter (PCI Bus, DV, W)
	1704  ISDN Adapter (PCI Bus, D, C)
0721  Sapphire, Inc.
0777  Ubiquiti Networks, Inc.
0795  Wired Inc.
	6663  Butane II (MPEG2 encoder board)
	6666  MediaPress (MPEG2 encoder board)
07d1  D-Link System Inc
0925  VIA Technologies, Inc. (Wrong ID)
0a89  BREA Technologies Inc
0b0b  Rhino Equipment Corp.
	0105  R1T1
	0205  R4FXO
	0206  RCB4FXO 4-channel FXO analog telephony card
	0305  R4T1
	0405  R8FXX
	0406  RCB8FXX 8-channel modular analog telephony card
	0505  R24FXX
	0506  RCB24FXS 24-Channel FXS analog telephony card
	0605  R2T1
	0705  R24FXS
	0706  RCB24FXO 24-Channel FXO analog telephony card
	0905  R1T3 Single T3 Digital Telephony Card
	0906  RCB24FXX 24-channel modular analog telephony card
	0a06  RCB672FXX 672-channel modular analog telephony card
0e11  Compaq Computer Corporation
	0001  PCI to EISA Bridge
	0002  PCI to ISA Bridge
	0046  Smart Array 64xx
	0049  NC7132 Gigabit Upgrade Module
	004a  NC6136 Gigabit Server Adapter
	005a  Remote Insight II board - Lights-Out
	007c  NC7770 1000BaseTX
	007d  NC6770 1000BaseTX
	0085  NC7780 1000BaseTX
	00b1  Remote Insight II board - PCI device
	00bb  NC7760
	00ca  NC7771
	00cb  NC7781
	00cf  NC7772
	00d0  NC7782
	00d1  NC7783
	00e3  NC7761
	0508  Netelligent 4/16 Token Ring
	1000  Triflex/Pentium Bridge, Model 1000
	2000  Triflex/Pentium Bridge, Model 2000
	3032  QVision 1280/p
	3033  QVision 1280/p
	3034  QVision 1280/p
	4000  4000 [Triflex]
	4040  Integrated Array
	4048  Compaq Raid LC2
	4050  Smart Array 4200
	4051  Smart Array 4250ES
	4058  Smart Array 431
	4070  Smart Array 5300
	4080  Smart Array 5i
	4082  Smart Array 532
	4083  Smart Array 5312
	4091  Smart Array 6i
	409a  Smart Array 641
	409b  Smart Array 642
	409c  Smart Array 6400
	409d  Smart Array 6400 EM
	6010  HotPlug PCI Bridge 6010
	7020  USB Controller
	a0ec  Fibre Channel Host Controller
	a0f0  Advanced System Management Controller
	a0f3  Triflex PCI to ISA Bridge
	a0f7  PCI Hotplug Controller
	a0f8  ZFMicro Chipset USB
	a0fc  FibreChannel HBA Tachyon
	ae10  Smart-2/P RAID Controller
	ae29  MIS-L
	ae2a  MPC
	ae2b  MIS-E
	ae31  System Management Controller
	ae32  Netelligent 10/100 TX PCI UTP
	ae33  Triflex Dual EIDE Controller
	ae34  Netelligent 10 T PCI UTP
	ae35  Integrated NetFlex-3/P
	ae40  Netelligent Dual 10/100 TX PCI UTP
	ae43  Netelligent Integrated 10/100 TX UTP
	ae69  CETUS-L
	ae6c  Northstar
	ae6d  NorthStar CPU to PCI Bridge
	b011  Netelligent 10/100 TX Embedded UTP
	b012  Netelligent 10 T/2 PCI UTP/Coax
	b01e  NC3120 Fast Ethernet NIC
	b01f  NC3122 Fast Ethernet NIC
	b02f  NC1120 Ethernet NIC
	b030  Netelligent 10/100 TX UTP
	b04a  10/100 TX PCI Intel WOL UTP Controller
	b060  Smart Array 5300 Controller
	b0c6  NC3161 Fast Ethernet NIC
	b0c7  NC3160 Fast Ethernet NIC
	b0d7  NC3121 Fast Ethernet NIC
	b0dd  NC3131 Fast Ethernet NIC
	b0de  NC3132 Fast Ethernet Module
	b0df  NC6132 Gigabit Module
	b0e0  NC6133 Gigabit Module
	b0e1  NC3133 Fast Ethernet Module
	b123  NC6134 Gigabit NIC
	b134  NC3163 Fast Ethernet NIC
	b13c  NC3162 Fast Ethernet NIC
	b144  NC3123 Fast Ethernet NIC
	b163  NC3134 Fast Ethernet NIC
	b164  NC3165 Fast Ethernet Upgrade Module
	b178  Smart Array 5i/532
	b1a4  NC7131 Gigabit Server Adapter
	b200  Memory Hot-Plug Controller
	b203  Integrated Lights Out Controller
	b204  Integrated Lights Out  Processor
	c000  Remote Insight Lights-Out Edition
	f130  NetFlex-3/P ThunderLAN 1.0
	f150  NetFlex-3/P ThunderLAN 2.3
0e55  HaSoTec GmbH
0eac  SHF Communication Technologies AG
	0008  Ethernet Powerlink Managing Node 01
0f62  Acrox Technologies Co., Ltd.
1000  LSI Logic / Symbios Logic
	0001  53c810
	0002  53c820
	0003  53c825
	0004  53c815
	0005  53c810AP
	0006  53c860
	000a  53c1510
	000b  53C896/897
	000c  53c895
	000d  53c885
	000f  53c875
	0010  53C1510
	0012  53c895a
	0013  53c875a
	0014  MegaRAID Tri-Mode SAS3516
	0015  MegaRAID Tri-Mode SAS3416
	0016  MegaRAID Tri-Mode SAS3508
	0017  MegaRAID Tri-Mode SAS3408
	001b  MegaRAID Tri-Mode SAS3504
	001c  MegaRAID Tri-Mode SAS3404
	0020  53c1010 Ultra3 SCSI Adapter
	0021  53c1010 66MHz  Ultra3 SCSI Adapter
	002f  MegaRAID SAS 2208 IOV [Thunderbolt]
	0030  53c1030 PCI-X Fusion-MPT Dual Ultra320 SCSI
	0031  53c1030ZC PCI-X Fusion-MPT Dual Ultra320 SCSI
	0032  53c1035 PCI-X Fusion-MPT Dual Ultra320 SCSI
	0033  1030ZC_53c1035 PCI-X Fusion-MPT Dual Ultra320 SCSI
	0040  53c1035 PCI-X Fusion-MPT Dual Ultra320 SCSI
	0041  53C1035ZC PCI-X Fusion-MPT Dual Ultra320 SCSI
	0050  SAS1064 PCI-X Fusion-MPT SAS
	0052  MegaRAID SAS-3 3216/3224 [Cutlass]
	0053  MegaRAID SAS-3 3216/3224 [Cutlass]
	0054  SAS1068 PCI-X Fusion-MPT SAS
	0055  SAS1068 PCI-X Fusion-MPT SAS
	0056  SAS1064ET PCI-Express Fusion-MPT SAS
	0057  M1064E MegaRAID SAS
	0058  SAS1068E PCI-Express Fusion-MPT SAS
	0059  MegaRAID SAS 8208ELP/8208ELP
	005a  SAS1066E PCI-Express Fusion-MPT SAS
	005b  MegaRAID SAS 2208 [Thunderbolt]
	005c  SAS1064A PCI-X Fusion-MPT SAS
	005d  MegaRAID SAS-3 3108 [Invader]
	005e  SAS1066 PCI-X Fusion-MPT SAS
	005f  MegaRAID SAS-3 3008 [Fury]
	0060  MegaRAID SAS 1078
	0062  SAS1078 PCI-Express Fusion-MPT SAS
	0064  SAS2116 PCI-Express Fusion-MPT SAS-2 [Meteor]
	0065  SAS2116 PCI-Express Fusion-MPT SAS-2 [Meteor]
	006e  SAS2308 PCI-Express Fusion-MPT SAS-2
	0070  SAS2004 PCI-Express Fusion-MPT SAS-2 [Spitfire]
	0071  MR SAS HBA 2004
	0072  SAS2008 PCI-Express Fusion-MPT SAS-2 [Falcon]
	0073  MegaRAID SAS 2008 [Falcon]
	0074  SAS2108 PCI-Express Fusion-MPT SAS-2 [Liberator]
	0076  SAS2108 PCI-Express Fusion-MPT SAS-2 [Liberator]
	0077  SAS2108 PCI-Express Fusion-MPT SAS-2 [Liberator]
	0079  MegaRAID SAS 2108 [Liberator]
	007c  MegaRAID SAS 1078DE
	007e  SSS6200 PCI-Express Flash SSD
	0080  SAS2208 PCI-Express Fusion-MPT SAS-2
	0081  SAS2208 PCI-Express Fusion-MPT SAS-2
	0082  SAS2208 PCI-Express Fusion-MPT SAS-2
	0083  SAS2208 PCI-Express Fusion-MPT SAS-2
	0084  SAS2208 PCI-Express Fusion-MPT SAS-2
	0085  SAS2208 PCI-Express Fusion-MPT SAS-2
	0086  SAS2308 PCI-Express Fusion-MPT SAS-2
	0087  SAS2308 PCI-Express Fusion-MPT SAS-2
	008f  53c875J
	0090  SAS3108 PCI-Express Fusion-MPT SAS-3
	0091  SAS3108 PCI-Express Fusion-MPT SAS-3
	0094  SAS3108 PCI-Express Fusion-MPT SAS-3
	0095  SAS3108 PCI-Express Fusion-MPT SAS-3
	0096  SAS3004 PCI-Express Fusion-MPT SAS-3
	0097  SAS3008 PCI-Express Fusion-MPT SAS-3
	00ab  SAS3516 Fusion-MPT Tri-Mode RAID On Chip (ROC)
	00ac  SAS3416 Fusion-MPT Tri-Mode I/O Controller Chip (IOC)
	00ae  SAS3508 Fusion-MPT Tri-Mode RAID On Chip (ROC)
	00af  SAS3408 Fusion-MPT Tri-Mode I/O Controller Chip (IOC)
	00be  SAS3504 Fusion-MPT Tri-Mode RAID On Chip (ROC)
	00bf  SAS3404 Fusion-MPT Tri-Mode I/O Controller Chip (IOC)
	00c0  SAS3324 PCI-Express Fusion-MPT SAS-3
	00c1  SAS3324 PCI-Express Fusion-MPT SAS-3
	00c2  SAS3324 PCI-Express Fusion-MPT SAS-3
	00c3  SAS3324 PCI-Express Fusion-MPT SAS-3
	00c4  SAS3224 PCI-Express Fusion-MPT SAS-3
	00c5  SAS3316 PCI-Express Fusion-MPT SAS-3
	00c6  SAS3316 PCI-Express Fusion-MPT SAS-3
	00c7  SAS3316 PCI-Express Fusion-MPT SAS-3
	00c8  SAS3316 PCI-Express Fusion-MPT SAS-3
	00c9  SAS3216 PCI-Express Fusion-MPT SAS-3
	00ce  MegaRAID SAS-3 3316 [Intruder]
	00cf  MegaRAID SAS-3 3324 [Intruder]
	00d0  SAS3716 Fusion-MPT Tri-Mode RAID Controller Chip (ROC)
	00d1  SAS3616 Fusion-MPT Tri-Mode I/O Controller Chip (IOC)
	00d3  MegaRAID Tri-Mode SAS3716W
	02b0  Virtual Endpoint on PCIe Switch
	0407  MegaRAID
	0408  MegaRAID
	0409  MegaRAID
	0411  MegaRAID SAS 1068
	0413  MegaRAID SAS 1068 [Verde ZCR]
	0621  FC909 Fibre Channel Adapter
	0622  FC929 Fibre Channel Adapter
	0623  FC929 LAN
	0624  FC919 Fibre Channel Adapter
	0625  FC919 LAN
	0626  FC929X Fibre Channel Adapter
	0627  FC929X LAN
	0628  FC919X Fibre Channel Adapter
	0629  FC919X LAN
	0640  FC949X Fibre Channel Adapter
	0642  FC939X Fibre Channel Adapter
	0646  FC949ES Fibre Channel Adapter
	0701  83C885 NT50 DigitalScape Fast Ethernet
	0702  Yellowfin G-NIC gigabit ethernet
	0804  SA2010
	0805  SA2010ZC
	0806  SA2020
	0807  SA2020ZC
	0901  61C102
	1000  63C815
	1960  MegaRAID
	3050  SAS2008 PCI-Express Fusion-MPT SAS-2
	6001  DX1 Multiformat Broadcast HD/SD Encoder/Decoder
1001  Kolter Electronic
	0010  PCI 1616 Measurement card with 32 digital I/O lines
	0011  OPTO-PCI Opto-Isolated digital I/O board
	0012  PCI-AD/DA Analogue I/O board
	0013  PCI-OPTO-RELAIS Digital I/O board with relay outputs
	0014  PCI-Counter/Timer Counter Timer board
	0015  PCI-DAC416 Analogue output board
	0016  PCI-MFB Analogue I/O board
	0017  PROTO-3 PCI Prototyping board
	9100  INI-9100/9100W SCSI Host
1002  Advanced Micro Devices, Inc. [AMD/ATI]
	1304  Kaveri
	1305  Kaveri
	1306  Kaveri
	1307  Kaveri
	1308  Kaveri HDMI/DP Audio Controller
	1309  Kaveri [Radeon R6/R7 Graphics]
	130a  Kaveri [Radeon R6 Graphics]
	130b  Kaveri [Radeon R4 Graphics]
	130c  Kaveri [Radeon R7 Graphics]
	130d  Kaveri [Radeon R6 Graphics]
	130e  Kaveri [Radeon R5 Graphics]
	130f  Kaveri [Radeon R7 Graphics]
	1310  Kaveri
	1311  Kaveri
	1312  Kaveri
	1313  Kaveri [Radeon R7 Graphics]
	1314  Wrestler HDMI Audio
	1315  Kaveri [Radeon R5 Graphics]
	1316  Kaveri [Radeon R5 Graphics]
	1317  Kaveri
	1318  Kaveri [Radeon R5 Graphics]
	131b  Kaveri [Radeon R4 Graphics]
	131c  Kaveri [Radeon R7 Graphics]
	131d  Kaveri [Radeon R6 Graphics]
	1714  BeaverCreek HDMI Audio [Radeon HD 6500D and 6400G-6600G series]
	3150  RV380/M24 [Mobility Radeon X600]
	3151  RV380 GL [FireMV 2400]
	3152  RV370/M22 [Mobility Radeon X300]
	3154  RV380/M24 GL [Mobility FireGL V3200]
	3155  RV380 GL [FireMV 2400]
	3171  RV380 GL [FireMV 2400] (Secondary)
	3e50  RV380 [Radeon X600]
	3e54  RV380 GL [FireGL V3200]
	3e70  RV380 [Radeon X600] (Secondary)
	4136  RS100 [Mobility IGP 320M]
	4137  RS200 [Radeon IGP 340]
	4144  R300 [Radeon 9500]
	4146  R300 [Radeon 9700 PRO]
	4147  R300 GL [FireGL Z1]
	4148  R350 [Radeon 9800/9800 SE]
	4150  RV350 [Radeon 9550/9600/X1050 Series]
	4151  RV350 [Radeon 9600 Series]
	4152  RV360 [Radeon 9600/X1050 Series]
	4153  RV350 [Radeon 9550]
	4154  RV350 GL [FireGL T2]
	4155  RV350 [Radeon 9600]
	4157  RV350 GL [FireGL T2]
	4158  68800AX [Graphics Ultra Pro PCI]
	4164  R300 [Radeon 9500 PRO] (Secondary)
	4165  R300 [Radeon 9700 PRO] (Secondary)
	4166  R300 [Radeon 9700 PRO] (Secondary)
	4168  RV350 [Radeon 9800 SE] (Secondary)
	4170  RV350 [Radeon 9550/9600/X1050 Series] (Secondary)
	4171  RV350 [Radeon 9600] (Secondary)
	4172  RV350 [Radeon 9600/X1050 Series] (Secondary)
	4173  RV350 [Radeon 9550] (Secondary)
	4242  R200 [All-In-Wonder Radeon 8500 DV]
	4243  R200 PCI Bridge [All-in-Wonder Radeon 8500DV]
	4336  RS100 [Radeon IGP 320M]
	4337  RS200M [Radeon IGP 330M/340M/345M/350M]
	4341  IXP150 AC'97 Audio Controller
	4342  IXP200 3COM 3C920B Ethernet Controller
	4345  EHCI USB Controller
	4347  OHCI USB Controller #1
	4348  OHCI USB Controller #2
	4349  Dual Channel Bus Master PCI IDE Controller
	434d  IXP AC'97 Modem
	4353  SMBus
	4354  215CT [Mach64 CT PCI]
	4358  210888CX [Mach64 CX]
	4361  IXP SB300 AC'97 Audio Controller
	4363  SMBus
	436e  436E Serial ATA Controller
	4370  IXP SB400 AC'97 Audio Controller
	4371  IXP SB4x0 PCI-PCI Bridge
	4372  IXP SB4x0 SMBus Controller
	4373  IXP SB4x0 USB2 Host Controller
	4374  IXP SB4x0 USB Host Controller
	4375  IXP SB4x0 USB Host Controller
	4376  IXP SB4x0 IDE Controller
	4377  IXP SB4x0 PCI-ISA Bridge
	4378  IXP SB400 AC'97 Modem Controller
	4379  IXP SB4x0 Serial ATA Controller
	437a  IXP SB400 Serial ATA Controller
	437b  IXP SB4x0 High Definition Audio Controller
	4380  SB600 Non-Raid-5 SATA
	4381  SB600 SATA Controller (RAID 5 mode)
	4382  SB600 AC97 Audio
	4383  SBx00 Azalia (Intel HDA)
	4384  SBx00 PCI to PCI Bridge
	4385  SBx00 SMBus Controller
	4386  SB600 USB Controller (EHCI)
	4387  SB600 USB (OHCI0)
	4388  SB600 USB (OHCI1)
	4389  SB600 USB (OHCI2)
	438a  SB600 USB (OHCI3)
	438b  SB600 USB (OHCI4)
	438c  SB600 IDE
	438d  SB600 PCI to LPC Bridge
	438e  SB600 AC97 Modem
	4390  SB7x0/SB8x0/SB9x0 SATA Controller [IDE mode]
	4391  SB7x0/SB8x0/SB9x0 SATA Controller [AHCI mode]
	4392  SB7x0/SB8x0/SB9x0 SATA Controller [Non-RAID5 mode]
	4393  SB7x0/SB8x0/SB9x0 SATA Controller [RAID5 mode]
	4394  SB7x0/SB8x0/SB9x0 SATA Controller [AHCI mode]
	4395  SB8x0/SB9x0 SATA Controller [Storage mode]
	4396  SB7x0/SB8x0/SB9x0 USB EHCI Controller
	4397  SB7x0/SB8x0/SB9x0 USB OHCI0 Controller
	4398  SB7x0 USB OHCI1 Controller
	4399  SB7x0/SB8x0/SB9x0 USB OHCI2 Controller
	439c  SB7x0/SB8x0/SB9x0 IDE Controller
	439d  SB7x0/SB8x0/SB9x0 LPC host controller
	43a0  SB700/SB800/SB900 PCI to PCI bridge (PCIE port 0)
	43a1  SB700/SB800/SB900 PCI to PCI bridge (PCIE port 1)
	43a2  SB900 PCI to PCI bridge (PCIE port 2)
	43a3  SB900 PCI to PCI bridge (PCIE port 3)
	4437  RS250 [Mobility Radeon 7000 IGP]
	4554  210888ET [Mach64 ET]
	4654  Mach64 VT
	4742  3D Rage PRO AGP 2X
	4744  3D Rage PRO AGP 1X
	4749  3D Rage PRO PCI
	474d  Rage XL AGP 2X
	474e  Rage XC AGP
	474f  Rage XL
	4750  3D Rage Pro PCI
	4752  Rage XL PCI
	4753  Rage XC
	4754  3D Rage II/II+ PCI [Mach64 GT]
	4755  Mach64 GTB [3D Rage II+ DVD]
	4756  3D Rage IIC PCI [Mach64 GT IIC]
	4757  3D Rage IIC AGP
	4758  210888GX [Mach64 GX PCI]
	4759  3D Rage IIC PCI
	475a  3D Rage IIC AGP
	4966  RV250 [Radeon 9000 Series]
	496e  RV250 [Radeon 9000] (Secondary)
	4a49  R420 [Radeon X800 PRO/GTO AGP]
	4a4a  R420 [Radeon X800 GT AGP]
	4a4b  R420 [Radeon X800 AGP Series]
	4a4d  R420 GL [FireGL X3-256]
	4a4e  RV420/M18 [Mobility Radeon 9800]
	4a4f  R420 [Radeon X850 AGP]
	4a50  R420 [Radeon X800 XT Platinum Edition AGP]
	4a54  R420 [Radeon X800 VE AGP]
	4a69  R420 [Radeon X800 PRO/GTO] (Secondary)
	4a6a  R420 [Radeon X800] (Secondary)
	4a6b  R420 [Radeon X800 XT AGP] (Secondary)
	4a70  R420 [Radeon X800 XT Platinum Edition AGP] (Secondary)
	4a74  R420 [Radeon X800 VE] (Secondary)
	4b49  R481 [Radeon X850 XT AGP]
	4b4b  R481 [Radeon X850 PRO AGP]
	4b4c  R481 [Radeon X850 XT Platinum Edition AGP]
	4b69  R481 [Radeon X850 XT AGP] (Secondary)
	4b6b  R481 [Radeon X850 PRO AGP] (Secondary)
	4b6c  R481 [Radeon X850 XT Platinum Edition AGP] (Secondary)
	4c42  3D Rage LT PRO AGP 2X
	4c46  Rage Mobility 128 AGP 2X/Mobility M3
	4c47  3D Rage IIC PCI / Mobility Radeon 7500/7500C
	4c49  3D Rage LT PRO PCI
	4c4d  Rage Mobility AGP 2x Series
	4c50  3D Rage LT PRO PCI
	4c52  Rage Mobility-M1 PCI
	4c54  264LT [Mach64 LT]
	4c57  RV200/M7 [Mobility Radeon 7500]
	4c58  RV200/M7 GL [Mobility FireGL 7800]
	4c59  RV100/M6 [Rage/Radeon Mobility Series]
	4c66  RV250/M9 GL [Mobility FireGL 9000/Radeon 9000]
	4c6e  RV250/M9 [Mobility Radeon 9000] (Secondary)
	4d46  Rage Mobility 128 AGP 4X/Mobility M4
	4d52  Theater 550 PRO PCI [ATI TV Wonder 550]
	4d53  Theater 550 PRO PCIe
	4e44  R300 [Radeon 9700/9700 PRO]
	4e45  R300 [Radeon 9500 PRO/9700]
	4e46  R300 [Radeon 9600 TX]
	4e47  R300 GL [FireGL X1]
	4e48  R350 [Radeon 9800 Series]
	4e49  R350 [Radeon 9800]
	4e4a  R360 [Radeon 9800 XXL/XT]
	4e4b  R350 GL [FireGL X2 AGP Pro]
	4e50  RV350/M10 / RV360/M11 [Mobility Radeon 9600 (PRO) / 9700]
	4e51  RV350 [Radeon 9550/9600/X1050 Series]
	4e52  RV350/M10 [Mobility Radeon 9500/9700 SE]
	4e54  RV350/M10 GL [Mobility FireGL T2]
	4e56  RV360/M12 [Mobility Radeon 9550]
	4e64  R300 [Radeon 9700 PRO] (Secondary)
	4e65  R300 [Radeon 9500 PRO] (Secondary)
	4e66  RV350 [Radeon 9600] (Secondary)
	4e67  R300 GL [FireGL X1] (Secondary)
	4e68  R350 [Radeon 9800 PRO] (Secondary)
	4e69  R350 [Radeon 9800] (Secondary)
	4e6a  RV350 [Radeon 9800 XT] (Secondary)
	4e71  RV350/M10 [Mobility Radeon 9600] (Secondary)
	4f72  RV250 [Radeon 9000 Series]
	4f73  RV250 [Radeon 9000 Series] (Secondary)
	5044  All-In-Wonder 128 PCI
	5046  Rage 128 PRO AGP 4x TMDS
	5050  Rage128 [Xpert 128 PCI]
	5052  Rage 128 PRO AGP 4X TMDS
	5144  R100 [Radeon 7200 / All-In-Wonder Radeon]
	5148  R200 GL [FireGL 8800]
	514c  R200 [Radeon 8500/8500 LE]
	514d  R200 [Radeon 9100]
	5157  RV200 [Radeon 7500/7500 LE]
	5159  RV100 [Radeon 7000 / Radeon VE]
	515e  ES1000
	5245  Rage 128 GL PCI
	5246  Rage Fury/Xpert 128/Xpert 2000 AGP 2x
	524b  Rage 128 VR PCI
	524c  Rage 128 VR AGP
	5346  Rage 128 SF/4x AGP 2x
	534d  Rage 128 4X AGP 4x
	5354  Mach 64 VT
	5446  Rage 128 PRO Ultra AGP 4x
	5452  Rage 128 PRO Ultra4XL VR-R AGP
	5460  RV370/M22 [Mobility Radeon X300]
	5461  RV370/M22 [Mobility Radeon X300]
	5462  RV380/M24C [Mobility Radeon X600 SE]
	5464  RV370/M22 GL [Mobility FireGL V3100]
	5549  R423 [Radeon X800 GTO]
	554a  R423 [Radeon X800 XT Platinum Edition]
	554b  R423 [Radeon X800 GT/SE]
	554d  R430 [Radeon X800 XL]
	554e  R430 [All-In-Wonder X800 GT]
	554f  R430 [Radeon X800]
	5550  R423 GL [FireGL V7100]
	5551  R423 GL [FireGL V5100]
	5569  R423 [Radeon X800 PRO] (Secondary)
	556b  R423 [Radeon X800 GT] (Secondary)
	556d  R430 [Radeon X800 XL] (Secondary)
	556f  R430 [Radeon X800] (Secondary)
	5571  R423 GL [FireGL V5100] (Secondary)
	564b  RV410/M26 GL [Mobility FireGL V5000]
	564f  RV410/M26 [Mobility Radeon X700 XL]
	5652  RV410/M26 [Mobility Radeon X700]
	5653  RV410/M26 [Mobility Radeon X700]
	5654  264VT [Mach64 VT]
	5655  264VT3 [Mach64 VT3]
	5656  264VT4 [Mach64 VT4]
	5657  RV410 [Radeon X550 XTX / X700]
	5830  RS300 Host Bridge
	5831  RS300 Host Bridge
	5832  RS300 Host Bridge
	5833  RS300 Host Bridge
	5834  RS300 [Radeon 9100 IGP]
	5835  RS300M [Mobility Radeon 9100 IGP]
	5838  RS300 AGP Bridge
	5854  RS480 [Radeon Xpress 200 Series] (Secondary)
	5874  RS480 [Radeon Xpress 1150] (Secondary)
	5940  RV280 [Radeon 9200 PRO] (Secondary)
	5941  RV280 [Radeon 9200] (Secondary)
	5944  RV280 [Radeon 9200 SE PCI]
	5950  RS480/RS482/RS485 Host Bridge
	5951  RX480/RX482 Host Bridge
	5952  RD580 Host Bridge
	5954  RS480 [Radeon Xpress 200 Series]
	5955  RS480M [Mobility Radeon Xpress 200]
	5956  RD790 Host Bridge
	5957  RX780/RX790 Host Bridge
	5958  RD780 Host Bridge
	5960  RV280 [Radeon 9200 PRO]
	5961  RV280 [Radeon 9200]
	5962  RV280 [Radeon 9200]
	5964  RV280 [Radeon 9200 SE]
	5965  RV280 GL [FireMV 2200 PCI]
	5974  RS482/RS485 [Radeon Xpress 1100/1150]
	5975  RS482M [Mobility Radeon Xpress 200]
	5978  RX780/RD790 PCI to PCI bridge (external gfx0 port A)
	5979  RD790 PCI to PCI bridge (external gfx0 port B)
	597a  RD790 PCI to PCI bridge (PCI express gpp port A)
	597b  RX780/RD790 PCI to PCI bridge (PCI express gpp port B)
	597c  RD790 PCI to PCI bridge (PCI express gpp port C)
	597d  RX780/RD790 PCI to PCI bridge (PCI express gpp port D)
	597e  RD790 PCI to PCI bridge (PCI express gpp port E)
	597f  RD790 PCI to PCI bridge (PCI express gpp port F)
	5980  RD790 PCI to PCI bridge (external gfx1 port A)
	5981  RD790 PCI to PCI bridge (external gfx1 port B)
	5982  RD790 PCI to PCI bridge (NB-SB link)
	5a10  RD890 Northbridge only dual slot (2x16) PCI-e GFX Hydra part
	5a11  RD890 Northbridge only single slot PCI-e GFX Hydra part
	5a12  RD890 Northbridge only dual slot (2x8) PCI-e GFX Hydra part
	5a13  RD890S/SR5650 Host Bridge
	5a14  RD9x0/RX980 Host Bridge
	5a15  RD890 PCI to PCI bridge (PCI express gpp port A)
	5a16  RD890/RD9x0/RX980 PCI to PCI bridge (PCI Express GFX port 0)
	5a17  RD890/RD9x0 PCI to PCI bridge (PCI Express GFX port 1)
	5a18  RD890/RD9x0/RX980 PCI to PCI bridge (PCI Express GPP Port 0)
	5a19  RD890/RD9x0/RX980 PCI to PCI bridge (PCI Express GPP Port 1)
	5a1a  RD890/RD9x0/RX980 PCI to PCI bridge (PCI Express GPP Port 2)
	5a1b  RD890/RD9x0/RX980 PCI to PCI bridge (PCI Express GPP Port 3)
	5a1c  RD890/RD9x0/RX980 PCI to PCI bridge (PCI Express GPP Port 4)
	5a1d  RD890/RD9x0/RX980 PCI to PCI bridge (PCI Express GPP Port 5)
	5a1e  RD890/RD9x0/RX980 PCI to PCI bridge (PCI Express GPP2 Port 0)
	5a1f  RD890/RD990 PCI to PCI bridge (PCI Express GFX2 port 0)
	5a20  RD890/RD990 PCI to PCI bridge (PCI Express GFX2 port 1)
	5a23  RD890S/RD990 I/O Memory Management Unit (IOMMU)
	5a31  RC410 Host Bridge
	5a33  RS400 Host Bridge
	5a34  RS4xx PCI Express Port [ext gfx]
	5a36  RC4xx/RS4xx PCI Express Port 1
	5a37  RC4xx/RS4xx PCI Express Port 2
	5a38  RC4xx/RS4xx PCI Express Port 3
	5a39  RC4xx/RS4xx PCI Express Port 4
	5a3f  RC4xx/RS4xx PCI Bridge [int gfx]
	5a41  RS400 [Radeon Xpress 200]
	5a42  RS400M [Radeon Xpress 200M]
	5a61  RC410 [Radeon Xpress 200/1100]
	5a62  RC410M [Mobility Radeon Xpress 200M]
	5b60  RV370 [Radeon X300]
	5b62  RV370 [Radeon X600/X600 SE]
	5b63  RV370 [Radeon X300/X550/X1050 Series]
	5b64  RV370 GL [FireGL V3100]
	5b65  RV370 GL [FireMV 2200]
	5b66  RV370X
	5b70  RV370 [Radeon X300 SE]
	5b72  RV380 [Radeon X300/X550/X1050 Series] (Secondary)
	5b73  RV370 [Radeon X300/X550/X1050 Series] (Secondary)
	5b74  RV370 GL [FireGL V3100] (Secondary)
	5b75  RV370 GL [FireMV 2200] (Secondary)
	5c61  RV280/M9+ [Mobility Radeon 9200 AGP]
	5c63  RV280/M9+ [Mobility Radeon 9200 AGP]
	5d44  RV280 [Radeon 9200 SE] (Secondary)
	5d45  RV280 GL [FireMV 2200 PCI] (Secondary)
	5d48  R423/M28 [Mobility Radeon X800 XT]
	5d49  R423/M28 GL [Mobility FireGL V5100]
	5d4a  R423/M28 [Mobility Radeon X800]
	5d4d  R480 [Radeon X850 XT Platinum Edition]
	5d4e  R480 [Radeon X850 SE]
	5d4f  R480 [Radeon X800 GTO]
	5d50  R480 GL [FireGL V7200]
	5d52  R480 [Radeon X850 XT]
	5d57  R423 [Radeon X800 XT]
	5d6d  R480 [Radeon X850 XT Platinum Edition] (Secondary)
	5d6f  R480 [Radeon X800 GTO] (Secondary)
	5d72  R480 [Radeon X850 XT] (Secondary)
	5d77  R423 [Radeon X800 XT] (Secondary)
	5e48  RV410 GL [FireGL V5000]
	5e49  RV410 [Radeon X700 Series]
	5e4a  RV410 [Radeon X700 XT]
	5e4b  RV410 [Radeon X700 PRO]
	5e4c  RV410 [Radeon X700 SE]
	5e4d  RV410 [Radeon X700]
	5e4f  RV410 [Radeon X700]
	5e6b  RV410 [Radeon X700 PRO] (Secondary)
	5e6d  RV410 [Radeon X700] (Secondary)
	5f57  R423 [Radeon X800 XT]
	6600  Mars [Radeon HD 8670A/8670M/8750M]
	6601  Mars [Radeon HD 8730M]
	6602  Mars
	6603  Mars
	6604  Opal XT [Radeon R7 M265]
	6605  Opal PRO [Radeon R7 M260]
	6606  Mars XTX [Radeon HD 8790M]
	6607  Mars LE [Radeon HD 8530M / R5 M240]
	6608  Oland GL [FirePro W2100]
	6610  Oland XT [Radeon HD 8670 / R7 250/350]
	6611  Oland [Radeon HD 8570 / R7 240/340 OEM]
	6613  Oland PRO [Radeon R7 240/340]
	6620  Mars
	6621  Mars PRO
	6623  Mars
	6631  Oland
	6640  Saturn XT [FirePro M6100]
	6641  Saturn PRO [Radeon HD 8930M]
	6646  Bonaire XT [Radeon R9 M280X]
	6647  Bonaire PRO [Radeon R9 M270X]
	6649  Bonaire [FirePro W5100]
	6650  Bonaire
	6651  Bonaire
	6658  Bonaire XTX [Radeon R7 260X/360]
	665c  Bonaire XT [Radeon HD 7790/8770 / R7 360 / R9 260/360 OEM]
	665d  Bonaire [Radeon R7 200 Series]
	665f  Tobago PRO [Radeon R7 360 / R9 360 OEM]
	6660  Sun XT [Radeon HD 8670A/8670M/8690M / R5 M330 / M430 / R7 M520]
	6663  Sun PRO [Radeon HD 8570A/8570M]
	6664  Jet XT [Radeon R5 M240]
	6665  Jet PRO [Radeon R5 M230]
	6667  Jet ULT [Radeon R5 M230]
	666f  Sun LE [Radeon HD 8550M / R5 M230]
	6704  Cayman PRO GL [FirePro V7900]
	6707  Cayman LE GL [FirePro V5900]
	6718  Cayman XT [Radeon HD 6970]
	6719  Cayman PRO [Radeon HD 6950]
	671c  Antilles [Radeon HD 6990]
	671d  Antilles [Radeon HD 6990]
	671f  Cayman CE [Radeon HD 6930]
	6720  Blackcomb [Radeon HD 6970M/6990M]
	6738  Barts XT [Radeon HD 6870]
	6739  Barts PRO [Radeon HD 6850]
	673e  Barts LE [Radeon HD 6790]
	6740  Whistler [Radeon HD 6730M/6770M/7690M XT]
	6741  Whistler [Radeon HD 6630M/6650M/6750M/7670M/7690M]
	6742  Whistler LE [Radeon HD 6610M/7610M]
	6743  Whistler [Radeon E6760]
	6749  Turks GL [FirePro V4900]
	674a  Turks GL [FirePro V3900]
	6750  Onega [Radeon HD 6650A/7650A]
	6751  Turks [Radeon HD 7650A/7670A]
	6758  Turks XT [Radeon HD 6670/7670]
	6759  Turks PRO [Radeon HD 6570/7570/8550]
	675b  Turks [Radeon HD 7600 Series]
	675d  Turks PRO [Radeon HD 7570]
	675f  Turks LE [Radeon HD 5570/6510/7510/8510]
	6760  Seymour [Radeon HD 6400M/7400M Series]
	6761  Seymour LP [Radeon HD 6430M]
	6763  Seymour [Radeon E6460]
	6764  Seymour [Radeon HD 6400M Series]
	6765  Seymour [Radeon HD 6400M Series]
	6766  Caicos
	6767  Caicos
	6768  Caicos
	6770  Caicos [Radeon HD 6450A/7450A]
	6771  Caicos XTX [Radeon HD 8490 / R5 235X OEM]
	6772  Caicos [Radeon HD 7450A]
	6778  Caicos XT [Radeon HD 7470/8470 / R5 235/310 OEM]
	6779  Caicos [Radeon HD 6450/7450/8450 / R5 230 OEM]
	677b  Caicos PRO [Radeon HD 7450]
	6780  Tahiti XT GL [FirePro W9000]
	6784  Tahiti [FirePro Series Graphics Adapter]
	6788  Tahiti [FirePro Series Graphics Adapter]
	678a  Tahiti PRO GL [FirePro Series]
	6798  Tahiti XT [Radeon HD 7970/8970 OEM / R9 280X]
	679a  Tahiti PRO [Radeon HD 7950/8950 OEM / R9 280]
	679b  Malta [Radeon HD 7990/8990 OEM]
	679e  Tahiti LE [Radeon HD 7870 XT]
	679f  Tahiti
	67a0  Hawaii XT GL [FirePro W9100]
	67a1  Hawaii PRO GL [FirePro W8100]
	67a2  Hawaii GL
	67a8  Hawaii
	67a9  Hawaii
	67aa  Hawaii
	67b0  Hawaii XT / Grenada XT [Radeon R9 290X/390X]
	67b1  Hawaii PRO [Radeon R9 290/390]
	67b9  Vesuvius [Radeon R9 295X2]
	67be  Hawaii LE
	67c0  Ellesmere [Radeon Pro WX 7100]
	67c4  Ellesmere [Radeon Pro WX 7100]
	67c7  Ellesmere [Radeon Pro WX 5100]
	67ca  Ellesmere [Polaris10]
	67cc  Ellesmere [Polaris10]
	67cf  Ellesmere [Polaris10]
	67df  Ellesmere [Radeon RX 470/480/570/580]
	67e0  Baffin [Polaris11]
	67e1  Baffin [Polaris11]
	67e3  Baffin [Radeon Pro WX 4100]
	67e8  Baffin [Polaris11]
	67e9  Baffin [Polaris11]
	67eb  Baffin [Polaris11]
	67ef  Baffin [Radeon RX 460/560D / Pro 450/455/460/560]
	67ff  Baffin [Radeon RX 560]
	6800  Wimbledon XT [Radeon HD 7970M]
	6801  Neptune XT [Radeon HD 8970M]
	6802  Wimbledon
	6806  Neptune
	6808  Pitcairn XT GL [FirePro W7000]
	6809  Pitcairn LE GL [FirePro W5000]
	6810  Curacao XT / Trinidad XT [Radeon R7 370 / R9 270X/370X]
	6811  Curacao PRO [Radeon R7 370 / R9 270/370 OEM]
	6816  Pitcairn
	6817  Pitcairn
	6818  Pitcairn XT [Radeon HD 7870 GHz Edition]
	6819  Pitcairn PRO [Radeon HD 7850 / R7 265 / R9 270 1024SP]
	6820  Venus XTX [Radeon HD 8890M / R9 M275X/M375X]
	6821  Venus XT [Radeon HD 8870M / R9 M270X/M370X]
	6822  Venus PRO [Radeon E8860]
	6823  Venus PRO [Radeon HD 8850M / R9 M265X]
	6825  Heathrow XT [Radeon HD 7870M]
	6826  Chelsea LP [Radeon HD 7700M Series]
	6827  Heathrow PRO [Radeon HD 7850M/8850M]
	6828  Cape Verde PRO [FirePro W600]
	6829  Cape Verde
	682a  Venus PRO
	682b  Venus LE [Radeon HD 8830M]
	682c  Cape Verde GL [FirePro W4100]
	682d  Chelsea XT GL [FirePro M4000]
	682f  Chelsea LP [Radeon HD 7730M]
	6830  Cape Verde [Radeon HD 7800M Series]
	6831  Cape Verde [AMD Radeon HD 7700M Series]
	6835  Cape Verde PRX [Radeon R9 255 OEM]
	6837  Cape Verde LE [Radeon HD 7730/8730]
	683d  Cape Verde XT [Radeon HD 7770/8760 / R7 250X]
	683f  Cape Verde PRO [Radeon HD 7750/8740 / R7 250E]
	6840  Thames [Radeon HD 7500M/7600M Series]
	6841  Thames [Radeon HD 7550M/7570M/7650M]
	6842  Thames LE [Radeon HD 7000M Series]
	6843  Thames [Radeon HD 7670M]
	6863  Vega 10 XTX [Radeon Vega Frontier Edition]
	687f  Vega 10 XT [Radeon RX Vega 64]
	6888  Cypress XT [FirePro V8800]
	6889  Cypress PRO [FirePro V7800]
	688a  Cypress XT [FirePro V9800]
	688c  Cypress XT GL [FireStream 9370]
	688d  Cypress PRO GL [FireStream 9350]
	6898  Cypress XT [Radeon HD 5870]
	6899  Cypress PRO [Radeon HD 5850]
	689b  Cypress PRO [Radeon HD 6800 Series]
	689c  Hemlock [Radeon HD 5970]
	689d  Hemlock [Radeon HD 5970]
	689e  Cypress LE [Radeon HD 5830]
	68a0  Broadway XT [Mobility Radeon HD 5870]
	68a1  Broadway PRO [Mobility Radeon HD 5850]
	68a8  Granville [Radeon HD 6850M/6870M]
	68a9  Juniper XT [FirePro V5800]
	68b8  Juniper XT [Radeon HD 5770]
	68b9  Juniper LE [Radeon HD 5670 640SP Edition]
	68ba  Juniper XT [Radeon HD 6770]
	68be  Juniper PRO [Radeon HD 5750]
	68bf  Juniper PRO [Radeon HD 6750]
	68c0  Madison [Mobility Radeon HD 5730 / 6570M]
	68c1  Madison [Mobility Radeon HD 5650/5750 / 6530M/6550M]
	68c7  Madison [Mobility Radeon HD 5570/6550A]
	68c8  Redwood XT GL [FirePro V4800]
	68c9  Redwood PRO GL [FirePro V3800]
	68d8  Redwood XT [Radeon HD 5670/5690/5730]
	68d9  Redwood PRO [Radeon HD 5550/5570/5630/6510/6610/7570]
	68da  Redwood LE [Radeon HD 5550/5570/5630/6390/6490/7570]
	68de  Redwood
	68e0  Park [Mobility Radeon HD 5430/5450/5470]
	68e1  Park [Mobility Radeon HD 5430]
	68e4  Robson CE [Radeon HD 6370M/7370M]
	68e5  Robson LE [Radeon HD 6330M]
	68e8  Cedar
	68e9  Cedar [ATI FirePro (FireGL) Graphics Adapter]
	68f1  Cedar GL [FirePro 2460]
	68f2  Cedar GL [FirePro 2270]
	68f8  Cedar [Radeon HD 7300 Series]
	68f9  Cedar [Radeon HD 5000/6000/7350/8350 Series]
	68fa  Cedar [Radeon HD 7350/8350 / R5 220]
	68fe  Cedar LE
	6900  Topaz XT [Radeon R7 M260/M265 / M340/M360 / M440/M445]
	6901  Topaz PRO [Radeon R5 M255]
	6907  Meso XT [Radeon R5 M315]
	6921  Amethyst XT [Radeon R9 M295X]
	6929  Tonga XT GL [FirePro S7150]
	692b  Tonga PRO GL [FirePro W7100]
	692f  Tonga XTV GL [FirePro S7150V]
	6938  Tonga XT / Amethyst XT [Radeon R9 380X / R9 M295X]
	6939  Tonga PRO [Radeon R9 285/380]
	6980  Polaris12
	6981  Polaris12
	6985  Lexa XT [Radeon PRO WX 3100]
	6986  Polaris12
	6987  Polaris12
	6995  Lexa XT [Radeon PRO WX 2100]
	699f  Lexa PRO [Radeon RX 550]
	700f  RS100 AGP Bridge
	7010  RS200/RS250 AGP Bridge
	7100  R520 [Radeon X1800 XT]
	7101  R520/M58 [Mobility Radeon X1800 XT]
	7102  R520/M58 [Mobility Radeon X1800]
	7104  R520 GL [FireGL V7200]
	7109  R520 [Radeon X1800 XL]
	710a  R520 [Radeon X1800 GTO]
	710b  R520 [Radeon X1800 GTO]
	7120  R520 [Radeon X1800] (Secondary)
	7124  R520 GL [FireGL V7200] (Secondary)
	7129  R520 [Radeon X1800] (Secondary)
	7140  RV515 [Radeon X1300/X1550/X1600 Series]
	7142  RV515 PRO [Radeon X1300/X1550 Series]
	7143  RV505 [Radeon X1300/X1550 Series]
	7145  RV515/M54 [Mobility Radeon X1400]
	7146  RV515 [Radeon X1300/X1550]
	7147  RV505 [Radeon X1550 64-bit]
	7149  RV515/M52 [Mobility Radeon X1300]
	714a  RV515/M52 [Mobility Radeon X1300]
	7152  RV515 GL [FireGL V3300]
	7153  RV515 GL [FireGL V3350]
	715f  RV505 CE [Radeon X1550 64-bit]
	7162  RV515 PRO [Radeon X1300/X1550 Series] (Secondary)
	7163  RV505 [Radeon X1550 Series] (Secondary)
	7166  RV515 [Radeon X1300/X1550 Series] (Secondary)
	7167  RV515 [Radeon X1550 64-bit] (Secondary)
	7172  RV515 GL [FireGL V3300] (Secondary)
	7173  RV515 GL [FireGL V3350] (Secondary)
	7181  RV516 [Radeon X1600/X1650 Series]
	7183  RV516 [Radeon X1300/X1550 Series]
	7186  RV516/M64 [Mobility Radeon X1450]
	7187  RV516 [Radeon X1300/X1550 Series]
	7188  RV516/M64-S [Mobility Radeon X2300]
	718a  RV516/M64 [Mobility Radeon X2300]
	718b  RV516/M62 [Mobility Radeon X1350]
	718c  RV516/M62-CSP64 [Mobility Radeon X1350]
	718d  RV516/M64-CSP128 [Mobility Radeon X1450]
	7193  RV516 [Radeon X1550 Series]
	7196  RV516/M62-S [Mobility Radeon X1350]
	719b  RV516 GL [FireMV 2250]
	719f  RV516 [Radeon X1550 Series]
	71a0  RV516 [Radeon X1300/X1550 Series] (Secondary)
	71a1  RV516 [Radeon X1600/X1650 Series] (Secondary)
	71a3  RV516 [Radeon X1300/X1550 Series] (Secondary)
	71a7  RV516 [Radeon X1300/X1550 Series] (Secondary)
	71bb  RV516 GL [FireMV 2250] (Secondary)
	71c0  RV530 [Radeon X1600 XT/X1650 GTO]
	71c1  RV535 [Radeon X1650 PRO]
	71c2  RV530 [Radeon X1600 PRO]
	71c4  RV530/M56 GL [Mobility FireGL V5200]
	71c5  RV530/M56-P [Mobility Radeon X1600]
	71c6  RV530LE [Radeon X1600/X1650 PRO]
	71c7  RV535 [Radeon X1650 PRO]
	71ce  RV530 [Radeon X1300 XT/X1600 PRO]
	71d2  RV530 GL [FireGL V3400]
	71d4  RV530/M66 GL [Mobility FireGL V5250]
	71d5  RV530/M66-P [Mobility Radeon X1700]
	71d6  RV530/M66-XT [Mobility Radeon X1700]
	71de  RV530/M66 [Mobility Radeon X1700/X2500]
	71e0  RV530 [Radeon X1600] (Secondary)
	71e1  RV535 [Radeon X1650 PRO] (Secondary)
	71e2  RV530 [Radeon X1600] (Secondary)
	71e6  RV530 [Radeon X1650] (Secondary)
	71e7  RV535 [Radeon X1650 PRO] (Secondary)
	71f2  RV530 GL [FireGL V3400] (Secondary)
	7210  RV550/M71 [Mobility Radeon HD 2300]
	7211  RV550/M71 [Mobility Radeon X2300 HD]
	7240  R580+ [Radeon X1950 XTX]
	7244  R580+ [Radeon X1950 XT]
	7248  R580 [Radeon X1950]
	7249  R580 [Radeon X1900 XT]
	724b  R580 [Radeon X1900 GT]
	724e  R580 GL [FireGL V7350]
	7269  R580 [Radeon X1900 XT] (Secondary)
	726b  R580 [Radeon X1900 GT] (Secondary)
	726e  R580 [AMD Stream Processor] (Secondary)
	7280  RV570 [Radeon X1950 PRO]
	7288  RV570 [Radeon X1950 GT]
	7291  RV560 [Radeon X1650 XT]
	7293  RV560 [Radeon X1650 GT]
	72a0  RV570 [Radeon X1950 PRO] (Secondary)
	72a8  RV570 [Radeon X1950 GT] (Secondary)
	72b1  RV560 [Radeon X1650 XT] (Secondary)
	72b3  RV560 [Radeon X1650 GT] (Secondary)
	7300  Fiji [Radeon R9 FURY / NANO Series]
	7833  RS350 Host Bridge
	7834  RS350 [Radeon 9100 PRO/XT IGP]
	7835  RS350M [Mobility Radeon 9000 IGP]
	7838  RS350 AGP Bridge
	7910  RS690 Host Bridge
	7911  RS690/RS740 Host Bridge
	7912  RS690/RS740 PCI to PCI Bridge (Internal gfx)
	7913  RS690 PCI to PCI Bridge (PCI Express Graphics Port 0)
	7915  RS690 PCI to PCI Bridge (PCI Express Port 1)
	7916  RS690 PCI to PCI Bridge (PCI Express Port 2)
	7917  RS690 PCI to PCI Bridge (PCI Express Port 3)
	7919  RS690 HDMI Audio [Radeon Xpress 1200 Series]
	791e  RS690 [Radeon X1200]
	791f  RS690M [Radeon Xpress 1200/1250/1270]
	7930  RS600 Host Bridge
	7932  RS600 PCI to PCI Bridge (Internal gfx)
	7933  RS600 PCI to PCI Bridge (PCI Express Graphics Port 0)
	7935  RS600 PCI to PCI Bridge (PCI Express Port 1)
	7936  RS600 PCI to PCI Bridge (PCI Express Port 2)
	7937  RS690 PCI to PCI Bridge (PCI Express Port 3)
	793b  RS600 HDMI Audio [Radeon Xpress 1250]
	793f  RS690M [Radeon Xpress 1200/1250/1270] (Secondary)
	7941  RS600 [Radeon Xpress 1250]
	7942  RS600M [Radeon Xpress 1250]
	796e  RS740 [Radeon 2100]
	9400  R600 [Radeon HD 2900 PRO/XT]
	9401  R600 [Radeon HD 2900 XT]
	9403  R600 [Radeon HD 2900 PRO]
	9405  R600 [Radeon HD 2900 GT]
	940a  R600 GL [FireGL V8650]
	940b  R600 GL [FireGL V8600]
	940f  R600 GL [FireGL V7600]
	9440  RV770 [Radeon HD 4870]
	9441  R700 [Radeon HD 4870 X2]
	9442  RV770 [Radeon HD 4850]
	9443  R700 [Radeon HD 4850 X2]
	9444  RV770 GL [FirePro V8750]
	9446  RV770 GL [FirePro V7760]
	944a  RV770/M98L [Mobility Radeon HD 4850]
	944b  RV770/M98 [Mobility Radeon HD 4850 X2]
	944c  RV770 LE [Radeon HD 4830]
	944e  RV770 CE [Radeon HD 4710]
	9450  RV770 GL [FireStream 9270]
	9452  RV770 GL [FireStream 9250]
	9456  RV770 GL [FirePro V8700]
	945a  RV770/M98-XT [Mobility Radeon HD 4870]
	9460  RV790 [Radeon HD 4890]
	9462  RV790 [Radeon HD 4860]
	946a  RV770 GL [FirePro M7750]
	9480  RV730/M96 [Mobility Radeon HD 4650/5165]
	9488  RV730/M96-XT [Mobility Radeon HD 4670]
	9489  RV730/M96 GL [Mobility FireGL V5725]
	9490  RV730 XT [Radeon HD 4670]
	9491  RV730/M96-CSP [Radeon E4690]
	9495  RV730 [Radeon HD 4600 AGP Series]
	9498  RV730 PRO [Radeon HD 4650]
	949c  RV730 GL [FirePro V7750]
	949e  RV730 GL [FirePro V5700]
	949f  RV730 GL [FirePro V3750]
	94a0  RV740/M97 [Mobility Radeon HD 4830]
	94a1  RV740/M97-XT [Mobility Radeon HD 4860]
	94a3  RV740/M97 GL [FirePro M7740]
	94b3  RV740 PRO [Radeon HD 4770]
	94b4  RV740 PRO [Radeon HD 4750]
	94c1  RV610 [Radeon HD 2400 PRO/XT]
	94c3  RV610 [Radeon HD 2400 PRO]
	94c4  RV610 LE [Radeon HD 2400 PRO AGP]
	94c5  RV610 [Radeon HD 2400 LE]
	94c7  RV610 [Radeon HD 2350]
	94c8  RV610/M74 [Mobility Radeon HD 2400 XT]
	94c9  RV610/M72-S [Mobility Radeon HD 2400]
	94cb  RV610 [Radeon E2400]
	94cc  RV610 LE [Radeon HD 2400 PRO PCI]
	9500  RV670 [Radeon HD 3850 X2]
	9501  RV670 [Radeon HD 3870]
	9504  RV670/M88 [Mobility Radeon HD 3850]
	9505  RV670 [Radeon HD 3690/3850]
	9506  RV670/M88 [Mobility Radeon HD 3850 X2]
	9507  RV670 [Radeon HD 3830]
	9508  RV670/M88-XT [Mobility Radeon HD 3870]
	9509  RV670/M88 [Mobility Radeon HD 3870 X2]
	950f  R680 [Radeon HD 3870 X2]
	9511  RV670 GL [FireGL V7700]
	9513  RV670 [Radeon HD 3850 X2]
	9515  RV670 PRO [Radeon HD 3850 AGP]
	9519  RV670 GL [FireStream 9170]
	9540  RV710 [Radeon HD 4550]
	954f  RV710 [Radeon HD 4350/4550]
	9552  RV710/M92 [Mobility Radeon HD 4330/4350/4550]
	9553  RV710/M92 [Mobility Radeon HD 4530/4570/545v]
	9555  RV710/M92 [Mobility Radeon HD 4350/4550]
	9557  RV711 GL [FirePro RG220]
	955f  RV710/M92 [Mobility Radeon HD 4330]
	9580  RV630 [Radeon HD 2600 PRO]
	9581  RV630/M76 [Mobility Radeon HD 2600]
	9583  RV630/M76 [Mobility Radeon HD 2600 XT/2700]
	9586  RV630 XT [Radeon HD 2600 XT AGP]
	9587  RV630 PRO [Radeon HD 2600 PRO AGP]
	9588  RV630 XT [Radeon HD 2600 XT]
	9589  RV630 PRO [Radeon HD 2600 PRO]
	958a  RV630 [Radeon HD 2600 X2]
	958b  RV630/M76 [Mobility Radeon HD 2600 XT]
	958c  RV630 GL [FireGL V5600]
	958d  RV630 GL [FireGL V3600]
	9591  RV635/M86 [Mobility Radeon HD 3650]
	9593  RV635/M86 [Mobility Radeon HD 3670]
	9595  RV635/M86 GL [Mobility FireGL V5700]
	9596  RV635 PRO [Radeon HD 3650 AGP]
	9597  RV635 PRO [Radeon HD 3650 AGP]
	9598  RV635 [Radeon HD 3650/3750/4570/4580]
	9599  RV635 PRO [Radeon HD 3650 AGP]
	95c0  RV620 PRO [Radeon HD 3470]
	95c2  RV620/M82 [Mobility Radeon HD 3410/3430]
	95c4  RV620/M82 [Mobility Radeon HD 3450/3470]
	95c5  RV620 LE [Radeon HD 3450]
	95c6  RV620 LE [Radeon HD 3450 AGP]
	95c9  RV620 LE [Radeon HD 3450 PCI]
	95cc  RV620 GL [FirePro V3700]
	95cd  RV620 [FirePro 2450]
	95cf  RV620 GL [FirePro 2260]
	960f  RS780 HDMI Audio [Radeon 3000/3100 / HD 3200/3300]
	9610  RS780 [Radeon HD 3200]
	9611  RS780C [Radeon 3100]
	9612  RS780M [Mobility Radeon HD 3200]
	9613  RS780MC [Mobility Radeon HD 3100]
	9614  RS780D [Radeon HD 3300]
	9616  RS780L [Radeon 3000]
	9640  BeaverCreek [Radeon HD 6550D]
	9641  BeaverCreek [Radeon HD 6620G]
	9642  Sumo [Radeon HD 6370D]
	9643  Sumo [Radeon HD 6380G]
	9644  Sumo [Radeon HD 6410D]
	9645  Sumo [Radeon HD 6410D]
	9647  BeaverCreek [Radeon HD 6520G]
	9648  Sumo [Radeon HD 6480G]
	9649  Sumo [Radeon HD 6480G]
	964a  BeaverCreek [Radeon HD 6530D]
	964b  Sumo
	964c  Sumo
	964e  Sumo
	964f  Sumo
	970f  RS880 HDMI Audio [Radeon HD 4200 Series]
	9710  RS880 [Radeon HD 4200]
	9712  RS880M [Mobility Radeon HD 4225/4250]
	9713  RS880M [Mobility Radeon HD 4100]
	9714  RS880 [Radeon HD 4290]
	9715  RS880 [Radeon HD 4250]
	9802  Wrestler [Radeon HD 6310]
	9803  Wrestler [Radeon HD 6310]
	9804  Wrestler [Radeon HD 6250]
	9805  Wrestler [Radeon HD 6250]
	9806  Wrestler [Radeon HD 6320]
	9807  Wrestler [Radeon HD 6290]
	9808  Wrestler [Radeon HD 7340]
	9809  Wrestler [Radeon HD 7310]
	980a  Wrestler [Radeon HD 7290]
	9830  Kabini [Radeon HD 8400 / R3 Series]
	9831  Kabini [Radeon HD 8400E]
	9832  Kabini [Radeon HD 8330]
	9833  Kabini [Radeon HD 8330E]
	9834  Kabini [Radeon HD 8210]
	9835  Kabini [Radeon HD 8310E]
	9836  Kabini [Radeon HD 8280 / R3 Series]
	9837  Kabini [Radeon HD 8280E]
	9838  Kabini [Radeon HD 8240 / R3 Series]
	9839  Kabini [Radeon HD 8180]
	983d  Temash [Radeon HD 8250/8280G]
	9840  Kabini HDMI/DP Audio
	9850  Mullins [Radeon R3 Graphics]
	9851  Mullins [Radeon R4/R5 Graphics]
	9852  Mullins [Radeon R2 Graphics]
	9853  Mullins [Radeon R2 Graphics]
	9854  Mullins [Radeon R3E Graphics]
	9855  Mullins [Radeon R6 Graphics]
	9856  Mullins [Radeon R1E/R2E Graphics]
	9857  Mullins [Radeon APU XX-2200M with R2 Graphics]
	9858  Mullins
	9859  Mullins
	985a  Mullins
	985b  Mullins
	985c  Mullins
	985d  Mullins
	985e  Mullins
	985f  Mullins
	9874  Carrizo
	9900  Trinity [Radeon HD 7660G]
	9901  Trinity [Radeon HD 7660D]
	9902  Trinity HDMI Audio Controller
	9903  Trinity [Radeon HD 7640G]
	9904  Trinity [Radeon HD 7560D]
	9905  Trinity [FirePro A300 Series Graphics]
	9906  Trinity [FirePro A300 Series Graphics]
	9907  Trinity [Radeon HD 7620G]
	9908  Trinity [Radeon HD 7600G]
	9909  Trinity [Radeon HD 7500G]
	990a  Trinity [Radeon HD 7500G]
	990b  Richland [Radeon HD 8650G]
	990c  Richland [Radeon HD 8670D]
	990d  Richland [Radeon HD 8550G]
	990e  Richland [Radeon HD 8570D]
	990f  Richland [Radeon HD 8610G]
	9910  Trinity [Radeon HD 7660G]
	9913  Trinity [Radeon HD 7640G]
	9917  Trinity [Radeon HD 7620G]
	9918  Trinity [Radeon HD 7600G]
	9919  Trinity [Radeon HD 7500G]
	9920  Liverpool [Playstation 4 APU]
	9921  Liverpool HDMI/DP Audio Controller
	9990  Trinity [Radeon HD 7520G]
	9991  Trinity [Radeon HD 7540D]
	9992  Trinity [Radeon HD 7420G]
	9993  Trinity [Radeon HD 7480D]
	9994  Trinity [Radeon HD 7400G]
	9995  Richland [Radeon HD 8450G]
	9996  Richland [Radeon HD 8470D]
	9997  Richland [Radeon HD 8350G]
	9998  Richland [Radeon HD 8370D]
	9999  Richland [Radeon HD 8510G]
	999a  Richland [Radeon HD 8410G]
	999b  Richland [Radeon HD 8310G]
	999c  Richland
	999d  Richland [Radeon HD 8550D]
	99a0  Trinity [Radeon HD 7520G]
	99a2  Trinity [Radeon HD 7420G]
	99a4  Trinity [Radeon HD 7400G]
	aa00  R600 HDMI Audio [Radeon HD 2900 GT/PRO/XT]
	aa01  RV635 HDMI Audio [Radeon HD 3650/3730/3750]
	aa08  RV630 HDMI Audio [Radeon HD 2600 PRO/XT / HD 3610]
	aa10  RV610 HDMI Audio [Radeon HD 2350 PRO / 2400 PRO/XT / HD 3410]
	aa18  RV670/680 HDMI Audio [Radeon HD 3690/3800 Series]
	aa20  RV635 HDMI Audio [Radeon HD 3650/3730/3750]
	aa28  RV620 HDMI Audio [Radeon HD 3450/3470/3550/3570]
	aa30  RV770 HDMI Audio [Radeon HD 4850/4870]
	aa38  RV710/730 HDMI Audio [Radeon HD 4000 series]
	aa50  Cypress HDMI Audio [Radeon HD 5830/5850/5870 / 6850/6870 Rebrand]
	aa58  Juniper HDMI Audio [Radeon HD 5700 Series]
	aa60  Redwood HDMI Audio [Radeon HD 5000 Series]
	aa68  Cedar HDMI Audio [Radeon HD 5400/6300/7300 Series]
	aa80  Cayman/Antilles HDMI Audio [Radeon HD 6930/6950/6970/6990]
	aa88  Barts HDMI Audio [Radeon HD 6790/6850/6870 / 7720 OEM]
	aa90  Turks HDMI Audio [Radeon HD 6500/6600 / 6700M Series]
	aa98  Caicos HDMI Audio [Radeon HD 6450 / 7450/8450/8490 OEM / R5 230/235/235X OEM]
	aaa0  Tahiti HDMI Audio [Radeon HD 7870 XT / 7950/7970]
	aab0  Cape Verde/Pitcairn HDMI Audio [Radeon HD 7700/7800 Series]
	aac0  Tobago HDMI Audio [Radeon R7 360 / R9 360 OEM]
	aac8  Hawaii HDMI Audio [Radeon R9 290/290X / 390/390X]
	aad8  Tonga HDMI Audio [Radeon R9 285/380]
	aae8  Fiji HDMI/DP Audio [Radeon R9 Nano / FURY/FURY X]
	aaf0  Ellesmere [Radeon RX 580]
	ac00  Theater 600 Pro
	ac02  TV Wonder HD 600 PCIe
	ac12  Theater HD T507 (DVB-T) TV tuner/capture device
	cab0  RS100 Host Bridge
	cab2  RS200 Host Bridge
	cab3  RS250 Host Bridge
	cbb2  RS200 Host Bridge
1003  ULSI Systems
	0201  US201
1004  VLSI Technology Inc
	0005  82C592-FC1
	0006  82C593-FC1
	0007  82C594-AFC2
	0008  82C596/7 [Wildcat]
	0009  82C597-AFC2
	000c  82C541 [Lynx]
	000d  82C543 [Lynx]
	0101  82C532
	0102  82C534 [Eagle]
	0103  82C538
	0104  82C535
	0105  82C147
	0200  82C975
	0280  82C925
	0304  QSound ThunderBird PCI Audio
	0305  QSound ThunderBird PCI Audio Gameport
	0306  QSound ThunderBird PCI Audio Support Registers
	0307  SAA7785 ThunderBird PCI Audio
	0308  SAA7785 ThunderBird PCI Audio Gameport
	0702  VAS96011 [Golden Gate II]
	0703  Tollgate
1005  Avance Logic Inc. [ALI]
	2064  ALG2032/2064
	2128  ALG2364A
	2301  ALG2301
	2302  ALG2302
	2364  ALG2364
	2464  ALG2364A
	2501  ALG2564A/25128A
1006  Reply Group
1007  NetFrame Systems Inc
1008  Epson
100a  Phoenix Technologies
100b  National Semiconductor Corporation
	0001  DP83810
	0002  87415/87560 IDE
	000e  87560 Legacy I/O
	000f  FireWire Controller
	0011  NS87560 National PCI System I/O
	0012  USB Controller
	0020  DP83815 (MacPhyter) Ethernet Controller
	0021  PC87200 PCI to ISA Bridge
	0022  DP83820 10/100/1000 Ethernet Controller
	0028  Geode GX2 Host Bridge
	002a  CS5535 South Bridge
	002b  CS5535 ISA bridge
	002d  CS5535 IDE
	002e  CS5535 Audio
	002f  CS5535 USB
	0030  Geode GX2 Graphics Processor
	0035  DP83065 [Saturn] 10/100/1000 Ethernet Controller
	0500  SCx200 Bridge
	0501  SCx200 SMI
	0502  SCx200, SC1100 IDE controller
	0503  SCx200, SC1100 Audio Controller
	0504  SCx200 Video
	0505  SCx200 XBus
	0510  SC1100 Bridge
	0511  SC1100 SMI & ACPI
	0515  SC1100 XBus
	d001  87410 IDE
100c  Tseng Labs Inc
	3202  ET4000/W32p rev A
	3205  ET4000/W32p rev B
	3206  ET4000/W32p rev C
	3207  ET4000/W32p rev D
	3208  ET6000
	4702  ET6300
100d  AST Research Inc
100e  Weitek
	9000  P9000 Viper
	9001  P9000 Viper
	9002  P9000 Viper
	9100  P9100 Viper Pro/SE
1010  Video Logic, Ltd.
1011  Digital Equipment Corporation
	0001  DECchip 21050
	0002  DECchip 21040 [Tulip]
	0004  DECchip 21030 [TGA]
	0007  NVRAM [Zephyr NVRAM]
	0008  KZPSA [KZPSA]
	0009  DECchip 21140 [FasterNet]
	000a  21230 Video Codec
	000d  PBXGB [TGA2]
	000f  DEFPA FDDI PCI-to-PDQ Interface Chip [PFI]
	0014  DECchip 21041 [Tulip Pass 3]
	0016  DGLPB [OPPO]
	0017  PV-PCI Graphics Controller (ZLXp-L)
	0018  Memory Channel interface
	0019  DECchip 21142/43
	001a  Farallon PN9000SX Gigabit Ethernet
	0021  DECchip 21052
	0022  DECchip 21150
	0023  DECchip 21150
	0024  DECchip 21152
	0025  DECchip 21153
	0026  DECchip 21154
	0034  56k Modem Cardbus
	0045  DECchip 21553
	0046  DECchip 21554
	1065  StrongARM DC21285
1012  Micronics Computers Inc
1013  Cirrus Logic
	0038  GD 7548
	0040  GD 7555 Flat Panel GUI Accelerator
	004c  GD 7556 Video/Graphics LCD/CRT Ctrlr
	00a0  GD 5430/40 [Alpine]
	00a2  GD 5432 [Alpine]
	00a4  GD 5434-4 [Alpine]
	00a8  GD 5434-8 [Alpine]
	00ac  GD 5436 [Alpine]
	00b0  GD 5440
	00b8  GD 5446
	00bc  GD 5480
	00d0  GD 5462
	00d2  GD 5462 [Laguna I]
	00d4  GD 5464 [Laguna]
	00d5  GD 5464 BD [Laguna]
	00d6  GD 5465 [Laguna]
	00e8  GD 5436U
	1100  CL 6729
	1110  PD 6832 PCMCIA/CardBus Ctrlr
	1112  PD 6834 PCMCIA/CardBus Ctrlr
	1113  PD 6833 PCMCIA/CardBus Ctrlr
	1200  GD 7542 [Nordic]
	1202  GD 7543 [Viking]
	1204  GD 7541 [Nordic Light]
	4000  MD 5620 [CLM Data Fax Voice]
	4400  CD 4400
	6001  CS 4610/11 [CrystalClear SoundFusion Audio Accelerator]
	6003  CS 4614/22/24/30 [CrystalClear SoundFusion Audio Accelerator]
	6004  CS 4614/22/24 [CrystalClear SoundFusion Audio Accelerator]
	6005  Crystal CS4281 PCI Audio
1014  IBM
	0002  PCI to MCA Bridge
	0005  Processor to I/O Controller [Alta Lite]
	0007  Processor to I/O Controller [Alta MP]
	000a  PCI to ISA Bridge (IBM27-82376) [Fire Coral]
	0017  CPU to PCI Bridge
	0018  TR Auto LANstreamer
	001b  GXT-150P
	001c  Carrera
	001d  SCSI-2 FAST PCI Adapter (82G2675)
	0020  GXT1000 Graphics Adapter
	0022  PCI to PCI Bridge (IBM27-82351)
	002d  Processor to I/O Controller [Python]
	002e  SCSI RAID Adapter [ServeRAID]
	0031  2 Port Serial Adapter
	0036  PCI to 32-bit LocalBus Bridge [Miami]
	0037  PowerPC to PCI Bridge (IBM27-82660)
	003a  CPU to PCI Bridge
	003c  GXT250P/GXT255P Graphics Adapter
	003e  16/4 Token ring UTP/STP controller
	0045  SSA Adapter
	0046  MPIC interrupt controller
	0047  PCI to PCI Bridge
	0048  PCI to PCI Bridge
	0049  Warhead SCSI Controller
	004e  ATM Controller (14104e00)
	004f  ATM Controller (14104f00)
	0050  ATM Controller (14105000)
	0053  25 MBit ATM Controller
	0054  GXT500P/GXT550P Graphics Adapter
	0057  MPEG PCI Bridge
	0058  SSA Adapter [Advanced SerialRAID/X]
	005e  GXT800P Graphics Adapter
	007c  ATM Controller (14107c00)
	007d  3780IDSP [MWave]
	008b  EADS PCI to PCI Bridge
	008e  GXT3000P Graphics Adapter
	0090  GXT 3000P
	0091  SSA Adapter
	0095  20H2999 PCI Docking Bridge
	0096  Chukar chipset SCSI controller
	009f  PCI 4758 Cryptographic Accelerator
	00a5  ATM Controller (1410a500)
	00a6  ATM 155MBPS MM Controller (1410a600)
	00b7  GXT2000P Graphics Adapter
	00b8  GXT2000P Graphics Adapter
	00be  ATM 622MBPS Controller (1410be00)
	00dc  Advanced Systems Management Adapter (ASMA)
	00fc  CPC710 Dual Bridge and Memory Controller (PCI-64)
	0105  CPC710 Dual Bridge and Memory Controller (PCI-32)
	010f  Remote Supervisor Adapter (RSA)
	0142  Yotta Video Compositor Input
	0144  Yotta Video Compositor Output
	0156  405GP PLB to PCI Bridge
	015e  622Mbps ATM PCI Adapter
	0160  64bit/66MHz PCI ATM 155 MMF
	016e  GXT4000P Graphics Adapter
	0170  GXT6000P Graphics Adapter
	017d  GXT300P Graphics Adapter
	0180  Snipe chipset SCSI controller
	0188  EADS-X PCI-X to PCI-X Bridge
	01a7  PCI-X to PCI-X Bridge
	01bd  ServeRAID Controller
	01c1  64bit/66MHz PCI ATM 155 UTP
	01e6  Cryptographic Accelerator
	01ef  PowerPC 440GP PCI Bridge
	01ff  10/100 Mbps Ethernet
	0219  Multiport Serial Adapter
	021b  GXT6500P Graphics Adapter
	021c  GXT4500P Graphics Adapter
	0233  GXT135P Graphics Adapter
	028c  Citrine chipset SCSI controller
	02a1  Calgary PCI-X Host Bridge
	02bd  Obsidian chipset SCSI controller
	0302  Winnipeg PCI-X Host Bridge
	0308  CalIOC2 PCI-E Root Port
	0311  FC 5740/1954 4-Port 10/100/1000 Base-TX PCI-X Adapter for POWER
	0314  ZISC 036 Neural accelerator card
	032d  Axon - Cell Companion Chip
	0339  Obsidian-E PCI-E SCSI controller
	033d  PCI-E IPR SAS Adapter (FPGA)
	034a  PCI-E IPR SAS Adapter (ASIC)
	03dc  POWER8 Host Bridge (PHB3)
	044b  GenWQE Accelerator Adapter
	04aa  Flash Adapter 90 (PCIe2 0.9TB)
	04da  PCI-E IPR SAS+ Adapter (ASIC)
	04ed  Internal Shared Memory (ISM) virtual PCI device
	3022  QLA3022 Network Adapter
	4022  QLA3022 Network Adapter
	ffff  MPIC-2 interrupt controller
1015  LSI Logic Corp of Canada
1016  ICL Personal Systems
1017  SPEA Software AG
	5343  SPEA 3D Accelerator
1018  Unisys Systems
1019  Elitegroup Computer Systems
101a  AT&T GIS (NCR)
	0005  100VG ethernet
	0007  BYNET BIC4G/2C/2G
	0009  PQS Memory Controller
	000a  BYNET BPCI Adapter
	000b  BYNET 4 Port BYA Switch (BYA4P)
	000c  BYNET 4 Port BYA Switch (BYA4G)
	0010  NCR AMC Memory Controller
	1dc1  BYNET BIC2M/BIC4M/BYA4M
	1fa8  BYNET Multi-port BIC Adapter (XBIC Based)
101b  Vitesse Semiconductor
	0452  VSC452 [SuperBMC]
101c  Western Digital
	0193  33C193A
	0196  33C196A
	0197  33C197A
	0296  33C296A
	3193  7193
	3197  7197
	3296  33C296A
	4296  34C296
	9710  Pipeline 9710
	9712  Pipeline 9712
	c24a  90C
101d  Maxim Integrated Products
101e  American Megatrends Inc.
	0009  MegaRAID 428 Ultra RAID Controller (rev 03)
	1960  MegaRAID
	9010  MegaRAID 428 Ultra RAID Controller
	9030  EIDE Controller
	9031  EIDE Controller
	9032  EIDE & SCSI Controller
	9033  SCSI Controller
	9040  Multimedia card
	9060  MegaRAID 434 Ultra GT RAID Controller
	9063  MegaRAC
101f  PictureTel
1020  Hitachi Computer Products
1021  OKI Electric Industry Co. Ltd.
1022  Advanced Micro Devices, Inc. [AMD]
	1100  K8 [Athlon64/Opteron] HyperTransport Technology Configuration
	1101  K8 [Athlon64/Opteron] Address Map
	1102  K8 [Athlon64/Opteron] DRAM Controller
	1103  K8 [Athlon64/Opteron] Miscellaneous Control
	1200  Family 10h Processor HyperTransport Configuration
	1201  Family 10h Processor Address Map
	1202  Family 10h Processor DRAM Controller
	1203  Family 10h Processor Miscellaneous Control
	1204  Family 10h Processor Link Control
	1300  Family 11h Processor HyperTransport Configuration
	1301  Family 11h Processor Address Map
	1302  Family 11h Processor DRAM Controller
	1303  Family 11h Processor Miscellaneous Control
	1304  Family 11h Processor Link Control
	1400  Family 15h (Models 10h-1fh) Processor Function 0
	1401  Family 15h (Models 10h-1fh) Processor Function 1
	1402  Family 15h (Models 10h-1fh) Processor Function 2
	1403  Family 15h (Models 10h-1fh) Processor Function 3
	1404  Family 15h (Models 10h-1fh) Processor Function 4
	1405  Family 15h (Models 10h-1fh) Processor Function 5
	1410  Family 15h (Models 10h-1fh) Processor Root Complex
	1412  Family 15h (Models 10h-1fh) Processor Root Port
	1413  Family 15h (Models 10h-1fh) Processor Root Port
	1414  Family 15h (Models 10h-1fh) Processor Root Port
	1415  Family 15h (Models 10h-1fh) Processor Root Port
	1416  Family 15h (Models 10h-1fh) Processor Root Port
	1417  Family 15h (Models 10h-1fh) Processor Root Port
	1418  Family 15h (Models 10h-1fh) Processor Root Port
	1419  Family 15h (Models 10h-1fh) I/O Memory Management Unit
	141a  Family 15h (Models 30h-3fh) Processor Function 0
	141b  Family 15h (Models 30h-3fh) Processor Function 1
	141c  Family 15h (Models 30h-3fh) Processor Function 2
	141d  Family 15h (Models 30h-3fh) Processor Function 3
	141e  Family 15h (Models 30h-3fh) Processor Function 4
	141f  Family 15h (Models 30h-3fh) Processor Function 5
	1422  Family 15h (Models 30h-3fh) Processor Root Complex
	1423  Family 15h (Models 30h-3fh) I/O Memory Management Unit
	1426  Family 15h (Models 30h-3fh) Processor Root Port
	1436  Liverpool Processor Root Complex
	1437  Liverpool I/O Memory Management Unit
	1438  Liverpool Processor Root Port
	1439  Family 16h Processor Functions 5:1
	1450  Family 17h (Models 00h-0fh) Root Complex
	1451  Family 17h (Models 00h-0fh) I/O Memory Management Unit
	1452  Family 17h (Models 00h-0fh) PCIe Dummy Host Bridge
	1454  Family 17h (Models 00h-0fh) Internal PCIe GPP Bridge 0 to Bus B
	1456  Family 17h (Models 00h-0fh) Platform Security Processor
	1457  Family 17h (Models 00h-0fh) HD Audio Controller
	145b  Zeppelin Non-Transparent Bridge
	145c  Family 17h (Models 00h-0fh) USB 3.0 Host Controller
	145f  USB 3.0 Host controller
	1460  Family 17h (Models 00h-0fh) Data Fabric: Device 18h; Function 0
	1461  Family 17h (Models 00h-0fh) Data Fabric: Device 18h; Function 1
	1462  Family 17h (Models 00h-0fh) Data Fabric: Device 18h; Function 2
	1463  Family 17h (Models 00h-0fh) Data Fabric: Device 18h; Function 3
	1464  Family 17h (Models 00h-0fh) Data Fabric: Device 18h; Function 4
	1465  Family 17h (Models 00h-0fh) Data Fabric: Device 18h; Function 5
	1466  Family 17h (Models 00h-0fh) Data Fabric Device 18h Function 6
	1467  Family 17h (Models 00h-0fh) Data Fabric: Device 18h; Function 7
	1510  Family 14h Processor Root Complex
	1512  Family 14h Processor Root Port
	1513  Family 14h Processor Root Port
	1514  Family 14h Processor Root Port
	1515  Family 14h Processor Root Port
	1516  Family 14h Processor Root Port
	1530  Family 16h Processor Function 0
	1531  Family 16h Processor Function 1
	1532  Family 16h Processor Function 2
	1533  Family 16h Processor Function 3
	1534  Family 16h Processor Function 4
	1535  Family 16h Processor Function 5
	1536  Family 16h Processor Root Complex
	1538  Family 16h Processor Function 0
	1600  Family 15h Processor Function 0
	1601  Family 15h Processor Function 1
	1602  Family 15h Processor Function 2
	1603  Family 15h Processor Function 3
	1604  Family 15h Processor Function 4
	1605  Family 15h Processor Function 5
	1700  Family 12h/14h Processor Function 0
	1701  Family 12h/14h Processor Function 1
	1702  Family 12h/14h Processor Function 2
	1703  Family 12h/14h Processor Function 3
	1704  Family 12h/14h Processor Function 4
	1705  Family 12h Processor Root Complex
	1707  Family 12h Processor Root Port
	1708  Family 12h Processor Root Port
	1709  Family 12h Processor Root Port
	170a  Family 12h Processor Root Port
	170b  Family 12h Processor Root Port
	170c  Family 12h Processor Root Port
	170d  Family 12h Processor Root Port
	1716  Family 12h/14h Processor Function 5
	1718  Family 12h/14h Processor Function 6
	1719  Family 12h/14h Processor Function 7
	2000  79c970 [PCnet32 LANCE]
	2001  79c978 [HomePNA]
	2003  Am 1771 MBW [Alchemy]
	2020  53c974 [PCscsi]
	2040  79c974
	2080  CS5536 [Geode companion] Host Bridge
	2081  Geode LX Video
	2082  Geode LX AES Security Block
	208f  CS5536 GeodeLink PCI South Bridge
	2090  CS5536 [Geode companion] ISA
	2091  CS5536 [Geode companion] FLASH
	2093  CS5536 [Geode companion] Audio
	2094  CS5536 [Geode companion] OHC
	2095  CS5536 [Geode companion] EHC
	2096  CS5536 [Geode companion] UDC
	2097  CS5536 [Geode companion] UOC
	209a  CS5536 [Geode companion] IDE
	3000  ELanSC520 Microcontroller
	43a0  Hudson PCI to PCI bridge (PCIE port 0)
	43a1  Hudson PCI to PCI bridge (PCIE port 1)
	43a2  Hudson PCI to PCI bridge (PCIE port 2)
	43a3  Hudson PCI to PCI bridge (PCIE port 3)
	43b4  300 Series Chipset PCIe Port
	43b7  300 Series Chipset SATA Controller
	43bb  300 Series Chipset USB 3.1 xHCI Controller
	7006  AMD-751 [Irongate] System Controller
	7007  AMD-751 [Irongate] AGP Bridge
	700a  AMD-IGR4 AGP Host to PCI Bridge
	700b  AMD-IGR4 PCI to PCI Bridge
	700c  AMD-760 MP [IGD4-2P] System Controller
	700d  AMD-760 MP [IGD4-2P] AGP Bridge
	700e  AMD-760 [IGD4-1P] System Controller
	700f  AMD-760 [IGD4-1P] AGP Bridge
	7400  AMD-755 [Cobra] ISA
	7401  AMD-755 [Cobra] IDE
	7403  AMD-755 [Cobra] ACPI
	7404  AMD-755 [Cobra] USB
	7408  AMD-756 [Viper] ISA
	7409  AMD-756 [Viper] IDE
	740b  AMD-756 [Viper] ACPI
	740c  AMD-756 [Viper] USB
	7410  AMD-766 [ViperPlus] ISA
	7411  AMD-766 [ViperPlus] IDE
	7413  AMD-766 [ViperPlus] ACPI
	7414  AMD-766 [ViperPlus] USB
	7440  AMD-768 [Opus] ISA
	7441  AMD-768 [Opus] IDE
	7443  AMD-768 [Opus] ACPI
	7445  AMD-768 [Opus] Audio
	7446  AMD-768 [Opus] MC97 Modem
	7448  AMD-768 [Opus] PCI
	7449  AMD-768 [Opus] USB
	7450  AMD-8131 PCI-X Bridge
	7451  AMD-8131 PCI-X IOAPIC
	7454  AMD-8151 System Controller
	7455  AMD-8151 AGP Bridge
	7458  AMD-8132 PCI-X Bridge
	7459  AMD-8132 PCI-X IOAPIC
	7460  AMD-8111 PCI
	7461  AMD-8111 USB
	7462  AMD-8111 Ethernet
	7463  AMD-8111 USB EHCI
	7464  AMD-8111 USB OHCI
	7468  AMD-8111 LPC
	7469  AMD-8111 IDE
	746a  AMD-8111 SMBus 2.0
	746b  AMD-8111 ACPI
	746d  AMD-8111 AC97 Audio
	746e  AMD-8111 MC97 Modem
	756b  AMD-8111 ACPI
	7800  FCH SATA Controller [IDE mode]
	7801  FCH SATA Controller [AHCI mode]
	7802  FCH SATA Controller [RAID mode]
	7803  FCH SATA Controller [RAID mode]
	7804  FCH SATA Controller [AHCI mode]
	7805  FCH SATA Controller [RAID mode]
	7806  FCH SD Flash Controller
	7807  FCH USB OHCI Controller
	7808  FCH USB EHCI Controller
	7809  FCH USB OHCI Controller
	780b  FCH SMBus Controller
	780c  FCH IDE Controller
	780d  FCH Azalia Controller
	780e  FCH LPC Bridge
	780f  FCH PCI Bridge
	7812  FCH USB XHCI Controller
	7813  FCH SD Flash Controller
	7814  FCH USB XHCI Controller
	7900  FCH SATA Controller [IDE mode]
	7901  FCH SATA Controller [AHCI mode]
	7902  FCH SATA Controller [RAID mode]
	7903  FCH SATA Controller [RAID mode]
	7904  FCH SATA Controller [AHCI mode]
	7906  FCH SD Flash Controller
	7908  FCH USB EHCI Controller
	790b  FCH SMBus Controller
	790e  FCH LPC Bridge
	790f  FCH PCI Bridge
	7914  FCH USB XHCI Controller
	9600  RS780 Host Bridge
	9601  RS880 Host Bridge
	9602  RS780/RS880 PCI to PCI bridge (int gfx)
	9603  RS780 PCI to PCI bridge (ext gfx port 0)
	9604  RS780/RS880 PCI to PCI bridge (PCIE port 0)
	9605  RS780/RS880 PCI to PCI bridge (PCIE port 1)
	9606  RS780 PCI to PCI bridge (PCIE port 2)
	9607  RS780/RS880 PCI to PCI bridge (PCIE port 3)
	9608  RS780/RS880 PCI to PCI bridge (PCIE port 4)
	9609  RS780/RS880 PCI to PCI bridge (PCIE port 5)
	960a  RS780 PCI to PCI bridge (NB-SB link)
	960b  RS780 PCI to PCI bridge (ext gfx port 1)
1023  Trident Microsystems
	0194  82C194
	2000  4DWave DX
	2001  4DWave NX
	2100  CyberBlade XP4m32
	2200  XGI Volari XP5
	8400  CyberBlade/i7
	8420  CyberBlade/i7d
	8500  CyberBlade/i1
	8520  CyberBlade i1
	8620  CyberBlade/i1
	8820  CyberBlade XPAi1
	9320  TGUI 9320
	9350  GUI Accelerator
	9360  Flat panel GUI Accelerator
	9382  Cyber 9382 [Reference design]
	9383  Cyber 9383 [Reference design]
	9385  Cyber 9385 [Reference design]
	9386  Cyber 9386
	9388  Cyber 9388
	9397  Cyber 9397
	939a  Cyber 9397DVD
	9420  TGUI 9420
	9430  TGUI 9430
	9440  TGUI 9440
	9460  TGUI 9460
	9470  TGUI 9470
	9520  Cyber 9520
	9525  Cyber 9525
	9540  Cyber 9540
	9660  TGUI 9660/938x/968x
	9680  TGUI 9680
	9682  TGUI 9682
	9683  TGUI 9683
	9685  ProVIDIA 9685
	9750  3DImage 9750
	9753  TGUI 9753
	9754  TGUI 9754
	9759  TGUI 975
	9783  TGUI 9783
	9785  TGUI 9785
	9850  3DImage 9850
	9880  Blade 3D PCI/AGP
	9910  CyberBlade/XP
	9930  CyberBlade/XPm
	9960  CyberBlade XP2
1024  Zenith Data Systems
1025  Acer Incorporated [ALI]
	1435  M1435
	1445  M1445
	1449  M1449
	1451  M1451
	1461  M1461
	1489  M1489
	1511  M1511
	1512  ALI M1512 Aladdin
	1513  M1513
	1521  ALI M1521 Aladdin III CPU Bridge
	1523  ALI M1523 ISA Bridge
	1531  M1531 Northbridge [Aladdin IV/IV+]
	1533  M1533 PCI-to-ISA Bridge
	1535  M1535 PCI Bridge + Super I/O + FIR
	1541  M1541 Northbridge [Aladdin V]
	1542  M1542 Northbridge [Aladdin V]
	1543  M1543 PCI-to-ISA Bridge + Super I/O + FIR
	1561  M1561 Northbridge [Aladdin 7]
	1621  M1621 Northbridge [Aladdin-Pro II]
	1631  M1631 Northbridge+3D Graphics [Aladdin TNT2]
	1641  M1641 Northbridge [Aladdin-Pro IV]
	1647  M1647 [MaGiK1] PCI North Bridge
	1671  M1671 Northbridge [ALADDiN-P4]
	1672  Northbridge [CyberALADDiN-P4]
	3141  M3141
	3143  M3143
	3145  M3145
	3147  M3147
	3149  M3149
	3151  M3151
	3307  M3307 MPEG-I Video Controller
	3309  M3309 MPEG-II Video w/ Software Audio Decoder
	3321  M3321 MPEG-II Audio/Video Decoder
	5212  M4803
	5215  ALI PCI EIDE Controller
	5217  M5217H
	5219  M5219
	5225  M5225
	5229  M5229
	5235  M5235
	5237  M5237 PCI USB Host Controller
	5240  EIDE Controller
	5241  PCMCIA Bridge
	5242  General Purpose Controller
	5243  PCI to PCI Bridge Controller
	5244  Floppy Disk Controller
	5247  M1541 PCI to PCI Bridge
	5251  M5251 P1394 Controller
	5427  PCI to AGP Bridge
	5451  M5451 PCI AC-Link Controller Audio Device
	5453  M5453 PCI AC-Link Controller Modem Device
	7101  M7101 PCI PMU Power Management Controller
	9602  AMD RS780/RS880 PCI to PCI bridge (int gfx)
1028  Dell
	0001  PowerEdge Expandable RAID Controller 2/Si
	0002  PowerEdge Expandable RAID Controller 3/Di
	0003  PowerEdge Expandable RAID Controller 3/Si
	0004  PowerEdge Expandable RAID Controller 3/Di [Iguana]
	0006  PowerEdge Expandable RAID Controller 3/Di
	0007  Remote Access Card III
	0008  Remote Access Card III
	0009  Remote Access Card III: BMC/SMIC device not present
	000a  PowerEdge Expandable RAID Controller 3/Di
	000c  Embedded Remote Access or ERA/O
	000d  Embedded Remote Access: BMC/SMIC device
	000e  PowerEdge Expandable RAID controller 4/Di
	000f  PowerEdge Expandable RAID controller 4/Di
	0010  Remote Access Card 4
	0011  Remote Access Card 4 Daughter Card
	0012  Remote Access Card 4 Daughter Card Virtual UART
	0013  PowerEdge Expandable RAID controller 4
	0014  Remote Access Card 4 Daughter Card SMIC interface
	0015  PowerEdge Expandable RAID controller 5
	0016  PowerEdge Expandable RAID controller S300
	0073  NV-RAM Adapter
1029  Siemens Nixdorf IS
102a  LSI Logic
	0000  HYDRA
	0010  ASPEN
	001f  AHA-2940U2/U2W /7890/7891 SCSI Controllers
	00c5  AIC-7899 U160/m SCSI Controller
	00cf  AIC-7899P U160/m
102b  Matrox Electronics Systems Ltd.
	0010  MGA-I [Impression?]
	0100  MGA 1064SG [Mystique]
	0518  MGA-II [Athena]
	0519  MGA 2064W [Millennium]
	051a  MGA 1064SG [Mystique]
	051b  MGA 2164W [Millennium II]
	051e  MGA 1064SG [Mystique] AGP
	051f  MGA 2164W [Millennium II] AGP
	0520  MGA G200
	0521  MGA G200 AGP
	0522  MGA G200e [Pilot] ServerEngines (SEP1)
	0525  MGA G400/G450
	0527  Parhelia
	0528  Parhelia
	0530  MGA G200EV
	0532  MGA G200eW WPCM450
	0533  MGA G200EH
	0534  G200eR2
	0536  Integrated Matrox G200eW3 Graphics Controller
	0538  G200eH
	0540  M91XX
	0550  SV2
	0d10  MGA Ultima/Impression
	1000  MGA G100 [Productiva]
	1001  MGA G100 [Productiva] AGP
	2007  MGA Mistral
	2527  Millennium G550
	2537  Millennium P650/P750
	2538  Millennium P650 PCIe
	2539  Millennium P690
	4164  Morphis QxT frame grabber
	43b4  Morphis Qxt encoding engine
	4510  Morphis COM port
	4536  VIA Framegrabber
	4686  Concord GX (customized Intel 82541)
	475b  Solios eCL/XCL-B frame grabber
	475d  Vio frame grabber family
	475f  Solios (single-Full) CL frame grabber
	47a1  Solios eA/XA frame grabber
	47a2  Solios COM port
	47c1  Solios (dual-Base/single-Medium) CL frame grabber
	47c2  Solios COM port
	4949  Radient frame grabber family
	4cdc  Morphis JPEG2000 accelerator
	4f54  Morphis (e)Quad frame grabber
	4fc5  Morphis (e)Dual frame grabber
	5e10  Morphis aux I/O
	6573  Shark 10/100 Multiport SwitchNIC
102c  Chips and Technologies
	00b8  F64310
	00c0  F69000 HiQVideo
	00d0  F65545
	00d8  F65545
	00dc  F65548
	00e0  F65550
	00e4  F65554
	00e5  F65555 HiQVPro
	00f0  F68554
	00f4  F68554 HiQVision
	00f5  F68555
	0c30  F69030
102d  Wyse Technology Inc.
	50dc  3328 Audio
102e  Olivetti Advanced Technology
102f  Toshiba America
	0009  r4x00
	000a  TX3927 MIPS RISC PCI Controller
	0020  ATM Meteor 155
	0030  TC35815CF PCI 10/100 Mbit Ethernet Controller
	0031  TC35815CF PCI 10/100 Mbit Ethernet Controller with WOL
	0032  TC35815CF PCI 10/100 Mbit Ethernet Controller on TX4939
	0105  TC86C001 [goku-s] IDE
	0106  TC86C001 [goku-s] USB 1.1 Host
	0107  TC86C001 [goku-s] USB Device Controller
	0108  TC86C001 [goku-s] I2C/SIO/GPIO Controller
	0180  TX4927/38 MIPS RISC PCI Controller
	0181  TX4925 MIPS RISC PCI Controller
	0182  TX4937 MIPS RISC PCI Controller
	01b4  Celleb platform IDE interface
	01b5  SCC USB 2.0 EHCI controller
	01b6  SCC USB 1.1 OHCI controller
1030  TMC Research
1031  Miro Computer Products AG
	5601  DC20 ASIC
	5607  Video I/O & motion JPEG compressor
	5631  Media 3D
	6057  MiroVideo DC10/DC30+
1032  Compaq
1033  NEC Corporation
	0000  Vr4181A USB Host or Function Control Unit
	0001  PCI to 486-like bus Bridge
	0002  PCI to VL98 Bridge
	0003  ATM Controller
	0004  R4000 PCI Bridge
	0005  PCI to 486-like bus Bridge
	0006  PC-9800 Graphic Accelerator
	0007  PCI to UX-Bus Bridge
	0008  PC-9800 Graphic Accelerator
	0009  PCI to PC9800 Core-Graph Bridge
	0016  PCI to VL Bridge
	001a  [Nile II]
	0021  Vrc4373 [Nile I]
	0029  PowerVR PCX1
	002a  PowerVR 3D
	002c  Star Alpha 2
	002d  PCI to C-bus Bridge
	0035  OHCI USB Controller
	003b  PCI to C-bus Bridge
	003e  NAPCCARD Cardbus Controller
	0046  PowerVR PCX2 [midas]
	005a  Vrc5074 [Nile 4]
	0063  uPD72862 [Firewarden] IEEE1394 OHCI 1.0 Link Controller
	0067  PowerVR Neon 250 Chipset
	0072  uPD72874 IEEE1394 OHCI 1.1 3-port PHY-Link Ctrlr
	0074  56k Voice Modem
	009b  Vrc5476
	00a5  VRC4173
	00a6  VRC5477 AC97
	00cd  uPD72870 [Firewarden] IEEE1394a OHCI 1.0 Link/3-port PHY Controller
	00ce  uPD72871 [Firewarden] IEEE1394a OHCI 1.0 Link/1-port PHY Controller
	00df  Vr4131
	00e0  uPD72010x USB 2.0 Controller
	00e7  uPD72873 [Firewarden] IEEE1394a OHCI 1.1 Link/2-port PHY Controller
	00f2  uPD72874 [Firewarden] IEEE1394a OHCI 1.1 Link/3-port PHY Controller
	00f3  uPD6113x Multimedia Decoder/Processor [EMMA2]
	010c  VR7701
	0125  uPD720400 PCI Express - PCI/PCI-X Bridge
	013a  Dual Tuner/MPEG Encoder
	0194  uPD720200 USB 3.0 Host Controller
	01e7  uPD72873 [Firewarden] IEEE1394a OHCI 1.1 Link/2-port PHY Controller
	01f2  uPD72874 [Firewarden] IEEE1394a OHCI 1.1 Link/3-port PHY Controller
1034  Framatome Connectors USA Inc.
1035  Comp. & Comm. Research Lab
1036  Future Domain Corp.
	0000  TMC-18C30 [36C70]
1037  Hitachi Micro Systems
1038  AMP, Inc
1039  Silicon Integrated Systems [SiS]
	0001  AGP Port (virtual PCI-to-PCI bridge)
	0002  AGP Port (virtual PCI-to-PCI bridge)
	0003  AGP Port (virtual PCI-to-PCI bridge)
	0004  PCI-to-PCI bridge
	0006  85C501/2/3
	0008  SiS85C503/5513 (LPC Bridge)
	0009  5595 Power Management Controller
	000a  PCI-to-PCI bridge
	0016  SiS961/2/3 SMBus controller
	0018  SiS85C503/5513 (LPC Bridge)
	0163  163 802.11b/g Wireless LAN Adapter
	0180  RAID bus controller 180 SATA/PATA  [SiS]
	0181  SATA
	0182  182 SATA/RAID Controller
	0186  AHCI Controller (0106)
	0190  190 Ethernet Adapter
	0191  191 Gigabit Ethernet Adapter
	0200  5597/5598/6326 VGA
	0204  82C204
	0205  SG86C205
	0300  300/305 PCI/AGP VGA Display Adapter
	0310  315H PCI/AGP VGA Display Adapter
	0315  315 PCI/AGP VGA Display Adapter
	0325  315PRO PCI/AGP VGA Display Adapter
	0330  330 [Xabre] PCI/AGP VGA Display Adapter
	0406  85C501/2
	0496  85C496
	0530  530 Host
	0540  540 Host
	0550  550 Host
	0597  5513C
	0601  85C601
	0620  620 Host
	0630  630 Host
	0633  633 Host
	0635  635 Host
	0645  SiS645 Host & Memory & AGP Controller
	0646  SiS645DX Host & Memory & AGP Controller
	0648  645xx
	0649  SiS649 Host
	0650  650/M650 Host
	0651  651 Host
	0655  655 Host
	0660  660 Host
	0661  661FX/M661FX/M661MX Host
	0662  662 Host
	0671  671MX
	0730  730 Host
	0733  733 Host
	0735  735 Host
	0740  740 Host
	0741  741/741GX/M741 Host
	0745  745 Host
	0746  746 Host
	0755  755 Host
	0760  760/M760 Host
	0761  761/M761 Host
	0900  SiS900 PCI Fast Ethernet
	0961  SiS961 [MuTIOL Media IO]
	0962  SiS962 [MuTIOL Media IO] LPC Controller
	0963  SiS963 [MuTIOL Media IO] LPC Controller
	0964  SiS964 [MuTIOL Media IO] LPC Controller
	0965  SiS965 [MuTIOL Media IO]
	0966  SiS966 [MuTIOL Media IO]
	0968  SiS968 [MuTIOL Media IO]
	1180  SATA Controller / IDE mode
	1182  SATA Controller / RAID mode
	1183  SATA Controller / IDE mode
	1184  AHCI Controller / RAID mode
	1185  AHCI IDE Controller (0106)
	3602  83C602
	5107  5107
	5300  SiS540 PCI Display Adapter
	5315  550 PCI/AGP VGA Display Adapter
	5401  486 PCI Chipset
	5511  5511/5512
	5513  5513 IDE Controller
	5517  5517
	5571  5571
	5581  5581 Pentium Chipset
	5582  5582
	5591  5591/5592 Host
	5596  5596 Pentium Chipset
	5597  5597 [SiS5582]
	5600  5600 Host
	6204  Video decoder & MPEG interface
	6205  VGA Controller
	6236  6236 3D-AGP
	6300  630/730 PCI/AGP VGA Display Adapter
	6306  530/620 PCI/AGP VGA Display Adapter
	6325  65x/M650/740 PCI/AGP VGA Display Adapter
	6326  86C326 5598/6326
	6330  661/741/760 PCI/AGP or 662/761Gx PCIE VGA Display Adapter
	6350  770/670 PCIE VGA Display Adapter
	6351  771/671 PCIE VGA Display Adapter
	7001  USB 1.1 Controller
	7002  USB 2.0 Controller
	7007  FireWire Controller
	7012  SiS7012 AC'97 Sound Controller
	7013  AC'97 Modem Controller
	7016  SiS7016 PCI Fast Ethernet Adapter
	7018  SiS PCI Audio Accelerator
	7019  SiS7019 Audio Accelerator
	7502  Azalia Audio Controller
103a  Seiko Epson Corporation
103b  Tatung Corp. Of America
103c  Hewlett-Packard Company
	1005  A4977A Visualize EG
	1008  Visualize FX
	1028  Tach TL Fibre Channel Host Adapter
	1029  Tach XL2 Fibre Channel Host Adapter
	102a  Tach TS Fibre Channel Host Adapter
	1030  J2585A DeskDirect 10/100VG NIC
	1031  J2585B HP 10/100VG PCI LAN Adapter
	1040  J2973A DeskDirect 10BaseT NIC
	1041  J2585B DeskDirect 10/100 NIC
	1042  J2970A DeskDirect 10BaseT/2 NIC
	1048  Diva Serial [GSP] Multiport UART
	1054  PCI Local Bus Adapter
	1064  79C970 PCnet Ethernet Controller
	108b  Visualize FXe
	10c1  NetServer Smart IRQ Router
	10ed  TopTools Remote Control
	10f0  rio System Bus Adapter
	10f1  rio I/O Controller
	1219  NetServer PCI Hot-Plug Controller
	121a  NetServer SMIC Controller
	121b  NetServer Legacy COM Port Decoder
	121c  NetServer PCI COM Port Decoder
	1229  zx1 System Bus Adapter
	122a  zx1 I/O Controller
	122e  PCI-X Local Bus Adapter
	127b  sx1000 System Bus Adapter
	127c  sx1000 I/O Controller
	1290  Auxiliary Diva Serial Port
	1291  Auxiliary Diva Serial Port
	12b4  zx1 QuickSilver AGP8x Local Bus Adapter
	12eb  sx2000 System Bus Adapter
	12ec  sx2000 I/O Controller
	12ee  PCI-X 2.0 Local Bus Adapter
	1302  RMP-3 Shared Memory Driver
	1303  RMP-3 (Remote Management Processor)
	22f6  iLO5 Virtual USB Controller
	2910  E2910A PCIBus Exerciser
	2925  E2925A 32 Bit, 33 MHzPCI Exerciser & Analyzer
	3206  Adaptec Embedded Serial ATA HostRAID
	3220  Smart Array P600
	3230  Smart Array Controller
	3238  Smart Array E200i (SAS Controller)
	3239  Smart Array Gen9 Controllers
	323a  Smart Array G6 controllers
	323b  Smart Array Gen8 Controllers
	323c  Smart Array Gen8+ Controllers
	3300  Integrated Lights-Out Standard Virtual USB Controller
	3301  Integrated Lights-Out Standard Serial Port
	3302  Integrated Lights-Out Standard KCS Interface
	3305  Integrated Lights-Out (iLO2) Controller
	3306  Integrated Lights-Out Standard Slave Instrumentation & System Support
	3307  Integrated Lights-Out Standard Management Processor Support and Messaging
	3308  Integrated Lights-Out Standard MS Watchdog Timer
	4030  zx2 System Bus Adapter
	4031  zx2 I/O Controller
	4037  PCIe Local Bus Adapter
103e  Solliday Engineering
103f  Synopsys/Logic Modeling Group
1040  Accelgraphics Inc.
1041  Computrend
1042  Micron
	1000  PC Tech RZ1000
	1001  PC Tech RZ1001
	3000  Samurai_0
	3010  Samurai_1
	3020  Samurai_IDE
1043  ASUSTeK Computer Inc.
	0464  Radeon R9 270x GPU
	0675  ISDNLink P-IN100-ST-D
	9602  AMD RS780/RS880 PCI to PCI bridge (int gfx)
1044  Adaptec (formerly DPT)
	1012  Domino RAID Engine
	a400  SmartCache/Raid I-IV Controller
	a500  PCI Bridge
	a501  SmartRAID V Controller
	a511  SmartRAID V Controller
	c066  3010S Ultra3 Dual Channel
1045  OPTi Inc.
	a0f8  82C750 [Vendetta] USB Controller
	c101  92C264
	c178  92C178
	c556  82X556 [Viper]
	c557  82C557 [Viper-M]
	c558  82C558 [Viper-M ISA+IDE]
	c567  82C750 [Vendetta], device 0
	c568  82C750 [Vendetta], device 1
	c569  82C579 [Viper XPress+ Chipset]
	c621  82C621 [Viper-M/N+]
	c700  82C700 [FireStar]
	c701  82C701 [FireStar Plus]
	c814  82C814 [Firebridge 1]
	c822  82C822
	c824  82C824
	c825  82C825 [Firebridge 2]
	c832  82C832
	c861  82C861
	c881  82C881 [FireLink] 1394 OHCI Link Controller
	c895  82C895
	c935  EV1935 ECTIVA MachOne PCIAudio
	d568  82C825 [Firebridge 2]
	d721  IDE [FireStar]
1046  IPC Corporation, Ltd.
1047  Genoa Systems Corp
1048  Elsa AG
	0c60  Gladiac MX
	0d22  Quadro4 900XGL [ELSA GLoria4 900XGL]
	1000  QuickStep 1000
	3000  QuickStep 3000
	8901  Gloria XL
1049  Fountain Technologies, Inc.
104a  STMicroelectronics
	0000  STLS2F Host Bridge
	0008  STG 2000X
	0009  STG 1764X
	0010  STG4000 [3D Prophet Kyro Series]
	0201  STPC Vega Northbridge
	0209  STPC Consumer/Industrial North- and Southbridge
	020a  STPC Atlas/ConsumerS/Consumer IIA Northbridge
	020b  STPC Consumer II ISA Bridge
	0210  STPC Atlas ISA Bridge
	021a  STPC Consumer S Southbridge
	021b  STPC Consumer IIA Southbridge
	0220  STPC Industrial PCI to PCCard bridge
	0228  STPC Atlas IDE
	0229  STPC Vega IDE
	0230  STPC Atlas/Vega OHCI USB Controller
	0238  STPC Vega LAN
	0500  ST70137 [Unicorn] ADSL DMT Transceiver
	0564  STPC Client Northbridge
	0981  21x4x DEC-Tulip compatible 10/100 Ethernet
	1746  STG 1764X
	2774  21x4x DEC-Tulip compatible 10/100 Ethernet
	3520  MPEG-II decoder card
	55cc  STPC Client Southbridge
104b  BusLogic
	0140  BT-946C (old) [multimaster  01]
	1040  BT-946C (BA80C30) [MultiMaster 10]
	8130  Flashpoint LT
104c  Texas Instruments
	0500  100 MBit LAN Controller
	0508  TMS380C2X Compressor Interface
	1000  Eagle i/f AS
	104c  PCI1510 PC card Cardbus Controller
	3d04  TVP4010 [Permedia]
	3d07  TVP4020 [Permedia 2]
	8000  PCILynx/PCILynx2 IEEE 1394 Link Layer Controller
	8009  TSB12LV22 IEEE-1394 Controller
	8017  PCI4410 FireWire Controller
	8019  TSB12LV23 IEEE-1394 Controller
	8020  TSB12LV26 IEEE-1394 Controller (Link)
	8021  TSB43AA22 IEEE-1394 Controller (PHY/Link Integrated)
	8022  TSB43AB22 IEEE-1394a-2000 Controller (PHY/Link) [iOHCI-Lynx]
	8023  TSB43AB22A IEEE-1394a-2000 Controller (PHY/Link) [iOHCI-Lynx]
	8024  TSB43AB23 IEEE-1394a-2000 Controller (PHY/Link)
	8025  TSB82AA2 IEEE-1394b Link Layer Controller
	8026  TSB43AB21 IEEE-1394a-2000 Controller (PHY/Link)
	8027  PCI4451 IEEE-1394 Controller
	8029  PCI4510 IEEE-1394 Controller
	802b  PCI7410,7510,7610 OHCI-Lynx Controller
	802e  PCI7x20 1394a-2000 OHCI Two-Port PHY/Link-Layer Controller
	8031  PCIxx21/x515 Cardbus Controller
	8032  OHCI Compliant IEEE 1394 Host Controller
	8033  PCIxx21 Integrated FlashMedia Controller
	8034  PCI6411/6421/6611/6621/7411/7421/7611/7621 Secure Digital Controller
	8035  PCI6411/6421/6611/6621/7411/7421/7611/7621 Smart Card Controller
	8036  PCI6515 Cardbus Controller
	8038  PCI6515 SmartCard Controller
	8039  PCIxx12 Cardbus Controller
	803a  PCIxx12 OHCI Compliant IEEE 1394 Host Controller
	803b  5-in-1 Multimedia Card Reader (SD/MMC/MS/MS PRO/xD)
	803c  PCIxx12 SDA Standard Compliant SD Host Controller
	803d  PCIxx12 GemCore based SmartCard controller
	8101  TSB43DB42 IEEE-1394a-2000 Controller (PHY/Link)
	8201  PCI1620 Firmware Loading Function
	8204  PCI7410/7510/7610 PCI Firmware Loading Function
	8231  XIO2000(A)/XIO2200A PCI Express-to-PCI Bridge
	8232  XIO3130 PCI Express Switch (Upstream)
	8233  XIO3130 PCI Express Switch (Downstream)
	8235  XIO2200A IEEE-1394a-2000 Controller (PHY/Link)
	823e  XIO2213A/B/XIO2221 PCI Express to PCI Bridge [Cheetah Express]
	823f  XIO2213A/B/XIO2221 IEEE-1394b OHCI Controller [Cheetah Express]
	8240  XIO2001 PCI Express-to-PCI Bridge
	8241  TUSB73x0 SuperSpeed USB 3.0 xHCI Host Controller
	8400  ACX 100 22Mbps Wireless Interface
	8401  ACX 100 22Mbps Wireless Interface
	8888  Multicore DSP+ARM KeyStone II SOC
	9000  Wireless Interface (of unknown type)
	9065  TMS320DM642
	9066  ACX 111 54Mbps Wireless Interface
	a001  TDC1570
	a100  TDC1561
	a102  TNETA1575 HyperSAR Plus w/PCI Host i/f & UTOPIA i/f
	a106  TMS320C6414 TMS320C6415 TMS320C6416
	ac10  PCI1050
	ac11  PCI1053
	ac12  PCI1130
	ac13  PCI1031
	ac15  PCI1131
	ac16  PCI1250
	ac17  PCI1220
	ac18  PCI1260
	ac19  PCI1221
	ac1a  PCI1210
	ac1b  PCI1450
	ac1c  PCI1225
	ac1d  PCI1251A
	ac1e  PCI1211
	ac1f  PCI1251B
	ac20  TI 2030
	ac21  PCI2031
	ac22  PCI2032 PCI Docking Bridge
	ac23  PCI2250 PCI-to-PCI Bridge
	ac28  PCI2050 PCI-to-PCI Bridge
	ac2c  PCI2060 PCI-to-PCI Bridge
	ac30  PCI1260 PC card Cardbus Controller
	ac40  PCI4450 PC card Cardbus Controller
	ac41  PCI4410 PC card Cardbus Controller
	ac42  PCI4451 PC card Cardbus Controller
	ac44  PCI4510 PC card Cardbus Controller
	ac46  PCI4520 PC card Cardbus Controller
	ac47  PCI7510 PC card Cardbus Controller
	ac48  PCI7610 PC Card Cardbus Controller
	ac49  PCI7410 PC Card Cardbus Controller
	ac4a  PCI7510/7610 CardBus Bridge
	ac4b  PCI7610 SD/MMC controller
	ac4c  PCI7610 Memory Stick controller
	ac50  PCI1410 PC card Cardbus Controller
	ac51  PCI1420 PC card Cardbus Controller
	ac52  PCI1451 PC card Cardbus Controller
	ac53  PCI1421 PC card Cardbus Controller
	ac54  PCI1620 PC Card Controller
	ac55  PCI1520 PC card Cardbus Controller
	ac56  PCI1510 PC card Cardbus Controller
	ac60  PCI2040 PCI to DSP Bridge Controller
	ac8d  PCI 7620
	ac8e  PCI7420 CardBus Controller
	ac8f  PCI7420/7620 SD/MS-Pro Controller
	b001  TMS320C6424
	fe00  FireWire Host Controller
	fe03  12C01A FireWire Host Controller
104d  Sony Corporation
	8004  DTL-H2500 [Playstation development board]
	8009  CXD1947Q i.LINK Controller
	8039  CXD3222 i.LINK Controller
	8056  Rockwell HCF 56K modem
	808a  Memory Stick Controller
	81ce  SxS Pro memory card
	905c  SxS Pro memory card
	907f  SxS Pro+ memory card
	908f  Aeolia ACPI
	909e  Aeolia Ethernet Controller (Marvell Yukon 2 Family)
	909f  Aeolia SATA AHCI Controller
	90a0  Aeolia SD/MMC Host Controller
	90a1  Aeolia PCI Express Glue and Miscellaneous Devices
	90a2  Aeolia DMA Controller
	90a3  Aeolia Memory (DDR3/SPM)
	90a4  Aeolia USB 3.0 xHCI Host Controller
	90bc  SxS Pro+ memory card
104e  Oak Technology, Inc
	0017  OTI-64017
	0107  OTI-107 [Spitfire]
	0109  Video Adapter
	0111  OTI-64111 [Spitfire]
	0217  OTI-64217
	0317  OTI-64317
104f  Co-time Computer Ltd
1050  Winbond Electronics Corp
	0000  NE2000
	0001  W83769F
	0033  W89C33D 802.11 a/b/g BB/MAC
	0105  W82C105
	0840  W89C840
	0940  W89C940
	5a5a  W89C940F
	6692  W6692
	9921  W99200F MPEG-1 Video Encoder
	9922  W99200F/W9922PF MPEG-1/2 Video Encoder
	9970  W9970CF
1051  Anigma, Inc.
1052  ?Young Micro Systems
1053  Young Micro Systems
1054  Hitachi, Ltd
	3009  2Gbps Fibre Channel to PCI HBA 3009
	300a  4Gbps Fibre Channel to PCI-X HBA 300a
	300b  4Gbps Fibre Channel to PCI-X HBA 300b
	300f  ColdFusion 3 Chipset Processor to I/O Controller
	3010  ColdFusion 3 Chipset Memory Controller Hub
	3011  ColdFusion 3e Chipset Processor to I/O Controller
	3012  ColdFusion 3e Chipset Memory Controller Hub
	3017  Unassigned Hitachi Shared FC Device 3017
	301b  Virtual VGA Device
	301d  PCIe-to-PCIe Bridge with Virtualization IO Assist Feature
	3020  FIVE-EX based Fibre Channel to PCIe HBA
	302c  M001 PCI Express Switch Upstream Port
	302d  M001 PCI Express Switch Downstream Port
	3070  Hitachi FIVE-FX Fibre Channel to PCIe HBA
	3505  SH7751 PCI Controller (PCIC)
	350e  SH7751R PCI Controller (PCIC)
1055  Efar Microsystems
	9130  SLC90E66 [Victory66] IDE
	9460  SLC90E66 [Victory66] ISA
	9462  SLC90E66 [Victory66] USB
	9463  SLC90E66 [Victory66] ACPI
	e420  LAN9420/LAN9420i
1056  ICL
1057  Motorola
	0001  MPC105 [Eagle]
	0002  MPC106 [Grackle]
	0003  MPC8240 [Kahlua]
	0004  MPC107
	0006  MPC8245 [Unity]
	0008  MPC8540
	0009  MPC8560
	0012  MPC8548 [PowerQUICC III]
	0100  MC145575 [HFC-PCI]
	0431  KTI829c 100VG
	1073  Nokia N770
	1219  Nokia N800
	1801  DSP56301 Digital Signal Processor
	18c0  MPC8265A/8266/8272
	18c1  MPC8271/MPC8272
	3052  SM56 Data Fax Modem
	3410  DSP56361 Digital Signal Processor
	4801  Raven
	4802  Falcon
	4803  Hawk
	4806  CPX8216
	4d68  20268
	5600  SM56 PCI Modem
	5608  Wildcard X100P
	5803  MPC5200
	5806  MCF54 Coldfire
	5808  MPC8220
	5809  MPC5200B
	6400  MPC190 Security Processor (S1 family, encryption)
	6405  MPC184 Security Processor (S1 family)
1058  Electronics & Telecommunications RSH
1059  Kontron
105a  Promise Technology, Inc.
	0d30  PDC20265 (FastTrak100 Lite/Ultra100)
	0d38  20263
	1275  20275
	3318  PDC20318 (SATA150 TX4)
	3319  PDC20319 (FastTrak S150 TX4)
	3371  PDC20371 (FastTrak S150 TX2plus)
	3373  PDC20378 (FastTrak 378/SATA 378)
	3375  PDC20375 (SATA150 TX2plus)
	3376  PDC20376 (FastTrak 376)
	3515  PDC40719 [FastTrak TX4300/TX4310]
	3519  PDC40519 (FastTrak TX4200)
	3570  PDC20771 [FastTrak TX2300]
	3571  PDC20571 (FastTrak TX2200)
	3574  PDC20579 SATAII 150 IDE Controller
	3577  PDC40779 (SATA 300 779)
	3d17  PDC40718 (SATA 300 TX4)
	3d18  PDC20518/PDC40518 (SATAII 150 TX4)
	3d73  PDC40775 (SATA 300 TX2plus)
	3d75  PDC20575 (SATAII150 TX2plus)
	3f20  PDC42819 [FastTrak TX2650/TX4650]
	4302  80333 [SuperTrak EX4350]
	4d30  PDC20267 (FastTrak100/Ultra100)
	4d33  20246
	4d38  PDC20262 (FastTrak66/Ultra66)
	4d68  PDC20268 [Ultra100 TX2]
	4d69  20269
	5275  PDC20276 (MBFastTrak133 Lite)
	5300  DC5300
	6268  PDC20270 (FastTrak100 LP/TX2/TX4)
	6269  PDC20271 (FastTrak TX2000)
	6300  PDC81731 [FastTrak SX8300]
	6621  PDC20621 (FastTrak S150 SX4/FastTrak SX4000 lite)
	6622  PDC20621 [SATA150 SX4] 4 Channel IDE RAID Controller
	6624  PDC20621 [FastTrak SX4100]
	6626  PDC20618 (Ultra 618)
	6629  PDC20619 (FastTrak TX4000)
	7275  PDC20277 (SBFastTrak133 Lite)
	8002  SATAII150 SX8
	8350  80333 [SuperTrak EX8350/EX16350], 80331 [SuperTrak EX8300/EX16300]
	8650  81384 [SuperTrak EX SAS and SATA RAID Controller]
	8760  PM8010 [SuperTrak EX SAS and SATA 6G RAID Controller]
	c350  80333 [SuperTrak EX12350]
	e350  80333 [SuperTrak EX24350]
105b  Foxconn International, Inc.
105c  Wipro Infotech Limited
105d  Number 9 Computer Company
	2309  Imagine 128
	2339  Imagine 128-II
	493d  Imagine 128 T2R [Ticket to Ride]
	5348  Revolution 4
105e  Vtech Computers Ltd
105f  Infotronic America Inc
1060  United Microelectronics [UMC]
	0001  UM82C881
	0002  UM82C886
	0101  UM8673F
	0881  UM8881
	0886  UM8886F
	0891  UM8891A
	1001  UM886A
	673a  UM8886BF
	673b  EIDE Master/DMA
	8710  UM8710
	886a  UM8886A
	8881  UM8881F
	8886  UM8886F
	888a  UM8886A
	8891  UM8891A
	9017  UM9017F
	9018  UM9018
	9026  UM9026
	e881  UM8881N
	e886  UM8886N
	e88a  UM8886N
	e891  UM8891N
1061  I.I.T.
	0001  AGX016
	0002  IIT3204/3501
1062  Maspar Computer Corp
1063  Ocean Office Automation
1064  Alcatel
	1102  Dynamite 2840 (ADSL PCI modem)
1065  Texas Microsystems
1066  PicoPower Technology
	0000  PT80C826
	0001  PT86C521 [Vesuvius v1] Host Bridge
	0002  PT86C523 [Vesuvius v3] PCI-ISA Bridge Master
	0003  PT86C524 [Nile] PCI-to-PCI Bridge
	0004  PT86C525 [Nile-II] PCI-to-PCI Bridge
	0005  National PC87550 System Controller
	8002  PT86C523 [Vesuvius v3] PCI-ISA Bridge Slave
1067  Mitsubishi Electric
	0301  AccelGraphics AccelECLIPSE
	0304  AccelGALAXY A2100 [OEM Evans & Sutherland]
	0308  Tornado 3000 [OEM Evans & Sutherland]
	1002  VG500 [VolumePro Volume Rendering Accelerator]
1068  Diversified Technology
1069  Mylex Corporation
	0001  DAC960P
	0002  DAC960PD
	0010  DAC960PG
	0020  DAC960LA
	0050  AcceleRAID 352/170/160 support Device
	b166  AcceleRAID 600/500/400/Sapphire support Device
	ba55  eXtremeRAID 1100 support Device
	ba56  eXtremeRAID 2000/3000 support Device
	ba57  eXtremeRAID 4000/5000 support Device
106a  Aten Research Inc
106b  Apple Inc.
	0001  Bandit PowerPC host bridge
	0002  Grand Central I/O
	0003  Control Video
	0004  PlanB Video-In
	0007  O'Hare I/O
	000c  DOS on Mac
	000e  Hydra Mac I/O
	0010  Heathrow Mac I/O
	0017  Paddington Mac I/O
	0018  UniNorth FireWire
	0019  KeyLargo USB
	001e  UniNorth Internal PCI
	001f  UniNorth PCI
	0020  UniNorth AGP
	0021  UniNorth GMAC (Sun GEM)
	0022  KeyLargo Mac I/O
	0024  UniNorth/Pangea GMAC (Sun GEM)
	0025  KeyLargo/Pangea Mac I/O
	0026  KeyLargo/Pangea USB
	0027  UniNorth/Pangea AGP
	0028  UniNorth/Pangea PCI
	0029  UniNorth/Pangea Internal PCI
	002d  UniNorth 1.5 AGP
	002e  UniNorth 1.5 PCI
	002f  UniNorth 1.5 Internal PCI
	0030  UniNorth/Pangea FireWire
	0031  UniNorth 2 FireWire
	0032  UniNorth 2 GMAC (Sun GEM)
	0033  UniNorth 2 ATA/100
	0034  UniNorth 2 AGP
	0035  UniNorth 2 PCI
	0036  UniNorth 2 Internal PCI
	003b  UniNorth/Intrepid ATA/100
	003e  KeyLargo/Intrepid Mac I/O
	003f  KeyLargo/Intrepid USB
	0040  K2 KeyLargo USB
	0041  K2 KeyLargo Mac/IO
	0042  K2 FireWire
	0043  K2 ATA/100
	0045  K2 HT-PCI Bridge
	0046  K2 HT-PCI Bridge
	0047  K2 HT-PCI Bridge
	0048  K2 HT-PCI Bridge
	0049  K2 HT-PCI Bridge
	004a  CPC945 HT Bridge
	004b  U3 AGP
	004c  K2 GMAC (Sun GEM)
	004f  Shasta Mac I/O
	0050  Shasta IDE
	0051  Shasta (Sun GEM)
	0052  Shasta Firewire
	0053  Shasta PCI Bridge
	0054  Shasta PCI Bridge
	0055  Shasta PCI Bridge
	0056  U4 PCIe
	0057  U3 HT Bridge
	0058  U3L AGP Bridge
	0059  U3H AGP Bridge
	005b  CPC945 PCIe Bridge
	0066  Intrepid2 AGP Bridge
	0067  Intrepid2 PCI Bridge
	0068  Intrepid2 PCI Bridge
	0069  Intrepid2 ATA/100
	006a  Intrepid2 Firewire
	006b  Intrepid2 GMAC (Sun GEM)
	0074  U4 HT Bridge
	1645  Broadcom NetXtreme BCM5701 Gigabit Ethernet
	2001  S1X NVMe Controller
	2002  S3ELab NVMe Controller
	2003  S3X NVMe Controller
	2005  ANS2 NVMe Controller
106c  Hynix Semiconductor
	8139  8139c 100BaseTX Ethernet Controller
	8801  Dual Pentium ISA/PCI Motherboard
	8802  PowerPC ISA/PCI Motherboard
	8803  Dual Window Graphics Accelerator
	8804  LAN Controller
	8805  100-BaseT LAN
106d  Sequent Computer Systems
106e  DFI, Inc
106f  City Gate Development Ltd
1070  Daewoo Telecom Ltd
1071  Mitac
	8160  Mitac 8060B Mobile Platform
1072  GIT Co Ltd
1073  Yamaha Corporation
	0001  3D GUI Accelerator
	0002  YGV615 [RPA3 3D-Graphics Controller]
	0003  YMF-740
	0004  YMF-724
	0005  DS1 Audio
	0006  DS1 Audio
	0008  DS1 Audio
	000a  DS1L Audio
	000c  YMF-740C [DS-1L Audio Controller]
	000d  YMF-724F [DS-1 Audio Controller]
	0010  YMF-744B [DS-1S Audio Controller]
	0012  YMF-754 [DS-1E Audio Controller]
	0020  DS-1 Audio
	1000  SW1000XG [XG Factory]
	2000  DS2416 Digital Mixing Card
1074  NexGen Microsystems
	4e78  82c500/1
1075  Advanced Integrations Research
1076  Chaintech Computer Co. Ltd
1077  QLogic Corp.
	1016  ISP10160 Single Channel Ultra3 SCSI Processor
	1020  ISP1020 Fast-wide SCSI
	1022  ISP1022 Fast-wide SCSI
	1080  ISP1080 SCSI Host Adapter
	1216  ISP12160 Dual Channel Ultra3 SCSI Processor
	1240  ISP1240 SCSI Host Adapter
	1280  ISP1280 SCSI Host Adapter
	1634  FastLinQ QL45000 Series 40GbE Controller
	1644  FastLinQ QL45000 Series 100GbE Controller
	1654  FastLinQ QL45000 Series 50GbE Controller
	1656  FastLinQ QL45000 Series 25GbE Controller
	165c  FastLinQ QL45000 Series 10/25/40/50GbE Controller (FCoE)
	165e  FastLinQ QL45000 Series 10/25/40/50GbE Controller (iSCSI)
	1664  FastLinQ QL45000 Series Gigabit Ethernet Controller (SR-IOV VF)
	2020  ISP2020A Fast!SCSI Basic Adapter
	2031  ISP8324-based 16Gb Fibre Channel to PCI Express Adapter
	2071  ISP2714-based 16/32Gb Fibre Channel to PCIe Adapter
	2100  QLA2100 64-bit Fibre Channel Adapter
	2200  QLA2200 64-bit Fibre Channel Adapter
	2261  ISP2722-based 16/32Gb Fibre Channel to PCIe Adapter
	2300  QLA2300 64-bit Fibre Channel Adapter
	2312  ISP2312-based 2Gb Fibre Channel to PCI-X HBA
	2322  ISP2322-based 2Gb Fibre Channel to PCI-X HBA
	2422  ISP2422-based 4Gb Fibre Channel to PCI-X HBA
	2432  ISP2432-based 4Gb Fibre Channel to PCI Express HBA
	2532  ISP2532-based 8Gb Fibre Channel to PCI Express HBA
	2971  ISP2684
	3022  ISP4022-based Ethernet NIC
	3032  ISP4032-based Ethernet IPv6 NIC
	4010  ISP4010-based iSCSI TOE HBA
	4022  ISP4022-based iSCSI TOE HBA
	4032  ISP4032-based iSCSI TOE IPv6 HBA
	5432  SP232-based 4Gb Fibre Channel to PCI Express HBA
	6312  SP202-based 2Gb Fibre Channel to PCI-X HBA
	6322  SP212-based 2Gb Fibre Channel to PCI-X HBA
	7220  IBA7220 InfiniBand HCA
	7322  IBA7322 QDR InfiniBand HCA
	8000  10GbE Converged Network Adapter (TCP/IP Networking)
	8001  10GbE Converged Network Adapter (FCoE)
	8020  cLOM8214 1/10GbE Controller
	8021  8200 Series 10GbE Converged Network Adapter (FCoE)
	8022  8200 Series 10GbE Converged Network Adapter (iSCSI)
	8030  ISP8324 1/10GbE Converged Network Controller
	8031  8300 Series 10GbE Converged Network Adapter (FCoE)
	8032  8300 Series 10GbE Converged Network Adapter (iSCSI)
	8070  FastLinQ QL41000 Series 10/25/40/50GbE Controller
	8080  FastLinQ QL41000 Series 10/25/40/50GbE Controller (FCoE)
	8084  FastLinQ QL41000 Series 10/25/40/50GbE Controller (iSCSI)
	8090  FastLinQ QL41000 Series Gigabit Ethernet Controller (SR-IOV VF)
	8430  ISP8324 1/10GbE Converged Network Controller (NIC VF)
	8431  8300 Series 10GbE Converged Network Adapter (FCoE VF)
	8432  ISP2432M-based 10GbE Converged Network Adapter (CNA)
1078  Cyrix Corporation
	0000  5510 [Grappa]
	0001  PCI Master
	0002  5520 [Cognac]
	0100  5530 Legacy [Kahlua]
	0101  5530 SMI [Kahlua]
	0102  5530 IDE [Kahlua]
	0103  5530 Audio [Kahlua]
	0104  5530 Video [Kahlua]
	0400  ZFMicro PCI Bridge
	0401  ZFMicro Chipset SMI
	0402  ZFMicro Chipset IDE
	0403  ZFMicro Expansion Bus
1079  I-Bus
107a  NetWorth
107b  Gateway, Inc.
107c  LG Electronics [Lucky Goldstar Co. Ltd]
107d  LeadTek Research Inc.
	0000  P86C850
107e  Interphase Corporation
	0001  5515 ATM Adapter [Flipper]
	0002  100 VG AnyLan Controller
	0004  5526 Fibre Channel Host Adapter
	0005  x526 Fibre Channel Host Adapter
	0008  5525/5575 ATM Adapter (155 Mbit) [Atlantic]
	9003  5535-4P-BRI-ST
	9007  5535-4P-BRI-U
	9008  5535-1P-SR
	900c  5535-1P-SR-ST
	900e  5535-1P-SR-U
	9011  5535-1P-PRI
	9013  5535-2P-PRI
	9023  5536-4P-BRI-ST
	9027  5536-4P-BRI-U
	9031  5536-1P-PRI
	9033  5536-2P-PRI
107f  Data Technology Corporation
	0802  SL82C105
1080  Contaq Microsystems
	0600  82C599
	c691  Cypress CY82C691
	c693  82c693
1081  Supermac Technology
	0d47  Radius PCI to NuBUS Bridge
1082  EFA Corporation of America
1083  Forex Computer Corporation
	0001  FR710
1084  Parador
1086  J. Bond Computer Systems
1087  Cache Computer
1088  Microcomputer Systems (M) Son
1089  Data General Corporation
108a  SBS Technologies
	0001  VME Bridge Model 617
	0010  VME Bridge Model 618
	0040  dataBLIZZARD
	3000  VME Bridge Model 2706
108c  Oakleigh Systems Inc.
108d  Olicom
	0001  Token-Ring 16/4 PCI Adapter (3136/3137)
	0002  16/4 Token Ring
	0004  RapidFire OC-3139/3140 Token-Ring 16/4 PCI Adapter
	0005  GoCard 3250 Token-Ring 16/4 CardBus PC Card
	0006  OC-3530 RapidFire Token-Ring 100
	0007  RapidFire 3141 Token-Ring 16/4 PCI Fiber Adapter
	0008  RapidFire 3540 HSTR 100/16/4 PCI Adapter
	0011  OC-2315
	0012  OC-2325
	0013  OC-2183/2185
	0014  OC-2326
	0019  OC-2327/2250 10/100 Ethernet Adapter
	0021  OC-6151/6152 [RapidFire ATM 155]
	0022  ATM Adapter
108e  Oracle/SUN
	0001  EBUS
	1000  EBUS
	1001  Happy Meal 10/100 Ethernet [hme]
	1100  RIO EBUS
	1101  RIO 10/100 Ethernet [eri]
	1102  RIO 1394
	1103  RIO USB
	1647  Broadcom 570x 10/100/1000 Ethernet [bge]
	1648  Broadcom 570x 10/100/1000 Ethernet [bge]
	16a7  Broadcom 570x 10/100/1000 Ethernet [bge]
	16a8  Broadcom 570x 10/100/1000 Ethernet [bge]
	2bad  GEM 10/100/1000 Ethernet [ge]
	5000  Simba Advanced PCI Bridge
	5043  SunPCI Co-processor
	5ca0  Crypto Accelerator 6000 [mca]
	6300  Intel 21554 PCI-PCI bus bridge [db21554]
	6301  Intel 21554 PCI-PCI bus bridge [db21554]
	6302  Intel 21554 PCI-PCI bus bridge [db21554]
	6303  Intel 21554 PCI-PCI bus bridge [db21554]
	6310  Intel 21554 PCI-PCI bus bridge [db21554]
	6311  Intel 21554 PCI-PCI bus bridge [db21554]
	6312  Intel 21554 PCI-PCI bus bridge [db21554]
	6313  Intel 21554 PCI-PCI bus bridge [db21554]
	6320  Intel 21554 PCI-PCI bus bridge [db21554]
	6323  Intel 21554 PCI-PCI bus bridge [db21554]
	6330  Intel 21554 PCI-PCI bus bridge [db21554]
	6331  Intel 21554 PCI-PCI bus bridge [db21554]
	6332  Intel 21554 PCI-PCI bus bridge [db21554]
	6333  Intel 21554 PCI-PCI bus bridge [db21554]
	6340  Intel 21554 PCI-PCI bus bridge [db21554]
	6343  Intel 21554 PCI-PCI bus bridge [db21554]
	6350  Intel 21554 PCI-PCI bus bridge [db21554]
	6353  Intel 21554 PCI-PCI bus bridge [db21554]
	6722  Intel 21554 PCI-PCI bus bridge [db21554]
	676e  SunPCiIII
	7063  SunPCiII / SunPCiIIpro
	8000  Psycho PCI Bus Module
	8001  Schizo PCI Bus Module
	8002  Schizo+ PCI Bus Module
	80f0  PCIe switch [px]
	80f8  PCIe switch [px]
	9010  PCIe/PCI bridge switch [pxb_plx]
	9020  PCIe/PCI bridge switch [pxb_plx]
	9102  Davicom Fast Ethernet driver for Davicom DM9102A [dmfe]
	a000  Psycho UPA-PCI Bus Module [pcipsy]
	a001  Psycho UPA-PCI Bus Module [pcipsy]
	a801  Schizo Fireplane-PCI bus bridge module [pcisch]
	aaaa  Multithreaded Shared 10GbE Ethernet Network Controller
	abba  Cassini 10/100/1000
	abcd  Multithreaded 10-Gigabit Ethernet Network Controller
	c416  Sun Fire System/System Controller Interface chip [sbbc]
108f  Systemsoft
1090  Compro Computer Services, Inc.
	4610  PCI RTOM
	4620  GPIO HSD
1091  Intergraph Corporation
	0020  3D graphics processor
	0021  3D graphics processor w/Texturing
	0040  3D graphics frame buffer
	0041  3D graphics frame buffer
	0060  Proprietary bus bridge
	00e4  Powerstorm 4D50T
	0720  Motion JPEG codec
	0780  Intense3D Wildcat 3410 (MSMT496)
	07a0  Sun Expert3D-Lite Graphics Accelerator
	1091  Sun Expert3D Graphics Accelerator
1092  Diamond Multimedia Systems
	0028  Viper V770
	00a0  Speedstar Pro SE
	00a8  Speedstar 64
	0550  Viper V550
	08d4  Supra 2260 Modem
	094c  SupraExpress 56i Pro
	1001  Video Crunch It 1001 capture card
	1092  Viper V330
	6120  Maximum DVD
	8810  Stealth SE
	8811  Stealth 64/SE
	8880  Stealth
	8881  Stealth
	88b0  Stealth 64
	88b1  Stealth 64
	88c0  Stealth 64
	88c1  Stealth 64
	88d0  Stealth 64
	88d1  Stealth 64
	88f0  Stealth 64
	88f1  Stealth 64
	9999  DMD-I0928-1 "Monster sound" sound chip
1093  National Instruments
	0160  PCI-DIO-96
	0162  PCI-MIO-16XE-50
	0fe1  PXI-8320
	1150  PCI-6533 (PCI-DIO-32HS)
	1170  PCI-MIO-16XE-10
	1180  PCI-MIO-16E-1
	1190  PCI-MIO-16E-4
	11b0  PXI-6070E
	11c0  PXI-6040E
	11d0  PXI-6030E
	1270  PCI-6032E
	1290  PCI-6704
	12b0  PCI-6534
	1310  PCI-6602
	1320  PXI-6533
	1330  PCI-6031E
	1340  PCI-6033E
	1350  PCI-6071E
	1360  PXI-6602
	13c0  PXI-6508
	1490  PXI-6534
	14e0  PCI-6110
	14f0  PCI-6111
	1580  PXI-6031E
	15b0  PXI-6071E
	1710  PXI-6509
	17c0  PXI-5690
	17d0  PCI-6503
	1870  PCI-6713
	1880  PCI-6711
	18b0  PCI-6052E
	18c0  PXI-6052E
	1920  PXI-6704
	1930  PCI-6040E
	19c0  PCI-4472
	1aa0  PXI-4110
	1ad0  PCI-6133
	1ae0  PXI-6133
	1e30  PCI-6624
	1e40  PXI-6624
	1e50  PXI-5404
	2410  PCI-6733
	2420  PXI-6733
	2430  PCI-6731
	2470  PCI-4474
	24a0  PCI-4065
	24b0  PXI-4200
	24f0  PXI-4472
	2510  PCI-4472
	2520  PCI-4474
	27a0  PCI-6123
	27b0  PXI-6123
	2880  DAQCard-6601
	2890  PCI-6036E
	28a0  PXI-4461
	28b0  PCI-6013
	28c0  PCI-6014
	28d0  PCI-5122
	28e0  PXI-5122
	29f0  PXI-7334
	2a00  PXI-7344
	2a60  PCI-6023E
	2a70  PCI-6024E
	2a80  PCI-6025E
	2ab0  PXI-6025E
	2b10  PXI-6527
	2b20  PCI-6527
	2b80  PXI-6713
	2b90  PXI-6711
	2c60  PCI-6601
	2c70  PXI-6601
	2c80  PCI-6035E
	2c90  PCI-6703
	2ca0  PCI-6034E
	2cb0  PCI-7344
	2cc0  PXI-6608
	2d20  PXI-5600
	2db0  PCI-6608
	2dc0  PCI-4070
	2dd0  PXI-4070
	2eb0  PXI-4472
	2ec0  PXI-6115
	2ed0  PCI-6115
	2ee0  PXI-6120
	2ef0  PCI-6120
	2fd1  PCI-7334
	2fd2  PCI-7350
	2fd3  PCI-7342
	2fd5  PXI-7350
	2fd6  PXI-7342
	7003  PCI-6551
	7004  PXI-6551
	700b  PXI-5421
	700c  PCI-5421
	701a  VXIpc-87xB
	701b  VXIpc-770
	7023  PXI-2593
	7027  PCI-MXI-2 Universal
	702c  PXI-7831R
	702d  PCI-7831R
	702e  PXI-7811R
	702f  PCI-7811R
	7030  PCI-CAN (Series 2)
	7031  PCI-CAN/2 (Series 2)
	7032  PCI-CAN/LS (Series 2)
	7033  PCI-CAN/LS2 (Series 2)
	7034  PCI-CAN/DS (Series 2)
	7035  PXI-8460 (Series 2, 1 port)
	7036  PXI-8460 (Series 2, 2 ports)
	7037  PXI-8461 (Series 2, 1 port)
	7038  PXI-8461 (Series 2, 2 ports)
	7039  PXI-8462 (Series 2)
	703f  PXI-2566
	7040  PXI-2567
	7044  MXI-4 Connection Monitor
	7047  PXI-6653
	704c  PXI-2530
	704f  PXI-4220
	7050  PXI-4204
	7055  PXI-7830R
	7056  PCI-7830R
	705a  PCI-CAN/XS (Series 2)
	705b  PCI-CAN/XS2 (Series 2)
	705c  PXI-8464 (Series 2, 1 port)
	705d  PXI-8464 (Series 2, 2 ports)
	705e  cRIO-9102
	7060  PXI-5610
	7064  PXI-1045 Trigger Routing Module
	7065  PXI-6652
	7066  PXI-6651
	7067  PXI-2529
	7068  PCI-CAN/SW (Series 2)
	7069  PCI-CAN/SW2 (Series 2)
	706a  PXI-8463 (Series 2, 1 port)
	706b  PXI-8463 (Series 2, 2 ports)
	7073  PCI-6723
	7074  PXI-7833R
	7075  PXI-6552
	7076  PCI-6552
	707c  PXI-1428
	707e  PXI-4462
	7080  PXI-8430/2 (RS-232) Interface
	7081  PXI-8431/2 (RS-485) Interface
	7083  PCI-7833R
	7085  PCI-6509
	7086  PXI-6528
	7087  PCI-6515
	7088  PCI-6514
	708c  PXI-2568
	708d  PXI-2569
	70a9  PCI-6528
	70aa  PCI-6229
	70ab  PCI-6259
	70ac  PCI-6289
	70ad  PXI-6251
	70ae  PXI-6220
	70af  PCI-6221
	70b0  PCI-6220
	70b1  PXI-6229
	70b2  PXI-6259
	70b3  PXI-6289
	70b4  PCI-6250
	70b5  PXI-6221
	70b6  PCI-6280
	70b7  PCI-6254
	70b8  PCI-6251
	70b9  PXI-6250
	70ba  PXI-6254
	70bb  PXI-6280
	70bc  PCI-6284
	70bd  PCI-6281
	70be  PXI-6284
	70bf  PXI-6281
	70c0  PCI-6143
	70c3  PCI-6511
	70c4  PXI-7330
	70c5  PXI-7340
	70c6  PCI-7330
	70c7  PCI-7340
	70c8  PCI-6513
	70c9  PXI-6515
	70ca  PCI-1405
	70cc  PCI-6512
	70cd  PXI-6514
	70ce  PXI-1405
	70cf  PCIe-GPIB
	70d0  PXI-2570
	70d1  PXI-6513
	70d2  PXI-6512
	70d3  PXI-6511
	70d4  PCI-6722
	70d6  PXI-4072
	70d7  PXI-6541
	70d8  PXI-6542
	70d9  PCI-6541
	70da  PCI-6542
	70db  PCI-8430/2 (RS-232) Interface
	70dc  PCI-8431/2 (RS-485) Interface
	70dd  PXI-8430/4 (RS-232) Interface
	70de  PXI-8431/4 (RS-485) Interface
	70df  PCI-8430/4 (RS-232) Interface
	70e0  PCI-8431/4 (RS-485) Interface
	70e1  PXI-2532
	70e2  PXI-8430/8 (RS-232) Interface
	70e3  PXI-8431/8 (RS-485) Interface
	70e4  PCI-8430/8 (RS-232) Interface
	70e5  PCI-8431/8 (RS-485) Interface
	70e6  PXI-8430/16 (RS-232) Interface
	70e7  PCI-8430/16 (RS-232) Interface
	70e8  PXI-8432/2 (Isolated RS-232) Interface
	70e9  PXI-8433/2 (Isolated RS-485) Interface
	70ea  PCI-8432/2 (Isolated RS-232) Interface
	70eb  PCI-8433/2 (Isolated RS-485) Interface
	70ec  PXI-8432/4 (Isolated RS-232) Interface
	70ed  PXI-8433/4 (Isolated RS-485) Interface
	70ee  PCI-8432/4 (Isolated RS-232) Interface
	70ef  PCI-8433/4 (Isolated RS-485) Interface
	70f0  PXI-5922
	70f1  PCI-5922
	70f2  PCI-6224
	70f3  PXI-6224
	70f6  cRIO-9101
	70f7  cRIO-9103
	70f8  cRIO-9104
	70ff  PXI-6723
	7100  PXI-6722
	7104  PCIx-1429
	7105  PCIe-1429
	710a  PXI-4071
	710d  PXI-6143
	710e  PCIe-GPIB
	710f  PXI-5422
	7110  PCI-5422
	7111  PXI-5441
	7119  PXI-6561
	711a  PXI-6562
	711b  PCI-6561
	711c  PCI-6562
	7120  PCI-7390
	7121  PXI-5122EX
	7122  PCI-5122EX
	7123  PXIe-5653
	7124  PCI-6510
	7125  PCI-6516
	7126  PCI-6517
	7127  PCI-6518
	7128  PCI-6519
	7137  PXI-2575
	713c  PXI-2585
	713d  PXI-2586
	7142  PXI-4224
	7144  PXI-5124
	7145  PCI-5124
	7146  PCI-6132
	7147  PXI-6132
	7148  PCI-6122
	7149  PXI-6122
	714c  PXI-5114
	714d  PCI-5114
	7150  PXI-2564
	7152  PCI-5640R
	7156  PXI-1044 Trigger Routing Module
	715d  PCI-1426
	7167  PXI-5412
	7168  PCI-5412
	716b  PCI-6230
	716c  PCI-6225
	716d  PXI-6225
	716f  PCI-4461
	7170  PCI-4462
	7171  PCI-6010
	7174  PXI-8360
	7177  PXI-6230
	717d  PCIe-6251
	717f  PCIe-6259
	7187  PCI-1410
	718b  PCI-6521
	718c  PXI-6521
	7191  PCI-6154
	7193  PXI-7813R
	7194  PCI-7813R
	7195  PCI-8254R
	7197  PXI-5402
	7198  PCI-5402
	719f  PCIe-6535
	71a0  PCIe-6536
	71a3  PXI-5650
	71a4  PXI-5652
	71a5  PXI-2594
	71a7  PXI-2595
	71a9  PXI-2596
	71aa  PXI-2597
	71ab  PXI-2598
	71ac  PXI-2599
	71ad  PCI-GPIB+
	71ae  PCIe-1430
	71b7  PXI-1056 Trigger Routing Module
	71b8  PXI-1045 Trigger Routing Module
	71b9  PXI-1044 Trigger Routing Module
	71bb  PXI-2584
	71bc  PCI-6221 (37-pin)
	71bf  PCIe-1427
	71c5  PCI-6520
	71c6  PXI-2576
	71c7  cRIO-9072
	71dc  PCI-1588
	71e0  PCI-6255
	71e1  PXI-6255
	71e2  PXI-5406
	71e3  PCI-5406
	71fc  PXI-4022
	7209  PCI-6233
	720a  PXI-6233
	720b  PCI-6238
	720c  PXI-6238
	7260  PXI-5142
	7261  PCI-5142
	726d  PXI-5651
	7273  PXI-4461
	7274  PXI-4462
	7279  PCI-6232
	727a  PXI-6232
	727b  PCI-6239
	727c  PXI-6239
	727e  SMBus Controller
	7281  PCI-6236
	7282  PXI-6236
	7283  PXI-2554
	7288  PXIe-5611
	7293  PCIe-8255R
	729d  cRIO-9074
	72a4  PCIe-4065
	72a7  PCIe-6537
	72a8  PXI-5152
	72a9  PCI-5152
	72aa  PXI-5105
	72ab  PCI-5105
	72b8  PXI-6682
	72d0  PXI-2545
	72d1  PXI-2546
	72d2  PXI-2547
	72d3  PXI-2548
	72d4  PXI-2549
	72d5  PXI-2555
	72d6  PXI-2556
	72d7  PXI-2557
	72d8  PXI-2558
	72d9  PXI-2559
	72e8  PXIe-6251
	72e9  PXIe-6259
	72ef  PXI-4498
	72f0  PXI-4496
	72fb  PXIe-6672
	730e  PXI-4130
	730f  PXI-5922EX
	7310  PCI-5922EX
	731c  PXI-2535
	731d  PXI-2536
	7322  PXIe-6124
	7327  PXI-6529
	732c  VXI-8360T
	7331  PXIe-5602
	7332  PXIe-5601
	7333  PXI-5900
	7335  PXI-2533
	7336  PXI-2534
	7342  PXI-4461
	7349  PXI-5154
	734a  PCI-5154
	7357  PXI-4065
	7359  PXI-4495
	7370  PXI-4461
	7373  sbRIO-9601
	7374  IOtech-9601
	7375  sbRIO-9602
	7378  sbRIO-9641
	737d  PXI-5124EX
	7384  PXI-7851R
	7385  PXI-7852R
	7386  PCIe-7851R
	7387  PCIe-7852R
	7390  PXI-7841R
	7391  PXI-7842R
	7392  PXI-7853R
	7393  PCIe-7841R
	7394  PCIe-7842R
	7397  sbRIO-9611
	7398  sbRIO-9612
	7399  sbRIO-9631
	739a  sbRIO-9632
	739b  sbRIO-9642
	73a1  PXIe-4498
	73a2  PXIe-4496
	73a5  PXIe-5641R
	73a7  PXI-8250 Chassis Monitor Module
	73a8  PXI-8511 CAN/LS
	73a9  PXI-8511 CAN/LS
	73aa  PXI-8512 CAN/HS
	73ab  PXI-8512 CAN/HS
	73ac  PXI-8513 CAN/XS
	73ad  PXI-8513 CAN/XS
	73af  PXI-8516 LIN
	73b1  PXI-8517 FlexRay
	73b2  PXI-8531 CANopen
	73b3  PXI-8531 CANopen
	73b4  PXI-8532 DeviceNet
	73b5  PXI-8532 DeviceNet
	73b6  PCI-8511 CAN/LS
	73b7  PCI-8511 CAN/LS
	73b8  PCI-8512 CAN/HS
	73b9  PCI-8512 CAN/HS
	73ba  PCI-8513 CAN/XS
	73bb  PCI-8513 CAN/XS
	73bd  PCI-8516 LIN
	73bf  PCI-8517 FlexRay
	73c0  PCI-8531 CANopen
	73c1  PCI-8531 CANopen
	73c2  PCI-8532 DeviceNet
	73c3  PCI-8532 DeviceNet
	73c5  PXIe-2527
	73c6  PXIe-2529
	73c8  PXIe-2530
	73c9  PXIe-2532
	73ca  PXIe-2569
	73cb  PXIe-2575
	73cc  PXIe-2593
	73d5  PXI-7951R
	73d6  PXI-7952R
	73d7  PXI-7953R
	73e1  PXI-7854R
	73ec  PXI-7954R
	73ed  cRIO-9073
	73f0  PXI-5153
	73f1  PCI-5153
	73f4  PXI-2515
	73f6  cRIO-9111
	73f7  cRIO-9112
	73f8  cRIO-9113
	73f9  cRIO-9114
	73fa  cRIO-9116
	73fb  cRIO-9118
	7404  PXI-4132
	7405  PXIe-6674T
	7406  PXIe-6674
	740e  PCIe-8430/16 (RS-232) Interface
	740f  PCIe-8430/8 (RS-232) Interface
	7410  PCIe-8431/16 (RS-485) Interface
	7411  PCIe-8431/8 (RS-485) Interface
	7414  PCIe-GPIB+
	741c  PXI-5691
	741d  PXI-5695
	743c  CSC-3059
	7448  PXI-2510
	7454  PXI-2512
	7455  PXI-2514
	7456  PXIe-2512
	7457  PXIe-2514
	745a  PXI-6682H
	745e  PXI-5153EX
	745f  PCI-5153EX
	7460  PXI-5154EX
	7461  PCI-5154EX
	746d  PXIe-5650
	746e  PXIe-5651
	746f  PXIe-5652
	7472  PXI-2800
	7495  PXIe-5603
	7497  PXIe-5605
	74ae  PXIe-2515
	74b4  PXI-2531
	74b5  PXIe-2531
	74c1  PXIe-8430/16 (RS-232) Interface
	74c2  PXIe-8430/8 (RS-232) Interface
	74c3  PXIe-8431/16 (RS-485) Interface
	74c4  PXIe-8431/8 (RS-485) Interface
	74d5  PXIe-5630
	74d9  PCIe-8432/2 (Isolated RS-232) Interface
	74da  PCIe-8433/2 (Isolated RS-485) Interface
	74db  PCIe-8432/4 (Isolated RS-232) Interface
	74dc  PCIe-8433/4 (Isolated RS-485) Interface
	74e8  NI 9148
	7515  PCIe-8430/2 (RS-232) Interface
	7516  PCIe-8430/4 (RS-232) Interface
	7517  PCIe-8431/2 (RS-485) Interface
	7518  PCIe-8431/4 (RS-485) Interface
	751b  cRIO-9081
	751c  cRIO-9082
	7528  PXIe-4497
	7529  PXIe-4499
	752a  PXIe-4492
	7539  NI 9157
	753a  NI 9159
	7598  PXI-2571
	75a4  PXI-4131A
	75b1  PCIe-7854R
	75ba  PXI-2543
	75bb  PXIe-2543
	75e5  PXI-6683
	75e6  PXI-6683H
	75ef  PXIe-5632
	761c  VXI-8360LT
	761f  PXI-2540
	7620  PXIe-2540
	7621  PXI-2541
	7622  PXIe-2541
	7626  NI 9154
	7627  NI 9155
	7638  PXI-2720
	7639  PXI-2722
	763a  PXIe-2725
	763b  PXIe-2727
	763c  PXI-4465
	764b  PXIe-2790
	764c  PXI-2520
	764d  PXI-2521
	764e  PXI-2522
	764f  PXI-2523
	7654  PXI-2796
	7655  PXI-2797
	7656  PXI-2798
	7657  PXI-2799
	765d  PXI-2542
	765e  PXIe-2542
	765f  PXI-2544
	7660  PXIe-2544
	766d  PCIe-6535B
	766e  PCIe-6536B
	766f  PCIe-6537B
	76a3  PXIe-6535B
	76a4  PXIe-6536B
	76a5  PXIe-6537B
	783e  PXI-8368
	9020  PXI-2501
	9030  PXI-2503
	9040  PXI-2527
	9050  PXI-2565
	9060  PXI-2590
	9070  PXI-2591
	9080  PXI-2580
	9090  PCI-4021
	90a0  PXI-4021
	a001  PCI-MXI-2
	b001  PCI-1408
	b011  PXI-1408
	b021  PCI-1424
	b022  PXI-1424
	b031  PCI-1413
	b041  PCI-1407
	b051  PXI-1407
	b061  PCI-1411
	b071  PCI-1422
	b081  PXI-1422
	b091  PXI-1411
	b0b1  PCI-1409
	b0c1  PXI-1409
	b0e1  PCI-1428
	c4c4  PXIe/PCIe Device
	c801  PCI-GPIB
	c811  PCI-GPIB+
	c821  PXI-GPIB
	c831  PMC-GPIB
	c840  PCI-GPIB
	d130  PCI-232/2 Interface
	d140  PCI-232/4 Interface
	d150  PCI-232/8 Interface
	d160  PCI-485/2 Interface
	d170  PCI-485/4 Interface
	d190  PXI-8422/2 (Isolated RS-232) Interface
	d1a0  PXI-8422/4 (Isolated RS-232) Interface
	d1b0  PXI-8423/2 (Isolated RS-485) Interface
	d1c0  PXI-8423/4 (Isolated RS-485) Interface
	d1d0  PXI-8420/2 (RS-232) Interface
	d1e0  PXI-8420/4 (RS-232) Interface
	d1f0  PXI-8420/8 (RS-232) Interface
	d1f1  PXI-8420/16 (RS-232) Interface
	d230  PXI-8421/2 (RS-485) Interface
	d240  PXI-8421/4 (RS-485) Interface
	d250  PCI-232/2 (Isolated) Interface
	d260  PCI-485/2 (Isolated) Interface
	d270  PCI-232/4 (Isolated) Interface
	d280  PCI-485/4 (Isolated) Interface
	d290  PCI-485/8 Interface
	d2a0  PXI-8421/8 (RS-485) Interface
	d2b0  PCI-232/16 Interface
	e111  PCI-CAN
	e131  PXI-8461 (1 port)
	e141  PCI-CAN/LS
	e151  PXI-8460 (1 port)
	e211  PCI-CAN/2
	e231  PXI-8461 (2 ports)
	e241  PCI-CAN/LS2
	e251  PXI-8460 (2 ports)
	e261  PCI-CAN/DS
	e271  PXI-8462
	f110  VMEpc-650
	f120  VXIpc-650
	fe00  VXIpc-87x
	fe41  VXIpc-860
	fe51  VXIpc-74x
	fe61  VXIpc-850
	fe70  VXIpc-880
1094  First International Computers [FIC]
1095  Silicon Image, Inc.
	0240  Adaptec AAR-1210SA SATA HostRAID Controller
	0640  PCI0640
	0643  PCI0643
	0646  PCI0646
	0647  PCI0647
	0648  PCI0648
	0649  SiI 0649 Ultra ATA/100 PCI to ATA Host Controller
	0650  PBC0650A
	0670  USB0670
	0673  USB0673
	0680  PCI0680 Ultra ATA-133 Host Controller
	3112  SiI 3112 [SATALink/SATARaid] Serial ATA Controller
	3114  SiI 3114 [SATALink/SATARaid] Serial ATA Controller
	3124  SiI 3124 PCI-X Serial ATA Controller
	3132  SiI 3132 Serial ATA Raid II Controller
	3512  SiI 3512 [SATALink/SATARaid] Serial ATA Controller
	3531  SiI 3531 [SATALink/SATARaid] Serial ATA Controller
1096  Alacron
1097  Appian Technology
1098  Quantum Designs (H.K.) Ltd
	0001  QD-8500
	0002  QD-8580
1099  Samsung Electronics Co., Ltd
109a  Packard Bell
109b  Gemlight Computer Ltd.
109c  Megachips Corporation
109d  Zida Technologies Ltd.
109e  Brooktree Corporation
	0310  Bt848 Video Capture
	032e  Bt878 Video Capture
	0350  Bt848 Video Capture
	0351  Bt849A Video capture
	0369  Bt878 Video Capture
	036c  Bt879(??) Video Capture
	036e  Bt878 Video Capture
	036f  Bt879 Video Capture
	0370  Bt880 Video Capture
	0878  Bt878 Audio Capture
	0879  Bt879 Audio Capture
	0880  Bt880 Audio Capture
	2115  BtV 2115 Mediastream controller
	2125  BtV 2125 Mediastream controller
	2164  BtV 2164
	2165  BtV 2165
	8230  Bt8230 ATM Segment/Reassembly Ctrlr (SRC)
	8472  Bt8472
	8474  Bt8474
109f  Trigem Computer Inc.
10a0  Meidensha Corporation
10a1  Juko Electronics Ind. Co. Ltd
10a2  Quantum Corporation
10a3  Everex Systems Inc
10a4  Globe Manufacturing Sales
10a5  Smart Link Ltd.
	3052  SmartPCI562 56K Modem
	5449  SmartPCI561 modem
10a6  Informtech Industrial Ltd.
10a7  Benchmarq Microelectronics
10a8  Sierra Semiconductor
	0000  STB Horizon 64
10a9  Silicon Graphics Intl. Corp.
	0001  Crosstalk to PCI Bridge
	0002  Linc I/O controller
	0003  IOC3 I/O controller
	0004  O2 MACE
	0005  RAD Audio
	0006  HPCEX
	0007  RPCEX
	0008  DiVO VIP
	0009  AceNIC Gigabit Ethernet
	0010  AMP Video I/O
	0011  GRIP
	0012  SGH PSHAC GSN
	0208  SSIM1 SAS Adapter
	1001  Magic Carpet
	1002  Lithium
	1003  Dual JPEG 1
	1004  Dual JPEG 2
	1005  Dual JPEG 3
	1006  Dual JPEG 4
	1007  Dual JPEG 5
	1008  Cesium
	100a  IOC4 I/O controller
	1504  SSIM1 Fibre Channel Adapter
	2001  Fibre Channel
	2002  ASDE
	4001  TIO-CE PCI Express Bridge
	4002  TIO-CE PCI Express Port
	8001  O2 1394
	8002  G-net NT
	802b  REACT external interrupt controller
10aa  ACC Microelectronics
	0000  ACCM 2188
	2051  2051 CPU bridge
	5842  2051 ISA bridge
10ab  Digicom
10ac  Honeywell IAC
10ad  Symphony Labs
	0001  W83769F
	0003  SL82C103
	0005  SL82C105
	0103  SL82c103
	0105  SL82c105
	0565  W83C553F/W83C554F
10ae  Cornerstone Technology
10af  Micro Computer Systems Inc
10b0  CardExpert Technology
10b1  Cabletron Systems Inc
10b2  Raytheon Company
10b3  Databook Inc
	3106  DB87144
	b106  DB87144
10b4  STB Systems Inc
	1b1d  Velocity 128 3D
10b5  PLX Technology, Inc.
	0001  i960 PCI bus interface
	0557  PCI9030 32-bit 33MHz PCI <-> IOBus Bridge
	1000  PCI9030 32-bit 33MHz PCI <-> IOBus Bridge
	1024  Acromag, Inc. IndustryPack Carrier Card
	1042  Brandywine / jxi2, Inc. - PMC-SyncClock32, IRIG A & B, Nasa 36
	106a  Dual OX16C952 4 port serial adapter [Megawolf Romulus/4]
	1076  VScom 800 8 port serial adaptor
	1077  VScom 400 4 port serial adaptor
	1078  VScom 210 2 port serial and 1 port parallel adaptor
	1103  VScom 200 2 port serial adaptor
	1146  VScom 010 1 port parallel adaptor
	1147  VScom 020 2 port parallel adaptor
	2000  PCI9030 32-bit 33MHz PCI <-> IOBus Bridge
	2540  IXXAT CAN-Interface PC-I 04/PCI
	2724  Thales PCSM Security Card
	3376  Cosateq 4 Port CAN Card
	4000  PCI9030 32-bit 33MHz PCI <-> IOBus Bridge
	4001  PCI9030 32-bit 33MHz PCI <-> IOBus Bridge
	4002  PCI9030 32-bit 33MHz PCI <-> IOBus Bridge
	6140  PCI6140 32-bit 33MHz PCI-to-PCI Bridge
	6150  PCI6150 32-bit 33MHz PCI-to-PCI Bridge
	6152  PCI6152 32-bit 66MHz PCI-to-PCI Bridge
	6154  PCI6154 64-bit 66MHz PCI-to-PCI Bridge
	6254  PCI6254 64-bit 66MHz PCI-to-PCI Bridge
	6466  PCI6466 64-bit 66MHz PCI-to-PCI Bridge
	6520  PCI6520 64-bit 133MHz PCI-X-to-PCI-X Bridge
	6540  PCI6540 64-bit 133MHz PCI-X-to-PCI-X Bridge
	6541  PCI6540/6466 PCI-PCI bridge (non-transparent mode, primary side)
	6542  PCI6540/6466 PCI-PCI bridge (non-transparent mode, secondary side)
	8111  PEX 8111 PCI Express-to-PCI Bridge
	8112  PEX8112 x1 Lane PCI Express-to-PCI Bridge
	8114  PEX 8114 PCI Express-to-PCI/PCI-X Bridge
	8311  PEX8311 x1 Lane PCI Express-to-Generic Local Bus Bridge
	8505  PEX 8505 5-lane, 5-port PCI Express Switch
	8508  PEX 8508 8-lane, 5-port PCI Express Switch
	8509  PEX 8509 8-lane, 8-port PCI Express Switch
	8512  PEX 8512 12-lane, 5-port PCI Express Switch
	8516  PEX 8516  Versatile PCI Express Switch
	8517  PEX 8517 16-lane, 5-port PCI Express Switch
	8518  PEX 8518 16-lane, 5-port PCI Express Switch
	8524  PEX 8524 24-lane, 6-port PCI Express Switch
	8525  PEX 8525 24-lane, 5-port PCI Express Switch
	8532  PEX 8532  Versatile PCI Express Switch
	8533  PEX 8533 32-lane, 6-port PCI Express Switch
	8547  PEX 8547 48-lane, 3-port PCI Express Switch
	8548  PEX 8548 48-lane, 9-port PCI Express Switch
	8603  PEX 8603 3-lane, 3-Port PCI Express Gen 2 (5.0 GT/s) Switch
	8604  PEX 8604 4-lane, 4-Port PCI Express Gen 2 (5.0 GT/s) Switch
	8605  PEX 8605 PCI Express 4-port Gen2 Switch
	8606  PEX 8606 6 Lane, 6 Port PCI Express Gen 2 (5.0 GT/s) Switch
	8608  PEX 8608 8-lane, 8-Port PCI Express Gen 2 (5.0 GT/s) Switch
	8609  PEX 8609 8-lane, 8-Port PCI Express Gen 2 (5.0 GT/s) Switch with DMA
	8612  PEX 8612 12-lane, 4-Port PCI Express Gen 2 (5.0 GT/s) Switch
	8613  PEX 8613 12-lane, 3-Port PCI Express Gen 2 (5.0 GT/s) Switch
	8614  PEX 8614 12-lane, 12-Port PCI Express Gen 2 (5.0 GT/s) Switch
	8615  PEX 8615 12-lane, 12-Port PCI Express Gen 2 (5.0 GT/s) Switch with DMA
	8616  PEX 8616 16-lane, 4-Port PCI Express Gen 2 (5.0 GT/s) Switch
	8617  PEX 8617 16-lane, 4-Port PCI Express Gen 2 (5.0 GT/s) Switch with P2P
	8618  PEX 8618 16-lane, 16-Port PCI Express Gen 2 (5.0 GT/s) Switch
	8619  PEX 8619 16-lane, 16-Port PCI Express Gen 2 (5.0 GT/s) Switch with DMA
	8624  PEX 8624 24-lane, 6-Port PCI Express Gen 2 (5.0 GT/s) Switch [ExpressLane]
	8625  PEX 8625 24-lane, 24-Port PCI Express Gen 2 (5.0 GT/s) Switch
	8632  PEX 8632 32-lane, 12-Port PCI Express Gen 2 (5.0 GT/s) Switch
	8636  PEX 8636 36-lane, 24-Port PCI Express Gen 2 (5.0 GT/s) Switch
	8647  PEX 8647 48-Lane, 3-Port PCI Express Gen 2 (5.0 GT/s) Switch
	8648  PEX 8648 48-lane, 12-Port PCI Express Gen 2 (5.0 GT/s) Switch
	8649  PEX 8649 48-lane, 12-Port PCI Express Gen 2 (5.0 GT/s) Switch
	8664  PEX 8664 64-lane, 16-Port PCI Express Gen 2 (5.0 GT/s) Switch
	8680  PEX 8680 80-lane, 20-Port PCI Express Gen 2 (5.0 GT/s) Multi-Root Switch
	8696  PEX 8696 96-lane, 24-Port PCI Express Gen 2 (5.0 GT/s) Multi-Root Switch
	8717  PEX 8717 16-lane, 8-Port PCI Express Gen 3 (8.0 GT/s) Switch with DMA
	8718  PEX 8718 16-Lane, 5-Port PCI Express Gen 3 (8.0 GT/s) Switch
	8724  PEX 8724 24-Lane, 6-Port PCI Express Gen 3 (8 GT/s) Switch, 19 x 19mm FCBGA
	8732  PEX 8732 32-lane, 8-Port PCI Express Gen 3 (8.0 GT/s) Switch
	8734  PEX 8734 32-lane, 8-Port PCI Express Gen 3 (8.0GT/s) Switch
	8747  PEX 8747 48-Lane, 5-Port PCI Express Gen 3 (8.0 GT/s) Switch
	8748  PEX 8748 48-Lane, 12-Port PCI Express Gen 3 (8 GT/s) Switch, 27 x 27mm FCBGA
	87b0  PEX 8732 32-lane, 8-Port PCI Express Gen 3 (8.0 GT/s) Switch
	9016  PLX 9016 8-port serial controller
	9030  PCI9030 32-bit 33MHz PCI <-> IOBus Bridge
	9036  9036
	9050  PCI <-> IOBus Bridge
	9052  PCI9052 PCI <-> IOBus Bridge
	9054  PCI9054 32-bit 33MHz PCI <-> IOBus Bridge
	9056  PCI9056 32-bit 66MHz PCI <-> IOBus Bridge
	9060  PCI9060 32-bit 33MHz PCI <-> IOBus Bridge
	906d  9060SD
	906e  9060ES
	9080  PCI9080 32-bit; 33MHz PCI <-> IOBus Bridge
	9656  PCI9656 PCI <-> IOBus Bridge
	9733  PEX 9733 33-lane, 9-port PCI Express Gen 3 (8.0 GT/s) Switch
	9749  PEX 9749 49-lane, 13-port PCI Express Gen 3 (8.0 GT/s) Switch
	a100  Blackmagic Design DeckLink
	bb04  B&B 3PCIOSD1A Isolated PCI Serial
	c001  CronyxOmega-PCI (8-port RS232)
	d00d  PCI9030 32-bit 33MHz PCI <-> IOBus Bridge
	d33d  PCI9030 32-bit 33MHz PCI <-> IOBus Bridge
	d44d  PCI9030 32-bit 33MHz PCI <-> IOBus Bridge
10b6  Madge Networks
	0001  Smart 16/4 PCI Ringnode
	0002  Smart 16/4 PCI Ringnode Mk2
	0003  Smart 16/4 PCI Ringnode Mk3
	0004  Smart 16/4 PCI Ringnode Mk1
	0006  16/4 Cardbus Adapter
	0007  Presto PCI Adapter
	0009  Smart 100/16/4 PCI-HS Ringnode
	000a  Token Ring 100/16/4 Ringnode/Ringrunner
	000b  16/4 CardBus Adapter Mk2
	000c  RapidFire 3140V2 16/4 TR Adapter
	1000  Collage 25/155 ATM Client Adapter
	1001  Collage 155 ATM Server Adapter
10b7  3Com Corporation
	0001  3c985 1000BaseSX (SX/TX)
	0013  AR5212 802.11abg NIC (3CRDAG675)
	0910  3C910-A01
	1006  MINI PCI type 3B Data Fax Modem
	1007  Mini PCI 56k Winmodem
	1201  3c982-TXM 10/100baseTX Dual Port A [Hydra]
	1202  3c982-TXM 10/100baseTX Dual Port B [Hydra]
	1700  3c940 10/100/1000Base-T [Marvell]
	3390  3c339 TokenLink Velocity
	3590  3c359 TokenLink Velocity XL
	4500  3c450 HomePNA [Tornado]
	5055  3c555 Laptop Hurricane
	5057  3c575 Megahertz 10/100 LAN CardBus [Boomerang]
	5157  3cCFE575BT Megahertz 10/100 LAN CardBus [Cyclone]
	5257  3cCFE575CT CardBus [Cyclone]
	5900  3c590 10BaseT [Vortex]
	5920  3c592 EISA 10mbps Demon/Vortex
	5950  3c595 100BaseTX [Vortex]
	5951  3c595 100BaseT4 [Vortex]
	5952  3c595 100Base-MII [Vortex]
	5970  3c597 EISA Fast Demon/Vortex
	5b57  3c595 Megahertz 10/100 LAN CardBus [Boomerang]
	6000  3CRSHPW796 [OfficeConnect Wireless CardBus]
	6001  3com 3CRWE154G72 [Office Connect Wireless LAN Adapter]
	6055  3c556 Hurricane CardBus [Cyclone]
	6056  3c556B CardBus [Tornado]
	6560  3cCFE656 CardBus [Cyclone]
	6561  3cCFEM656 10/100 LAN+56K Modem CardBus
	6562  3cCFEM656B 10/100 LAN+Winmodem CardBus [Cyclone]
	6563  3cCFEM656B 10/100 LAN+56K Modem CardBus
	6564  3cXFEM656C 10/100 LAN+Winmodem CardBus [Tornado]
	7646  3cSOHO100-TX Hurricane
	7770  3CRWE777 PCI Wireless Adapter [Airconnect]
	7940  3c803 FDDILink UTP Controller
	7980  3c804 FDDILink SAS Controller
	7990  3c805 FDDILink DAS Controller
	80eb  3c940B 10/100/1000Base-T
	8811  Token ring
	9000  3c900 10BaseT [Boomerang]
	9001  3c900 10Mbps Combo [Boomerang]
	9004  3c900B-TPO Etherlink XL [Cyclone]
	9005  3c900B-Combo Etherlink XL [Cyclone]
	9006  3c900B-TPC Etherlink XL [Cyclone]
	900a  3c900B-FL 10base-FL [Cyclone]
	9050  3c905 100BaseTX [Boomerang]
	9051  3c905 100BaseT4 [Boomerang]
	9054  3C905B-TX Fast Etherlink XL PCI
	9055  3c905B 100BaseTX [Cyclone]
	9056  3c905B-T4 Fast EtherLink XL [Cyclone]
	9058  3c905B Deluxe Etherlink 10/100/BNC [Cyclone]
	905a  3c905B-FX Fast Etherlink XL FX 100baseFx [Cyclone]
	9200  3c905C-TX/TX-M [Tornado]
	9201  3C920B-EMB Integrated Fast Ethernet Controller [Tornado]
	9202  3Com 3C920B-EMB-WNM Integrated Fast Ethernet Controller
	9210  3C920B-EMB-WNM Integrated Fast Ethernet Controller
	9300  3CSOHO100B-TX 910-A01 [tulip]
	9800  3c980-TX Fast Etherlink XL Server Adapter [Cyclone]
	9805  3c980-C 10/100baseTX NIC [Python-T]
	9900  3C990-TX [Typhoon]
	9902  3CR990-TX-95 [Typhoon 56-bit]
	9903  3CR990-TX-97 [Typhoon 168-bit]
	9904  3C990B-TX-M/3C990BSVR [Typhoon2]
	9905  3CR990-FX-95/97/95 [Typhon Fiber]
	9908  3CR990SVR95 [Typhoon Server 56-bit]
	9909  3CR990SVR97 [Typhoon Server 168-bit]
	990a  3C990SVR [Typhoon Server]
	990b  3C990SVR [Typhoon Server]
10b8  Standard Microsystems Corp [SMC]
	0005  83c170 EPIC/100 Fast Ethernet Adapter
	0006  83c175 EPIC/100 Fast Ethernet Adapter
	1000  FDC 37c665
	1001  FDC 37C922
	a011  83C170QF
	b106  SMC34C90
10b9  ULi Electronics Inc.
	0101  CMI8338/C3DX PCI Audio Device
	0111  C-Media CMI8738/C3DX Audio Device (OEM)
	0780  Multi-IO Card
	0782  Multi-IO Card
	1435  M1435
	1445  M1445
	1449  M1449
	1451  M1451
	1461  M1461
	1489  M1489
	1511  M1511 [Aladdin]
	1512  M1512 [Aladdin]
	1513  M1513 [Aladdin]
	1521  M1521 [Aladdin III]
	1523  M1523
	1531  M1531 [Aladdin IV]
	1533  M1533/M1535/M1543 PCI to ISA Bridge [Aladdin IV/V/V+]
	1541  M1541
	1543  M1543
	1563  M1563 HyperTransport South Bridge
	1573  PCI to LPC Controller
	1575  M1575 South Bridge
	1621  M1621
	1631  ALI M1631 PCI North Bridge Aladdin Pro III
	1632  M1632M Northbridge+Trident
	1641  ALI M1641 PCI North Bridge Aladdin Pro IV
	1644  M1644/M1644T Northbridge+Trident
	1646  M1646 Northbridge+Trident
	1647  M1647 Northbridge [MAGiK 1 / MobileMAGiK 1]
	1651  M1651/M1651T Northbridge [Aladdin-Pro 5/5M,Aladdin-Pro 5T/5TM]
	1671  M1671 Super P4 Northbridge [AGP4X,PCI and SDR/DDR]
	1672  M1672 Northbridge [CyberALADDiN-P4]
	1681  M1681 P4 Northbridge [AGP8X,HyperTransport and SDR/DDR]
	1687  M1687 K8 Northbridge [AGP8X and HyperTransport]
	1689  M1689 K8 Northbridge [Super K8 Single Chip]
	1695  M1695 Host Bridge
	1697  M1697 HTT Host Bridge
	3141  M3141
	3143  M3143
	3145  M3145
	3147  M3147
	3149  M3149
	3151  M3151
	3307  M3307
	3309  M3309
	3323  M3325 Video/Audio Decoder
	5212  M4803
	5215  MS4803
	5217  M5217H
	5219  M5219
	5225  M5225
	5228  M5228 ALi ATA/RAID Controller
	5229  M5229 IDE
	5235  M5225
	5237  USB 1.1 Controller
	5239  USB 2.0 Controller
	5243  M1541 PCI to AGP Controller
	5246  AGP8X Controller
	5247  PCI to AGP Controller
	5249  M5249 HTT to PCI Bridge
	524b  PCI Express Root Port
	524c  PCI Express Root Port
	524d  PCI Express Root Port
	524e  PCI Express Root Port
	5251  M5251 P1394 OHCI 1.0 Controller
	5253  M5253 P1394 OHCI 1.1 Controller
	5261  M5261 Ethernet Controller
	5263  ULi 1689,1573 integrated ethernet.
	5281  ALi M5281 Serial ATA / RAID Host Controller
	5287  ULi 5287 SATA
	5288  ULi M5288 SATA
	5289  ULi 5289 SATA
	5450  Lucent Technologies Soft Modem AMR
	5451  M5451 PCI AC-Link Controller Audio Device
	5453  M5453 PCI AC-Link Controller Modem Device
	5455  M5455 PCI AC-Link Controller Audio Device
	5457  M5457 AC'97 Modem Controller
	5459  SmartLink SmartPCI561 56K Modem
	545a  SmartLink SmartPCI563 56K Modem
	5461  HD Audio Controller
	5471  M5471 Memory Stick Controller
	5473  M5473 SD-MMC Controller
	7101  M7101 Power Management Controller [PMU]
10ba  Mitsubishi Electric Corp.
	0301  AccelGraphics AccelECLIPSE
	0304  AccelGALAXY A2100 [OEM Evans & Sutherland]
	0308  Tornado 3000 [OEM Evans & Sutherland]
	1002  VG500 [VolumePro Volume Rendering Accelerator]
10bb  Dapha Electronics Corporation
10bc  Advanced Logic Research
10bd  Surecom Technology
	0e34  NE-34
10be  Tseng Labs International Co.
10bf  Most Inc
10c0  Boca Research Inc.
10c1  ICM Co., Ltd.
10c2  Auspex Systems Inc.
10c3  Samsung Semiconductors, Inc.
10c4  Award Software International Inc.
10c5  Xerox Corporation
10c6  Rambus Inc.
10c7  Media Vision
10c8  Neomagic Corporation
	0001  NM2070 [MagicGraph 128]
	0002  NM2090 [MagicGraph 128V]
	0003  NM2093 [MagicGraph 128ZV]
	0004  NM2160 [MagicGraph 128XD]
	0005  NM2200 [MagicGraph 256AV]
	0006  NM2360 [MagicMedia 256ZX]
	0016  NM2380 [MagicMedia 256XL+]
	0025  NM2230 [MagicGraph 256AV+]
	0083  NM2093 [MagicGraph 128ZV+]
	8005  NM2200 [MagicMedia 256AV Audio]
	8006  NM2360 [MagicMedia 256ZX Audio]
	8016  NM2380 [MagicMedia 256XL+ Audio]
10c9  Dataexpert Corporation
10ca  Fujitsu Microelectr., Inc.
10cb  Omron Corporation
10cc  Mai Logic Incorporated
	0660  Articia S Host Bridge
	0661  Articia S PCI Bridge
10cd  Advanced System Products, Inc
	1100  ASC1100
	1200  ASC1200 [(abp940) Fast SCSI-II]
	1300  ASC1300 / ASC3030 [ABP940-U / ABP960-U / ABP3925]
	2300  ABP940-UW
	2500  ABP940-U2W
	2700  ABP3950-U3W
10ce  Radius
10cf  Fujitsu Limited.
	01ef  PCEA4 PCI-Express Dual Port ESCON Adapter
	1414  On-board USB 1.1 companion controller
	1415  On-board USB 2.0 EHCI controller
	1422  E8410 nVidia graphics adapter
	142d  HD audio (Realtek ALC262)
	1430  82566MM Intel 1Gb copper LAN interface
	1623  PCEA4 PCI-Express Dual Port ESCON Adapter
	2001  mb86605
	200c  MB86613L IEEE1394 OHCI 1.0 Controller
	2010  MB86613S IEEE1394 OHCI 1.1 Controller
	2019  MB86295S [CORAL P]
	201e  MB86296S [CORAL PA]
	202b  MB86297A [Carmine Graphics Controller]
10d1  FuturePlus Systems Corp.
10d2  Molex Incorporated
10d3  Jabil Circuit Inc
10d4  Hualon Microelectronics
10d5  Autologic Inc.
10d6  Cetia
10d7  BCM Advanced Research
10d8  Advanced Peripherals Labs
10d9  Macronix, Inc. [MXIC]
	0431  MX98715
	0512  MX98713
	0531  MX987x5
	8625  MX86250
	8626  Macronix MX86251 + 3Dfx Voodoo Rush
	8888  MX86200
10da  Compaq IPG-Austin
	0508  TC4048 Token Ring 4/16
	3390  Tl3c3x9
10db  Rohm LSI Systems, Inc.
10dc  CERN/ECP/EDU
	0001  STAR/RD24 SCI-PCI (PMC)
	0002  TAR/RD24 SCI-PCI (PMC)
	0021  HIPPI destination
	0022  HIPPI source
	10dc  ATT2C15-3 FPGA
10dd  Evans & Sutherland
	0100  Lightning 1200
10de  NVIDIA Corporation
	0008  NV1 [EDGE 3D]
	0009  NV1 [EDGE 3D]
	0020  NV4 [Riva TNT]
	0028  NV5 [Riva TNT2 / TNT2 Pro]
	0029  NV5 [Riva TNT2 Ultra]
	002a  NV5 [Riva TNT2]
	002b  NV5 [Riva TNT2]
	002c  NV5 [Vanta / Vanta LT]
	002d  NV5 [Riva TNT2 Model 64 / Model 64 Pro]
	0034  MCP04 SMBus
	0035  MCP04 IDE
	0036  MCP04 Serial ATA Controller
	0037  MCP04 Ethernet Controller
	0038  MCP04 Ethernet Controller
	003a  MCP04 AC'97 Audio Controller
	003b  MCP04 USB Controller
	003c  MCP04 USB Controller
	003d  MCP04 PCI Bridge
	003e  MCP04 Serial ATA Controller
	0040  NV40 [GeForce 6800 Ultra]
	0041  NV40 [GeForce 6800]
	0042  NV40 [GeForce 6800 LE]
	0043  NV40 [GeForce 6800 XE]
	0044  NV40 [GeForce 6800 XT]
	0045  NV40 [GeForce 6800 GT]
	0047  NV40 [GeForce 6800 GS]
	0048  NV40 [GeForce 6800 XT]
	004e  NV40GL [Quadro FX 4000]
	0050  CK804 ISA Bridge
	0051  CK804 ISA Bridge
	0052  CK804 SMBus
	0053  CK804 IDE
	0054  CK804 Serial ATA Controller
	0055  CK804 Serial ATA Controller
	0056  CK804 Ethernet Controller
	0057  CK804 Ethernet Controller
	0058  CK804 AC'97 Modem
	0059  CK804 AC'97 Audio Controller
	005a  CK804 USB Controller
	005b  CK804 USB Controller
	005c  CK804 PCI Bridge
	005d  CK804 PCIE Bridge
	005e  CK804 Memory Controller
	005f  CK804 Memory Controller
	0060  nForce2 ISA Bridge
	0064  nForce2 SMBus (MCP)
	0065  nForce2 IDE
	0066  nForce2 Ethernet Controller
	0067  nForce2 USB Controller
	0068  nForce2 USB Controller
	006a  nForce2 AC97 Audio Controler (MCP)
	006b  nForce Audio Processing Unit
	006c  nForce2 External PCI Bridge
	006d  nForce2 PCI Bridge
	006e  nForce2 FireWire (IEEE 1394) Controller
	0080  MCP2A ISA bridge
	0084  MCP2A SMBus
	0085  MCP2A IDE
	0086  MCP2A Ethernet Controller
	0087  MCP2A USB Controller
	0088  MCP2A USB Controller
	008a  MCP2S AC'97 Audio Controller
	008b  MCP2A PCI Bridge
	008c  MCP2A Ethernet Controller
	008e  nForce2 Serial ATA Controller
	0090  G70 [GeForce 7800 GTX]
	0091  G70 [GeForce 7800 GTX]
	0092  G70 [GeForce 7800 GT]
	0093  G70 [GeForce 7800 GS]
	0095  G70 [GeForce 7800 SLI]
	0097  G70 [GeForce GTS 250]
	0098  G70M [GeForce Go 7800]
	0099  G70M [GeForce Go 7800 GTX]
	009d  G70GL [Quadro FX 4500]
	00a0  NV5 [Aladdin TNT2]
	00c0  NV41 [GeForce 6800 GS]
	00c1  NV41 [GeForce 6800]
	00c2  NV41 [GeForce 6800 LE]
	00c3  NV41 [GeForce 6800 XT]
	00c5  NV41
	00c6  NV41
	00c7  NV41
	00c8  NV41M [GeForce Go 6800]
	00c9  NV41M [GeForce Go 6800 Ultra]
	00cc  NV41GLM [Quadro FX Go1400]
	00cd  NV42GL [Quadro FX 3450/4000 SDI]
	00ce  NV41GL [Quadro FX 1400]
	00cf  NV41
	00d0  nForce3 LPC Bridge
	00d1  nForce3 Host Bridge
	00d2  nForce3 AGP Bridge
	00d3  CK804 Memory Controller
	00d4  nForce3 SMBus
	00d5  nForce3 IDE
	00d6  nForce3 Ethernet
	00d7  nForce3 USB 1.1
	00d8  nForce3 USB 2.0
	00d9  nForce3 Audio
	00da  nForce3 Audio
	00dd  nForce3 PCI Bridge
	00df  CK8S Ethernet Controller
	00e0  nForce3 250Gb LPC Bridge
	00e1  nForce3 250Gb Host Bridge
	00e2  nForce3 250Gb AGP Host to PCI Bridge
	00e3  nForce3 Serial ATA Controller
	00e4  nForce 250Gb PCI System Management
	00e5  CK8S Parallel ATA Controller (v2.5)
	00e6  CK8S Ethernet Controller
	00e7  CK8S USB Controller
	00e8  nForce3 EHCI USB 2.0 Controller
	00ea  nForce3 250Gb AC'97 Audio Controller
	00ed  nForce3 250Gb PCI-to-PCI Bridge
	00ee  nForce3 Serial ATA Controller 2
	00f1  NV43 [GeForce 6600 GT]
	00f2  NV43 [GeForce 6600]
	00f3  NV43 [GeForce 6200]
	00f4  NV43 [GeForce 6600 LE]
	00f5  G71 [GeForce 7800 GS]
	00f6  NV43 [GeForce 6800 GS/XT]
	00f8  NV40GL [Quadro FX 3400/4400]
	00f9  NV40 [GeForce 6800 GT/GTO/Ultra]
	00fa  NV36 [GeForce PCX 5750]
	00fb  NV38 [GeForce PCX 5900]
	00fc  NV37GL [Quadro FX 330/GeForce PCX 5300]
	00fd  NV37GL [Quadro PCI-E Series]
	00fe  NV38GL [Quadro FX 1300]
	00ff  NV18 [GeForce PCX 4300]
	0100  NV10 [GeForce 256 SDR]
	0101  NV10 [GeForce 256 DDR]
	0103  NV10GL [Quadro]
	0110  NV11 [GeForce2 MX/MX 400]
	0111  NV11 [GeForce2 MX200]
	0112  NV11M [GeForce2 Go]
	0113  NV11GL [Quadro2 MXR/EX/Go]
	0140  NV43 [GeForce 6600 GT]
	0141  NV43 [GeForce 6600]
	0142  NV43 [GeForce 6600 LE]
	0143  NV43 [GeForce 6600 VE]
	0144  NV43M [GeForce Go 6600]
	0145  NV43 [GeForce 6610 XL]
	0146  NV43M [GeForce Go6200 TE / 6600 TE]
	0147  NV43 [GeForce 6700 XL]
	0148  NV43M [GeForce Go 6600]
	0149  NV43M [GeForce Go 6600 GT]
	014a  NV43 [Quadro NVS 440]
	014b  NV43
	014d  NV43GL [Quadro FX 550]
	014e  NV43GL [Quadro FX 540]
	014f  NV43 [GeForce 6200]
	0150  NV15 [GeForce2 GTS/Pro]
	0151  NV15 [GeForce2 Ti]
	0152  NV15 [GeForce2 Ultra]
	0153  NV15GL [Quadro2 Pro]
	0160  NV44 [GeForce 6500]
	0161  NV44 [GeForce 6200 TurboCache]
	0162  NV44 [GeForce 6200 SE TurboCache]
	0163  NV44 [GeForce 6200 LE]
	0164  NV44M [GeForce Go 6200]
	0165  NV44 [Quadro NVS 285]
	0166  NV44M [GeForce Go 6400]
	0167  NV44M [GeForce Go 6200]
	0168  NV44M [GeForce Go 6400]
	0169  NV44 [GeForce 6250]
	016a  NV44 [GeForce 7100 GS]
	016d  NV44
	016e  NV44
	016f  NV44
	0170  NV17 [GeForce4 MX 460]
	0171  NV17 [GeForce4 MX 440]
	0172  NV17 [GeForce4 MX 420]
	0173  NV17 [GeForce4 MX 440-SE]
	0174  NV17M [GeForce4 440 Go]
	0175  NV17M [GeForce4 420 Go]
	0176  NV17M [GeForce4 420 Go 32M]
	0177  NV17M [GeForce4 460 Go]
	0178  NV17GL [Quadro4 550 XGL]
	0179  NV17M [GeForce4 440 Go 64M]
	017a  NV17GL [Quadro NVS]
	017b  NV17GL [Quadro4 550 XGL]
	017c  NV17GL [Quadro4 500 GoGL]
	017f  NV17
	0181  NV18 [GeForce4 MX 440 AGP 8x]
	0182  NV18 [GeForce4 MX 440SE AGP 8x]
	0183  NV18 [GeForce4 MX 420 AGP 8x]
	0184  NV18 [GeForce4 MX]
	0185  NV18 [GeForce4 MX 4000]
	0186  NV18M [GeForce4 448 Go]
	0187  NV18M [GeForce4 488 Go]
	0188  NV18GL [Quadro4 580 XGL]
	0189  NV18 [GeForce4 MX with AGP8X (Mac)]
	018a  NV18GL [Quadro NVS 280 SD]
	018b  NV18GL [Quadro4 380 XGL]
	018c  NV18GL [Quadro NVS 50 PCI]
	018d  NV18M [GeForce4 448 Go]
	018f  NV18
	0190  G80 [GeForce 8800 GTS / 8800 GTX]
	0191  G80 [GeForce 8800 GTX]
	0192  G80 [GeForce 8800 GTS]
	0193  G80 [GeForce 8800 GTS]
	0194  G80 [GeForce 8800 Ultra]
	0197  G80GL [Tesla C870]
	019d  G80GL [Quadro FX 5600]
	019e  G80GL [Quadro FX 4600]
	01a0  nForce 220/420 NV11 [GeForce2 MX]
	01a4  nForce CPU bridge
	01ab  nForce 420 Memory Controller (DDR)
	01ac  nForce 220/420 Memory Controller
	01ad  nForce 220/420 Memory Controller
	01b0  nForce Audio Processing Unit
	01b1  nForce AC'97 Audio Controller
	01b2  nForce ISA Bridge
	01b4  nForce PCI System Management
	01b7  nForce AGP to PCI Bridge
	01b8  nForce PCI-to-PCI bridge
	01bc  nForce IDE
	01c1  nForce AC'97 Modem Controller
	01c2  nForce USB Controller
	01c3  nForce Ethernet Controller
	01d0  G72 [GeForce 7350 LE]
	01d1  G72 [GeForce 7300 LE]
	01d2  G72 [GeForce 7550 LE]
	01d3  G72 [GeForce 7200 GS / 7300 SE]
	01d5  G72
	01d6  G72M [GeForce Go 7200]
	01d7  G72M [Quadro NVS 110M/GeForce Go 7300]
	01d8  G72M [GeForce Go 7400]
	01d9  G72M [GeForce Go 7450]
	01da  G72M [Quadro NVS 110M]
	01db  G72M [Quadro NVS 120M]
	01dc  G72GLM [Quadro FX 350M]
	01dd  G72 [GeForce 7500 LE]
	01de  G72GL [Quadro FX 350]
	01df  G72 [GeForce 7300 GS]
	01e0  nForce2 IGP2
	01e8  nForce2 AGP
	01ea  nForce2 Memory Controller 0
	01eb  nForce2 Memory Controller 1
	01ec  nForce2 Memory Controller 2
	01ed  nForce2 Memory Controller 3
	01ee  nForce2 Memory Controller 4
	01ef  nForce2 Memory Controller 5
	01f0  C17 [GeForce4 MX IGP]
	0200  NV20 [GeForce3]
	0201  NV20 [GeForce3 Ti 200]
	0202  NV20 [GeForce3 Ti 500]
	0203  NV20GL [Quadro DCC]
	0211  NV48 [GeForce 6800]
	0212  NV48 [GeForce 6800 LE]
	0215  NV48 [GeForce 6800 GT]
	0218  NV48 [GeForce 6800 XT]
	0221  NV44A [GeForce 6200]
	0222  NV44 [GeForce 6200 A-LE]
	0224  NV44
	0240  C51PV [GeForce 6150]
	0241  C51 [GeForce 6150 LE]
	0242  C51G [GeForce 6100]
	0243  C51 PCI Express Bridge
	0244  C51 [GeForce Go 6150]
	0245  C51 [Quadro NVS 210S/GeForce 6150LE]
	0246  C51 PCI Express Bridge
	0247  C51 [GeForce Go 6100]
	0248  C51 PCI Express Bridge
	0249  C51 PCI Express Bridge
	024a  C51 PCI Express Bridge
	024b  C51 PCI Express Bridge
	024c  C51 PCI Express Bridge
	024d  C51 PCI Express Bridge
	024e  C51 PCI Express Bridge
	024f  C51 PCI Express Bridge
	0250  NV25 [GeForce4 Ti 4600]
	0251  NV25 [GeForce4 Ti 4400]
	0252  NV25 [GeForce4 Ti]
	0253  NV25 [GeForce4 Ti 4200]
	0258  NV25GL [Quadro4 900 XGL]
	0259  NV25GL [Quadro4 750 XGL]
	025b  NV25GL [Quadro4 700 XGL]
	0260  MCP51 LPC Bridge
	0261  MCP51 LPC Bridge
	0262  MCP51 LPC Bridge
	0263  MCP51 LPC Bridge
	0264  MCP51 SMBus
	0265  MCP51 IDE
	0266  MCP51 Serial ATA Controller
	0267  MCP51 Serial ATA Controller
	0268  MCP51 Ethernet Controller
	0269  MCP51 Ethernet Controller
	026a  MCP51 MCI
	026b  MCP51 AC97 Audio Controller
	026c  MCP51 High Definition Audio
	026d  MCP51 USB Controller
	026e  MCP51 USB Controller
	026f  MCP51 PCI Bridge
	0270  MCP51 Host Bridge
	0271  MCP51 PMU
	0272  MCP51 Memory Controller 0
	027e  C51 Memory Controller 2
	027f  C51 Memory Controller 3
	0280  NV28 [GeForce4 Ti 4800]
	0281  NV28 [GeForce4 Ti 4200 AGP 8x]
	0282  NV28 [GeForce4 Ti 4800 SE]
	0286  NV28M [GeForce4 Ti 4200 Go AGP 8x]
	0288  NV28GL [Quadro4 980 XGL]
	0289  NV28GL [Quadro4 780 XGL]
	028c  NV28GLM [Quadro4 Go700]
	0290  G71 [GeForce 7900 GTX]
	0291  G71 [GeForce 7900 GT/GTO]
	0292  G71 [GeForce 7900 GS]
	0293  G71 [GeForce 7900 GX2]
	0294  G71 [GeForce 7950 GX2]
	0295  G71 [GeForce 7950 GT]
	0297  G71M [GeForce Go 7950 GTX]
	0298  G71M [GeForce Go 7900 GS]
	0299  G71M [GeForce Go 7900 GTX]
	029a  G71GLM [Quadro FX 2500M]
	029b  G71GLM [Quadro FX 1500M]
	029c  G71GL [Quadro FX 5500]
	029d  G71GL [Quadro FX 3500]
	029e  G71GL [Quadro FX 1500]
	029f  G71GL [Quadro FX 4500 X2]
	02a0  NV2A [XGPU]
	02a5  MCPX CPU Bridge
	02a6  MCPX Memory Controller
	02e0  G73 [GeForce 7600 GT]
	02e1  G73 [GeForce 7600 GS]
	02e2  G73 [GeForce 7300 GT]
	02e3  G71 [GeForce 7900 GS]
	02e4  G71 [GeForce 7950 GT]
	02f0  C51 Host Bridge
	02f1  C51 Host Bridge
	02f2  C51 Host Bridge
	02f3  C51 Host Bridge
	02f4  C51 Host Bridge
	02f5  C51 Host Bridge
	02f6  C51 Host Bridge
	02f7  C51 Host Bridge
	02f8  C51 Memory Controller 5
	02f9  C51 Memory Controller 4
	02fa  C51 Memory Controller 0
	02fb  C51 PCI Express Bridge
	02fc  C51 PCI Express Bridge
	02fd  C51 PCI Express Bridge
	02fe  C51 Memory Controller 1
	02ff  C51 Host Bridge
	0300  NV30 [GeForce FX]
	0301  NV30 [GeForce FX 5800 Ultra]
	0302  NV30 [GeForce FX 5800]
	0308  NV30GL [Quadro FX 2000]
	0309  NV30GL [Quadro FX 1000]
	0311  NV31 [GeForce FX 5600 Ultra]
	0312  NV31 [GeForce FX 5600]
	0314  NV31 [GeForce FX 5600XT]
	0316  NV31M
	0318  NV31GL
	031a  NV31M [GeForce FX Go5600]
	031b  NV31M [GeForce FX Go5650]
	031c  NV31GLM [Quadro FX Go700]
	0320  NV34 [GeForce FX 5200]
	0321  NV34 [GeForce FX 5200 Ultra]
	0322  NV34 [GeForce FX 5200]
	0323  NV34 [GeForce FX 5200LE]
	0324  NV34M [GeForce FX Go5200 64M]
	0325  NV34M [GeForce FX Go5250]
	0326  NV34 [GeForce FX 5500]
	0327  NV34 [GeForce FX 5100]
	0328  NV34M [GeForce FX Go5200 32M/64M]
	0329  NV34M [GeForce FX Go5200]
	032a  NV34GL [Quadro NVS 280 PCI]
	032b  NV34GL [Quadro FX 500/600 PCI]
	032c  NV34M [GeForce FX Go5300 / Go5350]
	032d  NV34M [GeForce FX Go5100]
	032e  NV34
	032f  NV34 [GeForce FX 5200]
	0330  NV35 [GeForce FX 5900 Ultra]
	0331  NV35 [GeForce FX 5900]
	0332  NV35 [GeForce FX 5900XT]
	0333  NV38 [GeForce FX 5950 Ultra]
	0334  NV35 [GeForce FX 5900ZT]
	0338  NV35GL [Quadro FX 3000]
	033f  NV35GL [Quadro FX 700]
	0341  NV36 [GeForce FX 5700 Ultra]
	0342  NV36 [GeForce FX 5700]
	0343  NV36 [GeForce FX 5700LE]
	0344  NV36 [GeForce FX 5700VE]
	0347  NV36M [GeForce FX Go5700]
	0348  NV36M [GeForce FX Go5700]
	034c  NV36 [Quadro FX Go1000]
	034d  NV36
	034e  NV36GL [Quadro FX 1100]
	0360  MCP55 LPC Bridge
	0361  MCP55 LPC Bridge
	0362  MCP55 LPC Bridge
	0363  MCP55 LPC Bridge
	0364  MCP55 LPC Bridge
	0365  MCP55 LPC Bridge
	0366  MCP55 LPC Bridge
	0367  MCP55 LPC Bridge
	0368  MCP55 SMBus Controller
	0369  MCP55 Memory Controller
	036a  MCP55 Memory Controller
	036b  MCP55 SMU
	036c  MCP55 USB Controller
	036d  MCP55 USB Controller
	036e  MCP55 IDE
	0370  MCP55 PCI bridge
	0371  MCP55 High Definition Audio
	0372  MCP55 Ethernet
	0373  MCP55 Ethernet
	0374  MCP55 PCI Express bridge
	0375  MCP55 PCI Express bridge
	0376  MCP55 PCI Express bridge
	0377  MCP55 PCI Express bridge
	0378  MCP55 PCI Express bridge
	037a  MCP55 Memory Controller
	037e  MCP55 SATA Controller
	037f  MCP55 SATA Controller
	038b  G73 [GeForce 7650 GS]
	0390  G73 [GeForce 7650 GS]
	0391  G73 [GeForce 7600 GT]
	0392  G73 [GeForce 7600 GS]
	0393  G73 [GeForce 7300 GT]
	0394  G73 [GeForce 7600 LE]
	0395  G73 [GeForce 7300 GT]
	0396  G73
	0397  G73M [GeForce Go 7700]
	0398  G73M [GeForce Go 7600]
	0399  G73M [GeForce Go 7600 GT]
	039a  G73M [Quadro NVS 300M]
	039b  G73M [GeForce Go 7900 SE]
	039c  G73GLM [Quadro FX 550M]
	039d  G73
	039e  G73GL [Quadro FX 560]
	039f  G73
	03a0  C55 Host Bridge
	03a1  C55 Host Bridge
	03a2  C55 Host Bridge
	03a3  C55 Host Bridge
	03a4  C55 Host Bridge
	03a5  C55 Host Bridge
	03a6  C55 Host Bridge
	03a7  C55 Host Bridge
	03a8  C55 Memory Controller
	03a9  C55 Memory Controller
	03aa  C55 Memory Controller
	03ab  C55 Memory Controller
	03ac  C55 Memory Controller
	03ad  C55 Memory Controller
	03ae  C55 Memory Controller
	03af  C55 Memory Controller
	03b0  C55 Memory Controller
	03b1  C55 Memory Controller
	03b2  C55 Memory Controller
	03b3  C55 Memory Controller
	03b4  C55 Memory Controller
	03b5  C55 Memory Controller
	03b6  C55 Memory Controller
	03b7  C55 PCI Express bridge
	03b8  C55 PCI Express bridge
	03b9  C55 PCI Express bridge
	03ba  C55 Memory Controller
	03bb  C55 PCI Express bridge
	03bc  C55 Memory Controller
	03d0  C61 [GeForce 6150SE nForce 430]
	03d1  C61 [GeForce 6100 nForce 405]
	03d2  C61 [GeForce 6100 nForce 400]
	03d5  C61 [GeForce 6100 nForce 420]
	03d6  C61 [GeForce 7025 / nForce 630a]
	03e0  MCP61 LPC Bridge
	03e1  MCP61 LPC Bridge
	03e2  MCP61 Host Bridge
	03e3  MCP61 LPC Bridge
	03e4  MCP61 High Definition Audio
	03e5  MCP61 Ethernet
	03e6  MCP61 Ethernet
	03e7  MCP61 SATA Controller
	03e8  MCP61 PCI Express bridge
	03e9  MCP61 PCI Express bridge
	03ea  MCP61 Memory Controller
	03eb  MCP61 SMBus
	03ec  MCP61 IDE
	03ee  MCP61 Ethernet
	03ef  MCP61 Ethernet
	03f0  MCP61 High Definition Audio
	03f1  MCP61 USB 1.1 Controller
	03f2  MCP61 USB 2.0 Controller
	03f3  MCP61 PCI bridge
	03f4  MCP61 SMU
	03f5  MCP61 Memory Controller
	03f6  MCP61 SATA Controller
	03f7  MCP61 SATA Controller
	0400  G84 [GeForce 8600 GTS]
	0401  G84 [GeForce 8600 GT]
	0402  G84 [GeForce 8600 GT]
	0403  G84 [GeForce 8600 GS]
	0404  G84 [GeForce 8400 GS]
	0405  G84M [GeForce 9500M GS]
	0406  G84 [GeForce 8300 GS]
	0407  G84M [GeForce 8600M GT]
	0408  G84M [GeForce 9650M GS]
	0409  G84M [GeForce 8700M GT]
	040a  G84GL [Quadro FX 370]
	040b  G84GLM [Quadro NVS 320M]
	040c  G84GLM [Quadro FX 570M]
	040d  G84GLM [Quadro FX 1600M]
	040e  G84GL [Quadro FX 570]
	040f  G84GL [Quadro FX 1700]
	0410  G92 [GeForce GT 330]
	0414  G92 [GeForce 9800 GT]
	0420  G86 [GeForce 8400 SE]
	0421  G86 [GeForce 8500 GT]
	0422  G86 [GeForce 8400 GS]
	0423  G86 [GeForce 8300 GS]
	0424  G86 [GeForce 8400 GS]
	0425  G86M [GeForce 8600M GS]
	0426  G86M [GeForce 8400M GT]
	0427  G86M [GeForce 8400M GS]
	0428  G86M [GeForce 8400M G]
	0429  G86M [Quadro NVS 140M]
	042a  G86M [Quadro NVS 130M]
	042b  G86M [Quadro NVS 135M]
	042c  G86 [GeForce 9400 GT]
	042d  G86GLM [Quadro FX 360M]
	042e  G86M [GeForce 9300M G]
	042f  G86 [Quadro NVS 290]
	0440  MCP65 LPC Bridge
	0441  MCP65 LPC Bridge
	0442  MCP65 LPC Bridge
	0443  MCP65 LPC Bridge
	0444  MCP65 Memory Controller
	0445  MCP65 Memory Controller
	0446  MCP65 SMBus
	0447  MCP65 SMU
	0448  MCP65 IDE
	0449  MCP65 PCI bridge
	044a  MCP65 High Definition Audio
	044b  MCP65 High Definition Audio
	044c  MCP65 AHCI Controller
	044d  MCP65 AHCI Controller
	044e  MCP65 AHCI Controller
	044f  MCP65 AHCI Controller
	0450  MCP65 Ethernet
	0451  MCP65 Ethernet
	0452  MCP65 Ethernet
	0453  MCP65 Ethernet
	0454  MCP65 USB 1.1 OHCI Controller
	0455  MCP65 USB 2.0 EHCI Controller
	0456  MCP65 USB Controller
	0457  MCP65 USB Controller
	0458  MCP65 PCI Express bridge
	0459  MCP65 PCI Express bridge
	045a  MCP65 PCI Express bridge
	045b  MCP65 PCI Express bridge
	045c  MCP65 SATA Controller
	045d  MCP65 SATA Controller
	045e  MCP65 SATA Controller
	045f  MCP65 SATA Controller
	0531  C67 [GeForce 7150M / nForce 630M]
	0533  C67 [GeForce 7000M / nForce 610M]
	053a  C68 [GeForce 7050 PV / nForce 630a]
	053b  C68 [GeForce 7050 PV / nForce 630a]
	053e  C68 [GeForce 7025 / nForce 630a]
	0541  MCP67 Memory Controller
	0542  MCP67 SMBus
	0543  MCP67 Co-processor
	0547  MCP67 Memory Controller
	0548  MCP67 ISA Bridge
	054c  MCP67 Ethernet
	054d  MCP67 Ethernet
	054e  MCP67 Ethernet
	054f  MCP67 Ethernet
	0550  MCP67 AHCI Controller
	0554  MCP67 AHCI Controller
	0555  MCP67 SATA Controller
	055c  MCP67 High Definition Audio
	055d  MCP67 High Definition Audio
	055e  MCP67 OHCI USB 1.1 Controller
	055f  MCP67 EHCI USB 2.0 Controller
	0560  MCP67 IDE Controller
	0561  MCP67 PCI Bridge
	0562  MCP67 PCI Express Bridge
	0563  MCP67 PCI Express Bridge
	0568  MCP78S [GeForce 8200] Memory Controller
	0569  MCP78S [GeForce 8200] PCI Express Bridge
	056a  MCP73 [nForce 630i] USB 2.0 Controller (EHCI)
	056c  MCP73 IDE Controller
	056d  MCP73 PCI Express bridge
	056e  MCP73 PCI Express bridge
	056f  MCP73 PCI Express bridge
	05b1  NF200 PCIe 2.0 switch
	05b8  NF200 PCIe 2.0 switch for GTX 295
	05be  NF200 PCIe 2.0 switch for Quadro Plex S4 / Tesla S870 / Tesla S1070 / Tesla S2050
	05e0  GT200b [GeForce GTX 295]
	05e1  GT200 [GeForce GTX 280]
	05e2  GT200 [GeForce GTX 260]
	05e3  GT200b [GeForce GTX 285]
	05e6  GT200b [GeForce GTX 275]
	05e7  GT200GL [Tesla C1060 / M1060]
	05ea  GT200 [GeForce GTX 260]
	05eb  GT200 [GeForce GTX 295]
	05ed  GT200GL [Quadro Plex 2200 D2]
	05f1  GT200 [GeForce GTX 280]
	05f2  GT200 [GeForce GTX 260]
	05f8  GT200GL [Quadro Plex 2200 S4]
	05f9  GT200GL [Quadro CX]
	05fd  GT200GL [Quadro FX 5800]
	05fe  GT200GL [Quadro FX 4800]
	05ff  GT200GL [Quadro FX 3800]
	0600  G92 [GeForce 8800 GTS 512]
	0601  G92 [GeForce 9800 GT]
	0602  G92 [GeForce 8800 GT]
	0603  G92 [GeForce GT 230 OEM]
	0604  G92 [GeForce 9800 GX2]
	0605  G92 [GeForce 9800 GT]
	0606  G92 [GeForce 8800 GS]
	0607  G92 [GeForce GTS 240]
	0608  G92M [GeForce 9800M GTX]
	0609  G92M [GeForce 8800M GTS]
	060a  G92M [GeForce GTX 280M]
	060b  G92M [GeForce 9800M GT]
	060c  G92M [GeForce 8800M GTX]
	060d  G92 [GeForce 8800 GS]
	060f  G92M [GeForce GTX 285M]
	0610  G92 [GeForce 9600 GSO]
	0611  G92 [GeForce 8800 GT]
	0612  G92 [GeForce 9800 GTX / 9800 GTX+]
	0613  G92 [GeForce 9800 GTX+]
	0614  G92 [GeForce 9800 GT]
	0615  G92 [GeForce GTS 250]
	0617  G92M [GeForce 9800M GTX]
	0618  G92M [GeForce GTX 260M]
	0619  G92GL [Quadro FX 4700 X2]
	061a  G92GL [Quadro FX 3700]
	061b  G92GL [Quadro VX 200]
	061c  G92GLM [Quadro FX 3600M]
	061d  G92GLM [Quadro FX 2800M]
	061e  G92GLM [Quadro FX 3700M]
	061f  G92GLM [Quadro FX 3800M]
	0620  G94 [GeForce 9800 GT]
	0621  G94 [GeForce GT 230]
	0622  G94 [GeForce 9600 GT]
	0623  G94 [GeForce 9600 GS]
	0624  G94 [GeForce 9600 GT Green Edition]
	0625  G94 [GeForce 9600 GSO 512]
	0626  G94 [GeForce GT 130]
	0627  G94 [GeForce GT 140]
	0628  G94M [GeForce 9800M GTS]
	062a  G94M [GeForce 9700M GTS]
	062b  G94M [GeForce 9800M GS]
	062c  G94M [GeForce 9800M GTS]
	062d  G94 [GeForce 9600 GT]
	062e  G94 [GeForce 9600 GT]
	062f  G94 [GeForce 9800 S]
	0630  G94 [GeForce 9600 GT]
	0631  G94M [GeForce GTS 160M]
	0632  G94M [GeForce GTS 150M]
	0633  G94 [GeForce GT 220]
	0635  G94 [GeForce 9600 GSO]
	0637  G94 [GeForce 9600 GT]
	0638  G94GL [Quadro FX 1800]
	063a  G94GLM [Quadro FX 2700M]
	063f  G94 [GeForce 9600 GE]
	0640  G96 [GeForce 9500 GT]
	0641  G96 [GeForce 9400 GT]
	0642  G96 [D9M-10]
	0643  G96 [GeForce 9500 GT]
	0644  G96 [GeForce 9500 GS]
	0645  G96 [GeForce 9500 GS]
	0646  G96 [GeForce GT 120]
	0647  G96M [GeForce 9600M GT]
	0648  G96M [GeForce 9600M GS]
	0649  G96M [GeForce 9600M GT]
	064a  G96M [GeForce 9700M GT]
	064b  G96M [GeForce 9500M G]
	064c  G96M [GeForce 9650M GT]
	064d  G96 [GeForce 9600 GT]
	064e  G96 [GeForce 9600 GT / 9800 GT]
	0651  G96M [GeForce G 110M]
	0652  G96M [GeForce GT 130M]
	0653  G96M [GeForce GT 120M]
	0654  G96M [GeForce GT 220M]
	0655  G96 [GeForce GT 120]
	0656  G96 [GeForce 9650 S]
	0658  G96GL [Quadro FX 380]
	0659  G96GL [Quadro FX 580]
	065a  G96GLM [Quadro FX 1700M]
	065b  G96 [GeForce 9400 GT]
	065c  G96GLM [Quadro FX 770M]
	065d  G96 [GeForce 9500 GA / 9600 GT / GTS 250]
	065f  G96 [GeForce G210]
	06c0  GF100 [GeForce GTX 480]
	06c4  GF100 [GeForce GTX 465]
	06ca  GF100M [GeForce GTX 480M]
	06cb  GF100 [GeForce GTX 480]
	06cd  GF100 [GeForce GTX 470]
	06d1  GF100GL [Tesla C2050 / C2070]
	06d2  GF100GL [Tesla M2070]
	06d8  GF100GL [Quadro 6000]
	06d9  GF100GL [Quadro 5000]
	06da  GF100GLM [Quadro 5000M]
	06dc  GF100GL [Quadro 6000]
	06dd  GF100GL [Quadro 4000]
	06de  GF100GL [Tesla T20 Processor]
	06df  GF100GL [Tesla M2070-Q]
	06e0  G98 [GeForce 9300 GE]
	06e1  G98 [GeForce 9300 GS]
	06e2  G98 [GeForce 8400]
	06e3  G98 [GeForce 8300 GS]
	06e4  G98 [GeForce 8400 GS Rev. 2]
	06e5  G98M [GeForce 9300M GS]
	06e6  G98 [GeForce G 100]
	06e7  G98 [GeForce 9300 SE]
	06e8  G98M [GeForce 9200M GS]
	06e9  G98M [GeForce 9300M GS]
	06ea  G98M [Quadro NVS 150M]
	06eb  G98M [Quadro NVS 160M]
	06ec  G98M [GeForce G 105M]
	06ed  G98 [GeForce 9600 GT / 9800 GT]
	06ee  G98 [GeForce 9600 GT / 9800 GT]
	06ef  G98M [GeForce G 103M]
	06f1  G98M [GeForce G 105M]
	06f8  G98 [Quadro NVS 420]
	06f9  G98GL [Quadro FX 370 LP]
	06fa  G98 [Quadro NVS 450]
	06fb  G98GLM [Quadro FX 370M]
	06fd  G98 [Quadro NVS 295]
	06ff  G98 [HICx16 + Graphics]
	0751  MCP78S [GeForce 8200] Memory Controller
	0752  MCP78S [GeForce 8200] SMBus
	0753  MCP78S [GeForce 8200] Co-Processor
	0754  MCP78S [GeForce 8200] Memory Controller
	0759  MCP78S [GeForce 8200] IDE
	075a  MCP78S [GeForce 8200] PCI Bridge
	075b  MCP78S [GeForce 8200] PCI Express Bridge
	075c  MCP78S [GeForce 8200] LPC Bridge
	075d  MCP78S [GeForce 8200] LPC Bridge
	0760  MCP77 Ethernet
	0761  MCP77 Ethernet
	0762  MCP77 Ethernet
	0763  MCP77 Ethernet
	0774  MCP72XE/MCP72P/MCP78U/MCP78S High Definition Audio
	0778  MCP78S [GeForce 8200] PCI Express Bridge
	077a  MCP78S [GeForce 8200] PCI Bridge
	077b  MCP78S [GeForce 8200] OHCI USB 1.1 Controller
	077c  MCP78S [GeForce 8200] EHCI USB 2.0 Controller
	077d  MCP78S [GeForce 8200] OHCI USB 1.1 Controller
	077e  MCP78S [GeForce 8200] EHCI USB 2.0 Controller
	07c0  MCP73 Host Bridge
	07c1  MCP73 Host Bridge
	07c2  MCP73 Host Bridge
	07c3  MCP73 Host Bridge
	07c5  MCP73 Host Bridge
	07c8  MCP73 Memory Controller
	07cb  nForce 610i/630i memory controller
	07cd  nForce 610i/630i memory controller
	07ce  nForce 610i/630i memory controller
	07cf  nForce 610i/630i memory controller
	07d0  nForce 610i/630i memory controller
	07d1  nForce 610i/630i memory controller
	07d2  nForce 610i/630i memory controller
	07d3  nForce 610i/630i memory controller
	07d6  nForce 610i/630i memory controller
	07d7  MCP73 LPC Bridge
	07d8  MCP73 SMBus
	07d9  MCP73 Memory Controller
	07da  MCP73 Co-processor
	07dc  MCP73 Ethernet
	07dd  MCP73 Ethernet
	07de  MCP73 Ethernet
	07df  MCP73 Ethernet
	07e0  C73 [GeForce 7150 / nForce 630i]
	07e1  C73 [GeForce 7100 / nForce 630i]
	07e2  C73 [GeForce 7050 / nForce 630i]
	07e3  C73 [GeForce 7050 / nForce 610i]
	07e5  C73 [GeForce 7100 / nForce 620i]
	07f0  MCP73 SATA Controller (IDE mode)
	07f4  GeForce 7100/nForce 630i SATA
	07f8  MCP73 SATA RAID Controller
	07fc  MCP73 High Definition Audio
	07fe  MCP73 OHCI USB 1.1 Controller
	0840  C77 [GeForce 8200M]
	0844  C77 [GeForce 9100M G]
	0845  C77 [GeForce 8200M G]
	0846  C77 [GeForce 9200]
	0847  C78 [GeForce 9100]
	0848  C77 [GeForce 8300]
	0849  C77 [GeForce 8200]
	084a  C77 [nForce 730a]
	084b  C77 [GeForce 8200]
	084c  C77 [nForce 780a/980a SLI]
	084d  C77 [nForce 750a SLI]
	084f  C77 [GeForce 8100 / nForce 720a]
	0860  C79 [GeForce 9300]
	0861  C79 [GeForce 9400]
	0862  C79 [GeForce 9400M G]
	0863  C79 [GeForce 9400M]
	0864  C79 [GeForce 9300]
	0865  C79 [GeForce 9300/ION]
	0866  C79 [GeForce 9400M G]
	0867  C79 [GeForce 9400]
	0868  C79 [nForce 760i SLI]
	0869  MCP7A [GeForce 9400]
	086a  C79 [GeForce 9400]
	086c  C79 [GeForce 9300 / nForce 730i]
	086d  C79 [GeForce 9200]
	086e  C79 [GeForce 9100M G]
	086f  MCP79 [GeForce 8200M G]
	0870  C79 [GeForce 9400M]
	0871  C79 [GeForce 9200]
	0872  C79 [GeForce G102M]
	0873  C79 [GeForce G102M]
	0874  C79 [ION]
	0876  ION VGA [GeForce 9400M]
	087a  C79 [GeForce 9400]
	087d  ION VGA
	087e  ION LE VGA
	087f  ION LE VGA
	08a0  MCP89 [GeForce 320M]
	08a2  MCP89 [GeForce 320M]
	08a3  MCP89 [GeForce 320M]
	08a4  MCP89 [GeForce 320M]
	08a5  MCP89 [GeForce 320M]
	0a20  GT216 [GeForce GT 220]
	0a21  GT216M [GeForce GT 330M]
	0a22  GT216 [GeForce 315]
	0a23  GT216 [GeForce 210]
	0a26  GT216 [GeForce 405]
	0a27  GT216 [GeForce 405]
	0a28  GT216M [GeForce GT 230M]
	0a29  GT216M [GeForce GT 330M]
	0a2a  GT216M [GeForce GT 230M]
	0a2b  GT216M [GeForce GT 330M]
	0a2c  GT216M [NVS 5100M]
	0a2d  GT216M [GeForce GT 320M]
	0a30  GT216 [GeForce 505]
	0a32  GT216 [GeForce GT 415]
	0a34  GT216M [GeForce GT 240M]
	0a35  GT216M [GeForce GT 325M]
	0a38  GT216GL [Quadro 400]
	0a3c  GT216GLM [Quadro FX 880M]
	0a60  GT218 [GeForce G210]
	0a62  GT218 [GeForce 205]
	0a63  GT218 [GeForce 310]
	0a64  GT218 [ION]
	0a65  GT218 [GeForce 210]
	0a66  GT218 [GeForce 310]
	0a67  GT218 [GeForce 315]
	0a68  GT218M [GeForce G 105M]
	0a69  GT218M [GeForce G 105M]
	0a6a  GT218M [NVS 2100M]
	0a6c  GT218M [NVS 3100M]
	0a6e  GT218M [GeForce 305M]
	0a6f  GT218 [ION]
	0a70  GT218M [GeForce 310M]
	0a71  GT218M [GeForce 305M]
	0a72  GT218M [GeForce 310M]
	0a73  GT218M [GeForce 305M]
	0a74  GT218M [GeForce G210M]
	0a75  GT218M [GeForce 310M]
	0a76  GT218 [ION 2]
	0a78  GT218GL [Quadro FX 380 LP]
	0a7a  GT218M [GeForce 315M]
	0a7b  GT218 [GeForce 505]
	0a7c  GT218GLM [Quadro FX 380M]
	0a80  MCP79 Host Bridge
	0a81  MCP79 Host Bridge
	0a82  MCP79 Host Bridge
	0a83  MCP79 Host Bridge
	0a84  MCP79 Host Bridge
	0a85  MCP79 Host Bridge
	0a86  MCP79 Host Bridge
	0a87  MCP79 Host Bridge
	0a88  MCP79 Memory Controller
	0a89  MCP79 Memory Controller
	0a98  MCP79 Memory Controller
	0aa0  MCP79 PCI Express Bridge
	0aa2  MCP79 SMBus
	0aa3  MCP79 Co-processor
	0aa4  MCP79 Memory Controller
	0aa5  MCP79 OHCI USB 1.1 Controller
	0aa6  MCP79 EHCI USB 2.0 Controller
	0aa7  MCP79 OHCI USB 1.1 Controller
	0aa8  MCP79 OHCI USB 1.1 Controller
	0aa9  MCP79 EHCI USB 2.0 Controller
	0aaa  MCP79 EHCI USB 2.0 Controller
	0aab  MCP79 PCI Bridge
	0aac  MCP79 LPC Bridge
	0aad  MCP79 LPC Bridge
	0aae  MCP79 LPC Bridge
	0aaf  MCP79 LPC Bridge
	0ab0  MCP79 Ethernet
	0ab1  MCP79 Ethernet
	0ab2  MCP79 Ethernet
	0ab3  MCP79 Ethernet
	0ab4  MCP79 SATA Controller
	0ab5  MCP79 SATA Controller
	0ab6  MCP79 SATA Controller
	0ab7  MCP79 SATA Controller
	0ab8  MCP79 AHCI Controller
	0ab9  MCP79 AHCI Controller
	0aba  MCP79 AHCI Controller
	0abb  MCP79 AHCI Controller
	0abc  MCP79 RAID Controller
	0abd  MCP79 RAID Controller
	0abe  MCP79 RAID Controller
	0abf  MCP79 RAID Controller
	0ac0  MCP79 High Definition Audio
	0ac1  MCP79 High Definition Audio
	0ac2  MCP79 High Definition Audio
	0ac3  MCP79 High Definition Audio
	0ac4  MCP79 PCI Express Bridge
	0ac5  MCP79 PCI Express Bridge
	0ac6  MCP79 PCI Express Bridge
	0ac7  MCP79 PCI Express Bridge
	0ac8  MCP79 PCI Express Bridge
	0ad0  MCP78S [GeForce 8200] SATA Controller (non-AHCI mode)
	0ad4  MCP78S [GeForce 8200] AHCI Controller
	0ad8  MCP78S [GeForce 8200] SATA Controller (RAID mode)
	0be2  GT216 HDMI Audio Controller
	0be3  High Definition Audio Controller
	0be4  High Definition Audio Controller
	0be5  GF100 High Definition Audio Controller
	0be9  GF106 High Definition Audio Controller
	0bea  GF108 High Definition Audio Controller
	0beb  GF104 High Definition Audio Controller
	0bee  GF116 High Definition Audio Controller
	0bf0  Tegra2 PCIe x4 Bridge
	0bf1  Tegra2 PCIe x2 Bridge
	0ca0  GT215 [GeForce GT 330]
	0ca2  GT215 [GeForce GT 320]
	0ca3  GT215 [GeForce GT 240]
	0ca4  GT215 [GeForce GT 340]
	0ca5  GT215 [GeForce GT 220]
	0ca7  GT215 [GeForce GT 330]
	0ca8  GT215M [GeForce GTS 260M]
	0ca9  GT215M [GeForce GTS 250M]
	0cac  GT215 [GeForce GT 220/315]
	0caf  GT215M [GeForce GT 335M]
	0cb0  GT215M [GeForce GTS 350M]
	0cb1  GT215M [GeForce GTS 360M]
	0cbc  GT215GLM [Quadro FX 1800M]
	0d60  MCP89 HOST Bridge
	0d68  MCP89 Memory Controller
	0d69  MCP89 Memory Controller
	0d76  MCP89 PCI Express Bridge
	0d79  MCP89 SMBus
	0d7a  MCP89 Co-Processor
	0d7b  MCP89 Memory Controller
	0d7d  MCP89 Ethernet
	0d80  MCP89 LPC Bridge
	0d85  MCP89 SATA Controller
	0d88  MCP89 SATA Controller (AHCI mode)
	0d89  MCP89 SATA Controller (AHCI mode)
	0d8d  MCP89 SATA Controller (RAID mode)
	0d94  MCP89 High Definition Audio
	0d9c  MCP89 OHCI USB 1.1 Controller
	0d9d  MCP89 EHCI USB 2.0 Controller
	0dc0  GF106 [GeForce GT 440]
	0dc4  GF106 [GeForce GTS 450]
	0dc5  GF106 [GeForce GTS 450 OEM]
	0dc6  GF106 [GeForce GTS 450 OEM]
	0dcd  GF106M [GeForce GT 555M]
	0dce  GF106M [GeForce GT 555M]
	0dd1  GF106M [GeForce GTX 460M]
	0dd2  GF106M [GeForce GT 445M]
	0dd3  GF106M [GeForce GT 435M]
	0dd6  GF106M [GeForce GT 550M]
	0dd8  GF106GL [Quadro 2000]
	0dda  GF106GLM [Quadro 2000M]
	0de0  GF108 [GeForce GT 440]
	0de1  GF108 [GeForce GT 430]
	0de2  GF108 [GeForce GT 420]
	0de3  GF108M [GeForce GT 635M]
	0de4  GF108 [GeForce GT 520]
	0de5  GF108 [GeForce GT 530]
	0de7  GF108 [GeForce GT 610]
	0de8  GF108M [GeForce GT 620M]
	0de9  GF108M [GeForce GT 620M/630M/635M/640M LE]
	0dea  GF108M [GeForce 610M]
	0deb  GF108M [GeForce GT 555M]
	0dec  GF108M [GeForce GT 525M]
	0ded  GF108M [GeForce GT 520M]
	0dee  GF108M [GeForce GT 415M]
	0def  GF108M [NVS 5400M]
	0df0  GF108M [GeForce GT 425M]
	0df1  GF108M [GeForce GT 420M]
	0df2  GF108M [GeForce GT 435M]
	0df3  GF108M [GeForce GT 420M]
	0df4  GF108M [GeForce GT 540M]
	0df5  GF108M [GeForce GT 525M]
	0df6  GF108M [GeForce GT 550M]
	0df7  GF108M [GeForce GT 520M]
	0df8  GF108GL [Quadro 600]
	0df9  GF108GLM [Quadro 500M]
	0dfa  GF108GLM [Quadro 1000M]
	0dfc  GF108GLM [NVS 5200M]
	0e08  GF119 HDMI Audio Controller
	0e09  GF110 High Definition Audio Controller
	0e0a  GK104 HDMI Audio Controller
	0e0b  GK106 HDMI Audio Controller
	0e0c  GF114 HDMI Audio Controller
	0e0f  GK208 HDMI/DP Audio Controller
	0e12  TegraK1 PCIe x4 Bridge
	0e13  TegraK1 PCIe x1 Bridge
	0e1a  GK110 HDMI Audio
	0e1b  GK107 HDMI Audio Controller
	0e1c  Tegra3+ PCIe x4 Bridge
	0e1d  Tegra3+ PCIe x2 Bridge
	0e22  GF104 [GeForce GTX 460]
	0e23  GF104 [GeForce GTX 460 SE]
	0e24  GF104 [GeForce GTX 460 OEM]
	0e30  GF104M [GeForce GTX 470M]
	0e31  GF104M [GeForce GTX 485M]
	0e3a  GF104GLM [Quadro 3000M]
	0e3b  GF104GLM [Quadro 4000M]
	0f00  GF108 [GeForce GT 630]
	0f01  GF108 [GeForce GT 620]
	0f02  GF108 [GeForce GT 730]
	0f06  GF108 [GeForce GT 730]
	0fb0  GM200 High Definition Audio
	0fb8  GP108 High Definition Audio Controller
	0fb9  GP107GL High Definition Audio Controller
	0fbb  GM204 High Definition Audio Controller
	0fc0  GK107 [GeForce GT 640 OEM]
	0fc1  GK107 [GeForce GT 640]
	0fc2  GK107 [GeForce GT 630 OEM]
	0fc6  GK107 [GeForce GTX 650]
	0fc8  GK107 [GeForce GT 740]
	0fc9  GK107 [GeForce GT 730]
	0fcd  GK107M [GeForce GT 755M]
	0fce  GK107M [GeForce GT 640M LE]
	0fd1  GK107M [GeForce GT 650M]
	0fd2  GK107M [GeForce GT 640M]
	0fd3  GK107M [GeForce GT 640M LE]
	0fd4  GK107M [GeForce GTX 660M]
	0fd5  GK107M [GeForce GT 650M Mac Edition]
	0fd8  GK107M [GeForce GT 640M Mac Edition]
	0fd9  GK107M [GeForce GT 645M]
	0fdb  GK107M
	0fdf  GK107M [GeForce GT 740M]
	0fe0  GK107M [GeForce GTX 660M Mac Edition]
	0fe1  GK107M [GeForce GT 730M]
	0fe2  GK107M [GeForce GT 745M]
	0fe3  GK107M [GeForce GT 745M]
	0fe4  GK107M [GeForce GT 750M]
	0fe5  GK107 [GeForce K340 USM]
	0fe6  GK107 [GRID K1 NVS USM]
	0fe7  GK107GL [GRID K100 vGPU]
	0fe9  GK107M [GeForce GT 750M Mac Edition]
	0fea  GK107M [GeForce GT 755M Mac Edition]
	0fec  GK107M [GeForce 710A]
	0fed  GK107M [GeForce 820M]
	0fee  GK107M [GeForce 810M]
	0fef  GK107GL [GRID K340]
	0ff1  GK107 [NVS 1000]
	0ff2  GK107GL [GRID K1]
	0ff3  GK107GL [Quadro K420]
	0ff5  GK107GL [GRID K1 Tesla USM]
	0ff6  GK107GLM [Quadro K1100M]
	0ff7  GK107GL [GRID K140Q vGPU]
	0ff8  GK107GLM [Quadro K500M]
	0ff9  GK107GL [Quadro K2000D]
	0ffa  GK107GL [Quadro K600]
	0ffb  GK107GLM [Quadro K2000M]
	0ffc  GK107GLM [Quadro K1000M]
	0ffd  GK107 [NVS 510]
	0ffe  GK107GL [Quadro K2000]
	0fff  GK107GL [Quadro 410]
	1001  GK110B [GeForce GTX TITAN Z]
	1003  GK110 [GeForce GTX Titan LE]
	1004  GK110 [GeForce GTX 780]
	1005  GK110 [GeForce GTX TITAN]
	1007  GK110 [GeForce GTX 780 Rev. 2]
	1008  GK110 [GeForce GTX 780 Ti Rev. 2]
	100a  GK110B [GeForce GTX 780 Ti]
	100c  GK110B [GeForce GTX TITAN Black]
	101e  GK110GL [Tesla K20X]
	101f  GK110GL [Tesla K20]
	1020  GK110GL [Tesla K20X]
	1021  GK110GL [Tesla K20Xm]
	1022  GK110GL [Tesla K20c]
	1023  GK110BGL [Tesla K40m]
	1024  GK110BGL [Tesla K40c]
	1026  GK110GL [Tesla K20s]
	1027  GK110BGL [Tesla K40st]
	1028  GK110GL [Tesla K20m]
	1029  GK110BGL [Tesla K40s]
	102a  GK110BGL [Tesla K40t]
	102d  GK210GL [Tesla K80]
	102e  GK110BGL [Tesla K40d]
	103a  GK110GL [Quadro K6000]
	103c  GK110GL [Quadro K5200]
	1040  GF119 [GeForce GT 520]
	1042  GF119 [GeForce 510]
	1048  GF119 [GeForce 605]
	1049  GF119 [GeForce GT 620 OEM]
	104a  GF119 [GeForce GT 610]
	104b  GF119 [GeForce GT 625 OEM]
	104c  GF119 [GeForce GT 705]
	104d  GF119 [GeForce GT 710]
	1050  GF119M [GeForce GT 520M]
	1051  GF119M [GeForce GT 520MX]
	1052  GF119M [GeForce GT 520M]
	1054  GF119M [GeForce 410M]
	1055  GF119M [GeForce 410M]
	1056  GF119M [NVS 4200M]
	1057  GF119M [Quadro NVS 4200M]
	1058  GF119M [GeForce 610M]
	1059  GF119M [GeForce 610M]
	105a  GF119M [GeForce 610M]
	105b  GF119M [GeForce 705M]
	107c  GF119 [NVS 315]
	107d  GF119 [NVS 310]
	1080  GF110 [GeForce GTX 580]
	1081  GF110 [GeForce GTX 570]
	1082  GF110 [GeForce GTX 560 Ti OEM]
	1084  GF110 [GeForce GTX 560 OEM]
	1086  GF110 [GeForce GTX 570 Rev. 2]
	1087  GF110 [GeForce GTX 560 Ti 448 Cores]
	1088  GF110 [GeForce GTX 590]
	1089  GF110 [GeForce GTX 580 Rev. 2]
	108b  GF110 [GeForce GTX 580]
	108e  GF110GL [Tesla C2090]
	1091  GF110GL [Tesla M2090]
	1094  GF110GL [Tesla M2075]
	1096  GF110GL [Tesla C2050 / C2075]
	109a  GF100GLM [Quadro 5010M]
	109b  GF100GL [Quadro 7000]
	10c0  GT218 [GeForce 9300 GS Rev. 2]
	10c3  GT218 [GeForce 8400 GS Rev. 3]
	10c5  GT218 [GeForce 405]
	10d8  GT218 [NVS 300]
	10ef  GP102 HDMI Audio Controller
	10f0  GP104 High Definition Audio Controller
	10f1  GP106 High Definition Audio Controller
	1140  GF117M [GeForce 610M/710M/810M/820M / GT 620M/625M/630M/720M]
	1180  GK104 [GeForce GTX 680]
	1182  GK104 [GeForce GTX 760 Ti]
	1183  GK104 [GeForce GTX 660 Ti]
	1184  GK104 [GeForce GTX 770]
	1185  GK104 [GeForce GTX 660 OEM]
	1187  GK104 [GeForce GTX 760]
	1188  GK104 [GeForce GTX 690]
	1189  GK104 [GeForce GTX 670]
	118a  GK104GL [GRID K520]
	118b  GK104GL [GRID K2 GeForce USM]
	118c  GK104 [GRID K2 NVS USM]
	118d  GK104GL [GRID K200 vGPU]
	118e  GK104 [GeForce GTX 760 OEM]
	118f  GK104GL [Tesla K10]
	1191  GK104 [GeForce GTX 760 Rev. 2]
	1193  GK104 [GeForce GTX 760 Ti OEM]
	1194  GK104GL [Tesla K8]
	1195  GK104 [GeForce GTX 660 Rev. 2]
	1198  GK104M [GeForce GTX 880M]
	1199  GK104M [GeForce GTX 870M]
	119a  GK104M [GeForce GTX 860M]
	119d  GK104M [GeForce GTX 775M Mac Edition]
	119e  GK104M [GeForce GTX 780M Mac Edition]
	119f  GK104M [GeForce GTX 780M]
	11a0  GK104M [GeForce GTX 680M]
	11a1  GK104M [GeForce GTX 670MX]
	11a2  GK104M [GeForce GTX 675MX Mac Edition]
	11a3  GK104M [GeForce GTX 680MX]
	11a7  GK104M [GeForce GTX 675MX]
	11b0  GK104GL [GRID K240Q\K260Q vGPU]
	11b1  GK104GL [GRID K2 Tesla USM]
	11b4  GK104GL [Quadro K4200]
	11b6  GK104GLM [Quadro K3100M]
	11b7  GK104GLM [Quadro K4100M]
	11b8  GK104GLM [Quadro K5100M]
	11ba  GK104GL [Quadro K5000]
	11bb  GK104GL [Quadro 4100]
	11bc  GK104GLM [Quadro K5000M]
	11bd  GK104GLM [Quadro K4000M]
	11be  GK104GLM [Quadro K3000M]
	11bf  GK104GL [GRID K2]
	11c0  GK106 [GeForce GTX 660]
	11c2  GK106 [GeForce GTX 650 Ti Boost]
	11c3  GK106 [GeForce GTX 650 Ti OEM]
	11c4  GK106 [GeForce GTX 645 OEM]
	11c5  GK106 [GeForce GT 740]
	11c6  GK106 [GeForce GTX 650 Ti]
	11c7  GK106 [GeForce GTX 750 Ti]
	11c8  GK106 [GeForce GTX 650 OEM]
	11cb  GK106 [GeForce GT 740]
	11e0  GK106M [GeForce GTX 770M]
	11e1  GK106M [GeForce GTX 765M]
	11e2  GK106M [GeForce GTX 765M]
	11e3  GK106M [GeForce GTX 760M]
	11e7  GK106M
	11fa  GK106GL [Quadro K4000]
	11fc  GK106GLM [Quadro K2100M]
	1200  GF114 [GeForce GTX 560 Ti]
	1201  GF114 [GeForce GTX 560]
	1202  GF114 [GeForce GTX 560 Ti OEM]
	1203  GF114 [GeForce GTX 460 SE v2]
	1205  GF114 [GeForce GTX 460 v2]
	1206  GF114 [GeForce GTX 555]
	1207  GF114 [GeForce GT 645 OEM]
	1208  GF114 [GeForce GTX 560 SE]
	1210  GF114M [GeForce GTX 570M]
	1211  GF114M [GeForce GTX 580M]
	1212  GF114M [GeForce GTX 675M]
	1213  GF114M [GeForce GTX 670M]
	1241  GF116 [GeForce GT 545 OEM]
	1243  GF116 [GeForce GT 545]
	1244  GF116 [GeForce GTX 550 Ti]
	1245  GF116 [GeForce GTS 450 Rev. 2]
	1246  GF116M [GeForce GT 550M]
	1247  GF116M [GeForce GT 555M/635M]
	1248  GF116M [GeForce GT 555M/635M]
	1249  GF116 [GeForce GTS 450 Rev. 3]
	124b  GF116 [GeForce GT 640 OEM]
	124d  GF116M [GeForce GT 555M/635M]
	1251  GF116M [GeForce GT 560M]
	1280  GK208 [GeForce GT 635]
	1281  GK208 [GeForce GT 710]
	1282  GK208 [GeForce GT 640 Rev. 2]
	1284  GK208 [GeForce GT 630 Rev. 2]
	1286  GK208 [GeForce GT 720]
	1287  GK208B [GeForce GT 730]
	1288  GK208B [GeForce GT 720]
	1289  GK208 [GeForce GT 710]
	128b  GK208B [GeForce GT 710]
	1290  GK208M [GeForce GT 730M]
	1291  GK208M [GeForce GT 735M]
	1292  GK208M [GeForce GT 740M]
	1293  GK208M [GeForce GT 730M]
	1294  GK208M [GeForce GT 740M]
	1295  GK208M [GeForce 710M]
	1296  GK208M [GeForce 825M]
	1298  GK208M [GeForce GT 720M]
	1299  GK208BM [GeForce 920M]
	129a  GK208BM [GeForce 910M]
	12a0  GK208
	12b9  GK208GLM [Quadro K610M]
	12ba  GK208GLM [Quadro K510M]
	1340  GM108M [GeForce 830M]
	1341  GM108M [GeForce 840M]
	1344  GM108M [GeForce 845M]
	1346  GM108M [GeForce 930M]
	1347  GM108M [GeForce 940M]
	1348  GM108M [GeForce 945M / 945A]
	1349  GM108M [GeForce 930M]
	134b  GM108M [GeForce 940MX]
	134d  GM108M [GeForce 940MX]
	134e  GM108M [GeForce 930MX]
	134f  GM108M [GeForce 920MX]
	137a  GM108GLM [Quadro K620M / Quadro M500M]
	137b  GM108GLM [Quadro M520 Mobile]
	137d  GM108M [GeForce 940A]
	1380  GM107 [GeForce GTX 750 Ti]
	1381  GM107 [GeForce GTX 750]
	1382  GM107 [GeForce GTX 745]
	1389  GM107GL [GRID M30]
	1390  GM107M [GeForce 845M]
	1391  GM107M [GeForce GTX 850M]
	1392  GM107M [GeForce GTX 860M]
	1393  GM107M [GeForce 840M]
	1398  GM107M [GeForce 845M]
	139a  GM107M [GeForce GTX 950M]
	139b  GM107M [GeForce GTX 960M]
	139c  GM107M [GeForce 940M]
	139d  GM107M [GeForce GTX 750 Ti]
	13b0  GM107GLM [Quadro M2000M]
	13b1  GM107GLM [Quadro M1000M]
	13b2  GM107GLM [Quadro M600M]
	13b3  GM107GLM [Quadro K2200M]
	13b4  GM107GLM [Quadro M620 Mobile]
	13b6  GM107GLM [Quadro M1200 Mobile]
	13b9  GM107GL [NVS 810]
	13ba  GM107GL [Quadro K2200]
	13bb  GM107GL [Quadro K620]
	13bc  GM107GL [Quadro K1200]
	13bd  GM107GL [Tesla M10]
	13c0  GM204 [GeForce GTX 980]
	13c1  GM204
	13c2  GM204 [GeForce GTX 970]
	13c3  GM204
	13d7  GM204M [GeForce GTX 980M]
	13d8  GM204M [GeForce GTX 970M]
	13d9  GM204M [GeForce GTX 965M]
	13da  GM204M [GeForce GTX 980 Mobile]
	13e7  GM204 [GeForce GTX 980 Engineering Sample]
	13f0  GM204GL [Quadro M5000]
	13f1  GM204GL [Quadro M4000]
	13f2  GM204GL [Tesla M60]
	13f3  GM204GL [Tesla M6]
	13f8  GM204GLM [Quadro M5000M / M5000 SE]
	13f9  GM204GLM [Quadro M4000M]
	13fa  GM204GLM [Quadro M3000M]
	13fb  GM204GLM [Quadro M5500]
	1401  GM206 [GeForce GTX 960]
	1402  GM206 [GeForce GTX 950]
	1406  GM206 [GeForce GTX 960 OEM]
	1407  GM206 [GeForce GTX 750 v2]
	1427  GM206M [GeForce GTX 965M]
	1430  GM206GL [Quadro M2000]
	1431  GM206GL [Tesla M4]
	1436  GM206GLM [Quadro M2200 Mobile]
	15f0  GP100GL [Quadro GP100]
	15f1  GP100GL
	15f7  GP100GL [Tesla P100 PCIe 12GB]
	15f8  GP100GL [Tesla P100 PCIe 16GB]
	15f9  GP100GL [Tesla P100 SXM2 16GB]
	1617  GM204M [GeForce GTX 980M]
	1618  GM204M [GeForce GTX 970M]
	1619  GM204M [GeForce GTX 965M]
	161a  GM204M [GeForce GTX 980 Mobile]
	1667  GM204M [GeForce GTX 965M]
	1725  GP100
	172e  GP100
	172f  GP100
	17c2  GM200 [GeForce GTX TITAN X]
	17c8  GM200 [GeForce GTX 980 Ti]
	17f0  GM200GL [Quadro M6000]
	17f1  GM200GL [Quadro M6000 24GB]
	17fd  GM200GL [Tesla M40]
	1b00  GP102 [TITAN X]
	1b01  GP102
	1b02  GP102 [TITAN Xp]
	1b06  GP102 [GeForce GTX 1080 Ti]
	1b30  GP102GL [Quadro P6000]
	1b38  GP102GL [Tesla P40]
	1b70  GP102GL
	1b78  GP102GL
	1b80  GP104 [GeForce GTX 1080]
	1b81  GP104 [GeForce GTX 1070]
	1b82  GP104
	1b83  GP104
	1b84  GP104 [GeForce GTX 1060 3GB]
	1b87  GP104 [P104-100]
	1ba0  GP104M [GeForce GTX 1080 Mobile]
	1ba1  GP104M [GeForce GTX 1070 Mobile]
	1bad  GP104 [GeForce GTX 1070 Engineering Sample]
	1bb0  GP104GL [Quadro P5000]
	1bb1  GP104GL [Quadro P4000]
	1bb3  GP104GL [Tesla P4]
	1bb4  GP104GL
	1bb5  GP104GLM [Quadro P5200 Mobile]
	1bb6  GP104GLM [Quadro P5000 Mobile]
	1bb7  GP104GLM [Quadro P4000 Mobile]
	1bb8  GP104GLM [Quadro P3000 Mobile]
	1be0  GP104M [GeForce GTX 1080 Mobile]
	1be1  GP104M [GeForce GTX 1070 Mobile]
	1c00  GP106
	1c01  GP106
	1c02  GP106 [GeForce GTX 1060 3GB]
	1c03  GP106 [GeForce GTX 1060 6GB]
	1c07  GP106 [P106-100]
	1c09  GP106 [P106-090]
	1c20  GP106M [GeForce GTX 1060 Mobile]
	1c21  GP106M [GeForce GTX 1050 Ti Mobile]
	1c22  GP106M [GeForce GTX 1050 Mobile]
	1c30  GP106GL [Quadro P2000]
	1c35  GP106
	1c60  GP106M [GeForce GTX 1060 Mobile 6GB]
	1c61  GP106M [GeForce GTX 1050 Ti Mobile]
	1c62  GP106M [GeForce GTX 1050 Mobile]
	1c70  GP106GL
	1c80  GP107
	1c81  GP107 [GeForce GTX 1050]
	1c82  GP107 [GeForce GTX 1050 Ti]
	1c8c  GP107M [GeForce GTX 1050 Ti Mobile]
	1c8d  GP107M [GeForce GTX 1050 Mobile]
	1c8e  GP107M
	1ca7  GP107GL
	1ca8  GP107GL
	1caa  GP107GL
	1cb1  GP107GL [Quadro P1000]
	1cb2  GP107GL [Quadro P600]
	1cb3  GP107GL [Quadro P400]
	1d01  GP108 [GeForce GT 1030]
	1d10  GP108M [GeForce MX150]
	1d81  GV100
	1db1  GV100 [Tesla V100 SXM2]
	1db4  GV100 [Tesla V100 PCIe]
10df  Emulex Corporation
	0720  OneConnect NIC (Skyhawk)
	0722  OneConnect iSCSI Initiator (Skyhawk)
	0723  OneConnect iSCSI Initiator + Target (Skyhawk)
	0724  OneConnect FCoE Initiator (Skyhawk)
	0728  OneConnect NIC (Skyhawk-VF)
	072a  OneConnect iSCSI Initiator (Skyhawk-VF)
	072b  OneConnect iSCSI Initiator + Target (Skyhawk-VF)
	072c  OneConnect FCoE Initiator (Skyhawk-VF)
	1ae5  LP6000 Fibre Channel Host Adapter
	e100  Proteus-X: LightPulse IOV Fibre Channel Host Adapter
	e131  LightPulse 8Gb/s PCIe Shared I/O Fibre Channel Adapter
	e180  Proteus-X: LightPulse IOV Fibre Channel Host Adapter
	e200  LightPulse LPe16002
	e208  LightPulse 16Gb Fibre Channel Host Adapter (Lancer-VF)
	e220  OneConnect NIC (Lancer)
	e240  OneConnect iSCSI Initiator (Lancer)
	e260  OneConnect FCoE Initiator (Lancer)
	e268  OneConnect 10Gb FCoE Converged Network Adapter (Lancer-VF)
	e300  Lancer Gen6: LPe32000 Fibre Channel Host Adapter
	f011  Saturn: LightPulse Fibre Channel Host Adapter
	f015  Saturn: LightPulse Fibre Channel Host Adapter
	f085  LP850 Fibre Channel Host Adapter
	f095  LP952 Fibre Channel Host Adapter
	f098  LP982 Fibre Channel Host Adapter
	f0a1  Thor LightPulse Fibre Channel Host Adapter
	f0a5  Thor LightPulse Fibre Channel Host Adapter
	f0b5  Viper LightPulse Fibre Channel Host Adapter
	f0d1  Helios LightPulse Fibre Channel Host Adapter
	f0d5  Helios LightPulse Fibre Channel Host Adapter
	f0e1  Zephyr LightPulse Fibre Channel Host Adapter
	f0e5  Zephyr LightPulse Fibre Channel Host Adapter
	f0f5  Neptune LightPulse Fibre Channel Host Adapter
	f100  Saturn-X: LightPulse Fibre Channel Host Adapter
	f111  Saturn-X LightPulse Fibre Channel Host Adapter
	f112  Saturn-X LightPulse Fibre Channel Host Adapter
	f180  LPSe12002 EmulexSecure Fibre Channel Adapter
	f400  LPe36000 Fibre Channel Host Adapter [Prism]
	f700  LP7000 Fibre Channel Host Adapter
	f701  LP7000 Fibre Channel Host Adapter Alternate ID (JX1:2-3, JX2:1-2)
	f800  LP8000 Fibre Channel Host Adapter
	f801  LP8000 Fibre Channel Host Adapter Alternate ID (JX1:2-3, JX2:1-2)
	f900  LP9000 Fibre Channel Host Adapter
	f901  LP9000 Fibre Channel Host Adapter Alternate ID (JX1:2-3, JX2:1-2)
	f980  LP9802 Fibre Channel Host Adapter
	f981  LP9802 Fibre Channel Host Adapter Alternate ID
	f982  LP9802 Fibre Channel Host Adapter Alternate ID
	fa00  Thor-X LightPulse Fibre Channel Host Adapter
	fb00  Viper LightPulse Fibre Channel Host Adapter
	fc00  Thor-X LightPulse Fibre Channel Host Adapter
	fc10  Helios-X LightPulse Fibre Channel Host Adapter
	fc20  Zephyr-X LightPulse Fibre Channel Host Adapter
	fc40  Saturn-X: LightPulse Fibre Channel Host Adapter
	fc50  Proteus-X: LightPulse IOV Fibre Channel Host Adapter
	fd00  Helios-X LightPulse Fibre Channel Host Adapter
	fd11  Helios-X LightPulse Fibre Channel Host Adapter
	fd12  Helios-X LightPulse Fibre Channel Host Adapter
	fe00  Zephyr-X LightPulse Fibre Channel Host Adapter
	fe05  Zephyr-X: LightPulse FCoE Adapter
	fe11  Zephyr-X LightPulse Fibre Channel Host Adapter
	fe12  Zephyr-X LightPulse FCoE Adapter
	ff00  Neptune LightPulse Fibre Channel Host Adapter
10e0  Integrated Micro Solutions Inc.
	5026  IMS5026/27/28
	5027  IMS5027
	5028  IMS5028
	8849  IMS8849
	8853  IMS8853
	9128  IMS9128 [Twin turbo 128]
10e1  Tekram Technology Co.,Ltd.
	0391  TRM-S1040
	690c  DC-690c
	dc29  DC-290
10e2  Aptix Corporation
10e3  Tundra Semiconductor Corp.
	0000  CA91C042 [Universe]
	0108  Tsi108 Host Bridge for Single PowerPC
	0148  Tsi148 [Tempe]
	0860  CA91C860 [QSpan]
	0862  CA91C862A [QSpan-II]
	8260  CA91L8200B [Dual PCI PowerSpan II]
	8261  CA91L8260B [Single PCI PowerSpan II]
	a108  Tsi109 Host Bridge for Dual PowerPC
10e4  Tandem Computers
	8029  Realtek 8029 Network Card
10e5  Micro Industries Corporation
10e6  Gainbery Computer Products Inc.
10e7  Vadem
10e8  Applied Micro Circuits Corp.
	1072  INES GPIB-PCI (AMCC5920 based)
	2011  Q-Motion Video Capture/Edit board
	4750  S5930 [Matchmaker]
	5920  S5920
	8043  LANai4.x [Myrinet LANai interface chip]
	8062  S5933_PARASTATION
	807d  S5933 [Matchmaker]
	8088  Kongsberg Spacetec Format Synchronizer
	8089  Kongsberg Spacetec Serial Output Board
	809c  S5933_HEPC3
	80b9  Harmonix Hi-Card P8 (4x active ISDN BRI)
	80d7  PCI-9112
	80d8  PCI-7200
	80d9  PCI-9118
	80da  PCI-9812
	80fc  APCI1500 Signal processing controller (16 dig. inputs + 16 dig. outputs)
	811a  PCI-IEEE1355-DS-DE Interface
	814c  Fastcom ESCC-PCI (Commtech, Inc.)
	8170  S5933 [Matchmaker] (Chipset Development Tool)
	81e6  Multimedia video controller
	828d  APCI3001 Signal processing controller (up to 16 analog inputs)
	8291  Fastcom 232/8-PCI (Commtech, Inc.)
	82c4  Fastcom 422/4-PCI (Commtech, Inc.)
	82c5  Fastcom 422/2-PCI (Commtech, Inc.)
	82c6  Fastcom IG422/1-PCI (Commtech, Inc.)
	82c7  Fastcom IG232/2-PCI (Commtech, Inc.)
	82ca  Fastcom 232/4-PCI (Commtech, Inc.)
	82db  AJA HDNTV HD SDI Framestore
	82e2  Fastcom DIO24H-PCI (Commtech, Inc.)
	8406  PCIcanx/PCIcan CAN interface [Kvaser AB]
	8407  PCIcan II CAN interface (A1021, PCB-07, PCB-08) [Kvaser AB]
	8851  S5933 on Innes Corp FM Radio Capture card
	e004  X-Gene PCIe bridge
10e9  Alps Electric Co., Ltd.
10ea  Integraphics
	1680  IGA-1680
	1682  IGA-1682
	1683  IGA-1683
	2000  CyberPro 2000
	2010  CyberPro 2000A
	5000  CyberPro 5000
	5050  CyberPro 5050
	5202  CyberPro 5202
	5252  CyberPro5252
10eb  Artists Graphics
	0101  3GA
	8111  Twist3 Frame Grabber
10ec  Realtek Semiconductor Co., Ltd.
	0139  RTL-8139/8139C/8139C+ Ethernet Controller
	5208  RTS5208 PCI Express Card Reader
	5209  RTS5209 PCI Express Card Reader
	5227  RTS5227 PCI Express Card Reader
	5229  RTS5229 PCI Express Card Reader
	522a  RTS522A PCI Express Card Reader
	5249  RTS5249 PCI Express Card Reader
	524a  RTS524A PCI Express Card Reader
	5250  RTS5250 PCI Express Card Reader
	525a  RTS525A PCI Express Card Reader
	5286  RTS5286 PCI Express Card Reader
	5287  RTL8411B PCI Express Card Reader
	5288  RTS5288 PCI Express Card Reader
	5289  RTL8411 PCI Express Card Reader
	8029  RTL-8029(AS)
	8129  RTL-8129
	8136  RTL8101/2/6E PCI Express Fast Ethernet controller
	8138  RT8139 (B/C) Cardbus Fast Ethernet Adapter
	8139  RTL-8100/8101L/8139 PCI Fast Ethernet Adapter
	8167  RTL-8110SC/8169SC Gigabit Ethernet
	8168  RTL8111/8168/8411 PCI Express Gigabit Ethernet Controller
	8169  RTL8169 PCI Gigabit Ethernet Controller
	8171  RTL8191SEvA Wireless LAN Controller
	8172  RTL8191SEvB Wireless LAN Controller
	8173  RTL8192SE Wireless LAN Controller
	8174  RTL8192SE Wireless LAN Controller
	8176  RTL8188CE 802.11b/g/n WiFi Adapter
	8177  RTL8191CE PCIe Wireless Network Adapter
	8178  RTL8192CE PCIe Wireless Network Adapter
	8179  RTL8188EE Wireless Network Adapter
	8180  RTL8180L 802.11b MAC
	8185  RTL-8185 IEEE 802.11a/b/g Wireless LAN Controller
	818b  RTL8192EE PCIe Wireless Network Adapter
	8190  RTL8190 802.11n PCI Wireless Network Adapter
	8191  RTL8192CE PCIe Wireless Network Adapter
	8192  RTL8192E/RTL8192SE Wireless LAN Controller
	8193  RTL8192DE Wireless LAN Controller
	8196  RTL8196 Integrated PCI-e Bridge
	8197  SmartLAN56 56K Modem
	8199  RTL8187SE Wireless LAN Controller
	8723  RTL8723AE PCIe Wireless Network Adapter
	8812  RTL8812AE 802.11ac PCIe Wireless Network Adapter
	8813  RTL8813AE 802.11ac PCIe Wireless Network Adapter
	8821  RTL8821AE 802.11ac PCIe Wireless Network Adapter
	b723  RTL8723BE PCIe Wireless Network Adapter
10ed  Ascii Corporation
	7310  V7310
10ee  Xilinx Corporation
	0001  EUROCOM for PCI (ECOMP)
	0002  Octal E1/T1 for PCI ETP Card
	0007  Default PCIe endpoint ID
	0205  Wildcard TE205P
	0210  Wildcard TE210P
	0300  Spartan 3 Designs (Xilinx IP)
	0314  Wildcard TE405P/TE410P (1st Gen)
	0405  Wildcard TE405P (2nd Gen)
	0410  Wildcard TE410P (2nd Gen)
	0600  Xilinx 6 Designs (Xilinx IP)
	3fc0  RME Digi96
	3fc1  RME Digi96/8
	3fc2  RME Digi96/8 Pro
	3fc3  RME Digi96/8 Pad
	3fc4  RME Digi9652 (Hammerfall)
	3fc5  RME Hammerfall DSP
	3fc6  RME Hammerfall DSP MADI
	7038  FPGA Card XC7VX690T
	8380  Ellips ProfiXpress Profibus Master
	8381  Ellips Santos Frame Grabber
	d154  Copley Controls CAN card (PCI-CAN-02)
	ebf0  SED Systems Modulator/Demodulator
	ebf1  SED Systems Audio Interface Card
	ebf2  SED Systems Common PCI Interface
10ef  Racore Computer Products, Inc.
	8154  M815x Token Ring Adapter
10f0  Peritek Corporation
10f1  Tyan Computer
	2865  Tyan Thunder K8E S2865
	5300  Tyan S5380 Mainboard
10f2  Achme Computer, Inc.
10f3  Alaris, Inc.
10f4  S-MOS Systems, Inc.
10f5  NKK Corporation
	a001  NDR4000 [NR4600 Bridge]
10f6  Creative Electronic Systems SA
10f7  Matsushita Electric Industrial Co., Ltd.
10f8  Altos India Ltd
10f9  PC Direct
10fa  Truevision
	000c  TARGA 1000
10fb  Thesys Gesellschaft fuer Mikroelektronik mbH
	186f  TH 6255
10fc  I-O Data Device, Inc.
	0003  Cardbus IDE Controller
	0005  Cardbus SCSI CBSC II
10fd  Soyo Computer, Inc
10fe  Fast Multimedia AG
10ff  NCube
1100  Jazz Multimedia
1101  Initio Corporation
	0002  INI-920 Ultra SCSI Adapter
	1060  INI-A100U2W
	1622  INI-1623 PCI SATA-II Controller
	9100  INI-9100/9100W
	9400  INI-940 Fast Wide SCSI Adapter
	9401  INI-935 Fast Wide SCSI Adapter
	9500  INI-950 SCSI Adapter
	9502  INI-950P Ultra Wide SCSI Adapter
1102  Creative Labs
	0002  EMU10k1 [Sound Blaster Live! Series]
	0003  SB AWE64(D)
	0004  EMU10k2/CA0100/CA0102/CA10200 [Sound Blaster Audigy Series]
	0005  EMU20k1 [Sound Blaster X-Fi Series]
	0006  EMU10k1X [SB Live! Value/OEM Series]
	0007  CA0106/CA0111 [SB Live!/Audigy/X-Fi Series]
	0008  CA0108/CA10300 [Sound Blaster Audigy Series]
	0009  CA0110 [Sound Blaster X-Fi Xtreme Audio]
	000b  EMU20k2 [Sound Blaster X-Fi Titanium Series]
	0012  Sound Core3D [Sound Blaster Recon3D / Z-Series]
	4001  SB Audigy FireWire Port
	7002  SB Live! Game Port
	7003  SB Audigy Game Port
	7004  [SB Live! Value] Input device controller
	7005  SB Audigy LS Game Port
	7006  [SB X-Fi Xtreme Audio] CA0110-IBG PCIe to PCI Bridge
	8938  Ectiva EV1938
1103  HighPoint Technologies, Inc.
	0003  HPT343/345/346/363
	0004  HPT366/368/370/370A/372/372N
	0005  HPT372A/372N
	0006  HPT302/302N
	0007  HPT371/371N
	0008  HPT374
	0009  HPT372N
	0620  RocketRAID 620 2 Port SATA-III Controller
	0622  RocketRAID 622 2 Port SATA-III Controller
	0640  RocketRAID 640 4 Port SATA-III Controller
	0641  RocketRAID 640L 4 Port SATA-III Controller
	0642  RocketRAID 642L SATA-III Controller (2 eSATA ports + 2 internal SATA ports)
	0644  RocketRAID 644 4 Port SATA-III Controller (eSATA)
	0645  RocketRAID 644L 4 Port SATA-III Controller (eSATA)
	0646  RocketRAID 644LS SATA-III Controller (4 eSATA devices connected by 1 SAS cable)
	1720  RocketRAID 1720 (2x SATA II RAID Controller)
	1740  RocketRAID 1740
	1742  RocketRAID 1742
	2210  RocketRAID 2210 SATA-II Controller
	2300  RocketRAID 230x 4 Port SATA-II Controller
	2310  RocketRAID 2310 4 Port SATA-II Controller
	2320  RocketRAID 2320 SATA-II Controller
	2322  RocketRAID 2322 SATA-II Controller
	2340  RocketRAID 2340 16 Port SATA-II Controller
	2640  RocketRAID 2640 SAS/SATA Controller
	2722  RocketRAID 2722
	2740  RocketRAID 2740
	2744  RocketRaid 2744
	2782  RocketRAID 2782
	3120  RocketRAID 3120
	3220  RocketRAID 3220
	3320  RocketRAID 3320
	4310  RocketRaid 4310
1104  RasterOps Corp.
1105  Sigma Designs, Inc.
	1105  REALmagic Xcard MPEG 1/2/3/4 DVD Decoder
	8300  REALmagic Hollywood Plus DVD Decoder
	8400  EM840x REALmagic DVD/MPEG-2 Audio/Video Decoder
	8401  EM8401 REALmagic DVD/MPEG-2 A/V Decoder
	8470  EM8470 REALmagic DVD/MPEG-4 A/V Decoder
	8471  EM8471 REALmagic DVD/MPEG-4 A/V Decoder
	8475  EM8475 REALmagic DVD/MPEG-4 A/V Decoder
	8476  EM8476 REALmagic DVD/MPEG-4 A/V Decoder
	8485  EM8485 REALmagic DVD/MPEG-4 A/V Decoder
	8486  EM8486 REALmagic DVD/MPEG-4 A/V Decoder
	c621  EM8621L Digital Media Processor
	c622  EM8622L MPEG-4.10 (H.264) and SMPTE 421M (VC-1) A/V Decoder
1106  VIA Technologies, Inc.
	0102  Embedded VIA Ethernet Controller
	0130  VT6305 1394.A Controller
	0198  P4X600 Host Bridge
	0204  K8M800 Host Bridge
	0208  PT890 Host Bridge
	0238  K8T890 Host Bridge
	0258  PT880 Host Bridge
	0259  CN333/CN400/PM880 Host Bridge
	0269  KT880 Host Bridge
	0282  K8T800Pro Host Bridge
	0290  K8M890 Host Bridge
	0293  PM896 Host Bridge
	0296  P4M800 Host Bridge
	0305  VT8363/8365 [KT133/KM133]
	0308  PT880 Ultra/PT894 Host Bridge
	0314  CN700/VN800/P4M800CE/Pro Host Bridge
	0324  CX700/VX700 Host Bridge
	0327  P4M890 Host Bridge
	0336  K8M890CE Host Bridge
	0340  PT900 Host Bridge
	0351  K8T890CF Host Bridge
	0353  VX800 Host Bridge
	0364  CN896/VN896/P4M900 Host Bridge
	0391  VT8371 [KX133]
	0409  VX855/VX875 Host Bridge: Host Control
	0410  VX900 Host Bridge: Host Control
	0415  VT6415 PATA IDE Host Controller
	0501  VT8501 [Apollo MVP4]
	0505  VT82C505
	0561  VT82C576MV
	0571  VT82C586A/B/VT82C686/A/B/VT823x/A/C PIPC Bus Master IDE
	0576  VT82C576 3V [Apollo Master]
	0581  CX700/VX700 RAID Controller
	0585  VT82C585VP [Apollo VP1/VPX]
	0586  VT82C586/A/B PCI-to-ISA [Apollo VP]
	0591  VT8237A SATA 2-Port Controller
	0595  VT82C595 [Apollo VP2]
	0596  VT82C596 ISA [Mobile South]
	0597  VT82C597 [Apollo VP3]
	0598  VT82C598 [Apollo MVP3]
	0601  VT8601 [Apollo ProMedia]
	0605  VT8605 [ProSavage PM133]
	0680  VT82C680 [Apollo P6]
	0686  VT82C686 [Apollo Super South]
	0691  VT82C693A/694x [Apollo PRO133x]
	0693  VT82C693 [Apollo Pro Plus]
	0698  VT82C693A [Apollo Pro133 AGP]
	0709  VX11 Standard Host Bridge
	070a  VX11 PCI Express Root Port
	070b  VX11 PCI Express Root Port
	070c  VX11 PCI Express Root Port
	070d  VX11 PCI Express Root Port
	070e  VX11 PCI Express Root Port
	0926  VT82C926 [Amazon]
	1000  VT82C570MV
	1106  VT82C570MV
	1122  VX800/VX820 Chrome 9 HC3 Integrated Graphics
	1204  K8M800 Host Bridge
	1208  PT890 Host Bridge
	1238  K8T890 Host Bridge
	1258  PT880 Host Bridge
	1259  CN333/CN400/PM880 Host Bridge
	1269  KT880 Host Bridge
	1282  K8T800Pro Host Bridge
	1290  K8M890 Host Bridge
	1293  PM896 Host Bridge
	1296  P4M800 Host Bridge
	1308  PT894 Host Bridge
	1314  CN700/VN800/P4M800CE/Pro Host Bridge
	1324  CX700/VX700 Host Bridge
	1327  P4M890 Host Bridge
	1336  K8M890CE Host Bridge
	1340  PT900 Host Bridge
	1351  VT3351 Host Bridge
	1353  VX800/VX820 Error Reporting
	1364  CN896/VN896/P4M900 Host Bridge
	1409  VX855/VX875 Error Reporting
	1410  VX900 Error Reporting
	1571  VT82C576M/VT82C586
	1595  VT82C595/97 [Apollo VP2/97]
	1732  VT1732 [Envy24 II] PCI Multi-Channel Audio Controller
	2106  VIA Rhine Family Fast Ethernet Adapter (VT6105)
	2204  K8M800 Host Bridge
	2208  PT890 Host Bridge
	2238  K8T890 Host Bridge
	2258  PT880 Host Bridge
	2259  CN333/CN400/PM880 CPU Host Bridge
	2269  KT880 Host Bridge
	2282  K8T800Pro Host Bridge
	2290  K8M890 Host Bridge
	2293  PM896 Host Bridge
	2296  P4M800 Host Bridge
	2308  PT894 Host Bridge
	2314  CN700/VN800/P4M800CE/Pro Host Bridge
	2324  CX700/VX700 Host Bridge
	2327  P4M890 Host Bridge
	2336  K8M890CE Host Bridge
	2340  PT900 Host Bridge
	2351  VT3351 Host Bridge
	2353  VX800/VX820 Host Bus Control
	2364  CN896/VN896/P4M900 Host Bridge
	2409  VX855/VX875 Host Bus Control
	2410  VX900 CPU Bus Controller
	287a  VT8251 PCI to PCI Bridge
	287b  VT8251 Host Bridge
	287c  VT8251 PCIE Root Port
	287d  VT8251 PCIE Root Port
	287e  VT8237/8251 Ultra VLINK Controller
	3022  CLE266
	3038  VT82xx/62xx UHCI USB 1.1 Controller
	3040  VT82C586B ACPI
	3043  VT86C100A [Rhine]
	3044  VT6306/7/8 [Fire II(M)] IEEE 1394 OHCI Controller
	3050  VT82C596 Power Management
	3051  VT82C596 Power Management
	3053  VT6105M [Rhine-III]
	3057  VT82C686 [Apollo Super ACPI]
	3058  VT82C686 AC97 Audio Controller
	3059  VT8233/A/8235/8237 AC97 Audio Controller
	3065  VT6102/VT6103 [Rhine-II]
	3068  AC'97 Modem Controller
	3074  VT8233 PCI to ISA Bridge
	3091  VT8633 [Apollo Pro266]
	3099  VT8366/A/7 [Apollo KT266/A/333]
	3101  VT8653 Host Bridge
	3102  VT8662 Host Bridge
	3103  VT8615 Host Bridge
	3104  USB 2.0
	3106  VT6105/VT6106S [Rhine-III]
	3108  K8M800/K8N800/K8N800A [S3 UniChrome Pro]
	3109  VT8233C PCI to ISA Bridge
	3112  VT8361 [KLE133] Host Bridge
	3113  VPX/VPX2 PCI to PCI Bridge Controller
	3116  VT8375 [KM266/KL266] Host Bridge
	3118  CN400/PM800/PM880/PN800/PN880 [S3 UniChrome Pro]
	3119  VT6120/VT6121/VT6122 Gigabit Ethernet Adapter
	3122  VT8623 [Apollo CLE266] integrated CastleRock graphics
	3123  VT8623 [Apollo CLE266]
	3128  VT8753 [P4X266 AGP]
	3133  VT3133 Host Bridge
	3142  VT6651 WiFi Adapter, 802.11b
	3147  VT8233A ISA Bridge
	3148  P4M266 Host Bridge
	3149  VIA VT6420 SATA RAID Controller
	3156  P/KN266 Host Bridge
	3157  CX700/VX700 [S3 UniChrome Pro]
	3164  VT6410 ATA133 RAID controller
	3168  P4X333/P4X400/PT800 AGP Bridge
	3177  VT8235 ISA Bridge
	3178  ProSavageDDR P4N333 Host Bridge
	3188  VT8385 [K8T800 AGP] Host Bridge
	3189  VT8377 [KT400/KT600 AGP] Host Bridge
	31b0  VX11 Standard Host Bridge
	31b1  VX11 Standard Host Bridge
	31b2  VX11 DRAM Controller
	31b3  VX11 Power Management Controller
	31b4  VX11 I/O APIC
	31b5  VX11 Scratch Device
	31b7  VX11 Standard Host Bridge
	31b8  VX11 PCI to PCI Bridge
	3204  K8M800 Host Bridge
	3205  VT8378 [KM400/A] Chipset Host Bridge
	3208  PT890 Host Bridge
	3213  VPX/VPX2 PCI to PCI Bridge Controller
	3218  K8T800M Host Bridge
	3227  VT8237 ISA bridge [KT600/K8T800/K8T890 South]
	3230  K8M890CE/K8N890CE [Chrome 9]
	3238  K8T890 Host Bridge
	3249  VT6421 IDE/SATA Controller
	324a  CX700/VX700 PCI to PCI Bridge
	324b  CX700/VX700 Host Bridge
	324e  CX700/VX700 Internal Module Bus
	3253  VT6655 WiFi Adapter, 802.11a/b/g
	3258  PT880 Host Bridge
	3259  CN333/CN400/PM880 Host Bridge
	3260  VIA Chrome9 HC IGP
	3269  KT880 Host Bridge
	3282  K8T800Pro Host Bridge
	3287  VT8251 PCI to ISA Bridge
	3288  VT8237A/VT8251 HDA Controller
	3290  K8M890 Host Bridge
	3296  P4M800 Host Bridge
	3324  CX700/VX700 Host Bridge
	3327  P4M890 Host Bridge
	3336  K8M890CE Host Bridge
	3337  VT8237A PCI to ISA Bridge
	3340  PT900 Host Bridge
	3343  P4M890 [S3 UniChrome Pro]
	3344  CN700/P4M800 Pro/P4M800 CE/VN800 Graphics [S3 UniChrome Pro]
	3349  VT8251 AHCI/SATA 4-Port Controller
	3351  VT3351 Host Bridge
	3353  VX800 PCI to PCI Bridge
	3364  CN896/VN896/P4M900 Host Bridge
	3371  CN896/VN896/P4M900 [Chrome 9 HC]
	3372  VT8237S PCI to ISA Bridge
	337a  VT8237A PCI to PCI Bridge
	337b  VT8237A Host Bridge
	3403  VT6315 Series Firewire Controller
	3409  VX855/VX875 DRAM Bus Control
	3410  VX900 DRAM Bus Control
	3432  VL80x xHCI USB 3.0 Controller
	3456  VX11 Standard Host Bridge
	345b  VX11 Miscellaneous Bus
	3483  VL805 USB 3.0 Host Controller
	3a01  VX11 Graphics [Chrome 645/640]
	4149  VIA VT6420 (ATA133) Controller
	4204  K8M800 Host Bridge
	4208  PT890 Host Bridge
	4238  K8T890 Host Bridge
	4258  PT880 Host Bridge
	4259  CN333/CN400/PM880 Host Bridge
	4269  KT880 Host Bridge
	4282  K8T800Pro Host Bridge
	4290  K8M890 Host Bridge
	4293  PM896 Host Bridge
	4296  P4M800 Host Bridge
	4308  PT894 Host Bridge
	4314  CN700/VN800/P4M800CE/Pro Host Bridge
	4324  CX700/VX700 Host Bridge
	4327  P4M890 Host Bridge
	4336  K8M890CE Host Bridge
	4340  PT900 Host Bridge
	4351  VT3351 Host Bridge
	4353  VX800/VX820 Power Management Control
	4364  CN896/VN896/P4M900 Host Bridge
	4409  VX855/VX875 Power Management Control
	4410  VX900 Power Management and Chip Testing Control
	5030  VT82C596 ACPI [Apollo PRO]
	5122  VX855/VX875 Chrome 9 HCM Integrated Graphics
	5208  PT890 I/O APIC Interrupt Controller
	5238  K8T890 I/O APIC Interrupt Controller
	5287  VT8251 Serial ATA Controller
	5290  K8M890 I/O APIC Interrupt Controller
	5308  PT894 I/O APIC Interrupt Controller
	5324  VX800 Serial ATA and EIDE Controller
	5327  P4M890 I/O APIC Interrupt Controller
	5336  K8M890CE I/O APIC Interrupt Controller
	5340  PT900 I/O APIC Interrupt Controller
	5351  VT3351 I/O APIC Interrupt Controller
	5353  VX800/VX820 APIC and Central Traffic Control
	5364  CN896/VN896/P4M900 I/O APIC Interrupt Controller
	5372  VT8237/8251 Serial ATA Controller
	5409  VX855/VX875 APIC and Central Traffic Control
	5410  VX900 APIC and Central Traffic Control
	6100  VT85C100A [Rhine II]
	6287  SATA RAID Controller
	6290  K8M890CE Host Bridge
	6327  P4M890 Security Device
	6353  VX800/VX820 Scratch Registers
	6364  CN896/VN896/P4M900 Security Device
	6409  VX855/VX875 Scratch Registers
	6410  VX900 Scratch Registers
	7122  VX900 Graphics [Chrome9 HD]
	7204  K8M800 Host Bridge
	7205  KM400/KN400/P4M800 [S3 UniChrome]
	7208  PT890 Host Bridge
	7238  K8T890 Host Bridge
	7258  PT880 Host Bridge
	7259  CN333/CN400/PM880 Host Bridge
	7269  KT880 Host Bridge
	7282  K8T800Pro Host Bridge
	7290  K8M890 Host Bridge
	7293  PM896 Host Bridge
	7296  P4M800 Host Bridge
	7308  PT894 Host Bridge
	7314  CN700/VN800/P4M800CE/Pro Host Bridge
	7324  CX700/VX700 Host Bridge
	7327  P4M890 Host Bridge
	7336  K8M890CE Host Bridge
	7340  PT900 Host Bridge
	7351  VT3351 Host Bridge
	7353  VX800/VX820 North-South Module Interface Control
	7364  CN896/VN896/P4M900 Host Bridge
	7409  VX855/VX875 North-South Module Interface Control
	7410  VX900 North-South Module Interface Control
	8231  VT8231 [PCI-to-ISA Bridge]
	8235  VT8235 ACPI
	8305  VT8363/8365 [KT133/KM133 AGP]
	8324  CX700/VX700 PCI to ISA Bridge
	8353  VX800/VX820 Bus Control and Power Management
	8391  VT8371 [KX133 AGP]
	8400  MVP4
	8409  VX855/VX875 Bus Control and Power Management
	8410  VX900 Bus Control and Power Management
	8500  KLE133/PLE133/PLE133T
	8501  VT8501 [Apollo MVP4 AGP]
	8596  VT82C596 [Apollo PRO AGP]
	8597  VT82C597 [Apollo VP3 AGP]
	8598  VT82C598/694x [Apollo MVP3/Pro133x AGP]
	8601  VT8601 [Apollo ProMedia AGP]
	8605  VT8605 [PM133 AGP]
	8691  VT82C691 [Apollo Pro]
	8693  VT82C693 [Apollo Pro Plus] PCI Bridge
	8a25  PL133/PL133T [S3 ProSavage]
	8a26  KL133/KL133A/KM133/KM133A [S3 ProSavage]
	8d01  PN133/PN133T [S3 Twister]
	8d04  KM266/P4M266/P4M266A/P4N266 [S3 ProSavageDDR]
	9001  VX900 Serial ATA Controller
	9082  Standard AHCI 1.0 SATA Controller
	9140  HDMI Audio Device
	9201  USB3.0 Controller
	9530  Secure Digital Memory Card Controller
	95d0  SDIO Host Controller
	a208  PT890 PCI to PCI Bridge Controller
	a238  K8T890 PCI to PCI Bridge Controller
	a327  P4M890 PCI to PCI Bridge Controller
	a353  VX8xx South-North Module Interface Control
	a364  CN896/VN896/P4M900 PCI to PCI Bridge Controller
	a409  VX855/VX875 USB Device Controller
	a410  VX900 PCI Express Root Port 0
	b091  VT8633 [Apollo Pro266 AGP]
	b099  VT8366/A/7 [Apollo KT266/A/333 AGP]
	b101  VT8653 AGP Bridge
	b102  VT8362 AGP Bridge
	b103  VT8615 AGP Bridge
	b112  VT8361 [KLE133] AGP Bridge
	b113  VPX/VPX2 I/O APIC Interrupt Controller
	b115  VT8363/8365 [KT133/KM133] PCI Bridge
	b168  VT8235 PCI Bridge
	b188  VT8237/8251 PCI bridge [K8M890/K8T800/K8T890 South]
	b198  VT8237/VX700 PCI Bridge
	b213  VPX/VPX2 I/O APIC Interrupt Controller
	b353  VX855/VX875/VX900 PCI to PCI Bridge
	b410  VX900 PCI Express Root Port 1
	b999  [K8T890 North / VT8237 South] PCI Bridge
	c208  PT890 PCI to PCI Bridge Controller
	c238  K8T890 PCI to PCI Bridge Controller
	c327  P4M890 PCI to PCI Bridge Controller
	c340  PT900 PCI to PCI Bridge Controller
	c353  VX800/VX820 PCI Express Root Port
	c364  CN896/VN896/P4M900 PCI to PCI Bridge Controller
	c409  VX855/VX875 EIDE Controller
	c410  VX900 PCI Express Root Port 2
	d104  VT8237R USB UDCI Controller
	d208  PT890 PCI to PCI Bridge Controller
	d213  VPX/VPX2 PCI to PCI Bridge Controller
	d238  K8T890 PCI to PCI Bridge Controller
	d340  PT900 PCI to PCI Bridge Controller
	d410  VX900 PCI Express Root Port 3
	e208  PT890 PCI to PCI Bridge Controller
	e238  K8T890 PCI to PCI Bridge Controller
	e340  PT900 PCI to PCI Bridge Controller
	e353  VX800/VX820 PCI Express Root Port
	e410  VX900 PCI Express Physical Layer Electrical Sub-block
	f208  PT890 PCI to PCI Bridge Controller
	f238  K8T890 PCI to PCI Bridge Controller
	f340  PT900 PCI to PCI Bridge Controller
	f353  VX800/VX820 PCI Express Root Port
1107  Stratus Computers
	0576  VIA VT82C570MV [Apollo] (Wrong vendor ID!)
1108  Proteon, Inc.
	0100  p1690plus_AA
	0101  p1690plus_AB
	0105  P1690Plus
	0108  P1690Plus
	0138  P1690Plus
	0139  P1690Plus
	013c  P1690Plus
	013d  P1690Plus
1109  Cogent Data Technologies, Inc.
	1400  EM110TX [EX110TX]
110a  Siemens AG
	0002  Pirahna 2-port
	0005  Tulip controller, power management, switch extender
	0006  FSC PINC (I/O-APIC)
	0015  FSC Multiprocessor Interrupt Controller
	001d  FSC Copernicus Management Controller
	007b  FSC Remote Service Controller, mailbox device
	007c  FSC Remote Service Controller, shared memory device
	007d  FSC Remote Service Controller, SMIC device
	2101  HST SAPHIR V Primary PCI (ISDN/PMx)
	2102  DSCC4 PEB/PEF 20534 DMA Supported Serial Communication Controller with 4 Channels
	2104  Eicon Diva 2.02 compatible passive ISDN card
	3141  SIMATIC NET CP 5611 / 5621
	3142  SIMATIC NET CP 5613 / 5614
	3143  SIMATIC NET CP 1613
	4021  SIMATIC NET CP 5512 (Profibus and MPI Cardbus Adapter)
	4029  SIMATIC NET CP 5613 A2
	4035  SIMATIC NET CP 1613 A2
	4036  SIMATIC NET CP 1616
	4038  SIMATIC NET CP 1604
	4069  SIMATIC NET CP 5623
	407c  SIMATIC NET CP 5612
	407d  SIMATIC NET CP 5613 A3
	407e  SIMATIC NET CP 5622
	4083  SIMATIC NET CP 5614 A3
	4084  SIMATIC NET CP 1626
	4942  FPGA I-Bus Tracer for MBD
	6120  SZB6120
110b  Chromatic Research Inc.
	0001  Mpact Media Processor
	0004  Mpact 2
110c  Mini-Max Technology, Inc.
110d  Znyx Advanced Systems
110e  CPU Technology
110f  Ross Technology
1110  Powerhouse Systems
	6037  Firepower Powerized SMP I/O ASIC
	6073  Firepower Powerized SMP I/O ASIC
1111  Santa Cruz Operation
1112  Osicom Technologies Inc
	2200  FDDI Adapter
	2300  Fast Ethernet Adapter
	2340  4 Port Fast Ethernet Adapter
	2400  ATM Adapter
1113  Accton Technology Corporation
	1211  SMC2-1211TX
	1216  EN-1216 Ethernet Adapter
	1217  EN-1217 Ethernet Adapter
	5105  10Mbps Network card
	9211  EN-1207D Fast Ethernet Adapter
	9511  21x4x DEC-Tulip compatible Fast Ethernet
	d301  CPWNA100 (Philips wireless PCMCIA)
	ec02  SMC 1244TX v3
	ee23  SMCWPCIT-G 108Mbps Wireless PCI adapter
1114  Atmel Corporation
	0506  at76c506 802.11b Wireless Network Adaptor
1115  3D Labs
1116  Data Translation
	0022  DT3001
	0023  DT3002
	0024  DT3003
	0025  DT3004
	0026  DT3005
	0027  DT3001-PGL
	0028  DT3003-PGL
	0051  DT322
	0060  DT340
	0069  DT332
	80c2  DT3162
1117  Datacube, Inc
	9500  Max-1C SVGA card
	9501  Max-1C image processing
1118  Berg Electronics
1119  ICP Vortex Computersysteme GmbH
	0000  GDT 6000/6020/6050
	0001  GDT 6000B/6010
	0002  GDT 6110/6510
	0003  GDT 6120/6520
	0004  GDT 6530
	0005  GDT 6550
	0006  GDT 6117/6517
	0007  GDT 6127/6527
	0008  GDT 6537
	0009  GDT 6557/6557-ECC
	000a  GDT 6115/6515
	000b  GDT 6125/6525
	000c  GDT 6535
	000d  GDT 6555/6555-ECC
	0100  GDT 6117RP/6517RP
	0101  GDT 6127RP/6527RP
	0102  GDT 6537RP
	0103  GDT 6557RP
	0104  GDT 6111RP/6511RP
	0105  GDT 6121RP/6521RP
	0110  GDT 6117RD/6517RD
	0111  GDT 6127RD/6527RD
	0112  GDT 6537RD
	0113  GDT 6557RD
	0114  GDT 6111RD/6511RD
	0115  GDT 6121RD/6521RD
	0118  GDT 6118RD/6518RD/6618RD
	0119  GDT 6128RD/6528RD/6628RD
	011a  GDT 6538RD/6638RD
	011b  GDT 6558RD/6658RD
	0120  GDT 6117RP2/6517RP2
	0121  GDT 6127RP2/6527RP2
	0122  GDT 6537RP2
	0123  GDT 6557RP2
	0124  GDT 6111RP2/6511RP2
	0125  GDT 6121RP2/6521RP2
	0136  GDT 6113RS/6513RS
	0137  GDT 6123RS/6523RS
	0138  GDT 6118RS/6518RS/6618RS
	0139  GDT 6128RS/6528RS/6628RS
	013a  GDT 6538RS/6638RS
	013b  GDT 6558RS/6658RS
	013c  GDT 6533RS/6633RS
	013d  GDT 6543RS/6643RS
	013e  GDT 6553RS/6653RS
	013f  GDT 6563RS/6663RS
	0166  GDT 7113RN/7513RN/7613RN
	0167  GDT 7123RN/7523RN/7623RN
	0168  GDT 7118RN/7518RN/7518RN
	0169  GDT 7128RN/7528RN/7628RN
	016a  GDT 7538RN/7638RN
	016b  GDT 7558RN/7658RN
	016c  GDT 7533RN/7633RN
	016d  GDT 7543RN/7643RN
	016e  GDT 7553RN/7653RN
	016f  GDT 7563RN/7663RN
	01d6  GDT 4x13RZ
	01d7  GDT 4x23RZ
	01f6  GDT 8x13RZ
	01f7  GDT 8x23RZ
	01fc  GDT 8x33RZ
	01fd  GDT 8x43RZ
	01fe  GDT 8x53RZ
	01ff  GDT 8x63RZ
	0210  GDT 6519RD/6619RD
	0211  GDT 6529RD/6629RD
	0260  GDT 7519RN/7619RN
	0261  GDT 7529RN/7629RN
	02ff  GDT MAXRP
	0300  GDT NEWRX
	0301  GDT NEWRX2
111a  Efficient Networks, Inc
	0000  155P-MF1 (FPGA)
	0002  155P-MF1 (ASIC)
	0003  ENI-25P ATM
	0005  SpeedStream (LANAI)
	0007  SpeedStream ADSL
	1020  SpeedStream PCI 10/100 Network Card
	1203  SpeedStream 1023 Wireless PCI Adapter
111b  Teledyne Electronic Systems
111c  Tricord Systems Inc.
	0001  Powerbis Bridge
111d  Integrated Device Technology, Inc. [IDT]
	0001  IDT77201/77211 155Mbps ATM SAR Controller [NICStAR]
	0003  IDT77222/77252 155Mbps ATM MICRO ABR SAR Controller
	0004  IDT77V252 155Mbps ATM MICRO ABR SAR Controller
	0005  IDT77V222 155Mbps ATM MICRO ABR SAR Controller
	8018  PES12N3A PCI Express Switch
	801c  PES24N3A PCI Express Switch
	8028  PES4T4 PCI Express Switch
	802b  PES8T5A PCI Express Switch
	802c  PES16T4 PCI Express Switch
	802d  PES16T7 PCI Express Switch
	802e  PES24T6 PCI Express Switch
	802f  PES32T8 PCI Express Switch
	8032  PES48T12 PCI Express Switch
	8034  PES16/22/34H16 PCI Express Switch
	8035  PES32H8 PCI Express Switch
	8036  PES48H12 PCI Express Switch
	8037  PES64H16 PCI Express Switch
	8039  PES3T3 PCI Express Switch
	803a  PES4T4 PCI Express Switch
	803c  PES5T5 PCI Express Switch
	803d  PES6T5 PCI Express Switch
	8048  PES8NT2 PCI Express Switch
	8049  PES8NT2 PCI Express Switch
	804a  PES8NT2 PCI Express Internal NTB
	804b  PES8NT2 PCI Express External NTB
	804c  PES16NT2 PCI Express Switch
	804d  PES16NT2 PCI Express Switch
	804e  PES16NT2 PCI Express Internal NTB
	804f  PES16NT2 PCI Express External NTB
	8058  PES12NT3 PCI Express Switch
	8059  PES12NT3 PCI Express Switch
	805a  PES12NT3 PCI Express Internal NTB
	805b  PES12NT3 PCI Express External NTB
	805c  PES24NT3 PCI Express Switch
	805d  PES24NT3 PCI Express Switch
	805e  PES24NT3 PCI Express Internal NTB
	805f  PES24NT3 PCI Express External NTB
	8060  PES16T4G2 PCI Express Gen2 Switch
	8061  PES12T3G2 PCI Express Gen2 Switch
	8068  PES6T6G2 PCI Express Gen2 Switch
	806a  PES24T3G2 PCI Express Gen2 Switch
	806c  PES16T4A/4T4G2 PCI Express Gen2 Switch
	806e  PES24T6G2 PCI Express Gen2 Switch
	806f  HIO524G2 PCI Express Gen2 Switch
	8088  PES32NT8BG2 PCI Express Switch
	808e  PES24NT24G2 PCI Express Switch
	808f  PES32NT8AG2
	80cf  F32P08xG3 [PCIe boot mode]
	80d2  F32P08xG3 NVMe controller
111e  Eldec
111f  Precision Digital Images
	4a47  Precision MX Video engine interface
	5243  Frame capture bus interface
1120  Dell EMC
	2306  Unity Fibre Channel Controller
	2501  Unity Ethernet Controller
	2505  Unity Fibre Channel Controller
1121  Zilog
1122  Multi-tech Systems, Inc.
1123  Excellent Design, Inc.
1124  Leutron Vision AG
	2581  Picport Monochrome
1125  Eurocore
1126  Vigra
1127  FORE Systems Inc
	0200  ForeRunner PCA-200 ATM
	0210  PCA-200PC
	0250  ATM
	0300  ForeRunner PCA-200EPC ATM
	0310  ATM
	0400  ForeRunnerHE ATM Adapter
1129  Firmworks
112a  Hermes Electronics Company, Ltd.
112b  Linotype - Hell AG
112c  Zenith Data Systems
112d  Ravicad
112e  Infomedia Microelectronics Inc.
112f  Dalsa Inc.
	0000  MVC IC-PCI
	0001  MVC IM-PCI Video frame grabber/processor
	0008  PC-CamLink PCI framegrabber
1130  Computervision
1131  Philips Semiconductors
	1561  USB 1.1 Host Controller
	1562  USB 2.0 Host Controller
	3400  SmartPCI56(UCB1500) 56K Modem
	5400  TriMedia TM1000/1100
	5402  TriMedia TM1300
	5405  TriMedia TM1500
	5406  TriMedia TM1700
	540b  PNX1005 Media Processor
	7130  SAA7130 Video Broadcast Decoder
	7133  SAA7131/SAA7133/SAA7135 Video Broadcast Decoder
	7134  SAA7134/SAA7135HL Video Broadcast Decoder
	7145  SAA7145
	7146  SAA7146
	7160  SAA7160
	7162  SAA7162
	7164  SAA7164
	7231  SAA7231
	9730  SAA9730 Integrated Multimedia and Peripheral Controller
1132  Mitel Corp.
1133  Dialogic Corporation
	7701  Eiconcard C90
	7711  Eiconcard C91
	7901  EiconCard S90
	7902  EiconCard S90
	7911  EiconCard S91
	7912  EiconCard S91
	7921  Eiconcard S92
	7941  EiconCard S94
	7942  EiconCard S94
	7943  EiconCard S94
	7944  EiconCard S94
	7945  Eiconcard S94
	7948  Eiconcard S94 64bit/66MHz
	9711  Eiconcard S91 V2
	9911  Eiconcard S91 V2
	9941  Eiconcard S94 V2
	9a41  Eiconcard S94 PCIe
	b921  EiconCard P92
	b922  EiconCard P92
	b923  EiconCard P92
	e001  Diva Pro 2.0 S/T
	e002  Diva 2.0 S/T PCI
	e003  Diva Pro 2.0 U
	e004  Diva 2.0 U PCI
	e005  Diva 2.01 S/T PCI
	e006  Diva CT S/T PCI
	e007  Diva CT U PCI
	e008  Diva CT Lite S/T PCI
	e009  Diva CT Lite U PCI
	e00a  Diva ISDN+V.90 PCI
	e00b  Diva ISDN PCI 2.02
	e00c  Diva 2.02 PCI U
	e00d  Diva Pro 3.0 PCI
	e00e  Diva ISDN+CT S/T PCI Rev 2
	e010  Diva Server BRI-2M PCI
	e011  Diva Server BRI S/T Rev 2
	e012  Diva Server 4BRI-8M PCI
	e013  4BRI
	e014  Diva Server PRI-30M PCI
	e015  Diva PRI PCI v2
	e016  Diva Server Voice 4BRI PCI
	e017  Diva Server Voice 4BRI Rev 2
	e018  BRI
	e019  Diva Server Voice PRI Rev 2
	e01a  Diva BRI-2FX PCI v2
	e01b  Diva Server Voice BRI-2M 2.0 PCI
	e01c  PRI
	e01e  2PRI
	e020  4PRI
	e022  Analog-2
	e024  Analog-4
	e028  Analog-8
	e02a  Diva IPM-300 PCI v1
	e02c  Diva IPM-600 PCI v1
	e02e  4BRI
	e032  BRI
	e034  Diva BRI-CTI PCI v2
1134  Mercury Computer Systems
	0001  Raceway Bridge
	0002  Dual PCI to RapidIO Bridge
	000b  POET Serial RapidIO Bridge
	000d  POET PSDMS Device
1135  Fuji Xerox Co Ltd
	0001  Printer controller
1136  Momentum Data Systems
	0002  PCI-JTAG
1137  Cisco Systems Inc
	0023  VIC 81 PCIe Upstream Port
	0040  VIC PCIe Upstream Port
	0041  VIC PCIe Downstream Port
	0042  VIC Management Controller
	0043  VIC Ethernet NIC
	0044  VIC Ethernet NIC Dynamic
	0045  VIC FCoE HBA
	0046  VIC SCSI Controller
	004e  VIC 82 PCIe Upstream Port
	0071  VIC SR-IOV VF
	007a  VIC 1300 PCIe Upstream Port
	00cf  VIC Userspace NIC
1138  Ziatech Corporation
	8905  8905 [STD 32 Bridge]
1139  Dynamic Pictures, Inc
	0001  VGA Compatible 3D Graphics
113a  FWB Inc
113b  Network Computing Devices
113c  Cyclone Microsystems, Inc.
	0000  PCI-9060 i960 Bridge
	0001  PCI-SDK [PCI i960 Evaluation Platform]
	0911  PCI-911 [i960Jx-based Intelligent I/O Controller]
	0912  PCI-912 [i960CF-based Intelligent I/O Controller]
	0913  PCI-913
	0914  PCI-914 [I/O Controller w/ secondary PCI bus]
113d  Leading Edge Products Inc
113e  Sanyo Electric Co - Computer Engineering Dept
113f  Equinox Systems, Inc.
	0808  SST-64P Adapter
	1010  SST-128P Adapter
	80c0  SST-16P DB Adapter
	80c4  SST-16P RJ Adapter
	80c8  SST-16P Adapter
	8888  SST-4P Adapter
	9090  SST-8P Adapter
1140  Intervoice Inc
1141  Crest Microsystem Inc
1142  Alliance Semiconductor Corporation
	3210  AP6410
	6422  ProVideo 6422
	6424  ProVideo 6424
	6425  ProMotion AT25
	643d  ProMotion AT3D
1143  NetPower, Inc
1144  Cincinnati Milacron
	0001  Noservo controller
1145  Workbit Corporation
	8007  NinjaSCSI-32 Workbit
	f007  NinjaSCSI-32 KME
	f010  NinjaSCSI-32 Workbit
	f012  NinjaSCSI-32 Logitec
	f013  NinjaSCSI-32 Logitec
	f015  NinjaSCSI-32 Melco
	f020  NinjaSCSI-32 Sony PCGA-DVD51
	f021  NinjaPATA-32 Delkin Cardbus UDMA
	f024  NinjaPATA-32 Delkin Cardbus UDMA
	f103  NinjaPATA-32 Delkin Cardbus UDMA
1146  Force Computers
1147  Interface Corp
1148  SysKonnect
	4000  FDDI Adapter
	4200  Token Ring adapter
	4300  SK-9872 Gigabit Ethernet Server Adapter (SK-NET GE-ZX dual link)
	4320  SK-9871 V2.0 Gigabit Ethernet 1000Base-ZX Adapter, PCI64, Fiber ZX/SC
	4400  SK-9Dxx Gigabit Ethernet Adapter
	4500  SK-9Mxx Gigabit Ethernet Adapter
	9000  SK-9S21 10/100/1000Base-T Server Adapter, PCI-X, Copper RJ-45
	9843  [Fujitsu] Gigabit Ethernet
	9e00  SK-9E21D 10/100/1000Base-T Adapter, Copper RJ-45
	9e01  SK-9E21M 10/100/1000Base-T Adapter
1149  Win System Corporation
114a  VMIC
	5565  GE-IP PCI5565,PMC5565 Reflective Memory Node
	5579  VMIPCI-5579 (Reflective Memory Card)
	5587  VMIPCI-5587 (Reflective Memory Card)
	6504  VMIC PCI 7755 FPGA
	7587  VMIVME-7587
114b  Canopus Co., Ltd
114c  Annabooks
114d  IC Corporation
114e  Nikon Systems Inc
114f  Digi International
	0002  AccelePort EPC
	0003  RightSwitch SE-6
	0004  AccelePort Xem
	0005  AccelePort Xr
	0006  AccelePort Xr,C/X
	0009  AccelePort Xr/J
	000a  AccelePort EPC/J
	000c  DataFirePRIme T1 (1-port)
	000d  SyncPort 2-Port (x.25/FR)
	0011  AccelePort 8r EIA-232 (IBM)
	0012  AccelePort 8r EIA-422
	0013  AccelePort Xr
	0014  AccelePort 8r EIA-422
	0015  AccelePort Xem
	0016  AccelePort EPC/X
	0017  AccelePort C/X
	001a  DataFirePRIme E1 (1-port)
	001b  AccelePort C/X (IBM)
	001c  AccelePort Xr (SAIP)
	001d  DataFire RAS T1/E1/PRI
	0023  AccelePort RAS
	0024  DataFire RAS B4 ST/U
	0026  AccelePort 4r 920
	0027  AccelePort Xr 920
	0028  ClassicBoard 4
	0029  ClassicBoard 8
	0034  AccelePort 2r 920
	0035  DataFire DSP T1/E1/PRI cPCI
	0040  AccelePort Xp
	0042  AccelePort 2p
	0043  AccelePort 4p
	0044  AccelePort 8p
	0045  AccelePort 16p
	004e  AccelePort 32p
	0070  Datafire Micro V IOM2 (Europe)
	0071  Datafire Micro V (Europe)
	0072  Datafire Micro V IOM2 (North America)
	0073  Datafire Micro V (North America)
	00b0  Digi Neo 4
	00b1  Digi Neo 8
	00c8  Digi Neo 2 DB9
	00c9  Digi Neo 2 DB9 PRI
	00ca  Digi Neo 2 RJ45
	00cb  Digi Neo 2 RJ45 PRI
	00cc  Digi Neo 1 422
	00cd  Digi Neo 1 422 485
	00ce  Digi Neo 2 422 485
	00d0  ClassicBoard 4 422
	00d1  ClassicBoard 8 422
	00f1  Digi Neo PCI-E 4 port
	00f4  Digi Neo 4 (IBM version)
	6001  Avanstar
1150  Thinking Machines Corp
1151  JAE Electronics Inc.
1152  Megatek
1153  Land Win Electronic Corp
1154  Melco Inc
1155  Pine Technology Ltd
1156  Periscope Engineering
1157  Avsys Corporation
1158  Voarx R & D Inc
	3011  Tokenet/vg 1001/10m anylan
	9050  Lanfleet/Truevalue
	9051  Lanfleet/Truevalue
1159  Mutech Corp
	0001  MV-1000
	0002  MV-1500
115a  Harlequin Ltd
115b  Parallax Graphics
115c  Photron Ltd.
115d  Xircom
	0003  Cardbus Ethernet 10/100
	0005  Cardbus Ethernet 10/100
	0007  Cardbus Ethernet 10/100
	000b  Cardbus Ethernet 10/100
	000c  Mini-PCI V.90 56k Modem
	000f  Cardbus Ethernet 10/100
	00d4  Mini-PCI K56Flex Modem
	0101  Cardbus 56k modem
	0103  Cardbus Ethernet + 56k Modem
115e  Peer Protocols Inc
115f  Maxtor Corporation
1160  Megasoft Inc
1161  PFU Limited
1162  OA Laboratory Co Ltd
1163  Rendition
	0001  Verite 1000
	2000  Verite V2000/V2100/V2200
1164  Advanced Peripherals Technologies
1165  Imagraph Corporation
	0001  Motion TPEG Recorder/Player with audio
1166  Broadcom
	0000  CMIC-LE
	0005  CNB20-LE Host Bridge
	0006  CNB20HE Host Bridge
	0007  CNB20-LE Host Bridge
	0008  CNB20HE Host Bridge
	0009  CNB20LE Host Bridge
	0010  CIOB30
	0011  CMIC-HE
	0012  CMIC-WS Host Bridge (GC-LE chipset)
	0013  CNB20-HE Host Bridge
	0014  CMIC-LE Host Bridge (GC-LE chipset)
	0015  CMIC-GC Host Bridge
	0016  CMIC-GC Host Bridge
	0017  GCNB-LE Host Bridge
	0031  HT1100 HPX0 HT Host Bridge
	0036  BCM5785 [HT1000] PCI/PCI-X Bridge
	0101  CIOB-X2 PCI-X I/O Bridge
	0103  EPB PCI-Express to PCI-X Bridge
	0104  BCM5785 [HT1000] PCI/PCI-X Bridge
	0110  CIOB-E I/O Bridge with Gigabit Ethernet
	0130  BCM5780 [HT2000] PCI-X bridge
	0132  BCM5780 [HT2000] PCI-Express Bridge
	0140  HT2100 PCI-Express Bridge
	0141  HT2100 PCI-Express Bridge
	0142  HT2100 PCI-Express Bridge
	0144  HT2100 PCI-Express Bridge
	0200  OSB4 South Bridge
	0201  CSB5 South Bridge
	0203  CSB6 South Bridge
	0205  BCM5785 [HT1000] Legacy South Bridge
	0211  OSB4 IDE Controller
	0212  CSB5 IDE Controller
	0213  CSB6 RAID/IDE Controller
	0214  BCM5785 [HT1000] IDE
	0217  CSB6 IDE Controller
	021b  HT1100 HD Audio
	0220  OSB4/CSB5 OHCI USB Controller
	0221  CSB6 OHCI USB Controller
	0223  BCM5785 [HT1000] USB
	0225  CSB5 LPC bridge
	0227  GCLE-2 Host Bridge
	0230  CSB5 LPC bridge
	0234  BCM5785 [HT1000] LPC
	0235  BCM5785 [HT1000] XIOAPIC0-2
	0238  BCM5785 [HT1000] WDTimer
	0240  K2 SATA
	0241  RAIDCore RC4000
	0242  RAIDCore BC4000
	024a  BCM5785 [HT1000] SATA (Native SATA Mode)
	024b  BCM5785 [HT1000] SATA (PATA/IDE Mode)
	0406  HT1100 PCI-X Bridge
	0408  HT1100 Legacy Device
	040a  HT1100 ISA-LPC Bridge
	0410  HT1100 SATA Controller (Native SATA Mode)
	0411  HT1100 SATA Controller (PATA / IDE Mode)
	0412  HT1100 USB OHCI Controller
	0414  HT1100 USB EHCI Controller
	0416  HT1100 USB EHCI Controller (with Debug Port)
	0420  HT1100 PCI-Express Bridge
	0421  HT1100 SAS/SATA Controller
	0422  HT1100 PCI-Express Bridge
1167  Mutoh Industries Inc
1168  Thine Electronics Inc
1169  Centre for Development of Advanced Computing
116a  Luminex Software, Inc.
	6100  Bus/Tag Channel
	6800  Escon Channel
	7100  Bus/Tag Channel
	7800  Escon Channel
116b  Connectware Inc
116c  Intelligent Resources Integrated Systems
116d  Martin-Marietta
116e  Electronics for Imaging
116f  Workstation Technology
1170  Inventec Corporation
1171  Loughborough Sound Images Plc
1172  Altera Corporation
1173  Adobe Systems, Inc
1174  Bridgeport Machines
1175  Mitron Computer Inc.
1176  SBE Incorporated
1177  Silicon Engineering
1178  Alfa, Inc.
	afa1  Fast Ethernet Adapter
1179  Toshiba America Info Systems
	0102  Extended IDE Controller
	0103  EX-IDE Type-B
	010f  NVMe Controller
	0404  DVD Decoder card
	0406  Tecra Video Capture device
	0407  DVD Decoder card (Version 2)
	0601  CPU to PCI bridge
	0602  PCI to ISA bridge
	0603  ToPIC95 PCI to CardBus Bridge for Notebooks
	0604  PCI-Docking Host bridge
	060a  ToPIC95
	060f  ToPIC97
	0617  ToPIC100 PCI to Cardbus Bridge with ZV Support
	0618  CPU to PCI and PCI to ISA bridge
	0701  FIR Port Type-O
	0803  TC6371AF SD Host Controller
	0804  TC6371AF SmartMedia Controller
	0805  SD TypA Controller
	0d01  FIR Port Type-DO
117a  A-Trend Technology
117b  L G Electronics, Inc.
117c  ATTO Technology, Inc.
	002c  ExpressSAS R380
	002d  ExpressSAS R348
	0030  Ultra320 SCSI Host Adapter
	0033  SAS Adapter
	0041  ExpressSAS R30F
	8013  ExpressPCI UL4D
	8014  ExpressPCI UL4S
	8027  ExpressPCI UL5D
117d  Becton & Dickinson
117e  T/R Systems
117f  Integrated Circuit Systems
1180  Ricoh Co Ltd
	0465  RL5c465
	0466  RL5c466
	0475  RL5c475
	0476  RL5c476 II
	0477  RL5c477
	0478  RL5c478
	0511  R5C511
	0522  R5C522 IEEE 1394 Controller
	0551  R5C551 IEEE 1394 Controller
	0552  R5C552 IEEE 1394 Controller
	0554  R5C554
	0575  R5C575 SD Bus Host Adapter
	0576  R5C576 SD Bus Host Adapter
	0592  R5C592 Memory Stick Bus Host Adapter
	0811  R5C811
	0822  R5C822 SD/SDIO/MMC/MS/MSPro Host Adapter
	0832  R5C832 IEEE 1394 Controller
	0841  R5C841 CardBus/SD/SDIO/MMC/MS/MSPro/xD/IEEE1394
	0843  R5C843 MMC Host Controller
	0852  xD-Picture Card Controller
	e230  R5U2xx (R5U230 / R5U231 / R5U241) [Memory Stick Host Controller]
	e476  CardBus bridge
	e822  MMC/SD Host Controller
	e823  PCIe SDXC/MMC Host Controller
	e832  R5C832 PCIe IEEE 1394 Controller
	e852  PCIe xD-Picture Card Controller
1181  Telmatics International
1183  Fujikura Ltd
1184  Forks Inc
1185  Dataworld International Ltd
1186  D-Link System Inc
	1002  DL10050 Sundance Ethernet
	1025  AirPlus Xtreme G DWL-G650 Adapter
	1026  AirXpert DWL-AG650 Wireless Cardbus Adapter
	1043  AirXpert DWL-AG650 Wireless Cardbus Adapter
	1300  RTL8139 Ethernet
	1340  DFE-690TXD CardBus PC Card
	1540  DFE-680TX
	1541  DFE-680TXD CardBus PC Card
	1561  DRP-32TXD Cardbus PC Card
	3300  DWL-510 / DWL-610 802.11b [Realtek RTL8180L]
	3a10  AirXpert DWL-AG650 Wireless Cardbus Adapter(rev.B)
	3a11  AirXpert DWL-AG520 Wireless PCI Adapter(rev.B)
	4000  DL2000-based Gigabit Ethernet
	4001  DGE-550SX PCI-X Gigabit Ethernet Adapter
	4200  DFE-520TX Fast Ethernet PCI Adapter
	4300  DGE-528T Gigabit Ethernet Adapter
	4302  DGE-530T Gigabit Ethernet Adapter (rev.C1) [Realtek RTL8169]
	4b00  DGE-560T PCI Express Gigabit Ethernet Adapter
	4b01  DGE-530T Gigabit Ethernet Adapter (rev 11)
	4b02  DGE-560SX PCI Express Gigabit Ethernet Adapter
	4b03  DGE-550T Gigabit Ethernet Adapter V.B1
	4c00  Gigabit Ethernet Adapter
	8400  D-Link DWL-650+ CardBus PC Card
1187  Advanced Technology Laboratories, Inc.
1188  Shima Seiki Manufacturing Ltd.
1189  Matsushita Electronics Co Ltd
118a  Hilevel Technology
118b  Hypertec Pty Limited
118c  Corollary, Inc
	0014  PCIB [C-bus II to PCI bus host bridge chip]
	1117  Intel 8-way XEON Profusion Chipset [Cache Coherency Filter]
118d  BitFlow Inc
	0001  Raptor-PCI framegrabber
	0012  Model 12 Road Runner Frame Grabber
	0014  Model 14 Road Runner Frame Grabber
	0024  Model 24 Road Runner Frame Grabber
	0044  Model 44 Road Runner Frame Grabber
	0112  Model 12 Road Runner Frame Grabber
	0114  Model 14 Road Runner Frame Grabber
	0124  Model 24 Road Runner Frame Grabber
	0144  Model 44 Road Runner Frame Grabber
	0212  Model 12 Road Runner Frame Grabber
	0214  Model 14 Road Runner Frame Grabber
	0224  Model 24 Road Runner Frame Grabber
	0244  Model 44 Road Runner Frame Grabber
	0312  Model 12 Road Runner Frame Grabber
	0314  Model 14 Road Runner Frame Grabber
	0324  Model 24 Road Runner Frame Grabber
	0344  Model 44 Road Runner Frame Grabber
118e  Hermstedt GmbH
118f  Green Logic
1190  Tripace
	c731  TP-910/920/940 PCI Ultra(Wide) SCSI Adapter
1191  Artop Electronic Corp
	0003  SCSI Cache Host Adapter
	0004  ATP8400
	0005  ATP850UF
	0006  ATP860 NO-BIOS
	0007  ATP860
	0008  ATP865 NO-ROM
	0009  ATP865
	000a  ATP867-A
	000b  ATP867-B
	000d  ATP8620
	000e  ATP8620
	8002  AEC6710 SCSI-2 Host Adapter
	8010  AEC6712UW SCSI
	8020  AEC6712U SCSI
	8030  AEC6712S SCSI
	8040  AEC6712D SCSI
	8050  AEC6712SUW SCSI
	8060  AEC6712 SCSI
	8080  AEC67160 SCSI
	8081  AEC67160S SCSI
	808a  AEC67162 2-ch. LVD SCSI
1192  Densan Company Ltd
1193  Zeitnet Inc.
	0001  1221
	0002  1225
1194  Toucan Technology
1195  Ratoc System Inc
1196  Hytec Electronics Ltd
1197  Gage Applied Sciences, Inc.
	010c  CompuScope 82G 8bit 2GS/s Analog Input Card
1198  Lambda Systems Inc
1199  Attachmate Corporation
	0101  Advanced ISCA/PCI Adapter
119a  Mind Share, Inc.
119b  Omega Micro Inc.
	1221  82C092G
119c  Information Technology Inst.
119d  Bug, Inc. Sapporo Japan
119e  Fujitsu Microelectronics Ltd.
	0001  FireStream 155
	0003  FireStream 50
119f  Bull HN Information Systems
	1081  BXI Host Channel Adapter
11a0  Convex Computer Corporation
11a1  Hamamatsu Photonics K.K.
11a2  Sierra Research and Technology
11a3  Deuretzbacher GmbH & Co. Eng. KG
11a4  Barco Graphics NV
11a5  Microunity Systems Eng. Inc
11a6  Pure Data Ltd.
11a7  Power Computing Corp.
11a8  Systech Corp.
11a9  InnoSys Inc.
	4240  AMCC S933Q Intelligent Serial Card
11aa  Actel
11ab  Marvell Technology Group Ltd.
	0146  GT-64010/64010A System Controller
	0f53  88E6318 Link Street network controller
	11ab  MV88SE614x SATA II PCI-E controller
	138f  W8300 802.11 Adapter (rev 07)
	1fa6  Marvell W8300 802.11 Adapter
	1fa7  88W8310 and 88W8000G [Libertas] 802.11g client chipset
	1faa  88w8335 [Libertas] 802.11b/g Wireless
	2211  88SB2211 PCI Express to PCI Bridge
	2a01  88W8335 [Libertas] 802.11b/g Wireless
	2a02  88W8361 [TopDog] 802.11n Wireless
	2a08  88W8362e [TopDog] 802.11a/b/g/n Wireless
	2a0a  88W8363 [TopDog] 802.11n Wireless
	2a0c  88W8363 [TopDog] 802.11n Wireless
	2a24  88W8363 [TopDog] 802.11n Wireless
	2a2b  88W8687 [TopDog] 802.11b/g Wireless
	2a30  88W8687 [TopDog] 802.11b/g Wireless
	2a40  88W8366 [TopDog] 802.11n Wireless
	2a41  88W8366 [TopDog] 802.11n Wireless
	2a42  88W8366 [TopDog] 802.11n Wireless
	2a43  88W8366 [TopDog] 802.11n Wireless
	2a55  88W8864 [Avastar] 802.11ac Wireless
	2b36  88W8764 [Avastar] 802.11n Wireless
	2b38  88W8897 [AVASTAR] 802.11ac Wireless
	2b40  88W8964 [Avastar] 802.11ac Wireless
	4101  OLPC Cafe Controller Secure Digital Controller
	4320  88E8001 Gigabit Ethernet Controller
	4340  88E8021 PCI-X IPMI Gigabit Ethernet Controller
	4341  88E8022 PCI-X IPMI Gigabit Ethernet Controller
	4342  88E8061 PCI-E IPMI Gigabit Ethernet Controller
	4343  88E8062 PCI-E IPMI Gigabit Ethernet Controller
	4344  88E8021 PCI-X IPMI Gigabit Ethernet Controller
	4345  88E8022 PCI-X IPMI Gigabit Ethernet Controller
	4346  88E8061 PCI-E IPMI Gigabit Ethernet Controller
	4347  88E8062 PCI-E IPMI Gigabit Ethernet Controller
	4350  88E8035 PCI-E Fast Ethernet Controller
	4351  88E8036 PCI-E Fast Ethernet Controller
	4352  88E8038 PCI-E Fast Ethernet Controller
	4353  88E8039 PCI-E Fast Ethernet Controller
	4354  88E8040 PCI-E Fast Ethernet Controller
	4355  88E8040T PCI-E Fast Ethernet Controller
	4356  88EC033 PCI-E Fast Ethernet Controller
	4357  88E8042 PCI-E Fast Ethernet Controller
	435a  88E8048 PCI-E Fast Ethernet Controller
	4360  88E8052 PCI-E ASF Gigabit Ethernet Controller
	4361  88E8050 PCI-E ASF Gigabit Ethernet Controller
	4362  88E8053 PCI-E Gigabit Ethernet Controller
	4363  88E8055 PCI-E Gigabit Ethernet Controller
	4364  88E8056 PCI-E Gigabit Ethernet Controller
	4365  88E8070 based Ethernet Controller
	4366  88EC036 PCI-E Gigabit Ethernet Controller
	4367  88EC032 Ethernet Controller
	4368  88EC034 Ethernet Controller
	4369  88EC042 Ethernet Controller
	436a  88E8058 PCI-E Gigabit Ethernet Controller
	436b  88E8071 PCI-E Gigabit Ethernet Controller
	436c  88E8072 PCI-E Gigabit Ethernet Controller
	436d  88E8055 PCI-E Gigabit Ethernet Controller
	4370  88E8075 PCI-E Gigabit Ethernet Controller
	4380  88E8057 PCI-E Gigabit Ethernet Controller
	4381  Yukon Optima 88E8059 [PCIe Gigabit Ethernet Controller with AVB]
	4611  GT-64115 System Controller
	4620  GT-64120/64120A/64121A System Controller
	4801  GT-48001
	5005  Belkin F5D5005 Gigabit Desktop Network PCI Card
	5040  MV88SX5040 4-port SATA I PCI-X Controller
	5041  MV88SX5041 4-port SATA I PCI-X Controller
	5080  MV88SX5080 8-port SATA I PCI-X Controller
	5081  MV88SX5081 8-port SATA I PCI-X Controller
	5181  88f5181 [Orion-1] ARM SoC
	5182  88f5182 [Orion-NAS] ARM SoC
	5281  88f5281 [Orion-2] ARM SoC
	6041  MV88SX6041 4-port SATA II PCI-X Controller
	6042  88SX6042 PCI-X 4-Port SATA-II
	6081  MV88SX6081 8-port SATA II PCI-X Controller
	6101  88SE6101/6102 single-port PATA133 interface
	6111  88SE6111 1-port PATA133(IDE) and 1-port SATA II Controllers
	6121  88SE6121 SATA II / PATA Controller
	6141  88SE614x SATA II PCI-E controller
	6145  88SE6145 SATA II PCI-E controller
	6180  88F6180 [Kirkwood] ARM SoC
	6192  88F6190/6192 [Kirkwood] ARM SoC
	6281  88F6281 [Kirkwood] ARM SoC
	6381  MV78xx0 [Discovery Innovation] ARM SoC
	6440  88SE6440 SAS/SATA PCIe controller
	6450  64560 System Controller
	6460  MV64360/64361/64362 System Controller
	6480  MV64460/64461/64462 System Controller
	6485  MV64460/64461/64462 System Controller, Revision B
	7042  88SX7042 PCI-e 4-port SATA-II
	7810  MV78100 [Discovery Innovation] ARM SoC
	7820  MV78200 [Discovery Innovation] ARM SoC
	7823  MV78230 [Armada XP] ARM SoC
	7846  88F6820 [Armada 385] ARM SoC
	f003  GT-64010 Primary Image Piranha Image Generator
11ac  Canon Information Systems Research Aust.
11ad  Lite-On Communications Inc
	0002  LNE100TX
	c115  LNE100TX [Linksys EtherFast 10/100]
11ae  Aztech System Ltd
11af  Avid Technology Inc.
	0001  Cinema
	ee40  Digidesign Audiomedia III
11b0  V3 Semiconductor Inc.
	0002  V300PSC
	0292  V292PBC [Am29030/40 Bridge]
	0960  V96xPBC
	880a  Deltacast Delta-HD-22
	c960  V96DPC
11b1  Apricot Computers
11b2  Eastman Kodak
11b3  Barr Systems Inc.
11b4  Leitch Technology International
11b5  Radstone Technology Plc
11b6  United Video Corp
11b7  Motorola
11b8  XPoint Technologies, Inc
	0001  Quad PeerMaster
11b9  Pathlight Technology Inc.
	c0ed  SSA Controller
11ba  Videotron Corp
11bb  Pyramid Technology
11bc  Network Peripherals Inc
	0001  NP-PCI
11bd  Pinnacle Systems Inc.
	002e  PCTV 40i
	0040  Royal TS Function 1
	0041  RoyalTS Function 2
	0042  Royal TS Function 3
	0051  PCTV HD 800i
	bede  AV/DV Studio Capture Card
11be  International Microcircuits Inc
11bf  Astrodesign, Inc.
11c0  Hewlett Packard
11c1  LSI Corporation
	0440  56k WinModem
	0441  56k WinModem
	0442  56k WinModem
	0443  LT WinModem
	0444  LT WinModem
	0445  LT WinModem
	0446  LT WinModem
	0447  LT WinModem
	0448  WinModem 56k
	0449  L56xM+S [Mars-2] WinModem 56k
	044a  F-1156IV WinModem (V90, 56KFlex)
	044b  LT WinModem
	044c  LT WinModem
	044d  LT WinModem
	044e  LT WinModem
	044f  V90 WildWire Modem
	0450  LT WinModem
	0451  LT WinModem
	0452  LT WinModem
	0453  LT WinModem
	0454  LT WinModem
	0455  LT WinModem
	0456  LT WinModem
	0457  LT WinModem
	0458  LT WinModem
	0459  LT WinModem
	045a  LT WinModem
	045c  LT WinModem
	0461  V90 WildWire Modem
	0462  V90 WildWire Modem
	0480  Venus Modem (V90, 56KFlex)
	048c  V.92 56K WinModem
	048f  V.92 56k WinModem
	0620  Lucent V.92 Data/Fax Modem
	2600  StarPro26XX family (SP2601, SP2603, SP2612) DSP
	5400  OR3TP12 FPSC
	5656  Venus Modem
	5801  USB
	5802  USS-312 USB Controller
	5803  USS-344S USB Controller
	5811  FW322/323 [TrueFire] 1394a Controller
	5901  FW643 [TrueFire] PCIe 1394b Controller
	5903  FW533 [TrueFire] PCIe 1394a Controller
	8110  T8110 H.100/H.110 TDM switch
	ab10  WL60010 Wireless LAN MAC
	ab11  WL60040 Multimode Wireles LAN MAC
	ab20  ORiNOCO PCI Adapter
	ab21  Agere Wireless PCI Adapter
	ab30  Hermes2 Mini-PCI WaveLAN a/b/g
	ed00  ET-131x PCI-E Ethernet Controller
	ed01  ET-131x PCI-E Ethernet Controller
11c2  Sand Microelectronics
11c3  NEC Corporation
11c4  Document Technologies, Inc
11c5  Shiva Corporation
11c6  Dainippon Screen Mfg. Co. Ltd
11c7  D.C.M. Data Systems
11c8  Dolphin Interconnect Solutions AS
	0658  PSB32 SCI-Adapter D31x
	d665  PSB64 SCI-Adapter D32x
	d667  PSB66 SCI-Adapter D33x
11c9  Magma
	0010  16-line serial port w/- DMA
	0011  4-line serial port w/- DMA
11ca  LSI Systems, Inc
11cb  Specialix Research Ltd.
	2000  PCI_9050
	4000  SUPI_1
	8000  T225
11cc  Michels & Kleberhoff Computer GmbH
11cd  HAL Computer Systems, Inc.
11ce  Netaccess
11cf  Pioneer Electronic Corporation
11d0  Lockheed Martin Federal Systems-Manassas
11d1  Auravision
	01f7  VxP524
	01f9  VxP951
11d2  Intercom Inc.
11d3  Trancell Systems Inc
11d4  Analog Devices
	1535  Blackfin BF535 processor
	1805  SM56 PCI modem
11d5  Ikon Corporation
	0115  10115
	0117  10117
11d6  Tekelec Telecom
11d7  Trenton Technology, Inc.
11d8  Image Technologies Development
11d9  TEC Corporation
11da  Novell
11db  Sega Enterprises Ltd
11dc  Questra Corporation
11dd  Crosfield Electronics Limited
11de  Zoran Corporation
	6017  miroVIDEO DC30
	6057  ZR36057PQC Video cutting chipset
	6120  ZR36120
11df  New Wave PDG
11e0  Cray Communications A/S
11e1  GEC Plessey Semi Inc.
11e2  Samsung Information Systems America
11e3  Quicklogic Corporation
	0001  COM-ON-AIR Dosch&Amand DECT
	0560  QL5064 Companion Design Demo Board
	5030  PC Watchdog
	8417  QL5064 [QuickPCI] PCI v2.2 bridge for SMT417 Dual TMS320C6416T PMC Module
11e4  Second Wave Inc
11e5  IIX Consulting
11e6  Mitsui-Zosen System Research
11e7  Toshiba America, Elec. Company
11e8  Digital Processing Systems Inc.
11e9  Highwater Designs Ltd.
11ea  Elsag Bailey
11eb  Formation Inc.
11ec  Coreco Inc
	000d  Oculus-F/64P
	1800  Cobra/C6
11ed  Mediamatics
11ee  Dome Imaging Systems Inc
11ef  Nicolet Technologies B.V.
11f0  Compu-Shack
	4231  FDDI
	4232  FASTline UTP Quattro
	4233  FASTline FO
	4234  FASTline UTP
	4235  FASTline-II UTP
	4236  FASTline-II FO
	4731  GIGAline
11f1  Symbios Logic Inc
11f2  Picture Tel Japan K.K.
11f3  Keithley Metrabyte
	0011  KPCI-PIO24
11f4  Kinetic Systems Corporation
	2915  CAMAC controller
11f5  Computing Devices International
11f6  Compex
	0112  ENet100VG4
	0113  FreedomLine 100
	1401  ReadyLink 2000
	2011  RL100-ATX 10/100
	2201  ReadyLink 100TX (Winbond W89C840)
	9881  RL100TX Fast Ethernet
11f7  Scientific Atlanta
11f8  PMC-Sierra Inc.
	5220  BR522x [PMC-Sierra maxRAID SAS Controller]
	7364  PM7364 [FREEDM - 32 Frame Engine & Datalink Mgr]
	7375  PM7375 [LASAR-155 ATM SAR]
	7384  PM7384 [FREEDM - 84P672 Frm Engine & Datalink Mgr]
	8000  PM8000  [SPC - SAS Protocol Controller]
	8009  PM8009 SPCve 8x6G
	8032  ATTO Celerity FC8xEN
	8053  PM8053 SXP 12G 24-port SAS/SATA expander
	8054  PM8054 SXP 12G 36-port SAS/SATA expander
	8055  PM8055 SXP 12G 48-port SAS/SATA expander
	8056  PM8056 SXP 12G 68-port SAS/SATA expander
	8060  PM8060 SRCv 12G eight-port SAS/SATA RoC
	8063  PM8063 SRCv 12G 16-port SAS/SATA RoC
	8070  PM8070 Tachyon SPCv 12G eight-port SAS/SATA controller
	8071  PM8071 Tachyon SPCve 12G eight-port SAS/SATA controller
	8072  PM8072 Tachyon SPCv 12G 16-port SAS/SATA controller
	8073  PM8073 Tachyon SPCve 12G 16-port SAS/SATA controller
	8531  PM8531 PFX 24xG3 Fanout PCIe Switches
	8546  PM8546 B-FEIP PSX 96xG3 PCIe Storage Switch
11f9  I-Cube Inc
11fa  Kasan Electronics Company, Ltd.
11fb  Datel Inc
11fc  Silicon Magic
11fd  High Street Consultants
11fe  Comtrol Corporation
	0001  RocketPort 32 port w/external I/F
	0002  RocketPort 8 port w/external I/F
	0003  RocketPort 16 port w/external I/F
	0004  RocketPort 4 port w/quad cable
	0005  RocketPort 8 port w/octa cable
	0006  RocketPort 8 port w/RJ11 connectors
	0007  RocketPort 4 port w/RJ11 connectors
	0008  RocketPort 8 port w/ DB78 SNI (Siemens) connector
	0009  RocketPort 16 port w/ DB78 SNI (Siemens) connector
	000a  RocketPort Plus 4 port
	000b  RocketPort Plus 8 port
	000c  RocketModem 6 port
	000d  RocketModem 4-port
	000e  RocketPort Plus 2 port RS232
	000f  RocketPort Plus 2 port RS422
	0040  RocketPort Infinity Octa, 8port, RJ45
	0041  RocketPort Infinity 32port, External Interface
	0042  RocketPort Infinity 8port, External Interface
	0043  RocketPort Infinity 16port, External Interface
	0044  RocketPort Infinity Quad, 4port, DB
	0045  RocketPort Infinity Octa, 8port, DB
	0047  RocketPort Infinity 4port, RJ45
	004f  RocketPort Infinity 2port, SMPTE
	0052  RocketPort Infinity Octa, 8port, SMPTE
	0801  RocketPort UPCI 32 port w/external I/F
	0802  RocketPort UPCI 8 port w/external I/F
	0803  RocketPort UPCI 16 port w/external I/F
	0805  RocketPort UPCI 8 port w/octa cable
	080c  RocketModem III 8 port
	080d  RocketModem III 4 port
	0810  RocketPort UPCI Plus 4 port RS232
	0811  RocketPort UPCI Plus 8 port RS232
	0812  RocketPort UPCI Plus 8 port RS422
	0903  RocketPort Compact PCI 16 port w/external I/F
	8015  RocketPort 4-port UART 16954
11ff  Scion Corporation
	0003  AG-5
1200  CSS Corporation
1201  Vista Controls Corp
1202  Network General Corp.
	4300  Gigabit Ethernet Adapter
1203  Bayer Corporation, Agfa Division
1204  Lattice Semiconductor Corporation
	1965  SB6501 802.11ad Wireless Network Adapter
1205  Array Corporation
1206  Amdahl Corporation
1208  Parsytec GmbH
	4853  HS-Link Device
1209  SCI Systems Inc
120a  Synaptel
120b  Adaptive Solutions
120c  Technical Corp.
120d  Compression Labs, Inc.
120e  Cyclades Corporation
	0100  Cyclom-Y below first megabyte
	0101  Cyclom-Y above first megabyte
	0102  Cyclom-4Y below first megabyte
	0103  Cyclom-4Y above first megabyte
	0104  Cyclom-8Y below first megabyte
	0105  Cyclom-8Y above first megabyte
	0200  Cyclades-Z below first megabyte
	0201  Cyclades-Z above first megabyte
	0300  PC300/RSV or /X21 (2 ports)
	0301  PC300/RSV or /X21 (1 port)
	0310  PC300/TE (2 ports)
	0311  PC300/TE (1 port)
	0320  PC300/TE-M (2 ports)
	0321  PC300/TE-M (1 port)
	0400  PC400
120f  Essential Communications
	0001  Roadrunner serial HIPPI
1210  Hyperparallel Technologies
1211  Braintech Inc
1212  Kingston Technology Corp.
1213  Applied Intelligent Systems, Inc.
1214  Performance Technologies, Inc.
1215  Interware Co., Ltd
1216  Purup Prepress A/S
1217  O2 Micro, Inc.
	00f7  Firewire (IEEE 1394)
	10f7  1394 OHCI Compliant Host Controller
	11f7  OZ600 1394a-2000 Controller
	13f7  1394 OHCI Compliant Host Controller
	6729  OZ6729
	673a  OZ6730
	6832  OZ6832/6833 CardBus Controller
	6836  OZ6836/6860 CardBus Controller
	6872  OZ6812 CardBus Controller
	6925  OZ6922 CardBus Controller
	6933  OZ6933/711E1 CardBus/SmartCardBus Controller
	6972  OZ601/6912/711E0 CardBus/SmartCardBus Controller
	7110  OZ711Mx 4-in-1 MemoryCardBus Accelerator
	7112  OZ711EC1/M1 SmartCardBus/MemoryCardBus Controller
	7113  OZ711EC1 SmartCardBus Controller
	7114  OZ711M1/MC1 4-in-1 MemoryCardBus Controller
	7120  Integrated MMC/SD Controller
	7130  Integrated MS/xD Controller
	7134  OZ711MP1/MS1 MemoryCardBus Controller
	7135  Cardbus bridge
	7136  OZ711SP1 Memory CardBus Controller
	71e2  OZ711E2 SmartCardBus Controller
	7212  OZ711M2 4-in-1 MemoryCardBus Controller
	7213  OZ6933E CardBus Controller
	7223  OZ711M3/MC3 4-in-1 MemoryCardBus Controller
	7233  OZ711MP3/MS3 4-in-1 MemoryCardBus Controller
	8120  Integrated MMC/SD Controller
	8130  Integrated MS/MSPRO/xD Controller
	8220  OZ600FJ1/OZ900FJ1 SD/MMC Card Reader Controller
	8221  OZ600FJ0/OZ900FJ0/OZ600FJS SD/MMC Card Reader Controller
	8320  OZ600RJ1/OZ900RJ1 SD/MMC Card Reader Controller
	8321  OZ600RJ0/OZ900RJ0/OZ600RJS SD/MMC Card Reader Controller
	8330  OZ600 MS/xD Controller
	8331  O2 Flash Memory Card
	8520  SD/MMC Card Reader Controller
	8621  SD/MMC Card Reader Controller
1218  Hybricon Corp.
1219  First Virtual Corporation
121a  3Dfx Interactive, Inc.
	0001  Voodoo
	0002  Voodoo 2
	0003  Voodoo Banshee
	0004  Voodoo Banshee [Velocity 100]
	0005  Voodoo 3
	0009  Voodoo 4 / Voodoo 5
	0057  Voodoo 3/3000 [Avenger]
121b  Advanced Telecommunications Modules
121c  Nippon Texaco., Ltd
121d  LiPPERT ADLINK Technology GmbH
121e  CSPI
	0201  Myrinet 2000 Scalable Cluster Interconnect
121f  Arcus Technology, Inc.
1220  Ariel Corporation
	1220  AMCC 5933 TMS320C80 DSP/Imaging board
1221  Contec Co., Ltd
	9172  PO-64L(PCI)H [Isolated Digital Output Board for PCI]
	91a2  PO-32L(PCI)H [Isolated Digital Output Board for PCI]
	91c3  DA16-16(LPCI)L [Un-insulated highly precise analog output board for Low Profile PCI]
	b152  DIO-96D2-LPCI
	c103  ADA16-32/2(PCI)F [High-Speed Analog I/O Board for PCI]
1222  Ancor Communications, Inc.
1223  Artesyn Communication Products
	0003  PM/Link
	0004  PM/T1
	0005  PM/E1
	0008  PM/SLS
	0009  BajaSpan Resource Target
	000a  BajaSpan Section 0
	000b  BajaSpan Section 1
	000c  BajaSpan Section 2
	000d  BajaSpan Section 3
	000e  PM/PPC
1224  Interactive Images
1225  Power I/O, Inc.
1227  Tech-Source
	0006  Raptor GFX 8P
	0023  Raptor GFX [1100T]
	0045  Raptor 4000-L [Linux version]
	004a  Raptor 4000-LR-L [Linux version]
1228  Norsk Elektro Optikk A/S
1229  Data Kinesis Inc.
122a  Integrated Telecom
122b  LG Industrial Systems Co., Ltd
122c  Sican GmbH
122d  Aztech System Ltd
	1206  368DSP
	1400  Trident PCI288-Q3DII (NX)
	50dc  3328 Audio
	80da  3328 Audio
122e  Xyratex
	7722  Napatech XL1
	7724  Napatech XL2/XA
	7729  Napatech XD
122f  Andrew Corporation
1230  Fishcamp Engineering
1231  Woodward McCoach, Inc.
	04e1  Desktop PCI Telephony 4
	05e1  Desktop PCI Telephony 5/6
	0d00  LightParser
	0d02  LightParser 2
	0d13  Desktop PCI L1/L3 Telephony
1232  GPT Limited
1233  Bus-Tech, Inc.
1235  Risq Modular Systems, Inc.
1236  Sigma Designs Corporation
	0000  RealMagic64/GX
	6401  REALmagic 64/GX (SD 6425)
1237  Alta Technology Corporation
1238  Adtran
1239  3DO Company
123a  Visicom Laboratories, Inc.
123b  Seeq Technology, Inc.
123c  Century Systems, Inc.
123d  Engineering Design Team, Inc.
	0000  EasyConnect 8/32
	0002  EasyConnect 8/64
	0003  EasyIO
123e  Simutech, Inc.
123f  LSI Logic
	00e4  MPEG
	8120  DVxplore Codec
	8888  Cinemaster C 3.0 DVD Decoder
1240  Marathon Technologies Corp.
1241  DSC Communications
1242  JNI Corporation
	1560  JNIC-1560 PCI-X Fibre Channel Controller
	4643  FCI-1063 Fibre Channel Adapter
	6562  FCX2-6562 Dual Channel PCI-X Fibre Channel Adapter
	656a  FCX-6562 PCI-X Fibre Channel Adapter
1243  Delphax
1244  AVM GmbH
	0700  B1 ISDN
	0800  C4 ISDN
	0a00  A1 ISDN [Fritz]
	0e00  Fritz!Card PCI v2.0 ISDN
	0e80  Fritz!Card PCI v2.1 ISDN
	1100  C2 ISDN
	1200  T1 ISDN
	2700  Fritz!Card DSL SL
	2900  Fritz!Card DSL v2.0
1245  A.P.D., S.A.
1246  Dipix Technologies, Inc.
1247  Xylon Research, Inc.
1248  Central Data Corporation
1249  Samsung Electronics Co., Ltd.
124a  AEG Electrocom GmbH
124b  SBS/Greenspring Modular I/O
	0040  PCI-40A or cPCI-200 Quad IndustryPack carrier
124c  Solitron Technologies, Inc.
124d  Stallion Technologies, Inc.
	0000  EasyConnection 8/32
	0002  EasyConnection 8/64
	0003  EasyIO
	0004  EasyConnection/RA
124e  Cylink
124f  Infortrend Technology, Inc.
	0041  IFT-2000 Series RAID Controller
1250  Hitachi Microcomputer System Ltd
1251  VLSI Solutions Oy
1253  Guzik Technical Enterprises
1254  Linear Systems Ltd.
	0065  DVB Master FD
	007c  DVB Master Quad/o
1255  Optibase Ltd
	1110  MPEG Forge
	1210  MPEG Fusion
	2110  VideoPlex
	2120  VideoPlex CC
	2130  VideoQuest
1256  Perceptive Solutions, Inc.
	4201  PCI-2220I
	4401  PCI-2240I
	5201  PCI-2000
1257  Vertex Networks, Inc.
1258  Gilbarco, Inc.
1259  Allied Telesis
	2560  AT-2560 Fast Ethernet Adapter (i82557B)
	2801  AT-2801FX (RTL-8139)
	a117  RTL81xx Fast Ethernet
	a11e  RTL81xx Fast Ethernet
	a120  21x4x DEC-Tulip compatible 10/100 Ethernet
125a  ABB Power Systems
125b  Asix Electronics Corporation
	1400  AX88141 Fast Ethernet Controller
125c  Aurora Technologies, Inc.
	0101  Saturn 4520P
	0640  Aries 16000P
125d  ESS Technology
	0000  ES336H Fax Modem (Early Model)
	1948  ES1948 Maestro-1
	1968  ES1968 Maestro 2
	1969  ES1938/ES1946/ES1969 Solo-1 Audiodrive
	1978  ES1978 Maestro 2E
	1988  ES1988 Allegro-1
	1989  ESS Modem
	1998  ES1983S Maestro-3i PCI Audio Accelerator
	1999  ES1983S Maestro-3i PCI Modem Accelerator
	199a  ES1983S Maestro-3i PCI Audio Accelerator
	199b  ES1983S Maestro-3i PCI Modem Accelerator
	2808  ES336H Fax Modem (Later Model)
	2838  ES2838/2839 SuperLink Modem
	2898  ES2898 Modem
125e  Specialvideo Engineering SRL
125f  Concurrent Technologies, Inc.
	2071  CC PMC/232
	2084  CC PMC/23P
	2091  CC PMC/422
1260  Intersil Corporation
	3872  ISL3872 [Prism 3]
	3873  ISL3874 [Prism 2.5]/ISL3872 [Prism 3]
	3877  ISL3877 [Prism Indigo]
	3886  ISL3886 [Prism Javelin/Prism Xbow]
	3890  ISL3890 [Prism GT/Prism Duette]/ISL3886 [Prism Javelin/Prism Xbow]
	8130  HMP8130 NTSC/PAL Video Decoder
	8131  HMP8131 NTSC/PAL Video Decoder
	ffff  ISL3886IK
1261  Matsushita-Kotobuki Electronics Industries, Ltd.
1262  ES Computer Company, Ltd.
1263  Sonic Solutions
1264  Aval Nagasaki Corporation
1265  Casio Computer Co., Ltd.
1266  Microdyne Corporation
	0001  NE10/100 Adapter (i82557B)
	1910  NE2000Plus (RT8029) Ethernet Adapter
1267  S. A. Telecommunications
	5352  PCR2101
	5a4b  Telsat Turbo
1268  Tektronix
1269  Thomson-CSF/TTM
126a  Lexmark International, Inc.
126b  Adax, Inc.
126c  Northern Telecom
	1211  10/100BaseTX [RTL81xx]
	126c  802.11b Wireless Ethernet Adapter
126d  Splash Technology, Inc.
126e  Sumitomo Metal Industries, Ltd.
126f  Silicon Motion, Inc.
	0501  SM501 VoyagerGX Rev. AA
	0510  SM501 VoyagerGX Rev. B
	0710  SM710 LynxEM
	0712  SM712 LynxEM+
	0718  SM718 LynxSE+
	0720  SM720 Lynx3DM
	0730  SM731 Cougar3DR
	0750  SM750
	0810  SM810 LynxE
	0811  SM811 LynxE
	0820  SM820 Lynx3D
	0910  SM910
1270  Olympus Optical Co., Ltd.
1271  GW Instruments
1272  Telematics International
1273  Hughes Network Systems
	0002  DirecPC
1274  Ensoniq
	1171  ES1373 / Creative Labs CT5803 [AudioPCI]
	1371  ES1371/ES1373 / Creative Labs CT2518
	5000  ES1370 [AudioPCI]
	5880  5880B / Creative Labs CT5880
	8001  CT5880 [AudioPCI]
	8002  5880A [AudioPCI]
1275  Network Appliance Corporation
1276  Switched Network Technologies, Inc.
1277  Comstream
1278  Transtech Parallel Systems Ltd.
	0701  TPE3/TM3 PowerPC Node
	0710  TPE5 PowerPC PCI board
	1100  PMC-FPGA02
	1101  TS-C43 card with 4 ADSP-TS101 processors
1279  Transmeta Corporation
	0060  TM8000 Northbridge
	0061  TM8000 AGP bridge
	0295  Northbridge
	0395  LongRun Northbridge
	0396  SDRAM controller
	0397  BIOS scratchpad
127a  Rockwell International
	1002  HCF 56k Data/Fax Modem
	1003  HCF 56k Data/Fax Modem
	1004  HCF 56k Data/Fax/Voice Modem
	1005  HCF 56k Data/Fax/Voice/Spkp (w/Handset) Modem
	1022  HCF 56k Modem
	1023  HCF 56k Data/Fax Modem
	1024  HCF 56k Data/Fax/Voice Modem
	1025  HCF 56k Data/Fax/Voice/Spkp (w/Handset) Modem
	1026  HCF 56k PCI Speakerphone Modem
	1032  HCF 56k Modem
	1033  HCF 56k Modem
	1034  HCF 56k Modem
	1035  HCF 56k PCI Speakerphone Modem
	1036  HCF 56k Modem
	1085  HCF 56k Volcano PCI Modem
	2004  HSF 56k Data/Fax/Voice/Spkp (w/Handset) Modem
	2005  HCF 56k Data/Fax Modem
	2013  HSF 56k Data/Fax Modem
	2014  HSF 56k Data/Fax/Voice Modem
	2015  HSF 56k Data/Fax/Voice/Spkp (w/Handset) Modem
	2016  HSF 56k Data/Fax/Voice/Spkp Modem
	4311  Riptide HSF 56k PCI Modem
	4320  Riptide PCI Audio Controller
	4321  Riptide HCF 56k PCI Modem
	4322  Riptide PCI Game Controller
	8234  RapidFire 616X ATM155 Adapter
127b  Pixera Corporation
127c  Crosspoint Solutions, Inc.
127d  Vela Research
127e  Winnov, L.P.
	0010  Videum 1000 Plus
127f  Fujifilm
1280  Photoscript Group Ltd.
1281  Yokogawa Electric Corporation
1282  Davicom Semiconductor, Inc.
	6585  DM562P V90 Modem
	9009  Ethernet 100/10 MBit
	9100  21x4x DEC-Tulip compatible 10/100 Ethernet
	9102  21x4x DEC-Tulip compatible 10/100 Ethernet
	9132  Ethernet 100/10 MBit
1283  Integrated Technology Express, Inc.
	673a  IT8330G
	8152  IT8152F/G Advanced RISC-to-PCI Companion Chip
	8211  ITE 8211F Single Channel UDMA 133
	8212  IT8212 Dual channel ATA RAID controller
	8213  IT8213 IDE Controller
	8330  IT8330G
	8872  IT887xF PCI to ISA I/O chip with SMB, GPIO, Serial or Parallel Port
	8888  IT8888F/G PCI to ISA Bridge with SMB [Golden Gate]
	8889  IT8889F PCI to ISA Bridge
	8893  IT8893E PCIe to PCI Bridge
	e886  IT8330G
1284  Sahara Networks, Inc.
1285  Platform Technologies, Inc.
	0100  AGOGO sound chip (aka ESS Maestro 1)
1286  Mazet GmbH
1287  M-Pact, Inc.
	001e  LS220D DVD Decoder
	001f  LS220C DVD Decoder
1288  Timestep Corporation
1289  AVC Technology, Inc.
128a  Asante Technologies, Inc.
128b  Transwitch Corporation
128c  Retix Corporation
128d  G2 Networks, Inc.
	0021  ATM155 Adapter
128e  Hoontech Corporation/Samho Multi Tech Ltd.
	0008  ST128 WSS/SB
	0009  ST128 SAM9407
	000a  ST128 Game Port
	000b  ST128 MPU Port
	000c  ST128 Ctrl Port
128f  Tateno Dennou, Inc.
1290  Sord Computer Corporation
1291  NCS Computer Italia
1292  Tritech Microelectronics Inc
	fc02  Pyramid3D TR25202
1293  Media Reality Technology
1294  Rhetorex, Inc.
1295  Imagenation Corporation
	0800  PXR800
	1000  PXD1000
1296  Kofax Image Products
1297  Holco Enterprise Co, Ltd/Shuttle Computer
1298  Spellcaster Telecommunications Inc.
1299  Knowledge Technology Lab.
129a  VMetro, inc.
	0615  PBT-615 PCI-X Bus Analyzer
	1100  PMC-FPGA05
	1106  XMC-FPGA05F, PCI interface
	1107  XMC-FPGA05F, PCIe interface
	1108  XMC-FPGA05D, PCI interface
	1109  XMC-FPGA05D, PCIe interface
129b  Image Access
129c  Jaycor
129d  Compcore Multimedia, Inc.
129e  Victor Company of Japan, Ltd.
129f  OEC Medical Systems, Inc.
12a0  Allen-Bradley Company
12a1  Simpact Associates, Inc.
12a2  Newgen Systems Corporation
12a3  Lucent Technologies
	8105  T8105 H100 Digital Switch
12a4  NTT Electronics Technology Company
12a5  Vision Dynamics Ltd.
12a6  Scalable Networks, Inc.
12a7  AMO GmbH
12a8  News Datacom
12a9  Xiotech Corporation
12aa  SDL Communications, Inc.
12ab  YUAN High-Tech Development Co., Ltd.
	0000  MPG160/Kuroutoshikou ITVC15-STVLP
	0002  AU8830 [Vortex2] Based Sound Card With A3D Support
	0003  T507 (DVB-T) TV tuner/capture device
	2300  Club-3D Zap TV2100
	3000  MPG-200C PCI DVD Decoder Card
	4789  MPC788 MiniPCI Hybrid TV Tuner
	fff3  MPG600/Kuroutoshikou ITVC16-STVLP
	ffff  MPG600/Kuroutoshikou ITVC16-STVLP
12ac  Measurex Corporation
12ad  Multidata GmbH
12ae  Alteon Networks Inc.
	0001  AceNIC Gigabit Ethernet
	0002  AceNIC Gigabit Ethernet (Copper)
	00fa  Farallon PN9100-T Gigabit Ethernet
12af  TDK USA Corp
12b0  Jorge Scientific Corp
12b1  GammaLink
12b2  General Signal Networks
12b3  Inter-Face Co Ltd
12b4  FutureTel Inc
12b5  Granite Systems Inc.
12b6  Natural Microsystems
12b7  Cognex Modular Vision Systems Div. - Acumen Inc.
12b8  Korg
12b9  3Com Corp, Modem Division
	1006  WinModem
	1007  USR 56k Internal WinModem
	1008  56K FaxModem Model 5610
12ba  BittWare, Inc.
12bb  Nippon Unisoft Corporation
12bc  Array Microsystems
12bd  Computerm Corp.
12be  Anchor Chips Inc.
	3041  AN3041Q CO-MEM
	3042  AN3042Q CO-MEM Lite
12bf  Fujifilm Microdevices
12c0  Infimed
12c1  GMM Research Corp
12c2  Mentec Limited
12c3  Holtek Microelectronics Inc
	0058  PCI NE2K Ethernet
	5598  PCI NE2K Ethernet
12c4  Connect Tech Inc
	0001  Blue HEAT/PCI 8 (RS232/CL/RJ11)
	0002  Blue HEAT/PCI 4 (RS232)
	0003  Blue HEAT/PCI 2 (RS232)
	0004  Blue HEAT/PCI 8 (UNIV, RS485)
	0005  Blue HEAT/PCI 4+4/6+2 (UNIV, RS232/485)
	0006  Blue HEAT/PCI 4 (OPTO, RS485)
	0007  Blue HEAT/PCI 2+2 (RS232/485)
	0008  Blue HEAT/PCI 2 (OPTO, Tx, RS485)
	0009  Blue HEAT/PCI 2+6 (RS232/485)
	000a  Blue HEAT/PCI 8 (Tx, RS485)
	000b  Blue HEAT/PCI 4 (Tx, RS485)
	000c  Blue HEAT/PCI 2 (20 MHz, RS485)
	000d  Blue HEAT/PCI 2 PTM
	0100  NT960/PCI
	0201  cPCI Titan - 2 Port
	0202  cPCI Titan - 4 Port
	0300  CTI PCI UART 2 (RS232)
	0301  CTI PCI UART 4 (RS232)
	0302  CTI PCI UART 8 (RS232)
	0310  CTI PCI UART 1+1 (RS232/485)
	0311  CTI PCI UART 2+2 (RS232/485)
	0312  CTI PCI UART 4+4 (RS232/485)
	0320  CTI PCI UART 2
	0321  CTI PCI UART 4
	0322  CTI PCI UART 8
	0330  CTI PCI UART 2 (RS485)
	0331  CTI PCI UART 4 (RS485)
	0332  CTI PCI UART 8 (RS485)
12c5  Picture Elements Incorporated
	007e  Imaging/Scanning Subsystem Engine
	007f  Imaging/Scanning Subsystem Engine
	0081  PCIVST [Grayscale Thresholding Engine]
	0085  Video Simulator/Sender
	0086  THR2 Multi-scale Thresholder
12c6  Mitani Corporation
12c7  Dialogic Corp
	0546  Springware D/120JCT-LS
	0647  Springware D/240JCT-T1
	0676  Springware D/41JCT-LS
	0685  Springware D/480JCT-2T1
12c8  G Force Co, Ltd
12c9  Gigi Operations
12ca  Integrated Computing Engines
12cb  Antex Electronics Corporation
	0027  SC4 (StudioCard)
	002e  StudioCard 2000
12cc  Pluto Technologies International
12cd  Aims Lab
12ce  Netspeed Inc.
12cf  Prophet Systems, Inc.
12d0  GDE Systems, Inc.
12d1  PSITech
12d2  NVidia / SGS Thomson (Joint Venture)
	0008  NV1
	0009  DAC64
	0018  Riva128
	0019  Riva128ZX
	0020  TNT
	0028  TNT2
	0029  UTNT2
	002c  VTNT2
	00a0  ITNT2
12d3  Vingmed Sound A/S
12d4  Ulticom (Formerly DGM&S)
	0200  T1 Card
12d5  Equator Technologies Inc
	0003  BSP16
	1000  BSP15
12d6  Analogic Corp
12d7  Biotronic SRL
12d8  Pericom Semiconductor
	01a7  7C21P100 2-port PCI-X to PCI-X Bridge
	2608  PI7C9X2G608GP PCIe2 6-Port/8-Lane Packet Switch
	400a  PI7C9X442SL PCI Express Bridge Port
	400e  PI7C9X442SL USB OHCI Controller
	400f  PI7C9X442SL USB EHCI Controller
	71e2  PI7C7300A/PI7C7300D PCI-to-PCI Bridge
	71e3  PI7C7300A/PI7C7300D PCI-to-PCI Bridge (Secondary Bus 2)
	8140  PI7C8140A PCI-to-PCI Bridge
	8148  PI7C8148A/PI7C8148B PCI-to-PCI Bridge
	8150  PCI to PCI Bridge
	8152  PI7C8152A/PI7C8152B/PI7C8152BI PCI-to-PCI Bridge
	8154  PI7C8154A/PI7C8154B/PI7C8154BI PCI-to-PCI Bridge
	e110  PI7C9X110 PCI Express to PCI bridge
	e111  PI7C9X111SL PCIe-to-PCI Reversible Bridge
	e130  PCI Express to PCI-XPI7C9X130 PCI-X Bridge
12d9  Aculab PLC
	0002  PCI Prosody
	0004  cPCI Prosody
	0005  Aculab E1/T1 PCI card
	1078  Prosody X class e1000 device
12da  True Time Inc.
12db  Annapolis Micro Systems, Inc
12dc  Symicron Computer Communication Ltd.
12dd  Management Graphics
12de  Rainbow Technologies
	0200  CryptoSwift CS200
12df  SBS Technologies Inc
12e0  Chase Research
	0010  ST16C654 Quad UART
	0020  ST16C654 Quad UART
	0030  ST16C654 Quad UART
12e1  Nintendo Co, Ltd
12e2  Datum Inc. Bancomm-Timing Division
12e3  Imation Corp - Medical Imaging Systems
12e4  Brooktrout Technology Inc
12e5  Apex Semiconductor Inc
12e6  Cirel Systems
12e7  Sunsgroup Corporation
12e8  Crisc Corp
12e9  GE Spacenet
12ea  Zuken
12eb  Aureal Semiconductor
	0001  Vortex 1
	0002  Vortex 2
	0003  AU8810 Vortex Digital Audio Processor
	8803  Vortex 56k Software Modem
12ec  3A International, Inc.
12ed  Optivision Inc.
12ee  Orange Micro
12ef  Vienna Systems
12f0  Pentek
12f1  Sorenson Vision Inc
12f2  Gammagraphx, Inc.
12f3  Radstone Technology
12f4  Megatel
12f5  Forks
12f6  Dawson France
12f7  Cognex
12f8  Electronic Design GmbH
	0002  VideoMaker
12f9  Four Fold Ltd
12fb  Spectrum Signal Processing
	0001  PMC-MAI
	00f5  F5 Dakar
	02ad  PMC-2MAI
	2adc  ePMC-2ADC
	3100  PRO-3100
	3500  PRO-3500
	4d4f  Modena
	8120  ePMC-8120
	da62  Daytona C6201 PCI (Hurricane)
	db62  Ingliston XBIF
	dc62  Ingliston PLX9054
	dd62  Ingliston JTAG/ISP
	eddc  ePMC-MSDDC
	fa01  ePMC-FPGA
12fc  Capital Equipment Corp
12fd  I2S
12fe  ESD Electronic System Design GmbH
12ff  Lexicon
1300  Harman International Industries Inc
1302  Computer Sciences Corp
1303  Innovative Integration
	0030  X3-SDF 4-channel XMC acquisition board
1304  Juniper Networks
1305  Netphone, Inc
1306  Duet Technologies
1307  Measurement Computing
	0001  PCI-DAS1602/16
	000b  PCI-DIO48H
	000c  PCI-PDISO8
	000d  PCI-PDISO16
	000f  PCI-DAS1200
	0010  PCI-DAS1602/12
	0014  PCI-DIO24H
	0015  PCI-DIO24H/CTR3
	0016  PCI-DIO48H/CTR15
	0017  PCI-DIO96H
	0018  PCI-CTR05
	0019  PCI-DAS1200/JR
	001a  PCI-DAS1001
	001b  PCI-DAS1002
	001c  PCI-DAS1602JR/16
	001d  PCI-DAS6402/16
	001e  PCI-DAS6402/12
	001f  PCI-DAS16/M1
	0020  PCI-DDA02/12
	0021  PCI-DDA04/12
	0022  PCI-DDA08/12
	0023  PCI-DDA02/16
	0024  PCI-DDA04/16
	0025  PCI-DDA08/16
	0026  PCI-DAC04/12-HS
	0027  PCI-DAC04/16-HS
	0028  PCI-DIO24
	0029  PCI-DAS08
	002c  PCI-INT32
	0033  PCI-DUAL-AC5
	0034  PCI-DAS-TC
	0035  PCI-DAS64/M1/16
	0036  PCI-DAS64/M2/16
	0037  PCI-DAS64/M3/16
	004c  PCI-DAS1000
	004d  PCI-QUAD04
	0052  PCI-DAS4020/12
	0053  PCIM-DDA06/16
	0054  PCI-DIO96
	005d  PCI-DAS6023
	005e  PCI-DAS6025
	005f  PCI-DAS6030
	0060  PCI-DAS6031
	0061  PCI-DAS6032
	0062  PCI-DAS6033
	0063  PCI-DAS6034
	0064  PCI-DAS6035
	0065  PCI-DAS6040
	0066  PCI-DAS6052
	0067  PCI-DAS6070
	0068  PCI-DAS6071
	006f  PCI-DAS6036
	0070  PCI-DAC6702
	0078  PCI-DAS6013
	0079  PCI-DAS6014
	0115  PCIe-DAS1602/16
1308  Jato Technologies Inc.
	0001  NetCelerator Adapter
1309  AB Semiconductor Ltd
130a  Mitsubishi Electric Microcomputer
130b  Colorgraphic Communications Corp
130c  Ambex Technologies, Inc
130d  Accelerix Inc
130e  Yamatake-Honeywell Co. Ltd
130f  Advanet Inc
1310  Gespac
1311  Videoserver, Inc
1312  Acuity Imaging, Inc
1313  Yaskawa Electric Co.
1315  Wavesat
1316  Teradyne Inc
1317  ADMtek
	0981  21x4x DEC-Tulip compatible 10/100 Ethernet
	0985  NC100 Network Everywhere Fast Ethernet 10/100
	1985  21x4x DEC-Tulip compatible 10/100 Ethernet
	2850  HSP MicroModem 56
	5120  ADM5120 OpenGate System-on-Chip
	8201  ADM8211 802.11b Wireless Interface
	8211  ADM8211 802.11b Wireless Interface
	9511  21x4x DEC-Tulip compatible 10/100 Ethernet
1318  Packet Engines Inc.
	0911  GNIC-II PCI Gigabit Ethernet [Hamachi]
1319  Fortemedia, Inc
	0801  Xwave QS3000A [FM801]
	0802  Xwave QS3000A [FM801 game port]
	1000  FM801 PCI Audio
	1001  FM801 PCI Joystick
131a  Finisar Corp.
131c  Nippon Electro-Sensory Devices Corp
131d  Sysmic, Inc.
131e  Xinex Networks Inc
131f  Siig Inc
	1000  CyberSerial (1-port) 16550
	1001  CyberSerial (1-port) 16650
	1002  CyberSerial (1-port) 16850
	1010  Duet 1S(16550)+1P
	1011  Duet 1S(16650)+1P
	1012  Duet 1S(16850)+1P
	1020  CyberParallel (1-port)
	1021  CyberParallel (2-port)
	1030  CyberSerial (2-port) 16550
	1031  CyberSerial (2-port) 16650
	1032  CyberSerial (2-port) 16850
	1034  Trio 2S(16550)+1P
	1035  Trio 2S(16650)+1P
	1036  Trio 2S(16850)+1P
	1050  CyberSerial (4-port) 16550
	1051  CyberSerial (4-port) 16650
	1052  CyberSerial (4-port) 16850
	2000  CyberSerial (1-port) 16550
	2001  CyberSerial (1-port) 16650
	2002  CyberSerial (1-port) 16850
	2010  Duet 1S(16550)+1P
	2011  Duet 1S(16650)+1P
	2012  Duet 1S(16850)+1P
	2020  CyberParallel (1-port)
	2021  CyberParallel (2-port)
	2030  CyberSerial (2-port) 16550
	2031  CyberSerial (2-port) 16650
	2032  CyberSerial (2-port) 16850
	2040  Trio 1S(16550)+2P
	2041  Trio 1S(16650)+2P
	2042  Trio 1S(16850)+2P
	2050  CyberSerial (4-port) 16550
	2051  CyberSerial (4-port) 16650
	2052  CyberSerial (4-port) 16850
	2060  Trio 2S(16550)+1P
	2061  Trio 2S(16650)+1P
	2062  Trio 2S(16850)+1P
	2081  CyberSerial (8-port) ST16654
1320  Crypto AG
1321  Arcobel Graphics BV
1322  MTT Co., Ltd
1323  Dome Inc
1324  Sphere Communications
1325  Salix Technologies, Inc
1326  Seachange international
1327  Voss scientific
1328  quadrant international
1329  Productivity Enhancement
132a  Microcom Inc.
132b  Broadband Technologies
132c  Micrel Inc
132d  Integrated Silicon Solution, Inc.
1330  MMC Networks
1331  RadiSys Corporation
	0030  ENP-2611
	8200  82600 Host Bridge
	8201  82600 IDE
	8202  82600 USB
	8210  82600 PCI Bridge
1332  Micro Memory
	5415  MM-5415CN PCI Memory Module with Battery Backup
	5425  MM-5425CN PCI 64/66 Memory Module with Battery Backup
	6140  MM-6140D
1334  Redcreek Communications, Inc
1335  Videomail, Inc
1337  Third Planet Publishing
1338  BT Electronics
133a  Vtel Corp
133b  Softcom Microsystems
133c  Holontech Corp
133d  SS Technologies
133e  Virtual Computer Corp
133f  SCM Microsystems
1340  Atalla Corp
1341  Kyoto Microcomputer Co
1342  Promax Systems Inc
1343  Phylon Communications Inc
1344  Micron Technology Inc
	5150  RealSSD P320h
	5151  RealSSD P320m
	5152  RealSSD P320s
	5153  RealSSD P325m
	5160  RealSSD P420h
	5161  RealSSD P420m
	5163  RealSSD P425m
	5180  9100 PRO NVMe SSD
	5181  9100 MAX NVMe SSD
1345  Arescom Inc
1347  Odetics
1349  Sumitomo Electric Industries, Ltd.
134a  DTC Technology Corp.
	0001  Domex 536
	0002  Domex DMX3194UP SCSI Adapter
134b  ARK Research Corp.
134c  Chori Joho System Co. Ltd
134d  PCTel Inc
	2189  HSP56 MicroModem
	2486  2304WT V.92 MDC Modem
	7890  HSP MicroModem 56
	7891  HSP MicroModem 56
	7892  HSP MicroModem 56
	7893  HSP MicroModem 56
	7894  HSP MicroModem 56
	7895  HSP MicroModem 56
	7896  HSP MicroModem 56
	7897  HSP MicroModem 56
134e  CSTI
134f  Algo System Co Ltd
1350  Systec Co. Ltd
1351  Sonix Inc
1353  Vierling Communication SAS
	0002  Proserver
	0003  PCI-FUT
	0004  PCI-S0
	0005  PCI-FUT-S0
1354  Dwave System Inc
1355  Kratos Analytical Ltd
1356  The Logical Co
1359  Prisa Networks
135a  Brain Boxes
	0a61  UC-324 [VELOCITY RS422/485]
135b  Giganet Inc
135c  Quatech Inc
	0010  QSC-100
	0020  DSC-100
	0030  DSC-200/300
	0040  QSC-200/300
	0050  ESC-100D
	0060  ESC-100M
	00f0  MPAC-100 Synchronous Serial Card (Zilog 85230)
	0170  QSCLP-100
	0180  DSCLP-100
	0190  SSCLP-100
	01a0  QSCLP-200/300
	01b0  DSCLP-200/300
	01c0  SSCLP-200/300
	0258  DSPSX-200/300
135d  ABB Network Partner AB
135e  Sealevel Systems Inc
	5101  Route 56.PCI - Multi-Protocol Serial Interface (Zilog Z16C32)
	7101  Single Port RS-232/422/485/530
	7201  Dual Port RS-232/422/485 Interface
	7202  Dual Port RS-232 Interface
	7401  Four Port RS-232 Interface
	7402  Four Port RS-422/485 Interface
	7801  Eight Port RS-232 Interface
	7804  Eight Port RS-232/422/485 Interface
	8001  8001 Digital I/O Adapter
135f  I-Data International A-S
1360  Meinberg Funkuhren
	0101  PCI32 DCF77 Radio Clock
	0102  PCI509 DCF77 Radio Clock
	0103  PCI510 DCF77 Radio Clock
	0104  PCI511 DCF77 Radio Clock
	0105  PEX511 DCF77 Radio Clock (PCI Express)
	0106  PZF180PEX High Precision DCF77 Radio Clock (PCI Express)
	0201  GPS167PCI GPS Receiver
	0202  GPS168PCI GPS Receiver
	0203  GPS169PCI GPS Receiver
	0204  GPS170PCI GPS Receiver
	0205  GPS170PEX GPS Receiver (PCI Express)
	0206  GPS180PEX GPS Receiver (PCI Express)
	0207  GLN180PEX GPS/GLONASS receiver (PCI Express)
	0208  GPS180AMC GPS Receiver (PCI Express / MicroTCA / AdvancedMC)
	0209  GNS181PEX GPS/Galileo/GLONASS/BEIDOU receiver (PCI Express)
	0301  TCR510PCI IRIG Timecode Reader
	0302  TCR167PCI IRIG Timecode Reader
	0303  TCR511PCI IRIG Timecode Reader
	0304  TCR511PEX IRIG Timecode Reader (PCI Express)
	0305  TCR170PEX IRIG Timecode Reader (PCI Express)
	0306  TCR180PEX IRIG Timecode Reader (PCI Express)
	0501  PTP270PEX PTP/IEEE1588 slave card (PCI Express)
	0601  FRC511PEX Free Running Clock (PCI Express)
1361  Soliton Systems K.K.
1362  Fujifacom Corporation
1363  Phoenix Technology Ltd
1364  ATM Communications Inc
1365  Hypercope GmbH
1366  Teijin Seiki Co. Ltd
1367  Hitachi Zosen Corporation
1368  Skyware Corporation
1369  Digigram
136a  High Soft Tech
	0004  HST Saphir VII mini PCI
	0007  HST Saphir III E MultiLink 4
	0008  HST Saphir III E MultiLink 8
	000a  HST Saphir III E MultiLink 2
136b  Kawasaki Steel Corporation
	ff01  KL5A72002 Motion JPEG
136c  Adtek System Science Co Ltd
136d  Gigalabs Inc
136f  Applied Magic Inc
1370  ATL Products
1371  CNet Technology Inc
	434e  GigaCard Network Adapter
1373  Silicon Vision Inc
1374  Silicom Ltd.
	0024  Silicom Dual port Giga Ethernet BGE Bypass Server Adapter
	0025  Silicom Quad port Giga Ethernet BGE Bypass Server Adapter
	0026  Silicom Dual port Fiber Giga Ethernet 546 Bypass Server Adapter
	0027  Silicom Dual port Fiber LX Giga Ethernet 546 Bypass Server Adapter
	0029  Silicom Dual port Copper Giga Ethernet 546GB Bypass Server Adapter
	002a  Silicom Dual port Fiber Giga Ethernet 546 TAP/Bypass Server Adapter
	002b  Silicom Dual port Copper Fast Ethernet 546 TAP/Bypass Server Adapter (PXE2TBI)
	002c  Silicom Quad port Copper Giga Ethernet 546GB Bypass Server Adapter (PXG4BPI)
	002d  Silicom Quad port Fiber-SX Giga Ethernet 546GB Bypass Server Adapter (PXG4BPFI)
	002e  Silicom Quad port Fiber-LX Giga Ethernet 546GB Bypass Server Adapter (PXG4BPFI-LX)
	002f  Silicom Dual port Fiber-SX Giga Ethernet 546GB Low profile Bypass Server Adapter (PXG2BPFIL)
	0030  Silicom Dual port Fiber-LX Giga Ethernet 546GB Low profile Bypass Server Adapter
	0031  Silicom Quad port Copper Giga Ethernet PCI-E Bypass Server Adapter
	0032  Silicom Dual port Copper Fast Ethernet 546 TAP/Bypass Server Adapter
	0034  Silicom Dual port Copper Giga Ethernet PCI-E BGE Bypass Server Adapter
	0035  Silicom Quad port Copper Giga Ethernet PCI-E BGE Bypass Server Adapter
	0036  Silicom Dual port Fiber Giga Ethernet PCI-E BGE Bypass Server Adapter
	0037  Silicom Dual port Copper Ethernet PCI-E Intel based Bypass Server Adapter
	0038  Silicom Quad port Copper Ethernet PCI-E Intel based Bypass Server Adapter
	0039  Silicom Dual port Fiber-SX Ethernet PCI-E Intel based Bypass Server Adapter
	003a  Silicom Dual port Fiber-LX Ethernet PCI-E Intel based Bypass Server Adapter
	003b  Silicom Dual port Fiber Ethernet PMC Intel based Bypass Server Adapter (PMCX2BPFI)
	003c  Silicom Dual port Copper Ethernet PCI-X BGE based Bypass Server Adapter (PXG2BPRB)
	003d  2-port Copper GBE Bypass with Caviume 1010 PCI-X
	003e  Silicom Dual port Fiber Giga Ethernet PCI-E 571 TAP/Bypass Server Adapter (PEG2TBFI)
	003f  Silicom Dual port Copper Giga Ethernet PCI-X 546 TAP/Bypass Server Adapter (PXG2TBI)
	0040  Silicom Quad port Fiber-SX Giga Ethernet 571 Bypass Server Adapter (PEG4BPFI)
	0042  4-port Copper GBE PMC-X Bypass
	0043  Silicom Quad port Fiber-SX Giga Ethernet 546 Bypass Server Adapter (PXG4BPFID)
	0045  Silicom 6 port Copper Giga Ethernet 546 Bypass Server Adapter (PXG6BPI)
	0046  4-port bypass PCI-E w disconnect low profile
	0047  Silicom Dual port Fiber-SX Giga Ethernet 571 Bypass Disconnect Server Adapter (PEG2BPFID)
	004a  Silicom Quad port Fiber-LX Giga Ethernet 571 Bypass Server Adapter (PEG4BPFI-LX)
	004d  Dual port Copper Giga Ethernet PCI-E Bypass Server Adapter
	0401  Gigabit Ethernet ExpressModule Bypass Server Adapter
	0420  Gigabit Ethernet ExpressModule Bypass Server Adapter
	0460  Gigabit Ethernet Express Module Bypass Server Adapter
	0461  Gigabit Ethernet ExpressModule Bypass Server Adapter
	0462  Gigabit Ethernet ExpressModule Bypass Server Adapter
	0470  Octal-port Copper Gigabit Ethernet Express Module Bypass Server Adapter
	0482  Dual-port Fiber (SR) 10 Gigabit Ethernet ExpressModule Bypass Server Adapter
	0483  Dual-port Fiber (LR) 10 Gigabit Ethernet ExpressModule Bypass Server Adapter
1375  Argosystems Inc
1376  LMC
1377  Electronic Equipment Production & Distribution GmbH
1378  Telemann Co. Ltd
1379  Asahi Kasei Microsystems Co Ltd
137a  Mark of the Unicorn Inc
	0001  PCI-324 Audiowire Interface
137b  PPT Vision
137c  Iwatsu Electric Co Ltd
137d  Dynachip Corporation
137e  Patriot Scientific Corporation
137f  Japan Satellite Systems Inc
1380  Sanritz Automation Co Ltd
1381  Brains Co. Ltd
1382  Marian - Electronic & Software
	0001  ARC88 audio recording card
	2008  Prodif 96 Pro sound system
	2048  Prodif Plus sound system
	2088  Marc 8 Midi sound system
	20c8  Marc A sound system
	4008  Marc 2 sound system
	4010  Marc 2 Pro sound system
	4048  Marc 4 MIDI sound system
	4088  Marc 4 Digi sound system
	4248  Marc X sound system
	4424  TRACE D4 Sound System
1383  Controlnet Inc
1384  Reality Simulation Systems Inc
1385  Netgear
	006b  WA301 802.11b Wireless PCI Adapter
	4100  MA301 802.11b Wireless PCI Adapter
	4601  WAG511 802.11a/b/g Dual Band Wireless PC Card
	620a  GA620 Gigabit Ethernet
	630a  GA630 Gigabit Ethernet
1386  Video Domain Technologies
1387  Systran Corp
1388  Hitachi Information Technology Co Ltd
1389  Applicom International
	0001  PCI1500PFB [Intelligent fieldbus adaptor]
138a  Fusion Micromedia Corp
	003d  VFS491 Validity Sensor
138b  Tokimec Inc
138c  Silicon Reality
138d  Future Techno Designs pte Ltd
138e  Basler GmbH
138f  Patapsco Designs Inc
1390  Concept Development Inc
1391  Development Concepts Inc
1392  Medialight Inc
1393  Moxa Technologies Co Ltd
	0001  UC7000 Serial
	1020  CP102 (2-port RS-232 PCI)
	1021  CP102UL (2-port RS-232 Universal PCI)
	1022  CP102U (2-port RS-232 Universal PCI)
	1023  CP-102UF
	1024  CP-102E (2-port RS-232 Smart PCI Express Serial Board)
	1025  CP-102EL (2-port RS-232 Smart PCI Express Serial Board)
	1040  Smartio C104H/PCI
	1041  CP104U (4-port RS-232 Universal PCI)
	1042  CP104JU (4-port RS-232 Universal PCI)
	1043  CP104EL (4-port RS-232 Smart PCI Express)
	1044  POS104UL (4-port RS-232 Universal PCI)
	1045  CP-104EL-A (4-port RS-232 PCI Express Serial Board)
	1080  CB108 (8-port RS-232 PC/104-plus Module)
	1140  CT-114 series
	1141  Industrio CP-114
	1142  CB114 (4-port RS-232/422/485 PC/104-plus Module)
	1143  CP-114UL (4-port RS-232/422/485 Smart Universal PCI Serial Board)
	1144  CP-114EL (4-port RS-232/422/485 Smart PCI Express Serial Board)
	1180  CP118U (8-port RS-232/422/485 Smart Universal PCI)
	1181  CP118EL (8-port RS-232/422/485 Smart PCI Express)
	1182  CP-118EL-A (8-port RS-232/422/485 PCI Express Serial Board)
	1320  CP132 (2-port RS-422/485 PCI)
	1321  CP132U (2-Port RS-422/485 Universal PCI)
	1322  CP-132EL (2-port RS-422/485 Smart PCI Express Serial Board)
	1340  CP134U (4-Port RS-422/485 Universal PCI)
	1341  CB134I (4-port RS-422/485 PC/104-plus Module)
	1380  CP138U (8-port RS-232/422/485 Smart Universal PCI)
	1680  Smartio C168H/PCI
	1681  CP-168U V2 Smart Serial Board (8-port RS-232)
	1682  CP168EL (8-port RS-232 Smart PCI Express)
	1683  CP-168EL-A (8-port RS-232 PCI Express Serial Board)
	2040  Intellio CP-204J
	2180  Intellio C218 Turbo PCI
	3200  Intellio C320 Turbo PCI
1394  Level One Communications
	0001  LXT1001 Gigabit Ethernet
1395  Ambicom Inc
1396  Cipher Systems Inc
1397  Cologne Chip Designs GmbH
	08b4  ISDN network Controller [HFC-4S]
	16b8  ISDN network Controller [HFC-8S]
	2bd0  ISDN network controller [HFC-PCI]
	30b1  ISDN network Controller [HFC-E1]
	b700  ISDN network controller PrimuX S0 [HFC-PCI]
	f001  GSM Network Controller [HFC-4GSM]
1398  Clarion co. Ltd
1399  Rios systems Co Ltd
139a  Alacritech Inc
	0001  Quad Port 10/100 Server Accelerator
	0003  Single Port 10/100 Server Accelerator
	0005  Single Port Gigabit Server Accelerator
139b  Mediasonic Multimedia Systems Ltd
139c  Quantum 3d Inc
139d  EPL limited
139e  Media4
139f  Aethra s.r.l.
13a0  Crystal Group Inc
13a1  Kawasaki Heavy Industries Ltd
13a2  Ositech Communications Inc
13a3  Hifn Inc.
	0005  7751 Security Processor
	0006  6500 Public Key Processor
	0007  7811 Security Processor
	0012  7951 Security Processor
	0014  78XX Security Processor
	0016  8065 Security Processor
	0017  8165 Security Processor
	0018  8154 Security Processor
	001d  7956 Security Processor
	001f  7855 Security Processor
	0020  7955 Security Processor
	0026  8155 Security Processor
	002e  9630 Compression Processor
	002f  9725 Compression and Security Processor
	0033  8201 Acceleration Processor
	0034  8202 Acceleration Processor
	0035  8203 Acceleration Processor
	0037  8204 Acceleration Processor
13a4  Rascom Inc
13a5  Audio Digital Imaging Inc
13a6  Videonics Inc
13a7  Teles AG
13a8  Exar Corp.
	0152  XR17C/D152 Dual PCI UART
	0154  XR17C154 Quad UART
	0158  XR17C158 Octal UART
	0252  XR17V252 Dual UART PCI controller
	0254  XR17V254 Quad UART PCI controller
	0258  XR17V258 Octal UART PCI controller
13a9  Siemens Medical Systems, Ultrasound Group
13aa  Broadband Networks Inc
13ab  Arcom Control Systems Ltd
13ac  Motion Media Technology Ltd
13ad  Nexus Inc
13ae  ALD Technology Ltd
13af  T.Sqware
13b0  Maxspeed Corp
13b1  Tamura corporation
13b2  Techno Chips Co. Ltd
13b3  Lanart Corporation
13b4  Wellbean Co Inc
13b5  ARM
13b6  Dlog GmbH
13b7  Logic Devices Inc
13b8  Nokia Telecommunications oy
13b9  Elecom Co Ltd
13ba  Oxford Instruments
13bb  Sanyo Technosound Co Ltd
13bc  Bitran Corporation
13bd  Sharp corporation
13be  Miroku Jyoho Service Co. Ltd
13bf  Sharewave Inc
13c0  Microgate Corporation
	0010  SyncLink Adapter v1
	0020  SyncLink SCC Adapter
	0030  SyncLink Multiport Adapter
	0070  SyncLink GT Adapter
	0080  SyncLink GT4 Adapter
	00a0  SyncLink GT2 Adapter
	0210  SyncLink Adapter v2
13c1  3ware Inc
	1000  5xxx/6xxx-series PATA-RAID
	1001  7xxx/8xxx-series PATA/SATA-RAID
	1002  9xxx-series SATA-RAID
	1003  9550SX SATA-II RAID PCI-X
	1004  9650SE SATA-II RAID PCIe
	1005  9690SA SAS/SATA-II RAID PCIe
	1010  9750 SAS2/SATA-II RAID PCIe
13c2  Technotrend Systemtechnik GmbH
	000e  Technotrend/Hauppauge DVB card rev2.3
	1019  TTechnoTrend-budget DVB S2-3200
13c3  Janz Computer AG
13c4  Phase Metrics
13c5  Alphi Technology Corp
13c6  Condor Engineering Inc
	0520  CEI-520 A429 Card
	0620  CEI-620 A429 Card
	0820  CEI-820 A429 Card
	0830  CEI-830 A429 Card
	1004  P-SER Multi-channel PMC to RS-485/422/232 adapter
13c7  Blue Chip Technology Ltd
	0adc  PCI-ADC
	0b10  PCI-PIO
	0d10  PCI-DIO
	524c  PCI-RLY
	5744  PCI-WDT
13c8  Apptech Inc
13c9  Eaton Corporation
13ca  Iomega Corporation
13cb  Yano Electric Co Ltd
13cc  Metheus Corporation
13cd  Compatible Systems Corporation
13ce  Cocom A/S
13cf  Studio Audio & Video Ltd
13d0  Techsan Electronics Co Ltd
	2103  B2C2 FlexCopII DVB chip / Technisat SkyStar2 DVB card
	2104  B2C2 FlexCopIII DVB chip / Technisat SkyStar2 DVB card (rev 01)
	2200  B2C2 FlexCopIII DVB chip / Technisat SkyStar2 DVB card
13d1  Abocom Systems Inc
	ab02  ADMtek Centaur-C rev 17 [D-Link DFE-680TX] CardBus Fast Ethernet Adapter
	ab03  21x4x DEC-Tulip compatible 10/100 Ethernet
	ab06  RTL8139 [FE2000VX] CardBus Fast Ethernet Attached Port Adapter
	ab08  21x4x DEC-Tulip compatible 10/100 Ethernet
13d2  Shark Multimedia Inc
13d4  Graphics Microsystems Inc
13d5  Media 100 Inc
13d6  K.I. Technology Co Ltd
13d7  Toshiba Engineering Corporation
13d8  Phobos corporation
13d9  Apex PC Solutions Inc
13da  Intresource Systems pte Ltd
13db  Janich & Klass Computertechnik GmbH
13dc  Netboost Corporation
13dd  Multimedia Bundle Inc
13de  ABB Robotics Products AB
13df  E-Tech Inc
	0001  PCI56RVP Modem
13e0  GVC Corporation
13e1  Silicom Multimedia Systems Inc
13e2  Dynamics Research Corporation
13e3  Nest Inc
13e4  Calculex Inc
13e5  Telesoft Design Ltd
13e6  Argosy research Inc
13e7  NAC Incorporated
13e8  Chip Express Corporation
13e9  Intraserver Technology Inc
13ea  Dallas Semiconductor
13eb  Hauppauge Computer Works Inc
13ec  Zydacron Inc
	000a  NPC-RC01 Remote control receiver
13ed  Raytheion E-Systems
13ee  Hayes Microcomputer Products Inc
13ef  Coppercom Inc
13f0  Sundance Technology Inc / IC Plus Corp
	0200  IC Plus IP100A Integrated 10/100 Ethernet MAC + PHY
	0201  ST201 Sundance Ethernet
	1021  TC902x Gigabit Ethernet
	1023  IP1000 Family Gigabit Ethernet
13f1  Oce' - Technologies B.V.
13f2  Ford Microelectronics Inc
13f3  Mcdata Corporation
13f4  Troika Networks, Inc.
	1401  Zentai Fibre Channel Adapter
13f5  Kansai Electric Co. Ltd
13f6  C-Media Electronics Inc
	0011  CMI8738
	0100  CM8338A
	0101  CM8338B
	0111  CMI8738/CMI8768 PCI Audio
	0211  CM8738
	5011  CM8888 [Oxygen Express]
	8788  CMI8788 [Oxygen HD Audio]
13f7  Wildfire Communications
13f8  Ad Lib Multimedia Inc
13f9  NTT Advanced Technology Corp.
13fa  Pentland Systems Ltd
13fb  Aydin Corp
13fc  Computer Peripherals International
13fd  Micro Science Inc
13fe  Advantech Co. Ltd
	1240  PCI-1240 4-channel stepper motor controller card
	1600  PCI-16xx series PCI multiport serial board (function 0)
	1603  PCI-1603 2-port isolated RS-232/current loop
	1604  PCI-1604 2-port RS-232
	16ff  PCI-16xx series PCI multiport serial board (function 1: RX/TX steering CPLD)
	1711  PCI-1711 16-channel data acquisition card 12-bit, 100kS/s
	1733  PCI-1733 32-channel isolated digital input card
	1752  PCI-1752
	1754  PCI-1754
	1756  PCI-1756
	c302  MIOe-3680 2-Port CAN-Bus MIOe Module with Isolation Protection
13ff  Silicon Spice Inc
1400  Artx Inc
	1401  9432 TX
1401  CR-Systems A/S
1402  Meilhaus Electronic GmbH
	0630  ME-630
	0940  ME-94
	0950  ME-95
	0960  ME-96
	1000  ME-1000
	100a  ME-1000
	100b  ME-1000
	1400  ME-1400
	140a  ME-1400A
	140b  ME-1400B
	140c  ME-1400C
	140d  ME-1400D
	140e  ME-1400E
	14ea  ME-1400EA
	14eb  ME-1400EB
	1604  ME-1600/4U
	1608  ME-1600/8U
	160c  ME-1600/12U
	160f  ME-1600/16U
	168f  ME-1600/16U8I
	4610  ME-4610
	4650  ME-4650
	4660  ME-4660
	4661  ME-4660I
	4662  ME-4660
	4663  ME-4660I
	4670  ME-4670
	4671  ME-4670I
	4672  ME-4670S
	4673  ME-4670IS
	4680  ME-4680
	4681  ME-4680I
	4682  ME-4680S
	4683  ME-4680IS
	6004  ME-6000/4
	6008  ME-6000/8
	600f  ME-6000/16
	6014  ME-6000I/4
	6018  ME-6000I/8
	601f  ME-6000I/16
	6034  ME-6000ISLE/4
	6038  ME-6000ISLE/8
	603f  ME-6000ISLE/16
	6044  ME-6000/4/DIO
	6048  ME-6000/8/DIO
	604f  ME-6000/16/DIO
	6054  ME-6000I/4/DIO
	6058  ME-6000I/8/DIO
	605f  ME-6000I/16/DIO
	6074  ME-6000ISLE/4/DIO
	6078  ME-6000ISLE/8/DIO
	607f  ME-6000ISLE/16/DIO
	6104  ME-6100/4
	6108  ME-6100/8
	610f  ME-6100/16
	6114  ME-6100I/4
	6118  ME-6100I/8
	611f  ME-6100I/16
	6134  ME-6100ISLE/4
	6138  ME-6100ISLE/8
	613f  ME-6100ISLE/16
	6144  ME-6100/4/DIO
	6148  ME-6100/8/DIO
	614f  ME-6100/16/DIO
	6154  ME-6100I/4/DIO
	6158  ME-6100I/8/DIO
	615f  ME-6100I/16/DIO
	6174  ME-6100ISLE/4/DIO
	6178  ME-6100ISLE/8/DIO
	617f  ME-6100ISLE/16/DIO
	6259  ME-6200I/9/DIO
	6359  ME-6300I/9/DIO
	810a  ME-8100A
	810b  ME-8100B
	820a  ME-8200A
	820b  ME-8200B
1403  Ascor Inc
1404  Fundamental Software Inc
1405  Excalibur Systems Inc
1406  Oce' Printing Systems GmbH
1407  Lava Computer mfg Inc
	0100  Lava Dual Serial
	0101  Lava Quatro A
	0102  Lava Quatro B
	0110  Lava DSerial-PCI Port A
	0111  Lava DSerial-PCI Port B
	0120  Quattro-PCI A
	0121  Quattro-PCI B
	0180  Lava Octo A
	0181  Lava Octo B
	0200  Lava Port Plus
	0201  Lava Quad A
	0202  Lava Quad B
	0220  Lava Quattro PCI Ports A/B
	0221  Lava Quattro PCI Ports C/D
	0400  Lava 8255-PIO-PCI
	0500  Lava Single Serial
	0520  Lava RS422-SS-PCI
	0600  Lava Port 650
	8000  Lava Parallel
	8001  Dual parallel port controller A
	8002  Lava Dual Parallel port A
	8003  Lava Dual Parallel port B
	8800  BOCA Research IOPPAR
1408  Aloka Co. Ltd
1409  Timedia Technology Co Ltd
	7168  PCI2S550 (Dual 16550 UART)
	7268  SUN1888 (Dual IEEE1284 parallel port)
140a  DSP Research Inc
140b  Abaco Systems, Inc.
140c  Elmic Systems Inc
140d  Matsushita Electric Works Ltd
140e  Goepel Electronic GmbH
140f  Salient Systems Corp
1410  Midas lab Inc
1411  Ikos Systems Inc
1412  VIA Technologies Inc.
	1712  ICE1712 [Envy24] PCI Multi-Channel I/O Controller
	1724  VT1720/24 [Envy24PT/HT] PCI Multi-Channel Audio Controller
1413  Addonics
1414  Microsoft Corporation
	0001  MN-120 (ADMtek Centaur-C based)
	0002  MN-130 (ADMtek Centaur-P based)
	5353  Hyper-V virtual VGA
	5801  XMA Decoder (Xenon)
	5802  SATA Controller - CdRom (Xenon)
	5803  SATA Controller - Disk (Xenon)
	5804  OHCI Controller 0 (Xenon)
	5805  EHCI Controller 0 (Xenon)
	5806  OHCI Controller 1 (Xenon)
	5807  EHCI Controller 1 (Xenon)
	580a  Fast Ethernet Adapter (Xenon)
	580b  Secure Flash Controller (Xenon)
	580d  System Management Controller (Xenon)
	5811  Xenos GPU (Xenon)
1415  Oxford Semiconductor Ltd
	8401  OX9162 Mode 1 (8-bit bus)
	8403  OX9162 Mode 0 (parallel port)
	9500  OX16PCI954 (Quad 16950 UART) function 0 (Disabled)
	9501  OX16PCI954 (Quad 16950 UART) function 0 (Uart)
	9505  OXuPCI952 (Dual 16C950 UART)
	950a  EXSYS EX-41092 Dual 16950 Serial adapter
	950b  OXCB950 Cardbus 16950 UART
	9510  OX16PCI954 (Quad 16950 UART) function 1 (Disabled)
	9511  OX16PCI954 (Quad 16950 UART) function 1 (8bit bus)
	9512  OX16PCI954 (Quad 16950 UART) function 1 (32bit bus)
	9513  OX16PCI954 (Quad 16950 UART) function 1 (parallel port)
	9521  OX16PCI952 (Dual 16950 UART)
	9523  OX16PCI952 Integrated Parallel Port
	c158  OXPCIe952 Dual 16C950 UART
	c308  EX-44016 16-port serial
1416  Multiwave Innovation pte Ltd
1417  Convergenet Technologies Inc
1418  Kyushu electronics systems Inc
1419  Excel Switching Corp
141a  Apache Micro Peripherals Inc
141b  Zoom Telephonics Inc
141d  Digitan Systems Inc
141e  Fanuc Ltd
141f  Visiontech Ltd
1420  Psion Dacom plc
	8002  Gold Card NetGlobal 56k+10/100Mb CardBus (Ethernet part)
	8003  Gold Card NetGlobal 56k+10/100Mb CardBus (Modem part)
1421  Ads Technologies Inc
1422  Ygrec Systems Co Ltd
1423  Custom Technology Corp.
1424  Videoserver Connections
1425  Chelsio Communications Inc
	000b  T210 Protocol Engine
	000c  T204 Protocol Engine
	0022  10GbE Ethernet Adapter
	0030  T310 10GbE Single Port Adapter
	0031  T320 10GbE Dual Port Adapter
	0032  T302 1GbE Dual Port Adapter
	0033  T304 1GbE Quad Port Adapter
	0034  B320 10GbE Dual Port Adapter
	0035  S310-CR 10GbE Single Port Adapter
	0036  S320-LP-CR 10GbE Dual Port Adapter
	0037  N320-G2-CR 10GbE Dual Port Adapter
	4001  T420-CR Unified Wire Ethernet Controller
	4002  T422-CR Unified Wire Ethernet Controller
	4003  T440-CR Unified Wire Ethernet Controller
	4004  T420-BCH Unified Wire Ethernet Controller
	4005  T440-BCH Unified Wire Ethernet Controller
	4006  T440-CH Unified Wire Ethernet Controller
	4007  T420-SO Unified Wire Ethernet Controller
	4008  T420-CX Unified Wire Ethernet Controller
	4009  T420-BT Unified Wire Ethernet Controller
	400a  T404-BT Unified Wire Ethernet Controller
	400b  B420-SR Unified Wire Ethernet Controller
	400c  B404-BT Unified Wire Ethernet Controller
	400d  T480 Unified Wire Ethernet Controller
	400e  T440-LP-CR Unified Wire Ethernet Controller
	400f  T440 [Amsterdam] Unified Wire Ethernet Controller
	4080  T480-4080 T480 Unified Wire Ethernet Controller
	4081  T440F-4081 T440-FCoE Unified Wire Ethernet Controller
	4082  T420-4082  Unified Wire Ethernet Controller
	4083  T420X-4083 Unified Wire Ethernet Controller
	4084  T440-4084 Unified Wire Ethernet Controller
	4085  T420-4085 SFP+ Unified Wire Ethernet Controller
	4086  T440-4086 10Gbase-T Unified Wire Ethernet Controller
	4087  T440T-4087 Unified Wire Ethernet Controller
	4088  T440-4088 Unified Wire Ethernet Controller
	4401  T420-CR Unified Wire Ethernet Controller
	4402  T422-CR Unified Wire Ethernet Controller
	4403  T440-CR Unified Wire Ethernet Controller
	4404  T420-BCH Unified Wire Ethernet Controller
	4405  T440-BCH Unified Wire Ethernet Controller
	4406  T440-CH Unified Wire Ethernet Controller
	4407  T420-SO Unified Wire Ethernet Controller
	4408  T420-CX Unified Wire Ethernet Controller
	4409  T420-BT Unified Wire Ethernet Controller
	440a  T404-BT Unified Wire Ethernet Controller
	440b  B420-SR Unified Wire Ethernet Controller
	440c  B404-BT Unified Wire Ethernet Controller
	440d  T480 Unified Wire Ethernet Controller
	440e  T440-LP-CR Unified Wire Ethernet Controller
	440f  T440 [Amsterdam] Unified Wire Ethernet Controller
	4480  T480-4080 T480 Unified Wire Ethernet Controller
	4481  T440F-4081 T440-FCoE Unified Wire Ethernet Controller
	4482  T420-4082  Unified Wire Ethernet Controller
	4483  T420X-4083 Unified Wire Ethernet Controller
	4484  T440-4084 Unified Wire Ethernet Controller
	4485  T420-4085 SFP+ Unified Wire Ethernet Controller
	4486  T440-4086 10Gbase-T Unified Wire Ethernet Controller
	4487  T440T-4087 Unified Wire Ethernet Controller
	4488  T440-4088 Unified Wire Ethernet Controller
	4501  T420-CR Unified Wire Storage Controller
	4502  T422-CR Unified Wire Storage Controller
	4503  T440-CR Unified Wire Storage Controller
	4504  T420-BCH Unified Wire Storage Controller
	4505  T440-BCH Unified Wire Storage Controller
	4506  T440-CH Unified Wire Storage Controller
	4507  T420-SO Unified Wire Storage Controller
	4508  T420-CX Unified Wire Storage Controller
	4509  T420-BT Unified Wire Storage Controller
	450a  T404-BT Unified Wire Storage Controller
	450b  B420-SR Unified Wire Storage Controller
	450c  B404-BT Unified Wire Storage Controller
	450d  T480 Unified Wire Storage Controller
	450e  T440-LP-CR Unified Wire Storage Controller
	450f  T440 [Amsterdam] Unified Wire Storage Controller
	4580  T480-4080 T480 Unified Wire Storage Controller
	4581  T440F-4081 T440-FCoE Unified Wire Storage Controller
	4582  T420-4082  Unified Wire Storage Controller
	4583  T420X-4083 Unified Wire Storage Controller
	4584  T440-4084 Unified Wire Storage Controller
	4585  T420-4085 SFP+ Unified Wire Storage Controller
	4586  T440-4086 10Gbase-T Unified Wire Storage Controller
	4587  T440T-4087 Unified Wire Storage Controller
	4588  T440-4088 Unified Wire Storage Controller
	4601  T420-CR Unified Wire Storage Controller
	4602  T422-CR Unified Wire Storage Controller
	4603  T440-CR Unified Wire Storage Controller
	4604  T420-BCH Unified Wire Storage Controller
	4605  T440-BCH Unified Wire Storage Controller
	4606  T440-CH Unified Wire Storage Controller
	4607  T420-SO Unified Wire Storage Controller
	4608  T420-CX Unified Wire Storage Controller
	4609  T420-BT Unified Wire Storage Controller
	460a  T404-BT Unified Wire Storage Controller
	460b  B420-SR Unified Wire Storage Controller
	460c  B404-BT Unified Wire Storage Controller
	460d  T480 Unified Wire Storage Controller
	460e  T440-LP-CR Unified Wire Storage Controller
	460f  T440 [Amsterdam] Unified Wire Storage Controller
	4680  T480-4080 T480 Unified Wire Storage Controller
	4681  T440F-4081 T440-FCoE Unified Wire Storage Controller
	4682  T420-4082  Unified Wire Storage Controller
	4683  T420X-4083 Unified Wire Storage Controller
	4684  T440-4084 Unified Wire Storage Controller
	4685  T420-4085 SFP+ Unified Wire Storage Controller
	4686  T440-4086 10Gbase-T Unified Wire Storage Controller
	4687  T440T-4087 Unified Wire Storage Controller
	4688  T440-4088 Unified Wire Storage Controller
	4701  T420-CR Unified Wire Ethernet Controller
	4702  T422-CR Unified Wire Ethernet Controller
	4703  T440-CR Unified Wire Ethernet Controller
	4704  T420-BCH Unified Wire Ethernet Controller
	4705  T440-BCH Unified Wire Ethernet Controller
	4706  T440-CH Unified Wire Ethernet Controller
	4707  T420-SO Unified Wire Ethernet Controller
	4708  T420-CX Unified Wire Ethernet Controller
	4709  T420-BT Unified Wire Ethernet Controller
	470a  T404-BT Unified Wire Ethernet Controller
	470b  B420-SR Unified Wire Ethernet Controller
	470c  B404-BT Unified Wire Ethernet Controller
	470d  T480 Unified Wire Ethernet Controller
	470e  T440-LP-CR Unified Wire Ethernet Controller
	470f  T440 [Amsterdam] Unified Wire Ethernet Controller
	4780  T480-4080 T480 Unified Wire Ethernet Controller
	4781  T440F-4081 T440-FCoE Unified Wire Ethernet Controller
	4782  T420-4082  Unified Wire Ethernet Controller
	4783  T420X-4083 Unified Wire Ethernet Controller
	4784  T440-4084 Unified Wire Ethernet Controller
	4785  T420-4085 SFP+ Unified Wire Ethernet Controller
	4786  T440-4086 10Gbase-T Unified Wire Ethernet Controller
	4787  T440T-4087 Unified Wire Ethernet Controller
	4788  T440-4088 Unified Wire Ethernet Controller
	4801  T420-CR Unified Wire Ethernet Controller [VF]
	4802  T422-CR Unified Wire Ethernet Controller [VF]
	4803  T440-CR Unified Wire Ethernet Controller [VF]
	4804  T420-BCH Unified Wire Ethernet Controller [VF]
	4805  T440-BCH Unified Wire Ethernet Controller [VF]
	4806  T440-CH Unified Wire Ethernet Controller [VF]
	4807  T420-SO Unified Wire Ethernet Controller [VF]
	4808  T420-CX Unified Wire Ethernet Controller [VF]
	4809  T420-BT Unified Wire Ethernet Controller [VF]
	480a  T404-BT Unified Wire Ethernet Controller [VF]
	480b  B420-SR Unified Wire Ethernet Controller [VF]
	480c  B404-BT Unified Wire Ethernet Controller [VF]
	480d  T480 Unified Wire Ethernet Controller [VF]
	480e  T440-LP-CR Unified Wire Ethernet Controller [VF]
	480f  T440 [Amsterdam] Unified Wire Ethernet Controller [VF]
	4880  T480-4080 T480 Unified Wire Ethernet Controller [VF]
	4881  T440F-4081 T440-FCoE Unified Wire Ethernet Controller [VF]
	4882  T420-4082 Unified Wire Ethernet Controller [VF]
	4883  T420X-4083 Unified Wire Ethernet Controller [VF]
	4884  T440-4084 Unified Wire Ethernet Controller [VF]
	4885  T420-4085 SFP+ Unified Wire Ethernet Controller [VF]
	4886  T440-4086 10Gbase-T Unified Wire Ethernet Controller [VF]
	4887  T440T-4087 Unified Wire Ethernet Controller [VF]
	4888  T440-4088 Unified Wire Ethernet Controller [VF]
	5001  T520-CR Unified Wire Ethernet Controller
	5002  T522-CR Unified Wire Ethernet Controller
	5003  T540-CR Unified Wire Ethernet Controller
	5004  T520-BCH Unified Wire Ethernet Controller
	5005  T540-BCH Unified Wire Ethernet Controller
	5006  T540-CH Unified Wire Ethernet Controller
	5007  T520-SO Unified Wire Ethernet Controller
	5008  T520-CX Unified Wire Ethernet Controller
	5009  T520-BT Unified Wire Ethernet Controller
	500a  T504-BT Unified Wire Ethernet Controller
	500b  B520-SR Unified Wire Ethernet Controller
	500c  B504-BT Unified Wire Ethernet Controller
	500d  T580-CR Unified Wire Ethernet Controller
	500e  T540-LP-CR Unified Wire Ethernet Controller
	500f  T540 [Amsterdam] Unified Wire Ethernet Controller
	5010  T580-LP-CR Unified Wire Ethernet Controller
	5011  T520-LL-CR Unified Wire Ethernet Controller
	5012  T560-CR Unified Wire Ethernet Controller
	5013  T580-CHR Unified Wire Ethernet Controller
	5014  T580-SO-CR Unified Wire Ethernet Controller
	5015  T502-BT Unified Wire Ethernet Controller
	5016  T580-OCP-SO Unified Wire Ethernet Controller
	5017  T520-OCP-SO Unified Wire Ethernet Controller
	5018  T540-BT Unified Wire Ethernet Controller
	5080  T540-5080 Unified Wire Ethernet Controller
	5081  T540-5081 Unified Wire Ethernet Controller
	5082  T504-5082 Unified Wire Ethernet Controller
	5083  T540-5083 Unified Wire Ethernet Controller
	5084  T540-5084 Unified Wire Ethernet Controller
	5085  T580-5085 Unified Wire Ethernet Controller
	5086  T580-5086 Unified Wire Ethernet Controller
	5087  T580-5087 Unified Wire Ethernet Controller
	5088  T570-5088 Unified Wire Ethernet Controller
	5089  T520-5089 Unified Wire Ethernet Controller
	5090  T540-5090 Unified Wire Ethernet Controller
	5091  T522-5091 Unified Wire Ethernet Controller
	5092  T520-5092 Unified Wire Ethernet Controller
	5093  T580-5093 Unified Wire Ethernet Controller
	5094  T540-5094 Unified Wire Ethernet Controller
	5095  T540-5095 Unified Wire Ethernet Controller
	5096  T580-5096 Unified Wire Ethernet Controller
	5097  T520-5097 Unified Wire Ethernet Controller
	5098  T580-5098 Unified Wire Ethernet Controller
	5099  T580-5099 Unified Wire Ethernet Controller
	509a  T520-509A Unified Wire Ethernet Controller
	509b  T540-509B Unified Wire Ethernet Controller
	509c  T520-509C Unified Wire Ethernet Controller
	509d  T540-509D Unified Wire Ethernet Controller
	509e  T520-509E Unified Wire Ethernet Controller
	509f  T540-509F Unified Wire Ethernet Controller
	50a0  T540-50A0 Unified Wire Ethernet Controller
	50a1  T540-50A1 Unified Wire Ethernet Controller
	50a2  T580-50A2 Unified Wire Ethernet Controller
	50a3  T580-50A3 Unified Wire Ethernet Controller
	50a4  T540-50A4 Unified Wire Ethernet Controller
	50a5  T522-50A5 Unified Wire Ethernet Controller
	50a6  T522-50A6 Unified Wire Ethernet Controller
	50a7  T580-50A7 Unified Wire Ethernet Controller
	50a8  T580-50A8 Unified Wire Ethernet Controller
	50a9  T580-50A9 Unified Wire Ethernet Controller
	50aa  T580-50AA Unified Wire Ethernet Controller
	50ab  T520-50AB Unified Wire Ethernet Controller
	5401  T520-CR Unified Wire Ethernet Controller
	5402  T522-CR Unified Wire Ethernet Controller
	5403  T540-CR Unified Wire Ethernet Controller
	5404  T520-BCH Unified Wire Ethernet Controller
	5405  T540-BCH Unified Wire Ethernet Controller
	5406  T540-CH Unified Wire Ethernet Controller
	5407  T520-SO Unified Wire Ethernet Controller
	5408  T520-CX Unified Wire Ethernet Controller
	5409  T520-BT Unified Wire Ethernet Controller
	540a  T504-BT Unified Wire Ethernet Controller
	540b  B520-SR Unified Wire Ethernet Controller
	540c  B504-BT Unified Wire Ethernet Controller
	540d  T580-CR Unified Wire Ethernet Controller
	540e  T540-LP-CR Unified Wire Ethernet Controller
	540f  T540 [Amsterdam] Unified Wire Ethernet Controller
	5410  T580-LP-CR Unified Wire Ethernet Controller
	5411  T520-LL-CR Unified Wire Ethernet Controller
	5412  T560-CR Unified Wire Ethernet Controller
	5413  T580-CHR Unified Wire Ethernet Controller
	5414  T580-SO-CR Unified Wire Ethernet Controller
	5415  T502-BT Unified Wire Ethernet Controller
	5416  T580-OCP-SO Unified Wire Ethernet Controller
	5417  T520-OCP-SO Unified Wire Ethernet Controller
	5418  T540-BT Unified Wire Ethernet Controller
	5480  T540-5080 Unified Wire Ethernet Controller
	5481  T540-5081 Unified Wire Ethernet Controller
	5482  T504-5082 Unified Wire Ethernet Controller
	5483  T540-5083 Unified Wire Ethernet Controller
	5484  T540-5084 Unified Wire Ethernet Controller
	5485  T580-5085 Unified Wire Ethernet Controller
	5486  T580-5086 Unified Wire Ethernet Controller
	5487  T580-5087 Unified Wire Ethernet Controller
	5488  T570-5088 Unified Wire Ethernet Controller
	5489  T520-5089 Unified Wire Ethernet Controller
	5490  T540-5090 Unified Wire Ethernet Controller
	5491  T522-5091 Unified Wire Ethernet Controller
	5492  T520-5092 Unified Wire Ethernet Controller
	5493  T580-5093 Unified Wire Ethernet Controller
	5494  T540-5094 Unified Wire Ethernet Controller
	5495  T540-5095 Unified Wire Ethernet Controller
	5496  T580-5096 Unified Wire Ethernet Controller
	5497  T520-5097 Unified Wire Ethernet Controller
	5498  T580-5098 Unified Wire Ethernet Controller
	5499  T580-5099 Unified Wire Ethernet Controller
	549a  T520-509A Unified Wire Ethernet Controller
	549b  T540-509B Unified Wire Ethernet Controller
	549c  T520-509C Unified Wire Ethernet Controller
	549d  T540-509D Unified Wire Ethernet Controller
	549e  T520-509E Unified Wire Ethernet Controller
	549f  T540-509F Unified Wire Ethernet Controller
	54a0  T540-50A0 Unified Wire Ethernet Controller
	54a1  T540-50A1 Unified Wire Ethernet Controller
	54a2  T580-50A2 Unified Wire Ethernet Controller
	54a3  T580-50A3 Unified Wire Ethernet Controller
	54a4  T540-50A4 Unified Wire Ethernet Controller
	54a5  T522-50A5 Unified Wire Ethernet Controller
	54a6  T522-50A6 Unified Wire Ethernet Controller
	54a7  T580-50A7 Unified Wire Ethernet Controller
	54a8  T580-50A8 Unified Wire Ethernet Controller
	54a9  T580-50A9 Unified Wire Ethernet Controller
	54aa  T580-50AA Unified Wire Ethernet Controller
	54ab  T520-50AB Unified Wire Ethernet Controller
	5501  T520-CR Unified Wire Storage Controller
	5502  T522-CR Unified Wire Storage Controller
	5503  T540-CR Unified Wire Storage Controller
	5504  T520-BCH Unified Wire Storage Controller
	5505  T540-BCH Unified Wire Storage Controller
	5506  T540-CH Unified Wire Storage Controller
	5507  T520-SO Unified Wire Storage Controller
	5508  T520-CX Unified Wire Storage Controller
	5509  T520-BT Unified Wire Storage Controller
	550a  T504-BT Unified Wire Storage Controller
	550b  B520-SR Unified Wire Storage Controller
	550c  B504-BT Unified Wire Storage Controller
	550d  T580-CR Unified Wire Storage Controller
	550e  T540-LP-CR Unified Wire Storage Controller
	550f  T540 [Amsterdam] Unified Wire Storage Controller
	5510  T580-LP-CR Unified Wire Storage Controller
	5511  T520-LL-CR Unified Wire Storage Controller
	5512  T560-CR Unified Wire Storage Controller
	5513  T580-CHR Unified Wire Storage Controller
	5514  T580-SO-CR Unified Wire Storage Controller
	5515  T502-BT Unified Wire Storage Controller
	5516  T580-OCP-SO Unified Wire Storage Controller
	5517  T520-OCP-SO Unified Wire Storage Controller
	5518  T540-BT Unified Wire Storage Controller
	5580  T540-5080 Unified Wire Storage Controller
	5581  T540-5081 Unified Wire Storage Controller
	5582  T504-5082 Unified Wire Storage Controller
	5583  T540-5083 Unified Wire Storage Controller
	5584  T540-5084 Unified Wire Storage Controller
	5585  T580-5085 Unified Wire Storage Controller
	5586  T580-5086 Unified Wire Storage Controller
	5587  T580-5087 Unified Wire Storage Controller
	5588  T570-5088 Unified Wire Storage Controller
	5589  T520-5089 Unified Wire Storage Controller
	5590  T540-5090 Unified Wire Storage Controller
	5591  T522-5091 Unified Wire Storage Controller
	5592  T520-5092 Unified Wire Storage Controller
	5593  T580-5093 Unified Wire Storage Controller
	5594  T540-5094 Unified Wire Storage Controller
	5595  T540-5095 Unified Wire Storage Controller
	5596  T580-5096 Unified Wire Storage Controller
	5597  T520-5097 Unified Wire Storage Controller
	5598  T580-5098 Unified Wire Storage Controller
	5599  T580-5099 Unified Wire Storage Controller
	559a  T520-509A Unified Wire Storage Controller
	559b  T540-509B Unified Wire Storage Controller
	559c  T520-509C Unified Wire Storage Controller
	559d  T540-509D Unified Wire Storage Controller
	559e  T520-509E Unified Wire Storage Controller
	559f  T540-509F Unified Wire Storage Controller
	55a0  T540-50A0 Unified Wire Storage Controller
	55a1  T540-50A1 Unified Wire Storage Controller
	55a2  T580-50A2 Unified Wire Storage Controller
	55a3  T580-50A3 Unified Wire Storage Controller
	55a4  T540-50A4 Unified Wire Storage Controller
	55a5  T522-50A5 Unified Wire Storage Controller
	55a6  T522-50A6 Unified Wire Storage Controller
	55a7  T580-50A7 Unified Wire Storage Controller
	55a8  T580-50A8 Unified Wire Storage Controller
	55a9  T580-50A9 Unified Wire Storage Controller
	5601  T520-CR Unified Wire Storage Controller
	5602  T522-CR Unified Wire Storage Controller
	5603  T540-CR Unified Wire Storage Controller
	5604  T520-BCH Unified Wire Storage Controller
	5605  T540-BCH Unified Wire Storage Controller
	5606  T540-CH Unified Wire Storage Controller
	5607  T520-SO Unified Wire Storage Controller
	5608  T520-CX Unified Wire Storage Controller
	5609  T520-BT Unified Wire Storage Controller
	560a  T504-BT Unified Wire Storage Controller
	560b  B520-SR Unified Wire Storage Controller
	560c  B504-BT Unified Wire Storage Controller
	560d  T580-CR Unified Wire Storage Controller
	560e  T540-LP-CR Unified Wire Storage Controller
	560f  T540 [Amsterdam] Unified Wire Storage Controller
	5610  T580-LP-CR Unified Wire Storage Controller
	5611  T520-LL-CR Unified Wire Storage Controller
	5612  T560-CR Unified Wire Storage Controller
	5613  T580-CHR Unified Wire Storage Controller
	5614  T580-SO-CR Unified Wire Storage Controller
	5615  T502-BT Unified Wire Storage Controller
	5616  T580-OCP-SO Unified Wire Storage Controller
	5617  T520-OCP-SO Unified Wire Storage Controller
	5618  T540-BT Unified Wire Storage Controller
	5680  T540-5080 Unified Wire Storage Controller
	5681  T540-5081 Unified Wire Storage Controller
	5682  T504-5082 Unified Wire Storage Controller
	5683  T540-5083 Unified Wire Storage Controller
	5684  T540-5084 Unified Wire Storage Controller
	5685  T580-5085 Unified Wire Storage Controller
	5686  T580-5086 Unified Wire Storage Controller
	5687  T580-5087 Unified Wire Storage Controller
	5688  T570-5088 Unified Wire Storage Controller
	5689  T520-5089 Unified Wire Storage Controller
	5690  T540-5090 Unified Wire Storage Controller
	5691  T522-5091 Unified Wire Storage Controller
	5692  T520-5092 Unified Wire Storage Controller
	5693  T580-5093 Unified Wire Storage Controller
	5694  T540-5094 Unified Wire Storage Controller
	5695  T540-5095 Unified Wire Storage Controller
	5696  T580-5096 Unified Wire Storage Controller
	5697  T520-5097 Unified Wire Storage Controller
	5698  T580-5098 Unified Wire Storage Controller
	5699  T580-5099 Unified Wire Storage Controller
	569a  T520-509A Unified Wire Storage Controller
	569b  T540-509B Unified Wire Storage Controller
	569c  T520-509C Unified Wire Storage Controller
	569d  T540-509D Unified Wire Storage Controller
	569e  T520-509E Unified Wire Storage Controller
	569f  T540-509F Unified Wire Storage Controller
	56a0  T540-50A0 Unified Wire Storage Controller
	56a1  T540-50A1 Unified Wire Storage Controller
	56a2  T580-50A2 Unified Wire Storage Controller
	56a3  T580-50A3 Unified Wire Storage Controller
	56a4  T540-50A4 Unified Wire Storage Controller
	56a5  T522-50A5 Unified Wire Storage Controller
	56a6  T522-50A6 Unified Wire Storage Controller
	56a7  T580-50A7 Unified Wire Storage Controller
	56a8  T580-50A8 Unified Wire Storage Controller
	56a9  T580-50A9 Unified Wire Storage Controller
	56aa  T580-50AA Unified Wire Storage Controller
	56ab  T520-50AB Unified Wire Storage Controller
	5701  T520-CR Unified Wire Ethernet Controller
	5702  T522-CR Unified Wire Ethernet Controller
	5703  T540-CR Unified Wire Ethernet Controller
	5704  T520-BCH Unified Wire Ethernet Controller
	5705  T540-BCH Unified Wire Ethernet Controller
	5706  T540-CH Unified Wire Ethernet Controller
	5707  T520-SO Unified Wire Ethernet Controller
	5708  T520-CX Unified Wire Ethernet Controller
	5709  T520-BT Unified Wire Ethernet Controller
	570a  T504-BT Unified Wire Ethernet Controller
	570b  B520-SR Unified Wire Ethernet Controller
	570c  B504-BT Unified Wire Ethernet Controller
	570d  T580-CR Unified Wire Ethernet Controller
	570e  T540-LP-CR Unified Wire Ethernet Controller
	570f  T540 [Amsterdam] Unified Wire Ethernet Controller
	5710  T580-LP-CR Unified Wire Ethernet Controller
	5711  T520-LL-CR Unified Wire Ethernet Controller
	5712  T560-CR Unified Wire Ethernet Controller
	5713  T580-CR Unified Wire Ethernet Controller
	5714  T580-SO-CR Unified Wire Ethernet Controller
	5715  T502-BT Unified Wire Ethernet Controller
	5780  T540-5080 Unified Wire Ethernet Controller
	5781  T540-5081 Unified Wire Ethernet Controller
	5782  T504-5082 Unified Wire Ethernet Controller
	5783  T540-5083 Unified Wire Ethernet Controller
	5784  T580-5084 Unified Wire Ethernet Controller
	5785  T580-5085 Unified Wire Ethernet Controller
	5786  T580-5086 Unified Wire Ethernet Controller
	5787  T580-5087 Unified Wire Ethernet Controller
	5788  T570-5088 Unified Wire Ethernet Controller
	5789  T520-5089 Unified Wire Ethernet Controller
	5790  T540-5090 Unified Wire Ethernet Controller
	5791  T522-5091 Unified Wire Ethernet Controller
	5792  T520-5092 Unified Wire Ethernet Controller
	5793  T580-5093 Unified Wire Ethernet Controller
	5794  T540-5094 Unified Wire Ethernet Controller
	5795  T540-5095 Unified Wire Ethernet Controller
	5796  T580-5096 Unified Wire Ethernet Controller
	5797  T520-5097 Unified Wire Ethernet Controller
	5801  T520-CR Unified Wire Ethernet Controller [VF]
	5802  T522-CR Unified Wire Ethernet Controller [VF]
	5803  T540-CR Unified Wire Ethernet Controller [VF]
	5804  T520-BCH Unified Wire Ethernet Controller [VF]
	5805  T540-BCH Unified Wire Ethernet Controller [VF]
	5806  T540-CH Unified Wire Ethernet Controller [VF]
	5807  T520-SO Unified Wire Ethernet Controller [VF]
	5808  T520-CX Unified Wire Ethernet Controller [VF]
	5809  T520-BT Unified Wire Ethernet Controller [VF]
	580a  T504-BT Unified Wire Ethernet Controller [VF]
	580b  B520-SR Unified Wire Ethernet Controller [VF]
	580c  B504-BT Unified Wire Ethernet Controller [VF]
	580d  T580-CR Unified Wire Ethernet Controller [VF]
	580e  T540-LP-CR Unified Wire Ethernet Controller [VF]
	580f  T540 [Amsterdam] Unified Wire Ethernet Controller [VF]
	5810  T580-LP-CR Unified Wire Ethernet Controller [VF]
	5811  T520-LL-CR Unified Wire Ethernet Controller [VF]
	5812  T560-CR Unified Wire Ethernet Controller [VF]
	5813  T580-CHR Unified Wire Ethernet Controller [VF]
	5814  T580-SO-CR Unified Wire Ethernet Controller [VF]
	5815  T502-BT Unified Wire Ethernet Controller [VF]
	5816  T580-OCP-SO Unified Wire Ethernet Controller [VF]
	5817  T520-OCP-SO Unified Wire Ethernet Controller [VF]
	5818  T540-BT Unified Wire Ethernet Controller [VF]
	5880  T540-5080 Unified Wire Ethernet Controller [VF]
	5881  T540-5081 Unified Wire Ethernet Controller [VF]
	5882  T504-5082 Unified Wire Ethernet Controller [VF]
	5883  T540-5083 Unified Wire Ethernet Controller [VF]
	5884  T540-5084 Unified Wire Ethernet Controller [VF]
	5885  T580-5085 Unified Wire Ethernet Controller [VF]
	5886  T580-5086 Unified Wire Ethernet Controller [VF]
	5887  T580-5087 Unified Wire Ethernet Controller [VF]
	5888  T570-5088 Unified Wire Ethernet Controller [VF]
	5889  T520-5089 Unified Wire Ethernet Controller [VF]
	5890  T540-5090 Unified Wire Ethernet Controller [VF]
	5891  T522-5091 Unified Wire Ethernet Controller [VF]
	5892  T520-5092 Unified Wire Ethernet Controller [VF]
	5893  T580-5093 Unified Wire Ethernet Controller [VF]
	5894  T540-5094 Unified Wire Ethernet Controller [VF]
	5895  T540-5095 Unified Wire Ethernet Controller [VF]
	5896  T580-5096 Unified Wire Ethernet Controller [VF]
	5897  T520-5097 Unified Wire Ethernet Controller [VF]
	5898  T580-5098 Unified Wire Ethernet Controller [VF]
	5899  T580-5099 Unified Wire Ethernet Controller [VF]
	589a  T520-509A Unified Wire Ethernet Controller [VF]
	589b  T540-509B Unified Wire Ethernet Controller [VF]
	589c  T520-509C Unified Wire Ethernet Controller [VF]
	589d  T540-509D Unified Wire Ethernet Controller [VF]
	589e  T520-509E Unified Wire Ethernet Controller [VF]
	589f  T540-509F Unified Wire Ethernet Controller [VF]
	58a0  T540-50A0 Unified Wire Ethernet Controller [VF]
	58a1  T540-50A1 Unified Wire Ethernet Controller [VF]
	58a2  T580-50A2 Unified Wire Ethernet Controller [VF]
	58a3  T580-50A3 Unified Wire Ethernet Controller [VF]
	58a4  T540-50A4 Unified Wire Ethernet Controller [VF]
	58a5  T522-50A5 Unified Wire Ethernet Controller [VF]
	58a6  T522-50A6 Unified Wire Ethernet Controller [VF]
	58a7  T580-50A7 Unified Wire Ethernet Controller [VF]
	58a8  T580-50A8 Unified Wire Ethernet Controller [VF]
	58a9  T580-50A9 Unified Wire Ethernet Controller [VF]
	58aa  T580-50AA Unified Wire Ethernet Controller [VF]
	58ab  T520-50AB Unified Wire Ethernet Controller [VF]
	6001  T6225-CR Unified Wire Ethernet Controller
	6002  T6225-SO-CR Unified Wire Ethernet Controller
	6003  T6425-CR Unified Wire Ethernet Controller
	6004  T6425-SO-CR Unified Wire Ethernet Controller
	6005  T6225-OCP-SO Unified Wire Ethernet Controller
	6006  T62100-OCP-SO Unified Wire Ethernet Controller
	6007  T62100-LP-CR Unified Wire Ethernet Controller
	6008  T62100-SO-CR Unified Wire Ethernet Controller
	6009  T6210-BT Unified Wire Ethernet Controller
	600d  T62100-CR Unified Wire Ethernet Controller
	6011  T6225-LL-CR Unified Wire Ethernet Controller
	6014  T61100-OCP-SO Unified Wire Ethernet Controller
	6015  T6201-BT Unified Wire Ethernet Controller
	6080  T6225-6080 Unified Wire Ethernet Controller
	6081  T62100-6081 Unified Wire Ethernet Controller
	6082  T6225-6082 Unified Wire Ethernet Controller
	6083  T62100-6083 Unified Wire Ethernet Controller
	6084  T64100-6084 Unified Wire Ethernet Controller
	6085  T6240-6085 Unified Wire Ethernet Controller
	6401  T6225-CR Unified Wire Ethernet Controller
	6402  T6225-SO-CR Unified Wire Ethernet Controller
	6403  T6425-CR Unified Wire Ethernet Controller
	6404  T6425-SO-CR Unified Wire Ethernet Controller
	6405  T6225-OCP-SO Unified Wire Ethernet Controller
	6406  T62100-OCP-SO Unified Wire Ethernet Controller
	6407  T62100-LP-CR Unified Wire Ethernet Controller
	6408  T62100-SO-CR Unified Wire Ethernet Controller
	6409  T6210-BT Unified Wire Ethernet Controller
	640d  T62100-CR Unified Wire Ethernet Controller
	6411  T6225-LL-CR Unified Wire Ethernet Controller
	6414  T61100-OCP-SO Unified Wire Ethernet Controller
	6415  T6201-BT Unified Wire Ethernet Controller
	6480  T6225-6080 Unified Wire Ethernet Controller
	6481  T62100-6081 Unified Wire Ethernet Controller
	6482  T6225-6082 Unified Wire Ethernet Controller
	6483  T62100-6083 Unified Wire Ethernet Controller
	6484  T64100-6084 Unified Wire Ethernet Controller
	6485  T6240-6085 Unified Wire Ethernet Controller
	6501  T6225-CR Unified Wire Storage Controller
	6502  T6225-SO-CR Unified Wire Storage Controller
	6503  T6425-CR Unified Wire Storage Controller
	6504  T6425-SO-CR Unified Wire Storage Controller
	6505  T6225-OCP-SO Unified Wire Storage Controller
	6506  T62100-OCP-SO Unified Wire Storage Controller
	6507  T62100-LP-CR Unified Wire Storage Controller
	6508  T62100-SO-CR Unified Wire Storage Controller
	6509  T6210-BT Unified Wire Storage Controller
	650d  T62100-CR Unified Wire Storage Controller
	6511  T6225-LL-CR Unified Wire Storage Controller
	6514  T61100-OCP-SO Unified Wire Storage Controller
	6515  T6201-BT Unified Wire Storage Controller
	6580  T6225-6080 Unified Wire Storage Controller
	6581  T62100-6081 Unified Wire Storage Controller
	6582  T6225-6082 Unified Wire Storage Controller
	6583  T62100-6083 Unified Wire Storage Controller
	6584  T64100-6084 Unified Wire Storage Controller
	6585  T6240-6085 Unified Wire Storage Controller
	6601  T6225-CR Unified Wire Storage Controller
	6602  T6225-SO-CR Unified Wire Storage Controller
	6603  T6425-CR Unified Wire Storage Controller
	6604  T6425-SO-CR Unified Wire Storage Controller
	6605  T6225-OCP-SO Unified Wire Storage Controller
	6606  T62100-OCP-SO Unified Wire Storage Controller
	6607  T62100-LP-CR Unified Wire Storage Controller
	6608  T62100-SO-CR Unified Wire Storage Controller
	6609  T6210-BT Unified Wire Storage Controller
	660d  T62100-CR Unified Wire Storage Controller
	6611  T6225-LL-CR Unified Wire Storage Controller
	6614  T61100-OCP-SO Unified Wire Storage Controller
	6615  T6201-BT Unified Wire Storage Controller
	6680  T6225-6080 Unified Wire Storage Controller
	6681  T62100-6081 Unified Wire Storage Controller
	6682  T6225-6082 Unified Wire Storage Controller
	6683  T62100-6083 Unified Wire Storage Controller
	6684  T64100-6084 Unified Wire Storage Controller
	6685  T6240-6085 Unified Wire Storage Controller
	6801  T6225-CR Unified Wire Ethernet Controller [VF]
	6802  T6225-SO-CR Unified Wire Ethernet Controller [VF]
	6803  T6425-CR Unified Wire Ethernet Controller [VF]
	6804  T6425-SO-CR Unified Wire Ethernet Controller [VF]
	6805  T6225-OCP-SO Unified Wire Ethernet Controller [VF]
	6806  T62100-OCP-SO Unified Wire Ethernet Controller [VF]
	6807  T62100-LP-CR Unified Wire Ethernet Controller [VF]
	6808  T62100-SO-CR Unified Wire Ethernet Controller [VF]
	6809  T6210-BT Unified Wire Ethernet Controller [VF]
	680d  T62100-CR Unified Wire Ethernet Controller [VF]
	6811  T6225-LL-CR Unified Wire Ethernet Controller [VF]
	6814  T61100-OCP-SO Unified Wire Ethernet Controller [VF]
	6815  T6201-BT Unified Wire Ethernet Controller [VF]
	6880  T6225-6080 Unified Wire Ethernet Controller [VF]
	6881  T62100-6081 Unified Wire Ethernet Controller [VF]
	6882  T6225-6082 Unified Wire Ethernet Controller [VF]
	6883  T62100-6083 Unified Wire Ethernet Controller [VF]
	6884  T64100-6084 Unified Wire Ethernet Controller [VF]
	6885  T6240-6085 Unified Wire Ethernet Controller [VF]
	a000  PE10K Unified Wire Ethernet Controller
1426  Storage Technology Corp.
1427  Better On-Line Solutions
1428  Edec Co Ltd
1429  Unex Technology Corp.
142a  Kingmax Technology Inc
142b  Radiolan
142c  Minton Optic Industry Co Ltd
142d  Pix stream Inc
142e  Vitec Multimedia
	4020  VM2-2 [Video Maker 2] MPEG1/2 Encoder
	4337  VM2-2-C7 [Video Maker 2 rev. C7] MPEG1/2 Encoder
142f  Radicom Research Inc
1430  ITT Aerospace/Communications Division
1431  Gilat Satellite Networks
1432  Edimax Computer Co.
	9130  RTL81xx Fast Ethernet
1433  Eltec Elektronik GmbH
1435  RTD Embedded Technologies, Inc.
	4520  PCI4520
	6020  SPM6020
	6030  SPM6030
	6420  SPM186420
	6430  SPM176430
	6431  SPM176431
	7520  DM7520
	7540  SDM7540
	7820  DM7820
1436  CIS Technology Inc
1437  Nissin Inc Co
1438  Atmel-dream
1439  Outsource Engineering & Mfg. Inc
143a  Stargate Solutions Inc
143b  Canon Research Center, America
143c  Amlogic Inc
143d  Tamarack Microelectronics Inc
143e  Jones Futurex Inc
143f  Lightwell Co Ltd - Zax Division
1440  ALGOL Corp.
1441  AGIE Ltd
1442  Phoenix Contact GmbH & Co.
1443  Unibrain S.A.
1444  TRW
1445  Logical DO Ltd
1446  Graphin Co Ltd
1447  AIM GmBH
1448  Alesis Studio Electronics
1449  TUT Systems Inc
144a  Adlink Technology
	6208  PCI-6208V
	7250  PCI-7250
	7296  PCI-7296
	7432  PCI-7432
	7433  PCI-7433
	7434  PCI-7434
	7841  PCI-7841
	8133  PCI-8133
	8164  PCI-8164
	8554  PCI-8554
	9111  PCI-9111
	9113  PCI-9113
	9114  PCI-9114
	a001  ADi-BSEC
144b  Verint Systems Inc.
144c  Catalina Research Inc
144d  Samsung Electronics Co Ltd
	1600  Apple PCIe SSD
	a800  XP941 PCIe SSD
	a802  NVMe SSD Controller SM951/PM951
	a804  NVMe SSD Controller SM961/PM961
	a820  NVMe SSD Controller 171X
	a821  NVMe SSD Controller 172X
	a822  NVMe SSD Controller 172Xa
144e  OLITEC
144f  Askey Computer Corp.
1450  Octave Communications Ind.
1451  SP3D Chip Design GmBH
1453  MYCOM Inc
1454  Altiga Networks
1455  Logic Plus Plus Inc
1456  Advanced Hardware Architectures
1457  Nuera Communications Inc
1458  Gigabyte Technology Co., Ltd
1459  DOOIN Electronics
145a  Escalate Networks Inc
145b  PRAIM SRL
145c  Cryptek
145d  Gallant Computer Inc
145e  Aashima Technology B.V.
145f  Baldor Electric Company
	0001  NextMove PCI
1460  DYNARC INC
1461  Avermedia Technologies Inc
	a3ce  M179
	a3cf  M179
	a836  M115 DVB-T, PAL/SECAM/NTSC Tuner
	e836  M115S Hybrid Analog/DVB PAL/SECAM/NTSC Tuner
	f436  AVerTV Hybrid+FM
1462  Micro-Star International Co., Ltd. [MSI]
	aaf0  Radeon RX 580 Gaming X 8G
1463  Fast Corporation
1464  Interactive Circuits & Systems Ltd
1465  GN NETTEST Telecom DIV.
1466  Designpro Inc.
1467  DIGICOM SPA
1468  AMBIT Microsystem Corp.
1469  Cleveland Motion Controls
146a  Aeroflex
	3010  3010 RF Synthesizer
	3a11  3011A PXI RF Synthesizer
146b  Parascan Technologies Ltd
146c  Ruby Tech Corp.
	1430  FE-1430TX Fast Ethernet PCI Adapter
146d  Tachyon, INC.
146e  Williams Electronics Games, Inc.
146f  Multi Dimensional Consulting Inc
1470  Bay Networks
1471  Integrated Telecom Express Inc
1472  DAIKIN Industries, Ltd
1473  ZAPEX Technologies Inc
1474  Doug Carson & Associates
1475  PICAZO Communications
1476  MORTARA Instrument Inc
1477  Net Insight
1478  DIATREND Corporation
1479  TORAY Industries Inc
147a  FORMOSA Industrial Computing
147b  ABIT Computer Corp.
	1084  IP35 [Dark Raider]
147c  AWARE, Inc.
147d  Interworks Computer Products
147e  Matsushita Graphic Communication Systems, Inc.
147f  NIHON UNISYS, Ltd.
1480  SCII Telecom
1481  BIOPAC Systems Inc
1482  ISYTEC - Integrierte Systemtechnik GmBH
	0001  PCI-16 Host Interface for ITC-16
1483  LABWAY Corporation
1484  Logic Corporation
1485  ERMA - Electronic GmBH
1486  L3 Communications Telemetry & Instrumentation
1487  MARQUETTE Medical Systems
1489  KYE Systems Corporation
148a  OPTO
148b  INNOMEDIALOGIC Inc.
148c  Tul Corporation / PowerColor
148d  DIGICOM Systems, Inc.
	1003  HCF 56k Data/Fax Modem
148e  OSI Plus Corporation
148f  Plant Equipment, Inc.
1490  Stone Microsystems PTY Ltd.
1491  ZEAL Corporation
1492  Time Logic Corporation
1493  MAKER Communications
1494  WINTOP Technology, Inc.
1495  TOKAI Communications Industry Co. Ltd
1496  JOYTECH Computer Co., Ltd.
1497  SMA Regelsysteme GmBH
	1497  SMA Technologie AG
1498  TEWS Technologies GmbH
	0330  TPMC816 2 Channel CAN bus controller.
	035d  TPMC861 4-Channel Isolated Serial Interface RS422/RS485
	0385  TPMC901 Extended CAN bus with 2/4/6 CAN controller
	21cc  TCP460 CompactPCI 16 Channel Serial Interface RS232/RS422
	21cd  TCP461 CompactPCI 8 Channel Serial Interface RS232/RS422
	3064  TPCI100 (2 Slot IndustryPack PCI Carrier)
	30c8  TPCI200 4 Slot IndustryPack PCI Carrier
	70c8  TPCE200 4 Slot IndustryPack PCIe Carrier
	9177  TXMC375 8 channel RS232/RS422/RS485 programmable serial interface
1499  EMTEC CO., Ltd
149a  ANDOR Technology Ltd
149b  SEIKO Instruments Inc
149c  OVISLINK Corp.
149d  NEWTEK Inc
	0001  Video Toaster for PC
149e  Mapletree Networks Inc.
149f  LECTRON Co Ltd
14a0  SOFTING GmBH
14a1  Systembase Co Ltd
14a2  Millennium Engineering Inc
14a3  Maverick Networks
14a4  Lite-On Technology Corporation
	22f1  M8Pe Series NVMe SSD
	4318  Broadcom BCM4318 [AirForce One 54g] 802.11g WLAN Controller
14a5  XIONICS Document Technologies Inc
14a6  INOVA Computers GmBH & Co KG
14a7  MYTHOS Systems Inc
14a8  FEATRON Technologies Corporation
14a9  HIVERTEC Inc
14aa  Advanced MOS Technology Inc
14ab  Mentor Graphics Corp.
14ac  Novaweb Technologies Inc
14ad  Time Space Radio AB
14ae  CTI, Inc
14af  Guillemot Corporation
	7102  3D Prophet II MX
14b0  BST Communication Technology Ltd
14b1  Nextcom K.K.
14b2  ENNOVATE Networks Inc
14b3  XPEED Inc
	0000  DSL NIC
14b4  PHILIPS Business Electronics B.V.
14b5  Creamware GmBH
	0200  Scope
	0300  Pulsar
	0400  PulsarSRB
	0600  Pulsar2
	0800  DSP-Board
	0900  DSP-Board
	0a00  DSP-Board
	0b00  DSP-Board
14b6  Quantum Data Corp.
14b7  PROXIM Inc
	0001  Symphony 4110
14b8  Techsoft Technology Co Ltd
14b9  Cisco Aironet Wireless Communications
	0001  PC4800
	0340  PC4800
	0350  350 series 802.11b Wireless LAN Adapter
	4500  PC4500
	4800  Cisco Aironet 340 802.11b Wireless LAN Adapter/Aironet PC4800
	a504  Cisco Aironet Wireless 802.11b
	a505  Cisco Aironet CB20a 802.11a Wireless LAN Adapter
	a506  Cisco Aironet Mini PCI b/g
14ba  INTERNIX Inc.
	0600  ARC-PCI/22
14bb  SEMTECH Corporation
14bc  Globespan Semiconductor Inc.
	d002  Pulsar [PCI ADSL Card]
	d00f  Pulsar [PCI ADSL Card]
14bd  CARDIO Control N.V.
14be  L3 Communications
14bf  SPIDER Communications Inc.
14c0  COMPAL Electronics Inc
14c1  MYRICOM Inc.
	0008  Myri-10G Dual-Protocol NIC
	8043  Myrinet 2000 Scalable Cluster Interconnect
14c2  DTK Computer
14c3  MEDIATEK Corp.
	7630  MT7630e 802.11bgn Wireless Network Adapter
	7662  MT7662E 802.11ac PCI Express Wireless Network Adapter
14c4  IWASAKI Information Systems Co Ltd
14c5  Automation Products AB
14c6  Data Race Inc
14c7  Modular Technology Holdings Ltd
14c8  Turbocomm Tech. Inc.
14c9  ODIN Telesystems Inc
14ca  PE Logic Corp.
14cb  Billionton Systems Inc
14cc  NAKAYO Telecommunications Inc
14cd  Universal Global Scientific Industrial Co.,Ltd
	0001  USI-1514-1GbaseT [OCP1]
	0002  USI-4227-SFP [OCP2]
	0003  USI-X557-10GbaseT [OCP3]
14ce  Whistle Communications
14cf  TEK Microsystems Inc.
14d0  Ericsson Axe R & D
14d1  Computer Hi-Tech Co Ltd
14d2  Titan Electronics Inc
	8001  VScom 010L 1 port parallel adaptor
	8002  VScom 020L 2 port parallel adaptor
	8010  VScom 100L 1 port serial adaptor
	8011  VScom 110L 1 port serial and 1 port parallel adaptor
	8020  VScom 200L 1 or 2 port serial adaptor
	8021  VScom 210L 2 port serial and 1 port parallel adaptor
	8028  VScom 200I/200I-SI 2-port serial adapter
	8040  VScom 400L 4 port serial adaptor
	8043  VScom 430L 4-port serial and 3-port parallel adapter
	8048  VScom 400I 4-port serial adapter
	8080  VScom 800L 8 port serial adaptor
	8088  VScom 800I 8-port serial adapter
	a000  VScom 010H 1 port parallel adaptor
	a001  VScom 100H 1 port serial adaptor
	a003  VScom 400H 4 port serial adaptor
	a004  VScom 400HF1 4 port serial adaptor
	a005  VScom 200H 2 port serial adaptor
	a007  VScom PCI800EH (PCIe) 8-port serial adapter Port 1-4
	a008  VScom PCI800EH (PCIe) 8-port serial adapter Port 5-8
	a009  VScom PCI400EH (PCIe) 4-port serial adapter
	e001  VScom 010HV2 1 port parallel adaptor
	e010  VScom 100HV2 1 port serial adaptor
	e020  VScom 200HV2 2 port serial adaptor
14d3  CIRTECH (UK) Ltd
14d4  Panacom Technology Corp
14d5  Nitsuko Corporation
14d6  Accusys Inc
	6101  ACS-61xxx, PCIe to SAS/SATA RAID HBA
	6201  ACS-62xxx, External PCIe to SAS/SATA RAID controller
14d7  Hirakawa Hewtech Corp
14d8  HOPF Elektronik GmBH
14d9  Alliance Semiconductor Corporation
	0010  AP1011/SP1011 HyperTransport-PCI Bridge [Sturgeon]
	9000  AS90L10204/10208 HyperTransport to PCI-X Bridge
14da  National Aerospace Laboratories
14db  AFAVLAB Technology Inc
	2120  TK9902
	2182  AFAVLAB Technology Inc. 8-port serial card
14dc  Amplicon Liveline Ltd
	0000  PCI230
	0001  PCI242
	0002  PCI244
	0003  PCI247
	0004  PCI248
	0005  PCI249
	0006  PCI260
	0007  PCI224
	0008  PCI234
	0009  PCI236
	000a  PCI272
	000b  PCI215
14dd  Boulder Design Labs Inc
14de  Applied Integration Corporation
14df  ASIC Communications Corp
14e1  INVERTEX
14e2  INFOLIBRIA
14e3  AMTELCO
14e4  Broadcom Limited
	0576  BCM43224 802.11a/b/g/n
	0800  Sentry5 Chipcommon I/O Controller
	0804  Sentry5 PCI Bridge
	0805  Sentry5 MIPS32 CPU
	0806  Sentry5 Ethernet Controller
	080b  Sentry5 Crypto Accelerator
	080f  Sentry5 DDR/SDR RAM Controller
	0811  Sentry5 External Interface Core
	0816  BCM3302 Sentry5 MIPS32 CPU
	1570  720p FaceTime HD Camera
	1600  NetXtreme BCM5752 Gigabit Ethernet PCI Express
	1601  NetXtreme BCM5752M Gigabit Ethernet PCI Express
	1612  BCM70012 Video Decoder [Crystal HD]
	1615  BCM70015 Video Decoder [Crystal HD]
	1639  NetXtreme II BCM5709 Gigabit Ethernet
	163a  NetXtreme II BCM5709S Gigabit Ethernet
	163b  NetXtreme II BCM5716 Gigabit Ethernet
	163c  NetXtreme II BCM5716S Gigabit Ethernet
	163d  NetXtreme II BCM57811 10-Gigabit Ethernet
	163e  NetXtreme II BCM57811 10 Gigabit Ethernet Multi Function
	163f  NetXtreme II BCM57811 10-Gigabit Ethernet Virtual Function
	1641  NetXtreme BCM57787 Gigabit Ethernet PCIe
	1642  NetXtreme BCM57764 Gigabit Ethernet PCIe
	1643  NetXtreme BCM5725 Gigabit Ethernet PCIe
	1644  NetXtreme BCM5700 Gigabit Ethernet
	1645  NetXtreme BCM5701 Gigabit Ethernet
	1646  NetXtreme BCM5702 Gigabit Ethernet
	1647  NetXtreme BCM5703 Gigabit Ethernet
	1648  NetXtreme BCM5704 Gigabit Ethernet
	1649  NetXtreme BCM5704S_2 Gigabit Ethernet
	164a  NetXtreme II BCM5706 Gigabit Ethernet
	164c  NetXtreme II BCM5708 Gigabit Ethernet
	164d  NetXtreme BCM5702FE Gigabit Ethernet
	164e  NetXtreme II BCM57710 10-Gigabit PCIe [Everest]
	164f  NetXtreme II BCM57711 10-Gigabit PCIe
	1650  NetXtreme II BCM57711E 10-Gigabit PCIe
	1653  NetXtreme BCM5705 Gigabit Ethernet
	1654  NetXtreme BCM5705_2 Gigabit Ethernet
	1655  NetXtreme BCM5717 Gigabit Ethernet PCIe
	1656  NetXtreme BCM5718 Gigabit Ethernet PCIe
	1657  NetXtreme BCM5719 Gigabit Ethernet PCIe
	1659  NetXtreme BCM5721 Gigabit Ethernet PCI Express
	165a  NetXtreme BCM5722 Gigabit Ethernet PCI Express
	165b  NetXtreme BCM5723 Gigabit Ethernet PCIe
	165c  NetXtreme BCM5724 Gigabit Ethernet PCIe
	165d  NetXtreme BCM5705M Gigabit Ethernet
	165e  NetXtreme BCM5705M_2 Gigabit Ethernet
	165f  NetXtreme BCM5720 Gigabit Ethernet PCIe
	1662  NetXtreme II BCM57712 10 Gigabit Ethernet
	1663  NetXtreme II BCM57712 10 Gigabit Ethernet Multi Function
	1665  NetXtreme BCM5717 Gigabit Ethernet PCIe
	1668  NetXtreme BCM5714 Gigabit Ethernet
	1669  NetXtreme 5714S Gigabit Ethernet
	166a  NetXtreme BCM5780 Gigabit Ethernet
	166b  NetXtreme BCM5780S Gigabit Ethernet
	166e  570x 10/100 Integrated Controller
	166f  NetXtreme II BCM57712 10 Gigabit Ethernet Virtual Function
	1672  NetXtreme BCM5754M Gigabit Ethernet PCI Express
	1673  NetXtreme BCM5755M Gigabit Ethernet PCI Express
	1674  NetXtreme BCM5756ME Gigabit Ethernet PCI Express
	1677  NetXtreme BCM5751 Gigabit Ethernet PCI Express
	1678  NetXtreme BCM5715 Gigabit Ethernet
	1679  NetXtreme BCM5715S Gigabit Ethernet
	167a  NetXtreme BCM5754 Gigabit Ethernet PCI Express
	167b  NetXtreme BCM5755 Gigabit Ethernet PCI Express
	167d  NetXtreme BCM5751M Gigabit Ethernet PCI Express
	167e  NetXtreme BCM5751F Fast Ethernet PCI Express
	167f  NetLink BCM5787F Fast Ethernet PCI Express
	1680  NetXtreme BCM5761e Gigabit Ethernet PCIe
	1681  NetXtreme BCM5761 Gigabit Ethernet PCIe
	1682  NetXtreme BCM57762 Gigabit Ethernet PCIe
	1683  NetXtreme BCM57767 Gigabit Ethernet PCIe
	1684  NetXtreme BCM5764M Gigabit Ethernet PCIe
	1685  NetXtreme II BCM57500S Gigabit Ethernet
	1686  NetXtreme BCM57766 Gigabit Ethernet PCIe
	1687  NetXtreme BCM5762 Gigabit Ethernet PCIe
	1688  NetXtreme BCM5761 10/100/1000BASE-T Ethernet
	168a  NetXtreme II BCM57800 1/10 Gigabit Ethernet
	168d  NetXtreme II BCM57840 10/20 Gigabit Ethernet
	168e  NetXtreme II BCM57810 10 Gigabit Ethernet
	1690  NetXtreme BCM57760 Gigabit Ethernet PCIe
	1691  NetLink BCM57788 Gigabit Ethernet PCIe
	1692  NetLink BCM57780 Gigabit Ethernet PCIe
	1693  NetLink BCM5787M Gigabit Ethernet PCI Express
	1694  NetLink BCM57790 Gigabit Ethernet PCIe
	1696  NetXtreme BCM5782 Gigabit Ethernet
	1698  NetLink BCM5784M Gigabit Ethernet PCIe
	1699  NetLink BCM5785 Gigabit Ethernet
	169a  NetLink BCM5786 Gigabit Ethernet PCI Express
	169b  NetLink BCM5787 Gigabit Ethernet PCI Express
	169c  NetXtreme BCM5788 Gigabit Ethernet
	169d  NetLink BCM5789 Gigabit Ethernet PCI Express
	16a0  NetLink BCM5785 Fast Ethernet
	16a1  BCM57840 NetXtreme II 10 Gigabit Ethernet
	16a2  BCM57840 NetXtreme II 10/20-Gigabit Ethernet
	16a3  NetXtreme BCM57786 Gigabit Ethernet PCIe
	16a4  BCM57840 NetXtreme II Ethernet Multi Function
	16a5  NetXtreme II BCM57800 1/10 Gigabit Ethernet Multi Function
	16a6  NetXtreme BCM5702X Gigabit Ethernet
	16a7  NetXtreme BCM5703X Gigabit Ethernet
	16a8  NetXtreme BCM5704S Gigabit Ethernet
	16a9  NetXtreme II BCM57800 1/10 Gigabit Ethernet Virtual Function
	16aa  NetXtreme II BCM5706S Gigabit Ethernet
	16ab  NetXtreme II BCM57840 10/20 Gigabit Ethernet Multi Function
	16ac  NetXtreme II BCM5708S Gigabit Ethernet
	16ad  NetXtreme II BCM57840 10/20 Gigabit Ethernet Virtual Function
	16ae  NetXtreme II BCM57810 10 Gigabit Ethernet Multi Function
	16af  NetXtreme II BCM57810 10 Gigabit Ethernet Virtual Function
	16b0  NetXtreme BCM57761 Gigabit Ethernet PCIe
	16b1  NetLink BCM57781 Gigabit Ethernet PCIe
	16b2  NetLink BCM57791 Gigabit Ethernet PCIe
	16b3  NetXtreme BCM57786 Gigabit Ethernet PCIe
	16b4  NetXtreme BCM57765 Gigabit Ethernet PCIe
	16b5  NetLink BCM57785 Gigabit Ethernet PCIe
	16b6  NetLink BCM57795 Gigabit Ethernet PCIe
	16b7  NetXtreme BCM57782 Gigabit Ethernet PCIe
	16bc  BCM57765/57785 SDXC/MMC Card Reader
	16be  BCM57765/57785 MS Card Reader
	16bf  BCM57765/57785 xD-Picture Card Reader
	16c1  NetXtreme-E RDMA Virtual Function
	16c6  NetXtreme BCM5702A3 Gigabit Ethernet
	16c7  NetXtreme BCM5703 Gigabit Ethernet
	16c8  BCM57301 NetXtreme-C 10Gb Ethernet Controller
	16c9  BCM57302 NetXtreme-C 10Gb/25Gb Ethernet Controller
	16ca  BCM57304 NetXtreme-C 10Gb/25Gb/40Gb/50Gb Ethernet Controller
	16cb  NetXtreme-C Ethernet Virtual Function
	16cc  BCM57417 NetXtreme-E Ethernet Partition
	16ce  BCM57311 NetXtreme-C 10Gb RDMA Ethernet Controller
	16cf  BCM57312 NetXtreme-C 10Gb/25Gb RDMA Ethernet Controller
	16d0  BCM57402 NetXtreme-E 10Gb Ethernet Controller
	16d1  BCM57404 NetXtreme-E 10Gb/25Gb Ethernet Controller
	16d2  BCM57406 NetXtreme-E 10GBASE-T Ethernet Controller
	16d3  NetXtreme-E Ethernet Virtual Function
	16d4  BCM57402 NetXtreme-E Ethernet Partition
	16d5  BCM57407 NetXtreme-E 10GBase-T Ethernet Controller
	16d6  BCM57412 NetXtreme-E 10Gb RDMA Ethernet Controller
	16d7  BCM57414 NetXtreme-E 10Gb/25Gb RDMA Ethernet Controller
	16d8  BCM57416 NetXtreme-E 10GBase-T RDMA Ethernet Controller
	16d9  BCM57417 NetXtreme-E 10GBASE-T RDMA Ethernet Controller
	16dc  NetXtreme-E Ethernet Virtual Function
	16dd  NetLink BCM5781 Gigabit Ethernet PCI Express
	16de  BCM57412 NetXtreme-E Ethernet Partition
	16df  BCM57314 NetXtreme-C 10Gb/25Gb/40Gb/50Gb RDMA Ethernet Controller
	16e1  NetXtreme-C Ethernet Virtual Function
	16e2  BCM57417 NetXtreme-E 10Gb/25Gb RDMA Ethernet Controller
	16e3  BCM57416 NetXtreme-E 10Gb RDMA Ethernet Controller
	16e5  NetXtreme-C RDMA Virtual Function
	16e7  BCM57404 NetXtreme-E Ethernet Partition
	16e8  BCM57406 NetXtreme-E Ethernet Partition
	16e9  BCM57407 NetXtreme-E 25Gb Ethernet Controller
	16ec  BCM57414 NetXtreme-E Ethernet Partition
	16ed  BCM57414 NetXtreme-E RDMA Partition
	16ee  BCM57416 NetXtreme-E Ethernet Partition
	16ef  BCM57416 NetXtreme-E RDMA Partition
	16f3  NetXtreme BCM5727 Gigabit Ethernet PCIe
	16f7  NetXtreme BCM5753 Gigabit Ethernet PCI Express
	16fd  NetXtreme BCM5753M Gigabit Ethernet PCI Express
	16fe  NetXtreme BCM5753F Fast Ethernet PCI Express
	170c  BCM4401-B0 100Base-TX
	170d  NetXtreme BCM5901 100Base-TX
	170e  NetXtreme BCM5901 100Base-TX
	1712  NetLink BCM5906 Fast Ethernet PCI Express
	1713  NetLink BCM5906M Fast Ethernet PCI Express
	3352  BCM3352
	3360  BCM3360
	4210  BCM4210 iLine10 HomePNA 2.0
	4211  BCM4211 iLine10 HomePNA 2.0 + V.90 56k modem
	4212  BCM4212 v.90 56k modem
	4220  802-11b/g Wireless PCI controller, packaged as a Linksys WPC54G ver 1.2 PCMCIA card
	4222  NetXtreme BCM5753M Gigabit Ethernet PCI Express
	4301  BCM4301 802.11b Wireless LAN Controller
	4305  BCM4307 V.90 56k Modem
	4306  BCM4306 802.11bg Wireless LAN controller
	4307  BCM4306 802.11bg Wireless LAN Controller
	4310  BCM4310 Chipcommon I/OController
	4311  BCM4311 802.11b/g WLAN
	4312  BCM4311 802.11a/b/g
	4313  BCM4311 802.11a
	4315  BCM4312 802.11b/g LP-PHY
	4318  BCM4318 [AirForce One 54g] 802.11g Wireless LAN Controller
	4319  BCM4318 [AirForce 54g] 802.11a/b/g PCI Express Transceiver
	4320  BCM4306 802.11b/g Wireless LAN Controller
	4321  BCM4321 802.11a Wireless Network Controller
	4322  BCM4322 802.11bgn Wireless Network Controller
	4324  BCM4309 802.11abg Wireless Network Controller
	4325  BCM4306 802.11bg Wireless Network Controller
	4326  BCM4307 Chipcommon I/O Controller?
	4328  BCM4321 802.11a/b/g/n
	4329  BCM4321 802.11b/g/n
	432a  BCM4321 802.11an Wireless Network Controller
	432b  BCM4322 802.11a/b/g/n Wireless LAN Controller
	432c  BCM4322 802.11b/g/n
	432d  BCM4322 802.11an Wireless Network Controller
	4331  BCM4331 802.11a/b/g/n
	4333  Serial (EDGE/GPRS modem part of Option GT Combo Edge)
	4344  EDGE/GPRS data and 802.11b/g combo cardbus [GC89]
	4350  BCM43222 Wireless Network Adapter
	4351  BCM43222 802.11abgn Wireless Network Adapter
	4353  BCM43224 802.11a/b/g/n
	4357  BCM43225 802.11b/g/n
	4358  BCM43227 802.11b/g/n
	4359  BCM43228 802.11a/b/g/n
	4360  BCM4360 802.11ac Wireless Network Adapter
	4365  BCM43142 802.11b/g/n
	43a0  BCM4360 802.11ac Wireless Network Adapter
	43a1  BCM4360 802.11ac Wireless Network Adapter
	43a2  BCM4360 802.11ac Wireless Network Adapter
	43a3  BCM4350 802.11ac Wireless Network Adapter
	43a9  BCM43217 802.11b/g/n
	43aa  BCM43131 802.11b/g/n
	43ae  BCM43162 802.11ac Wireless Network Adapter
	43b1  BCM4352 802.11ac Wireless Network Adapter
	43ba  BCM43602 802.11ac Wireless LAN SoC
	43bb  BCM43602 802.11ac Wireless LAN SoC
	43bc  BCM43602 802.11ac Wireless LAN SoC
	43d3  BCM43567 802.11ac Wireless Network Adapter
	43d9  BCM43570 802.11ac Wireless Network Adapter
	43df  BCM4354 802.11ac Wireless LAN SoC
	43e9  BCM4358 802.11ac Wireless LAN SoC
	43ec  BCM4356 802.11ac Wireless Network Adapter
	4401  BCM4401 100Base-T
	4402  BCM4402 Integrated 10/100BaseT
	4403  BCM4402 V.90 56k Modem
	4410  BCM4413 iLine32 HomePNA 2.0
	4411  BCM4413 V.90 56k modem
	4412  BCM4412 10/100BaseT
	4430  BCM44xx CardBus iLine32 HomePNA 2.0
	4432  BCM4432 CardBus 10/100BaseT
	4610  BCM4610 Sentry5 PCI to SB Bridge
	4611  BCM4610 Sentry5 iLine32 HomePNA 1.0
	4612  BCM4610 Sentry5 V.90 56k Modem
	4613  BCM4610 Sentry5 Ethernet Controller
	4614  BCM4610 Sentry5 External Interface
	4615  BCM4610 Sentry5 USB Controller
	4704  BCM4704 PCI to SB Bridge
	4705  BCM4704 Sentry5 802.11b Wireless LAN Controller
	4706  BCM4704 Sentry5 Ethernet Controller
	4707  BCM4704 Sentry5 USB Controller
	4708  BCM4704 Crypto Accelerator
	4710  BCM4710 Sentry5 PCI to SB Bridge
	4711  BCM47xx Sentry5 iLine32 HomePNA 2.0
	4712  BCM47xx V.92 56k modem
	4713  Sentry5 Ethernet Controller
	4714  BCM47xx Sentry5 External Interface
	4715  BCM47xx Sentry5 USB / Ethernet Controller
	4716  BCM47xx Sentry5 USB Host Controller
	4717  BCM47xx Sentry5 USB Device Controller
	4718  Sentry5 Crypto Accelerator
	4719  BCM47xx/53xx RoboSwitch Core
	4720  BCM4712 MIPS CPU
	4727  BCM4313 802.11bgn Wireless Network Adapter
	5365  BCM5365P Sentry5 Host Bridge
	5600  BCM5600 StrataSwitch 24+2 Ethernet Switch Controller
	5605  BCM5605 StrataSwitch 24+2 Ethernet Switch Controller
	5615  BCM5615 StrataSwitch 24+2 Ethernet Switch Controller
	5625  BCM5625 StrataSwitch 24+2 Ethernet Switch Controller
	5645  BCM5645 StrataSwitch 24+2 Ethernet Switch Controller
	5670  BCM5670 8-Port 10GE Ethernet Switch Fabric
	5680  BCM5680 G-Switch 8 Port Gigabit Ethernet Switch Controller
	5690  BCM5690 12-port Multi-Layer Gigabit Ethernet Switch
	5691  BCM5691 GE/10GE 8+2 Gigabit Ethernet Switch Controller
	5692  BCM5692 12-port Multi-Layer Gigabit Ethernet Switch
	5695  BCM5695 12-port + HiGig Multi-Layer Gigabit Ethernet Switch
	5698  BCM5698 12-port Multi-Layer Gigabit Ethernet Switch
	5820  BCM5820 Crypto Accelerator
	5821  BCM5821 Crypto Accelerator
	5822  BCM5822 Crypto Accelerator
	5823  BCM5823 Crypto Accelerator
	5824  BCM5824 Crypto Accelerator
	5840  BCM5840 Crypto Accelerator
	5841  BCM5841 Crypto Accelerator
	5850  BCM5850 Crypto Accelerator
	8602  BCM7400/BCM7405 Serial ATA Controller
	a8d8  BCM43224/5 Wireless Network Adapter
	aa52  BCM43602 802.11ac Wireless LAN SoC
	b302  BCM56302 StrataXGS 24x1GE 2x10GE Switch Controller
	b334  BCM56334 StrataXGS 24x1GE 4x10GE Switch Controller
	b800  BCM56800 StrataXGS 10GE Switch Controller
	b842  BCM56842 Trident 10GE Switch Controller
	b850  Broadcom BCM56850 Switch ASIC
	b960  Broadcom BCM56960 Switch ASIC
14e5  Pixelfusion Ltd
14e6  SHINING Technology Inc
14e7  3CX
14e8  RAYCER Inc
14e9  GARNETS System CO Ltd
14ea  Planex Communications, Inc
	ab06  FNW-3603-TX CardBus Fast Ethernet
	ab07  RTL81xx RealTek Ethernet
	ab08  FNW-3602-TX CardBus Fast Ethernet
14eb  SEIKO EPSON Corp
14ec  Agilent Technologies
	0000  Aciris Digitizer (malformed ID)
14ed  DATAKINETICS Ltd
14ee  MASPRO KENKOH Corp
14ef  CARRY Computer ENG. CO Ltd
14f0  CANON RESEACH CENTRE FRANCE
14f1  Conexant Systems, Inc.
	1002  HCF 56k Modem
	1003  HCF 56k Modem
	1004  HCF 56k Modem
	1005  HCF 56k Modem
	1006  HCF 56k Modem
	1022  HCF 56k Modem
	1023  HCF 56k Modem
	1024  HCF 56k Modem
	1025  HCF 56k Modem
	1026  HCF 56k Modem
	1032  HCF 56k Modem
	1033  HCF 56k Data/Fax Modem
	1034  HCF 56k Data/Fax/Voice Modem
	1035  HCF 56k Data/Fax/Voice/Spkp (w/Handset) Modem
	1036  HCF 56k Data/Fax/Voice/Spkp Modem
	1052  HCF 56k Data/Fax Modem (Worldwide)
	1053  HCF 56k Data/Fax Modem (Worldwide)
	1054  HCF 56k Data/Fax/Voice Modem (Worldwide)
	1055  HCF 56k Data/Fax/Voice/Spkp (w/Handset) Modem (Worldwide)
	1056  HCF 56k Data/Fax/Voice/Spkp Modem (Worldwide)
	1057  HCF 56k Data/Fax/Voice/Spkp Modem (Worldwide)
	1059  HCF 56k Data/Fax/Voice Modem (Worldwide)
	1063  HCF 56k Data/Fax Modem
	1064  HCF 56k Data/Fax/Voice Modem
	1065  HCF 56k Data/Fax/Voice/Spkp (w/Handset) Modem
	1066  HCF 56k Data/Fax/Voice/Spkp Modem
	1085  HCF V90 56k Data/Fax/Voice/Spkp PCI Modem
	10b6  CX06834-11 HCF V.92 56k Data/Fax/Voice/Spkp Modem
	1433  HCF 56k Data/Fax Modem
	1434  HCF 56k Data/Fax/Voice Modem
	1435  HCF 56k Data/Fax/Voice/Spkp (w/Handset) Modem
	1436  HCF 56k Data/Fax Modem
	1453  HCF 56k Data/Fax Modem
	1454  HCF 56k Data/Fax/Voice Modem
	1455  HCF 56k Data/Fax/Voice/Spkp (w/Handset) Modem
	1456  HCF 56k Data/Fax/Voice/Spkp Modem
	1610  ADSL AccessRunner PCI Arbitration Device
	1611  AccessRunner PCI ADSL Interface Device
	1620  AccessRunner V2 PCI ADSL Arbitration Device
	1621  AccessRunner V2 PCI ADSL Interface Device
	1622  AccessRunner V2 PCI ADSL Yukon WAN Adapter
	1803  HCF 56k Modem
	1811  MiniPCI Network Adapter
	1815  HCF 56k Modem
	1830  CX861xx Integrated Host Bridge
	2003  HSF 56k Data/Fax Modem
	2004  HSF 56k Data/Fax/Voice Modem
	2005  HSF 56k Data/Fax/Voice/Spkp (w/Handset) Modem
	2006  HSF 56k Data/Fax/Voice/Spkp Modem
	2013  HSF 56k Data/Fax Modem
	2014  HSF 56k Data/Fax/Voice Modem
	2015  HSF 56k Data/Fax/Voice/Spkp (w/Handset) Modem
	2016  HSF 56k Data/Fax/Voice/Spkp Modem
	2043  HSF 56k Data/Fax Modem (WorldW SmartDAA)
	2044  HSF 56k Data/Fax/Voice Modem (WorldW SmartDAA)
	2045  HSF 56k Data/Fax/Voice/Spkp (w/Handset) Modem (WorldW SmartDAA)
	2046  HSF 56k Data/Fax/Voice/Spkp Modem (WorldW SmartDAA)
	2063  HSF 56k Data/Fax Modem (SmartDAA)
	2064  HSF 56k Data/Fax/Voice Modem (SmartDAA)
	2065  HSF 56k Data/Fax/Voice/Spkp (w/Handset) Modem (SmartDAA)
	2066  HSF 56k Data/Fax/Voice/Spkp Modem (SmartDAA)
	2093  HSF 56k Modem
	2143  HSF 56k Data/Fax/Cell Modem (Mob WorldW SmartDAA)
	2144  HSF 56k Data/Fax/Voice/Cell Modem (Mob WorldW SmartDAA)
	2145  HSF 56k Data/Fax/Voice/Spkp (w/HS)/Cell Modem (Mob WorldW SmartDAA)
	2146  HSF 56k Data/Fax/Voice/Spkp/Cell Modem (Mob WorldW SmartDAA)
	2163  HSF 56k Data/Fax/Cell Modem (Mob SmartDAA)
	2164  HSF 56k Data/Fax/Voice/Cell Modem (Mob SmartDAA)
	2165  HSF 56k Data/Fax/Voice/Spkp (w/HS)/Cell Modem (Mob SmartDAA)
	2166  HSF 56k Data/Fax/Voice/Spkp/Cell Modem (Mob SmartDAA)
	2343  HSF 56k Data/Fax CardBus Modem (Mob WorldW SmartDAA)
	2344  HSF 56k Data/Fax/Voice CardBus Modem (Mob WorldW SmartDAA)
	2345  HSF 56k Data/Fax/Voice/Spkp (w/HS) CardBus Modem (Mob WorldW SmartDAA)
	2346  HSF 56k Data/Fax/Voice/Spkp CardBus Modem (Mob WorldW SmartDAA)
	2363  HSF 56k Data/Fax CardBus Modem (Mob SmartDAA)
	2364  HSF 56k Data/Fax/Voice CardBus Modem (Mob SmartDAA)
	2365  HSF 56k Data/Fax/Voice/Spkp (w/HS) CardBus Modem (Mob SmartDAA)
	2366  HSF 56k Data/Fax/Voice/Spkp CardBus Modem (Mob SmartDAA)
	2443  HSF 56k Data/Fax Modem (Mob WorldW SmartDAA)
	2444  HSF 56k Data/Fax/Voice Modem (Mob WorldW SmartDAA)
	2445  HSF 56k Data/Fax/Voice/Spkp (w/HS) Modem (Mob WorldW SmartDAA)
	2446  HSF 56k Data/Fax/Voice/Spkp Modem (Mob WorldW SmartDAA)
	2463  HSF 56k Data/Fax Modem (Mob SmartDAA)
	2464  HSF 56k Data/Fax/Voice Modem (Mob SmartDAA)
	2465  HSF 56k Data/Fax/Voice/Spkp (w/HS) Modem (Mob SmartDAA)
	2466  HSF 56k Data/Fax/Voice/Spkp Modem (Mob SmartDAA)
	2702  HSFi modem RD01-D270
	2f00  HSF 56k HSFi Modem
	2f02  HSF 56k HSFi Data/Fax
	2f11  HSF 56k HSFi Modem
	2f20  HSF 56k Data/Fax Modem
	2f30  SoftV92 SpeakerPhone SoftRing Modem with SmartSP
	2f50  Conexant SoftK56 Data/Fax Modem
	5b7a  CX23418 Single-Chip MPEG-2 Encoder with Integrated Analog Video/Broadcast Audio Decoder
	8200  CX25850
	8234  RS8234 ATM SAR Controller [ServiceSAR Plus]
	8800  CX23880/1/2/3 PCI Video and Audio Decoder
	8801  CX23880/1/2/3 PCI Video and Audio Decoder [Audio Port]
	8802  CX23880/1/2/3 PCI Video and Audio Decoder [MPEG Port]
	8804  CX23880/1/2/3 PCI Video and Audio Decoder [IR Port]
	8811  CX23880/1/2/3 PCI Video and Audio Decoder [Audio Port]
	8852  CX23885 PCI Video and Audio Decoder
	8880  CX23887/8 PCIe Broadcast Audio and Video Decoder with 3D Comb
14f2  MOBILITY Electronics
	0120  EV1000 bridge
	0121  EV1000 Parallel port
	0122  EV1000 Serial port
	0123  EV1000 Keyboard controller
	0124  EV1000 Mouse controller
14f3  BroadLogic
	2030  2030 DVB-S Satellite Receiver
	2035  2035 DVB-S Satellite Receiver
	2050  2050 DVB-T Terrestrial (Cable) Receiver
	2060  2060 ATSC Terrestrial (Cable) Receiver
14f4  TOKYO Electronic Industry CO Ltd
14f5  SOPAC Ltd
14f6  COYOTE Technologies LLC
14f7  WOLF Technology Inc
14f8  AUDIOCODES Inc
	2077  TP-240 dual span E1 VoIP PCI card
14f9  AG COMMUNICATIONS
14fa  WANDEL & GOLTERMANN
14fb  TRANSAS MARINE (UK) Ltd
14fc  Quadrics Ltd
	0000  QsNet Elan3 Network Adapter
	0001  QsNetII Elan4 Network Adapter
	0002  QsNetIII Elan5 Network Adapter
14fd  JAPAN Computer Industry Inc
14fe  ARCHTEK TELECOM Corp
14ff  TWINHEAD INTERNATIONAL Corp
1500  DELTA Electronics, Inc
	1360  RTL81xx RealTek Ethernet
1501  BANKSOFT CANADA Ltd
1502  MITSUBISHI ELECTRIC LOGISTICS SUPPORT Co Ltd
1503  KAWASAKI LSI USA Inc
1504  KAISER Electronics
1505  ITA INGENIEURBURO FUR TESTAUFGABEN GmbH
1506  CHAMELEON Systems Inc
1507  Motorola ?? / HTEC
	0001  MPC105 [Eagle]
	0002  MPC106 [Grackle]
	0003  MPC8240 [Kahlua]
	0100  MC145575 [HFC-PCI]
	0431  KTI829c 100VG
	4801  Raven
	4802  Falcon
	4803  Hawk
	4806  CPX8216
1508  HONDA CONNECTORS/MHOTRONICS Inc
1509  FIRST INTERNATIONAL Computer Inc
150a  FORVUS RESEARCH Inc
150b  YAMASHITA Systems Corp
150c  KYOPAL CO Ltd
150d  WARPSPPED Inc
150e  C-PORT Corp
150f  INTEC GmbH
1510  BEHAVIOR TECH Computer Corp
1511  CENTILLIUM Technology Corp
1512  ROSUN Technologies Inc
1513  Raychem
1514  TFL LAN Inc
1515  Advent design
1516  MYSON Technology Inc
	0800  MTD-8xx 100/10M Ethernet PCI Adapter
	0803  SURECOM EP-320X-S 100/10M Ethernet PCI Adapter
	0891  MTD-8xx 100/10M Ethernet PCI Adapter
1517  ECHOTEK Corp
1518  Kontron
1519  TELEFON AKTIEBOLAGET LM Ericsson
151a  Globetek
	1002  PCI-1002
	1004  PCI-1004
	1008  PCI-1008
151b  COMBOX Ltd
151c  DIGITAL AUDIO LABS Inc
	0003  Prodif T 2496
	4000  Prodif 88
151d  Fujitsu Computer Products Of America
151e  MATRIX Corp
151f  TOPIC SEMICONDUCTOR Corp
	0000  TP560 Data/Fax/Voice 56k modem
1520  CHAPLET System Inc
1521  BELL Corp
1522  MainPine Ltd
	0100  PCI <-> IOBus Bridge
	4000  PCI Express UART
1523  MUSIC Semiconductors
1524  ENE Technology Inc
	0510  CB710 Memory Card Reader Controller
	0520  FLASH memory: ENE Technology Inc:
	0530  ENE PCI Memory Stick Card Reader Controller
	0550  ENE PCI Secure Digital Card Reader Controller
	0551  SD/MMC Card Reader Controller
	0610  PCI Smart Card Reader Controller
	0720  Memory Stick Card Reader Controller
	0730  ENE PCI Memory Stick Card Reader Controller
	0750  ENE PCI SmartMedia / xD Card Reader Controller
	0751  ENE PCI Secure Digital / MMC Card Reader Controller
	1211  CB1211 Cardbus Controller
	1225  CB1225 Cardbus Controller
	1410  CB1410 Cardbus Controller
	1411  CB-710/2/4 Cardbus Controller
	1412  CB-712/4 Cardbus Controller
	1420  CB1420 Cardbus Controller
	1421  CB-720/2/4 Cardbus Controller
	1422  CB-722/4 Cardbus Controller
1525  IMPACT Technologies
1526  ISS, Inc
1527  SOLECTRON
1528  ACKSYS
1529  AMERICAN MICROSystems Inc
152a  QUICKTURN DESIGN Systems
152b  FLYTECH Technology CO Ltd
152c  MACRAIGOR Systems LLC
152d  QUANTA Computer Inc
152e  MELEC Inc
152f  PHILIPS - CRYPTO
1530  ACQIS Technology Inc
1531  CHRYON Corp
1532  ECHELON Corp
	0020  LonWorks PCLTA-20 PCI LonTalk Adapter
1533  BALTIMORE
1534  ROAD Corp
1535  EVERGREEN Technologies Inc
1536  ACTIS Computer
1537  DATALEX COMMUNCATIONS
1538  ARALION Inc
	0303  ARS106S Ultra ATA 133/100/66 Host Controller
1539  ATELIER INFORMATIQUES et ELECTRONIQUE ETUDES S.A.
153a  ONO SOKKI
153b  TERRATEC Electronic GmbH
	1144  Aureon 5.1
	1147  Aureon 5.1 Sky
	1158  Philips Semiconductors SAA7134 (rev 01) [Terratec Cinergy 600 TV]
153c  ANTAL Electronic
153d  FILANET Corp
153e  TECHWELL Inc
153f  MIPS Technologies, Inc.
	0001  SOC-it 101 System Controller
1540  PROVIDEO MULTIMEDIA Co Ltd
1541  MACHONE Communications
1542  Concurrent Real-Time
	9260  RCIM-II Real-Time Clock & Interrupt Module
	9271  RCIM-III Real-Time Clock & Interrupt Module (PCIe)
	9272  Pulse Width Modulator Card
	9277  5 Volt Delta Sigma Converter Card
	9278  10 Volt Delta Sigma Converter Card
	9287  Analog Output Card
	9290  FPGA Card
1543  SILICON Laboratories
	3052  Intel 537 [Winmodem]
	4c22  Si3036 MC'97 DAA
1544  DCM DATA Systems
1545  VISIONTEK
1546  IOI Technology Corp
1547  MITUTOYO Corp
1548  JET PROPULSION Laboratory
1549  INTERCONNECT Systems Solutions
154a  MAX Technologies Inc
154b  COMPUTEX Co Ltd
154c  VISUAL Technology Inc
154d  PAN INTERNATIONAL Industrial Corp
154e  SERVOTEST Ltd
154f  STRATABEAM Technology
1550  OPEN NETWORK Co Ltd
1551  SMART Electronic DEVELOPMENT GmBH
1552  RACAL AIRTECH Ltd
1553  CHICONY Electronics Co Ltd
1554  PROLINK Microsystems Corp
1555  GESYTEC GmBH
1556  PLDA
	1100  PCI Express Core Reference Design
	110f  PCI Express Core Reference Design Virtual Function
1557  MEDIASTAR Co Ltd
1558  CLEVO/KAPOK Computer
1559  SI LOGIC Ltd
155a  INNOMEDIA Inc
155b  PROTAC INTERNATIONAL Corp
155c  Cemax-Icon Inc
155d  Mac System Co Ltd
155e  LP Elektronik GmbH
155f  Perle Systems Ltd
1560  Terayon Communications Systems
1561  Viewgraphics Inc
1562  Symbol Technologies
1563  A-Trend Technology Co Ltd
1564  Yamakatsu Electronics Industry Co Ltd
1565  Biostar Microtech Int'l Corp
1566  Ardent Technologies Inc
1567  Jungsoft
1568  DDK Electronics Inc
1569  Palit Microsystems Inc.
156a  Avtec Systems
156b  2wire Inc
156c  Vidac Electronics GmbH
156d  Alpha-Top Corp
156e  Alfa Inc
156f  M-Systems Flash Disk Pioneers Ltd
1570  Lecroy Corp
1571  Contemporary Controls
	a001  CCSI PCI20-485 ARCnet
	a002  CCSI PCI20-485D ARCnet
	a003  CCSI PCI20-485X ARCnet
	a004  CCSI PCI20-CXB ARCnet
	a005  CCSI PCI20-CXS ARCnet
	a006  CCSI PCI20-FOG-SMA ARCnet
	a007  CCSI PCI20-FOG-ST ARCnet
	a008  CCSI PCI20-TB5 ARCnet
	a009  CCSI PCI20-5-485 5Mbit ARCnet
	a00a  CCSI PCI20-5-485D 5Mbit ARCnet
	a00b  CCSI PCI20-5-485X 5Mbit ARCnet
	a00c  CCSI PCI20-5-FOG-ST 5Mbit ARCnet
	a00d  CCSI PCI20-5-FOG-SMA 5Mbit ARCnet
	a201  CCSI PCI22-485 10Mbit ARCnet
	a202  CCSI PCI22-485D 10Mbit ARCnet
	a203  CCSI PCI22-485X 10Mbit ARCnet
	a204  CCSI PCI22-CHB 10Mbit ARCnet
	a205  CCSI PCI22-FOG_ST 10Mbit ARCnet
	a206  CCSI PCI22-THB 10Mbit ARCnet
1572  Otis Elevator Company
1573  Lattice - Vantis
1574  Fairchild Semiconductor
1575  Voltaire Advanced Data Security Ltd
1576  Viewcast COM
1578  HITT
	4d34  VPMK4 [Video Processor Mk IV]
	5615  VPMK3 [Video Processor Mk III]
1579  Dual Technology Corp
157a  Japan Elecronics Ind Inc
157b  Star Multimedia Corp
157c  Eurosoft (UK)
	8001  Fix2000 PCI Y2K Compliance Card
157d  Gemflex Networks
157e  Transition Networks
157f  PX Instruments Technology Ltd
1580  Primex Aerospace Co
1581  SEH Computertechnik GmbH
1582  Cytec Corp
1583  Inet Technologies Inc
1584  Uniwill Computer Corp
1585  Logitron
1586  Lancast Inc
1587  Konica Corp
1588  Solidum Systems Corp
1589  Atlantek Microsystems Pty Ltd
	0008  Leutron Vision PicPortExpress CL
	0009  Leutron Vision PicPortExpress CL Stereo
158a  Digalog Systems Inc
158b  Allied Data Technologies
158c  Hitachi Semiconductor & Devices Sales Co Ltd
158d  Point Multimedia Systems
158e  Lara Technology Inc
158f  Ditect Coop
1590  Hewlett Packard Enterprise
	0001  Eagle Cluster Manager
	0002  Osprey Cluster Manager
	0003  Harrier Cluster Manager
	a01d  FC044X Fibre Channel HBA
1591  ARN
1592  Syba Tech Ltd
	0781  Multi-IO Card
	0782  Parallel Port Card 2xEPP
	0783  Multi-IO Card
	0785  Multi-IO Card
	0786  Multi-IO Card
	0787  Multi-IO Card
	0788  Multi-IO Card
	078a  Multi-IO Card
1593  Bops Inc
1594  Netgame Ltd
1595  Diva Systems Corp
1596  Folsom Research Inc
1597  Memec Design Services
1598  Granite Microsystems
1599  Delta Electronics Inc
159a  General Instrument
159b  Faraday Technology Corp
	4321  StorLink SL3516 (Gemini) Host Bridge
159c  Stratus Computer Systems
159d  Ningbo Harrison Electronics Co Ltd
159e  A-Max Technology Co Ltd
159f  Galea Network Security
15a0  Compumaster SRL
15a1  Geocast Network Systems
15a2  Catalyst Enterprises Inc
	0001  TA700 PCI Bus Analyzer/Exerciser
15a3  Italtel
15a4  X-Net OY
15a5  Toyota Macs Inc
15a6  Sunlight Ultrasound Technologies Ltd
15a7  SSE Telecom Inc
15a8  Shanghai Communications Technologies Center
15aa  Moreton Bay
15ab  Bluesteel Networks Inc
15ac  North Atlantic Instruments
	6893  3U OpenVPX Multi-function I/O Board [Model 68C3]
15ad  VMware
	0405  SVGA II Adapter
	0710  SVGA Adapter
	0720  VMXNET Ethernet Controller
	0740  Virtual Machine Communication Interface
	0770  USB2 EHCI Controller
	0774  USB1.1 UHCI Controller
	0778  USB3 xHCI 0.96 Controller
	0779  USB3 xHCI 1.0 Controller
	0790  PCI bridge
	07a0  PCI Express Root Port
	07b0  VMXNET3 Ethernet Controller
	07c0  PVSCSI SCSI Controller
	07e0  SATA AHCI controller
	0801  Virtual Machine Interface
	0820  Paravirtual RDMA controller
	1977  HD Audio Controller
15ae  Amersham Pharmacia Biotech
15b0  Zoltrix International Ltd
15b1  Source Technology Inc
15b2  Mosaid Technologies Inc
15b3  Mellanox Technologies
	0191  MT25408 [ConnectX IB Flash Recovery]
	01f6  MT27500 Family [ConnectX-3 Flash Recovery]
	01f8  MT27520 Family [ConnectX-3 Pro Flash Recovery]
	01ff  MT27600 Family [Connect-IB Flash Recovery]
	0209  MT27700 Family [ConnectX-4 Flash Recovery]
	020b  MT27710 Family [ConnectX-4 Lx Flash Recovery]
	020d  MT28800 Family [ConnectX-5 Flash Recovery]
	020f  MT28908A0 Family [ConnectX-6 Flash Recovery]
	0211  MT416842 Family [BlueField SoC Flash Recovery]
	024e  MT53100 [Spectrum-2, Flash recovery mode]
	024f  MT53100 [Spectrum-2, Flash recovery mode]
	0262  MT27710 [ConnectX-4 Lx Programmable] EN
	0263  MT27710 [ConnectX-4 Lx Programmable Virtual Function] EN
	0281  NPS-600 Flash Recovery
	1002  MT25400 Family [ConnectX-2 Virtual Function]
	1003  MT27500 Family [ConnectX-3]
	1004  MT27500/MT27520 Family [ConnectX-3/ConnectX-3 Pro Virtual Function]
	1005  MT27510 Family
	1006  MT27511 Family
	1007  MT27520 Family [ConnectX-3 Pro]
	1009  MT27530 Family
	100a  MT27531 Family
	100b  MT27540 Family
	100c  MT27541 Family
	100d  MT27550 Family
	100e  MT27551 Family
	100f  MT27560 Family
	1010  MT27561 Family
	1011  MT27600 [Connect-IB]
	1012  MT27600 Family [Connect-IB Virtual Function]
	1013  MT27700 Family [ConnectX-4]
	1014  MT27700 Family [ConnectX-4 Virtual Function]
	1015  MT27710 Family [ConnectX-4 Lx]
	1016  MT27710 Family [ConnectX-4 Lx Virtual Function]
	1017  MT27800 Family [ConnectX-5]
	1018  MT27800 Family [ConnectX-5 Virtual Function]
	1019  MT28800 Family [ConnectX-5 Ex]
	101a  MT28800 Family [ConnectX-5 Ex Virtual Function]
	101b  MT28908 Family [ConnectX-6]
	101c  MT28908 Family [ConnectX-6 Virtual Function]
	101d  MT28841
	101e  MT28850
	101f  MT28851
	1020  MT28860
	1021  MT28861
	1974  MT28800 Family [ConnectX-5 PCIe Bridge]
	1975  MT416842 Family [BlueField SoC PCIe Bridge]
	5274  MT21108 InfiniBridge
	5a44  MT23108 InfiniHost
	5a45  MT23108 [Infinihost HCA Flash Recovery]
	5a46  MT23108 PCI Bridge
	5e8c  MT24204 [InfiniHost III Lx HCA]
	5e8d  MT25204 [InfiniHost III Lx HCA Flash Recovery]
	6274  MT25204 [InfiniHost III Lx HCA]
	6278  MT25208 InfiniHost III Ex (Tavor compatibility mode)
	6279  MT25208 [InfiniHost III Ex HCA Flash Recovery]
	6282  MT25208 [InfiniHost III Ex]
	6340  MT25408 [ConnectX VPI - IB SDR / 10GigE]
	634a  MT25418 [ConnectX VPI PCIe 2.0 2.5GT/s - IB DDR / 10GigE]
	6368  MT25448 [ConnectX EN 10GigE, PCIe 2.0 2.5GT/s]
	6372  MT25408 [ConnectX EN 10GigE 10GBaseT, PCIe 2.0 2.5GT/s]
	6732  MT26418 [ConnectX VPI PCIe 2.0 5GT/s - IB DDR / 10GigE]
	673c  MT26428 [ConnectX VPI PCIe 2.0 5GT/s - IB QDR / 10GigE]
	6746  MT26438 [ConnectX VPI PCIe 2.0 5GT/s - IB QDR / 10GigE Virtualization+]
	6750  MT26448 [ConnectX EN 10GigE, PCIe 2.0 5GT/s]
	675a  MT25408 [ConnectX EN 10GigE 10GBaseT, PCIe Gen2 5GT/s]
	6764  MT26468 [ConnectX EN 10GigE, PCIe 2.0 5GT/s Virtualization+]
	676e  MT26478 [ConnectX EN 40GigE, PCIe 2.0 5GT/s]
	6778  MT26488 [ConnectX VPI PCIe 2.0 5GT/s - IB DDR / 10GigE Virtualization+]
	7101  NPS-400 configuration and management interface
	7102  NPS-400 network interface PF
	7103  NPS-400 network interface VF
	7121  NPS-600 configuration and management interface
	7122  NPS-600 network interface PF
	7123  NPS-600 network interface VF
	a2d0  MT416842 BlueField SoC Crypto enabled
	a2d1  MT416842 BlueField SoC Crypto disabled
	a2d2  MT416842 BlueField integrated ConnectX-5 network controller
	a2d3  MT416842 BlueField multicore SoC family VF
	c738  MT51136
	c739  MT51136 GW
	c838  MT52236
	c839  MT52236 router
	caf1  ConnectX-4 CAPI Function
	cb84  MT52100
	cf08  MT53236
	cf6c  MT53100 [Spectrum-2, 64 x 100GbE switch]
	d2f0  Switch-IB 3 HDR (200Gbps) switch
15b4  CCI/TRIAD
15b5  Cimetrics Inc
15b6  Texas Memory Systems Inc
	0001  XP15 DSP Accelerator
	0002  XP30 DSP Accelerator
	0003  XP00 Data Acquisition Device
	0004  XP35 DSP Accelerator
	0007  XP100 DSP Accelerator [XP100-T0]
	0008  XP100 DSP Accelerator [XP100-T1]
	0009  XP100 DSP Accelerator [XP100-E0]
	000a  XP100 DSP Accelerator [XP100-E1]
	000e  XP100 DSP Accelerator [XP100-0]
	000f  XP100 DSP Accelerator [XP100-1]
	0010  XP100 DSP Accelerator [XP100-P0]
	0011  XP100 DSP Accelerator [XP100-P1]
	0012  XP100 DSP Accelerator [XP100-P2]
	0013  XP100 DSP Accelerator [XP100-P3]
	0014  RamSan Flash SSD
	0015  ZBox
15b7  Sandisk Corp
	2001  Skyhawk Series NVME SSD
15b8  ADDI-DATA GmbH
	1001  APCI1516 SP controller (16 digi outputs)
	1003  APCI1032 SP controller (32 digi inputs w/ opto coupler)
	1004  APCI2032 SP controller (32 digi outputs)
	1005  APCI2200 SP controller (8/16 digi outputs (relay))
	1006  APCI1564 SP controller (32 digi ins, 32 digi outs)
	100a  APCI1696 SP controller (96 TTL I/Os)
	3001  APCI3501 SP controller (analog output board)
	300f  APCI3600 Noise and vibration measurement board
	7001  APCI7420 2-port Serial Controller
	7002  APCI7300 Serial Controller
15b9  Maestro Digital Communications
15ba  Impacct Technology Corp
15bb  Portwell Inc
15bc  Agilent Technologies
	0100  HPFC-5600 Tachyon DX2+ FC
	0103  QX4 PCI Express quad 4-gigabit Fibre Channel controller
	0105  Celerity FC-44XS/FC-42XS/FC-41XS/FC-44ES/FC-42ES/FC-41ES
	1100  E8001-66442 PCI Express CIC
	2922  64 Bit, 133MHz PCI-X Exerciser & Protocol Checker
	2928  64 Bit, 66MHz PCI Exerciser & Analyzer
	2929  64 Bit, 133MHz PCI-X Analyzer & Exerciser
15bd  DFI Inc
15be  Sola Electronics
15bf  High Tech Computer Corp (HTC)
15c0  BVM Ltd
15c1  Quantel
15c2  Newer Technology Inc
15c3  Taiwan Mycomp Co Ltd
15c4  EVSX Inc
15c5  Procomp Informatics Ltd
	8010  1394b - 1394 Firewire 3-Port Host Adapter Card
15c6  Technical University of Budapest
15c7  Tateyama System Laboratory Co Ltd
	0349  Tateyama C-PCI PLC/NC card Rev.01A
15c8  Penta Media Co Ltd
15c9  Serome Technology Inc
15ca  Bitboys OY
15cb  AG Electronics Ltd
15cc  Hotrail Inc
15cd  Dreamtech Co Ltd
15ce  Genrad Inc
15cf  Hilscher GmbH
	0000  CIFX 50E-DP(M/S)
15d1  Infineon Technologies AG
15d2  FIC (First International Computer Inc)
15d3  NDS Technologies Israel Ltd
15d4  Iwill Corp
15d5  Tatung Co
15d6  Entridia Corp
15d7  Rockwell-Collins Inc
15d8  Cybernetics Technology Co Ltd
15d9  Super Micro Computer Inc
15da  Cyberfirm Inc
15db  Applied Computing Systems Inc
15dc  Litronic Inc
	0001  Argus 300 PCI Cryptography Module
15dd  Sigmatel Inc
15de  Malleable Technologies Inc
15df  Infinilink Corp
15e0  Cacheflow Inc
15e1  Voice Technologies Group Inc
15e2  Quicknet Technologies Inc
	0500  PhoneJack-PCI
15e3  Networth Technologies Inc
15e4  VSN Systemen BV
15e5  Valley technologies Inc
15e6  Agere Inc
15e7  Get Engineering Corp
15e8  National Datacomm Corp
	0130  Wireless PCI Card
	0131  NCP130A2 Wireless NIC
15e9  Pacific Digital Corp
	1841  ADMA-100 DiscStaQ ATA Controller
15ea  Tokyo Denshi Sekei K.K.
15eb  DResearch Digital Media Systems GmbH
15ec  Beckhoff GmbH
	3101  FC3101 Profibus DP 1 Channel PCI
	5102  FC5102
15ed  Macrolink Inc
15ee  In Win Development Inc
15ef  Intelligent Paradigm Inc
15f0  B-Tree Systems Inc
15f1  Times N Systems Inc
15f2  Diagnostic Instruments Inc
15f3  Digitmedia Corp
15f4  Valuesoft
15f5  Power Micro Research
15f6  Extreme Packet Device Inc
15f7  Banctec
15f8  Koga Electronics Co
15f9  Zenith Electronics Corp
15fa  J.P. Axzam Corp
15fb  Zilog Inc
15fc  Techsan Electronics Co Ltd
15fd  N-CUBED.NET
15fe  Kinpo Electronics Inc
15ff  Fastpoint Technologies Inc
1600  Northrop Grumman - Canada Ltd
1601  Tenta Technology
1602  Prosys-tec Inc
1603  Nokia Wireless Communications
1604  Central System Research Co Ltd
1605  Pairgain Technologies
1606  Europop AG
1607  Lava Semiconductor Manufacturing Inc
1608  Automated Wagering International
1609  Scimetric Instruments Inc
1612  Telesynergy Research Inc.
1618  Stone Ridge Technology
	0001  RDX 11
	0002  HFT-01
	0400  FarSync T2P (2 port X.21/V.35/V.24)
	0440  FarSync T4P (4 port X.21/V.35/V.24)
	0610  FarSync T1U (1 port X.21/V.35/V.24)
	0620  FarSync T2U (2 port X.21/V.35/V.24)
	0640  FarSync T4U (4 port X.21/V.35/V.24)
	1610  FarSync TE1 (T1,E1)
	2610  FarSync DSL-S1 (SHDSL)
	3640  FarSync T4E (4-port X.21/V.35/V.24)
	4620  FarSync T2Ue PCI Express (2-port X.21/V.35/V.24)
	4640  FarSync T4Ue PCI Express (4-port X.21/V.35/V.24)
1619  FarSite Communications Ltd
	0400  FarSync T2P (2 port X.21/V.35/V.24)
	0440  FarSync T4P (4 port X.21/V.35/V.24)
	0610  FarSync T1U (1 port X.21/V.35/V.24)
	0620  FarSync T2U (2 port X.21/V.35/V.24)
	0640  FarSync T4U (4 port X.21/V.35/V.24)
	1610  FarSync TE1 (T1,E1)
	1612  FarSync TE1 PCI Express (T1,E1)
	2610  FarSync DSL-S1 (SHDSL)
	3640  FarSync T4E (4-port X.21/V.35/V.24)
	4620  FarSync T2Ue PCI Express (2-port X.21/V.35/V.24)
	4640  FarSync T4Ue PCI Express (4-port X.21/V.35/V.24)
	5621  FarSync T2Ee PCI Express (2 port X.21/V.35/V.24)
	5641  FarSync T4Ee PCI Express (4 port X.21/V.35/V.24)
	6620  FarSync T2U-PMC PCI Express (2 port X.21/V.35/V.24)
161f  Rioworks
1626  TDK Semiconductor Corp.
	8410  RTL81xx Fast Ethernet
1629  Kongsberg Spacetec AS
	1003  Format synchronizer v3.0
	1006  Format synchronizer, model 10500
	1007  Format synchronizer, model 21000
	2002  Fast Universal Data Output
1631  Packard Bell B.V.
1638  Standard Microsystems Corp [SMC]
	1100  SMC2602W EZConnect / Addtron AWA-100 / Eumitcom PCI WL11000
163c  Smart Link Ltd.
	3052  SmartLink SmartPCI562 56K Modem
	5449  SmartPCI561 Modem
1641  MKNet Corp.
1642  Bitland(ShenZhen) Information Technology Co., Ltd.
1657  Brocade Communications Systems, Inc.
	0013  425/825/42B/82B 4Gbps/8Gbps PCIe dual port FC HBA
	0014  1010/1020/1007/1741 10Gbps CNA
	0017  415/815/41B/81B 4Gbps/8Gbps PCIe single port FC HBA
	0021  804 8Gbps FC HBA for HP Bladesystem c-class
	0022  1860 16Gbps/10Gbps Fabric Adapter
	0023  1867/1869 16Gbps FC HBA
	0646  400 4Gbps PCIe FC HBA
165a  Epix Inc
	c100  PIXCI(R) CL1 Camera Link Video Capture Board [custom QL5232]
	d200  PIXCI(R) D2X Digital Video Capture Board [custom QL5232]
	d300  PIXCI(R) D3X Digital Video Capture Board [custom QL5232]
	eb01  PIXCI(R) EB1 PCI Camera Link Video Capture Board
165c  Gidel Ltd.
	5361  PROCStarII60-1
	5362  PROCStarII60-2
	5364  PROCStarII60-4
	5435  ProcSparkII
	5661  ProcE60
	56e1  ProcE180
	5911  ProcStarIII110-1
	5912  ProcStarIII110-2
	5913  ProcStarIII110-3
	5914  ProcStarIII110-4
	5921  ProcStarIII150-1
	5922  ProcStarIII150-2
	5923  ProcStarIII150-3
	5924  ProcStarIII150-4
	5931  ProcStarIII260-1
	5932  ProcStarIII260-2
	5933  ProcStarIII260-3
	5934  ProcStarIII260-4
	5941  ProcStarIII340-1
	5942  ProcStarIII340-2
	5943  ProcStarIII340-3
	5944  ProcStarIII340-4
	5a01  ProceIII80
	5a11  ProceIII110
	5a21  ProceIII150
	5a31  ProceIII260
	5a41  ProceIII340
	5b51  ProceIV360
	5b61  ProceIV530
	5b71  ProceIV820
	5c01  ProcStarIV80-1
	5c02  ProcStarIV80-2
	5c03  ProcStarIV80-3
	5c04  ProcStarIV80-4
	5c11  ProcStarIV110-1
	5c12  ProcStarIV110-2
	5c13  ProcStarIV110-3
	5c14  ProcStarIV110-4
	5c51  ProcStarIV360-1
	5c52  ProcStarIV360-2
	5c53  ProcStarIV360-3
	5c54  ProcStarIV360-4
	5c61  ProcStarIV530-1
	5c62  ProcStarIV530-2
	5c63  ProcStarIV530-3
	5c64  ProcStarIV530-4
	5c71  ProcStarIV820-1
	5c72  ProcStarIV820-2
	5c73  ProcStarIV820-3
	5c74  ProcStarIV820-4
	5d01  Proc10480
	5d11  Proc104110
	5f01  ProceV_A3
	5f11  ProceV_A7
	5f21  ProceV_AB
	5f31  ProceV_D5
	5f41  ProceV_D8
	6732  Proc6M
	6832  Proc12M
	7101  Proc10a_27
	7111  Proc10a_48
	7121  Proc10a_66
	7141  Proc10a_115
	7181  Proc10a_27S
	7191  Proc10a_48S
	71a1  Proc10a_66S
	71b1  Proc10A
165d  Hsing Tech. Enterprise Co., Ltd.
165f  Linux Media Labs, LLC
	1020  LMLM4 MPEG-4 encoder
1661  Worldspace Corp.
1668  Actiontec Electronics Inc
	0100  Mini-PCI bridge
166d  Broadcom Corporation
	0001  SiByte BCM1125/1125H/1250 System-on-a-Chip PCI
	0002  SiByte BCM1125H/1250 System-on-a-Chip HyperTransport
	0012  SiByte BCM1280/BCM1480 System-on-a-Chip PCI-X
	0014  Sibyte BCM1280/BCM1480 System-on-a-Chip HyperTransport
1677  Bernecker + Rainer
	104e  5LS172.6 B&R Dual CAN Interface Card
	12d7  5LS172.61 B&R Dual CAN Interface Card
	20ad  5ACPCI.MFIO-K01 Profibus DP / K-Feldbus / COM
1678  NetEffect
	0100  NE020 10Gb Accelerated Ethernet Adapter (iWARP RNIC)
1679  Tokyo Electron Device Ltd.
	3000  SD Standard host controller [Ellen]
167b  ZyDAS Technology Corp.
	2102  ZyDAS ZD1202
	2116  ZD1212B Wireless Adapter
167d  Samsung Electro-Mechanics Co., Ltd.
	a000  MagicLAN SWL-2210P 802.11b [Intersil ISL3874]
167e  ONNTO Corp.
1681  Hercules
1682  XFX Pine Group Inc.
1688  CastleNet Technology Inc.
	1170  WLAN 802.11b card
168c  Qualcomm Atheros
	0007  AR5210 Wireless Network Adapter [AR5000 802.11a]
	0011  AR5211 Wireless Network Adapter [AR5001A 802.11a]
	0012  AR5211 Wireless Network Adapter [AR5001X 802.11ab]
	0013  AR5212/5213/2414 Wireless Network Adapter
	001a  AR2413/AR2414 Wireless Network Adapter [AR5005G(S) 802.11bg]
	001b  AR5413/AR5414 Wireless Network Adapter [AR5006X(S) 802.11abg]
	001c  AR242x / AR542x Wireless Network Adapter (PCI-Express)
	001d  AR2417 Wireless Network Adapter [AR5007G 802.11bg]
	0020  AR5513 802.11abg Wireless NIC
	0023  AR5416 Wireless Network Adapter [AR5008 802.11(a)bgn]
	0024  AR5418 Wireless Network Adapter [AR5008E 802.11(a)bgn] (PCI-Express)
	0027  AR9160 Wireless Network Adapter [AR9001 802.11(a)bgn]
	0029  AR922X Wireless Network Adapter
	002a  AR928X Wireless Network Adapter (PCI-Express)
	002b  AR9285 Wireless Network Adapter (PCI-Express)
	002c  AR2427 802.11bg Wireless Network Adapter (PCI-Express)
	002d  AR9227 Wireless Network Adapter
	002e  AR9287 Wireless Network Adapter (PCI-Express)
	0030  AR93xx Wireless Network Adapter
	0032  AR9485 Wireless Network Adapter
	0033  AR958x 802.11abgn Wireless Network Adapter
	0034  AR9462 Wireless Network Adapter
	0036  QCA9565 / AR9565 Wireless Network Adapter
	0037  AR9485 Wireless Network Adapter
	003c  QCA986x/988x 802.11ac Wireless Network Adapter
	003e  QCA6174 802.11ac Wireless Network Adapter
	0040  QCA9980/9990 802.11ac Wireless Network Adapter
	0041  QCA6164 802.11ac Wireless Network Adapter
	0042  QCA9377 802.11ac Wireless Network Adapter
	0050  QCA9887 802.11ac Wireless Network Adapter
	0207  AR5210 Wireless Network Adapter [AR5000 802.11a]
	1014  AR5212 802.11abg NIC
	9013  AR5002X Wireless Network Adapter
	ff19  AR5006X Wireless Network Adapter
	ff1b  AR2425 Wireless Network Adapter [AR5007EG 802.11bg]
	ff1c  AR5008 Wireless Network Adapter
	ff1d  AR922x Wireless Network Adapter
1695  EPoX Computer Co., Ltd.
169c  Netcell Corporation
	0044  Revolution Storage Processing Card
169d  Club-3D VB (Wrong ID)
16a5  Tekram Technology Co.,Ltd.
16ab  Global Sun Technology Inc
	1100  GL24110P
	1101  PLX9052 PCMCIA-to-PCI Wireless LAN
	1102  PCMCIA-to-PCI Wireless Network Bridge
	8501  WL-8305 Wireless LAN PCI Adapter
16ae  SafeNet Inc
	0001  SafeXcel 1140
	000a  SafeXcel 1841
	1141  SafeXcel 1141
	1841  SafeXcel 1842
16af  SparkLAN Communications, Inc.
16b4  Aspex Semiconductor Ltd
16b8  Sonnet Technologies, Inc.
16be  Creatix Polymedia GmbH
16c3  Synopsys, Inc.
16c6  Micrel-Kendin
	8695  Centaur KS8695 ARM processor
	8842  KSZ8842-PMQL 2-Port Ethernet Switch
16c8  Octasic Inc.
16c9  EONIC B.V. The Netherlands
16ca  CENATEK Inc
	0001  Rocket Drive DL
16cd  Advantech Co. Ltd
	0101  DirectPCI SRAM for DPX-11x series
	0102  DirectPCI SRAM for DPX-S/C/E-series
	0103  DirectPCI ROM for DPX-11x series
	0104  DirectPCI ROM for DPX-S/C/E-series
	0105  DirectPCI I/O for DPX-114/DPX-115
	0106  DirectPCI I/O for DPX-116
	0107  DirectPCI I/O for DPX-116U
	0108  DirectPCI I/O for DPX-117
	0109  DirectPCI I/O for DPX-112
	010a  DirectPCI I/O for DPX-C/E-series
	010b  DirectPCI I/O for DPX-S series
16ce  Roland Corp.
16d5  Acromag, Inc.
	0504  PMC-DX504 Reconfigurable FPGA with LVDS I/O
	0520  PMC520 Serial Communication, 232 Octal
	0521  PMC521 Serial Communication, 422/485 Octal
	1020  PMC-AX1020 Reconfigurable FPGA with A/D & D/A
	1065  PMC-AX1065 Reconfigurable FPGA with A/D & D/A
	2004  PMC-DX2004 Reconfigurable FPGA with LVDS I/O
	2020  PMC-AX2020 Reconfigurable FPGA with A/D & D/A
	2065  PMC-AX2065 Reconfigurable FPGA with A/D & D/A
	3020  PMC-AX3020 Reconfigurable FPGA with A/D & D/A
	3065  PMC-AX3065 Reconfigurable FPGA with A/D & D/A
	4243  PMC424, APC424, AcPC424 Digital I/O and Counter Timer Module
	4248  PMC464, APC464, AcPC464 Digital I/O and Counter Timer Module
	424b  PMC-DX2002 Reconfigurable FPGA with Differential I/O
	4253  PMC-DX503 Reconfigurable FPGA with TTL and Differential I/O
	4312  PMC-CX1002 Reconfigurable Conduction-Cooled FPGA Virtex-II with Differential I/O
	4313  PMC-CX1003 Reconfigurable Conduction-Cooled FPGA Virtex-II with CMOS and Differential I/O
	4322  PMC-CX2002 Reconfigurable Conduction-Cooled FPGA Virtex-II with Differential I/O
	4323  PMC-CX2003 Reconfigurable Conduction-Cooled FPGA Virtex-II with CMOS and Differential I/O
	4350  PMC-DX501 Reconfigurable Digital I/O Module
	4353  PMC-DX2003 Reconfigurable FPGA with TTL and Differential I/O
	4357  PMC-DX502 Reconfigurable Differential I/O Module
	4457  PMC730, APC730, AcPC730 Multifunction Module
	464d  PMC408 32-Channel Digital Input/Output Module
	4850  PMC220-16 12-Bit Analog Output Module
	4a42  PMC483, APC483, AcPC483 Counter Timer Module
	4a50  PMC484, APC484, AcPC484 Counter Timer Module
	4a56  PMC230 16-Bit Analog Output Module
	4b47  PMC330, APC330, AcPC330 Analog Input Module, 16-bit A/D
	4c40  PMC-LX40 Reconfigurable Virtex-4 FPGA with plug-in I/O
	4c60  PMC-LX60 Reconfigurable Virtex-4 FPGA with plug-in I/O
	4d4d  PMC341, APC341, AcPC341 Analog Input Module, Simultaneous Sample & Hold
	4d4e  PMC482, APC482, AcPC482 Counter Timer Board
	524d  PMC-DX2001 Reconfigurable FPGA with TTL I/O
	5335  PMC-SX35 Reconfigurable Virtex-4 FPGA with plug-in I/O
	5456  PMC470 48-Channel Digital Input/Output Module
	5601  PMC-VLX85 Reconfigurable Virtex-5 FPGA with plug-in I/O
	5602  PMC-VLX110 Reconfigurable Virtex-5 FPGA with plug-in I/O
	5603  PMC-VSX95 Reconfigurable Virtex-5 FPGA with plug-in I/O
	5604  PMC-VLX155 Reconfigurable Virtex-5 FPGA with plug-in I/O
	5605  PMC-VFX70 Reconfigurable Virtex-5 FPGA with plug-in I/O
	5606  PMC-VLX155-1M Reconfigurable Virtex-5 FPGA with plug-in I/O
	5701  PMC-SLX150: Reconfigurable Spartan-6 FPGA with plug-in I/O
	5702  PMC-SLX150-1M: Reconfigurable Spartan-6 FPGA with plug-in I/O
	5801  XMC-VLX85 Reconfigurable Virtex-5 FPGA with plug-in I/O
	5802  XMC-VLX110 Reconfigurable Virtex-5 FPGA with plug-in I/O
	5803  XMC-VSX95 Reconfigurable Virtex-5 FPGA with plug-in I/O
	5804  XMC-VLX155 Reconfigurable Virtex-5 FPGA with plug-in I/O
	5807  XMC-SLX150: Reconfigurable Spartan-6 FPGA with plug-in I/O
	5808  XMC-SLX150-1M: Reconfigurable Spartan-6 FPGA with plug-in I/O
	5901  APCe8650 PCI Express IndustryPack Carrier Card
	6301  XMC Module with user-configurable Virtex-6 FPGA, 240k logic cells, SFP front I/O
	6302  XMC Module with user-configurable Virtex-6 FPGA, 365k logic cells, SFP front I/O
	6303  XMC Module with user-configurable Virtex-6 FPGA, 240k logic cells, no front I/O
	6304  XMC Module with user-configurable Virtex-6 FPGA, 365k logic cells, no front I/O
	7000  XMC-7K325F: User-configurable Kintex-7 FPGA, 325k logic cells plus SFP front I/O
	7001  XMC-7K410F: User-configurable Kintex-7 FPGA, 410k logic cells plus SFP front I/O
	7002  XMC-7K325AX: User-Configurable Kintex-7 FPGA, 325k logic cells with AXM Plug-In I/O
	7003  XMC-7K410AX: User-Configurable Kintex-7 FPGA, 410k logic cells with AXM Plug-In I/O
	7004  XMC-7K325CC: User-Configurable Kintex-7 FPGA, 325k logic cells, conduction-cooled
	7005  XMC-7K410CC: User-Configurable Kintex-7 FPGA, 410k logic cells, conduction-cooled
	7006  XMC-7A200: User-Configurable Artix-7 FPGA, 200k logic cells with Plug-In I/O
	7007  XMC-7A200CC: User-Configurable Conduction-Cooled Artix-7 FPGA, with 200k logic cells
	7011  AP440-1: 32-Channel Isolated Digital Input Module
	7012  AP440-2: 32-Channel Isolated Digital Input Module
	7013  AP440-3: 32-Channel Isolated Digital Input Module
	7014  AP445: 32-Channel Isolated Digital Output Module
	7016  AP470 48-Channel TTL Level Digital Input/Output Module
	7017  AP323 16-bit, 20 or 40 Channel Analog Input Module
	7018  AP408: 32-Channel Digital I/O Module
	7019  AP341 14-bit, 16-Channel Simultaneous Conversion Analog Input Module
	701a  AP220-16 12-Bit, 16-Channel Analog Output Module
	701b  AP231-16 16-Bit, 16-Channel Analog Output Module
	7021  APA7-201 Reconfigurable Artix-7 FPGA module 48 TTL channels
	7022  APA7-202 Reconfigurable Artix-7 FPGA module 24 RS485 channels
	7023  APA7-203 Reconfigurable Artix-7 FPGA module 24 TTL & 12 RS485 channels
	7024  APA7-204 Reconfigurable Artix-7 FPGA module 24 LVDS channels
	7027  AP418 16-Channel High Voltage Digital Input/Output Module
	7042  AP482 Counter Timer Module with TTL Level Input/Output
	7043  AP483 Counter Timer Module with TTL Level and RS422 Input/Output
	7044  AP484 Counter Timer Module with RS422 Input/Output
16da  Advantech Co., Ltd.
	0011  INES GPIB-PCI
16df  PIKA Technologies Inc.
16e2  Geotest-MTS
16e3  European Space Agency
	1e0f  LEON2FT Processor
16e5  Intellon Corp.
	6000  INT6000 Ethernet-to-Powerline Bridge [HomePlug AV]
	6300  INT6300 Ethernet-to-Powerline Bridge [HomePlug AV]
16ec  U.S. Robotics
	00ed  USR997900
	0116  USR997902 10/100/1000 Mbps PCI Network Card
	2f00  USR5660A (USR265660A, USR5660A-BP) 56K PCI Faxmodem
	3685  Wireless Access PCI Adapter Model 022415
	4320  USR997904 10/100/1000 64-bit NIC (Marvell Yukon)
	ab06  USR997901A 10/100 Cardbus NIC
16ed  Sycron N. V.
	1001  UMIO communication card
16f2  ETAS GmbH
	0200  I/O board
16f3  Jetway Information Co., Ltd.
16f4  Vweb Corp
	8000  VW2010
16f6  VideoTele.com, Inc.
1702  Internet Machines Corporation (IMC)
1705  Digital First, Inc.
170b  NetOctave
	0100  NSP2000-SSL crypto accelerator
170c  YottaYotta Inc.
1719  EZChip Technologies
	1000  NPA Access Network Processor Family
1725  Vitesse Semiconductor
	7174  VSC7174 PCI/PCI-X Serial ATA Host Bus Controller
172a  Accelerated Encryption
	13c8  AEP SureWare Runner 1000V3
1734  Fujitsu Technology Solutions
1735  Aten International Co. Ltd.
1737  Linksys
	0029  WPG54G ver. 4 PCI Card
	1032  Gigabit Network Adapter
	1064  Gigabit Network Adapter
	ab08  21x4x DEC-Tulip compatible 10/100 Ethernet
	ab09  21x4x DEC-Tulip compatible 10/100 Ethernet
173b  Altima (nee Broadcom)
	03e8  AC1000 Gigabit Ethernet
	03e9  AC1001 Gigabit Ethernet
	03ea  AC9100 Gigabit Ethernet
	03eb  AC1003 Gigabit Ethernet
1743  Peppercon AG
	8139  ROL/F-100 Fast Ethernet Adapter with ROL
1745  ViXS Systems, Inc.
	2020  XCode II Series
	2100  XCode 2100 Series
1749  RLX Technologies
174b  PC Partner Limited / Sapphire Technology
174d  WellX Telecom SA
175c  AudioScience Inc
175e  Sanera Systems, Inc.
1760  TEDIA spol. s r. o.
	0101  PCD-7004 Digital Bi-Directional Ports PCI Card
	0102  PCD-7104 Digital Input & Output PCI Card
	0303  PCD-7006C Digital Input & Output PCI Card
1771  InnoVISION Multimedia Ltd.
1775  GE Intelligent Platforms
177d  Cavium, Inc.
	0001  Nitrox XL N1
	0003  Nitrox XL N1 Lite
	0004  Octeon (and older) FIPS
	0005  Octeon CN38XX Network Processor Pass 3.x
	0006  RoHS
	0010  Nitrox XL NPX
	0020  Octeon CN31XX Network Processor
	0030  Octeon CN30XX Network Processor
	0040  Octeon CN58XX Network Processor
	0050  Octeon CN57XX Network Processor (CN54XX/CN55XX/CN56XX)
	0070  Octeon CN50XX Network Processor
	0080  Octeon CN52XX Network Processor
	0090  Octeon II CN63XX Network Processor
	0091  Octeon II CN68XX Network Processor
	0092  Octeon II CN65XX Network Processor
	0093  Octeon II CN61XX Network Processor
	0094  Octeon Fusion CNF71XX Cell processor
	0095  Octeon III CN78XX Network Processor
	0096  Octeon III CN70XX Network Processor
	9700  Octeon III CN73XX Network Processor
	9702  CN23XX [LiquidIO II] Intelligent Adapter
	9703  CN23XX [LiquidIO II] NVMe Controller
	9712  CN23XX [LiquidIO II] SRIOV Virtual Function
	9713  CN23XX [LiquidIO II] NVMe SRIOV Virtual Function
	9800  Octeon Fusion CNF75XX Processor
	a001  ThunderX MRML(Master RML Bridge to RSL devices)
	a002  THUNDERX PCC Bridge
	a008  THUNDERX SMMU
	a009  THUNDERX Generic Interrupt Controller
	a00a  THUNDERX GPIO Controller
	a00b  THUNDERX MPI / SPI Controller
	a00c  THUNDERX MIO-PTP Controller
	a00d  THUNDERX MIX Network Controller
	a00e  THUNDERX Reset Controller
	a00f  THUNDERX UART Controller
	a010  THUNDERX eMMC/SD Controller
	a011  THUNDERX MIO-BOOT Controller
	a012  THUNDERX TWSI / I2C Controller
	a013  THUNDERX CCPI (Multi-node connect)
	a014  THUNDERX Voltage Regulator Module
	a015  THUNDERX PCIe Switch Logic Interface
	a016  THUNDERX Key Memory
	a017  THUNDERX GTI (Global System Timers)
	a018  THUNDERX Random Number Generator
	a019  THUNDERX DFA
	a01a  THUNDERX Zip Coprocessor
	a01b  THUNDERX xHCI USB Controller
	a01c  THUNDERX AHCI SATA Controller
	a01d  THUNDERX RAID Coprocessor
	a01e  THUNDERX Network Interface Controller
	a01f  THUNDERX Traffic Network Switch
	a020  THUNDERX PEM (PCI Express Interface)
	a021  THUNDERX L2C (Level-2 Cache Controller)
	a022  THUNDERX LMC (DRAM Controller)
	a023  THUNDERX OCLA (On-Chip Logic Analyzer)
	a024  THUNDERX OSM
	a025  THUNDERX GSER (General Serializer/Deserializer)
	a026  THUNDERX BGX (Common Ethernet Interface)
	a027  THUNDERX IOBN
	a029  THUNDERX NCSI (Network Controller Sideband Interface)
	a02a  ThunderX SGPIO (Serial GPIO controller for SATA disk lights)
	a02b  THUNDERX SMI / MDIO Controller
	a02c  THUNDERX DAP (Debug Access Port)
	a02d  THUNDERX PCIERC (PCIe Root Complex)
	a02e  ThunderX L2C-TAD (Level 2 cache tag and data)
	a02f  THUNDERX L2C-CBC
	a030  THUNDERX L2C-MCI
	a031  THUNDERX MIO-FUS (Fuse Access Controller)
	a032  THUNDERX FUSF (Fuse Controller)
	a033  THUNDERX Random Number Generator virtual function
	a034  THUNDERX Network Interface Controller virtual function
	a035  THUNDERX Parallel Bus
	a036  ThunderX RAD (RAID acceleration engine) virtual function
	a037  THUNDERX ZIP virtual function
	a040  THUNDERX CPT Cryptographic Accelerator
	a100  THUNDERX CN88XX 48 core SoC
	a200  OCTEON TX CN81XX/CN80XX
	a300  OCTEON TX CN83XX
1787  Hightech Information System Ltd.
1789  Ennyah Technologies Corp.
1796  Research Centre Juelich
	0001  SIS1100 [Gigabit link]
	0002  HOTlink
	0003  Counter Timer
	0004  CAMAC Controller
	0005  PROFIBUS
	0006  AMCC HOTlink
	000d  Synchronisation Slave
	000e  SIS1100-eCMC
	000f  TDC (GPX)
	0010  PCIe Counter Timer
	0011  SIS1100-e single link
	0012  SIS1100-e quad link
	0015  SIS8100 [Gigabit link, MicroTCA]
1797  Intersil Techwell
	5864  TW5864 multimedia video controller
	6801  TW6802 multimedia video card
	6802  TW6802 multimedia other device
	6810  TW6816 multimedia video controller
	6811  TW6816 multimedia video controller
	6812  TW6816 multimedia video controller
	6813  TW6816 multimedia video controller
	6814  TW6816 multimedia video controller
	6815  TW6816 multimedia video controller
	6816  TW6816 multimedia video controller
	6817  TW6816 multimedia video controller
	6864  TW6864 multimedia video controller
1799  Belkin
	6001  F5D6001 Wireless PCI Card [Realtek RTL8180]
	6020  F5D6020 v3000 Wireless PCMCIA Card [Realtek RTL8180]
	6060  F5D6060 Wireless PDA Card
	700f  F5D7000 v7000 Wireless G Desktop Card [Realtek RTL8185]
	701f  F5D7010 v7000 Wireless G Notebook Card [Realtek RTL8185]
179a  id Quantique
	0001  Quantis PCI 16Mbps
179c  Data Patterns
	0557  DP-PCI-557 [PCI 1553B]
	0566  DP-PCI-566 [Intelligent PCI 1553B]
	1152  DP-cPCI-1152 (8-channel Isolated ADC Module)
	5031  DP-CPCI-5031-Synchro Module
	5112  DP-cPCI-5112 [MM-Carrier]
	5121  DP-CPCI-5121-IP Carrier
	5211  DP-CPCI-5211-IP Carrier
	5679  AGE Display Module
17a0  Genesys Logic, Inc
	7163  GL9701 PCIe to PCI Bridge
	8083  GL880 USB 1.1 UHCI controller
	8084  GL880 USB 2.0 EHCI controller
17aa  Lenovo
	402b  Intel 82599ES 10Gb 2-port Server Adapter X520-2
17ab  Phillips Components
17af  Hightech Information System Ltd.
17b3  Hawking Technologies
	ab08  PN672TX 10/100 Ethernet
17b4  Indra Networks, Inc.
	0011  WebEnhance 100 GZIP Compression Card
	0012  WebEnhance 200 GZIP Compression Card
	0015  WebEnhance 300 GZIP Compression Card
	0016  StorCompress 300 GZIP Compression Card
	0017  StorSecure 300 GZIP Compression and AES Encryption Card
17c0  Wistron Corp.
17c2  Newisys, Inc.
17cb  Qualcomm
	0001  AGN100 802.11 a/b/g True MIMO Wireless Card
	0002  AGN300 802.11 a/b/g True MIMO Wireless Card
	0400  Datacenter Technologies QDF2432 PCI Express Root Port
	0401  Datacenter Technologies QDF2400 PCI Express Root Port
17cc  NetChip Technology, Inc
	2280  USB 2.0
17cd  Cadence Design Systems, Inc.
17cf  Z-Com, Inc.
17d3  Areca Technology Corp.
	1110  ARC-1110 4-Port PCI-X to SATA RAID Controller
	1120  ARC-1120 8-Port PCI-X to SATA RAID Controller
	1130  ARC-1130 12-Port PCI-X to SATA RAID Controller
	1160  ARC-1160 16-Port PCI-X to SATA RAID Controller
	1170  ARC-1170 24-Port PCI-X to SATA RAID Controller
	1201  ARC-1200 2-Port PCI-Express to SATA II RAID Controller
	1203  ARC-1203 2/4/8 Port PCIe 2.0 to SATA 6Gb RAID Controller
	1210  ARC-1210 4-Port PCI-Express to SATA RAID Controller
	1214  ARC-12x4 PCIe 2.0 to SAS/SATA 6Gb RAID Controller
	1220  ARC-1220 8-Port PCI-Express to SATA RAID Controller
	1222  ARC-1222 8-Port PCI-Express to SAS/SATA II RAID Controller
	1230  ARC-1230 12-Port PCI-Express to SATA RAID Controller
	1260  ARC-1260 16-Port PCI-Express to SATA RAID Controller
	1280  ARC-1280/1280ML 24-Port PCI-Express to SATA II RAID Controller
	1300  ARC-1300ix-16 16-Port PCI-Express to SAS Non-RAID Host Adapter
	1320  ARC-1320 8/16 Port PCIe 2.0 to SAS/SATA 6Gb Non-RAID Host Adapter
	1330  ARC-1330 16 Port PCIe 3.0 to SAS/SATA 12Gb Non-RAID Host Adapter
	1680  ARC-1680 series PCIe to SAS/SATA 3Gb RAID Controller
	1880  ARC-188x series PCIe 2.0/3.0 to SAS/SATA 6/12Gb RAID Controller
	1884  ARC-1884 series PCIe 3.0 to SAS/SATA 12/6Gb RAID Controller
17d5  Exar Corp.
	5731  Xframe 10-Gigabit Ethernet PCI-X
	5732  Xframe II 10-Gigabit Ethernet PCI-X 2.0
	5831  Xframe 10-Gigabit Ethernet PCI-X
	5832  Xframe II 10-Gigabit Ethernet PCI-X 2.0
	5833  X3100 Series 10 Gigabit Ethernet PCIe
17db  Cray Inc
	0101  XT Series [Seastar] 3D Toroidal Router
17de  KWorld Computer Co. Ltd.
17df  Dini Group
	1864  Virtex4 PCI Board w/ QL5064 Bridge [DN7000K10PCI/DN8000K10PCI/DN8000K10PSX/NOTUS]
	1865  Virtex4 ASIC Emulator [DN8000K10PCIe]
	1866  Virtex4 ASIC Emulator Cable Connection [DN8000K10PCI]
	1867  Virtex4 ASIC Emulator Cable Connection [DN8000K10PCIe]
	1868  Virtex4 ASIC Emulator [DN8000K10PCIe-8]
	1900  Virtex5 PCIe ASIC Emulator [DN9000K10PCIe8T/DN9002K10PCIe8T/DN9200K10PCIe8T/DN7006K10PCIe8T/DN7406K10PCIe8T]
	1901  Virtex5 PCIe ASIC Emulator Large BARs [DN9000K10PCIe8T/DN9002K10PCIe8T/DN9200K10PCIe8T/DN7006K10PCIe8T/DN7406K10PCIe8T]
	1902  Virtex5 PCIe ASIC Emulator Low Power [Interceptor]
	1903  Spartan6 PCIe FPGA Accelerator Board [DNBFCS12PCIe]
	1904  Virtex6 PCIe ASIC Emulation Board [DNDUALV6_PCIe4]
	1905  Virtex6 PCIe ASIC Emulation Board [DNV6F6PCIe]
	1906  Virtex6 PCIe ASIC Emulation Board [DN2076K10]
	1907  Virtex6 PCIe ASIC Emulation Board [DNV6F2PCIe]
	1908  Virtex6 PCIe ASIC Emulation Board Large BARs[DNV6F2PCIe]
	1909  Kintex7 PCIe FPGA Accelerator Board [DNK7F5PCIe]
	190a  Virtex7 PCIe ASIC Emulation Board [DNV7F1A]
	190b  Stratix5 PCIe ASIC Emulation Board [DNS5GXF2]
	190c  Virtex7 PCIe ASIC Emulation Board [DNV7F2A]
	190d  Virtex7 PCIe ASIC Emulation Board [DNV7F4A]
	190e  Virtex7 PCIe ASIC Emulation Board [DNV7F2B]
	190f  KintexUS PCIe MainRef Design [DNPCIE_40G_KU_LL]
	1910  VirtexUS ASIC Emulation Board [DNVUF4A]
	1911  VirtexUS PCIe ASIC Emulation Board [DNVU_F2PCIe]
	1912  KintexUS PCIe MainRef Design [DNPCIe_40G_KU_LL_QSFP]
	1913  VirtexUS ASIC Emulation Board [DNVUF1A]
	1914  VirtexUS ASIC Emulation Board [DNVUF2A]
	1915  Arria10 PCIe MainRef Design [DNPCIe_80G_A10_LL]
	1916  VirtexUS PCIe Accelerator Board [DNVUF2_HPC_PCIe]
	1a00  Virtex6 PCIe DMA Netlist Design
	1a01  Virtex6 PCIe Darklite Design [DNPCIe_HXT_10G_LL]
	1a02  Virtex7 PCIe DMA Netlist Design
	1a03  Kintex7 PCIe Darklite Design [DNPCIe_K7_10G_LL]
	1a05  Stratix5 PCIe Darklite Design [DNS5GX_F2]
	1a06  VirtexUS PCIe DMA Netlist Design
	1a07  KintexUS PCIe Darklite Design [DNPCIe_40G_KU_LL]
	1a08  KintexUS PCIe Darklite Design [DNPCIe_40G_KU_LL_QSFP]
	1a09  Arria10 PCIe Darklite Design [DNPCIe_80G_A10_LL]
	1a0a  VirtexUS PCIe Darklite Design [DNVUF2_HPC_PCIe]
17e4  Sectra AB
	0001  KK671 Cardbus encryption board
	0002  KK672 Cardbus encryption board
17e6  MaxLinear
	0010  EN2010 [c.Link] MoCA Network Controller (Coax, PCI interface)
	0011  EN2010 [c.Link] MoCA Network Controller (Coax, MPEG interface)
	0021  EN2210 [c.Link] MoCA Network Controller (Coax)
	0025  EN2510 [c.Link] MoCA Network Controller (Coax, PCIe interface)
	0027  EN2710 [c.Link] MoCA 2.0 Network Controller (Coax, PCIe interface)
	3700  MoCA 2.0 Network Controller (Coax, PCIe interface)
	3710  MoCA 2.5 Network Controller (Coax, PCIe interface)
17ee  Connect Components Ltd
17f2  Albatron Corp.
17f3  RDC Semiconductor, Inc.
	1010  R1010 IDE Controller
	2012  M2012/R3308 VGA-compatible graphics adapter
	6020  R6020 North Bridge
	6021  R6021 Host Bridge
	6030  R6030 ISA Bridge
	6031  R6031 ISA Bridge
	6040  R6040 MAC Controller
	6060  R6060 USB 1.1 Controller
	6061  R6061 USB 2.0 Controller
17f7  Topdek Semiconductor Inc.
17f9  Gemtek Technology Co., Ltd
17fc  IOGEAR, Inc.
17fe  InProComm Inc.
	2120  IPN 2120 802.11b
	2220  IPN 2220 802.11g
17ff  Benq Corporation
1800  Qualcore Logic Inc.
	1100  Nanospeed Trading Gateway
1803  ProdaSafe GmbH
1805  Euresys S.A.
1809  Lumanate, Inc.
180c  IEI Integration Corp
1813  Ambient Technologies Inc
	4000  HaM controllerless modem
	4100  HaM plus Data Fax Modem
1814  Ralink corp.
	0101  Wireless PCI Adapter RT2400 / RT2460
	0200  RT2500 802.11g PCI [PC54G2]
	0201  RT2500 Wireless 802.11bg
	0300  Wireless Adapter Canyon CN-WF511
	0301  RT2561/RT61 802.11g PCI
	0302  RT2561/RT61 rev B 802.11g
	0401  RT2600 802.11 MIMO
	0601  RT2800 802.11n PCI
	0681  RT2890 Wireless 802.11n PCIe
	0701  RT2760 Wireless 802.11n 1T/2R
	0781  RT2790 Wireless 802.11n 1T/2R PCIe
	3060  RT3060 Wireless 802.11n 1T/1R
	3062  RT3062 Wireless 802.11n 2T/2R
	3090  RT3090 Wireless 802.11n 1T/1R PCIe
	3091  RT3091 Wireless 802.11n 1T/2R PCIe
	3092  RT3092 Wireless 802.11n 2T/2R PCIe
	3290  RT3290 Wireless 802.11n 1T/1R PCIe
	3298  RT3290 Bluetooth
	3592  RT3592 Wireless 802.11abgn 2T/2R PCIe
	359f  RT3592 PCIe Wireless Network Adapter
	5360  RT5360 Wireless 802.11n 1T/1R
	5362  RT5362 PCI 802.11n Wireless Network Adapter
	5390  RT5390 Wireless 802.11n 1T/1R PCIe
	5392  RT5392 PCIe Wireless Network Adapter
	539b  RT5390R 802.11bgn PCIe Wireless Network Adapter
	539f  RT5390 [802.11 b/g/n 1T1R G-band PCI Express Single Chip]
	5592  RT5592 PCIe Wireless Network Adapter
	e932  RT2560F 802.11 b/g PCI
1815  Devolo AG
1820  InfiniCon Systems Inc.
1822  Twinhan Technology Co. Ltd
	4e35  Mantis DTV PCI Bridge Controller [Ver 1.0]
182d  SiteCom Europe BV
	3069  ISDN PCI DC-105V2
	9790  WL-121 Wireless Network Adapter 100g+ [Ver.3]
182e  Raza Microelectronics, Inc.
	0008  XLR516 Processor
182f  Broadcom
	000b  BCM5785 [HT1000] SATA (RAID Mode)
1830  Credence Systems Corporation
183b  MikroM GmbH
	08a7  MVC100 DVI
	08a8  MVC101 SDI
	08a9  MVC102 DVI+Audio
	08b0  MVC200-DC
1846  Alcatel-Lucent
1849  ASRock Incorporation
184a  Thales Computers
	1100  MAX II cPLD
1850  Advantest Corporation
	0048  EK220-66401 Computer Interface Card
1851  Microtune, Inc.
1852  Anritsu Corp.
1853  SMSC Automotive Infotainment System Group
1854  LG Electronics, Inc.
185b  Compro Technology, Inc.
	1489  VideoMate Vista T100
185f  Wistron NeWeb Corp.
1864  SilverBack
	2110  ISNAP 2110
1867  Topspin Communications
	5a44  MT23108 InfiniHost HCA
	5a45  MT23108 InfiniHost HCA flash recovery
	5a46  MT23108 InfiniHost HCA bridge
	6278  MT25208 InfiniHost III Ex (Tavor compatibility mode)
	6282  MT25208 InfiniHost III Ex
186c  Humusoft, s.r.o.
	0612  AD612 Data Acquisition Device
	0614  MF614 Multifunction I/O Card
	0622  AD622 Data Acquisition Device
	0624  MF624 Multifunction I/O PCI Card
	0625  MF625 3-phase Motor Driver
	0634  MF634 Multifunction I/O PCIe Card
186f  WiNRADiO Communications
1876  L-3 Communications
	a101  VigraWATCH PCI
	a102  VigraWATCH PMC
	a103  Vigra I/O
187e  ZyXEL Communications Corporation
	3403  ZyAir G-110 802.11g
	340e  M-302 802.11g XtremeMIMO
1885  Avvida Systems Inc.
1888  Varisys Ltd
	0301  VMFX1 FPGA PMC module
	0601  VSM2 dual PMC carrier
	0710  VS14x series PowerPC PCI board
	0720  VS24x series PowerPC PCI board
188a  Ample Communications, Inc
1890  Egenera, Inc.
1894  KNC One
1896  B&B Electronics Manufacturing Company, Inc.
	4202  MIport 3PCIU2 2-port Serial
	4204  MIport 3PCIU4 4-port Serial
	4208  MIport 3PCIU8 8-port Serial
	4211  MIport 3PCIOU1 1-port Isolated Serial
	4212  MIport 3PCIOU2 2-port Isolated Serial
	4214  MIport 3PCIOU4 4-port Isolated Serial
	bb10  3PCI2 2-Port Serial
	bb11  3PCIO1 1-Port Isolated Serial
1897  AMtek
18a1  Astute Networks Inc.
18a2  Stretch Inc.
	0002  VRC6016 16-Channel PCIe DVR Card
18a3  AT&T
18ac  DViCO Corporation
	d500  FusionHDTV 5
	d800  FusionHDTV 3 Gold
	d810  FusionHDTV 3 Gold-Q
	d820  FusionHDTV 3 Gold-T
	db30  FusionHDTV DVB-T Pro
	db40  FusionHDTV DVB-T Hybrid
	db78  FusionHDTV DVB-T Dual Express
18b8  Ammasso
	b001  AMSO 1100 iWARP/RDMA Gigabit Ethernet Coprocessor
18bc  GeCube Technologies, Inc.
18c3  Micronas Semiconductor Holding AG
	0720  nGene PCI-Express Multimedia Controller
18c8  Cray Inc
18c9  ARVOO Engineering BV
18ca  XGI Technology Inc. (eXtreme Graphics Innovation)
	0020  Z7/Z9 (XG20 core)
	0021  Z9s/Z9m (XG21 core)
	0027  Z11/Z11M
	0040  Volari V3XT/V5/V8
	0047  Volari 8300 (chip: XP10, codename: XG47)
18d2  Sitecom Europe BV (Wrong ID)
	3069  DC-105v2 ISDN controller
18d4  Celestica
18d8  Dialogue Technology Corp.
18dd  Artimi Inc
	4c6f  Artimi RTMI-100 UWB adapter
18df  LeWiz Communications
18e6  MPL AG
	0001  OSCI [Octal Serial Communication Interface]
18eb  Advance Multimedia Internet Technology, Inc.
18ec  Cesnet, z.s.p.o.
	6d05  ML555
	c006  COMBO6
	c032  COMBO-LXT110
	c045  COMBO6E
	c050  COMBO-PTM
	c058  COMBO6X
	c132  COMBO-LXT155
	c232  COMBO-FXT100
18ee  Chenming Mold Ind. Corp.
18f1  Spectrum GmbH
18f4  Napatech A/S
	0031  NT20X Network Adapter
	0051  NT20X Capture Card
	0061  NT20E Capture Card
	0064  NT20E Inline Card
	0071  NT4E Capture Card
	0074  NT4E Inline Card
	0081  NT4E 4-port Expansion Card
	0091  NT20X Capture Card [New Rev]
	00a1  NT4E-STD Capture Card
	00a4  NT4E-STD Inline Card
	00b1  NTBPE Optical Bypass Adapter
	00c5  NT20E2 Network Adapter 2x10Gb
	00d5  NT40E2-4 Network Adapter 4x10Gb
	00e5  NT40E2-1 Network Adapter 1x40Gb
	00f5  NT4E2-4T-BP Network Adapter 4x1Gb with Electrical Bypass
	0105  NT4E2-4-PTP Network Adapter 4x1Gb
	0115  NT20E2-PTP Network Adapter 2x10Gb
	0125  NT4E2-4-PTP Network Adapter 4x1Gb
	0135  NT20E2-PTP Network Adapter 2x10Gb
	0145  NT40E3-4-PTP Network Adapter 4x10Gb
	0155  NT100E3-1-PTP Network Adapter 1x100Gb
	0165  NT80E3-2-PTP Network Adapter 2x40Gb
	0175  NT20E3-2-PTP Network Adapter 2x10Gb
	0185  NT40A01 Network Adapter
	01a5  NT200A01 Network Adapter
18f6  NextIO
	1000  [Nexsis] Switch Virtual P2P PCIe Bridge
	1001  [Texsis] Switch Virtual P2P PCIe Bridge
	1050  [Nexsis] Switch Virtual P2P PCI Bridge
	1051  [Texsis] Switch Virtual P2P PCI Bridge
	2000  [Nexsis] Switch Integrated Mgmt. Endpoint
	2001  [Texsis] Switch Integrated Mgmt. Endpoint
18f7  Commtech, Inc.
	0001  ESCC-PCI-335 Serial PCI Adapter [Fastcom]
	0002  422/4-PCI-335 Serial PCI Adapter [Fastcom]
	0003  232/4-1M-PCI Serial PCI Adapter [Fastcom]
	0004  422/2-PCI-335 Serial PCI Adapter [Fastcom]
	0005  IGESCC-PCI-ISO/1 Serial PCI Adapter [Fastcom]
	000a  232/4-PCI-335 Serial PCI Adapter [Fastcom]
	000b  232/8-PCI-335 Serial PCI Adapter [Fastcom]
	000f  FSCC Serial PCI Adapter [Fastcom]
	0010  GSCC Serial PCI Adapter [Fastcom]
	0011  QSSB Serial PCI Adapter [Fastcom]
	0014  SuperFSCC Serial PCI Adapter [Fastcom]
	0015  SuperFSCC-104-LVDS Serial PC/104+ Adapter [Fastcom]
	0016  FSCC-232 RS-232 Serial PCI Adapter [Fastcom]
	0017  SuperFSCC-104 Serial PC/104+ Adapter [Fastcom]
	0018  SuperFSCC/4 Serial PCI Adapter [Fastcom]
	0019  SuperFSCC Serial PCI Adapter [Fastcom]
	001a  SuperFSCC-LVDS Serial PCI Adapter [Fastcom]
	001b  FSCC/4 Serial PCI Adapter [Fastcom]
	001c  SuperFSCC/4-LVDS Serial PCI Adapter [Fastcom]
	001d  FSCC Serial PCI Adapter [Fastcom]
	001e  SuperFSCC/4 Serial PCIe Adapter [Fastcom]
	001f  SuperFSCC/4 Serial cPCI Adapter [Fastcom]
	0020  422/4-PCIe Serial PCIe Adapter [Fastcom]
	0021  422/8-PCIe Serial PCIe Adapter [Fastcom]
	0022  SuperFSCC/4-LVDS Serial PCIe Adapter [Fastcom]
	0023  SuperFSCC/4 Serial cPCI Adapter [Fastcom]
	0025  SuperFSCC/4-LVDS Serial PCI Adapter [Fastcom]
	0026  SuperFSCC-LVDS Serial PCI Adapter [Fastcom]
	0027  FSCC/4 Serial PCIe Adapter [Fastcom]
18fb  Resilience Corporation
1904  Hangzhou Silan Microelectronics Co., Ltd.
	2031  SC92031 PCI Fast Ethernet Adapter
	8139  RTL8139D [Realtek] PCI 10/100BaseTX ethernet adaptor
1905  Micronas USA, Inc.
1912  Renesas Technology Corp.
	0002  SH7780 PCI Controller (PCIC)
	0011  SH7757 PCIe End-Point [PBI]
	0012  SH7757 PCIe-PCI Bridge [PPB]
	0013  SH7757 PCIe Switch [PS]
	0014  uPD720201 USB 3.0 Host Controller
	0015  uPD720202 USB 3.0 Host Controller
	001a  SH7758 PCIe-PCI Bridge [PPB]
	001b  SH7758 PCIe End-Point [PBI]
	001d  SH7758 PCIe Switch [PS]
1919  Soltek Computer Inc.
1923  Sangoma Technologies Corp.
	0040  A200/Remora FXO/FXS Analog AFT card
	0100  A104d QUAD T1/E1 AFT card
	0300  A101 single-port T1/E1
	0400  A104u Quad T1/E1 AFT
1924  Solarflare Communications
	0703  SFC4000 rev A net [Solarstorm]
	0710  SFC4000 rev B [Solarstorm]
	0803  SFC9020 10G Ethernet Controller
	0813  SFL9021 10GBASE-T Ethernet Controller
	0903  SFC9120 10G Ethernet Controller
	0923  SFC9140 10/40G Ethernet Controller
	0a03  SFC9220 10/40G Ethernet Controller
	1803  SFC9020 10G Ethernet Controller (Virtual Function)
	1813  SFL9021 10GBASE-T Ethernet Controller (Virtual Function)
	1903  SFC9120 10G Ethernet Controller (Virtual Function)
	1923  SFC9140 10/40G Ethernet Controller (Virtual Function)
	1a03  SFC9220 10/40G Ethernet Controller (Virtual Function)
	6703  SFC4000 rev A iSCSI/Onload [Solarstorm]
	c101  EF1-21022T [EtherFabric]
192a  BiTMICRO Networks Inc.
192e  TransDimension
1931  Option N.V.
	000c  Qualcomm MSM6275 UMTS chip
1932  DiBcom
193c  MAXIM Integrated Products
193f  AHA Products Group
	0001  AHA36x-PCIX
	0360  AHA360-PCIe
	0363  AHA363-PCIe
	0364  AHA364-PCIe
	0367  AHA367-PCIe
	0370  AHA370-PCIe
	0604  AHA604
	0605  AHA605
	3641  AHA3641
	3642  AHA3642
	6101  AHA6101
	6102  AHA6102
1942  ClearSpeed Technology plc
	e511  Advance X620 accelerator card
	e521  Advance e620 accelerator card
1947  C-guys, Inc.
	4743  CG200 Dual SD/SDIO Host controller device
1948  Alpha Networks Inc.
194a  DapTechnology B.V.
	1111  FireSpy3850
	1112  FireSpy450b
	1113  FireSpy450bT
	1114  FireSpy850
	1115  FireSpy850bT
	1200  FireTrac 3460bT
	1201  FireTrac 3460bT (fallback firmware)
	1202  FireTrac 3460bT
	1203  FireTrac 3460bT (fallback firmware)
1954  One Stop Systems, Inc.
1957  Freescale Semiconductor Inc
	0012  MPC8548E
	0013  MPC8548
	0014  MPC8543E
	0015  MPC8543
	0018  MPC8547E
	0019  MPC8545E
	001a  MPC8545
	0020  MPC8568E
	0021  MPC8568
	0022  MPC8567E
	0023  MPC8567
	0030  MPC8533E
	0031  MPC8533
	0032  MPC8544E
	0033  MPC8544
	0040  MPC8572E
	0041  MPC8572
	0050  MPC8536E
	0051  MPC8536
	0052  MPC8535E
	0053  MPC8535
	0060  MPC8569
	0061  MPC8569E
	0070  P2020E
	0071  P2020
	0078  P2010E
	0079  P2010
	0080  MPC8349E
	0081  MPC8349
	0082  MPC8347E TBGA
	0083  MPC8347 TBGA
	0084  MPC8347E PBGA
	0085  MPC8347 PBGA
	0086  MPC8343E
	0087  MPC8343
	00b4  MPC8315E
	00b6  MPC8314E
	00c2  MPC8379E
	00c3  MPC8379
	00c4  MPC8378E
	00c5  MPC8378
	00c6  MPC8377E
	00c7  MPC8377
	0100  P1020E
	0101  P1020
	0102  P1021E
	0103  P1021
	0108  P1011E
	0109  P1011
	010a  P1012E
	010b  P1012
	0110  P1022E
	0111  P1022
	0118  P1013E
	0119  P1013
	0128  P1010
	0400  P4080E
	0401  P4080
	0408  P4040E
	0409  P4040
	041f  P3041
	0440  T4240 with security
	0441  T4240 without security
	0446  T4160 with security
	0447  T4160 without security
	0830  T2080 with security
	0831  T2080 without security
	0838  T2081 with security
	0839  T2081 without security
	580c  MPC5121e
	7010  MPC8641 PCI Host Bridge
	7011  MPC8641D PCI Host Bridge
	7018  MPC8610
	c006  MPC8308
	fc02  RedStone
	fc03  CFI
1958  Faster Technology, LLC.
1959  PA Semi, Inc
	a000  PA6T Core
	a001  PWRficient Host Bridge
	a002  PWRficient PCI-Express Port
	a003  PWRficient SMBus Controller
	a004  PWRficient 16550 UART
	a005  PWRficient Gigabit Ethernet
	a006  PWRficient 10-Gigabit Ethernet
	a007  PWRficient DMA Controller
	a008  PWRficient LPC/Localbus Interface
	a009  PWRficient L2 Cache
	a00a  PWRficient DDR2 Memory Controller
	a00b  PWRficient SERDES
	a00c  PWRficient System/Debug Controller
	a00d  PWRficient PCI-Express Internal Endpoint
1966  Orad Hi-Tec Systems
	1975  DVG64 family
	1977  DVG128 family
1969  Qualcomm Atheros
	1026  AR8121/AR8113/AR8114 Gigabit or Fast Ethernet
	1048  Attansic L1 Gigabit Ethernet
	1062  AR8132 Fast Ethernet
	1063  AR8131 Gigabit Ethernet
	1066  Attansic L2c Gigabit Ethernet
	1067  Attansic L1c Gigabit Ethernet
	1073  AR8151 v1.0 Gigabit Ethernet
	1083  AR8151 v2.0 Gigabit Ethernet
	1090  AR8162 Fast Ethernet
	1091  AR8161 Gigabit Ethernet
	10a0  QCA8172 Fast Ethernet
	10a1  QCA8171 Gigabit Ethernet
	2048  Attansic L2 Fast Ethernet
	2060  AR8152 v1.1 Fast Ethernet
	2062  AR8152 v2.0 Fast Ethernet
	e091  Killer E220x Gigabit Ethernet Controller
	e0a1  Killer E2400 Gigabit Ethernet Controller
	e0b1  Killer E2500 Gigabit Ethernet Controller
196a  Sensory Networks Inc.
	0101  NodalCore C-1000 Content Classification Accelerator
	0102  NodalCore C-2000 Content Classification Accelerator
	0105  NodalCore C-3000 Content Classification Accelerator
196d  Club-3D BV
1971  AGEIA Technologies, Inc.
	1011  Physics Processing Unit [PhysX]
1974  Eberspaecher Electronics
1976  TRENDnet
1977  Parsec
197b  JMicron Technology Corp.
	0250  JMC250 PCI Express Gigabit Ethernet Controller
	0260  JMC260 PCI Express Fast Ethernet Controller
	0368  JMB368 IDE controller
	2360  JMB360 AHCI Controller
	2361  JMB361 AHCI/IDE
	2362  JMB362 SATA Controller
	2363  JMB363 SATA/IDE Controller
	2364  JMB364 AHCI Controller
	2365  JMB365 AHCI/IDE
	2366  JMB366 AHCI/IDE
	2368  JMB368 IDE controller
	2369  JMB369 Serial ATA Controller
	2380  IEEE 1394 Host Controller
	2381  Standard SD Host Controller
	2382  SD/MMC Host Controller
	2383  MS Host Controller
	2384  xD Host Controller
	2386  Standard SD Host Controller
	2387  SD/MMC Host Controller
	2388  MS Host Controller
	2389  xD Host Controller
	2391  Standard SD Host Controller
	2392  SD/MMC Host Controller
	2393  MS Host Controller
	2394  xD Host Controller
1982  Distant Early Warning Communications Inc
	1600  OX16C954 HOST-A
	16ff  OX16C954 HOST-B
1989  Montilio Inc.
	0001  RapidFile Bridge
	8001  RapidFile
198a  Nallatech Ltd.
1993  Innominate Security Technologies AG
1999  A-Logics
	a900  AM-7209 Video Processor
199a  Pulse-LINK, Inc.
199d  Xsigo Systems
	8209  Virtual NIC Device
	890a  Virtual HBA Device
199f  Auvitek
	8501  AU85X1 PCI REV1.1
	8521  AU8521 TV card
19a2  Emulex Corporation
	0120  x1 PCIe Gen2 Bridge[Pilot4]
	0200  BladeEngine 10Gb PCI-E iSCSI adapter
	0201  BladeEngine 10Gb PCIe Network Adapter
	0211  BladeEngine2 10Gb Gen2 PCIe Network Adapter
	0212  BladeEngine2 10Gb Gen2 PCIe iSCSI Adapter
	0221  BladeEngine3 10Gb Gen2 PCIe Network Adapter
	0222  BladeEngine3 10Gb Gen2 PCIe iSCSI Adapter
	0700  OneConnect OCe10100/OCe10102 Series 10 GbE
	0702  OneConnect 10Gb iSCSI Initiator
	0704  OneConnect OCe10100/OCe10102 Series 10 GbE CNA
	0710  OneConnect 10Gb NIC (be3)
	0712  OneConnect 10Gb iSCSI Initiator (be3)
	0714  OneConnect 10Gb FCoE Initiator (be3)
	0800  ServerView iRMC HTI
19a8  DAQDATA GmbH
19ac  Kasten Chase Applied Research
	0001  ACA2400 Crypto Accelerator
19ae  Progeny Systems Corporation
	0520  4135 HFT Interface Controller
	0521  Decimator
19ba  ZyXEL Communications Corp.
	2330  ZyWALL Turbo Card
19c1  Exegy Inc.
19d1  Motorola Expedience
19d4  Quixant Limited
19da  ZOTAC International (MCO) Ltd.
19de  Pico Computing
19e2  Vector Informatik GmbH
19e3  DDRdrive LLC
	5801  DDRdrive X1
	5808  DDRdrive X8
	dd52  DDRdrive X1-30
19e5  Huawei Technologies Co., Ltd.
	1711  Hi1710 [iBMC Intelligent Management system chip w/VGA support]
19e7  NET (Network Equipment Technologies)
	1001  STIX DSP Card
	1002  STIX - 1 Port T1/E1 Card
	1003  STIX - 2 Port T1/E1 Card
	1004  STIX - 4 Port T1/E1 Card
	1005  STIX - 4 Port FXS Card
19ee  Netronome Systems, Inc.
19f1  BFG Tech
19ff  Eclipse Electronic Systems, Inc.
1a03  ASPEED Technology, Inc.
	1150  AST1150 PCI-to-PCI Bridge
	2000  ASPEED Graphics Family
1a07  Kvaser AB
	0006  CAN interface PC104+ HS/HS
	0007  CAN interface PCIcanx II HS or HS/HS
	0008  CAN interface PCIEcan HS or HS/HS
	0009  CAN interface PCI104 HS/HS
1a08  Sierra semiconductor
	0000  SC15064
1a0e  DekTec Digital Video B.V.
	083f  DTA-2111 VHF/UHF Modulator
1a17  Force10 Networks, Inc.
	8002  PB-10GE-2P 10GbE Security Card
1a1d  GFaI e.V.
	1a17  Meta Networks MTP-1G IDPS NIC
1a1e  3Leaf Systems, Inc.
1a22  Ambric Inc.
1a29  Fortinet, Inc.
	4338  CP8 Content Processor ASIC
	4e36  NP6 Network Processor
1a2b  Ascom AG
	0000  GESP v1.2
	0001  GESP v1.3
	0002  ECOMP v1.3
	0005  ETP v1.4
	000a  ETP-104 v1.1
	000e  DSLP-104 v1.1
1a30  Lantiq
	0680  MtW8171 [Hyperion II]
	0700  Wave300 PSB8224 [Hyperion III]
	0710  Wave300 PSB8231 [Hyperion III]
1a32  Quanta Microsystems, Inc
1a3b  AzureWave
	1112  AR9285 Wireless Network Adapter (PCI-Express)
1a41  Tilera Corp.
	0001  TILE64 processor
	0002  TILEPro processor
	0200  TILE-Gx processor
	0201  TILE-Gx Processor Virtual Function
	2000  TILE-Gx PCI Express Root Port
1a4a  SLAC National Accelerator Lab PPA-REG
	1000  MCOR Power Supply Controller
	1010  AMC EVR - Stockholm Timing Board
	2000  PGPCard - 4 Lane
	2001  PGPCard - 8 Lane Plus EVR
	2010  PCI-Express EVR
1a51  Hectronic AB
1a55  Rohde & Schwarz DVS GmbH
	0010  SDStationOEM
	0011  SDStationOEM II
	0020  Centaurus
	0021  Centaurus II
	0022  Centaurus II LT
	0030  CLIPSTER-VPU 1.x (Hugo)
	0040  Hydra Cinema (JPEG)
	0050  CLIPSTER-VPU 2.x (DigiLab)
	0060  CLIPSTER-DCI 2.x (HydraX)
	0061  Atomix
	0062  Atomix LT
	0063  Atomix HDMI
	0064  Atomix STAN
	0065  Atomix HDMI STAN
	0070  RED Rocket
	0090  CinePlay
1a56  Bigfoot Networks, Inc.
1a57  Highly Reliable Systems
1a58  Razer USA Ltd.
1a5d  Celoxica
1a5e  Aprius Inc.
1a5f  System TALKS Inc.
1a68  VirtenSys Limited
1a71  XenSource, Inc.
1a73  Violin Memory, Inc
	0001  Mozart [Memory Appliance 1010]
1a76  Wavesat
1a77  Lightfleet Corporation
1a78  Virident Systems Inc.
	0031  FlashMAX Drive
	0040  FlashMAX II
	0041  FlashMAX II
	0042  FlashMAX II
	0050  FlashMAX III
1a84  Commex Technologies
	0001  Vulcan SP HT6210 10-Gigabit Ethernet (rev 02)
1a88  MEN Mikro Elektronik
	4d45  Multifunction IP core
1a8a  StarBridge, Inc.
1a8c  Verigy Pte. Ltd.
	1100  E8001-66443 PCI Express CIC
1a8e  DRS Technologies
	2090  Model 2090 PCI Express
1aa8  Ciprico, Inc.
	0009  RAIDCore Controller
	000a  RAIDCore Controller
1aae  Global Velocity, Inc.
1ab6  CalDigit, Inc.
	6201  RAID Card
1ab8  Parallels, Inc.
	4000  Virtual Machine Communication Interface
	4005  Accelerated Virtual Video Adapter
	4006  Memory Ballooning Controller
1ab9  Espia Srl
1ac8  Aeroflex Gaisler
1acc  Point of View BV
1ad7  Spectracom Corporation
	8000  TSync-PCIe Time Code Processor
	9100  TPRO-PCI-66U Timecode Reader/Generator
1ade  Spin Master Ltd.
	1501  Swipetech barcode scanner
	3038  PCIe Video Bridge
1ae0  Google, Inc.
1ae7  First Wise Media GmbH
	0520  HFC-S PCI A [X-TENSIONS XC-520]
1ae8  Silicon Software GmbH
	0a40  microEnable IV-BASE x1
	0a41  microEnable IV-FULL x1
	0a44  microEnable IV-FULL x4
	0e44  microEnable IV-GigE x4
1ae9  Wilocity Ltd.
	0101  Wil6200 PCI Express Root Port
	0200  Wil6200 PCI Express Port
	0201  Wil6200 Wireless PCI Express Port
	0301  Wil6200 802.11ad Wireless Network Adapter
	0302  Wil6200 802.11ad Wireless Network Adapter
	0310  Wil6200 802.11ad Wireless Network Adapter
1aea  Alcor Micro
	6601  AU6601 PCI-E Flash card reader controller
1aec  Wolfson Microelectronics
1aed  SanDisk
	1003  ioDimm3 (v1.2)
	1005  ioDimm3
	1006  ioXtreme
	1007  ioXtreme Pro
	1008  ioXtreme-2
	2001  ioDrive2
	3001  ioMemory FHHL
	3002  ioMemory HHHL
	3003  ioMemory Mezzanine
1aee  Caustic Graphics Inc.
1af4  Red Hat, Inc.
	1000  Virtio network device
	1001  Virtio block device
	1002  Virtio memory balloon
	1003  Virtio console
	1004  Virtio SCSI
	1005  Virtio RNG
	1009  Virtio filesystem
	1041  Virtio network device
	1042  Virtio block device
	1043  Virtio console
	1044  Virtio RNG
	1045  Virtio memory balloon
	1048  Virtio SCSI
	1049  Virtio filesystem
	1050  Virtio GPU
	1052  Virtio input
	1110  Inter-VM shared memory
1af5  Netezza Corp.
1afa  J & W Electronics Co., Ltd.
1b03  Magnum Semiconductor, Inc,
	6100  DXT/DXTPro Multiformat Broadcast HD/SD Encoder/Decoder/Transcoder
	7000  D7 Multiformat Broadcast HD/SD Encoder/Decoder/Transcoder
1b08  MSC Technologies GmbH
1b0a  Pegatron
1b13  Jaton Corp
1b1a  K&F Computing Research Co.
	0e70  GRAPE
1b21  ASMedia Technology Inc.
	0611  ASM1061 SATA IDE Controller
	0612  ASM1062 Serial ATA Controller
	1042  ASM1042 SuperSpeed USB Host Controller
	1080  ASM1083/1085 PCIe to PCI Bridge
	1142  ASM1042A USB 3.0 Host Controller
	1242  ASM1142 USB 3.1 Host Controller
1b2c  Opal-RT Technologies Inc.
1b36  Red Hat, Inc.
	0001  QEMU PCI-PCI bridge
	0002  QEMU PCI 16550A Adapter
	0003  QEMU PCI Dual-port 16550A Adapter
	0004  QEMU PCI Quad-port 16550A Adapter
	0005  QEMU PCI Test Device
	0006  PCI Rocker Ethernet switch device
	0007  PCI SD Card Host Controller Interface
	0008  QEMU PCIe Host bridge
	0009  QEMU PCI Expander bridge
	000a  PCI-PCI bridge (multiseat)
	000b  QEMU PCIe Expander bridge
	000c  QEMU PCIe Root port
	000d  QEMU XHCI Host Controller
	0100  QXL paravirtual graphic card
1b37  Signal Processing Devices Sweden AB
	0001  ADQ214
	0003  ADQ114
	0005  ADQ112
	000e  ADQ108
	000f  ADQDSP
	0014  ADQ412
	0015  ADQ212
	001b  SDR14
	001c  ADQ1600
	001e  ADQ208
	001f  DSU
	0020  ADQ14
	0023  ADQ7
	2014  TX320
	2019  S6000
1b39  sTec, Inc.
	0001  S1120 PCIe Accelerator SSD
1b3a  Westar Display Technologies
	7589  HRED J2000 - JPEG 2000 Video Codec Device
1b3e  Teradata Corp.
	1fa8  BYNET BIC2SE/X
1b40  Schooner Information Technology, Inc.
1b47  Numascale AS
	0601  NumaChip N601
	0602  NumaChip N602
1b4b  Marvell Technology Group Ltd.
	0640  88SE9128 SATA III 6Gb/s RAID Controller
	9120  88SE9120 SATA 6Gb/s Controller
	9123  88SE9123 PCIe SATA 6.0 Gb/s controller
	9125  88SE9125 PCIe SATA 6.0 Gb/s controller
	9128  88SE9128 PCIe SATA 6 Gb/s RAID controller
	9130  88SE9128 PCIe SATA 6 Gb/s RAID controller with HyperDuo
	9172  88SE9172 SATA 6Gb/s Controller
	9178  88SE9170 PCIe SATA 6Gb/s Controller
	917a  88SE9172 SATA III 6Gb/s RAID Controller
	9183  88SS9183 PCIe SSD Controller
	9192  88SE9172 SATA III 6Gb/s RAID Controller
	91a0  88SE912x SATA 6Gb/s Controller [IDE mode]
	91a4  88SE912x IDE Controller
	9220  88SE9220 PCIe 2.0 x2 2-port SATA 6 Gb/s RAID Controller
	9230  88SE9230 PCIe SATA 6Gb/s Controller
	9235  88SE9235 PCIe 2.0 x2 4-port SATA 6 Gb/s Controller
	9445  88SE9445 PCIe 2.0 x4 4-Port SAS/SATA 6 Gbps RAID Controller
	9480  88SE9480 SAS/SATA 6Gb/s RAID controller
	9485  88SE9485 SAS/SATA 6Gb/s controller
1b55  NetUP Inc.
	18f6  Dual DVB Universal CI card
	18f7  Dual DVB Universal CI card rev 1.4
	2a2c  Dual DVB-S2-CI card
	e2e4  Dual DVB-T/C-CI RF card
	e5f4  MPEG2 and H264 Encoder-Transcoder
	f1c4  Dual ASI-RX/TX-CI card
1b66  Deltacast
	0007  Delta-3G-elp-11 SDI I/O Board
1b6f  Etron Technology, Inc.
	7023  EJ168 USB 3.0 Host Controller
	7052  EJ188/EJ198 USB 3.0 Host Controller
1b73  Fresco Logic
	1000  FL1000G USB 3.0 Host Controller
	1009  FL1009 USB 3.0 Host Controller
	1100  FL1100 USB 3.0 Host Controller
1b74  OpenVox Communication Co. Ltd.
	0115  D115P/D115E Single-port E1/T1 card
	d130  D130P/D130E Single-port E1/T1 card (3rd GEN)
	d210  D210P/D210E Dual-port E1/T1 card(2nd generation)
	d230  D230 Dual-port E1/T1 card (2nd generation)
	d410  D410/430 Quad-port E1/T1 card
	d430  D410/430 Quad-port E1/T1 card
1b79  Absolute Analysis
1b85  OCZ Technology Group, Inc.
	1041  RevoDrive 3 X2 PCI-Express SSD 240 GB (Marvell Controller)
	6018  RD400/400A SSD
	8788  RevoDrive Hybrid
1b94  Signatec / Dynamic Signals Corp
	e400  PX14400 Dual Xilinx Virtex5 based Digitizer
1b96  Western Digital
1b9a  XAVi Technologies Corp.
1bad  ReFLEX CES
1bb0  SimpliVity Corporation
	0002  OmniCube Accelerator OA-3000
	0010  OmniCube Accelerator OA-3000-2
1bb1  Seagate Technology PLC
	005d  Nytro PCIe Flash Storage
	0100  Nytro Flash Storage
1bb3  Bluecherry
	4304  BC-04120A MPEG4 4 port video encoder / decoder
	4309  BC-08240A MPEG4 4 port video encoder / decoder
	4310  BC-16480A MPEG4 16 port video encoder / decoder
	4e04  BC-04120A 4 port MPEG4 video encoder / decoder
	4e09  BC-08240A 8 port MPEG4 video encoder / decoder
	4e10  BC-16480A 16 port MPEG4 video encoder / decoder
	5304  BC-H04120A 4 port H.264 video and audio encoder / decoder
	5308  BC-H08240A 8 port H.264 video and audio encoder / decoder
	5310  BC-H16480A 16 port H.264 video and audio encoder / decoder
1bb5  Quantenna Communications, Inc.
1bbf  Maxeler Technologies Ltd.
	0003  MAX3
	0004  MAX4
1bd0  Astronics Corporation
	1001  Mx5 PMC/XMC Databus Interface Card
	1002  PM1553-5 (PC/104+ MIL-STD-1553 Interface Card)
	1004  AB3000 Series Rugged Computer
	1005  PE1000 (Multi-Protocol PCIe/104 Interface Card)
	1101  OmniBus II PCIe Multi-Protocol Interface Card
	1102  OmniBusBox II Multi-Protocol Interface Core
	1103  OmniBus II cPCIe/PXIe Multi-Protocol Interface Card
1bd4  Inspur Electronic Information Industry Co., Ltd.
1bee  IXXAT Automation GmbH
	0003  CAN-IB200/PCIe
1bef  Lantiq
	0011  MIPS SoC PCI Express Port
1bf4  VTI Instruments Corporation
	0001  SentinelEX
1bfd  EeeTOP
1c09  CSP, Inc.
	4254  10G-PCIE3-8D-2S
	4255  10G-PCIE3-8D-Q
	4256  10G-PCIE3-8D-2S
	4258  10G-PCIE3-8E-2S Network Adapter
	4260  10G-PCIE3-8E-4S Network Adapter
	4261  10G-PCIE3-8E-4S Network Adapter
	4262  10G-PCIE3-8E-4S Network Adapter
	4263  10G-PCIE3-8E-4S Network Adapter
	4264  10G-PCIE3-8E-2S Network Adapter
	4265  10G-PCIE3-8E-2S Network Adapter
1c1c  Symphony
	0001  82C101
1c28  Lite-On IT Corp. / Plextor
	0122  M6e PCI Express SSD [Marvell 88SS9183]
1c2c  Fiberblaze
	000a  Capture
	000f  SmartNIC
	00a0  FBC4G Capture 4x1Gb
	00a1  FBC4XG Capture 4x10Gb
	00a2  FBC8XG Capture 8x10Gb
	00a3  FBC2XG Capture 2x10Gb
	00a4  FBC4XGG3 Capture 4x10Gb
	00a5  FBC2XLG Capture 2x40Gb
	00a6  FBC1CG Capture 1x100Gb
	00a9  FBC2XGHH Capture 2x10Gb
	00ad  FBC2CGG3HL Capture 2x200Gb
	00af  Capture slave device
	a001  FBC2CGG3 Capture 2x200Gb
1c32  Highland Technology, Inc.
1c33  Daktronics, Inc
1c3b  Accensus, LLC
	0200  Telas2
	0300  Telas 2.V
1c44  Enmotus Inc
	8000  8000 Storage IO Controller
1c58  HGST, Inc.
	0003  Ultrastar SN100 Series NVMe SSD
	0023  Ultrastar SN200 Series NVMe SSD
1c5f  Beijing Memblaze Technology Co. Ltd.
	0540  PBlaze4 NVMe SSD
1c63  Science and Research Centre of Computer Technology (JSC "NICEVT")
	0008  K1927BB1Ya [EC8430] Angara Interconnection Network Adapter
1c7e  TTTech Computertechnik AG
	0200  zFAS Debug Port
1c7f  Elektrobit Austria GmbH
	5100  EB5100
1c8a  TSF5 Corporation
	0001  Hunter PCI Express
1cb1  Collion UG & Co.KG
1cb8  Dawning Information Industry Co., Ltd.
1cc5  Embedded Intelligence, Inc.
	0100  CAN-PCIe-02
1cc7  Radian Memory Systems Inc.
	0200  RMS-200
	0250  RMS-250
1ccf  Zoom Corporation
	0001  TAC-2 Thunderbolt Audio Converter
1cd2  SesKion GmbH
	0301  Simulyzer-RT CompactPCI Serial DIO-1 card
	0302  Simulyzer-RT CompactPCI Serial PSI5-ECU-1 card
	0303  Simulyzer-RT CompactPCI Serial PSI5-SIM-1 card
	0304  Simulyzer-RT CompactPCI Serial PWR-ANA-1 card
	0305  Simulyzer-RT CompactPCI Serial CAN-1 card
1cd7  Nanjing Magewell Electronics Co., Ltd.
	0010  Pro Capture Endpoint
1cdd  secunet Security Networks AG
1ce4  Exablaze
	0001  ExaNIC X4
	0002  ExaNIC X2
	0003  ExaNIC X10
	0004  ExaNIC X10-GM
	0005  ExaNIC X40
	0006  ExaNIC X10-HPT
	0007  ExaNIC X40
	0008  ExaNIC V5P
1cf7  Subspace Dynamics
1d00  Pure Storage
1d0f  Amazon.com, Inc.
	cd01  NVMe SSD Controller
	ec20  Elastic Network Adapter (ENA)
1d17  Zhaoxin
	070f  ZX-100 PCI Express Root Port
	0710  ZX-100/ZX-200 PCI Express Root Port
	0711  ZX-100/ZX-200 PCI Express Root Port
	0712  ZX-100/ZX-200 PCI Express Root Port
	0713  ZX-100/ZX-200 PCI Express Root Port
	0714  ZX-100/ZX-200 PCI Express Root Port
	0715  ZX-100/ZX-200 PCI Express Root Port
	0716  ZX-D PCI Express Root Port
	0717  ZX-D PCI Express Root Port
	0718  ZX-D PCI Express Root Port
	0719  ZX-D PCI Express Root Port
	071a  ZX-D PCI Express Root Port
	071b  ZX-D PCI Express Root Port
	071c  ZX-D PCI Express Root Port
	071d  ZX-D PCI Express Root Port
	071e  ZX-D PCI Express Root Port
	071f  ZX-200 Upstream Port of PCI Express Switch
	0720  ZX-200 PCIE RC6 controller
	0721  ZX-200 Downstream Port of PCI Express Switch
	0722  ZX-200 PCIE P2C bridge
	1000  ZX-D Standard Host Bridge
	1001  ZX-D Miscellaneous Bus
	3001  ZX-100 Standard Host Bridge
	300a  ZX-100 Miscellaneous Bus
	3038  ZX-100/ZX-200 Standard Universal PCI to USB Host Controller
	3104  ZX-100/ZX-200 Standard Enhanced PCI to USB Host Controller
	31b0  ZX-100/ZX-D Standard Host Bridge
	31b1  ZX-100/ZX-D Standard Host Bridge
	31b2  ZX-100/ZX-D DRAM Controller
	31b3  ZX-100/ZX-D Power Management Controller
	31b4  ZX-100/ZX-D I/O APIC
	31b5  ZX-100/ZX-D Scratch Device
	31b7  ZX-100/ZX-D Standard Host Bridge
	31b8  ZX-100/ZX-D PCI to PCI Bridge
	3288  ZX-100/ZX-D High Definition Audio Controller
	345b  ZX-100/ZX-D Miscellaneous Bus
	3a02  ZX-100 C-320 GPU
	3a03  ZX-D C-860 GPU
	9002  ZX-100/ZX-200 EIDE Controller
	9003  ZX-100 EIDE Controller
	9045  ZX-100/ZX-D RAID Accelerator
	9046  ZX-D RAID Accelerator
	9083  ZX-100/ZX-200 StorX AHCI Controller
	9084  ZX-100 StorX AHCI Controller
	9100  ZX-200 Cross bus
	9101  ZX-200 Traffic Controller
	9141  ZX-100 High Definition Audio Controller
	9142  ZX-D High Definition Audio Controller
	9180  ZX-200 Networking Gigabit Ethernet Adapter
	9202  ZX-100 USB eXtensible Host Controller
	9203  ZX-200 USB eXtensible Host Controller
	9286  ZX-D eMMC Host Controller
	9300  ZX-D eSPI Host Controller
	95d0  ZX-100 Universal SD Host Controller
	f410  ZX-100/ZX-D PCI Com Port
1d18  RME
	0001  Fireface UFX+
1d1d  CNEX Labs
	1f1f  QEMU NVM Express LightNVM Controller
	2807  8800 series NVMe SSD
1d21  Allo
1d26  Kalray Inc.
	0040  Turbocard2 Accelerator
	0080  Open Network Interface Card 80G
	00c0  Turbocard3 Accelerator
	e004  AB01/EMB01 Development Board
1d40  Techman Electronics (Changshu) Co., Ltd.
1d44  DPT
	a400  PM2x24/PM3224
1d49  Lenovo
1d4c  Diamanti, Inc.
1d5c  Fantasia Trading LLC
1d61  Technobox, Inc.
1d62  Nebbiolo Technologies
1d65  Imagine Communications Corp.
	04de  Taurus/McKinley
1d6a  Aquantia Corp.
	d107  AQC107 NBase-T/IEEE 802.3bz Ethernet Controller [AQtion]
1d6c  Atomic Rules LLC
	1001  A5PL-E1
	1002  A5PL-E7
	1003  S5PEDS-AB
	1004  KC705-K325
	1005  ZC706-Z045
	1006  KCU105-KU040
	1007  XUSP3S-VU095 [Jasper]
	1008  XUSPL4-VU065 [Mustang UltraScale]
	1009  XUSPL4-VU3P [Mustang UltraScale+]
	100a  A10PL4-A10GX115
	100b  K35-2SFP
	100c  K35-4SFP
	100d  AR-ARKA-FX0 [Arkville 32B DPDK Data Mover]
	100e  AR-ARKA-FX1 [Arkville 64B DPDK Data Mover]
	4200  A5PL-E1-10GETI [10 GbE Ethernet Traffic Instrument]
1d78  DERA
1d7c  Aerotech, Inc.
1d87  Fuzhou Rockchip Electronics Co., Ltd
1d8f  Enyx
1d95  Graphcore Ltd
1da1  Teko Telecom S.r.l.
1da2  Sapphire Technology Limited
1de1  Tekram Technology Co.,Ltd.
	0391  TRM-S1040 [DC-315 / DC-395 series]
	2020  DC-390
	690c  690c
	dc29  DC290
1de5  Eideticom, Inc
	1000  IO Memory Controller
	2000  NoLoad Hardware Development Kit
1fc0  Ascom (Finland) Oy
	0300  E2200 Dual E1/Rawpipe Card
	0301  C5400 SHDSL/E1 Card
1fc1  QLogic, Corp.
	000d  IBA6110 InfiniBand HCA
	0010  IBA6120 InfiniBand HCA
1fc9  Tehuti Networks Ltd.
	3009  10-Giga TOE SmartNIC
	3010  10-Giga TOE SmartNIC
	3014  10-Giga TOE SmartNIC 2-Port
	3110  10-Giga TOE Single Port SmartNIC
	3114  10-Giga TOE Dual Port Low Profile SmartNIC
	3310  10-Giga TOE SFP+ Single Port SmartNIC
	3314  10-Giga TOE Dual Port Low Profile SmartNIC
	4010  TN4010 Clean SROM
	4020  TN9030 10GbE CX4 Ethernet Adapter
	4022  TN9310 10GbE SFP+ Ethernet Adapter
	4024  TN9210 10GBase-T Ethernet Adapter
	4025  TN9510 10GBase-T/NBASE-T Ethernet Adapter
	4026  TN9610 10GbE SFP+ Ethernet Adapter
	4027  TN9710P 10GBase-T/NBASE-T Ethernet Adapter
	4527  TN9710Q 5GBase-T/NBASE-T Ethernet Adapter
1fcc  StreamLabs
	f416  MS416
	fb01  MH4LM
1fce  Cognio Inc.
	0001  Spectrum Analyzer PC Card (SAgE)
1fd4  SUNIX Co., Ltd.
	0001  Matrix multiport serial adapter
	1999  Multiport serial controller
2000  Smart Link Ltd.
	2800  SmartPCI2800 V.92 PCI Soft DFT
2001  Temporal Research Ltd
2003  Smart Link Ltd.
	8800  LM-I56N
2004  Smart Link Ltd.
20f4  TRENDnet
2116  ZyDAS Technology Corp.
21c3  21st Century Computer Corp.
2304  Colorgraphic Communications Corp.
2348  Racore
	2010  8142 100VG/AnyLAN
2646  Kingston Technologies
270b  Xantel Corporation
270f  Chaintech Computer Co. Ltd
2711  AVID Technology Inc.
2955  Connectix Virtual PC
	6e61  OHCI USB 1.1 controller
2a15  3D Vision(???)
2bd8  ROPEX Industrie-Elektronik GmbH
3000  Hansol Electronics Inc.
3112  Satelco Ingenieria S.A.
3130  AUDIOTRAK
3142  Post Impression Systems.
31ab  Zonet
	1faa  ZEW1602 802.11b/g Wireless Adapter
3388  Hint Corp
	0013  HiNT HC4 PCI to ISDN bridge, Multimedia audio controller
	0014  HiNT HC4 PCI to ISDN bridge, Network controller
	0020  HB6 Universal PCI-PCI bridge (transparent mode)
	0021  HB6 Universal PCI-PCI bridge (non-transparent mode)
	0022  HiNT HB4 PCI-PCI Bridge (PCI6150)
	0026  HB2 PCI-PCI Bridge
	1014  AudioTrak Maya
	1018  Audiotrak INCA88
	1019  Miditrak 2120
	101a  E.Band [AudioTrak Inca88]
	101b  E.Band [AudioTrak Inca88]
	8011  VXPro II Chipset
	8012  VXPro II Chipset
	8013  VXPro II IDE
	a103  Blackmagic Design DeckLink HD Pro
3411  Quantum Designs (H.K.) Inc
3442  Bihl+Wiedemann GmbH
	1783  AS-i 3.0 cPCI Master
	1922  AS-i 3.0 PCI Master
3475  Arastra Inc.
3513  ARCOM Control Systems Ltd
37d9  ITD Firm ltd.
	1138  SCHD-PH-8 Phase detector
	1140  VR-12-PCI
	1141  PCI-485(422)
	1142  PCI-CAN2
3842  eVga.com. Corp.
38ef  4Links
3d3d  3DLabs
	0001  GLINT 300SX
	0002  GLINT 500TX
	0003  GLINT Delta
	0004  Permedia
	0005  Permedia
	0006  GLINT MX
	0007  3D Extreme
	0008  GLINT Gamma G1
	0009  Permedia II 2D+3D
	000a  GLINT R3
	000c  GLINT R3 [Oxygen VX1]
	000d  GLint R4 rev A
	000e  GLINT Gamma G2
	0011  GLint R4 rev B
	0012  GLint R5 rev A
	0013  GLint R5 rev B
	0020  VP10 visual processor
	0022  VP10 visual processor
	0024  VP9 visual processor
	002c  Wildcat Realizm 100/200
	0030  Wildcat Realizm 800
	0032  Wildcat Realizm 500
	0100  Permedia II 2D+3D
	07a1  Wildcat III 6210
	07a2  Sun XVR-500 Graphics Accelerator
	07a3  Wildcat IV 7210
	1004  Permedia
	3d04  Permedia
	ffff  Glint VGA
4005  Avance Logic Inc.
	0300  ALS300 PCI Audio Device
	0308  ALS300+ PCI Audio Device
	0309  PCI Input Controller
	1064  ALG-2064
	2064  ALG-2064i
	2128  ALG-2364A GUI Accelerator
	2301  ALG-2301
	2302  ALG-2302
	2303  AVG-2302 GUI Accelerator
	2364  ALG-2364A
	2464  ALG-2464
	2501  ALG-2564A/25128A
	4000  ALS4000 Audio Chipset
	4710  ALC200/200P
4033  Addtron Technology Co, Inc.
	1360  RTL8139 Ethernet
4040  NetXen Incorporated
	0001  NXB-10GXSR 10-Gigabit Ethernet PCIe Adapter with SR-XFP optical interface
	0002  NXB-10GCX4 10-Gigabit Ethernet PCIe Adapter with CX4 copper interface
	0003  NXB-4GCU Quad Gigabit Ethernet PCIe Adapter with 1000-BASE-T interface
	0004  BladeCenter-H 10-Gigabit Ethernet High Speed Daughter Card
	0005  NetXen Dual Port 10GbE Multifunction Adapter for c-Class
	0024  XG Mgmt
	0025  XG Mgmt
	0100  NX3031 Multifunction 1/10-Gigabit Server Adapter
4143  Digital Equipment Corp
4144  Alpha Data
	0044  ADM-XRCIIPro
4150  ONA Electroerosion
	0001  PCI32TLITE FILSTRUP1 PCI to VME Bridge Controller
	0006  PCI32TLITE UART 16550 Opencores
	0007  PCI32TLITE CAN Controller Opencores
415a  Auzentech, Inc.
416c  Aladdin Knowledge Systems
	0100  AladdinCARD
	0200  CPC
4254  DVBSky
4321  Tata Power Strategic Electronics Division
4348  WCH.CN
	2273  CH351 PCI Dual Serial Port Controller
	3253  CH352 PCI Dual Serial Port Controller
	3453  CH353 PCI Quad Serial Port Controller
	5053  CH352 PCI Serial and Parallel Port Controller
	7053  CH353 PCI Dual Serial and Parallel Ports Controller
	7073  CH356 PCI Quad Serial and Parallel Ports Controller
	7173  CH355 PCI Quad Serial Port Controller
434e  CAST Navigation LLC
4444  Internext Compression Inc
	0016  iTVC16 (CX23416) Video Decoder
	0803  iTVC15 (CX23415) Video Decoder
4468  Bridgeport machines
4594  Cogetec Informatique Inc
45fb  Baldor Electric Company
4624  Budker Institute of Nuclear Physics
	adc1  ADC200ME High speed ADC
	de01  DL200ME High resolution delay line PCI based card
	de02  DL200ME Middle resolution delay line PCI based card
4651  TXIC
4680  Umax Computer Corp
4843  Hercules Computer Technology Inc
4916  RedCreek Communications Inc
	1960  RedCreek PCI adapter
4943  Growth Networks
494f  ACCES I/O Products, Inc.
	0508  PCI-IDO-16A FET Output Card
	0518  PCI-IDO-32A FET Output Card
	0520  PCI-IDO-48 FET Output Card
	0521  PCI-IDO-48A FET Output Card
	0703  PCIe-RO-4 Electromechanical Relay Output Card
	07d0  PCIe-IDO-24 FET Output Card
	0920  PCI-IDI-48 Isolated Digital Input Card
	0bd0  PCIe-IDI-24 Isolated Digital Input Card
	0c50  PCI-DIO-24H 1x 8255 Digital Input / Output Card
	0c51  PCI-DIO-24D 1x 8255 Digital Input / Output Card
	0c52  PCIe-DIO-24 1x 8255 Digital Input / Output Card
	0c53  PCIe-DIO-24H 8255 Digital Input / Output Card
	0c57  mPCIe-DIO-24 8255 Digital Input / Output Card
	0c60  PCI-DIO-48H 8255 Digital Input / Output Card
	0c61  PCIe-DIO-48 8255 Digital Input / Output Card
	0c62  P104-DIO-48 8255 Digital Input / Output Card
	0c68  PCI-DIO-72 8255 Digital Input / Output Card
	0c69  P104-DIO-96 8255 Digital Input / Output Card
	0c70  PCI-DIO-96 8255 Digital Input / Output Card
	0c78  PCI-DIO-120 8255 Digital Input / Output Card
	0dc8  PCI-IDIO-16 Isolated Digital Input / FET Output Card
	0e50  PCI-DIO-24S 8255 Digital Input / Output Card
	0e51  PCI-DIO-24H(C) 8255 Digital Input / Output Card
	0e52  PCI-DIO-24D(C) 8255 Digital Input / Output Card
	0e53  PCIe-DIO-24S 8255 Digital Input / Output Card
	0e54  PCIe-DIO-24HS 8255 Digital Input / Output Card
	0e55  PCIe-DIO-24DC 8255 Digital Input / Output Card
	0e56  PCIe-DIO-24DCS 8255 Digital Input / Output Card
	0e57  mPCIe-DIO-24S 8255 Digital Input / Output Card
	0e60  PCI-DIO-48S 2x 8255 Digital Input / Output Card
	0e61  PCIe-DIO-48S 2x 8255 Digital Input / Output Card
	0e62  P104-DIO-48S 2x 8255 Digital Input / Output Card
	0f00  PCI-IIRO-8 Isolated Digital / Relay Output Card
	0f01  LPCI-IIRO-8 Isolated Digital / Relay Output Card
	0f02  PCIe-IIRO-8 Isolated Digital / Relay Output Card
	0f08  PCI-IIRO-16 Isolated Digital / Relay Output Card
	0f09  PCIe-IIRO-16 Isolated Digital / Relay Output Card
	0fc0  PCIe-IDIO-12 Isolated Digital Input / FET Output Card
	0fc1  PCIe-IDI-12 Isolated Digital Input Card
	0fc2  PCIe-IDO-12 FET Output Card
	0fd0  PCIe-IDIO-24 Isolated Digital Input / FET Output Card
	1050  PCI-422/485-2 2x RS422/RS484 Card
	1051  PCIe-COM-2SRJ 2x RS422/RS484 Card w/RJ45 Connectors
	1052  104I-COM-2S 2x RS422/RS484 PCI/104 Board
	1053  mPCIe-COM-2S 2x RS422/RS484 PCI Express Mini Card
	1058  PCI-COM422/4 4x RS422 Card
	1059  PCI-COM485/4 4x RS485 Card
	105a  PCIe-COM422-4 4x RS422 Card
	105b  PCIe-COM485-4 4x RS485 Card
	105c  PCIe-COM-4SRJ 4x RS422/RS485 Card w/RJ45 Connectors
	105d  104I-COM-4S 4x RS422/RS484 PCI/104 Board
	105e  mPCIe-COM-4S 4x RS422/RS484 PCI Express Mini Card
	1068  PCI-COM422/8 8x RS422 Card
	1069  PCI-COM485/8 8x RS485 Card
	106a  PCIe-COM422-8 8x RS422 Card
	106b  PCIe-COM485-8 8x RS485 Card
	106c  104I-COM-8S 8x RS422/RS485 PCI/104 Board
	1088  PCI-COM232/1 1x RS232 Card
	1090  PCI-COM232/2 2x RS232 Card
	1091  PCIe-COM232-2RJ 2x RS232 Card w/RJ45 Connectors
	1093  mPCIe-COM232-2 2x RS232 PCI Express Mini Card
	1098  PCIe-COM232-4 4x RS232 Card
	1099  PCIe-COM232-4RJ 4x RS232 Card w/RJ45 Connectors
	109b  mPCIe-COM232-4 4x RS232 PCI Express Mini Card
	10a8  P104-COM232-8 8x RS232 PC-104+ Board
	10a9  PCIe-COM232-8 8x RS232 Card
	10c9  PCI-COM-1S 1x RS422/RS485 Card
	10d0  PCI-COM2S 2x RS422/RS485 Card
	10d1  PCIe-COM-2SMRJ 2x RS232/RS422/RS485 Card w/RJ45 Connectors
	10d2  104I-COM-2SM 2x RS232/RS422/RS485 PCI/104 Board
	10d3  mPCIe-COM-2SM 2x RS232/RS422/RS485 PCI Express Mini Card
	10d8  PCI-COM-4SM 4x RS232/RS422/RS485 Card
	10d9  PCIe-COM-4SM 4x RS232/RS422/RS485 Card
	10da  PCIe-COM-4SMRJ 4x RS232/RS422/RS485 Card w/RJ45 Connectors
	10db  104I-COM-4SM 4x RS232/RS422/RS485 PCI/104 Board
	10dc  mPCIe-COM-4SM 4x RS232/RS422/RS485 PCI Express Mini Card
	10e8  PCI-COM-8SM 8x RS232/RS422/RS485 Card
	10e9  PCIe-COM-8SM 8x RS232/RS422/RS485 Card
	10ea  104I-COM-8SM 8x RS232/RS422/RS485 PCI-104 Board
	1108  mPCIe-ICM485-1 1x Isolated RS485 PCI Express Mini Card
	1110  mPCIe-ICM422-2 2x Isolated RS422 PCI Express Mini Card
	1111  mPCIe-ICM485-2 2x Isolated RS485 PCI Express Mini Card
	1118  mPCIe-ICM422-4 4x Isolated RS422 PCI Express Mini Card
	1119  mPCIe-ICM485-4 4x Isolated RS485 PCI Express Mini Card
	1148  PCI-ICM-1S 1x Isolated RS422/RS485 Card
	1150  PCI-ICM-2S 2x Isolated RS422/RS485 Card
	1152  PCIe-ICM-2S 2x Isolated RS422/RS485 Card
	1158  PCI-ICM422/4 4x Isolated RS422 Card
	1159  PCI-ICM485/4 4x Isolated RS485 Card
	115a  PCIe-ICM-4S 4x Isolated RS422/RS485 Card
	1190  PCIe-ICM232-2 2x Isolated RS232 Card
	1191  mPCIe-ICM232-2 2x Isolated RS232 PCI Express Mini Card
	1198  PCIe-ICM232-4 4x Isolated RS232 Card
	1199  mPCIe-ICM232-4 4x Isolated RS422 PCI Express Mini Card
	11d0  PCIe-ICM-2SM 2x Isolated RS232/RS422/RS485 Card
	11d8  PCIe-ICM-4SM 4x Isolated RS232/RS422/RS485 Card
	1250  PCI-WDG-2S Watchdog and 2x Serial Card
	12d0  PCI-WDG-IMPAC
	2230  PCI-QUAD-8 8x Quadrature Input Card
	2231  PCI-QUAD-4 4x Quadrature Input Card
	22c0  PCI-WDG-CSM Watchdog Card
	25c0  P104-WDG-E Watchdog PC/104+ Board
	2c50  PCI-DIO-96CT 96x Digital Input / Output Card
	2c58  PCI-DIO-96C3 96x Digital Input / Output Card w/3x 8254 Counter Card
	2ee0  PCIe-DIO24S-CTR12 24x Digital Input / Output Card w/4x 8254 Counter Card
	2fc0  P104-WDG-CSM Watchdog PC/104+ Board
	2fc1  P104-WDG-CSMA Advanced Watchdog PC/104+ Board
	5ed0  PCI-DAC
	6c90  PCI-DA12-2 2x 12-bit Analog Output Card
	6c98  PCI-DA12-4 4x 12-bit Analog Output Card
	6ca0  PCI-DA12-6 6x 12-bit Analog Output Card
	6ca8  PCI-DA12-8 8x 12-bit Analog Output Card
	6ca9  PCI-DA12-8V
	6cb0  PCI-DA12-16 16x 12-bit Analog Output Card
	6cb1  PCI-DA12-16V
	8ef0  P104-FAS16-16
	aca8  PCI-AI12-16 12-bit 100kHz Analog Input Card
	aca9  PCI-AI12-16A 12-bit 100kHz Analog Input w/FIFO Card
	eca8  PCI-AIO12-16 12-bit 100kHz Analog Input w/2x Analog Output and FIFO Card
	ecaa  PCI-A12-16A 12-bit 100kHz Analog Input w/2x Analog Output and FIFO Card
	ece8  LPCI-A16-16A 16-bit 500kHz Analog Input low-profile Card
	ece9  LPCI-AIO16A 16-bit 500kHz Analog Input low-profile Card
4978  Axil Computer Inc
4a14  NetVin
	5000  NV5000SC
4b10  Buslogic Inc.
4c48  LUNG HWA Electronics
4c53  SBS Technologies
	0000  PLUSTEST device
	0001  PLUSTEST-MM device
4ca1  Seanix Technology Inc
4d51  MediaQ Inc.
	0200  MQ-200
4d54  Microtechnica Co Ltd
4d56  MATRIX VISION GmbH
	0000  [mvHYPERION-CLe/CLb] CameraLink PCI Express x1 Frame Grabber
	0001  [mvHYPERION-CLf/CLm] CameraLink PCI Express x4 Frame Grabber
	0010  [mvHYPERION-16R16/-32R16] 16 Video Channel PCI Express x4 Frame Grabber
	0020  [mvHYPERION-HD-SDI] HD-SDI PCI Express x4 Frame Grabber
	0030  [mvHYPERION-HD-SDI-Merger] HD-SDI PCI Express x4 Frame Grabber
4ddc  ILC Data Device Corp
	0100  DD-42924I5-300 (ARINC 429 Data Bus)
	0801  BU-65570I1 MIL-STD-1553 Test and Simulation
	0802  BU-65570I2 MIL-STD-1553 Test and Simulation
	0811  BU-65572I1 MIL-STD-1553 Test and Simulation
	0812  BU-65572I2 MIL-STD-1553 Test and Simulation
	0881  BU-65570T1 MIL-STD-1553 Test and Simulation
	0882  BU-65570T2 MIL-STD-1553 Test and Simulation
	0891  BU-65572T1 MIL-STD-1553 Test and Simulation
	0892  BU-65572T2 MIL-STD-1553 Test and Simulation
	0901  BU-65565C1 MIL-STD-1553 Data Bus
	0902  BU-65565C2 MIL-STD-1553 Data Bus
	0903  BU-65565C3 MIL-STD-1553 Data Bus
	0904  BU-65565C4 MIL-STD-1553 Data Bus
	0b01  BU-65569I1 MIL-STD-1553 Data Bus
	0b02  BU-65569I2 MIL-STD-1553 Data Bus
	0b03  BU-65569I3 MIL-STD-1553 Data Bus
	0b04  BU-65569I4 MIL-STD-1553 Data Bus
5045  University of Toronto
	4243  BLASTbus PCI Interface Card v1
5046  GemTek Technology Corporation
	1001  PCI Radio
5053  Voyetra Technologies
	2010  Daytona Audio Adapter
50b2  TerraTec Electronic GmbH
5136  S S Technologies
5143  Qualcomm Inc
5145  Ensoniq (Old)
	3031  Concert AudioPCI
5168  Animation Technologies Inc.
	0300  FlyDVB-S
	0301  FlyDVB-T
5301  Alliance Semiconductor Corp.
	0001  ProMotion aT3D
5333  S3 Graphics Ltd.
	0551  Plato/PX (system)
	5631  86c325 [ViRGE]
	8800  86c866 [Vision 866]
	8801  86c964 [Vision 964]
	8810  86c764_0 [Trio 32 vers 0]
	8811  86c764/765 [Trio32/64/64V+]
	8812  86cM65 [Aurora64V+]
	8813  86c764_3 [Trio 32/64 vers 3]
	8814  86c767 [Trio 64UV+]
	8815  86cM65 [Aurora 128]
	883d  86c988 [ViRGE/VX]
	8870  FireGL
	8880  86c868 [Vision 868 VRAM] vers 0
	8881  86c868 [Vision 868 VRAM] vers 1
	8882  86c868 [Vision 868 VRAM] vers 2
	8883  86c868 [Vision 868 VRAM] vers 3
	88b0  86c928 [Vision 928 VRAM] vers 0
	88b1  86c928 [Vision 928 VRAM] vers 1
	88b2  86c928 [Vision 928 VRAM] vers 2
	88b3  86c928 [Vision 928 VRAM] vers 3
	88c0  86c864 [Vision 864 DRAM] vers 0
	88c1  86c864 [Vision 864 DRAM] vers 1
	88c2  86c864 [Vision 864-P DRAM] vers 2
	88c3  86c864 [Vision 864-P DRAM] vers 3
	88d0  86c964 [Vision 964 VRAM] vers 0
	88d1  86c964 [Vision 964 VRAM] vers 1
	88d2  86c964 [Vision 964-P VRAM] vers 2
	88d3  86c964 [Vision 964-P VRAM] vers 3
	88f0  86c968 [Vision 968 VRAM] rev 0
	88f1  86c968 [Vision 968 VRAM] rev 1
	88f2  86c968 [Vision 968 VRAM] rev 2
	88f3  86c968 [Vision 968 VRAM] rev 3
	8900  86c755 [Trio 64V2/DX]
	8901  86c775/86c785 [Trio 64V2/DX or /GX]
	8902  Plato/PX
	8903  Trio 3D business multimedia
	8904  86c365, 86c366 [Trio 3D]
	8905  Trio 64V+ family
	8906  Trio 64V+ family
	8907  Trio 64V+ family
	8908  Trio 64V+ family
	8909  Trio 64V+ family
	890a  Trio 64V+ family
	890b  Trio 64V+ family
	890c  Trio 64V+ family
	890d  Trio 64V+ family
	890e  Trio 64V+ family
	890f  Trio 64V+ family
	8a01  86c375 [ViRGE/DX] or 86c385 [ViRGE/GX]
	8a10  ViRGE/GX2
	8a13  86c360 [Trio 3D/1X], 86c362, 86c368 [Trio 3D/2X]
	8a20  86c794 [Savage 3D]
	8a21  86c390 [Savage 3D/MV]
	8a22  Savage 4
	8a23  Savage 4
	8a25  ProSavage PM133
	8a26  ProSavage KM133
	8c00  ViRGE/M3
	8c01  ViRGE/MX
	8c02  ViRGE/MX+
	8c03  ViRGE/MX+MV
	8c10  86C270-294 [SavageMX-MV]
	8c11  82C270-294 [SavageMX]
	8c12  86C270-294 [SavageIX-MV]
	8c13  86C270-294 [SavageIX]
	8c22  SuperSavage MX/128
	8c24  SuperSavage MX/64
	8c26  SuperSavage MX/64C
	8c2a  SuperSavage IX/128 SDR
	8c2b  SuperSavage IX/128 DDR
	8c2c  SuperSavage IX/64 SDR
	8c2d  SuperSavage IX/64 DDR
	8c2e  SuperSavage IX/C SDR
	8c2f  SuperSavage IX/C DDR
	8d01  86C380 [ProSavageDDR K4M266]
	8d02  VT8636A [ProSavage KN133] AGP4X VGA Controller (TwisterK)
	8d03  VT8751 [ProSavageDDR P4M266]
	8d04  VT8375 [ProSavage8 KM266/KL266]
	8e00  DeltaChrome
	8e26  ProSavage
	8e40  2300E Graphics Processor
	8e48  Matrix [Chrome S25 / S27]
	9043  Chrome 430 GT
	9045  Chrome 430 ULP / 435 ULP / 440 GTX
	9060  Chrome 530 GT
	9102  86C410 [Savage 2000]
	ca00  SonicVibes
5431  AuzenTech, Inc.
544c  Teralogic Inc
	0350  TL880-based HDTV/ATSC tuner
544d  TBS Technologies
	6178  DVB-S2 4 Tuner PCIe Card
5452  SCANLAB AG
	3443  RTC4
5455  Technische University Berlin
	4458  S5933
5456  GoTView
5519  Cnet Technologies, Inc.
5544  Dunord Technologies
	0001  I-30xx Scanner Interface
5555  Genroco, Inc
	0003  TURBOstor HFP-832 [HiPPI NIC]
	3b00  Epiphan DVI2PCIe video capture card
5646  Vector Fabrics BV
5654  VoiceTronix Pty Ltd
5678  Dawicontrol Computersysteme GmbH
5700  Netpower
5845  X-ES, Inc.
584d  AuzenTech Co., Ltd.
5851  Exacq Technologies
	8008  tDVR8008 8-port video capture card
	8016  tDVR8016 16-chan video capture card
	8032  tDVR8032 32-chan video capture card
5853  XenSource, Inc.
	0001  Xen Platform Device
	c000  Citrix XenServer PCI Device for Windows Update
	c110  Virtualized HID
	c147  Virtualized Graphics Device
5854  GoTView
5ace  Beholder International Ltd.
6205  TBS Technologies (wrong ID)
6209  TBS Technologies (wrong ID)
631c  SmartInfra Ltd
	1652  PXI-1652 Signal Generator
	2504  PXI-2504 Signal Interrogator
6356  UltraStor
6374  c't Magazin fuer Computertechnik
	6773  GPPCI
6409  Logitec Corp.
6549  Teradici Corp.
	1200  TERA1200 PC-over-IP Host
6666  Decision Computer International Co.
	0001  PCCOM4
	0002  PCCOM8
	0004  PCCOM2
	0101  PCI 8255/8254 I/O Card
	0200  12-bit AD/DA Card
	0201  14-bit AD/DA Card
	1011  Industrial Card
	1021  8 photo couple 8 relay Card
	1022  4 photo couple 4 relay Card
	1025  16 photo couple 16 relay Card
	4000  WatchDog Card
6688  Zycoo Co., Ltd
	1200  CooVox TDM Analog Module
	1400  CooVOX TDM GSM Module
	1600  CooVOX TDM E1/T1 Module
	1800  CooVOX TDM BRI Module
6900  Red Hat, Inc.
7063  pcHDTV
	2000  HD-2000
	3000  HD-3000
	5500  HD5500 HDTV
7284  HT OMEGA Inc.
7401  EndRun Technologies
	e100  PTP3100 PCIe PTP Slave Clock
7470  TP-LINK Technologies Co., Ltd.
7604  O.N. Electronic Co Ltd.
7bde  MIDAC Corporation
7fed  PowerTV
8008  Quancom Electronic GmbH
	0010  WDOG1 [PCI-Watchdog 1]
	0011  PWDOG2 [PCI-Watchdog 2]
	0015  Clock77/PCI & Clock77/PCIe (DCF-77 receiver)
807d  Asustek Computer, Inc.
8086  Intel Corporation
	0007  82379AB
	0008  Extended Express System Support Controller
	0039  21145 Fast Ethernet
	0040  Core Processor DRAM Controller
	0041  Core Processor PCI Express x16 Root Port
	0042  Core Processor Integrated Graphics Controller
	0043  Core Processor Secondary PCI Express Root Port
	0044  Core Processor DRAM Controller
	0045  Core Processor PCI Express x16 Root Port
	0046  Core Processor Integrated Graphics Controller
	0047  Core Processor Secondary PCI Express Root Port
	0048  Core Processor DRAM Controller
	0049  Core Processor PCI Express x16 Root Port
	004a  Core Processor Integrated Graphics Controller
	004b  Core Processor Secondary PCI Express Root Port
	0050  Core Processor Thermal Management Controller
	0069  Core Processor DRAM Controller
	0082  Centrino Advanced-N 6205 [Taylor Peak]
	0083  Centrino Wireless-N 1000 [Condor Peak]
	0084  Centrino Wireless-N 1000 [Condor Peak]
	0085  Centrino Advanced-N 6205 [Taylor Peak]
	0087  Centrino Advanced-N + WiMAX 6250 [Kilmer Peak]
	0089  Centrino Advanced-N + WiMAX 6250 [Kilmer Peak]
	008a  Centrino Wireless-N 1030 [Rainbow Peak]
	008b  Centrino Wireless-N 1030 [Rainbow Peak]
	0090  Centrino Advanced-N 6230 [Rainbow Peak]
	0091  Centrino Advanced-N 6230 [Rainbow Peak]
	0100  2nd Generation Core Processor Family DRAM Controller
	0101  Xeon E3-1200/2nd Generation Core Processor Family PCI Express Root Port
	0102  2nd Generation Core Processor Family Integrated Graphics Controller
	0104  2nd Generation Core Processor Family DRAM Controller
	0105  Xeon E3-1200/2nd Generation Core Processor Family PCI Express Root Port
	0106  2nd Generation Core Processor Family Integrated Graphics Controller
	0108  Xeon E3-1200 Processor Family DRAM Controller
	0109  Xeon E3-1200/2nd Generation Core Processor Family PCI Express Root Port
	010a  Xeon E3-1200 Processor Family Integrated Graphics Controller
	010b  Xeon E3-1200/2nd Generation Core Processor Family Integrated Graphics Controller
	010c  Xeon E3-1200/2nd Generation Core Processor Family DRAM Controller
	010d  Xeon E3-1200/2nd Generation Core Processor Family PCI Express Root Port
	010e  Xeon E3-1200/2nd Generation Core Processor Family Integrated Graphics Controller
	0112  2nd Generation Core Processor Family Integrated Graphics Controller
	0116  2nd Generation Core Processor Family Integrated Graphics Controller
	0122  2nd Generation Core Processor Family Integrated Graphics Controller
	0126  2nd Generation Core Processor Family Integrated Graphics Controller
	0150  Xeon E3-1200 v2/3rd Gen Core processor DRAM Controller
	0151  Xeon E3-1200 v2/3rd Gen Core processor PCI Express Root Port
	0152  Xeon E3-1200 v2/3rd Gen Core processor Graphics Controller
	0153  3rd Gen Core Processor Thermal Subsystem
	0154  3rd Gen Core processor DRAM Controller
	0155  Xeon E3-1200 v2/3rd Gen Core processor PCI Express Root Port
	0156  3rd Gen Core processor Graphics Controller
	0158  Xeon E3-1200 v2/Ivy Bridge DRAM Controller
	0159  Xeon E3-1200 v2/3rd Gen Core processor PCI Express Root Port
	015a  Xeon E3-1200 v2/Ivy Bridge Graphics Controller
	015c  Xeon E3-1200 v2/3rd Gen Core processor DRAM Controller
	015d  Xeon E3-1200 v2/3rd Gen Core processor PCI Express Root Port
	015e  Xeon E3-1200 v2/3rd Gen Core processor Graphics Controller
	0162  Xeon E3-1200 v2/3rd Gen Core processor Graphics Controller
	0166  3rd Gen Core processor Graphics Controller
	016a  Xeon E3-1200 v2/3rd Gen Core processor Graphics Controller
	0172  Xeon E3-1200 v2/3rd Gen Core processor Graphics Controller
	0176  3rd Gen Core processor Graphics Controller
	0309  80303 I/O Processor PCI-to-PCI Bridge
	030d  80312 I/O Companion Chip PCI-to-PCI Bridge
	0326  6700/6702PXH I/OxAPIC Interrupt Controller A
	0327  6700PXH I/OxAPIC Interrupt Controller B
	0329  6700PXH PCI Express-to-PCI Bridge A
	032a  6700PXH PCI Express-to-PCI Bridge B
	032c  6702PXH PCI Express-to-PCI Bridge A
	0330  80332 [Dobson] I/O processor (A-Segment Bridge)
	0331  80332 [Dobson] I/O processor (A-Segment IOAPIC)
	0332  80332 [Dobson] I/O processor (B-Segment Bridge)
	0333  80332 [Dobson] I/O processor (B-Segment IOAPIC)
	0334  80332 [Dobson] I/O processor (ATU)
	0335  80331 [Lindsay] I/O processor (PCI-X Bridge)
	0336  80331 [Lindsay] I/O processor (ATU)
	0340  41210 [Lanai] Serial to Parallel PCI Bridge (A-Segment Bridge)
	0341  41210 [Lanai] Serial to Parallel PCI Bridge (B-Segment Bridge)
	0370  80333 Segment-A PCIe Express to PCI-X bridge
	0371  80333 A-Bus IOAPIC
	0372  80333 Segment-B PCIe Express to PCI-X bridge
	0373  80333 B-Bus IOAPIC
	0374  80333 Address Translation Unit
	0402  Xeon E3-1200 v3/4th Gen Core Processor Integrated Graphics Controller
	0406  4th Gen Core Processor Integrated Graphics Controller
	040a  Xeon E3-1200 v3 Processor Integrated Graphics Controller
	0412  Xeon E3-1200 v3/4th Gen Core Processor Integrated Graphics Controller
	0416  4th Gen Core Processor Integrated Graphics Controller
	041a  Xeon E3-1200 v3 Processor Integrated Graphics Controller
	041e  4th Generation Core Processor Family Integrated Graphics Controller
	0434  DH89XXCC Series QAT
	0435  DH895XCC Series QAT
	0436  DH8900CC Null Device
	0438  DH8900CC Series Gigabit Network Connection
	043a  DH8900CC Series Gigabit Fiber Network Connection
	043c  DH8900CC Series Gigabit Backplane Network Connection
	0440  DH8900CC Series Gigabit SFP Network Connection
	0442  DH89XXCC Series QAT Virtual Function
	0443  DH895XCC Series QAT Virtual Function
	0482  82375EB/SB PCI to EISA Bridge
	0483  82424TX/ZX [Saturn] CPU to PCI bridge
	0484  82378ZB/IB, 82379AB (SIO, SIO.A) PCI to ISA Bridge
	0486  82425EX/ZX [Aries] PCIset with ISA bridge
	04a3  82434LX/NX [Mercury/Neptune] Processor to PCI bridge
	04d0  82437FX [Triton FX]
	0500  E8870 Processor bus control
	0501  E8870 Memory controller
	0502  E8870 Scalability Port 0
	0503  E8870 Scalability Port 1
	0510  E8870IO Hub Interface Port 0 registers (8-bit compatibility port)
	0511  E8870IO Hub Interface Port 1 registers
	0512  E8870IO Hub Interface Port 2 registers
	0513  E8870IO Hub Interface Port 3 registers
	0514  E8870IO Hub Interface Port 4 registers
	0515  E8870IO General SIOH registers
	0516  E8870IO RAS registers
	0530  E8870SP Scalability Port 0 registers
	0531  E8870SP Scalability Port 1 registers
	0532  E8870SP Scalability Port 2 registers
	0533  E8870SP Scalability Port 3 registers
	0534  E8870SP Scalability Port 4 registers
	0535  E8870SP Scalability Port 5 registers
	0536  E8870SP Interleave registers 0 and 1
	0537  E8870SP Interleave registers 2 and 3
	0600  RAID Controller
	061f  80303 I/O Processor
	0700  CE Media Processor A/V Bridge
	0701  CE Media Processor NAND Flash Controller
	0703  CE Media Processor Media Control Unit 1
	0704  CE Media Processor Video Capture Interface
	0707  CE Media Processor SPI Slave
	0708  CE Media Processor 4100
	0800  Moorestown SPI Ctrl 0
	0801  Moorestown SPI Ctrl 1
	0802  Moorestown I2C 0
	0803  Moorestown I2C 1
	0804  Moorestown I2C 2
	0805  Moorestown Keyboard Ctrl
	0806  Moorestown USB Ctrl
	0807  Moorestown SD Host Ctrl 0
	0808  Moorestown SD Host Ctrl 1
	0809  Moorestown NAND Ctrl
	080a  Moorestown Audio Ctrl
	080b  Moorestown ISP
	080c  Moorestown Security Controller
	080d  Moorestown External Displays
	080e  Moorestown SCU IPC
	080f  Moorestown GPIO Controller
	0810  Moorestown Power Management Unit
	0811  Moorestown OTG Ctrl
	0812  Moorestown SPI Ctrl 2
	0813  Moorestown SC DMA
	0814  Moorestown LPE DMA
	0815  Moorestown SSP0
	0817  Medfield Serial IO I2C Controller #3
	0818  Medfield Serial IO I2C Controller #4
	0819  Medfield Serial IO I2C Controller #5
	081a  Medfield GPIO Controller [Core]
	081b  Medfield Serial IO HSUART Controller #1
	081c  Medfield Serial IO HSUART Controller #2
	081d  Medfield Serial IO HSUART Controller #3
	081e  Medfield Serial IO HSUART DMA Controller
	081f  Medfield GPIO Controller [AON]
	0820  Medfield SD Host Controller
	0821  Medfield SDIO Controller #1
	0822  Medfield SDIO Controller #2
	0823  Medfield eMMC Controller #0
	0824  Medfield eMMC Controller #1
	0827  Medfield Serial IO DMA Controller
	0828  Medfield Power Management Unit
	0829  Medfield USB Device Controller (OTG)
	082a  Medfield SCU IPC
	082c  Medfield Serial IO I2C Controller #0
	082d  Medfield Serial IO I2C Controller #1
	082e  Medfield Serial IO I2C Controller #2
	0885  Centrino Wireless-N + WiMAX 6150
	0886  Centrino Wireless-N + WiMAX 6150
	0887  Centrino Wireless-N 2230
	0888  Centrino Wireless-N 2230
	088e  Centrino Advanced-N 6235
	088f  Centrino Advanced-N 6235
	0890  Centrino Wireless-N 2200
	0891  Centrino Wireless-N 2200
	0892  Centrino Wireless-N 135
	0893  Centrino Wireless-N 135
	0894  Centrino Wireless-N 105
	0895  Centrino Wireless-N 105
	0896  Centrino Wireless-N 130
	0897  Centrino Wireless-N 130
	08ae  Centrino Wireless-N 100
	08af  Centrino Wireless-N 100
	08b1  Wireless 7260
	08b2  Wireless 7260
	08b3  Wireless 3160
	08b4  Wireless 3160
	08cf  Atom Processor Z2760 Integrated Graphics Controller
	0953  PCIe Data Center SSD
	095a  Wireless 7265
	095b  Wireless 7265
	0960  80960RP (i960RP) Microprocessor/Bridge
	0962  80960RM (i960RM) Bridge
	0964  80960RP (i960RP) Microprocessor/Bridge
	0a03  Haswell-ULT Thermal Subsystem
	0a04  Haswell-ULT DRAM Controller
	0a06  Haswell-ULT Integrated Graphics Controller
	0a0c  Haswell-ULT HD Audio Controller
	0a16  Haswell-ULT Integrated Graphics Controller
	0a22  Haswell-ULT Integrated Graphics Controller
	0a26  Haswell-ULT Integrated Graphics Controller
	0a2a  Haswell-ULT Integrated Graphics Controller
	0a2e  Haswell-ULT Integrated Graphics Controller
	0a53  DC P3520 SSD
	0a54  Express Flash NVMe P4500
	0a55  Express Flash NVMe P4600
	0be0  Atom Processor D2xxx/N2xxx Integrated Graphics Controller
	0be1  Atom Processor D2xxx/N2xxx Integrated Graphics Controller
	0be2  Atom Processor D2xxx/N2xxx Integrated Graphics Controller
	0be3  Atom Processor D2xxx/N2xxx Integrated Graphics Controller
	0be4  Atom Processor D2xxx/N2xxx Integrated Graphics Controller
	0be5  Atom Processor D2xxx/N2xxx Integrated Graphics Controller
	0be6  Atom Processor D2xxx/N2xxx Integrated Graphics Controller
	0be7  Atom Processor D2xxx/N2xxx Integrated Graphics Controller
	0be8  Atom Processor D2xxx/N2xxx Integrated Graphics Controller
	0be9  Atom Processor D2xxx/N2xxx Integrated Graphics Controller
	0bea  Atom Processor D2xxx/N2xxx Integrated Graphics Controller
	0beb  Atom Processor D2xxx/N2xxx Integrated Graphics Controller
	0bec  Atom Processor D2xxx/N2xxx Integrated Graphics Controller
	0bed  Atom Processor D2xxx/N2xxx Integrated Graphics Controller
	0bee  Atom Processor D2xxx/N2xxx Integrated Graphics Controller
	0bef  Atom Processor D2xxx/N2xxx Integrated Graphics Controller
	0bf0  Atom Processor D2xxx/N2xxx DRAM Controller
	0bf1  Atom Processor D2xxx/N2xxx DRAM Controller
	0bf2  Atom Processor D2xxx/N2xxx DRAM Controller
	0bf3  Atom Processor D2xxx/N2xxx DRAM Controller
	0bf4  Atom Processor D2xxx/N2xxx DRAM Controller
	0bf5  Atom Processor D2xxx/N2xxx DRAM Controller
	0bf6  Atom Processor D2xxx/N2xxx DRAM Controller
	0bf7  Atom Processor D2xxx/N2xxx DRAM Controller
	0c00  4th Gen Core Processor DRAM Controller
	0c01  Xeon E3-1200 v3/4th Gen Core Processor PCI Express x16 Controller
	0c04  Xeon E3-1200 v3/4th Gen Core Processor DRAM Controller
	0c05  Xeon E3-1200 v3/4th Gen Core Processor PCI Express x8 Controller
	0c08  Xeon E3-1200 v3 Processor DRAM Controller
	0c09  Xeon E3-1200 v3/4th Gen Core Processor PCI Express x4 Controller
	0c0c  Xeon E3-1200 v3/4th Gen Core Processor HD Audio Controller
	0c46  Atom Processor S1200 PCI Express Root Port 1
	0c47  Atom Processor S1200 PCI Express Root Port 2
	0c48  Atom Processor S1200 PCI Express Root Port 3
	0c49  Atom Processor S1200 PCI Express Root Port 4
	0c4e  Atom Processor S1200 NTB Primary
	0c50  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QuickData Technology Device
	0c51  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QuickData Technology Device
	0c52  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QuickData Technology Device
	0c53  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QuickData Technology Device
	0c54  Atom Processor S1200 Internal
	0c55  Atom Processor S1200 DFX 1
	0c56  Atom Processor S1200 DFX 2
	0c59  Atom Processor S1200 SMBus 2.0 Controller 0
	0c5a  Atom Processor S1200 SMBus 2.0 Controller 1
	0c5b  Atom Processor S1200 SMBus Controller 2
	0c5c  Atom Processor S1200 SMBus Controller 3
	0c5d  Atom Processor S1200 SMBus Controller 4
	0c5e  Atom Processor S1200 SMBus Controller 5
	0c5f  Atom Processor S1200 UART
	0c60  Atom Processor S1200 Integrated Legacy Bus
	0c70  Atom Processor S1200 Internal
	0c71  Atom Processor S1200 Internal
	0c72  Atom Processor S1200 Internal
	0c73  Atom Processor S1200 Internal
	0c74  Atom Processor S1200 Internal
	0c75  Atom Processor S1200 Internal
	0c76  Atom Processor S1200 Internal
	0c77  Atom Processor S1200 Internal
	0c78  Atom Processor S1200 Internal
	0c79  Atom Processor S1200 Internal
	0c7a  Atom Processor S1200 Internal
	0c7b  Atom Processor S1200 Internal
	0c7c  Atom Processor S1200 Internal
	0c7d  Atom Processor S1200 Internal
	0c7e  Atom Processor S1200 Internal
	0c7f  Atom Processor S1200 Internal
	0d00  Crystal Well DRAM Controller
	0d01  Crystal Well PCI Express x16 Controller
	0d04  Crystal Well DRAM Controller
	0d05  Crystal Well PCI Express x8 Controller
	0d09  Crystal Well PCI Express x4 Controller
	0d0c  Crystal Well HD Audio Controller
	0d16  Crystal Well Integrated Graphics Controller
	0d26  Crystal Well Integrated Graphics Controller
	0d36  Crystal Well Integrated Graphics Controller
	0e00  Xeon E7 v2/Xeon E5 v2/Core i7 DMI2
	0e01  Xeon E7 v2/Xeon E5 v2/Core i7 PCI Express Root Port in DMI2 Mode
	0e02  Xeon E7 v2/Xeon E5 v2/Core i7 PCI Express Root Port 1a
	0e03  Xeon E7 v2/Xeon E5 v2/Core i7 PCI Express Root Port 1b
	0e04  Xeon E7 v2/Xeon E5 v2/Core i7 PCI Express Root Port 2a
	0e05  Xeon E7 v2/Xeon E5 v2/Core i7 PCI Express Root Port 2b
	0e06  Xeon E7 v2/Xeon E5 v2/Core i7 PCI Express Root Port 2c
	0e07  Xeon E7 v2/Xeon E5 v2/Core i7 PCI Express Root Port 2d
	0e08  Xeon E7 v2/Xeon E5 v2/Core i7 PCI Express Root Port 3a
	0e09  Xeon E7 v2/Xeon E5 v2/Core i7 PCI Express Root Port 3b
	0e0a  Xeon E7 v2/Xeon E5 v2/Core i7 PCI Express Root Port 3c
	0e0b  Xeon E7 v2/Xeon E5 v2/Core i7 PCI Express Root Port 3d
	0e10  Xeon E7 v2/Xeon E5 v2/Core i7 IIO Configuration Registers
	0e13  Xeon E7 v2/Xeon E5 v2/Core i7 IIO Configuration Registers
	0e17  Xeon E7 v2/Xeon E5 v2/Core i7 IIO Configuration Registers
	0e18  Xeon E7 v2/Xeon E5 v2/Core i7 IIO Configuration Registers
	0e1c  Xeon E7 v2/Xeon E5 v2/Core i7 IIO Configuration Registers
	0e1d  Xeon E7 v2/Xeon E5 v2/Core i7 R2PCIe
	0e1e  Xeon E7 v2/Xeon E5 v2/Core i7 UBOX Registers
	0e1f  Xeon E7 v2/Xeon E5 v2/Core i7 UBOX Registers
	0e20  Xeon E7 v2/Xeon E5 v2/Core i7 Crystal Beach DMA Channel 0
	0e21  Xeon E7 v2/Xeon E5 v2/Core i7 Crystal Beach DMA Channel 1
	0e22  Xeon E7 v2/Xeon E5 v2/Core i7 Crystal Beach DMA Channel 2
	0e23  Xeon E7 v2/Xeon E5 v2/Core i7 Crystal Beach DMA Channel 3
	0e24  Xeon E7 v2/Xeon E5 v2/Core i7 Crystal Beach DMA Channel 4
	0e25  Xeon E7 v2/Xeon E5 v2/Core i7 Crystal Beach DMA Channel 5
	0e26  Xeon E7 v2/Xeon E5 v2/Core i7 Crystal Beach DMA Channel 6
	0e27  Xeon E7 v2/Xeon E5 v2/Core i7 Crystal Beach DMA Channel 7
	0e28  Xeon E7 v2/Xeon E5 v2/Core i7 VTd/Memory Map/Misc
	0e29  Xeon E7 v2/Xeon E5 v2/Core i7 Memory Hotplug
	0e2a  Xeon E7 v2/Xeon E5 v2/Core i7 IIO RAS
	0e2c  Xeon E7 v2/Xeon E5 v2/Core i7 IOAPIC
	0e2e  Xeon E7 v2/Xeon E5 v2/Core i7 CBDMA
	0e2f  Xeon E7 v2/Xeon E5 v2/Core i7 CBDMA
	0e30  Xeon E7 v2/Xeon E5 v2/Core i7 Home Agent 0
	0e32  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Link 0
	0e33  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Link 1
	0e34  Xeon E7 v2/Xeon E5 v2/Core i7 R2PCIe
	0e36  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Ring Performance Ring Monitoring
	0e37  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Ring Performance Ring Monitoring
	0e38  Xeon E7 v2/Xeon E5 v2/Core i7 Home Agent 1
	0e3a  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Link 2
	0e3e  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Ring Performance Ring Monitoring
	0e3f  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Ring Performance Ring Monitoring
	0e40  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Link 2
	0e41  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Ring Registers
	0e43  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Link Reut 2
	0e44  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Link Reut 2
	0e45  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Link Agent Register
	0e47  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Link Agent Register
	0e60  Xeon E7 v2/Xeon E5 v2/Core i7 Home Agent 1
	0e68  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 1 Target Address/Thermal Registers
	0e6a  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 1 Channel Target Address Decoder Registers
	0e6b  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 1 Channel Target Address Decoder Registers
	0e6c  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 1 Channel Target Address Decoder Registers
	0e6d  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 1 Channel Target Address Decoder Registers
	0e71  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 0 RAS Registers
	0e74  Xeon E7 v2/Xeon E5 v2/Core i7 R2PCIe
	0e75  Xeon E7 v2/Xeon E5 v2/Core i7 R2PCIe
	0e77  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Ring Registers
	0e79  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 1 RAS Registers
	0e7d  Xeon E7 v2/Xeon E5 v2/Core i7 UBOX Registers
	0e7f  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Ring Registers
	0e80  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Link 0
	0e81  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Ring Registers
	0e83  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Link Reut 0
	0e84  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Link Reut 0
	0e85  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Link Agent Register
	0e87  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Registers
	0e90  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Link 1
	0e93  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Link 1
	0e94  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Link Reut 1
	0e95  Xeon E7 v2/Xeon E5 v2/Core i7 QPI Link Agent Register
	0ea0  Xeon E7 v2/Xeon E5 v2/Core i7 Home Agent 0
	0ea8  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 0 Target Address/Thermal Registers
	0eaa  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 0 Channel Target Address Decoder Registers
	0eab  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 0 Channel Target Address Decoder Registers
	0eac  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 0 Channel Target Address Decoder Registers
	0ead  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 0 Channel Target Address Decoder Registers
	0eae  Xeon E7 v2/Xeon E5 v2/Core i7 DDRIO Registers
	0eaf  Xeon E7 v2/Xeon E5 v2/Core i7 DDRIO Registers
	0eb0  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 1 Channel 0-3 Thermal Control 0
	0eb1  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 1 Channel 0-3 Thermal Control 1
	0eb2  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 1 Channel 0-3 ERROR Registers 0
	0eb3  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 1 Channel 0-3 ERROR Registers 1
	0eb4  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 1 Channel 0-3 Thermal Control 2
	0eb5  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 1 Channel 0-3 Thermal Control 3
	0eb6  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 1 Channel 0-3 ERROR Registers 2
	0eb7  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 1 Channel 0-3 ERROR Registers 3
	0ebc  Xeon E7 v2/Xeon E5 v2/Core i7 DDRIO Registers
	0ebe  Xeon E7 v2/Xeon E5 v2/Core i7 DDRIO Registers
	0ebf  Xeon E7 v2/Xeon E5 v2/Core i7 DDRIO Registers
	0ec0  Xeon E7 v2/Xeon E5 v2/Core i7 Power Control Unit 0
	0ec1  Xeon E7 v2/Xeon E5 v2/Core i7 Power Control Unit 1
	0ec2  Xeon E7 v2/Xeon E5 v2/Core i7 Power Control Unit 2
	0ec3  Xeon E7 v2/Xeon E5 v2/Core i7 Power Control Unit 3
	0ec4  Xeon E7 v2/Xeon E5 v2/Core i7 Power Control Unit 4
	0ec8  Xeon E7 v2/Xeon E5 v2/Core i7 System Address Decoder
	0ec9  Xeon E7 v2/Xeon E5 v2/Core i7 Broadcast Registers
	0eca  Xeon E7 v2/Xeon E5 v2/Core i7 Broadcast Registers
	0ed8  Xeon E7 v2/Xeon E5 v2/Core i7 DDRIO
	0ed9  Xeon E7 v2/Xeon E5 v2/Core i7 DDRIO
	0edc  Xeon E7 v2/Xeon E5 v2/Core i7 DDRIO
	0edd  Xeon E7 v2/Xeon E5 v2/Core i7 DDRIO
	0ede  Xeon E7 v2/Xeon E5 v2/Core i7 DDRIO
	0edf  Xeon E7 v2/Xeon E5 v2/Core i7 DDRIO
	0ee0  Xeon E7 v2/Xeon E5 v2/Core i7 Unicast Registers
	0ee1  Xeon E7 v2/Xeon E5 v2/Core i7 Unicast Registers
	0ee2  Xeon E7 v2/Xeon E5 v2/Core i7 Unicast Registers
	0ee3  Xeon E7 v2/Xeon E5 v2/Core i7 Unicast Registers
	0ee4  Xeon E7 v2/Xeon E5 v2/Core i7 Unicast Registers
	0ee5  Xeon E7 v2/Xeon E5 v2/Core i7 Unicast Registers
	0ee6  Xeon E7 v2/Xeon E5 v2/Core i7 Unicast Registers
	0ee7  Xeon E7 v2/Xeon E5 v2/Core i7 Unicast Registers
	0ee8  Xeon E7 v2/Xeon E5 v2/Core i7 Unicast Registers
	0ee9  Xeon E7 v2/Xeon E5 v2/Core i7 Unicast Registers
	0eea  Xeon E7 v2/Xeon E5 v2/Core i7 Unicast Registers
	0eeb  Xeon E7 v2/Xeon E5 v2/Core i7 Unicast Registers
	0eec  Xeon E7 v2/Xeon E5 v2/Core i7 Unicast Registers
	0eed  Xeon E7 v2/Xeon E5 v2/Core i7 Unicast Registers
	0eee  Xeon E7 v2/Xeon E5 v2/Core i7 Unicast Registers
	0ef0  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 0 Channel 0-3 Thermal Control 0
	0ef1  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 0 Channel 0-3 Thermal Control 1
	0ef2  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 0 Channel 0-3 ERROR Registers 0
	0ef3  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 0 Channel 0-3 ERROR Registers 1
	0ef4  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 0 Channel 0-3 Thermal Control 2
	0ef5  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 0 Channel 0-3 Thermal Control 3
	0ef6  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 0 Channel 0-3 ERROR Registers 2
	0ef7  Xeon E7 v2/Xeon E5 v2/Core i7 Integrated Memory Controller 0 Channel 0-3 ERROR Registers 3
	0ef8  Xeon E7 v2/Xeon E5 v2/Core i7 DDRIO
	0ef9  Xeon E7 v2/Xeon E5 v2/Core i7 DDRIO
	0efa  Xeon E7 v2/Xeon E5 v2/Core i7 DDRIO
	0efb  Xeon E7 v2/Xeon E5 v2/Core i7 DDRIO
	0efc  Xeon E7 v2/Xeon E5 v2/Core i7 DDRIO
	0efd  Xeon E7 v2/Xeon E5 v2/Core i7 DDRIO
	0f00  Atom Processor Z36xxx/Z37xxx Series SoC Transaction Register
	0f04  Atom Processor Z36xxx/Z37xxx Series High Definition Audio Controller
	0f06  Atom Processor Z36xxx/Z37xxx Series LPIO1 DMA Controller
	0f08  Atom Processor Z36xxx/Z37xxx Series LPIO1 PWM Controller
	0f09  Atom Processor Z36xxx/Z37xxx Series LPIO1 PWM Controller
	0f0a  Atom Processor Z36xxx/Z37xxx Series LPIO1 HSUART Controller #1
	0f0c  Atom Processor Z36xxx/Z37xxx Series LPIO1 HSUART Controller #2
	0f0e  Atom Processor Z36xxx/Z37xxx Series LPIO1 SPI Controller
	0f12  Atom Processor E3800 Series SMBus Controller
	0f14  Atom Processor Z36xxx/Z37xxx Series SDIO Controller
	0f15  Atom Processor Z36xxx/Z37xxx Series SDIO Controller
	0f16  Atom Processor Z36xxx/Z37xxx Series SDIO Controller
	0f18  Atom Processor Z36xxx/Z37xxx Series Trusted Execution Engine
	0f1c  Atom Processor Z36xxx/Z37xxx Series Power Control Unit
	0f20  Atom Processor E3800 Series SATA IDE Controller
	0f21  Atom Processor E3800 Series SATA IDE Controller
	0f22  Atom Processor E3800 Series SATA AHCI Controller
	0f23  Atom Processor E3800 Series SATA AHCI Controller
	0f28  Atom Processor Z36xxx/Z37xxx Series LPE Audio Controller
	0f31  Atom Processor Z36xxx/Z37xxx Series Graphics & Display
	0f34  Atom Processor Z36xxx/Z37xxx Series USB EHCI
	0f35  Atom Processor Z36xxx/Z37xxx, Celeron N2000 Series USB xHCI
	0f37  Atom Processor Z36xxx/Z37xxx Series OTG USB Device
	0f38  Atom Processor Z36xxx/Z37xxx Series Camera ISP
	0f40  Atom Processor Z36xxx/Z37xxx Series LPIO2 DMA Controller
	0f41  Atom Processor Z36xxx/Z37xxx Series LPIO2 I2C Controller #1
	0f42  Atom Processor Z36xxx/Z37xxx Series LPIO2 I2C Controller #2
	0f43  Atom Processor Z36xxx/Z37xxx Series LPIO2 I2C Controller #3
	0f44  Atom Processor Z36xxx/Z37xxx Series LPIO2 I2C Controller #4
	0f45  Atom Processor Z36xxx/Z37xxx Series LPIO2 I2C Controller #5
	0f46  Atom Processor Z36xxx/Z37xxx Series LPIO2 I2C Controller #6
	0f47  Atom Processor Z36xxx/Z37xxx Series LPIO2 I2C Controller #7
	0f48  Atom Processor E3800 Series PCI Express Root Port 1
	0f4a  Atom Processor E3800 Series PCI Express Root Port 2
	0f4c  Atom Processor E3800 Series PCI Express Root Port 3
	0f4e  Atom Processor E3800 Series PCI Express Root Port 4
	0f50  Atom Processor E3800 Series eMMC 4.5 Controller
	1000  82542 Gigabit Ethernet Controller (Fiber)
	1001  82543GC Gigabit Ethernet Controller (Fiber)
	1002  Pro 100 LAN+Modem 56 Cardbus II
	1004  82543GC Gigabit Ethernet Controller (Copper)
	1008  82544EI Gigabit Ethernet Controller (Copper)
	1009  82544EI Gigabit Ethernet Controller (Fiber)
	100a  82540EM Gigabit Ethernet Controller
	100c  82544GC Gigabit Ethernet Controller (Copper)
	100d  82544GC Gigabit Ethernet Controller (LOM)
	100e  82540EM Gigabit Ethernet Controller
	100f  82545EM Gigabit Ethernet Controller (Copper)
	1010  82546EB Gigabit Ethernet Controller (Copper)
	1011  82545EM Gigabit Ethernet Controller (Fiber)
	1012  82546EB Gigabit Ethernet Controller (Fiber)
	1013  82541EI Gigabit Ethernet Controller
	1014  82541ER Gigabit Ethernet Controller
	1015  82540EM Gigabit Ethernet Controller (LOM)
	1016  82540EP Gigabit Ethernet Controller (Mobile)
	1017  82540EP Gigabit Ethernet Controller
	1018  82541EI Gigabit Ethernet Controller
	1019  82547EI Gigabit Ethernet Controller
	101a  82547EI Gigabit Ethernet Controller (Mobile)
	101d  82546EB Gigabit Ethernet Controller
	101e  82540EP Gigabit Ethernet Controller (Mobile)
	1026  82545GM Gigabit Ethernet Controller
	1027  82545GM Gigabit Ethernet Controller
	1028  82545GM Gigabit Ethernet Controller
	1029  82559 Ethernet Controller
	1030  82559 InBusiness 10/100
	1031  82801CAM (ICH3) PRO/100 VE (LOM) Ethernet Controller
	1032  82801CAM (ICH3) PRO/100 VE Ethernet Controller
	1033  82801CAM (ICH3) PRO/100 VM (LOM) Ethernet Controller
	1034  82801CAM (ICH3) PRO/100 VM Ethernet Controller
	1035  82801CAM (ICH3)/82562EH (LOM) Ethernet Controller
	1036  82801CAM (ICH3) 82562EH Ethernet Controller
	1037  82801CAM (ICH3) Chipset Ethernet Controller
	1038  82801CAM (ICH3) PRO/100 VM (KM) Ethernet Controller
	1039  82801DB PRO/100 VE (LOM) Ethernet Controller
	103a  82801DB PRO/100 VE (CNR) Ethernet Controller
	103b  82801DB PRO/100 VM (LOM) Ethernet Controller
	103c  82801DB PRO/100 VM (CNR) Ethernet Controller
	103d  82801DB PRO/100 VE (MOB) Ethernet Controller
	103e  82801DB PRO/100 VM (MOB) Ethernet Controller
	1040  536EP Data Fax Modem
	1043  PRO/Wireless LAN 2100 3B Mini PCI Adapter
	1048  82597EX 10GbE Ethernet Controller
	1049  82566MM Gigabit Network Connection
	104a  82566DM Gigabit Network Connection
	104b  82566DC Gigabit Network Connection
	104c  82562V 10/100 Network Connection
	104d  82566MC Gigabit Network Connection
	1050  82562EZ 10/100 Ethernet Controller
	1051  82801EB/ER (ICH5/ICH5R) integrated LAN Controller
	1052  PRO/100 VM Network Connection
	1053  PRO/100 VM Network Connection
	1054  PRO/100 VE Network Connection
	1055  PRO/100 VM Network Connection
	1056  PRO/100 VE Network Connection
	1057  PRO/100 VE Network Connection
	1059  82551QM Ethernet Controller
	105b  82546GB Gigabit Ethernet Controller (Copper)
	105e  82571EB Gigabit Ethernet Controller
	105f  82571EB Gigabit Ethernet Controller
	1060  82571EB Gigabit Ethernet Controller
	1064  82562ET/EZ/GT/GZ - PRO/100 VE (LOM) Ethernet Controller
	1065  82562ET/EZ/GT/GZ - PRO/100 VE Ethernet Controller
	1066  82562 EM/EX/GX - PRO/100 VM (LOM) Ethernet Controller
	1067  82562 EM/EX/GX - PRO/100 VM Ethernet Controller
	1068  82562ET/EZ/GT/GZ - PRO/100 VE (LOM) Ethernet Controller Mobile
	1069  82562EM/EX/GX - PRO/100 VM (LOM) Ethernet Controller Mobile
	106a  82562G - PRO/100 VE (LOM) Ethernet Controller
	106b  82562G - PRO/100 VE Ethernet Controller Mobile
	1075  82547GI Gigabit Ethernet Controller
	1076  82541GI Gigabit Ethernet Controller
	1077  82541GI Gigabit Ethernet Controller
	1078  82541ER Gigabit Ethernet Controller
	1079  82546GB Gigabit Ethernet Controller
	107a  82546GB Gigabit Ethernet Controller
	107b  82546GB Gigabit Ethernet Controller
	107c  82541PI Gigabit Ethernet Controller
	107d  82572EI Gigabit Ethernet Controller (Copper)
	107e  82572EI Gigabit Ethernet Controller (Fiber)
	107f  82572EI Gigabit Ethernet Controller
	1080  FA82537EP 56K V.92 Data/Fax Modem PCI
	1081  631xESB/632xESB LAN Controller Copper
	1082  631xESB/632xESB LAN Controller fiber
	1083  631xESB/632xESB LAN Controller SERDES
	1084  631xESB/632xESB IDE Redirection
	1085  631xESB/632xESB Serial Port Redirection
	1086  631xESB/632xESB IPMI/KCS0
	1087  631xESB/632xESB UHCI Redirection
	1089  631xESB/632xESB BT
	108a  82546GB Gigabit Ethernet Controller
	108b  82573V Gigabit Ethernet Controller (Copper)
	108c  82573E Gigabit Ethernet Controller (Copper)
	108e  82573E KCS (Active Management)
	108f  Active Management Technology - SOL
	1091  PRO/100 VM Network Connection
	1092  PRO/100 VE Network Connection
	1093  PRO/100 VM Network Connection
	1094  PRO/100 VE Network Connection
	1095  PRO/100 VE Network Connection
	1096  80003ES2LAN Gigabit Ethernet Controller (Copper)
	1097  631xESB/632xESB DPT LAN Controller (Fiber)
	1098  80003ES2LAN Gigabit Ethernet Controller (Serdes)
	1099  82546GB Gigabit Ethernet Controller (Copper)
	109a  82573L Gigabit Ethernet Controller
	109b  82546GB PRO/1000 GF Quad Port Server Adapter
	109e  82597EX 10GbE Ethernet Controller
	10a0  82571EB PRO/1000 AT Quad Port Bypass Adapter
	10a1  82571EB PRO/1000 AF Quad Port Bypass Adapter
	10a4  82571EB Gigabit Ethernet Controller
	10a5  82571EB Gigabit Ethernet Controller (Fiber)
	10a6  82599EB 10-Gigabit Dummy Function
	10a7  82575EB Gigabit Network Connection
	10a9  82575EB Gigabit Backplane Connection
	10b0  82573L PRO/1000 PL Network Connection
	10b2  82573V PRO/1000 PM Network Connection
	10b3  82573E PRO/1000 PM Network Connection
	10b4  82573L PRO/1000 PL Network Connection
	10b5  82546GB Gigabit Ethernet Controller (Copper)
	10b6  82598 10GbE PCI-Express Ethernet Controller
	10b9  82572EI Gigabit Ethernet Controller (Copper)
	10ba  80003ES2LAN Gigabit Ethernet Controller (Copper)
	10bb  80003ES2LAN Gigabit Ethernet Controller (Serdes)
	10bc  82571EB Gigabit Ethernet Controller (Copper)
	10bd  82566DM-2 Gigabit Network Connection
	10bf  82567LF Gigabit Network Connection
	10c0  82562V-2 10/100 Network Connection
	10c2  82562G-2 10/100 Network Connection
	10c3  82562GT-2 10/100 Network Connection
	10c4  82562GT 10/100 Network Connection
	10c5  82562G 10/100 Network Connection
	10c6  82598EB 10-Gigabit AF Dual Port Network Connection
	10c7  82598EB 10-Gigabit AF Network Connection
	10c8  82598EB 10-Gigabit AT Network Connection
	10c9  82576 Gigabit Network Connection
	10ca  82576 Virtual Function
	10cb  82567V Gigabit Network Connection
	10cc  82567LM-2 Gigabit Network Connection
	10cd  82567LF-2 Gigabit Network Connection
	10ce  82567V-2 Gigabit Network Connection
	10d3  82574L Gigabit Network Connection
	10d4  Matrox Concord GE (customized Intel 82574)
	10d5  82571PT Gigabit PT Quad Port Server ExpressModule
	10d6  82575GB Gigabit Network Connection
	10d8  82599EB 10 Gigabit Unprogrammed
	10d9  82571EB Dual Port Gigabit Mezzanine Adapter
	10da  82571EB Quad Port Gigabit Mezzanine Adapter
	10db  82598EB 10-Gigabit Dual Port Network Connection
	10dd  82598EB 10-Gigabit AT CX4 Network Connection
	10de  82567LM-3 Gigabit Network Connection
	10df  82567LF-3 Gigabit Network Connection
	10e1  82598EB 10-Gigabit AF Dual Port Network Connection
	10e2  82575GB Gigabit Network Connection
	10e5  82567LM-4 Gigabit Network Connection
	10e6  82576 Gigabit Network Connection
	10e7  82576 Gigabit Network Connection
	10e8  82576 Gigabit Network Connection
	10ea  82577LM Gigabit Network Connection
	10eb  82577LC Gigabit Network Connection
	10ec  82598EB 10-Gigabit AT CX4 Network Connection
	10ed  82599 Ethernet Controller Virtual Function
	10ef  82578DM Gigabit Network Connection
	10f0  82578DC Gigabit Network Connection
	10f1  82598EB 10-Gigabit AF Dual Port Network Connection
	10f4  82598EB 10-Gigabit AF Network Connection
	10f5  82567LM Gigabit Network Connection
	10f6  82574L Gigabit Network Connection
	10f7  10 Gigabit BR KX4 Dual Port Network Connection
	10f8  82599 10 Gigabit Dual Port Backplane Connection
	10f9  82599 10 Gigabit Dual Port Network Connection
	10fb  82599ES 10-Gigabit SFI/SFP+ Network Connection
	10fc  82599 10 Gigabit Dual Port Network Connection
	10fe  82552 10/100 Network Connection
	1107  PRO/1000 MF Server Adapter (LX)
	1130  82815 815 Chipset Host Bridge and Memory Controller Hub
	1131  82815 815 Chipset AGP Bridge
	1132  82815 Chipset Graphics Controller (CGC)
	1161  82806AA PCI64 Hub Advanced Programmable Interrupt Controller
	1162  Xscale 80200 Big Endian Companion Chip
	1190  Merrifield SD/SDIO/eMMC Controller
	1191  Merrifield Serial IO HSUART Controller
	1192  Merrifield Serial IO HSUART DMA Controller
	1194  Merrifield Serial IO SPI Controller
	1195  Merrifield Serial IO I2C Controller
	1196  Merrifield Serial IO I2C Controller
	1199  Merrifield GPIO Controller
	119e  Merrifield USB Device Controller (OTG)
	11a0  Merrifield SCU IPC
	11a1  Merrifield Power Management Unit
	11a2  Merrifield Serial IO DMA Controller
	11a5  Merrifield Serial IO PWM Controller
	1200  IXP1200 Network Processor
	1209  8255xER/82551IT Fast Ethernet Controller
	1221  82092AA PCI to PCMCIA Bridge
	1222  82092AA IDE Controller
	1223  SAA7116
	1225  82452KX/GX [Orion]
	1226  82596 PRO/10 PCI
	1227  82865 EtherExpress PRO/100A
	1228  82556 EtherExpress PRO/100 Smart
	1229  82557/8/9/0/1 Ethernet Pro 100
	122d  430FX - 82437FX TSC [Triton I]
	122e  82371FB PIIX ISA [Triton I]
	1230  82371FB PIIX IDE [Triton I]
	1231  DSVD Modem
	1234  430MX - 82371MX Mobile PCI I/O IDE Xcelerator (MPIIX)
	1235  430MX - 82437MX Mob. System Ctrlr (MTSC) & 82438MX Data Path (MTDP)
	1237  440FX - 82441FX PMC [Natoma]
	1239  82371FB PIIX IDE Interface
	123b  82380PB PCI to PCI Docking Bridge
	123c  82380AB (MISA) Mobile PCI-to-ISA Bridge
	123d  683053 Programmable Interrupt Device
	123e  82466GX (IHPC) Integrated Hot-Plug Controller (hidden mode)
	123f  82466GX Integrated Hot-Plug Controller (IHPC)
	1240  82752 (752) AGP Graphics Accelerator
	124b  82380FB (MPCI2) Mobile Docking Controller
	1250  430HX - 82439HX TXC [Triton II]
	1360  82806AA PCI64 Hub PCI Bridge
	1361  82806AA PCI64 Hub Controller (HRes)
	1460  82870P2 P64H2 Hub PCI Bridge
	1461  82870P2 P64H2 I/OxAPIC
	1462  82870P2 P64H2 Hot Plug Controller
	1501  82567V-3 Gigabit Network Connection
	1502  82579LM Gigabit Network Connection (Lewisville)
	1503  82579V Gigabit Network Connection
	1507  Ethernet Express Module X520-P2
	1508  82598EB Gigabit BX Network Connection
	150a  82576NS Gigabit Network Connection
	150b  82598EB 10-Gigabit AT2 Server Adapter
	150c  82583V Gigabit Network Connection
	150d  82576 Gigabit Backplane Connection
	150e  82580 Gigabit Network Connection
	150f  82580 Gigabit Fiber Network Connection
	1510  82580 Gigabit Backplane Connection
	1511  82580 Gigabit SFP Connection
	1513  CV82524 Thunderbolt Controller [Light Ridge 4C 2010]
	1514  Ethernet X520 10GbE Dual Port KX4 Mezz
	1515  X540 Ethernet Controller Virtual Function
	1516  82580 Gigabit Network Connection
	1517  82599ES 10 Gigabit Network Connection
	1518  82576NS SerDes Gigabit Network Connection
	151a  DSL2310 Thunderbolt Controller [Eagle Ridge 2C 2011]
	151b  CVL2510 Thunderbolt Controller [Light Peak 2C 2010]
	151c  82599 10 Gigabit TN Network Connection
	1520  I350 Ethernet Controller Virtual Function
	1521  I350 Gigabit Network Connection
	1522  I350 Gigabit Fiber Network Connection
	1523  I350 Gigabit Backplane Connection
	1524  I350 Gigabit Connection
	1525  82567V-4 Gigabit Network Connection
	1526  82576 Gigabit Network Connection
	1527  82580 Gigabit Fiber Network Connection
	1528  Ethernet Controller 10-Gigabit X540-AT2
	1529  82599 10 Gigabit Dual Port Network Connection with FCoE
	152a  82599 10 Gigabit Dual Port Backplane Connection with FCoE
	152e  82599 Virtual Function
	152f  I350 Virtual Function
	1530  X540 Virtual Function
	1533  I210 Gigabit Network Connection
	1536  I210 Gigabit Fiber Network Connection
	1537  I210 Gigabit Backplane Connection
	1538  I210 Gigabit Network Connection
	1539  I211 Gigabit Network Connection
	153a  Ethernet Connection I217-LM
	153b  Ethernet Connection I217-V
	1547  DSL3510 Thunderbolt Controller [Cactus Ridge 4C 2012]
	1548  DSL3310 Thunderbolt Controller [Cactus Ridge 2C 2012]
	1549  DSL2210 Thunderbolt Controller [Port Ridge 1C 2011]
	154a  Ethernet Server Adapter X520-4
	154c  Ethernet Virtual Function 700 Series
	154d  Ethernet 10G 2P X520 Adapter
	1557  82599 10 Gigabit Network Connection
	1558  Ethernet Converged Network Adapter X520-Q1
	1559  Ethernet Connection I218-V
	155a  Ethernet Connection I218-LM
	155c  Ethernet Server Bypass Adapter
	155d  Ethernet Server Bypass Adapter
	1560  Ethernet Controller X540
	1563  Ethernet Controller 10G X550T
	1564  X550 Virtual Function
	1565  X550 Virtual Function
	1566  DSL4410 Thunderbolt NHI [Redwood Ridge 2C 2013]
	1567  DSL4410 Thunderbolt Bridge [Redwood Ridge 2C 2013]
	1568  DSL4510 Thunderbolt NHI [Redwood Ridge 4C 2013]
	1569  DSL4510 Thunderbolt Bridge [Redwood Ridge 4C 2013]
	156a  DSL5320 Thunderbolt 2 NHI [Falcon Ridge 2C 2013]
	156b  DSL5320 Thunderbolt 2 Bridge [Falcon Ridge 2C 2013]
	156c  DSL5520 Thunderbolt 2 NHI [Falcon Ridge 4C 2013]
	156d  DSL5520 Thunderbolt 2 Bridge [Falcon Ridge 4C 2013]
	156f  Ethernet Connection I219-LM
	1570  Ethernet Connection I219-V
	1571  Ethernet Virtual Function 700 Series
	1572  Ethernet Controller X710 for 10GbE SFP+
	1575  DSL6340 Thunderbolt 3 NHI [Alpine Ridge 2C 2015]
	1576  DSL6340 Thunderbolt 3 Bridge [Alpine Ridge 2C 2015]
	1577  DSL6540 Thunderbolt 3 NHI [Alpine Ridge 4C 2015]
	1578  DSL6540 Thunderbolt 3 Bridge [Alpine Ridge 4C 2015]
	157b  I210 Gigabit Network Connection
	157c  I210 Gigabit Backplane Connection
	157d  DSL5110 Thunderbolt 2 NHI (Low Power) [Win Ridge 2C 2014]
	157e  DSL5110 Thunderbolt 2 Bridge (Low Power) [Win Ridge 2C 2014]
	1580  Ethernet Controller XL710 for 40GbE backplane
	1581  Ethernet Controller X710 for 10GbE backplane
	1583  Ethernet Controller XL710 for 40GbE QSFP+
	1584  Ethernet Controller XL710 for 40GbE QSFP+
	1585  Ethernet Controller X710 for 10GbE QSFP+
	1586  Ethernet Controller X710 for 10GBASE-T
	1587  Ethernet Controller XL710 for 20GbE backplane
	1588  Ethernet Controller XL710 for 20GbE backplane
	1589  Ethernet Controller X710/X557-AT 10GBASE-T
	158a  Ethernet Controller XXV710 for 25GbE backplane
	158b  Ethernet Controller XXV710 for 25GbE SFP28
	15a0  Ethernet Connection (2) I218-LM
	15a1  Ethernet Connection (2) I218-V
	15a2  Ethernet Connection (3) I218-LM
	15a3  Ethernet Connection (3) I218-V
	15a4  Ethernet Switch FM10000 Host Interface
	15a5  Ethernet Switch FM10000 Host Virtual Interface
	15a8  Ethernet Connection X552 Virtual Function
	15a9  X552 Virtual Function
	15aa  Ethernet Connection X552 10 GbE Backplane
	15ab  Ethernet Connection X552 10 GbE Backplane
	15ac  Ethernet Connection X552 10 GbE SFP+
	15ad  Ethernet Connection X552/X557-AT 10GBASE-T
	15ae  Ethernet Connection X552 1000BASE-T
	15b0  Ethernet Connection X552 Backplane
	15b4  X553 Virtual Function
	15b5  DSL6340 USB 3.1 Controller [Alpine Ridge]
	15b6  DSL6540 USB 3.1 Controller [Alpine Ridge]
	15b7  Ethernet Connection (2) I219-LM
	15b8  Ethernet Connection (2) I219-V
	15b9  Ethernet Connection (3) I219-LM
	15bb  Ethernet Connection (7) I219-LM
	15bc  Ethernet Connection (7) I219-V
	15bd  Ethernet Connection (6) I219-LM
	15be  Ethernet Connection (6) I219-V
	15bf  JHL6240 Thunderbolt 3 NHI (Low Power) [Alpine Ridge LP 2016]
	15c0  JHL6240 Thunderbolt 3 Bridge (Low Power) [Alpine Ridge LP 2016]
	15c2  Ethernet Connection X553 Backplane
	15c3  Ethernet Connection X553 Backplane
	15c4  Ethernet Connection X553 10 GbE SFP+
	15c5  X553 Virtual Function
	15c6  Ethernet Connection X553 1GbE
	15c7  Ethernet Connection X553 1GbE
	15c8  Ethernet Connection X553/X557-AT 10GBASE-T
	15ce  Ethernet Connection X553 10 GbE SFP+
	15d0  Ethernet SDI Adapter FM10420-100GbE-QDA2
	15d1  Ethernet Controller 10G X550T
	15d2  JHL6540 Thunderbolt 3 NHI (C step) [Alpine Ridge 4C 2016]
	15d3  JHL6540 Thunderbolt 3 Bridge (C step) [Alpine Ridge 4C 2016]
	15d4  JHL6540 Thunderbolt 3 USB Controller (C step) [Alpine Ridge 4C 2016]
	15d5  Ethernet SDI Adapter FM10420-25GbE-DA2
	15d6  Ethernet Connection (5) I219-V
	15d7  Ethernet Connection (4) I219-LM
	15d8  Ethernet Connection (4) I219-V
	15d9  JHL6340 Thunderbolt 3 NHI (C step) [Alpine Ridge 2C 2016]
	15da  JHL6340 Thunderbolt 3 Bridge (C step) [Alpine Ridge 2C 2016]
	15df  Ethernet Connection (8) I219-LM
	15e0  Ethernet Connection (8) I219-V
	15e1  Ethernet Connection (9) I219-LM
	15e2  Ethernet Connection (9) I219-V
	15e3  Ethernet Connection (5) I219-LM
	15e4  Ethernet Connection X553 1GbE
	15e5  Ethernet Connection X553 1GbE
	1600  Broadwell-U Host Bridge -OPI
	1601  Broadwell-U PCI Express x16 Controller
	1602  Broadwell-U Integrated Graphics
	1603  Broadwell-U Processor Thermal Subsystem
	1604  Broadwell-U Host Bridge -OPI
	1605  Broadwell-U PCI Express x8 Controller
	1606  HD Graphics
	1607  Broadwell-U CHAPS Device
	1608  Broadwell-U Host Bridge -OPI
	1609  Broadwell-U x4 PCIe
	160a  Broadwell-U Integrated Graphics
	160b  Broadwell-U Integrated Graphics
	160c  Broadwell-U Audio Controller
	160d  Broadwell-U Integrated Graphics
	160e  Broadwell-U Integrated Graphics
	160f  Broadwell-U SoftSKU
	1610  Broadwell-U Host Bridge - DMI
	1612  HD Graphics 5600
	1614  Broadwell-U Host Bridge - DMI
	1616  HD Graphics 5500
	1618  Broadwell-U Host Bridge - DMI
	161a  Broadwell-U Integrated Graphics
	161b  Broadwell-U Integrated Graphics
	161d  Broadwell-U Integrated Graphics
	161e  HD Graphics 5300
	1622  Iris Pro Graphics 6200
	1626  HD Graphics 6000
	162a  Iris Pro Graphics P6300
	162b  Iris Graphics 6100
	162d  Broadwell-U Integrated Graphics
	162e  Broadwell-U Integrated Graphics
	1632  Broadwell-U Integrated Graphics
	1636  Broadwell-U Integrated Graphics
	163a  Broadwell-U Integrated Graphics
	163b  Broadwell-U Integrated Graphics
	163d  Broadwell-U Integrated Graphics
	163e  Broadwell-U Integrated Graphics
	1889  Ethernet Adaptive Virtual Function
	1900  Xeon E3-1200 v5/E3-1500 v5/6th Gen Core Processor Host Bridge/DRAM Registers
	1901  Xeon E3-1200 v5/E3-1500 v5/6th Gen Core Processor PCIe Controller (x16)
	1902  HD Graphics 510
	1903  Xeon E3-1200 v5/E3-1500 v5/6th Gen Core Processor Thermal Subsystem
	1904  Xeon E3-1200 v5/E3-1500 v5/6th Gen Core Processor Host Bridge/DRAM Registers
	1905  Xeon E3-1200 v5/E3-1500 v5/6th Gen Core Processor PCIe Controller (x8)
	1906  HD Graphics 510
	1908  Xeon E3-1200 v5/E3-1500 v5/6th Gen Core Processor Host Bridge/DRAM Registers
	1909  Xeon E3-1200 v5/E3-1500 v5/6th Gen Core Processor PCIe Controller (x4)
	190c  Xeon E3-1200 v5/E3-1500 v5/6th Gen Core Processor Host Bridge/DRAM Registers
	190f  Xeon E3-1200 v5/E3-1500 v5/6th Gen Core Processor Host Bridge/DRAM Registers
	1910  Xeon E3-1200 v5/E3-1500 v5/6th Gen Core Processor Host Bridge/DRAM Registers
	1911  Xeon E3-1200 v5/v6 / E3-1500 v5 / 6th/7th Gen Core Processor Gaussian Mixture Model
	1912  HD Graphics 530
	1916  Skylake GT2 [HD Graphics 520]
	1918  Xeon E3-1200 v5/E3-1500 v5/6th Gen Core Processor Host Bridge/DRAM Registers
	1919  Xeon E3-1200 v5/E3-1500 v5/6th Gen Core Processor Imaging Unit
	191b  HD Graphics 530
	191d  HD Graphics P530
	191e  HD Graphics 515
	191f  Xeon E3-1200 v5/E3-1500 v5/6th Gen Core Processor Host Bridge/DRAM Registers
	1921  HD Graphics 520
	1926  Iris Graphics 540
	1927  Iris Graphics 550
	192b  Iris Graphics 555
	192d  Iris Graphics P555
	1932  Iris Pro Graphics 580
	193a  Iris Pro Graphics P580
	193b  Iris Pro Graphics 580
	193d  Iris Pro Graphics P580
	1960  80960RP (i960RP) Microprocessor
	1962  80960RM (i960RM) Microprocessor
	19ac  DNV SMBus Contoller - Host
	19b0  DNV SATA Controller 0
	19b1  DNV SATA Controller 0
	19b2  DNV SATA Controller 0
	19b3  DNV SATA Controller 0
	19b4  DNV SATA Controller 0
	19b5  DNV SATA Controller 0
	19b6  DNV SATA Controller 0
	19b7  DNV SATA Controller 0
	19be  DNV SATA Controller 0
	19bf  DNV SATA Controller 0
	19c0  DNV SATA Controller 1
	19c1  DNV SATA Controller 1
	19c2  DNV SATA Controller 1
	19c3  DNV SATA Controller 1
	19c4  DNV SATA Controller 1
	19c5  DNV SATA Controller 1
	19c6  DNV SATA Controller 1
	19c7  DNV SATA Controller 1
	19ce  DNV SATA Controller 1
	19cf  DNV SATA Controller 1
	19dc  DNV LPC or eSPI
	19df  DNV SMBus controller
	19e0  DNV SPI Controller
	1a21  82840 840 [Carmel] Chipset Host Bridge (Hub A)
	1a23  82840 840 [Carmel] Chipset AGP Bridge
	1a24  82840 840 [Carmel] Chipset PCI Bridge (Hub B)
	1a30  82845 845 [Brookdale] Chipset Host Bridge
	1a31  82845 845 [Brookdale] Chipset AGP Bridge
	1a38  5000 Series Chipset DMA Engine
	1a48  82597EX 10GbE Ethernet Controller
	1b48  82597EX 10GbE Ethernet Controller
	1c00  6 Series/C200 Series Chipset Family Desktop SATA Controller (IDE mode, ports 0-3)
	1c01  6 Series/C200 Series Chipset Family Mobile SATA Controller (IDE mode, ports 0-3)
	1c02  6 Series/C200 Series Chipset Family 6 port Desktop SATA AHCI Controller
	1c03  6 Series/C200 Series Chipset Family 6 port Mobile SATA AHCI Controller
	1c04  6 Series/C200 Series Desktop SATA RAID Controller
	1c05  6 Series/C200 Series Mobile SATA RAID Controller
	1c06  Z68 Express Chipset SATA RAID Controller
	1c08  6 Series/C200 Series Chipset Family Desktop SATA Controller (IDE mode, ports 4-5)
	1c09  6 Series/C200 Series Chipset Family Mobile SATA Controller (IDE mode, ports 4-5)
	1c10  6 Series/C200 Series Chipset Family PCI Express Root Port 1
	1c12  6 Series/C200 Series Chipset Family PCI Express Root Port 2
	1c14  6 Series/C200 Series Chipset Family PCI Express Root Port 3
	1c16  6 Series/C200 Series Chipset Family PCI Express Root Port 4
	1c18  6 Series/C200 Series Chipset Family PCI Express Root Port 5
	1c1a  6 Series/C200 Series Chipset Family PCI Express Root Port 6
	1c1c  6 Series/C200 Series Chipset Family PCI Express Root Port 7
	1c1e  6 Series/C200 Series Chipset Family PCI Express Root Port 8
	1c20  6 Series/C200 Series Chipset Family High Definition Audio Controller
	1c22  6 Series/C200 Series Chipset Family SMBus Controller
	1c24  6 Series/C200 Series Chipset Family Thermal Management Controller
	1c25  6 Series/C200 Series Chipset Family DMI to PCI Bridge
	1c26  6 Series/C200 Series Chipset Family USB Enhanced Host Controller #1
	1c27  6 Series/C200 Series Chipset Family USB Universal Host Controller #1
	1c2c  6 Series/C200 Series Chipset Family USB Universal Host Controller #5
	1c2d  6 Series/C200 Series Chipset Family USB Enhanced Host Controller #2
	1c33  6 Series/C200 Series Chipset Family LAN Controller
	1c35  6 Series/C200 Series Chipset Family VECI Controller
	1c3a  6 Series/C200 Series Chipset Family MEI Controller #1
	1c3b  6 Series/C200 Series Chipset Family MEI Controller #2
	1c3c  6 Series/C200 Series Chipset Family IDE-r Controller
	1c3d  6 Series/C200 Series Chipset Family KT Controller
	1c40  6 Series/C200 Series Chipset Family LPC Controller
	1c41  Mobile SFF 6 Series Chipset Family LPC Controller
	1c42  6 Series/C200 Series Chipset Family LPC Controller
	1c43  Mobile 6 Series Chipset Family LPC Controller
	1c44  Z68 Express Chipset Family LPC Controller
	1c45  6 Series/C200 Series Chipset Family LPC Controller
	1c46  P67 Express Chipset Family LPC Controller
	1c47  UM67 Express Chipset Family LPC Controller
	1c48  6 Series/C200 Series Chipset Family LPC Controller
	1c49  HM65 Express Chipset Family LPC Controller
	1c4a  H67 Express Chipset Family LPC Controller
	1c4b  HM67 Express Chipset Family LPC Controller
	1c4c  Q65 Express Chipset Family LPC Controller
	1c4d  QS67 Express Chipset Family LPC Controller
	1c4e  Q67 Express Chipset Family LPC Controller
	1c4f  QM67 Express Chipset Family LPC Controller
	1c50  B65 Express Chipset Family LPC Controller
	1c51  6 Series/C200 Series Chipset Family LPC Controller
	1c52  C202 Chipset Family LPC Controller
	1c53  6 Series/C200 Series Chipset Family LPC Controller
	1c54  C204 Chipset Family LPC Controller
	1c55  6 Series/C200 Series Chipset Family LPC Controller
	1c56  C206 Chipset Family LPC Controller
	1c57  6 Series/C200 Series Chipset Family LPC Controller
	1c58  Upgraded B65 Express Chipset Family LPC Controller
	1c59  Upgraded HM67 Express Chipset Family LPC Controller
	1c5a  Upgraded Q67 Express Chipset Family LPC Controller
	1c5b  6 Series/C200 Series Chipset Family LPC Controller
	1c5c  H61 Express Chipset Family LPC Controller
	1c5d  6 Series/C200 Series Chipset Family LPC Controller
	1c5e  6 Series/C200 Series Chipset Family LPC Controller
	1c5f  6 Series/C200 Series Chipset Family LPC Controller
	1d00  C600/X79 series chipset 4-Port SATA IDE Controller
	1d02  C600/X79 series chipset 6-Port SATA AHCI Controller
	1d04  C600/X79 series chipset SATA RAID Controller
	1d06  C600/X79 series chipset SATA Premium RAID Controller
	1d08  C600/X79 series chipset 2-Port SATA IDE Controller
	1d10  C600/X79 series chipset PCI Express Root Port 1
	1d11  C600/X79 series chipset PCI Express Root Port 1
	1d12  C600/X79 series chipset PCI Express Root Port 2
	1d13  C600/X79 series chipset PCI Express Root Port 2
	1d14  C600/X79 series chipset PCI Express Root Port 3
	1d15  C600/X79 series chipset PCI Express Root Port 3
	1d16  C600/X79 series chipset PCI Express Root Port 4
	1d17  C600/X79 series chipset PCI Express Root Port 4
	1d18  C600/X79 series chipset PCI Express Root Port 5
	1d19  C600/X79 series chipset PCI Express Root Port 5
	1d1a  C600/X79 series chipset PCI Express Root Port 6
	1d1b  C600/X79 series chipset PCI Express Root Port 6
	1d1c  C600/X79 series chipset PCI Express Root Port 7
	1d1d  C600/X79 series chipset PCI Express Root Port 7
	1d1e  C600/X79 series chipset PCI Express Root Port 8
	1d1f  C600/X79 series chipset PCI Express Root Port 8
	1d20  C600/X79 series chipset High Definition Audio Controller
	1d22  C600/X79 series chipset SMBus Host Controller
	1d24  C600/X79 series chipset Thermal Management Controller
	1d25  C600/X79 series chipset DMI to PCI Bridge
	1d26  C600/X79 series chipset USB2 Enhanced Host Controller #1
	1d2d  C600/X79 series chipset USB2 Enhanced Host Controller #2
	1d33  C600/X79 series chipset LAN Controller
	1d35  C600/X79 series chipset VECI Controller
	1d3a  C600/X79 series chipset MEI Controller #1
	1d3b  C600/X79 series chipset MEI Controller #2
	1d3c  C600/X79 series chipset IDE-r Controller
	1d3d  C600/X79 series chipset KT Controller
	1d3e  C600/X79 series chipset PCI Express Virtual Root Port
	1d3f  C608/C606/X79 series chipset PCI Express Virtual Switch Port
	1d40  C600/X79 series chipset LPC Controller
	1d41  C600/X79 series chipset LPC Controller
	1d50  C608 chipset Dual 4-Port SATA/SAS Storage Control Unit
	1d54  C600/X79 series chipset Dual 4-Port SATA/SAS Storage Control Unit
	1d55  C600/X79 series chipset 4-Port SATA/SAS Storage Control Unit
	1d58  C606 chipset Dual 4-Port SATA/SAS Storage Control Unit
	1d59  C604/X79 series chipset 4-Port SATA/SAS Storage Control Unit
	1d5a  C600/X79 series chipset Dual 4-Port SATA Storage Control Unit
	1d5b  C602 chipset 4-Port SATA Storage Control Unit
	1d5c  C600/X79 series chipset Dual 4-Port SATA/SAS Storage Control Unit
	1d5d  C600/X79 series chipset 4-Port SATA/SAS Storage Control Unit
	1d5e  C600/X79 series chipset Dual 4-Port SATA Storage Control Unit
	1d5f  C600/X79 series chipset 4-Port SATA Storage Control Unit
	1d60  C608 chipset Dual 4-Port SATA/SAS Storage Control Unit
	1d64  C600/X79 series chipset Dual 4-Port SATA/SAS Storage Control Unit
	1d65  C600/X79 series chipset 4-Port SATA/SAS Storage Control Unit
	1d68  C606 chipset Dual 4-Port SATA/SAS Storage Control Unit
	1d69  C604/X79 series chipset 4-Port SATA/SAS Storage Control Unit
	1d6a  C600/X79 series chipset Dual 4-Port SATA Storage Control Unit
	1d6b  C602 chipset 4-Port SATA Storage Control Unit
	1d6c  C600/X79 series chipset Dual 4-Port SATA/SAS Storage Control Unit
	1d6d  C600/X79 series chipset 4-Port SATA/SAS Storage Control Unit
	1d6e  C600/X79 series chipset Dual 4-Port SATA Storage Control Unit
	1d6f  C600/X79 series chipset 4-Port SATA Storage Control Unit
	1d70  C600/X79 series chipset SMBus Controller 0
	1d71  C608/C606/X79 series chipset SMBus Controller 1
	1d72  C608 chipset SMBus Controller 2
	1d74  C608/C606/X79 series chipset PCI Express Upstream Port
	1d76  C600/X79 series chipset Multi-Function Glue
	1e00  7 Series/C210 Series Chipset Family 4-port SATA Controller [IDE mode]
	1e01  7 Series Chipset Family 4-port SATA Controller [IDE mode]
	1e02  7 Series/C210 Series Chipset Family 6-port SATA Controller [AHCI mode]
	1e03  7 Series Chipset Family 6-port SATA Controller [AHCI mode]
	1e04  7 Series/C210 Series Chipset Family SATA Controller [RAID mode]
	1e05  7 Series Chipset SATA Controller [RAID mode]
	1e06  7 Series/C210 Series Chipset Family SATA Controller [RAID mode]
	1e07  7 Series Chipset Family SATA Controller [RAID mode]
	1e08  7 Series/C210 Series Chipset Family 2-port SATA Controller [IDE mode]
	1e09  7 Series Chipset Family 2-port SATA Controller [IDE mode]
	1e0e  7 Series/C210 Series Chipset Family SATA Controller [RAID mode]
	1e10  7 Series/C216 Chipset Family PCI Express Root Port 1
	1e12  7 Series/C210 Series Chipset Family PCI Express Root Port 2
	1e14  7 Series/C210 Series Chipset Family PCI Express Root Port 3
	1e16  7 Series/C216 Chipset Family PCI Express Root Port 4
	1e18  7 Series/C210 Series Chipset Family PCI Express Root Port 5
	1e1a  7 Series/C210 Series Chipset Family PCI Express Root Port 6
	1e1c  7 Series/C210 Series Chipset Family PCI Express Root Port 7
	1e1e  7 Series/C210 Series Chipset Family PCI Express Root Port 8
	1e20  7 Series/C216 Chipset Family High Definition Audio Controller
	1e22  7 Series/C216 Chipset Family SMBus Controller
	1e24  7 Series/C210 Series Chipset Family Thermal Management Controller
	1e25  7 Series/C210 Series Chipset Family DMI to PCI Bridge
	1e26  7 Series/C216 Chipset Family USB Enhanced Host Controller #1
	1e2d  7 Series/C216 Chipset Family USB Enhanced Host Controller #2
	1e31  7 Series/C210 Series Chipset Family USB xHCI Host Controller
	1e33  7 Series/C210 Series Chipset Family LAN Controller
	1e3a  7 Series/C216 Chipset Family MEI Controller #1
	1e3b  7 Series/C210 Series Chipset Family MEI Controller #2
	1e3c  7 Series/C210 Series Chipset Family IDE-r Controller
	1e3d  7 Series/C210 Series Chipset Family KT Controller
	1e41  7 Series Chipset Family LPC Controller
	1e42  7 Series Chipset Family LPC Controller
	1e43  7 Series Chipset Family LPC Controller
	1e44  Z77 Express Chipset LPC Controller
	1e45  7 Series Chipset Family LPC Controller
	1e46  Z75 Express Chipset LPC Controller
	1e47  Q77 Express Chipset LPC Controller
	1e48  Q75 Express Chipset LPC Controller
	1e49  B75 Express Chipset LPC Controller
	1e4a  H77 Express Chipset LPC Controller
	1e4b  7 Series Chipset Family LPC Controller
	1e4c  7 Series Chipset Family LPC Controller
	1e4d  7 Series Chipset Family LPC Controller
	1e4e  7 Series Chipset Family LPC Controller
	1e4f  7 Series Chipset Family LPC Controller
	1e50  7 Series Chipset Family LPC Controller
	1e51  7 Series Chipset Family LPC Controller
	1e52  7 Series Chipset Family LPC Controller
	1e53  C216 Series Chipset LPC Controller
	1e54  7 Series Chipset Family LPC Controller
	1e55  QM77 Express Chipset LPC Controller
	1e56  QS77 Express Chipset LPC Controller
	1e57  HM77 Express Chipset LPC Controller
	1e58  UM77 Express Chipset LPC Controller
	1e59  HM76 Express Chipset LPC Controller
	1e5a  7 Series Chipset Family LPC Controller
	1e5b  UM77 Express Chipset LPC Controller
	1e5c  7 Series Chipset Family LPC Controller
	1e5d  HM75 Express Chipset LPC Controller
	1e5e  7 Series Chipset Family LPC Controller
	1e5f  7 Series Chipset Family LPC Controller
	1f00  Atom processor C2000 SoC Transaction Router
	1f01  Atom processor C2000 SoC Transaction Router
	1f02  Atom processor C2000 SoC Transaction Router
	1f03  Atom processor C2000 SoC Transaction Router
	1f04  Atom processor C2000 SoC Transaction Router
	1f05  Atom processor C2000 SoC Transaction Router
	1f06  Atom processor C2000 SoC Transaction Router
	1f07  Atom processor C2000 SoC Transaction Router
	1f08  Atom processor C2000 SoC Transaction Router
	1f09  Atom processor C2000 SoC Transaction Router
	1f0a  Atom processor C2000 SoC Transaction Router
	1f0b  Atom processor C2000 SoC Transaction Router
	1f0c  Atom processor C2000 SoC Transaction Router
	1f0d  Atom processor C2000 SoC Transaction Router
	1f0e  Atom processor C2000 SoC Transaction Router
	1f0f  Atom processor C2000 SoC Transaction Router
	1f10  Atom processor C2000 PCIe Root Port 1
	1f11  Atom processor C2000 PCIe Root Port 2
	1f12  Atom processor C2000 PCIe Root Port 3
	1f13  Atom processor C2000 PCIe Root Port 4
	1f14  Atom processor C2000 RAS
	1f15  Atom processor C2000 SMBus 2.0
	1f16  Atom processor C2000 RCEC
	1f18  Atom processor C2000 QAT
	1f19  Atom processor C2000 QAT
	1f20  Atom processor C2000 4-Port IDE SATA2 Controller
	1f21  Atom processor C2000 4-Port IDE SATA2 Controller
	1f22  Atom processor C2000 AHCI SATA2 Controller
	1f23  Atom processor C2000 AHCI SATA2 Controller
	1f24  Atom processor C2000 RAID SATA2 Controller
	1f25  Atom processor C2000 RAID SATA2 Controller
	1f26  Atom processor C2000 RAID SATA2 Controller
	1f27  Atom processor C2000 RAID SATA2 Controller
	1f2c  Atom processor C2000 USB Enhanced Host Controller
	1f2e  Atom processor C2000 RAID SATA2 Controller
	1f2f  Atom processor C2000 RAID SATA2 Controller
	1f30  Atom processor C2000 2-Port IDE SATA3 Controller
	1f31  Atom processor C2000 2-Port IDE SATA3 Controller
	1f32  Atom processor C2000 AHCI SATA3 Controller
	1f33  Atom processor C2000 AHCI SATA3 Controller
	1f34  Atom processor C2000 RAID SATA3 Controller
	1f35  Atom processor C2000 RAID SATA3 Controller
	1f36  Atom processor C2000 RAID SATA3 Controller
	1f37  Atom processor C2000 RAID SATA3 Controller
	1f38  Atom processor C2000 PCU
	1f39  Atom processor C2000 PCU
	1f3a  Atom processor C2000 PCU
	1f3b  Atom processor C2000 PCU
	1f3c  Atom processor C2000 PCU SMBus
	1f3e  Atom processor C2000 RAID SATA3 Controller
	1f3f  Atom processor C2000 RAID SATA3 Controller
	1f40  Ethernet Connection I354 1.0 GbE Backplane
	1f41  Ethernet Connection I354
	1f42  Atom processor C2000 GbE
	1f44  Atom processor C2000 GbE Virtual Function
	1f45  Ethernet Connection I354 2.5 GbE Backplane
	2014  Sky Lake-E Ubox Registers
	2015  Sky Lake-E Ubox Registers
	2016  Sky Lake-E Ubox Registers
	2018  Sky Lake-E M2PCI Registers
	201a  Sky Lake-E Non-Transparent Bridge Registers
	201c  Sky Lake-E Non-Transparent Bridge Registers
	2020  Sky Lake-E DMI3 Registers
	2021  Sky Lake-E CBDMA Registers
	2024  Sky Lake-E MM/Vt-d Configuration Registers
	2030  Sky Lake-E PCI Express Root Port A
	2031  Sky Lake-E PCI Express Root Port B
	2032  Sky Lake-E PCI Express Root Port C
	2033  Sky Lake-E PCI Express Root Port D
	2035  Sky Lake-E RAS Configuration Registers
	204c  Sky Lake-E M3KTI Registers
	204d  Sky Lake-E M3KTI Registers
	204e  Sky Lake-E M3KTI Registers
	2054  Sky Lake-E CHA Registers
	2055  Sky Lake-E CHA Registers
	2056  Sky Lake-E CHA Registers
	2057  Sky Lake-E CHA Registers
	2068  Sky Lake-E DDRIO Registers
	2069  Sky Lake-E DDRIO Registers
	206a  Sky Lake-E IOxAPIC Configuration Registers
	206e  Sky Lake-E DDRIO Registers
	206f  Sky Lake-E DDRIO Registers
	2078  Sky Lake-E PCU Registers
	207a  Sky Lake-E PCU Registers
	2080  Sky Lake-E PCU Registers
	2081  Sky Lake-E PCU Registers
	2082  Sky Lake-E PCU Registers
	2083  Sky Lake-E PCU Registers
	2084  Sky Lake-E PCU Registers
	2085  Sky Lake-E PCU Registers
	2086  Sky Lake-E PCU Registers
	208d  Sky Lake-E CHA Registers
	208e  Sky Lake-E CHA Registers
	2250  Xeon Phi coprocessor 5100 series
	225c  Xeon Phi coprocessor SE10/7120 series
	225d  Xeon Phi coprocessor 3120 series
	225e  Xeon Phi coprocessor 31S1
	2280  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series SoC Transaction Register
	2284  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series High Definition Audio Controller
	2286  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series LPIO1 DMA Controller
	228a  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series LPIO1 HSUART Controller #1
	228c  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series LPIO1 HSUART Controller #2
	2292  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx SMBus Controller
	2294  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series MMC Controller
	2295  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series SDIO Controller
	2296  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series SD Controller
	2298  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series Trusted Execution Engine
	229c  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series PCU
	22a3  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series SATA Controller
	22a4  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series SATA AHCI Controller
	22a8  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series Low Power Engine Audio
	22b0  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series PCI Configuration Registers
	22b1  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Integrated Graphics Controller
	22b5  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series USB xHCI Controller
	22b8  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series Imaging Unit
	22c0  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series LPIO2 DMA Controller
	22c1  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series LPIO2 I2C Controller #1
	22c2  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series LPIO2 I2C Controller #2
	22c3  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series LPIO2 I2C Controller #3
	22c4  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series LPIO2 I2C Controller #4
	22c5  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series LPIO2 I2C Controller #5
	22c6  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series LPIO2 I2C Controller #6
	22c7  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series LPIO2 I2C Controller #7
	22c8  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series PCI Express Port #1
	22ca  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series PCI Express Port #2
	22cc  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series PCI Express Port #3
	22ce  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series PCI Express Port #4
	22dc  Atom/Celeron/Pentium Processor x5-E8000/J3xxx/N3xxx Series Power Management Controller
	2310  DH89xxCC LPC Controller
	2323  DH89xxCC 4 Port SATA AHCI Controller
	2330  DH89xxCC SMBus Controller
	2331  DH89xxCC Chap Counter
	2332  DH89xxCC Thermal Subsystem
	2334  DH89xxCC USB2 Enhanced Host Controller #1
	2335  DH89xxCC USB2 Enhanced Host Controller #1
	2342  DH89xxCC PCI Express Root Port #1
	2343  DH89xxCC PCI Express Root Port #1
	2344  DH89xxCC PCI Express Root Port #2
	2345  DH89xxCC PCI Express Root Port #2
	2346  DH89xxCC PCI Express Root Port #3
	2347  DH89xxCC PCI Express Root Port #3
	2348  DH89xxCC PCI Express Root Port #4
	2349  DH89xxCC PCI Express Root Port #4
	2360  DH89xxCC Watchdog Timer
	2364  DH89xxCC MEI 0
	2365  DH89xxCC MEI 1
	2390  DH895XCC Series LPC Controller
	23a1  DH895XCC Series 2-Port SATA Controller [IDE Mode]
	23a3  DH895XCC Series 4-Port SATA Controller [AHCI Mode]
	23a6  DH895XCC Series 2-Port SATA Controller [IDE Mode]
	23b0  DH895XCC Series SMBus Controller
	23b1  DH895XCC Series CHAP Counter
	23b2  DH895XCC Series Thermal Management Controller
	23b4  DH895XCC Series USB2 Enhanced Host Controller #1
	23b5  DH895XCC Series USB2 Enhanced Host Controller #1
	23c2  DH895XCC Series PCI Express Root Port #1
	23c3  DH895XCC Series PCI Express Root Port #1
	23c4  DH895XCC Series PCI Express Root Port #2
	23c5  DH895XCC Series PCI Express Root Port #2
	23c6  CDH895XCC Series PCI Express Root Port #3
	23c7  DH895XCC Series PCI Express Root Port #3
	23c8  DH895XCC Series PCI Express Root Port #4
	23c9  DH895XCC Series PCI Express Root Port #4
	23e0  DH895XCC Series Watchdog Timer
	23e4  DH895XCC Series MEI Controller #1
	23e5  DH895XCC Series MEI Controller #2
	2410  82801AA ISA Bridge (LPC)
	2411  82801AA IDE Controller
	2412  82801AA USB Controller
	2413  82801AA SMBus Controller
	2415  82801AA AC'97 Audio Controller
	2416  82801AA AC'97 Modem Controller
	2418  82801AA PCI Bridge
	2420  82801AB ISA Bridge (LPC)
	2421  82801AB IDE Controller
	2422  82801AB USB Controller
	2423  82801AB SMBus Controller
	2425  82801AB AC'97 Audio Controller
	2426  82801AB AC'97 Modem Controller
	2428  82801AB PCI Bridge
	2440  82801BA ISA Bridge (LPC)
	2442  82801BA/BAM UHCI USB 1.1 Controller #1
	2443  82801BA/BAM SMBus Controller
	2444  82801BA/BAM UHCI USB 1.1 Controller #2
	2445  82801BA/BAM AC'97 Audio Controller
	2446  82801BA/BAM AC'97 Modem Controller
	2448  82801 Mobile PCI Bridge
	2449  82801BA/BAM/CA/CAM Ethernet Controller
	244a  82801BAM IDE U100 Controller
	244b  82801BA IDE U100 Controller
	244c  82801BAM ISA Bridge (LPC)
	244e  82801 PCI Bridge
	2450  82801E ISA Bridge (LPC)
	2452  82801E USB Controller
	2453  82801E SMBus Controller
	2459  82801E Ethernet Controller 0
	245b  82801E IDE U100 Controller
	245d  82801E Ethernet Controller 1
	245e  82801E PCI Bridge
	2480  82801CA LPC Interface Controller
	2482  82801CA/CAM USB Controller #1
	2483  82801CA/CAM SMBus Controller
	2484  82801CA/CAM USB Controller #2
	2485  82801CA/CAM AC'97 Audio Controller
	2486  82801CA/CAM AC'97 Modem Controller
	2487  82801CA/CAM USB Controller #3
	248a  82801CAM IDE U100 Controller
	248b  82801CA Ultra ATA Storage Controller
	248c  82801CAM ISA Bridge (LPC)
	24c0  82801DB/DBL (ICH4/ICH4-L) LPC Interface Bridge
	24c1  82801DBL (ICH4-L) IDE Controller
	24c2  82801DB/DBL/DBM (ICH4/ICH4-L/ICH4-M) USB UHCI Controller #1
	24c3  82801DB/DBL/DBM (ICH4/ICH4-L/ICH4-M) SMBus Controller
	24c4  82801DB/DBL/DBM (ICH4/ICH4-L/ICH4-M) USB UHCI Controller #2
	24c5  82801DB/DBL/DBM (ICH4/ICH4-L/ICH4-M) AC'97 Audio Controller
	24c6  82801DB/DBL/DBM (ICH4/ICH4-L/ICH4-M) AC'97 Modem Controller
	24c7  82801DB/DBL/DBM (ICH4/ICH4-L/ICH4-M) USB UHCI Controller #3
	24ca  82801DBM (ICH4-M) IDE Controller
	24cb  82801DB (ICH4) IDE Controller
	24cc  82801DBM (ICH4-M) LPC Interface Bridge
	24cd  82801DB/DBM (ICH4/ICH4-M) USB2 EHCI Controller
	24d0  82801EB/ER (ICH5/ICH5R) LPC Interface Bridge
	24d1  82801EB (ICH5) SATA Controller
	24d2  82801EB/ER (ICH5/ICH5R) USB UHCI Controller #1
	24d3  82801EB/ER (ICH5/ICH5R) SMBus Controller
	24d4  82801EB/ER (ICH5/ICH5R) USB UHCI Controller #2
	24d5  82801EB/ER (ICH5/ICH5R) AC'97 Audio Controller
	24d6  82801EB/ER (ICH5/ICH5R) AC'97 Modem Controller
	24d7  82801EB/ER (ICH5/ICH5R) USB UHCI Controller #3
	24db  82801EB/ER (ICH5/ICH5R) IDE Controller
	24dc  82801EB (ICH5) LPC Interface Bridge
	24dd  82801EB/ER (ICH5/ICH5R) USB2 EHCI Controller
	24de  82801EB/ER (ICH5/ICH5R) USB UHCI Controller #4
	24df  82801ER (ICH5R) SATA Controller
	24f0  Omni-Path HFI Silicon 100 Series [discrete]
	24f1  Omni-Path HFI Silicon 100 Series [integrated]
	24f3  Wireless 8260
	24f4  Wireless 8260
	24fd  Wireless 8265 / 8275
	2500  82820 820 (Camino) Chipset Host Bridge (MCH)
	2501  82820 820 (Camino) Chipset Host Bridge (MCH)
	250b  82820 820 (Camino) Chipset Host Bridge
	250f  82820 820 (Camino) Chipset AGP Bridge
	2520  82805AA MTH Memory Translator Hub
	2521  82804AA MRH-S Memory Repeater Hub for SDRAM
	2530  82850 850 (Tehama) Chipset Host Bridge (MCH)
	2531  82860 860 (Wombat) Chipset Host Bridge (MCH)
	2532  82850 850 (Tehama) Chipset AGP Bridge
	2533  82860 860 (Wombat) Chipset AGP Bridge
	2534  82860 860 (Wombat) Chipset PCI Bridge
	2540  E7500 Memory Controller Hub
	2541  E7500/E7501 Host RASUM Controller
	2543  E7500/E7501 Hub Interface B PCI-to-PCI Bridge
	2544  E7500/E7501 Hub Interface B RASUM Controller
	2545  E7500/E7501 Hub Interface C PCI-to-PCI Bridge
	2546  E7500/E7501 Hub Interface C RASUM Controller
	2547  E7500/E7501 Hub Interface D PCI-to-PCI Bridge
	2548  E7500/E7501 Hub Interface D RASUM Controller
	254c  E7501 Memory Controller Hub
	2550  E7505 Memory Controller Hub
	2551  E7505/E7205 Series RAS Controller
	2552  E7505/E7205 PCI-to-AGP Bridge
	2553  E7505 Hub Interface B PCI-to-PCI Bridge
	2554  E7505 Hub Interface B PCI-to-PCI Bridge RAS Controller
	255d  E7205 Memory Controller Hub
	2560  82845G/GL[Brookdale-G]/GE/PE DRAM Controller/Host-Hub Interface
	2561  82845G/GL[Brookdale-G]/GE/PE Host-to-AGP Bridge
	2562  82845G/GL[Brookdale-G]/GE Chipset Integrated Graphics Device
	2570  82865G/PE/P DRAM Controller/Host-Hub Interface
	2571  82865G/PE/P AGP Bridge
	2572  82865G Integrated Graphics Controller
	2573  82865G/PE/P PCI to CSA Bridge
	2576  82865G/PE/P Processor to I/O Memory Interface
	2578  82875P/E7210 Memory Controller Hub
	2579  82875P Processor to AGP Controller
	257b  82875P/E7210 Processor to PCI to CSA Bridge
	257e  82875P/E7210 Processor to I/O Memory Interface
	2580  82915G/P/GV/GL/PL/910GL Memory Controller Hub
	2581  82915G/P/GV/GL/PL/910GL PCI Express Root Port
	2582  82915G/GV/910GL Integrated Graphics Controller
	2584  82925X/XE Memory Controller Hub
	2585  82925X/XE PCI Express Root Port
	2588  E7220/E7221 Memory Controller Hub
	2589  E7220/E7221 PCI Express Root Port
	258a  E7221 Integrated Graphics Controller
	2590  Mobile 915GM/PM/GMS/910GML Express Processor to DRAM Controller
	2591  Mobile 915GM/PM Express PCI Express Root Port
	2592  Mobile 915GM/GMS/910GML Express Graphics Controller
	25a1  6300ESB LPC Interface Controller
	25a2  6300ESB PATA Storage Controller
	25a3  6300ESB SATA Storage Controller
	25a4  6300ESB SMBus Controller
	25a6  6300ESB AC'97 Audio Controller
	25a7  6300ESB AC'97 Modem Controller
	25a9  6300ESB USB Universal Host Controller
	25aa  6300ESB USB Universal Host Controller
	25ab  6300ESB Watchdog Timer
	25ac  6300ESB I/O Advanced Programmable Interrupt Controller
	25ad  6300ESB USB2 Enhanced Host Controller
	25ae  6300ESB 64-bit PCI-X Bridge
	25b0  6300ESB SATA RAID Controller
	25c0  5000X Chipset Memory Controller Hub
	25d0  5000Z Chipset Memory Controller Hub
	25d4  5000V Chipset Memory Controller Hub
	25d8  5000P Chipset Memory Controller Hub
	25e2  5000 Series Chipset PCI Express x4 Port 2
	25e3  5000 Series Chipset PCI Express x4 Port 3
	25e4  5000 Series Chipset PCI Express x4 Port 4
	25e5  5000 Series Chipset PCI Express x4 Port 5
	25e6  5000 Series Chipset PCI Express x4 Port 6
	25e7  5000 Series Chipset PCI Express x4 Port 7
	25f0  5000 Series Chipset FSB Registers
	25f1  5000 Series Chipset Reserved Registers
	25f3  5000 Series Chipset Reserved Registers
	25f5  5000 Series Chipset FBD Registers
	25f6  5000 Series Chipset FBD Registers
	25f7  5000 Series Chipset PCI Express x8 Port 2-3
	25f8  5000 Series Chipset PCI Express x8 Port 4-5
	25f9  5000 Series Chipset PCI Express x8 Port 6-7
	25fa  5000X Chipset PCI Express x16 Port 4-7
	2600  E8500/E8501 Hub Interface 1.5
	2601  E8500/E8501 PCI Express x4 Port D
	2602  E8500/E8501 PCI Express x4 Port C0
	2603  E8500/E8501 PCI Express x4 Port C1
	2604  E8500/E8501 PCI Express x4 Port B0
	2605  E8500/E8501 PCI Express x4 Port B1
	2606  E8500/E8501 PCI Express x4 Port A0
	2607  E8500/E8501 PCI Express x4 Port A1
	2608  E8500/E8501 PCI Express x8 Port C
	2609  E8500/E8501 PCI Express x8 Port B
	260a  E8500/E8501 PCI Express x8 Port A
	260c  E8500/E8501 IMI Registers
	2610  E8500/E8501 FSB Registers
	2611  E8500/E8501 Address Mapping Registers
	2612  E8500/E8501 RAS Registers
	2613  E8500/E8501 Reserved Registers
	2614  E8500/E8501 Reserved Registers
	2615  E8500/E8501 Miscellaneous Registers
	2617  E8500/E8501 Reserved Registers
	2618  E8500/E8501 Reserved Registers
	2619  E8500/E8501 Reserved Registers
	261a  E8500/E8501 Reserved Registers
	261b  E8500/E8501 Reserved Registers
	261c  E8500/E8501 Reserved Registers
	261d  E8500/E8501 Reserved Registers
	261e  E8500/E8501 Reserved Registers
	2620  E8500/E8501 eXternal Memory Bridge
	2621  E8500/E8501 XMB Miscellaneous Registers
	2622  E8500/E8501 XMB Memory Interleaving Registers
	2623  E8500/E8501 XMB DDR Initialization and Calibration
	2624  E8500/E8501 XMB Reserved Registers
	2625  E8500/E8501 XMB Reserved Registers
	2626  E8500/E8501 XMB Reserved Registers
	2627  E8500/E8501 XMB Reserved Registers
	2640  82801FB/FR (ICH6/ICH6R) LPC Interface Bridge
	2641  82801FBM (ICH6M) LPC Interface Bridge
	2642  82801FW/FRW (ICH6W/ICH6RW) LPC Interface Bridge
	2651  82801FB/FW (ICH6/ICH6W) SATA Controller
	2652  82801FR/FRW (ICH6R/ICH6RW) SATA Controller
	2653  82801FBM (ICH6M) SATA Controller
	2658  82801FB/FBM/FR/FW/FRW (ICH6 Family) USB UHCI #1
	2659  82801FB/FBM/FR/FW/FRW (ICH6 Family) USB UHCI #2
	265a  82801FB/FBM/FR/FW/FRW (ICH6 Family) USB UHCI #3
	265b  82801FB/FBM/FR/FW/FRW (ICH6 Family) USB UHCI #4
	265c  82801FB/FBM/FR/FW/FRW (ICH6 Family) USB2 EHCI Controller
	2660  82801FB/FBM/FR/FW/FRW (ICH6 Family) PCI Express Port 1
	2662  82801FB/FBM/FR/FW/FRW (ICH6 Family) PCI Express Port 2
	2664  82801FB/FBM/FR/FW/FRW (ICH6 Family) PCI Express Port 3
	2666  82801FB/FBM/FR/FW/FRW (ICH6 Family) PCI Express Port 4
	2668  82801FB/FBM/FR/FW/FRW (ICH6 Family) High Definition Audio Controller
	266a  82801FB/FBM/FR/FW/FRW (ICH6 Family) SMBus Controller
	266c  82801FB/FBM/FR/FW/FRW (ICH6 Family) LAN Controller
	266d  82801FB/FBM/FR/FW/FRW (ICH6 Family) AC'97 Modem Controller
	266e  82801FB/FBM/FR/FW/FRW (ICH6 Family) AC'97 Audio Controller
	266f  82801FB/FBM/FR/FW/FRW (ICH6 Family) IDE Controller
	2670  631xESB/632xESB/3100 Chipset LPC Interface Controller
	2680  631xESB/632xESB/3100 Chipset SATA IDE Controller
	2681  631xESB/632xESB SATA AHCI Controller
	2682  631xESB/632xESB SATA RAID Controller
	2683  631xESB/632xESB SATA RAID Controller
	2688  631xESB/632xESB/3100 Chipset UHCI USB Controller #1
	2689  631xESB/632xESB/3100 Chipset UHCI USB Controller #2
	268a  631xESB/632xESB/3100 Chipset UHCI USB Controller #3
	268b  631xESB/632xESB/3100 Chipset UHCI USB Controller #4
	268c  631xESB/632xESB/3100 Chipset EHCI USB2 Controller
	2690  631xESB/632xESB/3100 Chipset PCI Express Root Port 1
	2692  631xESB/632xESB/3100 Chipset PCI Express Root Port 2
	2694  631xESB/632xESB/3100 Chipset PCI Express Root Port 3
	2696  631xESB/632xESB/3100 Chipset PCI Express Root Port 4
	2698  631xESB/632xESB AC '97 Audio Controller
	2699  631xESB/632xESB AC '97 Modem Controller
	269a  631xESB/632xESB High Definition Audio Controller
	269b  631xESB/632xESB/3100 Chipset SMBus Controller
	269e  631xESB/632xESB IDE Controller
	2770  82945G/GZ/P/PL Memory Controller Hub
	2771  82945G/GZ/P/PL PCI Express Root Port
	2772  82945G/GZ Integrated Graphics Controller
	2774  82955X Memory Controller Hub
	2775  82955X PCI Express Root Port
	2776  82945G/GZ Integrated Graphics Controller
	2778  E7230/3000/3010 Memory Controller Hub
	2779  E7230/3000/3010 PCI Express Root Port
	277a  82975X/3010 PCI Express Root Port
	277c  82975X Memory Controller Hub
	277d  82975X PCI Express Root Port
	2782  82915G Integrated Graphics Controller
	2792  Mobile 915GM/GMS/910GML Express Graphics Controller
	27a0  Mobile 945GM/PM/GMS, 943/940GML and 945GT Express Memory Controller Hub
	27a1  Mobile 945GM/PM/GMS, 943/940GML and 945GT Express PCI Express Root Port
	27a2  Mobile 945GM/GMS, 943/940GML Express Integrated Graphics Controller
	27a6  Mobile 945GM/GMS/GME, 943/940GML Express Integrated Graphics Controller
	27ac  Mobile 945GSE Express Memory Controller Hub
	27ad  Mobile 945GSE Express PCI Express Root Port
	27ae  Mobile 945GSE Express Integrated Graphics Controller
	27b0  82801GH (ICH7DH) LPC Interface Bridge
	27b8  82801GB/GR (ICH7 Family) LPC Interface Bridge
	27b9  82801GBM (ICH7-M) LPC Interface Bridge
	27bc  NM10 Family LPC Controller
	27bd  82801GHM (ICH7-M DH) LPC Interface Bridge
	27c0  NM10/ICH7 Family SATA Controller [IDE mode]
	27c1  NM10/ICH7 Family SATA Controller [AHCI mode]
	27c3  82801GR/GDH (ICH7R/ICH7DH) SATA Controller [RAID mode]
	27c4  82801GBM/GHM (ICH7-M Family) SATA Controller [IDE mode]
	27c5  82801GBM/GHM (ICH7-M Family) SATA Controller [AHCI mode]
	27c6  82801GHM (ICH7-M DH) SATA Controller [RAID mode]
	27c8  NM10/ICH7 Family USB UHCI Controller #1
	27c9  NM10/ICH7 Family USB UHCI Controller #2
	27ca  NM10/ICH7 Family USB UHCI Controller #3
	27cb  NM10/ICH7 Family USB UHCI Controller #4
	27cc  NM10/ICH7 Family USB2 EHCI Controller
	27d0  NM10/ICH7 Family PCI Express Port 1
	27d2  NM10/ICH7 Family PCI Express Port 2
	27d4  NM10/ICH7 Family PCI Express Port 3
	27d6  NM10/ICH7 Family PCI Express Port 4
	27d8  NM10/ICH7 Family High Definition Audio Controller
	27da  NM10/ICH7 Family SMBus Controller
	27dc  NM10/ICH7 Family LAN Controller
	27dd  82801G (ICH7 Family) AC'97 Modem Controller
	27de  82801G (ICH7 Family) AC'97 Audio Controller
	27df  82801G (ICH7 Family) IDE Controller
	27e0  82801GR/GH/GHM (ICH7 Family) PCI Express Port 5
	27e2  82801GR/GH/GHM (ICH7 Family) PCI Express Port 6
	2810  82801HB/HR (ICH8/R) LPC Interface Controller
	2811  82801HEM (ICH8M-E) LPC Interface Controller
	2812  82801HH (ICH8DH) LPC Interface Controller
	2814  82801HO (ICH8DO) LPC Interface Controller
	2815  82801HM (ICH8M) LPC Interface Controller
	2820  82801H (ICH8 Family) 4 port SATA Controller [IDE mode]
	2821  82801HR/HO/HH (ICH8R/DO/DH) 6 port SATA Controller [AHCI mode]
	2822  SATA Controller [RAID mode]
	2823  C610/X99 series chipset sSATA Controller [RAID mode]
	2824  82801HB (ICH8) 4 port SATA Controller [AHCI mode]
	2825  82801HR/HO/HH (ICH8R/DO/DH) 2 port SATA Controller [IDE mode]
	2826  C600/X79 series chipset SATA RAID Controller
	2827  C610/X99 series chipset sSATA Controller [RAID mode]
	2828  82801HM/HEM (ICH8M/ICH8M-E) SATA Controller [IDE mode]
	2829  82801HM/HEM (ICH8M/ICH8M-E) SATA Controller [AHCI mode]
	282a  82801 Mobile SATA Controller [RAID mode]
	2830  82801H (ICH8 Family) USB UHCI Controller #1
	2831  82801H (ICH8 Family) USB UHCI Controller #2
	2832  82801H (ICH8 Family) USB UHCI Controller #3
	2833  82801H (ICH8 Family) USB UHCI Controller #4
	2834  82801H (ICH8 Family) USB UHCI Controller #4
	2835  82801H (ICH8 Family) USB UHCI Controller #5
	2836  82801H (ICH8 Family) USB2 EHCI Controller #1
	283a  82801H (ICH8 Family) USB2 EHCI Controller #2
	283e  82801H (ICH8 Family) SMBus Controller
	283f  82801H (ICH8 Family) PCI Express Port 1
	2841  82801H (ICH8 Family) PCI Express Port 2
	2843  82801H (ICH8 Family) PCI Express Port 3
	2845  82801H (ICH8 Family) PCI Express Port 4
	2847  82801H (ICH8 Family) PCI Express Port 5
	2849  82801H (ICH8 Family) PCI Express Port 6
	284b  82801H (ICH8 Family) HD Audio Controller
	284f  82801H (ICH8 Family) Thermal Reporting Device
	2850  82801HM/HEM (ICH8M/ICH8M-E) IDE Controller
	2912  82801IH (ICH9DH) LPC Interface Controller
	2914  82801IO (ICH9DO) LPC Interface Controller
	2916  82801IR (ICH9R) LPC Interface Controller
	2917  ICH9M-E LPC Interface Controller
	2918  82801IB (ICH9) LPC Interface Controller
	2919  ICH9M LPC Interface Controller
	2920  82801IR/IO/IH (ICH9R/DO/DH) 4 port SATA Controller [IDE mode]
	2921  82801IB (ICH9) 2 port SATA Controller [IDE mode]
	2922  82801IR/IO/IH (ICH9R/DO/DH) 6 port SATA Controller [AHCI mode]
	2923  82801IB (ICH9) 4 port SATA Controller [AHCI mode]
	2925  82801IR/IO (ICH9R/DO) SATA Controller [RAID mode]
	2926  82801I (ICH9 Family) 2 port SATA Controller [IDE mode]
	2928  82801IBM/IEM (ICH9M/ICH9M-E) 2 port SATA Controller [IDE mode]
	2929  82801IBM/IEM (ICH9M/ICH9M-E) 4 port SATA Controller [AHCI mode]
	292c  82801IEM (ICH9M-E) SATA Controller [RAID mode]
	292d  82801IBM/IEM (ICH9M/ICH9M-E) 2 port SATA Controller [IDE mode]
	2930  82801I (ICH9 Family) SMBus Controller
	2932  82801I (ICH9 Family) Thermal Subsystem
	2934  82801I (ICH9 Family) USB UHCI Controller #1
	2935  82801I (ICH9 Family) USB UHCI Controller #2
	2936  82801I (ICH9 Family) USB UHCI Controller #3
	2937  82801I (ICH9 Family) USB UHCI Controller #4
	2938  82801I (ICH9 Family) USB UHCI Controller #5
	2939  82801I (ICH9 Family) USB UHCI Controller #6
	293a  82801I (ICH9 Family) USB2 EHCI Controller #1
	293c  82801I (ICH9 Family) USB2 EHCI Controller #2
	293e  82801I (ICH9 Family) HD Audio Controller
	2940  82801I (ICH9 Family) PCI Express Port 1
	2942  82801I (ICH9 Family) PCI Express Port 2
	2944  82801I (ICH9 Family) PCI Express Port 3
	2946  82801I (ICH9 Family) PCI Express Port 4
	2948  82801I (ICH9 Family) PCI Express Port 5
	294a  82801I (ICH9 Family) PCI Express Port 6
	294c  82566DC-2 Gigabit Network Connection
	2970  82946GZ/PL/GL Memory Controller Hub
	2971  82946GZ/PL/GL PCI Express Root Port
	2972  82946GZ/GL Integrated Graphics Controller
	2973  82946GZ/GL Integrated Graphics Controller
	2974  82946GZ/GL HECI Controller
	2975  82946GZ/GL HECI Controller
	2976  82946GZ/GL PT IDER Controller
	2977  82946GZ/GL KT Controller
	2980  82G35 Express DRAM Controller
	2981  82G35 Express PCI Express Root Port
	2982  82G35 Express Integrated Graphics Controller
	2983  82G35 Express Integrated Graphics Controller
	2984  82G35 Express HECI Controller
	2990  82Q963/Q965 Memory Controller Hub
	2991  82Q963/Q965 PCI Express Root Port
	2992  82Q963/Q965 Integrated Graphics Controller
	2993  82Q963/Q965 Integrated Graphics Controller
	2994  82Q963/Q965 HECI Controller
	2995  82Q963/Q965 HECI Controller
	2996  82Q963/Q965 PT IDER Controller
	2997  82Q963/Q965 KT Controller
	29a0  82P965/G965 Memory Controller Hub
	29a1  82P965/G965 PCI Express Root Port
	29a2  82G965 Integrated Graphics Controller
	29a3  82G965 Integrated Graphics Controller
	29a4  82P965/G965 HECI Controller
	29a5  82P965/G965 HECI Controller
	29a6  82P965/G965 PT IDER Controller
	29a7  82P965/G965 KT Controller
	29b0  82Q35 Express DRAM Controller
	29b1  82Q35 Express PCI Express Root Port
	29b2  82Q35 Express Integrated Graphics Controller
	29b3  82Q35 Express Integrated Graphics Controller
	29b4  82Q35 Express MEI Controller
	29b5  82Q35 Express MEI Controller
	29b6  82Q35 Express PT IDER Controller
	29b7  82Q35 Express Serial KT Controller
	29c0  82G33/G31/P35/P31 Express DRAM Controller
	29c1  82G33/G31/P35/P31 Express PCI Express Root Port
	29c2  82G33/G31 Express Integrated Graphics Controller
	29c3  82G33/G31 Express Integrated Graphics Controller
	29c4  82G33/G31/P35/P31 Express MEI Controller
	29c5  82G33/G31/P35/P31 Express MEI Controller
	29c6  82G33/G31/P35/P31 Express PT IDER Controller
	29c7  82G33/G31/P35/P31 Express Serial KT Controller
	29cf  Virtual HECI Controller
	29d0  82Q33 Express DRAM Controller
	29d1  82Q33 Express PCI Express Root Port
	29d2  82Q33 Express Integrated Graphics Controller
	29d3  82Q33 Express Integrated Graphics Controller
	29d4  82Q33 Express MEI Controller
	29d5  82Q33 Express MEI Controller
	29d6  82Q33 Express PT IDER Controller
	29d7  82Q33 Express Serial KT Controller
	29e0  82X38/X48 Express DRAM Controller
	29e1  82X38/X48 Express Host-Primary PCI Express Bridge
	29e4  82X38/X48 Express MEI Controller
	29e5  82X38/X48 Express MEI Controller
	29e6  82X38/X48 Express PT IDER Controller
	29e7  82X38/X48 Express Serial KT Controller
	29e9  82X38/X48 Express Host-Secondary PCI Express Bridge
	29f0  3200/3210 Chipset DRAM Controller
	29f1  3200/3210 Chipset Host-Primary PCI Express Bridge
	29f4  3200/3210 Chipset MEI Controller
	29f5  3200/3210 Chipset MEI Controller
	29f6  3200/3210 Chipset PT IDER Controller
	29f7  3200/3210 Chipset Serial KT Controller
	29f9  3210 Chipset Host-Secondary PCI Express Bridge
	2a00  Mobile PM965/GM965/GL960 Memory Controller Hub
	2a01  Mobile PM965/GM965/GL960 PCI Express Root Port
	2a02  Mobile GM965/GL960 Integrated Graphics Controller (primary)
	2a03  Mobile GM965/GL960 Integrated Graphics Controller (secondary)
	2a04  Mobile PM965/GM965 MEI Controller
	2a05  Mobile PM965/GM965 MEI Controller
	2a06  Mobile PM965/GM965 PT IDER Controller
	2a07  Mobile PM965/GM965 KT Controller
	2a10  Mobile GME965/GLE960 Memory Controller Hub
	2a11  Mobile GME965/GLE960 PCI Express Root Port
	2a12  Mobile GME965/GLE960 Integrated Graphics Controller
	2a13  Mobile GME965/GLE960 Integrated Graphics Controller
	2a14  Mobile GME965/GLE960 MEI Controller
	2a15  Mobile GME965/GLE960 MEI Controller
	2a16  Mobile GME965/GLE960 PT IDER Controller
	2a17  Mobile GME965/GLE960 KT Controller
	2a40  Mobile 4 Series Chipset Memory Controller Hub
	2a41  Mobile 4 Series Chipset PCI Express Graphics Port
	2a42  Mobile 4 Series Chipset Integrated Graphics Controller
	2a43  Mobile 4 Series Chipset Integrated Graphics Controller
	2a44  Mobile 4 Series Chipset MEI Controller
	2a45  Mobile 4 Series Chipset MEI Controller
	2a46  Mobile 4 Series Chipset PT IDER Controller
	2a47  Mobile 4 Series Chipset AMT SOL Redirection
	2a50  Cantiga MEI Controller
	2a51  Cantiga MEI Controller
	2a52  Cantiga PT IDER Controller
	2a53  Cantiga AMT SOL Redirection
	2b00  Xeon Processor E7 Product Family System Configuration Controller 1
	2b02  Xeon Processor E7 Product Family System Configuration Controller 2
	2b04  Xeon Processor E7 Product Family Power Controller
	2b08  Xeon Processor E7 Product Family Caching Agent 0
	2b0c  Xeon Processor E7 Product Family Caching Agent 1
	2b10  Xeon Processor E7 Product Family QPI Home Agent 0
	2b13  Xeon Processor E7 Product Family Memory Controller 0c
	2b14  Xeon Processor E7 Product Family Memory Controller 0a
	2b16  Xeon Processor E7 Product Family Memory Controller 0b
	2b18  Xeon Processor E7 Product Family QPI Home Agent 1
	2b1b  Xeon Processor E7 Product Family Memory Controller 1c
	2b1c  Xeon Processor E7 Product Family Memory Controller 1a
	2b1e  Xeon Processor E7 Product Family Memory Controller 1b
	2b20  Xeon Processor E7 Product Family Last Level Cache Coherence Engine 0
	2b22  Xeon Processor E7 Product Family System Configuration Controller 3
	2b24  Xeon Processor E7 Product Family Last Level Cache Coherence Engine 1
	2b28  Xeon Processor E7 Product Family Last Level Cache Coherence Engine 2
	2b2a  Xeon Processor E7 Product Family System Configuration Controller 4
	2b2c  Xeon Processor E7 Product Family Last Level Cache Coherence Engine 3
	2b30  Xeon Processor E7 Product Family Last Level Cache Coherence Engine 4
	2b34  Xeon Processor E7 Product Family Last Level Cache Coherence Engine 5
	2b38  Xeon Processor E7 Product Family Last Level Cache Coherence Engine 6
	2b3c  Xeon Processor E7 Product Family Last Level Cache Coherence Engine 7
	2b40  Xeon Processor E7 Product Family QPI Router Port 0-1
	2b42  Xeon Processor E7 Product Family QPI Router Port 2-3
	2b44  Xeon Processor E7 Product Family QPI Router Port 4-5
	2b46  Xeon Processor E7 Product Family QPI Router Port 6-7
	2b48  Xeon Processor E7 Product Family Test and Debug 0
	2b4c  Xeon Processor E7 Product Family Test and Debug 1
	2b50  Xeon Processor E7 Product Family QPI Physical Port 0: REUT control/status
	2b52  Xeon Processor E7 Product Family QPI Physical Port 0: Misc. control/status
	2b54  Xeon Processor E7 Product Family QPI Physical Port 1: REUT control/status
	2b56  Xeon Processor E7 Product Family QPI Physical Port 1: Misc. control/status
	2b58  Xeon Processor E7 Product Family QPI Physical Port 2: REUT control/status
	2b5a  Xeon Processor E7 Product Family QPI Physical Port 2: Misc. control/status
	2b5c  Xeon Processor E7 Product Family QPI Physical Port 3: REUT control/status
	2b5e  Xeon Processor E7 Product Family QPI Physical Port 3: Misc. control/status
	2b60  Xeon Processor E7 Product Family SMI Physical Port 0: REUT control/status
	2b62  Xeon Processor E7 Product Family SMI Physical Port 0: Misc control/status
	2b64  Xeon Processor E7 Product Family SMI Physical Port 1: REUT control/status
	2b66  Xeon Processor E7 Product Family SMI Physical Port 1: Misc control/status
	2b68  Xeon Processor E7 Product Family Last Level Cache Coherence Engine 8
	2b6c  Xeon Processor E7 Product Family Last Level Cache Coherence Engine 9
	2c01  Xeon 5500/Core i7 QuickPath Architecture System Address Decoder
	2c10  Xeon 5500/Core i7 QPI Link 0
	2c11  Xeon 5500/Core i7 QPI Physical 0
	2c14  Xeon 5500/Core i7 QPI Link 1
	2c15  Xeon 5500/Core i7 QPI Physical 1
	2c18  Xeon 5500/Core i7 Integrated Memory Controller
	2c19  Xeon 5500/Core i7 Integrated Memory Controller Target Address Decoder
	2c1a  Xeon 5500/Core i7 Integrated Memory Controller RAS Registers
	2c1c  Xeon 5500/Core i7 Integrated Memory Controller Test Registers
	2c20  Xeon 5500/Core i7 Integrated Memory Controller Channel 0 Control Registers
	2c21  Xeon 5500/Core i7 Integrated Memory Controller Channel 0 Address Registers
	2c22  Xeon 5500/Core i7 Integrated Memory Controller Channel 0 Rank Registers
	2c23  Xeon 5500/Core i7 Integrated Memory Controller Channel 0 Thermal Control Registers
	2c28  Xeon 5500/Core i7 Integrated Memory Controller Channel 1 Control Registers
	2c29  Xeon 5500/Core i7 Integrated Memory Controller Channel 1 Address Registers
	2c2a  Xeon 5500/Core i7 Integrated Memory Controller Channel 1 Rank Registers
	2c2b  Xeon 5500/Core i7 Integrated Memory Controller Channel 1 Thermal Control Registers
	2c30  Xeon 5500/Core i7 Integrated Memory Controller Channel 2 Control Registers
	2c31  Xeon 5500/Core i7 Integrated Memory Controller Channel 2 Address Registers
	2c32  Xeon 5500/Core i7 Integrated Memory Controller Channel 2 Rank Registers
	2c33  Xeon 5500/Core i7 Integrated Memory Controller Channel 2 Thermal Control Registers
	2c40  Xeon 5500/Core i7 QuickPath Architecture Generic Non-Core Registers
	2c41  Xeon 5500/Core i7 QuickPath Architecture Generic Non-Core Registers
	2c50  Core Processor QuickPath Architecture Generic Non-Core Registers
	2c51  Core Processor QuickPath Architecture Generic Non-Core Registers
	2c52  Core Processor QuickPath Architecture Generic Non-Core Registers
	2c53  Core Processor QuickPath Architecture Generic Non-Core Registers
	2c54  Core Processor QuickPath Architecture Generic Non-Core Registers
	2c55  Core Processor QuickPath Architecture Generic Non-Core Registers
	2c56  Core Processor QuickPath Architecture Generic Non-Core Registers
	2c57  Core Processor QuickPath Architecture Generic Non-Core Registers
	2c58  Xeon C5500/C3500 QPI Generic Non-core Registers
	2c59  Xeon C5500/C3500 QPI Generic Non-core Registers
	2c5a  Xeon C5500/C3500 QPI Generic Non-core Registers
	2c5b  Xeon C5500/C3500 QPI Generic Non-core Registers
	2c5c  Xeon C5500/C3500 QPI Generic Non-core Registers
	2c5d  Xeon C5500/C3500 QPI Generic Non-core Registers
	2c5e  Xeon C5500/C3500 QPI Generic Non-core Registers
	2c5f  Xeon C5500/C3500 QPI Generic Non-core Registers
	2c61  Core Processor QuickPath Architecture Generic Non-core Registers
	2c62  Core Processor QuickPath Architecture Generic Non-core Registers
	2c70  Xeon 5600 Series QuickPath Architecture Generic Non-core Registers
	2c81  Core Processor QuickPath Architecture System Address Decoder
	2c90  Core Processor QPI Link 0
	2c91  Core Processor QPI Physical 0
	2c98  Core Processor Integrated Memory Controller
	2c99  Core Processor Integrated Memory Controller Target Address Decoder
	2c9a  Core Processor Integrated Memory Controller Test Registers
	2c9c  Core Processor Integrated Memory Controller Test Registers
	2ca0  Core Processor Integrated Memory Controller Channel 0 Control Registers
	2ca1  Core Processor Integrated Memory Controller Channel 0 Address Registers
	2ca2  Core Processor Integrated Memory Controller Channel 0 Rank Registers
	2ca3  Core Processor Integrated Memory Controller Channel 0 Thermal Control Registers
	2ca8  Core Processor Integrated Memory Controller Channel 1 Control Registers
	2ca9  Core Processor Integrated Memory Controller Channel 1 Address Registers
	2caa  Core Processor Integrated Memory Controller Channel 1 Rank Registers
	2cab  Core Processor Integrated Memory Controller Channel 1 Thermal Control Registers
	2cc1  Xeon C5500/C3500 QPI System Address Decoder
	2cd0  Xeon C5500/C3500 QPI Link 0
	2cd1  Xeon C5500/C3500 QPI Physical 0
	2cd4  Xeon C5500/C3500 QPI Link 1
	2cd5  Xeon C5500/C3500 QPI Physical 1
	2cd8  Xeon C5500/C3500 Integrated Memory Controller Registers
	2cd9  Xeon C5500/C3500 Integrated Memory Controller Target Address Decoder
	2cda  Xeon C5500/C3500 Integrated Memory Controller RAS Registers
	2cdc  Xeon C5500/C3500 Integrated Memory Controller Test Registers
	2ce0  Xeon C5500/C3500 Integrated Memory Controller Channel 0 Control
	2ce1  Xeon C5500/C3500 Integrated Memory Controller Channel 0 Address
	2ce2  Xeon C5500/C3500 Integrated Memory Controller Channel 0 Rank
	2ce3  Xeon C5500/C3500 Integrated Memory Controller Channel 0 Thermal Control
	2ce8  Xeon C5500/C3500 Integrated Memory Controller Channel 1 Control
	2ce9  Xeon C5500/C3500 Integrated Memory Controller Channel 1 Address
	2cea  Xeon C5500/C3500 Integrated Memory Controller Channel 1 Rank
	2ceb  Xeon C5500/C3500 Integrated Memory Controller Channel 1 Thermal Control
	2cf0  Xeon C5500/C3500 Integrated Memory Controller Channel 2 Control
	2cf1  Xeon C5500/C3500 Integrated Memory Controller Channel 2 Address
	2cf2  Xeon C5500/C3500 Integrated Memory Controller Channel 2 Rank
	2cf3  Xeon C5500/C3500 Integrated Memory Controller Channel 2 Thermal Control
	2d01  Core Processor QuickPath Architecture System Address Decoder
	2d10  Core Processor QPI Link 0
	2d11  1st Generation Core i3/5/7 Processor QPI Physical 0
	2d12  1st Generation Core i3/5/7 Processor Reserved
	2d13  1st Generation Core i3/5/7 Processor Reserved
	2d81  Xeon 5600 Series QuickPath Architecture System Address Decoder
	2d90  Xeon 5600 Series QPI Link 0
	2d91  Xeon 5600 Series QPI Physical 0
	2d92  Xeon 5600 Series Mirror Port Link 0
	2d93  Xeon 5600 Series Mirror Port Link 1
	2d94  Xeon 5600 Series QPI Link 1
	2d95  Xeon 5600 Series QPI Physical 1
	2d98  Xeon 5600 Series Integrated Memory Controller Registers
	2d99  Xeon 5600 Series Integrated Memory Controller Target Address Decoder
	2d9a  Xeon 5600 Series Integrated Memory Controller RAS Registers
	2d9c  Xeon 5600 Series Integrated Memory Controller Test Registers
	2da0  Xeon 5600 Series Integrated Memory Controller Channel 0 Control
	2da1  Xeon 5600 Series Integrated Memory Controller Channel 0 Address
	2da2  Xeon 5600 Series Integrated Memory Controller Channel 0 Rank
	2da3  Xeon 5600 Series Integrated Memory Controller Channel 0 Thermal Control
	2da8  Xeon 5600 Series Integrated Memory Controller Channel 1 Control
	2da9  Xeon 5600 Series Integrated Memory Controller Channel 1 Address
	2daa  Xeon 5600 Series Integrated Memory Controller Channel 1 Rank
	2dab  Xeon 5600 Series Integrated Memory Controller Channel 1 Thermal Control
	2db0  Xeon 5600 Series Integrated Memory Controller Channel 2 Control
	2db1  Xeon 5600 Series Integrated Memory Controller Channel 2 Address
	2db2  Xeon 5600 Series Integrated Memory Controller Channel 2 Rank
	2db3  Xeon 5600 Series Integrated Memory Controller Channel 2 Thermal Control
	2e00  4 Series Chipset DRAM Controller
	2e01  4 Series Chipset PCI Express Root Port
	2e02  4 Series Chipset Integrated Graphics Controller
	2e03  4 Series Chipset Integrated Graphics Controller
	2e04  4 Series Chipset HECI Controller
	2e05  4 Series Chipset HECI Controller
	2e06  4 Series Chipset PT IDER Controller
	2e07  4 Series Chipset Serial KT Controller
	2e10  4 Series Chipset DRAM Controller
	2e11  4 Series Chipset PCI Express Root Port
	2e12  4 Series Chipset Integrated Graphics Controller
	2e13  4 Series Chipset Integrated Graphics Controller
	2e14  4 Series Chipset HECI Controller
	2e15  4 Series Chipset HECI Controller
	2e16  4 Series Chipset PT IDER Controller
	2e17  4 Series Chipset Serial KT Controller
	2e20  4 Series Chipset DRAM Controller
	2e21  4 Series Chipset PCI Express Root Port
	2e22  4 Series Chipset Integrated Graphics Controller
	2e23  4 Series Chipset Integrated Graphics Controller
	2e24  4 Series Chipset HECI Controller
	2e25  4 Series Chipset HECI Controller
	2e26  4 Series Chipset PT IDER Controller
	2e27  4 Series Chipset Serial KT Controller
	2e29  4 Series Chipset PCI Express Root Port
	2e30  4 Series Chipset DRAM Controller
	2e31  4 Series Chipset PCI Express Root Port
	2e32  4 Series Chipset Integrated Graphics Controller
	2e33  4 Series Chipset Integrated Graphics Controller
	2e34  4 Series Chipset HECI Controller
	2e35  4 Series Chipset HECI Controller
	2e36  4 Series Chipset PT IDER Controller
	2e37  4 Series Chipset Serial KT Controller
	2e40  4 Series Chipset DRAM Controller
	2e41  4 Series Chipset PCI Express Root Port
	2e42  4 Series Chipset Integrated Graphics Controller
	2e43  4 Series Chipset Integrated Graphics Controller
	2e44  4 Series Chipset HECI Controller
	2e45  4 Series Chipset HECI Controller
	2e46  4 Series Chipset PT IDER Controller
	2e47  4 Series Chipset Serial KT Controller
	2e50  CE Media Processor CE3100
	2e52  CE Media Processor Clock and Reset Controller
	2e58  CE Media Processor Interrupt Controller
	2e5a  CE Media Processor CE3100 A/V Bridge
	2e5b  Graphics Media Accelerator 500 Graphics
	2e5c  CE Media Processor Video Decoder
	2e5d  CE Media Processor Transport Stream Interface
	2e5e  CE Media Processor Transport Stream Processor 0
	2e5f  CE Media Processor Audio DSP
	2e60  CE Media Processor Audio Interfaces
	2e61  CE Media Processor Video Display Controller
	2e62  CE Media Processor Video Processing Unit
	2e63  CE Media Processor HDMI Tx Interface
	2e65  CE Media Processor Expansion Bus Interface
	2e66  CE Media Processor UART
	2e67  CE Media Processor General Purpose I/Os
	2e68  CE Media Processor I2C Interface
	2e69  CE Media Processor Smart Card Interface
	2e6a  CE Media Processor SPI Master Interface
	2e6e  CE Media Processor Gigabit Ethernet Controller
	2e6f  CE Media Processor Media Timing Unit
	2e70  CE Media Processor USB
	2e71  CE Media Processor SATA
	2e73  CE Media Processor CE3100 PCI Express
	2e90  4 Series Chipset DRAM Controller
	2e91  4 Series Chipset PCI Express Root Port
	2e92  4 Series Chipset Integrated Graphics Controller
	2e93  4 Series Chipset Integrated Graphics Controller
	2e94  4 Series Chipset HECI Controller
	2e95  4 Series Chipset HECI Controller
	2e96  4 Series Chipset PT IDER Controller
	2f00  Xeon E7 v3/Xeon E5 v3/Core i7 DMI2
	2f01  Xeon E7 v3/Xeon E5 v3/Core i7 PCI Express Root Port 0
	2f02  Xeon E7 v3/Xeon E5 v3/Core i7 PCI Express Root Port 1
	2f03  Xeon E7 v3/Xeon E5 v3/Core i7 PCI Express Root Port 1
	2f04  Xeon E7 v3/Xeon E5 v3/Core i7 PCI Express Root Port 2
	2f05  Xeon E7 v3/Xeon E5 v3/Core i7 PCI Express Root Port 2
	2f06  Xeon E7 v3/Xeon E5 v3/Core i7 PCI Express Root Port 2
	2f07  Xeon E7 v3/Xeon E5 v3/Core i7 PCI Express Root Port 2
	2f08  Xeon E7 v3/Xeon E5 v3/Core i7 PCI Express Root Port 3
	2f09  Xeon E7 v3/Xeon E5 v3/Core i7 PCI Express Root Port 3
	2f0a  Xeon E7 v3/Xeon E5 v3/Core i7 PCI Express Root Port 3
	2f0b  Xeon E7 v3/Xeon E5 v3/Core i7 PCI Express Root Port 3
	2f0d  Haswell Xeon Non-Transparent Bridge (Back-to-back)
	2f0e  Haswell Xeon Non-Transparent Bridge (Primary Side)
	2f0f  Haswell Xeon Non-Transparent Bridge (Secondary Side)
	2f10  Xeon E7 v3/Xeon E5 v3/Core i7 IIO Debug
	2f11  Xeon E7 v3/Xeon E5 v3/Core i7 IIO Debug
	2f12  Xeon E7 v3/Xeon E5 v3/Core i7 IIO Debug
	2f13  Xeon E7 v3/Xeon E5 v3/Core i7 IIO Debug
	2f14  Xeon E7 v3/Xeon E5 v3/Core i7 IIO Debug
	2f15  Xeon E7 v3/Xeon E5 v3/Core i7 IIO Debug
	2f16  Xeon E7 v3/Xeon E5 v3/Core i7 IIO Debug
	2f17  Xeon E7 v3/Xeon E5 v3/Core i7 IIO Debug
	2f18  Xeon E7 v3/Xeon E5 v3/Core i7 IIO Debug
	2f19  Xeon E7 v3/Xeon E5 v3/Core i7 IIO Debug
	2f1a  Xeon E7 v3/Xeon E5 v3/Core i7 IIO Debug
	2f1b  Xeon E7 v3/Xeon E5 v3/Core i7 IIO Debug
	2f1c  Xeon E7 v3/Xeon E5 v3/Core i7 IIO Debug
	2f1d  Xeon E7 v3/Xeon E5 v3/Core i7 PCIe Ring Interface
	2f1e  Xeon E7 v3/Xeon E5 v3/Core i7 Scratchpad & Semaphore Registers
	2f1f  Xeon E7 v3/Xeon E5 v3/Core i7 Scratchpad & Semaphore Registers
	2f20  Xeon E7 v3/Xeon E5 v3/Core i7 DMA Channel 0
	2f21  Xeon E7 v3/Xeon E5 v3/Core i7 DMA Channel 1
	2f22  Xeon E7 v3/Xeon E5 v3/Core i7 DMA Channel 2
	2f23  Xeon E7 v3/Xeon E5 v3/Core i7 DMA Channel 3
	2f24  Xeon E7 v3/Xeon E5 v3/Core i7 DMA Channel 4
	2f25  Xeon E7 v3/Xeon E5 v3/Core i7 DMA Channel 5
	2f26  Xeon E7 v3/Xeon E5 v3/Core i7 DMA Channel 6
	2f27  Xeon E7 v3/Xeon E5 v3/Core i7 DMA Channel 7
	2f28  Xeon E7 v3/Xeon E5 v3/Core i7 Address Map, VTd_Misc, System Management
	2f29  Xeon E7 v3/Xeon E5 v3/Core i7 Hot Plug
	2f2a  Xeon E7 v3/Xeon E5 v3/Core i7 RAS, Control Status and Global Errors
	2f2c  Xeon E7 v3/Xeon E5 v3/Core i7 I/O APIC
	2f2e  Xeon E7 v3/Xeon E5 v3/Core i7 RAID 5/6
	2f2f  Xeon E7 v3/Xeon E5 v3/Core i7 RAID 5/6
	2f30  Xeon E7 v3/Xeon E5 v3/Core i7 Home Agent 0
	2f32  Xeon E7 v3/Xeon E5 v3/Core i7 QPI Link 0
	2f33  Xeon E7 v3/Xeon E5 v3/Core i7 QPI Link 1
	2f34  Xeon E7 v3/Xeon E5 v3/Core i7 PCIe Ring Interface
	2f36  Xeon E7 v3/Xeon E5 v3/Core i7 R3 QPI Link 0 & 1 Monitoring
	2f37  Xeon E7 v3/Xeon E5 v3/Core i7 R3 QPI Link 0 & 1 Monitoring
	2f38  Xeon E7 v3/Xeon E5 v3/Core i7 Home Agent 1
	2f39  Xeon E7 v3/Xeon E5 v3/Core i7 I/O Performance Monitoring
	2f3a  Xeon E7 v3/Xeon E5 v3/Core i7 QPI Link 2
	2f3e  Xeon E7 v3/Xeon E5 v3/Core i7 R3 QPI Link 2 Monitoring
	2f3f  Xeon E7 v3/Xeon E5 v3/Core i7 R3 QPI Link 2 Monitoring
	2f40  Xeon E7 v3/Xeon E5 v3/Core i7 QPI Link 2
	2f41  Xeon E7 v3/Xeon E5 v3/Core i7 R3 QPI Link 2 Monitoring
	2f43  Xeon E7 v3/Xeon E5 v3/Core i7 QPI Link 2
	2f45  Xeon E7 v3/Xeon E5 v3/Core i7 QPI Link 2 Debug
	2f46  Xeon E7 v3/Xeon E5 v3/Core i7 QPI Link 2 Debug
	2f47  Xeon E7 v3/Xeon E5 v3/Core i7 QPI Link 2 Debug
	2f60  Xeon E7 v3/Xeon E5 v3/Core i7 Home Agent 1
	2f68  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 1 Target Address, Thermal & RAS Registers
	2f6a  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 1 Channel Target Address Decoder
	2f6b  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 1 Channel Target Address Decoder
	2f6c  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 1 Channel Target Address Decoder
	2f6d  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 1 Channel Target Address Decoder
	2f6e  Xeon E7 v3/Xeon E5 v3/Core i7 DDRIO Channel 2/3 Broadcast
	2f6f  Xeon E7 v3/Xeon E5 v3/Core i7 DDRIO Global Broadcast
	2f70  Xeon E7 v3/Xeon E5 v3/Core i7 Home Agent 0 Debug
	2f71  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 0 Target Address, Thermal & RAS Registers
	2f76  Xeon E7 v3/Xeon E5 v3/Core i7 E3 QPI Link Debug
	2f78  Xeon E7 v3/Xeon E5 v3/Core i7 Home Agent 1 Debug
	2f79  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 1 Target Address, Thermal & RAS Registers
	2f7d  Xeon E7 v3/Xeon E5 v3/Core i7 Scratchpad & Semaphore Registers
	2f7e  Xeon E7 v3/Xeon E5 v3/Core i7 E3 QPI Link Debug
	2f80  Xeon E7 v3/Xeon E5 v3/Core i7 QPI Link 0
	2f81  Xeon E7 v3/Xeon E5 v3/Core i7 R3 QPI Link 0 & 1 Monitoring
	2f83  Xeon E7 v3/Xeon E5 v3/Core i7 QPI Link 0
	2f85  Xeon E7 v3/Xeon E5 v3/Core i7 QPI Link 0 Debug
	2f86  Xeon E7 v3/Xeon E5 v3/Core i7 QPI Link 0 Debug
	2f87  Xeon E7 v3/Xeon E5 v3/Core i7 QPI Link 0 Debug
	2f88  Xeon E7 v3/Xeon E5 v3/Core i7 VCU
	2f8a  Xeon E7 v3/Xeon E5 v3/Core i7 VCU
	2f90  Xeon E7 v3/Xeon E5 v3/Core i7 QPI Link 1
	2f93  Xeon E7 v3/Xeon E5 v3/Core i7 QPI Link 1
	2f95  Xeon E7 v3/Xeon E5 v3/Core i7 QPI Link 1 Debug
	2f96  Xeon E7 v3/Xeon E5 v3/Core i7 QPI Link 1 Debug
	2f98  Xeon E7 v3/Xeon E5 v3/Core i7 Power Control Unit
	2f99  Xeon E7 v3/Xeon E5 v3/Core i7 Power Control Unit
	2f9a  Xeon E7 v3/Xeon E5 v3/Core i7 Power Control Unit
	2f9c  Xeon E7 v3/Xeon E5 v3/Core i7 Power Control Unit
	2fa0  Xeon E7 v3/Xeon E5 v3/Core i7 Home Agent 0
	2fa8  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 0 Target Address, Thermal & RAS Registers
	2faa  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 0 Channel Target Address Decoder
	2fab  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 0 Channel Target Address Decoder
	2fac  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 0 Channel Target Address Decoder
	2fad  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 0 Channel Target Address Decoder
	2fae  Xeon E7 v3/Xeon E5 v3/Core i7 DDRIO Channel 0/1 Broadcast
	2faf  Xeon E7 v3/Xeon E5 v3/Core i7 DDRIO Global Broadcast
	2fb0  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 0 Channel 0 Thermal Control
	2fb1  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 0 Channel 1 Thermal Control
	2fb2  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 0 Channel 0 ERROR Registers
	2fb3  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 0 Channel 1 ERROR Registers
	2fb4  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 0 Channel 2 Thermal Control
	2fb5  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 0 Channel 3 Thermal Control
	2fb6  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 0 Channel 2 ERROR Registers
	2fb7  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 0 Channel 3 ERROR Registers
	2fb8  Xeon E7 v3/Xeon E5 v3/Core i7 DDRIO (VMSE) 2 & 3
	2fb9  Xeon E7 v3/Xeon E5 v3/Core i7 DDRIO (VMSE) 2 & 3
	2fba  Xeon E7 v3/Xeon E5 v3/Core i7 DDRIO (VMSE) 2 & 3
	2fbb  Xeon E7 v3/Xeon E5 v3/Core i7 DDRIO (VMSE) 2 & 3
	2fbc  Xeon E7 v3/Xeon E5 v3/Core i7 DDRIO (VMSE) 0 & 1
	2fbd  Xeon E7 v3/Xeon E5 v3/Core i7 DDRIO (VMSE) 0 & 1
	2fbe  Xeon E7 v3/Xeon E5 v3/Core i7 DDRIO (VMSE) 0 & 1
	2fbf  Xeon E7 v3/Xeon E5 v3/Core i7 DDRIO (VMSE) 0 & 1
	2fc0  Xeon E7 v3/Xeon E5 v3/Core i7 Power Control Unit
	2fc1  Xeon E7 v3/Xeon E5 v3/Core i7 Power Control Unit
	2fc2  Xeon E7 v3/Xeon E5 v3/Core i7 Power Control Unit
	2fc3  Xeon E7 v3/Xeon E5 v3/Core i7 Power Control Unit
	2fc4  Xeon E7 v3/Xeon E5 v3/Core i7 Power Control Unit
	2fc5  Xeon E7 v3/Xeon E5 v3/Core i7 Power Control Unit
	2fd0  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 1 Channel 0 Thermal Control
	2fd1  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 1 Channel 1 Thermal Control
	2fd2  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 1 Channel 0 ERROR Registers
	2fd3  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 1 Channel 1 ERROR Registers
	2fd4  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 1 Channel 2 Thermal Control
	2fd5  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 1 Channel 3 Thermal Control
	2fd6  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 1 Channel 2 ERROR Registers
	2fd7  Xeon E7 v3/Xeon E5 v3/Core i7 Integrated Memory Controller 1 Channel 3 ERROR Registers
	2fe0  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2fe1  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2fe2  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2fe3  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2fe4  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2fe5  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2fe6  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2fe7  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2fe8  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2fe9  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2fea  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2feb  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2fec  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2fed  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2fee  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2fef  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2ff0  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2ff1  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2ff2  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2ff3  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2ff4  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2ff5  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2ff6  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2ff7  Xeon E7 v3/Xeon E5 v3/Core i7 Unicast Registers
	2ff8  Xeon E7 v3/Xeon E5 v3/Core i7 Buffered Ring Agent
	2ff9  Xeon E7 v3/Xeon E5 v3/Core i7 Buffered Ring Agent
	2ffa  Xeon E7 v3/Xeon E5 v3/Core i7 Buffered Ring Agent
	2ffb  Xeon E7 v3/Xeon E5 v3/Core i7 Buffered Ring Agent
	2ffc  Xeon E7 v3/Xeon E5 v3/Core i7 System Address Decoder & Broadcast Registers
	2ffd  Xeon E7 v3/Xeon E5 v3/Core i7 System Address Decoder & Broadcast Registers
	2ffe  Xeon E7 v3/Xeon E5 v3/Core i7 System Address Decoder & Broadcast Registers
	3165  Wireless 3165
	3166  Dual Band Wireless-AC 3165 Plus Bluetooth
	3200  GD31244 PCI-X SATA HBA
	3310  IOP348 I/O Processor
	3313  IOP348 I/O Processor (SL8e) in IOC Mode SAS/SATA
	331b  IOP348 I/O Processor (SL8x) in IOC Mode SAS/SATA
	3331  IOC340 I/O Controller (VV8e) SAS/SATA
	3339  IOC340 I/O Controller (VV8x) SAS/SATA
	3340  82855PM Processor to I/O Controller
	3341  82855PM Processor to AGP Controller
	3363  IOC340 I/O Controller in IOC Mode SAS/SATA
	3382  81342 [Chevelon] I/O Processor (ATUe)
	33c3  IOP348 I/O Processor (SL8De) in IOC Mode SAS/SATA
	33cb  IOP348 I/O Processor (SL8Dx) in IOC Mode SAS/SATA
	3400  5520/5500/X58 I/O Hub to ESI Port
	3401  5520/5500/X58 I/O Hub to ESI Port
	3402  5520/5500/X58 I/O Hub to ESI Port
	3403  5500 I/O Hub to ESI Port
	3404  5520/5500/X58 I/O Hub to ESI Port
	3405  5520/5500/X58 I/O Hub to ESI Port
	3406  5520 I/O Hub to ESI Port
	3407  5520/5500/X58 I/O Hub to ESI Port
	3408  5520/5500/X58 I/O Hub PCI Express Root Port 1
	3409  5520/5500/X58 I/O Hub PCI Express Root Port 2
	340a  5520/5500/X58 I/O Hub PCI Express Root Port 3
	340b  5520/X58 I/O Hub PCI Express Root Port 4
	340c  5520/X58 I/O Hub PCI Express Root Port 5
	340d  5520/X58 I/O Hub PCI Express Root Port 6
	340e  5520/5500/X58 I/O Hub PCI Express Root Port 7
	340f  5520/5500/X58 I/O Hub PCI Express Root Port 8
	3410  7500/5520/5500/X58 I/O Hub PCI Express Root Port 9
	3411  7500/5520/5500/X58 I/O Hub PCI Express Root Port 10
	3418  7500/5520/5500/X58 Physical Layer Port 0
	3419  7500/5520/5500 Physical Layer Port 1
	3420  7500/5520/5500/X58 I/O Hub PCI Express Root Port 0
	3421  7500/5520/5500/X58 I/O Hub PCI Express Root Port 0
	3422  7500/5520/5500/X58 I/O Hub GPIO and Scratch Pad Registers
	3423  7500/5520/5500/X58 I/O Hub Control Status and RAS Registers
	3425  7500/5520/5500/X58 Physical and Link Layer Registers Port 0
	3426  7500/5520/5500/X58 Routing and Protocol Layer Registers Port 0
	3427  7500/5520/5500 Physical and Link Layer Registers Port 1
	3428  7500/5520/5500 Routing & Protocol Layer Register Port 1
	3429  5520/5500/X58 Chipset QuickData Technology Device
	342a  5520/5500/X58 Chipset QuickData Technology Device
	342b  5520/5500/X58 Chipset QuickData Technology Device
	342c  5520/5500/X58 Chipset QuickData Technology Device
	342d  7500/5520/5500/X58 I/O Hub I/OxAPIC Interrupt Controller
	342e  7500/5520/5500/X58 I/O Hub System Management Registers
	342f  7500/5520/5500/X58 Trusted Execution Technology Registers
	3430  5520/5500/X58 Chipset QuickData Technology Device
	3431  5520/5500/X58 Chipset QuickData Technology Device
	3432  5520/5500/X58 Chipset QuickData Technology Device
	3433  5520/5500/X58 Chipset QuickData Technology Device
	3438  7500/5520/5500/X58 I/O Hub Throttle Registers
	3500  6311ESB/6321ESB PCI Express Upstream Port
	3501  6310ESB PCI Express Upstream Port
	3504  6311ESB/6321ESB I/OxAPIC Interrupt Controller
	3505  6310ESB I/OxAPIC Interrupt Controller
	350c  6311ESB/6321ESB PCI Express to PCI-X Bridge
	350d  6310ESB PCI Express to PCI-X Bridge
	3510  6311ESB/6321ESB PCI Express Downstream Port E1
	3511  6310ESB PCI Express Downstream Port E1
	3514  6311ESB/6321ESB PCI Express Downstream Port E2
	3515  6310ESB PCI Express Downstream Port E2
	3518  6311ESB/6321ESB PCI Express Downstream Port E3
	3519  6310ESB PCI Express Downstream Port E3
	3575  82830M/MG/MP Host Bridge
	3576  82830M/MP AGP Bridge
	3577  82830M/MG Integrated Graphics Controller
	3578  82830M/MG/MP Host Bridge
	3580  82852/82855 GM/GME/PM/GMV Processor to I/O Controller
	3581  82852/82855 GM/GME/PM/GMV Processor to AGP Controller
	3582  82852/855GM Integrated Graphics Device
	3584  82852/82855 GM/GME/PM/GMV Processor to I/O Controller
	3585  82852/82855 GM/GME/PM/GMV Processor to I/O Controller
	358c  82854 GMCH
	358e  82854 GMCH Integrated Graphics Device
	3590  E7520 Memory Controller Hub
	3591  E7525/E7520 Error Reporting Registers
	3592  E7320 Memory Controller Hub
	3593  E7320 Error Reporting Registers
	3594  E7520 DMA Controller
	3595  E7525/E7520/E7320 PCI Express Port A
	3596  E7525/E7520/E7320 PCI Express Port A1
	3597  E7525/E7520 PCI Express Port B
	3598  E7520 PCI Express Port B1
	3599  E7520 PCI Express Port C
	359a  E7520 PCI Express Port C1
	359b  E7525/E7520/E7320 Extended Configuration Registers
	359e  E7525 Memory Controller Hub
	35b0  3100 Chipset Memory I/O Controller Hub
	35b1  3100 DRAM Controller Error Reporting Registers
	35b5  3100 Chipset Enhanced DMA Controller
	35b6  3100 Chipset PCI Express Port A
	35b7  3100 Chipset PCI Express Port A1
	35c8  3100 Extended Configuration Test Overflow Registers
	3600  7300 Chipset Memory Controller Hub
	3604  7300 Chipset PCI Express Port 1
	3605  7300 Chipset PCI Express Port 2
	3606  7300 Chipset PCI Express Port 3
	3607  7300 Chipset PCI Express Port 4
	3608  7300 Chipset PCI Express Port 5
	3609  7300 Chipset PCI Express Port 6
	360a  7300 Chipset PCI Express Port 7
	360b  7300 Chipset QuickData Technology Device
	360c  7300 Chipset FSB Registers
	360d  7300 Chipset Snoop Filter Registers
	360e  7300 Chipset Debug and Miscellaneous Registers
	360f  7300 Chipset FBD Branch 0 Registers
	3610  7300 Chipset FBD Branch 1 Registers
	3700  Xeon C5500/C3500 DMI
	3701  Xeon C5500/C3500 DMI
	3702  Xeon C5500/C3500 DMI
	3703  Xeon C5500/C3500 DMI
	3704  Xeon C5500/C3500 DMI
	3705  Xeon C5500/C3500 DMI
	3706  Xeon C5500/C3500 DMI
	3707  Xeon C5500/C3500 DMI
	3708  Xeon C5500/C3500 DMI
	3709  Xeon C5500/C3500 DMI
	370a  Xeon C5500/C3500 DMI
	370b  Xeon C5500/C3500 DMI
	370c  Xeon C5500/C3500 DMI
	370d  Xeon C5500/C3500 DMI
	370e  Xeon C5500/C3500 DMI
	370f  Xeon C5500/C3500 DMI
	3710  Xeon C5500/C3500 CB3 DMA
	3711  Xeon C5500/C3500 CB3 DMA
	3712  Xeon C5500/C3500 CB3 DMA
	3713  Xeon C5500/C3500 CB3 DMA
	3714  Xeon C5500/C3500 CB3 DMA
	3715  Xeon C5500/C3500 CB3 DMA
	3716  Xeon C5500/C3500 CB3 DMA
	3717  Xeon C5500/C3500 CB3 DMA
	3718  Xeon C5500/C3500 CB3 DMA
	3719  Xeon C5500/C3500 CB3 DMA
	371a  Xeon C5500/C3500 QPI Link
	371b  Xeon C5500/C3500 QPI Routing and Protocol
	371d  Xeon C5500/C3500 QPI Routing and Protocol
	3720  Xeon C5500/C3500 PCI Express Root Port 0
	3721  Xeon C5500/C3500 PCI Express Root Port 1
	3722  Xeon C5500/C3500 PCI Express Root Port 2
	3723  Xeon C5500/C3500 PCI Express Root Port 3
	3724  Xeon C5500/C3500 PCI Express Root Port 4
	3725  Xeon C5500/C3500 NTB Primary
	3726  Xeon C5500/C3500 NTB Primary
	3727  Xeon C5500/C3500 NTB Secondary
	3728  Xeon C5500/C3500 Core
	3729  Xeon C5500/C3500 Core
	372a  Xeon C5500/C3500 Core
	372b  Xeon C5500/C3500 Core
	372c  Xeon C5500/C3500 Reserved
	373f  Xeon C5500/C3500 IOxAPIC
	37cd  Ethernet Virtual Function 700 Series
	37ce  Ethernet Connection X722 for 10GbE backplane
	37cf  Ethernet Connection X722 for 10GbE QSFP+
	37d0  Ethernet Connection X722 for 10GbE SFP+
	37d1  Ethernet Connection X722 for 1GbE
	37d2  Ethernet Connection X722 for 10GBASE-T
	37d3  Ethernet Connection X722 for 10GbE SFP+
	37d4  Ethernet Connection X722 for 10GbE QSFP+
	37d9  X722 Hyper-V Virtual Function
	3a00  82801JD/DO (ICH10 Family) 4-port SATA IDE Controller
	3a02  82801JD/DO (ICH10 Family) SATA AHCI Controller
	3a05  82801JD/DO (ICH10 Family) SATA RAID Controller
	3a06  82801JD/DO (ICH10 Family) 2-port SATA IDE Controller
	3a14  82801JDO (ICH10DO) LPC Interface Controller
	3a16  82801JIR (ICH10R) LPC Interface Controller
	3a18  82801JIB (ICH10) LPC Interface Controller
	3a1a  82801JD (ICH10D) LPC Interface Controller
	3a20  82801JI (ICH10 Family) 4 port SATA IDE Controller #1
	3a22  82801JI (ICH10 Family) SATA AHCI Controller
	3a25  82801JIR (ICH10R) SATA RAID Controller
	3a26  82801JI (ICH10 Family) 2 port SATA IDE Controller #2
	3a30  82801JI (ICH10 Family) SMBus Controller
	3a32  82801JI (ICH10 Family) Thermal Subsystem
	3a34  82801JI (ICH10 Family) USB UHCI Controller #1
	3a35  82801JI (ICH10 Family) USB UHCI Controller #2
	3a36  82801JI (ICH10 Family) USB UHCI Controller #3
	3a37  82801JI (ICH10 Family) USB UHCI Controller #4
	3a38  82801JI (ICH10 Family) USB UHCI Controller #5
	3a39  82801JI (ICH10 Family) USB UHCI Controller #6
	3a3a  82801JI (ICH10 Family) USB2 EHCI Controller #1
	3a3c  82801JI (ICH10 Family) USB2 EHCI Controller #2
	3a3e  82801JI (ICH10 Family) HD Audio Controller
	3a40  82801JI (ICH10 Family) PCI Express Root Port 1
	3a42  82801JI (ICH10 Family) PCI Express Port 2
	3a44  82801JI (ICH10 Family) PCI Express Root Port 3
	3a46  82801JI (ICH10 Family) PCI Express Root Port 4
	3a48  82801JI (ICH10 Family) PCI Express Root Port 5
	3a4a  82801JI (ICH10 Family) PCI Express Root Port 6
	3a4c  82801JI (ICH10 Family) Gigabit Ethernet Controller
	3a51  82801JDO (ICH10DO) VECI Controller
	3a55  82801JD/DO (ICH10 Family) Virtual SATA Controller
	3a60  82801JD/DO (ICH10 Family) SMBus Controller
	3a62  82801JD/DO (ICH10 Family) Thermal Subsystem
	3a64  82801JD/DO (ICH10 Family) USB UHCI Controller #1
	3a65  82801JD/DO (ICH10 Family) USB UHCI Controller #2
	3a66  82801JD/DO (ICH10 Family) USB UHCI Controller #3
	3a67  82801JD/DO (ICH10 Family) USB UHCI Controller #4
	3a68  82801JD/DO (ICH10 Family) USB UHCI Controller #5
	3a69  82801JD/DO (ICH10 Family) USB UHCI Controller #6
	3a6a  82801JD/DO (ICH10 Family) USB2 EHCI Controller #1
	3a6c  82801JD/DO (ICH10 Family) USB2 EHCI Controller #2
	3a6e  82801JD/DO (ICH10 Family) HD Audio Controller
	3a70  82801JD/DO (ICH10 Family) PCI Express Port 1
	3a72  82801JD/DO (ICH10 Family) PCI Express Port 2
	3a74  82801JD/DO (ICH10 Family) PCI Express Port 3
	3a76  82801JD/DO (ICH10 Family) PCI Express Port 4
	3a78  82801JD/DO (ICH10 Family) PCI Express Port 5
	3a7a  82801JD/DO (ICH10 Family) PCI Express Port 6
	3a7c  82801JD/DO (ICH10 Family) Gigabit Ethernet Controller
	3b00  5 Series/3400 Series Chipset LPC Interface Controller
	3b01  Mobile 5 Series Chipset LPC Interface Controller
	3b02  P55 Chipset LPC Interface Controller
	3b03  PM55 Chipset LPC Interface Controller
	3b04  5 Series Chipset LPC Interface Controller
	3b05  Mobile 5 Series Chipset LPC Interface Controller
	3b06  H55 Chipset LPC Interface Controller
	3b07  QM57 Chipset LPC Interface Controller
	3b08  H57 Chipset LPC Interface Controller
	3b09  HM55 Chipset LPC Interface Controller
	3b0a  Q57 Chipset LPC Interface Controller
	3b0b  HM57 Chipset LPC Interface Controller
	3b0c  5 Series Chipset LPC Interface Controller
	3b0d  5 Series/3400 Series Chipset LPC Interface Controller
	3b0e  5 Series/3400 Series Chipset LPC Interface Controller
	3b0f  QS57 Chipset LPC Interface Controller
	3b10  5 Series/3400 Series Chipset LPC Interface Controller
	3b11  5 Series/3400 Series Chipset LPC Interface Controller
	3b12  3400 Series Chipset LPC Interface Controller
	3b13  5 Series/3400 Series Chipset LPC Interface Controller
	3b14  3420 Chipset LPC Interface Controller
	3b15  5 Series/3400 Series Chipset LPC Interface Controller
	3b16  3450 Chipset LPC Interface Controller
	3b17  5 Series/3400 Series Chipset LPC Interface Controller
	3b18  5 Series/3400 Series Chipset LPC Interface Controller
	3b19  5 Series/3400 Series Chipset LPC Interface Controller
	3b1a  5 Series/3400 Series Chipset LPC Interface Controller
	3b1b  5 Series/3400 Series Chipset LPC Interface Controller
	3b1c  5 Series/3400 Series Chipset LPC Interface Controller
	3b1d  5 Series/3400 Series Chipset LPC Interface Controller
	3b1e  5 Series/3400 Series Chipset LPC Interface Controller
	3b1f  5 Series/3400 Series Chipset LPC Interface Controller
	3b20  5 Series/3400 Series Chipset 4 port SATA IDE Controller
	3b21  5 Series/3400 Series Chipset 2 port SATA IDE Controller
	3b22  5 Series/3400 Series Chipset 6 port SATA AHCI Controller
	3b23  5 Series/3400 Series Chipset 4 port SATA AHCI Controller
	3b25  5 Series/3400 Series Chipset SATA RAID Controller
	3b26  5 Series/3400 Series Chipset 2 port SATA IDE Controller
	3b28  5 Series/3400 Series Chipset 4 port SATA IDE Controller
	3b29  5 Series/3400 Series Chipset 4 port SATA AHCI Controller
	3b2c  5 Series/3400 Series Chipset SATA RAID Controller
	3b2d  5 Series/3400 Series Chipset 2 port SATA IDE Controller
	3b2e  5 Series/3400 Series Chipset 4 port SATA IDE Controller
	3b2f  5 Series/3400 Series Chipset 6 port SATA AHCI Controller
	3b30  5 Series/3400 Series Chipset SMBus Controller
	3b32  5 Series/3400 Series Chipset Thermal Subsystem
	3b34  5 Series/3400 Series Chipset USB2 Enhanced Host Controller
	3b36  5 Series/3400 Series Chipset USB Universal Host Controller
	3b37  5 Series/3400 Series Chipset USB Universal Host Controller
	3b38  5 Series/3400 Series Chipset USB Universal Host Controller
	3b39  5 Series/3400 Series Chipset USB Universal Host Controller
	3b3a  5 Series/3400 Series Chipset USB Universal Host Controller
	3b3b  5 Series/3400 Series Chipset USB Universal Host Controller
	3b3c  5 Series/3400 Series Chipset USB2 Enhanced Host Controller
	3b3e  5 Series/3400 Series Chipset USB Universal Host Controller
	3b3f  5 Series/3400 Series Chipset USB Universal Host Controller
	3b40  5 Series/3400 Series Chipset USB Universal Host Controller
	3b41  5 Series/3400 Series Chipset LAN Controller
	3b42  5 Series/3400 Series Chipset PCI Express Root Port 1
	3b44  5 Series/3400 Series Chipset PCI Express Root Port 2
	3b46  5 Series/3400 Series Chipset PCI Express Root Port 3
	3b48  5 Series/3400 Series Chipset PCI Express Root Port 4
	3b4a  5 Series/3400 Series Chipset PCI Express Root Port 5
	3b4c  5 Series/3400 Series Chipset PCI Express Root Port 6
	3b4e  5 Series/3400 Series Chipset PCI Express Root Port 7
	3b50  5 Series/3400 Series Chipset PCI Express Root Port 8
	3b53  5 Series/3400 Series Chipset VECI Controller
	3b56  5 Series/3400 Series Chipset High Definition Audio
	3b57  5 Series/3400 Series Chipset High Definition Audio
	3b64  5 Series/3400 Series Chipset HECI Controller
	3b65  5 Series/3400 Series Chipset HECI Controller
	3b66  5 Series/3400 Series Chipset PT IDER Controller
	3b67  5 Series/3400 Series Chipset KT Controller
	3c00  Xeon E5/Core i7 DMI2
	3c01  Xeon E5/Core i7 DMI2 in PCI Express Mode
	3c02  Xeon E5/Core i7 IIO PCI Express Root Port 1a
	3c03  Xeon E5/Core i7 IIO PCI Express Root Port 1b
	3c04  Xeon E5/Core i7 IIO PCI Express Root Port 2a
	3c05  Xeon E5/Core i7 IIO PCI Express Root Port 2b
	3c06  Xeon E5/Core i7 IIO PCI Express Root Port 2c
	3c07  Xeon E5/Core i7 IIO PCI Express Root Port 2d
	3c08  Xeon E5/Core i7 IIO PCI Express Root Port 3a in PCI Express Mode
	3c09  Xeon E5/Core i7 IIO PCI Express Root Port 3b
	3c0a  Xeon E5/Core i7 IIO PCI Express Root Port 3c
	3c0b  Xeon E5/Core i7 IIO PCI Express Root Port 3d
	3c0d  Xeon E5/Core i7 Non-Transparent Bridge
	3c0e  Xeon E5/Core i7 Non-Transparent Bridge
	3c0f  Xeon E5/Core i7 Non-Transparent Bridge
	3c20  Xeon E5/Core i7 DMA Channel 0
	3c21  Xeon E5/Core i7 DMA Channel 1
	3c22  Xeon E5/Core i7 DMA Channel 2
	3c23  Xeon E5/Core i7 DMA Channel 3
	3c24  Xeon E5/Core i7 DMA Channel 4
	3c25  Xeon E5/Core i7 DMA Channel 5
	3c26  Xeon E5/Core i7 DMA Channel 6
	3c27  Xeon E5/Core i7 DMA Channel 7
	3c28  Xeon E5/Core i7 Address Map, VTd_Misc, System Management
	3c2a  Xeon E5/Core i7 Control Status and Global Errors
	3c2c  Xeon E5/Core i7 I/O APIC
	3c2e  Xeon E5/Core i7 DMA
	3c2f  Xeon E5/Core i7 DMA
	3c40  Xeon E5/Core i7 IIO Switch and IRP Performance Monitor
	3c43  Xeon E5/Core i7 Ring to PCI Express Performance Monitor
	3c44  Xeon E5/Core i7 Ring to QuickPath Interconnect Link 0 Performance Monitor
	3c45  Xeon E5/Core i7 Ring to QuickPath Interconnect Link 1 Performance Monitor
	3c46  Xeon E5/Core i7 Processor Home Agent Performance Monitoring
	3c71  Xeon E5/Core i7 Integrated Memory Controller RAS Registers
	3c80  Xeon E5/Core i7 QPI Link 0
	3c83  Xeon E5/Core i7 QPI Link Reut 0
	3c84  Xeon E5/Core i7 QPI Link Reut 0
	3c90  Xeon E5/Core i7 QPI Link 1
	3c93  Xeon E5/Core i7 QPI Link Reut 1
	3c94  Xeon E5/Core i7 QPI Link Reut 1
	3ca0  Xeon E5/Core i7 Processor Home Agent
	3ca8  Xeon E5/Core i7 Integrated Memory Controller Registers
	3caa  Xeon E5/Core i7 Integrated Memory Controller Target Address Decoder 0
	3cab  Xeon E5/Core i7 Integrated Memory Controller Target Address Decoder 1
	3cac  Xeon E5/Core i7 Integrated Memory Controller Target Address Decoder 2
	3cad  Xeon E5/Core i7 Integrated Memory Controller Target Address Decoder 3
	3cae  Xeon E5/Core i7 Integrated Memory Controller Target Address Decoder 4
	3cb0  Xeon E5/Core i7 Integrated Memory Controller Channel 0-3 Thermal Control 0
	3cb1  Xeon E5/Core i7 Integrated Memory Controller Channel 0-3 Thermal Control 1
	3cb2  Xeon E5/Core i7 Integrated Memory Controller ERROR Registers 0
	3cb3  Xeon E5/Core i7 Integrated Memory Controller ERROR Registers 1
	3cb4  Xeon E5/Core i7 Integrated Memory Controller Channel 0-3 Thermal Control 2
	3cb5  Xeon E5/Core i7 Integrated Memory Controller Channel 0-3 Thermal Control 3
	3cb6  Xeon E5/Core i7 Integrated Memory Controller ERROR Registers 2
	3cb7  Xeon E5/Core i7 Integrated Memory Controller ERROR Registers 3
	3cb8  Xeon E5/Core i7 DDRIO
	3cc0  Xeon E5/Core i7 Power Control Unit 0
	3cc1  Xeon E5/Core i7 Power Control Unit 1
	3cc2  Xeon E5/Core i7 Power Control Unit 2
	3cd0  Xeon E5/Core i7 Power Control Unit 3
	3ce0  Xeon E5/Core i7 Interrupt Control Registers
	3ce3  Xeon E5/Core i7 Semaphore and Scratchpad Configuration Registers
	3ce4  Xeon E5/Core i7 R2PCIe
	3ce6  Xeon E5/Core i7 QuickPath Interconnect Agent Ring Registers
	3ce8  Xeon E5/Core i7 Unicast Register 0
	3ce9  Xeon E5/Core i7 Unicast Register 5
	3cea  Xeon E5/Core i7 Unicast Register 1
	3ceb  Xeon E5/Core i7 Unicast Register 6
	3cec  Xeon E5/Core i7 Unicast Register 3
	3ced  Xeon E5/Core i7 Unicast Register 7
	3cee  Xeon E5/Core i7 Unicast Register 4
	3cef  Xeon E5/Core i7 Unicast Register 8
	3cf4  Xeon E5/Core i7 Integrated Memory Controller System Address Decoder 0
	3cf5  Xeon E5/Core i7 Integrated Memory Controller System Address Decoder 1
	3cf6  Xeon E5/Core i7 System Address Decoder
	4000  5400 Chipset Memory Controller Hub
	4001  5400 Chipset Memory Controller Hub
	4003  5400 Chipset Memory Controller Hub
	4021  5400 Chipset PCI Express Port 1
	4022  5400 Chipset PCI Express Port 2
	4023  5400 Chipset PCI Express Port 3
	4024  5400 Chipset PCI Express Port 4
	4025  5400 Chipset PCI Express Port 5
	4026  5400 Chipset PCI Express Port 6
	4027  5400 Chipset PCI Express Port 7
	4028  5400 Chipset PCI Express Port 8
	4029  5400 Chipset PCI Express Port 9
	402d  5400 Chipset IBIST Registers
	402e  5400 Chipset IBIST Registers
	402f  5400 Chipset QuickData Technology Device
	4030  5400 Chipset FSB Registers
	4031  5400 Chipset CE/SF Registers
	4032  5400 Chipset IOxAPIC
	4035  5400 Chipset FBD Registers
	4036  5400 Chipset FBD Registers
	4100  Moorestown Graphics and Video
	4108  Atom Processor E6xx Integrated Graphics Controller
	4109  Atom Processor E6xx Integrated Graphics Controller
	410a  Atom Processor E6xx Integrated Graphics Controller
	410b  Atom Processor E6xx Integrated Graphics Controller
	410c  Atom Processor E6xx Integrated Graphics Controller
	410d  Atom Processor E6xx Integrated Graphics Controller
	410e  Atom Processor E6xx Integrated Graphics Controller
	410f  Atom Processor E6xx Integrated Graphics Controller
	4114  Atom Processor E6xx PCI Host Bridge #1
	4115  Atom Processor E6xx PCI Host Bridge #2
	4116  Atom Processor E6xx PCI Host Bridge #3
	4117  Atom Processor E6xx PCI Host Bridge #4
	4220  PRO/Wireless 2200BG [Calexico2] Network Connection
	4222  PRO/Wireless 3945ABG [Golan] Network Connection
	4223  PRO/Wireless 2915ABG [Calexico2] Network Connection
	4224  PRO/Wireless 2915ABG [Calexico2] Network Connection
	4227  PRO/Wireless 3945ABG [Golan] Network Connection
	4229  PRO/Wireless 4965 AG or AGN [Kedron] Network Connection
	422b  Centrino Ultimate-N 6300
	422c  Centrino Advanced-N 6200
	4230  PRO/Wireless 4965 AG or AGN [Kedron] Network Connection
	4232  WiFi Link 5100
	4235  Ultimate N WiFi Link 5300
	4236  Ultimate N WiFi Link 5300
	4237  PRO/Wireless 5100 AGN [Shiloh] Network Connection
	4238  Centrino Ultimate-N 6300
	4239  Centrino Advanced-N 6200
	423a  PRO/Wireless 5350 AGN [Echo Peak] Network Connection
	423b  PRO/Wireless 5350 AGN [Echo Peak] Network Connection
	423c  WiMAX/WiFi Link 5150
	423d  WiMAX/WiFi Link 5150
	444e  Turbo Memory Controller
	5001  LE80578
	5002  LE80578 Graphics Processor Unit
	5009  LE80578 Video Display Controller
	500d  LE80578 Expansion Bus
	500e  LE80578 UART Controller
	500f  LE80578 General Purpose IO
	5010  LE80578 I2C Controller
	5012  LE80578 Serial Peripheral Interface Bus
	5020  EP80579 Memory Controller Hub
	5021  EP80579 DRAM Error Reporting Registers
	5023  EP80579 EDMA Controller
	5024  EP80579 PCI Express Port PEA0
	5025  EP80579 PCI Express Port PEA1
	5028  EP80579 S-ATA IDE
	5029  EP80579 S-ATA AHCI
	502a  EP80579 S-ATA Reserved
	502b  EP80579 S-ATA Reserved
	502c  EP80579 Integrated Processor ASU
	502d  EP80579 Integrated Processor with QuickAssist ASU
	502e  EP80579 Reserved
	502f  EP80579 Reserved
	5030  EP80579 Reserved
	5031  EP80579 LPC Bus
	5032  EP80579 SMBus Controller
	5033  EP80579 USB 1.1 Controller
	5035  EP80579 USB 2.0 Controller
	5037  EP80579 PCI-PCI Bridge (transparent mode)
	5039  EP80579 Controller Area Network (CAN) interface #1
	503a  EP80579 Controller Area Network (CAN) interface #2
	503b  EP80579 Synchronous Serial Port (SPP)
	503c  EP80579 IEEE 1588 Hardware Assist
	503d  EP80579 Local Expansion Bus
	503e  EP80579 Global Control Unit (GCU)
	503f  EP80579 Reserved
	5040  EP80579 Integrated Processor Gigabit Ethernet MAC
	5041  EP80579 Integrated Processor with QuickAssist Gigabit Ethernet MAC
	5042  EP80579 Reserved
	5043  EP80579 Reserved
	5044  EP80579 Integrated Processor Gigabit Ethernet MAC
	5045  EP80579 Integrated Processor with QuickAssist Gigabit Ethernet MAC
	5046  EP80579 Reserved
	5047  EP80579 Reserved
	5048  EP80579 Integrated Processor Gigabit Ethernet MAC
	5049  EP80579 Integrated Processor with QuickAssist Gigabit Ethernet MAC
	504a  EP80579 Reserved
	504b  EP80579 Reserved
	504c  EP80579 Integrated Processor with QuickAssist TDM
	5200  EtherExpress PRO/100 Intelligent Server PCI Bridge
	5201  EtherExpress PRO/100 Intelligent Server Fast Ethernet Controller
	530d  80310 (IOP) IO Processor
	5845  QEMU NVM Express Controller
	5902  HD Graphics 610
	5904  Xeon E3-1200 v6/7th Gen Core Processor Host Bridge/DRAM Registers
	590f  Xeon E3-1200 v6/7th Gen Core Processor Host Bridge/DRAM Registers
	5910  Xeon E3-1200 v6/7th Gen Core Processor Host Bridge/DRAM Registers
	5912  HD Graphics 630
	5916  HD Graphics 620
	591d  HD Graphics P630
	591f  Intel Kaby Lake Host Bridge
	5a84  Celeron N3350/Pentium N4200/Atom E3900 Series Integrated Graphics Controller
	5a88  Celeron N3350/Pentium N4200/Atom E3900 Series Imaging Unit
	5a98  Celeron N3350/Pentium N4200/Atom E3900 Series Audio Cluster
	5a9a  Celeron N3350/Pentium N4200/Atom E3900 Series Trusted Execution Engine
	5aa2  Celeron N3350/Pentium N4200/Atom E3900 Series Integrated Sensor Hub
	5aa8  Celeron N3350/Pentium N4200/Atom E3900 Series USB xHCI
	5aac  Celeron N3350/Pentium N4200/Atom E3900 Series I2C Controller #1
	5aae  Celeron N3350/Pentium N4200/Atom E3900 Series I2C Controller #2
	5ab0  Celeron N3350/Pentium N4200/Atom E3900 Series I2C Controller #3
	5ab2  Celeron N3350/Pentium N4200/Atom E3900 Series I2C Controller #4
	5ab4  Celeron N3350/Pentium N4200/Atom E3900 Series I2C Controller #5
	5ab6  Celeron N3350/Pentium N4200/Atom E3900 Series I2C Controller #6
	5ab8  Celeron N3350/Pentium N4200/Atom E3900 Series I2C Controller #7
	5aba  Celeron N3350/Pentium N4200/Atom E3900 Series I2C Controller #8
	5abc  Celeron N3350/Pentium N4200/Atom E3900 Series HSUART Controller #1
	5abe  Celeron N3350/Pentium N4200/Atom E3900 Series HSUART Controller #2
	5ac0  Celeron N3350/Pentium N4200/Atom E3900 Series HSUART Controller #3
	5ac2  Celeron N3350/Pentium N4200/Atom E3900 Series SPI Controller #1
	5ac4  Celeron N3350/Pentium N4200/Atom E3900 Series SPI Controller #2
	5ac6  Celeron N3350/Pentium N4200/Atom E3900 Series SPI Controller #3
	5ac8  Celeron N3350/Pentium N4200/Atom E3900 Series PWM Pin Controller
	5aca  Celeron N3350/Pentium N4200/Atom E3900 Series SDXC/MMC Host Controller
	5acc  Celeron N3350/Pentium N4200/Atom E3900 Series eMMC Controller
	5ad0  Celeron N3350/Pentium N4200/Atom E3900 Series SDIO Controller
	5ad4  Celeron N3350/Pentium N4200/Atom E3900 Series SMBus Controller
	5ad6  Celeron N3350/Pentium N4200/Atom E3900 Series PCI Express Port B #1
	5ad7  Celeron N3350/Pentium N4200/Atom E3900 Series PCI Express Port B #2
	5ad8  Celeron N3350/Pentium N4200/Atom E3900 Series PCI Express Port A #1
	5ad9  Celeron N3350/Pentium N4200/Atom E3900 Series PCI Express Port A #2
	5ada  Celeron N3350/Pentium N4200/Atom E3900 Series PCI Express Port A #3
	5adb  Celeron N3350/Pentium N4200/Atom E3900 Series PCI Express Port A #4
	5ae3  Celeron N3350/Pentium N4200/Atom E3900 Series SATA AHCI Controller
	5ae8  Celeron N3350/Pentium N4200/Atom E3900 Series Low Pin Count Interface
	5aee  Celeron N3350/Pentium N4200/Atom E3900 Series HSUART Controller #4
	5af0  Celeron N3350/Pentium N4200/Atom E3900 Series Host Bridge
	65c0  5100 Chipset Memory Controller Hub
	65e2  5100 Chipset PCI Express x4 Port 2
	65e3  5100 Chipset PCI Express x4 Port 3
	65e4  5100 Chipset PCI Express x4 Port 4
	65e5  5100 Chipset PCI Express x4 Port 5
	65e6  5100 Chipset PCI Express x4 Port 6
	65e7  5100 Chipset PCI Express x4 Port 7
	65f0  5100 Chipset FSB Registers
	65f1  5100 Chipset Reserved Registers
	65f3  5100 Chipset Reserved Registers
	65f5  5100 Chipset DDR Channel 0 Registers
	65f6  5100 Chipset DDR Channel 1 Registers
	65f7  5100 Chipset PCI Express x8 Port 2-3
	65f8  5100 Chipset PCI Express x8 Port 4-5
	65f9  5100 Chipset PCI Express x8 Port 6-7
	65fa  5100 Chipset PCI Express x16 Port 4-7
	65ff  5100 Chipset DMA Engine
	6f00  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D DMI2
	6f01  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D PCI Express Root Port 0
	6f02  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D PCI Express Root Port 1
	6f03  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D PCI Express Root Port 1
	6f04  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D PCI Express Root Port 2
	6f05  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D PCI Express Root Port 2
	6f06  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D PCI Express Root Port 2
	6f07  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D PCI Express Root Port 2
	6f08  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D PCI Express Root Port 3
	6f09  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D PCI Express Root Port 3
	6f0a  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D PCI Express Root Port 3
	6f0b  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D PCI Express Root Port 3
	6f10  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D IIO Debug
	6f11  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D IIO Debug
	6f12  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D IIO Debug
	6f13  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D IIO Debug
	6f14  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D IIO Debug
	6f15  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D IIO Debug
	6f16  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D IIO Debug
	6f17  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D IIO Debug
	6f18  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D IIO Debug
	6f19  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D IIO Debug
	6f1a  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D IIO Debug
	6f1b  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D IIO Debug
	6f1c  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D IIO Debug
	6f1d  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D R2PCIe Agent
	6f1e  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Ubox
	6f1f  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Ubox
	6f20  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Crystal Beach DMA Channel 0
	6f21  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Crystal Beach DMA Channel 1
	6f22  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Crystal Beach DMA Channel 2
	6f23  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Crystal Beach DMA Channel 3
	6f24  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Crystal Beach DMA Channel 4
	6f25  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Crystal Beach DMA Channel 5
	6f26  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Crystal Beach DMA Channel 6
	6f27  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Crystal Beach DMA Channel 7
	6f28  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Map/VTd_Misc/System Management
	6f29  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D IIO Hot Plug
	6f2a  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D IIO RAS/Control Status/Global Errors
	6f2c  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D I/O APIC
	6f30  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Home Agent 0
	6f32  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QPI Link 0
	6f33  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QPI Link 1
	6f34  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D R2PCIe Agent
	6f36  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D R3 QPI Link 0/1
	6f37  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D R3 QPI Link 0/1
	6f38  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Home Agent 1
	6f39  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D IO Performance Monitoring
	6f3a  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QPI Link 2
	6f3e  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D R3 QPI Link 2
	6f3f  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D R3 QPI Link 2
	6f40  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QPI Link 2
	6f41  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D R3 QPI Link 2
	6f43  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QPI Link 2
	6f45  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QPI Link 2 Debug
	6f46  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QPI Link 2 Debug
	6f47  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QPI Link 2 Debug
	6f50  Xeon Processor D Family QuickData Technology Register DMA Channel 0
	6f51  Xeon Processor D Family QuickData Technology Register DMA Channel 1
	6f52  Xeon Processor D Family QuickData Technology Register DMA Channel 2
	6f53  Xeon Processor D Family QuickData Technology Register DMA Channel 3
	6f60  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Home Agent 1
	6f68  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Target Address/Thermal/RAS
	6f6a  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Channel Target Address Decoder
	6f6b  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Channel Target Address Decoder
	6f6c  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Channel Target Address Decoder
	6f6d  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Channel Target Address Decoder
	6f6e  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D DDRIO Channel 2/3 Broadcast
	6f6f  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D DDRIO Global Broadcast
	6f70  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Home Agent 0 Debug
	6f71  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 0 - Target Address/Thermal/RAS
	6f76  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D R3 QPI Link Debug
	6f78  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Home Agent 1 Debug
	6f79  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Target Address/Thermal/RAS
	6f7d  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Ubox
	6f7e  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D R3 QPI Link Debug
	6f80  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QPI Link 0
	6f81  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D R3 QPI Link 0/1
	6f83  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QPI Link 0
	6f85  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QPI Link 0 Debug
	6f86  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QPI Link 0 Debug
	6f87  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QPI Link 0 Debug
	6f88  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6f8a  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6f90  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QPI Link 1
	6f93  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QPI Link 1
	6f95  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QPI Link 1 Debug
	6f96  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D QPI Link 1 Debug
	6f98  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6f99  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6f9a  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6f9c  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6fa0  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Home Agent 0
	6fa8  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 0 - Target Address/Thermal/RAS
	6faa  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 0 - Channel Target Address Decoder
	6fab  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 0 - Channel Target Address Decoder
	6fac  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 0 - Channel Target Address Decoder
	6fad  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 0 - Channel Target Address Decoder
	6fae  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D DDRIO Channel 0/1 Broadcast
	6faf  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D DDRIO Global Broadcast
	6fb0  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 0 - Channel 0 Thermal Control
	6fb1  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 0 - Channel 1 Thermal Control
	6fb2  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 0 - Channel 0 Error
	6fb3  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 0 - Channel 1 Error
	6fb4  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 0 - Channel 2 Thermal Control
	6fb5  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 0 - Channel 3 Thermal Control
	6fb6  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 0 - Channel 2 Error
	6fb7  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 0 - Channel 3 Error
	6fb8  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D DDRIO Channel 2/3 Interface
	6fb9  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D DDRIO Channel 2/3 Interface
	6fba  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D DDRIO Channel 2/3 Interface
	6fbb  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D DDRIO Channel 2/3 Interface
	6fbc  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D DDRIO Channel 0/1 Interface
	6fbd  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D DDRIO Channel 0/1 Interface
	6fbe  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D DDRIO Channel 0/1 Interface
	6fbf  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D DDRIO Channel 0/1 Interface
	6fc0  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6fc1  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6fc2  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6fc3  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6fc4  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6fc5  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6fc6  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6fc7  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6fc8  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6fc9  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6fca  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6fcb  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6fcc  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6fcd  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6fce  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6fcf  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Power Control Unit
	6fd0  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 1 - Channel 0 Thermal Control
	6fd1  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 1 - Channel 1 Thermal Control
	6fd2  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 1 - Channel 0 Error
	6fd3  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 1 - Channel 1 Error
	6fd4  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 1 - Channel 2 Thermal Control
	6fd5  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 1 - Channel 3 Thermal Control
	6fd6  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 1 - Channel 2 Error
	6fd7  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Memory Controller 1 - Channel 3 Error
	6fe0  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6fe1  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6fe2  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6fe3  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6fe4  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6fe5  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6fe6  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6fe7  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6fe8  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6fe9  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6fea  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6feb  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6fec  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6fed  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6fee  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6fef  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6ff0  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6ff1  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6ff8  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6ff9  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6ffa  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6ffb  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6ffc  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6ffd  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	6ffe  Xeon E7 v4/Xeon E5 v4/Xeon E3 v4/Xeon D Caching Agent
	7000  82371SB PIIX3 ISA [Natoma/Triton II]
	7010  82371SB PIIX3 IDE [Natoma/Triton II]
	7020  82371SB PIIX3 USB [Natoma/Triton II]
	7030  430VX - 82437VX TVX [Triton VX]
	7050  Intercast Video Capture Card
	7051  PB 642365-003 (Business Video Conferencing Card)
	7100  430TX - 82439TX MTXC
	7110  82371AB/EB/MB PIIX4 ISA
	7111  82371AB/EB/MB PIIX4 IDE
	7112  82371AB/EB/MB PIIX4 USB
	7113  82371AB/EB/MB PIIX4 ACPI
	7120  82810 GMCH (Graphics Memory Controller Hub)
	7121  82810 (CGC) Chipset Graphics Controller
	7122  82810 DC-100 (GMCH) Graphics Memory Controller Hub
	7123  82810 DC-100 (CGC) Chipset Graphics Controller
	7124  82810E DC-133 (GMCH) Graphics Memory Controller Hub
	7125  82810E DC-133 (CGC) Chipset Graphics Controller
	7126  82810 DC-133 System and Graphics Controller
	7128  82810-M DC-100 System and Graphics Controller
	712a  82810-M DC-133 System and Graphics Controller
	7180  440LX/EX - 82443LX/EX Host bridge
	7181  440LX/EX - 82443LX/EX AGP bridge
	7190  440BX/ZX/DX - 82443BX/ZX/DX Host bridge
	7191  440BX/ZX/DX - 82443BX/ZX/DX AGP bridge
	7192  440BX/ZX/DX - 82443BX/ZX/DX Host bridge (AGP disabled)
	7194  82440MX Host Bridge
	7195  82440MX AC'97 Audio Controller
	7196  82440MX AC'97 Modem Controller
	7198  82440MX ISA Bridge
	7199  82440MX EIDE Controller
	719a  82440MX USB Universal Host Controller
	719b  82440MX Power Management Controller
	71a0  440GX - 82443GX Host bridge
	71a1  440GX - 82443GX AGP bridge
	71a2  440GX - 82443GX Host bridge (AGP disabled)
	7600  82372FB PIIX5 ISA
	7601  82372FB PIIX5 IDE
	7602  82372FB PIIX5 USB
	7603  82372FB PIIX5 SMBus
	7800  82740 (i740) AGP Graphics Accelerator
	8002  Trusted Execution Technology Registers
	8003  Trusted Execution Technology Registers
	8100  System Controller Hub (SCH Poulsbo)
	8108  System Controller Hub (SCH Poulsbo) Graphics Controller
	8110  System Controller Hub (SCH Poulsbo) PCI Express Port 1
	8112  System Controller Hub (SCH Poulsbo) PCI Express Port 2
	8114  System Controller Hub (SCH Poulsbo) USB UHCI #1
	8115  System Controller Hub (SCH Poulsbo) USB UHCI #2
	8116  System Controller Hub (SCH Poulsbo) USB UHCI #3
	8117  System Controller Hub (SCH Poulsbo) USB EHCI #1
	8118  System Controller Hub (SCH Poulsbo) USB Client Controller
	8119  System Controller Hub (SCH Poulsbo) LPC Bridge
	811a  System Controller Hub (SCH Poulsbo) IDE Controller
	811b  System Controller Hub (SCH Poulsbo) HD Audio Controller
	811c  System Controller Hub (SCH Poulsbo) SDIO Controller #1
	811d  System Controller Hub (SCH Poulsbo) SDIO Controller #2
	811e  System Controller Hub (SCH Poulsbo) SDIO Controller #3
	8180  Atom Processor E6xx PCI Express Port 3
	8181  Atom Processor E6xx PCI Express Port 4
	8182  Atom Processor E6xx Integrated Graphics Controller
	8183  Atom Processor E6xx Configuration Unit
	8184  Atom Processor E6xx PCI Express Port 1
	8185  Atom Processor E6xx PCI Express Port 2
	8186  Atom Processor E6xx LPC Bridge
	84c4  450KX/GX [Orion] - 82454KX/GX PCI bridge
	84c5  450KX/GX [Orion] - 82453KX/GX Memory controller
	84ca  450NX - 82451NX Memory & I/O Controller
	84cb  450NX - 82454NX/84460GX PCI Expander Bridge
	84e0  460GX - 84460GX System Address Controller (SAC)
	84e1  460GX - 84460GX System Data Controller (SDC)
	84e2  460GX - 84460GX AGP Bridge (GXB function 2)
	84e3  460GX - 84460GX Memory Address Controller (MAC)
	84e4  460GX - 84460GX Memory Data Controller (MDC)
	84e6  460GX - 82466GX Wide and fast PCI eXpander Bridge (WXB)
	84ea  460GX - 84460GX AGP Bridge (GXB function 1)
	8500  IXP4XX Network Processor (IXP420/421/422/425/IXC1100)
	8800  Platform Controller Hub EG20T PCI Express Port
	8801  Platform Controller Hub EG20T Packet Hub
	8802  Platform Controller Hub EG20T Gigabit Ethernet Controller
	8803  Platform Controller Hub EG20T General Purpose IO Controller
	8804  Platform Controller Hub EG20T USB OHCI Controller #4
	8805  Platform Controller Hub EG20T USB OHCI Controller #5
	8806  Platform Controller Hub EG20T USB OHCI Controller #6
	8807  Platform Controller Hub EG20T USB2 EHCI Controller #2
	8808  Platform Controller Hub EG20T USB Client Controller
	8809  Platform Controller Hub EG20T SDIO Controller #1
	880a  Platform Controller Hub EG20T SDIO Controller #2
	880b  Platform Controller Hub EG20T SATA AHCI Controller
	880c  Platform Controller Hub EG20T USB OHCI Controller #1
	880d  Platform Controller Hub EG20T USB OHCI Controller #2
	880e  Platform Controller Hub EG20T USB OHCI Controller #3
	880f  Platform Controller Hub EG20T USB2 EHCI Controller #1
	8810  Platform Controller Hub EG20T DMA Controller #1
	8811  Platform Controller Hub EG20T UART Controller 0
	8812  Platform Controller Hub EG20T UART Controller 1
	8813  Platform Controller Hub EG20T UART Controller 2
	8814  Platform Controller Hub EG20T UART Controller 3
	8815  Platform Controller Hub EG20T DMA Controller #2
	8816  Platform Controller Hub EG20T Serial Peripheral Interface Bus
	8817  Platform Controller Hub EG20T I2C Controller
	8818  Platform Controller Hub EG20T Controller Area Network (CAN) Controller
	8819  Platform Controller Hub EG20T IEEE 1588 Hardware Assist
	8c00  8 Series/C220 Series Chipset Family 4-port SATA Controller 1 [IDE mode]
	8c01  8 Series Chipset Family 4-port SATA Controller 1 [IDE mode] - Mobile
	8c02  8 Series/C220 Series Chipset Family 6-port SATA Controller 1 [AHCI mode]
	8c03  8 Series/C220 Series Chipset Family 6-port SATA Controller 1 [AHCI mode]
	8c04  8 Series/C220 Series Chipset Family SATA Controller 1 [RAID mode]
	8c05  8 Series/C220 Series Chipset Family SATA Controller 1 [RAID mode]
	8c06  8 Series/C220 Series Chipset Family SATA Controller 1 [RAID mode]
	8c07  8 Series/C220 Series Chipset Family SATA Controller 1 [RAID mode]
	8c08  8 Series/C220 Series Chipset Family 2-port SATA Controller 2 [IDE mode]
	8c09  8 Series/C220 Series Chipset Family 2-port SATA Controller 2 [IDE mode]
	8c0e  8 Series/C220 Series Chipset Family SATA Controller 1 [RAID mode]
	8c0f  8 Series/C220 Series Chipset Family SATA Controller 1 [RAID mode]
	8c10  8 Series/C220 Series Chipset Family PCI Express Root Port #1
	8c11  8 Series/C220 Series Chipset Family PCI Express Root Port #1
	8c12  8 Series/C220 Series Chipset Family PCI Express Root Port #2
	8c13  8 Series/C220 Series Chipset Family PCI Express Root Port #2
	8c14  8 Series/C220 Series Chipset Family PCI Express Root Port #3
	8c15  8 Series/C220 Series Chipset Family PCI Express Root Port #3
	8c16  8 Series/C220 Series Chipset Family PCI Express Root Port #4
	8c17  8 Series/C220 Series Chipset Family PCI Express Root Port #4
	8c18  8 Series/C220 Series Chipset Family PCI Express Root Port #5
	8c19  8 Series/C220 Series Chipset Family PCI Express Root Port #5
	8c1a  8 Series/C220 Series Chipset Family PCI Express Root Port #6
	8c1b  8 Series/C220 Series Chipset Family PCI Express Root Port #6
	8c1c  8 Series/C220 Series Chipset Family PCI Express Root Port #7
	8c1d  8 Series/C220 Series Chipset Family PCI Express Root Port #7
	8c1e  8 Series/C220 Series Chipset Family PCI Express Root Port #8
	8c1f  8 Series/C220 Series Chipset Family PCI Express Root Port #8
	8c20  8 Series/C220 Series Chipset High Definition Audio Controller
	8c21  8 Series/C220 Series Chipset High Definition Audio Controller
	8c22  8 Series/C220 Series Chipset Family SMBus Controller
	8c23  8 Series Chipset Family CHAP Counters
	8c24  8 Series Chipset Family Thermal Management Controller
	8c26  8 Series/C220 Series Chipset Family USB EHCI #1
	8c2d  8 Series/C220 Series Chipset Family USB EHCI #2
	8c31  8 Series/C220 Series Chipset Family USB xHCI
	8c33  8 Series/C220 Series Chipset Family LAN Controller
	8c34  8 Series/C220 Series Chipset Family NAND Controller
	8c3a  8 Series/C220 Series Chipset Family MEI Controller #1
	8c3b  8 Series/C220 Series Chipset Family MEI Controller #2
	8c3c  8 Series/C220 Series Chipset Family IDE-r Controller
	8c3d  8 Series/C220 Series Chipset Family KT Controller
	8c40  8 Series/C220 Series Chipset Family LPC Controller
	8c41  8 Series Chipset Family Mobile Super SKU LPC Controller
	8c42  8 Series/C220 Series Chipset Family Desktop Super SKU LPC Controller
	8c43  8 Series/C220 Series Chipset Family LPC Controller
	8c44  Z87 Express LPC Controller
	8c45  8 Series/C220 Series Chipset Family LPC Controller
	8c46  Z85 Express LPC Controller
	8c47  8 Series/C220 Series Chipset Family LPC Controller
	8c48  8 Series/C220 Series Chipset Family LPC Controller
	8c49  HM86 Express LPC Controller
	8c4a  H87 Express LPC Controller
	8c4b  HM87 Express LPC Controller
	8c4c  Q85 Express LPC Controller
	8c4d  8 Series/C220 Series Chipset Family LPC Controller
	8c4e  Q87 Express LPC Controller
	8c4f  QM87 Express LPC Controller
	8c50  B85 Express LPC Controller
	8c51  8 Series/C220 Series Chipset Family LPC Controller
	8c52  C222 Series Chipset Family Server Essential SKU LPC Controller
	8c53  8 Series/C220 Series Chipset Family LPC Controller
	8c54  C224 Series Chipset Family Server Standard SKU LPC Controller
	8c55  8 Series/C220 Series Chipset Family LPC Controller
	8c56  C226 Series Chipset Family Server Advanced SKU LPC Controller
	8c57  8 Series/C220 Series Chipset Family LPC Controller
	8c58  8 Series/C220 Series Chipset Family WS SKU LPC Controller
	8c59  8 Series/C220 Series Chipset Family LPC Controller
	8c5a  8 Series/C220 Series Chipset Family LPC Controller
	8c5b  8 Series/C220 Series Chipset Family LPC Controller
	8c5c  C220 Series Chipset Family H81 Express LPC Controller
	8c5d  8 Series/C220 Series Chipset Family LPC Controller
	8c5e  8 Series/C220 Series Chipset Family LPC Controller
	8c5f  8 Series/C220 Series Chipset Family LPC Controller
	8c80  9 Series Chipset Family SATA Controller [IDE Mode]
	8c81  9 Series Chipset Family SATA Controller [IDE Mode]
	8c82  9 Series Chipset Family SATA Controller [AHCI Mode]
	8c83  9 Series Chipset Family SATA Controller [AHCI Mode]
	8c84  9 Series Chipset Family SATA Controller [RAID Mode]
	8c85  9 Series Chipset Family SATA Controller [RAID Mode]
	8c86  9 Series Chipset Family SATA Controller [RAID Mode]
	8c87  9 Series Chipset Family SATA Controller [RAID Mode]
	8c88  9 Series Chipset Family SATA Controller [IDE Mode]
	8c89  9 Series Chipset Family SATA Controller [IDE Mode]
	8c8e  9 Series Chipset Family SATA Controller [RAID Mode]
	8c8f  9 Series Chipset Family SATA Controller [RAID Mode]
	8c90  9 Series Chipset Family PCI Express Root Port 1
	8c92  9 Series Chipset Family PCI Express Root Port 2
	8c94  9 Series Chipset Family PCI Express Root Port 3
	8c96  9 Series Chipset Family PCI Express Root Port 4
	8c98  9 Series Chipset Family PCI Express Root Port 5
	8c9a  9 Series Chipset Family PCI Express Root Port 6
	8c9c  9 Series Chipset Family PCI Express Root Port 7
	8c9e  9 Series Chipset Family PCI Express Root Port 8
	8ca0  9 Series Chipset Family HD Audio Controller
	8ca2  9 Series Chipset Family SMBus Controller
	8ca4  9 Series Chipset Family Thermal Controller
	8ca6  9 Series Chipset Family USB EHCI Controller #1
	8cad  9 Series Chipset Family USB EHCI Controller #2
	8cb1  9 Series Chipset Family USB xHCI Controller
	8cb3  9 Series Chipset Family LAN Controller
	8cba  9 Series Chipset Family ME Interface #1
	8cbb  9 Series Chipset Family ME Interface #2
	8cbc  9 Series Chipset Family IDE-R Controller
	8cbd  9 Series Chipset Family KT Controller
	8cc1  9 Series Chipset Family LPC Controller
	8cc2  9 Series Chipset Family LPC Controller
	8cc3  9 Series Chipset Family HM97 LPC Controller
	8cc4  9 Series Chipset Family Z97 LPC Controller
	8cc6  9 Series Chipset Family H97 Controller
	8d00  C610/X99 series chipset 4-port SATA Controller [IDE mode]
	8d02  C610/X99 series chipset 6-Port SATA Controller [AHCI mode]
	8d04  C610/X99 series chipset SATA Controller [RAID mode]
	8d06  C610/X99 series chipset SATA Controller [RAID mode]
	8d08  C610/X99 series chipset 2-port SATA Controller [IDE mode]
	8d0e  C610/X99 series chipset SATA Controller [RAID mode]
	8d10  C610/X99 series chipset PCI Express Root Port #1
	8d11  C610/X99 series chipset PCI Express Root Port #1
	8d12  C610/X99 series chipset PCI Express Root Port #2
	8d13  C610/X99 series chipset PCI Express Root Port #2
	8d14  C610/X99 series chipset PCI Express Root Port #3
	8d15  C610/X99 series chipset PCI Express Root Port #3
	8d16  C610/X99 series chipset PCI Express Root Port #4
	8d17  C610/X99 series chipset PCI Express Root Port #4
	8d18  C610/X99 series chipset PCI Express Root Port #5
	8d19  C610/X99 series chipset PCI Express Root Port #5
	8d1a  C610/X99 series chipset PCI Express Root Port #6
	8d1b  C610/X99 series chipset PCI Express Root Port #6
	8d1c  C610/X99 series chipset PCI Express Root Port #7
	8d1d  C610/X99 series chipset PCI Express Root Port #7
	8d1e  C610/X99 series chipset PCI Express Root Port #8
	8d1f  C610/X99 series chipset PCI Express Root Port #8
	8d20  C610/X99 series chipset HD Audio Controller
	8d21  C610/X99 series chipset HD Audio Controller
	8d22  C610/X99 series chipset SMBus Controller
	8d24  C610/X99 series chipset Thermal Subsystem
	8d26  C610/X99 series chipset USB Enhanced Host Controller #1
	8d2d  C610/X99 series chipset USB Enhanced Host Controller #2
	8d31  C610/X99 series chipset USB xHCI Host Controller
	8d33  C610/X99 series chipset LAN Controller
	8d34  C610/X99 series chipset NAND Controller
	8d3a  C610/X99 series chipset MEI Controller #1
	8d3b  C610/X99 series chipset MEI Controller #2
	8d3c  C610/X99 series chipset IDE-r Controller
	8d3d  C610/X99 series chipset KT Controller
	8d40  C610/X99 series chipset LPC Controller
	8d41  C610/X99 series chipset LPC Controller
	8d42  C610/X99 series chipset LPC Controller
	8d43  C610/X99 series chipset LPC Controller
	8d44  C610/X99 series chipset LPC Controller
	8d45  C610/X99 series chipset LPC Controller
	8d46  C610/X99 series chipset LPC Controller
	8d47  C610/X99 series chipset LPC Controller
	8d48  C610/X99 series chipset LPC Controller
	8d49  C610/X99 series chipset LPC Controller
	8d4a  C610/X99 series chipset LPC Controller
	8d4b  C610/X99 series chipset LPC Controller
	8d4c  C610/X99 series chipset LPC Controller
	8d4d  C610/X99 series chipset LPC Controller
	8d4e  C610/X99 series chipset LPC Controller
	8d4f  C610/X99 series chipset LPC Controller
	8d60  C610/X99 series chipset sSATA Controller [IDE mode]
	8d62  C610/X99 series chipset sSATA Controller [AHCI mode]
	8d64  C610/X99 series chipset sSATA Controller [RAID mode]
	8d66  C610/X99 series chipset sSATA Controller [RAID mode]
	8d68  C610/X99 series chipset sSATA Controller [IDE mode]
	8d6e  C610/X99 series chipset sSATA Controller [RAID mode]
	8d7c  C610/X99 series chipset SPSR
	8d7d  C610/X99 series chipset MS SMBus 0
	8d7e  C610/X99 series chipset MS SMBus 1
	8d7f  C610/X99 series chipset MS SMBus 2
	9000  IXP2000 Family Network Processor
	9001  IXP2400 Network Processor
	9002  IXP2300 Network Processor
	9004  IXP2800 Network Processor
	9621  Integrated RAID
	9622  Integrated RAID
	9641  Integrated RAID
	96a1  Integrated RAID
	9c00  8 Series SATA Controller 1 [IDE mode]
	9c01  8 Series SATA Controller 1 [IDE mode]
	9c02  8 Series SATA Controller 1 [AHCI mode]
	9c03  8 Series SATA Controller 1 [AHCI mode]
	9c04  8 Series SATA Controller 1 [RAID mode]
	9c05  8 Series SATA Controller 1 [RAID mode]
	9c06  8 Series SATA Controller 1 [RAID mode]
	9c07  8 Series SATA Controller 1 [RAID mode]
	9c08  8 Series SATA Controller 2 [IDE mode]
	9c09  8 Series SATA Controller 2 [IDE mode]
	9c0a  8 Series SATA Controller [Reserved]
	9c0b  8 Series SATA Controller [Reserved]
	9c0c  8 Series SATA Controller [Reserved]
	9c0d  8 Series SATA Controller [Reserved]
	9c0e  8 Series SATA Controller 1 [RAID mode]
	9c0f  8 Series SATA Controller 1 [RAID mode]
	9c10  8 Series PCI Express Root Port 1
	9c11  8 Series PCI Express Root Port 1
	9c12  8 Series PCI Express Root Port 2
	9c13  8 Series PCI Express Root Port 2
	9c14  8 Series PCI Express Root Port 3
	9c15  8 Series PCI Express Root Port 3
	9c16  8 Series PCI Express Root Port 4
	9c17  8 Series PCI Express Root Port 4
	9c18  8 Series PCI Express Root Port 5
	9c19  8 Series PCI Express Root Port 5
	9c1a  8 Series PCI Express Root Port 6
	9c1b  8 Series PCI Express Root Port 6
	9c1c  8 Series PCI Express Root Port 7
	9c1d  8 Series PCI Express Root Port 7
	9c1e  8 Series PCI Express Root Port 8
	9c1f  8 Series PCI Express Root Port 8
	9c20  8 Series HD Audio Controller
	9c21  8 Series HD Audio Controller
	9c22  8 Series SMBus Controller
	9c23  8 Series CHAP Counters
	9c24  8 Series Thermal
	9c26  8 Series USB EHCI #1
	9c2d  8 Series USB EHCI #2
	9c31  8 Series USB xHCI HC
	9c35  8 Series SDIO Controller
	9c36  8 Series Audio DSP Controller
	9c3a  8 Series HECI #0
	9c3b  8 Series HECI #1
	9c3c  8 Series HECI IDER
	9c3d  8 Series HECI KT
	9c40  8 Series LPC Controller
	9c41  8 Series LPC Controller
	9c42  8 Series LPC Controller
	9c43  8 Series LPC Controller
	9c44  8 Series LPC Controller
	9c45  8 Series LPC Controller
	9c46  8 Series LPC Controller
	9c47  8 Series LPC Controller
	9c60  8 Series Low Power Sub-System DMA
	9c61  8 Series I2C Controller #0
	9c62  8 Series I2C Controller #1
	9c63  8 Series UART Controller #0
	9c64  8 Series UART Controller #1
	9c65  8 Series SPI Controller #0
	9c66  8 Series SPI Controller #1
	9c83  Wildcat Point-LP SATA Controller [AHCI Mode]
	9c85  Wildcat Point-LP SATA Controller [RAID Mode]
	9c87  Wildcat Point-LP SATA Controller [RAID Mode]
	9c8f  Wildcat Point-LP SATA Controller [RAID Mode]
	9c90  Wildcat Point-LP PCI Express Root Port #1
	9c92  Wildcat Point-LP PCI Express Root Port #2
	9c94  Wildcat Point-LP PCI Express Root Port #3
	9c96  Wildcat Point-LP PCI Express Root Port #4
	9c98  Wildcat Point-LP PCI Express Root Port #5
	9c9a  Wildcat Point-LP PCI Express Root Port #6
	9ca0  Wildcat Point-LP High Definition Audio Controller
	9ca2  Wildcat Point-LP SMBus Controller
	9ca4  Wildcat Point-LP Thermal Management Controller
	9ca6  Wildcat Point-LP USB EHCI Controller
	9cb1  Wildcat Point-LP USB xHCI Controller
	9cb5  Wildcat Point-LP Secure Digital IO Controller
	9cb6  Wildcat Point-LP Smart Sound Technology Controller
	9cba  Wildcat Point-LP MEI Controller #1
	9cbb  Wildcat Point-LP MEI Controller #2
	9cbc  Wildcat Point-LP IDE-r Controller
	9cbd  Wildcat Point-LP KT Controller
	9cc1  Wildcat Point-LP LPC Controller
	9cc2  Wildcat Point-LP LPC Controller
	9cc3  Wildcat Point-LP LPC Controller
	9cc5  Wildcat Point-LP LPC Controller
	9cc6  Wildcat Point-LP LPC Controller
	9cc7  Wildcat Point-LP LPC Controller
	9cc9  Wildcat Point-LP LPC Controller
	9ce0  Wildcat Point-LP Serial IO DMA Controller
	9ce1  Wildcat Point-LP Serial IO I2C Controller #0
	9ce2  Wildcat Point-LP Serial IO I2C Controller #1
	9ce3  Wildcat Point-LP Serial IO UART Controller #0
	9ce4  Wildcat Point-LP Serial IO UART Controller #1
	9ce5  Wildcat Point-LP Serial IO GSPI Controller #0
	9ce6  Wildcat Point-LP Serial IO GSPI Controller #1
	9d03  Sunrise Point-LP SATA Controller [AHCI mode]
	9d10  Sunrise Point-LP PCI Express Root Port #1
	9d12  Sunrise Point-LP PCI Express Root Port #3
	9d14  Sunrise Point-LP PCI Express Root Port #5
	9d15  Sunrise Point-LP PCI Express Root Port #6
	9d16  Sunrise Point-LP PCI Express Root Port #7
	9d17  Sunrise Point-LP PCI Express Root Port #8
	9d18  Sunrise Point-LP PCI Express Root Port #9
	9d19  Sunrise Point-LP PCI Express Root Port #10
	9d21  Sunrise Point-LP PMC
	9d23  Sunrise Point-LP SMBus
	9d27  Sunrise Point-LP Serial IO UART Controller #0
	9d28  Sunrise Point-LP Serial IO UART Controller #1
	9d29  Sunrise Point-LP Serial IO SPI Controller #0
	9d2a  Sunrise Point-LP Serial IO SPI Controller #1
	9d2d  Sunrise Point-LP Secure Digital IO Controller
	9d2f  Sunrise Point-LP USB 3.0 xHCI Controller
	9d31  Sunrise Point-LP Thermal subsystem
	9d35  Sunrise Point-LP Integrated Sensor Hub
	9d3a  Sunrise Point-LP CSME HECI #1
	9d43  Sunrise Point-LP LPC Controller
	9d48  Sunrise Point-LP LPC Controller
	9d56  Sunrise Point-LP LPC Controller
	9d58  Sunrise Point-LP LPC Controller
	9d60  Sunrise Point-LP Serial IO I2C Controller #0
	9d61  Sunrise Point-LP Serial IO I2C Controller #1
	9d62  Sunrise Point-LP Serial IO I2C Controller #2
	9d63  Sunrise Point-LP Serial IO I2C Controller #3
	9d64  Sunrise Point-LP Serial IO I2C Controller #4
	9d65  Sunrise Point-LP Serial IO I2C Controller #5
	9d66  Sunrise Point-LP Serial IO UART Controller #2
	9d70  Sunrise Point-LP HD Audio
	9d71  Sunrise Point-LP HD Audio
	a000  Atom Processor D4xx/D5xx/N4xx/N5xx DMI Bridge
	a001  Atom Processor D4xx/D5xx/N4xx/N5xx Integrated Graphics Controller
	a002  Atom Processor D4xx/D5xx/N4xx/N5xx Integrated Graphics Controller
	a003  Atom Processor D4xx/D5xx/N4xx/N5xx CHAPS counter
	a010  Atom Processor D4xx/D5xx/N4xx/N5xx DMI Bridge
	a011  Atom Processor D4xx/D5xx/N4xx/N5xx Integrated Graphics Controller
	a012  Atom Processor D4xx/D5xx/N4xx/N5xx Integrated Graphics Controller
	a013  Atom Processor D4xx/D5xx/N4xx/N5xx CHAPS counter
	a102  Sunrise Point-H SATA controller [AHCI mode]
	a103  Sunrise Point-H SATA Controller [AHCI mode]
	a105  Sunrise Point-H SATA Controller [RAID mode]
	a107  Sunrise Point-H SATA Controller [RAID mode]
	a10f  Sunrise Point-H SATA Controller [RAID mode]
	a110  Sunrise Point-H PCI Express Root Port #1
	a111  Sunrise Point-H PCI Express Root Port #2
	a112  Sunrise Point-H PCI Express Root Port #3
	a113  Sunrise Point-H PCI Express Root Port #4
	a114  Sunrise Point-H PCI Express Root Port #5
	a115  Sunrise Point-H PCI Express Root Port #6
	a116  Sunrise Point-H PCI Express Root Port #7
	a117  Sunrise Point-H PCI Express Root Port #8
	a118  Sunrise Point-H PCI Express Root Port #9
	a119  Sunrise Point-H PCI Express Root Port #10
	a11a  Sunrise Point-H PCI Express Root Port #11
	a11b  Sunrise Point-H PCI Express Root Port #12
	a11c  Sunrise Point-H PCI Express Root Port #13
	a11d  Sunrise Point-H PCI Express Root Port #14
	a11e  Sunrise Point-H PCI Express Root Port #15
	a11f  Sunrise Point-H PCI Express Root Port #16
	a120  Sunrise Point-H P2SB
	a121  Sunrise Point-H PMC
	a122  Sunrise Point-H cAVS
	a123  Sunrise Point-H SMBus
	a124  Sunrise Point-H SPI Controller
	a125  Sunrise Point-H Gigabit Ethernet Controller
	a126  Sunrise Point-H Northpeak
	a127  Sunrise Point-H Serial IO UART #0
	a128  Sunrise Point-H Serial IO UART #1
	a129  Sunrise Point-H Serial IO SPI #0
	a12a  Sunrise Point-H Serial IO SPI #1
	a12f  Sunrise Point-H USB 3.0 xHCI Controller
	a130  Sunrise Point-H USB Device Controller (OTG)
	a131  Sunrise Point-H Thermal subsystem
	a133  Sunrise Point-H Northpeak ACPI Function
	a135  Sunrise Point-H Integrated Sensor Hub
	a13a  Sunrise Point-H CSME HECI #1
	a13b  Sunrise Point-H CSME HECI #2
	a13c  Sunrise Point-H CSME IDE Redirection
	a13d  Sunrise Point-H KT Redirection
	a13e  Sunrise Point-H CSME HECI #3
	a140  Sunrise Point-H LPC Controller
	a141  Sunrise Point-H LPC Controller
	a142  Sunrise Point-H LPC Controller
	a143  Sunrise Point-H LPC Controller
	a144  Sunrise Point-H LPC Controller
	a145  Sunrise Point-H LPC Controller
	a146  Sunrise Point-H LPC Controller
	a147  Sunrise Point-H LPC Controller
	a148  Sunrise Point-H LPC Controller
	a149  Sunrise Point-H LPC Controller
	a14a  Sunrise Point-H LPC Controller
	a14b  Sunrise Point-H LPC Controller
	a14c  Sunrise Point-H LPC Controller
	a14d  Sunrise Point-H LPC Controller
	a14e  Sunrise Point-H LPC Controller
	a14f  Sunrise Point-H LPC Controller
	a150  Sunrise Point-H LPC Controller
	a151  Sunrise Point-H LPC Controller
	a152  Sunrise Point-H LPC Controller
	a153  Sunrise Point-H LPC Controller
	a154  Sunrise Point-H LPC Controller
	a155  Sunrise Point-H LPC Controller
	a156  Sunrise Point-H LPC Controller
	a157  Sunrise Point-H LPC Controller
	a158  Sunrise Point-H LPC Controller
	a159  Sunrise Point-H LPC Controller
	a15a  Sunrise Point-H LPC Controller
	a15b  Sunrise Point-H LPC Controller
	a15c  Sunrise Point-H LPC Controller
	a15d  Sunrise Point-H LPC Controller
	a15e  Sunrise Point-H LPC Controller
	a15f  Sunrise Point-H LPC Controller
	a160  Sunrise Point-H Serial IO I2C Controller #0
	a161  Sunrise Point-H Serial IO I2C Controller #1
	a166  Sunrise Point-H Serial IO UART Controller #2
	a167  Sunrise Point-H PCI Root Port #17
	a168  Sunrise Point-H PCI Root Port #18
	a169  Sunrise Point-H PCI Root Port #19
	a16a  Sunrise Point-H PCI Root Port #20
	a170  Sunrise Point-H HD Audio
	a171  CM238 HD Audio Controller
	a182  Lewisburg SATA Controller [AHCI mode]
	a186  Lewisburg SATA Controller [RAID mode]
	a190  Lewisburg PCI Express Root Port #1
	a191  Lewisburg PCI Express Root Port #2
	a192  Lewisburg PCI Express Root Port #3
	a193  Lewisburg PCI Express Root Port #4
	a194  Lewisburg PCI Express Root Port #5
	a195  Lewisburg PCI Express Root Port #6
	a196  Lewisburg PCI Express Root Port #7
	a197  Lewisburg PCI Express Root Port #8
	a198  Lewisburg PCI Express Root Port #9
	a199  Lewisburg PCI Express Root Port #10
	a19a  Lewisburg PCI Express Root Port #11
	a19b  Lewisburg PCI Express Root Port #12
	a19c  Lewisburg PCI Express Root Port #13
	a19d  Lewisburg PCI Express Root Port #14
	a19e  Lewisburg PCI Express Root Port #15
	a19f  Lewisburg PCI Express Root Port #16
	a1a0  Lewisburg P2SB
	a1a1  Lewisburg PMC
	a1a2  Lewisburg cAVS
	a1a3  Lewisburg SMBus
	a1a4  Lewisburg SPI Controller
	a1af  Lewisburg USB 3.0 xHCI Controller
	a1b1  Lewisburg Thermal Subsystem
	a1ba  Lewisburg CSME: HECI #1
	a1bb  Lewisburg CSME: HECI #2
	a1bc  Lewisburg CSME: IDE-r
	a1bd  Lewisburg CSME: KT Controller
	a1be  Lewisburg CSME: HECI #3
	a1c1  Lewisburg LPC Controller
	a1c2  Lewisburg LPC Controller
	a1c3  Lewisburg LPC Controller
	a1c4  Lewisburg LPC Controller
	a1c5  Lewisburg LPC Controller
	a1c6  Lewisburg LPC Controller
	a1c7  Lewisburg LPC Controller
	a1d2  Lewisburg SSATA Controller [AHCI mode]
	a1d6  Lewisburg SSATA Controller [RAID mode]
	a1e7  Lewisburg PCI Express Root Port #17
	a1e8  Lewisburg PCI Express Root Port #18
	a1e9  Lewisburg PCI Express Root Port #19
	a1ea  Lewisburg PCI Express Root Port #20
	a1f0  Lewisburg MROM 0
	a1f1  Lewisburg MROM 1
	a1f8  Lewisburg IE: HECI #1
	a1f9  Lewisburg IE: HECI #2
	a1fa  Lewisburg IE: IDE-r
	a1fb  Lewisburg IE: KT Controller
	a1fc  Lewisburg IE: HECI #3
	a202  Lewisburg SATA Controller [AHCI mode]
	a206  Lewisburg SATA Controller [RAID mode]
	a223  Lewisburg SMBus
	a224  Lewisburg SPI Controller
	a242  Lewisburg LPC or eSPI Controller
	a243  Lewisburg LPC or eSPI Controller
	a252  Lewisburg SSATA Controller [AHCI mode]
	a256  Lewisburg SSATA Controller [RAID mode]
	a282  200 Series PCH SATA controller [AHCI mode]
	a286  200 Series PCH SATA controller [RAID mode]
	a290  200 Series PCH PCI Express Root Port #1
	a291  200 Series PCH PCI Express Root Port #2
	a292  200 Series PCH PCI Express Root Port #3
	a293  200 Series PCH PCI Express Root Port #4
	a294  200 Series PCH PCI Express Root Port #5
	a295  200 Series PCH PCI Express Root Port #6
	a296  200 Series PCH PCI Express Root Port #7
	a297  200 Series PCH PCI Express Root Port #8
	a298  200 Series PCH PCI Express Root Port #9
	a299  200 Series PCH PCI Express Root Port #10
	a29a  200 Series PCH PCI Express Root Port #11
	a29b  200 Series PCH PCI Express Root Port #12
	a29c  200 Series PCH PCI Express Root Port #13
	a29d  200 Series PCH PCI Express Root Port #14
	a29e  200 Series PCH PCI Express Root Port #15
	a29f  200 Series PCH PCI Express Root Port #16
	a2a1  200 Series PCH PMC
	a2a3  200 Series PCH SMBus Controller
	a2a7  200 Series PCH Serial IO UART Controller #0
	a2a8  200 Series PCH Serial IO UART Controller #1
	a2a9  200 Series PCH Serial IO SPI Controller #0
	a2aa  200 Series PCH Serial IO SPI Controller #1
	a2af  200 Series PCH USB 3.0 xHCI Controller
	a2b1  200 Series PCH Thermal Subsystem
	a2ba  200 Series PCH CSME HECI #1
	a2bb  200 Series PCH CSME HECI #2
	a2c4  200 Series PCH LPC Controller (H270)
	a2c5  200 Series PCH LPC Controller (Z270)
	a2c6  200 Series PCH LPC Controller (Q270)
	a2c7  200 Series PCH LPC Controller (Q250)
	a2c8  200 Series PCH LPC Controller (B250)
	a2e0  200 Series PCH Serial IO I2C Controller #0
	a2e1  200 Series PCH Serial IO I2C Controller #1
	a2e2  200 Series PCH Serial IO I2C Controller #2
	a2e3  200 Series PCH Serial IO I2C Controller #3
	a2e6  200 Series PCH Serial IO UART Controller #2
	a2e7  200 Series PCH PCI Express Root Port #17
	a2e8  200 Series PCH PCI Express Root Port #18
	a2e9  200 Series PCH PCI Express Root Port #19
	a2ea  200 Series PCH PCI Express Root Port #20
	a2eb  200 Series PCH PCI Express Root Port #21
	a2ec  200 Series PCH PCI Express Root Port #22
	a2ed  200 Series PCH PCI Express Root Port #23
	a2ee  200 Series PCH PCI Express Root Port #24
	a2f0  200 Series PCH HD Audio
	a620  6400/6402 Advanced Memory Buffer (AMB)
	abc0  Omni-Path Fabric Switch Silicon 100 Series
	b152  21152 PCI-to-PCI Bridge
	b154  21154 PCI-to-PCI Bridge
	b555  21555 Non transparent PCI-to-PCI Bridge
	d130  Core Processor DMI
	d131  Core Processor DMI
	d132  Core Processor DMI
	d133  Core Processor DMI
	d134  Core Processor DMI
	d135  Core Processor DMI
	d136  Core Processor DMI
	d137  Core Processor DMI
	d138  Core Processor PCI Express Root Port 1
	d139  Core Processor PCI Express Root Port 2
	d13a  Core Processor PCI Express Root Port 3
	d13b  Core Processor PCI Express Root Port 4
	d150  Core Processor QPI Link
	d151  Core Processor QPI Routing and Protocol Registers
	d155  Core Processor System Management Registers
	d156  Core Processor Semaphore and Scratchpad Registers
	d157  Core Processor System Control and Status Registers
	d158  Core Processor Miscellaneous Registers
80ee  InnoTek Systemberatung GmbH
	beef  VirtualBox Graphics Adapter
	cafe  VirtualBox Guest Service
8322  Sodick America Corp.
8384  SigmaTel
8401  TRENDware International Inc.
8686  ScaleMP
	1010  vSMP Foundation controller [vSMP CTL]
	1011  vSMP Foundation MEX/FLX controller [vSMP CTL]
8800  Trigem Computer Inc.
	2008  Video assistant component
8866  T-Square Design Inc.
8888  Silicon Magic
8912  TRX
8c4a  Winbond
	1980  W89C940 misprogrammed [ne2k]
8e0e  Computone Corporation
8e2e  KTI
	3000  ET32P2
9004  Adaptec
	0078  AHA-2940U_CN
	1078  AIC-7810
	1160  AIC-1160 [Family Fibre Channel Adapter]
	2178  AIC-7821
	3860  AHA-2930CU
	3b78  AHA-4844W/4844UW
	5075  AIC-755x
	5078  AIC-7850T/7856T [AVA-2902/4/6 / AHA-2910]
	5175  AIC-755x
	5178  AIC-7851
	5275  AIC-755x
	5278  AIC-7852
	5375  AIC-755x
	5378  AIC-7850
	5475  AIC-755x
	5478  AIC-7850
	5575  AVA-2930
	5578  AIC-7855
	5647  ANA-7711 TCP Offload Engine
	5675  AIC-755x
	5678  AIC-7856
	5775  AIC-755x
	5778  AIC-7850
	5800  AIC-5800
	5900  ANA-5910/5930/5940 ATM155 & 25 LAN Adapter
	5905  ANA-5910A/5930A/5940A ATM Adapter
	6038  AIC-3860
	6075  AIC-1480 / APA-1480
	6078  AIC-7860
	6178  AIC-7861
	6278  AIC-7860
	6378  AIC-7860
	6478  AIC-786x
	6578  AIC-786x
	6678  AIC-786x
	6778  AIC-786x
	6915  ANA620xx/ANA69011A
	7078  AHA-294x / AIC-7870
	7178  AIC-7870P/7871 [AHA-2940/W/S76]
	7278  AHA-3940/3940W / AIC-7872
	7378  AHA-3985 / AIC-7873
	7478  AHA-2944/2944W / AIC-7874
	7578  AHA-3944/3944W / AIC-7875
	7678  AHA-4944W/UW / AIC-7876
	7710  ANA-7711F Network Accelerator Card (NAC) - Optical
	7711  ANA-7711C Network Accelerator Card (NAC) - Copper
	7778  AIC-787x
	7810  AIC-7810
	7815  AIC-7815 RAID+Memory Controller IC
	7850  AIC-7850
	7855  AHA-2930
	7860  AIC-7860
	7870  AIC-7870
	7871  AHA-2940
	7872  AHA-3940
	7873  AHA-3980
	7874  AHA-2944
	7880  AIC-7880P
	7890  AIC-7890
	7891  AIC-789x
	7892  AIC-789x
	7893  AIC-789x
	7894  AIC-789x
	7895  AHA-2940U/UW / AHA-39xx / AIC-7895
	7896  AIC-789x
	7897  AIC-789x
	8078  AIC-7880U
	8178  AIC-7870P/7881U [AHA-2940U/UW/D/S76]
	8278  AHA-3940U/UW/UWD / AIC-7882U
	8378  AHA-3940U/UW / AIC-7883U
	8478  AHA-2944UW / AIC-7884U
	8578  AHA-3944U/UWD / AIC-7885
	8678  AHA-4944UW / AIC-7886
	8778  AHA-2940UW Pro / AIC-788x
	8878  AHA-2930UW / AIC-7888
	8b78  ABA-1030
	ec78  AHA-4944W/UW
9005  Adaptec
	0010  AHA-2940U2/U2W
	0011  AHA-2930U2
	0013  78902
	001f  AHA-2940U2/U2W / 7890/7891
	0020  AIC-7890
	002f  AIC-7890
	0030  AIC-7890
	003f  AIC-7890
	0050  AHA-3940U2x/395U2x
	0051  AHA-3950U2D
	0053  AIC-7896 SCSI Controller
	005f  AIC-7896U2/7897U2
	0080  AIC-7892A U160/m
	0081  AIC-7892B U160/m
	0083  AIC-7892D U160/m
	008f  AIC-7892P U160/m
	0092  AVC-2010 [VideoH!]
	0093  AVC-2410 [VideoH!]
	00c0  AHA-3960D / AIC-7899A U160/m
	00c1  AIC-7899B U160/m
	00c3  AIC-7899D U160/m
	00c5  RAID subsystem HBA
	00cf  AIC-7899P U160/m
	0241  Serial ATA II RAID 1420SA
	0242  Serial ATA II RAID 1220SA
	0243  Serial ATA II RAID 1430SA
	0244  eSATA II RAID 1225SA
	0250  ServeRAID Controller
	0279  ServeRAID 6M
	0283  AAC-RAID
	0284  AAC-RAID
	0285  AAC-RAID
	0286  AAC-RAID (Rocket)
	028b  Series 6 - 6G SAS/PCIe 2
	028c  Series 7 6G SAS/PCIe 3
	028d  Series 8 12G SAS/PCIe 3
	028f  Smart Storage PQI 12G SAS/PCIe 3
	0410  AIC-9410W SAS (Razor HBA RAID)
	0412  AIC-9410W SAS (Razor HBA non-RAID)
	0415  ASC-58300 SAS (Razor-External HBA RAID)
	0416  ASC-58300 SAS (Razor-External HBA non-RAID)
	041e  AIC-9410W SAS (Razor ASIC non-RAID)
	041f  AIC-9410W SAS (Razor ASIC RAID)
	042f  VSC7250/7251 SAS (Aurora ASIC non-RAID)
	0430  AIC-9405W SAS (Razor-Lite HBA RAID)
	0432  AIC-9405W SAS (Razor-Lite HBA non-RAID)
	043e  AIC-9405W SAS (Razor-Lite ASIC non-RAID)
	043f  AIC-9405W SAS (Razor-Lite ASIC RAID)
	0450  ASC-1405 Unified Serial HBA
	0500  Obsidian chipset SCSI controller
	0503  Scamp chipset SCSI controller
	0910  AUA-3100B
	091e  AUA-3100B
	8000  ASC-29320A U320
	800f  AIC-7901 U320
	8010  ASC-39320 U320
	8011  ASC-39320D
	8012  ASC-29320 U320
	8013  ASC-29320B U320
	8014  ASC-29320LP U320
	8015  ASC-39320B U320
	8016  ASC-39320A U320
	8017  ASC-29320ALP U320
	801c  ASC-39320D U320
	801d  AIC-7902B U320
	801e  AIC-7901A U320
	801f  AIC-7902 U320
	8080  ASC-29320A U320 w/HostRAID
	8081  PMC-Sierra PM8001 SAS HBA [Series 6H]
	8088  PMC-Sierra PM8018 SAS HBA [Series 7H]
	8089  PMC-Sierra PM8019 SAS encryption HBA [Series 7He]
	808f  AIC-7901 U320 w/HostRAID
	8090  ASC-39320 U320 w/HostRAID
	8091  ASC-39320D U320 w/HostRAID
	8092  ASC-29320 U320 w/HostRAID
	8093  ASC-29320B U320 w/HostRAID
	8094  ASC-29320LP U320 w/HostRAID
	8095  ASC-39320(B) U320 w/HostRAID
	8096  ASC-39320A U320 w/HostRAID
	8097  ASC-29320ALP U320 w/HostRAID
	809c  ASC-39320D(B) U320 w/HostRAID
	809d  AIC-7902(B) U320 w/HostRAID
	809e  AIC-7901A U320 w/HostRAID
	809f  AIC-7902 U320 w/HostRAID
907f  Atronics
	2015  IDE-2015PL
919a  Gigapixel Corp
9412  Holtek
	6565  6565
9413  Softlogic Co., Ltd.
	6010  SOLO6010 MPEG-4 Video encoder/decoder
	6110  SOLO6110 H.264 Video encoder/decoder
9618  JusonTech Corporation
	0001  JusonTech Gigabit Ethernet Controller
9699  Omni Media Technology Inc
	6565  6565
9710  MosChip Semiconductor Technology Ltd.
	9250  PCI-to-PCI bridge [MCS9250]
	9805  PCI 1 port parallel adapter
	9815  PCI 9815 Multi-I/O Controller
	9820  PCI 9820 Multi-I/O Controller
	9835  PCI 9835 Multi-I/O Controller
	9845  PCI 9845 Multi-I/O Controller
	9855  PCI 9855 Multi-I/O Controller
	9865  PCI 9865 Multi-I/O Controller
	9901  PCIe 9901 Multi-I/O Controller
	9904  4-Port PCIe Serial Adapter
	9912  PCIe 9912 Multi-I/O Controller
	9922  MCS9922 PCIe Multi-I/O Controller
	9990  MCS9990 PCIe to 4‐Port USB 2.0 Host Controller
9850  3Com (wrong ID)
9902  Stargen Inc.
	0001  SG2010 PCI over Starfabric Bridge
	0002  SG2010 PCI to Starfabric Gateway
	0003  SG1010 Starfabric Switch and PCI Bridge
a0a0  AOPEN Inc.
a0f1  UNISYS Corporation
a200  NEC Corporation
a259  Hewlett Packard
a25b  Hewlett Packard GmbH PL24-MKT
a304  Sony
a727  3Com Corporation
	0013  3CRPAG175 Wireless PC Card
	6803  3CRDAG675B Wireless 11a/b/g Adapter
aa00  iTuner
aa01  iTuner
aa02  iTuner
aa03  iTuner
aa04  iTuner
aa05  iTuner
aa06  iTuner
aa07  iTuner
aa08  iTuner
aa09  iTuner
aa0a  iTuner
aa0b  iTuner
aa0c  iTuner
aa0d  iTuner
aa0e  iTuner
aa0f  iTuner
aa42  Scitex Digital Video
aa55  Ncomputing X300 PCI-Engine
aaaa  Adnaco Technology Inc.
	0001  H1 PCIe over fiber optic host controller
	0002  R1BP1 PCIe over fiber optic expansion chassis
abcd  Vadatech Inc.
ac1e  Digital Receiver Technology Inc
ac3d  Actuality Systems
ad00  Alta Data Technologies LLC
aecb  Adrienne Electronics Corporation
	6250  VITC/LTC Timecode Reader card [PCI-VLTC/RDR]
affe  Sirrix AG security technologies
	01e1  PCI1E1 1-port ISDN E1 interface
	02e1  PCI2E1 2-port ISDN E1 interface
	450e  PCI4S0EC 4-port ISDN S0 interface
	dead  Sirrix.PCI4S0 4-port ISDN S0 interface
b100  OpenVox Communication Co. Ltd.
b10b  Uakron PCI Project
b1b3  Shiva Europe Limited
b1d9  ATCOM Technology co., LTD.
bd11  Pinnacle Systems, Inc. (Wrong ID)
bdbd  Blackmagic Design
	a106  Multibridge Extreme
	a117  Intensity Pro
	a11a  DeckLink HD Extreme 2
	a11b  DeckLink SDI/Duo/Quad
	a11c  DeckLink HD Extreme 3
	a11d  DeckLink Studio
	a11e  DeckLink Optical Fibre
	a120  Decklink Studio 2
	a121  DeckLink HD Extreme 3D/3D+
	a124  Intensity Extreme
	a126  Intensity Shuttle
	a127  UltraStudio Express
	a129  UltraStudio Mini Monitor
	a12a  UltraStudio Mini Recorder
	a12d  UltraStudio 4K
	a12e  DeckLink 4K Extreme
	a12f  DeckLink Mini Monitor
	a130  DeckLink Mini Recorder
	a132  UltraStudio 4K
	a136  DeckLink 4K Extreme 12G
	a137  DeckLink Studio 4K
	a138  Decklink SDI 4K
	a139  Intensity Pro 4K
	a13b  DeckLink Micro Recorder
	a13d  DeckLink 4K Pro
	a13e  UltraStudio 4K Extreme
	a13f  DeckLink Quad 2
	a140  DeckLink Duo 2
c001  TSI Telsys
c0a9  Micron/Crucial Technology
c0de  Motorola
c0fe  Motion Engineering, Inc.
ca50  Varian Australia Pty Ltd
cace  CACE Technologies, Inc.
	0001  TurboCap Port A
	0002  TurboCap Port B
	0023  AirPcap N
caed  Canny Edge
cafe  Chrysalis-ITS
	0003  Luna K3 Hardware Security Module
	0006  Luna PCI-e 3000 Hardware Security Module
cccc  Catapult Communications
ccec  Curtiss-Wright Controls Embedded Computing
cddd  Tyzx, Inc.
	0101  DeepSea 1 High Speed Stereo Vision Frame Grabber
	0200  DeepSea 2 High Speed Stereo Vision Frame Grabber
ceba  KEBA AG
d161  Digium, Inc.
	0120  Wildcard TE120P single-span T1/E1/J1 card
	0205  Wildcard TE205P/TE207P dual-span T1/E1/J1 card 5.0V
	0210  Wildcard TE210P/TE212P dual-span T1/E1/J1 card 3.3V
	0220  Wildcard TE220 dual-span T1/E1/J1 card 3.3V (PCI-Express)
	0405  Wildcard TE405P/TE407P quad-span T1/E1/J1 card 5.0V
	0410  Wildcard TE410P/TE412P quad-span T1/E1/J1 card 3.3V
	0420  Wildcard TE420P quad-span T1/E1/J1 card 3.3V (PCI-Express)
	0800  Wildcard TDM800P 8-port analog card
	1205  Wildcard TE205P/TE207P dual-span T1/E1/J1 card 5.0V (u1)
	1220  Wildcard TE220 dual-span T1/E1/J1 card 3.3V (PCI-Express) (5th gen)
	1405  Wildcard TE405P/TE407P quad-span T1/E1/J1 card 5.0V (u1)
	1410  Wildcard TE410P quad-span T1/E1/J1 card 3.3V (5th Gen)
	1420  Wildcard TE420 quad-span T1/E1/J1 card 3.3V (PCI-Express) (5th gen)
	1820  Wildcard TE820 octal-span T1/E1/J1 card 3.3V (PCI-Express)
	2400  Wildcard TDM2400P 24-port analog card
	3400  Wildcard TC400P transcoder base card
	8000  Wildcard TE121 single-span T1/E1/J1 card (PCI-Express)
	8001  Wildcard TE122 single-span T1/E1/J1 card
	8002  Wildcard AEX800 8-port analog card (PCI-Express)
	8003  Wildcard AEX2400 24-port analog card (PCI-Express)
	8004  Wildcard TCE400P transcoder base card
	8005  Wildcard TDM410 4-port analog card
	8006  Wildcard AEX410 4-port analog card (PCI-Express)
	8007  Hx8 Series 8-port Base Card
	8008  Hx8 Series 8-port Base Card (PCI-Express)
	800a  Wildcard TE133 single-span T1/E1/J1 card (PCI Express)
	800b  Wildcard TE134 single-span T1/E1/J1 card
	800c  Wildcard A8A 8-port analog card
	800d  Wildcard A8B 8-port analog card (PCI-Express)
	800e  Wildcard TE235/TE435 quad-span T1/E1/J1 card (PCI-Express)
	800f  Wildcard A4A 4-port analog card
	8010  Wildcard A4B 4-port analog card (PCI-Express)
	8013  Wildcard TE236/TE436 quad-span T1/E1/J1 card
	b410  Wildcard B410 quad-BRI card
d4d4  Dy4 Systems Inc
	0601  PCI Mezzanine Card
d531  I+ME ACTIA GmbH
d84d  Exsys
dada  Datapath Limited
	0133  VisionRGB-X2
	0139  VisionRGB-E1
	0144  VisionSD8
	0150  VisionRGB-E2
	0151  VisionSD4+1
	0159  VisionAV
	0161  DGC161
	0165  DGC165
	0167  DGC167
	0168  DGC168
	1139  VisionRGB-E1S
	1150  VisionRGB-E2S
	1151  VisionSD4+1S
	1153  VisionDVI-DL
	1154  VisionSDI2
db10  Diablo Technologies
dc93  Dawicontrol GmbH
dcba  Dynamic Engineering
	0046  PCIe Altera Cyclone IV
	0047  VPX-RCB
	0048  PMC-Biserial-III-BAE9
	004e  PC104p-Biserial-III-NVY5
	004f  PC104p-Biserial-III-NVY6
	0052  PCIeBiSerialDb37 BA22 LVDS IO
dd01  Digital Devices GmbH
	0003  Octopus DVB Adapter
	0006  Cine V7
	0007  Max
	0011  Octopus CI DVB Adapter
	0201  Resi DVB-C Modulator
dead  Indigita Corporation
deaf  Middle Digital Inc.
	9050  PC Weasel Virtual VGA
	9051  PC Weasel Serial Port
	9052  PC Weasel Watchdog Timer
deda  XIMEA
	4001  Camera CB
	4021  Camera MT
e000  Winbond
	e000  W89C940
e159  Tiger Jet Network Inc.
	0001  Tiger3XX Modem/ISDN interface
	0002  Tiger100APC ISDN chipset
e1c5  Elcus
e4bf  EKF Elektronik GmbH
	0ccd  CCD-CALYPSO
	0cd1  CD1-OPERA
	0cd2  CD2-BEBOP
	0cd3  CD3-JIVE
	50c1  PC1-GROOVE
	50c2  PC2-LIMBO
	53c1  SC1-ALLEGRO
	cc47  CCG-RUMBA
	cc4d  CCM-BOOGIE
e4e4  Xorcom
e55e  Essence Technology, Inc.
ea01  Eagle Technology
	000a  PCI-773 Temperature Card
	0032  PCI-730 & PC104P-30 Card
	003e  PCI-762 Opto-Isolator Card
	0041  PCI-763 Reed Relay Card
	0043  PCI-769 Opto-Isolator Reed Relay Combo Card
	0046  PCI-766 Analog Output Card
	0052  PCI-703 Analog I/O Card
	0800  PCI-800 Digital I/O Card
ea60  RME
	9896  Digi32
	9897  Digi32 Pro
	9898  Digi32/8
eabb  Aashima Technology B.V.
eace  Endace Measurement Systems, Ltd
	3100  DAG 3.10 OC-3/OC-12
	3200  DAG 3.2x OC-3/OC-12
	320e  DAG 3.2E Fast Ethernet
	340e  DAG 3.4E Fast Ethernet
	341e  DAG 3.41E Fast Ethernet
	3500  DAG 3.5 OC-3/OC-12
	351c  DAG 3.5ECM Fast Ethernet
	360d  DAG 3.6D DS3
	360e  DAG 3.6E Fast Ethernet
	368e  DAG 3.6E Gig Ethernet
	3707  DAG 3.7T T1/E1/J1
	370d  DAG 3.7D DS3/E3
	378e  DAG 3.7G Gig Ethernet
	3800  DAG 3.8S OC-3/OC-12
	4100  DAG 4.10 OC-48
	4110  DAG 4.11 OC-48
	4220  DAG 4.2 OC-48
	422e  DAG 4.2GE Gig Ethernet
	4230  DAG 4.2S OC-48
	423e  DAG 4.2GE Gig Ethernet
	4300  DAG 4.3S OC-48
	430e  DAG 4.3GE Gig Ethernet
	452e  DAG 4.5G2 Gig Ethernet
	454e  DAG 4.5G4 Gig Ethernet
	45b8  DAG 4.5Z8 Gig Ethernet
	45be  DAG 4.5Z2 Gig Ethernet
	520e  DAG 5.2X 10G Ethernet
	521a  DAG 5.2SXA 10G Ethernet/OC-192
	5400  DAG 5.4S-12 OC-3/OC-12
	5401  DAG 5.4SG-48 Gig Ethernet/OC-3/OC-12/OC-48
	540a  DAG 5.4GA Gig Ethernet
	541a  DAG 5.4SA-12 OC-3/OC-12
	542a  DAG 5.4SGA-48 Gig Ethernet/OC-3/OC-12/OC-48
	6000  DAG 6.0SE 10G Ethernet/OC-192
	6100  DAG 6.1SE 10G Ethernet/OC-192
	6200  DAG 6.2SE 10G Ethernet/OC-192
	7100  DAG 7.1S OC-3/OC-12
	7400  DAG 7.4S OC-3/OC-12
	7401  DAG 7.4S48 OC-48
	752e  DAG 7.5G2 Gig Ethernet
	754e  DAG 7.5G4 Gig Ethernet
	8100  DAG 8.1X 10G Ethernet
	8101  DAG 8.1SX 10G Ethernet/OC-192
	8102  DAG 8.1X 10G Ethernet
	820e  DAG 8.2X 10G Ethernet
	820f  DAG 8.2X 10G Ethernet (2nd bus)
	8400  DAG 8.4I Infiniband x4 SDR
	8500  DAG 8.5I Infiniband x4 DDR
	9200  DAG 9.2SX2 10G Ethernet
	920e  DAG 9.2X2 10G Ethernet
	a120  DAG 10X2-P 10G Ethernet
	a12e  DAG 10X2-S 10G Ethernet
	a140  DAG 10X4-P 10G Ethernet
ec80  Belkin Corporation
	ec00  F5D6000
ecc0  Echo Digital Audio Corporation
edd8  ARK Logic Inc
	a091  1000PV [Stingray]
	a099  2000PV [Stingray]
	a0a1  2000MT
	a0a9  2000MI
f043  ASUSTeK Computer Inc. (Wrong ID)
f05b  Foxconn International, Inc. (Wrong ID)
f1d0  AJA Video
	c0fe  Xena HS/HD-R
	c0ff  Kona/Xena 2
	cafe  Kona SD
	cfee  Xena LS/SD-22-DA/SD-DA
	daff  KONA LHi
	db01  Corvid22
	db09  Corvid 24
	dcaf  Kona HD
	dfee  Xena HD-DA
	eb0e  Corvid 44
	efac  Xena SD-MM/SD-22-MM
	facd  Xena HD-MM
f5f5  F5 Networks, Inc.
f849  ASRock Incorporation (Wrong ID)
fa57  Interagon AS
	0001  PMC [Pattern Matching Chip]
fab7  Fabric7 Systems, Inc.
febd  Ultraview Corp.
feda  Broadcom Inc
	a0fa  BCM4210 iLine10 HomePNA 2.0
	a10e  BCM4230 iLine10 HomePNA 2.0
fede  Fedetec Inc.
	0003  TABIC PCI v3
fffd  XenSource, Inc.
	0101  PCI Event Channel Controller
fffe  VMWare Inc (temporary ID)
	0710  Virtual SVGA
ffff  Illegal Vendor ID
`)
	ids = parse(pciids)
	return ids
}
