package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"time"
)

const JOB_TIME_FORMAT = time.RFC3339Nano

type Job struct {
	Id string
	received time.Time
	buildConfig *BuildConfig
}

func joblog(prefix string, message string) {
	os.Stdout.WriteString(
		fmt.Sprintf("%s %s %s\n",
			time.Now().Format(time.RFC3339),
			prefix,
			message,
		),
	)
}

func NewJob(t time.Time, bc *BuildConfig) *Job {
	hasher := sha1.New()
	hasher.Write([]byte(t.Format(JOB_TIME_FORMAT)))
	hasher.Write([]byte(bc.Name))
	id := hex.EncodeToString(hasher.Sum(nil))
	return &Job{
		Id: id,
		received: t,
		buildConfig: bc,
	}
}

func (job *Job) Run() {
	joblog(job.Id, "STARTING " + job.buildConfig.Name)
	var err error
	err = os.Mkdir(job.Id, os.ModeDir)
	if err != nil {
		joblog(job.Id, fmt.Sprintf("ERROR %v", err))
		return
	}

	for _, s := range job.buildConfig.Steps {
		err = s.Step(job)
		if err != nil {
			joblog(job.Id, fmt.Sprintf("ERROR %v", err))
			return
		}
	}

	joblog(job.Id, "SUCCESS")
}

func (job *Job) DockerRun() {
	containerName := "ci-" + job.Id
	joblog(job.Id, "LAUNCHING " + containerName)

	cmd := exec.Command(
		"docker",
		"run",
		"-d",
		"--name=" + containerName,
		"-v",
		"/var/run/docker.sock:/var/run/docker.sock",
		os.Getenv(ENV_IMAGE),
		"/go/bin/goblin",
		"-run",
		job.buildConfig.Name,
		"-time",
		job.received.Format(JOB_TIME_FORMAT),
	)
	cmd.Run()
}
