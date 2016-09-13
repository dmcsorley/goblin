package cibuild

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"time"
)

const (
	GitCloneStepType = "git-clone"
	DockerBuildStepType = "docker-build"
	URLKey = "url"
	TypeKey = "type"
	ImageKey = "image"
	CloneDir = "clone"
)

type Stepper interface {
	Step(build *Build) error
	Cleanup(build *Build)
}

func asString(key string, i interface{}) (string, error) {
	switch value := i.(type) {
	case string:
		return value, nil
	default:
		return "", errors.New(fmt.Sprintf("expected a string for %v", key))
	}
}

func pipe(prefix string, rc io.ReadCloser) {
	s := bufio.NewScanner(rc)
	for s.Scan() {
		Log(prefix, s.Text())
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

func NewStep(index int, stepJson map[string]interface{}) (Stepper, error) {
	typeValue, err := asString(TypeKey, stepJson[TypeKey])
	if err != nil {
		return nil, err
	}

	switch typeValue {
	case GitCloneStepType:
		return newCloneStep(index, stepJson)
	case DockerBuildStepType:
		return newBuildStep(index, stepJson)
	default:
		return nil, errors.New(fmt.Sprintf("Unknown step %v", typeValue))
	}
}
