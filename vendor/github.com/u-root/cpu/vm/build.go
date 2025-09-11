//go:build generate

package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/tarutil"
)

type build struct {
	os        string
	arch      string
	version   string
	kernel    string
	container string
	env       []string
	cmd       []string
}

var builds = []build{
	{
		os:        "linux",
		arch:      "amd64",
		version:   "v1",
		kernel:    "bzImage",
		container: "ghcr.io/hugelgupf/vmtest/kernel-amd64:main",
		env:       []string{"GOARCH=amd64", "GOAMD64=v1"},
		cmd:       []string{},
	},
	{
		os:        "linux",
		arch:      "arm64",
		version:   "v1",
		kernel:    "Image",
		container: "ghcr.io/hugelgupf/vmtest/kernel-arm64:main",
		env:       []string{"GOARCH=arm64"},
		cmd:       []string{},
	},
	{
		os:        "linux",
		arch:      "arm",
		kernel:    "zImage",
		container: "ghcr.io/hugelgupf/vmtest/kernel-arm:main",
		env:       []string{"GOARCH=arm", "GOARM=5"},
		cmd:       []string{},
	},
	{
		os:        "linux",
		arch:      "riscv64",
		kernel:    "Image",
		container: "ghcr.io/hugelgupf/vmtest/kernel-riscv64:main",
		env:       []string{"GOARCH=riscv64"},
		cmd:       []string{},
	},
}

func main() {
	env := []string{"CGO_ENABLED=0"}
	for _, b := range builds {
		log.Printf("Build %v", b)
		n := fmt.Sprintf("initramfs_%s_%s.cpio", b.os, b.arch)
		cmd := []string{"-initcmd=/bbin/cpud", "-defaultsh=", "-o=" + n,
			"../cmds/cpud",
			"../../u-root/cmds/core/dhclient",
		}
		c := exec.Command("u-root", cmd...)
		c.Stdout, c.Stderr = os.Stdout, os.Stderr
		c.Env = append(os.Environ(), append(env, b.env...)...)
		if err := c.Run(); err != nil {
			log.Fatal(err)
		}
		f, err := os.ReadFile(n)
		if err != nil {
			log.Fatal(err)
		}

		var newcpio bytes.Buffer
		rw := cpio.Newc.Writer(&newcpio)
		recs := cpio.Newc.Reader(bytes.NewReader(f))
		fixed := 0
		cpio.ForEachRecord(recs, func(r cpio.Record) error {
			switch r.Name {
			case "bbin/bb":
				fixed++
				r.Name = "bbin/cibb"
			case "bbin/init", "bbin/dhclient", "bbin/cpud":
				fixed++
				r.ReaderAt = bytes.NewReader([]byte("cibb"))
				r.Info.FileSize = 4
			}
			if err := rw.WriteRecord(r); err != nil {
				return fmt.Errorf("writing record %q failed: %w", r.Name, err)
			}
			return nil
		})

		if fixed < 2 {
			log.Fatal("Did not fix any entries in %q", n)
		}

		// because we modify the cpio, and one of the tests makes sure the cpio matches compressed,
		// write back the modified cpio.
		if err := os.WriteFile(n, newcpio.Bytes(), 0644); err != nil {
			log.Fatalf("writing back changed cpio %s:%v", n, err)
		}

		var out bytes.Buffer
		gz := gzip.NewWriter(&out)
		if _, err := gz.Write(newcpio.Bytes()); err != nil {
			log.Fatal(err)
		}
		if err := gz.Close(); err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(n+".gz", out.Bytes(), 0644); err != nil {
			log.Fatal(err)
		}

		if len(b.container) == 0 {
			continue
		}
		ref, err := name.ParseReference(b.container)
		if err != nil {
			log.Fatal(err)
		}

		img, err := crane.Pull(ref.Name())
		if err != nil {
			log.Fatal(err)
		}

		r := mutate.Extract(img)

		opts := &tarutil.Opts{}
		opts.Filters = []tarutil.Filter{tarutil.SafeFilter, func(h *tar.Header) bool {
			return h.Name == b.kernel
		},
		}

		if err := tarutil.ExtractDir(r, ".", opts); err != nil {
			log.Fatal(err)
		}

		// tarutil does not let us extract a file with a name into a []byte.Unfortunate.
		kname := fmt.Sprintf("kernel_%s_%s", b.os, b.arch)
		if err := os.Rename(b.kernel, kname); err != nil {
			log.Fatal(err)
		}
		if f, err = os.ReadFile(kname); err != nil {
			log.Fatal(err)
		}
		var kout bytes.Buffer
		gz = gzip.NewWriter(&kout)
		if _, err := gz.Write(f); err != nil {
			log.Fatal(err)
		}
		if err := gz.Close(); err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(kname+".gz", kout.Bytes(), 0644); err != nil {
			log.Fatal(err)
		}

	}
}
