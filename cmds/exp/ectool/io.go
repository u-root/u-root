// These should all implement io.ReadAt, with the address as the Offset; same for WriteAt.
package main

import (
	"encoding/binary"
	"os"
)

type ioaddr uint16

type ioport interface {
	Outb(ioaddr, uint8) error
	Outw(ioaddr, uint16) error
	Outl(ioaddr, uint32) error
	Outs(ioaddr, []uint8) (int, error)
	// For later, at some point, we may go with this.
	// Not yet.
	// Out(ioaddr, interface{}) (int, error)
	Inb(ioaddr) (uint8, error)
	Inw(ioaddr) (uint16, error)
	Inl(ioaddr) (uint32, error)
	Ins(ioaddr, int) ([]uint8, error)
	// In(ioaddr, interface{}) (int, error)
}

type debugf func(string, ...any)

func nodebugf(string, ...any) {}

type devports struct {
	*os.File
	debugf
}

func newDevPorts(d debugf) (ioport, error) {
	f, err := os.Create("/dev/port")
	if err != nil {
		return nil, err
	}
	p := &devports{File: f, debugf: nodebugf}
	if d != nil {
		p.debugf = d
	}

	return p, nil
}

func (p *devports) Outb(i ioaddr, d uint8) error {
	_, err := p.WriteAt([]byte{d}, int64(i))
	p.debugf("Write 0x%x @ 0x%x\n", d, i)
	return err
}

func (p *devports) Outs(i ioaddr, d []uint8) (int, error) {
	amt := 0
	for n, b := range d {
		err := p.Outb(i+ioaddr(n), b)
		if err != nil {
			return amt, err
		}
		amt++
	}
	return amt, nil
}

func (p *devports) Inb(i ioaddr) (uint8, error) {
	d := [1]byte{}
	_, err := p.ReadAt(d[:], int64(i))
	p.debugf("Read 0x%x @ 0x%x\n", d[0], i)
	return d[0], err
}

func (p *devports) Ins(i ioaddr, amt int) ([]uint8, error) {
	d := make([]byte, amt)
	var err error
	for n := range d {
		d[n], err = p.Inb(i + ioaddr(n))
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func (p *devports) Outw(i ioaddr, d uint16) error {
	b := [2]byte{}
	binary.LittleEndian.PutUint16(b[:], d)
	_, err := p.WriteAt(b[:], int64(i))
	p.debugf("Write 0x%x @ 0x%x\n", d, i)
	return err
}

func (p *devports) Inw(i ioaddr) (uint16, error) {
	b := [2]byte{}
	_, err := p.ReadAt(b[:], int64(i))
	d := binary.LittleEndian.Uint16(b[:])
	p.debugf("Read 0x%x @ 0x%x\n", d, i)
	return d, err
}

func (p *devports) Outl(i ioaddr, d uint32) error {
	b := [4]byte{}
	binary.LittleEndian.PutUint32(b[:], d)
	_, err := p.WriteAt(b[:], int64(i))
	p.debugf("Write 0x%x @ 0x%x\n", d, i)
	return err
}

func (p *devports) Inl(i ioaddr) (uint32, error) {
	b := [4]byte{}
	_, err := p.ReadAt(b[:], int64(i))
	d := binary.LittleEndian.Uint32(b[:])
	p.debugf("Read 0x%x @ 0x%x\n", d, i)
	return d, err
}
