// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"github.com/dmcsorley/goblin/config"
	"testing"
)

func TestDockerRunStepRequiresImage(t *testing.T) {
	expectStepConstructorFailure(
		newRunStep,
		&config.StepRecord{Type: string(DockerRun)},
		t,
		"docker-run step should have failed with no image",
	)
}

func TestDockerRunStepFailsForBadImageValue(t *testing.T) {
	expectStepConstructorFailure(
		newRunStep,
		&config.StepRecord{
			Type:          string(DockerRun),
			Image:         "${badexample}",
			DecodedFields: []string{"image"},
		},
		t,
		"docker-run step should have failed for bad image value",
	)
}

func TestDockerRunStepFailsForBadCmdValue(t *testing.T) {
	expectStepConstructorFailure(
		newRunStep,
		&config.StepRecord{
			Type:          string(DockerRun),
			Image:         "valid",
			Cmd:           "${badexample}",
			DecodedFields: []string{"image", "cmd"},
		},
		t,
		"docker-run step should have failed for bad cmd value",
	)
}

func TestDockerRunStepFailsForBadDirValue(t *testing.T) {
	expectStepConstructorFailure(
		newRunStep,
		&config.StepRecord{
			Type:          string(DockerRun),
			Image:         "valid",
			Dir:           "${badexample}",
			DecodedFields: []string{"image", "dir"},
		},
		t,
		"docker-run step should have failed for bad dir value",
	)
}
