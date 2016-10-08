// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/dmcsorley/goblin/gobdocker"
	"os"
	"os/exec"
	"time"
)

const (
	TimeFormat           = time.RFC3339Nano
	WorkDir              = "/tmp/workdir"
	BuildContainerPrefix = "goblin-build-"
	VolumePrefix         = "goblin-volume-"
)

type BuildConfig struct {
	Name  string
	Steps []Stepper
}

type Build struct {
	Id       string
	received time.Time
	config   *BuildConfig
}

func buildHash(timestamp string, buildName string) string {
	hasher := sha1.New()
	hasher.Write([]byte(timestamp))
	hasher.Write([]byte(buildName))
	return hex.EncodeToString(hasher.Sum(nil))
}

func New(t time.Time, bc *BuildConfig) *Build {
	id := buildHash(t.Format(TimeFormat), bc.Name)
	return &Build{
		Id:       id,
		received: t,
		config:   bc,
	}
}

func (build *Build) Run() {
	fmt.Println(build.Id, "STARTING", build.config.Name)

	for _, s := range build.config.Steps {
		defer s.Cleanup(build)
		err := s.Step(build)
		if err != nil {
			fmt.Println(build.Id, "ERROR", err)
			os.Exit(1)
		}
	}

	fmt.Println(build.Id, "SUCCESS")
}

func (build *Build) VolumeName() string {
	return VolumePrefix + build.Id
}

func (build *Build) createVolume() (string, error) {
	return gobdocker.CreateVolume(build.VolumeName())
}

func (build *Build) DockerRun(image string) {
	volumeName, err := build.createVolume()
	if err != nil {
		fmt.Println(build.Id, "ERROR", err)
		return
	}

	containerName := BuildContainerPrefix + build.Id
	fmt.Println(build.Id, "LAUNCHING", containerName)

	ts := build.received.Format(TimeFormat)

	cmd := exec.Command(
		"docker",
		"run",
		"-d",
		"--label=goblin.build="+build.config.Name,
		"--label=goblin.id="+build.Id,
		"--label=goblin.time="+ts,
		"--name="+containerName,
		"-v",
		volumeName+":"+WorkDir,
		"-v",
		"/var/run/docker.sock:/var/run/docker.sock",
		image,
		"goblin",
		"-run",
		build.config.Name,
		"-time",
		ts,
	)
	cmd.Run()
}

func Cleanup(eb *gobdocker.ExitedBuild) {
	gobdocker.RemoveContainer(eb.ContainerId)
	gobdocker.RemoveVolume(VolumePrefix + eb.Id)
}
