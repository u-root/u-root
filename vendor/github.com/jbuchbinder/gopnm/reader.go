// Copyright 2012 Harry de Boer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pnm implements a PBM, PGM and PPM image decoder and encoder.
//
// The decoder can read files in both plain and raw format with 8 or 16 bits
// per channel. The encoder can only write files in plain format with 8 bits
// per channel.
//
// To only be able to load pnm images using image.Decode, use
//	import _ "github.com/jbuchbinder/gopnm"
//
// Not implemented are:
//	- Writing pnm files in raw format.
//	- Writing images with 16 bits per channel.
//	- Writing images with a custom Maxvalue.
//	- Reading/and writing PAM images.
// (I would be happy to accept patches for these.)
//
// Specifications can be found at http://netpbm.sourceforge.net/doc/#formats.
package pnm

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"unicode"
)

type PNMConfig struct {
	Width  int
	Height int
	Maxval int
	magic  string
}

func decodePlainBW(r io.Reader, c PNMConfig) (image.Image, error) {
	m := image.NewGray(image.Rect(0, 0, c.Width, c.Height))
	pixelCount := len(m.Pix)

	for i := 0; i < pixelCount; i++ {
		if _, err := fmt.Fscan(r, &m.Pix[i]); err != nil {
			return nil, err
		}
		if m.Pix[i] == 0 {
			m.Pix[i] = 255
		} else {
			m.Pix[i] = 0
		}
	}

	return m, nil
}

func decodePlainGray(r io.Reader, c PNMConfig) (image.Image, error) {
	m := image.NewGray(image.Rect(0, 0, c.Width, c.Height))
	pixelCount := len(m.Pix)

	for i := 0; i < pixelCount; i++ {
		if _, err := fmt.Fscan(r, &m.Pix[i]); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func decodePlainGray16(r io.Reader, c PNMConfig) (image.Image, error) {
	m := image.NewGray16(image.Rect(0, 0, c.Width, c.Height))
	var col uint16

	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			if _, err := fmt.Fscan(r, &col); err != nil {
				return nil, err
			}
			m.Set(x, y, color.Gray16{col})
		}
	}

	return m, nil
}

func decodePlainRGB(r io.Reader, c PNMConfig) (image.Image, error) {
	m := image.NewRGBA(image.Rect(0, 0, c.Width, c.Height))
	count := len(m.Pix)

	for i := 0; i < count; i += 4 {
		if _, err := fmt.Fscan(r, &m.Pix[i]); err != nil {
			return nil, err
		}
		if _, err := fmt.Fscan(r, &m.Pix[i+1]); err != nil {
			return nil, err
		}
		if _, err := fmt.Fscan(r, &m.Pix[i+2]); err != nil {
			return nil, err
		}
		m.Pix[i+3] = 0xff
	}

	return m, nil
}

func decodePlainRGB64(r io.Reader, c PNMConfig) (image.Image, error) {
	m := image.NewRGBA64(image.Rect(0, 0, c.Width, c.Height))
	var cr, cg, cb uint16

	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			if _, err := fmt.Fscan(r, &cr); err != nil {
				return nil, err
			}
			if _, err := fmt.Fscan(r, &cg); err != nil {
				return nil, err
			}
			if _, err := fmt.Fscan(r, &cb); err != nil {
				return nil, err
			}
			m.Set(x, y, color.RGBA64{cr, cg, cb, 0xffff})
		}
	}

	return m, nil
}

// unpackByte unpacks 8 one bit pixels from byte b into slice bit.
//
// The bits are unpacked such that the most significant bit becomes the
// first value in the slice. If there are less than 8 values in bit,the
// remaining bits are ignored. If there are more than 8 values in bit,
// these remain unchanged.
func unpackByte(bit []uint8, b byte) {
	n := len(bit)
	if n > 8 {
		n = 8
	}
	for i := 0; i < n; i++ {
		if b&128 == 0 {
			bit[i] = 255
		}
		b = b << 1
	}
}

func decodeRawBW(r io.Reader, c PNMConfig) (image.Image, error) {
	m := image.NewGray(image.Rect(0, 0, c.Width, c.Height))

	byteCount := c.Width / 8
	if c.Width%8 != 0 {
		byteCount += 1
	}
	row := make([]byte, byteCount)
	pos := 0

	for y := 0; y < c.Height; y++ {
		if _, err := io.ReadFull(r, row); err != nil {
			return nil, err
		}
		bitsLeft := c.Width
		for _, b := range row {
			n := bitsLeft
			if n > 8 {
				n = 8
			}
			unpackByte(m.Pix[pos:pos+n], b)
			bitsLeft -= n
			pos += n
		}
	}

	return m, nil
}

func decodeRawGray(r io.Reader, c PNMConfig) (image.Image, error) {
	m := image.NewGray(image.Rect(0, 0, c.Width, c.Height))
	_, err := io.ReadFull(r, m.Pix)
	return m, err
}

func decodeRawGray16(r io.Reader, c PNMConfig) (image.Image, error) {
	m := image.NewGray16(image.Rect(0, 0, c.Width, c.Height))
	_, err := io.ReadFull(r, m.Pix)
	return m, err
}

func decodeRawRGB(r io.Reader, c PNMConfig) (image.Image, error) {
	m := image.NewRGBA(image.Rect(0, 0, c.Width, c.Height))
	numPixels := c.Width * c.Height

	// Do a large read for the rgb data.
	if _, err := io.ReadFull(r, m.Pix[0:numPixels*3]); err != nil {
		return nil, err
	}

	// Repack to RGBA form.
	dstPos := (numPixels * 4) - 1
	srcPos := (numPixels * 3) - 1

	for dstPos > 0 {
		m.Pix[dstPos] = 0xff
		dstPos--
		m.Pix[dstPos] = m.Pix[srcPos]
		dstPos--
		srcPos--
		m.Pix[dstPos] = m.Pix[srcPos]
		dstPos--
		srcPos--
		m.Pix[dstPos] = m.Pix[srcPos]
		dstPos--
		srcPos--
	}

	return m, nil
}

func decodeRawRGB64(r io.Reader, c PNMConfig) (image.Image, error) {
	m := image.NewRGBA(image.Rect(0, 0, c.Width, c.Height))
	count := len(m.Pix)

	for i := 0; i < count; i += 8 {
		pixel := m.Pix[i : i+6]
		m.Pix[i+6] = 0xff
		m.Pix[i+7] = 0xff

		if _, err := io.ReadFull(r, pixel); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func decodePAM(r io.Reader, c PNMConfig) (image.Image, error) {
	return nil, errors.New("pnm: reading PAM images is not supported yet.")
}

// Decode reads a PNM image from r and returns it as an image.Image.
//
// The type of Image returned depends on the PNM contents:
//  - PBM: image.Gray with black = 0 and white = 255
//  - PGM: image.Gray or image.Gray16, values as in the file
//  - PPM: image.RGBA or image.RGBA64, values as in the file
//  - PAM: not supported (yet)
func Decode(r io.Reader) (image.Image, error) {
	br := bufio.NewReader(r)
	c, err := DecodeConfigPNM(br)

	if err != nil {
		err = fmt.Errorf("pnm: parsing header failed: %v", err)
		return nil, err
	}

	switch c.magic {
	case "P1":
		return decodePlainBW(br, c)
	case "P2":
		if c.Maxval < 256 {
			return decodePlainGray(br, c)
		} else {
			return decodePlainGray16(br, c)
		}
	case "P3":
		if c.Maxval < 256 {
			return decodePlainRGB(br, c)
		} else {
			return decodePlainRGB64(br, c)
		}
	case "P4":
		return decodeRawBW(br, c)
	case "P5":
		if c.Maxval < 256 {
			return decodeRawGray(br, c)
		} else {
			return decodeRawGray16(br, c)
		}
	case "P6":
		if c.Maxval < 256 {
			return decodeRawRGB(br, c)
		} else {
			return decodeRawRGB64(br, c)
		}
	case "P7":
		return decodePAM(br, c)
	}

	return nil, fmt.Errorf("pnm: could not decode, invalid magic value %s", c.magic[0:2])
}

// skipComments skips all comments (and whitespace) that may occur between PNM
// header tokens.
//
// The singleSpace argument is used to scan comments between the header and the
// raster data where only a single whitespace delimiter is allowed. This
// prevents scanning the image data.
func skipComments(r *bufio.Reader, singleSpace bool) (err error) {
	var c byte

	for {
		// Skip whitespace
		c, err = r.ReadByte()
		for unicode.IsSpace(rune(c)) {
			if c, err = r.ReadByte(); err != nil {
				return err
			}
			if singleSpace {
				break
			}
		}
		// If there are no more comments, unread the last byte and return.
		if c != '#' {
			r.UnreadByte()
			return nil
		}
		// A comment ends with a newline or carriage return.
		for c != '\n' && c != '\r' {
			if c, err = r.ReadByte(); err != nil {
				return
			}
		}
	}
}

// DecodeConfigPNM reads and returns header data of PNM files.
//
// This may be useful to obtain the actual file type and for files that have a
// Maxval other than the maximum supported Maxval. To apply gamma correction
// this value is needed. Note that gamma correction is not performed by the
// decoder.
func DecodeConfigPNM(r *bufio.Reader) (c PNMConfig, err error) {
	// PNM magic number
	if _, err = fmt.Fscan(r, &c.magic); err != nil {
		return
	}
	switch c.magic {
	case "P1", "P2", "P3", "P4", "P5", "P6":
	case "P7":
		return c, errors.New("pnm: reading PAM images is not supported (yet).")
	default:
		return c, errors.New("pnm: invalid format " + c.magic[0:2])
	}

	// Image width
	if err = skipComments(r, false); err != nil {
		return
	}
	if _, err = fmt.Fscan(r, &c.Width); err != nil {
		return c, errors.New("pnm: could not read image width, " + err.Error())
	}
	// Image height
	if err = skipComments(r, false); err != nil {
		return
	}
	if _, err = fmt.Fscan(r, &c.Height); err != nil {
		return c, errors.New("pnm: could not read image height, " + err.Error())
	}
	// Number of colors, only for gray and color images.
	// For black and white images this is 2, obviously.
	if c.magic == "P1" || c.magic == "P4" {
		c.Maxval = 2
	} else {
		if err = skipComments(r, false); err != nil {
			return
		}
		if _, err = fmt.Fscan(r, &c.Maxval); err != nil {
			return c, errors.New("pnm: could not read number of colors, " + err.Error())
		}
	}

	if c.Maxval > 65535 || c.Maxval <= 0 {
		err = fmt.Errorf("pnm: maximum depth is 16 bit (65,535) colors but %d colors found", c.Maxval)
		return
	}

	// Skip comments after header.
	if err = skipComments(r, true); err != nil {
		return
	}

	return c, nil
}

// DecodeConfig returns the color model and dimensions of a PNM image without
// decoding the entire image.
func DecodeConfig(r io.Reader) (image.Config, error) {
	br := bufio.NewReader(r)
	c, err := DecodeConfigPNM(br)
	if err != nil {
		return image.Config{}, err
	}

	var cm color.Model
	switch c.magic {
	case "P1", "P4":
		cm = color.GrayModel
	case "P2", "P5":
		if c.Maxval < 256 {
			cm = color.GrayModel
		} else {
			cm = color.Gray16Model
		}
	case "P3", "P6":
		if c.Maxval < 256 {
			cm = color.RGBAModel
		} else {
			cm = color.RGBA64Model
		}
	}

	return image.Config{
		ColorModel: cm,
		Width:      c.Width,
		Height:     c.Height,
	}, nil
}

func init() {
	image.RegisterFormat("pbm ascii (black/white)", "P1", Decode, DecodeConfig)
	image.RegisterFormat("pgm ascii (grayscale)", "P2", Decode, DecodeConfig)
	image.RegisterFormat("ppm ascii (rgb)", "P3", Decode, DecodeConfig)
	image.RegisterFormat("pbm raw (black/white)", "P4", Decode, DecodeConfig)
	image.RegisterFormat("pgm raw (grayscale)", "P5", Decode, DecodeConfig)
	image.RegisterFormat("ppm raw (rgb)", "P6", Decode, DecodeConfig)
	//image.RegisterFormat("pam", "P7", Decode, DecodeConfig)
}
