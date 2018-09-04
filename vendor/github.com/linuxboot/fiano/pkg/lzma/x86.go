// Copyright 2018 the LinuxBoot Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lzma

// DecodeX86 decodes LZMA data with the x86 extension.
func DecodeX86(encodedData []byte) ([]byte, error) {
	decodedData, err := Decode(encodedData)
	if err != nil {
		return nil, err
	}
	var x86State uint32
	x86Convert(decodedData, uint(len(decodedData)), 0, &x86State, false)
	return decodedData, nil
}

// EncodeX86 encodes LZMA data with the x86 extension.
func EncodeX86(decodedData []byte) ([]byte, error) {
	// x86Convert modifies the input, so a copy is recommened.
	decodedDataCpy := make([]byte, len(decodedData))
	copy(decodedDataCpy, decodedData)

	var x86State uint32
	x86Convert(decodedDataCpy, uint(len(decodedDataCpy)), 0, &x86State, true)
	return Encode(decodedDataCpy)
}

// Adapted from: https://github.com/tianocore/edk2/blob/00f5e11913a8706a1733da2b591502d59f848a99/BaseTools/Source/C/LzmaCompress/Sdk/C/Bra86.c
func x86Convert(data []byte, size uint, ip uint32, state *uint32, encoding bool) uint {
	var pos uint
	mask := *state & 7
	if size < 5 {
		return 0
	}
	size -= 4
	ip += 5

	for {
		p := pos
		for ; p < size; p++ {
			if data[p]&0xFE == 0xE8 {
				break
			}
		}

		{
			d := p - pos
			pos = p
			if p >= size {
				if d > 2 {
					*state = 0
				} else {
					*state = mask >> d
				}
				return pos
			}
			if d > 2 {
				mask = 0
			} else {
				mask >>= d
				if mask != 0 && (mask > 4 || mask == 3 || test86MSByte(data[p+uint(mask>>1)+1])) {
					mask = (mask >> 1) | 4
					pos++
					continue
				}
			}
		}

		if test86MSByte(data[p+4]) {
			v := (uint32(data[p+4]) << 24) + (uint32(data[p+3]) << 16) + (uint32(data[p+2]) << 8) + uint32(data[p+1])
			cur := ip + uint32(pos)
			pos += 5
			if encoding {
				v += cur
			} else {
				v -= cur
			}
			if mask != 0 {
				sh := uint((mask & 6) << 2)
				if test86MSByte(uint8(v >> sh)) {
					v ^= (uint32(0x100) << sh) - 1
					if encoding {
						v += cur
					} else {
						v -= cur
					}
				}
				mask = 0
			}
			data[p+1] = uint8(v)
			data[p+2] = uint8(v >> 8)
			data[p+3] = uint8(v >> 16)
			data[p+4] = uint8(0 - ((v >> 24) & 1))
		} else {
			mask = (mask >> 1) | 4
			pos++
		}
	}
}

func test86MSByte(b byte) bool {
	return (b+1)&0xFE == 0
}
