// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"fmt"
	"github.com/dmcsorley/goblin/command"
	"github.com/dmcsorley/goblin/config"
	"os/exec"
)

type DockerLoginStep struct {
	index    int
	username string
	password string
}

func newLoginStep(index int, sr *config.StepRecord, vv ValueValidator) (Stepper, error) {
	if !sr.HasParameter(UsernameKey) {
		return stepParamRequired(DockerLogin, UsernameKey)
	}

	err := vv.ValidateValue(sr.Username)
	if err != nil {
		return stepParamError(DockerLogin, UsernameKey, err)
	}

	if !sr.HasParameter(PasswordKey) {
		return stepParamRequired(DockerLogin, PasswordKey)
	}

	err = vv.ValidateValue(sr.Password)
	if err != nil {
		return stepParamError(DockerLogin, PasswordKey, err)
	}

	return &DockerLoginStep{index: index, username: sr.Username, password: sr.Password}, nil
}

func (dls *DockerLoginStep) Step(se StepEnvironment) error {
	pfx := se.StepPrefix(dls.index)

	username, err := se.ResolveValues(dls.username)
	if err != nil {
		return err
	}

	password, err := se.ResolveValues(dls.password)
	if err != nil {
		return err
	}

	fmt.Println(pfx, DockerLogin, username)
	cmd := exec.Command(
		"docker",
		"login",
		"-p",
		password,
		"-u",
		username,
	)
	return command.Run(cmd, pfx)
}

func (dls *DockerLoginStep) Cleanup(se StepEnvironment) {
	// intentionally left blank
}
