package main

import (
	"io/ioutil"
	"os"
)

type SiteMeta struct {
	Name string
}

func Stat(path string) *SiteMeta {
	meta := new(SiteMeta)
	f, err := os.Stat(path)

	if os.IsNotExist(err) {
		return nil
	}

	if f != nil && f.IsDir() {
		meta.Name = f.Name()
		return meta
	}

	return nil
}

func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func ReadContent(path string) []byte {
	var data []byte
	if data, _ = ioutil.ReadFile(path); data != nil {
		return data
	}
	return nil
}