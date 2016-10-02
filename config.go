// import github.com/dmcsorley/goblin
package main

import (
	"errors"
	"fmt"
	"github.com/dmcsorley/goblin/cibuild"
	"github.com/dmcsorley/goblin/config"
	"io/ioutil"
)

type Goblin struct {
	Builds map[string]*cibuild.BuildConfig
}

func loadConfig(filename string) (*Goblin, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cr, err := config.LoadBytes(b)
	if err != nil {
		return nil, err
	}

	return configRecordAsGoblin(cr)
}

func configRecordAsGoblin(cr *config.Record) (*Goblin, error) {
	if len(cr.Builds) == 0 {
		return nil, errors.New("config has no builds")
	}

	values := make(map[string]*config.ValueRecord)
	for _, v := range cr.Values {
		if values[v.Name] != nil {
			return nil, fmt.Errorf("duplicate value '%s'", v.Name)
		}
		if !v.HasField("Literal") {
			return nil, fmt.Errorf("no value for '%s'", v.Name)
		}
		values[v.Name] = v
	}

	sc := &Goblin{Builds: map[string]*cibuild.BuildConfig{}}
	for name, br := range cr.Builds {
		bc, err := newBuild(name, br)
		if err != nil {
			return nil, err
		}
		sc.Builds[name] = bc
	}

	return sc, nil
}

func newBuild(name string, br *config.BuildRecord) (*cibuild.BuildConfig, error) {
	if len(br.Steps) == 0 {
		return nil, errors.New("build has no steps")
	}

	bc := &cibuild.BuildConfig{Name: name}
	for i, sr := range br.Steps {
		step, err := cibuild.NewStep(i, sr)
		if err != nil {
			return nil, err
		}
		bc.Steps = append(bc.Steps, step)
	}

	return bc, nil
}
