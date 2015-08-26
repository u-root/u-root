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
type kexec_segment struct {
	buf     uintptr
	bufsz   uintptr
	mem     uintptr
	memsize uintptr
}

func (k *kexec_segment)String() string{
	return fmt.Sprintf("[0x%x 0x%x 0x%x 0x%x]", k.buf, k.bufsz, k.mem, k.memsize)
}

type loader func(b[]byte) (uintptr, []kexec_segment, error)
var (
	kernel = flag.String("k", "vmlinux", "Kernel to load")
	loaders = []loader{elfexec, rawexec}
)

/* Load a new kernel image as described by the kexec_segment array
 * consisting of passed number of segments at the entry-point address.
 * The flags allow different useage types.
 */
//extern int kexec_load(void *, size_t, struct kexec_segment *,
//		unsigned long int);

// rawexec either succeeds or Fatals out. It's the last chance.
func rawexec(b []byte) (uintptr, []kexec_segment, error) {
	var entry uintptr
	segs := []kexec_segment{{buf: uintptr(unsafe.Pointer(&b[0])), bufsz: uintptr(len(b)), mem: 0x1000, memsize: uintptr(len(b))}}
	return entry, segs, nil
}

func elfexec(b []byte) (uintptr, []kexec_segment, error) {
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
	segs := make([]kexec_segment, scount)
	for i, v := range f.Progs {
		if v.Type.String() == "PT_LOAD" {
			f := v.Open()
			// kexec is stupid. Requires the length to be page aligned.
			// And, guess what -- that's not always a known thing. But of course
			// we all know it's 4K. BAH!
			l := uintptr((v.Memsz + 4095)) & (uintptr(^uintptr(4095)))
			b := make([]byte, l)
			if _, err := f.Read(b[:v.Filesz]); err != nil {
				log.Fatalf("Reading %d bytes of program header %d: %v", v.Filesz, i, err)
			}
			segs[i].buf = uintptr(unsafe.Pointer(&b[0]))
			segs[i].bufsz = uintptr(l)
			segs[i].mem = uintptr(v.Paddr)
			segs[i].memsize = uintptr(l)
		}
	}
	return uintptr(f.Entry), segs, nil
}
func main() {
	var err error
	var segs []kexec_segment
	var entry uintptr
	flag.Parse()
	b, err := ioutil.ReadFile(*kernel)
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("Loading %v\n", *kernel)
	for i := range loaders {
		if entry, segs, err = loaders[i](b); err == nil {
			break
		}
	}
	for _, s := range segs {
		log.Printf("%v", s.String())
	}

	log.Printf("%v %v %v %v %v", syscall.SYS_KEXEC_LOAD, entry, uintptr(len(segs)), uintptr(unsafe.Pointer(&segs[0])), uintptr(0))
	e1, e2, err := syscall.Syscall6(syscall.SYS_KEXEC_LOAD, entry, uintptr(len(segs)), uintptr(unsafe.Pointer(&segs[0])), 0, 0, 0)
	log.Printf("a %v b %v err %v", e1, e2, err)
}
