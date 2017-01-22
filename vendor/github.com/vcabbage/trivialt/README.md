# **trivialt**

[![Go Report Card](https://goreportcard.com/badge/vcabbage/trivialt)](https://goreportcard.com/report/vcabbage/trivialt)
[![Coverage Status](https://coveralls.io/repos/github/vcabbage/trivialt/badge.svg?branch=master)](https://coveralls.io/github/vcabbage/trivialt?branch=master)
[![Build Status](https://travis-ci.org/vcabbage/trivialt.svg?branch=master)](https://travis-ci.org/vcabbage/trivialt)
[![Build status](https://ci.appveyor.com/api/projects/status/0sxw1t6jjoe4yc9p/branch/master?svg=true)](https://ci.appveyor.com/project/vcabbage/trivialt/branch/master)
[![GoDoc](https://godoc.org/github.com/vcabbage/trivialt?status.svg)](http://godoc.org/github.com/vcabbage/trivialt)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/vcabbage/trivialt/master/LICENSE)


trivialt is a cross-platform, concurrent TFTP server and client. It can be used as a standalone executable or included in a Go project as a library.


### Standards Implemented

- [X] Binary Transfer ([RFC 1350](https://tools.ietf.org/html/rfc1350))
- [X] Netascii Transfer ([RFC 1350](https://tools.ietf.org/html/rfc1350))
- [X] Option Extension ([RFC 2347](https://tools.ietf.org/html/rfc2347))
- [X] Blocksize Option ([RFC 2348](https://tools.ietf.org/html/rfc2348))
- [X] Timeout Interval Option ([RFC 2349](https://tools.ietf.org/html/rfc2349))
- [X] Transfer Size Option ([RFC 2349](https://tools.ietf.org/html/rfc2349))
- [X] Windowsize Option ([RFC 7440](https://tools.ietf.org/html/rfc7440))

### Unique Features

- __Single Port Mode__

    TL;DR: It allows TFTP to work through firewalls.

    A standard TFTP server implementation receives requests on port 69 and allocates a new high port (over 1024) dedicated to that request.
    In single port mode, trivialt receives and responds to requests on the same port. If trivialt is started on port 69, all communication will
    be done on port 69.
    
    The primary use case of this feature is to play nicely with firewalls. Most firewalls will prevent the typical case where the server responds
    back on a random port because they have no way of knowing that it is in response to a request that went out on port 69. In single port mode,
    the firewall will see a request go out to a server on port 69 and that server respond back on the same port, which most firewalls will allow.
    
    Of course if the firewall in question is configured to block TFTP connections, this setting won't help you.
    
    Enable single port mode with the `--single-port` flag. This is currently marked experimental as is diverges from the TFTP standard.

## Installation

If you have the Go toolchain installed you can simply `go get` the packages. This will download the source into your `$GOPATH` and install the binary to `$GOPATH/bin/trivialt`.

``` bash
go get -u github.com/vcabbage/trivialt/...
```

Pre-built binaries can be downloaded from the [release page](https://github.com/vcabbage/trivialt/releases).

## Command Usage

Running as a server:
```
# trivialt serve --help
NAME:
   trivialt serve - Serve files from the filesystem.

USAGE:
   trivialt serve [bind address] [root directory]

DESCRIPTION:
   Serves files from the local file systemd.

   Bind address is in form "ip:port". Omitting the IP will listen on all interfaces.
   If not specified the server will listen on all interfaces, port 69.app

   Files will be served from root directory. If omitted files will be served from
   the current directory.

OPTIONS:
   --writeable, -w	    Enable file upload.
   --single-port, --sp	Enable single port mode. [Experimental]
```

```
# trivialt serve :6900 /tftproot --writable
Starting TFTP Server on ":6900", serving "/tftproot"
Read Request from 127.0.0.1:61877 for "ubuntu-16.04-server-amd64.iso"
Write Request from 127.0.0.1:51205 for "ubuntu-16.04-server-amd64.iso"

```

Downloading a file:
```
# trivialt get --help
NAME:
   trivialt get - Download file from a server.

USAGE:
   trivialt get [command options] [server:port] [file]

OPTIONS:
   --blksize, -b "512"      Number of data bytes to send per-packet.
   --windowsize, -w "1"     Number of packets to send before requiring an acknowledgement.
   --timeout, -t "10"       Number of seconds to wait before terminating a stalled connection.
   --tsize                  Enable the transfer size option. (default)
   --retransmit, -r "10"    Maximum number of back-to-back lost packets before terminating the connection.
   --netascii               Enable netascii transfer mode.
   --binary, --octet, -i    Enable binary transfer mode. (default)
   --quiet, -q              Don't display progress.
   --output, -o             Sets the output location to write the file. If not specified the
                            file will be written in the current directory.
                            Specifying "-" will write the file to stdout. ("-" implies "--quiet")
```

```
# trivialt get localhost:6900 ubuntu-16.04-server-amd64.iso
ubuntu-16.04-server-amd64.iso:
 655.00 MB / 655.00 MB [=====================================================] 100.00% 16.76 MB/s39s
```

Uploading a file:
```
# trivialt get --help
NAME:
   trivialt get - Download file from a server.

USAGE:
   trivialt get [command options] [server:port] [file]

OPTIONS:
   --blksize, -b "512"      Number of data bytes to send per-packet.
   --windowsize, -w "1"     Number of packets to send before requiring an acknowledgement.
   --timeout, -t "10"       Number of seconds to wait before terminating a stalled connection.
   --tsize                  Enable the transfer size option. (default)
   --retransmit, -r "10"	Maximum number of back-to-back lost packets before terminating the connection.
   --netascii               Enable netascii transfer mode.
   --binary, --octet, -i	Enable binary transfer mode. (default)
   --quiet, -q              Don't display progress.
   --output, -o             Sets the output location to write the file. If not specified the
                            file will be written in the current directory.
                            Specifying "-" will write the file to stdout. ("-" implies "--quiet")
```

```
# trivialt put localhost:6900 ubuntu-16.04-server-amd64.iso --blksize 1468 --windowsize 16
ubuntu-16.04-server-amd64.iso:
 655.00 MB / 655.00 MB [=====================================================] 100.00% 178.41 MB/s3s
```

## API

trivialt's API was inspired by Go's well-known net/http API. If you can write a net/http handler or middleware, you should have no problem doing the same for trivialt.

### Configuration Functions

One area that is noticeably different from net/http is the configuration of clients and servers. trivialt uses "configuration functions" rather than the direct modification of the
Client/Server struct or a configuration struct passed into the factory functions.

A few explanations of this pattern:
* [Self-referential functions and the design of options](http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html) by Rob Pike
* [Functional options for friendly APIs](https://www.youtube.com/watch?v=24lFtGHWxAQ) by Dave Cheney [video]

If this sounds complicated, don't worry, the public API is quiet simple. The `NewClient` and `NewServer` functions take zero or more configuration functions.

Want all defaults? Don't pass anything.

Want a Client configured for blocksize 9000 and windowsize 16? Pass in `ClientBlocksize(9000)` and `ClientWindowsize(16)`.

``` go
// Default Client
trivialt.NewClient()

// Client with blocksize 9000, windowsize 16
trivialt.NewClient(trivialt.ClientBlocksize(9000), trivialt.ClientWindowsize(16))

// Configuring with a slice of options
opts := []trivialt.ClientOpt{
    trivialt.ClientMode(trivialt.ModeOctet),
    trivialt.ClientBlocksize(9000),
    trivialt.ClientWindowsize(16),
    trivialt.ClientTimeout(1),
    trivialt.ClientTransferSize(true),
    trivialt.ClientRetransmit(3),
}

trivialt.NewClient(opts...)
```

### Examples

#### Read File From Server, Print to stdout

``` go
client := trivialt.NewClient()
resp, err := client.Get("myftp.local/myfile")
if err != nil {
    log.Fatalln(err)
}

err := io.Copy(os.Stdout, resp)
if err != nil {
    log.Fatalln(err)
}
```

#### Write File to Server

``` go

file, err := os.Open("myfile")
if err != nil {
    log.Fatalln(err)
}
defer file.Close()

// Get the file info se we can send size (not required)
fileInfo, err := file.Stat()
if err != nil {
    log.Println("error getting file size:", err)
}

client := trivialt.NewClient()
err := client.Put("myftp.local/myfile", file, fileInfo.Size())
if err != nil {
    log.Fatalln(err)
}
```


#### HTTP Proxy

This rather contrived example proxies an incoming GET request to GitHub's public API. A more realistic use case might be proxying to PXE boot files on an HTTP server.

``` go
const baseURL = "https://api.github.com/"

func proxyTFTP(w trivialt.ReadRequest) {
	// Append the requested path to the baseURL
	url := baseURL + w.Name()

	// Send the HTTP request
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		// This could send more specific errors, but here we'read
		// choosing to simply send "file not found"" with the error
		// message from the HTTP client back to the TFTP client.
		w.WriteError(trivialt.ErrCodeFileNotFound, err.Error())
		return
	}
	defer resp.Body.Close()

	// Copy the body of the response to the TFTP client.
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println(err)
	}
}
```


This function doesn't itself implement the required `ReadHandler` interface, but we can make it a `ReadHandler` with the `ReadHandlerFunc` adapter (much like `http.HandlerFunc`).

``` go
readHandler := trivialt.ReadHandlerFunc(proxyTFTP)

server.ReadHandler(readHandler)

server.ListenAndServe()
```

```
# trivialt get localhost:6900 repos/golang/go -o - | jq
{
  "id": 23096959,
  "name": "go",
  "full_name": "golang/go",
  ...
}
```

Full example in [examples/httpproxy/httpproxy.go](https://github.com/vcabbage/trivialt/blob/master/examples/httpproxy/httpproxy.go).

#### Save Files to Database

Here `tftpDB` implements the `WriteHandler` interface directly.

``` go
// tftpDB embeds a *sql.DB and implements the trivialt.ReadHandler interface.
type tftpDB struct {
	*sql.DB
}

func (db *tftpDB) ReceiveTFTP(w trivialt.WriteRequest) {
	// Read the data from the client into memory
	data, err := ioutil.ReadAll(w)
	if err != nil {
		log.Println(err)
		return
	}

	// Insert the IP address of the client and the data into the database
	res, err := db.Exec("INSERT INTO tftplogs (ip, log) VALUES (?, ?)", w.Addr().IP.String(), string(data))
	if err != nil {
		log.Println(err)
		return
	}

	// Log a message with the details
	id, _ := res.LastInsertId()
	log.Printf("Inserted %d bytes of data from %s. (ID=%d)", len(data), w.Addr().IP, id)
}
```

```
# go run examples/database/database.go
2016/04/30 11:20:27 Inserted 32 bytes of data from 127.0.0.1. (ID=13)
```

Full example including checking the size before accepting the request in [examples/database/database.go](https://github.com/vcabbage/trivialt/blob/master/examples/database/database.go).
