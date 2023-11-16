//nolint:dupl
package kvm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"unsafe"
)

// irqLevel defines an IRQ as Level? Not sure.
type irqLevel struct {
	IRQ   uint32
	Level uint32
}

// IRQLines sets the interrupt line for an IRQ.
func IRQLineStatus(vmFd uintptr, irq, level uint32) error {
	irqLev := irqLevel{
		IRQ:   irq,
		Level: level,
	}
	_, err := Ioctl(vmFd,
		IIOWR(kvmIRQLineStatus, unsafe.Sizeof(irqLevel{})),
		uintptr(unsafe.Pointer(&irqLev)))

	return err
}

// CreateIRQChip creates an IRQ device (chip) to which to attach interrupts?
func CreateIRQChip(vmFd uintptr) error {
	_, err := Ioctl(vmFd, IIO(kvmCreateIRQChip), 0)

	return err
}

// pitConfig defines properties of a programmable interrupt timer.
type pitConfig struct {
	Flags uint32
	_     [15]uint32
}

// CreatePIT2 creates a PIT type 2. Just having one was not enough.
func CreatePIT2(vmFd uintptr) error {
	pit := pitConfig{
		Flags: 0,
	}
	_, err := Ioctl(vmFd,
		IIOW(kvmCreatePIT2, unsafe.Sizeof(pitConfig{})),
		uintptr(unsafe.Pointer(&pit)))

	return err
}

type PITChannelState struct {
	Count         uint32
	LatchedCount  uint16
	CountLatched  uint8
	StatusLatched uint8
	Status        uint8
	ReadState     uint8
	WriteState    uint8
	WriteLatch    uint8
	RWMode        uint8
	Mode          uint8
	BCD           uint8
	Gate          uint8
	CountLoadTime int64
}

type PITState2 struct {
	Channels [3]PITChannelState
	Flags    uint32
	_        [9]uint32
}

// GetPIT2 retrieves the state of the in-kernel PIT model. Only valid after KVM_CREATE_PIT2.
func GetPIT2(vmFd uintptr, pstate *PITState2) error {
	_, err := Ioctl(vmFd,
		IIOR(kvmGetPIT2, unsafe.Sizeof(PITState2{})),
		uintptr(unsafe.Pointer(pstate)))

	return err
}

// SetPIT2 sets the state of the in-kernel PIT model. Only valid after KVM_CREATE_PIT2.
func SetPIT2(vmFd uintptr, pstate *PITState2) error {
	_, err := Ioctl(vmFd,
		IIOW(kvmSetPIT2, unsafe.Sizeof(PITState2{})),
		uintptr(unsafe.Pointer(pstate)))

	return err
}

type PICState struct {
	LastIRR                uint8 /* edge detection */
	IRR                    uint8 /* interrupt request register */
	IMR                    uint8 /* interrupt mask register */
	ISR                    uint8 /* interrupt service register */
	PriorityAdd            uint8 /* highest irq priority */
	IRQBase                uint8
	ReadRegSelect          uint8
	Poll                   uint8
	SpecialMask            uint8
	InitState              uint8
	AutoEOI                uint8
	RotateOnAutoEOI        uint8
	SpecialFullyNestedMode uint8
	Init4                  uint8 /* true if 4 byte init */
	ELCR                   uint8 /* PIIX edge/trigger selection */
	ELCRMask               uint8
}

type IRQChip struct {
	ChipID uint32
	_      uint32
	Chip   [512]byte
}

// GetIRQChip reads the state of a kernel interrupt controller created with
// KVM_CREATE_IRQCHIP into a buffer provided by the caller.
func GetIRQChip(vmFd uintptr, irqc *IRQChip) error {
	_, err := Ioctl(vmFd,
		IIOWR(kvmGetIRQChip, unsafe.Sizeof(IRQChip{})),
		uintptr(unsafe.Pointer(irqc)))

	return err
}

// SetIRQChip sets the state of a kernel interrupt controller created with
// KVM_CREATE_IRQCHIP from a buffer provided by the caller.
func SetIRQChip(vmFd uintptr, irqc *IRQChip) error {
	_, err := Ioctl(vmFd,
		IIOR(kvmSetIRQChip, unsafe.Sizeof(IRQChip{})), uintptr(unsafe.Pointer(irqc)))

	return err
}

type IRQRoutingIRQChip struct {
	IRQChip uint32
	Pin     uint32
}

type IRQRoutingEntry struct {
	GSI   uint32
	Type  uint32
	Flags uint32
	_     uint32
	IRQRoutingIRQChip
}

type IRQRouting struct {
	Nr      uint32
	Flags   uint32
	Entries []IRQRoutingEntry
}

func (r *IRQRouting) Bytes() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.LittleEndian, r.Nr); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, binary.LittleEndian, r.Flags); err != nil {
		return nil, err
	}

	for _, entry := range r.Entries {
		if err := binary.Write(&buf, binary.LittleEndian, entry); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func NewIRQRouting(data []byte) (*IRQRouting, error) {
	r := IRQRouting{}

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, data); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	if err := binary.Read(&buf, binary.LittleEndian, &r.Nr); err != nil {
		return nil, err
	}

	if err := binary.Read(&buf, binary.LittleEndian, &r.Flags); err != nil {
		return nil, err
	}

	r.Entries = make([]IRQRoutingEntry, r.Nr)

	if err := binary.Read(&buf, binary.LittleEndian, &r.Entries); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	return &r, nil
}

// SetGSIRouting sets the GSI routing table entries, overwriting any previously set entries.
func SetGSIRouting(vmFd uintptr, irqR *IRQRouting) error {
	data, err := irqR.Bytes()
	if err != nil {
		return err
	}

	_, err = Ioctl(vmFd,
		IIOW(kvmSetGSIRouting, unsafe.Sizeof(irqR)),
		uintptr(unsafe.Pointer(&data[0])))

	return err
}

// InjectInterrupt queues a hardware interrupt vector to be injected.
func InjectInterrupt(vcpuFd uintptr, intr uint32) error {
	_, err := Ioctl(vcpuFd,
		IIOW(kvmInterrupt, 4),
		uintptr(intr))

	return err
}

const LAPICRegSize = 0x400

type LAPICState struct {
	Regs [LAPICRegSize]byte
}

// GetLocalAPIC reads the Local APIC registers and copies them into the input argument.
func GetLocalAPIC(vcpuFd uintptr, lapic *LAPICState) error {
	_, err := Ioctl(vcpuFd,
		IIOR(kvmGetLAPIC, unsafe.Sizeof(LAPICState{})),
		uintptr(unsafe.Pointer(lapic)))

	return err
}

// SetLocalAPIC copies the input argument into the Local APIC registers.
func SetLocalAPIC(vcpuFd uintptr, lapic *LAPICState) error {
	_, err := Ioctl(vcpuFd,
		IIOW(kvmSetLAPIC, unsafe.Sizeof(LAPICState{})),
		uintptr(unsafe.Pointer(lapic)))

	return err
}

// ReinjectControl sets i8254 Inject mode.
func ReinjectControl(vmFd uintptr, mode uint8) error {
	tmp := struct {
		pitReinject uint8
		_           [31]byte
	}{
		pitReinject: mode,
	}
	_, err := Ioctl(vmFd,
		IIO(kvmReinjectControl), uintptr(unsafe.Pointer(&tmp)))

	return err
}
