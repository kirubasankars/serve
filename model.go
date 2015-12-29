package main

import (
	"io/ioutil"

	"github.com/metal"
)

func (site *Site) GetModel(path string) *metal.Metal {
	model := metal.NewMetal()
	var jsonPath = site.path + "/api" + path + ".json"
	var data = readContent(jsonPath)
	model.Parse(data)
	return model
}

func readContent(path string) []byte {
	var data []byte
	if data, _ = ioutil.ReadFile(path); data != nil {
		return data
	}
	return data
}
