package main

import "fmt"

func (f *File) String() string {

	return fmt.Sprintf("Ino 0x%x Mode 0%o UID %d GID %d Nlink %d Mtime %d FileSize %d Major %d Minor %d NameSize %d Name %s",
		f.Ino,
		f.Mode,
		f.UID,
		f.GID,
		f.Nlink,
		f.Mtime,
		f.FileSize,
		f.Major,
		f.Minor,
		//f.Rmajor,
		//f.Rminor,
		f.NameSize,
		f.Name)
}
