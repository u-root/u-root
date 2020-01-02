// Copyright 2013 Konstantin Kulikov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package framebuffer is an interface to linux framebuffer device.
package framebuffer

import (
	"os"
	"syscall"
	"unsafe"
)

// Framebuffer contains information about framebuffer.
type Framebuffer struct {
	dev      *os.File
	Finfo    fixedScreenInfo
	Vinfo    variableScreenInfo
	Data     []uint8
	restData []uint8
}

// Init opens framebuffer device, maps it to memory and saves its current contents.
func Init(dev string) (*Framebuffer, error) {
	var (
		fb  = new(Framebuffer)
		err error
	)

	fb.dev, err = os.OpenFile(dev, os.O_RDWR, os.ModeDevice)
	if err != nil {
		return nil, err
	}

	err = ioctl(fb.dev.Fd(), getFixedScreenInfo, unsafe.Pointer(&fb.Finfo))
	if err != nil {
		fb.dev.Close()
		return nil, err
	}

	err = ioctl(fb.dev.Fd(), getVariableScreenInfo, unsafe.Pointer(&fb.Vinfo))
	if err != nil {
		fb.dev.Close()
		return nil, err
	}

	fb.Data, err = syscall.Mmap(int(fb.dev.Fd()), 0, int(fb.Finfo.Smem_len+uint32(fb.Finfo.Smem_start&uint64(syscall.Getpagesize()-1))), protocolRead|protocolWrite, mapShared)
	if err != nil {
		fb.dev.Close()
		return nil, err
	}

	fb.restData = make([]byte, len(fb.Data))
	for i := range fb.Data {
		fb.restData[i] = fb.Data[i]
	}

	return fb, nil
}

// Close closes framebuffer device and restores its contents.
func (fb *Framebuffer) Close() {
	for i := range fb.restData {
		fb.Data[i] = fb.restData[i]
	}
	syscall.Munmap(fb.Data)
	fb.dev.Close()
}

// WritePixel changes pixel at x, y to specified color.
func (fb *Framebuffer) WritePixel(x, y int, red, green, blue, alpha uint8) {
	offset := (int(fb.Vinfo.Xoffset)+x)*(int(fb.Vinfo.Bits_per_pixel)/8) + (int(fb.Vinfo.Yoffset)+y)*int(fb.Finfo.Line_length)
	fb.Data[offset] = blue
	fb.Data[offset+1] = green
	fb.Data[offset+2] = red
	fb.Data[offset+3] = alpha
}

// Clear fills screen with specified color
func (fb *Framebuffer) Clear(red, green, blue, alpha uint8) {
	w, h := fb.Size()
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			fb.WritePixel(i, j, red, green, blue, alpha)
		}
	}
}

// Bpp returns bytes per pixel for a framebuffer.
func (fb *Framebuffer) Bpp() (bpp int) {
	return int(fb.Vinfo.Bits_per_pixel) / 8
}

// Stride returns line length of a framebuffer.
func (fb *Framebuffer) Stride() (stride int) {
	return int(fb.Finfo.Line_length) / fb.Bpp()
}

// Size returns dimensions of a framebuffer.
func (fb *Framebuffer) Size() (width, height int) {
	return int(fb.Vinfo.Xres), int(fb.Vinfo.Yres)
}

func ioctl(fd uintptr, cmd uintptr, data unsafe.Pointer) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, cmd, uintptr(data))
	if errno != 0 {
		return os.NewSyscallError("IOCTL", errno)
	}
	return nil
}
