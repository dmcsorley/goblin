package main

import (
	"bufio"
        "log"
        "net/http"
        "net/http/httputil"
	"os"
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

func main() {
	cfg, err := loadConfig()
	if err !=nil {
		log.Fatal(err)
	}

	log.Println(cfg)

        http.HandleFunc("/",  func(w http.ResponseWriter, r *http.Request) {
		bytes, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(bytes))
		if strings.ToUpper(r.Method) == "POST" && requestForValidPath(r, cfg) {
			w.WriteHeader(http.StatusOK)
		} else {
			log.Println("Bad request for " + r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
		}

        })
        log.Fatal(http.ListenAndServe(":80", nil))
}
