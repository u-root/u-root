package stconfig

import (
	"fmt"
	"os"

	"github.com/u-root/u-root/pkg/bootconfig"
)

// PackBootConfiguration packages a boot configuration containing different
// binaries and a manifest. The files to be included are taken from the
// path specified in the provided manifest.json
func PackBootConfiguration(createOutputFilename, createManifest string) error {
	if _, err := os.Stat(createManifest); os.IsNotExist(err) {
		return fmt.Errorf("manifest file does not exist: %v", err)
	}
	return bootconfig.ToZip(createOutputFilename, createManifest)
}
