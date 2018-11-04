// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/u-root/dhcp4/dhcp4server"
	"pack.ag/tftp"
)

var (
	selfIP    = flag.String("ip", "192.168.0.1", "IP of self")
	subnet    = flag.String("subnet", "192.168.1.0/24", "CIDR of network to assign to clients")
	directory = flag.String("dir", "", "Directory to serve")
)

func serve(w tftp.ReadRequest) {
	path := filepath.Join(*directory, filepath.Clean(w.Name()))

	file, err := os.Open(path)
	if err != nil {
		w.WriteError(tftp.ErrCodeFileNotFound, fmt.Sprintf("File %q does not exist", w.Name()))
		return
	}
	defer file.Close()

	finfo, _ := file.Stat()
	w.WriteSize(finfo.Size())
	if _, err = io.Copy(w, file); err != nil {
		log.Println(err)
	}
}

func main() {
	flag.Parse()

	var wg sync.WaitGroup
	if len(*directory) != 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			server, err := tftp.NewServer(":69")
			if err != nil {
				log.Fatalf("Could not start TFTP server: %v", err)
			}

			log.Println("starting file server")
			server.ReadHandler(tftp.ReadHandlerFunc(serve))
			log.Fatal(server.ListenAndServe())
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, sn, err := net.ParseCIDR(*subnet)
		if err != nil {
			log.Fatal(err)
		}

		l, err := net.ListenPacket("udp4", ":67")
		if err != nil {
			log.Fatal(err)
		}
		defer l.Close()

		self := net.ParseIP(*selfIP)
		log.Printf("Self IP: %v", self)
		s := dhcp4server.New(self, sn, "", "pxelinux.0")

		log.Fatal(s.Serve(log.New(os.Stdout, "", log.LstdFlags), l))
	}()

	wg.Wait()
}
