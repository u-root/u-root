package stconfig

import (
	"fmt"

	"github.com/u-root/u-root/pkg/bootconfig"
)

// UnpackBootConfiguration unpacks a boot configuration file and returns the
// file path of a directory containing the data
func UnpackBootConfiguration(unpackInputFilename string) error {
	_, outputDir, err := bootconfig.FromZip(unpackInputFilename)
	if err != nil {
		return err
	}

	fmt.Println("Boot configuration unpacked into: " + outputDir)

	return nil
}
