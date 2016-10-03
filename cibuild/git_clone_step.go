// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"errors"
	"fmt"
	"github.com/dmcsorley/goblin/command"
	"github.com/dmcsorley/goblin/config"
	"os/exec"
)

type GitCloneStep struct {
	index int
	url   string
}

func newCloneStep(index int, sr *config.StepRecord) (*GitCloneStep, error) {
	if !sr.HasParameter(UrlKey) {
		return nil, errors.New(GitCloneStepType + " requires " + UrlKey)
	}
	return &GitCloneStep{index: index, url: sr.Url}, nil
}

func (gcs *GitCloneStep) Step(build *Build) error {
	pfx := build.stepPrefix(gcs.index)
	fmt.Println(pfx, GitCloneStepType, gcs.url)
	cmd := exec.Command(
		"git",
		"clone",
		gcs.url,
		".",
	)
	cmd.Dir = WorkDir
	err := command.Run(cmd, pfx)

	if err != nil {
		return err
	}

	cmd = exec.Command(
		"git",
		"log",
		"-n",
		"1",
		"--pretty=oneline",
		"--no-color",
		"--decorate",
		"--abbrev-commit",
	)
	cmd.Dir = WorkDir
	return command.Run(cmd, pfx)
}

func (gcs *GitCloneStep) Cleanup(build *Build) {
	// intentionally left blank, will be cleaned up with the volume
}
