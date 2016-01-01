package main

import (
	"io/ioutil"
	"os"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func readContent(path string) []byte {
	var data []byte
	if data, _ = ioutil.ReadFile(path); data != nil {
		return data
	}
	return nil
}
