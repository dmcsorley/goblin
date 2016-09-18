package cibuild

import (
	"testing"
)

func TestUnknownStepTypeFails(t *testing.T) {
	_, err := NewStep(
		0,
		map[string]interface{}{
			"type": "unknown",
		},
	)

	if err == nil {
		t.Error("should have failed for unknown step type")
	}
}
