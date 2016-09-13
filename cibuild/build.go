package cibuild

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
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

func Log(prefix string, message string) {
	os.Stdout.WriteString(
		fmt.Sprintf("%s %s %s\n",
			time.Now().Format(time.RFC3339),
			prefix,
			message,
		),
	)
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
	Log(build.Id, "STARTING " + build.config.Name)
	var err error
	err = os.Mkdir(build.Id, os.ModeDir)
	if err != nil {
		Log(build.Id, fmt.Sprintf("ERROR %v", err))
		return
	}

	for _, s := range build.config.Steps {
		defer s.Cleanup(build)
		err = s.Step(build)
		if err != nil {
			Log(build.Id, fmt.Sprintf("ERROR %v", err))
			return
		}
	}

	Log(build.Id, "SUCCESS")
}

func (build *Build) DockerRun(image string) {
	containerName := "ci-" + build.Id
	Log(build.Id, "LAUNCHING " + containerName)

	cmd := exec.Command(
		"docker",
		"run",
		"-d",
		"--name=" + containerName,
		"-v",
		"/var/run/docker.sock:/var/run/docker.sock",
		image,
		"/go/bin/goblin",
		"-run",
		build.config.Name,
		"-time",
		build.received.Format(TimeFormat),
	)
	cmd.Run()
}
