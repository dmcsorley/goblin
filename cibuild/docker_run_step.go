// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"errors"
	"fmt"
	"github.com/dmcsorley/goblin/command"
	"github.com/dmcsorley/goblin/config"
	"github.com/dmcsorley/goblin/gobdocker"
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
	time.Sleep(5 * time.Second)

	workDir := WorkDir
	if drs.Dir != "" {
		workDir = drs.Dir
	}

	containerName := BuildContainerPrefix + pfx

	args := []string{
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
		args = append(args, "bash", "-c", drs.Cmd)
	}

	fmt.Println(pfx, DockerRunStepType, drs.Image)

	cmd := exec.Command("docker", args...)
	cmd.Dir = WorkDir
	err := command.Run(cmd, pfx)
	if err != nil {
		return err
	}

	i, err := gobdocker.WaitContainer(containerName)
	if err != nil {
		return err
	}

	if i != 0 {
		return fmt.Errorf("Container exited %v", i)
	}

	return nil
}

func (drs *DockerRunStep) Cleanup(build *Build) {
	pfx := build.stepPrefix(drs.Index)
	fmt.Println(pfx, "removing intermediate container")
	containerName := BuildContainerPrefix + pfx
	gobdocker.RemoveContainer(containerName)
}
