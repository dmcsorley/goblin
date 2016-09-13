package cibuild

import (
	"os/exec"
)

type GitCloneStep struct {
	Index int
	URL string
}

func newCloneStep(index int, stepJson map[string]interface{}) (*GitCloneStep, error) {
	url, err := asString(URLKey, stepJson[URLKey])
	if err != nil {
		return nil, err
	}
	return &GitCloneStep{Index:index, URL:url}, nil
}

func (gcs *GitCloneStep) Step(build *Build) error {
	pfx := build.stepPrefix(gcs.Index)
	workDir := build.Id
	Log(pfx, GitCloneStepType + " " + gcs.URL)
	cmd := exec.Command(
		"git",
		"clone",
		gcs.URL,
		CloneDir,
	)
	return runInDirAndPipe(cmd, workDir, pfx)
}
