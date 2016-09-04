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
	CONFIG_FILE = "config.txt"
)

func loadConfig() ([]string, error) {
	file, err := os.Open(CONFIG_FILE)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lines []string
	s := bufio.NewScanner(file)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	return lines, nil
}

func requestForValidPath(r *http.Request, config []string) bool {
	path := strings.TrimPrefix(r.URL.Path, "/")
	for _, s := range config {
		if s == path {
			return true
		}
	}

	return false
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

func runJob(job *Job) {
	joblog(job.Id, "STARTING", os.Stdout)
	direrr := os.Mkdir(job.Id, os.ModeDir)
	if direrr != nil {
		log.Println(direrr)
		return
	}

	cmdPrefix := job.Id[0:20]
	cmd := exec.Command("git", "clone", "--progress", "https://github.com/dmcsorley/simpleci")
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
		joblog(job.Id, "COMPLETE", os.Stdout)
	}
}

func isValidRequest(r *http.Request, cfg []string) bool {
	return strings.ToUpper(r.Method) == "POST" &&
		requestForValidPath(r, cfg)
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
	cfg, err := loadConfig()
	if err !=nil {
		log.Fatal(err)
	}

	log.Println(cfg)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		dumpRequest(r)
		if isValidRequest(r, cfg) {
			job := NewJob(r.URL.Path, now)
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
