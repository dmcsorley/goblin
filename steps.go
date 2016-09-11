package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
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
	url, err := asString("url", stepJson["url"])
	if err != nil {
		return nil, err
	}
	return &GitCloneStep{URL:url}, nil
}

func (gcs *GitCloneStep) Step(dir string, stepPrefix string) error {
	joblog(stepPrefix, "git-clone " + gcs.URL, os.Stdout)
	cmd := exec.Command(
		"git",
		"clone",
		"--progress",
		gcs.URL,
		"clone",
	)
	return runInDirAndPipe(cmd, dir, stepPrefix)
}

type DockerBuildStep struct {
	Image string
}

func newBuildStep(stepJson map[string]interface{}) (*DockerBuildStep, error) {
	image, err := asString("image", stepJson["image"])
	if err != nil {
		return nil, err
	}
	return &DockerBuildStep{Image:image}, nil
}

func (dbs *DockerBuildStep) Step(dir string, stepPrefix string) error {
	joblog(stepPrefix, "docker-build " + dbs.Image, os.Stdout)
	cmd := exec.Command(
		"docker",
		"build",
		 "-t",
		dbs.Image + ":" + stepPrefix,
		".",
	)
	return runInDirAndPipe(cmd, dir + "/clone", stepPrefix)
}

func newStep(stepJson map[string]interface{}) (Stepper, error) {
	typeValue, err := asString("type", stepJson["type"])
	if err != nil {
		return nil, err
	}

	switch typeValue {
	case "git-clone":
		return newCloneStep(stepJson)
	case "docker-build":
		return newBuildStep(stepJson)
	default:
		return nil, errors.New(fmt.Sprintf("Unknown step %v", typeValue))
	}
}
