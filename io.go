package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"serve/serve"
	"strings"

	"github.com/kirubasankars/metal"
)

type FileIO struct{}

func (fs *FileIO) IsSitePath(path string) bool {
	f, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	if f != nil && f.IsDir() {
		return true
	}
	return false
}

func (fs *FileIO) IsExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (fs *FileIO) ServeFile(w http.ResponseWriter, r *http.Request, path string) {
	http.ServeFile(w, r, path)
}

// func API(site *serve.Site, method string, path string, input *metal.Metal) {
// 	fmt.Println(path)
// 	req, _ := http.NewRequest(method, "http://localhost:5984/test/_design/type/_view/"+path+"?include_docs=true", nil)
// 	client := &http.Client{}
// 	res, err := client.Do(req)
// }

func (fs *FileIO) API(site *serve.Site, method string, path string, input *metal.Metal) *metal.Metal {
	var data []byte
	if data, _ = ioutil.ReadFile(site.Path() + "/api/" + path + "/" + strings.ToLower(method) + ".json"); data != nil {
		m := metal.NewMetal()
		m.Parse(data)
		return m
	}
	return nil
}

func (fs *FileIO) Template(path string) *[]byte {
	var data []byte
	if data, _ = ioutil.ReadFile(path); data != nil {
		return &data
	}
	return nil
}
