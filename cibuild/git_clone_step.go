// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"github.com/dmcsorley/goblin/goblog"
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
	workDir := WorkDir
	goblog.Log(pfx, GitCloneStepType + " " + gcs.URL)
	cmd := exec.Command(
		"git",
		"clone",
		gcs.URL,
	)
	return runInDirAndPipe(cmd, workDir, pfx)
}

func (gcs *GitCloneStep) Cleanup(build *Build) {
	// intentionally left blank, will be cleaned up with the volume
}
