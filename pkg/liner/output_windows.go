// SPDX-License-Identifier: MIT

// Copyright © 2012 Peter Harris
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice (including the next
// paragraph) shall be included in all copies or substantial portions of the
// Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package liner

import (
	"unsafe"
)

type coord struct {
	x, y int16
}
type smallRect struct {
	left, top, right, bottom int16
}

type consoleScreenBufferInfo struct {
	dwSize              coord
	dwCursorPosition    coord
	wAttributes         int16
	srWindow            smallRect
	dwMaximumWindowSize coord
}

func (s *State) cursorPos(x int) {
	var sbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(s.hOut), uintptr(unsafe.Pointer(&sbi)))
	procSetConsoleCursorPosition.Call(uintptr(s.hOut),
		uintptr(int(x)&0xFFFF|int(sbi.dwCursorPosition.y)<<16))
}

func (s *State) eraseLine() {
	var sbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(s.hOut), uintptr(unsafe.Pointer(&sbi)))
	var numWritten uint32
	procFillConsoleOutputCharacter.Call(uintptr(s.hOut), uintptr(' '),
		uintptr(sbi.dwSize.x-sbi.dwCursorPosition.x),
		uintptr(int(sbi.dwCursorPosition.x)&0xFFFF|int(sbi.dwCursorPosition.y)<<16),
		uintptr(unsafe.Pointer(&numWritten)))
}

func (s *State) eraseScreen() {
	var sbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(s.hOut), uintptr(unsafe.Pointer(&sbi)))
	var numWritten uint32
	procFillConsoleOutputCharacter.Call(uintptr(s.hOut), uintptr(' '),
		uintptr(sbi.dwSize.x)*uintptr(sbi.dwSize.y),
		0,
		uintptr(unsafe.Pointer(&numWritten)))
	procSetConsoleCursorPosition.Call(uintptr(s.hOut), 0)
}

func (s *State) moveUp(lines int) {
	var sbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(s.hOut), uintptr(unsafe.Pointer(&sbi)))
	procSetConsoleCursorPosition.Call(uintptr(s.hOut),
		uintptr(int(sbi.dwCursorPosition.x)&0xFFFF|(int(sbi.dwCursorPosition.y)-lines)<<16))
}

func (s *State) moveDown(lines int) {
	var sbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(s.hOut), uintptr(unsafe.Pointer(&sbi)))
	procSetConsoleCursorPosition.Call(uintptr(s.hOut),
		uintptr(int(sbi.dwCursorPosition.x)&0xFFFF|(int(sbi.dwCursorPosition.y)+lines)<<16))
}

func (s *State) emitNewLine() {
	// windows doesn't need to omit a new line
}

func (s *State) getColumns() {
	var sbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(s.hOut), uintptr(unsafe.Pointer(&sbi)))
	s.columns = int(sbi.dwSize.x)
}
