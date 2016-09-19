// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"errors"
	"github.com/dmcsorley/goblin/config"
	"github.com/dmcsorley/goblin/gobdocker"
	"github.com/dmcsorley/goblin/goblog"
	"os/exec"
	"time"
)

type DockerRunStep struct {
	Index int
	Image string
	Dir   string
	Cmd   string
}

func newRunStep(index int, sr *config.StepRecord) (*DockerRunStep, error) {
	if !sr.HasField(ImageKey) {
		return nil, errors.New(DockerRunStepType + " requires " + ImageKey)
	}

	drs := &DockerRunStep{Index: index, Image: sr.Image}

	if sr.HasField(DirKey) {
		drs.Dir = sr.Dir
	}

	if sr.HasField(CmdKey) {
		drs.Cmd = sr.Cmd
	}

	return drs, nil
}

func (drs *DockerRunStep) Step(build *Build) error {
	pfx := build.stepPrefix(drs.Index)
	time.Sleep(5*time.Second)

	workDir := WorkDir
	if drs.Dir != "" {
		workDir = drs.Dir
	}

	containerName := BuildContainerPrefix + pfx

	command := []string{
		"run",
		"-d",
		"--name",
		containerName,
		"-v",
		build.volumeName() + ":" + workDir,
		"-w",
		workDir,
		drs.Image,
	}

	if drs.Cmd != "" {
		command = append(command, "bash", "-c", drs.Cmd)
	}

	goblog.Log(pfx, DockerRunStepType+" "+drs.Image)

	cmd := exec.Command("docker", command...)
	err := runInDirAndPipe(cmd, WorkDir, pfx)
	if err != nil {
		return err
	}

	cmd = exec.Command("docker", "wait", containerName)
	return runInDirAndPipe(cmd, WorkDir, pfx)
}

func (drs *DockerRunStep) Cleanup(build *Build) {
	pfx := build.stepPrefix(drs.Index)
	goblog.Log(pfx, "removing intermediate container")
	containerName := BuildContainerPrefix + pfx
	gobdocker.RemoveContainer(containerName)
}
