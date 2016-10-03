package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/kirubasankars/serve/driver"
	"github.com/kirubasankars/serve/serve"
)

var port = flag.Int("port", 3000, "port number to listen")
var rootPath = flag.String("path", "../", "root path from serve file")

func main() {

	flag.Parse()

	var err error
	var path string

	if len(*rootPath) >= 0 {
		path, err = filepath.Abs(*rootPath)
	} else {
		path, err = filepath.Abs(filepath.Dir(filepath.Join(os.Args[0], "..")))
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	server := serve.NewServer(strconv.Itoa(*port), path, driver.NewFileSystem(driver.Stat, driver.LoadConfig))
	server.Start()
}
