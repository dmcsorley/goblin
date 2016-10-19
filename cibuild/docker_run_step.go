// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"fmt"
	"github.com/dmcsorley/goblin/command"
	"github.com/dmcsorley/goblin/config"
	"github.com/dmcsorley/goblin/gobdocker"
	"os/exec"
)

type DockerRunStep struct {
	index int
	image string
	cmd   string
	dir   string
}

func newRunStep(index int, sr *config.StepRecord, vv ValueValidator) (Stepper, error) {
	if !sr.HasParameter(ImageKey) {
		return stepParamRequired(DockerRun, ImageKey)
	}

	err := vv.ValidateValue(sr.Image)
	if err != nil {
		return stepParamError(DockerRun, ImageKey, err)
	}

	drs := &DockerRunStep{index: index, image: sr.Image}

	if sr.HasParameter(CmdKey) {
		err := vv.ValidateValue(sr.Cmd)
		if err != nil {
			return stepParamError(DockerRun, CmdKey, err)
		}
		drs.cmd = sr.Cmd
	}

	if sr.HasParameter(DirKey) {
		err := vv.ValidateValue(sr.Dir)
		if err != nil {
			return stepParamError(DockerRun, DirKey, err)
		}
		drs.dir = sr.Dir
	}

	return drs, nil
}

func (drs *DockerRunStep) Step(se StepEnvironment) error {
	pfx := se.StepPrefix(drs.index)

	workDir := WorkDir
	if drs.dir != "" {
		dir, err := se.ResolveValues(drs.dir)
		if err != nil {
			return err
		}
		workDir = dir
	}

	containerName := BuildContainerPrefix + pfx
	image, err := se.ResolveValues(drs.image)
	if err != nil {
		return err
	}

	args := []string{
		"run",
		"-d",
		"--name",
		containerName,
		"-v",
		se.VolumeName() + ":" + workDir,
		"-w",
		workDir,
		image,
	}

	if drs.cmd != "" {
		cmd, err := se.ResolveValues(drs.cmd)
		if err != nil {
			return err
		}
		args = append(args, "bash", "-c", cmd)
	}

	fmt.Println(pfx, DockerRun, image)

	cmd := exec.Command("docker", args...)
	err = command.Run(cmd, pfx)
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

func (drs *DockerRunStep) Cleanup(se StepEnvironment) {
	pfx := se.StepPrefix(drs.index)
	fmt.Println(pfx, "removing intermediate container")
	containerName := BuildContainerPrefix + pfx
	gobdocker.RemoveContainer(containerName)
}
