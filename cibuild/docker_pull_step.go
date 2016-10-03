// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"errors"
	"fmt"
	"github.com/dmcsorley/goblin/command"
	"github.com/dmcsorley/goblin/config"
	"os/exec"
)

type DockerPullStep struct {
	index int
	image string
}

func newPullStep(index int, sr *config.StepRecord) (*DockerPullStep, error) {
	if !sr.HasParameter(ImageKey) {
		return nil, errors.New(DockerPullStepType + " requires " + ImageKey)
	}
	return &DockerPullStep{index: index, image: sr.Image}, nil
}

func (dbs *DockerPullStep) Step(build *Build) error {
	pfx := build.stepPrefix(dbs.index)
	fmt.Println(pfx, DockerPullStepType, dbs.image)
	cmd := exec.Command(
		"docker",
		"pull",
		dbs.image,
	)
	return command.Run(cmd, pfx)
}

func (dbs *DockerPullStep) Cleanup(build *Build) {
	// intentionally left blank, un-pulling an image doesn't make sense
}
