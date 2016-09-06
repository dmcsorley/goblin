package main

import (
	"encoding/json"
	"io/ioutil"
)

type BuildConfig struct {
	Name string
	Steps []map[string]string
}

type ServerConfig struct {
	Builds []BuildConfig
}

func loadConfig(filename string) (*ServerConfig, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	sc := &ServerConfig{}

	err = json.Unmarshal(bytes, sc)
	if err != nil {
		return nil, err
	}

	return sc, nil
}
