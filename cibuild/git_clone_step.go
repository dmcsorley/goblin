// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"fmt"
	"github.com/dmcsorley/goblin/command"
	"github.com/dmcsorley/goblin/config"
	"os/exec"
)

type GitCloneStep struct {
	index int
	url   string
}

func newCloneStep(index int, sr *config.StepRecord, vv ValueValidator) (Stepper, error) {
	if !sr.HasParameter(UrlKey) {
		return stepParamRequired(GitClone, UrlKey)
	}

	err := vv.ValidateValue(sr.Url)
	if err != nil {
		return stepParamError(GitClone, UrlKey, err)
	}

	return &GitCloneStep{index: index, url: sr.Url}, nil
}

func (gcs *GitCloneStep) Step(se StepEnvironment) error {
	pfx := se.StepPrefix(gcs.index)
	url, err := se.ResolveValues(gcs.url)
	if err != nil {
		return err
	}

	fmt.Println(pfx, GitClone, url)
	cmd := exec.Command(
		"git",
		"clone",
		url,
		".",
	)
	cmd.Dir = WorkDir
	err = command.Run(cmd, pfx)

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

func (gcs *GitCloneStep) Cleanup(se StepEnvironment) {
	// intentionally left blank, will be cleaned up with the volume
}
