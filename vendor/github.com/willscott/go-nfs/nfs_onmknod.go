package nfs

import (
	"bytes"
	"context"
	"os"

	"github.com/go-git/go-billy/v5"
	"github.com/willscott/go-nfs-client/nfs/xdr"
)

type nfs_ftype int32

const (
	FTYPE_NF3REG  nfs_ftype = 1
	FTYPE_NF3DIR  nfs_ftype = 2
	FTYPE_NF3BLK  nfs_ftype = 3
	FTYPE_NF3CHR  nfs_ftype = 4
	FTYPE_NF3LNK  nfs_ftype = 5
	FTYPE_NF3SOCK nfs_ftype = 6
	FTYPE_NF3FIFO nfs_ftype = 7
)

// Backing billy.FS doesn't support creation of
// char, block, socket, or fifo pipe nodes
func onMknod(ctx context.Context, w *response, userHandle Handler) error {
	w.errorFmt = wccDataErrorFormatter
	obj := DirOpArg{}
	err := xdr.Read(w.req.Body, &obj)
	if err != nil {
		return &NFSStatusError{NFSStatusInval, err}
	}

	ftype, err := xdr.ReadUint32(w.req.Body)
	if err != nil {
		return &NFSStatusError{NFSStatusInval, err}
	}

	// see if the filesystem supports mknod
	fs, path, err := userHandle.FromHandle(obj.Handle)
	if err != nil {
		return &NFSStatusError{NFSStatusStale, err}
	}
	if !billy.CapabilityCheck(fs, billy.WriteCapability) {
		return &NFSStatusError{NFSStatusROFS, os.ErrPermission}
	}
	c := userHandle.Change(fs)
	if c == nil {
		return &NFSStatusError{NFSStatusAccess, os.ErrPermission}
	}
	cu, ok := c.(UnixChange)
	if !ok {
		return &NFSStatusError{NFSStatusAccess, os.ErrPermission}
	}

	if len(string(obj.Filename)) > PathNameMax {
		return &NFSStatusError{NFSStatusNameTooLong, os.ErrInvalid}
	}

	newFilePath := fs.Join(append(path, string(obj.Filename))...)
	if _, err := fs.Stat(newFilePath); err == nil {
		return &NFSStatusError{NFSStatusExist, os.ErrExist}
	}
	parent, err := fs.Stat(fs.Join(path...))
	if err != nil {
		return &NFSStatusError{NFSStatusAccess, err}
	} else if !parent.IsDir() {
		return &NFSStatusError{NFSStatusNotDir, nil}
	}
	fp := userHandle.ToHandle(fs, append(path, string(obj.Filename)))

	switch nfs_ftype(ftype) {
	case FTYPE_NF3CHR:
	case FTYPE_NF3BLK:
		// read devicedata3 = {sattr3, specdata3}
		attrs, err := ReadSetFileAttributes(w.req.Body)
		if err != nil {
			return &NFSStatusError{NFSStatusInval, err}
		}
		specData1, err := xdr.ReadUint32(w.req.Body)
		if err != nil {
			return &NFSStatusError{NFSStatusInval, err}
		}
		specData2, err := xdr.ReadUint32(w.req.Body)
		if err != nil {
			return &NFSStatusError{NFSStatusInval, err}
		}

		err = cu.Mknod(newFilePath, uint32(attrs.Mode(parent.Mode())), specData1, specData2)
		if err != nil {
			return &NFSStatusError{NFSStatusAccess, err}
		}
		if err = attrs.Apply(cu, fs, newFilePath); err != nil {
			return &NFSStatusError{NFSStatusServerFault, err}
		}

	case FTYPE_NF3SOCK:
		// read sattr3
		attrs, err := ReadSetFileAttributes(w.req.Body)
		if err != nil {
			return &NFSStatusError{NFSStatusInval, err}
		}
		if err := cu.Socket(newFilePath); err != nil {
			return &NFSStatusError{NFSStatusAccess, err}
		}
		if err = attrs.Apply(cu, fs, newFilePath); err != nil {
			return &NFSStatusError{NFSStatusServerFault, err}
		}

	case FTYPE_NF3FIFO:
		// read sattr3
		attrs, err := ReadSetFileAttributes(w.req.Body)
		if err != nil {
			return &NFSStatusError{NFSStatusInval, err}
		}
		err = cu.Mkfifo(newFilePath, uint32(attrs.Mode(parent.Mode())))
		if err != nil {
			return &NFSStatusError{NFSStatusAccess, err}
		}
		if err = attrs.Apply(cu, fs, newFilePath); err != nil {
			return &NFSStatusError{NFSStatusServerFault, err}
		}

	default:
		return &NFSStatusError{NFSStatusBadType, os.ErrInvalid}
		// end of input.
	}

	writer := bytes.NewBuffer([]byte{})
	if err := xdr.Write(writer, uint32(NFSStatusOk)); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}

	// "handle follows"
	if err := xdr.Write(writer, uint32(1)); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}
	// fh3
	if err := xdr.Write(writer, fp); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}
	// attr
	if err := WritePostOpAttrs(writer, tryStat(fs, append(path, string(obj.Filename)))); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}
	// wcc
	if err := WriteWcc(writer, nil, tryStat(fs, path)); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}

	if err := w.Write(writer.Bytes()); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}

	return nil
}
