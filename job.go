package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"
)

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
	hasher.Write([]byte(t.Format(time.RFC3339Nano)))
	hasher.Write([]byte(bc.Name))
	id := hex.EncodeToString(hasher.Sum(nil))
	return &Job{
		Id: id,
		received: t,
		buildConfig: bc,
	}
}

func runInDirAndPipe(cmd *exec.Cmd, dir string, stepPrefix string) error {
	cmd.Dir = dir
	cmdout, _ := cmd.StdoutPipe()
	cmderr, _ := cmd.StderrPipe()
	go pipe(stepPrefix, cmdout, os.Stdout)
	go pipe(stepPrefix, cmderr, os.Stderr)

	time.Sleep(time.Second)
	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

func cloneStep(dir string, stepPrefix string, cloneConfig map[string]string) error {
	cmd := exec.Command(
		"git",
		"clone",
		"--progress",
		cloneConfig["url"],
		"clone",
	)
	return runInDirAndPipe(cmd, dir, stepPrefix)
}

func buildStep(dir string, stepPrefix string, buildConfig map[string]string) error {
	cmd := exec.Command(
		"docker",
		"build",
		 "-t",
		buildConfig["image"] + ":" + stepPrefix,
		".",
	)
	return runInDirAndPipe(cmd, dir + "/clone", stepPrefix)
}

func (job *Job) Run() {
	joblog(job.Id, "STARTING", os.Stdout)
	var err error
	err = os.Mkdir(job.Id, os.ModeDir)
	if err != nil {
		joblog(job.Id, fmt.Sprintf("ERROR %v", err), os.Stdout)
		return
	}

	cmdPrefix := job.Id[0:20]

	Steps: for i, s := range job.buildConfig.Steps {
		stepPrefix := cmdPrefix + "-" + strconv.Itoa(i)
		joblog(stepPrefix, s["type"], os.Stdout)
		switch s["type"] {
		case "git-clone":
			err = cloneStep(job.Id, stepPrefix, s)
			if err != nil {
				break Steps
			}
		case "docker-build":
			err = buildStep(job.Id, stepPrefix, s)
			if err != nil {
				break Steps
			}
		default:
			err = errors.New(fmt.Sprintf("Unknown step %v", s))
			break Steps
		}
	}

	if err != nil {
		joblog(job.Id, fmt.Sprintf("ERROR %v", err), os.Stdout)
	} else {
		joblog(job.Id, "SUCCESS", os.Stdout)
	}
}
