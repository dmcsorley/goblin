// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"errors"
	"fmt"
	"github.com/dmcsorley/goblin/command"
	"github.com/dmcsorley/goblin/config"
	"os/exec"
)

type DockerBuildStep struct {
	index int
	image string
}

func newBuildStep(index int, sr *config.StepRecord) (*DockerBuildStep, error) {
	if !sr.HasParameter(ImageKey) {
		return nil, errors.New(DockerBuildStepType + " requires " + ImageKey)
	}
	return &DockerBuildStep{index: index, image: sr.Image}, nil
}

func (dbs *DockerBuildStep) Step(build *Build) error {
	pfx := build.stepPrefix(dbs.index)
	fmt.Println(pfx, DockerBuildStepType, dbs.image)
	cmd := exec.Command(
		"docker",
		"build",
		"--force-rm",
		"-t",
		dbs.image,
		".",
	)
	cmd.Dir = WorkDir
	return command.Run(cmd, pfx)
}

func (dbs *DockerBuildStep) Cleanup(build *Build) {
	// intentionally left blank
}
