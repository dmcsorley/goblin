// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"errors"
	"fmt"
	"github.com/dmcsorley/goblin/config"
	"strconv"
)

type StepType string

const (
	GitClone    StepType = "git-clone"
	DockerBuild StepType = "docker-build"
	DockerRun   StepType = "docker-run"
	DockerPull  StepType = "docker-pull"
	UrlKey               = "url"
	ImageKey             = "image"
	DirKey               = "dir"
	CmdKey               = "cmd"
)

type Stepper interface {
	Step(StepEnvironment) error
	Cleanup(StepEnvironment)
}

type StepEnvironment interface {
	StepPrefix(index int) string
	VolumeName() string
}

func (build *Build) StepPrefix(index int) string {
	return build.Id[0:20] + "-" + strconv.Itoa(index)
}

type ValueValidator interface {
	ValidateValue(string) error
}

func stepParamError(st StepType, param string, err interface{}) (Stepper, error) {
	return nil, fmt.Errorf("In %v %s: %v", st, param, err)
}

func stepParamRequired(st StepType, param string) (Stepper, error) {
	return stepParamError(st, param, "parameter required")
}

func NewStep(index int, sr *config.StepRecord, vv ValueValidator) (Stepper, error) {
	switch StepType(sr.Type) {
	case GitClone:
		return newCloneStep(index, sr, vv)
	case DockerBuild:
		return newBuildStep(index, sr, vv)
	case DockerRun:
		return newRunStep(index, sr, vv)
	case DockerPull:
		return newPullStep(index, sr, vv)
	default:
		return nil, errors.New(fmt.Sprintf("Unknown step %v", sr.Type))
	}
}
