// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"fmt"
	"github.com/dmcsorley/goblin/command"
	"github.com/dmcsorley/goblin/config"
	"os/exec"
)

type DockerBuildStep struct {
	index int
	image string
}

func newBuildStep(index int, sr *config.StepRecord, vv ValueValidator) (Stepper, error) {
	if !sr.HasParameter(ImageKey) {
		return stepParamRequired(DockerBuild, ImageKey)
	}

	err := vv.ValidateValue(sr.Image)
	if err != nil {
		return stepParamError(DockerBuild, ImageKey, err)
	}

	return &DockerBuildStep{index: index, image: sr.Image}, nil
}

func (dbs *DockerBuildStep) Step(se StepEnvironment) error {
	pfx := se.StepPrefix(dbs.index)
	fmt.Println(pfx, DockerBuild, dbs.image)
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

func (dbs *DockerBuildStep) Cleanup(se StepEnvironment) {
	// intentionally left blank
}
