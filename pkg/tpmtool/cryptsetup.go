package tpmtool

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"syscall"
)

const (
	// CryptsetupBinary name
	CryptsetupBinary = "cryptsetup"
	// DefaultFormatParams is a default cryptsetup secure option list
	DefaultFormatParams = "-c aes-xts-essiv:sha256 -s 512 -y --use-random -q"
	// DefaultKeyPath is the tmpfs directory for storing keys
	DefaultKeyPath = "/tmp/tpmtool"
	// TmpfsFsName is the linux tpmfs fs name
	TmpfsFsName = "tmpfs"
	// DefaultDevMapperPath is the standard Linux device mapper path
	DefaultDevMapperPath = "/dev/mapper/"
)

var (
	// TmpfsFsOptions are secure fs options
	TmpfsFsOptions string
)

func init() {
	processUser, err := user.Current()
	if err != nil {
		log.Fatalln("Couldn't set Tpmfs keystore options")
	}

	TmpfsFsOptions = path.Join("rw,size=1M,nr_inodes=5k,noexec,nodev,nosuid,uid=", processUser.Uid, ",gid=", processUser.Gid, ",mode=1700")
}

// MountKeystore mounts the tmpfs key store
func MountKeystore() (string, error) {
	flags := 0
	data := ""

	randBytes := make([]byte, 16)
	if _, err := rand.Read(randBytes); err != nil {
		return "", err
	}

	target := filepath.Join(DefaultKeyPath, hex.EncodeToString(randBytes))
	if err := os.MkdirAll(target, 0750); err != nil {
		return "", err
	}

	if err := syscall.Mount(TmpfsFsName, target, TmpfsFsName, uintptr(flags), data); err != nil {
		return "", err
	}

	return target, nil
}

// UnmountKeystore unmounts the tpmfs key store
func UnmountKeystore(target string) error {
	syscall.Sync()
	return syscall.Unmount(target, syscall.MNT_FORCE|syscall.MNT_DETACH)
}

// CryptsetupFormat formats a device with LUKS
func CryptsetupFormat(keyPath string, devicePath string) error {
	cryptsetup, err := exec.LookPath(CryptsetupBinary)
	if err != nil {
		return err
	}

	cmd := exec.Command(cryptsetup, DefaultFormatParams, "-d", keyPath, "luksFormat", devicePath)

	return cmd.Run()
}

// CryptsetupOpen opens a LUKS device
func CryptsetupOpen(keyPath string, devicePath string) (string, error) {
	cryptsetup, err := exec.LookPath(CryptsetupBinary)
	if err != nil {
		return "", err
	}

	randBytes := make([]byte, 16)
	if _, err := rand.Read(randBytes); err != nil {
		return "", err
	}

	deviceName := hex.EncodeToString(randBytes)
	cmd := exec.Command(cryptsetup, "-d", keyPath, "luksOpen", devicePath, deviceName)

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return deviceName, nil
}

// CryptsetupClose closes a LUKS device
func CryptsetupClose(deviceName string) error {
	cryptsetup, err := exec.LookPath(CryptsetupBinary)
	if err != nil {
		return err
	}

	devicePath := path.Join(DefaultDevMapperPath, deviceName)
	cmd := exec.Command(cryptsetup, "luksClose", devicePath)

	return cmd.Run()
}
