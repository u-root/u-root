// Copyright 2020-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package namespace

import (
	"io"
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		r io.Reader
	}
	_open := func(name string) func(t *testing.T) args {
		return func(t *testing.T) args {
			r, _ := os.Open(name)
			return args{r}
		}
	}

	tests := []struct {
		name string
		args func(t *testing.T) args

		want1      File
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name: "namespace",
			args: _open("testdata/namespace"),
			want1: File{
				cmd{
					syscall: MOUNT,
					flag:    REPL | AFTER | CACHE,
					args:    []string{"#s/boot", "/root", "$rootspec"},
				},
				cmd{
					syscall: BIND,
					flag:    REPL | AFTER,
					args:    []string{"$rootdir", "/"},
				},
				cmd{
					syscall: BIND,
					flag:    REPL | CREATE,
					args:    []string{"$rootdir/mnt", "/mnt"},
				},
				// kernel devices
				cmd{
					syscall: BIND,
					args:    []string{"#c", "/dev"},
				},
				cmd{
					syscall: BIND,
					args:    []string{"#d", "/fd"},
				},
				cmd{
					syscall: BIND,
					flag:    REPL | CREATE,
					args:    []string{"#e", "/env"},
				},
				cmd{
					syscall: BIND,
					args:    []string{"#p", "/proc"},
				},
				cmd{
					syscall: BIND,
					flag:    REPL | CREATE,
					args:    []string{"#s", "/srv"},
				},
				cmd{
					syscall: BIND,
					flag:    REPL | AFTER,
					args:    []string{"#¤", "/dev"},
				},
				cmd{
					syscall: BIND,
					flag:    REPL | AFTER,
					args:    []string{"#S", "/dev"},
				},
				cmd{
					syscall: BIND,
					flag:    REPL | BEFORE,
					args:    []string{"#k", "/dev"},
				},
				cmd{
					syscall: BIND,
					flag:    REPL | AFTER,
					args:    []string{"#κ", "/dev"},
				},
				cmd{
					syscall: BIND,
					flag:    REPL | AFTER,
					args:    []string{"#u", "/dev"},
				},
				cmd{
					syscall: BIND,
					flag:    REPL | BEFORE,
					args:    []string{"#P", "/dev"},
				},
				// mount points
				cmd{
					syscall: MOUNT,
					flag:    REPL | AFTER,
					args:    []string{"/srv/slashn", "/n"},
				},
				// authentication
				cmd{
					syscall: MOUNT,
					flag:    REPL | AFTER,
					args:    []string{"/srv/factotum", "/mnt"},
				},
				// standard bin
				cmd{
					syscall: BIND,
					flag:    REPL,
					args:    []string{"/$cputype/bin", "/bin"},
				},
				cmd{
					syscall: BIND,
					flag:    REPL | AFTER,
					args:    []string{"/rc/bin", "/bin"},
				},
				// internal networks
				cmd{
					syscall: BIND,
					flag:    REPL | AFTER,
					args:    []string{"#l", "/net"},
				},
				cmd{
					syscall: BIND,
					flag:    REPL | AFTER,
					args:    []string{"#I", "/net"},
				},
				cmd{
					syscall: MOUNT,
					flag:    REPL | AFTER,
					args:    []string{"/srv/cs", "/net"},
				},
				cmd{
					syscall: MOUNT,
					flag:    REPL | AFTER,
					args:    []string{"/srv/dns", "/net"},
				},
				cmd{
					syscall: MOUNT,
					flag:    REPL | AFTER,
					args:    []string{"/srv/net", "/net"},
				},
				cmd{
					syscall: MOUNT,
					flag:    REPL | BEFORE,
					args:    []string{"/srv/ssh", "/net"},
				},
				// usbd, mainly disks
				cmd{
					syscall: MOUNT,
					args:    []string{"/srv/usb", "/n/usb"},
				},
				cmd{
					syscall: MOUNT,
					flag:    REPL | AFTER,
					args:    []string{"/srv/usb", "/dev"},
				},
				cmd{
					syscall: BIND,
					flag:    REPL | CREATE,
					args:    []string{"/usr/$user/tmp", "/tmp"},
				},
				cmd{
					syscall: CHDIR,
					args:    []string{"/usr/$user"},
				},
				cmd{
					syscall: INCLUDE,
					args:    []string{"/lib/namespace.local"},
				},
				cmd{
					syscall: INCLUDE,
					args:    []string{"/lib/namespace.$sysname"},
				},
				cmd{
					syscall: INCLUDE,
					args:    []string{"/cfg/$sysname/namespace"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, err := Parse(tArgs.r)

			// assert.Equal(t, tt.want1, got1)
			if !reflect.DeepEqual(tt.want1, got1) {
				t.Errorf(`Parse(%v) = %v, _ , want %v`, tArgs.r, got1, tt.want1)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf(`Parse(%v) = _, %v, want not nil`, tArgs.r, err)
				}
				if tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			}
		})
	}
}

func printErr(err error, t *testing.T) { t.Log(err) }

func Test_parseLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1      cmd
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name: "mount",
			args: func(*testing.T) args { return args{"mount -aC #s/boot /root $rootspec"} },
			want1: cmd{
				syscall: MOUNT,
				flag:    REPL | AFTER | CACHE,
				args:    []string{"#s/boot", "/root", "$rootspec"},
			},
			wantErr:    false,
			inspectErr: printErr,
		},
		{
			name: "mount",
			args: func(*testing.T) args { return args{"mount -aC #s/boot /root"} },
			want1: cmd{
				syscall: MOUNT,
				flag:    REPL | AFTER | CACHE,
				args:    []string{"#s/boot", "/root"},
			},
			wantErr:    false,
			inspectErr: printErr,
		},
		{
			name: "mount",
			args: func(*testing.T) args { return args{"./mount -aC #s/boot /root"} },
			want1: cmd{
				syscall: MOUNT,
				flag:    REPL | AFTER | CACHE,
				args:    []string{"#s/boot", "/root"},
			},
			wantErr:    false,
			inspectErr: printErr,
		},
		{
			name: "bind",
			args: func(*testing.T) args { return args{"bind -a '#r' /dev"} },
			want1: cmd{
				syscall: BIND,
				args:    []string{"'#r'", "/dev"},
				flag:    REPL | AFTER,
			},
			wantErr: false,
		},
		{
			name: "include",
			args: func(*testing.T) args { return args{". /cfg/$sysname/namespace"} },
			want1: cmd{
				syscall: INCLUDE,
				args:    []string{"/cfg/$sysname/namespace"},
				flag:    REPL,
			},
			wantErr: false,
		},
		{
			name: "clear",
			args: func(*testing.T) args { return args{"clear"} },
			want1: cmd{
				syscall: RFORK,
				args:    []string{},
				flag:    REPL,
			},
			wantErr: false,
		},
		{
			name: "umount",
			args: func(*testing.T) args { return args{"unmount /dev"} },
			want1: cmd{
				syscall: UNMOUNT,
				args:    []string{"/dev"},
				flag:    REPL,
			},
			wantErr: false,
		},
		{
			name: "cd",
			args: func(*testing.T) args { return args{"cd /cfg"} },
			want1: cmd{
				syscall: CHDIR,
				args:    []string{"/cfg"},
				flag:    REPL,
			},
			wantErr: false,
		},
		{
			name: "import",
			args: func(*testing.T) args { return args{"import -a $host /srv /srv"} },
			want1: cmd{
				syscall: IMPORT,
				args:    []string{"$host", "/srv", "/srv"},
				flag:    AFTER,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, err := ParseLine(tArgs.line)

			if !reflect.DeepEqual(tt.want1, got1) {
				t.Errorf(`ParseLine(%v) = %v, _ , want %v`, tArgs.line, got1, tt.want1)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf(`ParseLine(%v) = _, %v, want not nil`, tArgs.line, err)
				}
				if tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			}
		})
	}
}
