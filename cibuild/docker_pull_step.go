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
	Index int
	Image string
}

func newPullStep(index int, sr *config.StepRecord) (*DockerPullStep, error) {
	if !sr.HasField(ImageKey) {
		return nil, errors.New(DockerPullStepType + " requires " + ImageKey)
	}
	return &DockerPullStep{Index: index, Image: sr.Image}, nil
}

func (dbs *DockerPullStep) Step(build *Build) error {
	pfx := build.stepPrefix(dbs.Index)
	fmt.Println(pfx, DockerPullStepType, dbs.Image)
	cmd := exec.Command(
		"docker",
		"pull",
		dbs.Image,
	)
	return command.Run(cmd, pfx)
}

func (dbs *DockerPullStep) Cleanup(build *Build) {
	// intentionally left blank, un-pulling an image doesn't make sense
}
