package main

import (
	"flag"
	"fmt"

	"github.com/rivik/go-aux/pkg/appver"
)

var (
	addr        string
	printAppVer bool
)

func init() {
	flag.BoolVar(&printAppVer, "version", false, "print app version and exit")

	flag.StringVar(&addr, "listen", "127.0.0.1:8080", "addr:port to listen on")
	flag.Parse()
}

func main() {
	if printAppVer {
		fmt.Printf("%+v\n", appver.Version)
		return
	}

	mainAPIInitAndServeForever(addr)
}
