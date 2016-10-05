// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"errors"
	"fmt"
	"github.com/dmcsorley/goblin/command"
	"os/exec"
)

type DockerBuildStep struct {
	Index      int
	stepConfig StepConfig
}

func newBuildStep(index int, sc StepConfig) (*DockerBuildStep, error) {
	if !sc.HasImage() {
		return nil, errors.New(DockerBuildStepType + " requires " + ImageKey)
	}
	return &DockerBuildStep{Index: index, stepConfig: sc}, nil
}

func (dbs *DockerBuildStep) Step(build *Build) error {
	pfx := build.stepPrefix(dbs.Index)
	image := dbs.stepConfig.ImageParam()
	fmt.Println(pfx, DockerBuildStepType, image)
	cmd := exec.Command(
		"docker",
		"build",
		"--force-rm",
		"-t",
		image,
		".",
	)
	cmd.Dir = WorkDir
	return command.Run(cmd, pfx)
}

func (dbs *DockerBuildStep) Cleanup(build *Build) {
	// intentionally left blank
}
