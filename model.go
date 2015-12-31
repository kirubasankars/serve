package main

import (
	"io/ioutil"

	"github.com/serve/lib/metal"
)

func (site *Site) GetModel(path string) *metal.Metal {
	model := metal.NewMetal()
	var apiPath = site.path + "/api" + path
	if data := readContent(apiPath + ".sql"); data != nil {
		getData(string(data), model)
		return model
	} else {
		var data = readContent(apiPath + ".json")
		model.Parse(data)
		return model
	}
}

func readContent(path string) []byte {
	var data []byte
	if data, _ = ioutil.ReadFile(path); data != nil {
		return data
	}
	return nil
}
