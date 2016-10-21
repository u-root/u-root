package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	host = flag.String("h", "127.0.0.1", "IP")
	port = flag.String("p", "8080", "port")
	dir  = flag.String("d", ".", "dir")
)

func main() {
	flag.Parse()
	log.Fatal(http.ListenAndServe(*host+":"+*port, http.FileServer(http.Dir(*dir))))
}
