// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"errors"
	"github.com/dmcsorley/goblin/config"
	"github.com/dmcsorley/goblin/goblog"
	"os/exec"
)

type GitCloneStep struct {
	Index int
	Url   string
}

func newCloneStep(index int, sr *config.StepRecord) (*GitCloneStep, error) {
	if !sr.HasField(UrlKey) {
		return nil, errors.New(GitCloneStepType + " requires " + UrlKey)
	}
	return &GitCloneStep{Index: index, Url: sr.Url}, nil
}

func (gcs *GitCloneStep) Step(build *Build) error {
	pfx := build.stepPrefix(gcs.Index)
	workDir := WorkDir
	goblog.Log(pfx, GitCloneStepType+" "+gcs.Url)
	cmd := exec.Command(
		"git",
		"clone",
		gcs.Url,
		".",
	)
	return runInDirAndPipe(cmd, workDir, pfx)
}

func (gcs *GitCloneStep) Cleanup(build *Build) {
	// intentionally left blank, will be cleaned up with the volume
}
