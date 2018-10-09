# **pack.ag/tftp**

[![Go Report Card](https://goreportcard.com/badge/vcabbage/go-tftp)](https://goreportcard.com/report/vcabbage/go-tftp)
[![Coverage Status](https://coveralls.io/repos/github/vcabbage/go-tftp/badge.svg?branch=master)](https://coveralls.io/github/vcabbage/go-tftp?branch=master)
[![Build Status](https://travis-ci.org/vcabbage/go-tftp.svg?branch=master)](https://travis-ci.org/vcabbage/go-tftp)
[![Build status](https://ci.appveyor.com/api/projects/status/0sxw1t6jjoe4yc9p/branch/master?svg=true)](https://ci.appveyor.com/project/vCabbage/trivialt/branch/master)
[![GoDoc](https://godoc.org/pack.ag/tftp?status.svg)](http://godoc.org/pack.ag/tftp)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/vcabbage/go-tftp/master/LICENSE)


pack.ag/tftp is a cross-platform, concurrent TFTP client and server implementation for Go.


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
    In single port mode, the same port is used for transmit and receive. If the server is started on port 69, all communication will
    be done on port 69.
    
    The primary use case of this feature is to play nicely with firewalls. Most firewalls will prevent the typical case where the server responds
    back on a random port because they have no way of knowing that it is in response to a request that went out on port 69. In single port mode,
    the firewall will see a request go out to a server on port 69 and that server respond back on the same port, which most firewalls will allow.
    
    Of course if the firewall in question is configured to block TFTP connections, this setting won't help you.
    
    Enable single port mode with the `--single-port` flag. This is currently marked experimental as is diverges from the TFTP standard.

## Installation

```
go get -u pack.ag/tftp
```

## API

The API was inspired by Go's well-known net/http API. If you can write a net/http handler or middleware, you should have no problem doing the same with pack.ag/tftp.

### Configuration Functions

One area that is noticeably different from net/http is the configuration of clients and servers. pack.ag/tftp uses "configuration functions" rather than the direct modification of the
Client/Server struct or a configuration struct passed into the factory functions.

A few explanations of this pattern:
* [Self-referential functions and the design of options](http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html) by Rob Pike
* [Functional options for friendly APIs](https://www.youtube.com/watch?v=24lFtGHWxAQ) by Dave Cheney [video]

If this sounds complicated, don't worry, the public API is quiet simple. The `NewClient` and `NewServer` functions take zero or more configuration functions.

Want all defaults? Don't pass anything.

Want a Client configured for blocksize 9000 and windowsize 16? Pass in `ClientBlocksize(9000)` and `ClientWindowsize(16)`.

``` go
// Default Client
tftp.NewClient()

// Client with blocksize 9000, windowsize 16
tftp.NewClient(tftp.ClientBlocksize(9000), tftp.ClientWindowsize(16))

// Configuring with a slice of options
opts := []tftp.ClientOpt{
    tftp.ClientMode(tftp.ModeOctet),
    tftp.ClientBlocksize(9000),
    tftp.ClientWindowsize(16),
    tftp.ClientTimeout(1),
    tftp.ClientTransferSize(true),
    tftp.ClientRetransmit(3),
}

tftp.NewClient(opts...)
```

### Examples

#### Read File From Server, Print to stdout

``` go
client := tftp.NewClient()
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

client := tftp.NewClient()
err := client.Put("myftp.local/myfile", file, fileInfo.Size())
if err != nil {
    log.Fatalln(err)
}
```


#### HTTP Proxy

This rather contrived example proxies an incoming GET request to GitHub's public API. A more realistic use case might be proxying to PXE boot files on an HTTP server.

``` go
const baseURL = "https://api.github.com/"

func proxyTFTP(w tftp.ReadRequest) {
	// Append the requested path to the baseURL
	url := baseURL + w.Name()

	// Send the HTTP request
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		// This could send more specific errors, but here we'read
		// choosing to simply send "file not found"" with the error
		// message from the HTTP client back to the TFTP client.
		w.WriteError(tftp.ErrCodeFileNotFound, err.Error())
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
readHandler := tftp.ReadHandlerFunc(proxyTFTP)

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

Full example in [examples/httpproxy/httpproxy.go](https://github.com/vcabbage/go-tftp/blob/master/examples/httpproxy/httpproxy.go).

#### Save Files to Database

Here `tftpDB` implements the `WriteHandler` interface directly.

``` go
// tftpDB embeds a *sql.DB and implements the tftp.ReadHandler interface.
type tftpDB struct {
	*sql.DB
}

func (db *tftpDB) ReceiveTFTP(w tftp.WriteRequest) {
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

Full example including checking the size before accepting the request in [examples/database/database.go](https://github.com/vcabbage/go-tftp/blob/master/examples/database/database.go).
