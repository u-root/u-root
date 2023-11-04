//nolint:dupl
package kvm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"unsafe"
)

// CPUID is the set of CPUID entries returned by GetCPUID.
type CPUID struct {
	Nent    uint32
	Padding uint32
	Entries []CPUIDEntry2
}

func (c *CPUID) Bytes() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.LittleEndian, c.Nent); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, binary.LittleEndian, c.Padding); err != nil {
		return nil, err
	}

	for _, entry := range c.Entries {
		if err := binary.Write(&buf, binary.LittleEndian, entry); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func NewCPUID(data []byte) (*CPUID, error) {
	c := CPUID{}

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, data); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	if err := binary.Read(&buf, binary.LittleEndian, &c.Nent); err != nil {
		return nil, err
	}

	if err := binary.Read(&buf, binary.LittleEndian, &c.Padding); err != nil {
		return nil, err
	}

	c.Entries = make([]CPUIDEntry2, c.Nent)

	if err := binary.Read(&buf, binary.LittleEndian, &c.Entries); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	return &c, nil
}

// CPUIDEntry2 is one entry for CPUID. It took 2 tries to get it right :-)
// Thanks x86 :-).
type CPUIDEntry2 struct {
	Function uint32
	Index    uint32
	Flags    uint32
	Eax      uint32
	Ebx      uint32
	Ecx      uint32
	Edx      uint32
	Padding  [3]uint32
}

// GetSupportedCPUID gets all supported CPUID entries for a vm.
func GetSupportedCPUID(kvmFd uintptr, kvmCPUID *CPUID) error {
	var c *CPUID

	data, err := kvmCPUID.Bytes()
	if err != nil {
		return err
	}

	if _, err = Ioctl(kvmFd,
		IIOWR(kvmGetSupportedCPUID, unsafe.Sizeof(kvmCPUID)),
		uintptr(unsafe.Pointer(&data[0]))); err != nil {
		return err
	}

	if c, err = NewCPUID(data); err != nil {
		return err
	}

	*kvmCPUID = *c

	return err
}

// SetCPUID2 sets entries for a vCPU.
// The progression is, hence, get the CPUID entries for a vm, then set them into
// individual vCPUs. This seems odd, but in fact lets code tailor CPUID entries
// as needed.
func SetCPUID2(vcpuFd uintptr, kvmCPUID *CPUID) error {
	data, err := kvmCPUID.Bytes()
	if err != nil {
		return err
	}

	if _, err := Ioctl(vcpuFd,
		IIOW(kvmSetCPUID2, unsafe.Sizeof(kvmCPUID)),
		uintptr(unsafe.Pointer(&data[0]))); err != nil {
		return err
	}

	return err
}

func GetCPUID2(vcpuFd uintptr, kvmCPUID *CPUID) error {
	var c *CPUID

	data, err := kvmCPUID.Bytes()
	if err != nil {
		return err
	}

	if _, err = Ioctl(vcpuFd,
		IIOWR(kvmGetCPUID2, 8),
		uintptr(unsafe.Pointer(&data[0]))); err != nil {
		return err
	}

	if c, err = NewCPUID(data); err != nil {
		return err
	}

	*kvmCPUID = *c

	return err
}

// GetEmulatedCPUID returns x86 cpuid features which are emulated by kvm.
func GetEmulatedCPUID(kvmFd uintptr, kvmCPUID *CPUID) error {
	var c *CPUID

	data, err := kvmCPUID.Bytes()
	if err != nil {
		return err
	}

	if _, err = Ioctl(kvmFd,
		IIOWR(kvmGetEmulatedCPUID, unsafe.Sizeof(kvmCPUID)),
		uintptr(unsafe.Pointer(&data[0]))); err != nil {
		return err
	}

	if c, err = NewCPUID(data); err != nil {
		return err
	}

	*kvmCPUID = *c

	return nil
}
