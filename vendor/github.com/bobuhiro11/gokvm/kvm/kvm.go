package kvm

import (
	"errors"
	"syscall"
	"unsafe"
)

const (
	kvmGetAPIVersion     = 0x00
	kvmCreateVM          = 0x1
	kvmGetMSRIndexList   = 0x02
	kvmCheckExtension    = 0x03
	kvmGetVCPUMMapSize   = 0x04
	kvmGetSupportedCPUID = 0x05

	kvmGetEmulatedCPUID       = 0x09
	kvmGetMSRFeatureIndexList = 0x0A

	kvmCreateVCPU          = 0x41
	kvmGetDirtyLog         = 0x42
	kvmSetNrMMUPages       = 0x44
	kvmGetNrMMUPages       = 0x45
	kvmSetUserMemoryRegion = 0x46
	kvmSetTSSAddr          = 0x47
	kvmSetIdentityMapAddr  = 0x48

	kvmCreateIRQChip = 0x60
	kvmGetIRQChip    = 0x62
	kvmSetIRQChip    = 0x63
	kvmIRQLineStatus = 0x67

	kvmResgisterCoalescedMMIO   = 0x67
	kvmUnResgisterCoalescedMMIO = 0x68

	kvmSetGSIRouting = 0x6A

	kvmReinjectControl = 0x71
	kvmCreatePIT2      = 0x77
	kvmSetClock        = 0x7B
	kvmGetClock        = 0x7C

	kvmRun       = 0x80
	kvmGetRegs   = 0x81
	kvmSetRegs   = 0x82
	kvmGetSregs  = 0x83
	kvmSetSregs  = 0x84
	kvmTranslate = 0x85
	kvmInterrupt = 0x86

	kvmGetMSRS = 0x88
	kvmSetMSRS = 0x89

	kvmGetLAPIC = 0x8e
	kvmSetLAPIC = 0x8f

	kvmSetCPUID2          = 0x90
	kvmGetCPUID2          = 0x91
	kvmTRPAccessReporting = 0x92

	kvmGetMPState = 0x98
	kvmSetMPState = 0x99

	kvmX86SetupMCE           = 0x9C
	kvmX86GetMCECapSupported = 0x9D

	kvmGetPIT2 = 0x9F
	kvmSetPIT2 = 0xA0

	kvmGetVCPUEvents = 0x9F
	kvmSetVCPUEvents = 0xA0

	kvmGetDebugRegs = 0xA1
	kvmSetDebugRegs = 0xA2

	kvmSetTSCKHz = 0xA2
	kvmGetTSCKHz = 0xA3

	kvmGetXCRS = 0xA6
	kvmSetXCRS = 0xA7

	kvmSMI = 0xB7

	kvmGetSRegs2 = 0xCC
	kvmSetSRegs2 = 0xCD

	kvmCreateDev = 0xE0
)

// ExitType is a virtual machine exit type.
//
//go:generate stringer -type=ExitType
type ExitType uint

const (
	EXITUNKNOWN       ExitType = 0
	EXITEXCEPTION     ExitType = 1
	EXITIO            ExitType = 2
	EXITHYPERCALL     ExitType = 3
	EXITDEBUG         ExitType = 4
	EXITHLT           ExitType = 5
	EXITMMIO          ExitType = 6
	EXITIRQWINDOWOPEN ExitType = 7
	EXITSHUTDOWN      ExitType = 8
	EXITFAILENTRY     ExitType = 9
	EXITINTR          ExitType = 10
	EXITSETTPR        ExitType = 11
	EXITTPRACCESS     ExitType = 12
	EXITS390SIEIC     ExitType = 13
	EXITS390RESET     ExitType = 14
	EXITDCR           ExitType = 15
	EXITNMI           ExitType = 16
	EXITINTERNALERROR ExitType = 17

	EXITIOIN  = 0
	EXITIOOUT = 1
)

const (
	numInterrupts   = 0x100
	CPUIDFeatures   = 0x40000001
	CPUIDSignature  = 0x40000000
	CPUIDFuncPerMon = 0x0A
)

var (
	// ErrUnexpectedExitReason is any error that we do not understand.
	ErrUnexpectedExitReason = errors.New("unexpected kvm exit reason")

	// ErrDebug is a debug exit, caused by single step or breakpoint.
	ErrDebug = errors.New("debug exit")
)

// RunData defines the data used to run a VM.
type RunData struct {
	RequestInterruptWindow     uint8
	ImmediateExit              uint8
	_                          [6]uint8
	ExitReason                 uint32
	ReadyForInterruptInjection uint8
	IfFlag                     uint8
	_                          [2]uint8
	CR8                        uint64
	ApicBase                   uint64
	Data                       [32]uint64
}

// IO interprets IO requests from a VM, by unpacking RunData.Data[0:1].
func (r *RunData) IO() (uint64, uint64, uint64, uint64, uint64) {
	direction := r.Data[0] & 0xFF
	size := (r.Data[0] >> 8) & 0xFF
	port := (r.Data[0] >> 16) & 0xFFFF
	count := (r.Data[0] >> 32) & 0xFFFFFFFF
	offset := r.Data[1]

	return direction, size, port, count, offset
}

// GetAPIVersion gets the qemu API version, which changes rarely if at all.
func GetAPIVersion(kvmFd uintptr) (uintptr, error) {
	return Ioctl(kvmFd, IIO(kvmGetAPIVersion), uintptr(0))
}

// CreateVM creates a KVM from the KVM device fd, i.e. /dev/kvm.
func CreateVM(kvmFd uintptr) (uintptr, error) {
	return Ioctl(kvmFd, IIO(kvmCreateVM), uintptr(0))
}

// CreateVCPU creates a single virtual CPU from the virtual machine FD.
// Thus, the progression:
// fd from opening /dev/kvm
// vmfd from creating a vm from the fd
// vcpu fd from the vmfd.
func CreateVCPU(vmFd uintptr, vcpuID int) (uintptr, error) {
	return Ioctl(vmFd, IIO(kvmCreateVCPU), uintptr(vcpuID))
}

// Run runs a single vcpu from the vcpufd from createvcpu.
func Run(vcpuFd uintptr) error {
	_, err := Ioctl(vcpuFd, IIO(kvmRun), uintptr(0))
	if err != nil {
		// refs: https://github.com/kvmtool/kvmtool/blob/415f92c33a227c02f6719d4594af6fad10f07abf/kvm-cpu.c#L44
		if errors.Is(err, syscall.EAGAIN) || errors.Is(err, syscall.EINTR) {
			return nil
		}
	}

	return err
}

// GetVCPUMmapSize returns the size of the VCPU region. This size is
// required for interacting with the vcpu, as the struct size can change
// over time.
func GetVCPUMMmapSize(kvmFd uintptr) (uintptr, error) {
	return Ioctl(kvmFd, IIO(kvmGetVCPUMMapSize), uintptr(0))
}

func SetTSCKHz(vcpuFd uintptr, freq uint64) error {
	_, err := Ioctl(vcpuFd,
		IIO(kvmSetTSCKHz), uintptr(freq))

	return err
}

func GetTSCKHz(vcpuFd uintptr) (uint64, error) {
	ret, err := Ioctl(vcpuFd,
		IIO(kvmGetTSCKHz), 0)
	if err != nil {
		return 0, err
	}

	return uint64(ret), nil
}

type ClockFlag uint32

const (
	TSCStable ClockFlag = 2
	Realtime  ClockFlag = (1 << 2)
	HostTSC   ClockFlag = (1 << 3)
)

type ClockData struct {
	Clock    uint64
	Flags    uint32
	_        uint32
	Realtime uint64
	HostTSC  uint64
	_        [4]uint32
}

// SetClock sets the current timestamp of kvmclock to the value specified in its parameter.
// In conjunction with GET_CLOCK, it is used to ensure monotonicity on scenarios such as migration.
func SetClock(vmFd uintptr, cd *ClockData) error {
	_, err := Ioctl(vmFd,
		IIOW(kvmSetClock, unsafe.Sizeof(ClockData{})),
		uintptr(unsafe.Pointer(cd)))

	return err
}

// GetClock gets the current timestamp of kvmclock as seen by the current guest.
// In conjunction with SET_CLOCK, it is used to ensure monotonicity on scenarios such as migration.
func GetClock(vmFd uintptr, cd *ClockData) error {
	_, err := Ioctl(vmFd,
		IIOR(kvmGetClock, unsafe.Sizeof(ClockData{})),
		uintptr(unsafe.Pointer(cd)))

	return err
}

type DevType uint32

const (
	DevFSLMPIC20 DevType = 1 + iota
	DevFSLMPIC42
	DevXICS
	DevVFIO
	_
	DevFLIC
	_
	_
	DevXIVE
	_
	DevMAX
)

type Device struct {
	Type  uint32
	Fd    uint32
	Flags uint32
}

// CreateDev creates an emulated device in the kernel.
// The file descriptor returned in fd can be used with SET/GET/HAS_DEVICE_ATTR.
func CreateDev(vmFd uintptr, dev *Device) error {
	_, err := Ioctl(vmFd,
		IIOWR(kvmCreateDev, unsafe.Sizeof(Device{})),
		uintptr(unsafe.Pointer(dev)))

	return err
}

// Translation is a struct for TRANSLATE queries.
type Translation struct {
	// LinearAddress is input.
	// Most people call this a "virtual address"
	// Intel has their own name.
	LinearAddress uint64

	// This is output
	PhysicalAddress uint64
	Valid           uint8
	Writeable       uint8
	Usermode        uint8
	_               [5]uint8
}

// Translate translates a virtual address according to the vcpu’s current address translation mode.
func Translate(vcpuFd uintptr, t *Translation) error {
	_, err := Ioctl(vcpuFd,
		IIOWR(kvmTranslate, unsafe.Sizeof(Translation{})),
		uintptr(unsafe.Pointer(t)))

	return err
}

type MPState struct {
	State uint32
}

const (
	MPStateRunnable      uint32 = 0 + iota // x86, arm64, riscv
	MPStateUninitialized                   // x86
	MPStateInitReceived                    // x86
	MPStateHalted                          // x86
	MPStateSipiReceived                    // x86
	MPStateStopped                         // x86
	MPStateCheckStop                       // s390, arm64, riscv
	MPStateOperating                       // s390
	MPStateLoad                            // s390
	MPStateApResetHold                     // s390
	MPStateSuspended                       // arm64
)

// GetMPState returns the vcpu’s current multiprocessing state.
func GetMPState(vcpuFd uintptr, mps *MPState) error {
	_, err := Ioctl(vcpuFd,
		IIOR(kvmGetMPState, unsafe.Sizeof(MPState{})),
		uintptr(unsafe.Pointer(mps)))

	return err
}

// SetMPState sets the vcpu’s current multiprocessing state.
func SetMPState(vcpuFd uintptr, mps *MPState) error {
	_, err := Ioctl(vcpuFd,
		IIOW(kvmSetMPState, unsafe.Sizeof(MPState{})),
		uintptr(unsafe.Pointer(mps)))

	return err
}

type Exception struct {
	Inject       uint8
	Nr           uint8
	HadErrorCode uint8
	Pending      uint8
	ErrorCode    uint32
}

type Interrupt struct {
	Inject uint8
	Nr     uint8
	Soft   uint8
	Shadow uint8
}

type NMI struct {
	Inject  uint8
	Pending uint8
	Masked  uint8
	_       uint8
}

type SMI struct {
	SMM          uint8
	Pening       uint8
	SMMInsideNMI uint8
	LatchedInit  uint8
}

type VCPUEvents struct {
	E                   Exception
	I                   Interrupt
	N                   NMI
	SipiVector          uint32
	Flags               uint32
	S                   SMI
	TripleFault         uint8
	_                   [26]uint8
	ExceptionHasPayload uint8
	ExceptionPayload    uint64
}

// GetVCPUEvents gets currently pending exceptions, interrupts, and NMIs as well as related states of the vcpu.
func GetVCPUEvents(vcpuFd uintptr, event *VCPUEvents) error {
	_, err := Ioctl(vcpuFd,
		IIOR(kvmGetVCPUEvents, unsafe.Sizeof(VCPUEvents{})),
		uintptr(unsafe.Pointer(event)))

	return err
}

// SetVCPUEvents sets spending exceptions, interrupts, and NMIs as well as related states of the vcpu.
func SetVCPUEvents(vcpuFd uintptr, event *VCPUEvents) error {
	_, err := Ioctl(vcpuFd,
		IIOW(kvmSetVCPUEvents, unsafe.Sizeof(VCPUEvents{})),
		uintptr(unsafe.Pointer(event)))

	return err
}

// SMI queues an SMI on the thread’s vcpu.
func PutSMI(vcpuFd uintptr) error {
	_, err := Ioctl(vcpuFd, IIO(kvmSMI), 0)

	return err
}
