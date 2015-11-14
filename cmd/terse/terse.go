package main

import (
	"fmt"
	"log"

	"github.com/brettlangdon/go-arg"
	"github.com/brettlangdon/terse"
)

var args struct {
	MaxEntries int    `arg:"-m,--max,help:max number of links to keep [default: 1000]"`
	Bind       string `arg:"-b,--bind,help:\"[host]:<port>\" to bind the server to [default: 127.0.0.1:5892]"`
	ServerURL  string `arg:"-s,--server,help:base server url to generate links as (e.g. \"https://short.domain.com\") [default: \"http://<bind>\"]"`
}

func main() {
	// Setup default args
	args.MaxEntries = 1000
	args.Bind = "127.0.0.1:5892"
	// Parse args from CLI
	arg.MustParse(&args)

	if args.ServerURL == "" {
		args.ServerURL = fmt.Sprintf("http://%s", args.Bind)
	}

	// Start the server
	server, err := terse.NewServer(args.Bind, args.MaxEntries, args.ServerURL)
	if err == nil {
		log.Printf("Listening on \"%s\"", args.Bind)
		err = server.ListenAndServe()
	}
	log.Fatal(err)
}
