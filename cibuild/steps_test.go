// import github.com/dmcsorley/goblin/config
package cibuild

import (
	"github.com/dmcsorley/goblin/config"
	"testing"
)

type StepConstructor func(int, *config.StepRecord, ValueValidator) (Stepper, error)

func expectStepConstructorFailure(sc StepConstructor, sr *config.StepRecord, t *testing.T, message string) {
	_, err := sc(0, sr, config.NewValueEngine())
	if err == nil {
		t.Error(message)
	}
}

func TestUnknownStepTypeFails(t *testing.T) {
	_, err := NewStep(
		0,
		&config.StepRecord{Type: "unknown"},
		config.NewValueEngine(),
	)

	if err == nil {
		t.Error("should have failed for unknown step type")
	}
}
