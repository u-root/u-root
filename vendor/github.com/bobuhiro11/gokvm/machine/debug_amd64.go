package machine

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/bobuhiro11/gokvm/kvm"
	"golang.org/x/arch/x86/x86asm"
)

// Debug is a normally empty function that enables debug prints.
// well too bad. var debug = log.Printf // func(string, ...interface{}) {}

// ErrBadRegister indicates a bad register was used.
var ErrBadRegister = errors.New("bad register")

// ErrBadArg indicates an argument number was out of range.
var ErrBadArg = errors.New("arg count must be in range 1..6")

// ErrBadArgType indicates an argument type is not correct,
// e.g. code expected a Mem but got an Imm.
var ErrBadArgType = errors.New("bad arg type")

// Args returns the top nargs args, going down the stack if needed. The max is 6.
// This is UEFI calling convention.
func (m *Machine) Args(cpu int, r *kvm.Regs, nargs int) ([]uintptr, error) {
	// We must always validate the cpu number, even if we don't absolutely need it.
	// i.e., for pure register args, the cpu is not needed, but it's
	// best to validate it.
	if _, err := m.CPUToFD(cpu); err != nil {
		return nil, err
	}

	sp := r.RSP

	switch nargs {
	case 6:
		w1, err := m.ReadWord(cpu, sp+0x28)
		if err != nil {
			return nil, err
		}

		w2, err := m.ReadWord(cpu, sp+0x30)
		if err != nil {
			return nil, err
		}

		return []uintptr{uintptr(r.RCX), uintptr(r.RDX), uintptr(r.R8), uintptr(r.R9), uintptr(w1), uintptr(w2)}, nil
	case 5:
		w1, err := m.ReadWord(cpu, sp+0x28)
		if err != nil {
			return nil, err
		}

		return []uintptr{uintptr(r.RCX), uintptr(r.RDX), uintptr(r.R8), uintptr(r.R9), uintptr(w1)}, nil
	case 4:
		return []uintptr{uintptr(r.RCX), uintptr(r.RDX), uintptr(r.R8), uintptr(r.R9)}, nil
	case 3:
		return []uintptr{uintptr(r.RCX), uintptr(r.RDX), uintptr(r.R8)}, nil
	case 2:
		return []uintptr{uintptr(r.RCX), uintptr(r.RDX)}, nil
	case 1:
		return []uintptr{uintptr(r.RCX)}, nil
	}

	return []uintptr{}, fmt.Errorf("args(%d):%w", nargs, ErrBadArg)
}

// Pointer returns the data pointed to by args[arg].
func (m *Machine) Pointer(inst *x86asm.Inst, r *kvm.Regs, arg uint) (uintptr, error) {
	if arg >= uint(len(inst.Args)) {
		return 0, fmt.Errorf("pointer(..,%d): only %d args:%w", arg, len(inst.Args), ErrBadArgType)
	}

	mem, ok := inst.Args[arg].(x86asm.Mem)
	if !ok {
		return 0, fmt.Errorf("arg %d is not a memory argument:%w", arg, ErrBadArgType)
	}
	// A Mem is a memory reference.
	// The general form is Segment:[Base+Scale*Index+Disp].
	/*
		type Mem struct {
			Segment Reg
			Base    Reg
			Scale   uint8
			Index   Reg
			Disp    int64
		}
	*/
	// debug("ARG[%d] %q m is %#x", inst.Args[arg], mem)

	b, err := GetReg(r, mem.Base)
	if err != nil {
		return 0, fmt.Errorf("base reg %v in %v:%w", mem.Base, mem, ErrBadRegister)
	}

	addr := *b + uint64(mem.Disp)

	x, err := GetReg(r, mem.Index)
	if err == nil {
		addr += uint64(mem.Scale) * (*x)
	}

	// if v, ok := inst.Args[0].(*x86asm.Mem); ok {
	// debug("computed addr is %#x", addr)

	return uintptr(addr), nil
}

// Pop pops the stack and returns what was at TOS.
// It is most often used to get the caller PC (cpc).
func (m *Machine) Pop(cpu int, r *kvm.Regs) (uint64, error) {
	cpc, err := m.ReadWord(cpu, r.RSP)
	if err != nil {
		return 0, err
	}

	r.RSP += 8

	return cpc, nil
}

// Inst retrieves an instruction from the guest, at RIP.
// It returns an x86asm.Inst, Ptraceregs, a string in GNU syntax,
// and error.
func (m *Machine) Inst(cpu int) (*x86asm.Inst, *kvm.Regs, string, error) {
	r, err := m.GetRegs(cpu)
	if err != nil {
		return nil, nil, "", fmt.Errorf("Inst:Getregs:%w", err)
	}

	pc := r.RIP

	// debug("Inst: pc %#x, sp %#x", pc, sp)
	// We know the PC; grab a bunch of bytes there, then decode and print
	insn := make([]byte, 16)
	if _, err := m.ReadBytes(cpu, insn, pc); err != nil {
		return nil, nil, "", fmt.Errorf("reading PC at #%x:%w", pc, err)
	}

	d, err := x86asm.Decode(insn, 64)
	if err != nil {
		return nil, nil, "", fmt.Errorf("decoding %#02x:%w", insn, err)
	}

	return &d, r, x86asm.GNUSyntax(d, r.RIP, nil), nil
}

// Asm returns a string for the given instruction at the given pc.
func Asm(d *x86asm.Inst, pc uint64) string {
	return "\"" + x86asm.GNUSyntax(*d, pc, nil) + "\""
}

// CallInfo provides calling info for a function.
func CallInfo(inst *x86asm.Inst, r *kvm.Regs) string {
	l := fmt.Sprintf("%s[", show("", r))
	for _, a := range inst.Args {
		l += fmt.Sprintf("%v,", a)
	}

	l += fmt.Sprintf("(%#x, %#x, %#x, %#x)", r.RCX, r.RDX, r.R8, r.R9)

	return l
}

// WriteWord writes the given word into the guest's virtual address space.
func (m *Machine) WriteWord(cpu int, vaddr uint64, word uint64) error {
	pa, err := m.VtoP(cpu, vaddr)
	if err != nil {
		return err
	}

	var b [8]byte

	binary.LittleEndian.PutUint64(b[:], word)
	_, err = m.WriteAt(b[:], pa)

	return err
}

// ReadBytes reads bytes from the CPUs virtual address space.
func (m *Machine) ReadBytes(cpu int, b []byte, vaddr uint64) (int, error) {
	pa, err := m.VtoP(cpu, vaddr)
	if err != nil {
		return -1, err
	}

	return m.ReadAt(b, pa)
}

// ReadWord reads the given word from the cpu's virtual address space.
func (m *Machine) ReadWord(cpu int, vaddr uint64) (uint64, error) {
	var b [8]byte
	if _, err := m.ReadBytes(cpu, b[:], vaddr); err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(b[:]), nil
}
