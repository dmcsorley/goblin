// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"github.com/dmcsorley/goblin/config"
	"testing"
)

func TestGitCloneStepRequiresUrl(t *testing.T) {
	expectStepConstructorFailure(
		newCloneStep,
		&config.StepRecord{Type: string(GitClone)},
		t,
		"git-clone step should have failed with no url",
	)
}

func TestGitCloneStepFailsForBadUrlValue(t *testing.T) {
	expectStepConstructorFailure(
		newCloneStep,
		&config.StepRecord{
			Type:          string(GitClone),
			Url:           "${badexample}",
			DecodedFields: []string{"url"},
		},
		t,
		"git-clone step should have failed for bad url value",
	)
}
