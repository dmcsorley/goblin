// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"errors"
	"fmt"
	"github.com/dmcsorley/goblin/config"
	"strconv"
)

const (
	GitCloneStepType    = "git-clone"
	DockerBuildStepType = "docker-build"
	DockerRunStepType   = "docker-run"
	UrlKey              = "url"
	ImageKey            = "image"
	DirKey              = "dir"
	CmdKey              = "cmd"
)

type Stepper interface {
	Step(build *Build) error
	Cleanup(build *Build)
}

func (build *Build) stepPrefix(index int) string {
	return build.Id[0:20] + "-" + strconv.Itoa(index)
}

func NewStep(index int, sr *config.StepRecord) (Stepper, error) {
	switch sr.Type {
	case GitCloneStepType:
		return newCloneStep(index, sr)
	case DockerBuildStepType:
		return newBuildStep(index, sr)
	case DockerRunStepType:
		return newRunStep(index, sr)
	default:
		return nil, errors.New(fmt.Sprintf("Unknown step %v", sr.Type))
	}
}
