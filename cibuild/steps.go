// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	GitCloneStepType    = "git-clone"
	DockerBuildStepType = "docker-build"
	DockerRunStepType   = "docker-run"
	DockerPullStepType  = "docker-pull"
	UrlKey              = "url"
	ImageKey            = "image"
	DirKey              = "dir"
	CmdKey              = "cmd"
)

type StepConfig interface {
	StepType() string
	HasUrl() bool
	UrlParam() string
	HasImage() bool
	ImageParam() string
	HasDir() bool
	DirParam() string
	HasCmd() bool
	CmdParam() string
}

type Stepper interface {
	Step(build *Build) error
	Cleanup(build *Build)
}

func (build *Build) stepPrefix(index int) string {
	return build.Id[0:20] + "-" + strconv.Itoa(index)
}

func NewStep(index int, sc StepConfig) (Stepper, error) {
	switch sc.StepType() {
	case GitCloneStepType:
		return newCloneStep(index, sc)
	case DockerBuildStepType:
		return newBuildStep(index, sc)
	case DockerRunStepType:
		return newRunStep(index, sc)
	case DockerPullStepType:
		return newPullStep(index, sc)
	default:
		return nil, errors.New(fmt.Sprintf("Unknown step %v", sc.StepType()))
	}
}
