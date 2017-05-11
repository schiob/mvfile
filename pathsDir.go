package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type dirJSON struct {
	Paths []dirPair `json:"paths"`
}

type dirPair struct {
	InPath  string `json:"in_path"`
	OutPath string `json:"out_path"`
	User    string `json:"user"`
}

var dirMap = make(map[string]dirPair)

func loadDirMap(jsonPath string) {
	file, e := ioutil.ReadFile(jsonPath)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var jsonObj dirJSON
	json.Unmarshal(file, &jsonObj)

	for _, v := range jsonObj.Paths {
		dirMap[v.InPath] = v
	}
}
