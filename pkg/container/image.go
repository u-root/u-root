// The MIT License (MIT)
//
// Copyright (c) 2018 The Genuinetools Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package container

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/storage"
	"github.com/docker/distribution/registry/storage/driver/filesystem"
	"github.com/docker/docker/pkg/archive"
	"github.com/genuinetools/reg/repoutils"
	bindata "github.com/jteeuwen/go-bindata"
)

// EmbedImage pulls a docker image locally. Creates a tarball of it's contents
// and then embeds the tarball as binary data into an output bindata.go file.
func EmbedImage(image string) error {
	// Get the current working directory.
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Create our output path
	output := filepath.Join(wd, "bindata.go")

	// Create the temporary directory for the image contents.
	tmpd, err := ioutil.TempDir("", "container-lib")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpd) // Cleanup on complete.

	// Create our tarball path.
	tarball := filepath.Join(tmpd, DefaultTarballPath)

	// Create our image root and state.
	root := filepath.Join(tmpd, "root")
	state := filepath.Join(tmpd, "state")

	// Create the rootfs
	if err := createRootFS(image, root, state); err != nil {
		return err
	}

	// Create the tar.
	tar, err := archive.Tar(root, archive.Gzip)
	if err != nil {
		return fmt.Errorf("create tar failed: %v", err)
	}

	// Create the tarball writer.
	writer, err := os.Create(tarball)
	if err != nil {
		return err
	}
	defer writer.Close() // Close the writer.

	if _, err := io.Copy(writer, tar); err != nil {
		return fmt.Errorf("copy tarball failed: %v", err)
	}

	// Create the bindata config.
	bc := bindata.NewConfig()
	bc.Input = []bindata.InputConfig{
		{
			Path:      tarball,
			Recursive: false,
		},
	}
	bc.Output = output
	bc.Package = "main"
	bc.NoMetadata = true
	bc.Prefix = filepath.Dir(tarball)

	if err := bindata.Translate(bc); err != nil {
		return fmt.Errorf("bindata failed: %v", err)
	}

	return nil
}

// createRootFS creates the base filesystem for a docker image.
// It will pull the base image if it does not exist locally.
// This function takes in a image name and the directory where the
// rootfs should be created.
func createRootFS(image, rootfs, state string) error {
	// Create the context.
	ctx := context.Background()

	// Create the new local registry storage.
	local, err := storage.NewRegistry(ctx, filesystem.New(filesystem.DriverParameters{
		RootDirectory: state,
		MaxThreads:    100,
	}))
	if err != nil {
		return fmt.Errorf("creating new registry storage failed: %v", err)
	}

	// Parse the repository name.
	name, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return fmt.Errorf("not a valid image %q: %v", image, err)
	}
	// Add latest to the image name if it is empty.
	name = reference.TagNameOnly(name)

	// Get the tag for the repo.
	_, tag, err := repoutils.GetRepoAndRef(image)
	if err != nil {
		return err
	}

	// Create the local repository.
	repo, err := local.Repository(ctx, name)
	if err != nil {
		return fmt.Errorf("creating local repository for %q failed: %v", reference.Path(name), err)
	}

	// Create the manifest service.
	ms, err := repo.Manifests(ctx)
	if err != nil {
		return fmt.Errorf("creating manifest service failed: %v", err)
	}

	// Get the specific tag.
	td, err := repo.Tags(ctx).Get(ctx, tag)
	// Check if we got an unknown error, that means the tag does not exist.
	if err != nil && strings.Contains(err.Error(), "unknown") {
		log.Println("image not found locally, pulling the image")

		// Pull the image.
		if err := pull(ctx, local, name, tag); err != nil {
			return fmt.Errorf("pulling failed: %v", err)
		}

		// Try to get the tag again.
		td, err = repo.Tags(ctx).Get(ctx, tag)
	}
	if err != nil {
		return fmt.Errorf("getting local repository tag %q failed: %v", tag, err)
	}

	// Get the specific manifest for the tag.
	manifest, err := ms.Get(ctx, td.Digest)
	if err != nil {
		return fmt.Errorf("getting local manifest for digest %q failed: %v", td.Digest.String(), err)
	}

	blobStore := repo.Blobs(ctx)
	for i, ref := range manifest.References() {
		if i == 0 {
			fmt.Printf("skipping config %v\n", ref.Digest.String())
			continue
		}
		fmt.Printf("unpacking %v\n", ref.Digest.String())
		layer, err := blobStore.Open(ctx, ref.Digest)
		if err != nil {
			return fmt.Errorf("getting blob %q failed: %v", ref.Digest.String(), err)
		}

		// Unpack the tarfile to the mount path.
		// FROM: https://godoc.org/github.com/moby/moby/pkg/archive#TarOptions
		if err := archive.Untar(layer, rootfs, &archive.TarOptions{
			NoLchown: true,
		}); err != nil {
			return fmt.Errorf("error extracting tar for %q: %v", ref.Digest.String(), err)
		}
	}

	return nil
}
