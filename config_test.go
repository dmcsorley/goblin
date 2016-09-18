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
