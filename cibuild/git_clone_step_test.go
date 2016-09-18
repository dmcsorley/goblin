package cibuild

import (
	"testing"
)

func TestGitCloneStepRequiresUrl(t *testing.T) {
	_, err := newCloneStep(
		0,
		map[string]interface{}{
			TypeKey: GitCloneStepType,
		},
	)

	if err == nil {
		t.Error("git-clone step should have failed with no url")
	}
}
