// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"errors"
	"fmt"
	"github.com/dmcsorley/goblin/config"
	"github.com/dmcsorley/goblin/goblog"
	"os/exec"
)

type DockerBuildStep struct {
	Index int
	Image string
}

func newBuildStep(index int, sr *config.StepRecord) (*DockerBuildStep, error) {
	if !sr.HasField(ImageKey) {
		return nil, errors.New(DockerBuildStepType + " requires " + ImageKey)
	}
	return &DockerBuildStep{Index:index, Image:sr.Image}, nil
}

func (dbs *DockerBuildStep) Step(build *Build) error {
	pfx := build.stepPrefix(dbs.Index)
	workDir := WorkDir
	goblog.Log(pfx, DockerBuildStepType + " " + dbs.Image)
	cmd := exec.Command(
		"docker",
		"build",
		 "-t",
		dbs.Image + ":" + pfx,
		".",
	)
	return runInDirAndPipe(cmd, workDir, pfx)
}

func (dbs *DockerBuildStep) Cleanup(build *Build) {
	pfx := build.stepPrefix(dbs.Index)
	goblog.Log(pfx, "removing intermediate image")
	cmd := exec.Command(
		"docker",
		"rmi",
		dbs.Image + ":" + pfx,
	)
	if err := cmd.Run(); err != nil {
		goblog.Log(pfx, fmt.Sprintf("%v", err))
	}
}
