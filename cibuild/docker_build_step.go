package cibuild

import (
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
	Log(pfx, DockerBuildStepType + " " + dbs.Image)
	cmd := exec.Command(
		"docker",
		"build",
		 "-t",
		dbs.Image + ":" + pfx,
		".",
	)
	return runInDirAndPipe(cmd, workDir, pfx)
}
