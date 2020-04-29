package gpt

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"unicode/utf16"
)

const standardHeaderSize = 92          // Size of standard GPT-header in bytes
const standardPartitionEntrySize = 128 // Size of standard GPT-partition entry in bytes

type Flags [8]byte
type Guid [16]byte

func (this Guid) String() string {
	return guidToString(this)
}

type PartType Guid

func (this PartType) String() string {
	return guidToString(this)
}

var (
	GUID_LVM = PartType([16]byte{0x79, 0xd3, 0xd6, 0xe6, 0x7, 0xf5, 0xc2, 0x44, 0xa2, 0x3c, 0x23, 0x8f, 0x2a, 0x3d, 0xf9, 0x28}) // E6D6D379-F507-44C2-A23C-238F2A3DF928
)

// https://en.wikipedia.org/wiki/GUID_Partition_Table#Partition_table_header_.28LBA_1.29
type Header struct {
	Signature               [8]byte // Offset  0. "EFI PART", 45h 46h 49h 20h 50h 41h 52h 54h
	Revision                uint32  // Offset  8
	Size                    uint32  // Offset 12
	CRC                     uint32  // Offset 16. Autocalc when save Header.
	Reserved                uint32  // Offset 20
	HeaderStartLBA          uint64  // Offset 24
	HeaderCopyStartLBA      uint64  // Offset 32
	FirstUsableLBA          uint64  // Offset 40
	LastUsableLBA           uint64  // Offset 48
	DiskGUID                Guid    // Offset 56
	PartitionsTableStartLBA uint64  // Offset 72
	PartitionsArrLen        uint32  // Offset 80
	PartitionEntrySize      uint32  // Offset 84
	PartitionsCRC           uint32  // Offset 88. Autocalc when save Table.
	TrailingBytes           []byte  // Offset 92
}

// https://en.wikipedia.org/wiki/GUID_Partition_Table#Partition_entries
type Partition struct {
	Type          PartType // Offset 0
	Id            Guid     // Offset 16
	FirstLBA      uint64   // Offset 32
	LastLBA       uint64   // Offset 40
	Flags         Flags    // Offset 68
	PartNameUTF16 [72]byte // Offset 56
	TrailingBytes []byte   // Offset 128. Usually it is empty
}

type Table struct {
	SectorSize uint64 // in bytes
	Header     Header
	Partitions []Partition
}

//////////////////////////////////////////////
////////////////// HEADER ////////////////////
//////////////////////////////////////////////

// Have to set to start of Header. Usually LBA1 for primary header.
func readHeader(reader io.Reader, sectorSize uint64) (res Header, err error) {
	read := func(data interface{}) {
		if err == nil {
			err = binary.Read(reader, binary.LittleEndian, data)
		}
	}

	read(&res.Signature)
	read(&res.Revision)
	read(&res.Size)
	read(&res.CRC)
	read(&res.Reserved)
	read(&res.HeaderStartLBA)
	read(&res.HeaderCopyStartLBA)
	read(&res.FirstUsableLBA)
	read(&res.LastUsableLBA)
	read(&res.DiskGUID)
	read(&res.PartitionsTableStartLBA)
	read(&res.PartitionsArrLen)
	read(&res.PartitionEntrySize)
	read(&res.PartitionsCRC)
	if err != nil {
		return
	}

	if string(res.Signature[:]) != "EFI PART" {
		return res, fmt.Errorf("Bad GPT signature")
	}
	trailingBytes := make([]byte, sectorSize-uint64(standardHeaderSize))
	reader.Read(trailingBytes)
	res.TrailingBytes = trailingBytes

	if res.calcCRC() != res.CRC {
		return res, fmt.Errorf("BAD GPT Header CRC")
	}

	return
}

func (this *Header) calcCRC() uint32 {
	buf := &bytes.Buffer{}
	this.write(buf, false)
	return crc32.ChecksumIEEE(buf.Bytes()[:this.Size])
}

func (this *Header) write(writer io.Writer, saveCRC bool) (err error) {
	write := func(data interface{}) {
		if err == nil {
			err = binary.Write(writer, binary.LittleEndian, data)
		}
	}

	write(&this.Signature)
	write(&this.Revision)
	write(&this.Size)

	if saveCRC {
		this.CRC = this.calcCRC()
		write(&this.CRC)
	} else {
		write(uint32(0))
	}

	write(&this.Reserved)
	write(&this.HeaderStartLBA)
	write(&this.HeaderCopyStartLBA)
	write(&this.FirstUsableLBA)
	write(&this.LastUsableLBA)
	write(&this.DiskGUID)
	write(&this.PartitionsTableStartLBA)
	write(&this.PartitionsArrLen)
	write(&this.PartitionEntrySize)
	write(&this.PartitionsCRC)
	if err != nil {
		return
	}
	write(this.TrailingBytes)
	return
}

//////////////////////////////////////////////
///////////////// PARTITION //////////////////
//////////////////////////////////////////////
func readPartition(reader io.Reader, size uint32) (p Partition, err error) {
	read := func(data interface{}) {
		if err == nil {
			err = binary.Read(reader, binary.LittleEndian, data)
		}
	}

	p.TrailingBytes = make([]byte, size-standardPartitionEntrySize)

	read(&p.Type)
	read(&p.Id)
	read(&p.FirstLBA)
	read(&p.LastLBA)
	read(&p.Flags)
	read(&p.PartNameUTF16)
	read(&p.TrailingBytes)

	return
}

func (this Partition) write(writer io.Writer, size uint32) (err error) {
	write := func(data interface{}) {
		if err == nil {
			err = binary.Write(writer, binary.LittleEndian, data)
		}
	}

	if size != uint32(standardPartitionEntrySize+len(this.TrailingBytes)) {
		return fmt.Errorf("Entry size(%v) != real entry size(%v)", size, standardPartitionEntrySize+len(this.TrailingBytes))
	}

	write(this.Type)
	write(this.Id)
	write(this.FirstLBA)
	write(this.LastLBA)
	write(this.Flags)
	write(this.PartNameUTF16)
	write(this.TrailingBytes)

	return
}

func (this Partition) IsEmpty() bool {
	return this.Type == [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
}

func (this Partition) Name() string {
	chars := make([]uint16, 0, 36)
	for i := 0; i < len(this.PartNameUTF16); i += 2 {
		byte1 := this.PartNameUTF16[i]
		byte2 := this.PartNameUTF16[i+1]
		if byte1 == 0 && byte2 == 0 {
			break
		}
		chars = append(chars, uint16(byte1)+uint16(byte2)<<8)
	}
	runes := utf16.Decode(chars)
	return string(runes)
}

//////////////////////////////////////////////
////////////////// TABLE /////////////////////
//////////////////////////////////////////////

// Read GPT partition
// Have to set to first byte of GPT Header (usually start of second sector on disk)
func ReadTable(reader io.ReadSeeker, SectorSize uint64) (table Table, err error) {
	table.SectorSize = SectorSize
	table.Header, err = readHeader(reader, SectorSize)
	if err != nil {
		return
	}
	if seekDest, ok := mul(int64(SectorSize), int64(table.Header.PartitionsTableStartLBA)); ok {
		reader.Seek(seekDest, 0)
	} else {
		err = fmt.Errorf("Seek overflow when read partition tables")
		return
	}
	for i := uint32(0); i < table.Header.PartitionsArrLen; i++ {
		var p Partition
		p, err = readPartition(reader, table.Header.PartitionEntrySize)
		if err != nil {
			return
		}
		table.Partitions = append(table.Partitions, p)
	}

	if table.Header.PartitionsCRC != table.calcPartitionsCRC() {
		err = fmt.Errorf("Bad partitions crc")
		return
	}
	return
}

func (this Table) CreateOtherSideTable() (res Table) {
	res = this.copy()

	// Copy of table on other side of disk
	tmpDest := res.Header.HeaderStartLBA
	res.Header.HeaderStartLBA = res.Header.HeaderCopyStartLBA
	res.Header.HeaderCopyStartLBA = tmpDest

	if res.Header.HeaderStartLBA == 1 {
		res.Header.PartitionsTableStartLBA = 2
	} else {
		// Partitions table on other side of disk
		res.Header.PartitionsTableStartLBA = res.Header.LastUsableLBA + 1
	}

	res.Header.CRC = res.Header.calcCRC()
	return res
}

// Create primary table for resized disk
// size - in sectors
func (this Table) CreateTableForNewDiskSize(size uint64) (res Table) {
	res = this.copy()

	// Always create primary table
	res.Header.HeaderStartLBA = 1
	res.Header.PartitionsTableStartLBA = 2
	res.Header.HeaderCopyStartLBA = size - 1 // Last sector

	partitionsTableSize := uint64(res.Header.PartitionEntrySize) * uint64(res.Header.PartitionsArrLen)
	partitionSizeInSector := partitionsTableSize / uint64(res.SectorSize)
	if partitionsTableSize%uint64(res.SectorSize) != 0 {
		partitionSizeInSector++
	}
	res.Header.LastUsableLBA = size - 1 - partitionSizeInSector - 1 // header in last sector and partitions table

	res.Header.CRC = res.Header.calcCRC()
	return res
}

func (this Table) copy() (res Table) {
	res = this

	res.Header.TrailingBytes = make([]byte, len(this.Header.TrailingBytes))
	copy(res.Header.TrailingBytes, this.Header.TrailingBytes)

	res.Partitions = make([]Partition, len(this.Partitions))
	copy(res.Partitions, this.Partitions)
	for i := range this.Partitions {
		res.Partitions[i].TrailingBytes = make([]byte, len(this.Partitions[i].TrailingBytes))
		copy(res.Partitions[i].TrailingBytes, this.Partitions[i].TrailingBytes)
	}

	return res
}

func (this Table) calcPartitionsCRC() uint32 {
	buf := &bytes.Buffer{}
	for _, part := range this.Partitions {
		part.write(buf, this.Header.PartitionEntrySize)
	}
	return crc32.ChecksumIEEE(buf.Bytes())
}

// Calc header and partitions CRC. Save Header and partition entries to the disk.
// It independent of start position: writer will be seek to position from Table.Header.
func (this Table) Write(writer io.WriteSeeker) (err error) {
	this.Header.PartitionsCRC = this.calcPartitionsCRC()
	if headerPos, ok := mul(int64(this.SectorSize), int64(this.Header.HeaderStartLBA)); ok {
		writer.Seek(headerPos, 0)
	}
	err = this.Header.write(writer, true)
	if err != nil {
		return
	}
	if partTablePos, ok := mul(int64(this.SectorSize), int64(this.Header.PartitionsTableStartLBA)); ok {
		writer.Seek(partTablePos, 0)
	}
	for _, part := range this.Partitions {
		err = part.write(writer, this.Header.PartitionEntrySize)
		if err != nil {
			return
		}
	}
	return
}

// Use for create guid predefined values in snippet http://play.golang.org/p/uOd_WQtiwE
func StringToGuid(guid string) (res [16]byte, err error) {
	byteOrder := [...]int{3, 2, 1, 0, -1, 5, 4, -1, 7, 6, -1, 8, 9, -1, 10, 11, 12, 13, 14, 15}
	if len(guid) != 36 {
		err = fmt.Errorf("BAD guid string length.")
		return
	}
	guidByteNum := 0
	for i := 0; i < len(guid); i += 2 {
		if byteOrder[guidByteNum] == -1 {
			if guid[i] == '-' {
				i++
				guidByteNum++
				if i >= len(guid)+1 {
					err = fmt.Errorf("BAD guid format minus")
					return
				}
			} else {
				err = fmt.Errorf("BAD guid char in minus pos")
				return
			}
		}

		sub := guid[i : i+2]
		var bt byte
		for pos, ch := range sub {
			var shift uint
			if pos == 0 {
				shift = 4
			} else {
				shift = 0
			}
			switch ch {
			case '0':
				bt |= 0 << shift
			case '1':
				bt |= 1 << shift
			case '2':
				bt |= 2 << shift
			case '3':
				bt |= 3 << shift
			case '4':
				bt |= 4 << shift
			case '5':
				bt |= 5 << shift
			case '6':
				bt |= 6 << shift
			case '7':
				bt |= 7 << shift
			case '8':
				bt |= 8 << shift
			case '9':
				bt |= 9 << shift
			case 'A', 'a':
				bt |= 10 << shift
			case 'B', 'b':
				bt |= 11 << shift
			case 'C', 'c':
				bt |= 12 << shift
			case 'D', 'd':
				bt |= 13 << shift
			case 'E', 'e':
				bt |= 14 << shift
			case 'F', 'f':
				bt |= 15 << shift
			default:
				err = fmt.Errorf("BAD guid char at pos %d: '%c'", i+pos, ch)
				return
			}
		}
		res[byteOrder[guidByteNum]] = bt
		guidByteNum++
	}
	return res, nil
}

// NewTableArgs - arguments NewTable creation.
type NewTableArgs struct {
	SectorSize uint64
	DiskGuid   Guid
}

// NewTable - return a valid empty Table for given sectorSize and diskSize
//    Note that a Protective MBR is needed for lots of software to read the GPT table.
func NewTable(diskSize uint64, args *NewTableArgs) Table {
	// CreateTableForNewdiskSize will update HeaderCopyStartLBA, LastUsableLBA, and CRC
	if args == nil {
		args = &NewTableArgs{}
	}
	if args.SectorSize == 0 {
		args.SectorSize = uint64(512)
	}
	var emptyGuid Guid
	if args.DiskGuid == emptyGuid {
		args.DiskGuid = NewGUID()
	}

	ptStartLBA := uint64(2)
	numParts := 128
	partitionsTableSize := uint64(standardPartitionEntrySize) * uint64(numParts)
	partitionSizeInSector := partitionsTableSize / uint64(args.SectorSize)
	if partitionsTableSize%uint64(args.SectorSize) != 0 {
		partitionSizeInSector++
	}

	return Table{
		SectorSize: args.SectorSize,
		Header: Header{
			Signature:               [8]byte{0x45, 0x46, 0x49, 0x20, 0x50, 0x41, 0x52, 0x54},
			Revision:                0x10000,
			Size:                    standardHeaderSize,
			CRC:                     0,
			Reserved:                0,
			HeaderStartLBA:          1,
			HeaderCopyStartLBA:      0,
			FirstUsableLBA:          ptStartLBA + partitionSizeInSector,
			LastUsableLBA:           0,
			DiskGUID:                args.DiskGuid,
			PartitionsTableStartLBA: ptStartLBA,
			PartitionsArrLen:        uint32(numParts),
			PartitionEntrySize:      uint32(standardPartitionEntrySize),
			PartitionsCRC:           0x0,
			TrailingBytes:           make([]byte, args.SectorSize-uint64(standardHeaderSize)),
		},
		Partitions: make([]Partition, numParts),
	}.CreateTableForNewDiskSize(diskSize / args.SectorSize)
}

//////////////////////////////////////////////
//////////////// INTERNALS ///////////////////
//////////////////////////////////////////////

// Multiply two int64 numbers with overflow check
// Algorithm from https://gist.github.com/areed/85d3614a58400e417027
func mul(a, b int64) (res int64, ok bool) {
	const mostPositive = 1<<63 - 1
	const mostNegative = -(mostPositive + 1)

	if a == 0 || b == 0 || a == 1 || b == 1 {
		return a * b, true
	}
	if a == mostNegative || b == mostNegative {
		return a * b, false
	}
	c := a * b
	return c, c/b == a
}

func guidToString(byteGuid [16]byte) string {
	byteToChars := func(b byte) (res []byte) {
		res = make([]byte, 0, 2)
		for i := 1; i >= 0; i-- {
			switch b >> uint(4*i) & 0x0F {
			case 0:
				res = append(res, '0')
			case 1:
				res = append(res, '1')
			case 2:
				res = append(res, '2')
			case 3:
				res = append(res, '3')
			case 4:
				res = append(res, '4')
			case 5:
				res = append(res, '5')
			case 6:
				res = append(res, '6')
			case 7:
				res = append(res, '7')
			case 8:
				res = append(res, '8')
			case 9:
				res = append(res, '9')
			case 10:
				res = append(res, 'A')
			case 11:
				res = append(res, 'B')
			case 12:
				res = append(res, 'C')
			case 13:
				res = append(res, 'D')
			case 14:
				res = append(res, 'E')
			case 15:
				res = append(res, 'F')
			}
		}
		return
	}
	s := make([]byte, 0, 36)
	byteOrder := [...]int{3, 2, 1, 0, -1, 5, 4, -1, 7, 6, -1, 8, 9, -1, 10, 11, 12, 13, 14, 15}
	for _, i := range byteOrder {
		if i == -1 {
			s = append(s, '-')
		} else {
			s = append(s, byteToChars(byteGuid[i])...)
		}
	}
	return string(s)
}

func NewGUID() Guid {
	var res Guid
	_, err := rand.Read(res[:])
	if err != nil {
		panic(err)
	}

	// set predefined bits for UUIDv4
	res[6] = (res[6] & 0x0f) | 0x40 // Version 4
	res[8] = (res[8] & 0x3f) | 0x80 // Variant 10
	return res
}
