// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"fmt"
	"github.com/dmcsorley/goblin/command"
	"github.com/dmcsorley/goblin/config"
	"os/exec"
)

type DockerPullStep struct {
	index int
	image string
}

func newPullStep(index int, sr *config.StepRecord, vv ValueValidator) (Stepper, error) {
	if !sr.HasParameter(ImageKey) {
		return stepParamRequired(DockerPull, ImageKey)
	}

	err := vv.ValidateValue(sr.Image)
	if err != nil {
		return stepParamError(DockerPull, ImageKey, err)
	}

	return &DockerPullStep{index: index, image: sr.Image}, nil
}

func (dps *DockerPullStep) Step(se StepEnvironment) error {
	pfx := se.StepPrefix(dps.index)
	image, err := se.ResolveValues(dps.image)
	if err != nil {
		return err
	}

	fmt.Println(pfx, DockerPull, image)
	cmd := exec.Command(
		"docker",
		"pull",
		image,
	)
	return command.Run(cmd, pfx)
}

func (dps *DockerPullStep) Cleanup(se StepEnvironment) {
	// intentionally left blank, un-pulling an image doesn't make sense
}
