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
	fmt.Println(pfx, GitCloneStepType, gcs.Url)
	cmd := exec.Command(
		"git",
		"clone",
		gcs.Url,
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
