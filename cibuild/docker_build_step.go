// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"fmt"
	"github.com/dmcsorley/goblin/goblog"
	"os/exec"
	"path/filepath"
)

type DockerBuildStep struct {
	Index int
	Image string
}

func newBuildStep(index int, stepJson map[string]interface{}) (*DockerBuildStep, error) {
	image, err := asString(ImageKey, stepJson[ImageKey])
	if err != nil {
		return nil, err
	}
	return &DockerBuildStep{Index:index, Image:image}, nil
}

func (dbs *DockerBuildStep) Step(build *Build) error {
	pfx := build.stepPrefix(dbs.Index)
	workDir := filepath.Join(build.Id, CloneDir)
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
	goblog.Log(pfx, "cleanup")
	cmd := exec.Command(
		"docker",
		"rmi",
		dbs.Image + ":" + pfx,
	)
	if err := cmd.Run(); err != nil {
		goblog.Log(pfx, fmt.Sprintf("%v", err))
	}
}
