package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	Step(dir string, prefix string) error
}

func asString(key string, i interface{}) (string, error) {
	switch value := i.(type) {
	case string:
		return value, nil
	default:
		return "", errors.New(fmt.Sprintf("expected a string for %v", key))
	}
}

func runInDirAndPipe(cmd *exec.Cmd, dir string, stepPrefix string) error {
	cmd.Dir = dir
	cmdout, _ := cmd.StdoutPipe()
	cmderr, _ := cmd.StderrPipe()
	go pipe(stepPrefix, cmdout, os.Stdout)
	go pipe(stepPrefix, cmderr, os.Stdout)

	time.Sleep(time.Second)
	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

type GitCloneStep struct {
	URL string
}

func newCloneStep(stepJson map[string]interface{}) (*GitCloneStep, error) {
	url, err := asString(URLKey, stepJson[URLKey])
	if err != nil {
		return nil, err
	}
	return &GitCloneStep{URL:url}, nil
}

func (gcs *GitCloneStep) Step(dir string, stepPrefix string) error {
	joblog(stepPrefix, GitCloneStepType + " " + gcs.URL, os.Stdout)
	cmd := exec.Command(
		"git",
		"clone",
		"--progress",
		gcs.URL,
		CloneDir,
	)
	return runInDirAndPipe(cmd, dir, stepPrefix)
}

type DockerBuildStep struct {
	Image string
}

func newBuildStep(stepJson map[string]interface{}) (*DockerBuildStep, error) {
	image, err := asString(ImageKey, stepJson[ImageKey])
	if err != nil {
		return nil, err
	}
	return &DockerBuildStep{Image:image}, nil
}

func (dbs *DockerBuildStep) Step(dir string, stepPrefix string) error {
	joblog(stepPrefix, DockerBuildStepType + " " + dbs.Image, os.Stdout)
	cmd := exec.Command(
		"docker",
		"build",
		 "-t",
		dbs.Image + ":" + stepPrefix,
		".",
	)
	return runInDirAndPipe(cmd, filepath.Join(dir, CloneDir), stepPrefix)
}

func newStep(stepJson map[string]interface{}) (Stepper, error) {
	typeValue, err := asString(TypeKey, stepJson[TypeKey])
	if err != nil {
		return nil, err
	}

	switch typeValue {
	case GitCloneStepType:
		return newCloneStep(stepJson)
	case DockerBuildStepType:
		return newBuildStep(stepJson)
	default:
		return nil, errors.New(fmt.Sprintf("Unknown step %v", typeValue))
	}
}
