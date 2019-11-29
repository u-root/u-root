package stconfig

import (
	"fmt"
	"os"

	"github.com/u-root/u-root/pkg/bootconfig"
)

// AddSignatureToBootConfiguration TODO:
func AddSignatureToBootConfiguration(signInputBootfile, signPrivKeyFile, signCertFile string) error {
	if _, err := os.Stat(signInputBootfile); os.IsNotExist(err) {
		return fmt.Errorf("boot config file does not exist: %v", err)
	}
	if _, err := os.Stat(signPrivKeyFile); os.IsNotExist(err) {
		return fmt.Errorf("private key file does not exist: %v", err)
	}
	if _, err := os.Stat(signCertFile); os.IsNotExist(err) {
		return fmt.Errorf("certifivate file does not exist: %v", err)
	}
	return bootconfig.AddSignature(signInputBootfile, signPrivKeyFile, signCertFile)
}
