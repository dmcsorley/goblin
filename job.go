package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
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

func (job *Job) Run() {
	joblog(job.Id, "STARTING", os.Stdout)
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
