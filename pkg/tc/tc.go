// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl

import (
	"errors"
	"io"

	"github.com/florianl/go-tc"
)

var (
	ErrNotEnoughArgs  = errors.New("not enough arguments")
	ErrInvalidArg     = errors.New("invalid argument in list")
	ErrNotImplemented = errors.New("not implemented")
	ErrOutOfBounds    = errors.New("integer argument out of bounds")
	ErrExitAfterHelp  = errors.New("exit after help message")
)

type Tctl interface {
	ShowQdisc(io.Writer, *Args) error
	AddQdisc(io.Writer, *Args) error
	DeleteQdisc(io.Writer, *Args) error
	ReplaceQdisc(io.Writer, *Args) error
	ChangeQdisc(io.Writer, *Args) error
	ShowClass(io.Writer, *Args) error
	AddClass(io.Writer, *Args) error
	DeleteClass(io.Writer, *Args) error
	ReplaceClass(io.Writer, *Args) error
	ChangeClass(io.Writer, *Args) error
	ShowFilter(io.Writer, *FArgs) error
	AddFilter(io.Writer, *FArgs) error
	DeleteFilter(io.Writer, *FArgs) error
	ReplaceFilter(io.Writer, *FArgs) error
	ChangeFilter(io.Writer, *FArgs) error
	GetFilter(io.Writer, *FArgs) error
}

type Trafficctl struct {
	*tc.Tc
}
