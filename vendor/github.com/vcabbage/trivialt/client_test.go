// Copyright (C) 2016 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package trivialt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	defaultOpts := map[string]string{
		optTransferSize: "0",
	}

	cases := map[string]struct {
		opts []ClientOpt

		expectedError      error
		expectedOpts       map[string]string
		expectedMode       transferMode
		expectedRetransmit int
	}{
		"default": {
			expectedOpts:       defaultOpts,
			expectedMode:       ModeOctet,
			expectedRetransmit: 10,
		},
		"mode": {
			opts: []ClientOpt{ClientMode(ModeNetASCII)},

			expectedOpts:       defaultOpts,
			expectedMode:       ModeNetASCII,
			expectedRetransmit: 10,
		},
		"blksize": {
			opts: []ClientOpt{ClientBlocksize(42)},

			expectedOpts: map[string]string{
				optTransferSize: "0",
				optBlocksize:    "42",
			},
			expectedMode:       ModeOctet,
			expectedRetransmit: 10,
		},
		"timeout": {
			opts: []ClientOpt{ClientTimeout(24)},

			expectedOpts: map[string]string{
				optTransferSize: "0",
				optTimeout:      "24",
			},
			expectedMode:       ModeOctet,
			expectedRetransmit: 10,
		},
		"windowsize": {
			opts: []ClientOpt{ClientWindowsize(13)},

			expectedOpts: map[string]string{
				optTransferSize: "0",
				optWindowSize:   "13",
			},
			expectedMode:       ModeOctet,
			expectedRetransmit: 10,
		},
		"tsize enabled": {
			opts: []ClientOpt{ClientTransferSize(true)},

			expectedOpts: map[string]string{
				optTransferSize: "0",
			},
			expectedMode:       ModeOctet,
			expectedRetransmit: 10,
		},
		"tsize disabled": {
			opts: []ClientOpt{ClientTransferSize(false)},

			expectedOpts:       map[string]string{},
			expectedMode:       ModeOctet,
			expectedRetransmit: 10,
		},
		"retransmit": {
			opts: []ClientOpt{ClientRetransmit(13)},

			expectedOpts:       defaultOpts,
			expectedMode:       ModeOctet,
			expectedRetransmit: 13,
		},
		"two opts": {
			opts: []ClientOpt{
				ClientWindowsize(13),
				ClientTimeout(24),
			},

			expectedOpts: map[string]string{
				optTransferSize: "0",
				optWindowSize:   "13",
				optTimeout:      "24",
			},
			expectedMode:       ModeOctet,
			expectedRetransmit: 10,
		},
		"bad mode": {
			opts: []ClientOpt{
				ClientMode("fast"),
			},

			expectedError: ErrInvalidMode,
		},
		"blocksize too small": {
			opts: []ClientOpt{
				ClientBlocksize(7),
			},

			expectedError: ErrInvalidBlocksize,
		},
		"blocksize too large": {
			opts: []ClientOpt{
				ClientBlocksize(65465),
			},

			expectedError: ErrInvalidBlocksize,
		},
		"timeout too small": {
			opts: []ClientOpt{
				ClientTimeout(0),
			},

			expectedError: ErrInvalidTimeout,
		},
		"timeout too large": {
			opts: []ClientOpt{
				ClientTimeout(256),
			},

			expectedError: ErrInvalidTimeout,
		},
		"windowsize too small": {
			opts: []ClientOpt{
				ClientWindowsize(0),
			},

			expectedError: ErrInvalidWindowsize,
		},
		"windowsize too large": {
			opts: []ClientOpt{
				ClientWindowsize(65536),
			},

			expectedError: ErrInvalidWindowsize,
		},
		"retransmit negative": {
			opts: []ClientOpt{
				ClientRetransmit(-1),
			},

			expectedError: ErrInvalidRetransmit,
		},
	}

	for label, c := range cases {
		client, err := NewClient(c.opts...)

		// Error
		if err != c.expectedError {
			t.Errorf("%s: Expected %#v to be %#v", label, err, c.expectedError)
		}

		if err != nil {
			continue // Skip remaining test if error, avoid nil dereference
		}

		// Options
		if !reflect.DeepEqual(client.opts, c.expectedOpts) {
			t.Errorf("%s: Expected opts to be %#v, but they were %#v", label, c.expectedOpts, client.opts)
		}

		// Mode
		if client.mode != c.expectedMode {
			t.Errorf("%s: Expected mode to be %s, but it was %s", label, c.expectedMode, client.mode)
		}

		// Retransmit
		if client.retransmit != c.expectedRetransmit {
			t.Errorf("%s: Expected retransmit to be %d, but it was %d", label, c.expectedRetransmit, client.retransmit)
		}
	}
}

func TestClient_Get(t *testing.T) {
	t.Parallel()

	random1MB := getTestData(t, "1MB-random")
	text := getTestData(t, "text")
	textWindows := getTestData(t, "text-windows")
	randomUnder1MB := random1MB[:len(random1MB)-3] // not divisible by 512

	cases := map[string]struct {
		url             string
		response        []byte
		opts            []ClientOpt
		omitSize        bool
		sendServerError bool
		windowsOnly     bool
		nixOnly         bool

		expectedResponse []byte
		expectedSize     int64
		expectedError    string
	}{
		"small data": {
			url:      "tftp://#host#:#port#/file",
			response: []byte("the data"),

			expectedResponse: []byte("the data"),
			expectedSize:     8,
		},
		"small data-netascii": {
			url:      "tftp://#host#:#port#/file",
			response: []byte("the data"),
			opts:     []ClientOpt{ClientMode(ModeNetASCII)},

			expectedResponse: []byte("the data"),
			expectedSize:     8,
		},
		"small-netascii": {
			url:      "tftp://#host#:#port#/file",
			response: []byte("the\r\x00data with\r\nnewline"),
			opts:     []ClientOpt{ClientMode(ModeNetASCII)},
			nixOnly:  true,

			expectedResponse: []byte("the\rdata with\nnewline"),
			expectedSize:     23, // Decoded size is larger than received
		},
		"small-netascii-windows": {
			url:         "tftp://#host#:#port#/file",
			response:    []byte("the\r\x00data with\r\nnewline"),
			opts:        []ClientOpt{ClientMode(ModeNetASCII)},
			windowsOnly: true,

			expectedResponse: []byte("the\rdata with\r\nnewline"),
			expectedSize:     23, // Decoded size is larger than received
		},
		"small data, don't send size": {
			url:      "tftp://#host#:#port#/file",
			response: []byte("thedata"),
			omitSize: true,

			expectedResponse: []byte("thedata"),
			expectedSize:     0,
		},
		"text": {
			url:      "tftp://#host#:#port#/file",
			response: text,

			expectedResponse: text,
			expectedSize:     810880,
		},
		"text-netascii-nix": {
			url:      "tftp://#host#:#port#/file",
			response: text,
			opts:     []ClientOpt{ClientMode(ModeNetASCII)},
			nixOnly:  true,

			expectedResponse: text,
			expectedSize:     810880, // TODO: Disable tsize for netascii?
		},
		"text-netascii-windows": {
			url:         "tftp://#host#:#port#/file",
			response:    text,
			opts:        []ClientOpt{ClientMode(ModeNetASCII)},
			windowsOnly: true,

			expectedResponse: textWindows,
			expectedSize:     810880, // TODO: Disable tsize for netascii?
		},
		"1MB": {
			url:      "tftp://#host#:#port#/file",
			response: random1MB,

			expectedResponse: random1MB,
			expectedSize:     1048576,
		},
		"1MB, don't send size": {
			url:      "tftp://#host#:#port#/file",
			response: random1MB,
			omitSize: true,

			expectedResponse: random1MB,
			expectedSize:     0,
		},
		"1MB-blksize9000": {
			url:      "tftp://#host#:#port#/file",
			response: random1MB,
			opts:     []ClientOpt{ClientBlocksize(9000)},

			expectedResponse: random1MB,
			expectedSize:     1048576,
		},
		"1MB-window5": {
			url:      "tftp://#host#:#port#/file",
			response: random1MB,
			opts:     []ClientOpt{ClientWindowsize(5)},

			expectedResponse: random1MB,
			expectedSize:     1048576,
		},
		"1MB-timeout5": {
			url:      "tftp://#host#:#port#/file",
			response: random1MB,
			opts:     []ClientOpt{ClientTimeout(5)},

			expectedResponse: random1MB,
			expectedSize:     1048576,
		},
		"under-1MB": {
			url:      "tftp://#host#:#port#/file",
			response: randomUnder1MB,

			expectedResponse: randomUnder1MB,
			expectedSize:     1048573,
		},
		"under-1MB, don't send size": {
			url:      "tftp://#host#:#port#/file",
			response: randomUnder1MB,
			omitSize: true,

			expectedResponse: randomUnder1MB,
			expectedSize:     0,
		},
		"under-1MB-blksize9000": {
			url:      "tftp://#host#:#port#/file",
			response: randomUnder1MB,
			opts:     []ClientOpt{ClientBlocksize(9000)},

			expectedResponse: randomUnder1MB,
			expectedSize:     1048573,
		},
		"under-1MB-window5": {
			url:      "tftp://#host#:#port#/file",
			response: randomUnder1MB,
			opts:     []ClientOpt{ClientWindowsize(5)},

			expectedResponse: randomUnder1MB,
			expectedSize:     1048573,
		},
		"under-1MB-timeout5": {
			url:      "tftp://#host#:#port#/file",
			response: randomUnder1MB,
			opts:     []ClientOpt{ClientTimeout(5)},

			expectedResponse: randomUnder1MB,
			expectedSize:     1048573,
		},
		"localhost": {
			url:      "tftp://localhost:#port#/file",
			response: []byte("the data"),

			expectedResponse: []byte("the data"),
			expectedSize:     8,
		},
		"bad url": {
			url: "host:#host#:#port#/file",

			expectedError: "invalid host/IP",
		},
		"cannot connect": {
			url: "thishostdoesnotexist/file",

			expectedError: "[Nn]o such host",
		},
		"server error": {
			url:             "tftp://#host#:#port#/file",
			response:        []byte("the data"),
			sendServerError: true,

			expectedError: `remote error: ERROR\[Code: ACCESS_VIOLATION; Message: \"server error\"\]`,
		},
	}

	for label, c := range cases {
		for _, singlePort := range []bool{true, false} {
			label := fmt.Sprintf("%s, single port mode: %t", label, singlePort)

			if (c.windowsOnly && runtime.GOOS != "windows") || (c.nixOnly && runtime.GOOS == "windows") {
				t.Logf("skipping case %q marked windowsOnly:%t; nixOnly:%t; GOOS: %q", label, c.windowsOnly, c.nixOnly, runtime.GOOS)
				continue
			}

			var mu sync.Mutex

			ip, port, close := newTestServer(t, singlePort, func(w ReadRequest) {
				mu.Lock()
				defer mu.Unlock()
				if c.sendServerError {
					w.WriteError(ErrCodeAccessViolation, "server error")
					return
				}

				if !c.omitSize {
					w.WriteSize(int64(len(c.response)))
				}
				w.Write([]byte(c.response))
			}, nil)
			defer close()

			client, err := NewClient(c.opts...)
			if err != nil {
				t.Fatal(err)
			}

			url := strings.Replace(c.url, "#host#", ip, 1)
			url = strings.Replace(url, "#port#", strconv.Itoa(port), 1)

			file, err := client.Get(url)
			if err != nil {
				if match, _ := regexp.MatchString(c.expectedError, ErrorCause(err).Error()); !match {
					t.Errorf("%s: expected error %q, got %q", label, c.expectedError, ErrorCause(err).Error())
				}
				mu.Lock()
				mu.Unlock()
				continue
			}

			response, err := ioutil.ReadAll(file)
			mu.Lock()
			mu.Unlock()
			if err != nil {
				t.Fatal(label, err)
			}

			// Data
			if !reflect.DeepEqual(response, c.expectedResponse) {
				if len(response) > 1000 || len(c.expectedResponse) > 1000 {
					t.Errorf("%s: Response didn't match (over 1000 characters, omitting)", label)
				} else {
					t.Errorf("%s: Expected response to be %q, but it was %q", label, c.expectedResponse, response)
				}
			}

			// Size
			if i, _ := file.Size(); i != c.expectedSize {
				t.Errorf("%s: Expected size to be %d, but it was %d", label, c.expectedSize, i)
			}
		}
	}
}

func TestClient_Put(t *testing.T) {
	t.Parallel()

	random1MB := getTestData(t, "1MB-random")
	text := getTestData(t, "text")
	textWindows := getTestData(t, "text-windows")
	randomUnder1MB := random1MB[:len(random1MB)-3] // not divisible by 512

	cases := map[string]struct {
		url             string
		send            []byte
		opts            []ClientOpt
		omitSize        bool
		sendServerError bool
		windowsOnly     bool
		nixOnly         bool

		expectedData  []byte
		expectedSize  int64
		expectedError string
	}{
		"small data": {
			url:  "tftp://#host#:#port#/file",
			send: []byte("the data"),

			expectedData: []byte("the data"),
			expectedSize: 8,
		},
		"small data-netascii": {
			url:  "tftp://#host#:#port#/file",
			send: []byte("the data"),
			opts: []ClientOpt{ClientMode(ModeNetASCII)},

			expectedData: []byte("the data"),
			expectedSize: 8,
		},
		"small-netascii": {
			url:     "tftp://#host#:#port#/file",
			send:    []byte("the\r\x00data with\r\nnewline"),
			opts:    []ClientOpt{ClientMode(ModeNetASCII)},
			nixOnly: true,

			expectedData: []byte("the\rdata with\nnewline"),
			expectedSize: 23, // Decoded size is larger than received
		},
		"small-netascii-windows": {
			url:         "tftp://#host#:#port#/file",
			send:        []byte("the\r\x00data with\r\nnewline"),
			opts:        []ClientOpt{ClientMode(ModeNetASCII)},
			windowsOnly: true,

			expectedData: []byte("the\rdata with\r\nnewline"),
			expectedSize: 23, // Decoded size is larger than received
		},
		"small data, don't send size": {
			url:      "tftp://#host#:#port#/file",
			send:     []byte("thedata"),
			omitSize: true,

			expectedData: []byte("thedata"),
			expectedSize: 0,
		},
		"text": {
			url:  "tftp://#host#:#port#/file",
			send: text,

			expectedData: text,
			expectedSize: 810880,
		},
		"text-netascii-nix": {
			url:     "tftp://#host#:#port#/file",
			send:    text,
			opts:    []ClientOpt{ClientMode(ModeNetASCII)},
			nixOnly: true,

			expectedData: text,
			expectedSize: 810880, // TODO: Disable tsize for netascii?
		},
		"text-netascii-windows": {
			url:         "tftp://#host#:#port#/file",
			send:        text,
			opts:        []ClientOpt{ClientMode(ModeNetASCII)},
			windowsOnly: true,

			expectedData: textWindows,
			expectedSize: 810880, // TODO: Disable tsize for netascii?
		},
		"1MB": {
			url:  "tftp://#host#:#port#/file",
			send: random1MB,

			expectedData: random1MB,
			expectedSize: 1048576,
		},
		"1MB, don't send size": {
			url:      "tftp://#host#:#port#/file",
			send:     random1MB,
			omitSize: true,

			expectedData: random1MB,
			expectedSize: 0,
		},
		"1MB-blksize9000": {
			url:  "tftp://#host#:#port#/file",
			send: random1MB,
			opts: []ClientOpt{ClientBlocksize(9000)},

			expectedData: random1MB,
			expectedSize: 1048576,
		},
		"1MB-window2": {
			url:  "tftp://#host#:#port#/file",
			send: random1MB,
			opts: []ClientOpt{ClientWindowsize(2)},

			expectedData: random1MB,
			expectedSize: 1048576,
		},
		"1MB-timeout5": {
			url:  "tftp://#host#:#port#/file",
			send: random1MB,
			opts: []ClientOpt{ClientTimeout(5)},

			expectedData: random1MB,
			expectedSize: 1048576,
		},
		"under-1MB": {
			url:  "tftp://#host#:#port#/file",
			send: randomUnder1MB,

			expectedData: randomUnder1MB,
			expectedSize: 1048573,
		},
		"under-1MB, don't send size": {
			url:      "tftp://#host#:#port#/file",
			send:     randomUnder1MB,
			omitSize: true,

			expectedData: randomUnder1MB,
			expectedSize: 0,
		},
		"under-1MB-blksize9000": {
			url:  "tftp://#host#:#port#/file",
			send: randomUnder1MB,
			opts: []ClientOpt{ClientBlocksize(9000)},

			expectedData: randomUnder1MB,
			expectedSize: 1048573,
		},
		"under-1MB-window5": {
			url:  "tftp://#host#:#port#/file",
			send: randomUnder1MB,
			opts: []ClientOpt{ClientWindowsize(2)},

			expectedData: randomUnder1MB,
			expectedSize: 1048573,
		},
		"under-1MB-timeout5": {
			url:  "tftp://#host#:#port#/file",
			send: randomUnder1MB,
			opts: []ClientOpt{ClientTimeout(5)},

			expectedData: randomUnder1MB,
			expectedSize: 1048573,
		},
		"bad url": {
			url: "host:#host#:#port#/file",

			expectedError: "invalid host/IP",
		},
		"cannot connect": {
			url: "thishostdoesnotexist/file",

			expectedError: "[Nn]o such host",
		},
		"server error": {
			url:             "tftp://#host#:#port#/file",
			sendServerError: true,

			expectedError: `remote error: ERROR\[Code: ACCESS_VIOLATION; Message: \"server error\"\]`,
		},
	}

	for label, c := range cases {
		for _, singlePort := range []bool{true, false} {
			label := fmt.Sprintf("%s, single port mode: %t", label, singlePort)

			if (c.windowsOnly && runtime.GOOS != "windows") || (c.nixOnly && runtime.GOOS == "windows") {
				t.Logf("skipping case %q marked windowsOnly:%t; nixOnly:%t; GOOS: %q", label, c.windowsOnly, c.nixOnly, runtime.GOOS)
				continue
			}

			var wr WriteRequest
			var data []byte
			var mu sync.Mutex

			ip, port, close := newTestServer(t, singlePort, nil, func(w WriteRequest) {
				mu.Lock()
				defer mu.Unlock()
				if c.sendServerError {
					w.WriteError(ErrCodeAccessViolation, "server error")
					return
				}
				wr = w

				d, err := ioutil.ReadAll(w)
				if err != nil {
					t.Fatal(err)
				}
				data = d
			})
			defer close()

			client, err := NewClient(c.opts...)
			if err != nil {
				t.Fatal(err)
			}

			size := 0
			if !c.omitSize {
				size = len(c.send)
			}

			url := strings.Replace(c.url, "#host#", ip, 1)
			url = strings.Replace(url, "#port#", strconv.Itoa(port), 1)

			err = client.Put(url, bytes.NewReader(c.send), int64(size))
			mu.Lock()
			mu.Unlock()
			if err != nil {
				if match, _ := regexp.MatchString(c.expectedError, ErrorCause(err).Error()); !match {
					t.Errorf("%s: expected error %q, got %q", label, c.expectedError, ErrorCause(err).Error())
				}
				continue
			}

			// Data
			if !reflect.DeepEqual(data, c.expectedData) {
				if len(data) > 1000 || len(c.expectedData) > 1000 {
					t.Errorf("%s: Response didn't match (over 1000 characters, omitting)", label)
				} else {
					t.Errorf("%s: Expected response to be %q, but it was %q", label, c.expectedData, data)
				}
			}

			// Size
			if size, _ := wr.Size(); size != c.expectedSize {
				t.Errorf("%s: Expected size to be %d, but it was %d", label, c.expectedSize, size)
			}
		}
	}
}

func TestClient_parseURL(t *testing.T) {
	cases := map[string]struct {
		url string

		expectedHost  string
		expectedFile  string
		expectedError error
	}{
		"host and file": {
			url: "myhost/myfile",

			expectedHost: "myhost:69",
			expectedFile: "myfile",
		},
		"host, port, and file": {
			url: "myhost:8345/myfile",

			expectedHost: "myhost:8345",
			expectedFile: "myfile",
		},
		"scheme, host, port, and file": {
			url: "tftp://myhost:8345/myfile",

			expectedHost: "myhost:8345",
			expectedFile: "myfile",
		},
		"port and file": {
			url: ":8345/myfile",

			expectedError: ErrInvalidHostIP,
		},
		"file onle": {
			url: "/myfile",

			expectedError: ErrInvalidHostIP,
		},
		"? in url": {
			url: "host:8345/myfile?path",

			expectedHost: "host:8345",
			expectedFile: "myfile?path",
		},
		"# in url": {
			url: "host:8345/myfile#path",

			expectedHost: "host:8345",
			expectedFile: "myfile#path",
		},
		"no file": {
			url: "localhost:69/",

			expectedError: ErrInvalidFile,
		},
		"empty": {
			url: "",

			expectedError: ErrInvalidURL,
		},
		"host is numeric": {
			url: "12345:69/file",

			expectedError: ErrInvalidHostIP,
		},
		"port is not numeric": {
			url: "host:a/file",

			expectedError: ErrInvalidHostIP,
		},
		"colons in hostname": {
			url: "my:host:a/file",

			expectedError: ErrInvalidHostIP,
		},
	}

	for label, c := range cases {
		u, err := parseURL(c.url)

		// Error
		if err != c.expectedError {
			t.Errorf("%s: Expected error %v, got %v", label, c.expectedError, err)
		}

		if err != nil {
			continue
		}

		// Host
		if u.host != c.expectedHost {
			t.Errorf("%s: Expected host %q, got %q", label, c.expectedHost, u.host)
		}

		// File
		if u.file != c.expectedFile {
			t.Errorf("%s: Expected file %q, got %q", label, c.expectedFile, u.file)
		}
	}
}

func newTestServer(t *testing.T, singlePort bool, rh ReadHandlerFunc, wh WriteHandlerFunc) (string, int, func()) {
	s, err := NewServer("127.0.0.1:0", ServerSinglePort(singlePort))

	if err != nil {
		t.Fatalf("newTestServer: %v\n", err)
	}
	s.ReadHandler(rh)
	s.WriteHandler(wh)

	go s.ListenAndServe()

	closer := func() {
		s.Close()
	}

	// Wait for server to start
	for !s.Connected() {
		runtime.Gosched() // Prevents gettting stuck here
	}

	// Check for IPv6
	addr, _ := s.Addr()
	ip := addr.IP.String()
	if addr.IP.To4() == nil {
		ip = fmt.Sprintf("[%s]", addr.IP)
	}

	return ip, addr.Port, closer
}

func getTestData(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("getTestData(%q): %v", name, err)
	}

	return data
}
