# Package pnm

[![Build Status](https://secure.travis-ci.org/jbuchbinder/gopnm.png)](http://travis-ci.org/jbuchbinder/gopnm)
[![GoDoc](https://godoc.org/github.com/jbuchbinder/gopnm?status.png)](https://godoc.org/github.com/jbuchbinder/gopnm)

Package pnm implements a PBM, PGM and PPM image decoder and encoder.

This package is compatible with Go version 1.


## Installation

```
	go install github.com/jbuchbinder/gopnm
```

## Limitations

Not implemented are:

* Writing pnm files in raw format.
* Writing images with 16 bits per channel.
* Writing images with a custom Maxvalue.
* Reading/and writing PAM images.

(I would be happy to accept patches for these.)

