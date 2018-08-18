// Created by cgo -godefs - DO NOT EDIT
// cgo.exe -godefs types_src_windows.go

package sys

type (
	Coord struct {
		X int16
		Y int16
	}
	InputRecord struct {
		EventType uint16
		Pad_cgo_0 [2]byte
		Event     [16]byte
	}

	KeyEvent struct {
		BKeyDown          int32
		WRepeatCount      uint16
		WVirtualKeyCode   uint16
		WVirtualScanCode  uint16
		UChar             [2]byte
		DwControlKeyState uint32
	}
	MouseEvent struct {
		DwMousePosition   Coord
		DwButtonState     uint32
		DwControlKeyState uint32
		DwEventFlags      uint32
	}
	WindowBufferSizeEvent struct {
		DwSize Coord
	}
	MenuEvent struct {
		DwCommandId uint32
	}
	FocusEvent struct {
		BSetFocus int32
	}
)

const (
	KEY_EVENT                = 0x1
	MOUSE_EVENT              = 0x2
	WINDOW_BUFFER_SIZE_EVENT = 0x4
	MENU_EVENT               = 0x8
	FOCUS_EVENT              = 0x10
)
