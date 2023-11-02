package machine

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"github.com/bobuhiro11/gokvm/bootparam"
	"github.com/bobuhiro11/gokvm/ebda"
	"github.com/bobuhiro11/gokvm/iodev"
	"github.com/bobuhiro11/gokvm/kvm"
	"github.com/bobuhiro11/gokvm/pci"
	"github.com/bobuhiro11/gokvm/pvh"
	"github.com/bobuhiro11/gokvm/serial"
	"github.com/bobuhiro11/gokvm/tap"
	"github.com/bobuhiro11/gokvm/virtio"
	"golang.org/x/arch/x86/x86asm"
)

const (
	bootParamAddr = 0x10000
	cmdlineAddr   = 0x20000

	initrdAddr  = 0xf000000
	highMemBase = 0x100000

	serialIRQ    = 4
	virtioNetIRQ = 9
	virtioBlkIRQ = 10

	pageTableBase = 0x30_000

	MinMemSize = 1 << 25
)

const (
	// These *could* be in kvm, but we'll see.

	// golangci-lint is completely wrong about these names.
	// Control Register Paging Enable for example:
	// golang style requires all letters in an acronym to be caps.
	// CR0 bits.
	CR0xPE = 1
	CR0xMP = (1 << 1)
	CR0xEM = (1 << 2)
	CR0xTS = (1 << 3)
	CR0xET = (1 << 4)
	CR0xNE = (1 << 5)
	CR0xWP = (1 << 16)
	CR0xAM = (1 << 18)
	CR0xNW = (1 << 29)
	CR0xCD = (1 << 30)
	CR0xPG = (1 << 31)

	// CR4 bits.
	CR4xVME        = 1
	CR4xPVI        = (1 << 1)
	CR4xTSD        = (1 << 2)
	CR4xDE         = (1 << 3)
	CR4xPSE        = (1 << 4)
	CR4xPAE        = (1 << 5)
	CR4xMCE        = (1 << 6)
	CR4xPGE        = (1 << 7)
	CR4xPCE        = (1 << 8)
	CR4xOSFXSR     = (1 << 8)
	CR4xOSXMMEXCPT = (1 << 10)
	CR4xUMIP       = (1 << 11)
	CR4xVMXE       = (1 << 13)
	CR4xSMXE       = (1 << 14)
	CR4xFSGSBASE   = (1 << 16)
	CR4xPCIDE      = (1 << 17)
	CR4xOSXSAVE    = (1 << 18)
	CR4xSMEP       = (1 << 20)
	CR4xSMAP       = (1 << 21)

	EFERxSCE = 1
	EFERxLME = (1 << 8)
	EFERxLMA = (1 << 10)
	EFERxNXE = (1 << 11)

	// 64-bit page * entry bits.
	PDE64xPRESENT  = 1
	PDE64xRW       = (1 << 1)
	PDE64xUSER     = (1 << 2)
	PDE64xACCESSED = (1 << 5)
	PDE64xDIRTY    = (1 << 6)
	PDE64xPS       = (1 << 7)
	PDE64xG        = (1 << 8)
)

const (
	// Poison is an instruction that should force a vmexit.
	// it fills memory to make catching guest errors easier.
	// vmcall, nop is this pattern
	// Poison = []byte{0x0f, 0x0b, } //0x01, 0xC1, 0x90}
	// Disassembly:
	// 0:  b8 be ba fe ca          mov    eax,0xcafebabe
	// 5:  90                      nop
	// 6:  0f 0b                   ud2
	Poison = "\xB8\xBE\xBA\xFE\xCA\x90\x0F\x0B"
)

var ErrZeroSizeKernel = errors.New("kernel is 0 bytes")

// ErrWriteToCF9 indicates a write to cf9, the standard x86 reset port.
var ErrWriteToCF9 = fmt.Errorf("power cycle via 0xcf9")

// ErrBadVA indicates a bad virtual address was used.
var ErrBadVA = fmt.Errorf("bad virtual address")

// ErrBadCPU indicates a cpu number is invalid.
var ErrBadCPU = fmt.Errorf("bad cpu number")

// ErrUnsupported indicates something we do not yet do.
var ErrUnsupported = fmt.Errorf("unsupported")

// ErrMemTooSmall indicates the requested memory size is too small.
var ErrMemTooSmall = fmt.Errorf("mem request must be at least 1<<20")

var ErrNotELF64File = fmt.Errorf("file is not ELF64")

var errPTNoteHasNoFSize = fmt.Errorf("elf programm PT_NOTE has file size equel zero")

type Machine struct {
	kvmFd, vmFd    uintptr
	vcpuFds        []uintptr
	mem            []byte
	runs           []*kvm.RunData
	pci            *pci.PCI
	serial         *serial.Serial
	devices        []iodev.Device
	ioportHandlers [0x10000][2]func(port uint64, bytes []byte) error
}

// New creates a new KVM. This includes opening the kvm device, creating VM, creating
// vCPUs, and attaching memory, disk (if needed), and tap (if needed).
func New(kvmPath string, nCpus int, memSize int) (*Machine, error) {
	if memSize < MinMemSize {
		return nil, fmt.Errorf("memory size %d:%w", memSize, ErrMemTooSmall)
	}

	m := &Machine{}

	m.pci = pci.New(pci.NewBridge())

	var err error

	m.kvmFd, m.vmFd, m.vcpuFds, m.runs, err = initVMandVCPU(kvmPath, nCpus)
	if err != nil {
		return nil, err
	}

	// initCPUIDs here manually
	for cpuNr := range m.runs {
		if err := m.initCPUID(cpuNr); err != nil {
			return nil, err
		}
	}

	// Another coding anti-pattern reguired by golangci-lint.
	// Would not pass review in Google.
	if m.mem, err = syscall.Mmap(-1, 0, memSize,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED|syscall.MAP_ANONYMOUS); err != nil {
		return m, err
	}

	err = kvm.SetUserMemoryRegion(m.vmFd, &kvm.UserspaceMemoryRegion{
		Slot: 0, Flags: 0, GuestPhysAddr: 0, MemorySize: uint64(memSize),
		UserspaceAddr: uint64(uintptr(unsafe.Pointer(&m.mem[0]))),
	})
	if err != nil {
		return m, err
	}

	// Poison memory.
	// 0 is valid instruction and if you start running in the middle of all those
	// 0's it is impossible to diagnore.
	for i := highMemBase; i < len(m.mem); i += len(Poison) {
		copy(m.mem[i:], Poison)
	}

	return m, nil
}

func (m *Machine) AddTapIf(tapIfName string) error {
	t, err := tap.New(tapIfName)
	if err != nil {
		return err
	}

	v := virtio.NewNet(virtioNetIRQ, m, t, m.mem)
	go v.TxThreadEntry()
	go v.RxThreadEntry()
	// 00:01.0 for Virtio net
	m.pci.Devices = append(m.pci.Devices, v)

	return nil
}

func (m *Machine) AddDisk(diskPath string) error {
	v, err := virtio.NewBlk(diskPath, virtioBlkIRQ, m, m.mem)
	if err != nil {
		return err
	}

	go v.IOThreadEntry()
	// 00:02.0 for Virtio blk
	m.pci.Devices = append(m.pci.Devices, v)

	return nil
}

// Translate translates a virtual address for all active CPUs
// and returns a []*Translate or error.
func (m *Machine) Translate(vaddr uint64) ([]*kvm.Translation, error) {
	t := make([]*kvm.Translation, 0, len(m.vcpuFds))

	for cpu := range m.vcpuFds {
		tr := &kvm.Translation{
			LinearAddress: vaddr,
		}
		if err := kvm.Translate(m.vcpuFds[cpu], tr); err != nil {
			return t, err
		}

		t = append(t, tr)
	}

	return t, nil
}

// SetupRegs sets up the general purpose registers,
// including a RIP and BP.
func (m *Machine) SetupRegs(rip, bp uint64, amd64 bool) error {
	for _, cpu := range m.vcpuFds {
		if err := m.initRegs(cpu, rip, bp); err != nil {
			return err
		}

		if err := m.initSregs(cpu, amd64); err != nil {
			return err
		}
	}

	return nil
}

// RunData returns the kvm.RunData for the VM.
func (m *Machine) RunData() []*kvm.RunData {
	return m.runs
}

func (m *Machine) LoadPVH(kern, initrd *os.File, cmdline string) error {
	// Set EDBA-Pointer
	edbaval := uint32(bootparam.EBDAStart >> 4)
	edbabytes := make([]byte, 4)

	// Convert EBDA-Address to bytes
	binary.LittleEndian.PutUint32(edbabytes, edbaval)

	// Copy EBDA-Address to memory at EBDAPointer-Address (0x40e)
	copy(m.mem[pvh.EBDAPointer:], edbabytes)

	// Create EBDA/mptables - Required for booting into Linux with PVH.
	e, err := ebda.New(len(m.vcpuFds))
	if err != nil {
		return err
	}

	// Convert EBDA/mptables to bytes
	eb, err := e.Bytes()
	if err != nil {
		return err
	}

	// Write EBDA/mptables to memory at EBDAStart (0x0009_FC00)
	copy(m.mem[bootparam.EBDAStart:], eb)

	// Create Global Descriptor Table
	gdt := pvh.CreateGDT()

	// Write GDT to memory
	copy(m.mem[pvh.BootGDTStart:], gdt.Bytes())

	// Write IDT to memory
	copy(m.mem[pvh.BootIDTStart:], []byte{0x0})

	// Load firmware as ELF
	fwElf, err := elf.NewFile(kern)
	if err != nil {
		// Abort if no ELF-File
		return err
	}

	ripAddr := fwElf.Entry

	for _, entry := range fwElf.Progs {
		if entry.Type == elf.PT_LOAD {
			_, err := entry.ReadAt(m.mem[entry.Paddr:], 0)
			if err != nil && !errors.Is(err, io.EOF) {
				return err
			}
		} else if entry.Type == elf.PT_NOTE {
			if entry.Filesz == 0 {
				return errPTNoteHasNoFSize
			}

			addr, _ := pvh.ParsePVHEntry(kern, entry)

			if fwElf.Entry != uint64(addr) {
				ripAddr = uint64(addr)
			}
		}

		continue
	}

	for _, cpu := range m.vcpuFds {
		if err := pvh.InitRegs(cpu, ripAddr); err != nil {
			return err
		}

		if err := pvh.InitSRegs(cpu, gdt); err != nil {
			return err
		}
	}

	pvhstartinfo := pvh.NewStartInfo(bootparam.EBDAStart, cmdlineAddr)

	if initrd != nil {
		initrdSize, err := initrd.ReadAt(m.mem[initrdAddr:], 0)
		if err != nil && initrdSize == 0 && !errors.Is(err, io.EOF) {
			return fmt.Errorf("initrd: (%v, %w)", initrdSize, err)
		}

		// Load kernel command-line parameters
		copy(m.mem[cmdlineAddr:], cmdline)
		m.mem[cmdlineAddr+len(cmdline)] = 0 // for null terminated string

		ramdiskmod := pvh.NewModListEntry(initrdAddr, uint64(initrdSize), 0)

		pvhstartinfo.NrModules += 1
		pvhstartinfo.ModlistPAddr = pvh.PVHModlistStart

		ramdiskmodbytes, err := ramdiskmod.Bytes()
		if err != nil {
			return err
		}

		copy(m.mem[pvh.PVHModlistStart:], ramdiskmodbytes)

		m.AddDevice(&iodev.Noop{Port: 0x80, Psize: 0x30}) // DMA Page Registers (Commonly 74L612 Chip)
	} else {
		m.AddDevice(&iodev.PostCode{}) // Port 0x80
	}

	memmapentries := make([]*pvh.HVMMemMapTableEntry, 0)

	entry0 := pvh.NewMemMapTableEntry(0,
		bootparam.EBDAStart,
		bootparam.E820Ram)

	memmapentries = append(memmapentries, entry0)

	entry := pvh.NewMemMapTableEntry(
		pvh.HighRAMStart,
		uint64(len(m.mem)-pvh.HighRAMStart),
		bootparam.E820Ram)

	memmapentries = append(memmapentries, entry)

	pvhstartinfo.MemMapEntries = uint32(len(memmapentries))

	memOffset := pvh.PVHMemMapStart

	// Copy the MEMMapEntries to memory one at a time.
	for _, entry := range memmapentries {
		b, err := entry.Bytes()
		if err != nil {
			return err
		}

		copy(m.mem[memOffset:], b)

		memOffset += len(b)
	}

	// Copy the PVHInfoStart struct to memory.
	pvhstartinfob, err := pvhstartinfo.Bytes()
	if err != nil {
		return err
	}

	copy(m.mem[pvh.PVHInfoStart:], pvhstartinfob)

	if m.serial, err = serial.New(m); err != nil {
		return err
	}

	m.AddDevice(&iodev.FWDebug{}) // Port 0x402
	m.AddDevice(iodev.NewCMOS(0xC000000, 0x0))
	m.AddDevice(iodev.NewACPIPMTimer())
	m.initIOPortHandlers()

	return nil
}

// LoadLinux loads a bzImage or ELF file, an optional initrd, and
// optional params.
func (m *Machine) LoadLinux(kernel, initrd io.ReaderAt, params string) error {
	var (
		DefaultKernelAddr = uint64(highMemBase)
		err               error
	)

	e, err := ebda.New(len(m.vcpuFds))
	if err != nil {
		return err
	}

	bytes, err := e.Bytes()
	if err != nil {
		return err
	}

	copy(m.mem[bootparam.EBDAStart:], bytes)

	// Load initrd
	var initrdSize int
	if initrd != nil {
		initrdSize, err = initrd.ReadAt(m.mem[initrdAddr:], 0)
		if err != nil && initrdSize == 0 && !errors.Is(err, io.EOF) {
			return fmt.Errorf("initrd: (%v, %w)", initrdSize, err)
		}
	}

	// Load kernel command-line parameters
	copy(m.mem[cmdlineAddr:], params)
	m.mem[cmdlineAddr+len(params)] = 0 // for null terminated string

	// try to read as ELF. If it fails, no problem,
	// next effort is to read as a bzimage.
	var isElfFile bool

	k, err := elf.NewFile(kernel)
	if err == nil {
		isElfFile = true
	}

	bootParam := &bootparam.BootParam{}

	// might be a bzimage
	if !isElfFile {
		// Load Boot Param
		bootParam, err = bootparam.New(kernel)
		if err != nil {
			return err
		}
	}

	// refs https://github.com/kvmtool/kvmtool/blob/0e1882a49f81cb15d328ef83a78849c0ea26eecc/x86/bios.c#L66-L86
	bootParam.AddE820Entry(
		bootparam.RealModeIvtBegin,
		bootparam.EBDAStart-bootparam.RealModeIvtBegin,
		bootparam.E820Ram,
	)
	bootParam.AddE820Entry(
		bootparam.EBDAStart,
		bootparam.VGARAMBegin-bootparam.EBDAStart,
		bootparam.E820Reserved,
	)
	bootParam.AddE820Entry(
		bootparam.MBBIOSBegin,
		bootparam.MBBIOSEnd-bootparam.MBBIOSBegin,
		bootparam.E820Reserved,
	)
	bootParam.AddE820Entry(
		highMemBase,
		uint64(len(m.mem)-highMemBase),
		bootparam.E820Ram,
	)

	bootParam.Hdr.VidMode = 0xFFFF                                                                  // Proto ALL
	bootParam.Hdr.TypeOfLoader = 0xFF                                                               // Proto 2.00+
	bootParam.Hdr.RamdiskImage = initrdAddr                                                         // Proto 2.00+
	bootParam.Hdr.RamdiskSize = uint32(initrdSize)                                                  // Proto 2.00+
	bootParam.Hdr.LoadFlags |= bootparam.CanUseHeap | bootparam.LoadedHigh | bootparam.KeepSegments // Proto 2.00+
	bootParam.Hdr.HeapEndPtr = 0xFE00                                                               // Proto 2.01+
	bootParam.Hdr.ExtLoaderVer = 0                                                                  // Proto 2.02+
	bootParam.Hdr.CmdlinePtr = cmdlineAddr                                                          // Proto 2.06+
	bootParam.Hdr.CmdlineSize = uint32(len(params) + 1)                                             // Proto 2.06+

	bytes, err = bootParam.Bytes()
	if err != nil {
		return err
	}

	copy(m.mem[bootParamAddr:], bytes)

	var (
		amd64    bool
		kernSize int
	)

	switch isElfFile {
	case false:
		// Load kernel
		// copy to g.mem with offset setupsz
		//
		// The 32-bit (non-real-mode) kernel starts at offset (setup_sects+1)*512 in
		// the kernel file (again, if setup_sects == 0 the real value is 4.) It should
		// be loaded at address 0x10000 for Image/zImage kernels and highMemBase for bzImage kernels.
		//
		// refs: https://www.kernel.org/doc/html/latest/x86/boot.html#loading-the-rest-of-the-kernel
		setupsz := int(bootParam.Hdr.SetupSects+1) * 512

		kernSize, err = kernel.ReadAt(m.mem[DefaultKernelAddr:], int64(setupsz))

		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("kernel: (%v, %w)", kernSize, err)
		}
	case true:
		if k.Class == elf.ELFCLASS64 {
			amd64 = true
		}

		DefaultKernelAddr = k.Entry

		for i, p := range k.Progs {
			if p.Type != elf.PT_LOAD {
				continue
			}

			log.Printf("Load elf segment @%#x from file %#x %#x bytes", p.Paddr, p.Off, p.Filesz)

			n, err := p.ReadAt(m.mem[p.Paddr:], 0)
			if !errors.Is(err, io.EOF) || uint64(n) != p.Filesz {
				return fmt.Errorf("reading ELF prog %d@%#x: %d/%d bytes, err %w", i, p.Paddr, n, p.Filesz, err)
			}

			kernSize += n
		}
	}

	if kernSize == 0 {
		return ErrZeroSizeKernel
	}

	if err := m.SetupRegs(DefaultKernelAddr, bootParamAddr, amd64); err != nil {
		return err
	}

	if m.serial, err = serial.New(m); err != nil {
		return err
	}

	m.AddDevice(iodev.NewCMOS(0xC000_0000, 0x0))
	m.AddDevice(&iodev.Noop{Port: 0x80, Psize: 0xA0})
	m.initIOPortHandlers()

	return nil
}

// GetInputChan returns a chan <- byte for serial.
func (m *Machine) GetInputChan() chan<- byte {
	return m.serial.GetInputChan()
}

// GetRegs gets regs for vCPU.
func (m *Machine) GetRegs(cpu int) (*kvm.Regs, error) {
	fd, err := m.CPUToFD(cpu)
	if err != nil {
		return nil, err
	}

	return kvm.GetRegs(fd)
}

// GetSRegs gets sregs for vCPU.
func (m *Machine) GetSRegs(cpu int) (*kvm.Sregs, error) {
	fd, err := m.CPUToFD(cpu)
	if err != nil {
		return nil, err
	}

	return kvm.GetSregs(fd)
}

// SetRegs sets regs for vCPU.
func (m *Machine) SetRegs(cpu int, r *kvm.Regs) error {
	fd, err := m.CPUToFD(cpu)
	if err != nil {
		return err
	}

	return kvm.SetRegs(fd, r)
}

// SetSRegs sets sregs for vCPU.
func (m *Machine) SetSRegs(cpu int, s *kvm.Sregs) error {
	fd, err := m.CPUToFD(cpu)
	if err != nil {
		return err
	}

	return kvm.SetSregs(fd, s)
}

func (m *Machine) initRegs(vcpufd uintptr, rip, bp uint64) error {
	regs, err := kvm.GetRegs(vcpufd)
	if err != nil {
		return err
	}

	// Clear all FLAGS bits, except bit 1 which is always set.
	regs.RFLAGS = 2
	regs.RIP = rip
	// Create stack which will grow down.
	regs.RSI = bp

	if err := kvm.SetRegs(vcpufd, regs); err != nil {
		return err
	}

	return nil
}

func (m *Machine) initSregs(vcpufd uintptr, amd64 bool) error {
	sregs, err := kvm.GetSregs(vcpufd)
	if err != nil {
		return err
	}

	if !amd64 {
		// set all segment flat
		sregs.CS.Base, sregs.CS.Limit, sregs.CS.G = 0, 0xFFFFFFFF, 1
		sregs.DS.Base, sregs.DS.Limit, sregs.DS.G = 0, 0xFFFFFFFF, 1
		sregs.FS.Base, sregs.FS.Limit, sregs.FS.G = 0, 0xFFFFFFFF, 1
		sregs.GS.Base, sregs.GS.Limit, sregs.GS.G = 0, 0xFFFFFFFF, 1
		sregs.ES.Base, sregs.ES.Limit, sregs.ES.G = 0, 0xFFFFFFFF, 1
		sregs.SS.Base, sregs.SS.Limit, sregs.SS.G = 0, 0xFFFFFFFF, 1

		sregs.CS.DB, sregs.SS.DB = 1, 1
		sregs.CR0 |= 1 // protected mode

		if err := kvm.SetSregs(vcpufd, sregs); err != nil {
			return err
		}

		return nil
	}

	high64k := m.mem[pageTableBase : pageTableBase+0x6000]

	// zero out the page tables.
	// but we might in fact want to poison them?
	// do we really want 1G, for example?
	for i := range high64k {
		high64k[i] = 0
	}

	// Set up page tables for long mode.
	// take the first six pages of an area it should not touch -- PageTableBase
	// present, read/write, page table at 0xffff0000
	// ptes[0] = PageTableBase + 0x1000 | 0x3
	// 3 in lowest 2 bits means present and read/write
	// 0x60 means accessed/dirty
	// 0x80 means the page size bit -- 0x80 | 0x60 = 0xe0
	// 0x10 here is making it point at the next page.
	// another go anti-pattern from golangci-lint.
	// golangci-lint claims this file has not been go-fumpt-ed
	// but it has.
	copy(high64k, []byte{
		0x03,
		0x10 | uint8((pageTableBase>>8)&0xff),
		uint8((pageTableBase >> 16) & 0xff),
		uint8((pageTableBase >> 24) & 0xff), 0, 0, 0, 0,
	})
	// need four pointers to 2M page tables -- PHYSICAL addresses:
	// 0x2000, 0x3000, 0x4000, 0x5000
	// experiment: set PS bit
	// Don't.
	for i := uint64(0); i < 4; i++ {
		ptb := pageTableBase + (i+2)*0x1000
		// Another coding anti-pattern
		copy(high64k[int(i*8)+0x1000:],
			[]byte{
				/*0x80 |*/ 0x63,
				uint8((ptb >> 8) & 0xff),
				uint8((ptb >> 16) & 0xff),
				uint8((ptb >> 24) & 0xff), 0, 0, 0, 0,
			})
	}
	// Now the 2M pages.
	for i := uint64(0); i < 0x1_0000_0000; i += 0x2_00_000 {
		ptb := i | 0xe3
		ix := int((i/0x2_00_000)*8 + 0x2000)
		// another coding anti-pattern from golangci-lint.
		copy(high64k[ix:], []byte{
			uint8(ptb),
			uint8((ptb >> 8) & 0xff),
			uint8((ptb >> 16) & 0xff),
			uint8((ptb >> 24) & 0xff), 0, 0, 0, 0,
		})
	}

	// set to true to debug.
	if false {
		log.Printf("Page tables: %s", hex.Dump(m.mem[pageTableBase:pageTableBase+0x3000]))
	}

	sregs.CR3 = uint64(pageTableBase)
	sregs.CR4 = CR4xPAE
	sregs.CR0 = CR0xPE | CR0xMP | CR0xET | CR0xNE | CR0xWP | CR0xAM | CR0xPG
	sregs.EFER = EFERxLME | EFERxLMA

	seg := kvm.Segment{
		Base:     0,
		Limit:    0xffffffff,
		Selector: 1 << 3,
		Typ:      11, /* Code: execute, read, accessed */
		Present:  1,
		DPL:      0,
		DB:       0,
		S:        1, /* Code/data */
		L:        1,
		G:        1, /* 4KB granularity */
		AVL:      0,
	}

	sregs.CS = seg

	seg.Typ = 3 /* Data: read/write, accessed */
	seg.Selector = 2 << 3
	sregs.DS, sregs.ES, sregs.FS, sregs.GS, sregs.SS = seg, seg, seg, seg, seg

	if err := kvm.SetSregs(vcpufd, sregs); err != nil {
		return err
	}

	return nil
}

func (m *Machine) initCPUID(cpu int) error {
	cpuid := kvm.CPUID{
		Nent:    100,
		Entries: make([]kvm.CPUIDEntry2, 100),
	}

	if err := kvm.GetSupportedCPUID(m.kvmFd, &cpuid); err != nil {
		return err
	}

	// https://www.kernel.org/doc/html/latest/virt/kvm/cpuid.html
	for i := 0; i < int(cpuid.Nent); i++ {
		if cpuid.Entries[i].Function == kvm.CPUIDFuncPerMon {
			cpuid.Entries[i].Eax = 0 // disable
		} else if cpuid.Entries[i].Function == kvm.CPUIDSignature {
			cpuid.Entries[i].Eax = kvm.CPUIDFeatures
			cpuid.Entries[i].Ebx = 0x4b4d564b // KVMK
			cpuid.Entries[i].Ecx = 0x564b4d56 // VMKV
			cpuid.Entries[i].Edx = 0x4d       // M
		}
	}

	if err := kvm.SetCPUID2(m.vcpuFds[cpu], &cpuid); err != nil {
		return err
	}

	return nil
}

// SingleStep enables single stepping the guest.
func (m *Machine) SingleStep(onoff bool) error {
	for cpu := range m.vcpuFds {
		if err := kvm.SingleStep(m.vcpuFds[cpu], onoff); err != nil {
			return fmt.Errorf("single step %d:%w", cpu, err)
		}
	}

	return nil
}

// RunInfiniteLoop runs the guest cpu until there is an error.
// If the error is ErrExitDebug, this function can be called again.
func (m *Machine) RunInfiniteLoop(cpu int) error {
	// https://www.kernel.org/doc/Documentation/virtual/kvm/api.txt
	// - vcpu ioctls: These query and set attributes that control the operation
	//   of a single virtual cpu.
	//
	//   vcpu ioctls should be issued from the same thread that was used to create
	//   the vcpu, except for asynchronous vcpu ioctl that are marked as such in
	//   the documentation.  Otherwise, the first ioctl after switching threads
	//   could see a performance impact.
	//
	// - device ioctls: These query and set attributes that control the operation
	//   of a single device.
	//
	//   device ioctls must be issued from the same process (address space) that
	//   was used to create the VM.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	for {
		isContinue, err := m.RunOnce(cpu)
		if isContinue {
			if err != nil {
				fmt.Printf("%v\r\n", err)
			}

			continue
		}

		if err != nil {
			return err
		}
	}
}

// RunOnce runs the guest vCPU until it exits.
func (m *Machine) RunOnce(cpu int) (bool, error) {
	fd, err := m.CPUToFD(cpu)
	if err != nil {
		return false, err
	}

	_ = kvm.Run(fd)
	exit := kvm.ExitType(m.runs[cpu].ExitReason)

	switch exit {
	case kvm.EXITHLT:
		return false, err
	case kvm.EXITIO:
		direction, size, port, count, offset := m.runs[cpu].IO()
		f := m.ioportHandlers[port][direction]

		bytes := (*(*[100]byte)(unsafe.Pointer(uintptr(unsafe.Pointer(m.runs[cpu])) + uintptr(offset))))[0:size]
		for i := 0; i < int(count); i++ {
			if err := f(port, bytes); err != nil {
				return false, err
			}
		}

		return true, err
	case kvm.EXITUNKNOWN:
		return true, err
	case kvm.EXITINTR:
		// When a signal is sent to the thread hosting the VM it will result in EINTR
		// refs https://gist.github.com/mcastelino/df7e65ade874f6890f618dc51778d83a
		return true, nil
	case kvm.EXITDEBUG:
		return false, kvm.ErrDebug

	case kvm.EXITDCR,
		kvm.EXITEXCEPTION,
		kvm.EXITFAILENTRY,
		kvm.EXITHYPERCALL,
		kvm.EXITINTERNALERROR,
		kvm.EXITIRQWINDOWOPEN,
		kvm.EXITMMIO,
		kvm.EXITNMI,
		kvm.EXITS390RESET,
		kvm.EXITS390SIEIC,
		kvm.EXITSETTPR,
		kvm.EXITSHUTDOWN,
		kvm.EXITTPRACCESS:
		if err != nil {
			return false, err
		}

		return false, fmt.Errorf("%w: %s", kvm.ErrUnexpectedExitReason, exit.String())
	default:
		if err != nil {
			return false, err
		}

		r, _ := m.GetRegs(cpu)
		s, _ := m.GetSRegs(cpu)
		// another coding anti-pattern from golangci-lint.
		return false, fmt.Errorf("%w: %v: regs:\n%s",
			kvm.ErrUnexpectedExitReason,
			kvm.ExitType(m.runs[cpu].ExitReason).String(), show("", &s, &r))
	}
}

func (m *Machine) registerIOPortHandler(
	start, end uint64,
	inHandler, outHandler func(port uint64, bytes []byte) error,
) {
	for i := start; i < end; i++ {
		m.ioportHandlers[i][kvm.EXITIOIN] = inHandler
		m.ioportHandlers[i][kvm.EXITIOOUT] = outHandler
	}
}

func (m *Machine) initIOPortHandlers() {
	funcNone := func(port uint64, bytes []byte) error {
		return nil
	}

	funcError := func(port uint64, bytes []byte) error {
		return fmt.Errorf("%w: unexpected io port 0x%x", kvm.ErrUnexpectedExitReason, port)
	}

	// 0xCF9 port can get three values for three types of reset:
	//
	// Writing 4 to 0xCF9:(INIT) Will INIT the CPU. Meaning it will jump
	// to the initial location of booting but it will keep many CPU
	// elements untouched. Most internal tables, chaches etc will remain
	// unchanged by the Init call (but may change during it).
	//
	// Writing 6 to 0xCF9:(RESET) Will RESET the CPU with all
	// internal tables caches etc cleared to initial state.
	//
	// Writing 0xE to 0xCF9:(RESTART) Will power cycle the mother board
	// with everything that comes with it.
	// For now, we will exit without regard to the value. Should we wish
	// to have more sophisticated cf9 handling, we will need to modify
	// gokvm a bit more.
	funcOutbCF9 := func(port uint64, bytes []byte) error {
		if len(bytes) == 1 && bytes[0] == 0xe {
			return fmt.Errorf("write 0xe to cf9: %w", ErrWriteToCF9)
		}

		return fmt.Errorf("write %#x to cf9: %w", bytes, ErrWriteToCF9)
	}

	// In ubuntu 20.04 on wsl2, the output to IO port 0x64 continued
	// infinitely. To deal with this issue, refer to kvmtool and
	// configure the input to the Status Register of the PS2 controller.
	//
	// refs:
	// https://github.com/kvmtool/kvmtool/blob/0e1882a49f81cb15d328ef83a78849c0ea26eecc/hw/i8042.c#L312
	// https://git.kernel.org/pub/scm/linux/kernel/git/will/kvmtool.git/tree/hw/i8042.c#n312
	// https://wiki.osdev.org/%228042%22_PS/2_Controller
	funcInbPS2 := func(port uint64, bytes []byte) error {
		bytes[0] = 0x20

		return nil
	}

	m.registerIOPortHandler(0, 0x10000, funcError, funcError)    // default handler
	m.registerIOPortHandler(0xcf9, 0xcfa, funcNone, funcOutbCF9) // CF9
	m.registerIOPortHandler(0x3c0, 0x3db, funcNone, funcNone)    // VGA
	m.registerIOPortHandler(0x3b4, 0x3b6, funcNone, funcNone)    // VGA
	m.registerIOPortHandler(0x2f8, 0x300, funcNone, funcNone)    // Serial port 2
	m.registerIOPortHandler(0x3e8, 0x3f0, funcNone, funcNone)    // Serial port 3
	m.registerIOPortHandler(0x2e8, 0x2f0, funcNone, funcNone)    // Serial port 4
	m.registerIOPortHandler(0xcfe, 0xcff, funcNone, funcNone)    // unknown
	m.registerIOPortHandler(0xcfa, 0xcfc, funcNone, funcNone)    // unknown
	m.registerIOPortHandler(0xc000, 0xd000, funcNone, funcNone)  // PCI Configuration Space Access Mechanism #2
	m.registerIOPortHandler(0x60, 0x70, funcInbPS2, funcNone)    // PS/2 Keyboard (Always 8042 Chip)
	m.registerIOPortHandler(0xed, 0xee, funcNone, funcNone)      // 0xed is the new standard delay port.

	// Serial port 1
	m.registerIOPortHandler(serial.COM1Addr, serial.COM1Addr+8, m.serial.In, m.serial.Out)

	// PCI configuration
	//
	// 0xcf8 for address register for PCI Config Space
	// 0xcfc + 0xcff for data for PCI Config Space
	// see https://github.com/torvalds/linux/blob/master/arch/x86/pci/direct.c for more detail.
	m.registerIOPortHandler(0xcf8, 0xcf9, m.pci.PciConfAddrIn, m.pci.PciConfAddrOut)
	m.registerIOPortHandler(0xcfc, 0xd00, m.pci.PciConfDataIn, m.pci.PciConfDataOut)

	// IO Devices - non PCI
	for _, dev := range m.devices {
		m.registerIOPortHandler(dev.IOPort(), dev.IOPort()+dev.Size(), dev.Read, dev.Write)
	}

	// PCI devices
	for _, dev := range m.pci.Devices {
		m.registerIOPortHandler(dev.IOPort(), dev.IOPort()+dev.Size(), dev.Read, dev.Write)
	}
}

// InjectSerialIRQ injects a serial interrupt.
func (m *Machine) InjectSerialIRQ() error {
	if err := kvm.IRQLineStatus(m.vmFd, serialIRQ, 0); err != nil {
		return err
	}

	if err := kvm.IRQLineStatus(m.vmFd, serialIRQ, 1); err != nil {
		return err
	}

	return nil
}

// InjectViortNetIRQ injects a virtio net interrupt.
func (m *Machine) InjectVirtioNetIRQ() error {
	if err := kvm.IRQLineStatus(m.vmFd, virtioNetIRQ, 0); err != nil {
		return err
	}

	if err := kvm.IRQLineStatus(m.vmFd, virtioNetIRQ, 1); err != nil {
		return err
	}

	return nil
}

// InjectViortNetIRQ injects a virtio block interrupt.
func (m *Machine) InjectVirtioBlkIRQ() error {
	if err := kvm.IRQLineStatus(m.vmFd, virtioBlkIRQ, 0); err != nil {
		return err
	}

	if err := kvm.IRQLineStatus(m.vmFd, virtioBlkIRQ, 1); err != nil {
		return err
	}

	return nil
}

// ReadAt implements io.ReadAt for the kvm guest pvh.
func (m *Machine) ReadAt(b []byte, off int64) (int, error) {
	mem := bytes.NewReader(m.mem)

	return mem.ReadAt(b, off)
}

// WriteAt implements io.WriteAt for the kvm guest pvh.
func (m *Machine) WriteAt(b []byte, off int64) (int, error) {
	if off > int64(len(m.mem)) {
		return 0, syscall.EFBIG
	}

	n := copy(m.mem[off:], b)

	return n, nil
}

func showone(indent string, in interface{}) string {
	var ret string

	s := reflect.ValueOf(in).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if f.Kind() == reflect.String {
			ret += fmt.Sprintf(indent+"%s %s = %s\n", typeOfT.Field(i).Name, f.Type(), f.Interface())
		} else {
			ret += fmt.Sprintf(indent+"%s %s = %#x\n", typeOfT.Field(i).Name, f.Type(), f.Interface())
		}
	}

	return ret
}

func show(indent string, l ...interface{}) string {
	var ret string
	for _, i := range l {
		ret += showone(indent, i)
	}

	return ret
}

// CPUToFD translates a CPU number to an fd.
func (m *Machine) CPUToFD(cpu int) (uintptr, error) {
	if cpu > len(m.vcpuFds) {
		return 0, fmt.Errorf("cpu %d out of range 0-%d:%w", cpu, len(m.vcpuFds), ErrBadCPU)
	}

	return m.vcpuFds[cpu], nil
}

// VtoP returns the physical address for a vCPU virtual address.
func (m *Machine) VtoP(cpu int, vaddr uint64) (int64, error) {
	fd, err := m.CPUToFD(cpu)
	if err != nil {
		return 0, err
	}

	t := &kvm.Translation{
		LinearAddress: vaddr,
	}
	if err := kvm.Translate(fd, t); err != nil {
		return -1, err
	}

	// There can exist a valid translation for memory that does not exist.
	// For now, we call that an error.
	if t.Valid == 0 || t.PhysicalAddress > uint64(len(m.mem)) {
		return -1, fmt.Errorf("%#x:valid not set:%w", vaddr, ErrBadVA)
	}

	return int64(t.PhysicalAddress), nil
}

// GetReg gets a pointer to a register in kvm.Regs, given
// a register number from reg. This used to be a comprehensive
// case, but golangci-lint disliked the cyclomatic complexity
// So we only show the few registers we support.
func GetReg(r *kvm.Regs, reg x86asm.Reg) (*uint64, error) {
	if reg == x86asm.RAX {
		return &r.RAX, nil
	}

	if reg == x86asm.RCX {
		return &r.RCX, nil
	}

	if reg == x86asm.RDX {
		return &r.RDX, nil
	}

	if reg == x86asm.RBX {
		return &r.RBX, nil
	}

	if reg == x86asm.RSP {
		return &r.RSP, nil
	}

	if reg == x86asm.RBP {
		return &r.RBP, nil
	}

	if reg == x86asm.RSI {
		return &r.RSI, nil
	}

	if reg == x86asm.RDI {
		return &r.RDI, nil
	}

	if reg == x86asm.R8 {
		return &r.R8, nil
	}

	if reg == x86asm.R9 {
		return &r.R9, nil
	}

	if reg == x86asm.R10 {
		return &r.R10, nil
	}

	if reg == x86asm.R11 {
		return &r.R11, nil
	}

	if reg == x86asm.R12 {
		return &r.R12, nil
	}

	if reg == x86asm.R13 {
		return &r.R13, nil
	}

	if reg == x86asm.R14 {
		return &r.R14, nil
	}

	if reg == x86asm.R15 {
		return &r.R15, nil
	}

	if reg == x86asm.RIP {
		return &r.RIP, nil
	}

	return nil, fmt.Errorf("register %v%w", reg, ErrUnsupported)
}

// InitKVM takes care of the general kvm setup without dependencies to runtime target.
func initVMandVCPU(
	kvmPath string,
	nCpus int,
) (uintptr, uintptr, []uintptr, []*kvm.RunData, error) {
	var err error

	devKVM, err := os.OpenFile(kvmPath, os.O_RDWR, 0o644)
	if err != nil {
		return 0, 0, nil, nil, err
	}

	kvmFd := devKVM.Fd()
	vmFd := uintptr(0)
	vcpuFds := make([]uintptr, nCpus)
	runs := make([]*kvm.RunData, nCpus)

	if vmFd, err = kvm.CreateVM(kvmFd); err != nil {
		return 0, 0, nil, nil, fmt.Errorf("CreateVM: %w", err)
	}

	if err := kvm.SetTSSAddr(vmFd, pvh.KVMTSSStart); err != nil {
		return 0, 0, nil, nil, err
	}

	if err := kvm.SetIdentityMapAddr(vmFd, pvh.KVMIdentityMapStart); err != nil {
		return 0, 0, nil, nil, err
	}

	if err := kvm.CreateIRQChip(vmFd); err != nil {
		return 0, 0, nil, nil, err
	}

	if err := kvm.CreatePIT2(vmFd); err != nil {
		return 0, 0, nil, nil, err
	}

	mmapSize, err := kvm.GetVCPUMMmapSize(kvmFd)
	if err != nil {
		return 0, 0, nil, nil, err
	}

	for cpu := 0; cpu < nCpus; cpu++ {
		// Create vCPU
		vcpuFds[cpu], err = kvm.CreateVCPU(vmFd, cpu)
		if err != nil {
			return 0, 0, nil, nil, err
		}

		// init kvm_run structure
		r, err := syscall.Mmap(int(vcpuFds[cpu]), 0, int(mmapSize),
			syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
		if err != nil {
			return 0, 0, nil, nil, err
		}

		runs[cpu] = (*kvm.RunData)(unsafe.Pointer(&r[0]))
	}

	return kvmFd, vmFd, vcpuFds, runs, nil
}

func (m *Machine) StartVCPU(cpu, traceCount int, wg *sync.WaitGroup) {
	trace := traceCount > 0

	go func(cpu int) {
		var err error
		// Consider ANOTHER option, maxInsCount, which would
		// exit this loop after a certain number of instructions
		// were run.
		for tc := 0; ; tc++ {
			err = m.RunInfiniteLoop(cpu)
			if err == nil {
				continue
			}

			if !errors.Is(err, kvm.ErrDebug) {
				fmt.Printf("err: %v\r\n", err)

				break
			}

			if err := m.SingleStep(trace); err != nil {
				fmt.Printf("Setting trace to %v:%v", trace, err)
			}

			if tc%traceCount != 0 {
				continue
			}

			_, r, s, err := m.Inst(cpu)
			if err != nil {
				fmt.Printf("disassembling after debug exit:%v", err)
			} else {
				fmt.Printf("%#x:%s\r\n", r.RIP, s)
			}
		}

		wg.Done()
		fmt.Printf("CPU %d exits\n\r", cpu)
	}(cpu)
}

func (m *Machine) GetSerial() *serial.Serial {
	return m.serial
}

func (m *Machine) AddDevice(dev iodev.Device) {
	m.devices = append(m.devices, dev)
}
