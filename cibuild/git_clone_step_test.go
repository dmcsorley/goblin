// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"github.com/dmcsorley/goblin/config"
	"testing"
)

func TestGitCloneStepRequiresUrl(t *testing.T) {
	_, err := newCloneStep(
		0,
		&config.StepRecord{Type: GitCloneStepType},
	)

	if err == nil {
		t.Error("git-clone step should have failed with no url")
	}
}
