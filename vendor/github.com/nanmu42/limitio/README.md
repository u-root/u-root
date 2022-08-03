# LimitIO

[![GoDoc](https://godoc.org/github.com/nanmu42/limitio?status.svg)](https://pkg.go.dev/github.com/nanmu42/limitio)
[![Build status](https://github.com/nanmu42/limitio/workflows/test/badge.svg)](https://github.com/nanmu42/limitio/actions)
[![codecov](https://codecov.io/gh/nanmu42/limitio/branch/master/graph/badge.svg)](https://codecov.io/gh/nanmu42/limitio)
[![Go Report Card](https://goreportcard.com/badge/github.com/nanmu42/limitio)](https://goreportcard.com/report/github.com/nanmu42/limitio)

`io.Reader` and `io.Writer` with limit.

```bash
go get github.com/nanmu42/limitio
```

## Rationale and Usage

There are times when a limited reader or writer comes in handy.

1. wrap upstream so that reading is metered and limited:

```go
// request is an incoming http.Request
request.Body = limitio.NewReadCloser(c.Request.Body, maxRequestBodySize, false)

// deal with the body now with easy mind. It's maximum size is assured.
```

Yes, `io.LimitReader` works the same way, but throws `EOF` on exceeding limit, which is confusing.

LimitIO provides error that can be identified.

```go
decoder := json.NewDecoder(request.Body)
err := decoder.Decode(&myStruct)
if err != nil {
    if errors.Is(err, limitio.ErrThresholdExceeded) {
        // oops, we reached the limit
    }

    err = fmt.Errorf("other error happened: %w", err)
    return
}
```

2. wrap downstream so that writing is metered and limited(or instead, just pretending writing):

```go
// request is an incoming http.Request.
// Say, we want to record its body somewhere in the middleware,
// but feeling uneasy since its body might be HUGE, which may
// result in OOM and a successful DDOS...

var reqBuf bytes.buffer

// a limited writer comes to rescue!
// `true` means after reaching `RequestBodyMaxLength`,
// `limitedReqBuf` will start pretending writing so that
// io.TeeReader continues working while reqBuf stays unmodified.
limitedReqBuf := limitio.NewWriter(&reqBuf, RequestBodyMaxLength, true)

request.Body = &readCloser{
    Reader: io.TeeReader(request.Body, limitedReqBuf), 
    Closer: c.Request.Body,
}
```

LimitIO provides Reader, Writer and their Closer versions, for details, see [docs](https://pkg.go.dev/github.com/nanmu42/limitio).

## Status: Stable

LimitIO is well battle tested under production environment.

APIs are subjected to change in backward compatible way during 1.x releases.

## License

MIT License

Copyright (c) 2021 LI Zhennan
