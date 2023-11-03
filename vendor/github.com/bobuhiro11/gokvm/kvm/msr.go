package kvm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"unsafe"
)

type MSR uint32

const (
	MSRIA32TSC            MSR = 0x10
	MSRIA32APICBASE       MSR = 0x1b
	MSRIA32FEATURECONTROL MSR = 0x0000003a
	MSRTSCADJUST          MSR = 0x0000003b
	MSRIA32SPECCTRL       MSR = 0x48
	MSRVIRTSSBD           MSR = 0xc001011f
	MSRIA32PREDCMD        MSR = 0x49
	MSRIA32UCODEREV       MSR = 0x8b
	MSRIA32CORECAPABILITY MSR = 0xcf

	MSRIA32ARCHCAPABILITIES MSR = 0x10a
	MSRIA32PERFCAPABILITIES MSR = 0x345
	MSRIA32TSXCTRL          MSR = 0x122
	MSRIA32TSCDEADLINE      MSR = 0x6e0
	MSRIA32PKRS             MSR = 0x6e1
	MSRARCHLBRCTL           MSR = 0x000014ce
	MSRARCHLBRDEPTH         MSR = 0x000014cf
	MSRARCHLBRFROM0         MSR = 0x00001500
	MSRARCHLBRTO0           MSR = 0x00001600
	MSRARCHLBRINFO0         MSR = 0x00001200
	MSRIA32SGXLEPUBKEYHASH0 MSR = 0x8c
	MSRIA32SGXLEPUBKEYHASH1 MSR = 0x8d
	MSRIA32SGXLEPUBKEYHASH2 MSR = 0x8e
	MSRIA32SGXLEPUBKEYHASH3 MSR = 0x8f

	MSRP6PERFCTR0 MSR = 0xc1

	MSRIA32SMBASE      MSR = 0x9e
	MSRSMICOUNT        MSR = 0x34
	MSRCORETHREADCOUNT MSR = 0x35
	MSRMTRRcap         MSR = 0xfe
	MSRIA32SYSENTERCS  MSR = 0x174
	MSRIA32SYSENTERESP MSR = 0x175
	MSRIA32SYSENTEREIP MSR = 0x176

	MSRMCGCAP    MSR = 0x179
	MSRMCGSTATUS MSR = 0x17a
	MSRMCGCTL    MSR = 0x17b
	MSRMCGEXTCTL MSR = 0x4d0

	MSRP6EVNTSEL0 MSR = 0x186

	MSRIA32PERFSTATUS MSR = 0x198

	MSRIA32MISCENABLE  MSR = 0x1a0
	MSRMTRRfix64K00000 MSR = 0x250
	MSRMTRRfix16K80000 MSR = 0x258
	MSRMTRRfix16KA0000 MSR = 0x259
	MSRMTRRfix4KC0000  MSR = 0x268
	MSRMTRRfix4KC8000  MSR = 0x269
	MSRMTRRfix4KD0000  MSR = 0x26a
	MSRMTRRfix4KD8000  MSR = 0x26b
	MSRMTRRfix4KE0000  MSR = 0x26c
	MSRMTRRfix4KE8000  MSR = 0x26d
	MSRMTRRfix4KF0000  MSR = 0x26e
	MSRMTRRfix4KF8000  MSR = 0x26f

	MSRPAT MSR = 0x277

	MSRMTRRdefType MSR = 0x2ff

	MSRCOREPERFFIXEDCTR0     MSR = 0x309
	MSRCOREPERFFIXEDCTR1     MSR = 0x30a
	MSRCOREPERFFIXEDCTR2     MSR = 0x30b
	MSRCOREPERFFIXEDCTRCTRL  MSR = 0x38d
	MSRCOREPERFGLOBALSTATUS  MSR = 0x38e
	MSRCOREPERFGLOBALCTRL    MSR = 0x38f
	MSRCOREPERFGLOBALOVFCTRL MSR = 0x390

	MSRMC0CTL    MSR = 0x400
	MSRMC0STATUS MSR = 0x401
	MSRMC0ADDR   MSR = 0x402
	MSRMC0MISC   MSR = 0x403

	MSRIA32RTITOUTPUTBASE MSR = 0x560
	MSRIA32RTITOUTPUTMASK MSR = 0x561
	MSRIA32RTITCTL        MSR = 0x570
	MSRIA32RTITSTATUS     MSR = 0x571
	MSRIA32RTITCR3MATCH   MSR = 0x572
	MSRIA32RTITADDR0A     MSR = 0x580
	MSRIA32RTITADDR0B     MSR = 0x581
	MSRIA32RTITADDR1A     MSR = 0x582
	MSRIA32RTITADDR1B     MSR = 0x583
	MSRIA32RTITADDR2A     MSR = 0x584
	MSRIA32RTITADDR2B     MSR = 0x585
	MSRIA32RTITADDR3A     MSR = 0x586
	MSRIA32RTITADDR3B     MSR = 0x587
	MSRSTAR               MSR = 0xc0000081
	MSRLSTAR              MSR = 0xc0000082
	MSRCSTAR              MSR = 0xc0000083
	MSRFMASK              MSR = 0xc0000084
	MSRFSBASE             MSR = 0xc0000100
	MSRGSBASE             MSR = 0xc0000101
	MSRKERNELGSBASE       MSR = 0xc0000102
	MSRTSCAUX             MSR = 0xc0000103
	MSRAMD64TSCRATIO      MSR = 0xc0000104

	MSRVMHSAVEPA MSR = 0xc0010117

	MSRIA32XFD    MSR = 0x000001c4
	MSRIA32XFDERR MSR = 0x000001c5

	MSRIA32BNDCFGS       MSR = 0x00000d90
	MSRIA32XSS           MSR = 0x00000da0
	MSRIA32UMWAITCONTROL MSR = 0xe1

	MSRIA32VMXBASIC             MSR = 0x00000480
	MSRIA32VMXPINBASEDCTLS      MSR = 0x00000481
	MSRIA32VMXPROCBASEDCTLS     MSR = 0x00000482
	MSRIA32VMXEXITCTLS          MSR = 0x00000483
	MSRIA32VMXENTRYCTLS         MSR = 0x00000484
	MSRIA32VMXMISC              MSR = 0x00000485
	MSRIA32VMXCR0FIXED0         MSR = 0x00000486
	MSRIA32VMXCR0FIXED1         MSR = 0x00000487
	MSRIA32VMXCR4FIXED0         MSR = 0x00000488
	MSRIA32VMXCR4FIXED1         MSR = 0x00000489
	MSRIA32VMXVMCSENUM          MSR = 0x0000048a
	MSRIA32VMXPROCBASEDCTLS2    MSR = 0x0000048b
	MSRIA32VMXEPTVPIDCAP        MSR = 0x0000048c
	MSRIA32VMXTRUEPINBASEDCTLS  MSR = 0x0000048d
	MSRIA32VMXTRUEPROCBASEDCTLS MSR = 0x0000048e
	MSRIA32VMXTRUEEXITCTLS      MSR = 0x0000048f
	MSRIA32VMXTRUEENTRYCTLS     MSR = 0x00000490
	MSRIA32VMXVMFUNC            MSR = 0x00000491
)

type MSRList struct {
	NMSRs uint32
	// Perhaps it could be generated dynamically,
	// but if it is large enough, well it would work.
	Indicies [1000]uint32
}

// GetMSRIndexList returns the guest msrs that are supported.
// The list varies by kvm version and host processor, but does not change otherwise.
func GetMSRIndexList(kvmFd uintptr, list *MSRList) error {
	_, err := Ioctl(kvmFd,
		IIOWR(kvmGetMSRIndexList, unsafe.Sizeof(list.NMSRs)),
		uintptr(unsafe.Pointer(list)))

	return err
}

// GetMSRFeatureIndexList returns the list of MSRs that can be passed to the KVMGETMSRS system ioctl.
// This lets userspace probe host capabilities and processor features that are exposed via MSRs
// (e.g., VMX capabilities). This list also varies by kvm version and host processor, but does not change otherwise.
func GetMSRFeatureIndexList(kvmFd uintptr, list *MSRList) error {
	_, err := Ioctl(kvmFd,
		IIOWR(kvmGetMSRFeatureIndexList, unsafe.Sizeof(list.NMSRs)),
		uintptr(unsafe.Pointer(list)))

	return err
}

type MSREntry struct {
	Index   uint32
	Padding uint32
	Data    uint64
}

type MSRS struct {
	NMSRs   uint32
	Padding uint32
	Entries []MSREntry
}

func (m *MSRS) Bytes() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.LittleEndian, m.NMSRs); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, binary.LittleEndian, m.Padding); err != nil {
		return nil, err
	}

	for _, entry := range m.Entries {
		if err := binary.Write(&buf, binary.LittleEndian, entry.Index); err != nil {
			return nil, err
		}

		if err := binary.Write(&buf, binary.LittleEndian, entry.Padding); err != nil {
			return nil, err
		}

		if err := binary.Write(&buf, binary.LittleEndian, entry.Data); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func NewMSRS(data []byte) (*MSRS, error) {
	m := MSRS{}

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, data); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	if err := binary.Read(&buf, binary.LittleEndian, &m.NMSRs); err != nil {
		return nil, err
	}

	if err := binary.Read(&buf, binary.LittleEndian, &m.Padding); err != nil {
		return nil, err
	}

	m.Entries = make([]MSREntry, m.NMSRs)

	if err := binary.Read(&buf, binary.LittleEndian, &m.Entries); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	return &m, nil
}

func SetMSRs(vcpuFd uintptr, msrs *MSRS) error {
	var m *MSRS

	data, err := msrs.Bytes()
	if err != nil {
		return err
	}

	if _, err := Ioctl(vcpuFd,
		IIOW(kvmSetMSRS, 8),
		uintptr(unsafe.Pointer(&data[0]))); err != nil {
		return err
	}

	if m, err = NewMSRS(data); err != nil {
		return err
	}

	*msrs = *m

	return err
}

func GetMSRs(vcpuFd uintptr, msrs *MSRS) error {
	var m *MSRS

	data, err := msrs.Bytes()
	if err != nil {
		return err
	}

	if _, err := Ioctl(vcpuFd,
		IIOWR(kvmGetMSRS, 8),
		uintptr(unsafe.Pointer(&data[0]))); err != nil {
		return err
	}

	if m, err = NewMSRS(data); err != nil {
		return err
	}

	*msrs = *m

	return err
}
