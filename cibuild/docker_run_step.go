// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"errors"
	"fmt"
	"github.com/dmcsorley/goblin/command"
	"github.com/dmcsorley/goblin/gobdocker"
	"os/exec"
	"time"
)

type DockerRunStep struct {
	Index      int
	stepConfig StepConfig
}

func newRunStep(index int, sc StepConfig) (*DockerRunStep, error) {
	if !sc.HasImage() {
		return nil, errors.New(DockerRunStepType + " requires " + ImageKey)
	}

	return &DockerRunStep{Index: index, stepConfig: sc}, nil
}

func (drs *DockerRunStep) Step(build *Build) error {
	pfx := build.stepPrefix(drs.Index)
	time.Sleep(5 * time.Second)

	workDir := WorkDir
	if drs.stepConfig.HasDir() {
		workDir = drs.stepConfig.DirParam()
	}

	containerName := BuildContainerPrefix + pfx
	image := drs.stepConfig.ImageParam()

	args := []string{
		"run",
		"-d",
		"--name",
		containerName,
		"-v",
		build.volumeName() + ":" + workDir,
		"-w",
		workDir,
		image,
	}

	if drs.stepConfig.HasCmd() {
		args = append(args, "bash", "-c", drs.stepConfig.CmdParam())
	}

	fmt.Println(pfx, DockerRunStepType, image)

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
