package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
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

func joblog(prefix string, message string, w io.Writer) {
	io.WriteString(
		w,
		fmt.Sprintf("%s %s %s\n",
			time.Now().Format(time.RFC3339),
			prefix,
			message,
		),
	)
}

func pipe(prefix string, rc io.ReadCloser, w io.Writer) {
	s := bufio.NewScanner(rc)
	for s.Scan() {
		joblog(prefix, s.Text(), w)
	}
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
	joblog(job.Id, "STARTING " + job.buildConfig.Name, os.Stdout)
	var err error
	err = os.Mkdir(job.Id, os.ModeDir)
	if err != nil {
		joblog(job.Id, fmt.Sprintf("ERROR %v", err), os.Stdout)
		return
	}

	for _, s := range job.buildConfig.Steps {
		err = s.Step(job)
		if err != nil {
			joblog(job.Id, fmt.Sprintf("ERROR %v", err), os.Stdout)
			return
		}
	}

	joblog(job.Id, "SUCCESS", os.Stdout)
}

func (job *Job) DockerRun() {
	containerName := "ci-" + job.Id
	joblog(job.Id, "LAUNCHING " + containerName, os.Stdout)

	cmd := exec.Command(
		"docker",
		"run",
		"-d",
		"--name=" + containerName,
		"-v",
		"/var/run/docker.sock:/var/run/docker.sock",
		os.Getenv(ENV_IMAGE),
		"app",
		"-run",
		job.buildConfig.Name,
		"-time",
		job.received.Format(JOB_TIME_FORMAT),
	)
	cmd.Run()
}
