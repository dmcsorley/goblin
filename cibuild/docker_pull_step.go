// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"errors"
	"fmt"
	"github.com/dmcsorley/goblin/command"
	"os/exec"
)

type DockerPullStep struct {
	Index      int
	stepConfig StepConfig
}

func newPullStep(index int, sc StepConfig) (*DockerPullStep, error) {
	if !sc.HasImage() {
		return nil, errors.New(DockerPullStepType + " requires " + ImageKey)
	}
	return &DockerPullStep{Index: index, stepConfig: sc}, nil
}

func (dbs *DockerPullStep) Step(build *Build) error {
	pfx := build.stepPrefix(dbs.Index)
	image := dbs.stepConfig.ImageParam()
	fmt.Println(pfx, DockerPullStepType, image)
	cmd := exec.Command(
		"docker",
		"pull",
		image,
	)
	return command.Run(cmd, pfx)
}

func (dbs *DockerPullStep) Cleanup(build *Build) {
	// intentionally left blank, un-pulling an image doesn't make sense
}
