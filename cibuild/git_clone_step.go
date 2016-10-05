// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"errors"
	"fmt"
	"github.com/dmcsorley/goblin/command"
	"os/exec"
)

type GitCloneStep struct {
	Index      int
	stepConfig StepConfig
}

func newCloneStep(index int, sc StepConfig) (*GitCloneStep, error) {
	if !sc.HasUrl() {
		return nil, errors.New(GitCloneStepType + " requires " + UrlKey)
	}
	return &GitCloneStep{Index: index, stepConfig: sc}, nil
}

func (gcs *GitCloneStep) Step(build *Build) error {
	pfx := build.stepPrefix(gcs.Index)
	url := gcs.stepConfig.UrlParam()
	fmt.Println(pfx, GitCloneStepType, url)
	cmd := exec.Command(
		"git",
		"clone",
		url,
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
