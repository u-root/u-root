/*
Copyright 2013-2014 Graham King

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

For full license details see <http://www.gnu.org/licenses/>.
*/

// our big change is that we just smash it all into one ether/IP/UDP header. 
// layers are for cakes.
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type EtherIPUDPHeader struct {
     Version uint8 // 4
     IHL     uint8 //4
     DSCP    uint8 //6
     ECN     uint8 //2
     TotalLength uint16
     Id uint16
     Flags uint8 // 3
    Fragoff uint16 //13
     TTL	   uint8
     Protocol	   uint8
     HCsum	   uint16	
     SIP	   uint32
     DIP	   uint32
     Options	   []uint8
     // UDP bits
     SPort  uint16
     DPort  uint16
     Length  uint16 // header + data
     Csum    uint16 // 0 if you're lazy
}

type TCPOption struct {
	Kind   uint8
	Length uint8
	Data   []byte
}

// Parse packet into TCPHeader structure
func NewEtherIPUDPHeader(data []byte) *EtherIPUDPHeader {
     var t8 uint8
     var t16 uint16
	u := &EtherIPUDPHeader{}
	r := bytes.NewReader(data)
/*
	binary.Read(r, binary.BigEndian, u.Dst)
	binary.Read(r, binary.BigEndian, u.Src)
	binary.Read(r, binary.BigEndian, u.Etype)
 */
	binary.Read(r, binary.BigEndian, &t8)
	u.Version = t8& 0xf
	u.IHL = t8>>4
	binary.Read(r, binary.BigEndian, &t8)
	u.DSCP = t8& 0x3f
	u.ECN = t8>> 6
	binary.Read(r, binary.BigEndian, &u.TotalLength)
	binary.Read(r, binary.BigEndian, &u.Id)
	binary.Read(r, binary.BigEndian, &t16)
	u.Flags = uint8(t16 & 3)
	u.Fragoff = t16 >>2
	binary.Read(r, binary.BigEndian, &u.TTL)
	binary.Read(r, binary.BigEndian, &u.Protocol)
	binary.Read(r, binary.BigEndian, &u.HCsum)
	binary.Read(r, binary.BigEndian, &u.SIP)
	binary.Read(r, binary.BigEndian, &u.DIP)

	binary.Read(r, binary.BigEndian, &u.SPort)
	binary.Read(r, binary.BigEndian, &u.DPort)
	binary.Read(r, binary.BigEndian, &u.Length)
	binary.Read(r, binary.BigEndian, &u.Csum)

	return u
}

func (u *EtherIPUDPHeader) Marshal(datapacket[]byte) []byte {

     var t8 uint8
     var t16 uint16
	buf := new(bytes.Buffer)
	// stupid code.
	u.HCsum = 0
	t8 = (u.Version <<4) | u.IHL
	binary.Write(buf, binary.BigEndian, t8)
	t8 = u.DSCP | (u.ECN << 6)
	binary.Write(buf, binary.BigEndian, t8)

	binary.Write(buf, binary.BigEndian, u.TotalLength)
	binary.Write(buf, binary.BigEndian, u.Id)
	t16 = uint16(u.Flags) | u.Fragoff << 13
	binary.Write(buf, binary.BigEndian, t16)
	binary.Write(buf, binary.BigEndian, u.TTL)
	binary.Write(buf, binary.BigEndian, u.Protocol)
	binary.Write(buf, binary.BigEndian, u.HCsum)
	binary.Write(buf, binary.BigEndian, u.SIP)
	binary.Write(buf, binary.BigEndian, u.DIP)

	binary.Write(buf, binary.BigEndian, u.SPort)
	binary.Write(buf, binary.BigEndian, u.DPort)
	binary.Write(buf, binary.BigEndian, u.Length)
	binary.Write(buf, binary.BigEndian, u.Csum)
	binary.Write(buf, binary.BigEndian, datapacket)
	out := buf.Bytes()
	u.HCsum = hcsum(out[0:20])
	out[10] = byte(u.HCsum >> 8)
	out[11] = byte(u.HCsum)

	// ah joy .. 
	cs := csum(out[16:], [4]byte{}, [4]byte{0xff,0xff,0xff,0xff}, byte(17))
	out[26] = byte(cs >> 8)
	out[27] = byte(cs)
	fmt.Printf("cs %v\n", cs)

	return out
}

// Checksum
func hcsum(sumThis []byte) uint16 {
	fmt.Printf("% x\n", sumThis)

	lenSumThis := len(sumThis)
	var nextWord uint16
	var sum uint32
	for i := 0; i+1 < lenSumThis; i += 2 {
		nextWord = uint16(sumThis[i])<<8 | uint16(sumThis[i+1])
		sum += uint32(nextWord)
	}
	if lenSumThis%2 != 0 {
		//fmt.Println("Odd byte")
		sum += uint32(sumThis[len(sumThis)-1])
	}

	// Add back any carry, and any carry from adding the carry
	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)

	// Bitwise complement
	return uint16(^sum)
}
func csum(data []byte, srcip, dstip [4]byte, proto byte) uint16 {
	pseudoHeader := []byte{
		srcip[0], srcip[1], srcip[2], srcip[3],
		dstip[0], dstip[1], dstip[2], dstip[3],
		0,                  // zero
		proto,              // protocol number (6 == TCP)
		0, byte(len(data)), // length (16 bits), not inc pseudo header
	}

	sumThis := make([]byte, 0, len(pseudoHeader)+len(data))
	sumThis = append(sumThis, pseudoHeader...)
	sumThis = append(sumThis, data...)
	//fmt.Printf("% x\n", sumThis)

	lenSumThis := len(sumThis)
	var nextWord uint16
	var sum uint32
	for i := 0; i+1 < lenSumThis; i += 2 {
		nextWord = uint16(sumThis[i])<<8 | uint16(sumThis[i+1])
		sum += uint32(nextWord)
	}
	if lenSumThis%2 != 0 {
		//fmt.Println("Odd byte")
		sum += uint32(sumThis[len(sumThis)-1])
	}

	// Add back any carry, and any carry from adding the carry
	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)

	// Bitwise complement
	return uint16(^sum)
}
