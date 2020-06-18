// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bzImage implements decoding for bzImage files.
//
// The bzImage struct contains all the information about the file and can
// be used to create a new bzImage.
package bzimage

// xz --check=crc32 $BCJ --lzma2=$LZMA2OPTS,dict=32MiB

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"reflect"
	"strings"

	"github.com/u-root/u-root/pkg/cpio"
)

var (
	xzmagic = [...]byte{0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00}
	// String of unknown meaning.
	// The build script has this value:
	//initRAMFStag = [4]byte{0250, 0362, 0156, 0x01}
	// The resultant bzd has this value:
	initRAMFStag = [4]byte{0xf8, 0x85, 0x21, 0x01}
	Debug        = func(string, ...interface{}) {}
)

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// For now, it hardwires the KernelBase to 0x100000.
// bzImages were created by a process of evilution, and they are wondrous to behold.
// bzImages are almost impossible to modify. They form a sandwich with
// the compressed kernel code in the middle. It's actually a BLT:
// MBR and bootparams first 512 bytes
//    the MBR includes 0xc0 bytes of boot code which is vestigial.
// Then there is "preamble" code which is the kernel decompressor; then the
// xz compressed kernel; then a library of sorts after the kernel which is called
// by the early uncompressed kernel code. This is all linked together and forms
// an essentially indivisible whole -- which we wish to divisible.
// Hence the groveling around for the xz header that you see here, and hence the checks
// to ensure that the kernel layout ends up largely the same before and after.
// That said, if you keep layout unchanged, you can modify the uncompressed kernel.
// For example, when you first build a kernel, you can:
// dd if=/dev/urandom of=x bs=1048576 count=8
// echo x | cpio -o > x.cpio
// and use that as an initrd, it's more or less an 8 MiB block you can replace
// as needed. Just make sure nothing grows. And make sure the initramfs is in
// the same place. Ah, joy.
func (b *BzImage) UnmarshalBinary(d []byte) error {
	Debug("Processing %d byte image", len(d))
	r := bytes.NewBuffer(d)
	if err := binary.Read(r, binary.LittleEndian, &b.Header); err != nil {
		return err
	}
	Debug("Header was %d bytes", len(d)-r.Len())
	Debug("magic %x switch %v", b.Header.HeaderMagic, b.Header.RealModeSwitch)
	if b.Header.HeaderMagic != HeaderMagic {
		return fmt.Errorf("not a bzImage: magic should be %02x, and is %02x", HeaderMagic, b.Header.HeaderMagic)
	}
	Debug("RamDisk image %x size %x", b.Header.RamDiskImage, b.Header.RamDiskSize)
	Debug("StartSys %x", b.Header.StartSys)
	Debug("Boot type: %s(%x)", LoaderType[boottype(b.Header.TypeOfLoader)], b.Header.TypeOfLoader)
	Debug("SetupSects %d", b.Header.SetupSects)

	off := len(d) - r.Len()
	b.KernelOffset = (uintptr(b.Header.SetupSects) + 1) * 512
	bclen := int(b.KernelOffset) - off
	Debug("Kernel offset is %d bytes, low1mcode is %d bytes", b.KernelOffset, bclen)
	b.BootCode = make([]byte, bclen)
	if _, err := r.Read(b.BootCode); err != nil {
		return err
	}
	Debug("%d bytes of BootCode", len(b.BootCode))

	Debug("Remaining length is %d bytes, PayloadSize %d", r.Len(), b.Header.PayloadSize)
	x := bytes.Index(r.Bytes(), xzmagic[:])
	if x == -1 {
		return fmt.Errorf("can't find xz header")
	}
	Debug("xz is at %d", x)
	b.HeadCode = make([]byte, x)
	if _, err := r.Read(b.HeadCode); err != nil {
		return fmt.Errorf("can't read HeadCode: %v", err)
	}
	// Now size up the kernel code. Is it just PayloadSize?
	b.compressed = make([]byte, b.Header.PayloadSize)
	if _, err := r.Read(b.compressed); err != nil {
		return fmt.Errorf("can't read HeadCode: %v", err)
	}
	var err error
	Debug("Uncompress %d bytes", len(b.compressed))
	if b.KernelCode, err = unpack(b.compressed); err != nil {
		return err
	}
	b.TailCode = make([]byte, r.Len())
	if _, err := r.Read(b.TailCode); err != nil {
		return fmt.Errorf("can't read TailCode: %v", err)
	}
	Debug("Kernel at %d, %d bytes", b.KernelOffset, len(b.KernelCode))
	b.KernelBase = uintptr(0x100000)
	if b.Header.RamDiskImage == 0 {
		return nil
	}
	if r.Len() != 0 {
		return fmt.Errorf("%d bytes left over", r.Len())
	}
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (b *BzImage) MarshalBinary() ([]byte, error) {
	// First step, make sure we can compress the kernel.
	dat, err := compress(b.KernelCode, "--lzma2=,dict=32MiB")
	if err != nil {
		return nil, err
	}
	dat = append(dat, initRAMFStag[:]...)
	if len(dat) > len(b.compressed) {
		return nil, fmt.Errorf("marshal: compressed KernelCode too big: was %d, now %d", len(b.compressed), len(dat))
	}
	Debug("b.compressed len %#x dat len %#x pad it out", len(b.compressed), len(dat))
	if len(dat) < len(b.compressed) {
		l := len(dat)
		n := make([]byte, len(b.compressed)-4)
		copy(n, dat[:l-4])
		n = append(n, initRAMFStag[:]...)
		dat = n
	}

	var w bytes.Buffer
	if err := binary.Write(&w, binary.LittleEndian, &b.Header); err != nil {
		return nil, err
	}
	Debug("Wrote %d bytes of header", w.Len())
	if _, err := w.Write(b.BootCode); err != nil {
		return nil, err
	}
	Debug("Wrote %d bytes of BootCode", w.Len())
	if _, err := w.Write(b.HeadCode); err != nil {
		return nil, err
	}
	Debug("Wrote %d bytes of HeadCode", w.Len())
	if _, err := w.Write(dat); err != nil {
		return nil, err
	}
	Debug("Last bytes %#02x", dat[len(dat)-4:])
	Debug("Last bytes %#o", dat[len(dat)-4:])
	Debug("Last bytes %#d", dat[len(dat)-4:])
	b.compressed = dat
	Debug("Wrote %d bytes of Compressed kernel", w.Len())
	if _, err := w.Write(b.TailCode); err != nil {
		return nil, err
	}
	Debug("Wrote %d bytes of header", w.Len())
	Debug("Finished writing, len is now %d bytes", w.Len())

	return w.Bytes(), nil
}

// unpack extracts the header code and data from the kernel part
// of the bzImage. It also uncompresses the kernel.
// It searches the Kernel []byte for an xz header. Where it begins
// is never certain. We only do relatively newer images, i.e. we only
// look for the xz magic.
func unpack(d []byte) ([]byte, error) {
	Debug("Kernel is %d bytes", len(d))
	Debug("Some kernel data: %#02x %#02x", d[:32], d[len(d)-8:])
	c := exec.Command("xzcat")
	stdout, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		return nil, err
	}
	c.Stdin = bytes.NewBuffer(d)
	if err := c.Start(); err != nil {
		return nil, err
	}

	dat, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	// fyi, the xz standard and code are shit. A shame.
	// You can enable this if you have a nasty bug from xz.
	// Just be aware that xz ALWAYS errors out even when nothing is wrong.
	if false {
		if e, err := ioutil.ReadAll(stderr); err != nil || len(e) > 0 {
			Debug("xz stderr: '%s', %v", string(e), err)
		}
	}
	Debug("Uncompressed kernel is %d bytes", len(dat))
	return dat, nil
}

// compress compresses a []byte via xz using the dictOps, collecting it from stdout
func compress(b []byte, dictOps string) ([]byte, error) {
	Debug("b is %d bytes", len(b))
	c := exec.Command("xz", "--check=crc32", "--x86", dictOps, "--stdout")
	stdout, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}
	c.Stdin = bytes.NewBuffer(b)
	if err := c.Start(); err != nil {
		return nil, err
	}

	dat, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	if err := c.Wait(); err != nil {
		return nil, err
	}
	Debug("Compressed data is %d bytes, starts with %#02x", len(dat), dat[:32])
	Debug("Last 16 bytes: %#02x", dat[len(dat)-16:])
	return dat, nil
}

// Extract extracts the KernelCode as an ELF.
func (b *BzImage) ELF() (*elf.File, error) {
	e, err := elf.NewFile(bytes.NewReader(b.KernelCode))
	if err != nil {
		return nil, err
	}
	return e, nil
}

func Equal(a, b []byte) error {
	if len(a) != len(b) {
		return fmt.Errorf("images differ in len: %d bytes and %d bytes", len(a), len(b))
	}
	var ba BzImage
	if err := ba.UnmarshalBinary(a); err != nil {
		return err
	}
	var bb BzImage
	if err := bb.UnmarshalBinary(b); err != nil {
		return err
	}
	if !reflect.DeepEqual(ba.Header, bb.Header) {
		return fmt.Errorf("headers do not match: %s", ba.Header.Diff(&bb.Header))
	}
	// this is overkill, I can't see any way it can happen.
	if len(ba.KernelCode) != len(bb.KernelCode) {
		return fmt.Errorf("kernel lengths differ: %d vs %d bytes", len(ba.KernelCode), len(bb.KernelCode))
	}
	if len(ba.BootCode) != len(bb.BootCode) {
		return fmt.Errorf("boot code lengths differ: %d vs %d bytes", len(ba.KernelCode), len(bb.KernelCode))
	}

	if !reflect.DeepEqual(ba.BootCode, bb.BootCode) {
		return fmt.Errorf("boot code does not match")
	}
	if !reflect.DeepEqual(ba.KernelCode, bb.KernelCode) {
		return fmt.Errorf("kernels do not match")
	}
	return nil
}

func (b *BzImage) AddInitRAMFS(name string) error {
	u, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	// Should we ever want to compress the initramfs this is one
	// way to do it.
	d := u
	if false {
		d, err = compress(u, "--lzma2=,dict=1MiB")
		if err != nil {
			return err
		}
	}
	s, e, err := b.InitRAMFS()
	if err != nil {
		return err
	}
	l := e - s

	if len(d) > l {
		return fmt.Errorf("new initramfs is %d bytes, won't fit in %d byte old one", len(d), l)
	}
	// Do this in a stupid way that is easy to read.
	// What's interesting: the kernel decompressor, if I read it right,
	// finds it easier to skip a bunch of leading nulls. So do that.
	n := make([]byte, l)
	Debug("Offset into n is %d\n", len(n)-len(d))
	copy(n[len(n)-len(d):], d)
	Debug("Install %d byte initramfs in %d bytes of kernel code, @ %d:%d", len(d), len(n), s, e)
	copy(b.KernelCode[s:e], n)
	return nil
}

// MakeLinuxHeader marshals a LinuxHeader into a []byte.
func MakeLinuxHeader(h *LinuxHeader) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, h)
	return buf.Bytes(), err
}

// Show stringifies a LinuxHeader into a []string
func (h *LinuxHeader) Show() []string {
	var s []string

	val := reflect.ValueOf(*h)
	for i := 0; i < val.NumField(); i++ {
		v := val.Field(i)
		k := reflect.ValueOf(v).Kind()
		n := val.Type().Field(i).Name
		switch k {
		case reflect.Slice:
			s = append(s, fmt.Sprintf("%s:%#02x", n, v))
		case reflect.Bool:
			s = append(s, fmt.Sprintf("%s:%v", n, v))
		default:
			s = append(s, fmt.Sprintf("%s:%#02x", n, v))
		}
	}
	return s
}

// Diff is a convenience function that returns a string showing
// differences between a header and another header.
func (h *LinuxHeader) Diff(i *LinuxHeader) string {
	var s string
	hs := h.Show()
	is := i.Show()
	for i := range hs {
		if hs[i] != is[i] {
			s += fmt.Sprintf("%s != %s", hs[i], is[i])
		}
	}
	return s
}

// Diff is a convenience function that returns a string showing
// differences between a bzImage and another bzImage
func (b *BzImage) Diff(b2 *BzImage) string {
	s := b.Header.Diff(&b2.Header)
	if len(b.BootCode) != len(b2.BootCode) {
		s = s + fmt.Sprintf("b Bootcode is %d; b2 BootCode is %d", len(b.BootCode), len(b2.BootCode))
	}
	if len(b.HeadCode) != len(b2.HeadCode) {
		s = s + fmt.Sprintf("b Headcode is %d; b2 HeadCode is %d", len(b.HeadCode), len(b2.HeadCode))
	}
	if len(b.KernelCode) != len(b2.KernelCode) {
		s = s + fmt.Sprintf("b Kernelcode is %d; b2 KernelCode is %d", len(b.KernelCode), len(b2.KernelCode))
	}
	if len(b.TailCode) != len(b2.TailCode) {
		s = s + fmt.Sprintf("b Tailcode is %d; b2 TailCode is %d", len(b.TailCode), len(b2.TailCode))
	}
	if b.KernelBase != b2.KernelBase {
		s = s + fmt.Sprintf("b KernelBase is %#x; b2 KernelBase is %#x", b.KernelBase, b2.KernelBase)
	}
	if b.KernelOffset != b2.KernelOffset {
		s = s + fmt.Sprintf("b KernelOffset is %#x; b2 KernelOffset is %#x", b.KernelOffset, b2.KernelOffset)
	}
	return s
}

// String stringifies a LinuxHeader into comma-separated parts
func (h *LinuxHeader) String() string {
	return strings.Join(h.Show(), ",")
}

// InitRAMFS returns a []byte from KernelCode which can be used to save or replace
// an existing InitRAMFS. The fun part is that there are no symbols; what we do instead
// is find the programs what are RW and look for the cpio magic in them. If we find it,
// we see if it can be read as a cpio and, if so, if there is a /dev or /init inside.
// We repeat until we succeed or there's nothing left.
func (b *BzImage) InitRAMFS() (int, int, error) {
	f, err := b.ELF()
	if err != nil {
		return -1, -1, err
	}
	// Find the program header with RWE.
	var dat []byte
	var prog *elf.Prog
	for _, p := range f.Progs {
		if p.Flags&(elf.PF_X|elf.PF_W|elf.PF_R) == elf.PF_X|elf.PF_W|elf.PF_R {
			dat, err = ioutil.ReadAll(p.Open())
			if err != nil {
				return -1, -1, err
			}
			prog = p
			break
		}
	}
	if dat == nil {
		return -1, -1, fmt.Errorf("can't find an RWE prog in kernel")
	}

	archiver, err := cpio.Format("newc")
	if err != nil {
		return -1, -1, fmt.Errorf("format newc not supported: %v", err)
	}
	var cur int
	for cur < len(dat) {
		x := bytes.Index(dat, []byte("070701"))
		if x == -1 {
			return -1, -1, fmt.Errorf("no newc cpio magic found")
		}
		if err != nil {
			return -1, -1, err
		}
		cur = x
		var r = bytes.NewReader(dat[cur:])
		rr := archiver.Reader(r)
		Debug("r.Len is %v", r.Len())
		var found bool
		var size int
		for {
			rec, err := rr.ReadRecord()
			Debug("Check %v", rec)
			if err == io.EOF {
				break
			}
			if err != nil {
				Debug("error reading records: %v", err)
				break
			}
			switch rec.Name {
			case "init", "dev", "bin", "usr":
				Debug("Found initramfs at %d, %d bytes", cur, len(dat)-r.Len())
				found = true
			}
			size = int(rec.FilePos) + int(rec.FileSize)
		}
		Debug("Size is %d", size)
		// Add the trailer size.
		y := x + size
		if found {
			// The slice consists of the bytes for cur to the length of initramfs.
			// We can derive the initramfs length by knowing how much is left of the reader.
			Debug("Return %d %#x slice %d:%d from %d byte dat", len(dat[x:y]), len(dat[x:y]), cur, y, len(dat))
			x += int(prog.Off)
			y += int(prog.Off)
			// We need to round y up to the end of the record. We have to do this after we
			// add the prog.Off value to it.
			y = (y + 3) &^ 3
			// and add the size of the trailer record.
			y += 120
			y += 4 // and add at least one word of null
			y = (y + 3) &^ 3
			Debug("InitRAMFS: return %d, %d", x, y)
			return x, y, nil
		}
		cur += 6
	}
	return -1, -1, fmt.Errorf("no cpio found")
}
