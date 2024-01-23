// Copyright 2019-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fb

import (
	"fmt"
	"image"
	"os"

	"github.com/orangecms/go-framebuffer/framebuffer"
)

const fbdev = "/dev/fb0"

func DrawOnBufAt(
	buf []byte,
	img image.Image,
	posx int,
	posy int,
	stride int,
	bpp int,
) {
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			offset := bpp * ((posy+y)*stride + posx + x)
			// framebuffer is BGR(A)
			buf[offset+0] = byte(b)
			buf[offset+1] = byte(g)
			buf[offset+2] = byte(r)
			if bpp >= 4 {
				buf[offset+3] = byte(a)
			}
		}
	}
}

// FbInit initializes a frambuffer by querying ioctls and returns the width and
// height in pixels, the stride, and the bytes per pixel
func FbInit() (int, int, int, int, error) {
	fbo, err := framebuffer.Init(fbdev)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	width, height := fbo.Size()
	stride := fbo.Stride()
	bpp := fbo.Bpp()
	fmt.Fprintf(os.Stdout, "Framebuffer resolution: %v %v %v %v\n", width, height, stride, bpp)
	return width, height, stride, bpp, nil
}

func DrawImageAt(img image.Image, posx int, posy int) error {
	width, height, stride, bpp, err := FbInit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Framebuffer init error: %v\n", err)
		// fallback values, probably a bad assumption
		width, height, stride, bpp = 1920, 1080, 1920*4, 4
	}
	buf := make([]byte, width*height*bpp)
	DrawOnBufAt(buf, img, posx, posy, stride, bpp)
	err = os.WriteFile(fbdev, buf, 0o600)
	if err != nil {
		return fmt.Errorf("error writing to framebuffer: %w", err)
	}
	return nil
}

func DrawScaledOnBufAt(
	buf []byte,
	img image.Image,
	posx int,
	posy int,
	factor int,
	stride int,
	bpp int,
) {
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			for sx := 1; sx <= factor; sx++ {
				for sy := 1; sy <= factor; sy++ {
					offset := bpp * ((posy+y*factor+sy)*stride + posx + x*factor + sx)
					buf[offset+0] = byte(b)
					buf[offset+1] = byte(g)
					buf[offset+2] = byte(r)
					if bpp == 4 {
						buf[offset+3] = byte(a)
					}
				}
			}
		}
	}
}

func DrawScaledImageAt(img image.Image, posx int, posy int, factor int) error {
	width, height, stride, bpp, err := FbInit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Framebuffer init error: %v\n", err)
	}
	buf := make([]byte, width*height*bpp)
	DrawScaledOnBufAt(buf, img, posx, posy, factor, stride, bpp)
	err = os.WriteFile(fbdev, buf, 0o600)
	if err != nil {
		return fmt.Errorf("error writing to framebuffer: %w", err)
	}
	return nil
}
