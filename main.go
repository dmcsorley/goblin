package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	CONFIG_FILE = "config.json"
)

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

func runJob(job *Job) {
	joblog(job.Id, "STARTING", os.Stdout)
	direrr := os.Mkdir(job.Id, os.ModeDir)
	if direrr != nil {
		log.Println(direrr)
		return
	}

	cmdPrefix := job.Id[0:20]
	cmd := exec.Command("git", "clone", "--progress", job.GitURL)
	cmd.Dir = job.Id
	cmdout, _ := cmd.StdoutPipe()
	cmderr, _ := cmd.StderrPipe()
	go pipe(cmdPrefix, cmdout, os.Stdout)
	go pipe(cmdPrefix, cmderr, os.Stderr)

	time.Sleep(time.Second)
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return
	}

	err := cmd.Wait()

	if err != nil {
		joblog(job.Id, fmt.Sprintf("ERROR %v", err), os.Stdout)
	} else {
		joblog(job.Id, "SUCCESS", os.Stdout)
	}
}

func isValidRequest(r *http.Request, cfg *ServerConfig) *BuildConfig {
	if strings.ToUpper(r.Method) != "POST" {
		return nil
	}

	return cfg.BuildConfigForPath(r.URL.Path)
}

func dumpRequest(r *http.Request) {
	bytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(bytes))
}

func main() {
	cfg, err := loadConfig(CONFIG_FILE)
	if err !=nil {
		log.Fatal(err)
	}

	log.Println(cfg)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		dumpRequest(r)
		bc := isValidRequest(r, cfg)
		if bc != nil {
			job := NewJob(r.URL.Path, now, bc)
			log.Println("Received request to " + r.URL.Path + " for job " + job.Id)
			w.WriteHeader(http.StatusOK)
			go runJob(job);
		} else {
			log.Println("Bad request for " + r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
		}

	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
