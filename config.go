package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
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

func (sc *ServerConfig) BuildConfigForPath(path string) *BuildConfig {
	name := strings.TrimPrefix(path, "/")
	for _, bc := range sc.Builds {
		if name == bc.Name {
			return &bc
		}
	}

	return nil
}
