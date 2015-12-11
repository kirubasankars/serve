package main

import (
	"github.com/serve/server"
	"log"
	"os"
	"path/filepath"
)

func main() {	
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        log.Fatal(err)
		return
    }
	var s = new(server.Server);
	s.SetConfig(dir, "3000")
	s.Start()
}
