/*
 * This file is part of the harvey project.
 *
 * Copyright (C) 2015 Henry Donnay
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

package main

import (
	"archive/tar"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"flag"
	"hash"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	ignore          = "*\n!.gitignore\n"
	dirPermissions   = 0755
)

type V struct {
	Upstream     string
	Digest       map[string]string
	Compress     string
	RemovePrefix bool
	Exclude      []string
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	justCheck := flag.Bool("check", false, "verify the code in upstream/")

	flag.Parse()

	if(*justCheck && !repositoryIsClean()){
		log.Fatal("cannot verify upstream/ files: working directory not clean")
	}

	f, err := ioutil.ReadFile("vendor.json")
	if err != nil {
		log.Fatal(err)
	}

	vendor := &V{}
	if err := json.Unmarshal(f, vendor); err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat("upstream"); err == nil {
		log.Println("recreating upstream")
		if(*justCheck){
			run("rm", "-r", "-f", "upstream")
		} else {
			run("git", "rm", "-r", "-f", "upstream")
		}
	} else {
		if(*justCheck){
			log.Fatalf("Cannot verify upstream/ as it does not exists.")
		}
		os.MkdirAll("patch", dirPermissions)
		os.MkdirAll("harvey", dirPermissions)
		ig, err := os.Create(path.Join("harvey", ".gitignore"))
		if err != nil {
			log.Fatal(err)
		}
		defer ig.Close()
		if _, err := ig.WriteString(ignore); err != nil {
			log.Fatal(err)
		}
		run("git", "add", ig.Name())
	}

	if err := do(vendor, *justCheck); err != nil {
		log.Fatal(err)
	}

	if(*justCheck){
		if(repositoryIsClean()){
			log.Printf("the files in upstream/ matches those in "+path.Base(vendor.Upstream))
		} else {
			log.Fatalf("the files in upstream/ does not match those in "+path.Base(vendor.Upstream))
		}
	} else {
		run("git", "add", "vendor.json")
		run("git", "commit", "-s", "-m", "vendor: pull in "+path.Base(vendor.Upstream))
	}
}

func repositoryIsClean() bool {
	out, err := exec.Command("git", "status", "--porcelain").Output()
    if err != nil {
        log.Fatal(err)
    }
	if(len(out) > 0){
		log.Println("git status --porcelain");
		log.Println(string(out))
		return false
	}
	return true
}

func do(v *V, justCheck bool) error {
	name := fetch(v)
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer os.Remove(name)

	var unZ io.Reader
	switch v.Compress {
	case "gzip":
		unZ, err = gzip.NewReader(f)
		if err != nil {
			return err
		}
	case "bzip2":
		unZ = bzip2.NewReader(f)
	default:
		unZ = f
	}

	ar := tar.NewReader(unZ)
	h, err := ar.Next()
untar:
	for ; err == nil; h, err = ar.Next() {
		n := h.Name
		if v.RemovePrefix {
			n = strings.SplitN(n, "/", 2)[1]
		}
		for _, ex := range v.Exclude {
			if strings.HasPrefix(n, ex) {
				continue untar
			}
		}
		n = path.Join("upstream", n)
		if h.FileInfo().IsDir() {
			os.MkdirAll(n, dirPermissions)
			continue
		}
		os.MkdirAll(path.Dir(n), dirPermissions)
		out, err := os.Create(n)
		if err != nil {
			log.Println(err)
			continue
		}

		if n, err := io.Copy(out, ar); n != h.Size || err != nil {
			return err
		}
		out.Close()
	}
	if err != io.EOF {
		return err
	}

	if(justCheck){
		return nil;
	}
	return run("git", "add", "upstream")
}

type match struct {
	hash.Hash
	Good []byte
	Name string
}

func (m match) OK() bool {
	return bytes.Equal(m.Good, m.Hash.Sum(nil))
}

func fetch(v *V) string {
	if len(v.Digest) == 0 {
		log.Fatal("no checksums specifed")
	}

	f, err := ioutil.TempFile("", "cmdVendor")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	req, err := http.NewRequest("GET", v.Upstream, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy:              http.ProxyFromEnvironment,
			DisableCompression: true,
		},
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var digests []match
	for k, v := range v.Digest {
		g, err := hex.DecodeString(v)
		if err != nil {
			log.Fatal(err)
		}
		switch k {
		case "sha1":
			digests = append(digests, match{sha1.New(), g, k})
		case "sha224":
			digests = append(digests, match{sha256.New224(), g, k})
		case "sha256":
			digests = append(digests, match{sha256.New(), g, k})
		case "sha384":
			digests = append(digests, match{sha512.New384(), g, k})
		case "sha512":
			digests = append(digests, match{sha512.New(), g, k})
		}
	}
	ws := make([]io.Writer, len(digests))
	for i := range digests {
		ws[i] = digests[i]
	}
	w := io.MultiWriter(ws...)

	if _, err := io.Copy(f, io.TeeReader(res.Body, w)); err != nil {
		log.Fatal(err)
	}
	for _, h := range digests {
		if !h.OK() {
			log.Fatalf("mismatched %q hash\n\tWanted %x\n\tGot %x\n", h.Name, h.Good, h.Hash.Sum(nil))
		}
	}
	return f.Name()
}

func run(exe string, arg ...string) error {
	cmd := exec.Command(exe, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
