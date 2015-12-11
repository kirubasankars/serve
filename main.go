package main

import (	
	"gravity/serve/server"
)

func main() {
	var s = new(server.Server);
	s.SetConfig("C:/workspace/go/bin", "3000")
	s.Start()
}
