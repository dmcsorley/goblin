package main

import (
	"bufio"
        "log"
        "net/http"
        "net/http/httputil"
	"os"
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
                w.WriteHeader(http.StatusOK)
        })
        log.Fatal(http.ListenAndServe(":80", nil))
}
