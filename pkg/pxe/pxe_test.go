package pxe

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"reflect"
	"testing"
)

func TestProbeFiles(t *testing.T) {
	// Anyone got some ideas for other test cases?
	for _, tt := range []struct {
		mac   net.HardwareAddr
		ip    net.IP
		files []string
	}{
		{
			mac: []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
			ip:  []byte{192, 168, 0, 1},
			files: []string{
				"01-aa-bb-cc-dd-ee-ff",
				"C0A80001",
				"C0A8000",
				"C0A800",
				"C0A80",
				"C0A8",
				"C0A",
				"C0",
				"C",
				"default",
			},
		},
		{
			mac: []byte{0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd},
			ip:  []byte{192, 168, 2, 91},
			files: []string{
				"01-88-99-aa-bb-cc-dd",
				"C0A8025B",
				"C0A8025",
				"C0A802",
				"C0A80",
				"C0A8",
				"C0A",
				"C0",
				"C",
				"default",
			},
		},
	} {
		got := probeFiles(tt.mac, tt.ip)
		if !reflect.DeepEqual(got, tt.files) {
			t.Errorf("probeFiles(%s, %s) = %v, want %v", tt.mac, tt.ip, got, tt.files)
		}
	}
}

func TestAppendFile(t *testing.T) {
	content1 := "1111"
	content2 := "2222"
	content3 := "3333"
	content4 := "4444"

	type label struct {
		kernel    string
		kernelErr error
		initrd    string
		initrdErr error
		cmdline   string
	}
	type config struct {
		defaultEntry string
		labels       map[string]label
	}

	for i, tt := range []struct {
		desc          string
		configFileURI string
		schemeFunc    func() Schemes
		wd            *url.URL
		config        *Config
		want          config
		err           error
	}{
		{
			desc:          "all files exist, simple config with cmdline initrd",
			configFileURI: "default",
			schemeFunc: func() Schemes {
				s := make(Schemes)
				fs := NewMockScheme("tftp")
				fs.Add("1.2.3.4", "/foobar/pxelinux.0", "")
				conf := `default foo
				label foo
				kernel ./pxefiles/kernel
				append initrd=./pxefiles/initrd`
				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				fs.Add("1.2.3.4", "/foobar/pxefiles/kernel", content1)
				fs.Add("1.2.3.4", "/foobar/pxefiles/initrd", content2)
				s.Register(fs.scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: config{
				defaultEntry: "foo",
				labels: map[string]label{
					"foo": label{
						kernel:  content1,
						initrd:  content2,
						cmdline: "initrd=./pxefiles/initrd",
					},
				},
			},
		},
		{
			desc:          "all files exist, simple config with directive initrd",
			configFileURI: "default",
			schemeFunc: func() Schemes {
				s := make(Schemes)
				fs := NewMockScheme("tftp")
				fs.Add("1.2.3.4", "/foobar/pxelinux.0", "")
				conf := `default foo
				label foo
				kernel ./pxefiles/kernel
				initrd ./pxefiles/initrd
				append foo=bar`
				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				fs.Add("1.2.3.4", "/foobar/pxefiles/kernel", content1)
				fs.Add("1.2.3.4", "/foobar/pxefiles/initrd", content2)
				s.Register(fs.scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: config{
				defaultEntry: "foo",
				labels: map[string]label{
					"foo": label{
						kernel:  content1,
						initrd:  content2,
						cmdline: "foo=bar",
					},
				},
			},
		},
		{
			desc:          "all files exist, simple config, no initrd",
			configFileURI: "default",
			schemeFunc: func() Schemes {
				s := make(Schemes)
				fs := NewMockScheme("tftp")
				fs.Add("1.2.3.4", "/foobar/pxelinux.0", "")
				conf := `default foo
				label foo
				kernel ./pxefiles/kernel`
				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				fs.Add("1.2.3.4", "/foobar/pxefiles/kernel", content1)
				s.Register(fs.scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: config{
				defaultEntry: "foo",
				labels: map[string]label{
					"foo": label{
						kernel:  content1,
						initrd:  "",
						cmdline: "",
					},
				},
			},
		},
		{
			desc:          "kernel does not exist, simple config",
			configFileURI: "default",
			schemeFunc: func() Schemes {
				s := make(Schemes)
				fs := NewMockScheme("tftp")
				fs.Add("1.2.3.4", "/foobar/pxelinux.0", "")
				conf := `default foo
				label foo
				kernel ./pxefiles/kernel`
				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				s.Register(fs.scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: config{
				defaultEntry: "foo",
				labels: map[string]label{
					"foo": label{
						kernelErr: errNoSuchFile,
						initrd:    "",
						cmdline:   "",
					},
				},
			},
		},
		{
			desc:          "config file does not exist",
			configFileURI: "default",
			schemeFunc: func() Schemes {
				s := make(Schemes)
				fs := NewMockScheme("tftp")
				s.Register(fs.scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			err: ErrConfigNotFound,
		},
		{
			desc:          "empty config",
			configFileURI: "default",
			schemeFunc: func() Schemes {
				s := make(Schemes)
				fs := NewMockScheme("tftp")
				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", "")
				s.Register(fs.scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: config{
				defaultEntry: "",
			},
		},
		{
			desc:          "valid config with two labels",
			configFileURI: "default",
			schemeFunc: func() Schemes {
				s := make(Schemes)
				fs := NewMockScheme("tftp")
				fs.Add("1.2.3.4", "/foobar/pxelinux.0", "")
				conf := `default foo

				label foo
				kernel ./pxefiles/fookernel
				append earlyprintk=ttyS0 printk=ttyS0

				label bar
				kernel ./pxefiles/barkernel
				append console=ttyS0`
				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				fs.Add("1.2.3.4", "/foobar/pxefiles/fookernel", content1)
				fs.Add("1.2.3.4", "/foobar/pxefiles/barkernel", content2)
				s.Register(fs.scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: config{
				defaultEntry: "foo",
				labels: map[string]label{
					"foo": label{
						kernel:  content1,
						cmdline: "earlyprintk=ttyS0 printk=ttyS0",
					},
					"bar": label{
						kernel:  content2,
						cmdline: "console=ttyS0",
					},
				},
			},
		},
		{
			desc:          "valid config with global APPEND directive",
			configFileURI: "default",
			schemeFunc: func() Schemes {
				s := make(Schemes)
				fs := NewMockScheme("tftp")
				fs.Add("1.2.3.4", "/foobar/pxelinux.0", "")
				conf := `default foo
				append foo=bar

				label foo
				kernel ./pxefiles/fookernel
				append earlyprintk=ttyS0 printk=ttyS0

				label bar
				kernel ./pxefiles/barkernel

				label baz
				kernel ./pxefiles/barkernel
				append -`
				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				fs.Add("1.2.3.4", "/foobar/pxefiles/fookernel", content1)
				fs.Add("1.2.3.4", "/foobar/pxefiles/barkernel", content2)
				s.Register(fs.scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: config{
				defaultEntry: "foo",
				labels: map[string]label{
					"foo": label{
						kernel: content1,
						// Does not contain global APPEND.
						cmdline: "earlyprintk=ttyS0 printk=ttyS0",
					},
					"bar": label{
						kernel: content2,
						// Contains only global APPEND.
						cmdline: "foo=bar",
					},
					"baz": label{
						kernel: content2,
						// "APPEND -" means ignore global APPEND.
						cmdline: "",
					},
				},
			},
		},
		{
			desc:          "valid config with global APPEND with initrd",
			configFileURI: "default",
			schemeFunc: func() Schemes {
				s := make(Schemes)
				fs := NewMockScheme("tftp")
				fs.Add("1.2.3.4", "/foobar/pxelinux.0", "")
				conf := `default mcnulty
				append initrd=./pxefiles/normal_person

				label mcnulty
				kernel ./pxefiles/copkernel
				append earlyprintk=ttyS0 printk=ttyS0

				label lester
				kernel ./pxefiles/copkernel

				label omar
				kernel ./pxefiles/drugkernel
				append -
				
				label stringer
				kernel ./pxefiles/drugkernel
				initrd ./pxefiles/criminal
				`
				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				fs.Add("1.2.3.4", "/foobar/pxefiles/copkernel", content1)
				fs.Add("1.2.3.4", "/foobar/pxefiles/drugkernel", content2)
				fs.Add("1.2.3.4", "/foobar/pxefiles/normal_person", content3)
				fs.Add("1.2.3.4", "/foobar/pxefiles/criminal", content4)
				s.Register(fs.scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: config{
				defaultEntry: "mcnulty",
				labels: map[string]label{
					"mcnulty": label{
						kernel: content1,
						// Does not contain global APPEND.
						cmdline: "earlyprintk=ttyS0 printk=ttyS0",
					},
					"lester": label{
						kernel: content1,
						initrd: content3,
						// Contains only global APPEND.
						cmdline: "initrd=./pxefiles/normal_person",
					},
					"omar": label{
						kernel: content2,
						// "APPEND -" means ignore global APPEND.
						cmdline: "",
					},
					"stringer": label{
						kernel: content2,
						// See TODO in pxe.go initrd handling.
						initrd:  content4,
						cmdline: "initrd=./pxefiles/normal_person",
					},
				},
			},
		},
		{
			desc:          "default label does not exist",
			configFileURI: "default",
			schemeFunc: func() Schemes {
				s := make(Schemes)
				fs := NewMockScheme("tftp")
				conf := `default avon`

				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				s.Register(fs.scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			err: ErrDefaultEntryNotFound,
			want: config{
				defaultEntry: "avon",
			},
		},
		{
			desc:          "multi-scheme valid config",
			configFileURI: "default",
			schemeFunc: func() Schemes {
				conf := `default sheeeit

				label sheeeit
				kernel ./pxefiles/kernel
				initrd http://someplace.com/someinitrd.gz`

				tftp := NewMockScheme("tftp")
				tftp.Add("1.2.3.4", "/foobar/pxelinux.0", "")
				tftp.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				tftp.Add("1.2.3.4", "/foobar/pxefiles/kernel", content2)

				http := NewMockScheme("http")
				http.Add("someplace.com", "/someinitrd.gz", content3)

				s := make(Schemes)
				s.Register(tftp.scheme, tftp)
				s.Register(http.scheme, http)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: config{
				defaultEntry: "sheeeit",
				labels: map[string]label{
					"sheeeit": label{
						kernel: content2,
						initrd: content3,
					},
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("Test [%02d] %s", i, tt.desc), func(t *testing.T) {
			s := tt.schemeFunc()
			c := NewConfig(tt.wd)
			c.schemes = s

			if err := c.AppendFile(tt.configFileURI); err != tt.err {
				t.Errorf("AppendFile() got %v, want %v", err, tt.err)
			} else if err != nil {
				return
			}

			if got, want := c.DefaultEntry, tt.want.defaultEntry; got != want {
				t.Errorf("DefaultEntry got %v, want %v", got, want)
			}

			for labelName, want := range tt.want.labels {
				t.Run(fmt.Sprintf("label %s", labelName), func(t *testing.T) {
					label, ok := c.Entries[labelName]
					if !ok {
						t.Errorf("Config label %v does not exist", labelName)
						return
					}

					// Same kernel?
					if label.Kernel == nil && (len(want.kernel) > 0 || want.kernelErr != nil) {
						t.Errorf("want kernel, got none")
					}
					if label.Kernel != nil {
						k, err := ioutil.ReadAll(label.Kernel)
						if err != want.kernelErr {
							t.Errorf("could not read kernel of label %q: %v, want %v", labelName, err, want.kernelErr)
						}
						if got, want := string(k), want.kernel; got != want {
							t.Errorf("got kernel %s, want %s", got, want)
						}
					}

					// Same initrd?
					if label.Initrd == nil && (len(want.initrd) > 0 || want.initrdErr != nil) {
						t.Errorf("want initrd, got none")
					}
					if label.Initrd != nil {
						i, err := ioutil.ReadAll(label.Initrd)
						if err != want.initrdErr {
							t.Errorf("could not read initrd of label %q: %v, want %v", labelName, err, want.initrdErr)
						}
						if got, want := string(i), want.initrd; got != want {
							t.Errorf("got initrd %s, want %s", got, want)
						}
					}

					// Same cmdline?
					if got, want := label.Cmdline, want.cmdline; got != want {
						t.Errorf("got cmdline %s, want %s", got, want)
					}
				})
			}

			// Check that the parser didn't make up labels.
			for labelName := range c.Entries {
				if _, ok := tt.want.labels[labelName]; !ok {
					t.Errorf("config has extra label %s, but not wanted", labelName)
				}
			}
		})
	}
}
