package driver

import (
	"io/ioutil"
	"os"
	"regexp"
)

// APPS name of the folder
const APPS string = "apps"

// MODULES name of the folder
const MODULES string = "modules"

var re = regexp.MustCompile("[^A-Za-z0-9/._]+")

type statFunction func(path string) bool
type getConfigFunction func(path string) *[]byte

// FileSystem dads
type FileSystem struct {
	stat      statFunction
	getConfig getConfigFunction
}

//GetConfig get config
func (fs *FileSystem) GetConfig(path string) *[]byte {
	strings.Trim()
	return fs.getConfig(path)
}

//LoadConfig used for load config
func LoadConfig(path string) *[]byte {
	if s, err := ioutil.ReadFile(path); err == nil {
		return &s
	}
	return nil
}

// NewFileSystem dadsasf
func NewFileSystem(stat statFunction, getConfig getConfigFunction) *FileSystem {
	fs := new(FileSystem)
	fs.stat = stat
	fs.getConfig = getConfig
	return fs
}

// Stat aas
func Stat(path string) bool {
	if f, _ := os.Stat(path); f != nil {
		return true
	}
	return false
}
