// Copyright (C) 2013 Patrick Georgi <patrick@georgi-clan.de>
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; version 2 of the License.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc.
//
// This is a very much cut down version of the coreboot version.
package main

const (
	TrampolinePointer  = 0x40000
	LinuxParamPointer  = 0x90000
	CommandLinePointer = 0x91000
)

var (
	trampoline = []byte{
		0xfa, //cli
		0xeb, 0xfe, // remove before flight ... jmp .
		0xb8, 0x00, 0x10, 0x40, 0x00, //mov    $0x401000,%eax
		0x0f, 0x01, 0x00, //sgdt   (%rax)
		0x8b, 0x58, 0x02, //mov    0x2(%rax),%ebx
		0xc7, 0x43, 0x10, 0xff, 0xff, 0x00, 0x00, //movl   $0xffff,0x10(%rbx)
		0xc7, 0x43, 0x14, 0x00, 0x9b, 0xcf, 0x00, //movl   $0xcf9b00,0x14(%rbx)
		0xc7, 0x43, 0x18, 0xff, 0xff, 0x00, 0x00, //movl   $0xffff,0x18(%rbx)
		0xc7, 0x43, 0x1c, 0x00, 0x93, 0xcf, 0x00, //movl   $0xcf9300,0x1c(%rbx)
		0xbe, 0x00, 0x00, 0x09, 0x00, //mov    $0x90000,%esi
		0xeb, 0xfe, // remove before flight ... jmp .
		0xea, 0x00, 0x00, 0x10, 0x00, 0x10, 0x00, // ljmp $0x10, $0x100000
		}
)
