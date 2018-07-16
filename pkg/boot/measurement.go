// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	_ "crypto/sha512"
	//"encoding/binary"
	"fmt"
	"io"

	"github.com/google/go-tpm/tpm"
	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/uio"
	"golang.org/x/sys/unix"
)

// MeasuringReader is a cpio.Reader that collects the signed data and compares
// it against the signature in the given cpio archive.
type MeasuringReader struct {
	r cpio.RecordReader

	signed    *bytes.Buffer
	signature *bytes.Buffer
}

// NewMeasuringReader returns a new measuring reader.
func NewMeasuringReader(r cpio.RecordReader) *MeasuringReader {
	return &MeasuringReader{
		r:         r,
		signed:    &bytes.Buffer{},
		signature: &bytes.Buffer{},
	}
}

// Verify verifies the contents of the archive as read so far.
//
// NOTE(UGH): Go crypto stuff is totally incompatible. ecdsa.PrivateKey.Sign
// does not output shit that is compatible with ecdsa.Verify -- COME ON. Only
// ecdsa.Sign does.
func (mr *MeasuringReader) Verify(pk *rsa.PublicKey) error {
	hashed := sha256.Sum256(mr.signed.Bytes())
	return rsa.VerifyPKCS1v15(pk, crypto.SHA256, hashed[:], mr.signature.Bytes())
}

// ExtendTPM extends the given tpm at pcrIndex with the content of the package.
func (mr *MeasuringReader) ExtendTPM(tpmRW io.ReadWriter, pcrIndex uint32) error {
	pcrValue := sha1.Sum(mr.signed.Bytes())
	_, err := tpm.PcrExtend(tpmRW, pcrIndex, pcrValue)
	return err
}

// ReadRecord wraps cpio.Reader.ReadRecord and adds the content to `signed` as
// necessary.
func (mr *MeasuringReader) ReadRecord() (cpio.Record, error) {
	for {
		rec, err := mr.r.ReadRecord()
		if err != nil {
			return rec, err
		}

		switch rec.Name {
		case "signature":
			_, err = mr.signature.ReadFrom(uio.Reader(rec))
			continue

		case "signature_algo":
			//err = binary.Read(uio.Reader(rec), binary.LittleEndian, &mr.algo)
			continue

		default:
			// Measure all regular files.
			if rec.Info.Mode&unix.S_IFMT == unix.S_IFREG {
				if _, err := mr.signed.WriteString(rec.Name); err != nil {
					return cpio.Record{}, err
				}
				if _, err := mr.signed.ReadFrom(uio.Reader(rec)); err != nil {
					return cpio.Record{}, err
				}
			}
			return rec, nil
		}
	}
}

// SigningWriter is a cpio.RecordWriter that collects digests as it writes
// files to the cpio archive.
type SigningWriter struct {
	w cpio.RecordWriter

	digest *bytes.Buffer
}

// NewSigningWriter returns a new signing cpio writer.
func NewSigningWriter(w cpio.RecordWriter) *SigningWriter {
	return &SigningWriter{
		w:      w,
		digest: &bytes.Buffer{},
	}
}

// WriteRecord implements cpio.RecordWriter.
func (sw *SigningWriter) WriteRecord(rec cpio.Record) error {
	rec = cpio.MakeReproducible(rec)
	if rec.Info.Name == "signature" || rec.Info.Name == "signature_algo" {
		return fmt.Errorf("cannot write signature or signature_algo files")
	}
	if rec.Info.Mode&unix.S_IFMT == unix.S_IFREG {
		if _, err := sw.digest.WriteString(rec.Info.Name); err != nil {
			return err
		}
		if _, err := sw.digest.ReadFrom(uio.Reader(rec)); err != nil {
			return err
		}
	}
	return sw.w.WriteRecord(rec)
}

// SHA1Sum returns the SHA1 sum of the collected digest.
func (sw *SigningWriter) SHA1Sum() [sha1.Size]byte {
	return sha1.Sum(sw.digest.Bytes())
}

// WriteSignature writes the signature and signature_algo files based on the
// collected digest.
//
// TODO(hugelgupf): stop hard-coding the private key and algorithm. Use
// crypto.Signer so TPM could be used to sign this if so desired.
func (sw *SigningWriter) WriteSignature(signer *rsa.PrivateKey) error {
	hashed := sha256.Sum256(sw.digest.Bytes())
	signature, err := signer.Sign(rand.Reader, hashed[:], crypto.SHA256)
	if err != nil {
		return err
	}
	if err := sw.w.WriteRecord(cpio.StaticFile("signature", string(signature), 0700)); err != nil {
		return err
	}

	return nil
	/*algo := &bytes.Buffer{}
	if err := binary.Write(algo, binary.LittleEndian, x509.ECDSAWithSHA512); err != nil {
		return err
	}
	// TODO(hugelgupf): use x509 package for all of this.
	// TODO(hugelgupf, later): no, please don't.
	return sw.w.WriteRecord(cpio.StaticFile("signature_algo", string(algo.Bytes()), 0700))*/
}
