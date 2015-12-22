package main

import(
	"bytes"
	"io"
	"net/http"
	"testing"
)

func Test_wget(t *testing.T) {

	var buf, buf2 bytes.Buffer
	
	if err := wget("http://example.com", &buf); err != nil {
		t.Fatalf("%v", err)
	}

	resp, err := http.Get("http://example.com")
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(&buf2, resp.Body)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if bytes.Compare(buf.Bytes(),buf2.Bytes()) != 0 {
		t.Fatalf("buffers differ")
	}
}
