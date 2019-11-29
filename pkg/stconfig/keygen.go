package stconfig

import (
	"fmt"
	"os"

	"github.com/u-root/u-root/pkg/crypto"
)

// GenKeys generates ED25519 keypair and stores it on the harddrive
func GenKeys(genkeysPrivateKeyFile, genkeysPublicKeyFile, genkeysPassphrase string) error {
	if _, err := os.Stat(genkeysPrivateKeyFile); os.IsNotExist(err) {
		return fmt.Errorf("private key file does not exist: %v", err)
	}
	if _, err := os.Stat(genkeysPublicKeyFile); os.IsNotExist(err) {
		return fmt.Errorf("public key file does not exist: %v", err)
	}
	return crypto.GeneratED25519Key([]byte(genkeysPassphrase), genkeysPrivateKeyFile, genkeysPublicKeyFile)
}
