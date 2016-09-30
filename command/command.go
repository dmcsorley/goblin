// import github.com/dmcsorley/goblin/command
package command

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)

func pipe(prefix string, rc io.ReadCloser, wg *sync.WaitGroup) {
	s := bufio.NewScanner(rc)
	for s.Scan() {
		fmt.Println(prefix, s.Text())
	}
	wg.Done()
}

func Run(cmd *exec.Cmd, prefix string) error {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	cmdout, _ := cmd.StdoutPipe()
	cmderr, _ := cmd.StderrPipe()
	go pipe(prefix, cmdout, wg)
	go pipe(prefix, cmderr, wg)

	time.Sleep(time.Second)
	if err := cmd.Start(); err != nil {
		return err
	}

	wg.Wait()
	return cmd.Wait()
}
