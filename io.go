package serve

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"fmt"

	//"github.com/kirubasankars/gravity"
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
	fmt.Println(path)
	http.ServeFile(w, r, path)
}

func (fs *FileIO) API(site *Site, method string, path string, input *metal.Metal) *metal.Metal {
	//g := gravity.NewGravity()
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
