package boot

import (
	"fmt"
	"log"

	"github.com/u-root/u-root/pkg/cpio"
)

// multibootImage is a multiboot-formated OSImage.
type multibootImage struct{}

var _ OSImage = &multibootImage{}

func newMultibootImage(a *cpio.Archive) (OSImage, error) {
	return nil, fmt.Errorf("multiboot images unimplemented")
}

// ExecutionInfo implements OSImage.ExecutionInfo.
func (multibootImage) ExecutionInfo(log *log.Logger) {
	log.Printf("Multiboot images are unsupported")
}

// Execute implements OSImage.Execute.
func (multibootImage) Execute() error {
	return fmt.Errorf("multiboot images unimplemented")
}

// Pack implements OSImage.Pack.
func (multibootImage) Pack(sw cpio.RecordWriter) error {
	return fmt.Errorf("multiboot images unimplemented")
}
