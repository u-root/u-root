// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sulog works around some log/slog inefficiencies.
package sulog

import (
	"context"
	"fmt"
	"log/slog"

	"golang.org/x/exp/constraints"
)

type LevelHandler struct {
	slog.Handler

	Level slog.Level
}

func (lh LevelHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= lh.Level
}

func SetLevelHandlerDefault(l slog.Level) {
	h := slog.Default().Handler()
	// Assign to slog.Default due to https://github.com/golang/go/issues/61892
	*slog.Default() = *slog.New(LevelHandler{Handler: h, Level: l})
}

func SetDebugDefault() {
	SetLevelHandlerDefault(slog.LevelDebug)
}

type Hexable interface {
	constraints.Integer | []byte
}

func HexValue[T Hexable](value T) slog.Value {
	return slog.StringValue(fmt.Sprintf("%#x", value))
}
