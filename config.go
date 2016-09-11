package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type ServerRecord struct {
	Builds []BuildRecord
}

type BuildRecord struct {
	Name string
	Steps []map[string]interface{}
}

type ServerConfig struct {
	Builds []BuildConfig
}

type BuildConfig struct {
	Name string
	Steps []Stepper
}

func loadConfig(filename string) (*ServerConfig, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return loadConfigBytes(bytes)
}

func loadConfigBytes(bytes []byte) (*ServerConfig, error) {
	sr := &ServerRecord{}

	err := json.Unmarshal(bytes, sr)
	if err != nil {
		return nil, err
	}

	if len(sr.Builds) == 0 {
		return nil, errors.New("server config has no builds")
	}

	sc := &ServerConfig{}
	for _, br := range sr.Builds {
		bc, err := newBuild(br)
		if err != nil {
			return nil, err
		}
		sc.Builds = append(sc.Builds, bc)
	}

	return sc, nil
}

func newBuild(br BuildRecord) (BuildConfig, error) {
	bc := BuildConfig{}
	if br.Name == "" {
		return bc, errors.New("build has no name")
	}

	if len(br.Steps) == 0 {
		return bc, errors.New("build has no steps")
	}

	bc.Name = br.Name
	for i, sjson := range br.Steps {
		step, err := newStep(i, sjson)
		if err != nil {
			return bc, err
		}
		bc.Steps = append(bc.Steps, step)
	}

	return bc, nil
}

func (sc *ServerConfig) FindBuildByName(name string) *BuildConfig {
	for _, bc := range sc.Builds {
		if bc.Name == name {
			return &bc
		}
	}
	return nil
}
