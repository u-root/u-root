package nfs

import (
	"bytes"
	"context"
	"os"

	"github.com/go-git/go-billy/v5"
	"github.com/willscott/go-nfs-client/nfs/xdr"
)

// Backing billy.FS doesn't support hard links
func onLink(ctx context.Context, w *response, userHandle Handler) error {
	w.errorFmt = wccDataErrorFormatter
	obj := DirOpArg{}
	err := xdr.Read(w.req.Body, &obj)
	if err != nil {
		return &NFSStatusError{NFSStatusInval, err}
	}
	attrs, err := ReadSetFileAttributes(w.req.Body)
	if err != nil {
		return &NFSStatusError{NFSStatusInval, err}
	}

	target, err := xdr.ReadOpaque(w.req.Body)
	if err != nil {
		return &NFSStatusError{NFSStatusInval, err}
	}

	fs, path, err := userHandle.FromHandle(obj.Handle)
	if err != nil {
		return &NFSStatusError{NFSStatusStale, err}
	}
	if !billy.CapabilityCheck(fs, billy.WriteCapability) {
		return &NFSStatusError{NFSStatusROFS, os.ErrPermission}
	}

	if len(string(obj.Filename)) > PathNameMax {
		return &NFSStatusError{NFSStatusNameTooLong, os.ErrInvalid}
	}

	newFilePath := fs.Join(append(path, string(obj.Filename))...)
	if _, err := fs.Stat(newFilePath); err == nil {
		return &NFSStatusError{NFSStatusExist, os.ErrExist}
	}
	if s, err := fs.Stat(fs.Join(path...)); err != nil {
		return &NFSStatusError{NFSStatusAccess, err}
	} else if !s.IsDir() {
		return &NFSStatusError{NFSStatusNotDir, nil}
	}

	fp := userHandle.ToHandle(fs, append(path, string(obj.Filename)))
	changer := userHandle.Change(fs)
	if changer == nil {
		return &NFSStatusError{NFSStatusAccess, err}
	}
	cos, ok := changer.(UnixChange)
	if !ok {
		return &NFSStatusError{NFSStatusAccess, err}
	}

	err = cos.Link(string(target), newFilePath)
	if err != nil {
		return &NFSStatusError{NFSStatusAccess, err}
	}
	if err := attrs.Apply(changer, fs, newFilePath); err != nil {
		return &NFSStatusError{NFSStatusIO, err}
	}

	writer := bytes.NewBuffer([]byte{})
	if err := xdr.Write(writer, uint32(NFSStatusOk)); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}

	// "handle follows"
	if err := xdr.Write(writer, uint32(1)); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}
	if err := xdr.Write(writer, fp); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}
	if err := WritePostOpAttrs(writer, tryStat(fs, append(path, string(obj.Filename)))); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}

	if err := WriteWcc(writer, nil, tryStat(fs, path)); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}

	if err := w.Write(writer.Bytes()); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}
	return nil
}
