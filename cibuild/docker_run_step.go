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
	index int
	image string
	cmd   string
	dir   string
}

func newRunStep(index int, sr *config.StepRecord) (*DockerRunStep, error) {
	if !sr.HasParameter(ImageKey) {
		return nil, errors.New(DockerRunStepType + " requires " + ImageKey)
	}

	drs := &DockerRunStep{index: index, image: sr.Image}

	if sr.HasParameter(CmdKey) {
		drs.cmd = sr.Cmd
	}

	if sr.HasParameter(DirKey) {
		drs.dir = sr.Dir
	}

	return drs, nil
}

func (drs *DockerRunStep) Step(build *Build) error {
	pfx := build.stepPrefix(drs.index)
	time.Sleep(5 * time.Second)

	workDir := WorkDir
	if drs.dir != "" {
		workDir = drs.dir
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
		drs.image,
	}

	if drs.cmd != "" {
		args = append(args, "bash", "-c", drs.cmd)
	}

	fmt.Println(pfx, DockerRunStepType, drs.image)

	cmd := exec.Command("docker", args...)
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
	pfx := build.stepPrefix(drs.index)
	fmt.Println(pfx, "removing intermediate container")
	containerName := BuildContainerPrefix + pfx
	gobdocker.RemoveContainer(containerName)
}
