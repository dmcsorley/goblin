// import github.com/dmcsorley/goblin
package main

import (
	"encoding/json"
	"errors"
	"github.com/dmcsorley/goblin/cibuild"
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
	Builds map[string]*cibuild.BuildConfig
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

	sc := &ServerConfig{Builds:map[string]*cibuild.BuildConfig{}}
	for _, br := range sr.Builds {
		bc, err := newBuild(br)
		if err != nil {
			return nil, err
		}
		sc.Builds[bc.Name] = bc
	}

	return sc, nil
}

func newBuild(br BuildRecord) (*cibuild.BuildConfig, error) {
	if br.Name == "" {
		return nil, errors.New("build has no name")
	}

	if len(br.Steps) == 0 {
		return nil, errors.New("build has no steps")
	}

	bc := &cibuild.BuildConfig{Name:br.Name}
	for i, sjson := range br.Steps {
		step, err := cibuild.NewStep(i, sjson)
		if err != nil {
			return bc, err
		}
		bc.Steps = append(bc.Steps, step)
	}

	return bc, nil
}
