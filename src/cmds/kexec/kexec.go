// Copyright 2015 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// kexec command in Go.
// Wow, Linux. Just Wow. Who designs this stuff? Because, basically, it's not very good.
// I considered keeping the C declarations and such stuff here but it has not changed
// much over the years, so that seemed pointless.
// For a good design, you can always see the console device in Plan 9.
// E.g.,
// echo reboot path/to/kernel > /dev/consctl
// There, that's it, done. Is it so hard? Especially given that Linux already has an
// elf parser?
package main

// N.B. /**/ comments are verbatim from uapi/linux/kexec.h.
/* kexec system call -  It loads the new kernel to boot into.
 * kexec does not sync, or unmount filesystems so if you need
 * that to happen you need to do that yourself.
 */

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"syscall"
	"unsafe"
)

/* kexec flags for different usage scenarios */
const (
	KEXEC_ON_CRASH         = 0x00000001
	KEXEC_PRESERVE_CONTEXT = 0x00000002
	KEXEC_ARCH_MASK        = 0xffff0000

	/* These values match the ELF architecture values.
	 * Unless there is a good reason that should continue to be the case.
	 */
	KEXEC_ARCH_DEFAULT = (0 << 16)
	KEXEC_ARCH_386     = (3 << 16)
	KEXEC_ARCH_68K     = (4 << 16)
	KEXEC_ARCH_X86_64  = (62 << 16)
	KEXEC_ARCH_PPC     = (20 << 16)
	KEXEC_ARCH_PPC64   = (21 << 16)
	KEXEC_ARCH_IA_64   = (50 << 16)
	KEXEC_ARCH_ARM     = (40 << 16)
	KEXEC_ARCH_S390    = (22 << 16)
	KEXEC_ARCH_SH      = (42 << 16)
	KEXEC_ARCH_MIPS_LE = (10 << 16)
	KEXEC_ARCH_MIPS    = (8 << 16)

	/* The artificial cap on the number of segments passed to kexec_load. */
	KEXEC_SEGMENT_MAX = 16
)

/*
 * This structure is used to hold the arguments that are used when
 * loading  kernel binaries.
 */
type KexecSegment struct {
	buf     uintptr
	bufsz   uintptr
	mem     uintptr
	memsize uintptr
}

func (k *KexecSegment) String() string {
	return fmt.Sprintf("[0x%x 0x%x 0x%x 0x%x]", k.buf, k.bufsz, k.mem, k.memsize)
}

type loader func(b []byte) (uintptr, []KexecSegment, error)

var (
	dryrun  = flag.Bool("dryrun", false, "Do not do kexec system calls")
	loaders = []loader{elfexec, bzImage, rawexec}
	tramp []byte
)

func pagealloc(len int) []byte {
	flags := syscall.MAP_PRIVATE | syscall.MAP_ANON
	prot := syscall.PROT_READ | syscall.PROT_WRITE
	len = ((len + 4095) >> 12) << 12
	b, err := syscall.Mmap(-1, 0, len, prot, flags)
	if err != nil {
		log.Fatal("Could not mmap %d bytes: %v", len, err)
	}
	return b
}

func makeseg(b []byte, paddr uintptr) KexecSegment {
	if b == nil {
		panic("bad b")
	}
	return KexecSegment{
		buf:     uintptr(unsafe.Pointer(&b[0])),
		bufsz:   uintptr(len(b)),
		mem:     paddr,
		memsize: uintptr(len(b)),
	}
}

/* kexec requires page-aligned, page-sized buffers. It also assumes that means 4k. Oh well. */
func pages(size uintptr) []byte {
	size = ((size + 4095) >> 12) << 12
	return make([]byte, size)
}

/* Load a new kernel image as described by the KexecSegment array
 * consisting of passed number of segments at the entry-point address.
 * The flags allow different useage types.
 */
//extern int kexec_load(void *, size_t, struct KexecSegment *,
//		unsigned long int);

// rawexec either succeeds or Fatals out. It's the last chance.
func rawexec(b []byte) (uintptr, []KexecSegment, error) {
	var entry uintptr
	segs := []KexecSegment{
		//	makeseg(b, 0x1000)
		makeseg(tramp, TrampolinePointer),
	}
	entry = TrampolinePointer
	log.Printf("Using raw image loader")
	return entry, segs, nil
}

func bzImage(b []byte) (uintptr, []KexecSegment, error) {
	var h LinuxHeader
	r := bytes.NewReader(b)
	binary.Read(r, binary.LittleEndian, &h)
	log.Printf("bzImage header %v", h)
	log.Printf("magic %x switch %v", h.HeaderMagic, h.RealModeSwitch)
	if h.HeaderMagic != HeaderMagic {
		return 0, nil, errors.New("Not a bzImage")
	}
	log.Printf("RamDisk image %x size %x", h.RamdiskImage, h.RamdiskSize)
	log.Printf("StartSys %x", h.StartSys)
	log.Printf("Boot type: %s(%x)", LoaderType[boottype(h.TypeOfLoader)], h.TypeOfLoader)
	log.Printf("SetupSects %d", h.SetupSects)

	kernel := pagealloc(len(b))
	copy(kernel, b[(int(h.SetupSects)+1)*512:])
	kernelBase := uintptr(0x100000)

	// Now for the good fun.
	l := &LinuxParams{
		MountRootReadonly: h.RootFlags,
		OrigRootDev:       h.RootDev,
		OrigVideoMode:     3,
		OrigVideoCols:     80,
		OrigVideoLines:    25,
		OrigVideoIsVGA:    1,
		OrigVideoPoints:   16,
		LoaderType:        0xff,
		CLPtr:             CommandLinePointer,
		KernelStart:       0x100000,
	}
	w := pagealloc(1)
	// binary.Write may reallocate the buffer, but we are pretty sure
	// it won't.
	binary.Write(bytes.NewBuffer(w), binary.LittleEndian, l)

	cmdline := pagealloc(1)
	segs := []KexecSegment{
		makeseg(w, LinuxParamPointer),
		makeseg(tramp, TrampolinePointer),
		makeseg(kernel, kernelBase),
		makeseg(cmdline, CommandLinePointer),
		//makeseg(initrd, initrd_base),
	}

	return TrampolinePointer, segs, nil
}

func elfexec(b []byte) (uintptr, []KexecSegment, error) {
	f, err := elf.NewFile(bytes.NewReader(b))
	if err != nil {
		return 0, nil, err
	}
	scount := 0
	for _, v := range f.Progs {
		if v.Type.String() == "PT_LOAD" {
			scount++
		}
	}
	if scount > KEXEC_SEGMENT_MAX {
		log.Fatalf("Too many segments: got %v, max is %v", scount, KEXEC_SEGMENT_MAX)
	}
	segs := make([]KexecSegment, scount)
	for i, v := range f.Progs {
		if v.Type.String() == "PT_LOAD" {
			f := v.Open()
			b := pages(uintptr(v.Memsz))
			if _, err := f.Read(b[:v.Filesz]); err != nil {
				log.Fatalf("Reading %d bytes of program header %d: %v", v.Filesz, i, err)
			}
			segs[i] = makeseg(b, uintptr(v.Paddr))
		}
	}
	log.Printf("Using ELF image loader")
	return uintptr(f.Entry), segs, nil
}
func main() {
	kernel := "bzImage"
	var err error
	var segs []KexecSegment
	var entry uintptr
	flag.Parse()
	if len(flag.Args()) == 1 {
		kernel = flag.Args()[0]
	}

	b, err := ioutil.ReadFile(kernel)
	if err != nil {
		log.Fatalf("%v", err)
	}

	tramp = pagealloc(len(trampoline))
	copy(tramp, trampoline)

	log.Printf("Loading %v\n", kernel)
	for i := range loaders {
		if entry, segs, err = loaders[i](b); err == nil {
			break
		}
	}
	for _, s := range segs {
		log.Printf("%v", s.String())
	}

	if *dryrun {
		log.Printf("Dry run -- exiting now")
		return
	}
	log.Printf("%v %v %v %v %v", syscall.SYS_KEXEC_LOAD, entry, uintptr(len(segs)), uintptr(unsafe.Pointer(&segs[0])), uintptr(0))
	e1, e2, err := syscall.Syscall6(syscall.SYS_KEXEC_LOAD, entry, uintptr(len(segs)), uintptr(unsafe.Pointer(&segs[0])), 0, 0, 0)
	log.Printf("a %v b %v err %v", e1, e2, err)

	e1, e2, err = syscall.Syscall6(syscall.SYS_REBOOT, syscall.LINUX_REBOOT_MAGIC1, syscall.LINUX_REBOOT_MAGIC2, syscall.LINUX_REBOOT_CMD_KEXEC, 0, 0, 0)

	log.Printf("a %v b %v err %v", e1, e2, err)
}
