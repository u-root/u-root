package gpt

import (
	"bytes"
	"testing"
)

type randomWriteBuffer struct {
	buf    []byte
	offset int
}

func (this *randomWriteBuffer) Seek(offset int64, whence int) (newOffset int64, err error) {
	switch whence {
	case 0:
		this.offset = int(offset)
	case 1:
		this.offset += int(offset)
	case 2:
		this.offset = len(this.buf) + int(offset)
	default:
		panic("Error whence")
	}

	if this.offset >= len(this.buf) {
		newBuf := make([]byte, this.offset)
		copy(newBuf, this.buf)
		this.buf = newBuf
	}
	return int64(this.offset), nil
}

func (this *randomWriteBuffer) Write(p []byte) (n int, err error) {
	needLen := this.offset + len(p)
	if needLen > len(this.buf) {
		newBuf := make([]byte, needLen)
		copy(newBuf, this.buf)
		this.buf = newBuf
	}
	copy(this.buf[this.offset:], p)
	this.offset += len(p)
	return len(p), nil
}

func TestHeaderRead(t *testing.T) {
	reader := bytes.NewReader(GPT_TEST_HEADER)
	h, err := readHeader(reader, 512)
	if err != nil {
		t.Errorf(err.Error())
	}
	if string(h.Signature[:]) != "EFI PART" {
		t.Error("Signature: ", string(h.Signature[:]))
	}
	if h.Revision != 0x00010000 { // v1.00 in hex
		t.Error("Revision: ", h.Revision)
	}
	if h.Size != 92 {
		t.Error("Header size: ", h.Size)
	}
	if h.CRC != h.calcCRC() {
		t.Error("CRC")
	}
	if h.Reserved != 0 {
		t.Error("Reserved")
	}
	if h.HeaderStartLBA != 1 {
		t.Error("CurrentLBA", h.HeaderStartLBA)
	}
	if h.HeaderCopyStartLBA != 1953525167 {
		t.Error("Other LBA", h.HeaderCopyStartLBA)
	}
	if h.FirstUsableLBA != 34 {
		t.Error("FirstUsable: ", h.FirstUsableLBA)
	}
	if h.LastUsableLBA != 1953525134 {
		t.Error("LastUsable: ", h.LastUsableLBA)
	}
	if !bytes.Equal(h.DiskGUID[:], []byte{190, 139, 78, 124, 58, 164, 159, 72, 142, 28, 5, 196, 90, 42, 168, 188}) {
		t.Error("Disk GUID: ", h.DiskGUID)
	}
	if h.PartitionsTableStartLBA != 2 {
		t.Error("Start partition entries: ", h.PartitionsTableStartLBA)
	}
	if h.PartitionsArrLen != 128 {
		t.Error("Partition arr len", h.PartitionsArrLen)
	}
	if h.PartitionEntrySize != 128 {
		t.Error("Partition entry size:", h.PartitionEntrySize)
	}
	if h.PartitionsCRC != 1233018821 {
		t.Error("Partitions CRC", h.PartitionsCRC)
	}
	if !bytes.Equal(h.TrailingBytes, make([]byte, 420)) {
		t.Error("Trailing bytes: ", h.TrailingBytes)
	}
}

func TestHeaderReadWrite(t *testing.T) {
	reader := bytes.NewReader(GPT_TEST_HEADER)
	h, err := readHeader(reader, 512)
	if err != nil {
		t.Errorf(err.Error())
	}
	writer := &bytes.Buffer{}
	h.write(writer, true)

	if !bytes.Equal(GPT_TEST_HEADER, writer.Bytes()) {
		t.Error("Read and write not equal")
	}
}

func TestEntryReadWrite(t *testing.T) {
	testEntry := make([]byte, 137)
	copy(testEntry, GPT_TEST_ENTRIES[0:128])
	testEntry[128] = 1
	testEntry[129] = 231
	testEntry[130] = 144
	testEntry[131] = 66
	testEntry[132] = 123
	testEntry[133] = 15
	testEntry[134] = 18
	testEntry[135] = 26
	testEntry[126] = 215

	p, err := readPartition(bytes.NewReader(testEntry), 137)
	if err != nil {
		t.Error(err)
	}

	writeBuf := &bytes.Buffer{}
	err = p.write(writeBuf, 137)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(testEntry, writeBuf.Bytes()) {
		t.Error("Read-write")
	}
}

func TestPartitionRead(t *testing.T) {
	p, err := readPartition(bytes.NewReader(GPT_TEST_ENTRIES), 128)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(p.Type[:], []byte{40, 115, 42, 193, 31, 248, 210, 17, 186, 75, 0, 160, 201, 62, 201, 59}) {
		// C12A7328-F81F-11D2-BA4B-00A0C93EC93B in little endian: first 3 parts in reverse order
		t.Error("Part type: ", p.Type)
	}
	if !bytes.Equal(p.Id[:], []byte{176, 80, 47, 220, 222, 152, 129, 70, 168, 104, 66, 233, 254, 189, 110, 62}) {
		t.Error("Partition GUID: ", p.Id)
	}
	if p.FirstLBA != 2048 {
		t.Error("First LBA: ", p.FirstLBA)
	}
	if p.LastLBA != 780287 {
		t.Error("Last LBA: ", p.LastLBA)
	}
	if !bytes.Equal(p.Flags[:], make([]byte, 8)) {
		t.Error("Flags: ", p.Flags)
	}
	if !bytes.Equal(p.PartNameUTF16[:], make([]byte, 72)) {
		t.Error("Name: ", p.PartNameUTF16)
	}
}

func TestPartitionBadWrite(t *testing.T) {
	var p Partition
	p.TrailingBytes = []byte{1, 2, 3}
	buf := &bytes.Buffer{}
	if p.write(buf, 130) == nil || p.write(buf, 132) == nil {
		t.Error("Write with bad entry size")
	}
	if p.write(buf, 131) != nil {
		t.Error("Write with ok entry size")
	}
}

func TestReadWriteTable(t *testing.T) {
	buf := make([]byte, 512+512+32*512)
	copy(buf[512:], GPT_TEST_HEADER)
	copy(buf[1024:], GPT_TEST_ENTRIES)
	reader := bytes.NewReader(buf)
	reader.Seek(512, 0)
	table, err := ReadTable(reader, 512)
	if err != nil {
		t.Error(err)
	}

	buf2 := &randomWriteBuffer{}
	table.Write(buf2)
	if !bytes.Equal(buf, buf2.buf) {
		t.Error("Bad read-write")
	}
}

func TestNewTable(t *testing.T) {
	guid := Guid{0xc5, 0x7f, 0x7e, 0x46, 0x36, 0x2b, 0x4e, 0x60, 0x9a, 0xa9, 0xa6, 0xe9, 0xdd, 0x85, 0x94, 0xa6}
	ssize := 4096
	diskSize := uint64(ssize * 1024 * 1024) // round 4TiB
	numSectors := diskSize / uint64(ssize)

	// sector size 4096 means header+ptdata fits in 5 sectors (33 for size 512)
	headerSizePlusPartData := uint64(5)
	expectedLastLBA := uint64(numSectors - headerSizePlusPartData - 1)

	table := NewTable(diskSize, &NewTableArgs{uint64(ssize), guid})
	h := table.Header
	if h.DiskGUID != guid {
		t.Errorf("found DiskGUID %v != %v", guid, h.DiskGUID)
	}

	if len(h.TrailingBytes) != ssize-standardHeaderSize {
		t.Errorf("Found TralingBytes len %d expected %d", len(h.TrailingBytes), ssize-standardHeaderSize)
	}

	if h.CRC != h.calcCRC() {
		t.Error("CRC")
	}
	if h.Reserved != 0 {
		t.Error("Reserved")
	}
	if h.HeaderStartLBA != 1 {
		t.Error("CurrentLBA", h.HeaderStartLBA)
	}
	if h.HeaderCopyStartLBA != numSectors-1 {
		t.Error("Other LBA", h.HeaderCopyStartLBA, numSectors-1)
	}
	if h.FirstUsableLBA != uint64(headerSizePlusPartData+1) {
		t.Error("FirstUsable: ", h.FirstUsableLBA)
	}
	if h.LastUsableLBA != expectedLastLBA {
		t.Errorf("LastUsable: expected %d found %d", expectedLastLBA, h.LastUsableLBA)
	}
	if h.PartitionsTableStartLBA != 2 {
		t.Error("Start partition entries: ", h.PartitionsTableStartLBA)
	}
	if h.PartitionsArrLen != 128 {
		t.Error("Partition arr len", h.PartitionsArrLen)
	}
	if h.PartitionEntrySize != 128 {
		t.Error("Partition entry size:", h.PartitionEntrySize)
	}
	if h.PartitionsCRC != 0 {
		t.Error("Partitions CRC", h.PartitionsCRC)
	}

}

func TestNewTableNil(t *testing.T) {
	emptyGuid := Guid{}
	ssize := 512
	diskSize := uint64(250000384) // round 512 sectors closest to 250gb
	numSectors := diskSize / uint64(ssize)

	// sector size 512 means header+ptdata fits in 33 sectors.
	headerSizePlusPartData := uint64(33)
	expectedLastLBA := uint64(numSectors - headerSizePlusPartData - 1)

	table := NewTable(diskSize, nil)

	h := table.Header
	if h.DiskGUID == emptyGuid {
		t.Errorf("DiskGUID did not get generated. still empty")
	}

	if len(h.TrailingBytes) != ssize-standardHeaderSize {
		t.Errorf("Found TralingBytes len %d expected %d", len(h.TrailingBytes), ssize-standardHeaderSize)
	}

	if h.CRC != h.calcCRC() {
		t.Error("CRC")
	}
	if h.Reserved != 0 {
		t.Error("Reserved")
	}
	if h.HeaderStartLBA != 1 {
		t.Error("CurrentLBA", h.HeaderStartLBA)
	}
	if h.HeaderCopyStartLBA != numSectors-1 {
		t.Error("Other LBA", h.HeaderCopyStartLBA, numSectors-1)
	}
	if h.FirstUsableLBA != headerSizePlusPartData+1 {
		t.Error("FirstUsable: ", h.FirstUsableLBA)
	}
	if h.LastUsableLBA != expectedLastLBA {
		t.Errorf("LastUsable: expected %d found %d", expectedLastLBA, h.LastUsableLBA)
	}
	if h.PartitionsTableStartLBA != 2 {
		t.Error("Start partition entries: ", h.PartitionsTableStartLBA)
	}
	if h.PartitionsArrLen != 128 {
		t.Error("Partition arr len", h.PartitionsArrLen)
	}
	if h.PartitionEntrySize != 128 {
		t.Error("Partition entry size:", h.PartitionEntrySize)
	}
	if h.PartitionsCRC != 0 {
		t.Error("Partitions CRC", h.PartitionsCRC)
	}
}

func TestTableCreateOtherSide(t *testing.T) {
	buf := make([]byte, 512+512+32*512)
	copy(buf[512:], GPT_TEST_HEADER)
	copy(buf[1024:], GPT_TEST_ENTRIES)
	reader := bytes.NewReader(buf)
	reader.Seek(512, 0)
	t1, err := ReadTable(reader, 512)
	if err != nil {
		t.Error(err)
	}

	t2 := t1.CreateOtherSideTable()
	if t2.Header.CRC != t2.Header.calcCRC() {
		t.Error("crc")
	}
	if t2.Header.Signature != t1.Header.Signature {
		t.Error("signature")
	}
	if t2.Header.Revision != t1.Header.Revision {
		t.Error("revision")
	}
	if t2.Header.Size != t1.Header.Size {
		t.Error("size")
	}
	if t2.Header.Reserved != t1.Header.Reserved {
		t.Error("reserved")
	}
	if t2.Header.HeaderStartLBA != t1.Header.HeaderCopyStartLBA {
		t.Error("header start")
	}
	if t2.Header.HeaderCopyStartLBA != t1.Header.HeaderStartLBA {
		t.Error("header copy")
	}
	if t2.Header.FirstUsableLBA != t1.Header.FirstUsableLBA {
		t.Error("first usable")
	}
	if t2.Header.LastUsableLBA != t1.Header.LastUsableLBA {
		t.Error("last usable")
	}
	if t2.Header.DiskGUID != t1.Header.DiskGUID {
		t.Error("disk guid")
	}
	if t2.Header.PartitionsTableStartLBA != t2.Header.LastUsableLBA+1 {
		t.Error("partitions table start")
	}
	if t2.Header.PartitionsArrLen != t1.Header.PartitionsArrLen {
		t.Error("partitions table len")
	}
	if t2.Header.PartitionEntrySize != t1.Header.PartitionEntrySize {
		t.Error("partitions entry")
	}
	if t2.Header.PartitionsCRC != t1.Header.PartitionsCRC {
		t.Error("partitions crc")
	}
	if !bytes.Equal(t2.Header.TrailingBytes, t1.Header.TrailingBytes) {
		t.Error("trailing bytes aren't equal")
	}

	if len(t1.Partitions) != len(t2.Partitions) {
		t.Fatal("partitions len are different")
	}

	for i := range t1.Partitions {
		l := &t1.Partitions[i]
		r := &t2.Partitions[i]
		if l.Type != r.Type {
			t.Error("part type")
		}
		if l.Id != r.Id {
			t.Error("part id")
		}
		if l.FirstLBA != r.FirstLBA {
			t.Error("part first lba")
		}
		if l.LastLBA != r.LastLBA {
			t.Error("part last lba")
		}
		if l.Flags != r.Flags {
			t.Error("part flags")
		}
		if l.PartNameUTF16 != r.PartNameUTF16 {
			t.Error("part name")
		}

		if !bytes.Equal(l.TrailingBytes, r.TrailingBytes) {
			t.Error("part trailing are different")
		}
	}
}

func TestTableCopy(t *testing.T) {
	buf := make([]byte, 512+512+32*512)
	copy(buf[512:], GPT_TEST_HEADER)
	copy(buf[1024:], GPT_TEST_ENTRIES)
	reader := bytes.NewReader(buf)
	reader.Seek(512, 0)
	t1, err := ReadTable(reader, 512)
	if err != nil {
		t.Error(err)
	}

	t1.Header.TrailingBytes = append(t1.Header.TrailingBytes, 1)
	t1.Partitions[0].TrailingBytes = append(t1.Partitions[0].TrailingBytes, 1)

	t2 := t1.copy()
	if t2.Header.CRC != t2.Header.calcCRC() {
		t.Error("crc")
	}
	if t2.Header.Signature != t1.Header.Signature {
		t.Error("signature")
	}
	if t2.Header.Revision != t1.Header.Revision {
		t.Error("revision")
	}
	if t2.Header.Size != t1.Header.Size {
		t.Error("size")
	}
	if t2.Header.Reserved != t1.Header.Reserved {
		t.Error("reserved")
	}
	if t2.Header.HeaderStartLBA != t1.Header.HeaderStartLBA {
		t.Error("header start")
	}
	if t2.Header.HeaderCopyStartLBA != t1.Header.HeaderCopyStartLBA {
		t.Error("header copy")
	}
	if t2.Header.FirstUsableLBA != t1.Header.FirstUsableLBA {
		t.Error("first usable")
	}
	if t2.Header.LastUsableLBA != t1.Header.LastUsableLBA {
		t.Error("last usable")
	}
	if t2.Header.DiskGUID != t1.Header.DiskGUID {
		t.Error("disk guid")
	}
	if t2.Header.PartitionsTableStartLBA != t1.Header.PartitionsTableStartLBA {
		t.Error("partitions table start")
	}
	if t2.Header.PartitionsArrLen != t1.Header.PartitionsArrLen {
		t.Error("partitions table len")
	}
	if t2.Header.PartitionEntrySize != t1.Header.PartitionEntrySize {
		t.Error("partitions entry")
	}
	if t2.Header.PartitionsCRC != t1.Header.PartitionsCRC {
		t.Error("partitions crc")
	}
	if !bytes.Equal(t2.Header.TrailingBytes, t1.Header.TrailingBytes) {
		t.Error("trailing bytes aren't equal")
	}
	t1.Header.TrailingBytes[0] += 1
	if t1.Header.TrailingBytes[0] == t2.Header.TrailingBytes[0] {
		t.Error("same trailing bytes")
	}

	if len(t1.Partitions) != len(t2.Partitions) {
		t.Fatal("partitions len are different")
	}

	for i := range t1.Partitions {
		l := &t1.Partitions[i]
		r := &t2.Partitions[i]
		if l.Type != r.Type {
			t.Error("part type")
		}
		if l.Id != r.Id {
			t.Error("part id")
		}
		if l.FirstLBA != r.FirstLBA {
			t.Error("part first lba")
		}
		if l.LastLBA != r.LastLBA {
			t.Error("part last lba")
		}
		if l.Flags != r.Flags {
			t.Error("part flags")
		}
		if l.PartNameUTF16 != r.PartNameUTF16 {
			t.Error("part name")
		}

		if !bytes.Equal(l.TrailingBytes, r.TrailingBytes) {
			t.Error("part trailing are different")
		}
	}
	t1.Partitions[0].TrailingBytes[0] += 1
	if t1.Partitions[0].TrailingBytes[0] == t2.Partitions[0].TrailingBytes[0] {
		t.Error("same partitions")
	}
}

func TestTableNewSize(t *testing.T) {
	buf := make([]byte, 512+512+32*512)
	copy(buf[512:], GPT_TEST_HEADER)
	copy(buf[1024:], GPT_TEST_ENTRIES)
	reader := bytes.NewReader(buf)
	reader.Seek(512, 0)
	t1, err := ReadTable(reader, 512)
	if err != nil {
		t.Error(err)
	}

	t1.Header.TrailingBytes = append(t1.Header.TrailingBytes, 1)
	t1.Partitions[0].TrailingBytes = append(t1.Partitions[0].TrailingBytes, 1)

	t2 := t1.CreateTableForNewDiskSize(1953525168) // Same size as original disk
	if t2.Header.CRC != t2.Header.calcCRC() {
		t.Error("crc")
	}
	if t2.Header.Signature != t1.Header.Signature {
		t.Error("signature")
	}
	if t2.Header.Revision != t1.Header.Revision {
		t.Error("revision")
	}
	if t2.Header.Size != t1.Header.Size {
		t.Error("size")
	}
	if t2.Header.Reserved != t1.Header.Reserved {
		t.Error("reserved")
	}
	if t2.Header.HeaderStartLBA != t1.Header.HeaderStartLBA {
		t.Error("header start")
	}
	if t2.Header.HeaderCopyStartLBA != t1.Header.HeaderCopyStartLBA {
		t.Error("header copy")
	}
	if t2.Header.FirstUsableLBA != t1.Header.FirstUsableLBA {
		t.Error("first usable")
	}
	if t2.Header.LastUsableLBA != t1.Header.LastUsableLBA {
		t.Error("last usable")
	}
	if t2.Header.DiskGUID != t1.Header.DiskGUID {
		t.Error("disk guid")
	}
	if t2.Header.PartitionsTableStartLBA != t2.Header.PartitionsTableStartLBA {
		t.Error("partitions table start")
	}
	if t2.Header.PartitionsArrLen != t1.Header.PartitionsArrLen {
		t.Error("partitions table len")
	}
	if t2.Header.PartitionEntrySize != t1.Header.PartitionEntrySize {
		t.Error("partitions entry")
	}
	if t2.Header.PartitionsCRC != t1.Header.PartitionsCRC {
		t.Error("partitions crc")
	}
	if !bytes.Equal(t2.Header.TrailingBytes, t1.Header.TrailingBytes) {
		t.Error("trailing bytes aren't equal")
	}
	t1.Header.TrailingBytes[0] += 1
	if t1.Header.TrailingBytes[0] == t2.Header.TrailingBytes[0] {
		t.Error("same trailing bytes")
	}

	if len(t1.Partitions) != len(t2.Partitions) {
		t.Fatal("partitions len are different")
	}

	for i := range t1.Partitions {
		l := &t1.Partitions[i]
		r := &t2.Partitions[i]
		if l.Type != r.Type {
			t.Error("part type")
		}
		if l.Id != r.Id {
			t.Error("part id")
		}
		if l.FirstLBA != r.FirstLBA {
			t.Error("part first lba")
		}
		if l.LastLBA != r.LastLBA {
			t.Error("part last lba")
		}
		if l.Flags != r.Flags {
			t.Error("part flags")
		}
		if l.PartNameUTF16 != r.PartNameUTF16 {
			t.Error("part name")
		}

		if !bytes.Equal(l.TrailingBytes, r.TrailingBytes) {
			t.Error("part trailing are different")
		}
	}
	t1.Partitions[0].TrailingBytes[0] += 1
	if t1.Partitions[0].TrailingBytes[0] == t2.Partitions[0].TrailingBytes[0] {
		t.Error("same partitions")
	}

	// Check about it is cacled state, not copy
	t2 = t1.CreateTableForNewDiskSize(100)
	if t2.Header.HeaderCopyStartLBA != 99 {
		t.Error("t2 calc header copy: ", t2.Header.HeaderCopyStartLBA)
	}
	if t2.Header.LastUsableLBA != 66 {
		t.Error("t2 calc last lba: ", t2.Header.LastUsableLBA)
	}
}

func TestGuidToString(t *testing.T) {
	guid := [...]byte{40, 115, 42, 193, 31, 248, 210, 17, 186, 75, 0, 160, 201, 62, 201, 59}
	guidS := guidToString(guid)
	if guidS != "C12A7328-F81F-11D2-BA4B-00A0C93EC93B" {
		t.Errorf("Error guid: %v != %v", guidS, "C12A7328-F81F-11D2-BA4B-00A0C93EC93B")
	}
}

func TestStringToGuid(t *testing.T) {
	guid, err := StringToGuid("C12A7328-F81F-11D2-BA4B-00A0C93EC93B")
	if err != nil {
		t.Error(err)
	}
	if guid != [...]byte{40, 115, 42, 193, 31, 248, 210, 17, 186, 75, 0, 160, 201, 62, 201, 59} {
		t.Errorf("Bad result. Expected:\n%v\nResult:\n%v", [...]byte{40, 115, 42, 193, 31, 248, 210, 17, 186, 75, 0, 160, 201, 62, 201, 59}, guid)
	}
	if _, err := StringToGuid(""); err == nil {
		t.Error("Must return error")
	}
	if _, err := StringToGuid("C12A7328-F81F-11D2-BA4B-00A0C93EC93BA"); err == nil {
		t.Error("Must return error")
	}
	if _, err := StringToGuid("C12A7328-F81F-11D2-BA4B!00A0C93EC93B"); err == nil {
		t.Error("Must return error")
	}
	if _, err := StringToGuid("C12A7328-F81F-11D2-BA4B-00A0C93EC93Z"); err == nil {
		t.Error("Must return error")
	}
}
