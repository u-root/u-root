// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements the ioport functions for linux 386 using tinygo.
// Since we write C and not assembly, the compiler will adjust the instructions
// to fit for for 32 and 64 bit architectures.

#include <stdint.h>

uint32_t archInl(uint16_t port)
{
	uint32_t data;
	__asm__ __volatile__("inl %1, %0" : "=a"(data) : "d"(port));
	return data;
}

uint16_t archInw(uint16_t port)
{
	uint16_t data;
	__asm__ __volatile__("inw %1, %0" : "=a"(data) : "d"(port));
	return data;
}

uint8_t archInb(uint16_t port)
{
	uint8_t data;
	__asm__ __volatile__("inb %1, %0" : "=a"(data) : "d"(port));
	return data;
}

void archOutl(uint16_t port, uint32_t data)
{
	__asm__ __volatile__("outl %0, %1" : : "a"(data), "d"(port));
}

void archOutw(uint16_t port, uint16_t data)
{
	__asm__ __volatile__("outw %0, %1" : : "a"(data), "d"(port));
}

void archOutb(uint16_t port, uint8_t data)
{
	__asm__ __volatile__("outb %0, %1" : : "a"(data), "d"(port));
}