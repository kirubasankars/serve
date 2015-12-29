package main

import (
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

	var s = new(Server)
	s.SetConfig(dir, "3000")
	s.Start()
}
