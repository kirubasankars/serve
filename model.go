package main

import (
	"io/ioutil"

	"github.com/serve/lib/metal"
)

func (site *Site) GetModel(path string) *metal.Metal {
	model := metal.NewMetal()
	var apiPath = site.path + "/api" + path
	var data = readContent(apiPath + ".json")
	model.Parse(data)
	return model
}

func readContent(path string) []byte {
	var data []byte
	if data, _ = ioutil.ReadFile(path); data != nil {
		return data
	}
	return nil
}
