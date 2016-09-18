// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/dmcsorley/goblin/config"
	"github.com/dmcsorley/goblin/goblog"
	"io"
	"os/exec"
	"strconv"
	"time"
)

const (
	GitCloneStepType = "git-clone"
	DockerBuildStepType = "docker-build"
	UrlKey = "url"
	ImageKey = "image"
)

type Stepper interface {
	Step(build *Build) error
	Cleanup(build *Build)
}

func pipe(prefix string, rc io.ReadCloser) {
	s := bufio.NewScanner(rc)
	for s.Scan() {
		goblog.Log(prefix, s.Text())
	}
}

func (build *Build) stepPrefix(index int) string {
	return build.Id[0:20] + "-" + strconv.Itoa(index)
}

func runInDirAndPipe(cmd *exec.Cmd, dir string, prefix string) error {
	cmd.Dir = dir
	cmdout, _ := cmd.StdoutPipe()
	cmderr, _ := cmd.StderrPipe()
	go pipe(prefix, cmdout)
	go pipe(prefix, cmderr)

	time.Sleep(time.Second)
	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

func NewStep(index int, sr *config.StepRecord) (Stepper, error) {
	switch sr.Type {
	case GitCloneStepType:
		return newCloneStep(index, sr)
	case DockerBuildStepType:
		return newBuildStep(index, sr)
	default:
		return nil, errors.New(fmt.Sprintf("Unknown step %v", sr.Type))
	}
}
