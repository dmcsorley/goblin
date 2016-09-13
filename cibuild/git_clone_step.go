// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"fmt"
	"github.com/dmcsorley/goblin/goblog"
	"os"
	"os/exec"
	"path/filepath"
)

type GitCloneStep struct {
	Index int
	URL string
}

func newCloneStep(index int, stepJson map[string]interface{}) (*GitCloneStep, error) {
	url, err := asString(URLKey, stepJson[URLKey])
	if err != nil {
		return nil, err
	}
	return &GitCloneStep{Index:index, URL:url}, nil
}

func (gcs *GitCloneStep) Step(build *Build) error {
	pfx := build.stepPrefix(gcs.Index)
	workDir := build.Id
	goblog.Log(pfx, GitCloneStepType + " " + gcs.URL)
	cmd := exec.Command(
		"git",
		"clone",
		gcs.URL,
		CloneDir,
	)
	return runInDirAndPipe(cmd, workDir, pfx)
}

func (gcs *GitCloneStep) Cleanup(build *Build) {
	pfx := build.stepPrefix(gcs.Index)
	goblog.Log(pfx, "cleanup")
	clonedDir := filepath.Join(build.Id, CloneDir)
	if err := os.RemoveAll(clonedDir); err != nil {
		goblog.Log(pfx, fmt.Sprintf("%v", err))
	}
}
