// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"fmt"
	"github.com/dmcsorley/goblin/command"
	"github.com/dmcsorley/goblin/config"
	"os/exec"
)

type DockerPushStep struct {
	index int
	image string
}

func newPushStep(index int, sr *config.StepRecord, vv ValueValidator) (Stepper, error) {
	if !sr.HasParameter(ImageKey) {
		return stepParamRequired(DockerPush, ImageKey)
	}

	err := vv.ValidateValue(sr.Image)
	if err != nil {
		return stepParamError(DockerPush, ImageKey, err)
	}

	return &DockerPushStep{index: index, image: sr.Image}, nil
}

func (dps *DockerPushStep) Step(se StepEnvironment) error {
	pfx := se.StepPrefix(dps.index)
	image, err := se.ResolveValues(dps.image)
	if err != nil {
		return err
	}

	fmt.Println(pfx, DockerPush, image)
	cmd := exec.Command(
		"docker",
		"push",
		image,
	)
	return command.Run(cmd, pfx)
}

func (dps *DockerPushStep) Cleanup(se StepEnvironment) {
	// intentionally left blank
}
