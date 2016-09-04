package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"strings"
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

func pipe(rc io.ReadCloser, w io.Writer) {
	s := bufio.NewScanner(rc)
	for s.Scan() {
		io.WriteString(w, s.Text() + "\n")
	}
}

func runJob() {
	cmd := exec.Command("git", "clone", "--progress", "https://github.com/dmcsorley/simpleci")
	cmdout, _ := cmd.StdoutPipe()
	cmderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return
	}

	go pipe(cmdout, os.Stdout)
	go pipe(cmderr, os.Stderr)
	cmd.Wait()
}

func main() {
	cfg, err := loadConfig()
	if err !=nil {
		log.Fatal(err)
	}

	log.Println(cfg)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(bytes))
		if strings.ToUpper(r.Method) == "POST" && requestForValidPath(r, cfg) {
			runJob();
			w.WriteHeader(http.StatusOK)
		} else {
			log.Println("Bad request for " + r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
		}

	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
