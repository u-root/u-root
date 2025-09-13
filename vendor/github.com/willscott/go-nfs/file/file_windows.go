//go:build windows

package file

import "os"

func getOSFileInfo(info os.FileInfo) *FileInfo {
	// https://godoc.org/golang.org/x/sys/windows#GetFileInformationByHandle
	// can be potentially used to populate Nlink

	return nil
}
