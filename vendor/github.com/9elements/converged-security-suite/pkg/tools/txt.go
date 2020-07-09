package tools

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/9elements/converged-security-suite/pkg/hwapi"
)

const (
	txtPublicSpace       = 0xFED30000
	txtEsts              = 0x8
	txtErrorCode         = 0x30
	txtBootStatus        = 0xa0
	txtVerFSBIF          = 0x100
	txtDIDVID            = 0x110
	txtVerQPIFF          = 0x200
	txtsInitBase         = 0x270
	txtsInitSize         = 0x278
	txtMLEJoin           = 0x290
	txtHeapBase          = 0x300
	txtHeapSize          = 0x308
	txtACMStatus         = 0x328
	txtDMAProtectedRange = 0x330
	txtPublicKey         = 0x400
	txtE2STS             = 0x8f0
)

//TXTStatus represents serveral configurations within the TXT config space
type TXTStatus struct {
	SenterDone bool // SENTER.DONE.STS (0)
	SexitDone  bool // SEXIT.DONE.STS (1)
	// Reserved (2-5)
	MemConfigLock bool // MEM-CONFIG-LOCK (6)
	PrivateOpen   bool // PRIVATE-OPEN.STS (7)
	// Reserved (8-14)
	Locality1Open bool // TXT.LOCALITY1.OPEN.STS (15)
	Locality2Open bool // TXT.LOCALITY1.OPEN.STS (16)
	// Reserved (17-63)
}

//TXTErrorCode holds the decoded ACM error code read from TXT config space
type TXTErrorCode struct {
	ModuleType        uint8 // 0: BIOS ACM, 1: Intel TXT
	ClassCode         uint8
	MajorErrorCode    uint8
	SoftwareSource    bool // 0: ACM, 1: MLE
	MinorErrorCode    uint16
	Type1Reserved     uint8
	ProcessorSoftware bool
	ValidInvalid      bool
}

//TXTRegisterSpace holds the decoded TXT config space
type TXTRegisterSpace struct {
	Sts          TXTStatus    // TXT.STS (0x0)
	TxtReset     bool         // TXT.ESTS (0x8)
	ErrorCode    TXTErrorCode // TXT.ERRORCODE
	ErrorCodeRaw uint32
	BootStatus   uint64                  // TXT.BOOTSTATUS
	FsbIf        uint32                  // TXT.VER.FSBIF
	Vid          uint16                  // TXT.DIDVID.VID
	Did          uint16                  // TXT.DIDVID.DID
	Rid          uint16                  // TXT.DIDVID.RID
	IDExt        uint16                  // TXT.DIDVID.ID-EXT
	QpiIf        uint32                  // TXT.VER.QPIIF
	SinitBase    uint32                  // TXT.SINIT.BASE
	SinitSize    uint32                  // TXT.SINIT.SIZE
	MleJoin      uint32                  // TXT.MLE.JOIN
	HeapBase     uint32                  // TXT.HEAP.BASE
	HeapSize     uint32                  // TXT.HEAP.SIZE
	Dpr          hwapi.DMAProtectedRange // TXT.DPR
	PublicKey    [4]uint64               // TXT.PUBLIC.KEY
	E2Sts        uint64                  // TXT.E2STS
}

//ACMStatus holds the decoded ACM run state
type ACMStatus struct {
	Valid          bool
	MinorErrorCode uint16
	ACMStarted     bool
	MajorErrorCode uint8
	ClassCode      uint8
	ModuleType     uint8
}

//TXTBiosData holds the decoded BIOSDATA regions as read from TXT config space
type TXTBiosData struct {
	Version       uint32
	BiosSinitSize uint32
	Reserved1     uint64
	Reserved2     uint64
	NumLogProcs   uint32
	SinitFlags    uint32
	MleFlags      *TXTBiosMLEFlags
}

//TXTBiosMLEFlags holds the decoded BIOSDATA region MLE flags as read from TXT config space
type TXTBiosMLEFlags struct {
	SupportsACPIPPI bool
	IsLegacyState   bool
	IsServerState   bool
	IsClientState   bool
}

//FetchTXTRegs returns a raw copy of the TXT config space
func FetchTXTRegs(txtAPI hwapi.APIInterfaces) ([]byte, error) {
	data := make([]byte, 0x1000)
	if err := txtAPI.ReadPhysBuf(txtPublicSpace, data); err != nil {
		return nil, err
	}
	return data, nil
}

//ParseTXTRegs decodes a raw copy of the TXT config space
func ParseTXTRegs(data []byte) (TXTRegisterSpace, error) {
	var regSpace TXTRegisterSpace
	var err error

	regSpace.Sts, err = readTXTStatus(data)
	if err != nil {
		return regSpace, err

	}

	regSpace.ErrorCode, regSpace.ErrorCodeRaw, err = readTXTErrorCode(data)
	if err != nil {
		return regSpace, err

	}

	regSpace.Dpr, err = readDMAProtectedRange(data)
	if err != nil {
		return regSpace, err

	}

	// TXT.ESTS (0x8)
	buf := bytes.NewReader(data)
	buf.Seek(int64(txtEsts), io.SeekStart)
	err = binary.Read(buf, binary.LittleEndian, &regSpace.TxtReset)

	if err != nil {
		return regSpace, err
	}

	// TXT.BootSTATUS (0xa0)
	buf.Seek(int64(txtBootStatus), io.SeekStart)
	err = binary.Read(buf, binary.LittleEndian, &regSpace.BootStatus)
	if err != nil {
		return regSpace, err
	}

	// TXT.VER.FSBIF
	buf.Seek(int64(txtVerFSBIF), io.SeekStart)
	err = binary.Read(buf, binary.LittleEndian, &regSpace.FsbIf)
	if err != nil {
		return regSpace, err
	}

	// TXT.DIDVID
	buf.Seek(int64(txtDIDVID), io.SeekStart)
	err = binary.Read(buf, binary.LittleEndian, &regSpace.Vid)
	if err != nil {
		return regSpace, err
	}
	err = binary.Read(buf, binary.LittleEndian, &regSpace.Did)
	if err != nil {
		return regSpace, err
	}
	err = binary.Read(buf, binary.LittleEndian, &regSpace.Rid)
	if err != nil {
		return regSpace, err
	}
	err = binary.Read(buf, binary.LittleEndian, &regSpace.IDExt)
	if err != nil {
		return regSpace, err
	}

	// TXT.VER.QPIIF
	buf.Seek(int64(txtVerQPIFF), io.SeekStart)
	err = binary.Read(buf, binary.LittleEndian, &regSpace.QpiIf)
	if err != nil {
		return regSpace, err
	}

	// TXT.SINIT.BASE
	buf.Seek(int64(txtsInitBase), io.SeekStart)
	err = binary.Read(buf, binary.LittleEndian, &regSpace.SinitBase)
	if err != nil {
		return regSpace, err
	}

	// TXT.SINIT.SIZE
	buf.Seek(int64(txtsInitSize), io.SeekStart)
	err = binary.Read(buf, binary.LittleEndian, &regSpace.SinitSize)
	if err != nil {
		return regSpace, err
	}

	// TXT.MLE.JOIN
	buf.Seek(int64(txtMLEJoin), io.SeekStart)
	err = binary.Read(buf, binary.LittleEndian, &regSpace.MleJoin)
	if err != nil {
		return regSpace, err
	}

	// TXT.HEAP.BASE
	buf.Seek(int64(txtHeapBase), io.SeekStart)
	err = binary.Read(buf, binary.LittleEndian, &regSpace.HeapBase)
	if err != nil {
		return regSpace, err
	}

	// TXT.HEAP.SIZE
	buf.Seek(int64(txtHeapSize), io.SeekStart)
	err = binary.Read(buf, binary.LittleEndian, &regSpace.HeapSize)
	if err != nil {
		return regSpace, err
	}

	// TXT.PUBLIC.KEY
	for i := 0; i < 4; i++ {
		buf.Seek(int64(txtPublicKey+int64(i)*8), io.SeekStart)
		err = binary.Read(buf, binary.LittleEndian, &regSpace.PublicKey[i])
		if err != nil {
			return regSpace, err
		}
	}

	// TXT.E2STS
	buf.Seek(int64(txtE2STS), io.SeekStart)
	err = binary.Read(buf, binary.LittleEndian, &regSpace.E2Sts)
	if err != nil {
		return regSpace, err
	}
	return regSpace, nil
}

//ParseBIOSDataRegion decodes a raw copy of the BIOSDATA region
func ParseBIOSDataRegion(heap []byte) (TXTBiosData, error) {
	var ret TXTBiosData
	var biosDataSize uint64

	buf := bytes.NewReader(heap)
	// TXT Heap Bios Data size
	err := binary.Read(buf, binary.LittleEndian, &biosDataSize)
	if err != nil {
		return ret, err
	}

	// Bios Data
	err = binary.Read(buf, binary.LittleEndian, &ret.Version)
	if err != nil {
		return ret, err
	}

	err = binary.Read(buf, binary.LittleEndian, &ret.BiosSinitSize)
	if err != nil {
		return ret, err
	}

	err = binary.Read(buf, binary.LittleEndian, &ret.Reserved1)
	if err != nil {
		return ret, err
	}

	err = binary.Read(buf, binary.LittleEndian, &ret.Reserved2)
	if err != nil {
		return ret, err
	}

	err = binary.Read(buf, binary.LittleEndian, &ret.NumLogProcs)
	if err != nil {
		return ret, err
	}

	if ret.Version >= 3 && ret.Version < 5 {
		err = binary.Read(buf, binary.LittleEndian, &ret.SinitFlags)
		if err != nil {
			return ret, err
		}
	}

	if ret.Version >= 5 {
		var mleFlags uint32
		var flags TXTBiosMLEFlags

		err = binary.Read(buf, binary.LittleEndian, &mleFlags)
		if err != nil {
			return ret, err
		}

		flags.SupportsACPIPPI = mleFlags&1 != 0
		flags.IsClientState = mleFlags&6 == 2
		flags.IsServerState = mleFlags&6 == 4
		ret.MleFlags = &flags
	}

	return ret, nil
}

func readTXTStatus(data []byte) (TXTStatus, error) {
	var ret TXTStatus
	var u64 uint64
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &u64)

	if err != nil {
		return ret, err
	}

	ret.SenterDone = u64&(1<<0) != 0
	ret.SexitDone = u64&(1<<1) != 0
	ret.MemConfigLock = u64&(1<<6) != 0
	ret.PrivateOpen = u64&(1<<7) != 0
	ret.Locality1Open = u64&(1<<15) != 0
	ret.Locality2Open = u64&(1<<16) != 0

	return ret, nil
}

func readTXTErrorCode(data []byte) (TXTErrorCode, uint32, error) {
	var ret TXTErrorCode
	var u32 uint32
	buf := bytes.NewReader(data[txtErrorCode:])
	err := binary.Read(buf, binary.LittleEndian, &u32)

	if err != nil {
		return ret, 0, err
	}

	ret.ModuleType = uint8((u32 >> 0) & 0x7)           // 3:0
	ret.ClassCode = uint8((u32 >> 4) & 0x3f)           // 9:4
	ret.MajorErrorCode = uint8((u32 >> 10) & 0x1f)     // 14:10
	ret.SoftwareSource = (u32>>15)&0x1 != 0            // 15
	ret.MinorErrorCode = uint16((u32 >> 16) & 0x3ffff) // 27:16
	ret.Type1Reserved = uint8((u32 >> 28) & 0x3)       // 29:28
	ret.ProcessorSoftware = (u32>>30)&0x1 != 0         // 30
	ret.ValidInvalid = (u32>>31)&0x1 != 0              // 31

	return ret, uint32(u32), nil
}

func readDMAProtectedRange(data []byte) (hwapi.DMAProtectedRange, error) {
	var ret hwapi.DMAProtectedRange
	var u32 uint32
	buf := bytes.NewReader(data[txtDMAProtectedRange:])
	err := binary.Read(buf, binary.LittleEndian, &u32)

	if err != nil {
		return ret, err
	}

	ret.Lock = u32&1 != 0
	ret.Size = uint8((u32 >> 4) & 0xff)   // 11:4
	ret.Top = uint16((u32 >> 20) & 0xfff) // 31:20

	return ret, nil
}

//ReadACMStatus decodes the raw ACM status register bits
func ReadACMStatus(data []byte) (ACMStatus, error) {
	var ret ACMStatus
	var u64 uint64
	buf := bytes.NewReader(data[txtACMStatus:])
	err := binary.Read(buf, binary.LittleEndian, &u64)
	if err != nil {
		return ret, err
	}

	ret.ModuleType = uint8(u64 & 0xF)
	ret.ClassCode = uint8((u64 >> 4) & 0x3f)
	ret.MajorErrorCode = uint8((u64 >> 10) & 0x1f)
	ret.ACMStarted = (u64>>15)&1 == 1
	ret.MinorErrorCode = uint16((u64 >> 16) & 0xfff)
	ret.Valid = (u64>>31)&1 == 1

	return ret, nil
}
