package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	Step(job *Job) error
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
	Index int
	URL string
}

func newCloneStep(index int, stepJson map[string]interface{}) (*GitCloneStep, error) {
	url, err := asString(URLKey, stepJson[URLKey])
	if err != nil {
		return nil, err
	}
	return &GitCloneStep{Index:index, URL:url}, nil
}

func (gcs *GitCloneStep) Step(job *Job) error {
	stepPrefix := job.Id[0:20] + "-" + strconv.Itoa(gcs.Index)
	workDir := job.Id
	joblog(stepPrefix, GitCloneStepType + " " + gcs.URL, os.Stdout)
	cmd := exec.Command(
		"git",
		"clone",
		gcs.URL,
		CloneDir,
	)
	return runInDirAndPipe(cmd, workDir, stepPrefix)
}

type DockerBuildStep struct {
	Index int
	Image string
}

func newBuildStep(index int, stepJson map[string]interface{}) (*DockerBuildStep, error) {
	image, err := asString(ImageKey, stepJson[ImageKey])
	if err != nil {
		return nil, err
	}
	return &DockerBuildStep{Index:index, Image:image}, nil
}

func (dbs *DockerBuildStep) Step(job *Job) error {
	stepPrefix := job.Id[0:20] + "-" + strconv.Itoa(dbs.Index)
	workDir := filepath.Join(job.Id, CloneDir)
	joblog(stepPrefix, DockerBuildStepType + " " + dbs.Image, os.Stdout)
	cmd := exec.Command(
		"docker",
		"build",
		 "-t",
		dbs.Image + ":" + stepPrefix,
		".",
	)
	return runInDirAndPipe(cmd, workDir, stepPrefix)
}

func newStep(index int, stepJson map[string]interface{}) (Stepper, error) {
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
