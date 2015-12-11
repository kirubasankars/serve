package main

import (	
	"github.com/serve/server"
	"github.com/serve/metal"
)

func main() {
	var m = metal.NewMetal();
	_ = m
	var s = new(server.Server);
	s.SetConfig("C:/workspace/go/bin", "3000")
	s.Start()
}
