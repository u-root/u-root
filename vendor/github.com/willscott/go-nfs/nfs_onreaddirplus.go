package nfs

import (
	"bytes"
	"context"
	"path"

	"github.com/willscott/go-nfs-client/nfs/xdr"
)

type readDirPlusArgs struct {
	Handle      []byte
	Cookie      uint64
	CookieVerif uint64
	DirCount    uint32
	MaxCount    uint32
}

type readDirPlusEntity struct {
	FileID     uint64
	Name       []byte
	Cookie     uint64
	Attributes *FileAttribute `xdr:"optional"`
	Handle     *[]byte        `xdr:"optional"`
	Next       bool
}

func joinPath(parent []string, elements ...string) []string {
	joinedPath := make([]string, 0, len(parent)+len(elements))
	joinedPath = append(joinedPath, parent...)
	joinedPath = append(joinedPath, elements...)
	return joinedPath
}

func onReadDirPlus(ctx context.Context, w *response, userHandle Handler) error {
	w.errorFmt = opAttrErrorFormatter
	obj := readDirPlusArgs{}
	if err := xdr.Read(w.req.Body, &obj); err != nil {
		return &NFSStatusError{NFSStatusInval, err}
	}

	// in case of test, nfs-client send:
	// DirCount = 512
	// MaxCount = 4096
	if obj.DirCount < 512 || obj.MaxCount < 4096 {
		return &NFSStatusError{NFSStatusTooSmall, nil}
	}

	fs, p, err := userHandle.FromHandle(obj.Handle)
	if err != nil {
		return &NFSStatusError{NFSStatusStale, err}
	}

	contents, verifier, err := getDirListingWithVerifier(userHandle, obj.Handle, obj.CookieVerif)
	if err != nil {
		return err
	}
	if obj.Cookie > 0 && obj.CookieVerif > 0 && verifier != obj.CookieVerif {
		return &NFSStatusError{NFSStatusBadCookie, nil}
	}

	entities := make([]readDirPlusEntity, 0)
	dirBytes := uint32(0)
	maxBytes := uint32(100) // conservative overhead measure

	started := obj.Cookie == 0
	if started {
		// add '.' and '..' to entities
		dotdotFileID := uint64(0)
		if len(p) > 0 {
			dda := tryStat(fs, p[0:len(p)-1])
			if dda != nil {
				dotdotFileID = dda.Fileid
			}
		}
		dotFileID := uint64(0)
		da := tryStat(fs, p)
		if da != nil {
			dotFileID = da.Fileid
		}
		entities = append(entities,
			readDirPlusEntity{Name: []byte("."), Cookie: 0, Next: true, FileID: dotFileID, Attributes: da},
			readDirPlusEntity{Name: []byte(".."), Cookie: 1, Next: true, FileID: dotdotFileID},
		)
	}

	eof := true
	maxEntities := userHandle.HandleLimit() / 2
	fb := 0
	fss := 0
	for i, c := range contents {
		// cookie equates to index within contents + 2 (for '.' and '..')
		cookie := uint64(i + 2)
		fb++
		if started {
			fss++
			dirBytes += uint32(len(c.Name()) + 20)
			maxBytes += 512 // TODO: better estimation.
			if dirBytes > obj.DirCount || maxBytes > obj.MaxCount || len(entities) > maxEntities {
				eof = false
				break
			}

			filePath := joinPath(p, c.Name())
			handle := userHandle.ToHandle(fs, filePath)
			attrs := ToFileAttribute(c, path.Join(filePath...))
			entities = append(entities, readDirPlusEntity{
				FileID:     attrs.Fileid,
				Name:       []byte(c.Name()),
				Cookie:     cookie,
				Attributes: attrs,
				Handle:     &handle,
				Next:       true,
			})
		} else if cookie == obj.Cookie {
			started = true
		}
	}

	writer := bytes.NewBuffer([]byte{})
	if err := xdr.Write(writer, uint32(NFSStatusOk)); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}
	if err := WritePostOpAttrs(writer, tryStat(fs, p)); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}
	if err := xdr.Write(writer, verifier); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}

	if err := xdr.Write(writer, len(entities) > 0); err != nil { // next
		return &NFSStatusError{NFSStatusServerFault, err}
	}
	if len(entities) > 0 {
		entities[len(entities)-1].Next = false
		// no next for last entity

		for _, e := range entities {
			if err := xdr.Write(writer, e); err != nil {
				return &NFSStatusError{NFSStatusServerFault, err}
			}
		}
	}
	if err := xdr.Write(writer, eof); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}
	// TODO: track writer size at this point to validate maxcount estimation and stop early if needed.

	if err := w.Write(writer.Bytes()); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}
	return nil
}
