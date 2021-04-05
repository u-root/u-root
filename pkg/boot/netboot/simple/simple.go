// Copyright 2017-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simple

import (
	"context"
	"fmt"
	"io"
	l "log"
	"math"
	"net/url"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/fit"
	"github.com/u-root/u-root/pkg/curl"
)

// FetchAndProbe fetches the file at the specified URL and checks if it is an
// Image file type rather than a config such as ipxe.
// TODO: detect nonFIT multiboot and bzImage Linux kernel files
func FetchAndProbe(ctx context.Context, u *url.URL, s curl.Schemes) ([]boot.OSImage, error) {
	file, err := s.Fetch(ctx, u)
	if err != nil {
		return nil, err
	}
	var images []boot.OSImage

	fimgs, err := fit.ParseConfig(io.NewSectionReader(file, 0, math.MaxInt64))
	if err == nil {
		for i := range fimgs {
			images = append(images, &fimgs[i])
		}
	} else {
		l.Printf("Parsing boot file as FIT image failed: %v", err)
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("exhausted all supported simple file types")
	}
	return images, nil
}
