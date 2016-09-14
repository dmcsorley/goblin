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
	Dir string
}

func newBuildStep(index int, stepJson map[string]interface{}) (*DockerBuildStep, error) {
	step := &DockerBuildStep{Index: index}
	image, err := asString(ImageKey, stepJson[ImageKey])
	if err != nil {
		return nil, err
	}
	step.Image = image

	if stepJson[DirKey] != nil {
		dir, err := asString(DirKey, stepJson[DirKey])
		if err != nil {
			return nil, err
		}
		step.Dir = dir
	}

	return step, nil
}

func (dbs *DockerBuildStep) Step(build *Build) error {
	pfx := build.stepPrefix(dbs.Index)
	workDir := WorkDir
	if dbs.Dir != "" {
		workDir = filepath.Join(workDir, dbs.Dir)
	}
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
