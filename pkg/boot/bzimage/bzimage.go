// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bzimage implements decoding for bzImage files.
//
// The bzImage struct contains all the information about the file and can
// be used to create a new bzImage.
package bzimage

// xz --check=crc32 $BCJ --lzma2=$LZMA2OPTS,dict=32MiB

import (
	"bytes"
	"compress/gzip"
	"debug/elf"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os/exec"
	"reflect"
	"strings"
	"unsafe"
)

const minBootParamLen = 616

// A decompressor is a function which reads compressed bytes via the io.Reader and
// writes the uncompressed bytes to the io.Writer.
type decompressor func(w io.Writer, r io.Reader) error

type magic struct {
	name          string
	signature     []byte
	decompressors []decompressor
}

// MSDOS tag used in .efi binaries.
// There are no words.
const MSDOS = "MZ"

var (
	// TODO(10000TB): remove dependency on cmds / programs.
	//
	// These are the magics, along with the command to run
	// it as a pipe. They need be the actual command than a
	// shell script, which won't work in u-root.
	magics = []*magic{
		// GZIP
		{"gunzip", []byte{0x1F, 0x8B}, []decompressor{gunzip}},
		// XZ
		// It would be nice to use a Go package instead of shelling out to 'unxz'.
		// https://github.com/ulikunitz/xz fails to decompress the payloads and returns an error: "unsupported filter count"
		{"unxz", []byte{0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00}, []decompressor{stripSize(unxz), stripSize(execer("unxz"))}},
		// LZMA
		{"unlzma", []byte{0x5D, 0x00, 0x00}, []decompressor{stripSize(unlzma)}},
		// LZO
		{"lzop", []byte{0x89, 0x4C, 0x5A, 0x4F, 0x00, 0x0D, 0x0A, 0x1A, 0x0A}, []decompressor{stripSize(execer("lzop", "-c", "-d"))}},
		// ZSTD
		{"unzstd", []byte{0x28, 0xB5, 0x2F, 0xFD}, []decompressor{stripSize(unzstd)}},
		// BZIP2
		{"unbzip2", []byte{0x42, 0x5A, 0x68}, []decompressor{stripSize(unbzip2)}},
		// LZ4 - Note that there are *two* file formats for LZ4 (http://fileformats.archiveteam.org/wiki/LZ4).
		// The Linux boot process uses the legacy 02 21 4C 18 magic bytes, while newer systems
		// use 04 22 4D 18
		{"unlz4", []byte{0x02, 0x21, 0x4C, 0x18}, []decompressor{stripSize(unlz4)}},
	}

	ErrNoMagic = errors.New("magic is not known")

	// Debug is a function used to log debug information. It
	// can be set to, for example, log.Printf.
	Debug = func(string, ...any) {}
)

// This is all C to Go and it reads like it, sorry
// unpacking bzimage is a mess, so for now, this is a mess.

// decompressor finds a decompressor by scanning a []byte for a tag.
func findDecompressors(b []byte) (*magic, error) {
	for _, m := range magics {
		if bytes.Index(b, m.signature) == 0 {
			return m, nil
		}
	}
	return nil, fmt.Errorf("%#x:%w", b[:16], ErrNoMagic)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// For now, it hardwires the KernelBase to 0x100000.
// bzImages were created by a process of evilution, and they are wondrous to behold.
// "Documentation" can be found at https://www.kernel.org/doc/html/latest/x86/boot.html.
// bzImages are almost impossible to modify. They form a sandwich with
// the compressed kernel code in the middle. It's actually a BLT:
// MBR and bootparams first 512 bytes
// the MBR includes 0xc0 bytes of boot code which is used for UEFI booting.
// Then there is "preamble" code which is the kernel decompressor; then the
// xz compressed kernel; then a library of sorts after the kernel which is called
// by the early uncompressed kernel code. This is all linked together and forms
// an essentially indivisible whole -- which we wish to divisible.
// That said, if you keep layout unchanged, you can modify the uncompressed kernel.
// For example, when you first build a kernel, you can:
// dd if=/dev/urandom of=x bs=1048576 count=8
// echo x | cpio -o > x.cpio
// and use that as an initrd, it's more or less an 8 MiB block you can replace
// as needed. Just make sure nothing grows. And make sure the initramfs is in
// the same place. Ah, joy.
//
// Important note for signed kernel images: The kernel signature is stripped away
// and ignored. Users of UnmarshalBinary must separately check the image signature,
// if required.
func (b *BzImage) UnmarshalBinary(d []byte) error {
	Debug("Processing %d byte image", len(d))

	stripped, err := stripSignature(d)
	if err != nil {
		return fmt.Errorf("error stripping kernel signature: %w", err)
	}
	d = stripped

	r := bytes.NewBuffer(d)
	if err := binary.Read(r, binary.LittleEndian, &b.Header); err != nil {
		return err
	}
	Debug("Header was %d bytes", len(d)-r.Len())
	Debug("magic %x switch %v", b.Header.HeaderMagic, b.Header.RealModeSwitch)
	if b.Header.HeaderMagic != HeaderMagic {
		return fmt.Errorf("not a bzImage: magic should be %02x, and is %02x", HeaderMagic, b.Header.HeaderMagic)
	}
	if b.Header.Protocolversion < 0x0208 {
		return fmt.Errorf("boot protocol version 0x%04x not supported, version 0x0208 or higher (Kernel 2.6.26) required", b.Header.Protocolversion)
	}
	Debug("RamDisk image %x size %x", b.Header.RamdiskImage, b.Header.RamdiskSize)
	Debug("StartSys %x", b.Header.StartSys)
	Debug("Boot type: %s(%x)", LoaderType[boottype(b.Header.TypeOfLoader)], b.Header.TypeOfLoader)

	if b.Header.SetupSects == 0 {
		// Per https://www.kernel.org/doc/html/latest/x86/boot.html?highlight=boot:
		// "For backwards compatibility, if the setup_sects field contains 0, the real value is 4."
		b.Header.SetupSects = 4
	}

	Debug("SetupSects %d", b.Header.SetupSects)

	off := len(d) - r.Len()
	// Per https://www.kernel.org/doc/html/v5.4/x86/boot.html#loading-the-rest-of-the-kernel:
	// "the 32-bit (non-real-mode) kernel starts at offset (setup_sects+1)*512 in the kernel file"
	// The +1 is because the MBR (1 sect) is always assumed. The logic calculating this
	// can be confirmed at: https://github.com/torvalds/linux/blob/master/arch/x86/boot/tools/build.c#L440
	b.KernelOffset = (uintptr(b.Header.SetupSects) + 1) * 512
	bclen := int(b.KernelOffset) - off
	Debug("Kernel offset is %d bytes, low1mcode is %d bytes", b.KernelOffset, bclen)
	b.BootCode = make([]byte, bclen)
	if _, err := r.Read(b.BootCode); err != nil {
		return err
	}
	Debug("%d bytes of BootCode", len(b.BootCode))

	b.HeadCode = make([]byte, b.Header.PayloadOffset)
	if _, err := r.Read(b.HeadCode); err != nil {
		return fmt.Errorf("can't read HeadCode: %w", err)
	}
	b.compressed = make([]byte, b.Header.PayloadSize)
	if _, err := r.Read(b.compressed); err != nil {
		return fmt.Errorf("can't read KernelCode: %w", err)
	}
	m, err := findDecompressors(b.compressed)
	if err != nil {
		return err
	}
	if b.NoDecompress {
		Debug("skipping code decompress")
	} else {
		Debug("Uncompress %d bytes", len(b.compressed))

		// The Linux boot process expects that the last 4 bytes of the compressed payload will
		// contain the size of the uncompressed payload. This works well for gzip, where the
		// last 4 bytes of the compressed payload contain the uncompressed size. However other
		// compression formats (bzip2, lzma, xz, lzo, lz4, zstd, etc) do not satisfy this
		// requirement, so the Makefile tacks on an extra 4 bytes for these compression formats
		// and expects that the decompression code will ignore them.
		// The authoritative list of compression formats that have the 4 byte size appended
		// can be found here: https://github.com/torvalds/linux/blob/master/arch/x86/boot/compressed/Makefile#L132-L145
		// (look for the entries ending in "_with_size", examples: bzip2_with_size, lzma_with_size.

		// Read the uncompressed length of the payload from the last 4 bytes of the payload.
		var uncompressedLength uint32
		last4Bytes := b.compressed[(len(b.compressed) - 4):]
		if err := binary.Read(bytes.NewBuffer(last4Bytes), binary.LittleEndian, &uncompressedLength); err != nil {
			return fmt.Errorf("error reading uncompressed kernel size: %w", err)
		}
		Debug("Original length of uncompressed kernel is: %d", uncompressedLength)

		// Use the decompressor and write the decompressed payload into b.KernelCode.
		var buf bytes.Buffer
		success := false
		var err error
		for _, decompressor := range m.decompressors {
			e := decompressor(&buf, bytes.NewBuffer(b.compressed))
			if e == nil {
				success = true
				b.KernelCode = buf.Bytes()
				break
			}
			err = errors.Join(err, fmt.Errorf("%s:%w", m.name, e))
		}
		if !success {
			return fmt.Errorf("error decompressing payload: %w", err)
		}

		// Verify that the length of the uncompressed payload matches the size read from the last 4 bytes of the compressed payload.
		if uint32(len(b.KernelCode)) != uncompressedLength {
			return fmt.Errorf("decompression failed, got size=%d bytes, expected size=%d bytes", len(b.KernelCode), uncompressedLength)
		}

		// Verify that the uncompressed payload is an ELF.
		elfMagic := []byte{0x7F, 0x45, 0x4C, 0x46}
		if bytes.Index(b.KernelCode, elfMagic) != 0 {
			return fmt.Errorf("decompressed payload must be an ELF with magic 0x%08x, found 0x%08x", elfMagic, b.KernelCode[0:4])
		}

		Debug("Kernel at %d, %d bytes", b.KernelOffset, len(b.KernelCode))
		Debug("KernelCode size: %d", len(b.KernelCode))
	}

	if err := binary.Read(r, binary.LittleEndian, &b.CRC32); err != nil {
		return fmt.Errorf("error reading CRC: %w", err)
	}
	Debug("CRC read from image is: 0x%08x", b.CRC32)

	b.TailCode = make([]byte, r.Len()) // Read all remaining bytes.
	if _, err := r.Read(b.TailCode); err != nil {
		return fmt.Errorf("can't read TailCode: %w", err)
	}

	b.KernelBase = uintptr(0x100000)
	if b.Header.RamdiskImage == 0 {
		return nil
	}
	if r.Len() != 0 {
		return fmt.Errorf("%d bytes left over", r.Len())
	}
	return nil
}

// stripSignature returns an image with the UEFI/PE signatures stripped.
//
// The linux kernel supports UEFI Stub booting, which allows the UEFI firmware to load the kernel as
// an executable. All UEFI images contain a PE/COFF header that defines the format of the executable
// code. The PE format is documented at: https://learn.microsoft.com/en-us/windows/win32/debug/pe-format.
//
// Signed kernels are problematic because the kernel signature process updates the boot code in the
// image, which in turn makes the CRC checksum of the image invalid. Specifically the `sbsigntools`
// package [1] used by Debian (and others) updates the "Certificate Table" information [2] and PE checksum.
// [1] https://github.com/phrack/sbsigntools
// [2] https://learn.microsoft.com/en-us/windows/win32/debug/pe-format#optional-header-data-directories-image-only
func stripSignature(image []byte) ([]byte, error) {
	// Clone the slice so that we do not modify the slice that is passed to this function.
	d := make([]byte, len(image))
	copy(d, image)

	dosMagic := []byte("MZ")
	peMagic := []byte("PE\x00\x00")
	peSignaturePtr := 0x3C

	// Verify that the image has a MS DOS Stub.
	if bytes.Index(d, dosMagic) != 0 {
		return d, nil
	}

	// Locate the PE signature.
	// The PE signature is located at the offset found in location 0x3C.
	if peSignaturePtr+4 > len(d) {
		// Image is not large enough to have a PE signature offset.
		return d, nil
	}
	peMagicOffset := uintptr(binary.LittleEndian.Uint32(d[peSignaturePtr:]))
	if peMagicOffset+uintptr(len(peMagic)) > uintptr(len(d)) {
		// Image is not large enough to have a PE signature.
		return d, nil
	}

	peImage := &PEImage{}
	if peMagicOffset+unsafe.Sizeof(peImage) > uintptr(len(d)) {
		// File is too small to have the PE headers.
		return d, nil
	}
	if err := binary.Read(bytes.NewReader(d[peMagicOffset:]), binary.LittleEndian, peImage); err != nil {
		return nil, fmt.Errorf("failed to read PE header: %w", err)
	}
	// Verify that the image has the PE magic number.
	if !bytes.Equal(peImage.PEMagic[:], peMagic) {
		return d, nil
	}

	Debug("Found a PE image")

	// TODO: Consider performing PE checksum and signature verification.
	// This is non trivial because we must decide what roots to trust, etc.
	// Existing code at https://github.com/saferwall/pe might be helpful in this process.

	optionalHeaderOffset := peMagicOffset + unsafe.Offsetof(peImage.OptionalHeader)
	Debug("Optional header offset: 0x%x", optionalHeaderOffset)

	// Zero out the PE Checksum.
	checksumOffset := uintptr(64)
	if checksumOffset+4 < uintptr(peImage.COFFHeader.SizeOfOptionalHeader) {
		Debug("Clearing checksum")
		binary.LittleEndian.PutUint32(d[optionalHeaderOffset+checksumOffset:], 0)
	}

	// Zero out the Certificate Table.
	var certificateTableOffset uintptr
	// Unfortunately the offset of the Certificate Table depends on whether the image is
	// PE32 or P32+ (https://learn.microsoft.com/en-us/windows/win32/debug/pe-format#optional-header-data-directories-image-only)
	switch peImage.OptionalHeader.Magic {
	case 0x10B: // PE32
		Debug("Found PE32 image")
		certificateTableOffset = 128
	case 0x20B: // PES32+
		Debug("Found PE32+ image")
		certificateTableOffset = 144
	default:
		return nil, fmt.Errorf("unknown Magic type: 0x%x", peImage.OptionalHeader.Magic)
	}
	if certificateTableOffset+8 < uintptr(peImage.COFFHeader.SizeOfOptionalHeader) {
		certificateTableAddress := optionalHeaderOffset + certificateTableOffset
		if binary.LittleEndian.Uint64(d[certificateTableAddress:]) > 0 {
			log.Printf("WARNING! The image is signed but the signature is being ignored.")
		}

		Debug("Clearing Certificate Table")
		binary.LittleEndian.PutUint64(d[certificateTableAddress:], 0)
	}

	return d, nil
}

// ErrKCodeMissing is returned if kernel code was not decompressed.
var ErrKCodeMissing = errors.New("no kernel code was decompressed")

// MarshalBinary implements the encoding.BinaryMarshaler interface.
// The marshal'd image is *not* signed.
func (b *BzImage) MarshalBinary() ([]byte, error) {
	if b.NoDecompress || b.KernelCode == nil {
		return nil, ErrKCodeMissing
	}
	// First step, make sure we can compress the kernel.
	dat, err := compress(b.KernelCode, "--lzma2=,dict=32MiB")
	if err != nil {
		return nil, err
	}
	if len(dat) > len(b.compressed) {
		return nil, fmt.Errorf("marshal: compressed KernelCode too big: was %d, now %d", len(b.compressed), len(dat))
	}
	Debug("b.compressed len %#x dat len %#x pad it out", len(b.compressed), len(dat))

	if len(dat) < len(b.compressed) {
		// If the new compressed payload fits in the existing compressed payload space then we
		// can fit the new payload in by putting it at the *end* of the original payload space
		// and updating `PayloadOffset` and `PayloadSize`. This is safer than placing the new
		// image at the start and padding with tailing NULLs because there's no guarantee about
		// how different decompressors will handle the trailing NULLs.

		diff := len(b.compressed) - len(dat)

		// Create the new payload with the length of the original payload and copy the new
		// payload to the end.
		newPayload := make([]byte, len(b.compressed))
		copy(newPayload[diff:], dat)

		// Update the headers with the new payload offset and size.
		b.Header.PayloadOffset += uint32(diff)
		b.Header.PayloadSize -= uint32(diff)

		// Swap in the new payload.
		dat = newPayload
	}

	b.compressed = dat

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
	if _, err := w.Write(b.compressed); err != nil {
		return nil, err
	}
	// b.TailCode is not written to the marshalled image. TailCode is used by signed images
	// and therefore likely to break because this code does not produce signed images.
	totalSize := (b.KernelOffset + uintptr(b.Header.Syssize)*16) - unsafe.Sizeof(b.CRC32)
	padding := int(totalSize) - w.Len()
	if padding > 0 {
		if _, err := w.Write(bytes.Repeat([]byte{0}, padding)); err != nil {
			return nil, fmt.Errorf("error writing padding")
		}
	}

	Debug("Wrote %d bytes of header", w.Len())
	generatedCRC := crc32.ChecksumIEEE(w.Bytes()) ^ (0xffffffff)
	if err := binary.Write(&w, binary.LittleEndian, generatedCRC); err != nil {
		return nil, err
	}
	Debug("Finished writing, len is now %d bytes", w.Len())

	return w.Bytes(), nil
}

// compress compresses a []byte via xz using the dictOps, collecting it from stdout
func compress(b []byte, dictOps string) ([]byte, error) {
	Debug("b is %d bytes", len(b))
	// TODO: Replace this use of `exec` with a proper Go package.
	c := exec.Command("xz", "--check=crc32", "--x86", dictOps, "--stdout")
	stdout, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}
	c.Stdin = bytes.NewBuffer(b)
	if err := c.Start(); err != nil {
		return nil, err
	}

	dat, err := io.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	if err := c.Wait(); err != nil {
		return nil, err
	}
	Debug("Compressed data is %d bytes, starts with %#02x", len(dat), dat[:32])
	Debug("Last 16 bytes: %#02x", dat[len(dat)-16:])

	// Append the original, uncompressed size of the payload.
	// HEAR YE, HEAR YE: The uncompressed size of the payload is appended to the payload because
	// the Linux boot process expects that the last 4 bytes of teh payload will contain the
	// uncompressed size. This appending is only required if the compression format does not
	// already satisfy this requirement. If this function is changed to use GZIP compression in
	// the future then this code is not required. This code is required for compression formats
	// such as bzip lzma xz lzo lz4 and zstd. See https://github.com/torvalds/linux/blob/master/arch/x86/boot/compressed/Makefile#L132-L145
	// for an authoritative list of which file formats require the extra 4 bytes appended (look for
	// "_with_size").
	buf := bytes.NewBuffer(dat)
	if binary.Write(buf, binary.LittleEndian, uint32(len(b))); err != nil {
		return nil, fmt.Errorf("failed to append the uncompressed size: %w", err)
	}
	return buf.Bytes(), nil
}

// ELF extracts the KernelCode.
func (b *BzImage) ELF() (*elf.File, error) {
	Debug("getting ELF...")
	if b.NoDecompress || b.KernelCode == nil {
		return nil, ErrKCodeMissing
	}
	Debug("creating a elf NewFile...")
	e, err := elf.NewFile(bytes.NewReader(b.KernelCode))
	if err != nil {
		return nil, err
	}
	return e, nil
}

// Equal compares two kernels and returns true if they are equal.
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

// MarshalBinary implements encoding.BinaryMarshaler
func (h *LinuxHeader) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, h)
	return buf.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryMarshaler
func (h *LinuxHeader) UnmarshalBinary(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, h)
}

// MarshalBinary implements encoding.BinaryMarshaler
func (h *LinuxParams) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, h)
	return buf.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryMarshaler
func (h *LinuxParams) UnmarshalBinary(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, h)
}

// Show stringifies a LinuxHeader into a []string.
func (h *LinuxHeader) Show() []string {
	var s []string

	val := reflect.ValueOf(*h)
	for i := 0; i < val.NumField(); i++ {
		v := val.Field(i)
		k := reflect.ValueOf(v).Kind()
		n := val.Type().Field(i).Name
		switch k {
		case reflect.Bool:
			s = append(s, fmt.Sprintf("%s:%v", n, v))
		default:
			s = append(s, fmt.Sprintf("%s:%#02x", n, v))
		}
	}
	return s
}

// Show stringifies a LinuxParams into a []string.
func (h *LinuxParams) Show() []string {
	var s []string

	val := reflect.ValueOf(*h)
	for i := 0; i < val.NumField(); i++ {
		v := val.Field(i)
		k := reflect.ValueOf(v).Kind()
		n := val.Type().Field(i).Name
		switch k {
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
	if b.CRC32 != b2.CRC32 {
		s = s + fmt.Sprintf("b CRC32 is 0x%08x; b2 CRC32 is 0x%08x", b.CRC32, b2.CRC32)
	}
	if b.KernelBase != b2.KernelBase {
		// NOTE: this is hardcoded to 0x100000
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

// String stringifies a LinuxParams into comma-separated parts
func (h *LinuxParams) String() string {
	return strings.Join(h.Show(), ",")
}

// ErrCfgNotFound is returned if embedded config is not found.
var ErrCfgNotFound = errors.New("embedded config not found")

// ReadConfig extracts embedded config from kernel
func (b *BzImage) ReadConfig() (string, error) {
	i := bytes.Index(b.KernelCode, []byte("IKCFG_ST\037\213\010"))
	if i == -1 {
		return "", ErrCfgNotFound
	}
	i += 8
	mb := 1024 * 1024 // read only 1 mb; arbitrary
	buf := bytes.NewReader(b.KernelCode[i : i+mb])
	gz, err := gzip.NewReader(buf)
	if err != nil {
		return "", err
	}
	// make it stop at end of stream, since we don't know the actual size
	gz.Multistream(false)
	cfg, err := io.ReadAll(gz)
	if err != nil {
		return "", err
	}
	return string(cfg), nil
}
