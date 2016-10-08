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
	builds map[string]*cibuild.BuildConfig
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

	values := config.NewValueEngine()
	for _, v := range cr.Values {
		if values.HasValue(v.Name) {
			return nil, fmt.Errorf("duplicate value '%s'", v.Name)
		}
		if !v.HasField("Literal") {
			return nil, fmt.Errorf("no value for '%s'", v.Name)
		}
		values.AddValue(v)
	}

	sc := &Goblin{builds: map[string]*cibuild.BuildConfig{}}
	for name, br := range cr.Builds {
		bc, err := newBuild(name, br, values)
		if err != nil {
			return nil, err
		}
		sc.builds[name] = bc
	}

	return sc, nil
}

func newBuild(name string, br *config.BuildRecord, ve *config.ValueEngine) (*cibuild.BuildConfig, error) {
	if len(br.Steps) == 0 {
		return nil, errors.New("build has no steps")
	}

	bc := &cibuild.BuildConfig{Name: name}
	for i, sr := range br.Steps {
		step, err := cibuild.NewStep(i, sr, ve)
		if err != nil {
			return nil, err
		}
		bc.Steps = append(bc.Steps, step)
	}

	return bc, nil
}
