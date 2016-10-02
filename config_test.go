// import github.com/dmcsorley/goblin
package main

import (
	"github.com/dmcsorley/goblin/config"
	"testing"
)

func TestServerConfigRequiresBuild(t *testing.T) {
	_, err := configRecordAsGoblin(&config.Record{})
	if err == nil {
		t.Error("should have failed for empty builds")
	}
}

func TestBuildConfigRequiresStep(t *testing.T) {
	_, err := configRecordAsGoblin(&config.Record{
		Builds: map[string]*config.BuildRecord{
			"foo": &config.BuildRecord{},
		},
	})
	if err == nil {
		t.Error("should have failed for empty steps")
	}
}

func TestServerConfigRejectsDuplicateValues(t *testing.T) {
	_, err := configRecordAsGoblin(&config.Record{
		Values: []*config.ValueRecord{
			&config.ValueRecord{
				Name:          "value1",
				Literal:       "example1",
				DecodedFields: []string{"Literal"},
			},
			&config.ValueRecord{
				Name:          "value1",
				Literal:       "example2",
				DecodedFields: []string{"Literal"},
			},
		},
		Builds: map[string]*config.BuildRecord{
			"build1": &config.BuildRecord{
				Steps: []*config.StepRecord{
					&config.StepRecord{
						Type:          "git-clone",
						Url:           "example",
						DecodedFields: []string{"Url"},
					},
				},
			},
		},
	})
	if err == nil {
		t.Error("should have failed for duplicate values")
	}
}

func TestServerConfigRejectsUnspecifiedValues(t *testing.T) {
	_, err := configRecordAsGoblin(&config.Record{
		Values: []*config.ValueRecord{
			&config.ValueRecord{Name: "value1"},
		},
		Builds: map[string]*config.BuildRecord{
			"build1": &config.BuildRecord{
				Steps: []*config.StepRecord{
					&config.StepRecord{
						Type:          "git-clone",
						Url:           "example",
						DecodedFields: []string{"Url"},
					},
				},
			},
		},
	})
	if err == nil {
		t.Error("should have failed for unspecified value")
	}
}
