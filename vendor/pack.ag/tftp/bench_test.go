package tftp // import "pack.ag/tftp"

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkGet_random(b *testing.B) {
	// random1MB := getTestData(b, "1MB-random")
	// text := getTestData(b, "text")

	cases := []struct {
		name     string
		url      string
		response []byte
		opts     []ClientOpt
	}{
		{
			name:     "small data",
			url:      "tftp://#host#:#port#/file",
			response: []byte("the data"),
		},
		// {
		// 	name:     "small data-netascii",
		// 	url:      "tftp://#host#:#port#/file",
		// 	response: []byte("the data"),
		// 	opts:     []ClientOpt{ClientMode(ModeNetASCII)},
		// },
		// {
		// 	name:     "small-netascii",
		// 	url:      "tftp://#host#:#port#/file",
		// 	response: []byte("the\r\x00data with\r\nnewline"),
		// 	opts:     []ClientOpt{ClientMode(ModeNetASCII)},
		// },
		// {
		// 	name:     "text",
		// 	url:      "tftp://#host#:#port#/file",
		// 	response: text,
		// },
		// {
		// 	name:     "text-netascii-nix",
		// 	url:      "tftp://#host#:#port#/file",
		// 	response: text,
		// 	opts:     []ClientOpt{ClientMode(ModeNetASCII)},
		// },
		// {
		// 	name:     "text-netascii-windows",
		// 	url:      "tftp://#host#:#port#/file",
		// 	response: text,
		// 	opts:     []ClientOpt{ClientMode(ModeNetASCII)},
		// },
		// {
		// 	name:     "1MB",
		// 	url:      "tftp://#host#:#port#/file",
		// 	response: random1MB,
		// },
		// {
		// 	name:     "1MB, don't send size",
		// 	url:      "tftp://#host#:#port#/file",
		// 	response: random1MB,
		// },
		// {
		// 	name:     "1MB-blksize9000",
		// 	url:      "tftp://#host#:#port#/file",
		// 	response: random1MB,
		// 	opts:     []ClientOpt{ClientBlocksize(9000)},
		// },
		// {
		// 	name:     "1MB-window5",
		// 	url:      "tftp://#host#:#port#/file",
		// 	response: random1MB,
		// 	opts:     []ClientOpt{ClientWindowsize(5)},
		// },
	}

	for _, c := range cases {
		for _, singlePort := range []bool{true, false} {
			name := fmt.Sprintf("%s, single port mode: %t", c.name, singlePort)
			b.Run(name, func(b *testing.B) {
				ip, port, close := newTestServer(b, singlePort, func(w ReadRequest) {
					w.WriteSize(int64(len(c.response)))
					w.Write([]byte(c.response))
				}, nil)
				defer close()

				url := strings.Replace(c.url, "#host#", ip, 1)
				url = strings.Replace(url, "#port#", strconv.Itoa(port), 1)

				for i := 0; i < b.N; i++ {
					client, err := NewClient(c.opts...)
					if err != nil {
						b.Fatal(err)
					}

					file, err := client.Get(url)
					if err != nil {
						b.Fatal(err)
					}
					b.ResetTimer()

					_, err = ioutil.ReadAll(file)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}
