// Copyright 2012 Harry de Boer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pnm

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
)

const (
	PBM int = 0
	PGM int = 1
	PPM int = 2
)

// packByte packs 8 pixels of bit depth 1 into a byte.
//
// The bits are packed with the first value as the most significant bit.
// If there are less than 8 values in bit, the remaining bits are 0.
// If there are more than 8 values in bit, these are ignored.
func packByte(bit []uint8) (b byte) {
	n := len(bit)
	if n > 8 {
		n = 8
	}
	for i := n - 1; i >= 0; i-- {
		b = b >> 1
		if bit[i] == 0 {
			b += 128
		}
	}
	return b
}

func encodePBM(w io.Writer, m image.Image) error {
	b := m.Bounds()
	// write header
	_, err := fmt.Fprintf(w, "P4\n%d %d\n", b.Dx(), b.Dy())
	if err != nil {
		return err
	}
	cm := make(color.Palette, 2)
	cm[0] = color.Gray{255}
	cm[1] = color.Gray{0}

	// write raster
	byteCount := b.Dx() / 8
	if b.Dx()%8 != 0 {
		byteCount += 1
	}
	row := make([]uint8, b.Dx())
	packedRow := make([]byte, byteCount)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		// Read row and convert to black/white.
		for x := b.Min.X; x < b.Max.X; x++ {
			c := cm.Convert(m.At(x, y)).(color.Gray)
			row[x-b.Min.X] = c.Y
		}

		// Pack values into and write
		i := 0
		x := 0
		for x < b.Dx() {
			n := b.Dx() - x
			if n > 8 {
				n = 8
			}
			packedRow[i] = packByte(row[x : x+n])
			x += n
			i++
		}
		if _, err := w.Write(packedRow); err != nil {
			return err
		}
	}
	return nil
}

func encodePGM(w io.Writer, m image.Image, maxvalue int) error {
	b := m.Bounds()
	// write header
	_, err := fmt.Fprintf(w, "P5\n%d %d\n%d\n", b.Dx(), b.Dy(), maxvalue)
	if err != nil {
		return err
	}

	// write raster
	cm := color.GrayModel
	row := make([]uint8, b.Dx())
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := cm.Convert(m.At(x, y)).(color.Gray)
			row[x-b.Min.X] = c.Y
		}
		if _, err := w.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func encodePPM(w io.Writer, m image.Image, maxvalue int) error {
	b := m.Bounds()
	// write header
	_, err := fmt.Fprintf(w, "P6\n%d %d\n%d\n", b.Dx(), b.Dy(), maxvalue)
	if err != nil {
		return err
	}

	// write raster
	cm := color.RGBAModel
	row := make([]uint8, b.Dx()*3)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		i := 0
		for x := b.Min.X; x < b.Max.X; x++ {
			c := cm.Convert(m.At(x, y)).(color.RGBA)
			row[i] = c.R
			row[i+1] = c.G
			row[i+2] = c.B
			i += 3
		}
		if _, err := w.Write(row); err != nil {
			return err
		}
	}
	return nil
}

// Encode writes an image.Image m to io.Writer w in PNM format.
//
// The specific format is determined by pnmType, this can be one of:
//  - pnm.PBM (black/white)
//  - pnm.PGM (grayscale)
//  - pnm.PPM (RGB)
// The image m is converted if necessary.
// Note that PGM/PPM always use 8 bits per channel at the moment and that
// maxvalue is always 255.
func Encode(w io.Writer, m image.Image, pnmType int) error {
	switch pnmType {
	case PBM:
		return encodePBM(w, m)
	case PGM:
		return encodePGM(w, m, 255)
	case PPM:
		return encodePPM(w, m, 255)
	}
	return errors.New("Invalid PNM type specified.")
}
