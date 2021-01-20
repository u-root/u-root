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
// Package container allows for running a process in a container.
package container

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/coreos/go-systemd/activation"
	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/libcontainer/cgroups/systemd"
	"github.com/opencontainers/runc/libcontainer/specconv"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// Container defines the behavior and settings for a container object.
type Container struct {
	ID               string
	Spec             *specs.Spec
	PIDFile          string
	ConsoleSocket    string
	Root             string
	Detach           bool
	UseSystemdCgroup bool
	NoPivotRoot      bool
	NoNewKeyring     bool
	Rootless         bool
}

// Run starts the container. It returns the exit status or -1 and an
// error. Signals sent to the current process will be forwarded to container.
func (c *Container) Run() (int, error) {
	var err error

	// Convert pid-file to an absolute path so we can write to the
	// right file after chdir to bundle.
	if c.PIDFile != "" {
		c.PIDFile, err = filepath.Abs(c.PIDFile)
		if err != nil {
			return -1, err
		}
	}

	// Get the absolute path to the root.
	c.Root, err = filepath.Abs(c.Root)
	if err != nil {
		return -1, err
	}

	notifySocket := newNotifySocket(c.ID, c.Root)
	if notifySocket != nil {
		// Setup the spec for the notify socket.
		notifySocket.setupSpec(c.Spec)
	}

	// Create the libcontainer config.
	config, err := specconv.CreateLibcontainerConfig(&specconv.CreateOpts{
		CgroupName:       c.ID,
		UseSystemdCgroup: c.UseSystemdCgroup,
		NoPivotRoot:      c.NoPivotRoot,
		NoNewKeyring:     c.NoNewKeyring,
		Spec:             c.Spec,
		Rootless:         c.Rootless,
	})
	if err != nil {
		return -1, err
	}

	// Setup the cgroups manager. Default is cgroupfs.
	cgroupManager := libcontainer.Cgroupfs
	if c.UseSystemdCgroup {
		if systemd.UseSystemd() {
			cgroupManager = libcontainer.SystemdCgroups
		} else {
			return -1, fmt.Errorf("systemd cgroup flag passed, but systemd support for managing cgroups is not available")
		}
	}

	// We resolve the paths for {newuidmap,newgidmap} from the context of runc,
	// to avoid doing a path lookup in the nsexec context. TODO: The binary
	// names are not currently configurable.
	newuidmap, err := exec.LookPath("newuidmap")
	if err != nil {
		newuidmap = ""
	}
	newgidmap, err := exec.LookPath("newgidmap")
	if err != nil {
		newgidmap = ""
	}

	// Create the new libcontainer factory.
	factory, err := libcontainer.New(c.Root, cgroupManager, nil, nil,
		libcontainer.NewuidmapPath(newuidmap),
		libcontainer.NewgidmapPath(newgidmap))
	if err != nil {
		return -1, err
	}

	// Create the factory.
	container, err := factory.Create(c.ID, config)
	if err != nil {
		return -1, err
	}

	if notifySocket != nil {
		// Setup the socket for the notify socket.
		err := notifySocket.setupSocket()
		if err != nil {
			return -1, err
		}
	}

	// Support on-demand socket activation by passing file descriptors into
	// the container init process.
	listenFDs := []*os.File{}
	if os.Getenv("LISTEN_FDS") != "" {
		listenFDs = activation.Files(false)
	}

	// Initialize the runner.
	r := &runner{
		enableSubreaper: true,
		shouldDestroy:   true,
		container:       container,
		listenFDs:       listenFDs,
		notifySocket:    notifySocket,
		consoleSocket:   c.ConsoleSocket,
		detach:          c.Detach,
		pidFile:         c.PIDFile,
	}

	// Run the process.
	return r.run(c.Spec.Process)
}
