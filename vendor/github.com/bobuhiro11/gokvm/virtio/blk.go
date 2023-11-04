package virtio

import (
	"bytes"
	"encoding/binary"
	"os"
	"unsafe"

	"github.com/bobuhiro11/gokvm/pci"
)

const (
	BlkIOPortStart = 0x6300
	BlkIOPortSize  = 0x100

	SectorSize = 512
)

type Blk struct {
	file *os.File
	Hdr  blkHdr

	VirtQueue    [1]*VirtQueue
	Mem          []byte
	LastAvailIdx [1]uint16

	kick chan interface{}

	irq         uint8
	IRQInjector IRQInjector
}

type blkHdr struct {
	commonHeader commonHeader
	blkHeader    blkHeader
}

func (h blkHdr) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, h); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

type blkHeader struct {
	capacity uint64
}

func (v Blk) GetDeviceHeader() pci.DeviceHeader {
	return pci.DeviceHeader{
		DeviceID:    0x1001,
		VendorID:    0x1AF4,
		HeaderType:  0,
		SubsystemID: 2, // Block Device
		Command:     1, // Enable IO port
		BAR: [6]uint32{
			BlkIOPortStart | 0x1,
		},
		// https://github.com/torvalds/linux/blob/fb3b0673b7d5b477ed104949450cd511337ba3c6/drivers/pci/setup-irq.c#L30-L55
		InterruptPin: 1,
		// https://www.webopedia.com/reference/irqnumbers/
		InterruptLine: v.irq,
	}
}

func (v Blk) Read(port uint64, bytes []byte) error {
	offset := int(port - BlkIOPortStart)

	b, err := v.Hdr.Bytes()
	if err != nil {
		return err
	}

	l := len(bytes)
	copy(bytes[:l], b[offset:offset+l])

	return nil
}

func (v *Blk) IOThreadEntry() {
	for range v.kick {
		for v.IO() == nil {
		}
	}
}

type BlkReq struct {
	Type   uint32
	_      uint32
	Sector uint64
}

func (v *Blk) IO() error {
	sel := uint16(0)
	// v.dumpDesc(sel)
	availRing := &v.VirtQueue[sel].AvailRing
	usedRing := &v.VirtQueue[sel].UsedRing

	if v.LastAvailIdx[sel] == availRing.Idx {
		return ErrNoTxPacket
	}

	for v.LastAvailIdx[sel] != availRing.Idx {
		descID := availRing.Ring[v.LastAvailIdx[sel]%QueueSize]

		// This structure is holding both the index of the descriptor chain and the
		// number of bytes that were written to the memory as part of serving the request.
		usedRing.Ring[usedRing.Idx%QueueSize].Idx = uint32(descID)
		usedRing.Ring[usedRing.Idx%QueueSize].Len = 0

		var buf [3][]byte

		for i := 0; i < 3; i++ {
			desc := v.VirtQueue[sel].DescTable[descID]
			buf[i] = v.Mem[desc.Addr : desc.Addr+uint64(desc.Len)]

			usedRing.Ring[usedRing.Idx%QueueSize].Len += desc.Len
			descID = desc.Next
		}

		// buf[0] contains type, reserved, and sector fields.
		// buf[1] contains raw io data.
		// buf[2] contains a status field.
		//
		// refs https://wiki.osdev.org/Virtio#Block_Device_Packets
		blkReq := *((*BlkReq)(unsafe.Pointer(&buf[0][0])))
		data := buf[1]

		var err error
		if blkReq.Type&0x1 == 0x1 {
			// write to file
			_, err = v.file.WriteAt(data, int64(blkReq.Sector*SectorSize))
		} else {
			// read from file
			_, err = v.file.ReadAt(data, int64(blkReq.Sector*SectorSize))
		}

		if err != nil {
			return err
		}

		if err = v.file.Sync(); err != nil {
			return err
		}

		usedRing.Idx++
		v.LastAvailIdx[sel]++
	}

	v.Hdr.commonHeader.isr = 0x1
	if err := v.IRQInjector.InjectVirtioBlkIRQ(); err != nil {
		return err
	}

	return nil
}

func (v *Blk) Write(port uint64, bytes []byte) error {
	offset := int(port - BlkIOPortStart)

	switch offset {
	case 8:
		// Queue PFN is aligned to page (4096 bytes)
		physAddr := uint32(pci.BytesToNum(bytes) * 4096)
		v.VirtQueue[v.Hdr.commonHeader.queueSEL] = (*VirtQueue)(unsafe.Pointer(&v.Mem[physAddr]))
	case 14:
		v.Hdr.commonHeader.queueSEL = uint16(pci.BytesToNum(bytes))
	case 16:
		v.Hdr.commonHeader.isr = 0x0
		v.kick <- true
	case 19:
	default:
	}

	return nil
}

func (v Blk) IOPort() uint64 {
	return BlkIOPortStart
}

func (v Blk) Size() uint64 {
	return BlkIOPortSize
}

func NewBlk(path string, irq uint8, irqInjector IRQInjector, mem []byte) (*Blk, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := uint64(fileInfo.Size())

	res := &Blk{
		Hdr: blkHdr{
			commonHeader: commonHeader{
				queueNUM: QueueSize,
				isr:      0x0,
			},
			blkHeader: blkHeader{
				capacity: fileSize / SectorSize,
			},
		},
		file:         file,
		irq:          irq,
		IRQInjector:  irqInjector,
		kick:         make(chan interface{}),
		Mem:          mem,
		VirtQueue:    [1]*VirtQueue{},
		LastAvailIdx: [1]uint16{0},
	}

	return res, nil
}
