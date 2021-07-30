// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs _defs-source.go

package framebuffer

type fixedScreenInfo struct {
	Id           [16]int8
	Smem_start   uint64
	Smem_len     uint32
	Type         uint32
	Type_aux     uint32
	Visual       uint32
	Xpanstep     uint16
	Ypanstep     uint16
	Ywrapstep    uint16
	Pad_cgo_0    [2]byte
	Line_length  uint32
	Pad_cgo_1    [4]byte
	Mmio_start   uint64
	Mmio_len     uint32
	Accel        uint32
	Capabilities uint16
	Reserved     [2]uint16
	Pad_cgo_2    [2]byte
}
type variableScreenInfo struct {
	Xres           uint32
	Yres           uint32
	Xres_virtual   uint32
	Yres_virtual   uint32
	Xoffset        uint32
	Yoffset        uint32
	Bits_per_pixel uint32
	Grayscale      uint32
	Red            bitField
	Green          bitField
	Blue           bitField
	Transp         bitField
	Nonstd         uint32
	Activate       uint32
	Height         uint32
	Width          uint32
	Accel_flags    uint32
	Pixclock       uint32
	Left_margin    uint32
	Right_margin   uint32
	Upper_margin   uint32
	Lower_margin   uint32
	Hsync_len      uint32
	Vsync_len      uint32
	Sync           uint32
	Vmode          uint32
	Rotate         uint32
	Colorspace     uint32
	Reserved       [4]uint32
}
type bitField struct {
	Offset uint32
	Length uint32
	Right  uint32
}

const (
	getFixedScreenInfo    uintptr = 0x4602
	getVariableScreenInfo uintptr = 0x4600
)

const (
	protocolRead  int = 0x1
	protocolWrite int = 0x2
	mapShared     int = 0x1
)
