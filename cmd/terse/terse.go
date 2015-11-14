package main

import (
	"log"

	"github.com/brettlangdon/go-arg"
	"github.com/brettlangdon/terse"
)

var args struct {
	MaxEntries int    `arg:"-m,--max,help:max number of links to keep [default: 1000]"`
	Bind       string `arg:"-b,--bind,help:\"[host]:<port>\" to bind the server to [default: 127.0.0.1:5892]"`
}

func main() {
	args.MaxEntries = 1000
	args.Bind = "127.0.0.1:5892"
	arg.MustParse(&args)

	server := terse.NewServer(args.Bind, args.MaxEntries)
	log.Fatal(server.ListenAndServe())
}
