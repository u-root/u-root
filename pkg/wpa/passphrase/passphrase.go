package passphrase

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

const ResultFormat = `network={
	ssid="%s"
	#psk="%s"
	psk=%s
}
`

func errorCheck(essid string, pass string) error {
	if len(pass) < 8 || len(pass) > 63 {
		return fmt.Errorf("Passphrase must be 8..63 characters\n")
	}
	if len(essid) == 0 {
		return fmt.Errorf("essid cannot be empty\n")
	}
	return nil
}

func Run(essid string, pass string) ([]byte, error) {
	err := errorCheck(essid, pass)
	if err != nil {
		return nil, err
	}

	pskBinary := pbkdf2.Key([]byte(pass), []byte(essid), 4096, 32, sha1.New)
	pskHexString := hex.EncodeToString(pskBinary)
	return []byte(fmt.Sprintf(ResultFormat, essid, pass, pskHexString)), nil
}
