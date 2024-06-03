// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"os"
)

type PutCmd struct {
	localfiles []string
	remotefile string
	remotedir  string
}

func executePut(client ClientIf, host, port string, files []string) error {
	ret := &PutCmd{}
	switch len(files) {
	case 1:
		// Put file to server
		ret.localfiles = append(ret.localfiles, files...)
	case 2:
		// files[0] == localfile
		ret.localfiles = append(ret.localfiles, files[0])
		// files[1] == remotefile
		ret.remotefile = files[1]
	default:
		// files[:len(files)-2] == localfiles,
		ret.localfiles = append(ret.localfiles, files[:len(files)-2]...)
		// files[len(files)-1] == remote-directory
		ret.remotedir = files[len(files)-1]
	}

	for _, file := range ret.localfiles {
		url := constructURL(host, port, "", file)

		if len(ret.localfiles) == 1 && ret.remotefile != "" {
			url = constructURL(host, port, "", ret.remotefile)
		} else if len(ret.localfiles) > 1 {
			url = constructURL(host, port, ret.remotedir, file)
		}

		locFile, err := os.Open(file)
		if err != nil {
			return err
		}

		fs, err := locFile.Stat()
		if err != nil {
			return err
		}
		if err := client.Put(url, locFile, fs.Size()); err != nil {
			return err
		}
	}

	return nil
}

type GetCmd struct {
	remotefiles []string
	localfile   string
}

var errSizeNoMatch = errors.New("data size of read and write mismatch")

func executeGet(client ClientIf, host, port string, files []string) error {
	ret := &GetCmd{}
	switch len(files) {
	case 1:
		// files[0] == remotefile
		ret.remotefiles = append(ret.remotefiles, files[0])
	case 2:
		// files[0] == remotefile
		ret.remotefiles = append(ret.remotefiles, files[0])
		// files[1] == localfile
		ret.localfile = files[1]
	default:
		// files... == remotefiles
		ret.remotefiles = append(ret.remotefiles, files...)
	}

	for _, file := range ret.remotefiles {
		resp, err := client.Get(constructURL(host, port, "", file))
		if err != nil {
			return err
		}

		localfile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0o666)
		if err != nil {
			return nil
		}
		defer localfile.Close()

		if ret.localfile != "" && len(ret.remotefiles) == 1 {
			localfile, err = os.OpenFile(ret.localfile, os.O_CREATE|os.O_WRONLY, 0o666)
			if err != nil {
				return err
			}
		}

		datalen, err := resp.Size()
		if err != nil {
			return err
		}

		data := make([]byte, datalen)
		nR, err := resp.Read(data)
		if err != nil {
			return err
		}

		nW, err := localfile.Write(data)
		if err != nil {
			return err
		}

		if nR != nW {
			return errSizeNoMatch
		}
	}

	return nil
}
