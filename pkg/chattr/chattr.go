// Package chattr implements the chattr unix command.
package chattr

import (
	"fmt"
	"log"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	// Attribute flags
	FS_IMMUTABLE_FL = 0x00000010 // Immutable file
	FS_APPEND_FL    = 0x00000020 // Append-only file
	// ... other flags as needed
)

// SetAttr sets the attributes of a file.
func SetAttr(file *os.File, attrStr string) error {
	var attr int32
	switch attrStr {
	case "+i":
		attr = FS_IMMUTABLE_FL
	case "+a":
		attr = FS_APPEND_FL
	case "-i":
		attr = -FS_IMMUTABLE_FL
	case "-a":
		attr = -FS_APPEND_FL
	default:
		return fmt.Errorf("Invalid attribute. Use +i, +a, -i, or -a")
	}

	currentAttr, err := GetAttr(file)
	log.Printf("currentAttr: %v %v", currentAttr, err)
	if err != nil {
		return err
	}

	if attr > 0 {
		currentAttr |= attr // Add attribute
	} else {
		currentAttr &= ^(-attr) // Remove attribute
	}
	ptr := (uintptr)(unsafe.Pointer(&currentAttr))
	return unix.IoctlSetInt(int(file.Fd()), unix.FS_IOC_SETFLAGS, int(ptr))
}

// GetAttr gets the attributes of a file.
func GetAttr(file *os.File) (int32, error) {
	attr, err := unix.IoctlGetInt(int(file.Fd()), unix.FS_IOC_GETFLAGS)
	return int32(attr), err
}
