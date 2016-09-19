// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/dmcsorley/goblin/gobdocker"
	"github.com/dmcsorley/goblin/goblog"
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
	goblog.Log(build.Id, "STARTING "+build.config.Name)
	var err error
	err = os.Mkdir(build.Id, os.ModeDir)
	if err != nil {
		goblog.Log(build.Id, fmt.Sprintf("ERROR %v", err))
		os.Exit(1)
	}

	for _, s := range build.config.Steps {
		defer s.Cleanup(build)
		err = s.Step(build)
		if err != nil {
			goblog.Log(build.Id, fmt.Sprintf("ERROR %v", err))
			os.Exit(1)
		}
	}

	goblog.Log(build.Id, "SUCCESS")
}

func (build *Build) volumeName() string {
	return VolumePrefix + build.Id
}

func (build *Build) createVolume() (string, error) {
	return gobdocker.CreateVolume(build.volumeName())
}

func (build *Build) DockerRun(image string) {
	volumeName, err := build.createVolume()
	if err != nil {
		goblog.Log(build.Id, fmt.Sprintf("ERROR %v", err))
		return
	}

	containerName := BuildContainerPrefix + build.Id
	goblog.Log(build.Id, "LAUNCHING "+containerName)

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
