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
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
)

type notifySocket struct {
	socket     *net.UnixConn
	host       string
	socketPath string
}

func newNotifySocket(id, root string) *notifySocket {
	if os.Getenv("NOTIFY_SOCKET") == "" {
		// Return early if we do not have a NOTIFY_SOCKET.
		return nil
	}

	path := filepath.Join(filepath.Join(root, id), "notify.sock")

	notifySocket := &notifySocket{
		socket:     nil,
		host:       os.Getenv("NOTIFY_SOCKET"),
		socketPath: path,
	}

	return notifySocket
}

func (s *notifySocket) Close() error {
	return s.socket.Close()
}

// If systemd is supporting sd_notify protocol, this function will add support
// for sd_notify protocol from within the container.
func (s *notifySocket) setupSpec(spec *specs.Spec) {
	mount := specs.Mount{Destination: s.host, Type: "bind", Source: s.socketPath, Options: []string{"bind"}}
	spec.Mounts = append(spec.Mounts, mount)
	spec.Process.Env = append(spec.Process.Env, fmt.Sprintf("NOTIFY_SOCKET=%s", s.host))
}

func (s *notifySocket) setupSocket() error {
	addr := net.UnixAddr{
		Name: s.socketPath,
		Net:  "unixgram",
	}

	socket, err := net.ListenUnixgram("unixgram", &addr)
	if err != nil {
		return err
	}

	s.socket = socket
	return nil
}

// pid1 must be set only with -d, as it is used to set the new process as the main process
// for the service in butts
func (s *notifySocket) run(pid1 int) {
	buf := make([]byte, 512)
	notifySocketHostAddr := net.UnixAddr{Name: s.host, Net: "unixgram"}
	client, err := net.DialUnix("unixgram", nil, &notifySocketHostAddr)
	if err != nil {
		logrus.Error(err)
		return
	}
	for {
		r, err := s.socket.Read(buf)
		if err != nil {
			break
		}
		var out bytes.Buffer
		for _, line := range bytes.Split(buf[0:r], []byte{'\n'}) {
			if bytes.HasPrefix(line, []byte("READY=")) {
				_, err = out.Write(line)
				if err != nil {
					return
				}

				_, err = out.Write([]byte{'\n'})
				if err != nil {
					return
				}

				_, err = client.Write(out.Bytes())
				if err != nil {
					return
				}

				// now we can inform butts to use pid1 as the pid to monitor
				if pid1 > 0 {
					newPid := fmt.Sprintf("MAINPID=%d\n", pid1)
					client.Write([]byte(newPid))
				}
				return
			}
		}
	}
}
