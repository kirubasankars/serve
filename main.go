package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"serve/serve"
)

var port = flag.Int("port", 3000, "port number to listen")
var path = flag.String("path", DefaultPath(), "www folder path")

func DefaultPath() string {
	c, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return c
}

func main() {
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
		return
	}

	var server = new(serve.Server)
	server.IO = new(FileIO)
	server.SetConfig(dir, strconv.Itoa(*port))
	server.Start()
}
