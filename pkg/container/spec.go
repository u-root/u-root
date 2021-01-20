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
	"github.com/containerd/containerd/contrib/seccomp"
	aaprofile "github.com/docker/docker/profiles/apparmor"
	"github.com/opencontainers/runc/libcontainer/apparmor"
	"github.com/opencontainers/runc/libcontainer/specconv"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

const (
	// DefaultApparmorProfile is the default apparmor profile for the containers.
	DefaultApparmorProfile = "docker-default"
)

// SpecOpts defines the options available for a spec.
type SpecOpts struct {
	Rootless bool
	Readonly bool
	Terminal bool
	Args     []string
	Mounts   []specs.Mount
	Hooks    *specs.Hooks
}

// Spec returns a default oci spec with some options being passed.
func Spec(opts SpecOpts) *specs.Spec {
	// Initialize the spec.
	spec := specconv.Example()

	// Set the spec to be rootless.
	if opts.Rootless {
		specconv.ToRootless(spec)
	}

	// Setup readonly fs in spec.
	spec.Root.Readonly = opts.Readonly

	// Setup tty in spec.
	spec.Process.Terminal = opts.Terminal

	// Pass in any hooks to the spec.
	spec.Hooks = opts.Hooks

	// Set the default seccomp profile.
	spec.Linux.Seccomp = seccomp.DefaultProfile(spec)

	// Install the default apparmor profile.
	if apparmor.IsEnabled() {
		// Check if we have the docker-default apparmor profile loaded.
		if _, err := aaprofile.IsLoaded(DefaultApparmorProfile); err == nil {
			spec.Process.ApparmorProfile = DefaultApparmorProfile
		}
	}

	if opts.Args != nil {
		spec.Process.Args = opts.Args
	}

	if opts.Mounts != nil {
		spec.Mounts = append(spec.Mounts, opts.Mounts...)
	}

	return spec
}
