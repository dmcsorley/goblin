// import github.com/dmcsorley/goblin/config
package cibuild

import (
	"github.com/dmcsorley/goblin/config"
	"testing"
)

func TestUnknownStepTypeFails(t *testing.T) {
	_, err := NewStep(
		0,
		&config.StepRecord{Type: "unknown"},
	)

	if err == nil {
		t.Error("should have failed for unknown step type")
	}
}
