// import github.com/dmcsorley/goblin/command
package command

import (
	"bufio"
	"github.com/dmcsorley/goblin/goblog"
	"io"
	"os/exec"
	"time"
)

func pipe(prefix string, rc io.ReadCloser) {
	s := bufio.NewScanner(rc)
	for s.Scan() {
		goblog.Log(prefix, s.Text())
	}
}

func Run(cmd *exec.Cmd, prefix string) error {
	cmdout, _ := cmd.StdoutPipe()
	cmderr, _ := cmd.StderrPipe()
	go pipe(prefix, cmdout)
	go pipe(prefix, cmderr)

	time.Sleep(time.Second)
	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}
