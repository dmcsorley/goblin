// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/dmcsorley/goblin/goblog"
	"os"
	"os/exec"
	"time"
)

const TimeFormat = time.RFC3339Nano

type BuildConfig struct {
	Name string
	Steps []Stepper
}

type Build struct {
	Id string
	received time.Time
	config *BuildConfig
}

func New(t time.Time, bc *BuildConfig) *Build {
	hasher := sha1.New()
	hasher.Write([]byte(t.Format(TimeFormat)))
	hasher.Write([]byte(bc.Name))
	id := hex.EncodeToString(hasher.Sum(nil))
	return &Build{
		Id: id,
		received: t,
		config: bc,
	}
}

func (build *Build) Run() {
	goblog.Log(build.Id, "STARTING " + build.config.Name)
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

func (build *Build) DockerRun(image string) {
	containerName := "ci-" + build.Id
	goblog.Log(build.Id, "LAUNCHING " + containerName)

	ts := build.received.Format(TimeFormat)

	cmd := exec.Command(
		"docker",
		"run",
		"-d",
		"--label=goblin.build=" + build.config.Name,
		"--label=goblin.id=" + build.Id,
		"--label=goblin.time=" + ts,
		"--name=" + containerName,
		"-v",
		"/var/run/docker.sock:/var/run/docker.sock",
		image,
		"/go/bin/goblin",
		"-run",
		build.config.Name,
		"-time",
		ts,
	)
	cmd.Run()
}
