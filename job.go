package main

import (
	"crypto/sha1"
	"encoding/hex"
	"time"
)

type Job struct {
	Id string
	job string
	received time.Time
}

func NewJob(job string, t time.Time) *Job {
	hasher := sha1.New()
	hasher.Write([]byte(t.Format(time.RFC3339Nano)))
	hasher.Write([]byte(job))
	id := hex.EncodeToString(hasher.Sum(nil))
	return &Job{
		Id: id,
		job: job,
		received: t,
	}
}
