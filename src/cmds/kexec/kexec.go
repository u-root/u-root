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
	"flag"
	"io"
	"log"
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
type kexec_segment struct {
	buf     []byte
	bufsz   uint64
	mem     uint64
	memsize uint64
}

var (
	kernel = flag.String("k", "vmlinux", "Kernel to load")
)

/* Load a new kernel image as described by the kexec_segment array
 * consisting of passed number of segments at the entry-point address.
 * The flags allow different useage types.
 */
//extern int kexec_load(void *, size_t, struct kexec_segment *,
//		unsigned long int);

func main() {
	flag.Parse()
	log.Printf("Loading %v\n", *kernel)
	f, err := elf.Open(*kernel)
	if err != nil {
		log.Fatalf("%v", err)
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
	segs := make([]kexec_segment, scount)
	scount = 0
	for _, v := range f.Progs {
		if v.Type.String() == "PT_LOAD" {
			f := v.Open()
			b := bytes.NewBuffer([]byte{})
			io.Copy(b, f)
			segs[scount].buf = b.Bytes()
			segs[scount].bufsz = uint64(b.Len())
			segs[scount].mem = v.Paddr
			segs[scount].memsize = v.Memsz
			scount++
		}
	}
}
