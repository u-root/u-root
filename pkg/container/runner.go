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
	"fmt"
	"os"

	"github.com/opencontainers/runc/libcontainer"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/sys/unix"
)

type runner struct {
	enableSubreaper bool
	detach          bool
	shouldDestroy   bool
	consoleSocket   string
	pidFile         string
	container       libcontainer.Container
	listenFDs       []*os.File
	notifySocket    *notifySocket
}

func (r *runner) run(config *specs.Process) (int, error) {
	// Check the terminal settings.
	if r.detach && config.Terminal && r.consoleSocket == "" {
		return -1, fmt.Errorf("cannot allocate tty if runc will detach without setting console socket")
	}
	if (!r.detach || !config.Terminal) && r.consoleSocket != "" {
		return -1, fmt.Errorf("cannot use console socket if runc will not detach or allocate tty")
	}

	// Create the process.
	process, err := newProcess(*config)
	if err != nil {
		r.destroy()
		return -1, err
	}

	// Setup the listen file descriptors.
	if len(r.listenFDs) > 0 {
		process.Env = append(process.Env, fmt.Sprintf("LISTEN_FDS=%d", len(r.listenFDs)), "LISTEN_PID=1")
		process.ExtraFiles = append(process.ExtraFiles, r.listenFDs...)
	}

	// Get the rootuid.
	rootuid, err := r.container.Config().HostRootUID()
	if err != nil {
		r.destroy()
		return -1, err
	}

	// Get the rootgid.
	rootgid, err := r.container.Config().HostRootGID()
	if err != nil {
		r.destroy()
		return -1, err
	}

	// Setting up IO is a two stage process. We need to modify process to deal
	// with detaching containers, and then we get a tty after the container has
	// started.
	handler := newSignalHandler(r.enableSubreaper, r.notifySocket)
	tty, err := setupIO(process, rootuid, rootgid, config.Terminal, r.detach, r.consoleSocket)
	if err != nil {
		r.destroy()
		return -1, err
	}
	defer tty.Close()

	// Run the container.
	if err := r.container.Run(process); err != nil {
		r.destroy()
		tty.Close()
		return -1, err
	}

	// Wait for the tty.
	if err := tty.waitConsole(); err != nil {
		r.terminate(process)
		r.destroy()
		tty.Close()
		return -1, err
	}

	// Close after start the tty.
	if err = tty.ClosePostStart(); err != nil {
		r.terminate(process)
		r.destroy()
		tty.Close()
		return -1, err
	}

	// Create the pid file.
	if r.pidFile != "" {
		if err := createPidFile(r.pidFile, process); err != nil {
			r.terminate(process)
			r.destroy()
			tty.Close()
			return -1, err
		}
	}

	// Forward the handler.
	status, err := handler.forward(process, tty, r.detach)
	if err != nil {
		r.terminate(process)
	}

	// Return early if we are detaching.
	if r.detach {
		return 0, nil
	}

	// Cleanup.
	r.destroy()

	return status, err
}

func (r *runner) destroy() {
	if r.shouldDestroy {
		destroy(r.container)
	}
}

func (r *runner) terminate(p *libcontainer.Process) {
	_ = p.Signal(unix.SIGKILL)
	_, _ = p.Wait()
}
