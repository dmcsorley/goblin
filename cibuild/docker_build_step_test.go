package cibuild

import (
	"testing"
)

func TestNewWithNoDir(t *testing.T) {
	step, err := newBuildStep(
		0,
		map[string]interface{}{
			TypeKey: DockerBuildStepType,
			ImageKey: "foo",
		},
	)

	if err != nil {
		t.Error(err)
	}

	if step.Dir != "" {
		t.Error("docker-build dir should be empty")
	}
}

func TestNewWithDir(t *testing.T) {
	step, err := newBuildStep(
		0,
		map[string]interface{}{
			TypeKey: DockerBuildStepType,
			ImageKey: "foo",
			DirKey: "bar",
		},
	)

	if err != nil {
		t.Error(err)
	}

	if step.Dir != "bar" {
		t.Error("docker-build dir should be bar")
	}
}

func TestDockerBuildStepRequiresImage(t *testing.T) {
	_, err := newBuildStep(
		0,
		map[string]interface{}{
			TypeKey: DockerBuildStepType,
		},
	)

	if err == nil {
		t.Error("docker-build step should have failed with no image")
	}
}
