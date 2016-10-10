// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"github.com/dmcsorley/goblin/config"
	"testing"
)

func TestDockerLoginStepParams(t *testing.T) {
	sr := &config.StepRecord{Type: string(DockerLogin)}
	sr.Password = "foo"
	sr.DecodedFields = []string{PasswordKey}

	expectStepConstructorFailure(
		newLoginStep,
		sr,
		t,
		"docker-login step should have failed with no username",
	)

	sr.Username = "foo"
	sr.Password = ""
	sr.DecodedFields = []string{UsernameKey}

	expectStepConstructorFailure(
		newLoginStep,
		sr,
		t,
		"docker-login step should have failed with no password",
	)

	sr.Username = "${foo}"
	sr.Password = "bar"
	sr.DecodedFields = []string{UsernameKey, PasswordKey}

	expectStepConstructorFailure(
		newLoginStep,
		sr,
		t,
		"docker-login step should have failed with invalid username",
	)

	sr.Username = "foo"
	sr.Password = "${bar}"

	expectStepConstructorFailure(
		newLoginStep,
		sr,
		t,
		"docker-login step should have failed with invalid password",
	)
}
