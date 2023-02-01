package boot

import (
	"fmt"
	"io"
	"log"
	"os"
	"unsafe"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/kexec"
	"github.com/u-root/u-root/pkg/pe"
	"github.com/u-root/u-root/pkg/uio"
)

const IMAGE_SCN_MEM_DISCARDABLE = 0x02000000

// PEImage is a PE-formated OSImage.
type PEImage struct {
	Kernel io.ReaderAt
}

var _ OSImage = &PEImage{}

func PEImageFromArchive(a *cpio.Archive) (OSImage, error) {
	return nil, fmt.Errorf("PE images unimplemented")
}

func PEImageFromFile(kernel *os.File) (OSImage, error) {
	k, err := uio.InMemFile(kernel)
	if err != nil {
		return nil, err
	}
	return &PEImage{
		Kernel: k,
	}, nil
}

// ExecutionInfo implements OSImage.ExecutionInfo.
func (PEImage) ExecutionInfo(log *log.Logger) {
	log.Printf("PE images are unsupported")
}

const M16 = 0x1000000
const KEXEC_RUN_PE = 0x00000004

// Execute implements OSImage.Execute.
func (p *PEImage) Execute() error {
	f, err := pe.NewFile(p.Kernel)
	if err != nil {
		return err
	}
	kernelBuf, err := uio.ReadAll(p.Kernel)
	if err != nil {
		return err
	}

	var segment []kexec.Segment

	// Windows loader expects to also fing the header of the EFI file.
	// We add everything before the first section as the first segment.
	var section_0 = f.Sections[0]
	s := kexec.Segment{
		Buf: kexec.Range{
			Start: uintptr(unsafe.Pointer(&kernelBuf[0])),
			Size:  uint(section_0.Offset),
		},
		Phys: kexec.Range{
			Start: M16,
			Size:  uint(uint64(section_0.VirtualAddress)),
		},
	}
	log.Printf("virt: %#x + %#x | phys: %#x + %#x", s.Buf.Start, s.Buf.Size, s.Phys.Start, s.Phys.Size)
	segment = append(segment, s)

	// Now add the actuall sections
	for _, section := range f.Sections {
		s := kexec.Segment{
			Buf: kexec.Range{
				Start: uintptr(unsafe.Pointer(&kernelBuf[section.Offset])),
				Size:  uint(section.Size),
			},
			Phys: kexec.Range{
				Start: M16 + uintptr(section.VirtualAddress),
				Size:  uint(section.VirtualSize),
			},
		}
		log.Printf("virt: %#x + %#x | phys: %#x + %#x", s.Buf.Start, s.Buf.Size, s.Phys.Start, s.Phys.Size)
		segment = append(segment, s)
	}

	var entry uintptr
	switch oh := f.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		entry = uintptr(oh.AddressOfEntryPoint)
	case *pe.OptionalHeader64:
		entry = uintptr(oh.AddressOfEntryPoint)
	}

	if err := kexec.Load(M16+entry, segment, KEXEC_RUN_PE); err != nil {
		return err
	}
	// kexec.Reboot()
	return nil
}

// Pack implements OSImage.Pack.
func (PEImage) Pack(sw cpio.RecordWriter) error {
	return fmt.Errorf("PE images unimplemented")
}
